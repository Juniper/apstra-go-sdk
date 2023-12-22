//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestCreateDeleteRack(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	rackTypeId := ObjectId("L2_Virtual")

	var rt *RackType
	for _, client := range clients {
		rt, err = client.client.GetRackType(ctx, rackTypeId)
		if err != nil {
			t.Fatal(err)
		}
		break
	}

	if rt == nil {
		t.Fatal("failed to collect rack type data")
	}

	testCases := map[string]TwoStageL3ClosRackRequest{
		"rack-by-id": {
			PodId:      "",
			RackTypeId: rackTypeId,
		},
		"rack-by-obj": {
			PodId: "",
			RackType: &RackTypeRequest{
				DisplayName:              "test",
				Description:              "test rack",
				FabricConnectivityDesign: FabricConnectivityDesignL3Clos,
				LeafSwitches: []RackElementLeafSwitchRequest{
					{
						Label:              "leaf",
						LinkPerSpineCount:  1,
						LinkPerSpineSpeed:  "10G",
						RedundancyProtocol: LeafRedundancyProtocolNone,
						Tags:               []ObjectId{"hypervisor", "firewall"},
						LogicalDeviceId:    "AOS-7x10-Leaf",
					},
				},
			},
		},
	}

	for clientName, client := range clients {
		bp, bpDel := testBlueprintD(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		for tName, tCase := range testCases {
			log.Printf("testing CreateRack(%s) against %s %s (%s)", tName, client.clientType, clientName, client.client.ApiVersion())
			id, err := bp.CreateRack(ctx, &tCase)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing DeleteRack() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bp.DeleteRack(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
