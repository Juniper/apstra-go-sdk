package goapstra

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"testing"
	"time"
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
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing listAllTemplateIds() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		templateIds, err := client.client.listAllTemplateIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("fetching %d templateIds...", len(templateIds))

		for _, id := range templateIds {
			log.Printf("testing getTemplate(%s) against %s %s (%s)", id, client.clientType, client.clientName, client.client.ApiVersion())
			x, err := client.client.getTemplate(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}

			var id ObjectId
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
				id = rbt.Id
				name = rbt.DisplayName
				rbt2, err := client.client.GetRackBasedTemplate(context.TODO(), id)
				if err != nil {
					log.Fatal(err)
				}
				if id != rbt2.Id {
					log.Fatalf("template ID mismatch: '%s' vs. '%s", id, rbt2.Id)
				}
				if name != rbt2.DisplayName {
					log.Fatalf("template ID mismatch: '%s' vs. '%s", name, rbt2.DisplayName)
				}
			case templateTypePodBased:
				rbt := &rawTemplatePodBased{}
				err = json.Unmarshal(x, rbt)
				if err != nil {
					t.Fatal(err)
				}
				id = rbt.Id
				name = rbt.DisplayName
				rbt2, err := client.client.GetPodBasedTemplate(context.TODO(), id)
				if err != nil {
					log.Fatal(err)
				}
				if id != rbt2.Id {
					log.Fatalf("template ID mismatch: '%s' vs. '%s", id, rbt2.Id)
				}
				if name != rbt2.DisplayName {
					log.Fatalf("template ID mismatch: '%s' vs. '%s", name, rbt2.DisplayName)
				}
			case templateTypeL3Collapsed:
				rbt := &rawTemplateL3Collapsed{}
				err = json.Unmarshal(x, rbt)
				if err != nil {
					t.Fatal(err)
				}
				id = rbt.Id
				name = rbt.DisplayName
				rbt2, err := client.client.GetL3CollapsedTemplate(context.TODO(), id)
				if err != nil {
					log.Fatal(err)
				}
				if id != rbt2.Id {
					log.Fatalf("template ID mismatch: '%s' vs. '%s", id, rbt2.Id)
				}
				if name != rbt2.DisplayName {
					log.Fatalf("template ID mismatch: '%s' vs. '%s", name, rbt2.DisplayName)
				}
			}
			log.Printf("template '%s' '%s'", id, name)
		}
	}
}

