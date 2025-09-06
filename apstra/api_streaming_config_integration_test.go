// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"testing"

	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestClient_GetAllStreamingConfigs(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			ids, err := client.Client.GetAllStreamingConfigIds(ctx)
			require.NoError(t, err)

			for _, id := range ids {
				streamingConfig, err := client.Client.GetStreamingConfig(ctx, id)
				require.NoError(t, err)
				require.Equal(t, id, streamingConfig.Id)
			}
		})
	}
}
