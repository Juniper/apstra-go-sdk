// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalTemplate(t *testing.T) {
	data := `{
  "anti_affinity_policy": {
    "max_links_per_port": 0,
    "algorithm": "heuristic",
    "max_per_system_links_per_port": 0,
    "max_links_per_slot": 0,
    "max_per_system_links_per_slot": 0,
    "mode": "disabled"
  },
  "display_name": "L2 Virtual Dual MLAG",
  "virtual_network_policy": {
    "overlay_control_protocol": null
  },
  "fabric_addressing_policy": {
    "spine_leaf_links": "ipv4",
    "spine_superspine_links": "ipv4"
  },
  "spine": {
    "count": 2,
    "link_per_superspine_count": 0,
    "tags": [],
    "logical_device": {
      "panels": [
        {
          "panel_layout": {
            "row_count": 2,
            "column_count": 16
          },
          "port_indexing": {
            "order": "T-B, L-R",
            "start_index": 1,
            "schema": "absolute"
          },
          "port_groups": [
            {
              "count": 24,
              "speed": {
                "unit": "G",
                "value": 10
              },
              "roles": [
                "superspine",
                "leaf"
              ]
            },
            {
              "count": 8,
              "speed": {
                "unit": "G",
                "value": 10
              },
              "roles": [
                "generic"
              ]
            }
          ]
        }
      ],
      "display_name": "AOS-32x10-Spine",
      "id": "AOS-32x10-Spine"
    },
    "link_per_superspine_speed": null
  },
  "created_at": "2022-04-22T06:08:57.993697Z",
  "rack_type_counts": [
    {
      "rack_type_id": "L2_Virtual_Dual_MLAG",
      "count": 2
    }
  ],
  "dhcp_service_intent": {
    "active": true
  },
  "last_modified_at": "2022-04-22T06:08:57.993697Z",
  "rack_types": [
    {
      "description": "",
      "tags": [],
      "logical_devices": [
        {
          "panels": [
            {
              "panel_layout": {
                "row_count": 1,
                "column_count": 7
              },
              "port_indexing": {
                "order": "T-B, L-R",
                "start_index": 1,
                "schema": "absolute"
              },
              "port_groups": [
                {
                  "count": 2,
                  "speed": {
                    "unit": "G",
                    "value": 10
                  },
                  "roles": [
                    "leaf",
                    "spine"
                  ]
                },
                {
                  "count": 2,
                  "speed": {
                    "unit": "G",
                    "value": 10
                  },
                  "roles": [
                    "peer"
                  ]
                },
                {
                  "count": 2,
                  "speed": {
                    "unit": "G",
                    "value": 10
                  },
                  "roles": [
                    "generic",
                    "access"
                  ]
                },
                {
                  "count": 1,
                  "speed": {
                    "unit": "G",
                    "value": 10
                  },
                  "roles": [
                    "generic"
                  ]
                }
              ]
            }
          ],
          "display_name": "AOS-7x10-Leaf",
          "id": "AOS-7x10-Leaf"
        },
        {
          "panels": [
            {
              "panel_layout": {
                "row_count": 2,
                "column_count": 4
              },
              "port_indexing": {
                "order": "T-B, L-R",
                "start_index": 1,
                "schema": "absolute"
              },
              "port_groups": [
                {
                  "count": 8,
                  "speed": {
                    "unit": "G",
                    "value": 10
                  },
                  "roles": [
                    "leaf",
                    "generic",
                    "peer",
                    "access"
                  ]
                }
              ]
            }
          ],
          "display_name": "AOS-8x10-1",
          "id": "AOS-8x10-1"
        }
      ],
      "generic_systems": [
        {
          "count": 1,
          "asn_domain": "disabled",
          "links": [
            {
              "tags": [],
              "link_per_switch_count": 2,
              "label": "link1",
              "link_speed": {
                "unit": "G",
                "value": 10
              },
              "target_switch_label": "leaf_pair_1",
              "attachment_type": "dualAttached",
              "lag_mode": "lacp_active"
            },
            {
              "tags": [],
              "link_per_switch_count": 2,
              "label": "link2",
              "link_speed": {
                "unit": "G",
                "value": 10
              },
              "target_switch_label": "leaf_pair_2",
              "attachment_type": "dualAttached",
              "lag_mode": "lacp_active"
            }
          ],
          "management_level": "unmanaged",
          "port_channel_id_min": 0,
          "port_channel_id_max": 0,
          "logical_device": "AOS-8x10-1",
          "loopback": "disabled",
          "tags": [],
          "label": "generic"
        }
      ],
      "servers": [],
      "leafs": [
        {
          "leaf_leaf_l3_link_speed": null,
          "redundancy_protocol": "mlag",
          "leaf_leaf_link_port_channel_id": 0,
          "leaf_leaf_l3_link_count": 0,
          "logical_device": "AOS-7x10-Leaf",
          "leaf_leaf_link_speed": {
            "unit": "G",
            "value": 10
          },
          "link_per_spine_count": 1,
          "leaf_leaf_link_count": 2,
          "tags": [],
          "link_per_spine_speed": {
            "unit": "G",
            "value": 10
          },
          "label": "leaf_pair_1",
          "mlag_vlan_id": 2999,
          "leaf_leaf_l3_link_port_channel_id": 0
        },
        {
          "leaf_leaf_l3_link_speed": null,
          "redundancy_protocol": "mlag",
          "leaf_leaf_link_port_channel_id": 0,
          "leaf_leaf_l3_link_count": 0,
          "logical_device": "AOS-7x10-Leaf",
          "leaf_leaf_link_speed": {
            "unit": "G",
            "value": 10
          },
          "link_per_spine_count": 1,
          "leaf_leaf_link_count": 2,
          "tags": [],
          "link_per_spine_speed": {
            "unit": "G",
            "value": 10
          },
          "label": "leaf_pair_2",
          "mlag_vlan_id": 2999,
          "leaf_leaf_l3_link_port_channel_id": 0
        }
      ],
      "access_switches": [],
      "id": "L2_Virtual_Dual_MLAG",
      "display_name": "L2 Virtual 2xMLAG",
      "fabric_connectivity_design": "l3clos",
      "created_at": "1970-01-01T00:00:00.000000Z",
      "last_modified_at": "1970-01-01T00:00:00.000000Z"
    }
  ],
  "capability": "blueprint",
  "asn_allocation_policy": {
    "spine_asn_scheme": "distinct"
  },
  "type": "rack_based",
  "id": "L2_Virtual_Dual_MLAG"
}`
	raw := &json.RawMessage{}
	err := json.Unmarshal([]byte(data), raw)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := template(*raw)
	tType, err := tmpl.templateType()
	if err != nil {
		t.Fatal(err)
	}

	log.Println(tType)
}

