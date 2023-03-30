//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestClient_GetAllStreamingConfigs(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing GetAllStreamingConfigIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ids, err := client.client.GetAllStreamingConfigIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		for _, id := range ids {
			log.Printf("testing GetStreamingConfig() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			streamingConfig, err := client.client.GetStreamingConfig(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("streaming config: %s, %s, %s:%d", streamingConfig.Protocol, streamingConfig.StreamingType, streamingConfig.Hostname, streamingConfig.Port)
		}
	}
}
