// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var rackTypeTestCollapsedSimple = RackType{
	id:                       "id__collapsed_simple",
	createdAt:                pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt:           pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:                    "Collapsed Simple",
	Description:              "DESCRIPTION",
	FabricConnectivityDesign: enum.FabricConnectivityDesignL3Collapsed,
	Status:                   pointer.To(enum.FFEConsistencyStatusInconsistent),
	LeafSwitches: []LeafSwitch{
		{
			Label: "leafy",
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 16,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
					},
				},
			},
			Tags: []Tag{
				{Label: "tag a", Description: "TAG A"},
				{Label: "tag b", Description: "TAG B"},
			},
		},
	},
}

const rackTypeTestCollapsedSimpleJSON = `{
  "id": "id__collapsed_simple",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "display_name": "Collapsed Simple",
  "description": "DESCRIPTION",
  "fabric_connectivity_design": "l3collapsed",
  "status": "inconsistent",
  "tags": [ { "label": "tag a", "description": "TAG A" }, { "label": "tag b", "description": "TAG B" } ],
  "logical_devices": [
    {
      "id": "a0e801da44420b1a0c915fc21f81d2d5",
      "display_name": "ld label",
      "panels": [
        {
          "panel_layout": { "row_count": 2, "column_count": 8 },
          "port_indexing": { "order": "T-B, L-R", "schema": "absolute", "start_index": 1 },
          "port_groups": [ { "count": 16, "speed": { "unit": "G", "value": 100 }, "roles": [ "spine", "leaf", "access", "generic" ] } ]
        }
      ]
    }
  ],
  "leafs": [
    { "label": "leafy", "logical_device": "a0e801da44420b1a0c915fc21f81d2d5", "tags": [ "tag a", "tag b" ] }
  ]
}`

var rackTypeTestCollapsedESI = RackType{
	id:                       "id__collapsed_esi",
	createdAt:                pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt:           pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:                    "Collapsed ESI",
	Description:              "DESCRIPTION",
	FabricConnectivityDesign: enum.FabricConnectivityDesignL3Collapsed,
	Status:                   pointer.To(enum.FFEConsistencyStatusInconsistent),
	LeafSwitches: []LeafSwitch{
		{
			Label: "leafy",
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 16,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
					},
				},
			},
			Tags: []Tag{
				{Label: "tag a", Description: "TAG A"},
				{Label: "tag b", Description: "TAG B"},
			},
			RedundancyProtocol: enum.LeafRedundancyProtocolESI,
		},
	},
}

const rackTypeTestCollapsedESIJSON = `{
  "id": "id__collapsed_esi",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "display_name": "Collapsed ESI",
  "description": "DESCRIPTION",
  "fabric_connectivity_design": "l3collapsed",
  "status": "inconsistent",
  "tags": [ { "label": "tag a", "description": "TAG A" }, { "label": "tag b", "description": "TAG B" } ],
  "logical_devices": [
    {
      "id": "a0e801da44420b1a0c915fc21f81d2d5",
      "display_name": "ld label",
      "panels": [
        {
          "panel_layout": { "row_count": 2, "column_count": 8 },
          "port_indexing": { "order": "T-B, L-R", "schema": "absolute", "start_index": 1 },
          "port_groups": [ { "count": 16, "speed": { "unit": "G", "value": 100 }, "roles": [ "spine", "leaf", "access", "generic" ] } ]
        }
      ]
    }
  ],
  "leafs": [
    { "label": "leafy", "logical_device": "a0e801da44420b1a0c915fc21f81d2d5", "redundancy_protocol": "esi", "tags": [ "tag a", "tag b" ] }
  ]
}`

