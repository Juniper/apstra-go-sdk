// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package compare

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/stretchr/testify/require"
)

func RackType(t testing.TB, req apstra.RackTypeRequest, data apstra.RackTypeData) {
	t.Helper()

	require.Equal(t, req.Description, data.Description)
	require.Equal(t, req.DisplayName, data.DisplayName)
	require.Equal(t, req.FabricConnectivityDesign, data.FabricConnectivityDesign)
	require.Equal(t, len(req.LeafSwitches), len(data.LeafSwitches))
	for i := range data.LeafSwitches {
		rackElementLeafSwitch(t, req.LeafSwitches[i], data.LeafSwitches[i])
	}
	require.Equal(t, len(req.AccessSwitches), len(data.AccessSwitches))
	for i := range data.AccessSwitches {
		rackElementAccessSwitch(t, req.AccessSwitches[i], data.AccessSwitches[i])
	}
	require.Equal(t, len(req.GenericSystems), len(data.GenericSystems))
	for i := range data.GenericSystems {
		rackElementGenericSystem(t, req.GenericSystems[i], data.GenericSystems[i])
	}
}

func rackElementLeafSwitch(t testing.TB, req apstra.RackElementLeafSwitchRequest, data apstra.RackElementLeafSwitch) {
	t.Helper()

	require.Equal(t, req.Label, data.Label)
	mlagInfo(t, req.MlagInfo, data.MlagInfo)
	require.Equal(t, req.LinkPerSpineCount, data.LinkPerSpineCount)
	require.Equal(t, req.LinkPerSpineSpeed, data.LinkPerSpineSpeed)
	require.Equal(t, req.RedundancyProtocol, data.RedundancyProtocol)
	require.Equal(t, len(req.Tags), len(data.Tags))
}

func rackElementAccessSwitch(t testing.TB, req apstra.RackElementAccessSwitchRequest, data apstra.RackElementAccessSwitch) {
	t.Helper()

	require.Equal(t, req.InstanceCount, data.InstanceCount)
	require.Equal(t, req.RedundancyProtocol, data.RedundancyProtocol)
	require.Equal(t, len(req.Links), len(data.Links))
	for i := range data.Links {
		rackLink(t, req.Links[i], data.Links[i])
	}
	require.Equal(t, req.Label, data.Label)
	// cannot compare logical device
	require.Equal(t, len(req.Tags), len(data.Tags))
	esiLagInfo(t, req.EsiLagInfo, data.EsiLagInfo)
}

func rackElementGenericSystem(t testing.TB, req apstra.RackElementGenericSystemRequest, data apstra.RackElementGenericSystem) {
	t.Helper()

	require.Equal(t, req.Count, data.Count)
	require.Equal(t, req.AsnDomain, data.AsnDomain)
	require.Equal(t, req.ManagementLevel, data.ManagementLevel)
	require.Equal(t, req.PortChannelIdMin, data.PortChannelIdMin)
	require.Equal(t, req.PortChannelIdMax, data.PortChannelIdMax)
	require.Equal(t, req.Loopback, data.Loopback)
	require.Equal(t, len(req.Tags), len(data.Tags))
	require.Equal(t, req.Label, data.Label)
	require.Equal(t, len(req.Links), len(data.Links))
	for i := range data.Links {
		rackLink(t, req.Links[i], data.Links[i])
	}
	// cannot compare logical device
}

func rackLink(t testing.TB, req apstra.RackLinkRequest, data apstra.RackLink) {
	t.Helper()

	require.Equal(t, req.Label, data.Label)
	require.Equal(t, len(req.Tags), len(data.Tags))
	require.Equal(t, req.LinkPerSwitchCount, data.LinkPerSwitchCount)
	require.Equal(t, req.LinkSpeed, data.LinkSpeed)
	require.Equal(t, req.TargetSwitchLabel, data.TargetSwitchLabel)
	require.Equal(t, req.AttachmentType, data.AttachmentType)
	require.Equal(t, req.LagMode, data.LagMode)
	require.Equal(t, req.SwitchPeer, data.SwitchPeer)
}

func esiLagInfo(t testing.TB, a, b *apstra.EsiLagInfo) {
	t.Helper()

	if a == nil {
		require.Nil(t, b)
		return
	}

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.Equal(t, a.AccessAccessLinkCount, b.AccessAccessLinkCount)
	require.Equal(t, a.AccessAccessLinkSpeed, b.AccessAccessLinkSpeed)
}

func mlagInfo(t testing.TB, a, b *apstra.LeafMlagInfo) {
	t.Helper()

	if a == nil {
		if b != nil {
			require.Zero(t, b.LeafLeafL3LinkCount)
			require.Zero(t, b.LeafLeafL3LinkPortChannelId)
			require.Zero(t, b.LeafLeafL3LinkSpeed)
			require.Zero(t, b.LeafLeafLinkCount)
			require.Zero(t, b.LeafLeafLinkPortChannelId)
			require.Zero(t, b.LeafLeafLinkSpeed)
		}
		return
	}

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.Equal(t, a.LeafLeafL3LinkCount, b.LeafLeafL3LinkCount)
	require.Equal(t, a.LeafLeafL3LinkPortChannelId, b.LeafLeafL3LinkPortChannelId)
	require.Equal(t, a.LeafLeafL3LinkCount, b.LeafLeafL3LinkCount)
	require.Equal(t, a.LeafLeafLinkPortChannelId, b.LeafLeafLinkPortChannelId)
	require.Equal(t, a.LeafLeafLinkSpeed, b.LeafLeafLinkSpeed)
	require.Equal(t, a.LeafLeafLinkSpeed, b.LeafLeafLinkSpeed)
}
