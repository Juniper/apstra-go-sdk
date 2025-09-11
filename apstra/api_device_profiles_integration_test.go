// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestListAndGetSampleDeviceProfiles(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			ids, err := client.Client.ListDeviceProfileIds(ctx)
			require.NoError(t, err)
			require.NotZero(t, len(ids))

			for _, i := range testutils.SampleIndexes(t, len(ids), 20) {
				id := ids[i]
				t.Run(fmt.Sprintf("GET_Device_Profile_ID_%s", id), func(t *testing.T) {
					ctx := testutils.ContextWithTestID(ctx, t)

					dp, err := client.Client.GetDeviceProfile(ctx, id)
					require.NoError(t, err)
					log.Printf("device profile id '%s' label '%s'\n", id, dp.Data.Label)
				})
			}

			profiles, err := client.Client.GetAllDeviceProfiles(ctx)
			require.NoError(t, err)

			log.Printf("list found %d, GetAll found %d", len(ids), len(profiles))
			require.Equal(t, len(ids), len(profiles))

			for _, i := range testutils.SampleIndexes(t, len(profiles), 5) {
				label := profiles[i].Data.Label
				t.Run(fmt.Sprintf("GET_DeviceProfile_Label_%s", label), func(t *testing.T) {
					// t.Parallel() // this seems to make things worse
					ctx := testutils.ContextWithTestID(ctx, t)

					dp, err := client.Client.GetDeviceProfileByName(ctx, label)
					require.NoError(t, err)
					require.Equal(t, label, dp.Data.Label)

					log.Printf("device profile id '%s' label '%s'\n", dp.Id, dp.Data.Label)
				})
			}
		})
	}
}

func TestGetDeviceProfile(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	desiredId := apstra.ObjectId("Cisco_3172PQ_NXOS")
	desiredLabel := "Cisco 3172PQ"

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			dp, err := client.Client.GetDeviceProfile(ctx, desiredId)
			require.NoError(t, err)
			require.Equal(t, dp.Data.Label, desiredLabel)
		})
	}
}

func TestGetDeviceProfileByName(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	desiredLabel := "Cisco 3172PQ"
	desiredId := apstra.ObjectId("Cisco_3172PQ_NXOS")

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			dp, err := client.Client.GetDeviceProfileByName(ctx, desiredLabel)
			require.NoError(t, err)
			require.Equal(t, desiredId, dp.Id)
		})
	}
}

func TestGetTransformCandidates(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	dpId := apstra.ObjectId("Juniper_QFX5120-48T_Junos")
	intfName := "et-0/0/48"
	intfSpeed := apstra.LogicalDevicePortSpeed("40G")

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			dp, err := client.Client.GetDeviceProfile(ctx, dpId)
			require.NoError(t, err)

			portInfo, err := dp.Data.PortByInterfaceName(intfName)
			require.NoError(t, err)

			candidates := portInfo.TransformationCandidates(intfName, intfSpeed)
			for k, v := range candidates {
				dump, err := json.MarshalIndent(&v, "", "  ")
				require.NoError(t, err)

				log.Printf("port %d (%s) transformations: %s", k, intfName, string(dump))
			}
		})
	}
}
