// Copyright (c) Juniper Networks, Inc., 2022-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/datacenter"
)

func TestGetAllPolicies(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintA(ctx, t, client.client)

		log.Printf("testing getAllPolicies() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		policies, err := bpClient.getAllPolicies(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for _, policy := range policies {
			log.Printf("testing getPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			p, err := bpClient.getPolicy(ctx, policy.Id)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("policy '%s'\t'%s'", p.Id, p.Label)
		}
	}
}

func comparePolicyRules(aName string, a PolicyRule, bName string, b PolicyRule, t *testing.T) {
	if a.Id != b.Id {
		t.Fatalf("Policy Rule IDs don't match: %s has %q, %s has %q", aName, a.Id, bName, b.Id)
	}

	aData := a.Data != nil
	bData := b.Data != nil

	if (aData || bData) && !(aData && bData) { // xor
		t.Fatalf("Policy Rule data presence mismatch -- a: %t vs. b: %t", aData, bData)
	}

	if aData && bData {
		comparePolicyRuleData(aName, a.Data, bName, b.Data, t)
	}
}

func comparePolicyRuleData(aName string, a *PolicyRuleData, bName string, b *PolicyRuleData, t *testing.T) {
	if a.Label != b.Label {
		t.Fatalf("Policy Rule Labels don't match: %s has %q, %s has %q", aName, a.Label, bName, b.Label)
	}

	if a.Description != b.Description {
		t.Fatalf("Policy Rule Descriptions don't match: %s has %q, %s has %q", aName, a.Description, bName, b.Description)
	}

	if a.Protocol != b.Protocol {
		t.Fatalf("Policy Rule Protocols don't match: %s has %q, %s has %q", aName, a.Protocol, bName, b.Protocol)
	}

	if a.Action != b.Action {
		t.Fatalf("Policy Rule Actions don't match: %s has %s, %s has %s", aName, a.Action, bName, b.Action)
	}

	comparedatacenter.PortRanges(t, a.SrcPort, b.SrcPort)
	comparedatacenter.PortRanges(t, a.DstPort, b.DstPort)

	aTcpStateQualifier := a.TcpStateQualifier != nil
	bTcpStateQualifier := b.TcpStateQualifier != nil
	if (aTcpStateQualifier || bTcpStateQualifier) && !(aTcpStateQualifier && bTcpStateQualifier) { // xor
		t.Fatalf("TCP state qualifier presence mismatch -- a: %t vs. b: %t", aTcpStateQualifier, bTcpStateQualifier)
	}

	if aTcpStateQualifier && bTcpStateQualifier && (a.TcpStateQualifier.Value != b.TcpStateQualifier.Value) {
		t.Fatalf("TCP state qualifier value mismatch -- a: %q vs. b: %q", a.TcpStateQualifier.Value, b.TcpStateQualifier.Value)
	}
}

func comparePolicies(a *Policy, aName string, b *Policy, bName string, t *testing.T) {
	if a.Id != b.Id {
		t.Fatalf("Policy IDs don't match: %s has %q, %s has %q", aName, a.Id, bName, b.Id)
	}

	comparePolicyData(a.Data, aName, b.Data, bName, t)
}

func comparePolicyData(a *PolicyData, aName string, b *PolicyData, bName string, t *testing.T) {
	if a.Enabled != b.Enabled {
		t.Fatalf("Policy enabled switches don't match: %s has %t, %s has %t", aName, a.Enabled, bName, b.Enabled)
	}

	if a.Label != b.Label {
		t.Fatalf("Policy Labels don't match: %s has %q, %s has %q", aName, a.Label, bName, b.Label)
	}

	if a.Description != b.Description {
		t.Fatalf("Policy Descriptions don't match: %s has %q, %s has %q", aName, a.Description, bName, b.Description)
	}

	if a.SrcApplicationPoint.Id != b.SrcApplicationPoint.Id {
		t.Fatalf("Policy SrcApplicationPoints don't match: %s has %q, %s has %q", aName, a.SrcApplicationPoint.Id, bName, b.SrcApplicationPoint.Id)
	}

	if a.DstApplicationPoint.Id != b.DstApplicationPoint.Id {
		t.Fatalf("Policy DstApplicationPoints don't match: %s has %q, %s has %q", aName, a.DstApplicationPoint.Id, bName, b.DstApplicationPoint.Id)
	}

	compareSlicesAsSets(t, a.Tags, b.Tags, fmt.Sprintf("%s tags: %v, %s tags %v", aName, a.Tags, bName, b.Tags))

	if len(a.Rules) != len(b.Rules) {
		t.Fatalf("Policy ruleset sizes don't match: %s has %d rules, %s has %d rules", aName, len(a.Rules), bName, len(b.Rules))
	}

	for i := 0; i < len(a.Rules); i++ {
		comparePolicyRules(fmt.Sprintf("%s rule %d", aName, i), a.Rules[i], fmt.Sprintf("%s rule %d", bName, i), b.Rules[i], t)
	}
}

