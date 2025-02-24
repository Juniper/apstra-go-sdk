// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/stretchr/testify/require"
)

func TestListAndGetAllLogicalDevices(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	ctx = wrapCtxWithTestId(t, ctx)
	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()
			ctx = wrapCtxWithTestId(t, ctx)

			log.Printf("testing listLogicalDeviceIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			ids, err := client.client.listLogicalDeviceIds(ctx)
			require.NoError(t, err)
			require.NotEqual(t, len(ids), 0)

			for _, i := range samples(t, len(ids)) {
				id := ids[i]
				t.Run(fmt.Sprintf("GET_%s", id), func(t *testing.T) {
					t.Parallel()
					ctx = wrapCtxWithTestId(t, ctx)

					ld, err := client.client.GetLogicalDevice(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, ld.Id)
				})
			}
		})
	}
}

func TestCreateGetUpdateDeleteLogicalDevice(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	indexingTypes := []string{
		PortIndexingVerticalFirst,
		PortIndexingHorizontalFirst,
	}

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			deviceConfigs := make([]LogicalDeviceData, len(indexingTypes))
			for i, indexing := range indexingTypes {
				deviceConfigs[i] = LogicalDeviceData{
					DisplayName: randString(6, "hex"),
					Panels: []LogicalDevicePanel{
						{
							PanelLayout: LogicalDevicePanelLayout{
								RowCount:    2,
								ColumnCount: 2,
							},
							PortIndexing: LogicalDevicePortIndexing{
								Order:      indexing,
								StartIndex: 0,
								Schema:     "absolute",
							},
							PortGroups: []LogicalDevicePortGroup{
								{
									Count: 4,
									Speed: "10G",
									Roles: LogicalDevicePortRoles{enum.PortRoleUnused},
								},
							},
						},
					},
				}
			}

			ids := make([]ObjectId, len(deviceConfigs))
			for i, devCfg := range deviceConfigs {
				log.Printf("testing createLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				ids[i], err = client.client.createLogicalDevice(ctx, devCfg.raw())
				require.NoError(t, err)

				log.Printf("testing GetLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				d, err := client.client.GetLogicalDevice(ctx, ids[i])
				require.NoError(t, err)

				log.Println(d.Id)
				devCfg.Panels[0].PortIndexing.StartIndex = 1
				log.Printf("testing updateLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				require.NoError(t, client.client.updateLogicalDevice(ctx, d.Id, devCfg.raw()))

				if i > 0 {
					log.Printf("testing GetLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					previous, err := client.client.GetLogicalDevice(ctx, ids[i-1])
					require.NoError(t, err)

					previous.Data.DisplayName = d.Data.DisplayName

					log.Printf("testing updateLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					require.NoError(t, client.client.updateLogicalDevice(ctx, ids[i], previous.Data.raw()))
				}

				log.Printf("testing GetLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				_, err = client.client.GetLogicalDevice(ctx, ids[i])
				require.NoError(t, err)

			}

			for _, id := range ids {
				t.Run(fmt.Sprintf("DELETE_logical_device_%s", id), func(t *testing.T) {
					t.Parallel()
					require.NoError(t, client.client.deleteLogicalDevice(ctx, id))
				})
			}
		})
	}
}

func TestGetLogicalDeviceByName(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			log.Printf("testing deleteLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			logicalDevices, err := client.client.getAllLogicalDevices(context.TODO())
			if err != nil {
				t.Fatal(err)
			}

			for _, i := range samples(t, len(logicalDevices)) {
				test := logicalDevices[i]
				t.Run(fmt.Sprintf("GET_LD_%s", test.Data.DisplayName), func(t *testing.T) {
					t.Parallel()

					logicalDevice, err := client.client.GetLogicalDeviceByName(context.TODO(), test.Data.DisplayName)
					if err != nil {
						var ace ClientErr
						require.ErrorAs(t, err, ace)
						require.Equal(t, ErrMultipleMatch, ace.Type())
					} else {
						require.Equal(t, test.Id, logicalDevice.Id)
						require.Equal(t, test.Data.DisplayName, logicalDevice.Data.DisplayName)
					}
				})
			}
		})
	}
}
