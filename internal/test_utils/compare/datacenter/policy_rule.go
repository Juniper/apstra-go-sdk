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

func PolicyRule(t testing.TB, req, resp datacenter.PolicyRule, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Policy Rule")

	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, req.Description, resp.Description, msg)
	require.Equal(t, req.Protocol, resp.Protocol, msg)
	require.Equal(t, req.Action, resp.Action, msg)
	PortRanges(t, req.SrcPort, resp.SrcPort, msg...)
	PortRanges(t, req.DstPort, resp.DstPort, msg...)
	if req.TcpStateQualifier == nil {
		require.Nil(t, resp.TcpStateQualifier)
	} else {
		require.Equal(t, *req.TcpStateQualifier, *resp.TcpStateQualifier)
	}
	if req.ID() != nil {
		require.Equal(t, *req.ID(), *resp.ID(), msg)
	}
}
