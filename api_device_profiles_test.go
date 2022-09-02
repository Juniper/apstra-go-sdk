//go:build integration
// +build integration

package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestListAndGetAllDeviceProfiles(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listDeviceProfileIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ids, err := client.client.listDeviceProfileIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if len(ids) <= 0 {
			t.Fatalf("only got %d ids", len(ids))
		}
		for _, i := range samples(len(ids)) {
			id := ids[i]
			log.Printf("testing getDeviceProfile(%s) against %s %s (%s)", id, client.clientType, clientName, client.client.ApiVersion())
			dp, err := client.client.getDeviceProfile(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("device profile id '%s' label '%s'\n", id, dp.Label)
		}
		log.Printf("testing getAllDeviceProfiles() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		profiles, err := client.client.getAllDeviceProfiles(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("list found %d, getAll found %d", len(ids), len(profiles))
	}
}
