// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
	"github.com/Juniper/apstra-go-sdk/policy"
)

var templateL3CollapsedACS = TemplateL3Collapsed{
	id:             "id__templateL3CollapsedACS",
	createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:          "templateL3CollapsedACS",
	MeshLinkCount:  1,
	MeshLinkSpeed:  "10G",
	Racks: []RackTypeWithCount{
		{
			Count: 1,
			RackType: RackType{
				Label:                    "Collapsed 1xleaf",
				FabricConnectivityDesign: enum.FabricConnectivityDesign{Value: "l3collapsed"},
				LeafSwitches: []LeafSwitch{
					{
						Label: "leaf",
						LogicalDevice: LogicalDevice{
							Label: "AOS-7x10-Leaf",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 7},
									PortGroups: []LogicalDevicePanelPortGroup{
										{Count: 2, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "leaf"}, enum.PortRole{Value: "spine"}}},
										{Count: 2, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "peer"}}},
										{Count: 2, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "generic"}, enum.PortRole{Value: "access"}}},
										{Count: 1, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "generic"}}},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexing{Value: "T-B, L-R"},
								},
							},
						},
					},
				},
				AccessSwitches: []AccessSwitch{
					{
						Count: 1,
						Label: "access",
						Links: []RackTypeLink{
							{
								Label:              "leaf_link",
								TargetSwitchLabel:  "leaf",
								LinkPerSwitchCount: 1,
								Speed:              "10G",
								AttachmentType:     enum.LinkAttachmentType{Value: "singleAttached"},
								LAGMode:            enum.LAGMode{Value: "lacp_active"},
							},
						},
						LogicalDevice: LogicalDevice{
							Label: "AOS-8x10-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 4},
									PortGroups: []LogicalDevicePanelPortGroup{
										{Count: 8, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "leaf"}, enum.PortRole{Value: "generic"}, enum.PortRole{Value: "peer"}, enum.PortRole{Value: "access"}}},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexing{Value: "T-B, L-R"},
								},
							},
						},
					},
				},
				GenericSystems: []GenericSystem{
					{
						AsnDomain: &enum.FeatureSwitchDisabled,
						Count:     2,
						Label:     "generic",
						Links: []RackTypeLink{
							{
								Label:              "link",
								TargetSwitchLabel:  "access",
								LinkPerSwitchCount: 1,
								Speed:              "10G",
								AttachmentType:     enum.LinkAttachmentType{Value: "singleAttached"},
							},
						},
						LogicalDevice: LogicalDevice{
							Label: "AOS-1x10-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 1},
									PortGroups: []LogicalDevicePanelPortGroup{
										{Count: 1, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "leaf"}, enum.PortRole{Value: "access"}}},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexing{Value: "T-B, L-R"},
								},
							},
						},
						Loopback:        &enum.FeatureSwitchDisabled,
						ManagementLevel: enum.SystemManagementLevel{Value: "unmanaged"},
					},
				},
				id:             "7423322458cb07c0f948023d39a2e8bf",
				createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "1970-01-01T00:00:00.000000Z")),
				lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "1970-01-01T00:00:00.000000Z")),
			},
		},
	},
	DHCPServiceIntent:    policy.DHCPServiceIntent{Active: true},
	VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolEVPN},
}

const templateL3CollapsedACSJSON = `{
  "id": "id__templateL3CollapsedACS",
  "display_name": "templateL3CollapsedACS",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "rack_type_counts": [
    {
      "rack_type_id": "7423322458cb07c0f948023d39a2e8bf",
      "count": 1
    }
  ],
  "type": "l3_collapsed",
  "capability": "blueprint",
  "virtual_network_policy": {
    "overlay_control_protocol": "evpn"
  },
  "mesh_link_speed": {
    "value": 10,
    "unit": "G"
  },
  "dhcp_service_intent": {
    "active": true
  },
  "rack_types": [
    {
      "description": "",
      "tags": [],
      "leafs": [
        {
          "tags": [],
          "label": "leaf",
          "logical_device": "7f62dc1877fb6a7da4ca8e9411b74180"
        }
      ],
      "logical_devices": [
        {
          "id": "3d06a0827e9a3969933edafd82d26387",
          "display_name": "AOS-8x10-1",
          "panels": [
            {
              "port_groups": [
                {
                  "roles": [
                    "leaf",
                    "generic",
                    "peer",
                    "access"
                  ],
                  "count": 8,
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
                "column_count": 4
              }
            }
          ]
        },
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
      "access_switches": [
        {
          "tags": [],
          "access_access_link_port_channel_id_max": 0,
          "access_access_link_count": 0,
          "label": "access",
          "access_access_link_port_channel_id_min": 0,
          "logical_device": "3d06a0827e9a3969933edafd82d26387",
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
              "label": "leaf_link",
              "lag_mode": "lacp_active"
            }
          ],
          "instance_count": 1,
          "access_access_link_speed": null
        }
      ],
      "fabric_connectivity_design": "l3collapsed",
      "id": "7423322458cb07c0f948023d39a2e8bf",
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
              "target_switch_label": "access",
              "link_per_switch_count": 1,
              "label": "link"
            }
          ],
          "port_channel_id_min": 0
        }
      ],
      "display_name": "Collapsed 1xleaf"
    }
  ],
  "mesh_link_count": 1
}`

