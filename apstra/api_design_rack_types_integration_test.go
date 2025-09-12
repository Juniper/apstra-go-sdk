// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestListGetOneRackType(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			rtIds, err := client.Client.ListRackTypeIds(ctx)
			require.NoError(t, err)
			require.NotZero(t, len(rtIds))

			id := rtIds[0]

			rt, err := client.Client.GetRackType(ctx, id)
			require.NoError(t, err)

			log.Println(rt.Id)
		})
	}
}

func TestListGetAllGetRackType(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			rackTypeIds, err := client.Client.ListRackTypeIds(ctx)
			require.NoError(t, err)
			require.NotZero(t, len(rackTypeIds))

			rackTypes, err := client.Client.GetAllRackTypes(ctx)
			require.NoError(t, err)
			require.NotZero(t, len(rackTypes))

			require.Equal(t, len(rackTypeIds), len(rackTypes))

			for _, i := range testutils.SampleIndexes(t, len(rackTypeIds)) {
				id := rackTypeIds[i]

				rt, err := client.Client.GetRackType(ctx, id)
				require.NoError(t, err)

				require.Contains(t, rackTypeIds, rt.Id)

				log.Println(rt.Id)
			}
		})
	}
}

func TestCreateGetRackDeleteRackType(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	leafLabel := "ll-" + testutils.RandString(10, "hex")

	testCases := map[string]apstra.RackTypeRequest{
		"leaf_only_no_tags": {
			DisplayName:              "rdn " + testutils.RandString(5, "hex"),
			FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
			LeafSwitches: []apstra.RackElementLeafSwitchRequest{
				{
					Label:             leafLabel,
					LogicalDeviceId:   "AOS-48x10_6x40-leaf_spine",
					LinkPerSpineCount: 2,
					LinkPerSpineSpeed: "10G",
				},
			},
		},
		"leaf_generic_with_tags": {
			DisplayName:              "rdn " + testutils.RandString(5, "hex"),
			FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
			LeafSwitches: []apstra.RackElementLeafSwitchRequest{
				{
					Label:             leafLabel,
					LogicalDeviceId:   "AOS-48x10_6x40-leaf_spine",
					LinkPerSpineCount: 2,
					LinkPerSpineSpeed: "10G",
					Tags:              []apstra.ObjectId{"hypervisor", "bare_metal"},
				},
			},
			GenericSystems: []apstra.RackElementGenericSystemRequest{
				{
					Count: 5,
					Label: "some generic system",
					Links: []apstra.RackLinkRequest{
						{
							Label:              "foo",
							LinkPerSwitchCount: 1,
							LinkSpeed:          "10G",
							TargetSwitchLabel:  leafLabel,
							AttachmentType:     apstra.RackLinkAttachmentTypeSingle,
							LagMode:            apstra.RackLinkLagModeNone,
							Tags:               []apstra.ObjectId{"firewall"},
						},
					},
					LogicalDeviceId: "AOS-1x10-1",
					Tags:            []apstra.ObjectId{"firewall"},
				},
			},
			AccessSwitches: nil,
		},
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					ctx := testutils.ContextWithTestID(ctx, t)

					id, err := client.Client.CreateRackType(ctx, &tCase)
					require.NoError(t, err)

					rt, err := client.Client.GetRackType(ctx, id)
					require.NoError(t, err)

					require.Equal(t, rt.Id, id)
					require.NotNil(t, rt.Data)
					compare.RackType(t, tCase, *rt.Data)

					err = client.Client.DeleteRackType(ctx, id)
					require.NoError(t, err)
				})
			}
		})
	}
}
