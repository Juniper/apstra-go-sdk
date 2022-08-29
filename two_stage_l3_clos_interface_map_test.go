package goapstra

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestGetSetInterfaceMapAssignments(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	skipMsg := make(map[string]string)
	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpIds, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(bpIds) == 0 {
			skipMsg[clientName] = fmt.Sprintf("cannot get interface map assignments - no blueprint in '%s'", clientName)
			continue
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
	if len(skipMsg) > 0 {
		sb := strings.Builder{}
		for _, msg := range skipMsg {
			sb.WriteString(msg + ";")
		}
		t.Skip(sb.String())
	}
}