var rackTypeTestCollapsedSimpleWithAccess = RackType{
	id:                       "id__collapsed_simple_with_access",
	createdAt:                pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt:           pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:                    "Collapsed Simple With Access",
	Description:              "DESCRIPTION",
	FabricConnectivityDesign: enum.FabricConnectivityDesignL3Collapsed,
	Status:                   pointer.To(enum.FFEConsistencyStatusInconsistent),
	LeafSwitches: []LeafSwitch{
		{
			Label: "leafy",
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 16,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
					},
				},
			},
			Tags: []Tag{
				{Label: "tag a", Description: "TAG A"},
				{Label: "tag b", Description: "TAG B"},
			},
		},
	},
	AccessSwitches: []AccessSwitch{
		{
			Count: 1,
			Label: "accessy",
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 2},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 4,
								Speed: "400G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
					},
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 32},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 16,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
							{
								Count: 16,
								Speed: "50G",
								Roles: LogicalDevicePortRoles{enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
					},
				},
			},
			Tags: []Tag{
				{Label: "tag c", Description: "TAG C"},
				{Label: "tag d", Description: "TAG D"},
			},
			Links: []RackTypeLink{
				{
					Label:              "linky",
					TargetSwitchLabel:  "leafy",
					LinkPerSwitchCount: 1,
					Speed:              "100G",
					Tags: []Tag{
						{Label: "tag e", Description: "TAG E"},
						{Label: "tag f", Description: "TAG F"},
					},
				},
			},
		},
	},
}

const rackTypeTestCollapsedSimpleWithAccessJSON = `{
  "id": "id__collapsed_simple_with_access",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "display_name": "Collapsed Simple With Access",
  "description": "DESCRIPTION",
  "fabric_connectivity_design": "l3collapsed",
  "status": "inconsistent",
  "tags": [
    { "label": "tag a", "description": "TAG A" }, { "label": "tag b", "description": "TAG B" },
    { "label": "tag c", "description": "TAG C" }, { "label": "tag d", "description": "TAG D" },
    { "label": "tag e", "description": "TAG E" }, { "label": "tag f", "description": "TAG F" }
  ],
  "logical_devices": [
    {
      "id": "707591024f5e44a088cb5ba03a6c33bf",
      "display_name": "ld label",
      "panels": [
        {
          "panel_layout": { "row_count": 2, "column_count": 2 },
          "port_indexing": { "order": "L-R, T-B", "schema": "absolute", "start_index": 1 },
          "port_groups": [ { "count": 4, "speed": { "unit": "G", "value": 400 }, "roles": [ "spine", "leaf", "access", "generic" ] } ]
        },
        {
          "panel_layout": { "row_count": 1, "column_count": 32 },
          "port_indexing": { "order": "L-R, T-B", "schema": "absolute", "start_index": 1 },
          "port_groups": [
            { "count": 16, "speed": { "unit": "G", "value": 100 }, "roles": [ "spine", "leaf", "access", "generic" ] },
            { "count": 16, "speed": { "unit": "G", "value": 50 }, "roles": [ "generic" ] }
          ]
        }
      ]
    },
    {
      "id": "a0e801da44420b1a0c915fc21f81d2d5",
      "display_name": "ld label",
      "panels": [
        {
          "panel_layout": { "row_count": 2, "column_count": 8 },
          "port_indexing": { "order": "T-B, L-R", "schema": "absolute", "start_index": 1 },
          "port_groups": [ { "count": 16, "speed": { "unit": "G", "value": 100 }, "roles": [ "spine", "leaf", "access", "generic" ] } ]
        }
      ]
    }
  ],
  "leafs": [
    { "label": "leafy", "logical_device": "a0e801da44420b1a0c915fc21f81d2d5", "tags": [ "tag a", "tag b" ] }
  ],
  "access_switches": [
    {
      "instance_count": 1,
      "label": "accessy",
	  "links" : [
        {
          "label" : "linky",
          "target_switch_label" : "leafy",
          "link_per_switch_count" : 1,
          "link_speed" : { "unit" : "G", "value" : 100 },
          "attachment_type" : "singleAttached",
          "tags" : [ "tag e", "tag f" ]
        }
      ],
      "logical_device": "707591024f5e44a088cb5ba03a6c33bf",
      "tags": [ "tag c", "tag d" ],
      "access_access_link_count": 0,
      "access_access_link_speed": null,
      "access_access_link_port_channel_id_max": 0,
      "access_access_link_port_channel_id_min": 0
    }
  ]
}`

