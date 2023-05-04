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
	expected := 18
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

func TestTwoStageL3ClosResourceStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		Int() int
	}

	type apiIotaString interface {
		parse() (int, error)
		string() string
	}

	type stringTestData struct {
		stringVal  string
		intType    apiStringIota
		stringType apiIotaString
	}
	testData := []stringTestData{
		{stringVal: "", intType: ResourceTypeNone, stringType: resourceTypeNone},
		{stringVal: "asn", intType: ResourceTypeAsnPool, stringType: resourceTypeAsnPool},
		{stringVal: "ip", intType: ResourceTypeIp4Pool, stringType: resourceTypeIp4Pool},
		{stringVal: "ipv6", intType: ResourceTypeIp6Pool, stringType: resourceTypeIp6Pool},
		{stringVal: "vni", intType: ResourceTypeVniPool, stringType: resourceTypeVniPool},

		{stringVal: "", intType: ResourceGroupNameNone, stringType: resourceGroupNameNone},
		{stringVal: "superspine_asns", intType: ResourceGroupNameSuperspineAsn, stringType: resourceGroupNameSuperspineAsn},
		{stringVal: "spine_asns", intType: ResourceGroupNameSpineAsn, stringType: resourceGroupNameSpineAsn},
		{stringVal: "leaf_asns", intType: ResourceGroupNameLeafAsn, stringType: resourceGroupNameLeafAsn},
		{stringVal: "access_asns", intType: ResourceGroupNameAccessAsn, stringType: resourceGroupNameAccessAsn},
		{stringVal: "superspine_loopback_ips", intType: ResourceGroupNameSuperspineIp4, stringType: resourceGroupNameSuperspineIp4},
		{stringVal: "spine_loopback_ips", intType: ResourceGroupNameSpineIp4, stringType: resourceGroupNameSpineIp4},
		{stringVal: "leaf_loopback_ips", intType: ResourceGroupNameLeafIp4, stringType: resourceGroupNameLeafIp4},
		{stringVal: "access_loopback_ips", intType: ResourceGroupNameAccessIp4, stringType: resourceGroupNameAccessIp4},
		{stringVal: "spine_superspine_link_ips", intType: ResourceGroupNameSuperspineSpineIp4, stringType: resourceGroupNameSuperspineSpineIp4},
		{stringVal: "ipv6_spine_superspine_link_ips", intType: ResourceGroupNameSuperspineSpineIp6, stringType: resourceGroupNameSuperspineSpineIp6},
		{stringVal: "spine_leaf_link_ips", intType: ResourceGroupNameSpineLeafIp4, stringType: resourceGroupNameSpineLeafIp4},
		{stringVal: "ipv6_spine_leaf_link_ips", intType: ResourceGroupNameSpineLeafIp6, stringType: resourceGroupNameSpineLeafIp6},
		{stringVal: "access_l3_peer_link_link_ips", intType: ResourceGroupNameAccessAccessIp4, stringType: resourceGroupNameAccessAccessIp4},
		{stringVal: "leaf_leaf_link_ips", intType: ResourceGroupNameLeafLeafIp4, stringType: resourceGroupNameLeafLeafIp4},
		{stringVal: "mlag_domain_svi_subnets", intType: ResourceGroupNameMlagDomainIp4, stringType: resourceGroupNameMlagDomainSviIp4},
		{stringVal: "vtep_ips", intType: ResourceGroupNameVtepIp4, stringType: resourceGroupNameVtepIp4},
		{stringVal: "evpn_l3_vnis", intType: ResourceGroupNameEvpnL3Vni, stringType: resourceGroupNameEvpnL3Vni},
		{stringVal: "virtual_network_svi_subnets", intType: ResourceGroupNameVirtualNetworkSviIpv4, stringType: resourceGroupNameVirtualNetworkSviIpv4},
	}

	for i, td := range testData {
		ii := td.intType.Int()
		is := td.intType.String()
		sp, err := td.stringType.parse()
		if err != nil {
			t.Fatal(err)
		}
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.Int() != sp ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}
