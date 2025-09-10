// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestListAndGetAllLogicalDevices(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())

	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			ids, err := client.Client.ListLogicalDeviceIds(ctx)
			require.NoError(t, err)
			require.NotZero(t, len(ids))

			for _, i := range testutils.Range(len(ids)) {
				id := ids[i]
				t.Run(fmt.Sprintf("GET_%s", id), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(t, ctx)

					ld, err := client.Client.GetLogicalDevice(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, ld.Id)
				})
			}
		})
	}
}

func TestCreateGetUpdateDeleteLogicalDevice(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())

	clients := testclient.GetTestClients(t, ctx)

	indexingTypes := []string{
		apstra.PortIndexingVerticalFirst,
		apstra.PortIndexingHorizontalFirst,
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			deviceConfigs := make([]apstra.LogicalDeviceData, len(indexingTypes))
			for i, indexing := range indexingTypes {
				deviceConfigs[i] = apstra.LogicalDeviceData{
					DisplayName: testutils.RandString(6, "hex"),
					Panels: []apstra.LogicalDevicePanel{
						{
							PanelLayout: apstra.LogicalDevicePanelLayout{
								RowCount:    2,
								ColumnCount: 2,
							},
							PortIndexing: apstra.LogicalDevicePortIndexing{
								Order:      indexing,
								StartIndex: 0,
								Schema:     "absolute",
							},
							PortGroups: []apstra.LogicalDevicePortGroup{
								{
									Count: 4,
									Speed: "10G",
									Roles: apstra.LogicalDevicePortRoles{enum.PortRoleUnused},
								},
							},
						},
					},
				}
			}

			ids := make([]apstra.ObjectId, len(deviceConfigs))
			var err error
			for i, devCfg := range deviceConfigs {
				ids[i], err = client.Client.CreateLogicalDevice(ctx, &devCfg)
				require.NoError(t, err)

				d, err := client.Client.GetLogicalDevice(ctx, ids[i])
				require.NoError(t, err)

				log.Println(d.Id)
				devCfg.Panels[0].PortIndexing.StartIndex = 1
				require.NoError(t, client.Client.UpdateLogicalDevice(ctx, d.Id, &devCfg))

				if i > 0 {
					previous, err := client.Client.GetLogicalDevice(ctx, ids[i-1])
					require.NoError(t, err)

					previous.Data.DisplayName = d.Data.DisplayName

					require.NoError(t, client.Client.UpdateLogicalDevice(ctx, ids[i], previous.Data))
				}

				_, err = client.Client.GetLogicalDevice(ctx, ids[i])
				require.NoError(t, err)

			}

			for _, id := range ids {
				t.Run(fmt.Sprintf("DELETE_logical_device_%s", id), func(t *testing.T) {
					t.Parallel()
					require.NoError(t, client.Client.DeleteLogicalDevice(ctx, id))
				})
			}
		})
	}
}

func TestGetLogicalDeviceByName(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())

	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			ldIDs, err := client.Client.ListLogicalDeviceIds(ctx)
			require.NoError(t, err)

			for _, i := range testutils.SampleIndexes(t, len(ldIDs)) {
				testLD, err := client.Client.GetLogicalDevice(ctx, ldIDs[i])
				require.NoError(t, err)

				t.Run(fmt.Sprintf("GET_LD_%s", testLD.Data.DisplayName), func(t *testing.T) {
					t.Parallel()

					resultLD, err := client.Client.GetLogicalDeviceByName(ctx, testLD.Data.DisplayName)
					if err != nil {
						var ace apstra.ClientErr
						require.ErrorAs(t, err, ace)
						require.Equal(t, apstra.ErrMultipleMatch, ace.Type())
					} else {
						require.Equal(t, testLD.Id, resultLD.Id)
						require.Equal(t, testLD.Data.DisplayName, resultLD.Data.DisplayName)
					}
				})
			}
		})
	}
}
