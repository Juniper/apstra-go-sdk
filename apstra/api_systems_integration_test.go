// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"log"
	"testing"

	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestListSystems(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			systems, err := client.Client.ListSystems(ctx)
			require.NoError(t, err)

			for _, system := range systems {
				log.Println(system)
			}

			systemInfos, err := client.Client.GetAllSystemsInfo(ctx)
			require.NoError(t, err)
			require.Equal(t, len(systems), len(systemInfos))
			for _, systemInfo := range systemInfos {
				require.Contains(t, systems, systemInfo.Id)
			}
		})
	}
}

func TestGetSystems(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			systems, err := client.Client.ListSystems(ctx)
			require.NoError(t, err)

			for _, s := range systems {
				system, err := client.Client.GetSystemInfo(ctx, s)
				require.NoError(t, err)

				log.Println(system.Facts.HwModel)
			}
		})
	}
}
