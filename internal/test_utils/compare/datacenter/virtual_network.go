// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedatacenter

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func VirtualNetwork(t testing.TB, req, resp datacenter.VirtualNetwork, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Virtual Network")

	// Convert slices of Bindings to maps for easy comparison.
	reqVNBindings := make(map[string]datacenter.VNBinding, len(req.Bindings))
	for _, vnBinding := range req.Bindings {
		reqVNBindings[vnBinding.SystemID] = vnBinding
	}
	respVNBindings := make(map[string]datacenter.VNBinding, len(resp.Bindings))
	for _, vnBinding := range resp.Bindings {
		respVNBindings[vnBinding.SystemID] = vnBinding
	}
	// Ensure both maps have the same keys.
	compare.SlicesAsSets(t, slices.Collect(maps.Keys(reqVNBindings)), slices.Collect(maps.Keys(respVNBindings)), "VN Bindings Leaf IDs")
	// Compare each map entry.
	for leafID := range reqVNBindings {
		VNBinding(t, reqVNBindings[leafID], respVNBindings[leafID], fmt.Sprintf("VN Binding for leaf %q", leafID))
	}

	if req.Description != nil {
		require.NotNil(t, resp.Description)
		require.Equal(t, *req.Description, *resp.Description)
	}
	if req.DHCPService != nil && len(req.Bindings) > 0 {
		require.NotNil(t, resp.DHCPService)
		require.Equal(t, *req.DHCPService, *resp.DHCPService) // only checked with bindings because this info is lost otherwise with releases prior to 6.2.0
	}

	if req.EncapsulateInnerVLAN != nil {
		require.NotNil(t, resp.EncapsulateInnerVLAN)
		require.Equal(t, *req.EncapsulateInnerVLAN, *resp.EncapsulateInnerVLAN)
	}

	require.NotNil(t, resp.IPv4Enabled)
	if req.IPv4Enabled != nil {
		require.Equal(t, *req.IPv4Enabled, *resp.IPv4Enabled)
	}
	if req.IPv4Subnet != nil {
		require.NotNil(t, resp.IPv4Subnet)
		require.Equal(t, req.IPv4Subnet.String(), resp.IPv4Subnet.String())
	}

	require.NotNil(t, resp.IPv6Enabled)
	if req.IPv6Enabled != nil {
		require.Equal(t, *req.IPv6Enabled, *resp.IPv6Enabled)
	}
	if req.IPv6Subnet != nil {
		require.NotNil(t, resp.IPv6Subnet)
		require.Equal(t, req.IPv6Subnet.String(), resp.IPv6Subnet.String())
	}

	require.NotNil(t, req.Label)
	require.NotNil(t, resp.Label)
	require.Equal(t, *req.Label, *resp.Label)

	if req.SecurityZoneID != "" {
		require.Equal(t, req.SecurityZoneID, resp.SecurityZoneID)
	}
	if req.SwitchingZoneID != "" {
		require.Equal(t, req.SwitchingZoneID, resp.SwitchingZoneID)
	}
	if req.Type.Value != "" {
		require.Equal(t, req.Type, resp.Type)
	}
	if req.VirtualGatewayIPv4 != nil {
		require.NotNil(t, resp.VirtualGatewayIPv4)
		require.Equal(t, req.VirtualGatewayIPv4.String(), resp.VirtualGatewayIPv4.String())
	}
	if req.VirtualGatewayIPv6 != nil {
		require.NotNil(t, resp.VirtualGatewayIPv6)
		require.Equal(t, req.VirtualGatewayIPv6.String(), resp.VirtualGatewayIPv6.String())
	}
	require.NotNil(t, resp.VirtualGatewayIPv4Enabled)
	if req.VirtualGatewayIPv4 != nil {
		require.Equal(t, *req.VirtualGatewayIPv4Enabled, *resp.VirtualGatewayIPv4Enabled)
	}
	require.NotNil(t, resp.VirtualGatewayIPv6Enabled)
	if req.VirtualGatewayIPv6 != nil {
		require.Equal(t, *req.VirtualGatewayIPv6Enabled, *resp.VirtualGatewayIPv6Enabled)
	}
	if req.VirtualMAC != nil {
		require.NotNil(t, resp.VirtualMAC)
		require.Equal(t, req.VirtualMAC, resp.VirtualMAC)
	}

	require.NotNil(t, resp.L3MTU)
	if req.L3MTU != nil {
		require.Equal(t, req.L3MTU, resp.L3MTU)
	}

	if req.ReservedVLAN != nil {
		if *req.ReservedVLAN == 0 { // zero is a special case, it clears the reserved value
			require.Nil(t, resp.ReservedVLAN)
		} else {
			require.NotNil(t, resp.ReservedVLAN)
			require.Equal(t, req.ReservedVLAN, resp.ReservedVLAN)
		}
	}

	if req.RTPolicy != nil && (req.RTPolicy.ExportRTs != nil || req.RTPolicy.ImportRTs != nil) {
		require.NotNil(t, req.RTPolicy)
		require.NotNil(t, resp.RTPolicy)
		RTPolicy(t, *req.RTPolicy, *resp.RTPolicy)
	}

	// Convert slices of SVI Addresses to maps for easy comparison.
	reqSVIAddressingMap := make(map[string]datacenter.SVIAddressing, len(req.SVIIPs)) // Collect request leaf IDs and build a map keyed by ID
	for _, sviIP := range req.SVIIPs {
		reqSVIAddressingMap[sviIP.SystemID] = sviIP
	}
	respSVIAddressingMap := make(map[string]datacenter.SVIAddressing, len(resp.SVIIPs)) // Collect response leaf IDs and build a map keyed by ID
	for _, sviIP := range req.SVIIPs {
		respSVIAddressingMap[sviIP.SystemID] = sviIP
	}
	// Ensure both maps have the same keys.
	compare.SlicesAsSets(t, slices.Collect(maps.Keys(reqSVIAddressingMap)), slices.Collect(maps.Keys(respSVIAddressingMap)), "SVI Addressing Leaf IDs")
	// Compare each map entry.
	for leafID := range reqSVIAddressingMap {
		SVIAddressing(t, reqSVIAddressingMap[leafID], respSVIAddressingMap[leafID], fmt.Sprintf("SVI Addressing for leaf %q", leafID))
	}

	compare.SlicesAsSets(t, req.Tags, req.Tags, "Tags")

	if req.VNI != nil {
		require.NotNil(t, resp.VNI)
		require.Equal(t, req.VNI, resp.VNI)
	}
}