func TestCreateDatacenterPolicy(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			bp := testBlueprintA(ctx, t, client.client)

			// collect leaf switch IDs
			leafIds, err := getSystemIdsByRole(ctx, bp, "leaf")
			if err != nil {
				t.Fatal(err)
			}

			// prep VN bindings
			bindings := make([]datacenter.VNBinding, len(leafIds))
			for i, leafId := range leafIds {
				bindings[i] = datacenter.VNBinding{SystemID: string(leafId)}
			}

			// create a security zone (VNs live here)
			szName := randString(5, "hex")
			szId, err := bp.CreateSecurityZone(ctx, datacenter.SecurityZone{
				Type:    enum.SecurityZoneTypeEVPN,
				Label:   szName,
				VRFName: szName,
			})
			if err != nil {
				t.Fatal(err)
			}

			// create a couple of virtual networks we'll use as policy rule endpoints
			vnIds := make([]string, 2)
			for i := range vnIds {
				vnId, err := bp.CreateVirtualNetwork(ctx, datacenter.VirtualNetwork{
					IPv4Enabled:    true,
					Label:          "vn_" + strconv.Itoa(i),
					SecurityZoneID: szId,
					Bindings:       bindings,
					Type:           enum.VnTypeVxlan,
				})
				if err != nil {
					t.Fatal(err)
				}
				vnIds[i] = vnId
			}

			tags := make([]string, rand.Intn(4))
			for i := range tags {
				tags[i] = randString(5, "hex")
			}

			policyDatas := []PolicyData{
				{
					Enabled:             randBool(),
					Label:               randString(5, "hex"),
					Description:         randString(5, "hex"),
					SrcApplicationPoint: &datacenter.PolicyApplicationPointData{Id: vnIds[0]},
					DstApplicationPoint: &datacenter.PolicyApplicationPointData{Id: vnIds[1]},
					Rules:               nil,
					Tags:                tags,
				},
				{
					Enabled:             randBool(),
					Label:               randString(5, "hex"),
					Description:         randString(5, "hex"),
					SrcApplicationPoint: &datacenter.PolicyApplicationPointData{Id: vnIds[1]},
					DstApplicationPoint: &datacenter.PolicyApplicationPointData{Id: vnIds[0]},
					Rules:               nil,
					Tags:                tags,
				},
			}

			var previousPolicy *Policy
			var previousPolicyId ObjectId
			for i, policyData := range policyDatas {
				policyData := policyData
				if previousPolicy == nil {
					log.Printf("testing CreatePolicy(%d) against %s %s (%s)", i, client.clientType, clientName, client.client.ApiVersion())
					previousPolicyId, err = bp.CreatePolicy(ctx, &policyData)
					if err != nil {
						t.Fatal(err)
					}

					created := Policy{
						Id:   previousPolicyId,
						Data: &policyData,
					}

					log.Printf("testing GetPolicy(%d) against %s %s (%s)", i, client.clientType, clientName, client.client.ApiVersion())
					previousPolicy, err = bp.GetPolicy(ctx, previousPolicyId)
					if err != nil {
						t.Fatal(err)
					}

					comparePolicies(&created, "created", previousPolicy, "fetched", t)
				}

				log.Printf("testing UpdatePolicy(%d) against %s %s (%s)", i, client.clientType, clientName, client.client.ApiVersion())
				err = bp.UpdatePolicy(ctx, previousPolicyId, &policyData)
				if err != nil {
					t.Fatal(err)
				}

				log.Printf("testing GetPolicy(%d) against %s %s (%s)", i, client.clientType, clientName, client.client.ApiVersion())
				previousPolicy, err = bp.GetPolicy(ctx, previousPolicyId)
				if err != nil {
					t.Fatal(err)
				}

				comparePolicies(&Policy{
					Id:   previousPolicyId,
					Data: &policyData,
				}, "updated", previousPolicy, "fetched", t)
			}
		})
	}
}

