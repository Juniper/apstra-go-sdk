//go:build integration
// +build integration

package apstra //

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestIbaPredefinedProbes(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for clientName, client := range clients {
		log.Printf("testing Predefined Probes against %s %s (%s)", client.clientType, clientName,
			client.client.ApiVersion())

		bpClient := testBlueprintA(ctx, t, client.client)
		pdps, err := bpClient.GetAllIbaPredefinedProbes(ctx)
		require.NoError(t, err)
		expectedToFail := map[string]bool{
			"external_ecmp_imbalance":            true,
			"evpn_vxlan_type5":                   true,
			"eastwest_traffic":                   true,
			"vxlan_floodlist":                    true,
			"fabric_hotcold_ifcounter":           true,
			"specific_interface_flapping":        true,
			"evpn_vxlan_type3":                   true,
			"specific_hotcold_ifcounter":         true,
			"spine_superspine_hotcold_ifcounter": true,
		}
		t.Logf("Try an obviously fake name : %s", "FAKE")
		_, err = bpClient.GetIbaPredefinedProbeByName(ctx, "FAKE")
		if err == nil {
			t.Fatal("FAKE name should have failed, but succeeded")
		} else {
			t.Log(err)
		}

		for _, p := range pdps {
			t.Logf("Get Predefined Probe By Name %s", p.Name)
			_, err := bpClient.GetIbaPredefinedProbeByName(ctx, p.Name)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(p.Description)
			t.Log(p.Schema)

			t.Logf("Instantiating Probe %s", p.Name)

			probeId, err := bpClient.InstantiateIbaPredefinedProbe(ctx, &IbaPredefinedProbeRequest{
				Name: p.Name,
				Data: json.RawMessage([]byte(`{"label":"` + p.Name + `"}`)),
			})
			if err != nil {
				if !expectedToFail[p.Name] {
					t.Fatal(err)
				} else {
					t.Logf("%s was expected to fail", p.Name)
					continue
				}
			}

			t.Logf("Got back Probe Id %s \n Now Make a Widget with it.", probeId)

			widgetId, err := bpClient.CreateIbaWidget(ctx, &IbaWidgetData{
				Type:      IbaWidgetTypeStage,
				ProbeId:   probeId,
				Label:     p.Name,
				StageName: p.Name,
			})
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Got back Widget Id %s \n Now fetch it.", widgetId)

			widget, err := bpClient.GetIbaWidget(ctx, widgetId)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Widget %s created", widget.Data.Label)

			t.Logf("Try to Delete Probe this should fail because a widget is using it")
			err = bpClient.DeleteIbaProbe(ctx, probeId)
			if err == nil {
				t.Fatal("Probe Deletion should have failed")
			} else {
				t.Log(err)
			}

			t.Logf("Delete Widget and then the probe this path should succeed")
			err = bpClient.DeleteIbaWidget(ctx, widgetId)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Delete probe")

			err = bpClient.DeleteIbaProbe(ctx, probeId)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
