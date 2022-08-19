package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestGetAllVirtualNetworks(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		blueprints, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(blueprints) == 0 {
			t.Skip("no blueprints, no test")
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		dcClient, err := client.client.NewTwoStageL3ClosClient(context.TODO(), blueprints[0])
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing listAllVirtualNetworkIds() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		vns, err := dcClient.listAllVirtualNetworkIds(context.TODO(), BlueprintTypeStaging)
		if err != nil {
			t.Fatal(err)
		}

		for _, id := range vns {
			log.Printf("testing getVirtualNetwork() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
			vn, err := dcClient.getVirtualNetwork(context.TODO(), id, BlueprintTypeStaging)
			if err != nil {
				t.Fatal(err)
			}
			if vn.Ipv4Subnet != nil {
				log.Printf("testing getVirtualNetworkBySubnet() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
				vn, err = dcClient.getVirtualNetworkBySubnet(context.TODO(), vn.Ipv4Subnet, vn.SecurityZoneId, BlueprintTypeStaging)
				if err != nil {
					t.Fatal(err)
				}
			}
			if vn.Ipv6Subnet != nil {
				log.Printf("testing getVirtualNetworkBySubnet() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
				vn, err = dcClient.getVirtualNetworkBySubnet(context.TODO(), vn.Ipv6Subnet, vn.SecurityZoneId, BlueprintTypeStaging)
				if err != nil {
					t.Fatal(err)
				}
			}
			log.Printf("vn: %s ipv4: %s, ipv6:%s\n", vn.Id, vn.Ipv4Subnet.String(), vn.Ipv6Subnet.String())
		}
	}
}
