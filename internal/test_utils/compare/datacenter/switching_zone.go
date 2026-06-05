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

func SwitchingZone(t testing.TB, req, resp datacenter.SwitchingZone, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Switching Zone")

	if req.Label != nil {
		require.NotNil(t, resp.Label, msg)
		require.Equal(t, *req.Label, *resp.Label, msg)
	}

	if req.MACVRFDescription != nil {
		require.NotNil(t, resp.MACVRFDescription, msg)
		require.Equal(t, *req.MACVRFDescription, *resp.MACVRFDescription, msg)
	}

	if req.MACVRFName != nil {
		require.NotNil(t, resp.MACVRFName, msg)
		require.Equal(t, *req.MACVRFName, *resp.MACVRFName, msg)
	}

	if req.MACVRFServiceType != nil {
		require.NotNil(t, resp.MACVRFServiceType, msg)
		require.Equal(t, *req.MACVRFServiceType, *resp.MACVRFServiceType, msg)
	}

	if req.RouteTarget != nil {
		require.NotNil(t, resp.RouteTarget, msg)
		require.Equal(t, *req.RouteTarget, *resp.RouteTarget, msg)
	}

	if req.Tags != nil {
		require.Equal(t, len(req.Tags), len(resp.Tags), msg)
		require.ElementsMatch(t, req.Tags, resp.Tags, msg)
	}

	if req.ID() != nil {
		require.NotNil(t, resp.ID(), msg)
		require.Equal(t, *req.ID(), *resp.ID(), msg)
	}
}
