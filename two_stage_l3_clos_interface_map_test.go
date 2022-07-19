package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestGetSetInterfaceMapAssignments(t *testing.T) {
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

	bpClient, err := client.NewTwoStageL3ClosClient(context.TODO(), "d7ff0cbb-3cba-48b6-9271-9c6d7aef8b46")

	ifMapAss, err := bpClient.GetInterfaceMapAssignments(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range ifMapAss {
		log.Println(i)
	}

	// todo check length before using in assignment

	err = bpClient.SetInterfaceMapAssignments(context.TODO(), ifMapAss)
	if err != nil {
		t.Fatal(err)
	}
}

//func TestSetInterfaceMapAssignments(t *testing.T) {
//
//}
