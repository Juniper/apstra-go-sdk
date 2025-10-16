// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
	"github.com/Juniper/apstra-go-sdk/policy"
	"github.com/Juniper/apstra-go-sdk/speed"
)

var templateRackBasedL2VirtualEVPN = TemplateRackBased{
	Label: "templateRackBasedL2VirtualEVPN",
	Racks: []RackTypeWithCount{
		{
			Count: 4,
			RackType: RackType{
				id:                       "1a976a730bf7dbc1e1e4f134820e28e8",
				Label:                    "L2 Virtual",
				FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
				LeafSwitches: []LeafSwitch{
					{
						Label:             "leaf",
						LinkPerSpineCount: pointer.To(1),
						LinkPerSpineSpeed: pointer.To(speed.Speed("10G")),
						LogicalDevice: LogicalDevice{
							Label: "AOS-7x10-Leaf",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout:  LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 7},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 2,
											Speed: "10G",
											Roles: []enum.PortRole{enum.PortRoleLeaf, enum.PortRoleSpine},
										},
										{
											Count: 2,
											Speed: "10G",
											Roles: []enum.PortRole{enum.PortRolePeer},
										},
										{
											Count: 2,
											Speed: "10G",
											Roles: []enum.PortRole{enum.PortRoleGeneric, enum.PortRoleAccess},
										},
										{
											Count: 1,
											Speed: "10G",
											Roles: []enum.PortRole{enum.PortRoleGeneric},
										},
									},
								},
							},
						},
						RedundancyProtocol: enum.LeafRedundancyProtocol{},
						Tags:               []Tag{},
					},
				},
				AccessSwitches: nil,
				GenericSystems: []GenericSystem{
					{
						Count:     2,
						Label:     "generic",
						AsnDomain: &enum.FeatureSwitchDisabled,
						Links: []RackTypeLink{
							{
								Label:              "link",
								TargetSwitchLabel:  "leaf",
								LinkPerSwitchCount: 1,
								Speed:              "10G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								LAGMode:            enum.LAGModeNone,
								Tags:               []Tag{},
							},
						},
						LogicalDevice: LogicalDevice{
							Label: "AOS-1x10-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout:  LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 1},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 1,
											Speed: "10G",
											Roles: []enum.PortRole{enum.PortRoleLeaf, enum.PortRoleAccess},
										},
									},
								},
							},
						},
						Loopback:        &enum.FeatureSwitchDisabled,
						ManagementLevel: enum.SystemManagementLevelUnmanaged,
						Tags:            []Tag{},
					},
				},
			},
		},
	},
	AsnAllocationPolicy: &AsnAllocationPolicy{SpineAsnScheme: enum.AsnAllocationSchemeDistinct},
	Capability:          &enum.TemplateCapabilityBlueprint,
	DHCPServiceIntent:   policy.DHCPServiceIntent{Active: true},
	Spine: Spine{
		Count: 2,
		LogicalDevice: LogicalDevice{
			id:    "f9f963360acfe0a5362ac9fad36020af",
			Label: "AOS-7x10-Spine",
			Panels: []LogicalDevicePanel{
				{
					PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 7},
					PortGroups: []LogicalDevicePanelPortGroup{
						{
							Count: 5,
							Speed: "10G",
							Roles: []enum.PortRole{enum.PortRoleSuperspine, enum.PortRoleLeaf},
						},
						{
							Count: 2,
							Speed: "10G",
							Roles: []enum.PortRole{enum.PortRoleGeneric},
						},
					},
					PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
				},
			},
		},
	},
	VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolEVPN},
	id:                   "id__templateRackBasedL2VirtualEVPN",
	createdAt:            pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt:       pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
}

