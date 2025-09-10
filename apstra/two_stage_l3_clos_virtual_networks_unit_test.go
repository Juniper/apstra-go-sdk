// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func TestTwoStageL3ClosVirtualNetworkStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		int() int
	}

	type apiIotaString interface {
		parse() (int, error)
		string() string
	}

	type stringTestData struct {
		stringVal  string
		intType    apiStringIota
		stringType apiIotaString
	}
	testData := []stringTestData{
		{stringVal: "", intType: SystemRoleNone, stringType: systemRoleNone},
		{stringVal: "access", intType: SystemRoleAccess, stringType: systemRoleAccess},
		{stringVal: "leaf", intType: SystemRoleLeaf, stringType: systemRoleLeaf},
	}

	for i, td := range testData {
		ii := td.intType.int()
		is := td.intType.String()
		sp, err := td.stringType.parse()
		if err != nil {
			t.Fatal(err)
		}
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.int() != sp ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}

func TestVirtualNetworkDataMarshalJson(t *testing.T) {
	mustParseIpNet := func(t testing.TB, s string) *net.IPNet {
		t.Helper()
		result, err := ipNetFromString(s)
		require.NoError(t, err)
		return result
	}

	mustParseIp := func(t testing.TB, s string) net.IP {
		t.Helper()
		result, err := ipFromString(s)
		require.NoError(t, err)
		return result
	}

	mustParseMac := func(t testing.TB, s string) net.HardwareAddr {
		t.Helper()
		result, err := macFromString(s)
		require.NoError(t, err)
		return result
	}

	type testCase struct {
		d VirtualNetworkData
		e string
	}

	testCases := map[string]testCase{
		"full_detail_vlan": {
			d: VirtualNetworkData{
				DhcpService:    true,
				Ipv4Enabled:    true,
				Ipv4Subnet:     mustParseIpNet(t, "192.0.2.0/24"),
				Ipv6Enabled:    true,
				Ipv6Subnet:     mustParseIpNet(t, "3fff::/64"),
				L3Mtu:          toPtr(9010),
				Label:          "a",
				ReservedVlanId: toPtr(Vlan(10)),
				SecurityZoneId: "dtUF3UAr4Cqfuoy6iII",
				SviIps: []SviIp{
					{
						SystemId: "UJoJhK-jXJkc5Mtarc8",
						Ipv4Addr: mustParseIpNet(t, "192.0.2.2/24"),
						Ipv4Mode: enum.SviIpv4ModeEnabled,
						Ipv6Addr: mustParseIpNet(t, "3fff::2/64"),
						Ipv6Mode: enum.SviIpv6ModeEnabled,
					},
				},
				VirtualGatewayIpv4:        mustParseIp(t, "192.0.2.1"),
				VirtualGatewayIpv6:        mustParseIp(t, "3fff::1"),
				VirtualGatewayIpv4Enabled: true,
				VirtualGatewayIpv6Enabled: true,
				VnBindings: []VnBinding{
					{
						AccessSwitchNodeIds: []ObjectId{"tFlHPPD766lj8g8PYsqw", "Ik82Xta17zkWGHNw6pbN"},
						SystemId:            "UJoJhK-jXJkc5Mtarc8",
						VlanId:              toPtr(Vlan(10)),
					},
				},
				VnId:       toPtr(VNI(10 * 1000)),
				VnType:     enum.VnTypeVlan,
				VirtualMac: mustParseMac(t, "08:00:20:01:02:03"),
			},
			e: `{
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
			d: VirtualNetworkData{
				DhcpService:    false,
				Ipv4Enabled:    false,
				Ipv4Subnet:     nil,
				Ipv6Enabled:    false,
				Ipv6Subnet:     nil,
				L3Mtu:          toPtr(9010),
				Label:          "a",
				ReservedVlanId: toPtr(Vlan(10)),
				SecurityZoneId: "dtUF3UAr4Cqfuoy6iII",
				SviIps: []SviIp{
					{
						SystemId: "UJoJhK-jXJkc5Mtarc8",
						Ipv4Addr: mustParseIpNet(t, "192.0.2.2/24"),
						Ipv4Mode: enum.SviIpv4ModeEnabled,
						Ipv6Addr: mustParseIpNet(t, "3fff::2/64"),
						Ipv6Mode: enum.SviIpv6ModeEnabled,
					},
					{
						SystemId: "iEbRhCoGzNlIgfO9DZ4r",
						Ipv4Addr: mustParseIpNet(t, "192.0.2.3/24"),
						Ipv4Mode: enum.SviIpv4ModeEnabled,
						Ipv6Addr: mustParseIpNet(t, "3fff::3/64"),
						Ipv6Mode: enum.SviIpv6ModeEnabled,
					},
				},
				VirtualGatewayIpv4:        mustParseIp(t, "192.0.2.1"),
				VirtualGatewayIpv6:        mustParseIp(t, "3fff::1"),
				VirtualGatewayIpv4Enabled: true,
				VirtualGatewayIpv6Enabled: true,
				VnBindings: []VnBinding{
					{
						AccessSwitchNodeIds: []ObjectId{"tFlHPPD766lj8g8PYsqw", "Ik82Xta17zkWGHNw6pbN"},
						SystemId:            "UJoJhK-jXJkc5Mtarc8",
						VlanId:              toPtr(Vlan(10)),
					},
					{
						AccessSwitchNodeIds: []ObjectId{"dxodxTyx6SAlMYbP45Bp", "7ZNWB3KlYGlxj87UdeaF", "v6mwfD43FJP3mhhvF7YN"},
						SystemId:            "iEbRhCoGzNlIgfO9DZ4r",
						VlanId:              toPtr(Vlan(10)),
					},
				},
				VnId:       toPtr(VNI(10 * 1000)),
				VnType:     enum.VnTypeVxlan,
				VirtualMac: mustParseMac(t, "08:00:20:01:02:03"),
			},
			e: `{
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

			var fc VirtualNetwork // fc: full circle test result
			err = json.Unmarshal(a, &fc)
			require.NoError(t, err)
			compareVirtualNetworkData(t, &tCase.d, fc.Data, true)
		})
	}
}
