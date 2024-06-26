//go:build integration
// +build integration

package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDRaLocalPools(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *FreeformRaLocalIntPoolData) {
		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Label, b.Label)
	}

	for _, client := range clients {
		ffc, systemIds := testFFBlueprintB(ctx, t, client.client, 2)

		cfg := FreeformRaLocalIntPoolData{
			ResourceType: FFResourceTypeVlan,
			Label:        randString(6, "hex"),
			OwnerId:      systemIds[0],
			Chunks: []FFLocalIntPoolChunk{
				{
					Start: 10,
					End:   20,
				},
			},
		}

		id, err := ffc.CreateRaLocalIntPool(ctx, &cfg)
		require.NoError(t, err)

		raGroup, err := ffc.GetRaLocalIntPool(ctx, id)
		require.NoError(t, err)

		compare(t, &cfg, raGroup.Data)

		require.NoError(t, err)
		cfg = FreeformRaLocalIntPoolData{
			ResourceType: FFResourceTypeVni,
			Label:        randString(6, "hex"),
			OwnerId:      systemIds[1],
			Chunks: []FFLocalIntPoolChunk{
				{
					Start: 5,
					End:   15,
				},
				{
					Start: 16,
					End:   25,
				},
			},
		}

		err = ffc.UpdateRaLocalIntPool(ctx, id, &cfg)
		require.NoError(t, err)

		raGroup, err = ffc.GetRaLocalIntPool(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, raGroup.Id)
		compare(t, &cfg, raGroup.Data)

		raGroups, err := ffc.GetAllRaLocalIntPools(ctx)
		require.NoError(t, err)
		ids := make([]ObjectId, len(raGroups))
		for i, template := range raGroups {
			ids[i] = template.Id
		}
		require.Contains(t, ids, id)

		err = ffc.DeleteRaLocalIntPool(ctx, id)
		require.NoError(t, err)

		_, err = ffc.GetRaLocalIntPool(ctx, id)
		require.Error(t, err)
		var ace ClientErr
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteRaLocalIntPool(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

	}
}
