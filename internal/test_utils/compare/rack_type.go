// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

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

func mlagInfo(t testing.TB, a, b *apstra.LeafMlagInfo) {
	t.Helper()

	if a == nil {
		require.Nil(t, b)
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
