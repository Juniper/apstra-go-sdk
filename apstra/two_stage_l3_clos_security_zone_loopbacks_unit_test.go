// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra_test

import (
	"encoding/json"
	"log"
	"net/netip"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/stretchr/testify/require"
)

func TestSecurityZoneLoopback_MarshalJSON(t *testing.T) {
	mustParsePrefixPtr := func(s string) *netip.Prefix {
		p := netip.MustParsePrefix(s)
		return &p
	}

	type testCase struct {
		data     apstra.SecurityZoneLoopback
		expected string
	}

	testCases := map[string]testCase{
		"both_have_values": {
			data: apstra.SecurityZoneLoopback{
				IPv4Addr: mustParsePrefixPtr("192.0.2.1/32"),
				IPv6Addr: mustParsePrefixPtr("3fff::/128"),
			},
			expected: `{"ipv4_addr":"192.0.2.1/32","ipv6_addr":"3fff::/128"}`,
		},
		"both_omit_values": {
			data:     apstra.SecurityZoneLoopback{},
			expected: `{}`,
		},
		"both_null_values": {
			data: apstra.SecurityZoneLoopback{
				IPv4Addr: &netip.Prefix{},
				IPv6Addr: &netip.Prefix{},
			},
			expected: `{"ipv4_addr":null,"ipv6_addr":null}`,
		},
		"omit_ipv4_value_ipv6": {
			data: apstra.SecurityZoneLoopback{
				IPv6Addr: mustParsePrefixPtr("3fff::/128"),
			},
			expected: `{"ipv6_addr":"3fff::/128"}`,
		},
		"omit_ipv6_value_ipv4": {
			data: apstra.SecurityZoneLoopback{
				IPv4Addr: mustParsePrefixPtr("192.0.2.1/32"),
			},
			expected: `{"ipv4_addr":"192.0.2.1/32"}`,
		},
		"null_ipv4_value_ipv6": {
			data: apstra.SecurityZoneLoopback{
				IPv4Addr: &netip.Prefix{},
				IPv6Addr: mustParsePrefixPtr("3fff::/128"),
			},
			expected: `{"ipv4_addr":null,"ipv6_addr":"3fff::/128"}`,
		},
		"value_ipv4_null_ipv6": {
			data: apstra.SecurityZoneLoopback{
				IPv4Addr: mustParsePrefixPtr("192.0.2.1/32"),
				IPv6Addr: &netip.Prefix{},
			},
			expected: `{"ipv4_addr":"192.0.2.1/32","ipv6_addr":null}`,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			actual, err := json.Marshal(tCase.data)
			require.NoError(t, err)
			log.Println(string(actual))
			require.JSONEq(t, tCase.expected, string(actual))
		})
	}
}
