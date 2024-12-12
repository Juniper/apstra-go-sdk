// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sort"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/stretchr/testify/require"
)

func TestSetVirtualNetworkLeafBindings(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	bindingSliceToMap := func(t testing.TB, in []VnBinding) map[ObjectId]*VnBinding {
		t.Helper()

		result := make(map[ObjectId]*VnBinding, len(in))
		for _, binding := range in {
			result[binding.SystemId] = &binding
		}

		if len(in) != len(result) {
			t.Fatalf("after converting binding slice to map, had %d bindings, expected %d", len(result), len(in))
		}

		return result
	}

	sviSliceToMap := func(t testing.TB, in []SviIp) map[ObjectId]*SviIp {
		t.Helper()

		result := make(map[ObjectId]*SviIp, len(in))
		for _, sviIp := range in {
			result[sviIp.SystemId] = &sviIp
		}

		if len(in) != len(result) {
			t.Fatalf("after converting svi slice to map, had %d svis, expected %d", len(result), len(in))
		}

		return result
	}

	compareBindingsMaps := func(t testing.TB, a, b map[ObjectId]*VnBinding) {
		t.Helper()

		if len(a) != len(b) {
			t.Fatalf("bindings count mismatch: %d vs %d", len(a), len(b))
		}

		for id, bindingA := range a {
			bindingB, ok := b[id]
			if !ok {
				t.Fatalf("binding %q from 'a' map not found in 'b' map", id)
			}

			compareVnBindings(t, *bindingA, *bindingB, true)
		}
	}

	compareSviMaps := func(t testing.TB, a, b map[ObjectId]*SviIp) {
		t.Helper()

		if len(a) != len(b) {
			t.Fatalf("svi count mismatch: %d vs %d", len(a), len(b))
		}

		for id, sviIpA := range a {
			sviIpB, ok := b[id]
			if !ok {
				t.Fatalf("svi %q from 'a' map not found in 'b' map", id)
			}

			comapareSviIps(t, *sviIpA, *sviIpB)
		}
	}

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			if !compatibility.EmptyVnBindingsOk.Check(client.client.apiVersion) {
				t.Skipf("test applies only to versions %q", compatibility.EmptyVnBindingsOk)
			}

			bp := testBlueprintC(ctx, t, client.client)

			leafIds, err := getSystemIdsByRole(ctx, bp, "leaf")
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(leafIds), 2, "test requires at least 2 leaf switches")
			require.LessOrEqual(t, len(leafIds), 40, "test requires no more than 40 leaf switches")
			sort.Slice(leafIds, func(i, j int) bool { return leafIds[i] < leafIds[j] })

			rzLabel := randString(6, "hex")
			rzId, err := bp.CreateSecurityZone(ctx, &SecurityZoneData{
				Label:   rzLabel,
				SzType:  SecurityZoneTypeEVPN,
				VrfName: rzLabel,
			})
			require.NoError(t, err)

			vnPrefix := randomPrefix(t, "10.0.0.0/8", 24)
			vnId, err := bp.CreateVirtualNetwork(ctx, &VirtualNetworkData{
				Ipv4Enabled:    true,
				Ipv4Subnet:     &vnPrefix,
				Label:          randString(6, "hex"),
				SecurityZoneId: rzId,
				VnType:         enum.VnTypeVxlan,
			})
			require.NoError(t, err)

			// set up an slice with leaf counts we'll test with. Assuming 4 leafs in the BP,
			// this slice will be: [0,1,2,3,4,3,2,1,0]
			testLeafCounts := make([]int, (2*len(leafIds))+1)
			for i := range (len(testLeafCounts) + 1) / 2 {
				testLeafCounts[i] = i
				testLeafCounts[len(testLeafCounts)-i-1] = i
			}

			for _, count := range testLeafCounts {
				bindings := make(map[ObjectId]*VnBinding, count)
				for j := range count {
					bindings[leafIds[j]] = &VnBinding{
						AccessSwitchNodeIds: nil,
						SystemId:            leafIds[j],
						VlanId:              toPtr(Vlan(100*(count) + rand.Intn(100))),
					}
				}

				sviIps := make(map[ObjectId]*SviIp, count)
				lastOctets := make([]int, count)
				randomIntsN(lastOctets, 253) // avoid .0, .1 and .255
				for j := range count {
					_, ipNet, err := net.ParseCIDR(vnPrefix.String())
					require.NoError(t, err)
					ipNet.IP[3] = byte(lastOctets[j])
					sviIps[leafIds[j]] = &SviIp{
						SystemId: leafIds[j],
						Ipv4Addr: ipNet,
						Ipv4Mode: oneOf(enum.SviIpv4ModeForced, enum.SviIpv4ModeEnabled),
						Ipv6Mode: enum.SviIpv6ModeDisabled,
					}
				}

				log.Printf("testing SetVirtualNetworkLeafBindings() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				err = bp.SetVirtualNetworkLeafBindings(ctx, VirtualNetworkBindingsRequest{
					VnId:               vnId,
					VnBindings:         bindings,
					SviIps:             sviIps,
					DhcpServiceEnabled: toPtr(DhcpServiceEnabled(false)),
				})
				require.NoError(t, err)

				vn, err := bp.GetVirtualNetwork(ctx, vnId)
				require.NoError(t, err)
				compareBindingsMaps(t, bindings, bindingSliceToMap(t, vn.Data.VnBindings))
				compareSviMaps(t, sviIps, sviSliceToMap(t, vn.Data.SviIps))
			}
		})
	}
}
