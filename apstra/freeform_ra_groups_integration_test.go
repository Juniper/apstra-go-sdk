//go:build integration
// +build integration

package apstra

import (
	"context"
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDRaGroups(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *FreeformRaGroupData) {
		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Label, b.Label)
		if a.ParentId != nil && b.ParentId != nil {
			require.Equal(t, *a.ParentId, *b.ParentId)
		} else {
			require.Nil(t, a.ParentId)
			require.Nil(t, b.ParentId)
		}
		compareSlicesAsSets(t, a.Tags, b.Tags, "Tags comparison mismatch")
		require.True(t, jsonEqual(t, a.Data, b.Data), "Data mismatch")
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
		if len(cfg.Data) == 0 {
			cfg.Data = json.RawMessage{'{', '}'}
		}
		compare(t, &cfg, raGroup.Data)
		data, err := json.Marshal(struct {
			Foo string `json:"foo"`
			Bar int    `json:"bar"`
		}{
			Foo: randString(6, "hex"),
			Bar: rand.Intn(40),
		})
		require.NoError(t, err)
		cfg = FreeformRaGroupData{
			Label: randString(6, "hex"),
			Tags:  []ObjectId{"tagA", "tagB"},
			Data:  data,
		}

		err = ffc.UpdateRaGroup(ctx, id, &cfg)
		require.NoError(t, err)

		raGroup, err = ffc.GetRaGroup(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, raGroup.Id)
		if len(cfg.Data) == 0 {
			cfg.Data = json.RawMessage{'{', '}'}
		}
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

		err = ffc.DeleteRaGroup(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

	}
}
