// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func TestCRUDRaLocalPools(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compareChunks := func(t *testing.T, a, b []FFLocalIntPoolChunk) {
		t.Helper()

		require.Equal(t, len(a), len(b))

		for i := range a {
			require.Equal(t, a[i].Start, b[i].Start)
			require.Equal(t, a[i].End, b[i].End)
		}
	}

	compare := func(t *testing.T, a, b *FreeformRaLocalIntPoolData) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.ResourceType, b.ResourceType)
		require.Equal(t, a.Label, b.Label)
		require.Equal(t, a.OwnerId, b.OwnerId)
		if a.GeneratorId != nil {
			require.NotNil(t, b.GeneratorId)
		}
		if b.GeneratorId != nil {
			require.NotNil(t, a.GeneratorId)
			require.Equal(t, *a.GeneratorId, b.GeneratorId)
		}
		compareChunks(t, a.Chunks, b.Chunks)
	}

	for _, client := range clients {
		ffc, systemIds, _ := testFFBlueprintB(ctx, t, client.client, 1, 0)
		require.Equal(t, len(systemIds), 1)

		cfg := FreeformRaLocalIntPoolData{
			ResourceType: enum.FFResourceTypeVlan,
			Label:        randString(6, "hex"),
			OwnerId:      systemIds[0],
			Chunks:       []FFLocalIntPoolChunk{{Start: 10, End: 20}},
		}

		id, err := ffc.CreateRaLocalIntPool(ctx, &cfg)
		require.NoError(t, err)

		localIntPool, err := ffc.GetRaLocalIntPool(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, localIntPool.Id)
		compare(t, &cfg, localIntPool.Data)

		cfg.Label = randString(6, "hex")
		cfg.Chunks = []FFLocalIntPoolChunk{{Start: 5, End: 15}, {Start: 16, End: 25}}

		err = ffc.UpdateRaLocalIntPool(ctx, id, &cfg)
		require.NoError(t, err)

		localIntPool, err = ffc.GetRaLocalIntPool(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, localIntPool.Id)
		compare(t, &cfg, localIntPool.Data)

		localIntPools, err := ffc.GetAllRaLocalIntPools(ctx)
		require.NoError(t, err)
		ids := make([]ObjectId, len(localIntPools))
		for i, template := range localIntPools {
			ids[i] = template.Id
		}
		require.Contains(t, ids, id)

		err = ffc.DeleteRaLocalIntPool(ctx, id)
		require.NoError(t, err)

		var ace ClientErr

		_, err = ffc.GetRaLocalIntPool(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteRaLocalIntPool(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

	}
}
