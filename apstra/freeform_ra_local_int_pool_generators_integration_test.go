// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func TestCRUDRaLocalGroupGenerators(t *testing.T) {
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

	compare := func(t *testing.T, a, b *FreeformRaLocalIntPoolGeneratorData) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.ResourceType, b.ResourceType)
		require.Equal(t, a.Label, b.Label)
		require.Equal(t, a.Scope, b.Scope, "Scope comparison mismatch")
		compareChunks(t, a.Chunks, b.Chunks)
	}

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)

		cfg := FreeformRaLocalIntPoolGeneratorData{
			ResourceType: enum.FFResourceTypeVlan,
			Label:        randString(6, "hex"),
			Scope:        "node('link', role='internal', name='target')",
			Chunks: []FFLocalIntPoolChunk{
				{Start: 100, End: 200},
				{Start: 10, End: 20},
			},
		}

		id, err := ffc.CreateRaLocalIntPoolGenerator(ctx, &cfg)
		require.NoError(t, err)

		localIntPoolGen, err := ffc.GetRaLocalIntPoolGenerator(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, localIntPoolGen.Id)
		compare(t, &cfg, localIntPoolGen.Data)

		cfg.Scope = "node('link', role='internal', name='target')"
		cfg.Label = randString(6, "hex")
		cfg.Chunks = []FFLocalIntPoolChunk{{Start: 5, End: 150}}

		err = ffc.UpdateRaLocalIntPoolGenerator(ctx, id, &cfg)
		require.NoError(t, err)

		localIntPoolGen, err = ffc.GetRaLocalIntPoolGenerator(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, localIntPoolGen.Id)
		compare(t, &cfg, localIntPoolGen.Data)

		localIntPoolGens, err := ffc.GetAllLocalIntPoolGenerators(ctx)
		require.NoError(t, err)
		ids := make([]ObjectId, len(localIntPoolGens))
		for i, template := range localIntPoolGens {
			ids[i] = template.Id
		}
		require.Contains(t, ids, id)

		err = ffc.DeleteRaLocalPoolGenerator(ctx, id)
		require.NoError(t, err)

		var ace ClientErr

		_, err = ffc.GetRaLocalIntPoolGenerator(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteRaLocalPoolGenerator(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}
