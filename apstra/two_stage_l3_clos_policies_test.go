//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
)

func TestPortRangeString(t *testing.T) {
	var tests []struct {
		data     PortRange
		expected string
	}
	tests = append(tests, struct {
		data     PortRange
		expected string
	}{data: PortRange{first: 10, last: 10}, expected: "10"})

	tests = append(tests, struct {
		data     PortRange
		expected string
	}{data: PortRange{first: 10, last: 20}, expected: "10-20"})

	tests = append(tests, struct {
		data     PortRange
		expected string
	}{data: PortRange{first: 20, last: 10}, expected: "10-20"})

	for _, test := range tests {
		if test.expected != test.data.string() {
			t.Fatalf("expected '%s', got '%s'", test.expected, test.data.string())
		}
	}
}

func portRangeSlicesMatch(a, b []PortRange) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i].first != b[i].first {
			return false
		}
		if a[i].last != b[i].last {
			return false
		}
	}

	return true
}

func TestRawPortRangesParse(t *testing.T) {
	var tests []struct {
		data     rawPortRanges
		expected []PortRange
	}
	tests = append(tests, struct {
		data     rawPortRanges
		expected []PortRange
	}{data: "10", expected: []PortRange{{first: 10, last: 10}}})

	tests = append(tests, struct {
		data     rawPortRanges
		expected []PortRange
	}{data: "10,11", expected: []PortRange{{first: 10, last: 10}, {first: 11, last: 11}}})

	tests = append(tests, struct {
		data     rawPortRanges
		expected []PortRange
	}{data: "12,11", expected: []PortRange{{first: 12, last: 12}, {first: 11, last: 11}}})

	tests = append(tests, struct {
		data     rawPortRanges
		expected []PortRange
	}{data: "10-11,12-13", expected: []PortRange{{first: 10, last: 11}, {first: 12, last: 13}}})

	for i, test := range tests {
		parsed, err := test.data.parse()
		if err != nil {
			t.Fatal(err)
		}
		if !portRangeSlicesMatch(parsed, test.expected) {
			t.Fatalf("no match test index %d", i)
		}
	}
}