var templateL3CollapsedACS420 = TemplateL3Collapsed{
	id:             "id__templateL3CollapsedACS",
	createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:          "templateL3CollapsedACS",
	MeshLinkCount:  1,
	MeshLinkSpeed:  "10G",
	Racks: []RackTypeWithCount{
		{
			Count: 1,
			RackType: RackType{
				Label:                    "Collapsed 1xleaf",
				FabricConnectivityDesign: enum.FabricConnectivityDesign{Value: "l3collapsed"},
				LeafSwitches: []LeafSwitch{
					{
						Label: "leaf",
						LogicalDevice: LogicalDevice{
							Label: "AOS-7x10-Leaf",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 7},
									PortGroups: []LogicalDevicePanelPortGroup{
										{Count: 2, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "leaf"}, enum.PortRole{Value: "spine"}}},
										{Count: 2, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "peer"}}},
										{Count: 2, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "generic"}, enum.PortRole{Value: "access"}}},
										{Count: 1, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "generic"}}},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexing{Value: "T-B, L-R"},
								},
							},
						},
					},
				},
				AccessSwitches: []AccessSwitch{
					{
						Count: 1,
						Label: "access",
						Links: []RackTypeLink{
							{
								Label:              "leaf_link",
								TargetSwitchLabel:  "leaf",
								LinkPerSwitchCount: 1,
								Speed:              "10G",
								AttachmentType:     enum.LinkAttachmentType{Value: "singleAttached"},
								LAGMode:            enum.LAGMode{Value: "lacp_active"},
							},
						},
						LogicalDevice: LogicalDevice{
							Label: "AOS-8x10-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 4},
									PortGroups: []LogicalDevicePanelPortGroup{
										{Count: 8, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "leaf"}, enum.PortRole{Value: "generic"}, enum.PortRole{Value: "peer"}, enum.PortRole{Value: "access"}}},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexing{Value: "T-B, L-R"},
								},
							},
						},
					},
				},
				GenericSystems: []GenericSystem{
					{
						AsnDomain: &enum.FeatureSwitchDisabled,
						Count:     2,
						Label:     "generic",
						Links: []RackTypeLink{
							{
								Label:              "link",
								TargetSwitchLabel:  "access",
								LinkPerSwitchCount: 1,
								Speed:              "10G",
								AttachmentType:     enum.LinkAttachmentType{Value: "singleAttached"},
							},
						},
						LogicalDevice: LogicalDevice{
							Label: "AOS-1x10-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 1},
									PortGroups: []LogicalDevicePanelPortGroup{
										{Count: 1, Speed: "10G", Roles: LogicalDevicePortRoles{enum.PortRole{Value: "leaf"}, enum.PortRole{Value: "access"}}},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexing{Value: "T-B, L-R"},
								},
							},
						},
						Loopback:        &enum.FeatureSwitchDisabled,
						ManagementLevel: enum.SystemManagementLevel{Value: "unmanaged"},
					},
				},
				id:             "7423322458cb07c0f948023d39a2e8bf",
				createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "1970-01-01T00:00:00.000000Z")),
				lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "1970-01-01T00:00:00.000000Z")),
			},
		},
	},
	DHCPServiceIntent:    policy.DHCPServiceIntent{Active: true},
	VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolEVPN},
	AntiAffinityPolicy: &policy.AntiAffinity{
		MaxLinksPerPort:          4,
		MaxLinksPerSlot:          8,
		MaxPerSystemLinksPerPort: 1,
		MaxPerSystemLinksPerSlot: 2,
		Mode:                     enum.AntiAffinityModeLoose,
	},
}

const templateL3CollapsedACS420JSON = `{
  "id": "id__templateL3CollapsedACS",
  "display_name": "templateL3CollapsedACS",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "rack_type_counts": [
    {
      "rack_type_id": "7423322458cb07c0f948023d39a2e8bf",
      "count": 1
    }
  ],
  "type": "l3_collapsed",
  "capability": "blueprint",
  "virtual_network_policy": {
    "overlay_control_protocol": "evpn"
  },
  "mesh_link_speed": {
    "value": 10,
    "unit": "G"
  },
  "dhcp_service_intent": {
    "active": true
  },
  "rack_types": [
    {
      "description": "",
      "tags": [],
      "leafs": [
        {
          "tags": [],
          "label": "leaf",
          "logical_device": "7f62dc1877fb6a7da4ca8e9411b74180"
        }
      ],
      "logical_devices": [
        {
          "id": "3d06a0827e9a3969933edafd82d26387",
          "display_name": "AOS-8x10-1",
          "panels": [
            {
              "port_groups": [
                {
                  "roles": [
                    "leaf",
                    "generic",
                    "peer",
                    "access"
                  ],
                  "count": 8,
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
                "column_count": 4
              }
            }
          ]
        },
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
      "access_switches": [
        {
          "tags": [],
          "access_access_link_port_channel_id_max": 0,
          "access_access_link_count": 0,
          "label": "access",
          "access_access_link_port_channel_id_min": 0,
          "logical_device": "3d06a0827e9a3969933edafd82d26387",
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
              "label": "leaf_link",
              "lag_mode": "lacp_active"
            }
          ],
          "instance_count": 1,
          "access_access_link_speed": null
        }
      ],
      "fabric_connectivity_design": "l3collapsed",
      "id": "7423322458cb07c0f948023d39a2e8bf",
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
              "target_switch_label": "access",
              "link_per_switch_count": 1,
              "label": "link"
            }
          ],
          "port_channel_id_min": 0
        }
      ],
      "display_name": "Collapsed 1xleaf"
    }
  ],
  "mesh_link_count": 1,
  "anti_affinity_policy": {
    "mode": "enabled_loose",
    "algorithm": "heuristic",
    "max_links_per_slot": 8,
    "max_links_per_port": 4,
    "max_per_system_links_per_slot": 2,
    "max_per_system_links_per_port": 1
  }
}`
