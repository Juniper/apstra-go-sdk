// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestEmptyAsnPool(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	asnRangeCount := rand.Intn(5) + 2 // random number of ASN ranges to add to new pool
	asnBeginEnds, err := testutils.GetRandInts(1, 100000000, asnRangeCount*2)
	require.NoError(t, err)

	sort.Ints(asnBeginEnds) // sort so that the ASN ranges will be ([0]...[1], [2]...[3], etc.)
	asnRanges := make([]apstra.IntfIntRange, asnRangeCount)
	for i := 0; i < asnRangeCount; i++ {
		asnRanges[i] = apstra.IntRangeRequest{
			First: uint32(asnBeginEnds[2*i]),
			Last:  uint32(asnBeginEnds[(2*i)+1]),
		}
	}

	poolName := "test-" + testutils.RandString(10, "hex")

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			newPoolId, err := client.Client.CreateAsnPool(ctx, &apstra.AsnPoolRequest{DisplayName: poolName})
			require.NoError(t, err)
			log.Printf("created ASN pool name %s id %s", poolName, newPoolId)

			newPool, err := client.Client.GetAsnPool(ctx, newPoolId)
			require.NoError(t, err)
			require.Equal(t, newPoolId, newPool.Id)
			require.Equal(t, newPool.DisplayName, newPool.DisplayName)
			require.Zero(t, len(newPool.Ranges))

			for i := range asnRanges {
				newName := fmt.Sprintf("%s-%d", poolName, i)
				err = client.Client.UpdateAsnPool(ctx, newPoolId, &apstra.AsnPoolRequest{
					DisplayName: newName,
					Ranges:      asnRanges[:i+1],
				})
				require.NoError(t, err)

				newPool, err = client.Client.GetAsnPool(ctx, newPoolId)
				require.NoError(t, err)
				require.Equal(t, newPoolId, newPool.Id)
				require.Equal(t, newName, newPool.DisplayName)
				require.Equal(t, i+1, len(newPool.Ranges))
			}

			err = client.Client.DeleteAsnPool(ctx, newPoolId)
			require.NoError(t, err)
		})
	}
}

func TestGetAsnPoolByName(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	poolName := testutils.RandString(10, "hex")

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			_, err := client.Client.GetAsnPoolByName(ctx, poolName)
			require.Error(t, err)
			var ace apstra.ClientErr
			require.ErrorAs(t, err, &ace)
			require.Equal(t, apstra.ErrNotfound, ace.Type())

			id, err := client.Client.CreateAsnPool(ctx, &apstra.AsnPoolRequest{DisplayName: poolName})
			require.NoError(t, err)

			p, err := client.Client.GetAsnPoolByName(ctx, poolName)
			require.NoError(t, err)
			require.Equal(t, id, p.Id)

			err = client.Client.DeleteAsnPool(ctx, id)
			require.NoError(t, err)
		})
	}
}

func TestListAsnPoolIds(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			poolIds, err := client.Client.ListAsnPoolIds(ctx)
			require.NoError(t, err)
			require.NotZero(t, len(poolIds), "no ASN pools on this system?")
		})
	}
}

func TestEmptyVniPool(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	vniRangeCount := rand.Intn(5) + 2 // random number of VNI ranges to add to new pool
	vniBeginEnds, err := testutils.GetRandInts(apstra.VniMin, apstra.VniMax, vniRangeCount*2)
	require.NoError(t, err)

	sort.Ints(vniBeginEnds) // sort so that the VNI ranges will be ([0]...[1], [2]...[3], etc.)
	vniRanges := make([]apstra.IntfIntRange, vniRangeCount)
	for i := 0; i < vniRangeCount; i++ {
		vniRanges[i] = apstra.IntRangeRequest{
			First: uint32(vniBeginEnds[2*i]),
			Last:  uint32(vniBeginEnds[(2*i)+1]),
		}
	}

	poolName := "test-" + testutils.RandString(10, "hex")

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			newPoolId, err := client.Client.CreateVniPool(ctx, &apstra.VniPoolRequest{DisplayName: poolName})
			require.NoError(t, err)
			log.Printf("created VNI pool name %s id %s", poolName, newPoolId)

			newPool, err := client.Client.GetVniPool(ctx, newPoolId)
			require.NoError(t, err)
			require.Equal(t, newPoolId, newPool.Id)
			require.Equal(t, newPool.DisplayName, newPool.DisplayName)
			require.Zero(t, len(newPool.Ranges))

			for i := range vniRanges {
				newName := fmt.Sprintf("%s-%d", poolName, i)
				err = client.Client.UpdateVniPool(ctx, newPoolId, &apstra.VniPoolRequest{
					DisplayName: newName,
					Ranges:      vniRanges[:i+1],
				})
				require.NoError(t, err)

				newPool, err = client.Client.GetVniPool(ctx, newPoolId)
				require.NoError(t, err)
				require.Equal(t, newPoolId, newPool.Id)
				require.Equal(t, newName, newPool.DisplayName)
				require.Equal(t, i+1, len(newPool.Ranges))
			}

			err = client.Client.DeleteVniPool(ctx, newPoolId)
			require.NoError(t, err)
		})
	}
}

func TestGetVniPoolByName(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	poolName := testutils.RandString(10, "hex")

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			_, err := client.Client.GetVniPoolByName(ctx, poolName)
			require.Error(t, err)
			var ace apstra.ClientErr
			require.ErrorAs(t, err, &ace)
			require.Equal(t, apstra.ErrNotfound, ace.Type())

			id, err := client.Client.CreateVniPool(ctx, &apstra.VniPoolRequest{DisplayName: poolName})
			require.NoError(t, err)

			p, err := client.Client.GetVniPoolByName(ctx, poolName)
			require.NoError(t, err)
			require.Equal(t, id, p.Id)

			err = client.Client.DeleteVniPool(ctx, id)
			require.NoError(t, err)
		})
	}
}

func TestListVniPoolIds(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			poolIds, err := client.Client.ListVniPoolIds(ctx)
			require.NoError(t, err)
			require.NotZero(t, len(poolIds))
		})
	}
}
