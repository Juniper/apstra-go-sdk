// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"log"
	"math/rand"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestListGetAllInterfaceMaps(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			iMapIds, err := client.Client.ListAllInterfaceMapIds(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, iMapIds)

			log.Println("all interface maps IDs: ", iMapIds)

			iMap, err := client.Client.GetInterfaceMap(ctx, iMapIds[rand.Intn(len(iMapIds))])
			require.NoError(t, err)
			log.Println("random interface map: ", iMap)
		})
	}
}

func TestCreateInterfaceMap(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	newMapInfo := apstra.InterfaceMapData{
		LogicalDeviceId: "AOS-1x1-1",
		DeviceProfileId: "Generic_Server_1RU_1x1G",
		Label:           "label-" + testutils.RandString(10, "hex"),
		Interfaces: []apstra.InterfaceMapInterface{
			{
				Name:  "eth0",
				Roles: apstra.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess},
				Mapping: apstra.InterfaceMapMapping{
					DPPortId:      1,
					DPTransformId: 1,
					DPInterfaceId: 1,
					LDPanel:       1,
					LDPort:        1,
				},
				ActiveState: true,
				Position:    1,
				Speed:       "1G",
				Setting:     apstra.InterfaceMapInterfaceSetting{Param: ""},
			},
		},
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			mapId, err := client.Client.CreateInterfaceMap(ctx, &newMapInfo)
			require.NoError(t, err)

			asCreated, err := client.Client.GetInterfaceMap(ctx, mapId)
			require.NoError(t, err)
			require.NotNil(t, asCreated)
			require.Equal(t, mapId, asCreated.Id)
			require.Equal(t, newMapInfo.LogicalDeviceId, asCreated.Data.LogicalDeviceId)
			require.Equal(t, newMapInfo.DeviceProfileId, asCreated.Data.DeviceProfileId)
			require.Equal(t, newMapInfo.Label, asCreated.Data.Label)
			require.Equal(t, asCreated.Data.Interfaces, newMapInfo.Interfaces)

			for i := range asCreated.Data.Interfaces {
				require.Equal(t, newMapInfo.Interfaces[i].Name, asCreated.Data.Interfaces[i].Name)
				require.Equal(t, newMapInfo.Interfaces[i].Roles, asCreated.Data.Interfaces[i].Roles)
				require.Equal(t, newMapInfo.Interfaces[i].ActiveState, asCreated.Data.Interfaces[i].ActiveState)
				require.Equal(t, newMapInfo.Interfaces[i].Setting.Param, asCreated.Data.Interfaces[i].Setting.Param)
				require.Equal(t, newMapInfo.Interfaces[i].Position, asCreated.Data.Interfaces[i].Position)
				require.Equal(t, newMapInfo.Interfaces[i].Speed.BitsPerSecond(), asCreated.Data.Interfaces[i].Speed.BitsPerSecond())
				require.Equal(t, newMapInfo.Interfaces[i].Mapping.DPInterfaceId, asCreated.Data.Interfaces[i].Mapping.DPInterfaceId)
				require.Equal(t, newMapInfo.Interfaces[i].Mapping.DPPortId, asCreated.Data.Interfaces[i].Mapping.DPPortId)
				require.Equal(t, newMapInfo.Interfaces[i].Mapping.DPTransformId, asCreated.Data.Interfaces[i].Mapping.DPTransformId)
				require.Equal(t, newMapInfo.Interfaces[i].Mapping.LDPanel, asCreated.Data.Interfaces[i].Mapping.LDPanel)
				require.Equal(t, newMapInfo.Interfaces[i].Mapping.LDPort, asCreated.Data.Interfaces[i].Mapping.LDPort)
			}

			log.Println("new interface map: ", mapId)

			err = client.Client.DeleteInterfaceMap(ctx, mapId)
			require.NoError(t, err)
		})
	}
}

func TestGetInterfaceMapByName(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	desired := "Juniper_QFX5120-32C_Junos__AOS-32x100-1"

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			interfaceMap, err := client.Client.GetInterfaceMapByName(ctx, desired)
			require.NoError(t, err)

			log.Printf("%s <---> %s", interfaceMap.Data.LogicalDeviceId, interfaceMap.Data.DeviceProfileId)
		})
	}
}
