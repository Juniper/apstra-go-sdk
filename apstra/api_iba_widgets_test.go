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
	ctx := context.Background()
	for clientName, client := range clients {
		log.Printf("testing IBA Widget Code against %s %s (%s)", client.clientType, clientName,
			client.client.ApiVersion())

		bpClient, bpDelete := testBlueprintA(ctx, t, client.client)
		defer bpDelete(ctx)

		widgetAId, widgetA, widgetBId, widgetB := testWidgets(ctx, t, bpClient)

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
	}
}