var rackTypeTestRackBasedESIWithAccessESI = RackType{
	id:                       "id__rack_based_esi_with_access_esi",
	createdAt:                pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt:           pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:                    "Collapsed Simple With Access",
	Description:              "DESCRIPTION",
	FabricConnectivityDesign: enum.FabricConnectivityDesignL3Collapsed,
	Status:                   pointer.To(enum.FFEConsistencyStatusInconsistent),
	LeafSwitches: []LeafSwitch{
		{
			Label:              "leafy",
			RedundancyProtocol: enum.LeafRedundancyProtocolESI,
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 16,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
					},
				},
			},
			Tags: []Tag{
				{Label: "tag a", Description: "TAG A"},
				{Label: "tag b", Description: "TAG B"},
			},
		},
	},
	AccessSwitches: []AccessSwitch{
		{
			Count: 1,
			Label: "accessy",
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 2},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 4,
								Speed: "400G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
					},
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 32},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 16,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
							{
								Count: 16,
								Speed: "50G",
								Roles: LogicalDevicePortRoles{enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
					},
				},
			},
			ESILAGInfo: &RackTypeAccessSwitchESILAGInfo{
				LinkCount:        2,
				LinkSpeed:        "100G",
				PortChannelIdMax: 10,
				PortChannelIdMin: 20,
			},
			Tags: []Tag{
				{Label: "tag c", Description: "TAG C"},
				{Label: "tag d", Description: "TAG D"},
			},
			Links: []RackTypeLink{
				{
					Label:              "linky",
					TargetSwitchLabel:  "leafy",
					LinkPerSwitchCount: 2,
					Speed:              "100G",
					AttachmentType:     enum.LinkAttachmentTypeDual,
					Tags: []Tag{
						{Label: "tag e", Description: "TAG E"},
						{Label: "tag f", Description: "TAG F"},
					},
				},
			},
		},
	},
}

const rackTypeTestRackBasedESIWithAccessESIJSON = `{
  "id": "id__rack_based_esi_with_access_esi",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "display_name": "Collapsed Simple With Access",
  "description": "DESCRIPTION",
  "fabric_connectivity_design": "l3collapsed",
  "status": "inconsistent",
  "tags": [
    { "label": "tag a", "description": "TAG A" }, { "label": "tag b", "description": "TAG B" },
    { "label": "tag c", "description": "TAG C" }, { "label": "tag d", "description": "TAG D" },
    { "label": "tag e", "description": "TAG E" }, { "label": "tag f", "description": "TAG F" }
  ],
  "logical_devices": [
    {
      "id": "707591024f5e44a088cb5ba03a6c33bf",
      "display_name": "ld label",
      "panels": [
        {
          "panel_layout": { "row_count": 2, "column_count": 2 },
          "port_indexing": { "order": "L-R, T-B", "schema": "absolute", "start_index": 1 },
          "port_groups": [ { "count": 4, "speed": { "unit": "G", "value": 400 }, "roles": [ "spine", "leaf", "access", "generic" ] } ]
        },
        {
          "panel_layout": { "row_count": 1, "column_count": 32 },
          "port_indexing": { "order": "L-R, T-B", "schema": "absolute", "start_index": 1 },
          "port_groups": [
            { "count": 16, "speed": { "unit": "G", "value": 100 }, "roles": [ "spine", "leaf", "access", "generic" ] },
            { "count": 16, "speed": { "unit": "G", "value": 50 }, "roles": [ "generic" ] }
          ]
        }
      ]
    },
    {
      "id": "a0e801da44420b1a0c915fc21f81d2d5",
      "display_name": "ld label",
      "panels": [
        {
          "panel_layout": { "row_count": 2, "column_count": 8 },
          "port_indexing": { "order": "T-B, L-R", "schema": "absolute", "start_index": 1 },
          "port_groups": [ { "count": 16, "speed": { "unit": "G", "value": 100 }, "roles": [ "spine", "leaf", "access", "generic" ] } ]
        }
      ]
    }
  ],
  "leafs": [
    { "label": "leafy", "redundancy_protocol": "esi", "logical_device": "a0e801da44420b1a0c915fc21f81d2d5", "tags": [ "tag a", "tag b" ] }
  ],
  "access_switches": [
    {
      "instance_count": 1,
      "label": "accessy",
      "redundancy_protocol": "esi",
	  "links" : [
        {
          "label" : "linky",
          "target_switch_label" : "leafy",
          "link_per_switch_count" : 2,
          "link_speed" : { "unit" : "G", "value" : 100 },
          "attachment_type" : "dualAttached",
          "tags" : [ "tag e", "tag f" ]
        }
      ],
      "logical_device": "707591024f5e44a088cb5ba03a6c33bf",
      "tags": [ "tag c", "tag d" ],
      "access_access_link_count": 2,
      "access_access_link_speed": {
        "unit": "G",
        "value": 100
      },
      "access_access_link_port_channel_id_max": 10,
      "access_access_link_port_channel_id_min": 20
    }
  ]
}`

