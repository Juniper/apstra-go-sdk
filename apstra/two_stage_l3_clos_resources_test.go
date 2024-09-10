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
			defer func() {
				require.NoError(t, client.client.DeleteBlueprint(ctx, bpClient.Id()))
				bpWait.Done() // signal that blueprint is deleted so that ASN pools can be removed
			}()

			poolIds := make([]ObjectId, poolCount)
			for i := range poolIds {
				poolId, err := client.client.CreateAsnPool(ctx, &AsnPoolRequest{
					DisplayName: label + "-" + strconv.Itoa(i),
					Ranges: []IntfIntRange{IntRange{
						First: uint32(1000 + (i * 1000)),
						Last:  uint32(1999 + (i * 1000)),
					}},
				})
				require.NoError(t, err)

				poolIds[i] = poolId
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
					Type: ResourceTypeAsnPool,
					Name: ResourceGroupNameSpineAsn,
				},
			}))

			log.Printf("testing GetResourceAllocation() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			rga, err := bpClient.GetResourceAllocation(ctx, &ResourceGroup{
				Type: ResourceTypeAsnPool,
				Name: ResourceGroupNameSpineAsn,
			})
			require.NoError(t, err)

			require.Nilf(t, rga.ResourceGroup.SecurityZoneId, "resource group security zone ID must be nil")
			require.Equalf(t, len(poolIds), len(rga.PoolIds), "expected pool ID count (%d) must equal actual pool ID count (%d)", len(poolIds), len(rga.PoolIds))
		})
	}
}

func TestAllResourceGroupNames(t *testing.T) {
	all := AllResourceGroupNames()
	expected := 32
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
		{stringVal: "generic_asns", intType: ResourceGroupNameGenericAsn, stringType: resourceGroupNameGenericAsn},
		{stringVal: "superspine_loopback_ips", intType: ResourceGroupNameSuperspineIp4, stringType: resourceGroupNameSuperspineIp4},
		{stringVal: "superspine_loopback_ips_ipv6", intType: ResourceGroupNameSuperspineIp6, stringType: resourceGroupNameSuperspineIp6},
		{stringVal: "spine_loopback_ips", intType: ResourceGroupNameSpineIp4, stringType: resourceGroupNameSpineIp4},
		{stringVal: "spine_loopback_ips_ipv6", intType: ResourceGroupNameSpineIp6, stringType: resourceGroupNameSpineIp6},
		{stringVal: "leaf_loopback_ips", intType: ResourceGroupNameLeafIp4, stringType: resourceGroupNameLeafIp4},
		{stringVal: "leaf_loopback_ips_ipv6", intType: ResourceGroupNameLeafIp6, stringType: resourceGroupNameLeafIp6},
		{stringVal: "access_loopback_ips", intType: ResourceGroupNameAccessIp4, stringType: resourceGroupNameAccessIp4},
		{stringVal: "generic_loopback_ips", intType: ResourceGroupNameGenericIp4, stringType: resourceGroupNameGenericIp4},
		{stringVal: "generic_loopback_ips_ipv6", intType: ResourceGroupNameGenericIp6, stringType: resourceGroupNameGenericIp6},
		{stringVal: "spine_superspine_link_ips", intType: ResourceGroupNameSuperspineSpineIp4, stringType: resourceGroupNameSuperspineSpineIp4},
		{stringVal: "ipv6_spine_superspine_link_ips", intType: ResourceGroupNameSuperspineSpineIp6, stringType: resourceGroupNameSuperspineSpineIp6},
		{stringVal: "spine_leaf_link_ips", intType: ResourceGroupNameSpineLeafIp4, stringType: resourceGroupNameSpineLeafIp4},
		{stringVal: "ipv6_spine_leaf_link_ips", intType: ResourceGroupNameSpineLeafIp6, stringType: resourceGroupNameSpineLeafIp6},
		{stringVal: "access_l3_peer_link_link_ips", intType: ResourceGroupNameAccessAccessIp4, stringType: resourceGroupNameAccessAccessIp4},
		{stringVal: "leaf_leaf_link_ips", intType: ResourceGroupNameLeafLeafIp4, stringType: resourceGroupNameLeafLeafIp4},
		{stringVal: "leaf_l3_peer_link_link_ips", intType: ResourceGroupNameLeafL3PeerLinkLinkIp4, stringType: resourceGroupNameLeafL3PeerLinkLinkIp4},
		{stringVal: "ipv6_leaf_l3_peer_link_link_ips", intType: ResourceGroupNameLeafL3PeerLinkLinkIp6, stringType: resourceGroupNameLeafL3PeerLinkLinkIp6},
		{stringVal: "mlag_domain_svi_subnets", intType: ResourceGroupNameMlagDomainIp4, stringType: resourceGroupNameMlagDomainSviIp4},
		{stringVal: "mlag_domain_svi_subnets_ipv6", intType: ResourceGroupNameMlagDomainIp6, stringType: resourceGroupNameMlagDomainSviIp6},
		{stringVal: "vtep_ips", intType: ResourceGroupNameVtepIp4, stringType: resourceGroupNameVtepIp4},
		{stringVal: "evpn_l3_vnis", intType: ResourceGroupNameEvpnL3Vni, stringType: resourceGroupNameEvpnL3Vni},
		{stringVal: "virtual_network_svi_subnets", intType: ResourceGroupNameVirtualNetworkSviIpv4, stringType: resourceGroupNameVirtualNetworkSviIpv4},
		{stringVal: "virtual_network_svi_subnets_ipv6", intType: ResourceGroupNameVirtualNetworkSviIpv6, stringType: resourceGroupNameVirtualNetworkSviIpv6},
		{stringVal: "vxlan_vn_ids", intType: ResourceGroupNameVxlanVnIds, stringType: resourceGroupNameVxlanVnIds},
		{stringVal: "to_generic_link_ips", intType: ResourceGroupNameToGenericLinkIpv4, stringType: resourceGroupNameToGenericLinkIpv4},
		{stringVal: "ipv6_to_generic_link_ips", intType: ResourceGroupNameToGenericLinkIpv6, stringType: resourceGroupNameToGenericLinkIpv6},
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
