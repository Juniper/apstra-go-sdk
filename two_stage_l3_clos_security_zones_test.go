package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestCreateUpdateDeleteRoutingZone(t *testing.T) {
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

	randStr := randString(5, "hex")

	label := "test-" + randStr
	vrfName := "test-" + randStr
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

	randStr2 := randString(5, "hex")
	vrfName2 := "test-" + randStr2
	label2 := "test-" + randStr2
	err = dcClient.UpdateSecurityZone(context.TODO(), zoneId, &CreateSecurityZoneCfg{
		SzType:  "evpn",
		VrfName: vrfName2,
		Label:   label2,
	})
	if err != nil {
		t.Fatal(err)
	}

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
