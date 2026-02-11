// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedatacenter

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func SecurityZone(t testing.TB, req, resp apstra.SecurityZone, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Security Zone")

	if req.ID() != nil {
		require.NotNil(t, resp.ID(), msg)
		require.Equal(t, *req.ID(), *resp.ID(), msg)
	}

	require.Equal(t, req.Label, resp.Label, msg)

	if req.Description == nil {
		require.Nil(t, resp.Description, msg)
	} else {
		require.NotNil(t, resp.Description, msg)
		require.Equal(t, *req.Description, *resp.Description, msg)
	}

	require.Equal(t, req.Description, resp.Description, msg, msg)
	require.Equal(t, req.Type, resp.Type, msg, msg)

	require.Equal(t, req.VRFName, resp.VRFName, msg, msg)
	if req.RoutingPolicyID != "" {
		require.Equal(t, req.RoutingPolicyID, resp.RoutingPolicyID, msg)
	}

	if req.RTPolicy == nil {
		require.Nil(t, resp.RTPolicy, msg)
	} else {
		require.NotNil(t, resp.RTPolicy, msg)
		require.Equal(t, *req.RTPolicy, *resp.RTPolicy, msg)
	}

	if req.VLAN != nil {
		require.NotNil(t, resp.VLAN, msg)
		require.Equal(t, *req.VLAN, *resp.VLAN, msg)
	}

	if req.VNI != nil {
		require.NotNil(t, resp.VNI, msg)
		require.Equal(t, *req.VNI, *resp.VNI, msg)
	}

	if req.JunosEVPNIRBMode != nil {
		require.NotNil(t, resp.JunosEVPNIRBMode, msg)
		require.Equal(t, *req.JunosEVPNIRBMode, *resp.JunosEVPNIRBMode, msg)
	}

	if req.AddressingSupport != nil {
		require.NotNil(t, resp.AddressingSupport, msg)
		require.Equal(t, *req.AddressingSupport, *resp.AddressingSupport, msg)
	}

	if req.DisableIPv4 != nil {
		require.NotNil(t, resp.DisableIPv4, msg)
		require.Equal(t, *req.DisableIPv4, *resp.DisableIPv4, msg)
	}
}