func TestTemplateStrings(t *testing.T) {
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
		{stringVal: "heuristic", intType: AlgorithmHeuristic, stringType: algorithmHeuristic},

		{stringVal: "", intType: TemplateTypeNone, stringType: templateTypeNone},
		{stringVal: "rack_based", intType: TemplateTypeRackBased, stringType: templateTypeRackBased},
		{stringVal: "pod_based", intType: TemplateTypePodBased, stringType: templateTypePodBased},
		{stringVal: "l3_collapsed", intType: TemplateTypeL3Collapsed, stringType: templateTypeL3Collapsed},

		{stringVal: "distinct", intType: AsnAllocationSchemeDistinct, stringType: asnAllocationSchemeDistinct},
		{stringVal: "single", intType: AsnAllocationSchemeSingle, stringType: asnAllocationSchemeSingle},

		{stringVal: "ipv4", intType: AddressingSchemeIp4, stringType: addressingSchemeIp4},
		{stringVal: "ipv6", intType: AddressingSchemeIp6, stringType: addressingSchemeIp6},
		{stringVal: "ipv4_ipv6", intType: AddressingSchemeIp46, stringType: addressingSchemeIp46},

		{stringVal: "", intType: OverlayControlProtocolNone, stringType: overlayControlProtocolNone},
		{stringVal: "evpn", intType: OverlayControlProtocolEvpn, stringType: overlayControlProtocolEvpn},

		{stringVal: "", intType: TemplateCapabilityNone, stringType: templateCapabilityNone},
		{stringVal: "blueprint", intType: TemplateCapabilityBlueprint, stringType: templateCapabilityBlueprint},
		{stringVal: "pod", intType: TemplateCapabilityPod, stringType: templateCapabilityPod},
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

func TestGetTemplate(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listAllTemplateIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		templateIds, err := client.client.listAllTemplateIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("fetching %d templateIds...", len(templateIds))

		for _, i := range sampleIndexes(t, len(templateIds)) {
			templateId := templateIds[i]
			log.Printf("testing getTemplate(%s) against %s %s (%s)", templateId, client.clientType, clientName, client.client.ApiVersion())
			x, err := client.client.getTemplate(context.TODO(), templateId)
			if err != nil {
				t.Fatal(err)
			}

			var name string
			tType, err := x.templateType()
			if err != nil {
				t.Fatal(err)
			}
			switch tType {
			case templateTypeRackBased:
				rbt := &rawTemplateRackBased{}
				err = json.Unmarshal(x, rbt)
				if err != nil {
					t.Fatal(err)
				}
				name = rbt.DisplayName
				rbt2, err := client.client.GetRackBasedTemplate(context.TODO(), templateId)
				if err != nil {
					t.Fatal(err)
				}
				if templateId != rbt2.Id {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", templateId, rbt2.Id)
				}
				if name != rbt2.Data.DisplayName {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", name, rbt2.Data.DisplayName)
				}
			case templateTypePodBased:
				rbt := &rawTemplatePodBased{}
				err = json.Unmarshal(x, rbt)
				if err != nil {
					t.Fatal(err)
				}
				name = rbt.DisplayName
				rbt2, err := client.client.GetPodBasedTemplate(context.TODO(), templateId)
				if err != nil {
					t.Fatal(err)
				}
				if templateId != rbt2.Id {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", templateId, rbt2.Id)
				}
				if name != rbt2.Data.DisplayName {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", name, rbt2.Data.DisplayName)
				}
			case templateTypeL3Collapsed:
				rbt := &rawTemplateL3Collapsed{}
				err = json.Unmarshal(x, rbt)
				if err != nil {
					t.Fatal(err)
				}
				name = rbt.DisplayName
				rbt2, err := client.client.GetL3CollapsedTemplate(context.TODO(), templateId)
				if err != nil {
					t.Fatal(err)
				}
				if templateId != rbt2.Id {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", templateId, rbt2.Id)
				}
				if name != rbt2.Data.DisplayName {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", name, rbt2.Data.DisplayName)
				}
			}
			log.Printf("template '%s' '%s'", templateId, name)
		}
	}
}