func TestGetAllPolicies(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintA(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

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

func comparePolicyPortRanges(a PortRange, aName string, b PortRange, bName string, t *testing.T) {
	if a.first != b.first {
		t.Fatalf("Policy Port Ranges 'first' field don't match: %s has %d, %s has %d", aName, a.first, bName, b.first)
	}

	if a.last != b.last {
		t.Fatalf("Policy Port Ranges 'last' field don't match: %s has %d, %s has %d", aName, a.last, bName, b.last)
	}
}

func comparePolicyRules(aName string, a PolicyRule, bName string, b PolicyRule, t *testing.T) {
	if a.Id != b.Id {
		t.Fatalf("Policy Rule IDs don't match: %s has %q, %s has %q", aName, a.Id, bName, b.Id)
	}

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

	if len(a.SrcPort) != len(b.SrcPort) {
		t.Fatalf("Policy Rule Src Port Ranges don't match %s has %d ranges, %s has %d ranges", aName, len(a.SrcPort), bName, b.SrcPort)
	}

	for i := 0; i < len(a.SrcPort); i++ {
		comparePolicyPortRanges(a.SrcPort[i], fmt.Sprintf("%s rule %d SrcPort", aName, i), b.SrcPort[i], fmt.Sprintf("%s rule %d SrcPort", bName, i), t)
	}

	if len(a.DstPort) != len(b.DstPort) {
		t.Fatalf("Policy Rule Dst Port Ranges don't match %s has %d ranges, %s has %d ranges", aName, len(a.DstPort), bName, b.DstPort)
	}

	for i := 0; i < len(a.DstPort); i++ {
		comparePolicyPortRanges(a.DstPort[i], fmt.Sprintf("%s rule %d DstPort", aName, i), b.DstPort[i], fmt.Sprintf("%s rule %d DstPort", bName, i), t)
	}
}

func comparePolicies(a *Policy, aName string, b *Policy, bName string, t *testing.T) {
	if a.Id != b.Id {
		t.Fatalf("Policy IDs don't match: %s has %q, %s has %q", aName, a.Id, bName, b.Id)
	}

	if a.Enabled != b.Enabled {
		t.Fatalf("Policy enabled switches don't match: %s has %t, %s has %t", aName, a.Enabled, bName, b.Enabled)
	}

	if a.Label != b.Label {
		t.Fatalf("Policy Labels don't match: %s has %q, %s has %q", aName, a.Label, bName, b.Label)
	}

	if a.Description != b.Description {
		t.Fatalf("Policy Descriptions don't match: %s has %q, %s has %q", aName, a.Description, bName, b.Description)
	}

	if a.SrcApplicationPoint.ObjectId() != b.SrcApplicationPoint.ObjectId() {
		t.Fatalf("Policy SrcApplicationPoints don't match: %s has %q, %s has %q", aName, a.SrcApplicationPoint.ObjectId(), bName, b.SrcApplicationPoint.ObjectId())
	}

	if a.DstApplicationPoint.ObjectId() != b.DstApplicationPoint.ObjectId() {
		t.Fatalf("Policy DstApplicationPoints don't match: %s has %q, %s has %q", aName, a.DstApplicationPoint.ObjectId(), bName, b.DstApplicationPoint.ObjectId())
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
		bp, bpDel := testBlueprintA(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		// collect leaf switch IDs
		leafIds, err := getSystemIdsByRole(ctx, bp, "leaf")
		if err != nil {
			t.Fatal(err)
		}

		// prep VN bindings
		bindings := make([]VnBinding, len(leafIds))
		for i, leafId := range leafIds {
			bindings[i] = VnBinding{SystemId: leafId}
		}

		// create a security zone (VNs live here)
		szName := randString(5, "hex")
		szId, err := bp.CreateSecurityZone(ctx, &SecurityZoneData{
			SzType:  SecurityZoneTypeEVPN,
			Label:   szName,
			VrfName: szName,
		})
		if err != nil {
			t.Fatal(err)
		}

		// create a couple of virtual networks we'll use a policy rule endpoints
		vnIds := make([]ObjectId, 2)
		for i := range vnIds {
			vnId, err := bp.CreateVirtualNetwork(ctx, &VirtualNetworkData{
				Ipv4Enabled:    true,
				Label:          "vn_" + strconv.Itoa(i),
				SecurityZoneId: szId,
				VnBindings:     bindings,
				VnType:         VnTypeVxlan,
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

		policy1 := Policy{
			Enabled:             randBool(),
			Label:               randString(5, "hex"),
			Description:         randString(5, "hex"),
			SrcApplicationPoint: vnIds[0],
			DstApplicationPoint: vnIds[1],
			Rules:               nil,
			Tags:                tags,
		}

		log.Printf("testing CreatePolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		policyId, err := bp.CreatePolicy(ctx, &policy1)
		if err != nil {
			t.Fatal(err)
		}

		policy1.Id = policyId

		policy2, err := bp.GetPolicy(ctx, policy1.Id)
		if err != nil {
			t.Fatal(err)
		}

		comparePolicies(&policy1, "as created", policy2, "as fetched", t)
	}
}

func TestAddDeletePolicyRule(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bp, bpDelete := testBlueprintA(ctx, t, client.client)
		defer func() {
			err := bpDelete(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		// collect leaf switch IDs
		leafIds, err := getSystemIdsByRole(ctx, bp, "leaf")
		if err != nil {
			t.Fatal(err)
		}

		// prep VN bindings
		bindings := make([]VnBinding, len(leafIds))
		for i, leafId := range leafIds {
			bindings[i] = VnBinding{SystemId: leafId}
		}

		// create a security zone (VNs live here)
		szName := randString(5, "hex")
		szId, err := bp.CreateSecurityZone(ctx, &SecurityZoneData{
			SzType:  SecurityZoneTypeEVPN,
			Label:   szName,
			VrfName: szName,
		})
		if err != nil {
			t.Fatal(err)
		}

		// create a couple of virtual networks we'll use a policy rule endpoints
		vnIds := make([]ObjectId, 2)
		for i := range vnIds {
			vnId, err := bp.CreateVirtualNetwork(ctx, &VirtualNetworkData{
				Ipv4Enabled:    true,
				Label:          "vn_" + strconv.Itoa(i),
				SecurityZoneId: szId,
				VnBindings:     bindings,
				VnType:         VnTypeVxlan,
			})
			if err != nil {
				t.Fatal(err)
			}
			vnIds[i] = vnId
		}

		// create a security policy
		policyId, err := bp.CreatePolicy(ctx, &Policy{
			Enabled:             false,
			Label:               randString(5, "hex"),
			SrcApplicationPoint: vnIds[0],
			DstApplicationPoint: vnIds[1],
		})
		if err != nil {
			t.Fatal(err)
		}

		newRule := &PolicyRule{
			Label:       randString(5, "hex"),
			Description: randString(5, "hex"),
			Protocol:    "TCP",
			Action:      PolicyRuleActionDenyLog,
			SrcPort:     PortRanges{{5, 6}},
			DstPort:     PortRanges{{7, 8}, {9, 10}},
		}

		p, err := bp.getPolicy(ctx, policyId)
		if err != nil {
			t.Fatal(err)
		}
		ruleCount := len(p.Rules)

		log.Printf("testing addPolicyRule() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ruleId, err := bp.addPolicyRule(context.TODO(), newRule, 0, policyId)
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
		err = bp.deletePolicyRuleById(context.TODO(), policyId, ruleId)
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
	}
}
