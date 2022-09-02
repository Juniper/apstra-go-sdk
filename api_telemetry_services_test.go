//go:build integration
// +build integration

package goapstra

import (
	"bytes"
	"context"
	"log"
	"testing"
)

func TestGetTelemetryServicesDeviceMapping(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing GetTelemetryServicesDeviceMapping() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		result, err := client.client.GetTelemetryServicesDeviceMapping(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		buf := bytes.NewBuffer([]byte{})
		err = pp(result, buf)
		if err != nil {
			t.Fatal(err)
		}
	}
}
