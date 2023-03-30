//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"
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
	rand.Seed(time.Now().UnixNano())
	clients, err := getTestClients(context.Background())
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
	clients, err := getTestClients(context.Background())
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
					Roles: LogicalDevicePortRoleLeaf | LogicalDevicePortRoleAccess,
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
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getInterfaceMap() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		asCreated, err := client.client.GetInterfaceMap(context.TODO(), mapId)
		if err != nil {
			t.Fatal(err)
		}

		if asCreated.Id != mapId {
			t.Fatalf("interface map id mismatch: '%s' vs. '%s'", asCreated.Id, mapId)
		}

		if asCreated.Data.LogicalDeviceId != newMapInfo.LogicalDeviceId {
			t.Fatalf("interface map logical device id mismatch: '%s' vs. '%s'", asCreated.Data.LogicalDeviceId, newMapInfo.LogicalDeviceId)
		}

		if asCreated.Data.DeviceProfileId != newMapInfo.DeviceProfileId {
			t.Fatalf("interface map device profile id mismatch: '%s' vs. '%s'", asCreated.Data.DeviceProfileId, newMapInfo.DeviceProfileId)
		}

		if asCreated.Data.Label != newMapInfo.Label {
			t.Fatalf("interface map label mismatch: '%s' vs. '%s'", asCreated.Data.Label, newMapInfo.Label)
		}

		if len(asCreated.Data.Interfaces) != len(newMapInfo.Interfaces) {
			t.Fatalf("interface map interface count mismatch: '%d' vs. '%d'", len(asCreated.Data.Interfaces), len(newMapInfo.Interfaces))
		}

		for i := 0; i < len(asCreated.Data.Interfaces); i++ {
			if asCreated.Data.Interfaces[i].Name != newMapInfo.Interfaces[i].Name {
				t.Fatalf("interface map interface [%d] name mistatch: '%s' vs. '%s'", i, asCreated.Data.Interfaces[i].Name, newMapInfo.Interfaces[i].Name)
			}
			if asCreated.Data.Interfaces[i].Roles != newMapInfo.Interfaces[i].Roles {
				t.Fatalf("interface map interface [%d] roles mistatch: '%s' vs. '%s'", i, asCreated.Data.Interfaces[i].Roles.Strings(), newMapInfo.Interfaces[i].Roles.Strings())
			}
			if asCreated.Data.Interfaces[i].ActiveState != newMapInfo.Interfaces[i].ActiveState {
				t.Fatalf("interface map interface [%d] state mistatch: '%s' vs. '%s'", i, asCreated.Data.Interfaces[i].ActiveState.raw(), newMapInfo.Interfaces[i].ActiveState.raw())
			}
			if asCreated.Data.Interfaces[i].Setting.Param != newMapInfo.Interfaces[i].Setting.Param {
				t.Fatalf("interface map interface [%d] setting param mistatch: '%s' vs. '%s'", i, asCreated.Data.Interfaces[i].Setting.Param, newMapInfo.Interfaces[i].Setting.Param)
			}
			if asCreated.Data.Interfaces[i].Position != newMapInfo.Interfaces[i].Position {
				t.Fatalf("interface map interface [%d] position mistatch: '%d' vs. '%d'", i, asCreated.Data.Interfaces[i].Position, newMapInfo.Interfaces[i].Position)
			}
			if asCreated.Data.Interfaces[i].Speed.BitsPerSecond() != newMapInfo.Interfaces[i].Speed.BitsPerSecond() {
				t.Fatalf("interface map interface [%d] speed mistatch: '%dbps' vs. '%dbps'", i, asCreated.Data.Interfaces[i].Speed.BitsPerSecond(), newMapInfo.Interfaces[i].Speed.BitsPerSecond())
			}
			if asCreated.Data.Interfaces[i].Mapping.DPInterfaceId != newMapInfo.Interfaces[i].Mapping.DPInterfaceId {
				t.Fatalf("interface map interface [%d] mapping device profile interface Id: '%d' vs. '%d'", i, asCreated.Data.Interfaces[i].Mapping.DPInterfaceId, newMapInfo.Interfaces[i].Mapping.DPInterfaceId)
			}
			if asCreated.Data.Interfaces[i].Mapping.DPPortId != newMapInfo.Interfaces[i].Mapping.DPPortId {
				t.Fatalf("interface map interface [%d] mapping device profile port Id: '%d' vs. '%d'", i, asCreated.Data.Interfaces[i].Mapping.DPPortId, newMapInfo.Interfaces[i].Mapping.DPPortId)
			}
			if asCreated.Data.Interfaces[i].Mapping.DPTransformId != newMapInfo.Interfaces[i].Mapping.DPTransformId {
				t.Fatalf("interface map interface [%d] mapping device profile transform Id: '%d' vs. '%d'", i, asCreated.Data.Interfaces[i].Mapping.DPTransformId, newMapInfo.Interfaces[i].Mapping.DPTransformId)
			}
			if asCreated.Data.Interfaces[i].Mapping.LDPanel != newMapInfo.Interfaces[i].Mapping.LDPanel {
				t.Fatalf("interface map interface [%d] mapping logical device panel Id: '%d' vs. '%d'", i, asCreated.Data.Interfaces[i].Mapping.LDPanel, newMapInfo.Interfaces[i].Mapping.LDPanel)
			}
			if asCreated.Data.Interfaces[i].Mapping.LDPort != newMapInfo.Interfaces[i].Mapping.LDPort {
				t.Fatalf("interface map interface [%d] mapping logical device port Id: '%d' vs. '%d'", i, asCreated.Data.Interfaces[i].Mapping.LDPort, newMapInfo.Interfaces[i].Mapping.LDPort)
			}
		}

		log.Println("new interface map: ", mapId)

		log.Printf("testing deleteInterfaceMap() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteInterfaceMap(context.TODO(), mapId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetInterfaceMapByName(t *testing.T) {
	clients, err := getTestClients(context.Background())
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
