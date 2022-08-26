package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestGetResourceAllocation(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		blueprintIds, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(blueprintIds) == 0 {
			t.Skip("cannot test resource allocation - no blueprints")
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpClient, err := client.client.NewTwoStageL3ClosClient(context.TODO(), blueprintIds[0])
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getResourceAllocation() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		spineAsns, err := bpClient.getResourceAllocation(context.TODO(), &ResourceGroupAllocation{
			Type: ResourceTypeAsnPool,
			Name: ResourceGroupNameSpineAsn,
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Println(spineAsns.PoolIds)
	}
}