func TestGetTemplateMethods(t *testing.T) {
	var n int
	rand.Seed(time.Now().UnixNano())

	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing getAllTemplates() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		templates, err := client.client.getAllTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("got %d templates", len(templates))

		// rack-based templates
		log.Printf("testing getAllRackBasedTemplates() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		rackBasedTemplates, err := client.client.getAllRackBasedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got %d rack-based templates\n", len(rackBasedTemplates))

		n = rand.Intn(len(rackBasedTemplates))
		log.Printf("testing getRackBasedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(rackBasedTemplates))
		rackBasedTemplate, err := client.client.getRackBasedTemplate(context.TODO(), rackBasedTemplates[n].Id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got template type '%s', ID '%s'\n", rackBasedTemplate.Type, rackBasedTemplate.Id)

		// pod-based templates
		log.Printf("testing getAllPodBasedTemplates() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		podBasedTemplates, err := client.client.getAllPodBasedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got %d pod-based templates\n", len(podBasedTemplates))

		n = rand.Intn(len(podBasedTemplates))
		log.Printf("testing getPodBasedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(podBasedTemplates))
		podBasedTemplate, err := client.client.getPodBasedTemplate(context.TODO(), podBasedTemplates[n].Id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got template type '%s', ID '%s'\n", podBasedTemplate.Type, podBasedTemplate.Id)

		// l3-collapsed templates
		log.Printf("testing getAllL3CollapsedTemplates() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		l3CollapsedTemplates, err := client.client.getAllL3CollapsedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("  got %d pod-based templates\n", len(l3CollapsedTemplates))

		n = rand.Intn(len(l3CollapsedTemplates))
		log.Printf("testing getL3CollapsedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(l3CollapsedTemplates))
		l3CollapsedTemplate, err := client.client.getL3CollapsedTemplate(context.TODO(), l3CollapsedTemplates[n].Id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got template type '%s', ID '%s'\n", l3CollapsedTemplate.Type, l3CollapsedTemplate.Id)
	}
}

func TestCreateGetDeleteRackBasedTemplate(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	dn := randString(5, "hex")
	req := CreateRackBasedTemplateRequest{
		DisplayName: dn,
		Capability:  TemplateCapabilityBlueprint,
		Spine: TemplateElementSpineRequest{
			Count:                   2,
			ExternalLinkSpeed:       "10G",
			LinkPerSuperspineSpeed:  "10G",
			LogicalDevice:           "AOS-7x10-Spine",
			LinkPerSuperspineCount:  1,
			Tags:                    nil,
			ExternalLinksPerNode:    0,
			ExternalFacingNodeCount: 0,
			ExternalLinkCount:       0,
		},
		RackTypeIds: []ObjectId{"evpn-single"},
		RackTypeCounts: []RackTypeCounts{{
			Count:      1,
			RackTypeId: "evpn-single",
		}},
		DhcpServiceIntent:      DhcpServiceIntent{Active: true},
		AntiAffinityPolicy:     AntiAffinityPolicy{Algorithm: AlgorithmHeuristic},
		AsnAllocationPolicy:    AsnAllocationPolicy{SpineAsnScheme: AsnAllocationSchemeSingle},
		FabricAddressingPolicy: FabricAddressingPolicy{},
		VirtualNetworkPolicy:   VirtualNetworkPolicy{},
	}

	for _, client := range clients {
		log.Printf("testing createRackBasedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		id, err := client.client.createRackBasedTemplate(context.TODO(), &req)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing createRackBasedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		template, err := client.client.GetRackBasedTemplate(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
		if template.DisplayName != dn {
			t.Fatalf("new template displayname mismatch: '%s' vs. '%s'", dn, template.DisplayName)
		}

		log.Printf("testing deleteTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		err = client.client.deleteTemplate(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCreateGetDeletePodBasedTemplate(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	dn := randString(5, "hex")

	rbtdn := "rbtr-" + dn
	rbtr := CreateRackBasedTemplateRequest{
		DisplayName: rbtdn,
		Capability:  TemplateCapabilityBlueprint,
		Spine: TemplateElementSpineRequest{
			Count:                   2,
			LinkPerSuperspineSpeed:  "10G",
			LogicalDevice:           "AOS-7x10-Spine",
			LinkPerSuperspineCount:  1,
			Tags:                    nil,
			ExternalLinksPerNode:    0,
			ExternalFacingNodeCount: 0,
			ExternalLinkCount:       0,
		},
		RackTypeIds: []ObjectId{"evpn-single"},
		RackTypeCounts: []RackTypeCounts{{
			Count:      1,
			RackTypeId: "evpn-single",
		}},
		DhcpServiceIntent:      DhcpServiceIntent{Active: true},
		AntiAffinityPolicy:     AntiAffinityPolicy{Algorithm: AlgorithmHeuristic},
		AsnAllocationPolicy:    AsnAllocationPolicy{SpineAsnScheme: AsnAllocationSchemeSingle},
		FabricAddressingPolicy: FabricAddressingPolicy{},
		VirtualNetworkPolicy:   VirtualNetworkPolicy{},
	}

	for _, client := range clients {
		log.Printf("testing createRackBasedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		rbtid, err := client.client.createRackBasedTemplate(context.TODO(), &rbtr)
		if err != nil {
			t.Fatal(err)
		}

		pbtdn := "pbtr-" + dn
		pbtr := CreatePodBasedTemplateRequest{
			DisplayName: pbtdn,
			Capability:  TemplateCapabilityPod,
			Superspine: TemplateElementSuperspineRequest{
				PlaneCount:         1,
				ExternalLinkCount:  0,
				ExternalLinkSpeed:  "10G",
				Tags:               nil,
				SuperspinePerPlane: 4,
				LogicalDeviceId:    "AOS-4x40_8x10-1",
			},
			RackBasedTemplateIds: []ObjectId{rbtid},
			RackBasedTemplateCounts: []RackBasedTemplateCounts{{
				RackBasedTemplateId: rbtid,
				Count:               1,
			}},
			AntiAffinityPolicy: AntiAffinityPolicy{
				Algorithm:                AlgorithmHeuristic,
				MaxLinksPerPort:          1,
				MaxLinksPerSlot:          1,
				MaxPerSystemLinksPerPort: 1,
				MaxPerSystemLinksPerSlot: 1,
				Mode:                     AntiAffinityModeDisabled,
			},
			FabricAddressingPolicy: FabricAddressingPolicy{
				SpineSuperspineLinks: AddressingSchemeIp4,
				SpineLeafLinks:       AddressingSchemeIp4,
			},
		}

		log.Printf("testing createPodBasedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		pbtid, err := client.client.createPodBasedTemplate(context.TODO(), &pbtr)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetPodBasedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		pbt, err := client.client.GetPodBasedTemplate(context.TODO(), pbtid)
		if err != nil {
			t.Fatal(err)
		}

		if pbt.DisplayName != pbtdn {
			t.Fatalf("new template displayname mismatch: '%s' vs. '%s'", dn, pbt.DisplayName)
		}

		log.Printf("testing deleteTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		err = client.client.deleteTemplate(context.TODO(), pbtid)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing deleteTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		err = client.client.deleteTemplate(context.TODO(), rbtid)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCreateGetDeleteL3CollapsedTemplate(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	dn := randString(5, "hex")

	req := &CreateL3CollapsedTemplateRequest{
		DisplayName:   dn,
		Capability:    TemplateCapabilityBlueprint,
		MeshLinkCount: 1,
		MeshLinkSpeed: "10G",
		RackTypeIds:   []ObjectId{"L3_collapsed_acs"},
		RackTypeCounts: []RackTypeCounts{{
			RackTypeId: "L3_collapsed_acs",
			Count:      1,
		}},
		//DhcpServiceIntent:    DhcpServiceIntent{},
		//AntiAffinityPolicy:   AntiAffinityPolicy{},
		VirtualNetworkPolicy: VirtualNetworkPolicy{OverlayControlProtocol: OverlayControlProtocolEvpn},
	}

	for _, client := range clients {
		log.Printf("testing CreateL3CollapsedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		id, err := client.client.CreateL3CollapsedTemplate(context.TODO(), req)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("testing DeleteTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		err = client.client.DeleteTemplate(context.TODO(), id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
