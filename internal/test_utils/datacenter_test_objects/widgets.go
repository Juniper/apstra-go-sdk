package dctestobj

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
	"github.com/stretchr/testify/require"
)

// TestWidgetsABC instantiates two predefined probes and creates widgets from them,
// returning the widget Object Id and the IbaWidget object used for creation
func TestWidgetsABC(t testing.TB, ctx context.Context, bpClient *apstra.TwoStageL3ClosClient) (apstra.IbaWidget, apstra.IbaWidget, apstra.IbaWidget) {
	t.Helper()

	probeAId, err := bpClient.InstantiateIbaPredefinedProbe(ctx, &apstra.IbaPredefinedProbeRequest{
		Name: "bgp_session",
		Data: []byte(`{
			"label":     "BGP Session Flapping"
		}`),
	})
	require.NoError(t, err)

	probeBId, err := bpClient.InstantiateIbaPredefinedProbe(ctx, &apstra.IbaPredefinedProbeRequest{
		Name: "drain_node_traffic_anomaly",
		Data: []byte(`{
			"label":     "Drain Traffic Anomaly"
		}`),
	})
	require.NoError(t, err)

	ap := timeutils.NewDurationInSecs(1)
	widgetA := apstra.IbaWidget{
		Label:              "BGP Session Flapping",
		ProbeId:            probeAId,
		StageName:          "BGP Session",
		Type:               enum.IbaWidgetTypeStage,
		AggregationPeriod:  ap,
		TimeSeriesDuration: ap,
	}
	// widgetAId, err := bpClient.CreateIbaWidget(ctx, &widgetA)
	// require.NoError(t, err)

	widgetB := apstra.IbaWidget{
		Label:              "Drain Traffic Anomaly",
		ProbeId:            probeBId,
		StageName:          "excess_range",
		Type:               enum.IbaWidgetTypeStage,
		AggregationPeriod:  ap,
		TimeSeriesDuration: ap,
	}
	// widgetBId, err := bpClient.CreateIbaWidget(ctx, &widgetB)
	// require.NoError(t, err)
	widgetC := apstra.IbaWidget{
		Label:     "Drain Traffic Anomaly 2",
		ProbeId:   probeBId,
		StageName: "excess_range",
		Type:      enum.IbaWidgetTypeStage,
	}
	return widgetA, widgetB, widgetC
}
