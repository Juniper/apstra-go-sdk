package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestGetInterfaceMapAssignments(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	bpIds, err := client.listAllBlueprintIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(bpIds) == 0 {
		t.Skip("cannot get interface map assignments with no blueprints")
	}

	bpClient, err := client.NewTwoStageL3ClosClient(context.TODO(), "1ed198c3-ccac-4adc-917b-4300eaab7f8e")

	ifMapAss, err := bpClient.GetInterfaceMapAssignments(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range ifMapAss {
		log.Println(i)
	}
}
