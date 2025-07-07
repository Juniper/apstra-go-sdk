// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"encoding/json"
	"log"
	"testing"
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
