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

var railCollapsedSmall = TemplateRailCollapsed{
	id:             "id__rail_collapsed_small",
	createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:          "Collapsed Fabric 128GPU",
	Racks: []RackTypeWithCount{
		{
			Count: 1,
			RackType: RackType{
				id:                       "305fc185294f6e42d24ccf2acedb7164",
				Label:                    "Collapsed 128GPU",
				FabricConnectivityDesign: enum.FabricConnectivityDesignRailCollapsed,
				LeafSwitches: []LeafSwitch{
					{
						Label: "leaf_1",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
						},
						Tags: []Tag{},
					},
				},
				GenericSystems: []GenericSystem{
					{
						ASNDomain: &enum.FeatureSwitchDisabled,
						Count:     16,
						Label:     "server",
						Links: []RackTypeLink{
							{
								Label:              "link_1",
								TargetSwitchLabel:  "leaf_1",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(1),
								Tags:               []Tag{},
							},
							{
								Label:              "link_2",
								TargetSwitchLabel:  "leaf_1",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(2),
								Tags:               []Tag{},
							},
							{
								Label:              "link_3",
								TargetSwitchLabel:  "leaf_1",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(3),
								Tags:               []Tag{},
							},
							{
								Label:              "link_4",
								TargetSwitchLabel:  "leaf_1",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(4),
								Tags:               []Tag{},
							},
							{
								Label:              "link_5",
								TargetSwitchLabel:  "leaf_1",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(5),
								Tags:               []Tag{},
							},
							{
								Label:              "link_6",
								TargetSwitchLabel:  "leaf_1",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(6),
								Tags:               []Tag{},
							},
							{
								Label:              "link_7",
								TargetSwitchLabel:  "leaf_1",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(7),
								Tags:               []Tag{},
							},
							{
								Label:              "link_8",
								TargetSwitchLabel:  "leaf_1",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(8),
								Tags:               []Tag{},
							},
						},
						LogicalDevice: LogicalDevice{
							Label: "AOS-8x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 4},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 8,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
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
	DHCPServiceIntent:    policy.DHCPServiceIntent{Active: true},
	VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
}

const railCollapsedSmallJSON = `{
  "id": "id__rail_collapsed_small",
  "display_name": "Collapsed Fabric 128GPU",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "rack_type_counts": [
    {
      "rack_type_id": "305fc185294f6e42d24ccf2acedb7164",
      "count": 1
    }
  ],
  "type": "rail_collapsed",
  "capability": "blueprint",
  "virtual_network_policy": {
    "overlay_control_protocol": null
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
          "label": "leaf_1",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        }
      ],
      "logical_devices": [
        {
          "id": "1cdfb7b53085622a86d796b64c7f0449",
          "display_name": "AOS-128x400-1",
          "panels": [
            {
              "port_groups": [
                {
                  "roles": [
                    "superspine",
                    "unused",
                    "leaf",
                    "generic",
                    "peer",
                    "access",
                    "spine"
                  ],
                  "count": 128,
                  "speed": {
                    "value": 400,
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
                "row_count": 4,
                "column_count": 32
              }
            }
          ]
        },
        {
          "id": "bf939cf9cdd183c562aaf8f09a356f11",
          "display_name": "AOS-8x400-1",
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
                    "value": 400,
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
        }
      ],
      "fabric_connectivity_design": "rail_collapsed",
      "id": "305fc185294f6e42d24ccf2acedb7164",
      "generic_systems": [
        {
          "tags": [],
          "loopback": "disabled",
          "asn_domain": "disabled",
          "port_channel_id_max": 0,
          "label": "server",
          "count": 16,
          "management_level": "unmanaged",
          "logical_device": "bf939cf9cdd183c562aaf8f09a356f11",
          "links": [
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_1",
              "link_per_switch_count": 1,
              "label": "link_1",
              "rail_index": 1
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_1",
              "link_per_switch_count": 1,
              "label": "link_2",
              "rail_index": 2
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_1",
              "link_per_switch_count": 1,
              "label": "link_3",
              "rail_index": 3
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_1",
              "link_per_switch_count": 1,
              "label": "link_4",
              "rail_index": 4
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_1",
              "link_per_switch_count": 1,
              "label": "link_5",
              "rail_index": 5
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_1",
              "link_per_switch_count": 1,
              "label": "link_6",
              "rail_index": 6
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_1",
              "link_per_switch_count": 1,
              "label": "link_7",
              "rail_index": 7
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_1",
              "link_per_switch_count": 1,
              "label": "link_8",
              "rail_index": 8
            }
          ],
          "port_channel_id_min": 0
        }
      ],
      "display_name": "Collapsed 128GPU"
    }
  ]
}`

var railCollapsedMedium = TemplateRailCollapsed{
	id:             "id__rail_collapsed_medium",
	createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:          "Collapsed Fabric 512GPU",
	Racks: []RackTypeWithCount{
		{
			Count: 1,
			RackType: RackType{
				id:                       "e14f2a94565c92f9a99e9d7bcf3e1dec",
				Label:                    "Collapsed 512GPU",
				FabricConnectivityDesign: enum.FabricConnectivityDesignRailCollapsed,
				LeafSwitches: []LeafSwitch{
					{
						Label: "leaf_1",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
					{
						Label: "leaf_2",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
					{
						Label: "leaf_3",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
					{
						Label: "leaf_4",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
				},
				GenericSystems: []GenericSystem{
					{
						ASNDomain: &enum.FeatureSwitchDisabled,
						Count:     64,
						Label:     "server",
						Links: []RackTypeLink{
							{
								Label:              "link_1",
								TargetSwitchLabel:  "leaf_1",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(1),
								Tags:               []Tag{},
							},
							{
								Label:              "link_2",
								TargetSwitchLabel:  "leaf_1",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(2),
								Tags:               []Tag{},
							},
							{
								Label:              "link_3",
								TargetSwitchLabel:  "leaf_2",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(3),
								Tags:               []Tag{},
							},
							{
								Label:              "link_4",
								TargetSwitchLabel:  "leaf_2",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(4),
								Tags:               []Tag{},
							},
							{
								Label:              "link_5",
								TargetSwitchLabel:  "leaf_3",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(5),
								Tags:               []Tag{},
							},
							{
								Label:              "link_6",
								TargetSwitchLabel:  "leaf_3",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(6),
								Tags:               []Tag{},
							},
							{
								Label:              "link_7",
								TargetSwitchLabel:  "leaf_4",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(7),
								Tags:               []Tag{},
							},
							{
								Label:              "link_8",
								TargetSwitchLabel:  "leaf_4",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(8),
								Tags:               []Tag{},
							},
						},
						LogicalDevice: LogicalDevice{
							Label: "AOS-8x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 4},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 8,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Loopback:        &enum.FeatureSwitchDisabled,
						ManagementLevel: enum.SystemManagementLevelUnmanaged,
						Tags:            []Tag{},
					},
				},
			},
		},
	},
	DHCPServiceIntent:    policy.DHCPServiceIntent{Active: true},
	VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
}

const railCollapsedMediumJSON = `{
  "id": "id__rail_collapsed_medium",
  "display_name": "Collapsed Fabric 512GPU",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "rack_type_counts": [
    {
      "rack_type_id": "e14f2a94565c92f9a99e9d7bcf3e1dec",
      "count": 1
    }
  ],
  "type": "rail_collapsed",
  "capability": "blueprint",
  "virtual_network_policy": {
    "overlay_control_protocol": null
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
          "label": "leaf_1",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        },
        {
          "tags": [],
          "label": "leaf_2",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        },
        {
          "tags": [],
          "label": "leaf_3",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        },
        {
          "tags": [],
          "label": "leaf_4",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        }
      ],
      "logical_devices": [
        {
          "id": "1cdfb7b53085622a86d796b64c7f0449",
          "display_name": "AOS-128x400-1",
          "panels": [
            {
              "port_groups": [
                {
                  "roles": [
                    "superspine",
                    "unused",
                    "leaf",
                    "generic",
                    "peer",
                    "access",
                    "spine"
                  ],
                  "count": 128,
                  "speed": {
                    "value": 400,
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
                "row_count": 4,
                "column_count": 32
              }
            }
          ]
        },
        {
          "id": "bf939cf9cdd183c562aaf8f09a356f11",
          "display_name": "AOS-8x400-1",
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
                    "value": 400,
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
        }
      ],
      "fabric_connectivity_design": "rail_collapsed",
      "id": "e14f2a94565c92f9a99e9d7bcf3e1dec",
      "generic_systems": [
        {
          "tags": [],
          "loopback": "disabled",
          "asn_domain": "disabled",
          "port_channel_id_max": 0,
          "label": "server",
          "count": 64,
          "management_level": "unmanaged",
          "logical_device": "bf939cf9cdd183c562aaf8f09a356f11",
          "links": [
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_1",
              "link_per_switch_count": 1,
              "label": "link_1",
              "rail_index": 1
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_1",
              "link_per_switch_count": 1,
              "label": "link_2",
              "rail_index": 2
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_2",
              "link_per_switch_count": 1,
              "label": "link_3",
              "rail_index": 3
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_2",
              "link_per_switch_count": 1,
              "label": "link_4",
              "rail_index": 4
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_3",
              "link_per_switch_count": 1,
              "label": "link_5",
              "rail_index": 5
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_3",
              "link_per_switch_count": 1,
              "label": "link_6",
              "rail_index": 6
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_4",
              "link_per_switch_count": 1,
              "label": "link_7",
              "rail_index": 7
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_4",
              "link_per_switch_count": 1,
              "label": "link_8",
              "rail_index": 8
            }
          ],
          "port_channel_id_min": 0
        }
      ],
      "display_name": "Collapsed 512GPU"
    }
  ]
}`

var railCollapsedLarge = TemplateRailCollapsed{
	id:             "id__rail_collapsed_large",
	createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:          "Collapsed Fabric 1024GPU",
	Racks: []RackTypeWithCount{
		{
			Count: 1,
			RackType: RackType{
				id:                       "4296e45fff3d84e0b13943d42d6ae390",
				Label:                    "Collapsed 1024GPU",
				FabricConnectivityDesign: enum.FabricConnectivityDesignRailCollapsed,
				LeafSwitches: []LeafSwitch{
					{
						Label: "leaf_1",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
					{
						Label: "leaf_2",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
					{
						Label: "leaf_3",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
					{
						Label: "leaf_4",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
					{
						Label: "leaf_5",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
					{
						Label: "leaf_6",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
					{
						Label: "leaf_7",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
					{
						Label: "leaf_8",
						LogicalDevice: LogicalDevice{
							Label: "AOS-128x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 128,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Tags: []Tag{},
					},
				},
				GenericSystems: []GenericSystem{
					{
						ASNDomain: &enum.FeatureSwitchDisabled,
						Count:     128,
						Label:     "server",
						Links: []RackTypeLink{
							{
								Label:              "link_1",
								TargetSwitchLabel:  "leaf_1",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(1),
								Tags:               []Tag{},
							},
							{
								Label:              "link_2",
								TargetSwitchLabel:  "leaf_2",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(2),
								Tags:               []Tag{},
							},
							{
								Label:              "link_3",
								TargetSwitchLabel:  "leaf_3",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(3),
								Tags:               []Tag{},
							},
							{
								Label:              "link_4",
								TargetSwitchLabel:  "leaf_4",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(4),
								Tags:               []Tag{},
							},
							{
								Label:              "link_5",
								TargetSwitchLabel:  "leaf_5",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(5),
								Tags:               []Tag{},
							},
							{
								Label:              "link_6",
								TargetSwitchLabel:  "leaf_6",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(6),
								Tags:               []Tag{},
							},
							{
								Label:              "link_7",
								TargetSwitchLabel:  "leaf_7",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(7),
								Tags:               []Tag{},
							},
							{
								Label:              "link_8",
								TargetSwitchLabel:  "leaf_8",
								LinkPerSwitchCount: 1,
								Speed:              "400G",
								AttachmentType:     enum.LinkAttachmentTypeSingle,
								RailIndex:          pointer.To(8),
								Tags:               []Tag{},
							},
						},
						LogicalDevice: LogicalDevice{
							Label: "AOS-8x400-1",
							Panels: []LogicalDevicePanel{
								{
									PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 4},
									PortGroups: []LogicalDevicePanelPortGroup{
										{
											Count: 8,
											Speed: "400G",
											Roles: LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess},
										},
									},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
								},
							},
							id: "",
						},
						Loopback:        &enum.FeatureSwitchDisabled,
						ManagementLevel: enum.SystemManagementLevelUnmanaged,
						Tags:            []Tag{},
					},
				},
			},
		},
	},
	DHCPServiceIntent:    policy.DHCPServiceIntent{Active: true},
	VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
}

const railCollapsedLargeJSON = `{
  "id": "id__rail_collapsed_large",
  "display_name": "Collapsed Fabric 1024GPU",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "rack_type_counts": [
    {
      "rack_type_id": "4296e45fff3d84e0b13943d42d6ae390",
      "count": 1
    }
  ],
  "type": "rail_collapsed",
  "capability": "blueprint",
  "virtual_network_policy": {
    "overlay_control_protocol": null
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
          "label": "leaf_1",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        },
        {
          "tags": [],
          "label": "leaf_2",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        },
        {
          "tags": [],
          "label": "leaf_3",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        },
        {
          "tags": [],
          "label": "leaf_4",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        },
        {
          "tags": [],
          "label": "leaf_5",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        },
        {
          "tags": [],
          "label": "leaf_6",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        },
        {
          "tags": [],
          "label": "leaf_7",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        },
        {
          "tags": [],
          "label": "leaf_8",
          "logical_device": "1cdfb7b53085622a86d796b64c7f0449"
        }
      ],
      "logical_devices": [
        {
          "id": "1cdfb7b53085622a86d796b64c7f0449",
          "display_name": "AOS-128x400-1",
          "panels": [
            {
              "port_groups": [
                {
                  "roles": [
                    "superspine",
                    "unused",
                    "leaf",
                    "generic",
                    "peer",
                    "access",
                    "spine"
                  ],
                  "count": 128,
                  "speed": {
                    "value": 400,
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
                "row_count": 4,
                "column_count": 32
              }
            }
          ]
        },
        {
          "id": "bf939cf9cdd183c562aaf8f09a356f11",
          "display_name": "AOS-8x400-1",
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
                    "value": 400,
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
        }
      ],
      "fabric_connectivity_design": "rail_collapsed",
      "id": "4296e45fff3d84e0b13943d42d6ae390",
      "generic_systems": [
        {
          "tags": [],
          "loopback": "disabled",
          "asn_domain": "disabled",
          "port_channel_id_max": 0,
          "label": "server",
          "count": 128,
          "management_level": "unmanaged",
          "logical_device": "bf939cf9cdd183c562aaf8f09a356f11",
          "links": [
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_1",
              "link_per_switch_count": 1,
              "label": "link_1",
              "rail_index": 1
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_2",
              "link_per_switch_count": 1,
              "label": "link_2",
              "rail_index": 2
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_3",
              "link_per_switch_count": 1,
              "label": "link_3",
              "rail_index": 3
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_4",
              "link_per_switch_count": 1,
              "label": "link_4",
              "rail_index": 4
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_5",
              "link_per_switch_count": 1,
              "label": "link_5",
              "rail_index": 5
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_6",
              "link_per_switch_count": 1,
              "label": "link_6",
              "rail_index": 6
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_7",
              "link_per_switch_count": 1,
              "label": "link_7",
              "rail_index": 7
            },
            {
              "tags": [],
              "attachment_type": "singleAttached",
              "link_speed": {
                "value": 400,
                "unit": "G"
              },
              "target_switch_label": "leaf_8",
              "link_per_switch_count": 1,
              "label": "link_8",
              "rail_index": 8
            }
          ],
          "port_channel_id_min": 0
        }
      ],
      "display_name": "Collapsed 1024GPU"
    }
  ]
}`
