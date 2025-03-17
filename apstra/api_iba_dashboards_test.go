// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/stretchr/testify/require"
)

func TestCreateReadUpdateDeleteIbaDashboards(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client

		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			if !compatibility.IbaDashboardSupported.Check(client.client.apiVersion) {
				t.Skipf("skipping test due to unsupported API changes in %s", client.client.apiVersion)
			}

			bpClient := testBlueprintA(ctx, t, client.client)
			widgetA, widgetB, widgetC := testWidgetsABC(ctx, t, bpClient)

			pds, err := bpClient.ListAllIbaPredefinedDashboardIds(ctx)
			require.NoError(t, err)
			t.Log(len(pds))

			for _, d := range pds {
				// Some predefined dashboards cannot be tested
				switch d {
				case "evpn_vxlan_route_summary": // This one requires the blueprint to be deployed
					continue
				}

				id, err := bpClient.InstantiateIbaPredefinedDashboard(ctx, d, d.String())
				require.NoError(t, err)

				t.Logf("Name :%s Created Id :%s", d, id)
				t.Log("Getting Dashboard")
				d1, err := bpClient.GetIbaDashboard(ctx, id)
				require.NoError(t, err)

				d1.Data.Label = randString(5, "hex")
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
			}

			ds, err := bpClient.GetAllIbaDashboards(ctx)
			require.NoError(t, err)
			require.Equalf(t, 0, len(ds), "Expected no dashboards, got %d.", len(ds))

			req1 := IbaDashboardData{
				Description:   "Test Dashboard",
				Default:       false,
				Label:         "Test Dash",
				IbaWidgetGrid: [][]IbaWidget{{widgetA, widgetB}, {widgetC}},
			}
			id, err := bpClient.CreateIbaDashboard(ctx, &req1)
			require.NoError(t, err)

			widgetA.Label = "label2A"
			widgetB.Label = "label2B"

			req2 := IbaDashboardData{
				Description:   "Test Dashboard Backup",
				Default:       false,
				Label:         "Test Dash B",
				IbaWidgetGrid: [][]IbaWidget{{widgetA, widgetB}},
			}
			_, err = bpClient.CreateIbaDashboard(ctx, &req2)
			require.NoError(t, err)

			ds, err = bpClient.GetAllIbaDashboards(ctx)
			require.NoError(t, err)

			require.Equalf(t, 2, len(ds), "expected %d dashboards, got %d", 2, len(ds))

			checkDashes := func() {
				d1, err := bpClient.GetIbaDashboard(ctx, id)
				require.NoError(t, err)

				d2, err := bpClient.GetIbaDashboardByLabel(ctx, d1.Data.Label)
				require.NoError(t, err)

				priorValue := req1.UpdatedBy
				req1.UpdatedBy = d1.Data.UpdatedBy // this wasn't part of the request
				if !compareDashboards(req1, *d1.Data) {
					t.Log(req1)
					t.Log(d1.Data)
					t.Fatal("Dashboard request doesn't match GetIbaDashboard.Data")
				}
				req1.UpdatedBy = priorValue // restore prior value

				if !compareDashboards(*d1.Data, *d2.Data) {
					t.Fatal("GetIbaDashboardByLabel gets different object than GetIbaDashboard")
				}
			}
			checkDashes()

			req1.Label = "Test Dash 2"
			req1.IbaWidgetGrid = append(req1.IbaWidgetGrid, []IbaWidget{widgetA, widgetB})
			req1.Description = "Test Dashboard 2"
			err = bpClient.UpdateIbaDashboard(ctx, id, &req1)
			require.NoError(t, err)
			checkDashes()

			err = bpClient.DeleteIbaDashboard(ctx, id)
			require.NoError(t, err)

			var ace ClientErr

			// attempt to fetch the deleted dashboard
			_, err = bpClient.GetIbaDashboard(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.Type(), ErrNotfound)

			// attempt to fetch the deleted dashboard by name
			_, err = bpClient.GetIbaDashboardByLabel(ctx, req1.Label)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.Type(), ErrNotfound)

			// attempt to delete the deleted dashboard
			err = bpClient.DeleteIbaDashboard(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.Type(), ErrNotfound)

			// ensure the deleted dashboard isn't among "all"
			ds, err = bpClient.GetAllIbaDashboards(ctx)
			require.NoError(t, err)
			ids := make([]ObjectId, len(ds))
			for i, d := range ds {
				ids[i] = d.Id
			}
			require.NotContains(t, ids, id)
		})
	}
}
