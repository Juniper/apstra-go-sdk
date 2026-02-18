// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetInterfaceMapDigest(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		allImd, err := client.Client.GetInterfaceMapDigests(context.Background())
		require.NoError(t, err)

		randId := allImd[rand.Intn(len(allImd))].Id

		imd, err := client.Client.GetInterfaceMapDigest(context.Background(), randId)
		require.NoError(t, err)

		log.Printf("%s: %s -> %s", imd.Label, imd.LogicalDevice.Label, imd.DeviceProfile.Label)
	}
}

func TestGetInterfaceMapDigestsByLogicalDevice(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		ldIDs, err := client.Client.ListLogicalDeviceIds(ctx)
		require.NoError(t, err)

		randId := ldIDs[rand.Intn(len(ldIDs))]
		imds, err := client.Client.GetInterfaceMapDigestsByLogicalDevice(ctx, randId)
		require.NoError(t, err)

		for _, imd := range imds {
			log.Printf("%s: %s -> %s", imd.Label, imd.LogicalDevice.Label, imd.DeviceProfile.Label)
		}
	}
}

func TestGetInterfaceMapDigestsByDeviceProfile(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			longCtx, cf := context.WithTimeout(ctx, time.Minute)
			dpIDs, err := client.Client.GetDeviceProfiles(longCtx)
			cf()
			require.NoError(t, err)

			randId := dpIDs[rand.Intn(len(dpIDs))].ID()
			require.NotNil(t, randId)

			t.Run(*randId, func(t *testing.T) {
				t.Parallel()
				ctx := testutils.ContextWithTestID(ctx, t)

				imds, err := client.Client.GetInterfaceMapDigestsByDeviceProfile(ctx, apstra.ObjectId(*randId))
				require.NoError(t, err)

				for _, imd := range imds {
					log.Printf("%s: %s -> %s", imd.Label, imd.LogicalDevice.Label, imd.DeviceProfile.Label)
				}
			})
		})
	}
}
