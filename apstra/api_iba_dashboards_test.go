//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"reflect"
	"testing"
)

func TestCreateReadUpdateDeleteIbaDashboards(t *testing.T) {

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for clientName, client := range clients {
		log.Printf("testing IBA Dashboard code against %s %s (%s)", client.clientType, clientName,
			client.client.ApiVersion())

		bpClient, bpDelete := testBlueprintA(ctx, t, client.client)
		defer bpDelete(ctx)
		widgetAId, _, widgetBId, _ := testWidgetsAB(ctx, t, bpClient)

		data := IbaDashboardData{
			Description:   "Test Dashboard",
			Default:       false,
			Label:         "Test Dash",
			IbaWidgetGrid: [][]ObjectId{{widgetAId, widgetBId}, {widgetAId, widgetBId}},
		}

		ds, err := bpClient.GetAllIbaDashboards(ctx)
		l := len(ds)
		if l != 0 {
			t.Fatalf("Expected no dashboards. got %d", l)
		}

		t.Logf("Test Create Dashboard")
		id, err := bpClient.CreateIbaDashboard(ctx, &data)
		if err != nil {
			t.Log(data)
			t.Fatal(err)
		}

		data2 := IbaDashboardData{
			Description:   "Test Dashboard Backup",
			Default:       false,
			Label:         "Test Dash B",
			IbaWidgetGrid: [][]ObjectId{{widgetAId, widgetBId}, {widgetAId, widgetBId}},
		}

		t.Logf("Test Create Second Dashboard")
		_, err = bpClient.CreateIbaDashboard(ctx, &data2)
		if err != nil {
			t.Log(data)
			t.Fatal(err)
		}

		ds, err = bpClient.GetAllIbaDashboards(ctx)
		l = len(ds)
		t.Logf("Found %d dashboards", l)
		if l != 2 {
			t.Fatalf("Expected 2 dashboards. got %d", l)
		}

		checkDashes := func() {
			t.Log("Test GetIbaDashboard")
			d1, err := bpClient.GetIbaDashboard(ctx, id)
			if err != nil {
				t.Log(id)
				t.Fatal(err)
			}
			t.Log("Test GetIbaDashboardByLabel")
			d2, err := bpClient.GetIbaDashboardByLabel(ctx, data.Label)
			if err != nil {
				t.Log(data.Label)
				t.Fatal(err)
			}
			t.Log("Dashboard Data")
			t.Log(data)
			t.Log("IBA Probe by Id")
			t.Log(d1)
			t.Log("IBA Dashboard by Name")
			t.Log(d2)

			if !reflect.DeepEqual(d1, d2) {
				t.Fatal("GetIbaDashboardByLabel gets different object than GetIbaDashboard")
			}
			t.Log("Ensure Data matches")
			t.Log(d1.Data)

			if d1.Data.Label != data.Label {
				t.Fatal("IBA Dashboard Label mismatch")
			}
			if d1.Data.Default != data.Default {
				t.Fatal("IBA Dasboard Default mismatch")
			}
			if d1.Data.Description != data.Description {
				t.Fatal("IBA Dashboard Description mismatch")
			}
			if !reflect.DeepEqual(d1.Data.IbaWidgetGrid, data.IbaWidgetGrid) {
				t.Fatal("Widget Grid mismatch")
			}
		}
		checkDashes()

		t.Log("Test Update Dashboard")
		data.Label = "Test Dash 2"
		data.IbaWidgetGrid = append(data.IbaWidgetGrid, []ObjectId{widgetAId, widgetBId})
		data.Description = "Test Dashboard 2"
		err = bpClient.UpdateIbaDashboard(ctx, id, &data)
		if err != nil {
			t.Log(data)
			t.Fatal(err)
		}
		checkDashes()

		t.Log("Test Delete Dashboard")
		err = bpClient.DeleteIbaDashboard(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		_, err = bpClient.GetIbaDashboard(ctx, id)
		if err == nil {
			t.Fatalf("Deleted but id %s is still available", id)
		}
	}
}
