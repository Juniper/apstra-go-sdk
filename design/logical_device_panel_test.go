// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func TestLogicalDevicePanel_MarshalUnmarshal(t *testing.T) {
	type testCase struct {
		panel LogicalDevicePanel
		json  string
	}

	var allRoles LogicalDevicePortRoles
	allRoles.IncludeAllUses()

	testCases := map[string]testCase{
		"a": {
			panel: LogicalDevicePanel{
				PanelLayout: LogicalDevicePanelLayout{
					RowCount:    2,
					ColumnCount: 4,
				},
				PortGroups: []LogicalDevicePanelPortGroup{
					{
						Count: 8,
						Speed: "100G",
						Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleLeaf},
					},
				},
				PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
			},
			json: `{"panel_layout":{"row_count":2,"column_count":4},"port_indexing":{"order":"T-B, L-R","schema":"absolute","start_index":1},"port_groups":[{"count":8,"speed":{"unit":"G","value":100},"roles":["superspine","leaf"]}]}`,
		},
		"b": {
			panel: LogicalDevicePanel{
				PanelLayout: LogicalDevicePanelLayout{
					RowCount:    4,
					ColumnCount: 8,
				},
				PortGroups: []LogicalDevicePanelPortGroup{
					{
						Count: 8,
						Speed: "100G",
						Roles: LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleLeaf},
					},
					{
						Count: 24,
						Speed: "25G",
						Roles: LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess},
					},
				},
				PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
			},
			json: `{"panel_layout":{"row_count":4,"column_count":8},"port_indexing":{"order":"L-R, T-B","schema":"absolute","start_index":1},"port_groups":[{"count":8,"speed":{"unit":"G","value":100},"roles":["superspine","leaf"]},{"count":24,"speed":{"unit":"G","value":25},"roles":["generic","access"]}]}`,
		},
		"c": {
			panel: LogicalDevicePanel{
				PanelLayout: LogicalDevicePanelLayout{
					RowCount:    1,
					ColumnCount: 1,
				},
				PortGroups: []LogicalDevicePanelPortGroup{
					{
						Count: 1,
						Speed: "10M",
						Roles: allRoles,
					},
				},
				PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
			},
			json: `{"panel_layout":{"row_count":1,"column_count":1},"port_indexing":{"order":"T-B, L-R","schema":"absolute","start_index":1},"port_groups":[{"count":1,"speed":{"unit":"M","value":10},"roles":["access", "generic", "leaf", "peer", "spine", "superspine", "unused"]}]}`,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			jsonResult, err := tCase.panel.MarshalJSON()
			require.NoError(t, err)

			var panelResult LogicalDevicePanel
			err = panelResult.UnmarshalJSON([]byte(tCase.json))
			require.NoError(t, err)

			require.JSONEq(t, tCase.json, string(jsonResult))
			require.Equal(t, tCase.panel, panelResult)
		})
	}
}
