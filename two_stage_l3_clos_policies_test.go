package goapstra

import (
	"context"
	"log"
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
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpIds, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(bpIds) == 0 {
			t.Skip()
		}

		for _, bpId := range bpIds {
			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			dcClient, err := client.client.NewTwoStageL3ClosClient(context.TODO(), bpId)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing getAllPolicies() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			policies, err := dcClient.getAllPolicies(context.TODO())
			if err != nil {
				t.Fatal(err)
			}

			for _, policy := range policies {
				log.Printf("testing getPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				p, err := dcClient.getPolicy(context.TODO(), policy.Id)
				if err != nil {
					t.Fatal(err)
				}
				log.Printf("policy '%s'\t'%s'", p.Id, p.Label)
			}
		}
	}
}

func TestCreatePolicy(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpIds, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(bpIds) == 0 {
			t.Skip()
		}

		bpId := bpIds[0]

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		dcClient, err := client.client.NewTwoStageL3ClosClient(context.TODO(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing listAllVirtualNetworkIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		vnIds, err := dcClient.listAllVirtualNetworkIds(context.TODO(), BlueprintTypeStaging)
		if err != nil {
			t.Fatal(err)
		}

		// find two virtual network IDs in the same routing zone
		var src, dst ObjectId
		rzToVnId := make(map[ObjectId]ObjectId)
		for _, vnId := range vnIds {
			log.Printf("testing getVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			vn, err := dcClient.getVirtualNetwork(context.TODO(), vnId, BlueprintTypeStaging)
			if err != nil {
				t.Fatal(err)
			}
			if dstVN, found := rzToVnId[vn.SecurityZoneId]; !found {
				rzToVnId[vn.SecurityZoneId] = vnId
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
		policyId, err := dcClient.createPolicy(context.TODO(), &Policy{
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
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpIds, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(bpIds) == 0 {
			t.Skip()
		}

		bpId := bpIds[0]

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		dcClient, err := client.client.NewTwoStageL3ClosClient(context.TODO(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		randStr := randString(5, "hex")
		label := "new-" + randStr
		log.Printf("adding new rule '%s'", label)

		newRule := &PolicyRule{
			Label:       label,
			Description: "new rule " + randStr,
			Protocol:    "TCP",
			Action:      PolicyRuleActionDenyLog,
			SrcPort:     PortRanges{{5, 6}},
			DstPort:     PortRanges{{7, 8}, {9, 10}},
		}

		policyId := ObjectId("lkmWBn_wM9ShK9VCCQk")

		log.Printf("testing addPolicyRule() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ruleId, err := dcClient.addPolicyRule(context.TODO(), newRule, 0, policyId)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("new rule id: '%s'", ruleId)
		log.Printf("testing deletePolicyRuleById() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = dcClient.deletePolicyRuleById(context.TODO(), policyId, ruleId)
		if err != nil {
			t.Fatal(err)
		}
	}
}
