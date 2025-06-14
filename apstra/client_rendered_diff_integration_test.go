// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"bufio"
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/stretchr/testify/require"
)

func TestGetNodeRenderedDiff(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bp := testBlueprintI(ctx, t, client.client)

			leafIds, err := getSystemIdsByRole(ctx, bp, "leaf")
			require.NoError(t, err)

			leafWg := new(sync.WaitGroup)
			leafWg.Add(len(leafIds))
			for _, leafId := range leafIds {
				t.Run("leaf_"+leafId.String()+"_without_diff", func(t *testing.T) {
					t.Parallel()

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
				leafWg.Wait() // wait for leafs to be verified diff-free

				// create a security zone
				szLabel := randString(6, "hex")
				szId, err := bp.CreateSecurityZone(ctx, &SecurityZoneData{
					Label:   szLabel,
					SzType:  SecurityZoneTypeEVPN,
					VrfName: szLabel,
				})
				require.NoError(t, err)

				err = bp.SetResourceAllocation(ctx, &ResourceGroupAllocation{
					ResourceGroup: ResourceGroup{
						Type:           ResourceTypeIp4Pool,
						Name:           ResourceGroupNameLeafIp4,
						SecurityZoneId: &szId,
					},
					PoolIds: []ObjectId{"Private-10_0_0_0-8"},
				})
				require.NoError(t, err)

				// prep VN bindings
				vlanId := Vlan(rand.IntN(vlanMax-2) + 2) // 2-4094
				vnBindings := make([]VnBinding, len(leafIds))
				for i, leafId := range leafIds {
					vnBindings[i] = VnBinding{
						SystemId: leafId,
						VlanId:   &vlanId,
					}
				}

				// create a VN within the security zone
				rip := randomPrefix(t, "172.16.0.0/12", 24)
				vnId, err := bp.CreateVirtualNetwork(ctx, &VirtualNetworkData{
					VirtualGatewayIpv4Enabled: true,
					Ipv4Enabled:               true,
					Ipv4Subnet:                &rip,
					Label:                     randString(6, "hex"),
					SecurityZoneId:            szId,
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

					leafWg.Wait()

					// commit the blueprint so our new VN shows up in deployed config
					status, err := bp.Client().GetBlueprintStatus(ctx, bp.Id())
					require.NoError(t, err)
					_, err = bp.Client().DeployBlueprint(ctx, &BlueprintDeployRequest{
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
