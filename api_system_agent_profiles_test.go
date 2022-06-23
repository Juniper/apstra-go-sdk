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

	var cfgs []*SystemAgentProfileConfig
	for _, p := range []string{
		apstraAgentPlatformEOS,
		apstraAgentPlatformJunos,
		apstraAgentPlatformNXOS,
	} {
		cfgs = append(cfgs, &SystemAgentProfileConfig{
			Label:    randString(10, "hex"),
			Username: randString(10, "hex"),
			Password: randString(10, "hex"),
			Platform: p,
		})
	}

	var newIds []ObjectId
	for _, c := range cfgs {
		id, err := client.createSystemAgentProfile(context.TODO(), c)
		if err != nil {
			t.Fatal(err)
		}
		newIds = append(newIds, id)
	}

	apiIds, err := client.listSystemAgentProfileIds(context.TODO())
	if err != nil {
		t.Fatal(err)
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

	for _, id := range apiIds {
		err := client.deleteSystemAgentProfile(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