var rackTypeTestRackBasedMLAGWithAccessPair = RackType{
	id:                       "id__rack_based_mlag_with_access_pair",
	createdAt:                pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt:           pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:                    "Collapsed Simple With Access",
	Description:              "DESCRIPTION",
	FabricConnectivityDesign: enum.FabricConnectivityDesignL3Collapsed,
	Status:                   pointer.To(enum.FFEConsistencyStatusInconsistent),
	LeafSwitches: []LeafSwitch{
		{
			Label:              "leafy",
			RedundancyProtocol: enum.LeafRedundancyProtocolMLAG,
			MLAGInfo: &RackTypeLeafSwitchMLAGInfo{
				LeafLeafL3LinkCount:         2,
				LeafLeafL3LinkSpeed:         "100G",
				LeafLeafL3LinkPortChannelId: 3,
				LeafLeafLinkCount:           4,
				LeafLeafLinkSpeed:           "100G",
				LeafLeafLinkPortChannelId:   5,
				MLAGVLAN:                    6,
			},
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 16,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
					},
				},
			},
			Tags: []Tag{
				{Label: "tag a", Description: "TAG A"},
				{Label: "tag b", Description: "TAG B"},
			},
		},
	},
	AccessSwitches: []AccessSwitch{
		{
			Count: 1,
			Label: "accessa",
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 2},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 4,
								Speed: "400G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
					},
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 32},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 16,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
							{
								Count: 16,
								Speed: "50G",
								Roles: LogicalDevicePortRoles{enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
					},
				},
			},
			Tags: []Tag{
				{Label: "tag c", Description: "TAG C"},
				{Label: "tag d", Description: "TAG D"},
			},
			Links: []RackTypeLink{
				{
					Label:              "linky",
					TargetSwitchLabel:  "leafy",
					LinkPerSwitchCount: 2,
					Speed:              "100G",
					AttachmentType:     enum.LinkAttachmentTypeSingle,
					SwitchPeer:         enum.LinkSwitchPeerFirst,
					Tags: []Tag{
						{Label: "tag e", Description: "TAG E"},
						{Label: "tag f", Description: "TAG F"},
					},
				},
			},
		},
		{
			Count: 1,
			Label: "accessb",
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 2},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 4,
								Speed: "400G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
					},
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 32},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 16,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
							{
								Count: 16,
								Speed: "50G",
								Roles: LogicalDevicePortRoles{enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
					},
				},
			},
			Tags: []Tag{
				{Label: "tag g", Description: "TAG G"},
				{Label: "tag h", Description: "TAG H"},
			},
			Links: []RackTypeLink{
				{
					Label:              "linky",
					TargetSwitchLabel:  "leafy",
					LinkPerSwitchCount: 3,
					Speed:              "100G",
					AttachmentType:     enum.LinkAttachmentTypeSingle,
					SwitchPeer:         enum.LinkSwitchPeerSecond,
					Tags: []Tag{
						{Label: "tag i", Description: "TAG I"},
						{Label: "tag j", Description: "TAG J"},
					},
				},
			},
		},
	},
}

