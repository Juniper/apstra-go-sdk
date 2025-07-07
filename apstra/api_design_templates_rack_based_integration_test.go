// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
	"log"
	"strings"
	"testing"
)

func TestCreateGetDeleteRackBasedTemplate(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	dn := randString(5, "hex")
	req := CreateRackBasedTemplateRequest{
		DisplayName: dn,
		Spine: &TemplateElementSpineRequest{
			Count:                  2,
			LinkPerSuperspineSpeed: "10G",
			LogicalDevice:          "AOS-7x10-Spine",
			LinkPerSuperspineCount: 1,
			Tags:                   []ObjectId{"firewall", "hypervisor"},
		},
		RackInfos: map[ObjectId]TemplateRackBasedRackInfo{
			"access_switch": {
				Count: 1,
			},
		},
		DhcpServiceIntent:    &DhcpServiceIntent{Active: true},
		AntiAffinityPolicy:   &AntiAffinityPolicy{Algorithm: AlgorithmHeuristic},
		AsnAllocationPolicy:  &AsnAllocationPolicy{SpineAsnScheme: AsnAllocationSchemeSingle},
		VirtualNetworkPolicy: &VirtualNetworkPolicy{},
	}

	for clientName, client := range clients {
		log.Printf("testing CreateRackBasedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.CreateRackBasedTemplate(context.TODO(), &req)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetRackBasedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		rbt, err := client.client.GetRackBasedTemplate(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
		if rbt.Data.DisplayName != dn {
			t.Fatalf("new template displayname mismatch: '%s' vs. '%s'", dn, rbt.Data.DisplayName)
		}

		log.Printf("testing DeleteTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteTemplate(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetRackBasedTemplateByName(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	name := "L2 Pod"

	for _, client := range clients {
		rbt, err := client.client.GetRackBasedTemplateByName(ctx, name)
		if err != nil {
			t.Fatal(err)
		}
		if rbt.templateType.String() != templateTypeRackBased.string() {
			t.Fatalf("expected '%s', got '%s'", rbt.templateType.String(), templateTypeRackBased)
		}
		if rbt.Data.DisplayName != name {
			t.Fatalf("expected '%s', got '%s'", name, rbt.Data.DisplayName)
		}
	}
}

func TestRackBasedTemplateMethods(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	compareSpine := func(req TemplateElementSpineRequest, rbt Spine) error {
		if rbt.Count != req.Count {
			return fmt.Errorf("spine count mismatch: expected %d got %d", rbt.Count, req.Count)
		}

		if rbt.LinkPerSuperspineCount != req.LinkPerSuperspineCount {
			return fmt.Errorf("spine link per superspine count mismatch: expected %d got %d", rbt.LinkPerSuperspineCount, req.LinkPerSuperspineCount)
		}

		if rbt.LinkPerSuperspineSpeed != req.LinkPerSuperspineSpeed {
			return fmt.Errorf("spine link per superspine speed mismatch: expected %q got %q", rbt.LinkPerSuperspineSpeed, req.LinkPerSuperspineSpeed)
		}

		reqTags := make(map[string]bool, len(req.Tags))
		for _, tag := range req.Tags {
			reqTags[strings.ToLower(string(tag))] = true
		}

		rbtTags := make(map[string]bool, len(rbt.Tags))
		for _, tag := range rbt.Tags {
			rbtTags[strings.ToLower(tag.Label)] = true
		}

		if len(reqTags) != len(rbtTags) {
			return fmt.Errorf("tag count mismatch: expected %d got %d", len(reqTags), len(rbtTags))
		}

		for reqTag := range reqTags {
			if !rbtTags[reqTag] {
				return fmt.Errorf("tag mismatch: expected tag %q not found", reqTag)
			}
		}

		return nil
	}

	compareRackInfo := func(req, rbt TemplateRackBasedRackInfo) error {
		if req.Count != rbt.Count {
			return fmt.Errorf("count mismatch: expected %d got %d", req.Count, rbt.Count)
		}

		return nil
	}

	compareRackInfos := func(req, rbt map[ObjectId]TemplateRackBasedRackInfo) error {
		if len(req) != len(rbt) {
			return fmt.Errorf("rack type length mismatch expected %d got %d", len(req), len(rbt))
		}

		for k, reqRI := range req {
			if rbtRI, ok := rbt[k]; ok {
				err = compareRackInfo(reqRI, rbtRI)
				if err != nil {
					return fmt.Errorf("rack infos %q mismatch - %w", k, err)
				}
			} else {
				return fmt.Errorf("rack type mismatch expected rack based info %q not found", k)
			}
		}

		return nil
	}

	compareRequestToTemplate := func(t testing.TB, req CreateRackBasedTemplateRequest, rbt TemplateRackBasedData) error {
		t.Helper()

		if req.DisplayName != rbt.DisplayName {
			return fmt.Errorf("displayname mismatch expected %q got %q", req.DisplayName, rbt.DisplayName)
		}

		err = compareSpine(*req.Spine, rbt.Spine)
		if err != nil {
			return err
		}

		err = compareRackInfos(req.RackInfos, rbt.RackInfo)
		if err != nil {
			return err
		}

		if req.DhcpServiceIntent.Active != rbt.DhcpServiceIntent.Active {
			return fmt.Errorf("dhcp service intent mismatch expected %t got %t", req.DhcpServiceIntent.Active, rbt.DhcpServiceIntent.Active)
		}

		if req.AntiAffinityPolicy != nil {
			require.NotNilf(t, rbt.AntiAffinityPolicy, "rbt.AntiAffinityPolicy is nil")
			compareAntiAffinityPolicy(t, *req.AntiAffinityPolicy, *rbt.AntiAffinityPolicy)
		}

		if req.AsnAllocationPolicy.SpineAsnScheme != rbt.AsnAllocationPolicy.SpineAsnScheme {
			return fmt.Errorf("asn allocation policy spine asn scheme mismatch expected %q got %q", req.AsnAllocationPolicy.SpineAsnScheme, rbt.AsnAllocationPolicy.SpineAsnScheme)
		}

		if req.VirtualNetworkPolicy.OverlayControlProtocol != rbt.VirtualNetworkPolicy.OverlayControlProtocol {
			return fmt.Errorf("virtual network policy overlay control policy mismatch expected %q got %q", req.VirtualNetworkPolicy.OverlayControlProtocol, rbt.VirtualNetworkPolicy.OverlayControlProtocol)
		}

		return nil
	}

	type testCase struct {
		request            CreateRackBasedTemplateRequest
		versionConstraints version.Constraints
	}

	spines := []TemplateElementSpineRequest{
		{
			Count:                  2,
			LinkPerSuperspineSpeed: "10G",
			LogicalDevice:          "AOS-7x10-Spine",
			LinkPerSuperspineCount: 1,
			Tags:                   []ObjectId{"firewall", "hypervisor"},
		},
		{
			Count:                  1,
			LinkPerSuperspineSpeed: "10G",
			LogicalDevice:          "AOS-7x10-Spine",
			LinkPerSuperspineCount: 2,
		},
	}

	rackInfos := []map[ObjectId]TemplateRackBasedRackInfo{
		{"access_switch": {Count: 1}},
		{"access_switch": {Count: 2}},
	}

	testCases := []testCase{
		{
			versionConstraints: compatibility.EqApstra420,
			request: CreateRackBasedTemplateRequest{
				DisplayName:          randString(5, "hex"),
				Spine:                &spines[0],
				RackInfos:            rackInfos[0],
				DhcpServiceIntent:    &DhcpServiceIntent{Active: true},
				AntiAffinityPolicy:   &AntiAffinityPolicy{Algorithm: AlgorithmHeuristic}, // 4.2.0 only?
				AsnAllocationPolicy:  &AsnAllocationPolicy{SpineAsnScheme: AsnAllocationSchemeSingle},
				VirtualNetworkPolicy: &VirtualNetworkPolicy{},
			},
		},
		{
			versionConstraints: compatibility.EqApstra420,
			request: CreateRackBasedTemplateRequest{
				DisplayName:          randString(5, "hex"),
				Spine:                &spines[1],
				RackInfos:            rackInfos[1],
				DhcpServiceIntent:    &DhcpServiceIntent{Active: false},
				AntiAffinityPolicy:   &AntiAffinityPolicy{Algorithm: AlgorithmHeuristic},
				AsnAllocationPolicy:  &AsnAllocationPolicy{SpineAsnScheme: AsnAllocationSchemeSingle},
				VirtualNetworkPolicy: &VirtualNetworkPolicy{},
			},
		},
		{
			versionConstraints: compatibility.GeApstra421,
			request: CreateRackBasedTemplateRequest{
				DisplayName:          randString(5, "hex"),
				Spine:                &spines[0],
				RackInfos:            rackInfos[0],
				DhcpServiceIntent:    &DhcpServiceIntent{Active: true},
				AsnAllocationPolicy:  &AsnAllocationPolicy{SpineAsnScheme: AsnAllocationSchemeSingle},
				VirtualNetworkPolicy: &VirtualNetworkPolicy{},
			},
		},
		{
			request: CreateRackBasedTemplateRequest{
				DisplayName:          randString(5, "hex"),
				Spine:                &spines[1],
				RackInfos:            rackInfos[1],
				DhcpServiceIntent:    &DhcpServiceIntent{Active: false},
				AsnAllocationPolicy:  &AsnAllocationPolicy{SpineAsnScheme: AsnAllocationSchemeSingle},
				VirtualNetworkPolicy: &VirtualNetworkPolicy{},
			},
			versionConstraints: compatibility.GeApstra421,
		},
	}

	for i, tc := range testCases {
		i, tc := i, tc

		t.Run(fmt.Sprintf("Test-%d", i), func(t *testing.T) {
			for clientName, client := range clients {
				clientName, client := clientName, client

				t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
					if !tc.versionConstraints.Check(client.client.apiVersion) {
						t.Skipf("skipping testcase %d because of versionConstraint %s, version %s", i, tc.versionConstraints, client.client.apiVersion)
					}

					log.Printf("testing CreateRackBasedTemplate(testCase[%d]) against %s %s (%s)", i, client.clientType, clientName, client.client.ApiVersion())
					id, err := client.client.CreateRackBasedTemplate(ctx, &tc.request)
					if err != nil {
						t.Fatal(err)
					}

					log.Printf("testing GetRackBasedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					rbt, err := client.client.GetRackBasedTemplate(ctx, id)
					if err != nil {
						t.Fatal(err)
					}

					if id != rbt.Id {
						t.Fatalf("test case %d template id mismatch expected %q got %q", i, id, rbt.Id)
					}

					err = compareRequestToTemplate(t, tc.request, *rbt.Data)
					if err != nil {
						t.Fatalf("test case %d template differed from request: %s", i, err.Error())
					}

					for j := i; j < i+len(testCases); j++ { // j counts up from i
						k := j % len(testCases) // k counts up from i, but loops back to zero

						if !testCases[k].versionConstraints.Check(client.client.apiVersion) {
							continue
						}

						req := testCases[k].request
						log.Printf("testing UpdateRackBasedTemplate(testCase[%d]) against %s %s (%s)", k, client.clientType, clientName, client.client.ApiVersion())
						err = client.client.UpdateRackBasedTemplate(ctx, id, &req)
						if err != nil {
							t.Fatal(err)
						}

						log.Printf("testing GetRackBasedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
						rbt, err = client.client.GetRackBasedTemplate(ctx, id)
						if err != nil {
							t.Fatal(err)
						}

						if id != rbt.Id {
							t.Fatalf("test case %d template id mismatch expected %q got %q", i, id, rbt.Id)
						}

						err = compareRequestToTemplate(t, req, *rbt.Data)
						if err != nil {
							t.Fatalf("test case %d template differed from request: %s", i, err.Error())
						}
					}

					log.Printf("testing DeleteTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					err = client.client.DeleteTemplate(ctx, id)
					if err != nil {
						t.Fatal(err)
					}
				})
			}
		})
	}
}
