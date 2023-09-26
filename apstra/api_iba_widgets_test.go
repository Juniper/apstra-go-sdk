//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestIbaWidgetsGet(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for clientName, client := range clients {
		log.Printf("testing GetAllIbaWidgets against %s %s (%s)", client.clientType, clientName,
			client.client.ApiVersion())

		bpClient, bpDelete := testBlueprintA(ctx, t, client.client)
		defer bpDelete(ctx)
		var idResponse objectIdResponse

		probeA := struct {
			Label     string `json:"label"`
			Duration  int    `json:"duration"`
			Threshold int    `json:"threshold"`
		}{
			Label:     "BGP Session Flapping",
			Duration:  300,
			Threshold: 40,
		}
		probeAUrlStr := fmt.Sprintf(apiUrlBlueprintByIdPrefix, bpClient.blueprintId) + "iba/predefined-probes/bgp_session"

		probeB := struct {
			Label     string `json:"label"`
			Threshold int    `json:"threshold"`
		}{
			Label:     "Drain Traffic Anomaly",
			Threshold: 100000,
		}
		probeBUrlStr := fmt.Sprintf(apiUrlBlueprintByIdPrefix, bpClient.blueprintId) + "iba/predefined-probes/drain_node_traffic_anomaly"

		err = client.client.talkToApstra(ctx, &talkToApstraIn{
			method:         http.MethodPost,
			urlStr:         probeAUrlStr,
			apiInput:       &probeA,
			apiResponse:    &idResponse,
			doNotLogin:     false,
			unsynchronized: false,
		})
		if err != nil {
			t.Fatal(err)
		}
		probeAId := idResponse.Id

		err = client.client.talkToApstra(ctx, &talkToApstraIn{
			method:         http.MethodPost,
			urlStr:         probeBUrlStr,
			apiInput:       &probeB,
			apiResponse:    &idResponse,
			doNotLogin:     false,
			unsynchronized: false,
		})
		if err != nil {
			t.Fatal(err)
		}
		probeBId := idResponse.Id

		widgetA := rawIbaWidget{
			Type:      "stage",
			Label:     probeA.Label,
			ProbeId:   probeAId.String(),
			StageName: "BGP Session",
		}

		widgetB := rawIbaWidget{
			Type:      "stage",
			Label:     probeB.Label,
			ProbeId:   probeBId.String(),
			StageName: "excess_range",
		}

		err = client.client.talkToApstra(ctx, &talkToApstraIn{
			method:      http.MethodPost,
			urlStr:      fmt.Sprintf(apiUrlBlueprintByIdPrefix, bpClient.blueprintId) + "iba/widgets",
			apiInput:    &widgetA,
			apiResponse: &idResponse,
		})
		if err != nil {
			t.Fatal(err)
		}
		widgetAId := idResponse.Id

		err = client.client.talkToApstra(ctx, &talkToApstraIn{
			method:      http.MethodPost,
			urlStr:      fmt.Sprintf(apiUrlBlueprintByIdPrefix, bpClient.blueprintId) + "iba/widgets",
			apiInput:    &widgetB,
			apiResponse: &idResponse,
		})
		if err != nil {
			t.Fatal(err)
		}
		widgetBId := idResponse.Id

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
