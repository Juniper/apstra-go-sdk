// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package deepcopy

import (
	"slices"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
)

func VirtualNetwork(in datacenter.VirtualNetwork) datacenter.VirtualNetwork {
	out := in

	out.IPv4Subnet = cloneIPNet(in.IPv4Subnet)
	out.IPv6Subnet = cloneIPNet(in.IPv6Subnet)
	out.VirtualGatewayIPv4 = slices.Clone(in.VirtualGatewayIPv4)
	out.VirtualGatewayIPv6 = slices.Clone(in.VirtualGatewayIPv6)
	out.VirtualMAC = slices.Clone(in.VirtualMAC)

	out.Tags = slices.Clone(in.Tags)
	out.Bindings = cloneVNBindings(in.Bindings)
	out.SVIIPs = cloneSVIAddressings(in.SVIIPs)
	out.VNI = pointer.ToCopyOfValue(in.VNI)
	out.L3MTU = pointer.ToCopyOfValue(in.L3MTU)
	out.ReservedVLAN = pointer.ToCopyOfValue(in.ReservedVLAN)
	out.RTPolicy = cloneRTPolicy(in.RTPolicy)

	return out
}

func cloneRTPolicy(in *datacenter.RTPolicy) *datacenter.RTPolicy {
	if in == nil {
		return nil
	}

	out := *in
	out.ImportRTs = slices.Clone(in.ImportRTs)
	out.ExportRTs = slices.Clone(in.ExportRTs)
	return &out
}

func cloneVNBindings(in []datacenter.VNBinding) []datacenter.VNBinding {
	if in == nil {
		return nil
	}

	out := make([]datacenter.VNBinding, len(in))
	for i := range in {
		out[i] = in[i]
		out[i].AccessSwitchNodeIDs = slices.Clone(in[i].AccessSwitchNodeIDs)
		out[i].VLAN = pointer.ToCopyOfValue(in[i].VLAN)
	}
	return out
}

func cloneSVIAddressings(in []datacenter.SVIAddressing) []datacenter.SVIAddressing {
	if in == nil {
		return nil
	}

	out := make([]datacenter.SVIAddressing, len(in))
	for i := range in {
		out[i] = in[i]
		out[i].IPv4Addr = cloneIPNet(in[i].IPv4Addr)
		out[i].IPv6Addr = cloneIPNet(in[i].IPv6Addr)
	}
	return out
}