const rackTypeTestRackBasedMLAGWithAccessPairJSON = `{
  "id": "id__rack_based_mlag_with_access_pair",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "display_name": "Collapsed Simple With Access",
  "description": "DESCRIPTION",
  "fabric_connectivity_design": "l3collapsed",
  "status": "inconsistent",
  "tags": [
    { "label": "tag a", "description": "TAG A" }, { "label": "tag b", "description": "TAG B" },
    { "label": "tag c", "description": "TAG C" }, { "label": "tag d", "description": "TAG D" },
    { "label": "tag e", "description": "TAG E" }, { "label": "tag f", "description": "TAG F" },
    { "label": "tag g", "description": "TAG G" }, { "label": "tag h", "description": "TAG H" },
    { "label": "tag i", "description": "TAG I" }, { "label": "tag j", "description": "TAG J" }
  ],
  "logical_devices": [
    {
      "id": "707591024f5e44a088cb5ba03a6c33bf",
      "display_name": "ld label",
      "panels": [
        {
          "panel_layout": { "row_count": 2, "column_count": 2 },
          "port_indexing": { "order": "L-R, T-B", "schema": "absolute", "start_index": 1 },
          "port_groups": [ { "count": 4, "speed": { "unit": "G", "value": 400 }, "roles": [ "spine", "leaf", "access", "generic" ] } ]
        },
        {
          "panel_layout": { "row_count": 1, "column_count": 32 },
          "port_indexing": { "order": "L-R, T-B", "schema": "absolute", "start_index": 1 },
          "port_groups": [
            { "count": 16, "speed": { "unit": "G", "value": 100 }, "roles": [ "spine", "leaf", "access", "generic" ] },
            { "count": 16, "speed": { "unit": "G", "value": 50 }, "roles": [ "generic" ] }
          ]
        }
      ]
    },
    {
      "id": "a0e801da44420b1a0c915fc21f81d2d5",
      "display_name": "ld label",
      "panels": [
        {
          "panel_layout": { "row_count": 2, "column_count": 8 },
          "port_indexing": { "order": "T-B, L-R", "schema": "absolute", "start_index": 1 },
          "port_groups": [ { "count": 16, "speed": { "unit": "G", "value": 100 }, "roles": [ "spine", "leaf", "access", "generic" ] } ]
        }
      ]
    }
  ],
  "leafs": [
      {
        "label": "leafy", "redundancy_protocol": "mlag", "logical_device": "a0e801da44420b1a0c915fc21f81d2d5", "tags": [ "tag a", "tag b" ],
        "leaf_leaf_l3_link_count": 2, "leaf_leaf_l3_link_port_channel_id": 3, "leaf_leaf_l3_link_speed": { "unit": "G", "value": 100 },
        "leaf_leaf_link_count": 4, "leaf_leaf_link_port_channel_id": 5, "leaf_leaf_link_speed": { "unit": "G", "value": 100 }, "mlag_vlan_id": 6
      }
  ],
  "access_switches": [
    {
      "instance_count": 1,
      "label": "accessa",
	  "links" : [
        {
          "label" : "linky",
          "target_switch_label" : "leafy",
          "link_per_switch_count" : 2,
          "link_speed" : { "unit" : "G", "value" : 100 },
          "attachment_type" : "singleAttached",
          "switch_peer": "first",
          "tags" : [ "tag e", "tag f" ]
        }
      ],
      "logical_device": "707591024f5e44a088cb5ba03a6c33bf",
      "tags": [ "tag c", "tag d" ],
      "access_access_link_count": 0,
      "access_access_link_speed": null,
      "access_access_link_port_channel_id_max": 0,
      "access_access_link_port_channel_id_min": 0
    },
    {
      "instance_count": 1,
      "label": "accessb",
	  "links" : [
        {
          "label" : "linky",
          "target_switch_label" : "leafy",
          "link_per_switch_count" : 3,
          "link_speed" : { "unit" : "G", "value" : 100 },
          "attachment_type" : "singleAttached",
          "switch_peer": "second",
          "tags" : [ "tag i", "tag j" ]
        }
      ],
      "logical_device": "707591024f5e44a088cb5ba03a6c33bf",
      "tags": [ "tag g", "tag h" ],
      "access_access_link_count": 0,
      "access_access_link_speed": null,
      "access_access_link_port_channel_id_max": 0,
      "access_access_link_port_channel_id_min": 0
    }
  ]
}`

