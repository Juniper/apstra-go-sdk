package goapstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestListGetAllGetRackType(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	rackTypeIds, err := client.listRackTypeIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	rackTypes, err := client.getAllRackTypes(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(rackTypeIds) != len(rackTypes) {
		t.Fatalf("got %d rack type IDs but %d rack types", len(rackTypeIds), len(rackTypes))
	}

	rand.Seed(time.Now().UnixNano())
	randRackid := rackTypeIds[rand.Intn(len(rackTypeIds))]
	rt, err := client.getRackType(context.TODO(), randRackid)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("randomly selected rack type '%s' (%s) has %d leaf switches, %d access switches, and %d generic systems",
		rt.DisplayName, rt.Id, len(rt.LeafSwitches), len(rt.AccessSwitches), len(rt.GenericSystems))
}

func TestRackTypeStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		Int() int
	}

	type apiIotaString interface {
		parse() (int, error)
		string() string
	}

	type stringTestData struct {
		stringVal  string
		intType    apiStringIota
		stringType apiIotaString
	}
	testData := []stringTestData{
		{stringVal: "", intType: LeafRedundancyProtocolNone, stringType: leafRedundancyProtocolNone},
		{stringVal: "mlag", intType: LeafRedundancyProtocolMlag, stringType: leafRedundancyProtocolMlag},
		{stringVal: "esi", intType: LeafRedundancyProtocolEsi, stringType: leafRedundancyProtocolEsi},

		{stringVal: "", intType: AccessRedundancyProtocolNone, stringType: accessRedundancyProtocolNone},
		{stringVal: "esi", intType: AccessRedundancyProtocolEsi, stringType: accessRedundancyProtocolEsi},

		{stringVal: "l3clos", intType: FabricConnectivityDesignL3Clos, stringType: fabricConnectivityDesignL3Clos},
		{stringVal: "l3collapsed", intType: FabricConnectivityDesignL3Collapsed, stringType: fabricConnectivityDesignL3Collapsed},

		{stringVal: "singleAttached", intType: RackLinkAttachmentTypeSingle, stringType: rackLinkAttachmentTypeSingle},
		{stringVal: "dualAttached", intType: RackLinkAttachmentTypeDual, stringType: rackLinkAttachmentTypeDual},

		{stringVal: "", intType: RackLinkLagModeNone, stringType: rackLinkLagModeNone},
		{stringVal: "lacp_active", intType: RackLinkLagModeActive, stringType: rackLinkLagModeActive},
		{stringVal: "lacp_passive", intType: RackLinkLagModePassive, stringType: rackLinkLagModePassive},
		{stringVal: "static_lag", intType: RackLinkLagModeStatic, stringType: rackLinkLagModeStatic},
	}

	for i, td := range testData {
		ii := td.intType.Int()
		is := td.intType.String()
		sp, err := td.stringType.parse()
		if err != nil {
			t.Fatal(err)
		}
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.Int() != sp ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}

func TestCreateGetRackType(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	leafLabel := "ll-" + randString(10, "hex")

	id, err := client.createRackType(context.TODO(), &RackType{
		DisplayName:              "rdn " + randString(5, "hex"),
		FabricConnectivityDesign: FabricConnectivityDesignL3Clos,
		LeafSwitches: []RackElementLeafSwitch{
			{
				Label:             leafLabel,
				LogicalDeviceId:   "virtual-7x10-1",
				LinkPerSpineCount: 2,
				LinkPerSpineSpeed: "10G",
			},
		},
		GenericSystems: []RackElementGenericSystem{
			{
				Count: 5,
				Label: "some generic system",
				Links: []RackLink{
					{
						Label:              "foo",
						LinkPerSwitchCount: 1,
						LinkSpeed:          "10G",
						TargetSwitchLabel:  leafLabel,
						AttachmentType:     RackLinkAttachmentTypeSingle,
						LagMode:            RackLinkLagModeNone,
					},
				},
				LogicalDeviceId: "5ed7ed07-7222-4d6c-a5cb-1e1aa6036dab",
			},
		},
		AccessSwitches: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("id: '%s'\n", id)
}
