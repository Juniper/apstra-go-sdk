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

func InterfaceMap(t testing.TB, req, resp design.InterfaceMap, msg ...string) {
	msg = addMsg(msg, "Comparing Interface Map")

	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, req.DeviceProfileID, resp.DeviceProfileID, msg)
	require.Equal(t, len(req.Interfaces), len(resp.Interfaces), msg)
	for i := range len(req.Interfaces) {
		InterfaceMapInterface(t, req.Interfaces[i], resp.Interfaces[i], addMsg(msg, "Comparing Interface %d", i)...)
	}

	if req.ID() != nil && resp.ID() != nil {
		require.Equal(t, *req.ID(), *resp.ID(), msg)
	}
}

func InterfaceMapInterface(t testing.TB, req, resp design.InterfaceMapInterface, msg ...string) {
	msg = addMsg(msg, "Comparing Interface Map Interface")

	require.Equal(t, req.Name, resp.Name, msg)
	require.Equal(t, req.Roles, resp.Roles, msg)
	require.Equal(t, req.Position, resp.Position, msg)
	require.Equal(t, req.State, resp.State, msg)
	require.Equal(t, req.Speed, resp.Speed, msg)
	require.Equal(t, req.Setting, resp.Setting, msg)
	InterfaceMapInterfaceMapping(t, req.Mapping, resp.Mapping, msg...)
}

func InterfaceMapInterfaceMapping(t testing.TB, req, resp design.InterfaceMapInterfaceMapping, msg ...string) {
	msg = addMsg(msg, "Comparing Mapping")

	require.Equal(t, req.DeviceProfilePortID, resp.DeviceProfilePortID, msg)
	require.Equal(t, req.DeviceProfileTransformID, resp.DeviceProfileTransformID, msg)
	require.Equal(t, req.DeviceProfileInterfaceID, resp.DeviceProfileInterfaceID, msg)
	if req.LogicalDevicePanel == nil {
		require.Nil(t, resp.LogicalDevicePanel, msg)
	} else {
		require.NotNil(t, resp.LogicalDevicePanel, msg)
		require.Equal(t, *req.LogicalDevicePanel, *resp.LogicalDevicePanel, msg)
	}
	if req.LogicalDevicePort == nil {
		require.Nil(t, resp.LogicalDevicePort, msg)
	} else {
		require.NotNil(t, resp.LogicalDevicePort, msg)
		require.Equal(t, *req.LogicalDevicePort, *resp.LogicalDevicePort, msg)
	}
}