func TestGetTemplateMethods(t *testing.T) {
	var n int

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getAllTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		templates, err := client.client.getAllTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("got %d templates", len(templates))

		// rack-based templates
		log.Printf("testing getAllRackBasedTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		rackBasedTemplates, err := client.client.getAllRackBasedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got %d rack-based templates\n", len(rackBasedTemplates))

		n = rand.Intn(len(rackBasedTemplates))
		log.Printf("testing getRackBasedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(rackBasedTemplates))
		rackBasedTemplate, err := client.client.getRackBasedTemplate(context.TODO(), rackBasedTemplates[n].Id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got template type '%s', ID '%s'\n", rackBasedTemplate.Type, rackBasedTemplate.Id)

		// pod-based templates
		log.Printf("testing getAllPodBasedTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		podBasedTemplates, err := client.client.getAllPodBasedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got %d pod-based templates\n", len(podBasedTemplates))

		n = rand.Intn(len(podBasedTemplates))
		log.Printf("testing getPodBasedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(podBasedTemplates))
		podBasedTemplate, err := client.client.getPodBasedTemplate(context.TODO(), podBasedTemplates[n].Id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got template type '%s', ID '%s'\n", podBasedTemplate.Type, podBasedTemplate.Id)

		// l3-collapsed templates
		log.Printf("testing getAllL3CollapsedTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		l3CollapsedTemplates, err := client.client.getAllL3CollapsedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("  got %d pod-based templates\n", len(l3CollapsedTemplates))

		n = rand.Intn(len(l3CollapsedTemplates))
		log.Printf("testing getL3CollapsedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(l3CollapsedTemplates))
		l3CollapsedTemplate, err := client.client.getL3CollapsedTemplate(context.TODO(), l3CollapsedTemplates[n].Id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got template type '%s', ID '%s'\n", l3CollapsedTemplate.Type, l3CollapsedTemplate.Id)
	}
}

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

func TestCreateGetDeleteL3CollapsedTemplate(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	dn := randString(5, "hex")

	req := &CreateL3CollapsedTemplateRequest{
		DisplayName:   dn,
		MeshLinkCount: 1,
		MeshLinkSpeed: "10G",
		RackTypeIds:   []ObjectId{"L3_collapsed_acs"},
		RackTypeCounts: []RackTypeCount{{
			RackTypeId: "L3_collapsed_acs",
			Count:      1,
		}},
		VirtualNetworkPolicy: VirtualNetworkPolicy{OverlayControlProtocol: OverlayControlProtocolEvpn},
	}

	for clientName, client := range clients {
		log.Printf("testing CreateL3CollapsedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.CreateL3CollapsedTemplate(ctx, req)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing DeleteTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteTemplate(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetL3CollapsedTemplateByName(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	name := "Collapsed Fabric ESI"

	for _, client := range clients {
		l3ct, err := client.client.GetL3CollapsedTemplateByName(ctx, name)
		if err != nil {
			t.Fatal(err)
		}
		if l3ct.templateType.String() != templateTypeL3Collapsed.string() {
			t.Fatalf("expected '%s', got '%s'", l3ct.templateType.String(), templateTypeL3Collapsed)
		}
		if l3ct.Data.DisplayName != name {
			t.Fatalf("expected '%s', got '%s'", name, l3ct.Data.DisplayName)
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

func TestGetTemplateType(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	type testData struct {
		templateId   ObjectId
		templateType templateType
	}

	data := []testData{
		{"pod1", templateTypeRackBased},
		{"L2_superspine_multi_plane", templateTypePodBased},
		{"L3_Collapsed_ACS", templateTypeL3Collapsed},
	}

	for clientName, client := range clients {
		for _, d := range data {
			log.Printf("testing getTemplateType(%s) against %s %s (%s)", d.templateType, client.clientType, clientName, client.client.ApiVersion())
			ttype, err := client.client.getTemplateType(ctx, d.templateId)
			if err != nil {
				t.Fatal(err)
			}
			if ttype != d.templateType {
				t.Fatalf("expected '%s', got '%s'", ttype.string(), d.templateType)
			}
		}
	}
}

func TestGetTemplateIdsTypesByName(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	templateName := randString(10, "hex")
	for clientName, client := range clients {
		// fetch all template IDs
		templateIds, err := client.client.listAllTemplateIds(ctx)
		if err != nil {
			t.Fatal(err)
		}

		// choose a random template for cloning
		cloneMeId := templateIds[rand.Intn(len(templateIds))]
		cloneMeType, err := client.client.getTemplateType(ctx, cloneMeId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("cloning template '%s' (%s) for this test", cloneMeId, cloneMeType)

		cloneCount := rand.Intn(5) + 2
		cloneIds := make([]ObjectId, cloneCount)
		for i := 0; i < cloneCount; i++ {
			switch cloneMeType {
			case templateTypeRackBased:
				cloneMe, err := client.client.getRackBasedTemplate(ctx, cloneMeId)
				if err != nil {
					t.Fatal(err)
				}
				id, err := client.client.createRackBasedTemplate(ctx, &rawCreateRackBasedTemplateRequest{
					Type:                 cloneMe.Type,
					DisplayName:          fmt.Sprintf("%s-%d", templateName, i),
					Spine:                cloneMe.Spine,
					RackTypes:            cloneMe.RackTypes,
					RackTypeCounts:       cloneMe.RackTypeCounts,
					DhcpServiceIntent:    cloneMe.DhcpServiceIntent,
					AntiAffinityPolicy:   cloneMe.AntiAffinityPolicy,
					AsnAllocationPolicy:  cloneMe.AsnAllocationPolicy,
					VirtualNetworkPolicy: cloneMe.VirtualNetworkPolicy,
				})
				if err != nil {
					t.Fatal(err)
				}
				cloneIds[i] = id
			case templateTypePodBased:
				cloneMe, err := client.client.getPodBasedTemplate(ctx, cloneMeId)
				if err != nil {
					t.Fatal(err)
				}
				id, err := client.client.createPodBasedTemplate(ctx, &rawCreatePodBasedTemplateRequest{
					Type:                    cloneMe.Type,
					DisplayName:             fmt.Sprintf("%s-%d", templateName, i),
					Superspine:              cloneMe.Superspine,
					RackBasedTemplates:      cloneMe.RackBasedTemplates,
					RackBasedTemplateCounts: cloneMe.RackBasedTemplateCounts,
					AntiAffinityPolicy:      cloneMe.AntiAffinityPolicy,
				})
				if err != nil {
					t.Fatal(err)
				}
				cloneIds[i] = id
			case templateTypeL3Collapsed:
				cloneMe, err := client.client.getL3CollapsedTemplate(ctx, cloneMeId)
				if err != nil {
					t.Fatal(err)
				}
				id, err := client.client.createL3CollapsedTemplate(ctx, &rawCreateL3CollapsedTemplateRequest{
					Type:                 cloneMe.Type,
					DisplayName:          fmt.Sprintf("%s-%d", templateName, i),
					MeshLinkCount:        cloneMe.MeshLinkCount,
					MeshLinkSpeed:        *cloneMe.MeshLinkSpeed,
					RackTypes:            cloneMe.RackTypes,
					RackTypeCounts:       cloneMe.RackTypeCounts,
					DhcpServiceIntent:    cloneMe.DhcpServiceIntent,
					AntiAffinityPolicy:   cloneMe.AntiAffinityPolicy,
					VirtualNetworkPolicy: cloneMe.VirtualNetworkPolicy,
				})
				if err != nil {
					t.Fatal(err)
				}
				cloneIds[i] = id
			}
		}
		clones := make([]string, len(cloneIds))
		for i, clone := range cloneIds {
			clones[i] = string(clone)
		}
		log.Printf("clone IDs: '%s'", strings.Join(clones, ", "))

		templateIdsToType := make(map[ObjectId]TemplateType)
		for i := 0; i < cloneCount; i++ {
			log.Printf("testing getTemplateIdsTypesByName(%s) against %s %s (%s)", templateName, client.clientType, clientName, client.client.ApiVersion())
			temp, err := client.client.getTemplateIdsTypesByName(ctx, fmt.Sprintf("%s-%d", templateName, i))
			if err != nil {
				t.Fatal(err)
			}
			for k, v := range temp {
				templateIdsToType[k] = v
			}
		}

		if cloneCount != len(templateIdsToType) {
			t.Fatalf("expected %d, got %d", cloneCount, len(templateIdsToType))
		}
		for _, v := range templateIdsToType {
			parsed, err := cloneMeType.parse()
			if err != nil {
				t.Fatal(err)
			}
			if parsed != v.Int() {
				t.Fatalf("expected %d, got %d", parsed, v.Int())
			}
		}

		for i, cloneId := range cloneIds {
			name := fmt.Sprintf("%s-%d", templateName, i)
			if i+1 == len(cloneIds) { // last one before they're all deleted
				log.Printf("testing getTemplateIdTypeByName(%s) against %s %s (%s)", name, client.clientType, clientName, client.client.ApiVersion())
				tId, tType, err := client.client.getTemplateIdTypeByName(ctx, name)
				if err != nil {
					t.Fatal(err)
				}
				if cloneId != tId {
					t.Fatalf("expected template id '%s', got '%s'", cloneIds, tId)
				}
				if cloneMeType != templateType(tType.String()) {
					t.Fatalf("expected template type '%s', got '%s'", cloneMeType, tType.String())
				}

			}
			log.Printf("deleting clone '%s'", cloneId)
			err = client.client.deleteTemplate(ctx, cloneId)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestAllTemplateTypes(t *testing.T) {
	all := AllTemplateTypes()
	expected := 3
	if len(all) != expected {
		t.Fatalf("expected %d template types, got %d", expected, len(all))
	}
}

func TestAllOverlayControlProtocols(t *testing.T) {
	all := AllOverlayControlProtocols()
	expected := 2
	if len(all) != expected {
		log.Println(all)
		t.Fatalf("expected %d overlay control protocols, got %d", expected, len(all))
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
