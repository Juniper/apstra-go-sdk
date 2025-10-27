// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedesign

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/stretchr/testify/require"
)

func LogicalDevice(t testing.TB, req, resp design.LogicalDevice, msg ...string) {
	msg = addMsg(msg, "Comparing Logical Device")

	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, len(req.Panels), len(resp.Panels), msg)
	for i := range len(req.Panels) {
		LogicalDevicePanel(t, req.Panels[i], resp.Panels[i], addMsg(msg, "Comparing Panel %d", i)...)
	}
	if req.ID() != nil && resp.ID() != nil {
		require.Equal(t, req.ID(), resp.ID(), msg)
	}
	if req.CreatedAt() != nil && resp.CreatedAt() != nil {
		require.Equal(t, req.CreatedAt(), resp.CreatedAt(), msg)
	}
	if req.LastModifiedAt() != nil && resp.LastModifiedAt() != nil {
		require.Equal(t, req.LastModifiedAt(), resp.LastModifiedAt(), msg)
	}
}

func LogicalDevicePanel(t testing.TB, req, resp design.LogicalDevicePanel, msg ...string) {
	msg = addMsg(msg, "Comparing Logical Device Panel")

	require.Equal(t, req.PanelLayout, resp.PanelLayout, msg)
	require.Equal(t, req.PortIndexing, resp.PortIndexing, msg)
	require.Equal(t, len(req.PortGroups), len(resp.PortGroups), msg)
	for i := range len(req.PortGroups) {
		LogicalDevicePanelPortGroup(t, req.PortGroups[i], resp.PortGroups[i], addMsg(msg, "Port Group %d", i)...)
	}
}

func LogicalDevicePanelPortGroup(t testing.TB, req, resp design.LogicalDevicePanelPortGroup, msg ...string) {
	msg = addMsg(msg, "Comparing Logical Device Panel Port Group")

	require.Equal(t, req.Count, resp.Count, msg)
	require.Equal(t, req.Speed, resp.Speed, msg)
	require.ElementsMatch(t, req.Roles, resp.Roles, msg)
}