func TestAddDeletePolicyRule(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			bp := testBlueprintA(ctx, t, client.client)

			// collect leaf switch IDs
			leafIds, err := getSystemIdsByRole(ctx, bp, "leaf")
			if err != nil {
				t.Fatal(err)
			}

			// prep VN bindings
			bindings := make([]datacenter.VNBinding, len(leafIds))
			for i, leafId := range leafIds {
				bindings[i] = datacenter.VNBinding{SystemID: string(leafId)}
			}

			// create a security zone (VNs live here)
			szName := randString(5, "hex")
			szId, err := bp.CreateSecurityZone(ctx, datacenter.SecurityZone{
				Type:    enum.SecurityZoneTypeEVPN,
				Label:   szName,
				VRFName: szName,
			})
			if err != nil {
				t.Fatal(err)
			}

			// create a couple of virtual networks we'll use a policy rule endpoints
			vnIds := make([]string, 2)
			for i := range vnIds {
				vnId, err := bp.CreateVirtualNetwork(ctx, datacenter.VirtualNetwork{
					IPv4Enabled:    true,
					Label:          "vn_" + strconv.Itoa(i),
					SecurityZoneID: szId,
					Bindings:       bindings,
					Type:           enum.VnTypeVxlan,
				})
				if err != nil {
					t.Fatal(err)
				}
				vnIds[i] = vnId
			}

			// create a security policy
			policyId, err := bp.CreatePolicy(ctx, &PolicyData{
				Enabled:             false,
				Label:               randString(5, "hex"),
				SrcApplicationPoint: &datacenter.PolicyApplicationPointData{Id: vnIds[0]},
				DstApplicationPoint: &datacenter.PolicyApplicationPointData{Id: vnIds[1]},
			})
			if err != nil {
				t.Fatal(err)
			}

			newRule := &PolicyRuleData{
				Label:             randString(5, "hex"),
				Description:       randString(5, "hex"),
				Protocol:          enum.PolicyRuleProtocolTcp,
				Action:            enum.PolicyRuleActionDenyLog,
				SrcPort:           datacenter.PortRanges{{5, 6}},
				DstPort:           datacenter.PortRanges{{7, 8}, {9, 10}},
				TcpStateQualifier: &enum.TcpStateQualifierEstablished,
			}

			p, err := bp.getPolicy(ctx, policyId)
			if err != nil {
				t.Fatal(err)
			}
			ruleCount := len(p.Rules)

			log.Printf("testing addPolicyRule() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			ruleId, err := bp.AddPolicyRule(ctx, newRule, 0, policyId)
			if err != nil {
				t.Fatal(err)
			}

			p, err = bp.getPolicy(ctx, policyId)
			if err != nil {
				t.Fatal(err)
			}
			if len(p.Rules) != ruleCount+1 {
				t.Fatalf("expected %d rules, got %d rules", ruleCount+1, len(p.Rules))
			}

			log.Printf("testing deletePolicyRuleById() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bp.deletePolicyRuleById(ctx, policyId, ruleId)
			if err != nil {
				t.Fatal(err)
			}

			p, err = bp.getPolicy(ctx, policyId)
			if err != nil {
				t.Fatal(err)
			}
			if len(p.Rules) != ruleCount {
				t.Fatalf("expected %d rules, got %d rules", ruleCount, len(p.Rules))
			}
		})
	}
}
