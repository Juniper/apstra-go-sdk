//go:build integration
// +build integration

package apstra

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDRaResource(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *FreeformRaResourceData) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.ResourceType.String(), b.ResourceType.String())
		require.Equal(t, a.Label, b.Label)
		require.Equal(t, a.GroupId, b.GroupId)

		if a.Value != nil && b.Value != nil {
			require.Equal(t, *a.Value, *b.Value)
		} else {
			require.Nil(t, a.Value)
			require.Nil(t, b.Value)
		}

		if a.AllocatedFrom != nil && b.AllocatedFrom != nil {
			require.Equal(t, *a.AllocatedFrom, *b.AllocatedFrom)
		} else {
			require.Nil(t, a.AllocatedFrom)
			require.Nil(t, b.AllocatedFrom)
		}

		if a.SubnetPrefixLen != nil && b.SubnetPrefixLen != nil {
			require.Equal(t, *a.SubnetPrefixLen, *b.SubnetPrefixLen)
		} else {
			require.Nil(t, a.SubnetPrefixLen)
			require.Nil(t, b.SubnetPrefixLen)
		}

		if a.GeneratorId != nil && b.GeneratorId != nil {
			require.Equal(t, *a.GeneratorId, *b.GeneratorId)
		} else {
			require.Nil(t, a.GeneratorId)
			require.Nil(t, b.GeneratorId)
		}
	}

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)

		// first thing to do here is to create a group
		groupCfg := FreeformRaGroupData{
			Label: randString(6, "hex"),
		}
		groupId, err := ffc.CreateRaGroup(ctx, &groupCfg)
		require.NoError(t, err)

		// now create our resource
		cfg := FreeformRaResourceData{
			Label:        randString(6, "hex"),
			GroupId:      groupId,
			ResourceType: FFResourceTypeInt,
		}

		id, err := ffc.CreateRaResource(ctx, &cfg)
		require.NoError(t, err)

		raResource, err := ffc.GetRaResource(ctx, id)
		require.Equal(t, id, raResource.Id)
		require.NoError(t, err)
		compare(t, &cfg, raResource.Data)

		_, prefix, err := net.ParseCIDR(randomIpv4().String() + "/24")
		require.NoError(t, err)

		cfg = FreeformRaResourceData{
			Label:           randString(6, "hex"),
			GroupId:         groupId,
			ResourceType:    FFResourceTypeIpv4,
			Value:           toPtr(prefix.String()),
			SubnetPrefixLen: toPtr(24),
		}

		err = ffc.UpdateRaResource(ctx, id, &cfg)
		require.NoError(t, err)

		resource, err := ffc.GetRaResource(ctx, raResource.Id)
		require.NoError(t, err)
		require.Equal(t, id, resource.Id)

		raResources, err := ffc.GetAllRaResources(ctx)
		require.NoError(t, err)
		ids := make([]ObjectId, len(raResources))
		for i, template := range raResources {
			ids[i] = template.Id
		}
		require.Contains(t, ids, id)

		err = ffc.DeleteRaResource(ctx, id)
		require.NoError(t, err)

		var ace ClientErr

		_, err = ffc.GetRaResource(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteRaResource(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}
