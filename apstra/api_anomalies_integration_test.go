// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"log"
	"testing"

	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetBlueprintAnomalies(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)

	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bpClient := dctestobj.TestBlueprintB(t, ctx, client.Client)

			anomalies, err := client.Client.GetBlueprintAnomalies(ctx, bpClient.Id())
			require.NoError(t, err)

			log.Printf("%d blueprint anomalies retrieved", len(anomalies))
		})
	}
}

func TestGetBlueprintNodeAnomalyCounts(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)

	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bpClient := dctestobj.TestBlueprintB(t, ctx, client.Client)

			anomalies, err := client.Client.GetBlueprintNodeAnomalyCounts(ctx, bpClient.Id())
			require.NoError(t, err)

			log.Printf("%d node anomaly counts retrieved", len(anomalies))
		})
	}
}

func TestGetBlueprintServiceAnomalyCounts(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)

	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bpClient := dctestobj.TestBlueprintB(t, ctx, client.Client)

			anomalies, err := client.Client.GetBlueprintServiceAnomalyCounts(ctx, bpClient.Id())
			require.NoError(t, err)

			log.Printf("%d service anomaly counts retrieved", len(anomalies))
		})
	}
}
