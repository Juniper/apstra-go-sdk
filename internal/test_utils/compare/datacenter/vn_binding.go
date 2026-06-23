// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedatacenter

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func VNBinding(t testing.TB, req, resp datacenter.VNBinding, msg ...string) {
	msg = testmessage.Add(msg, "Comparing VN Binding")

	require.Equal(t, len(req.AccessSwitchNodeIDs), len(resp.AccessSwitchNodeIDs)) // nil and [] are equal for our purposes
	if len(req.AccessSwitchNodeIDs) > 0 {
		compare.SlicesAsSets(t, req.AccessSwitchNodeIDs, resp.AccessSwitchNodeIDs, "VNBinding AccessSwitchNodeIDs")
	}
	require.Equal(t, req.SystemID, resp.SystemID)
	if req.VLAN != nil {
		require.NotNil(t, resp.VLAN)
		require.Equal(t, *req.VLAN, *resp.VLAN)
	}
}
