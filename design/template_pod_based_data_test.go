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

var l2SuperspineMultiPlane = TemplatePodBased{
	id:             "id__L2_superspine_multi_plane",
	createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:          "L2 superspine multi plane",
	Superspine: Superspine{
		PlaneCount:         4,
		SuperspinePerPlane: 4,
		Tags:               []Tag{},
		LogicalDevice: LogicalDevice{
			Label: "AOS-32x40-3",
			Panels: []LogicalDevicePanel{
				{
					PanelLayout:  LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 16},
					PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
					PortGroups: []LogicalDevicePanelPortGroup{
						{Count: 32, Speed: "40G", Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRoleSpine}},
					},
				},
			},
		},
	},
	Pods: []PodWithCount{
		{
			Count: 2,
			Pod: TemplateRackBased{
				Label: "L2 Pod",
				Racks: []RackTypeWithCount{
					{
						Count: 1,
						RackType: RackType{
							Label:                    "L2 One Leaf",
							FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
							LeafSwitches: []LeafSwitch{
								{
									Label:             "leaf",
									LinkPerSpineCount: pointer.To(1),
									LinkPerSpineSpeed: pointer.To(speed.Speed("40G")),
									LogicalDevice: LogicalDevice{
										Label: "AOS-64x10+16x40-2",
										Panels: []LogicalDevicePanel{
											{
												PanelLayout:  LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 32},
												PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
												PortGroups: []LogicalDevicePanelPortGroup{
													{Count: 64, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}},
												},
											},
											{
												PanelLayout:  LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
												PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
												PortGroups: []LogicalDevicePanelPortGroup{
													{Count: 16, Speed: "40G", Roles: LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleSpine}},
												},
											},
										},
									},
									Tags: []Tag{},
								},
							},
							GenericSystems: []GenericSystem{
								{
									ASNDomain: &enum.FeatureSwitchDisabled,
									Count:     48,
									Label:     "generic",
									Links: []RackTypeLink{
										{
											Label:              "link",
											TargetSwitchLabel:  "leaf",
											LinkPerSwitchCount: 1,
											Speed:              "10G",
											AttachmentType:     enum.LinkAttachmentTypeSingle,
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
													{Count: 1, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess}},
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
					{
						Count: 1,
						RackType: RackType{
							Label:                    "L2 Mlag Leaf",
							FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
							LeafSwitches: []LeafSwitch{
								{
									Label:             "leaf",
									LinkPerSpineCount: pointer.To(1),
									LinkPerSpineSpeed: pointer.To(speed.Speed("40G")),
									LogicalDevice: LogicalDevice{
										Label: "AOS-64x10+16x40-2",
										Panels: []LogicalDevicePanel{
											{
												PanelLayout:  LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 32},
												PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
												PortGroups: []LogicalDevicePanelPortGroup{
													{Count: 64, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}},
												},
											},
											{
												PanelLayout:  LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
												PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
												PortGroups: []LogicalDevicePanelPortGroup{
													{Count: 16, Speed: "40G", Roles: LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleSpine}},
												},
											},
										},
									},
									RedundancyProtocol: enum.LeafRedundancyProtocolMLAG,
									Tags:               []Tag{},
									MLAGInfo: &RackTypeLeafSwitchMLAGInfo{
										LeafLeafLinkCount: 4,
										LeafLeafLinkSpeed: "40G",
										MLAGVLAN:          2999,
									},
								},
							},
							GenericSystems: []GenericSystem{
								{
									ASNDomain: &enum.FeatureSwitchDisabled,
									Count:     24,
									Label:     "generic(1)",
									Links: []RackTypeLink{
										{
											Label:              "link",
											TargetSwitchLabel:  "leaf",
											LinkPerSwitchCount: 1,
											Speed:              "10G",
											AttachmentType:     enum.LinkAttachmentTypeSingle,
											SwitchPeer:         enum.LinkSwitchPeerFirst,
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
													{Count: 1, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess}},
												},
											},
										},
									},
									Loopback:        &enum.FeatureSwitchDisabled,
									ManagementLevel: enum.SystemManagementLevelUnmanaged,
									Tags:            []Tag{},
								},
								{
									ASNDomain: &enum.FeatureSwitchDisabled,
									Count:     24, Label: "generic(2)",
									Links: []RackTypeLink{
										{
											Label:              "link",
											TargetSwitchLabel:  "leaf",
											LinkPerSwitchCount: 1,
											Speed:              "10G",
											AttachmentType:     enum.LinkAttachmentTypeSingle,
											SwitchPeer:         enum.LinkSwitchPeerSecond,
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
													{Count: 1, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess}},
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
				AsnAllocationPolicy: &AsnAllocationPolicy{SpineAsnScheme: enum.AsnAllocationSchemeSingle},
				DHCPServiceIntent:   policy.DHCPServiceIntent{Active: true},
				Spine: Spine{
					Count:                  4,
					LinkPerSuperspineCount: 1,
					LinkPerSuperspineSpeed: "40G",
					LogicalDevice: LogicalDevice{
						Label: "AOS-32x40-3",
						Panels: []LogicalDevicePanel{
							{
								PanelLayout:  LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 16},
								PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								PortGroups: []LogicalDevicePanelPortGroup{
									{Count: 32, Speed: "40G", Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRoleSpine}},
								},
							},
						},
					},
					Tags: []Tag{},
				},
				VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
			},
		},
	},
}

