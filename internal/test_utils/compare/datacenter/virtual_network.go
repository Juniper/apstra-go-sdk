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

	require.Equal(t, req.Description, resp.Description)
	if len(req.Bindings) > 0 {
		require.Equal(t, req.DHCPService, resp.DHCPService) // only checked with bindings because this info is lost otherwise
	}
	require.Equal(t, req.IPv4Enabled, resp.IPv4Enabled)
	require.Equal(t, req.IPv4Subnet.String(), resp.IPv4Subnet.String())
	require.Equal(t, req.IPv6Enabled, resp.IPv6Enabled)
	require.Equal(t, req.IPv6Subnet.String(), resp.IPv6Subnet.String())
	require.Equal(t, req.Label, resp.Label)
	if req.SecurityZoneID != "" {
		require.Equal(t, req.SecurityZoneID, resp.SecurityZoneID)
	}
	require.Equal(t, req.Type, resp.Type)
	require.Equal(t, req.VirtualGatewayIPv4.String(), resp.VirtualGatewayIPv4.String())
	require.Equal(t, req.VirtualGatewayIPv6.String(), resp.VirtualGatewayIPv6.String())
	require.Equal(t, req.VirtualGatewayIPv4Enabled, resp.VirtualGatewayIPv4Enabled)
	require.Equal(t, req.VirtualGatewayIPv6Enabled, resp.VirtualGatewayIPv6Enabled)
	require.Equal(t, req.VirtualMAC, resp.VirtualMAC)

	require.NotNil(t, resp.L3MTU)
	if req.L3MTU != nil {
		require.Equal(t, req.L3MTU, resp.L3MTU)
	}

	if req.ReservedVLAN != nil || resp.ReservedVLAN != nil {
		require.NotNil(t, req.ReservedVLAN)
		require.NotNil(t, resp.ReservedVLAN)
		require.Equal(t, req.ReservedVLAN, resp.ReservedVLAN)
	}

	if req.RTPolicy != nil || resp.RTPolicy != nil {
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
