//go:build integration
// +build integration

package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDRaResource(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *FreeformRaResourceData) {
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

		groupCfg := FreeformRaGroupData{
			Label: randString(6, "hex"),
		}

		// first thing to do here is to create a group
		groupId, err := ffc.CreateRaGroup(ctx, &groupCfg)
		require.NoError(t, err)
		// now create our resource
		raResourceCfg := FreeformRaResourceData{
			Label:        randString(6, "hex"),
			GroupId:      groupId,
			ResourceType: FFResourceTypeInt,
		}
		resourceId, err := ffc.CreateRaResource(ctx, &raResourceCfg)
		require.NoError(t, err)
		raResource, err := ffc.GetRaResource(ctx, resourceId)
		require.NoError(t, err)

		compare(t, &raResourceCfg, raResource.Data)
		raResourceCfg = FreeformRaResourceData{
			Label:        randString(6, "hex"),
			GroupId:      groupId,
			ResourceType: FFResourceTypeIpv4,
		}

		err = ffc.UpdateRaResource(ctx, resourceId, &raResourceCfg)
		require.NoError(t, err)

		resource, err := ffc.GetRaResource(ctx, raResource.Id)
		require.NoError(t, err)
		require.Equal(t, resourceId, resource.Id)

		raResources, err := ffc.GetAllRaResources(ctx)
		require.NoError(t, err)
		ids := make([]ObjectId, len(raResources))
		for i, template := range raResources {
			ids[i] = template.Id
		}
		require.Contains(t, ids, resourceId)

		err = ffc.DeleteRaResource(ctx, resourceId)
		require.NoError(t, err)

		_, err = ffc.GetRaResource(ctx, resourceId)
		require.Error(t, err)
		var ace ClientErr
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteRaResource(ctx, resourceId)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}
