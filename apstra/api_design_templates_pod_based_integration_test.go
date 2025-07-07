// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestCreateGetDeletePodBasedTemplate(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	dn := randString(5, "hex")

	rbtdn := "rbtr-" + dn
	rbtr := CreateRackBasedTemplateRequest{
		DisplayName: rbtdn,
		Spine: &TemplateElementSpineRequest{
			Count:                  2,
			LinkPerSuperspineSpeed: "10G",
			LogicalDevice:          "AOS-7x10-Spine",
			LinkPerSuperspineCount: 1,
			Tags:                   nil,
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
		rbtid, err := client.client.CreateRackBasedTemplate(ctx, &rbtr)
		if err != nil {
			t.Fatal(err)
		}

		pbtdn := "pbtr-" + dn
		pbtr := CreatePodBasedTemplateRequest{
			DisplayName: pbtdn,
			Superspine: &TemplateElementSuperspineRequest{
				PlaneCount:         1,
				Tags:               nil,
				SuperspinePerPlane: 4,
				LogicalDeviceId:    "AOS-4x40_8x10-1",
			},
			PodInfos: map[ObjectId]TemplatePodBasedInfo{
				rbtid: {
					Count: 1,
				},
			},
			AntiAffinityPolicy: &AntiAffinityPolicy{
				Algorithm:                AlgorithmHeuristic,
				MaxLinksPerPort:          1,
				MaxLinksPerSlot:          1,
				MaxPerSystemLinksPerPort: 1,
				MaxPerSystemLinksPerSlot: 1,
				Mode:                     AntiAffinityModeDisabled,
			},
		}

		log.Printf("testing CreatePodBasedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pbtid, err := client.client.CreatePodBasedTemplate(ctx, &pbtr)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetPodBasedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pbt, err := client.client.GetPodBasedTemplate(ctx, pbtid)
		if err != nil {
			t.Fatal(err)
		}

		if pbt.Data.DisplayName != pbtdn {
			t.Fatalf("new template displayname mismatch: '%s' vs. '%s'", dn, pbt.Data.DisplayName)
		}

		log.Printf("testing DeleteTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteTemplate(ctx, pbtid)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing DeleteTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteTemplate(ctx, rbtid)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetPodBasedTemplateByName(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	name := "L2 superspine single plane"

	for _, client := range clients {
		pbt, err := client.client.GetPodBasedTemplateByName(ctx, name)
		if err != nil {
			t.Fatal(err)
		}
		if pbt.templateType.String() != templateTypePodBased.string() {
			t.Fatalf("expected '%s', got '%s'", pbt.templateType.String(), templateTypePodBased)
		}
		if pbt.Data.DisplayName != name {
			t.Fatalf("expected '%s', got '%s'", name, pbt.Data.DisplayName)
		}
	}
}