var rackTypeTestCollapsedESIWithGenericSystems = RackType{
	id:                       "id__collapsed_esi_with_generic_systems",
	createdAt:                pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt:           pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:                    "Collapsed ESI",
	Description:              "DESCRIPTION",
	FabricConnectivityDesign: enum.FabricConnectivityDesignL3Collapsed,
	Status:                   pointer.To(enum.FFEConsistencyStatusInconsistent),
	LeafSwitches: []LeafSwitch{
		{
			Label: "leafy",
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 16,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleLeaf, enum.PortRoleAccess, enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
					},
				},
			},
			Tags:               []Tag{{Label: "tag a", Description: "TAG A"}, {Label: "tag b", Description: "TAG B"}},
			RedundancyProtocol: enum.LeafRedundancyProtocolESI,
		},
	},
	GenericSystems: []GenericSystem{
		{
			AsnDomain: &enum.FeatureSwitchEnabled,
			Count:     2,
			Label:     "lefty",
			Links: []RackTypeLink{
				{
					Label:              "left",
					TargetSwitchLabel:  "leafy",
					LinkPerSwitchCount: 2,
					Speed:              "100G",
					AttachmentType:     enum.LinkAttachmentTypeSingle,
					LAGMode:            enum.LAGModePassiveLACP,
					SwitchPeer:         enum.LinkSwitchPeerFirst,
					Tags: []Tag{
						{Label: "tag e", Description: "TAG E"}, {Label: "tag f", Description: "TAG F"},
					},
				},
			},
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 1},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 1,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleLeaf},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
					},
				},
			},
			Loopback:         &enum.FeatureSwitchEnabled,
			ManagementLevel:  enum.SystemManagementLevelUnmanaged,
			PortChannelIDMax: 10,
			PortChannelIDMin: 19,
			Tags:             []Tag{{Label: "tag c", Description: "TAG C"}, {Label: "tag d", Description: "TAG D"}},
		},
		{
			AsnDomain: &enum.FeatureSwitchDisabled,
			Count:     1,
			Label:     "dually",
			Links: []RackTypeLink{
				{
					Label:              "dually",
					TargetSwitchLabel:  "leafy",
					LinkPerSwitchCount: 1,
					Speed:              "100G",
					AttachmentType:     enum.LinkAttachmentTypeDual,
					LAGMode:            enum.LAGModeStatic,
					Tags:               []Tag{{Label: "tag i", Description: "TAG I"}, {Label: "tag j", Description: "TAG J"}},
				},
			},
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 2},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 2,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
					},
				},
			},
			Loopback:         &enum.FeatureSwitchDisabled,
			ManagementLevel:  enum.SystemManagementLevelUnmanaged,
			PortChannelIDMax: 20,
			PortChannelIDMin: 29,
			Tags:             []Tag{{Label: "tag g", Description: "TAG G"}, {Label: "tag h", Description: "TAG H"}},
		},
		{
			AsnDomain: nil,
			Count:     3,
			Label:     "righty",
			Links: []RackTypeLink{
				{
					Label:              "right",
					TargetSwitchLabel:  "leafy",
					LinkPerSwitchCount: 3,
					Speed:              "100G",
					AttachmentType:     enum.LinkAttachmentTypeSingle,
					LAGMode:            enum.LAGModeActiveLACP,
					SwitchPeer:         enum.LinkSwitchPeerSecond,
					Tags:               []Tag{{Label: "tag m", Description: "TAG M"}, {Label: "tag n", Description: "TAG N"}},
				},
			},
			LogicalDevice: LogicalDevice{
				Label: "ld label",
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 1},
						PortGroups: []LogicalDevicePanelPortGroup{
							{
								Count: 1,
								Speed: "100G",
								Roles: LogicalDevicePortRoles{enum.PortRoleLeaf},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
					},
				},
			},
			ManagementLevel:  enum.SystemManagementLevelUnmanaged,
			PortChannelIDMax: 30,
			PortChannelIDMin: 39,
			Tags:             []Tag{{Label: "tag k", Description: "TAG K"}, {Label: "tag l", Description: "TAG L"}},
		},
	},
}

