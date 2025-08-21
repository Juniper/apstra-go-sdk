// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"
)

type VirtualNetworkBindingsRequest struct {
	VnId               ObjectId
	VnBindings         map[ObjectId]*VnBinding
	SviIps             map[ObjectId]*SviIp
	DhcpServiceEnabled *DhcpServiceEnabled
}

// SetVirtualNetworkLeafBindings replaces the current set of SVI and
// Binding info on the Virtual Network identified by vnId.
func (o *TwoStageL3ClosClient) SetVirtualNetworkLeafBindings(ctx context.Context, req VirtualNetworkBindingsRequest) error {
	var i int

	// turn the supplied map into a slice for the API, with some sanity checking along the way
	i = 0
	boundTo := make([]VnBinding, len(req.VnBindings))
	for k, v := range req.VnBindings {
		if v == nil {
			return fmt.Errorf("vbMap[%s] is nil", k) // map entries should not be nil
		}
		if k != v.SystemId {
			return fmt.Errorf("vbMap[%s] has system ID (%s)", k, v.SystemId) // map key should match payload
		}
		boundTo[i] = *v
		i++
	}

	// turn the supplied map into a slice for the API, with some sanity checking along the way
	i = 0
	sviIps := make([]SviIp, len(req.SviIps))
	for k, v := range req.SviIps {
		if v == nil {
			return fmt.Errorf("siMap[%s] is nil", k) // map entries should not be nil
		}
		if k != v.SystemId {
			return fmt.Errorf("siMap[%s] has system ID (%s)", k, v.SystemId) // map key should match payload
		}
		if _, ok := req.VnBindings[k]; !ok {
			// this is the big one - check that SVI info represents an active VN binding
			return fmt.Errorf("SVI requested for system %[1]s but %[1]s not among bound leaf IDs", k)
		}
		sviIps[i] = *v
		i++
	}

	o.client.Lock(o.blueprintId.String() + "_" + req.VnId.String())
	defer o.client.Unlock(o.blueprintId.String() + "_" + req.VnId.String())

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPatch,
		urlStr: fmt.Sprintf(apiUrlVirtualNetworkById, o.blueprintId, req.VnId),
		apiInput: struct {
			SviIps             []SviIp             `json:"svi_ips"`
			BoundTo            []VnBinding         `json:"bound_to"`
			DhcpServiceEnabled *DhcpServiceEnabled `json:"dhcp_service,omitempty"`
		}{
			SviIps:             sviIps,
			BoundTo:            boundTo,
			DhcpServiceEnabled: req.DhcpServiceEnabled,
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
	for k, v := range req.VnBindings {
		if v == nil {
			continue
		}

		if k != v.SystemId {
			return fmt.Errorf("vbMap[%s] has system ID (%s)", k, v.SystemId)
		}
	}

	for k, v := range req.SviIps {
		if v == nil {
			continue
		}

		if k != v.SystemId {
			return fmt.Errorf("siMap[%s] has system ID (%s)", k, v.SystemId)
		}
		if binding := req.VnBindings[k]; binding == nil {
			return fmt.Errorf("SVI requested for system %[1]s but %[1]s not among bound leaf IDs, or is being removed", k)
		}
	}

	o.client.Lock(o.blueprintId.String() + "_" + req.VnId.String())
	defer o.client.Unlock(o.blueprintId.String() + "_" + req.VnId.String())

	// collect current VN info (bindings and SVI IPs)
	vn, err := o.GetVirtualNetwork(ctx, req.VnId)
	if err != nil {
		return fmt.Errorf("update VN binding info failed to fetch current bindings - %w", err)
	}

	// drop any current binding which appears in the caller's binding map
	for i := len(vn.Data.VnBindings) - 1; i >= 0; i-- { // loop backward over current VN bindings
		if _, ok := req.VnBindings[vn.Data.VnBindings[i].SystemId]; ok {
			// current VN passed by caller; delete it; we'll use caller's data instead
			vn.Data.VnBindings[i] = vn.Data.VnBindings[len(vn.Data.VnBindings)-1] // copy last to index i
			vn.Data.VnBindings = vn.Data.VnBindings[:len(vn.Data.VnBindings)-1]   // trim last from slice
		}
	}

	// drop any current SVI IP referencing a binding mentioned in the caller's binding map
	for i := len(vn.Data.SviIps) - 1; i >= 0; i-- { // loop backward over current SVI IP info
		if _, ok := req.VnBindings[vn.Data.SviIps[i].SystemId]; ok {
			// current SVI IP info relates to VN passed by caller; delete it; we'll use caller's data instead
			vn.Data.SviIps[i] = vn.Data.SviIps[len(vn.Data.SviIps)-1] // copy last to index i
			vn.Data.SviIps = vn.Data.SviIps[:len(vn.Data.SviIps)-1]   // trim last from slice
		}
	}

	// copy non-nil (non-delete) bindings from the caller's map to the binding slice we'll send back
	for _, binding := range req.VnBindings {
		if binding != nil {
			vn.Data.VnBindings = append(vn.Data.VnBindings, *binding)
		}
	}

	// copy non-nil (non-delete) SVI IPs from the caller's map to the binding slice we'll send back
	for _, sviIp := range req.SviIps {
		if sviIp != nil {
			vn.Data.SviIps = append(vn.Data.SviIps, *sviIp)
		}
	}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPatch,
		urlStr: fmt.Sprintf(apiUrlVirtualNetworkById, o.blueprintId, req.VnId),
		apiInput: struct {
			SviIps             []SviIp             `json:"svi_ips"`
			BoundTo            []VnBinding         `json:"bound_to"`
			DhcpServiceEnabled *DhcpServiceEnabled `json:"dhcp_service,omitempty"`
		}{
			SviIps:             vn.Data.SviIps,
			BoundTo:            vn.Data.VnBindings,
			DhcpServiceEnabled: req.DhcpServiceEnabled,
		},
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
