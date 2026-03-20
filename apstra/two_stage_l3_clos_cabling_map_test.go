// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra_test

import (
	"encoding/json"
	"net/netip"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/stretchr/testify/require"
)

func TestCablingMapLinkEndpointInterface_MarshalJSON(t *testing.T) {
	type TestCase struct {
		data     apstra.CablingMapLinkEndpointInterface
		expected string
	}

	testCases := map[string]TestCase{
		"everything_null": {
			data: apstra.CablingMapLinkEndpointInterface{
				ID:       "abc",
				IfName:   pointer.To(""),
				IPv4Addr: new(netip.Prefix),
				IPv6Addr: new(netip.Prefix),
			},
			expected: `{"id":"abc","if_name":null,"ipv4_addr":null,"ipv6_addr":null}`,
		},
		"everything_populated": {
			data: apstra.CablingMapLinkEndpointInterface{
				ID:       "abc",
				IfName:   pointer.To("def"),
				IPv4Addr: pointer.To(netip.MustParsePrefix("192.0.2.55/24")),
				IPv6Addr: pointer.To(netip.MustParsePrefix("3fff::1:2:3/64")),
			},
			expected: `{"id":"abc","if_name":"def","ipv4_addr":"192.0.2.55/24","ipv6_addr":"3fff::1:2:3/64"}`,
		},
		"everyting_omitted": {
			data:     apstra.CablingMapLinkEndpointInterface{ID: "abc"},
			expected: `{"id":"abc"}`,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			result, err := json.Marshal(tCase.data)
			require.NoError(t, err)
			require.JSONEq(t, tCase.expected, string(result))
		})
	}
}
