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

func TestUpdateVirtualNetworkLeafBindings(t *testing.T) {
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

			// drop the first leaf from the slice, since we'll put a fixed binding and IP on it
			fixedLeafId := leafIds[0]
			leafIds = leafIds[1:]

			rzLabel := randString(6, "hex")
			rzId, err := bp.CreateSecurityZone(ctx, &SecurityZoneData{
				Label:   rzLabel,
				SzType:  SecurityZoneTypeEVPN,
				VrfName: rzLabel,
			})
			require.NoError(t, err)

			vnPrefix := randomPrefix(t, "10.0.0.0/8", 24)
			vnBinding := VnBinding{
				SystemId: fixedLeafId,
				VlanId:   toPtr(Vlan(rand.Intn(97) + 2)),
			}
			sviIp := SviIp{
				SystemId: fixedLeafId,
				Ipv4Addr: &net.IPNet{
					IP:   []byte{vnPrefix.IP[0], vnPrefix.IP[1], vnPrefix.IP[2], 2},
					Mask: vnPrefix.Mask,
				},
				Ipv4Mode: oneOf(enum.SviIpv4ModeEnabled, enum.SviIpv4ModeForced),
				Ipv6Mode: enum.SviIpv6ModeDisabled,
			}
			vnId, err := bp.CreateVirtualNetwork(ctx, &VirtualNetworkData{
				Ipv4Enabled:    true,
				Ipv4Subnet:     &vnPrefix,
				Label:          randString(6, "hex"),
				SecurityZoneId: rzId,
				VnType:         enum.VnTypeVxlan,
				VnBindings:     []VnBinding{vnBinding},
				SviIps:         []SviIp{sviIp},
			})
			require.NoError(t, err)

			// set up an slice with leaf counts we'll test with. Assuming 3 test leafs,
			// this slice will be: [0,1,2,3,2,1,0]
			testLeafCounts := make([]int, (2*len(leafIds))+1)
			for i := range (len(testLeafCounts) + 1) / 2 {
				testLeafCounts[i] = i
				testLeafCounts[len(testLeafCounts)-i-1] = i
			}

			for _, count := range testLeafCounts {
				requestBindings := make(map[ObjectId]*VnBinding, len(leafIds))
				for j, leafId := range leafIds {
					if j < count {
						requestBindings[leafId] = &VnBinding{
							AccessSwitchNodeIds: nil,
							SystemId:            leafIds[j],
							VlanId:              toPtr(Vlan(100*(count) + rand.Intn(100))),
						}
					} else {
						requestBindings[leafId] = nil
					}
				}

				requestSviIps := make(map[ObjectId]*SviIp, count)
				lastOctets := make([]int, count)
				randomIntsN(lastOctets, 253) // avoid .0, .1 and .255
				for j, leafId := range leafIds {
					if j < count {
						_, ipNet, err := net.ParseCIDR(vnPrefix.String())
						require.NoError(t, err)
						ipNet.IP[3] = byte(lastOctets[j])
						requestSviIps[leafId] = &SviIp{
							SystemId: leafIds[j],
							Ipv4Addr: ipNet,
							Ipv4Mode: oneOf(enum.SviIpv4ModeForced, enum.SviIpv4ModeEnabled),
							Ipv6Mode: enum.SviIpv6ModeDisabled,
						}
					} else {
						requestSviIps[leafId] = nil
					}
				}

				log.Printf("testing UpdateVirtualNetworkLeafBindings() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				err = bp.UpdateVirtualNetworkLeafBindings(ctx, VirtualNetworkBindingsRequest{
					VnId:               vnId,
					VnBindings:         requestBindings,
					SviIps:             requestSviIps,
					DhcpServiceEnabled: toPtr(DhcpServiceEnabled(false)),
				})
				require.NoError(t, err)

				vn, err := bp.GetVirtualNetwork(ctx, vnId)
				require.NoError(t, err)

				actualBindings := bindingSliceToMap(t, vn.Data.VnBindings)
				delete(actualBindings, fixedLeafId)
				for k, v := range requestBindings {
					if v == nil {
						delete(requestBindings, k)
					}
				}

				actualSviIps := sviSliceToMap(t, vn.Data.SviIps)
				delete(actualSviIps, fixedLeafId)
				for k, v := range requestSviIps {
					if v == nil {
						delete(requestSviIps, k)
					}
				}

				compareBindingsMaps(t, requestBindings, actualBindings)
				compareSviMaps(t, requestSviIps, actualSviIps)
			}
		})
	}
}
