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
	Label:                    "LABEL",
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
				{Label: "tag b", Description: "TAG B"},
				{Label: "tag a", Description: "TAG A"},
			},
		},
	},
}

const rackTypeTestCollapsedSimpleJSON = `{
  "id": "id__collapsed_simple",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "display_name": "LABEL",
  "description": "DESCRIPTION",
  "fabric_connectivity_design": "l3collapsed",
  "status": "inconsistent",
  "tags": [ { "label": "tag a", "description": "TAG A" }, { "label": "tag b", "description": "TAG B" }
  ],
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
	Label:                    "LABEL",
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
				{Label: "tag b", Description: "TAG B"},
				{Label: "tag a", Description: "TAG A"},
			},
			RedundancyProtocol: enum.LeafRedundancyProtocolESI,
		},
	},
}

const rackTypeTestCollapsedESIJSON = `{
  "id": "id__collapsed_esi",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "display_name": "LABEL",
  "description": "DESCRIPTION",
  "fabric_connectivity_design": "l3collapsed",
  "status": "inconsistent",
  "tags": [ { "label": "tag a", "description": "TAG A" }, { "label": "tag b", "description": "TAG B" }
  ],
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
