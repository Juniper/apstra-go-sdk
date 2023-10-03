//go:build integration
// +build integration

package apstra

import (
	"context"
	"encoding/json"
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

		probeAId, err := bpClient.InstantiateIbaPredefinedProbe(ctx, &IbaPredefinedProbeRequest{
			Name: "bgp_session",
			Data: json.RawMessage([]byte(`{
			"Label":     "BGP Session Flapping",
			"Duration":  300,
			"Threshold": 40
		}`)),
		})
		if err != nil {
			t.Fatal(err)
		}

		probeBId, err := bpClient.InstantiateIbaPredefinedProbe(ctx, &IbaPredefinedProbeRequest{
			Name: "drain_node_traffic_anomaly",
			Data: json.RawMessage([]byte(`{
			"Label":     "Drain Traffic Anomaly",
			"Threshold": 100000
		}`)),
		})

		if err != nil {
			t.Fatal(err)
		}

		widgetA := IbaWidgetData{
			Type:      IbaWidgetTypeStage,
			Label:     "BGP Session Flapping",
			ProbeId:   probeAId,
			StageName: "BGP Session",
		}
		widgetAId, err := bpClient.CreateIbaWidget(ctx, &widgetA)
		if err != nil {
			t.Fatal(err)
		}

		widgetB := IbaWidgetData{
			Type:      IbaWidgetTypeStage,
			Label:     "Drain Traffic Anomaly",
			ProbeId:   probeBId,
			StageName: "excess_range",
		}
		widgetBId, err := bpClient.CreateIbaWidget(ctx, &widgetB)

		if err != nil {
			t.Fatal(err)
		}

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
