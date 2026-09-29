// Copyright (c) Juniper Networks, Inc., 2022-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func TestSetGetResourceAllocation(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	poolCount := rand.Intn(5) + 2
	randStr := randString(5, "hex")
	label := "test-" + randStr

	for clientName, client := range clients {
		clientName, client := clientName, client // local copy of iterator variables safe for use in deferred function
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bpWait := sync.WaitGroup{}
			bpWait.Add(1)
			bpClient := testBlueprintB(ctx, t, client.client)

			poolIds := make([]string, poolCount)
			for i := range poolIds {
				poolId, err := client.client.CreateAsnPool(ctx, &AsnPoolRequest{
					DisplayName: label + "-" + strconv.Itoa(i),
					Ranges: []IntfIntRange{IntRange{
						First: uint32(1000 + (i * 1000)),
						Last:  uint32(1999 + (i * 1000)),
					}},
				})
				require.NoError(t, err)

				poolIds[i] = string(poolId)
				defer func() {
					go func() {
						bpWait.Wait()
						require.NoError(t, client.client.DeleteAsnPool(ctx, poolId))
					}()
				}()
			}

			log.Printf("testing SetResourceAllocation() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			require.NoError(t, bpClient.SetResourceAllocation(ctx, &ResourceGroupAllocation{
				PoolIds: poolIds,
				ResourceGroup: ResourceGroup{
					Type: enum.ResourceTypeASN,
					Name: enum.ResourceGroupSpineASN,
				},
			}))

			log.Printf("testing GetResourceAllocation() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			rga, err := bpClient.GetResourceAllocation(ctx, &ResourceGroup{
				Type: enum.ResourceTypeASN,
				Name: enum.ResourceGroupSpineASN,
			})
			require.NoError(t, err)

			require.Nilf(t, rga.ResourceGroup.SecurityZoneId, "resource group security zone ID must be nil")
			require.Equalf(t, len(poolIds), len(rga.PoolIds), "expected pool ID count (%d) must equal actual pool ID count (%d)", len(poolIds), len(rga.PoolIds))
		})
	}
}
