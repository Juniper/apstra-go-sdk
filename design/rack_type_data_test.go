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
          "speed" : { "unit" : "G", "value" : 100 },
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
