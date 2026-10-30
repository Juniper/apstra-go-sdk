// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/internal/urls"
)

// VirtualNetworkBindingsRequest is the input to SetVirtualNetworkLeafBindings and
// UpdateVirtualNetworkLeafBindings. Both Bindings and SVIAddressing maps are keyed
// by their respective SystemID values. Their map key and SystemID must match.
//
// Maps, rather than slices, are used for this payload to provide a mechanism which
// signals removal in UpdateVirtualNetworkLeafBindings().
//
// Nil map values are not permitted by SetVirtualNetworkLeafBindings().
//
// Note that Bindings map entries for members of a leaf switch redundancy group must
// refer to the group ID, not either of the individual switch system IDs. Both the
// map key and the value's SystemID must use the redundancy group ID for members of
// a leaf switch pair.
//
// The entries in the SVIAddressing map always refer to individual switch IDs, even
// for members of a redundancy group.
type VirtualNetworkBindingsRequest struct {
	VirtualNetworkID   string
	Bindings           map[string]*datacenter.VNBinding
	SVIAddressing      map[string]*datacenter.SVIAddressing
	DHCPServiceEnabled *datacenter.DHCPServiceEnabled
}

// SetVirtualNetworkLeafBindings replaces the current set of SVI and Binding info on
// the Virtual Network identified by vnId. Bindings and SVIAddresssing maps must never
// contain nil entries, and for both maps, the map key must match the SystemID in the
// map value. For members of a leaf switch redundancy group, the Bindings map key and
// value's SystemID must be the redundancy group ID, not either of the switch system
// IDs representing group members.
func (o *TwoStageL3ClosClient) SetVirtualNetworkLeafBindings(ctx context.Context, req VirtualNetworkBindingsRequest) error {
	// Set of system IDs which are being (re)bound - used to validate SVIAddresing data
	systemsToBind := make(map[string]struct{})

	// Validate each Bindings map entry and collect the system IDs for which SVIAddressing data is permitted.
	for k, v := range req.Bindings {
		if v == nil {
			return fmt.Errorf("binding request for target %q is invalid - value required, got nil", k)
		}
		if k != v.SystemID {
			return fmt.Errorf("binding request for target %q is invalid - key/SystemID mismatch - expected %q, got %q", k, k, v.SystemID)
		}
		keyNodeType, err := o.GetBindingNodeType(ctx, k)
		if err != nil {
			return fmt.Errorf("getting binding node type for %s: %w", k, err)
		}

		switch keyNodeType {
		case NodeTypeSystem:
			groupID, peerID, err := o.GetRedundancyGroupBySystem(ctx, k)
			if err != nil {
				return fmt.Errorf("getting redundancy group for system %s: %w", k, err)
			}
			if groupID != nil || peerID != nil {
				return fmt.Errorf("binding request for system %q is invalid - found redundancy group info - use group ID instead", k)
			}
			systemsToBind[k] = struct{}{} // stash system ID for subsequent validation of SVIAddressing data
		case NodeTypeRedundancyGroup:
			systems, err := o.GetSystemsByRedundancyGroup(ctx, k)
			if err != nil {
				return fmt.Errorf("getting group members for group ID %q: %w", k, err)
			}
			systemsToBind[systems[0]] = struct{}{} // stash system ID for subsequent validation of SVIAddressing data
			systemsToBind[systems[1]] = struct{}{} // stash system ID for subsequent validation of SVIAddressing data
		default: // womp
			return fmt.Errorf("unhandled node type %q for %q", keyNodeType, k)
		}
	}

	// Validate each SVIAddressing map entry.
	for k, v := range req.SVIAddressing {
		if v == nil {
			return fmt.Errorf("addressing request for leaf %q in VN %q is invalid - value required, got nil", k, req.VirtualNetworkID)
		}
		if k != v.SystemID {
			return fmt.Errorf("addressing request for leaf %q in VN %q is invalid - key/SystemID mismatch - expected %q, got %q", k, req.VirtualNetworkID, k, v.SystemID)
		}
		if _, ok := systemsToBind[k]; !ok {
			return fmt.Errorf("addressing request for leaf %q in VN %q is invalid - the request neglects to bind that leaf to the VN", k, req.VirtualNetworkID)
		}
	}

	// Make sure we send empty bindings slice ([] vs. null) if the caller has supplied a nil or empty map.
	bindings := slices.Collect(maps.Values(req.Bindings))
	if bindings == nil {
		bindings = make([]*datacenter.VNBinding, 0)
	}

	// Make sure we send empty SVIAddressing slice ([] vs. null) if the caller has supplied a nil or empty map.
	sviIPs := slices.Collect(maps.Values(req.SVIAddressing))
	if sviIPs == nil {
		sviIPs = make([]*datacenter.SVIAddressing, 0)
	}

	o.client.lock(o.blueprintId.String() + "_" + req.VirtualNetworkID)
	defer o.client.unlock(o.blueprintId.String() + "_" + req.VirtualNetworkID)

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPatch,
		urlStr: fmt.Sprintf(urls.DatacenterVirtualNetworkByID, o.blueprintId, req.VirtualNetworkID),
		apiInput: struct {
			SviIps             []*datacenter.SVIAddressing    `json:"svi_ips"`
			BoundTo            []*datacenter.VNBinding        `json:"bound_to"`
			DhcpServiceEnabled *datacenter.DHCPServiceEnabled `json:"dhcp_service,omitempty"`
		}{
			SviIps:             sviIPs,
			BoundTo:            bindings,
			DhcpServiceEnabled: req.DHCPServiceEnabled,
		},
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

// UpdateVirtualNetworkLeafBindings updates the supplied SVI and Binding info on
// the Virtual Network identified by vnId. The update requires two API calls:
// The current set of SVI and Binding info is fetched, the supplied data is
// merged, and then the combined set is sent back to the API. Binding and SVI
// configuration can be removed using a nil map entry.
//
// Note that Bindings map entries for members of a leaf switch redundancy group must
// refer to the group ID, not either of the individual switch system IDs. Both the
// map key and the value's SystemID must use the redundancy group ID for members of
// a leaf switch pair.
//
// The entries in the SVIAddressing map always refer to individual switch IDs, even
// for members of a redundancy group.
func (o *TwoStageL3ClosClient) UpdateVirtualNetworkLeafBindings(ctx context.Context, req VirtualNetworkBindingsRequest) error {
	// todo: Consider relaxing the requirement that every entry in the SVIAddressing map
	//  must be represented (either by group ID or by system ID) in the Bindings map. A
	//  caller may want to update SVIAddressing for a leaf switch which is already bound.

	// Set of system IDs which are being explicitly (re)bound
	systemsToBind := make(map[string]struct{})

	// Set of system IDs which are being explicitly unbound
	systemsToUnbind := make(map[string]struct{})

	// Set of group IDs which are being explicitly (re)bound
	groupsToBind := make(map[string]struct{})

	// Set of group IDs which are being explicitly unbound
	groupsToUnbind := make(map[string]struct{})

	for k, v := range req.Bindings {
		// Ensure key matches payload for non-nil entries. Nil entries are unbind requests and have no payload.
		if v != nil && k != v.SystemID {
			return fmt.Errorf("binding request for system %q is invalid - key/SystemID mismatch - expected %q, got %q", k, k, v.SystemID)
		}

		// To record per-system binding/unbinding, we need to know whether
		// the caller has supplied a system ID or a redundancy group ID.
		keyNodeType, err := o.GetBindingNodeType(ctx, k)
		if err != nil {
			return fmt.Errorf("getting binding node type for %s: %w", k, err)
		}

		if keyNodeType == NodeTypeSystem {
			groupID, peerID, err := o.GetRedundancyGroupBySystem(ctx, k)
			if err != nil {
				return fmt.Errorf("getting group info for system ID %q: %w", k, err)
			}
			if groupID != nil || peerID != nil {
				return fmt.Errorf("binding request for system %q is invalid - found redundancy group info - use group ID instead", k)
			}
		}

		if v == nil {
			// this is an unbind request. Record the system ID(s) for subsequent removal of SVIAddressing.

			switch keyNodeType {
			case NodeTypeSystem:
				systemsToUnbind[k] = struct{}{} // Caller is un-binding a leaf. Record the system ID for subsequent removal of SVIAddressing.
			case NodeTypeRedundancyGroup:
				groupsToUnbind[k] = struct{}{} // Record the group ID for subsequent removal.

				// Caller is unbinding a leaf pair. Get the IDs of the two leaf
				// systems in the pair so we can remove any SVIAddressing for them.
				systems, err := o.GetSystemsByRedundancyGroup(ctx, k)
				if err != nil {
					return fmt.Errorf("getting group members for group ID %q: %w", k, err)
				}
				systemsToUnbind[systems[0]] = struct{}{} // Record the first system ID for subsequent removal of SVIAddressing.
				systemsToUnbind[systems[1]] = struct{}{} // Record the second system ID for subsequent removal of SVIAddressing.
			default: // womp
				return fmt.Errorf("unhandled node type %q for %q", keyNodeType, k)
			}

			// Un-bind request have no value to process, so we can skip the rest of the loop and move on to the next map entry.
			continue
		}

		// This is a non-nil entry, so the caller is requesting a (re)binding.

		switch keyNodeType {
		case NodeTypeSystem:
			systemsToBind[k] = struct{}{} // Record the system ID in the set of systems which are being (re)bound to the VN.
		case NodeTypeRedundancyGroup:
			// Caller is binding a leaf pair. Retrieve the system IDs associated with the pair.
			systems, err := o.GetSystemsByRedundancyGroup(ctx, k)
			if err != nil {
				return fmt.Errorf("getting group members for group ID %q: %w", k, err)
			}

			// Record the group and system IDs to their respective sets.
			groupsToBind[k] = struct{}{}           // Record the group ID in the set of groups which are being (re)bound to the VN.
			systemsToBind[systems[0]] = struct{}{} // Record the first system ID in the set of systems which are being (re)bound to the VN.
			systemsToBind[systems[1]] = struct{}{} // Record the second system ID in the set of systems which are being (re)bound to the VN.
		default: // womp
			return fmt.Errorf("unhandled node type %q for %q", keyNodeType, k)
		}
	}

	for k, v := range req.SVIAddressing {
		if v == nil {
			continue // nil signals removal. We don't need to check the map key against the map value.
		}

		if k != v.SystemID { // Map key must match the map value's system ID.
			return fmt.Errorf("addressing entry for leaf %q is invalid - key/SystemID mismatch - expected %q, got %q", k, k, v.SystemID)
		}

		// Check that the caller isn't trying to set SVIAddressing for a system for which they've neglected to request a binding.
		if _, ok := systemsToBind[k]; !ok {
			return fmt.Errorf("addressing requested for system %q in VN %q but the request neglects to bind that leaf to the VN", k, req.VirtualNetworkID)
		}

		// Check that the caller isn't trying to set SVIAddressing for a system which is being unbound in the same request.
		if _, ok := systemsToUnbind[k]; ok {
			return fmt.Errorf("addressing requested for system %q in VN %q but the request un-binds that leaf from the VN", k, req.VirtualNetworkID)
		}
	}

	o.client.lock(o.blueprintId.String() + "_" + req.VirtualNetworkID)
	defer o.client.unlock(o.blueprintId.String() + "_" + req.VirtualNetworkID)

	// collect current VN info (bindings and SVI IPs)
	vn, err := o.GetVirtualNetwork(ctx, req.VirtualNetworkID)
	if err != nil {
		return fmt.Errorf("update VN binding info failed to fetch current bindings - %w", err)
	}

	// Convert current Bindings into a map keyed by binding target ID for easier merging with
	// caller-supplied data. The API returns group IDs for members of a redundancy group, we we
	// expect this map to be keyed appropriately by either group ID (for groups) or by system ID
	// (for individual leaf switches). The resulting map should be fully compatible with the
	// validated caller request map.
	bindingMap := make(map[string]datacenter.VNBinding, len(vn.Bindings))
	for _, binding := range vn.Bindings {
		bindingMap[binding.SystemID] = binding
	}

	// Convert the current SVIAddressing into a map keyed by system ID for easier merging
	// with caller-supplied data. The API returns system IDs for each SVIAddressing entry,
	// so the resulting map should be fully compatible with the validated caller request map.
	sviAddressingMap := make(map[string]datacenter.SVIAddressing, len(vn.SVIIPs))
	for _, sviAddressing := range vn.SVIIPs {
		sviAddressingMap[sviAddressing.SystemID] = sviAddressing
	}

	// Drop any current bindings for systems or groups which the caller has mentioned, whether
	// binding or unbinding. The caller's request takes precedence over the current state of the VN.
	for k := range bindingMap {
		if _, ok := systemsToBind[k]; ok {
			delete(bindingMap, k)
			continue
		}
		if _, ok := groupsToBind[k]; ok {
			delete(bindingMap, k)
			continue
		}
		if _, ok := systemsToUnbind[k]; ok {
			delete(bindingMap, k)
			continue
		}
		if _, ok := groupsToUnbind[k]; ok {
			delete(bindingMap, k)
			continue // end of function - keeping this for consistency with the checks above
		}
	}

	// Drop any current SVIAddressing for systems which the caller has either:
	//  - requested that the system be unbound from the VN (nil Bindings map entry)
	//  - mentioned in the SVIAddressing request map (nil or non-nil).
	for k := range sviAddressingMap {
		if _, ok := systemsToUnbind[k]; ok {
			delete(sviAddressingMap, k)
			continue
		}
		if _, ok := req.SVIAddressing[k]; ok {
			delete(sviAddressingMap, k)
			continue
		}
	}

	// Merge the caller's non-nil bindings with the surviving bindings retrieved from the API.
	bindingSlice := slices.Collect(maps.Values(bindingMap))
	for _, binding := range req.Bindings {
		if binding != nil {
			bindingSlice = append(bindingSlice, *binding)
		}
	}

	// Merge the caller's non-nil SVIAddressing with the surviving SVIAddressing retrieved from the API.
	sviAddressingSlice := slices.Collect(maps.Values(sviAddressingMap))
	for _, sviAddressing := range req.SVIAddressing {
		if sviAddressing != nil {
			sviAddressingSlice = append(sviAddressingSlice, *sviAddressing)
		}
	}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPatch,
		urlStr: fmt.Sprintf(urls.DatacenterVirtualNetworkByID, o.blueprintId, req.VirtualNetworkID),
		apiInput: struct {
			SviIps             []datacenter.SVIAddressing     `json:"svi_ips"`
			BoundTo            []datacenter.VNBinding         `json:"bound_to"`
			DhcpServiceEnabled *datacenter.DHCPServiceEnabled `json:"dhcp_service,omitempty"`
		}{
			SviIps:             sviAddressingSlice,
			BoundTo:            bindingSlice,
			DhcpServiceEnabled: req.DHCPServiceEnabled,
		},
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