const l2SuperspineMultiPlaneJSON = `{
  "created_at": "2025-05-15T18:11:11.087777Z",
  "last_modified_at": "2025-05-15T18:11:11.087777Z",
  "display_name": "L2 superspine multi plane",
  "superspine": {
    "plane_count": 4,
    "superspine_per_plane": 4,
    "logical_device": {
      "id": "4d1ae9eb01a0c918a815ee6b28f5c6af",
      "display_name": "AOS-32x40-3",
      "panels": [
        {
          "port_groups": [
            {
              "roles": [
                "superspine",
                "leaf",
                "generic",
                "spine"
              ],
              "count": 32,
              "speed": {
                "value": 40,
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
            "row_count": 2,
            "column_count": 16
          }
        }
      ]
    },
    "tags": []
  },
  "type": "pod_based",
  "capability": "blueprint",
  "rack_based_templates": [
    {
      "id": "6728b6cd95d37e89807cd18988f3c3e7",
      "display_name": "L2 Pod",
      "type": "rack_based",
      "rack_type_counts": [
        {
          "rack_type_id": "8b4d80a049caf462d174640282014198",
          "count": 1
        },
        {
          "rack_type_id": "80677cb9f0d4c35bf036a0701f8f5111",
          "count": 1
        }
      ],
      "rack_types": [
        {
          "description": "",
          "tags": [],
          "leafs": [
            {
              "tags": [],
              "link_per_spine_count": 1,
              "link_per_spine_speed": {
                "value": 40,
                "unit": "G"
              },
              "label": "leaf",
              "logical_device": "a9e9d630eac56c6a91c5f4b208f19d5d"
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
              "id": "a9e9d630eac56c6a91c5f4b208f19d5d",
              "display_name": "AOS-64x10+16x40-2",
              "panels": [
                {
                  "port_groups": [
                    {
                      "roles": [
                        "generic",
                        "access"
                      ],
                      "count": 64,
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
                    "row_count": 2,
                    "column_count": 32
                  }
                },
                {
                  "port_groups": [
                    {
                      "roles": [
                        "generic",
                        "peer",
                        "spine"
                      ],
                      "count": 16,
                      "speed": {
                        "value": 40,
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
                    "row_count": 2,
                    "column_count": 8
                  }
                }
              ]
            }
          ],
          "fabric_connectivity_design": "l3clos",
          "id": "8b4d80a049caf462d174640282014198",
          "generic_systems": [
            {
              "tags": [],
              "loopback": "disabled",
              "asn_domain": "disabled",
              "port_channel_id_max": 0,
              "label": "generic",
              "count": 48,
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
          "display_name": "L2 One Leaf"
        },
        {
          "description": "",
          "tags": [],
          "leafs": [
            {
              "tags": [],
              "link_per_spine_count": 1,
              "redundancy_protocol": "mlag",
              "leaf_leaf_link_speed": {
                "value": 40,
                "unit": "G"
              },
              "link_per_spine_speed": {
                "value": 40,
                "unit": "G"
              },
              "label": "leaf",
              "logical_device": "a9e9d630eac56c6a91c5f4b208f19d5d",
              "leaf_leaf_link_count": 4,
              "mlag_vlan_id": 2999
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
              "id": "a9e9d630eac56c6a91c5f4b208f19d5d",
              "display_name": "AOS-64x10+16x40-2",
              "panels": [
                {
                  "port_groups": [
                    {
                      "roles": [
                        "generic",
                        "access"
                      ],
                      "count": 64,
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
                    "row_count": 2,
                    "column_count": 32
                  }
                },
                {
                  "port_groups": [
                    {
                      "roles": [
                        "generic",
                        "peer",
                        "spine"
                      ],
                      "count": 16,
                      "speed": {
                        "value": 40,
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
                    "row_count": 2,
                    "column_count": 8
                  }
                }
              ]
            }
          ],
          "fabric_connectivity_design": "l3clos",
          "id": "80677cb9f0d4c35bf036a0701f8f5111",
          "generic_systems": [
            {
              "tags": [],
              "loopback": "disabled",
              "asn_domain": "disabled",
              "port_channel_id_max": 0,
              "label": "generic(1)",
              "count": 24,
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
                  "label": "link",
                  "switch_peer": "first"
                }
              ],
              "port_channel_id_min": 0
            },
            {
              "tags": [],
              "loopback": "disabled",
              "asn_domain": "disabled",
              "port_channel_id_max": 0,
              "label": "generic(2)",
              "count": 24,
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
                  "label": "link",
                  "switch_peer": "second"
                }
              ],
              "port_channel_id_min": 0
            }
          ],
          "display_name": "L2 Mlag Leaf"
        }
      ],
      "dhcp_service_intent": {
        "active": true
      },
      "virtual_network_policy": {
        "overlay_control_protocol": null
      },
      "asn_allocation_policy": {
        "spine_asn_scheme": "single"
      },
      "spine": {
        "count": 4,
        "logical_device": {
          "id": "4d1ae9eb01a0c918a815ee6b28f5c6af",
          "display_name": "AOS-32x40-3",
          "panels": [
            {
              "port_groups": [
                {
                  "roles": [
                    "superspine",
                    "leaf",
                    "generic",
                    "spine"
                  ],
                  "count": 32,
                  "speed": {
                    "value": 40,
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
                "row_count": 2,
                "column_count": 16
              }
            }
          ]
        },
        "link_per_superspine_count": 1,
        "link_per_superspine_speed": {
          "value": 40,
          "unit": "G"
        },
        "tags": []
      }
    }
  ],
  "rack_based_template_counts": [
    {
      "rack_based_template_id": "6728b6cd95d37e89807cd18988f3c3e7",
      "count": 2
    }
  ]
}`
