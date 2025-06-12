// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDTwoStageL3ClosTags(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	extraTagCount := 3

	compare := func(t *testing.T, a, b *TwoStageL3ClosTagData) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Label, b.Label)
		require.Equal(t, a.Description, b.Description)
	}

	type testStep struct {
		tagData TwoStageL3ClosTagData
	}

	type testCase struct {
		steps []testStep
	}

	testCases := map[string]testCase{
		"start_minimal": {
			steps: []testStep{
				{
					tagData: TwoStageL3ClosTagData{
						Label: "start_minimal",
					},
				},
				{
					tagData: TwoStageL3ClosTagData{
						Label:       "start_minimal",
						Description: randString(6, "hex"),
					},
				},
				{
					tagData: TwoStageL3ClosTagData{
						Label: "start_minimal",
					},
				},
			},
		},
		"start_maximal": {
			steps: []testStep{
				{
					tagData: TwoStageL3ClosTagData{
						Label:       "start_maximal",
						Description: randString(6, "hex"),
					},
				},
				{
					tagData: TwoStageL3ClosTagData{
						Label: "start_maximal",
					},
				},
				{
					tagData: TwoStageL3ClosTagData{
						Label:       "start_maximal",
						Description: randString(6, "hex"),
					},
				},
			},
		},
	}

	for clientName, client := range clients {
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bp := testBlueprintA(ctx, t, client.client)

			for range extraTagCount {
				_, err = bp.CreateTag(ctx, TwoStageL3ClosTagData{Label: randString(6, "hex")})
				require.NoError(t, err)
			}

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					// t.Parallel(x) // do not parallelize - we count the total number of tags in here

					require.Greater(t, len(tCase.steps), 0) // we will blindly refer to steps[0] - it better exist

					id, err := bp.CreateTag(ctx, tCase.steps[0].tagData)
					require.NoError(t, err)

					tag, err := bp.GetTag(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, tag.Id)
					compare(t, &tCase.steps[0].tagData, tag.Data)

					tags, err := bp.GetAllTags(ctx)
					require.NoError(t, err)
					require.Equal(t, extraTagCount+1, len(tags))
					require.Contains(t, tags, id)
					compare(t, tags[id].Data, &tCase.steps[0].tagData)

					for _, step := range tCase.steps {
						err = bp.UpdateTag(ctx, id, step.tagData)
						require.NoError(t, err)

						tag, err := bp.GetTag(ctx, id)
						require.NoError(t, err)
						require.Equal(t, id, tag.Id)
						compare(t, &step.tagData, tag.Data)

						tags, err := bp.GetAllTags(ctx)
						require.NoError(t, err)
						require.Equal(t, extraTagCount+1, len(tags))
						require.Contains(t, tags, id)
						compare(t, tags[id].Data, &step.tagData)
					}

					err = bp.DeleteTag(ctx, id)
					require.NoError(t, err)

					tags, err = bp.GetAllTags(ctx)
					require.NoError(t, err)
					require.NotContains(t, tags, id)

					var ace ClientErr

					_, err = bp.GetTag(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ErrNotfound, ace.Type())

					err = bp.DeleteTag(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ErrNotfound, ace.Type())
				})
			}
		})
	}
}
