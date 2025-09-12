// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"sort"
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func TestGetRedundancyGroupInfo(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *RedundancyGroupInfo) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Id, b.Id)
		require.Equal(t, a.Type, b.Type)
		require.Equal(t, a.SystemType, b.SystemType)
		require.Equal(t, a.SystemRole, b.SystemRole)

		aSystemIds := a.SystemIds[:]
		bSystemIds := b.SystemIds[:]
		sort.Slice(aSystemIds, func(i, j int) bool { return aSystemIds[i] < aSystemIds[j] })
		sort.Slice(bSystemIds, func(i, j int) bool { return bSystemIds[i] < bSystemIds[j] })
		require.Equal(t, aSystemIds, bSystemIds)
	}

	for _, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			bp := testBlueprintE(ctx, t, client.client)
			expectedAccessRgCount := 2
			expectedLeafRgCount := 2
			expectedTotalRgCount := expectedAccessRgCount + expectedLeafRgCount

			rgInfoMap, err := bp.GetAllRedundancyGroupInfo(ctx)
			require.NoError(t, err)
			require.Equal(t, expectedTotalRgCount, len(rgInfoMap)) // blueprint E has 4 RGs

			var accessRgCount, leafRgCount int
			systemIds := make(map[ObjectId]struct{}, len(rgInfoMap)*2)
			for k, v := range rgInfoMap {
				require.Equal(t, k, v.Id)
				require.Equal(t, enum.RedundancyGroupTypeEsi.String(), v.Type.String())
				require.Equal(t, enum.SystemTypeSwitch.String(), v.SystemType.String())

				switch v.SystemRole.String() {
				case enum.NodeRoleAccess.String():
					accessRgCount++
				case enum.NodeRoleLeaf.String():
					leafRgCount++
				}

				for _, id := range v.SystemIds {
					systemIds[id] = struct{}{}
				}
			}
			require.Equal(t, expectedTotalRgCount*2, len(systemIds))
			require.Equal(t, expectedAccessRgCount, accessRgCount)
			require.Equal(t, expectedLeafRgCount, leafRgCount)

			for k, v := range rgInfoMap {
				rgInfo, err := bp.GetRedundancyGroupInfo(ctx, k)
				require.NoError(t, err)
				compare(t, &v, rgInfo)

				for _, id := range v.SystemIds {
					rgInfo, err = bp.GetRedundancyGroupInfoBySystemId(ctx, id)
					require.NoError(t, err)
					compare(t, &v, rgInfo)
				}
			}
		})
	}
}
