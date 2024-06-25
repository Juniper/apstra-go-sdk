package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDRaLocalGroupGenerators(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *FreeformRaLocalIntPoolGeneratorData) {
		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Label, b.Label)
		require.Equal(t, a.Scope, b.Scope, "Scope comparison mismatch")
	}

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)

		cfg := FreeformRaLocalIntPoolGeneratorData{
			ResourceType: FFResourceTypeVlan,
			Label:        randString(6, "hex"),
			Scope:        "node('link', role='internal', name='target')",
			Chunks: []FFLocalIntPoolChunk{
				{Start: 10, End: 20},
			},
		}

		id, err := ffc.CreateRaLocalIntPoolGenerator(ctx, &cfg)
		require.NoError(t, err)
		raGroup := new(FreeformRaLocalIntPoolGenerator)
		raGroup.Data = new(FreeformRaLocalIntPoolGeneratorData)
		raGroup, err = ffc.GetRaLocalIntPoolGenerator(ctx, id)
		require.NoError(t, err)

		compare(t, &cfg, raGroup.Data)

		require.NoError(t, err)
		cfg = FreeformRaLocalIntPoolGeneratorData{
			ResourceType: FFResourceTypeVlan,
			Label:        randString(6, "hex"),
			Scope:        "node('link', role='internal', name='target')",
			Chunks: []FFLocalIntPoolChunk{
				{Start: 20, End: 30},
			},
		}

		err = ffc.UpdateRaLocalIntPoolGenerator(ctx, id, &cfg)
		require.NoError(t, err)

		raGroup, err = ffc.GetRaLocalIntPoolGenerator(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, raGroup.Id)
		// compare(t, &cfg, raGroup.Data)

		raGroups, err := ffc.GetAllLocalIntPoolGenerators(ctx)
		require.NoError(t, err)
		ids := make([]ObjectId, len(raGroups))
		for i, template := range raGroups {
			ids[i] = template.Id
		}
		require.Contains(t, ids, id)

		err = ffc.DeleteRaLocalPoolGenerator(ctx, id)
		require.NoError(t, err)

		_, err = ffc.GetRaLocalIntPoolGenerator(ctx, id)
		require.Error(t, err)
		var ace ClientErr
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteRaLocalPoolGenerator(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

	}
}
