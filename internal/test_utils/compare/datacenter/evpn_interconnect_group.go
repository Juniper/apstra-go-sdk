// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedatacenter

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func EVPNInterconnectGroup(t testing.TB, req, resp apstra.EVPNInterconnectGroup, msg ...string) {
	msg = testmessage.Add(msg, "Comparing EVPN Interconnect Group")

	if req.Label != nil {
		require.NotNil(t, resp.Label, msg)
		require.Equal(t, *req.Label, *resp.Label, msg)
	}

	if req.RouteTarget != nil {
		require.NotNil(t, resp.RouteTarget, msg)
		require.Equal(t, *req.RouteTarget, *resp.RouteTarget, msg)
	}

	if req.ESIMAC != nil {
		require.NotNil(t, resp.ESIMAC, msg)
		require.Equal(t, req.ESIMAC.String(), resp.ESIMAC.String(), msg)
	}

	for k, reqv := range req.InterconnectSecurityZones {
		require.Contains(t, resp.InterconnectSecurityZones, k, msg)
		respv := resp.InterconnectSecurityZones[k]
		InterconnectSecurityZone(t, reqv, respv, testmessage.Add(msg, "Comparing Interconnect Security Zone %s", k)...)
	}

	for k, reqv := range req.InterconnectVirtualNetworks {
		require.Contains(t, resp.InterconnectVirtualNetworks, k, msg)
		respv := resp.InterconnectVirtualNetworks[k]
		InterconnectVirtualNetwork(t, reqv, respv, testmessage.Add(msg, "Comparing Interconnect Virtual Network %s", k)...)
	}
}
