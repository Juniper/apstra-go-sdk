// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var logicalDeviceTest1x1 = LogicalDevice{
	id:             "id__test-1x1",
	createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:          "test-1x1",
	Panels: []LogicalDevicePanel{
		{
			PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 1},
			PortGroups: []LogicalDevicePanelPortGroup{
				{
					Count: 1,
					Speed: "1G",
					Roles: LogicalDevicePortRoles{enum.PortRoleAccess, enum.PortRoleLeaf},
				},
			},
			PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
		},
	},
}

const logicalDeviceTest1x1JSON = `{
  "id": "id__test-1x1",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "display_name": "test-1x1",
  "panels": [
    {
      "panel_layout": { "row_count": 1, "column_count": 1 },
      "port_indexing": {
        "schema": "absolute",
        "order": "T-B, L-R",
        "start_index": 1
      },
      "port_groups": [
        {
          "count": 1,
          "speed": { "value": 1, "unit": "G" },
          "roles": [ "access", "leaf" ]
        }
      ]
    }
  ]
}`

var logicalDeviceTest48x10plus4x100 = LogicalDevice{
	id:             "id__test-48x10+4x100",
	createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:          "test-48x10+4x100",
	Panels: []LogicalDevicePanel{
		{
			PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 24},
			PortGroups: []LogicalDevicePanelPortGroup{
				{
					Count: 48,
					Speed: "10G",
					Roles: LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess},
				},
			},
			PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
		},
		{
			PanelLayout: LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 2},
			PortGroups: []LogicalDevicePanelPortGroup{
				{
					Count: 2,
					Speed: "100G",
					Roles: LogicalDevicePortRoles{enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
				},
				{
					Count: 2,
					Speed: "100G",
					Roles: LogicalDevicePortRoles{enum.PortRolePeer},
				},
			},
			PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
		},
	},
}

const logicalDeviceTest48x10plus4x100JSON = `{
  "id": "id__test-48x10+4x100",
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "display_name": "test-48x10+4x100",
  "panels": [
    {
      "panel_layout": { "row_count": 2, "column_count": 24 },
      "port_indexing": {
        "schema": "absolute",
        "order": "T-B, L-R",
        "start_index": 1
      },
      "port_groups": [
        {
          "count": 48,
          "speed": { "value": 10, "unit": "G" },
          "roles": [ "generic", "peer", "access" ]
        }
      ]
    },
    {
      "panel_layout": { "row_count": 2, "column_count": 2 },
      "port_indexing": {
        "schema": "absolute",
        "order": "L-R, T-B",
        "start_index": 1
      },
      "port_groups": [
        {
          "count": 2,
          "speed": { "value": 100, "unit": "G" },
          "roles": [ "peer", "access", "spine" ]
        },
        {
          "count": 2,
          "speed": { "value": 100, "unit": "G" },
          "roles": [ "peer" ]
        }
      ]
    }
  ]
}`
