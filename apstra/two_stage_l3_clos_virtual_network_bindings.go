// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/internal/urls"
)

type VirtualNetworkBindingsRequest struct {
	VirtualNetworkID   string
	Bindings           map[string]*datacenter.VNBinding
	SVIAddressing      map[string]*datacenter.SVIAddressing
	DHCPServiceEnabled *datacenter.DHCPServiceEnabled
}

// SetVirtualNetworkLeafBindings replaces the current set of SVI and
// Binding info on the Virtual Network identified by vnId.
func (o *TwoStageL3ClosClient) SetVirtualNetworkLeafBindings(ctx context.Context, req VirtualNetworkBindingsRequest) error {
	var i int

	// turn the supplied map into a slice for the API, with some sanity checking along the way
	i = 0
	boundTo := make([]datacenter.VNBinding, len(req.Bindings))
	for k, v := range req.Bindings {
		if v == nil {
			return fmt.Errorf("vbMap[%s] is nil", k) // map entries should not be nil
		}
		if k != v.SystemID {
			return fmt.Errorf("vbMap[%s] has system ID (%s)", k, v.SystemID) // map key should match payload
		}
		boundTo[i] = *v
		i++
	}

	// turn the supplied map into a slice for the API, with some sanity checking along the way
	i = 0
	sviIps := make([]datacenter.SVIAddressing, len(req.SVIAddressing))
	for k, v := range req.SVIAddressing {
		if v == nil {
			return fmt.Errorf("siMap[%s] is nil", k) // map entries should not be nil
		}
		if k != v.SystemID {
			return fmt.Errorf("siMap[%s] has system ID (%s)", k, v.SystemID) // map key should match payload
		}
		groupID, _, err := o.GetRedundancyGroupBySystem(ctx, k)
		if err != nil {
			return fmt.Errorf("getting group info for system ID %q: %w", k, err)
		}

		// Check that SVI info represents an active VN binding.
		var systemBound, groupBound bool
		_, systemBound = req.Bindings[k]
		if groupID != nil {
			_, groupBound = req.Bindings[*groupID]
		}
		if !systemBound && !groupBound {
			return fmt.Errorf("SVI requested for system %[1]s but %[1]s not among bound leaf IDs", k)
		}
		sviIps[i] = *v
		i++
	}

	o.client.lock(o.blueprintId.String() + "_" + req.VirtualNetworkID)
	defer o.client.unlock(o.blueprintId.String() + "_" + req.VirtualNetworkID)

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPatch,
		urlStr: fmt.Sprintf(urls.DatacenterVirtualNetworkByID, o.blueprintId, req.VirtualNetworkID),
		apiInput: struct {
			SviIps             []datacenter.SVIAddressing     `json:"svi_ips"`
			BoundTo            []datacenter.VNBinding         `json:"bound_to"`
			DhcpServiceEnabled *datacenter.DHCPServiceEnabled `json:"dhcp_service,omitempty"`
		}{
			SviIps:             sviIps,
			BoundTo:            boundTo,
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
func (o *TwoStageL3ClosClient) UpdateVirtualNetworkLeafBindings(ctx context.Context, req VirtualNetworkBindingsRequest) error {
	for k, v := range req.Bindings {
		if v == nil {
			continue
		}

		if k != v.SystemID {
			return fmt.Errorf("entry at Bindings[%s] has system ID (%s)", k, v.SystemID)
		}
	}

	for k, v := range req.SVIAddressing {
		if v == nil {
			continue
		}

		if k != v.SystemID {
			return fmt.Errorf("entry at SVIAddressing[%s] has system ID (%s)", k, v.SystemID)
		}

		// Discover the Leaf Redudnancy Group ID, if any.
		groupID, _, err := o.GetRedundancyGroupBySystem(ctx, k)
		if err != nil {
			return fmt.Errorf("getting group info for system ID %q: %w", k, err)
		}

		// Discover whether the request is trying to bind or unbind the Leaf Switch and/or the its Redundancy Group.
		var systemBound, systemToBeUnbound, groupBound, groupToBeUnbound bool
		var val interface{}
		val, systemBound = req.Bindings[k]
		systemToBeUnbound = val == nil
		if groupID != nil {
			val, groupBound = req.Bindings[*groupID]
			groupToBeUnbound = val == nil
		}

		if (!systemBound && !groupBound) || systemToBeUnbound || groupToBeUnbound {
			return fmt.Errorf("SVI requested for system %[1]s but %[1]s not among bound leaf IDs, or is being removed", k)
		}
	}

	o.client.lock(o.blueprintId.String() + "_" + req.VirtualNetworkID)
	defer o.client.unlock(o.blueprintId.String() + "_" + req.VirtualNetworkID)

	// collect current VN info (bindings and SVI IPs)
	vn, err := o.GetVirtualNetwork(ctx, req.VirtualNetworkID)
	if err != nil {
		return fmt.Errorf("update VN binding info failed to fetch current bindings - %w", err)
	}

	// drop any current binding which appears in the caller's binding map
	for i := len(vn.Bindings) - 1; i >= 0; i-- { // loop backward over current VN bindings
		if _, ok := req.Bindings[vn.Bindings[i].SystemID]; ok {
			// current VN passed by caller; delete it; we'll use caller's data instead
			vn.Bindings[i] = vn.Bindings[len(vn.Bindings)-1] // copy last to index i
			vn.Bindings = vn.Bindings[:len(vn.Bindings)-1]   // trim last from slice
		}
	}

	// drop any current SVI IP referencing a binding mentioned in the caller's binding map
	for i := len(vn.SVIIPs) - 1; i >= 0; i-- { // loop backward over current SVI IP info
		if _, ok := req.Bindings[vn.SVIIPs[i].SystemID]; ok {
			// current SVI IP info relates to VN passed by caller; delete it; we'll use caller's data instead
			vn.SVIIPs[i] = vn.SVIIPs[len(vn.SVIIPs)-1] // copy last to index i
			vn.SVIIPs = vn.SVIIPs[:len(vn.SVIIPs)-1]   // trim last from slice
		}
	}

	// copy non-nil (non-delete) bindings from the caller's map to the binding slice we'll send back
	for _, binding := range req.Bindings {
		if binding != nil {
			vn.Bindings = append(vn.Bindings, *binding)
		}
	}

	// copy non-nil (non-delete) SVI IPs from the caller's map to the binding slice we'll send back
	for _, sviIp := range req.SVIAddressing {
		if sviIp != nil {
			vn.SVIIPs = append(vn.SVIIPs, *sviIp)
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
			SviIps:             vn.SVIIPs,
			BoundTo:            vn.Bindings,
			DhcpServiceEnabled: req.DHCPServiceEnabled,
		},
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
