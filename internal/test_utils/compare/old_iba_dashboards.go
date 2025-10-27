// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package compare

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
	"github.com/stretchr/testify/require"
)

func Dashboards(t testing.TB, d1, d2 apstra.IbaDashboardData) {
	t.Helper()

	defaultTimeSeriesDuration := 86400
	defaultAggregationPeriod := 300

	require.Equal(t, d1.Label, d2.Label)
	require.Equal(t, d1.Default, d2.Default)
	require.Equal(t, d1.PredefinedDashboard, d2.PredefinedDashboard)

	for k1, v1 := range d1.IbaWidgetGrid {
		for k2, v2 := range v1 {
			require.Equal(t, v2.Label, d2.IbaWidgetGrid[k1][k2].Label)
			require.Equal(t, v2.Description, d2.IbaWidgetGrid[k1][k2].Description)
			require.Equal(t, v2.ProbeId, d2.IbaWidgetGrid[k1][k2].ProbeId)

			if v2.TimeSeriesDuration == nil {
				v2.TimeSeriesDuration = timeutils.NewDurationInSecs(defaultTimeSeriesDuration)
			}
			require.Equal(t, v2.TimeSeriesDuration.TimeinSecs(), d2.IbaWidgetGrid[k1][k2].TimeSeriesDuration.TimeinSecs())

			if v2.AggregationPeriod == nil {
				v2.AggregationPeriod = timeutils.NewDurationInSecs(defaultAggregationPeriod)
			}
			require.Equal(t, v2.AggregationPeriod.TimeinSecs(), d2.IbaWidgetGrid[k1][k2].AggregationPeriod.TimeinSecs())
		}
	}
}
