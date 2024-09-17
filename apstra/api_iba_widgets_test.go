//go:build integration

package apstra

import (
	"context"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"log"
	"testing"
)

func TestIbaWidgets(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			if !compatibility.IbaWidgetSupported.Check(client.client.apiVersion) {
				t.Skipf("skipping due to IbaWidgetSupported --  need to implement.")
			}

			log.Printf("testing IBA Widget Code against %s %s (%s)", client.clientType, clientName,
				client.client.ApiVersion())

			bpClient := testBlueprintA(ctx, t, client.client)

			widgetAId, widgetA, widgetBId, widgetB := testWidgetsAB(ctx, t, bpClient)

			widgets, err := bpClient.GetAllIbaWidgets(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if len(widgets) != 2 {
				t.Fatalf("expected 2 widgets, got %d widgets", len(widgets))
			}

			wa, err := bpClient.GetIbaWidget(ctx, widgetAId)
			if err != nil {
				t.Fatal(err)
			}
			if wa.Id != widgetAId {
				t.Fatalf("expected wiget A ID %q, got %q", widgetAId, wa.Id)
			}
			if wa.Data.Label != widgetA.Label {
				t.Fatalf("expected wiget A Label %q, got %q", widgetA.Label, wa.Data.Label)
			}

			wb, err := bpClient.GetIbaWidget(ctx, widgetBId)
			if err != nil {
				t.Fatal(err)
			}
			if wb.Id != widgetBId {
				t.Fatalf("expected wiget B ID %q, got %q", widgetBId, wb.Id)
			}
			if wb.Data.Label != widgetB.Label {
				t.Fatalf("expected wiget B Label %q, got %q", widgetB.Label, wb.Data.Label)
			}

			for _, widget := range widgets {
				ws, err := bpClient.GetIbaWidgetByLabel(ctx, widget.Data.Label)
				if err != nil {
					t.Fatal(err)
				}
				if ws.Id != widget.Id {
					t.Fatalf("GetIbaWidgetsByLabel returned a different id than the original. Expected %s. Got %s",
						widget.Id, ws.Id)
				}
			}

			t.Logf("Update %s", wa.Data.Label)
			wa.Data.Description = "This widget now updated"
			err = bpClient.UpdateIbaWidget(ctx, widgetAId, wa.Data)
			if err != nil {
				t.Fatal(err)
			}
			wa1, err := bpClient.GetIbaWidget(ctx, widgetAId)
			if err != nil {
				t.Fatal(err)
			}
			if wa.Data.Description != wa1.Data.Description {
				t.Fatal("Looks like the update failed")
			}
			t.Logf("Test Deletion of %s", wa.Data.Label)
			err = bpClient.DeleteIbaWidget(ctx, widgetAId)
			if err != nil {
				t.Fatal(err)
			}

			wa, err = bpClient.GetIbaWidget(ctx, widgetAId)
			if err == nil {
				t.Fatalf("Widget with id %s should have been deleted", widgetAId)
			}
		})
	}
}
