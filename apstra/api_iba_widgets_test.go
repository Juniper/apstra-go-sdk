//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestIbaWidgetsGet(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	for clientName, client := range clients {
		log.Printf("testing GetAllIbaWidgets against %s %s (%s)", client.clientType, clientName,
			client.client.ApiVersion())

		bpids, err := client.client.ListAllBlueprintIds(ctx)
		if err != nil {
			t.Fatal(err)
		}

		// IBA probes will not exist until the blueprint is deployed. This test expects that there will be a blueprint
		// with existing IBA probes
		bpClient, err := client.client.NewTwoStageL3ClosClient(ctx, bpids[0])
		if err != nil {
			t.Fatal(err)
		}

		widgets, err := bpClient.GetAllIbaWidgets(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if len(widgets) <= 0 {
			t.Fatalf("only got %d widgets", len(widgets))
		}
		for _, w := range widgets {
			ws, err := bpClient.GetIbaWidgetsByLabel(ctx, w.Data.Label)
			if err != nil {
				t.Fatal(err)
			}
			if len(ws) > 1 {
				t.Fatalf("Was expecting only 1 widget with name %s got %d", w.Data.Label, len(ws))
			}
			if ws[0].Id != w.Id {
				t.Fatalf("GetIbaWidgetsByLabel returned a different id than the original. Expected %s. Got %s",
					w.Id, ws[0].Id)
			}
			t.Logf("Found Widget Label %s ID %s type %s", w.Id, w.Data.Label, w.Data.Type)
		}
	}
}
