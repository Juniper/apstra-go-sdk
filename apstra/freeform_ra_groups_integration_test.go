package apstra

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCRUDRaGroups(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *FreeformRaGroupData) {
		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Label, b.Label)
		require.Equal(t, &a.ParentId, &b.ParentId)
	}

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)

		cfg := FreeformRaGroupData{
			Label: randString(6, "hex"),
		}

		id, err := ffc.CreateRaGroup(ctx, &cfg)
		require.NoError(t, err)

		raGroup, err := ffc.GetRaGroup(ctx, id)
		require.NoError(t, err)
		compare(t, &cfg, raGroup.Data)

		cfg = FreeformRaGroupData{
			ParentId: nil,
			Label:    randString(6, "hex"),
			Tags:     []ObjectId{"tagA", "tagB"},
			Data:     nil,
		}

		cfg.Label = randString(6, "hex")

		err = ffc.UpdateRaGroup(ctx, id, &cfg)
		require.NoError(t, err)

		raGroup, err = ffc.GetRaGroup(ctx, id)
		require.NoError(t, err)
		compare(t, &cfg, raGroup.Data)

		raGroups, err := ffc.GetAllRaGroups(ctx)
		require.NoError(t, err)
		ids := make([]ObjectId, len(raGroups))
		for i, template := range raGroups {
			ids[i] = template.Id
		}
		require.Contains(t, ids, id)

		err = ffc.DeleteRaGroup(ctx, id)
		require.NoError(t, err)

		_, err = ffc.GetRaGroup(ctx, id)
		require.Error(t, err)
		var ace ClientErr
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}
