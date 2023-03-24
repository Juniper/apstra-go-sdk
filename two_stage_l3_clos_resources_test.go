//go:build integration
// +build integration

package goapstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestSetGetResourceAllocation(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	poolCount := rand.Intn(5) + 2
	randStr := randString(5, "hex")
	label := "test-" + randStr

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintA(ctx, t, client.client)
		defer func() {
			err := bpDel()
			if err != nil {
				t.Error(err)
			}
		}()

		poolIds := make([]ObjectId, poolCount)
		for i := range poolIds {
			poolId, err := client.client.CreateAsnPool(ctx, &AsnPoolRequest{
				DisplayName: label,
				Ranges:      []IntfIntRange{IntRange{First: 1000, Last: 1999}},
			})
			defer func() {
				err := client.client.DeleteAsnPool(ctx, poolId)
				if err != nil {
					t.Error(err)
				}
			}()
			if err != nil {
				t.Fatal(err)
			}
			poolIds[i] = poolId
		}

		log.Printf("testing SetResourceAllocation() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.SetResourceAllocation(ctx, &ResourceGroupAllocation{
			PoolIds: poolIds,
			ResourceGroup: ResourceGroup{
				Type: ResourceTypeAsnPool,
				Name: ResourceGroupNameSpineAsn,
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetResourceAllocation() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		rga, err := bpClient.GetResourceAllocation(ctx, &ResourceGroup{
			Type: ResourceTypeAsnPool,
			Name: ResourceGroupNameSpineAsn,
		})
		if err != nil {
			t.Fatal(err)
		}

		if rga.ResourceGroup.SecurityZoneId != nil {
			t.Fatal("resource group security zone ID is not nil")
		}

		if len(poolIds) != len(rga.PoolIds) {
			t.Fatalf("expected %d pool IDs, got %d pool IDs", len(poolIds), len(rga.PoolIds))
		}
		log.Println(rga.PoolIds)
	}
}

func TestAllResourceGroupNames(t *testing.T) {
	all := AllResourceGroupNames()
	expected := 16
	if len(all) != expected {
		t.Fatalf("expected %d resource group names, got %d", expected, len(all))
	}
}
