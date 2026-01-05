// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedesign

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func RackType(t testing.TB, req, resp design.RackType, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Rack Type")

	require.Equal(t, req.Description, resp.Description, msg)
	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, req.FabricConnectivityDesign, resp.FabricConnectivityDesign, msg)
	// req.Status -- ignoring this attribute because it seems like it's intended for use within a blueprint:
	//   As a result of flexible fabric expansion, a rack type of a blueprint rack may become inconsistent.
	//   Such rack type will be marked with inconsistent status in the blueprint. This field is accepted for
	//   global rack type as a guardrail against UI exporting inconsistent rack type from blueprint to the
	//   global catalog. The field itself will not be stored in the rack type, as it is assumed that every
	//   rack type in global catalog is consistent
	require.Equal(t, len(req.LeafSwitches), len(resp.LeafSwitches), msg)
	for i := range len(req.LeafSwitches) {
		RackTypeLeafSwitch(t, req.LeafSwitches[i], resp.LeafSwitches[i], testmessage.Add(msg, "Comparing Leaf Switch %d", i)...)
	}
	require.Equal(t, len(req.AccessSwitches), len(resp.AccessSwitches), msg)
	for i := range len(req.AccessSwitches) {
		RackTypeAccessSwitch(t, req.AccessSwitches[i], resp.AccessSwitches[i], testmessage.Add(msg, "Comparing Access Switch %d", i)...)
	}
	require.Equal(t, len(req.GenericSystems), len(resp.GenericSystems), msg)
	for i := range len(req.GenericSystems) {
		RackTypeGenericSystem(t, req.GenericSystems[i], resp.GenericSystems[i], testmessage.Add(msg, "Comparing Generic System %d", i)...)
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

func RackTypeLeafSwitch(t testing.TB, req, resp design.RackTypeLeafSwitch, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Leaf Switch")

	require.Equal(t, req.Label, resp.Label, msg)

	require.NotNil(t, resp.LinkPerSpineCount, msg)
	if req.LinkPerSpineCount == nil {
		require.Zero(t, *resp.LinkPerSpineCount, msg)
	} else {
		require.Equal(t, *req.LinkPerSpineCount, *resp.LinkPerSpineCount, msg)
	}
	if req.LinkPerSpineSpeed == nil {
		require.Nil(t, resp.LinkPerSpineSpeed, msg)
	} else {
		require.NotNil(t, resp.LinkPerSpineSpeed, msg)
		require.Equal(t, *req.LinkPerSpineSpeed, *resp.LinkPerSpineSpeed, msg)
	}
	LogicalDevice(t, req.LogicalDevice, resp.LogicalDevice, msg...)
	require.Equal(t, req.RedundancyProtocol, resp.RedundancyProtocol, msg)
	require.ElementsMatch(t, req.Tags, resp.Tags, msg)
	require.Equal(t, req.MLAGInfo, resp.MLAGInfo, msg)
}

func RackTypeAccessSwitch(t testing.TB, req, resp design.RackTypeAccessSwitch, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Access Switch")

	require.Equal(t, req.Count, resp.Count, msg)
	if req.ESILAGInfo == nil {
		require.Nil(t, resp.ESILAGInfo, msg)
	} else {
		require.NotNil(t, resp.ESILAGInfo, msg)
		require.Equal(t, *req.ESILAGInfo, *resp.ESILAGInfo, msg)
	}
	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, len(req.Links), len(resp.Links), msg)
	for i := range len(req.Links) {
		RackTypeLink(t, req.Links[i], resp.Links[i], testmessage.Add(msg, "Comparing Link %d", i)...)
	}
	LogicalDevice(t, req.LogicalDevice, resp.LogicalDevice, msg...)
	require.ElementsMatch(t, req.Tags, resp.Tags, msg)
}

func RackTypeGenericSystem(t testing.TB, req, resp design.RackTypeGenericSystem, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Generic System")

	if req.ASNDomain == nil {
		require.NotNil(t, resp.ASNDomain, msg)
		require.Equal(t, enum.FeatureSwitchDisabled, *resp.ASNDomain, msg)
	} else {
		require.NotNil(t, resp.ASNDomain, msg)
		require.Equal(t, *req.ASNDomain, *resp.ASNDomain, msg)
	}
	if req.Count == 0 {
		require.Equal(t, 1, req.Count, msg)
	} else {
		require.Equal(t, req.Count, resp.Count, msg)
	}
	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, len(req.Links), len(resp.Links), msg)
	for i := range len(req.Links) {
		RackTypeLink(t, req.Links[i], resp.Links[i], testmessage.Add(msg, "Comparing Link %d", i)...)
	}
	LogicalDevice(t, req.LogicalDevice, resp.LogicalDevice, msg...)
	if req.Loopback == nil {
		require.NotNil(t, resp.Loopback, msg)
		require.Equal(t, enum.FeatureSwitchDisabled, *resp.Loopback, msg)
	} else {
		require.NotNil(t, resp.Loopback, msg)
		require.Equal(t, *req.Loopback, *resp.Loopback, msg)
	}
	require.Equal(t, req.ManagementLevel, resp.ManagementLevel, msg)
	require.Equal(t, req.PortChannelIDMax, resp.PortChannelIDMax, msg)
	require.Equal(t, req.PortChannelIDMin, resp.PortChannelIDMin, msg)
	require.ElementsMatch(t, req.Tags, resp.Tags, msg)
}

func RackTypeLink(t testing.TB, req, resp design.RackTypeLink, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Link")

	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, req.TargetSwitchLabel, resp.TargetSwitchLabel, msg)
	require.Equal(t, req.LinkPerSwitchCount, resp.LinkPerSwitchCount, msg)
	require.Equal(t, req.Speed, resp.Speed, msg)
	require.Equal(t, req.AttachmentType, resp.AttachmentType, msg)
	require.Equal(t, req.LAGMode, resp.LAGMode, msg)
	require.Equal(t, req.SwitchPeer, resp.SwitchPeer, msg)
	if req.RailIndex == nil {
		require.Nil(t, resp.RailIndex, msg)
	} else {
		require.NotNil(t, resp.RailIndex, msg)
		require.Equal(t, *req.RailIndex, *resp.RailIndex, msg)
	}
	require.ElementsMatch(t, req.Tags, resp.Tags, msg)
}
