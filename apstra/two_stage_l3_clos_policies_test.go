//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
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

func TestCreateDatacenterPolicy(t *testing.T) {
	// todo: rewrite this test to create the required conditions
	t.Skip()
	return
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

		log.Printf("testing listAllVirtualNetworkIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		vnIds, err := bpClient.listAllVirtualNetworkIds(ctx)
		if err != nil {
			t.Fatal(err)
		}

		// find two virtual network IDs in the same routing zone
		var src, dst ObjectId
		rzToVnId := make(map[ObjectId]ObjectId)
		for _, vnId := range vnIds {
			log.Printf("testing getVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			vn, err := bpClient.GetVirtualNetwork(ctx, vnId)
			if err != nil {
				t.Fatal(err)
			}
			if dstVN, found := rzToVnId[vn.Data.SecurityZoneId]; !found {
				rzToVnId[vn.Data.SecurityZoneId] = vnId
			} else {
				src = vnId
				dst = dstVN
				break
			}
		}

		randStr := randString(5, "hex")
		label := "label-" + randStr
		description := "description of " + randStr
		log.Printf("testing createPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		policyId, err := bpClient.createPolicy(context.TODO(), &Policy{
			Enabled:             true,
			Label:               label,
			Description:         description,
			SrcApplicationPoint: src,
			DstApplicationPoint: dst,
			Rules:               nil,
			Tags:                nil,
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("created policy id: '%s'", policyId)
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
