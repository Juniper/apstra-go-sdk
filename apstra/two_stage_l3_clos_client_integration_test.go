//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestNodeIdsByType(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		testBlueprintANodeCount := 4
		bpClient := testBlueprintA(ctx, t, client.client)

		if len(bpClient.nodeIdsByType) != 0 {
			t.Fatal("nodeIdsByType should be empty with a new blueprint client")
		}

		log.Printf("testing NodeIdsByType() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		systemIds, err := bpClient.NodeIdsByType(ctx, NodeTypeSystem)
		if err != nil {
			t.Fatal(err)
		}

		if len(systemIds) != testBlueprintANodeCount {
			t.Fatalf("expected %d nodes, got %d nodes", testBlueprintANodeCount, len(systemIds))
		}

		log.Printf("purging system node slice")
		bpClient.nodeIdsByType[NodeTypeSystem] = []ObjectId{}

		log.Printf("testing NodeIdsByType() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		systemIds, err = bpClient.NodeIdsByType(ctx, NodeTypeSystem)
		if err != nil {
			t.Fatal(err)
		}

		if len(systemIds) != 0 {
			t.Fatalf("expected 0 nodes, got %d nodes", len(systemIds))
		}

		log.Printf("testing RefreshNodeIdsByType() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		systemIds, err = bpClient.RefreshNodeIdsByType(ctx, NodeTypeSystem)
		if err != nil {
			t.Fatal(err)
		}

		if len(systemIds) != testBlueprintANodeCount {
			t.Fatalf("expected %d nodes, got %d nodes", testBlueprintANodeCount, len(systemIds))
		}
	}
}
