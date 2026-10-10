package redundancycache_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	cache "github.com/Juniper/apstra-go-sdk/internal/datacenter_system_redundancy_cache"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestLookup(t *testing.T) {
	ctx := t.Context()

	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()

			bp := dctestobj.TestBlueprintJ(t, ctx, client.Client)
			bpID := string(bp.Id())

			// Values based on TestBlueprintJ()
			expectedGroupCount := 3   // 1 leaf pair, 2 access pair
			expectedSwitchCount := 10 // 1 spine, 2 leaf, 7 access

			t.Run("lookup_systems_using_bogus_group_id", func(t *testing.T) {
				// Begin by clearing the cache.
				cache.Drop(bpID)

				systems, err := cache.LookupSystems(ctx, bp, "bogus-group-id")
				require.Error(t, err, "expected error for bogus group ID, but got none")
				require.ErrorContains(t, err, cache.GroupNotFoundError)
				require.Empty(t, systems[0], "expected first member of bogus redundant system pair to have empty ID")
				require.Empty(t, systems[1], "expected second member of bogus redundant system pair to have empty ID")
				require.Equal(t, expectedGroupCount, cache.CountGroups(bpID))   // expectedGroupCount groups in the per-bp cache
				require.Equal(t, expectedSwitchCount, cache.CountSystems(bpID)) // expectedSwitchCount systems in the per-bp cache
			})

			t.Run("lookup_group_using_bogus_system_id", func(t *testing.T) {
				// Begin by clearing the cache.
				cache.Drop(bpID)

				group, peer, err := cache.LookupGroup(ctx, bp, "bogus-system-id")
				require.Nilf(t, group, "expected nil group for bogus system id")
				require.Nilf(t, peer, "expected nil group for bogus system id")
				require.Error(t, err, "expected error for bogus system ID, but got none")
				require.ErrorContains(t, err, cache.SystemNotFoundError)
				require.Equal(t, expectedGroupCount, cache.CountGroups(bpID))   // expectedGroupCount groups in the per-bp cache
				require.Equal(t, expectedSwitchCount, cache.CountSystems(bpID)) // expectedSwitchCount systems in the per-bp cache
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
				cache.Drop(bpID)

				groupCount := 0
				for _, groupID := range groupIDs(t, ctx, bp) {
					groupCount++
					groupIDSet[groupID] = struct{}{}
					systemIDs, err := cache.LookupSystems(ctx, bp, groupID)
					require.NoError(t, err)
					require.NotEmpty(t, systemIDs[0])
					require.NotEmpty(t, systemIDs[1])
					nodeType, err := cache.LookupNodeType(ctx, bp, systemIDs[0])
					require.NoError(t, err)
					require.Equal(t, apstra.NodeTypeSystem, nodeType)
					nodeType, err = cache.LookupNodeType(ctx, bp, systemIDs[1])
					require.NoError(t, err)
					require.Equal(t, apstra.NodeTypeSystem, nodeType)
					nodeType, err = cache.LookupNodeType(ctx, bp, groupID)
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
				cache.Drop(bpID)

				redundantSystemCount := 0
				for _, systemID := range switchIDs(t, ctx, bp) {
					systemIDSet[systemID] = struct{}{}
					group, peer, err := cache.LookupGroup(ctx, bp, systemID)
					require.NoError(t, err)
					if group == nil {
						require.Nil(t, peer)
					} else {
						require.NotNil(t, peer)
						require.NotEqual(t, systemID, *peer)
						groupIDSet[*group] = struct{}{}
						redundantSystemCount++
						nodeType, err := cache.LookupNodeType(ctx, bp, *group)
						require.NoError(t, err)
						require.Equal(t, apstra.NodeTypeRedundancyGroup, nodeType)
						nodeType, err = cache.LookupNodeType(ctx, bp, *peer)
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
				require.Equal(t, expectedGroupCount, cache.CountGroups(bpID))
				require.Equal(t, expectedSwitchCount, cache.CountSystems(bpID))
			}

			t.Run("concurrent_access", func(t *testing.T) {
				// Begin by clearing the cache.
				cache.Drop(bpID)

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
							group, peer, err := cache.LookupGroup(ctx, bp, fmt.Sprintf("bogus_system_%03d", i))
							require.Nil(t, group)
							require.Nil(t, peer)
							require.Error(t, err)
							require.ErrorContains(t, err, cache.SystemNotFoundError)
						case i%6 == 0: // Lookup using bogus group ID every 6th request.
							systems, err := cache.LookupSystems(ctx, bp, fmt.Sprintf("bogus_group_%03d", i))
							require.Error(t, err)
							require.ErrorContains(t, err, cache.GroupNotFoundError)
							require.Empty(t, systems[0])
							require.Empty(t, systems[1])
						case i%2 == 0: // Valid lookup in both directions (by system and by group, if any) on even numbers not divisible by 6 or 7.
							testSys := systemIDSlice[i%len(systemIDSlice)]
							group, peer, err := cache.LookupGroup(ctx, bp, testSys)
							require.NoError(t, err)
							if group == nil {
								require.Nil(t, peer)
							} else { // We got a group ID. Run it the other way.
								require.NotNil(t, peer)
								require.NotEqual(t, testSys, *peer)
								systems, err := cache.LookupSystems(ctx, bp, *group)
								require.NoError(t, err)
								require.Contains(t, systems, testSys)
								require.Contains(t, systems, *peer)
								nodeType, err := cache.LookupNodeType(ctx, bp, *group)
								require.NoError(t, err)
								require.Equal(t, apstra.NodeTypeRedundancyGroup, nodeType)
								nodeType, err = cache.LookupNodeType(ctx, bp, *peer)
								require.NoError(t, err)
								require.Equal(t, apstra.NodeTypeSystem, nodeType)
								group2, peer2, err := cache.LookupGroup(ctx, bp, *peer)
								require.NoError(t, err)
								require.NotNil(t, group2)
								require.NotNil(t, peer2)
								require.Equal(t, *group, *group2)
								require.Equal(t, testSys, *peer2)
							}
						case i%2 == 1: // Valid system lookup on odd numbers not divisible by 6 or 7.
							testGrp := groupIDSlice[i%len(groupIDSlice)]
							systems, err := cache.LookupSystems(ctx, bp, testGrp)
							require.NoError(t, err)
							for _, system := range systems { // look up the group associated with each returned sys ID
								group, peer, err := cache.LookupGroup(ctx, bp, system)
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
