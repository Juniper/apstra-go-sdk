package goapstra

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestCreateListGetDeleteSystemAgentProfile(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {

		var cfgs []*AgentProfileConfig
		for _, p := range []string{
			apstraAgentPlatformEOS,
			apstraAgentPlatformJunos,
			apstraAgentPlatformNXOS,
		} {
			cfgs = append(cfgs, &AgentProfileConfig{
				Label:    randString(10, "hex"),
				Username: randString(10, "hex"),
				Password: randString(10, "hex"),
				Platform: p,
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
			log.Printf("testing createAgentProfile() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
			id, err := client.client.createAgentProfile(context.TODO(), c)
			if err != nil {
				t.Fatal(err)
			}
			newIds = append(newIds, id)

			log.Printf("testing GetAgentProfileByLabel() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
			sap, err := client.client.GetAgentProfileByLabel(context.TODO(), c.Label)
			if err != nil {
				t.Fatal(err)
			}
			if id != sap.Id {
				t.Fatalf("error fetching System Agent Profile by label - '%s' != '%s'", id, sap.Id)
			}
		}

		log.Printf("testing listAgentProfileIds() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		apiIds, err := client.client.listAgentProfileIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetAllAgentProfiles() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
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
			log.Printf("testing deleteAgentProfile() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
			err := client.client.deleteAgentProfile(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
