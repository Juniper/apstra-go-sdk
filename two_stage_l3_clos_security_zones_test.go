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

	randString := randString(10, "hex")

	label := "label-test-" + randString
	zoneId, err := dcClient.CreateSecurityZone(context.TODO(), &CreateSecurityZoneCfg{
		SzType:  "evpn",
		VrfName: "test-" + randString,
		Label:   label,
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("created zone - id:'%s', label:'%s'", zoneId, label)

	log.Println("fetching by id...")
	zone, err := dcClient.getSecurityZone(context.TODO(), zoneId)
	if err != nil {
		t.Fatal(err)
	}
	if zone.Id != zoneId {
		t.Fatalf("created vs. fetched zone IDs don't match: '%s' and '%s'", zone.Id, zoneId)
	}

	log.Println("fetching by label...")
	zone, err = dcClient.getSecurityZoneByLabel(context.TODO(), label)
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
