//go:build integration
// +build integration

package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestCreateUpdateDeleteRoutingZone(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx)
	if err != nil {
		t.Fatal(err)
	}

	randStr := randString(5, "hex")
	label := "test-" + randStr
	vrfName := "test-" + randStr

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintA(ctx, t, client.client)
		defer func() {
			err := bpDel()
			if err != nil {
				t.Fatal(err)
			}
		}()

		log.Printf("testing CreateSecurityZone() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		zoneId, err := bpClient.CreateSecurityZone(ctx, &SecurityZoneData{
			SzType:  SecurityZoneTypeEVPN,
			VrfName: vrfName,
			Label:   label,
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("created zone - id:'%s', name: '%s', label:'%s'", zoneId, vrfName, label)

		log.Println("fetching by id...")
		log.Printf("testing GetSecurityZone() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		zone, err := bpClient.GetSecurityZone(ctx, zoneId)
		if err != nil {
			t.Fatal(err)
		}
		if zone.Id != zoneId {
			t.Fatalf("created vs. fetched zone IDs don't match: '%s' and '%s'", zone.Id, zoneId)
		}

		log.Println("fetching by vrf name...")
		log.Printf("testing getSecurityZoneByVrfName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		zone, err = bpClient.GetSecurityZoneByVrfName(ctx, vrfName)
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
		err = bpClient.UpdateSecurityZone(ctx, zoneId, &SecurityZoneData{
			SzType:  SecurityZoneTypeEVPN,
			VrfName: vrfName2,
			Label:   label2,
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetSecurityZoneByVrfName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		zone, err = bpClient.GetSecurityZoneByVrfName(ctx, vrfName2)
		if err != nil {
			t.Fatal(err)
		}
		if zone.Id != zoneId {
			t.Fatal()
		}
		if zone.Data.VrfName != vrfName2 {
			t.Fatal()
		}

		log.Printf("testing GetAllSecurityZones() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		zones, err := bpClient.GetAllSecurityZones(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if len(zones) != 2 {
			t.Fatalf("expected 2 security zones, got %d", len(zones))
		}

		log.Printf("testing DeleteSecurityZone() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeleteSecurityZone(ctx, zoneId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetDefaultRoutingZone(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintA(ctx, t, client.client)
		defer func() {
			err := bpDel()
			if err != nil {
				t.Fatal(err)
			}
		}()

		log.Printf("testing GetSecurityZoneByVrfName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		sz, err := bpClient.GetSecurityZoneByVrfName(ctx, "default")
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("blueprint: %s - default security zone: %s", bpClient.blueprintId, sz.Id)
	}
}
