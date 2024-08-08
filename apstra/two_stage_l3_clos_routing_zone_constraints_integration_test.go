//go:build integration

package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoutingZoneConstraints(t *testing.T) {
	ctx := context.Background()

	compare := func(t testing.TB, expected, actual *RoutingZoneConstraintData) {
		require.NotNil(t, expected)
		require.NotNil(t, actual)

		require.Equal(t, expected.Label, actual.Label)
		require.Equal(t, expected.Mode, actual.Mode)
		if expected.MaxRoutingZones == nil {
			require.Nil(t, actual.MaxRoutingZones)
		} else {
			require.Equal(t, expected.MaxRoutingZones, actual.MaxRoutingZones)
		}
		if expected.RoutingZoneIds != nil {
			compareSlicesAsSets(t, expected.RoutingZoneIds, actual.RoutingZoneIds, "routing zone IDs don't match")
		} else {
			require.Equal(t, 0, len(actual.RoutingZoneIds), "expected empty routing zone id list")
		}
	}

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		steps []RoutingZoneConstraintData
	}

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(client.client.apiVersion.String()+"_"+client.clientType+"_"+clientName, func(t *testing.T) {
			t.Parallel()

			bpClient := testBlueprintA(ctx, t, client.client)

			testCases := map[string]testCase{
				"change_label_only": {
					steps: []RoutingZoneConstraintData{
						{
							Label: randString(8, "hex"),
							Mode:  RoutingZoneConstraintModeAllow,
						},
						{
							Label: randString(8, "hex"),
							Mode:  RoutingZoneConstraintModeDeny,
						},
						{
							Label: randString(8, "hex"),
							Mode:  RoutingZoneConstraintModeNone,
						},
						{
							Label: randString(8, "hex"),
							Mode:  RoutingZoneConstraintModeAllow,
						},
					},
				},
				"change_mode_only": {
					steps: []RoutingZoneConstraintData{
						{
							Label: "change_mode_only",
							Mode:  RoutingZoneConstraintModeAllow,
						},
						{
							Label: "change_mode_only",
							Mode:  RoutingZoneConstraintModeDeny,
						},
						{
							Label: "change_mode_only",
							Mode:  RoutingZoneConstraintModeNone,
						},
						{
							Label: "change_mode_only",
							Mode:  RoutingZoneConstraintModeAllow,
						},
					},
				},
				"change_max_only": {
					steps: []RoutingZoneConstraintData{
						{
							Label:           "change_max_only",
							Mode:            RoutingZoneConstraintModeAllow,
							MaxRoutingZones: nil,
						},
						{
							Label:           "change_max_only",
							Mode:            RoutingZoneConstraintModeAllow,
							MaxRoutingZones: toPtr(0),
						},
						{
							Label:           "change_max_only",
							Mode:            RoutingZoneConstraintModeAllow,
							MaxRoutingZones: nil,
						},
						{
							Label:           "change_max_only",
							Mode:            RoutingZoneConstraintModeAllow,
							MaxRoutingZones: toPtr(1),
						},
						{
							Label:           "change_max_only",
							Mode:            RoutingZoneConstraintModeAllow,
							MaxRoutingZones: toPtr(2),
						},
					},
				},
				"change_rz_ids_only": {
					steps: []RoutingZoneConstraintData{
						{
							Label:          "change_rz_ids_only",
							Mode:           RoutingZoneConstraintModeAllow,
							RoutingZoneIds: nil,
						},
						{
							Label:          "change_rz_ids_only",
							Mode:           RoutingZoneConstraintModeAllow,
							RoutingZoneIds: []ObjectId{testSecurityZone(t, ctx, bpClient)},
						},
						{
							Label:          "change_rz_ids_only",
							Mode:           RoutingZoneConstraintModeAllow,
							RoutingZoneIds: []ObjectId{},
						},
						{
							Label:          "change_rz_ids_only",
							Mode:           RoutingZoneConstraintModeAllow,
							RoutingZoneIds: []ObjectId{testSecurityZone(t, ctx, bpClient), testSecurityZone(t, ctx, bpClient)},
						},
						{
							Label:          "change_rz_ids_only",
							Mode:           RoutingZoneConstraintModeAllow,
							RoutingZoneIds: nil,
						},
					},
				},
			}

			for tName, tCase := range testCases {
				tName, tCase := tName, tCase
				t.Run(tName, func(t *testing.T) {
					t.Parallel()

					var id ObjectId
					var rzc *RoutingZoneConstraint
					for i, step := range tCase.steps {
						if i == 0 {
							id, err = bpClient.CreateRoutingZoneConstraint(ctx, &step)
							require.NoError(t, err)
						} else {
							err = bpClient.UpdateRoutingZoneConstraint(ctx, id, &step)
							require.NoError(t, err)
						}

						rzc, err = bpClient.GetRoutingZoneConstraint(ctx, id)
						require.NoError(t, err)
						require.Equal(t, id, rzc.Id)
						compare(t, &step, rzc.Data)

						rzc, err = bpClient.GetRoutingZoneConstraintByName(ctx, step.Label)
						require.NoError(t, err)
						require.Equal(t, id, rzc.Id)
						compare(t, &step, rzc.Data)
					}

					all, err := bpClient.GetAllRoutingZoneConstraints(ctx)
					require.NoError(t, err)
					allIds := make([]ObjectId, len(all))
					for i, rzc := range all {
						allIds[i] = rzc.Id
					}
					require.Contains(t, allIds, id)

					err = bpClient.DeleteRoutingZoneConstraint(ctx, id)
					require.NoError(t, err)

					var ace ClientErr

					_, err = bpClient.GetRoutingZoneConstraint(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ErrNotfound, ace.Type())

					all, err = bpClient.GetAllRoutingZoneConstraints(ctx)
					require.NoError(t, err)
					allIds = make([]ObjectId, len(all))
					for i, rzc := range all {
						allIds[i] = rzc.Id
					}
					require.NotContains(t, allIds, id)

					_, err = bpClient.GetRoutingZoneConstraintByName(ctx, rzc.Data.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ErrNotfound, ace.Type())

					err = bpClient.DeleteRoutingZoneConstraint(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ErrNotfound, ace.Type())
				})
			}
		})
	}
}
