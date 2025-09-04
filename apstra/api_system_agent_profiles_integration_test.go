// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestCreateListGetDeleteSystemAgentProfile(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			var cfgs []*apstra.AgentProfileConfig
			for _, p := range []string{"eos", "junos", "nxos"} {
				platform := p
				cfgs = append(cfgs, &apstra.AgentProfileConfig{
					Label:    testutils.RandString(10, "hex"),
					Username: testutils.ToPtr(testutils.RandString(10, "hex")),
					Password: testutils.ToPtr(testutils.RandString(10, "hex")),
					Platform: &platform,
					Packages: map[string]string{
						testutils.RandString(10, "hex"): testutils.RandString(10, "hex"),
						testutils.RandString(10, "hex"): testutils.RandString(10, "hex"),
					},
					OpenOptions: map[string]string{
						testutils.RandString(10, "hex"): testutils.RandString(10, "hex"),
						testutils.RandString(10, "hex"): testutils.RandString(10, "hex"),
					},
				})
			}

			var newIds []apstra.ObjectId
			for _, c := range cfgs {
				id, err := client.Client.CreateAgentProfile(ctx, c)
				if err != nil {
					t.Fatal(err)
				}
				newIds = append(newIds, id)

				sap, err := client.Client.GetAgentProfileByLabel(ctx, c.Label)
				if err != nil {
					t.Fatal(err)
				}
				if id != sap.Id {
					t.Fatalf("error fetching System Agent Profile by label - '%s' != '%s'", id, sap.Id)
				}
			}

			apiIds, err := client.Client.ListAgentProfileIds(ctx)
			if err != nil {
				t.Fatal(err)
			}

			allProfiles, err := client.Client.GetAllAgentProfiles(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if len(allProfiles) != len(apiIds) {
				t.Fatalf("found %d profiles and %d profile IDs", len(allProfiles), len(apiIds))
			}

			apiIdsMap := make(map[apstra.ObjectId]struct{}, len(apiIds))
			for _, id := range apiIds {
				apiIdsMap[id] = struct{}{}
			}

			for _, id := range newIds {
				require.Contains(t, apiIdsMap, id)
			}

			for _, id := range newIds {
				err := client.Client.DeleteAgentProfile(ctx, id)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestClient_UpdateAgentProfile_ClearStringFields(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			id, err := client.Client.CreateAgentProfile(ctx, &apstra.AgentProfileConfig{
				Label:    testutils.RandString(5, "hex"),
				Username: testutils.ToPtr(testutils.RandString(5, "hex")),
				Password: testutils.ToPtr(testutils.RandString(5, "hex")),
				Platform: testutils.ToPtr("junos"),
			})
			require.NoError(t, err)

			err = client.Client.UpdateAgentProfile(ctx, id, &apstra.AgentProfileConfig{
				Username: testutils.ToPtr(""),
				Password: testutils.ToPtr(""),
				Platform: testutils.ToPtr(""),
			})
			require.NoError(t, err)

			ap, err := client.Client.GetAgentProfile(ctx, id)
			require.NoError(t, err)
			require.Equal(t, false, ap.HasUsername)
			require.Equal(t, false, ap.HasPassword)
			require.Equal(t, "", ap.Platform)

			err = client.Client.DeleteAgentProfile(ctx, id)
			require.NoError(t, err)
		})
	}
}

func TestClient_UpdateAgentProfile(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			agents, err := client.Client.GetAllSystemAgents(ctx)
			require.NoError(t, err)

			profiles, err := client.Client.GetAllAgentProfiles(ctx)
			require.NoError(t, err)
			profileMap := make(map[apstra.ObjectId]apstra.AgentProfile, len(profiles))
			for _, profile := range profiles {
				profileMap[profile.Id] = profile
			}

			for i := len(agents) - 1; i >= 0; i-- { // loop backward through agents
				// remove agents without profile
				if agents[i].Config.Profile == "" {
					agents[i] = agents[len(agents)-1] // copy last to index i
					agents = agents[:len(agents)-1]   // delete last item
					continue
				}

				// remove agents with platform configured
				if agents[i].Config.Platform != "" {
					agents[i] = agents[len(agents)-1] // copy last to index i
					agents = agents[:len(agents)-1]   // delete last item
					continue
				}

				_, ok := profileMap[agents[i].Config.Profile]
				if !ok {
					t.Fatal(fmt.Errorf("agent %s claims profile %s, which is not a valid profile", agents[i].Id, agents[i].Config.Profile))
				}

				// remove agents with un-acked systems
				systemInfo, err := client.Client.GetSystemInfo(ctx, agents[i].Status.SystemId)
				require.NoError(t, err)
				if !systemInfo.Status.IsAcknowledged {
					agents[i] = agents[len(agents)-1] // copy last to index i
					agents = agents[:len(agents)-1]   // delete last item
					continue
				}
			}

			if len(agents) == 0 {
				t.Skip("skipping because no system agents rely on agent profile to determine platform type")
			}

			// At this point, agents is full of system agents which rely on their associated agent
			// profile for platform info. We'll attempt to modify the agent profile to elicit an error.

			var ace apstra.ClientErr
			for _, agent := range agents {
				profile := profileMap[agent.Config.Profile]
				err = client.Client.UpdateAgentProfile(ctx, profile.Id, &apstra.AgentProfileConfig{
					Label:    profile.Label,
					Platform: testutils.ToPtr(""),
				})
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, apstra.ErrAgentProfilePlatformRequired, ace.Type())
			}
		})
	}
}
