package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestCreateDeleteRoutingZone(t *testing.T) {
	client, err := blueprintsTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	blueprints, err := client.listAllBlueprintIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(blueprints) == 0 {
		t.Skipf("cannot proceed without at least one blueprint")
	}

	dcClient, err := client.NewTwoStageL3ClosClient(context.TODO(), blueprints[0])
	if err != nil {
		t.Fatal(err)
	}

	randString := randString(5, "hex")

	label := "test-" + randString
	vrfName := "test-" + randString
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
	zone, err := dcClient.getSecurityZone(context.TODO(), zoneId)
	if err != nil {
		t.Fatal(err)
	}
	if zone.Id != zoneId {
		t.Fatalf("created vs. fetched zone IDs don't match: '%s' and '%s'", zone.Id, zoneId)
	}

	log.Println("fetching by vrf name...")
	zone, err = dcClient.getSecurityZoneByName(context.TODO(), vrfName)
	if err != nil {
		t.Fatal(err)
	}
	if zone.Id != zoneId {
		t.Fatalf("created vs. fetched zone IDs don't match: '%s' and '%s'", zone.Id, zoneId)
	}

	log.Println("deleting security zone...")

	err = dcClient.DeleteSecurityZone(context.TODO(), zoneId)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDefaultRoutingZone(t *testing.T) {
	client, err := blueprintsTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	blueprints, err := client.listAllBlueprintIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(blueprints) == 0 {
		t.Skipf("cannot proceed without at least one blueprint")
	}

	for _, bpId := range blueprints {
		dcClient, err := client.NewTwoStageL3ClosClient(context.TODO(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		sz, err := dcClient.GetSecurityZoneByName(context.TODO(), "default")
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("blueprint: %s - default security zone: %s", bpId, sz.Id)
	}
}
