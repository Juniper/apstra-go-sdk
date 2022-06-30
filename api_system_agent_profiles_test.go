package goapstra

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"
)

func systemAgentProfilesClient1() (*Client, error) {
	return NewClient(&ClientCfg{
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	})
}

func TestCreateListGetDeleteSystemAgentProfile(t *testing.T) {
	client, err := systemAgentProfilesClient1()
	if err != nil {
		t.Fatal(err)
	}

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
		id, err := client.createAgentProfile(context.TODO(), c)
		if err != nil {
			t.Fatal(err)
		}
		newIds = append(newIds, id)

		sap, err := client.GetAgentProfileByLabel(context.TODO(), c.Label)
		if err != nil {
			t.Fatal(err)
		}
		if id != sap.Id {
			t.Fatalf("error fetching System Agent Profile by label - '%s' != '%s'", id, sap.Id)
		}
	}

	apiIds, err := client.listAgentProfileIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	allProfiles, err := client.GetAllAgentProfiles(context.TODO())
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
		err := client.deleteAgentProfile(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