const templateRackBasedL2VirtualEVPNJSON = `{
  "id": "id__templateRackBasedL2VirtualEVPN",
  "display_name": "templateRackBasedL2VirtualEVPN",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "rack_type_counts": [
    {
      "rack_type_id": "1a976a730bf7dbc1e1e4f134820e28e8",
      "count": 4
    }
  ],
  "asn_allocation_policy": {
    "spine_asn_scheme": "distinct"
  },
  "type": "rack_based",
  "capability": "blueprint",
  "virtual_network_policy": {
    "overlay_control_protocol": "evpn"
  },
  "dhcp_service_intent": {
    "active": true
  },
  "spine": {
    "count": 2,
    "logical_device": {
      "id": "f9f963360acfe0a5362ac9fad36020af",
      "display_name": "AOS-7x10-Spine",
      "panels": [
        {
          "port_groups": [
            {
              "roles": [
                "superspine",
                "leaf"
              ],
              "count": 5,
              "speed": {
                "value": 10,
                "unit": "G"
              }
            },
            {
              "roles": [
                "generic"
              ],
              "count": 2,
              "speed": {
                "value": 10,
                "unit": "G"
              }
            }
          ],
          "port_indexing": {
            "schema": "absolute",
            "order": "T-B, L-R",
            "start_index": 1
          },
          "panel_layout": {
            "row_count": 1,
            "column_count": 7
          }
        }
      ]
    },
    "link_per_superspine_count": 0,
    "link_per_superspine_speed": null,
    "tags": null
  },
  "rack_types": [
    {
      "description": "",
      "tags": [],
      "leafs": [
        {
          "tags": [],
          "link_per_spine_count": 1,
          "link_per_spine_speed": {
            "value": 10,
            "unit": "G"
          },
          "label": "leaf",
          "logical_device": "7f62dc1877fb6a7da4ca8e9411b74180"
        }
      ],
      "logical_devices": [
        {
          "id": "7dc18619d1d57562a5ca0738535a7d04",
          "display_name": "AOS-1x10-1",
          "panels": [
            {
              "port_groups": [
                {
                  "roles": [
                    "leaf",
                    "access"
                  ],
                  "count": 1,
                  "speed": {
                    "value": 10,
                    "unit": "G"
                  }
                }
              ],
              "port_indexing": {
                "schema": "absolute",
                "order": "T-B, L-R",
                "start_index": 1
              },
              "panel_layout": {
                "row_count": 1,
                "column_count": 1
              }
            }
          ]
        },
        {
          "id": "7f62dc1877fb6a7da4ca8e9411b74180",
          "display_name": "AOS-7x10-Leaf",
          "panels": [
            {
              "port_groups": [
                {
                  "roles": [
                    "leaf",
                    "spine"
                  ],
                  "count": 2,
                  "speed": {
                    "value": 10,
                    "unit": "G"
                  }
                },
                {
                  "roles": [
                    "peer"
                  ],
                  "count": 2,
                  "speed": {
                    "value": 10,
                    "unit": "G"
                  }
                },
                {
                  "roles": [
                    "generic",
                    "access"
                  ],
                  "count": 2,
                  "speed": {
                    "value": 10,
                    "unit": "G"
                  }
                },
                {
                  "roles": [
                    "generic"
                  ],
                  "count": 1,
                  "speed": {
                    "value": 10,
                    "unit": "G"
                  }
                }
              ],
              "port_indexing": {
                "schema": "absolute",
                "order": "T-B, L-R",
                "start_index": 1
              },
              "panel_layout": {
                "row_count": 1,
                "column_count": 7
              }
            }
          ]
        }
      ],
      "fabric_connectivity_design": "l3clos",
      "id": "1a976a730bf7dbc1e1e4f134820e28e8",
      "generic_systems": [
        {
          "tags": [],
          "loopback": "disabled",
          "asn_domain": "disabled",
          "port_channel_id_max": 0,
          "label": "generic",
          "count": 2,
          "management_level": "unmanaged",
          "logical_device": "7dc18619d1d57562a5ca0738535a7d04",
          "links": [
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 10,
                "unit": "G"
              },
              "target_switch_label": "leaf",
              "link_per_switch_count": 1,
              "label": "link"
            }
          ],
          "port_channel_id_min": 0
        }
      ],
      "display_name": "L2 Virtual"
    }
  ]
}`
