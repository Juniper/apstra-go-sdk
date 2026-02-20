// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparefreeform

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func AggregateLinkEndpoint(t testing.TB, req, resp apstra.FreeformAggregateLinkEndpoint, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Aggregate Link Endpoint")

	if req.ID() != nil {
		require.NotNil(t, resp.ID(), msg)
		require.Equal(t, *req.ID(), *resp.ID(), msg)
	}

	require.Equal(t, req.SystemID, resp.SystemID, msg)
	require.Equal(t, req.IfName, resp.IfName, msg)

	if req.IPv4Addr != nil {
		require.NotNil(t, resp.IPv4Addr, msg)
		require.Equal(t, *req.IPv4Addr, *resp.IPv4Addr, msg)
	}

	if req.IPv6Addr != nil {
		require.NotNil(t, resp.IPv6Addr, msg)
		require.Equal(t, *req.IPv6Addr, *resp.IPv6Addr, msg)
	}

	require.Equal(t, req.PortChannelID, resp.PortChannelID, msg)
	require.ElementsMatch(t, req.Tags, resp.Tags, msg)
	require.Equal(t, req.LAGMode, resp.LAGMode, msg)
}
