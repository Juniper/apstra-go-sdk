// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
)

func TestListGetOneRackType(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listRackTypeIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		rtIds, err := client.client.listRackTypeIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		// id := rtIds[rand.Intn(len(rtIds))]
		id := rtIds[0]

		log.Printf("testing getRackType(%s) against %s %s (%s)", id, client.clientType, clientName, client.client.ApiVersion())
		rt, err := client.client.GetRackType(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}

		log.Println(rt.Id)
	}
}

func TestListGetAllGetRackType(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listRackTypeIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		rackTypeIds, err := client.client.listRackTypeIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		for _, i := range samples(t, len(rackTypeIds)) {
			id := rackTypeIds[i]
			log.Printf("testing getRackType(%s) against %s %s (%s)", id, client.clientType, clientName, client.client.ApiVersion())
			rt, err := client.client.GetRackType(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
			log.Println(rt.Id)
		}

		log.Printf("testing getAllRackTypes() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		rackTypes, err := client.client.getAllRackTypes(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(rackTypeIds) != len(rackTypes) {
			t.Fatalf("got %d rack type IDs but %d rack types", len(rackTypeIds), len(rackTypes))
		}

		randRackid := rackTypeIds[rand.Intn(len(rackTypeIds))]
		log.Printf("testing getRackType() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		rt, err := client.client.GetRackType(context.TODO(), randRackid)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("randomly selected rack type '%s' (%s) has %d leaf switches, %d access switches, and %d generic systems",
			rt.Data.DisplayName, rt.Id, len(rt.Data.LeafSwitches), len(rt.Data.AccessSwitches), len(rt.Data.GenericSystems))
	}
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

func TestCreateGetRackDeleteRackType(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	leafLabel := "ll-" + randString(10, "hex")

	testCases := map[string]RackTypeRequest{
		"leaf_only_no_tags": {
			DisplayName:              "rdn " + randString(5, "hex"),
			FabricConnectivityDesign: FabricConnectivityDesignL3Clos,
			LeafSwitches: []RackElementLeafSwitchRequest{
				{
					Label:             leafLabel,
					LogicalDeviceId:   "AOS-48x10_6x40-leaf_spine",
					LinkPerSpineCount: 2,
					LinkPerSpineSpeed: "10G",
				},
			},
		},
		"leaf_generic_with_tags": {
			DisplayName:              "rdn " + randString(5, "hex"),
			FabricConnectivityDesign: FabricConnectivityDesignL3Clos,
			LeafSwitches: []RackElementLeafSwitchRequest{
				{
					Label:             leafLabel,
					LogicalDeviceId:   "AOS-48x10_6x40-leaf_spine",
					LinkPerSpineCount: 2,
					LinkPerSpineSpeed: "10G",
					Tags:              []ObjectId{"hypervisor", "bare_metal"},
				},
			},
			GenericSystems: []RackElementGenericSystemRequest{
				{
					Count: 5,
					Label: "some generic system",
					Links: []RackLinkRequest{
						{
							Label:              "foo",
							LinkPerSwitchCount: 1,
							LinkSpeed:          "10G",
							TargetSwitchLabel:  leafLabel,
							AttachmentType:     RackLinkAttachmentTypeSingle,
							LagMode:            RackLinkLagModeNone,
							Tags:               []ObjectId{"firewall"},
						},
					},
					LogicalDeviceId: "AOS-1x10-1",
					Tags:            []ObjectId{"firewall"},
				},
			},
			AccessSwitches: nil,
		},
	}

	for clientName, client := range clients {
		for _, tCase := range testCases {
			log.Printf("testing CreateRackType() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			id, err := client.client.CreateRackType(context.TODO(), &tCase)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetRackType() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			rt, err := client.client.GetRackType(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}

			if id != rt.Id {
				t.Fatalf("expected %q, got %q", id, rt.Id)
			}

			log.Printf("testing DeleteRackType() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.DeleteRackType(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
