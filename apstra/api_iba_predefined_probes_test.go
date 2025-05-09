// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra //

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIbaPredefinedProbes(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			log.Printf("testing Predefined Probes against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())

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
						t.Log(err)
						t.Logf("%s was expected to fail", p.Name)
						continue
					}
				}
				t.Logf("Got back Probe Id %s \n Now Make a Widget with it.", probeId)
			}
		})
	}
}
