// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestLookup(t *testing.T) {
	ctx := t.Context()

	clients := testclient.GetTestClients(t, ctx)

	bogusGroupID := "bogus-group-id"
	bogusSystemID := "bogus-system-id"

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()

			bp := dctestobj.BlueprintJ(t, ctx, client.Client)

			// Values based on dctestobj.BlueprintJ()
			expectedGroupCount := 3   // 1 leaf pair, 2 access pair
			expectedSwitchCount := 10 // 1 spine, 2 leaf, 7 access

			t.Run("lookup_systems_using_bogus_group_id", func(t *testing.T) {
				// Begin by clearing the cache.
				bp.DropSysRedundancyCache()

				systems, err := bp.GetSystemsByRedundancyGroup(ctx, bogusGroupID)
				require.Error(t, err, "expected error for bogus group ID, but got none")
				var ace apstra.ClientErr
				require.ErrorAs(t, err, &ace)
				require.Equal(t, apstra.ErrGroupNotFoundInCache, ace.Type())
				require.Equal(t, bogusGroupID, ace.Detail())
				require.Empty(t, systems[0], "expected first member of bogus redundant system pair to have empty ID")
				require.Empty(t, systems[1], "expected second member of bogus redundant system pair to have empty ID")
				require.Equal(t, expectedGroupCount, bp.CountGroups())   // expectedGroupCount groups in the per-bp cache
				require.Equal(t, expectedSwitchCount, bp.CountSystems()) // expectedSwitchCount systems in the per-bp cache
			})

			t.Run("lookup_group_using_bogus_system_id", func(t *testing.T) {
				// Begin by clearing the cache.
				bp.DropSysRedundancyCache()

				group, peer, err := bp.GetRedundancyGroupBySystem(ctx, bogusSystemID)
				require.Nilf(t, group, "expected nil group for bogus system id")
				require.Nilf(t, peer, "expected nil group for bogus system id")
				require.Error(t, err, "expected error for bogus system ID, but got none")
				var ace apstra.ClientErr
				require.ErrorAs(t, err, &ace)
				require.Equal(t, apstra.ErrSystemNotFoundInCache, ace.Type())
				require.Equal(t, bogusSystemID, ace.Detail())
				require.Equal(t, expectedGroupCount, bp.CountGroups())   // expectedGroupCount groups in the per-bp cache
				require.Equal(t, expectedSwitchCount, bp.CountSystems()) // expectedSwitchCount systems in the per-bp cache
			})

			// Function which returns system IDs of switch nodes.
			switchIDs := func(t testing.TB, ctx context.Context, bp *apstra.TwoStageL3ClosClient) []string {
				q := new(apstra.PathQuery).
					SetClient(bp.Client()).
					SetBlueprintId(bp.Id()).
					Node([]apstra.QEEAttribute{
						apstra.NodeTypeSystem.QEEAttribute(),
						{Key: "system_type", Value: apstra.QEStringVal(apstra.SystemTypeSwitch.String())},
						{Key: "name", Value: apstra.QEStringVal("n_sys")},
					})

				var target struct {
					Items []struct {
						System struct {
							ID string `json:"id"`
						} `json:"n_sys"`
					} `json:"items"`
				}

				require.NoError(t, q.Do(ctx, &target))
				result := make([]string, len(target.Items))
				for i, item := range target.Items {
					result[i] = item.System.ID
				}
				return result
			}

			// Function which returns redundancy group node IDs.
			groupIDs := func(t testing.TB, ctx context.Context, bp *apstra.TwoStageL3ClosClient) []string {
				q := new(apstra.PathQuery).
					SetClient(bp.Client()).
					SetBlueprintId(bp.Id()).
					Node([]apstra.QEEAttribute{
						apstra.NodeTypeRedundancyGroup.QEEAttribute(),
						{Key: "name", Value: apstra.QEStringVal("n_grp")},
					})

				var target struct {
					Items []struct {
						Group struct {
							ID string `json:"id"`
						} `json:"n_grp"`
					} `json:"items"`
				}

				require.NoError(t, q.Do(ctx, &target))
				result := make([]string, len(target.Items))
				for i, item := range target.Items {
					result[i] = item.Group.ID
				}
				return result
			}

			// Create maps to keep track of every system and group we encounter during the following tests.
			systemIDSet := make(map[string]struct{})
			groupIDSet := make(map[string]struct{})
			var ranAllGroups, ranAllSystems bool

			t.Run("lookup_all_groups", func(t *testing.T) {
				ranAllGroups = true
				// Begin by clearing the cache.
				bp.DropSysRedundancyCache()

				groupCount := 0
				for _, groupID := range groupIDs(t, ctx, bp) {
					groupCount++
					groupIDSet[groupID] = struct{}{}
					systemIDs, err := bp.GetSystemsByRedundancyGroup(ctx, groupID)
					require.NoError(t, err)
					require.NotEmpty(t, systemIDs[0])
					require.NotEmpty(t, systemIDs[1])
					nodeType, err := bp.GetBindingNodeType(ctx, systemIDs[0])
					require.NoError(t, err)
					require.Equal(t, apstra.NodeTypeSystem, nodeType)
					nodeType, err = bp.GetBindingNodeType(ctx, systemIDs[1])
					require.NoError(t, err)
					require.Equal(t, apstra.NodeTypeSystem, nodeType)
					nodeType, err = bp.GetBindingNodeType(ctx, groupID)
					require.NoError(t, err)
					require.Equal(t, apstra.NodeTypeRedundancyGroup, nodeType)
					systemIDSet[systemIDs[0]] = struct{}{}
					systemIDSet[systemIDs[1]] = struct{}{}
				}
				require.Equal(t, expectedGroupCount, groupCount)
			})

			t.Run("lookup_all_systems", func(t *testing.T) {
				ranAllSystems = true
				// Begin by clearing the cache.
				bp.DropSysRedundancyCache()

				redundantSystemCount := 0
				for _, systemID := range switchIDs(t, ctx, bp) {
					systemIDSet[systemID] = struct{}{}
					group, peer, err := bp.GetRedundancyGroupBySystem(ctx, systemID)
					require.NoError(t, err)
					if group == nil {
						require.Nil(t, peer)
					} else {
						require.NotNil(t, peer)
						require.NotEqual(t, systemID, *peer)
						groupIDSet[*group] = struct{}{}
						redundantSystemCount++
						nodeType, err := bp.GetBindingNodeType(ctx, *group)
						require.NoError(t, err)
						require.Equal(t, apstra.NodeTypeRedundancyGroup, nodeType)
						nodeType, err = bp.GetBindingNodeType(ctx, *peer)
						require.NoError(t, err)
						require.Equal(t, apstra.NodeTypeSystem, nodeType)
					}
				}
				require.Equal(t, expectedGroupCount*2, redundantSystemCount)
			})

			// Aggregate totals assume both discovery subtests were selected.
			if ranAllGroups && ranAllSystems {
				require.Equal(t, expectedGroupCount, len(groupIDSet))
				require.Equal(t, expectedSwitchCount, len(systemIDSet))
				require.Equal(t, expectedGroupCount, bp.CountGroups())
				require.Equal(t, expectedSwitchCount, bp.CountSystems())
			}

			t.Run("concurrent_access", func(t *testing.T) {
				// Begin by clearing the cache.
				bp.DropSysRedundancyCache()

				// Discover inputs independently of which other subtests ran.
				systemIDSlice := switchIDs(t, ctx, bp)
				groupIDSlice := groupIDs(t, ctx, bp)
				require.Len(t, systemIDSlice, expectedSwitchCount)
				require.Len(t, groupIDSlice, expectedGroupCount)

				var wg sync.WaitGroup
				numGoRoutines := 100
				wg.Add(numGoRoutines)

				for i := range numGoRoutines {
					go func() {
						defer wg.Done()

						switch {
						case i%7 == 0: // Lookup using bogus system ID every 7th request.
							group, peer, err := bp.GetRedundancyGroupBySystem(ctx, fmt.Sprintf(bogusSystemID+"_%03d", i))
							require.Nil(t, group)
							require.Nil(t, peer)
							require.Error(t, err)
							var ace apstra.ClientErr
							require.ErrorAs(t, err, &ace)
							require.Equal(t, apstra.ErrSystemNotFoundInCache, ace.Type())
							require.Equal(t, fmt.Sprintf(bogusSystemID+"_%03d", i), ace.Detail())
						case i%6 == 0: // Lookup using bogus group ID every 6th request.
							systems, err := bp.GetSystemsByRedundancyGroup(ctx, fmt.Sprintf(bogusGroupID+"_%03d", i))
							require.Error(t, err)
							var ace apstra.ClientErr
							require.ErrorAs(t, err, &ace)
							require.Equal(t, apstra.ErrGroupNotFoundInCache, ace.Type())
							require.Equal(t, fmt.Sprintf(bogusGroupID+"_%03d", i), ace.Detail())
							require.Empty(t, systems[0])
							require.Empty(t, systems[1])
						case i%2 == 0: // Valid lookup in both directions (by system and by group, if any) on even numbers not divisible by 6 or 7.
							testSys := systemIDSlice[i%len(systemIDSlice)]
							group, peer, err := bp.GetRedundancyGroupBySystem(ctx, testSys)
							require.NoError(t, err)
							if group == nil {
								require.Nil(t, peer)
							} else { // We got a group ID. Run it the other way.
								require.NotNil(t, peer)
								require.NotEqual(t, testSys, *peer)
								systems, err := bp.GetSystemsByRedundancyGroup(ctx, *group)
								require.NoError(t, err)
								require.Contains(t, systems, testSys)
								require.Contains(t, systems, *peer)
								nodeType, err := bp.GetBindingNodeType(ctx, *group)
								require.NoError(t, err)
								require.Equal(t, apstra.NodeTypeRedundancyGroup, nodeType)
								nodeType, err = bp.GetBindingNodeType(ctx, *peer)
								require.NoError(t, err)
								require.Equal(t, apstra.NodeTypeSystem, nodeType)
								group2, peer2, err := bp.GetRedundancyGroupBySystem(ctx, *peer)
								require.NoError(t, err)
								require.NotNil(t, group2)
								require.NotNil(t, peer2)
								require.Equal(t, *group, *group2)
								require.Equal(t, testSys, *peer2)
							}
						case i%2 == 1: // Valid system lookup on odd numbers not divisible by 6 or 7.
							testGrp := groupIDSlice[i%len(groupIDSlice)]
							systems, err := bp.GetSystemsByRedundancyGroup(ctx, testGrp)
							require.NoError(t, err)
							for _, system := range systems { // look up the group associated with each returned sys ID
								group, peer, err := bp.GetRedundancyGroupBySystem(ctx, system)
								require.NoError(t, err)
								require.NotNil(t, group)
								require.NotNil(t, peer)
								require.Equal(t, testGrp, *group)
								require.Contains(t, systems, *peer)
							}
						}
					}()
				}

				wg.Wait()
			})
		})
	}
}
