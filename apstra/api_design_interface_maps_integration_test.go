// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/stretchr/testify/require"
)

func TestInterfaceSettingParam(t *testing.T) {
	expected := `{\"global\":{\"breakout\":false,\"fpc\":0,\"pic\":0,\"port\":0,\"speed\":\"100g\"},\"interface\":{\"speed\":\"\"}}`
	test := InterfaceSettingParam{
		Global: struct {
			Breakout bool   `json:"breakout"`
			Fpc      int    `json:"fpc"`
			Pic      int    `json:"pic"`
			Port     int    `json:"port"`
			Speed    string `json:"speed"`
		}{
			Breakout: false,
			Fpc:      0,
			Pic:      0,
			Port:     0,
			Speed:    "100g",
		},
		Interface: struct {
			Speed string `json:"speed"`
		}{},
	}
	result := test.String()
	if result != expected {
		t.Fatalf("expected '%s', got '%s'", expected, result)
	}
}

func TestListGetAllInterfaceMaps(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}
	for clientName, client := range clients {
		log.Printf("testing listAllInterfaceMapIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		iMapIds, err := client.client.listAllInterfaceMapIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(iMapIds) == 0 {
			t.Fatal("we should have gotten some interface maps here")
		}

		log.Println("all interface maps IDs: ", iMapIds)

		log.Printf("testing getInterfaceMap() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		iMap, err := client.client.getInterfaceMap(context.TODO(), iMapIds[rand.Intn(len(iMapIds))])
		if err != nil {
			t.Fatal(err)
		}
		log.Println("random interface map: ", iMap)
	}
}

func TestCreateInterfaceMap(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {

		ldId := ObjectId("AOS-1x1-1")
		dpId := ObjectId("Generic_Server_1RU_1x1G")
		label := "label-" + randString(10, "hex")

		newMapInfo := InterfaceMapData{
			LogicalDeviceId: ldId,
			DeviceProfileId: dpId,
			Label:           label,
			Interfaces: []InterfaceMapInterface{
				{
					Name:  "eth0",
					Roles: LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess},
					Mapping: InterfaceMapMapping{
						DPPortId:      1,
						DPTransformId: 1,
						DPInterfaceId: 1,
						LDPanel:       1,
						LDPort:        1,
					},
					ActiveState: true,
					Position:    1,
					Speed:       "1G",
					Setting:     InterfaceMapInterfaceSetting{Param: ""},
				},
			},
		}

		log.Printf("testing createInterfaceMap() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		mapId, err := client.client.createInterfaceMap(context.TODO(), &newMapInfo)
		require.NoError(t, err)

		log.Printf("testing getInterfaceMap() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		asCreated, err := client.client.GetInterfaceMap(context.TODO(), mapId)
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

		log.Printf("testing deleteInterfaceMap() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteInterfaceMap(context.TODO(), mapId)
		require.NoError(t, err)
	}
}

func TestGetInterfaceMapByName(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	desired := "Juniper_QFX5120-32C_Junos__AOS-32x100-1"

	for clientName, client := range clients {
		log.Printf("testing getInterfaceMapByName(%s) against %s %s (%s)", desired, client.clientType, clientName, client.client.ApiVersion())
		interfaceMap, err := client.client.GetInterfaceMapByName(context.Background(), desired)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("%s <---> %s", interfaceMap.Data.LogicalDeviceId, interfaceMap.Data.DeviceProfileId)
	}
}
