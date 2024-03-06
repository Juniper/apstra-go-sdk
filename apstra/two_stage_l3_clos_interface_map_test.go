//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestGetSetInterfaceMapAssignments(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintA(ctx, t, client.client)

		log.Printf("testing GetInterfaceMapAssignments() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ifMapAss, err := bpClient.GetInterfaceMapAssignments(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for k, v := range ifMapAss {
			if v == nil {
				v = "<nil>"
			}
			log.Printf("'%s' -> '%s'", k, v)
		}

		// todo check length before using in assignment

		log.Printf("testing SetInterfaceMapAssignments() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.SetInterfaceMapAssignments(ctx, ifMapAss)
		if err != nil {
			t.Fatal(err)
		}
	}
}
