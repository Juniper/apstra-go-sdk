//go:build integration
// +build integration

package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func compareFfAllocGroupData(t *testing.T, a, b *FreeformAllocGroupData) {
	t.Helper()

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.Equal(t, a.Name, b.Name)
	require.Equal(t, a.Type, b.Type)
	if a.PoolIds != nil && b.PoolIds != nil {
		compareSlicesAsSets(t, a.PoolIds, b.PoolIds, "pool IDs don't match")
	} else {
		require.Nil(t, a.PoolIds)
		require.Nil(t, b.PoolIds)
	}
}

func TestCRUDAsnAllocGroup(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)

		cfg := FreeformAllocGroupData{
			Name: randString(6, "hex"),
			Type: ResourcePoolTypeAsn,
			PoolIds: []ObjectId{
				testAsnPool(ctx, t, ffc.client),
			},
		}

		id, err := ffc.CreateAllocGroup(ctx, &cfg)
		require.NoError(t, err)

		allocGroup, err := ffc.GetAllocGroup(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, allocGroup.Id)
		compareFfAllocGroupData(t, &cfg, allocGroup.Data)

		cfg.PoolIds = []ObjectId{
			testAsnPool(ctx, t, ffc.client),
			testAsnPool(ctx, t, ffc.client),
		}

		err = ffc.UpdateAllocGroup(ctx, id, &cfg)
		require.NoError(t, err)

		allocGroup, err = ffc.GetAllocGroup(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, allocGroup.Id)
		compareFfAllocGroupData(t, &cfg, allocGroup.Data)

		err = ffc.DeleteAllocGroup(ctx, id)
		require.NoError(t, err)

		var ace ClientErr

		_, err = ffc.GetAllocGroup(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteSystem(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}

func TestCRUDIntAllocGroup(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)

		cfg := FreeformAllocGroupData{
			Name: randString(6, "hex"),
			Type: ResourcePoolTypeInt,
			PoolIds: []ObjectId{
				testIntPool(ctx, t, ffc.client),
			},
		}

		id, err := ffc.CreateAllocGroup(ctx, &cfg)
		require.NoError(t, err)

		allocGroup, err := ffc.GetAllocGroup(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, allocGroup.Id)
		compareFfAllocGroupData(t, &cfg, allocGroup.Data)

		cfg.PoolIds = []ObjectId{
			testIntPool(ctx, t, ffc.client),
			testIntPool(ctx, t, ffc.client),
		}

		err = ffc.UpdateAllocGroup(ctx, id, &cfg)
		require.NoError(t, err)

		allocGroup, err = ffc.GetAllocGroup(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, allocGroup.Id)
		compareFfAllocGroupData(t, &cfg, allocGroup.Data)

		err = ffc.DeleteAllocGroup(ctx, id)
		require.NoError(t, err)

		var ace ClientErr

		_, err = ffc.GetAllocGroup(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteSystem(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}

func TestCRUDVniAllocGroup(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)

		cfg := FreeformAllocGroupData{
			Name: randString(6, "hex"),
			Type: ResourcePoolTypeVni,
			PoolIds: []ObjectId{
				testVniPool(ctx, t, ffc.client),
			},
		}

		id, err := ffc.CreateAllocGroup(ctx, &cfg)
		require.NoError(t, err)

		allocGroup, err := ffc.GetAllocGroup(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, allocGroup.Id)
		compareFfAllocGroupData(t, &cfg, allocGroup.Data)

		cfg.PoolIds = []ObjectId{
			testVniPool(ctx, t, ffc.client),
			testVniPool(ctx, t, ffc.client),
		}

		err = ffc.UpdateAllocGroup(ctx, id, &cfg)
		require.NoError(t, err)

		allocGroup, err = ffc.GetAllocGroup(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, allocGroup.Id)
		compareFfAllocGroupData(t, &cfg, allocGroup.Data)

		err = ffc.DeleteAllocGroup(ctx, id)
		require.NoError(t, err)

		var ace ClientErr

		_, err = ffc.GetAllocGroup(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteSystem(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}

func TestCRUDIpv4AllocGroup(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)

		cfg := FreeformAllocGroupData{
			Name: randString(6, "hex"),
			Type: ResourcePoolTypeIpv4,
			PoolIds: []ObjectId{
				testIpv4Pool(ctx, t, ffc.client),
			},
		}

		id, err := ffc.CreateAllocGroup(ctx, &cfg)
		require.NoError(t, err)

		allocGroup, err := ffc.GetAllocGroup(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, allocGroup.Id)
		compareFfAllocGroupData(t, &cfg, allocGroup.Data)

		cfg.PoolIds = []ObjectId{
			testIpv4Pool(ctx, t, ffc.client),
			testIpv4Pool(ctx, t, ffc.client),
		}

		err = ffc.UpdateAllocGroup(ctx, id, &cfg)
		require.NoError(t, err)

		allocGroup, err = ffc.GetAllocGroup(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, allocGroup.Id)
		compareFfAllocGroupData(t, &cfg, allocGroup.Data)

		err = ffc.DeleteAllocGroup(ctx, id)
		require.NoError(t, err)

		var ace ClientErr

		_, err = ffc.GetAllocGroup(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteSystem(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}

func TestCRUDIpv6AllocGroup(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)

		cfg := FreeformAllocGroupData{
			Name: randString(6, "hex"),
			Type: ResourcePoolTypeIpv6,
			PoolIds: []ObjectId{
				testIpv6Pool(ctx, t, ffc.client),
			},
		}

		id, err := ffc.CreateAllocGroup(ctx, &cfg)
		require.NoError(t, err)

		allocGroup, err := ffc.GetAllocGroup(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, allocGroup.Id)
		compareFfAllocGroupData(t, &cfg, allocGroup.Data)

		cfg.PoolIds = []ObjectId{
			testIpv6Pool(ctx, t, ffc.client),
			testIpv6Pool(ctx, t, ffc.client),
		}

		err = ffc.UpdateAllocGroup(ctx, id, &cfg)
		require.NoError(t, err)

		allocGroup, err = ffc.GetAllocGroup(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, allocGroup.Id)
		compareFfAllocGroupData(t, &cfg, allocGroup.Data)

		err = ffc.DeleteAllocGroup(ctx, id)
		require.NoError(t, err)

		var ace ClientErr

		_, err = ffc.GetAllocGroup(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteSystem(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}
