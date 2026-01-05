// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"bufio"
	"math/rand/v2"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/query"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetNodeRenderedDiff(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bp := dctestobj.TestBlueprintI(t, ctx, client.Client)

			leafIds, err := query.SystemIdsByRole(ctx, bp, "leaf")
			require.NoError(t, err)

			leafWg := new(sync.WaitGroup)
			leafWg.Add(len(leafIds))
			for _, leafId := range leafIds {
				t.Run("leaf_"+leafId.String()+"_without_diff", func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					// staging config should have no diffs at this point
					diff, err := bp.Client().GetNodeRenderedConfigDiff(ctx, bp.Id(), leafId)
					require.NoError(t, err)
					require.NotNil(t, diff)
					require.Equal(t, "null", string(diff.PristineConfig)) // no pristine config, i guess
					require.Empty(t, diff.Config)                         // no diff
					require.False(t, diff.SupportsDiffConfig)             // whatever this is
					require.Greater(t, len(diff.Context), 4000)           // 4KB-ish of context?

					leafWg.Done()
				})
			}

			// make changes to the rendered config by deploying a virtual network to the switches
			t.Run("leafs_with_diffs", func(t *testing.T) {
				t.Parallel()
				ctx := testutils.ContextWithTestID(ctx, t)

				leafWg.Wait() // wait for leafs to be verified diff-free

				// create a security zone
				szLabel := testutils.RandString(6, "hex")
				szId, err := bp.CreateSecurityZone(ctx, apstra.SecurityZone{
					Label:   szLabel,
					Type:    enum.SecurityZoneTypeEVPN,
					VRFName: szLabel,
				})
				require.NoError(t, err)

				err = bp.SetResourceAllocation(ctx, &apstra.ResourceGroupAllocation{
					ResourceGroup: apstra.ResourceGroup{
						Type:           apstra.ResourceTypeIp4Pool,
						Name:           apstra.ResourceGroupNameLeafIp4,
						SecurityZoneId: (*apstra.ObjectId)(&szId),
					},
					PoolIds: []apstra.ObjectId{"Private-10_0_0_0-8"},
				})
				require.NoError(t, err)

				// prep VN bindings
				vlanId := apstra.VLAN(rand.IntN(apstra.VlanMax-2) + 2) // 2-4094
				vnBindings := make([]apstra.VnBinding, len(leafIds))
				for i, leafId := range leafIds {
					vnBindings[i] = apstra.VnBinding{
						SystemId: leafId,
						VlanId:   &vlanId,
					}
				}

				// create a VN within the security zone
				rip := testutils.RandomPrefix(t, "172.16.0.0/12", 24)
				vnId, err := bp.CreateVirtualNetwork(ctx, &apstra.VirtualNetworkData{
					VirtualGatewayIpv4Enabled: true,
					Ipv4Enabled:               true,
					Ipv4Subnet:                &rip,
					Label:                     testutils.RandString(6, "hex"),
					SecurityZoneId:            apstra.ObjectId(szId),
					VnBindings:                vnBindings,
					VnType:                    enum.VnTypeVxlan,
				})
				require.NoError(t, err)
				t.Log(vnId.String())

				time.Sleep(time.Second) // ensure time for config diffs to render

				leafWg.Add(len(leafIds))
				for _, leafId := range leafIds {
					t.Run("leaf_"+leafId.String(), func(t *testing.T) {
						t.Parallel()
						ctx := testutils.ContextWithTestID(ctx, t)

						// staging config should have diffs at this point
						diff, err := bp.Client().GetNodeRenderedConfigDiff(ctx, bp.Id(), leafId)
						require.NoError(t, err)
						require.NotNil(t, diff)
						require.Equal(t, "null", string(diff.PristineConfig)) // no pristine config, i guess
						require.False(t, diff.SupportsDiffConfig)             // whatever this is
						require.Greater(t, len(diff.Context), 4000)           // 4KB-ish of context?
						adds, dels := 0, 0
						scanner := bufio.NewScanner(strings.NewReader(diff.Config))
						for scanner.Scan() {
							switch {
							case strings.HasPrefix(scanner.Text(), "+"):
								adds++
							case strings.HasPrefix(scanner.Text(), "-"):
								dels++
							}
						}
						require.Greater(t, adds, 40)
						require.Equal(t, dels, 0)

						leafWg.Done() // there is a deadlock here if require above fails
					})
				}

				t.Run("test_config_withdrawal_diff", func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					leafWg.Wait()

					// commit the blueprint so our new VN shows up in deployed config
					status, err := bp.Client().GetBlueprintStatus(ctx, bp.Id())
					require.NoError(t, err)
					_, err = bp.Client().DeployBlueprint(ctx, &apstra.BlueprintDeployRequest{
						Id:          bp.Id(),
						Description: `commit so that we can generate "delete" diffs`,
						Version:     status.Version,
					})
					require.NoError(t, err)

					// delete the VN to generate config withdrawals
					err = bp.DeleteVirtualNetwork(ctx, vnId)
					require.NoError(t, err)

					time.Sleep(time.Second) // ensure time for config diffs to render

					for _, leafId := range leafIds {
						// staging config should have diffs at this point
						diff, err := bp.Client().GetNodeRenderedConfigDiff(ctx, bp.Id(), leafId)
						require.NoError(t, err)
						require.NotNil(t, diff)
						require.Equal(t, "null", string(diff.PristineConfig)) // no pristine config, i guess
						require.False(t, diff.SupportsDiffConfig)             // whatever this is
						require.Greater(t, len(diff.Context), 4000)           // 4KB-ish of context?
						adds, dels := 0, 0
						scanner := bufio.NewScanner(strings.NewReader(diff.Config))
						for scanner.Scan() {
							switch {
							case strings.HasPrefix(scanner.Text(), "+"):
								adds++
							case strings.HasPrefix(scanner.Text(), "-"):
								dels++
							}
						}
						require.Equal(t, adds, 0)
						require.Greater(t, dels, 40)
					}
				})
			})
		})
	}
}
