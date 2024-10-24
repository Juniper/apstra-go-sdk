// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
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

	type testCase struct {
		d VirtualNetworkData
		e string
	}

	testCases := map[string]testCase{
		"a": {
			d: VirtualNetworkData{
				Ipv4Enabled:    true,
				Ipv4Subnet:     mustParseIpNet(t, "192.0.2.0/24"),
				Ipv6Enabled:    true,
				Ipv6Subnet:     mustParseIpNet(t, "3fff::/64"),
				L3Mtu:          toPtr(9000),
				Label:          "a",
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
						AccessSwitchNodeIds: []ObjectId{},
						SystemId:            "UJoJhK-jXJkc5Mtarc8",
						VlanId:              toPtr(Vlan(3)),
					},
				},
				VnType: enum.VnTypeVxlan,
			},
			e: `{
                  "dhcp_service": "dhcpServiceDisabled",
                  "ipv4_enabled": true,
                  "ipv4_subnet": "192.0.2.0/24",
                  "ipv6_enabled": true,
                  "ipv6_subnet": "3fff::/64",
                  "l3_mtu": 9000,
                  "label": "a",
                  "rt_policy": null,
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
                      "access_switch_node_ids": [],
                      "system_id": "UJoJhK-jXJkc5Mtarc8",
                      "vlan_id": 3
                    }
                  ],
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
