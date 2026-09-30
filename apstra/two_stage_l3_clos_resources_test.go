// Copyright (c) Juniper Networks, Inc., 2022-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"sync"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	resourcetestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/resource_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestSetGetResourceAllocation(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	ipv4PoolCount := 4
	asnPoolCount := 4

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			ctx := testutils.ContextWithTestID(t.Context(), t)

			var bp *apstra.TwoStageL3ClosClient
			var rzID string

			bpWait := new(sync.WaitGroup)
			bpWait.Add(1)
			go func() {
				bp = dctestobj.TestBlueprintA(t, ctx, client.Client)
				rzID = dctestobj.TestSecurityZoneA(t, ctx, bp)
				dctestobj.TestVirtualNetworkA(t, ctx, bp, rzID) // Create the RZ to force the need for per-RZ resource allocation
				bpWait.Done()
			}()

			ipv4PoolIDS := make([]string, ipv4PoolCount)
			for i := range ipv4PoolIDS {
				ipv4PoolIDS[i] = resourcetestobj.RandomIPv4Pool(t, ctx, client.Client)
			}
			asnPoolIDS := make([]string, asnPoolCount)
			for i := range asnPoolIDS {
				asnPoolIDS[i] = resourcetestobj.RandomASNPool(t, ctx, client.Client)
			}

			bpWait.Wait()
			t.Log(bp.Id())
			t.Log(rzID)

			// Retrieve leaf ASN allocations -- it should be empty
			t.Run("check_empty_leaf_asn", func(t *testing.T) {
				rga, err := bp.GetResourceAllocation(ctx, apstra.ResourceGroup{Name: enum.ResourceGroupLeafASN})
				require.NoError(t, err)
				require.Equal(t, rga.ResourceGroup.Name, enum.ResourceGroupLeafASN)
				require.Nil(t, rga.ResourceGroup.SecurityZoneID)
				require.Empty(t, rga.PoolIds)
			})

			// Set and check leaf ASN allocations
			t.Run("set_leaf_asn_1", func(t *testing.T) {
				rg := apstra.ResourceGroup{Name: enum.ResourceGroupLeafASN}
				set := apstra.ResourceGroupAllocation{
					ResourceGroup: rg,
					PoolIds:       []string{asnPoolIDS[0], asnPoolIDS[2]},
				}
				err := bp.SetResourceAllocation(ctx, set)
				require.NoError(t, err)

				get, err := bp.GetResourceAllocation(ctx, rg)
				require.NoError(t, err)
				require.Equal(t, rg.String(), get.ResourceGroup.String())
				require.Nil(t, get.ResourceGroup.SecurityZoneID)
				compare.SlicesAsSets(t, set.PoolIds, get.PoolIds, "expected pool IDs must equal actual pool IDs")
			})

			// Set and check leaf ASN allocations
			t.Run("set_leaf_asn_2", func(t *testing.T) {
				rg := apstra.ResourceGroup{Name: enum.ResourceGroupLeafASN}
				set := apstra.ResourceGroupAllocation{
					ResourceGroup: rg,
					PoolIds:       asnPoolIDS,
				}
				err := bp.SetResourceAllocation(ctx, set)
				require.NoError(t, err)

				get, err := bp.GetResourceAllocation(ctx, rg)
				require.NoError(t, err)
				require.Equal(t, rg.String(), get.ResourceGroup.String())
				require.Nil(t, get.ResourceGroup.SecurityZoneID)
				compare.SlicesAsSets(t, set.PoolIds, get.PoolIds, "expected pool IDs must equal actual pool IDs")
			})

			// Clear and check leaf ASN allocations
			t.Run("clear_leaf_asn", func(t *testing.T) {
				rg := apstra.ResourceGroup{Name: enum.ResourceGroupLeafASN}
				err := bp.SetResourceAllocation(ctx, apstra.ResourceGroupAllocation{ResourceGroup: rg})
				require.NoError(t, err)

				get, err := bp.GetResourceAllocation(ctx, rg)
				require.NoError(t, err)
				require.Equal(t, rg.String(), get.ResourceGroup.String())
				require.Nil(t, get.ResourceGroup.SecurityZoneID)
				require.Empty(t, get.PoolIds)
			})

			// Retrieve leaf RZ IPv4 loopback allocations -- it should be empty
			t.Run("check_empty_leaf_rz_loopback_ipv4", func(t *testing.T) {
				rga, err := bp.GetResourceAllocation(ctx, apstra.ResourceGroup{Name: enum.ResourceGroupLeafIPv4, SecurityZoneID: &rzID})
				require.NoError(t, err)
				require.Equal(t, rga.ResourceGroup.Name, enum.ResourceGroupLeafIPv4)
				require.NotNil(t, rga.ResourceGroup.SecurityZoneID)
				require.Empty(t, rga.PoolIds)
			})

			// Set and check leaf RZ IPv4 loopback allocations
			t.Run("set_leaf_rz_loopback_ipv4_1", func(t *testing.T) {
				rg := apstra.ResourceGroup{Name: enum.ResourceGroupLeafIPv4, SecurityZoneID: &rzID}
				set := apstra.ResourceGroupAllocation{
					ResourceGroup: rg,
					PoolIds:       []string{ipv4PoolIDS[0], ipv4PoolIDS[2]},
				}
				err := bp.SetResourceAllocation(ctx, set)
				require.NoError(t, err)

				get, err := bp.GetResourceAllocation(ctx, rg)
				require.NoError(t, err)
				require.Equal(t, rg.String(), get.ResourceGroup.String())
				require.NotNil(t, get.ResourceGroup.SecurityZoneID)
				compare.SlicesAsSets(t, set.PoolIds, get.PoolIds, "expected pool IDs must equal actual pool IDs")
			})

			// Set and check leaf RZ IPv4 loopback allocations
			t.Run("set_leaf_rz_loopback_ipv4_2", func(t *testing.T) {
				rg := apstra.ResourceGroup{Name: enum.ResourceGroupLeafIPv4, SecurityZoneID: &rzID}
				set := apstra.ResourceGroupAllocation{
					ResourceGroup: rg,
					PoolIds:       ipv4PoolIDS,
				}
				err := bp.SetResourceAllocation(ctx, set)
				require.NoError(t, err)

				get, err := bp.GetResourceAllocation(ctx, rg)
				require.NoError(t, err)
				require.Equal(t, rg.String(), get.ResourceGroup.String())
				require.NotNil(t, get.ResourceGroup.SecurityZoneID)
				compare.SlicesAsSets(t, set.PoolIds, get.PoolIds, "expected pool IDs must equal actual pool IDs")
			})

			// Clear and check leaf RZ IPv4 loopback allocations
			t.Run("clear_leaf_rz_loopback_ipv4", func(t *testing.T) {
				rg := apstra.ResourceGroup{Name: enum.ResourceGroupLeafIPv4, SecurityZoneID: &rzID}
				err := bp.SetResourceAllocation(ctx, apstra.ResourceGroupAllocation{ResourceGroup: rg})
				require.NoError(t, err)

				get, err := bp.GetResourceAllocation(ctx, rg)
				require.NoError(t, err)
				require.Equal(t, rg.String(), get.ResourceGroup.String())
				require.NotNil(t, get.ResourceGroup.SecurityZoneID)
				require.Empty(t, get.PoolIds)
			})
		})
	}
}
