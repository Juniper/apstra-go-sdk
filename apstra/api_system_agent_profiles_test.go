// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateListGetDeleteSystemAgentProfile(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {

		var cfgs []*AgentProfileConfig
		for _, p := range []string{
			apstraAgentPlatformEOS,
			apstraAgentPlatformJunos,
			apstraAgentPlatformNXOS,
		} {
			platform := p
			cfgs = append(cfgs, &AgentProfileConfig{
				Label:    randString(10, "hex"),
				Username: toPtr(randString(10, "hex")),
				Password: toPtr(randString(10, "hex")),
				Platform: &platform,
				Packages: map[string]string{
					randString(10, "hex"): randString(10, "hex"),
					randString(10, "hex"): randString(10, "hex"),
				},
				OpenOptions: map[string]string{
					randString(10, "hex"): randString(10, "hex"),
					randString(10, "hex"): randString(10, "hex"),
				},
			})
		}

		var newIds []ObjectId
		for _, c := range cfgs {
			log.Printf("testing createAgentProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			id, err := client.client.createAgentProfile(context.TODO(), c)
			if err != nil {
				t.Fatal(err)
			}
			newIds = append(newIds, id)

			log.Printf("testing GetAgentProfileByLabel() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			sap, err := client.client.GetAgentProfileByLabel(context.TODO(), c.Label)
			if err != nil {
				t.Fatal(err)
			}
			if id != sap.Id {
				t.Fatalf("error fetching System Agent Profile by label - '%s' != '%s'", id, sap.Id)
			}
		}

		log.Printf("testing listAgentProfileIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		apiIds, err := client.client.listAgentProfileIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetAllAgentProfiles() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		allProfiles, err := client.client.GetAllAgentProfiles(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(allProfiles) != len(apiIds) {
			t.Fatalf("found %d profiles and %d profile IDs", len(allProfiles), len(apiIds))
		}

		apiIdsMap := make(map[ObjectId]struct{})
		for _, id := range apiIds {
			apiIdsMap[id] = struct{}{}
		}

		for _, id := range newIds {
			if _, found := apiIdsMap[id]; !found {
				t.Fatal(fmt.Errorf("created id %s, but didn't find it in the list returned by the API", id))
			}
		}

		for _, id := range newIds {
			log.Printf("testing deleteAgentProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err := client.client.deleteAgentProfile(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestClient_UpdateAgentProfile_ClearStringFields(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing CreateAgentProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.CreateAgentProfile(ctx, &AgentProfileConfig{
			Label:    randString(5, "hex"),
			Username: toPtr(randString(5, "hex")),
			Password: toPtr(randString(5, "hex")),
			Platform: toPtr("junos"),
		})
		require.NoError(t, err)

		log.Printf("testing UpdateAgentProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.UpdateAgentProfile(ctx, id, &AgentProfileConfig{
			Username: toPtr(""),
			Password: toPtr(""),
			Platform: toPtr(""),
		})
		require.NoError(t, err)

		ap, err := client.client.GetAgentProfile(ctx, id)
		require.NoError(t, err)
		require.Equal(t, false, ap.HasUsername)
		require.Equal(t, false, ap.HasPassword)
		require.Equal(t, "", ap.Platform)

		log.Printf("testing DeleteAgentProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteAgentProfile(ctx, id)
		require.NoError(t, err)
	}
}

func TestClient_UpdateAgentProfile(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			t.Logf("testing GetAllSystemAgents() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			agents, err := client.client.GetAllSystemAgents(ctx)
			require.NoError(t, err)

			profiles, err := client.client.GetAllAgentProfiles(ctx)
			require.NoError(t, err)
			profileMap := make(map[ObjectId]AgentProfile, len(profiles))
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
				systemInfo, err := client.client.GetSystemInfo(ctx, agents[i].Status.SystemId)
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

			for _, agent := range agents {
				profile := profileMap[agent.Config.Profile]
				t.Logf("trying to provoke an error with UpdateAgentProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				err = client.client.UpdateAgentProfile(ctx, profile.Id, &AgentProfileConfig{
					Label:    profile.Label,
					Platform: toPtr(""),
				})
				require.Error(t, err)
				var ace ClientErr
				if !(errors.As(err, &ace) && ace.errType == ErrAgentProfilePlatformRequired) {
					t.Fatalf("error should have been type %d, err is %q", ErrAgentProfilePlatformRequired, err.Error())
				}
			}
		})
	}
}
