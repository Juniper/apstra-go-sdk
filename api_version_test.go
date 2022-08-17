package goapstra

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

func TestGetVersion(t *testing.T) {
	clients, err := getCloudlabsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		ver, err := client.getVersion(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		result, err := json.Marshal(ver)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("%s %s", client.baseUrl.String(), string(result))
	}
}