const rackTypeTestCollapsedESIWithGenericSystemsJSON = `{
  "id": "id__collapsed_esi_with_generic_systems",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "display_name": "Collapsed ESI",
  "description": "DESCRIPTION",
  "fabric_connectivity_design": "l3collapsed",
  "status": "inconsistent",
  "tags": [
    { "label": "tag a", "description": "TAG A" }, { "label": "tag b", "description": "TAG B" },
    { "label": "tag c", "description": "TAG C" }, { "label": "tag d", "description": "TAG D" },
    { "label": "tag e", "description": "TAG E" }, { "label": "tag f", "description": "TAG F" },
    { "label": "tag g", "description": "TAG G" }, { "label": "tag h", "description": "TAG H" },
    { "label": "tag i", "description": "TAG I" }, { "label": "tag j", "description": "TAG J" },
    { "label": "tag k", "description": "TAG K" }, { "label": "tag l", "description": "TAG L" },
    { "label": "tag m", "description": "TAG M" }, { "label": "tag n", "description": "TAG N" }
  ],
  "logical_devices": [
    {
      "display_name": "ld label",
      "id": "3c4cf87f95f4507d72f89459a0116981",
      "panels": [
        {
          "panel_layout": { "column_count": 2, "row_count": 1 },
          "port_groups": [ { "count": 2, "roles": [ "leaf", "access" ], "speed": { "unit": "G", "value": 100 } } ],
          "port_indexing": { "order": "T-B, L-R", "schema": "absolute", "start_index": 1 }
        }
      ]
    },
    {
      "display_name": "ld label",
      "id": "9fc7b1730dda8dad1e5bb4c877647a7e",
      "panels": [
        {
          "panel_layout": { "column_count": 1, "row_count": 1 },
          "port_groups": [ { "count": 1, "roles": [ "leaf" ], "speed": { "unit": "G", "value": 100 } } ],
          "port_indexing": { "order": "T-B, L-R", "schema": "absolute", "start_index": 1 }
        }
      ]
    },
    {
      "id": "a0e801da44420b1a0c915fc21f81d2d5",
      "display_name": "ld label",
      "panels": [
        {
          "panel_layout": { "row_count": 2, "column_count": 8 },
          "port_indexing": { "order": "T-B, L-R", "schema": "absolute", "start_index": 1 },
          "port_groups": [ { "count": 16, "speed": { "unit": "G", "value": 100 }, "roles": [ "spine", "leaf", "access", "generic" ] } ]
        }
      ]
    }
  ],
  "leafs": [
    { "label": "leafy", "logical_device": "a0e801da44420b1a0c915fc21f81d2d5", "redundancy_protocol": "esi", "tags": [ "tag a", "tag b" ] }
  ],
  "generic_systems": [
    {
      "asn_domain": "enabled",
      "count": 2,
      "label": "lefty",
      "links": [
        {
          "label" : "left",
          "target_switch_label" : "leafy",
          "link_per_switch_count" : 2,
          "link_speed" : { "unit" : "G", "value" : 100 },
          "attachment_type" : "singleAttached",
          "switch_peer": "first",
          "tags" : [ "tag e", "tag f" ],
          "lag_mode": "lacp_passive"
        }
      ],
      "logical_device": "9fc7b1730dda8dad1e5bb4c877647a7e",
      "loopback": "enabled",
      "management_level": "unmanaged",
      "port_channel_id_max": 10,
      "port_channel_id_min": 19,
      "tags": [ "tag c", "tag d" ]
    },
    {
      "asn_domain": "disabled",
      "count": 1,
      "label": "dually",
      "links": [
        {
          "label" : "dually",
          "target_switch_label" : "leafy",
          "link_per_switch_count" : 1,
          "link_speed" : { "unit" : "G", "value" : 100 },
          "attachment_type" : "dualAttached",
          "tags" : [ "tag i", "tag j" ],
          "lag_mode": "static_lag"
        }
      ],
      "logical_device": "3c4cf87f95f4507d72f89459a0116981",
      "loopback": "disabled",
      "management_level": "unmanaged",
      "port_channel_id_max": 20,
      "port_channel_id_min": 29,
      "tags": [ "tag g", "tag h" ]
    },
    {
      "count": 3,
      "label": "righty",
      "links": [
        {
          "label" : "right",
          "target_switch_label" : "leafy",
          "link_per_switch_count" : 3,
          "link_speed" : { "unit" : "G", "value" : 100 },
          "attachment_type" : "singleAttached",
          "switch_peer": "second",
          "tags" : [ "tag m", "tag n" ],
          "lag_mode": "lacp_active"
        }
      ],
      "logical_device": "9fc7b1730dda8dad1e5bb4c877647a7e",
      "management_level": "unmanaged",
      "port_channel_id_max": 30,
      "port_channel_id_min": 39,
      "tags": [ "tag k", "tag l" ]
    }
  ]
}`
