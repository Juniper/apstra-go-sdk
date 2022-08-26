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

	skipped := true

	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		blueprints, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(blueprints) == 0 {
			continue
		}

		skipped = false

		for i := range samples(len(blueprints)) {
			id := blueprints[i]
			bpStatus, err := client.client.getBlueprintStatus(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
			switch bpStatus.Design {
			case RefDesignFreeform:
				log.Printf("'%s' design is '%s.", id, bpStatus.Design.String())
				// todo
			case RefDesignDatacenter:
				log.Printf("'%s' design is '%s.", id, bpStatus.Design.String())
				log.Printf("testing NewTwoStageL3ClosClient(%s) against %s %s (%s)", id, client.clientType, clientName, client.client.ApiVersion())
				dcClient, err := client.client.NewTwoStageL3ClosClient(context.TODO(), id)
				if err != nil {
					t.Fatal(err)
				}

				log.Printf("testing listAllVirtualNetworkIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				vns, err := dcClient.listAllVirtualNetworkIds(context.TODO(), BlueprintTypeStaging)
				if err != nil {
					t.Fatal(err)
				}

				for _, id := range vns {
					log.Printf("testing getVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					vn, err := dcClient.getVirtualNetwork(context.TODO(), id, BlueprintTypeStaging)
					if err != nil {
						t.Fatal(err)
					}
					if vn.Ipv4Subnet != nil {
						log.Printf("testing getVirtualNetworkBySubnet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
						vn, err = dcClient.getVirtualNetworkBySubnet(context.TODO(), vn.Ipv4Subnet, vn.SecurityZoneId, BlueprintTypeStaging)
						if err != nil {
							t.Fatal(err)
						}
					}
					if vn.Ipv6Subnet != nil {
						log.Printf("testing getVirtualNetworkBySubnet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
						vn, err = dcClient.getVirtualNetworkBySubnet(context.TODO(), vn.Ipv6Subnet, vn.SecurityZoneId, BlueprintTypeStaging)
						if err != nil {
							t.Fatal(err)
						}
					}
					log.Printf("vn: %s ipv4: %s, ipv6:%s\n", vn.Id, vn.Ipv4Subnet.String(), vn.Ipv6Subnet.String())
				}
			}
		}
	}
	if skipped {
		t.Skip("no blueprints found on any test system")
	}
}
