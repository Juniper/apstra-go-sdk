// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBlueprintAnomalies(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bpClient := testBlueprintB(ctx, t, client.client)

			log.Printf("testing GetBlueprintAnomalies() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			anomalies, err := client.client.GetBlueprintAnomalies(ctx, bpClient.Id())
			require.NoError(t, err)

			log.Printf("%d blueprint anomalies retrieved", len(anomalies))
		})
	}
}

func TestGetBlueprintNodeAnomalyCounts(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bpClient := testBlueprintB(ctx, t, client.client)

			log.Printf("testing GetBlueprintAnomalies() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			anomalies, err := client.client.GetBlueprintNodeAnomalyCounts(ctx, bpClient.Id())
			require.NoError(t, err)

			log.Printf("%d node anomaly counts retrieved", len(anomalies))
		})
	}
}

func TestGetBlueprintServiceAnomalyCounts(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bpClient := testBlueprintB(ctx, t, client.client)

			log.Printf("testing GetBlueprintAnomalies() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			anomalies, err := client.client.GetBlueprintServiceAnomalyCounts(ctx, bpClient.Id())
			require.NoError(t, err)

			log.Printf("%d service anomaly counts retrieved", len(anomalies))
		})
	}
}
