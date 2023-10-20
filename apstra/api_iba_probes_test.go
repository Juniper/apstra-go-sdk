//go:build integration
// +build integration

package apstra //

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

func TestIbaProbes(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for clientName, client := range clients {
		log.Printf("testing Predefined Probes against %s %s (%s)", client.clientType, clientName,
			client.client.ApiVersion())

		bpClient, _ := testBlueprintA(ctx, t, client.client)
		// defer bpDelete(ctx)
		pdps, err := bpClient.GetAllIbaPredefinedProbes(ctx)
		if err != nil {
			t.Fatal(err)
		}
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

			t.Logf("Got back Probe Id %s \n Now GET it.", probeId)

			p, err := bpClient.GetIbaProbe(ctx, probeId)

			t.Logf("Label %s", p.Label)
			t.Logf("Description %s", p.Description)
			t.Log(p)
			t.Logf("Delete probe")
			for _, i := range p.Stages {
				t.Logf("Stage name %s", i["name"])
			}
			err = bpClient.DeleteIbaProbe(ctx, probeId)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Delete Probe again, this should fail")
			err = bpClient.DeleteIbaProbe(ctx, probeId)
			if err == nil {
				t.Fatal("Probe Deletion should have failed")
			} else {
				t.Log(err)
			}
		}
	}
}
