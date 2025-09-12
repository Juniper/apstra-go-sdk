// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestListIp4Pools(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			poolIds, err := client.Client.ListIp4PoolIds(ctx)
			require.NoError(t, err)
			require.NotZero(t, len(poolIds))
		})
	}
}

func TestGetAllIp4Pools(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			pools, err := client.Client.GetIp4Pools(ctx)
			require.NoError(t, err)
			require.NotZero(t, len(pools))
		})
	}
}

func TestGetIp4PoolByName(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			pools, err := client.Client.GetIp4Pools(ctx)
			require.NoError(t, err)

			poolNames := make([]string, len(pools))
			for i, p := range pools {
				poolNames[i] = p.DisplayName
			}

			for _, name := range poolNames {
				t.Run(name, func(t *testing.T) {
					ctx := testutils.ContextWithTestID(ctx, t)

					pool, err := client.Client.GetIp4PoolByName(ctx, name)
					require.NoError(t, err)

					if pool.Used.Cmp(&pool.Total) == 0 {
						log.Fatal("every IP in the pool is in use? seems unlikely.")
					}

					for _, subnet := range pool.Subnets {
						if subnet.Used.Cmp(&subnet.Total) == 0 {
							log.Fatal("every IP in the subnet is in use? seems unlikely.")
						}
					}
				})
			}
		})
	}
}

func TestCreateGetDeleteIp4Pool(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			req := apstra.NewIpPoolRequest{
				DisplayName: testutils.RandString(10, "hex"),
				Tags:        []string{"tag one", "tag two"},
			}

			id, err := client.Client.CreateIp4Pool(ctx, &req)
			require.NoError(t, err)

			pool, err := client.Client.GetIp4Pool(ctx, id)
			require.NoError(t, err)
			require.NotNil(t, pool)
			compare.IpPool(t, req, *pool)

			pool, err = client.Client.GetIp4PoolByName(ctx, req.DisplayName)
			require.NoError(t, err)
			require.NotNil(t, pool)
			compare.IpPool(t, req, *pool)

			err = client.Client.DeleteIp4Pool(ctx, id)
			require.NoError(t, err)

			var ace apstra.ClientErr

			_, err = client.Client.GetIp4Pool(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, apstra.ErrNotfound, ace.Type())

			_, err = client.Client.GetIp4PoolByName(ctx, req.DisplayName)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, apstra.ErrNotfound, ace.Type())

			err = client.Client.DeleteIp4Pool(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, apstra.ErrNotfound, ace.Type())
		})
	}
}

func TestListIp6Pools(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			poolIds, err := client.Client.ListIp6PoolIds(ctx)
			require.NoError(t, err)
			require.NotZero(t, len(poolIds))
		})
	}
}

func TestGetAllIp6Pools(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			pools, err := client.Client.GetIp6Pools(ctx)
			require.NoError(t, err)
			require.NotZero(t, len(pools))
		})
	}
}

func TestGetIp6PoolByName(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			pools, err := client.Client.GetIp6Pools(ctx)
			require.NoError(t, err)

			poolNames := make([]string, len(pools))
			for i, p := range pools {
				poolNames[i] = p.DisplayName
			}

			for _, name := range poolNames {
				t.Run(name, func(t *testing.T) {
					ctx := testutils.ContextWithTestID(ctx, t)

					pool, err := client.Client.GetIp6PoolByName(ctx, name)
					require.NoError(t, err)

					if pool.Used.Cmp(&pool.Total) == 0 {
						log.Fatal("every IP in the pool is in use? seems unlikely.")
					}

					for _, subnet := range pool.Subnets {
						if subnet.Used.Cmp(&subnet.Total) == 0 {
							log.Fatal("every IP in the subnet is in use? seems unlikely.")
						}
					}
				})
			}
		})
	}
}

func TestCreateGetDeleteIp6Pool(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			req := apstra.NewIpPoolRequest{
				DisplayName: testutils.RandString(10, "hex"),
				Tags:        []string{"tag one", "tag two"},
			}

			id, err := client.Client.CreateIp6Pool(ctx, &req)
			require.NoError(t, err)

			pool, err := client.Client.GetIp6Pool(ctx, id)
			require.NoError(t, err)
			require.NotNil(t, pool)
			compare.IpPool(t, req, *pool)

			pool, err = client.Client.GetIp6PoolByName(ctx, req.DisplayName)
			require.NoError(t, err)
			require.NotNil(t, pool)
			compare.IpPool(t, req, *pool)

			err = client.Client.DeleteIp6Pool(ctx, id)
			require.NoError(t, err)

			var ace apstra.ClientErr

			_, err = client.Client.GetIp6Pool(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, apstra.ErrNotfound, ace.Type())

			_, err = client.Client.GetIp6PoolByName(ctx, req.DisplayName)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, apstra.ErrNotfound, ace.Type())

			err = client.Client.DeleteIp6Pool(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, apstra.ErrNotfound, ace.Type())
		})
	}
}
