package apstra

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRaGroupGenA(t *testing.T) {
	var x FreeformRaGroupGenerator
	x.Data = new(FreeformRaGroupGeneratorData)
	x.Id = "foo"
	x.Data.Scope = "node('link', role='internal', name='target')"
	x.Data.Label = "GroupGenTest"
	rawjson, err := json.Marshal(&x)
	require.NoError(t, err)
	log.Println(string(rawjson))
}

func TestCRUDRaGroupGenerators(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *FreeformRaGroupGeneratorData) {
		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Label, b.Label)
		if a.ParentId != nil && b.ParentId != nil {
			require.Equal(t, *a.ParentId, *b.ParentId)
		} else {
			require.Nil(t, a.ParentId)
			require.Nil(t, b.ParentId)
		}
		require.Equal(t, a.Scope, b.Scope, "Scope comparison mismatch")
	}

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)

		cfg := FreeformRaGroupGeneratorData{
			Label: randString(6, "hex"),
			Scope: "node('link', role='internal', name='target')",
		}

		id, err := ffc.CreateRaGroupGenerator(ctx, &cfg)
		require.NoError(t, err)

		raGroup, err := ffc.GetRaGroupGenerator(ctx, id)
		require.NoError(t, err)

		compare(t, &cfg, raGroup.Data)

		require.NoError(t, err)
		cfg = FreeformRaGroupGeneratorData{
			Label: randString(6, "hex"),
			Scope: "node('link', role='internal', name='target')",
		}

		err = ffc.UpdateRaGroupGenerator(ctx, id, &cfg)
		require.NoError(t, err)

		raGroup, err = ffc.GetRaGroupGenerator(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, raGroup.Id)
		compare(t, &cfg, raGroup.Data)

		raGroups, err := ffc.GetAllRaGroupGenerators(ctx)
		require.NoError(t, err)
		ids := make([]ObjectId, len(raGroups))
		for i, template := range raGroups {
			ids[i] = template.Id
		}
		require.Contains(t, ids, id)

		err = ffc.DeleteRaGroupGenerator(ctx, id)
		require.NoError(t, err)

		_, err = ffc.GetRaGroupGenerator(ctx, id)
		require.Error(t, err)
		var ace ClientErr
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteRaGroupGenerator(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

	}
}
