//go:build integration
// +build integration

package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestTwoStageL3ClosVirtualNetworkStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		int() int
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
		{stringVal: "", intType: SviIpRequirementNone, stringType: sviIpRequirementNone},
		{stringVal: "optional", intType: SviIpRequirementOptional, stringType: sviIpRequirementOptional},
		{stringVal: "forbidden", intType: SviIpRequirementForbidden, stringType: sviIpRequirementForbidden},
		{stringVal: "mandatory", intType: SviIpRequirementMandatory, stringType: sviIpRequirementMandatory},
		{stringVal: "intention_conflict", intType: SviIpRequirementIntentionConflict, stringType: sviIpRequirementIntentionConflict},

		{stringVal: "disabled", intType: Ipv4ModeDisabled, stringType: ipv4ModeDisabled},
		{stringVal: "enabled", intType: Ipv4ModeEnabled, stringType: ipv4ModeEnabled},
		{stringVal: "forced", intType: Ipv4ModeForced, stringType: ipv4ModeForced},

		{stringVal: "disabled", intType: Ipv6ModeDisabled, stringType: ipv6ModeDisabled},
		{stringVal: "enabled", intType: Ipv6ModeEnabled, stringType: ipv6ModeEnabled},
		{stringVal: "forced", intType: Ipv6ModeForced, stringType: ipv6ModeForced},
		{stringVal: "link_local", intType: Ipv6ModeLinkLocal, stringType: ipv6ModeLinkLocal},

		{stringVal: "vlan", intType: VnTypeVlan, stringType: vnTypeVlan},
		{stringVal: "vxlan", intType: VnTypeVxlan, stringType: vnTypeVxlan},
		{stringVal: "overlay", intType: VnTypeOverlay, stringType: vnTypeOverlay},
	}

	for i, td := range testData {
		ii := td.intType.int()
		is := td.intType.String()
		sp, err := td.stringType.parse()
		if err != nil {
			t.Fatal(err)
		}
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.int() != sp ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}

func TestGetAllVirtualNetworks(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	skipped := true

	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		blueprints, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(blueprints) == 0 {
			continue
		}

		skipped = false

		for i := range samples(len(blueprints)) {
			id := blueprints[i]
			bpStatus, err := client.client.getBlueprintStatus(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
			switch bpStatus.Design {
			case refDesignFreeform:
				log.Printf("'%s' design is '%s.", id, bpStatus.Design)
				// todo
			case refDesignDatacenter:
				log.Printf("'%s' design is '%s.", id, bpStatus.Design)
				log.Printf("testing NewTwoStageL3ClosClient(%s) against %s %s (%s)", id, client.clientType, clientName, client.client.ApiVersion())
				dcClient, err := client.client.NewTwoStageL3ClosClient(context.TODO(), id)
				if err != nil {
					t.Fatal(err)
				}

				log.Printf("testing listAllVirtualNetworkIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				vns, err := dcClient.listAllVirtualNetworkIds(context.TODO(), BlueprintTypeStaging)
				if err != nil {
					t.Fatal(err)
				}

				for _, id := range vns {
					log.Printf("testing getVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					vn, err := dcClient.getVirtualNetwork(context.TODO(), id, BlueprintTypeStaging)
					if err != nil {
						t.Fatal(err)
					}
					if vn.Ipv4Subnet != nil {
						log.Printf("testing getVirtualNetworkBySubnet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
						vn, err = dcClient.getVirtualNetworkBySubnet(context.TODO(), vn.Ipv4Subnet, vn.SecurityZoneId, BlueprintTypeStaging)
						if err != nil {
							t.Fatal(err)
						}
					}
					if vn.Ipv6Subnet != nil {
						log.Printf("testing getVirtualNetworkBySubnet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
						vn, err = dcClient.getVirtualNetworkBySubnet(context.TODO(), vn.Ipv6Subnet, vn.SecurityZoneId, BlueprintTypeStaging)
						if err != nil {
							t.Fatal(err)
						}
					}
					log.Printf("vn: %s ipv4: %s, ipv6:%s\n", vn.Id, vn.Ipv4Subnet.String(), vn.Ipv6Subnet.String())
				}
			}
		}
	}
	if skipped {
		t.Skip("no blueprints found on any test system")
	}
}
