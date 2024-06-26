//go:build integration
// +build integration

package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDPropSets(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *FreeformPropertySetData) {
		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Label, b.Label)
		require.Equal(t, a.Values, b.Values)
		if a.SystemId != nil {
			require.NotNil(t, b.SystemId)
		}
		if b.SystemId != nil {
			require.NotNil(t, a.SystemId)
			require.Equal(t, *a.SystemId, *b.SystemId)
		}
	}

	for _, client := range clients {
		ffc, systemIds := testFFBlueprintB(ctx, t, client.client, 1)
		require.Equal(t, len(systemIds), 1)

		cfg := FreeformPropertySetData{
			Label:  randString(6, "hex"),
			Values: make(map[string]string),
		}
		for i := 0; i < 5; i++ {
			cfg.Values["a"+randString(6, "hex")] = randString(6, "hex")
		}

		// todo: test CreatePropertySet with non-nil SystemId

		id, err := ffc.CreatePropertySet(ctx, &cfg)
		require.NoError(t, err)

		propertySet, err := ffc.GetPropertySet(ctx, id)
		require.NoError(t, err)
		compare(t, &cfg, propertySet.Data)

		cfg.Label = randString(6, "hex")
		cfg.SystemId = &systemIds[0]
		cfg.Values = map[string]string{}
		for i := 0; i < 5; i++ {
			cfg.Values["a"+randString(6, "hex")] = randString(6, "hex")
		}

		// todo: test CreatePropertySet with nil SystemId

		err = ffc.UpdatePropertySet(ctx, id, &cfg)
		require.NoError(t, err)

		propertySet, err = ffc.GetPropertySet(ctx, id)
		require.NoError(t, err)
		compare(t, &cfg, propertySet.Data)

		propertySets, err := ffc.GetAllPropertySets(ctx)
		require.NoError(t, err)

		ids := make([]ObjectId, len(propertySets))
		for i, template := range propertySets {
			ids[i] = template.Id
		}
		require.Contains(t, ids, id)

		err = ffc.DeletePropertySet(ctx, id)
		require.NoError(t, err)

		_, err = ffc.GetPropertySet(ctx, id)
		require.Error(t, err)
		var ace ClientErr
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeletePropertySet(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}
