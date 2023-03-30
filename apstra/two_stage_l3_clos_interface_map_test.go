//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestGetSetInterfaceMapAssignments(t *testing.T) {
	clients, err := getTestClients(context.Background())
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

		// reduce bpIds to just "datacenter" blueprints
		for i := len(bpIds) - 1; i >= 0; i-- {
			bpStatus, err := client.client.getBlueprintStatus(context.TODO(), bpIds[i])
			if err != nil {
				t.Fatal(err)
			}
			if bpStatus.Design != refDesignDatacenter {
				bpIds[i] = bpIds[len(bpIds)-1] // move last element to current position
				bpIds = bpIds[:len(bpIds)-1]   // remove last element
			}
		}

		if len(bpIds) == 0 {
			skipMsg[clientName] = fmt.Sprintf("cannot get interface map assignments - no 'datacenter' blueprint in '%s'", clientName)
			continue
		}

		for _, bpId := range bpIds {
			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bpClient, err := client.client.NewTwoStageL3ClosClient(context.TODO(), bpId)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetInterfaceMapAssignments() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			ifMapAss, err := bpClient.GetInterfaceMapAssignments(context.TODO())
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
			err = bpClient.SetInterfaceMapAssignments(context.TODO(), ifMapAss)
			if err != nil {
				t.Fatal(err)
			}
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
