// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter_test

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/parse"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	comparedatacenter "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/datacenter"
	"github.com/stretchr/testify/require"
)

func TestVirtualNetworkData_MarshalJson(t *testing.T) {
	mustParseIpNet := func(t testing.TB, s string) *net.IPNet {
		t.Helper()
		result, err := parse.IPNetFromString(s)
		require.NoError(t, err)
		return result
	}

	mustParseIp := func(t testing.TB, s string) net.IP {
		t.Helper()
		result, err := parse.IPFromString(s)
		require.NoError(t, err)
		return result
	}

	mustParseMac := func(t testing.TB, s string) net.HardwareAddr {
		t.Helper()
		result, err := parse.MACFromString(s)
		require.NoError(t, err)
		return result
	}

	type testCase struct {
		d datacenter.VirtualNetwork
		e string
	}

	testCases := map[string]testCase{
		"full_detail_vlan": {
			d: datacenter.VirtualNetwork{
				Description:    "My Virtual Network",
				DHCPService:    true,
				IPv4Enabled:    true,
				IPv4Subnet:     mustParseIpNet(t, "192.0.2.0/24"),
				IPv6Enabled:    true,
				IPv6Subnet:     mustParseIpNet(t, "3fff::/64"),
				L3MTU:          pointer.To(9010),
				Label:          "a",
				ReservedVLAN:   pointer.To(uint16(10)),
				SecurityZoneID: "dtUF3UAr4Cqfuoy6iII",
				SVIIPs: []datacenter.SVIAddressing{
					{
						SystemID: "UJoJhK-jXJkc5Mtarc8",
						IPv4Addr: mustParseIpNet(t, "192.0.2.2/24"),
						IPv4Mode: enum.IPv4SVIModeEnabled,
						IPv6Addr: mustParseIpNet(t, "3fff::2/64"),
						IPv6Mode: enum.IPv6SVIModeEnabled,
					},
				},
				Tags:                      []string{"a", "b"},
				VirtualGatewayIPv4:        mustParseIp(t, "192.0.2.1"),
				VirtualGatewayIPv6:        mustParseIp(t, "3fff::1"),
				VirtualGatewayIPv4Enabled: true,
				VirtualGatewayIPv6Enabled: true,
				Bindings: []datacenter.VNBinding{
					{
						AccessSwitchNodeIDs: []string{"tFlHPPD766lj8g8PYsqw", "Ik82Xta17zkWGHNw6pbN"},
						SystemID:            "UJoJhK-jXJkc5Mtarc8",
						VLAN:                pointer.To(uint16(10)),
					},
				},
				VNI:        pointer.To(uint32(10 * 1000)),
				Type:       enum.VnTypeVlan,
				VirtualMAC: mustParseMac(t, "08:00:20:01:02:03"),
			},
			e: `{
                  "description": "My Virtual Network",
                  "dhcp_service": "dhcpServiceEnabled",
                  "ipv4_enabled": true,
                  "ipv4_subnet": "192.0.2.0/24",
                  "ipv6_enabled": true,
                  "ipv6_subnet": "3fff::/64",
                  "l3_mtu": 9010,
                  "label": "a",
                  "rt_policy": null,
                  "reserved_vlan_id": 10,
                  "security_zone_id": "dtUF3UAr4Cqfuoy6iII",
                  "svi_ips": [
                    {
                      "ipv4_addr": "192.0.2.2/24",
                      "ipv4_mode": "enabled",
                      "ipv6_addr": "3fff::2/64",
                      "ipv6_mode": "enabled",
                      "system_id": "UJoJhK-jXJkc5Mtarc8"
                    }
                  ],
                  "tags": ["a", "b"],
                  "virtual_gateway_ipv4": "192.0.2.1",
                  "virtual_gateway_ipv6": "3fff::1",
                  "virtual_gateway_ipv6_enabled": true,
                  "virtual_gateway_ipv4_enabled": true,
                  "bound_to": [
                    {
                      "access_switch_node_ids": [
                        "tFlHPPD766lj8g8PYsqw",
                        "Ik82Xta17zkWGHNw6pbN"
                      ],
                      "system_id": "UJoJhK-jXJkc5Mtarc8",
                      "vlan_id": 10
                    }
                  ],
                  "virtual_mac": "08:00:20:01:02:03",
                  "vn_id": "10000",
                  "vn_type": "vlan"
                }`,
		},
		"full_detail_vxlan": {
			d: datacenter.VirtualNetwork{
				Description:    "My Virtual Network",
				DHCPService:    false,
				IPv4Enabled:    false,
				IPv4Subnet:     nil,
				IPv6Enabled:    false,
				IPv6Subnet:     nil,
				L3MTU:          pointer.To(9010),
				Label:          "a",
				ReservedVLAN:   pointer.To(uint16(10)),
				SecurityZoneID: "dtUF3UAr4Cqfuoy6iII",
				SVIIPs: []datacenter.SVIAddressing{
					{
						SystemID: "UJoJhK-jXJkc5Mtarc8",
						IPv4Addr: mustParseIpNet(t, "192.0.2.2/24"),
						IPv4Mode: enum.IPv4SVIModeEnabled,
						IPv6Addr: mustParseIpNet(t, "3fff::2/64"),
						IPv6Mode: enum.IPv6SVIModeEnabled,
					},
					{
						SystemID: "iEbRhCoGzNlIgfO9DZ4r",
						IPv4Addr: mustParseIpNet(t, "192.0.2.3/24"),
						IPv4Mode: enum.IPv4SVIModeEnabled,
						IPv6Addr: mustParseIpNet(t, "3fff::3/64"),
						IPv6Mode: enum.IPv6SVIModeEnabled,
					},
				},
				Tags:                      []string{"a", "b"},
				VirtualGatewayIPv4:        mustParseIp(t, "192.0.2.1"),
				VirtualGatewayIPv6:        mustParseIp(t, "3fff::1"),
				VirtualGatewayIPv4Enabled: true,
				VirtualGatewayIPv6Enabled: true,
				Bindings: []datacenter.VNBinding{
					{
						AccessSwitchNodeIDs: []string{"tFlHPPD766lj8g8PYsqw", "Ik82Xta17zkWGHNw6pbN"},
						SystemID:            "UJoJhK-jXJkc5Mtarc8",
						VLAN:                pointer.To(uint16(10)),
					},
					{
						AccessSwitchNodeIDs: []string{"dxodxTyx6SAlMYbP45Bp", "7ZNWB3KlYGlxj87UdeaF", "v6mwfD43FJP3mhhvF7YN"},
						SystemID:            "iEbRhCoGzNlIgfO9DZ4r",
						VLAN:                pointer.To(uint16(10)),
					},
				},
				VNI:        pointer.To(uint32(10 * 1000)),
				Type:       enum.VnTypeVxlan,
				VirtualMAC: mustParseMac(t, "08:00:20:01:02:03"),
			},
			e: `{
                  "description": "My Virtual Network",
                  "dhcp_service": "dhcpServiceDisabled",
                  "ipv4_enabled": false,
                  "ipv6_enabled": false,
                  "l3_mtu": 9010,
                  "label": "a",
                  "rt_policy": null,
                  "reserved_vlan_id": 10,
                  "security_zone_id": "dtUF3UAr4Cqfuoy6iII",
                  "svi_ips": [
                    {
                      "ipv4_addr": "192.0.2.2/24",
                      "ipv4_mode": "enabled",
                      "ipv6_addr": "3fff::2/64",
                      "ipv6_mode": "enabled",
                      "system_id": "UJoJhK-jXJkc5Mtarc8"
                    },
                    {
                      "ipv4_addr": "192.0.2.3/24",
                      "ipv4_mode": "enabled",
                      "ipv6_addr": "3fff::3/64",
                      "ipv6_mode": "enabled",
                      "system_id": "iEbRhCoGzNlIgfO9DZ4r"
                    }
                  ],
                  "tags": ["a", "b"],
                  "virtual_gateway_ipv4": "192.0.2.1",
                  "virtual_gateway_ipv6": "3fff::1",
                  "virtual_gateway_ipv6_enabled": true,
                  "virtual_gateway_ipv4_enabled": true,
                  "bound_to": [
                    {
                      "access_switch_node_ids": [
                        "tFlHPPD766lj8g8PYsqw",
                        "Ik82Xta17zkWGHNw6pbN"
                      ],
                      "system_id": "UJoJhK-jXJkc5Mtarc8",
                      "vlan_id": 10
                    },
                    {
                      "access_switch_node_ids": [
                        "dxodxTyx6SAlMYbP45Bp",
                        "7ZNWB3KlYGlxj87UdeaF",
                        "v6mwfD43FJP3mhhvF7YN"
                      ],
                      "system_id": "iEbRhCoGzNlIgfO9DZ4r",
                      "vlan_id": 10
                    }
                  ],
                  "virtual_mac": "08:00:20:01:02:03",
                  "vn_id": "10000",
                  "vn_type": "vxlan"
                }`,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			a, err := json.Marshal(tCase.d)
			require.NoError(t, err)
			require.JSONEq(t, tCase.e, string(a))

			var fc datacenter.VirtualNetwork // fc: full circle test result
			err = json.Unmarshal(a, &fc)
			require.NoError(t, err)
			comparedatacenter.VirtualNetwork(t, tCase.d, fc)
		})
	}
}
