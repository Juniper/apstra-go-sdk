// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package compareapstra

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func InterconnectSecurityZone(t testing.TB, req, resp apstra.InterconnectSecurityZone, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Interconnect Security Zone")

	require.Equal(t, req.L3Enabled, resp.L3Enabled, msg)
	if req.RouteTarget == nil {
		require.Nil(t, resp.RouteTarget, msg)
	} else {
		require.NotNil(t, resp.RouteTarget, msg)
		require.Equal(t, *req.RouteTarget, *resp.RouteTarget, msg)
	}
	if req.RoutingPolicyID == nil {
		require.Nil(t, resp.RoutingPolicyID, msg)
	} else {
		require.NotNil(t, resp.RoutingPolicyID, msg)
		require.Equal(t, *req.RoutingPolicyID, *resp.RoutingPolicyID, msg)
	}
}
