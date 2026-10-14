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

type SetVNToLeafBindingRequest struct {
	VNID    string
	Binding datacenter.VNBinding
	SVIIP   *datacenter.SVIAddressing
	Unbind  bool
}

// SetVNToLeafBinding binds a single VN to a leaf switch and its access switches.
func (o *TwoStageL3ClosClient) SetVNToLeafBinding(ctx context.Context, req SetVNToLeafBindingRequest) error {
	if req.SVIIP != nil && req.Binding.SystemID != req.SVIIP.SystemID {
		return fmt.Errorf("binding system ID (%s) does not match SVI IP system ID (%s)", req.Binding.SystemID, req.SVIIP.SystemID)
	}
	if req.Unbind && req.SVIIP != nil {
		return fmt.Errorf("cannot unbind and set SVI IP at the same time")
	}

	// Save the leaf ID, then convert any system IDs in the VNBinding to group IDs.
	leafID := req.Binding.SystemID
	err := o.convertVNBindingSystemIDsToGroupIDs(ctx, &req.Binding)
	if err != nil {
		return fmt.Errorf("converting VN binding system IDs to group IDs: %w", err)
	}

	// Establish exclusive access to the VN being modified. This is necessary because the API does not
	// support partial updates to the VN binding info. We need to retrieve the current binding info,
	// modify it, and then send the entire set back to the API. If multiple goroutines were to do this
	// concurrently, they could overwrite each other's changes. This mutex constrains only the current
	// client object, so it provides a limited but still useful level of protection.
	o.client.lock(o.blueprintId.String() + "_" + req.VNID)
	defer o.client.unlock(o.blueprintId.String() + "_" + req.VNID)

	// Collect the current VN info (bindings and SVI IPs) so we can merge in our new info.
	vn, err := o.GetVirtualNetwork(ctx, req.VNID)
	if err != nil {
		return fmt.Errorf("failed to fetch current bindings for VN %s: %w", req.VNID, err)
	}

	// Convert received VN bindings to a map.
	bindings := make(map[string]datacenter.VNBinding, len(vn.Bindings)+1)
	for _, binding := range vn.Bindings {
		bindings[binding.SystemID] = binding
	}

	if _, ok := bindings[req.Binding.SystemID]; !ok && req.Unbind {
		// Caller wants to un-bind and the binding doesn't exist. This is a no-op, so we can return early.
		return nil
	}

	// Merge (or delete) our binding as appropriate.
	if req.Unbind {
		delete(bindings, req.Binding.SystemID)
	} else {
		bindings[req.Binding.SystemID] = req.Binding
	}

	// Convert received SVI IPs to a map and merge (or delete) SVI info as necessary.
	sviIPs := make(map[string]datacenter.SVIAddressing, len(vn.SVIIPs)+1)
	for _, sviIP := range vn.SVIIPs {
		sviIPs[sviIP.SystemID] = sviIP
	}
	if req.SVIIP != nil {
		sviIPs[leafID] = *req.SVIIP
	}
	if req.Unbind {
		// Delete SVI info for our leaf switch.
		delete(sviIPs, leafID)

		// Delete SVI info for the peer leaf switch if it exists.
		_, peerLeafID, err := o.GetRedundancyGroupBySystem(ctx, leafID)
		if err != nil {
			return fmt.Errorf("getting redundancy group for leaf system id %q: %w", leafID, err)
		}
		if peerLeafID != nil {
			delete(sviIPs, *peerLeafID)
		}
	}

	vn.Bindings = slices.Collect(maps.Values(bindings))
	vn.SVIIPs = slices.Collect(maps.Values(sviIPs))

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPatch,
		urlStr: fmt.Sprintf(urls.DatacenterVirtualNetworkByID, o.blueprintId, req.VNID),
		apiInput: struct {
			Bindings []datacenter.VNBinding     `json:"bound_to"`
			SVIIPs   []datacenter.SVIAddressing `json:"svi_ips"`
		}{
			Bindings: slices.Collect(maps.Values(bindings)),
			SVIIPs:   slices.Collect(maps.Values(sviIPs)),
		},
	})
	if err != nil {
		return fmt.Errorf("updating VN %q: %w", req.VNID, err)
	}

	return nil
}

// convertVNBindingSystemIDsToGroupIDs converts the leaf SystemID and AccessSwitchNodeIDs
// in the VNBinding to their corresponding redundancy group IDs, if applicable.
func (o *TwoStageL3ClosClient) convertVNBindingSystemIDsToGroupIDs(ctx context.Context, b *datacenter.VNBinding) error {
	groupID, _, err := o.GetRedundancyGroupBySystem(ctx, b.SystemID)
	if err != nil {
		return fmt.Errorf("getting redundancy group for leaf system id %q: %w", b.SystemID, err)
	}
	if groupID != nil {
		b.SystemID = *groupID // swap leaf system ID for leaf redundancy group ID
	}

	accessIDSet := make(map[string]struct{})
	for _, accessID := range b.AccessSwitchNodeIDs {
		groupID, _, err = o.GetRedundancyGroupBySystem(ctx, accessID)
		if err != nil {
			return fmt.Errorf("getting redundancy group for access system id %q: %w", accessID, err)
		}
		if groupID == nil {
			accessIDSet[accessID] = struct{}{}
		} else {
			accessIDSet[*groupID] = struct{}{}
		}
	}

	b.AccessSwitchNodeIDs = slices.Collect(maps.Keys(accessIDSet))

	return nil
}
