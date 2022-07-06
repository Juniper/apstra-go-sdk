package goapstra

import (
	"context"
	"crypto/tls"
	"log"
	"testing"
)

func refDesignBlueprintTestclient1() (*Client, error) {
	return NewClient(&ClientCfg{
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	})
}

func TestGetResourceAllocation(t *testing.T) {
	client, err := refDesignBlueprintTestclient1()
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

	bp, err := client.GetBlueprint(context.TODO(), blueprintIds[0])
	spineAsns, err := client.getResourceAllocation(context.TODO(), bp.Id, &ResourceGroupAllocation{
		Type: ResourceTypeAsnPool,
		Name: ResourceGroupNameSpineAsn,
	})
	log.Println(spineAsns.PoolIds)
}
