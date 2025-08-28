// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"github.com/Juniper/apstra-go-sdk/apstra"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestCreateReadUpdateDeleteIbaDashboards(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {

		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			if !compatibility.IbaDashboardSupported.Check(client.APIVersion()) {
				t.Skipf("skipping test due to unsupported API changes in %s", client.APIVersion())
			}

			bpClient := dctestobj.TestBlueprintA(t, ctx, client.Client)
			widgetA, widgetB, widgetC := dctestobj.TestWidgetsABC(t, ctx, bpClient)

			predefinedDashboardIds, err := bpClient.ListAllIbaPredefinedDashboardIds(ctx)
			require.NoError(t, err)
			t.Logf("found %d predefined dashboards", len(predefinedDashboardIds))

			for _, predefinedDashboardId := range predefinedDashboardIds {
				t.Run(predefinedDashboardId.String(), func(t *testing.T) {
					t.Parallel()

					// Some predefined dashboards cannot be tested
					switch predefinedDashboardId {
					case "evpn_vxlan_route_summary": // This one requires the blueprint to be deployed
						t.Skip("skipping because it requires blueprint to be deployed")
					case "stripe_traffic_summary": // This one is an autoscaling dashboard, so it cannot be instantiated
						t.Skip("skipping because it is an autoscaling dashboard")
					}

					t.Logf("instantiating dashboard id: %s", predefinedDashboardId)
					id, err := bpClient.InstantiateIbaPredefinedDashboard(ctx, predefinedDashboardId, testutils.RandString(6, "hex"))
					require.NoError(t, err)

					t.Logf("Name :%s Created Id :%s", predefinedDashboardId, id)
					t.Log("Getting Dashboard")
					d1, err := bpClient.GetIbaDashboard(ctx, id)
					require.NoError(t, err)

					d1.Data.Label = testutils.RandString(5, "hex")
					t.Log("Updating Dashboard")
					d1.Data.UpdatedBy = ""
					d1.Data.PredefinedDashboard = ""

					err = bpClient.UpdateIbaDashboard(ctx, id, d1.Data)
					require.NoError(t, err)

					d2, err := bpClient.GetIbaDashboard(ctx, id)
					require.NoError(t, err)
					require.Equalf(t, d1.Data.Label, d2.Data.Label, "Update Seems to have failed. Label should have been %s is %s", d1.Data.Label, d2.Data.Label)

					t.Log("Deleting Dashboard")
					err = bpClient.DeleteIbaDashboard(ctx, id)
					require.NoError(t, err)
				})
			}

			ds, err := bpClient.GetAllIbaDashboards(ctx)
			require.NoError(t, err)
			require.Equalf(t, 0, len(ds), "Expected no dashboards, got %d.", len(ds))

			req1 := apstra.IbaDashboardData{
				Description:   "Test Dashboard",
				Default:       false,
				Label:         "Test Dash",
				IbaWidgetGrid: [][]apstra.IbaWidget{{widgetA, widgetB}, {widgetC}},
			}
			id, err := bpClient.CreateIbaDashboard(ctx, &req1)
			require.NoError(t, err)

			widgetA.Label = "label2A"
			widgetB.Label = "label2B"

			req2 := apstra.IbaDashboardData{
				Description:   "Test Dashboard Backup",
				Default:       false,
				Label:         "Test Dash B",
				IbaWidgetGrid: [][]apstra.IbaWidget{{widgetA, widgetB}},
			}
			_, err = bpClient.CreateIbaDashboard(ctx, &req2)
			require.NoError(t, err)

			ds, err = bpClient.GetAllIbaDashboards(ctx)
			require.NoError(t, err)

			require.Equalf(t, 2, len(ds), "expected %d dashboards, got %d", 2, len(ds))

			checkDashes := func() {
				d1, err := bpClient.GetIbaDashboard(ctx, id)
				require.NoError(t, err)
				require.NotNil(t, d1.Data)

				d2, err := bpClient.GetIbaDashboardByLabel(ctx, d1.Data.Label)
				require.NoError(t, err)
				require.NotNil(t, d2.Data)

				priorValue := req1.UpdatedBy
				req1.UpdatedBy = d1.Data.UpdatedBy // this wasn't part of the request

				compare.Dashboards(t, req1, *d1.Data)
				req1.UpdatedBy = priorValue // restore prior value

				compare.Dashboards(t, *d1.Data, *d2.Data)
			}
			checkDashes()

			req1.Label = "Test Dash 2"
			req1.IbaWidgetGrid = append(req1.IbaWidgetGrid, []apstra.IbaWidget{widgetA, widgetB})
			req1.Description = "Test Dashboard 2"
			err = bpClient.UpdateIbaDashboard(ctx, id, &req1)
			require.NoError(t, err)
			checkDashes()

			err = bpClient.DeleteIbaDashboard(ctx, id)
			require.NoError(t, err)

			var ace apstra.ClientErr

			// attempt to fetch the deleted dashboard
			_, err = bpClient.GetIbaDashboard(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.Type(), apstra.ErrNotfound)

			// attempt to fetch the deleted dashboard by name
			_, err = bpClient.GetIbaDashboardByLabel(ctx, req1.Label)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.Type(), apstra.ErrNotfound)

			// attempt to delete the deleted dashboard
			err = bpClient.DeleteIbaDashboard(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.Type(), apstra.ErrNotfound)

			// ensure the deleted dashboard isn't among "all"
			ds, err = bpClient.GetAllIbaDashboards(ctx)
			require.NoError(t, err)
			ids := make([]apstra.ObjectId, len(ds))
			for i, d := range ds {
				ids[i] = d.Id
			}
			require.NotContains(t, ids, id)
		})
	}
}
