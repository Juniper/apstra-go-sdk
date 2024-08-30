//go:build integration

package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ConfigTemplate_Assignment(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	intSysCount := 4
	ctCount := 4

	type testStep struct {
		ctId   ObjectId
		sysIds []ObjectId
	}
	type testCase struct {
		steps []testStep
	}

	for clientName, client := range clients {
		t.Run(clientName, func(t *testing.T) {
			t.Parallel()

			ffc, intSysIds, _ := testFFBlueprintB(ctx, t, client.client, intSysCount, 0)

			ctIds := make([]ObjectId, ctCount)
			for i := range ctIds {
				id, err := ffc.CreateConfigTemplate(ctx, &ConfigTemplateData{
					Label: randString(6, "hex") + ".jinja",
					Text:  randString(30, "hex"),
				})
				require.NoError(t, err)
				ctIds[i] = id
			}

			testCases := map[string]testCase{
				"single_sysID": {
					steps: []testStep{
						{
							ctId:   ctIds[0],
							sysIds: []ObjectId{intSysIds[0]},
						},
						{
							ctId:   ctIds[0],
							sysIds: []ObjectId{intSysIds[1]},
						},
						{
							ctId:   ctIds[0],
							sysIds: []ObjectId{intSysIds[2]},
						},
						{
							ctId:   ctIds[0],
							sysIds: []ObjectId{intSysIds[0]},
						},
					},
				},
				"set_clear_set_single_sysID": {
					steps: []testStep{
						{
							ctId:   ctIds[1],
							sysIds: []ObjectId{intSysIds[0]},
						},
						{
							ctId:   ctIds[1],
							sysIds: []ObjectId{},
						},
						{
							ctId:   ctIds[1],
							sysIds: []ObjectId{intSysIds[0]},
						},
					},
				},
				"clear_set_clear_single_sysID": {
					steps: []testStep{
						{
							ctId:   ctIds[2],
							sysIds: []ObjectId{},
						},
						{
							ctId:   ctIds[2],
							sysIds: []ObjectId{intSysIds[0]},
						},
						{
							ctId:   ctIds[2],
							sysIds: []ObjectId{},
						},
					},
				},
				"exercise_multiple_systems": {
					steps: []testStep{
						{
							ctId:   ctIds[3],
							sysIds: []ObjectId{},
						},
						{
							ctId:   ctIds[3],
							sysIds: []ObjectId{intSysIds[0], intSysIds[1]},
						},
						{
							ctId:   ctIds[2],
							sysIds: []ObjectId{intSysIds[1], intSysIds[2]},
						},
						{
							ctId:   ctIds[2],
							sysIds: []ObjectId{intSysIds[0], intSysIds[3]},
						},
					},
				},
			}

			for tName, tCase := range testCases {
				tName, tCase := tName, tCase
				t.Run(tName, func(t *testing.T) {
					for _, step := range tCase.steps {
						require.NoError(t, ffc.UpdateConfigTemplateAssignmentsByTemplate(ctx, step.ctId, step.sysIds))
						assigned, err := ffc.GetConfigTemplateAssignments(ctx, step.ctId)
						require.NoError(t, err)
						compareSlicesAsSets(t, step.sysIds, assigned, "assignment mismatch after update")
					}
				})
			}
		})
	}
}
