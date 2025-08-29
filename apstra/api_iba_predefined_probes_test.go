// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"encoding/json"
	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIbaPredefinedProbes(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

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

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			bpClient := dctestobj.TestBlueprintA(t, ctx, client.Client)
			pdps, err := bpClient.GetAllIbaPredefinedProbes(ctx)
			require.NoError(t, err)

			t.Logf("Try an obviously fake name : %s", "FAKE")
			_, err = bpClient.GetIbaPredefinedProbeByName(ctx, "FAKE")
			require.Error(t, err)
			t.Log(err)

			for _, p := range pdps {
				t.Logf("Get Predefined Probe By Name %s", p.Name)
				_, err := bpClient.GetIbaPredefinedProbeByName(ctx, p.Name)
				require.NoError(t, err)

				t.Log(p.Description)
				t.Logf("%s", p.Schema)
				t.Logf("Instantiating Probe %s", p.Name)

				probeId, err := bpClient.InstantiateIbaPredefinedProbe(ctx, &apstra.IbaPredefinedProbeRequest{
					Name: p.Name,
					Data: json.RawMessage(`{"label":"` + p.Name + `"}`),
				})
				if expectedToFail[p.Name] {
					t.Log(err)
					t.Logf("%s was expected to fail", p.Name)
					continue
				} else {
					require.NoError(t, err)
				}

				t.Logf("Got back Probe Id %s \n Now Make a Widget with it.", probeId)
			}
		})
	}
}
