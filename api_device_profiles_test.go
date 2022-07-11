package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestListAndGetAllDeviceProfiles(t *testing.T) {
	DebugLevel = 0
	clients, _, err := getTestClientsAndMockAPIs()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(len(clients))

	for clientName, client := range clients {
		if clientName == "mock" {
			continue // todo have I given up on mock testing?
		}
		ids, err := client.listDeviceProfileIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if len(ids) <= 0 {
			t.Fatalf("only got %d ids from %s client", len(ids), clientName)
		}
		for _, id := range ids {
			dp, err := client.getDeviceProfile(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("device profile id '%s' label '%s'\n", id, dp.Label)
		}
		profiles, err := client.getAllDeviceProfiles(context.TODO())
		log.Printf("list found %d, getAll found %d", len(ids), len(profiles))
	}
}
