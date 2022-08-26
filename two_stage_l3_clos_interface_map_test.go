package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestGetSetInterfaceMapAssignments(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpIds, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(bpIds) == 0 {
			t.Skip("cannot get interface map assignments with no blueprints")
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpClient, err := client.client.NewTwoStageL3ClosClient(context.TODO(), "d7ff0cbb-3cba-48b6-9271-9c6d7aef8b46")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetInterfaceMapAssignments() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ifMapAss, err := bpClient.GetInterfaceMapAssignments(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		for _, i := range ifMapAss {
			log.Println(i)
		}

		// todo check length before using in assignment

		log.Printf("testing SetInterfaceMapAssignments() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.SetInterfaceMapAssignments(context.TODO(), ifMapAss)
		if err != nil {
			t.Fatal(err)
		}
	}
}
