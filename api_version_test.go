//go:build integration
// +build integration

package goapstra

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

func TestGetVersion(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		ver, err := client.client.getVersion(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		result, err := json.Marshal(ver)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("%s %s", client.client.baseUrl.String(), string(result))
	}
}
