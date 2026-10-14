// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	comparedatacenter "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/datacenter"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestSetVNToLeafBinding(t *testing.T) {
	ctx := t.Context()

	type testCase struct {
		constraints []compatibility.Constraint
		steps       []apstra.SetVNToLeafBindingRequest
	}

	testCases := map[string]testCase{
		"simple": {
			steps: []apstra.SetVNToLeafBindingRequest{
				{
					Binding: datacenter.VNBinding{
						SystemID:            "rack1_leaf1",
						VLAN:                pointer.To(uint16(rand.IntN(4000) + 95)), // VLANs 95-4094
						AccessSwitchNodeIDs: []string{"rack1_access1", "rack1_access2"},
					},
				},
				{
					Binding: datacenter.VNBinding{
						SystemID:            "rack1_leaf2",
						VLAN:                pointer.To(uint16(rand.IntN(4000) + 95)), // VLANs 95-4094
						AccessSwitchNodeIDs: []string{"rack1_access3"},
					},
				},
				{
					Binding: datacenter.VNBinding{
						SystemID:            "rack1_leaf1",
						VLAN:                pointer.To(uint16(rand.IntN(4000) + 95)), // VLANs 95-4094
						AccessSwitchNodeIDs: []string{"rack1_access4", "rack1_access5", "rack1_access6", "rack1_access7"},
					},
				},
			},
		},
	}

	findAndConvertBinding := func(t testing.TB, ctx context.Context, leafID string, bindings []datacenter.VNBinding, bp *apstra.TwoStageL3ClosClient) datacenter.VNBinding {
		t.Helper()

		groupID, _, err := bp.GetRedundancyGroupBySystem(ctx, leafID)
		require.NoError(t, err)

		for _, binding := range bindings {
			if leafID == binding.SystemID || (groupID != nil && *groupID == binding.SystemID) {
				binding.SystemID = leafID                                    // Swap the actual leaf ID in place of the bound ID.
				for i := len(binding.AccessSwitchNodeIDs) - 1; i >= 0; i-- { // Loop backward to safely modify the slice.
					accessIDs, err := bp.GetSystemsByRedundancyGroup(ctx, binding.AccessSwitchNodeIDs[i])
					if err != nil { // A "not found" error is expected if the access switch is not a member of a redundancy group.
						var ace apstra.ClientErr
						require.ErrorAs(t, err, &ace)
						require.Equal(t, apstra.ErrGroupNotFoundInCache, ace.Type())
						continue
					}

					binding.AccessSwitchNodeIDs = append(binding.AccessSwitchNodeIDs[:i], binding.AccessSwitchNodeIDs[i+1:]...) // Remove the group ID.
					binding.AccessSwitchNodeIDs = append(binding.AccessSwitchNodeIDs, accessIDs[:]...)                          // Add the individual system IDs.
				}
				return binding
			}
		}

		t.Fatalf("couldn't find binding with system ID %s", leafID)
		return datacenter.VNBinding{}
	}

	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bp := dctestobj.BlueprintK(t, ctx, client.Client)
			rzID := dctestobj.SecurityZoneA(t, ctx, bp)
			vnID := dctestobj.VirtualNetworkB(t, ctx, bp, rzID, "192.0.2.0/24", "3fff::/20")

			query := new(apstra.PathQuery).
				SetBlueprintId(bp.Id()).
				SetClient(bp.Client()).
				Node([]apstra.QEEAttribute{
					apstra.NodeTypeSystem.QEEAttribute(),
					{Key: "name", Value: apstra.QEStringVal("switch")},
					{Key: "system_type", Value: apstra.QEStringVal("switch")},
					{Key: "role", Value: apstra.QEStringValIsIn{"access", "leaf"}},
				})
			var target struct {
				Items []struct {
					Switch struct {
						ID    string `json:"id"`
						Label string `json:"label"`
					} `json:"switch"`
				} `json:"items"`
			}
			require.NoError(t, query.Do(ctx, &target))

			// Build a map of simple name representing rack and switch (rack1_leaf1, rack2_access1, etc.)
			// to system ID for all access and leaf switches in the blueprint.
			sw := make(map[string]string)
			for _, item := range target.Items {
				parts := strings.Split(item.Switch.Label, "_")
				require.Len(t, parts, 3)
				rackNum, err := strconv.Atoi(parts[1])
				require.NoError(t, err)
				sw[fmt.Sprintf("rack%d_%s", rackNum, parts[2])] = item.Switch.ID
			}
			require.Len(t, sw, 18)                   // Blueprint K has 4 leaf switches and 14 access switches
			require.Contains(t, sw, "rack1_leaf1")   // pair member with rack1_leaf2
			require.Contains(t, sw, "rack1_leaf2")   // pair member with rack1_leaf1
			require.Contains(t, sw, "rack1_access1") // single-homed to rack1_leaf1
			require.Contains(t, sw, "rack1_access2") // single-homed to rack1_leaf2
			require.Contains(t, sw, "rack1_access3") // dual-homed to rack1_leaf1/2
			require.Contains(t, sw, "rack1_access4") // pair member with access5, single-homed to rack1_leaf1
			require.Contains(t, sw, "rack1_access5") // pair member with access4, single-homed to rack1_leaf2
			require.Contains(t, sw, "rack1_access6") // pair member with access7, dual-homed to rack1_leaf1/2
			require.Contains(t, sw, "rack1_access7") // pair member with access6, dual-homed to rack1_leaf1/2
			require.Contains(t, sw, "rack2_leaf1")   // pair member with rack2_leaf2
			require.Contains(t, sw, "rack2_leaf2")   // pair member with rack2_leaf1
			require.Contains(t, sw, "rack2_access1") // single-homed to rack2_leaf1
			require.Contains(t, sw, "rack2_access2") // single-homed to rack2_leaf2
			require.Contains(t, sw, "rack2_access3") // dual-homed to rack2_leaf1/2
			require.Contains(t, sw, "rack2_access4") // pair member with access5, single-homed to rack2_leaf1
			require.Contains(t, sw, "rack2_access5") // pair member with access4, single-homed to rack2_leaf2
			require.Contains(t, sw, "rack2_access6") // pair member with access7, dual-homed to rack1_leaf1/2
			require.Contains(t, sw, "rack2_access7") // pair member with access6, dual-homed to rack1_leaf1/2

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					for _, constraint := range tCase.constraints {
						if !constraint.Check(client.APIVersion()) {
							t.Skipf("skipping test case %q for API version %s: %s", tName, client.APIVersion(), constraint)
						}
					}

					for _, step := range tCase.steps {
						// Swap actual IDs in place of the simple names used in the test case definition.
						step.VNID = vnID
						step.Binding.SystemID = sw[step.Binding.SystemID]
						for i, v := range step.Binding.AccessSwitchNodeIDs {
							step.Binding.AccessSwitchNodeIDs[i] = sw[v]
						}

						err := bp.SetVNToLeafBinding(ctx, step)
						require.NoError(t, err)

						vn, err := bp.GetVirtualNetwork(ctx, vnID)
						require.NoError(t, err)
						binding := findAndConvertBinding(t, ctx, step.Binding.SystemID, vn.Bindings, bp)
						comparedatacenter.VNBinding(t, step.Binding, binding)

						//var foundBinding *datacenter.VNBinding
						//for _, binding := range vn.Bindings {
						//	bp.GetSystemsByRedundancyGroup(ctx, binding.SystemID)
						//}
					}

				})
			}

			//req := apstra.SetVNToLeafBindingRequest{
			//	VNID: vnID,
			//	Binding: datacenter.VNBinding{
			//		AccessSwitchNodeIDs: access,
			//		SystemID:            leafID,
			//		VLAN:                pointer.To(vlan),
			//	},
			//	//SVIIP: &datacenter.SVIAddressing{
			//	//	SystemID: leafID,
			//	//	IPv4Addr: parse(t, ip),
			//	//IPv4Mode: enum.IPv4SVIModeEnabled,
			//	//IPv6Addr: nil,
			//	//IPv6Mode: enum.IPv6SVIModeDisabled,
			//	//},
			//	Unbind: true,
			//}
			//
			//err := bp.SetVNToLeafBinding(ctx, req)
			//require.NoError(t, err)
		})
	}
}
