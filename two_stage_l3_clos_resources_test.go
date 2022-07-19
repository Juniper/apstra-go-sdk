package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestGetResourceAllocation(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	blueprintIds, err := client.listAllBlueprintIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(blueprintIds) == 0 {
		t.Skip("cannot test resource allocation - no blueprints")
	}

	bpClient, err := client.NewTwoStageL3ClosClient(context.TODO(), blueprintIds[0])
	if err != nil {
		t.Fatal(err)
	}

	spineAsns, err := bpClient.getResourceAllocation(context.TODO(), &ResourceGroupAllocation{
		Type: ResourceTypeAsnPool,
		Name: ResourceGroupNameSpineAsn,
	})
	log.Println(spineAsns.PoolIds)
}
