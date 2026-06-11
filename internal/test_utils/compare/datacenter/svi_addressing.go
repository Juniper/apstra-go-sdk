// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedatacenter

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func SVIAddressing(t testing.TB, req, resp datacenter.SVIAddressing, msg ...string) {
	msg = testmessage.Add(msg, "Comparing SVI Addressing")

	require.Equal(t, req.SystemID, resp.SystemID)
	require.Equal(t, req.IPv4Mode, resp.IPv4Mode)
	require.Equal(t, req.IPv6Mode, resp.IPv6Mode)
	if req.IPv4Addr != nil {
		require.NotNil(t, resp.IPv4Addr)
		require.Equal(t, req.IPv4Addr.String(), resp.IPv4Addr.String())
	}
	if req.IPv6Addr != nil {
		require.NotNil(t, resp.IPv6Addr)
		require.Equal(t, req.IPv6Addr.String(), resp.IPv6Addr.String())
	}
}
