// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package compare

import (
	"fmt"
	"testing"

	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/internal/str"
	"github.com/stretchr/testify/require"
)

func LogicalDevice(t testing.TB, a, b design.LogicalDevice) {
	t.Helper()

	if a.ID() != nil && b.ID() != nil {
		require.Equal(t, *a.ID(), *b.ID(), fmt.Sprintf("%s: IDs do not match", str.FuncName()))
	}

	require.Equal(t, len(a.Panels), len(b.Panels), fmt.Sprintf("%s: Panel counts do not match", str.FuncName()))
	for i := range a.Panels {
		LogicalDevicePanel(t, a.Panels[i], b.Panels[i])
	}

	require.Equal(t, a.Label, b.Label, fmt.Sprintf("%s: Labels do not match", str.FuncName()))
}

func LogicalDevicePanel(t testing.TB, a, b design.LogicalDevicePanel) {
	t.Helper()

	LogicalDevicePanelLayout(t, a.PanelLayout, b.PanelLayout)

	require.Equal(t, len(a.PortGroups), len(b.PortGroups), fmt.Sprintf("%s: PortGroup counts do not match", str.FuncName()))
	for i := range a.PortGroups {
		LogicalDevicePanelPortGroup(t, a.PortGroups[i], b.PortGroups[i])
	}

	require.Equal(t, a.PortIndexing, b.PortIndexing, fmt.Sprintf("%s: PortIndexings do not match", str.FuncName()))
}

func LogicalDevicePanelLayout(t testing.TB, a, b design.LogicalDevicePanelLayout) {
	t.Helper()

	require.Equal(t, a.ColumnCount, b.ColumnCount, fmt.Sprintf("%s: ColumnCounts do not match", str.FuncName()))
	require.Equal(t, a.RowCount, b.RowCount, fmt.Sprintf("%s: RowCounts do not match", str.FuncName()))
}

func LogicalDevicePanelPortGroup(t testing.TB, a, b design.LogicalDevicePanelPortGroup) {
	t.Helper()

	require.Equal(t, a.Count, b.Count, fmt.Sprintf("%s: Counts do not match", str.FuncName()))
	require.Equal(t, a.Speed, b.Speed, fmt.Sprintf("%s: Speeds do not match", str.FuncName()))
	require.ElementsMatch(t, a.Roles.Strings(), b.Roles.Strings(), fmt.Sprintf("%s: Roles do not match", str.FuncName()))
}
