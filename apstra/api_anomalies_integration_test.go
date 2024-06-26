//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestGetAnomalies(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getAnomalies() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())

		anomalies, err := client.client.GetAnomalies(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("%d anomalies retrieved", len(anomalies))
	}
}

func TestGetBlueprintAnomalies(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintB(ctx, t, client.client)

		log.Printf("testing GetBlueprintAnomalies() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		anomalies, err := client.client.GetBlueprintAnomalies(ctx, bpClient.Id())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("%d blueprint anomalies retrieved", len(anomalies))
	}
}

func TestGetBlueprintNodeAnomalyCounts(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintB(ctx, t, client.client)

		log.Printf("testing GetBlueprintAnomalies() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		anomalies, err := client.client.GetBlueprintNodeAnomalyCounts(ctx, bpClient.Id())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("%d node anomaly counts retrieved", len(anomalies))
	}
}

func TestGetBlueprintServiceAnomalyCounts(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintB(ctx, t, client.client)

		log.Printf("testing GetBlueprintAnomalies() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		anomalies, err := client.client.GetBlueprintServiceAnomalyCounts(ctx, bpClient.Id())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("%d service anomaly counts retrieved", len(anomalies))
	}
}
