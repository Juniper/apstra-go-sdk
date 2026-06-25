// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedatacenter

import (
	"fmt"
	"testing"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func Policy(t testing.TB, req, resp datacenter.Policy, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Policy")

	require.Equal(t, req.Enabled, resp.Enabled, msg)
	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, req.Description, resp.Description, msg)
	if req.SrcApplicationPoint == nil {
		require.Nil(t, resp.SrcApplicationPoint, msg)
	} else {
		require.NotNil(t, resp.SrcApplicationPoint, msg)
		require.Equal(t, *req.SrcApplicationPoint, *resp.SrcApplicationPoint, msg)
	}
	if req.DstApplicationPoint == nil {
		require.Nil(t, resp.DstApplicationPoint, msg)
	} else {
		require.NotNil(t, resp.DstApplicationPoint, msg)
		require.Equal(t, *req.DstApplicationPoint, *resp.DstApplicationPoint, msg)
	}
	require.Equal(t, len(req.Rules), len(resp.Rules), msg)
	for i := range req.Rules {
		PolicyRule(t, req.Rules[i], resp.Rules[i], testmessage.Add(msg, fmt.Sprintf("Comparing Policy Rule at index %d", i))...)
	}
	compare.SlicesAsSets(t, req.Tags, resp.Tags, "tags")
}
