package goapstra

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestCreateUpdateDeleteRoutingZone(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	skipMsg := make(map[string]string)
	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		blueprints, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(blueprints) == 0 {
			skipMsg[clientName] = fmt.Sprintf("cannot manipualte routing zone in '%s' without blueprints", clientName)
			continue
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		dcClient, err := client.client.NewTwoStageL3ClosClient(context.TODO(), blueprints[0])
		if err != nil {
			t.Fatal(err)
		}

		randStr := randString(5, "hex")

		label := "test-" + randStr
		vrfName := "test-" + randStr
		log.Printf("testing CreateSecurityZone() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		zoneId, err := dcClient.CreateSecurityZone(context.TODO(), &CreateSecurityZoneCfg{
			SzType:  "evpn",
			VrfName: vrfName,
			Label:   label,
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("created zone - id:'%s', name: '%s', label:'%s'", zoneId, vrfName, label)

		log.Println("fetching by id...")
		log.Printf("testing getSecurityZone() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		zone, err := dcClient.getSecurityZone(context.TODO(), zoneId)
		if err != nil {
			t.Fatal(err)
		}
		if zone.Id != zoneId {
			t.Fatalf("created vs. fetched zone IDs don't match: '%s' and '%s'", zone.Id, zoneId)
		}

		log.Println("fetching by vrf name...")
		log.Printf("testing getSecurityZoneByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		zone, err = dcClient.getSecurityZoneByName(context.TODO(), vrfName)
		if err != nil {
			t.Fatal(err)
		}
		if zone.Id != zoneId {
			t.Fatalf("created vs. fetched zone IDs don't match: '%s' and '%s'", zone.Id, zoneId)
		}

		randStr2 := randString(5, "hex")
		vrfName2 := "test-" + randStr2
		label2 := "test-" + randStr2
		log.Printf("testing UpdateSecurityZone() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = dcClient.UpdateSecurityZone(context.TODO(), zoneId, &CreateSecurityZoneCfg{
			SzType:  "evpn",
			VrfName: vrfName2,
			Label:   label2,
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetSecurityZoneByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		zone, err = dcClient.GetSecurityZoneByName(context.TODO(), vrfName2)
		if err != nil {
			t.Fatal(err)
		}
		if zone.Id != zoneId {
			t.Fatal()
		}
		if zone.VrfName != vrfName2 {
			t.Fatal()
		}

		log.Printf("testing DeleteSecurityZone() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = dcClient.DeleteSecurityZone(context.TODO(), zoneId)
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(skipMsg) > 0 {
		sb := strings.Builder{}
		for _, msg := range skipMsg {
			sb.WriteString(msg + ";")
		}
		t.Skip(sb.String())
	}
}

func TestGetDefaultRoutingZone(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	skipMsg := make(map[string]string)
	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		blueprints, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(blueprints) == 0 {
			skipMsg[clientName] = fmt.Sprintf("cannot fetch routing zone from '%s' with no blueprints", clientName)
			continue
		}

		for _, bpId := range blueprints {
			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			dcClient, err := client.client.NewTwoStageL3ClosClient(context.TODO(), bpId)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetSecurityZoneByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			sz, err := dcClient.GetSecurityZoneByName(context.TODO(), "default")
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("blueprint: %s - default security zone: %s", bpId, sz.Id)
		}
	}
	if len(skipMsg) > 0 {
		sb := strings.Builder{}
		for _, msg := range skipMsg {
			sb.WriteString(msg + ";")
		}
		t.Skip(sb.String())
	}
}
