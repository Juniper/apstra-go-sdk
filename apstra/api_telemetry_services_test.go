//go:build integration
// +build integration

package apstra

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestGetTelemetryServicesDeviceMapping(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing GetTelemetryServicesDeviceMapping() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		result, err := client.client.GetTelemetryServicesDeviceMapping(ctx)
		require.NoError(t, err)

		buf := bytes.NewBuffer([]byte{})
		err = pp(result, buf)
		require.NoError(t, err)
	}
}
