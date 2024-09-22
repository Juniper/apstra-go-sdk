// Copyright (c) Juniper Networks, Inc., 2024-2024.
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

func TestFfGroupGeneratorsCrud(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	resourceGroupCount := 2

	compare := func(t testing.TB, a, b *FreeformGroupGeneratorData) {
		require.NotNil(t, a)
		require.NotNil(t, b)

		if a.ParentId != nil || b.ParentId != nil {
			require.NotNil(t, a.ParentId)
			require.NotNil(t, b.ParentId)
			require.Equal(t, *a.ParentId, *b.ParentId)
		}

		require.Equal(t, a.Label, b.Label)
		require.Equal(t, a.Scope, b.Scope)
	}

	type testStep struct {
		config FreeformGroupGeneratorData
	}

	type testCase struct {
		steps []testStep
	}

	for clientName, client := range clients {
		t.Run(clientName, func(t *testing.T) {
			t.Parallel()

			bp := testFFBlueprintA(ctx, t, client.client)

			resourceGroupIds := make([]ObjectId, resourceGroupCount)
			for i := range resourceGroupCount {
				resourceGroupIds[i], err = bp.CreateRaGroup(ctx, &FreeformRaGroupData{
					Label: randString(6, "hex"),
				})
				require.NoError(t, err)
			}

			testCases := map[string]testCase{
				"root": {
					steps: []testStep{
						{
							config: FreeformGroupGeneratorData{
								ParentId: nil,
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
						{
							config: FreeformGroupGeneratorData{
								ParentId: nil,
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
					},
				},
				"group_member": {
					steps: []testStep{
						{
							config: FreeformGroupGeneratorData{
								ParentId: &resourceGroupIds[0],
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
						{
							config: FreeformGroupGeneratorData{
								ParentId: &resourceGroupIds[0],
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
					},
				},
				"change_group_membership": {
					steps: []testStep{
						{
							config: FreeformGroupGeneratorData{
								ParentId: &resourceGroupIds[0],
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
						{
							config: FreeformGroupGeneratorData{
								ParentId: &resourceGroupIds[1],
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
					},
				},
				"change_between_root_and_group_member": {
					steps: []testStep{
						{
							config: FreeformGroupGeneratorData{
								ParentId: nil,
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
						{
							config: FreeformGroupGeneratorData{
								ParentId: &resourceGroupIds[0],
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
						{
							config: FreeformGroupGeneratorData{
								ParentId: nil,
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
						{
							config: FreeformGroupGeneratorData{
								ParentId: &resourceGroupIds[1],
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
					},
				},
				"change_between_group_member_and_root": {
					steps: []testStep{
						{
							config: FreeformGroupGeneratorData{
								ParentId: &resourceGroupIds[0],
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
						{
							config: FreeformGroupGeneratorData{
								ParentId: nil,
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
						{
							config: FreeformGroupGeneratorData{
								ParentId: &resourceGroupIds[1],
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
						{
							config: FreeformGroupGeneratorData{
								ParentId: nil,
								Label:    randString(6, "hex"),
								Scope:    fmt.Sprintf("node(label='%s', name='target')", randString(6, "hex")),
							},
						},
					},
				},
			}

			for tName, tCase := range testCases {
				tName, tCase := tName, tCase

				t.Run(tName, func(t *testing.T) {
					var err error
					var id ObjectId

					for i, step := range tCase.steps {

						// create or update the generator
						if i == 0 {
							id, err = bp.CreateGroupGenerator(ctx, &step.config)
							require.NoError(t, err)
						} else {
							require.NoError(t, bp.UpdateGroupGenerator(ctx, id, &step.config))
						}

						// fetch the generator by ID
						generator, err := bp.GetGroupGenerator(ctx, id)
						require.NoError(t, err)
						require.Equal(t, id, generator.Id)
						compare(t, &step.config, generator.Data)

						// fetch the generator by name
						generator, err = bp.GetGroupGeneratorByName(ctx, step.config.Label)
						require.NoError(t, err)
						require.Equal(t, id, generator.Id)
						compare(t, &step.config, generator.Data)

						// fetch all generators
						generators, err := bp.GetAllGroupGenerators(ctx)
						ids := make([]ObjectId, len(generators))
						for i, g := range generators {
							ids[i] = g.Id
							if id == g.Id {
								compare(t, &step.config, g.Data)
							}
						}
						require.Contains(t, ids, id)
					}

					// delete the generator
					require.NoError(t, bp.DeleteGroupGenerator(ctx, id))

					var ace ClientErr

					// fail to fetch the generator after deletion
					_, err = bp.GetGroupGenerator(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), ErrNotfound)

					// fail to delete the generator after deletion
					err = bp.DeleteGroupGenerator(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), ErrNotfound)

					// fail to find the generator after deletion
					generators, err := bp.GetAllGroupGenerators(ctx)
					ids := make([]ObjectId, len(generators))
					for i, g := range generators {
						ids[i] = g.Id
					}
					require.NotContains(t, ids, id)
				})
			}
		})
	}
}
