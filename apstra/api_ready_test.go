//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestClientReady(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing ready() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.ready(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestClientWaitUntilReady(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		log.Printf("testing ready() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.ready(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}
}
