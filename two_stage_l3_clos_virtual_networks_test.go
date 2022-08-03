package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestGetAllVirtualNetworks(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	blueprints, err := client.listAllBlueprintIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(blueprints) == 0 {
		t.Skip("no blueprints, no test")
	}

	dcClient, err := client.NewTwoStageL3ClosClient(context.TODO(), blueprints[0])
	if err != nil {
		t.Fatal(err)
	}

	vns, err := dcClient.listAllVirtualNetworkIds(context.TODO(), BlueprintTypeStaging)
	if err != nil {
		t.Fatal(err)
	}

	for _, id := range vns {
		vn, err := dcClient.getVirtualNetwork(context.TODO(), id, BlueprintTypeStaging)
		if err != nil {
			t.Fatal(err)
		}
		if vn.Ipv4Subnet != nil {
			vn, err = dcClient.getVirtualNetworkBySubnet(context.TODO(), vn.Ipv4Subnet, vn.SecurityZoneId, BlueprintTypeStaging)
			if err != nil {
				t.Fatal(err)
			}
		}
		if vn.Ipv6Subnet != nil {
			vn, err = dcClient.getVirtualNetworkBySubnet(context.TODO(), vn.Ipv6Subnet, vn.SecurityZoneId, BlueprintTypeStaging)
			if err != nil {
				t.Fatal(err)
			}
		}
		log.Printf("vn: %s ipv4: %s, ipv6:%s\n", vn.Id, vn.Ipv4Subnet.String(), vn.Ipv6Subnet.String())
	}

}
