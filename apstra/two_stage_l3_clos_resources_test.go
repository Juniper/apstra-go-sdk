//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"sync"
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
		client := client // local copy of iterator variable safe for use in deferred function
		asnPoolWait := sync.WaitGroup{}
		bpWait := sync.WaitGroup{}
		bpWait.Add(1)
		bpClient, bpDel := testBlueprintB(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Error(err)
			}
			bpWait.Done()      // resource pool deletion is waiting for this
			asnPoolWait.Wait() // wait for resource pool deletion to complete
		}()
		poolIds := make([]ObjectId, poolCount)
		for i := range poolIds {
			asnPoolWait.Add(1)
			poolId, err := client.client.CreateAsnPool(ctx, &AsnPoolRequest{
				DisplayName: label + "-" + strconv.Itoa(i),
				Ranges: []IntfIntRange{IntRange{
					First: uint32(1000 + (i * 1000)),
					Last:  uint32(1999 + (i * 1000)),
				}},
			})
			if err != nil {
				t.Fatal(err)
			}
			poolIds[i] = poolId
			defer func() {
				go func() {
					log.Printf("waiting before deleting pool %q", poolId)
					bpWait.Wait()
					log.Printf("deleting pool %q now", poolId)
					err := client.client.DeleteAsnPool(ctx, poolId)
					if err != nil {
						t.Error(err)
					}
					asnPoolWait.Done()
				}()
			}()
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
	expected := 17
	if len(all) != expected {
		t.Fatalf("expected %d resource group names, got %d", expected, len(all))
	}
}

func TestAllResourceTypes(t *testing.T) {
	all := AllResourceTypes()
	expected := 5
	if len(all) != expected {
		t.Fatalf("expected %d resource types, got %d", expected, len(all))
	}
}
