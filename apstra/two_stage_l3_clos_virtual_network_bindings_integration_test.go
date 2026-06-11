// Copyright (c) Juniper Networks, Inc., 2024-2026.
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

	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/stretchr/testify/require"
)

func compareSVIAddressing(t testing.TB, a, b datacenter.SVIAddressing) {
	t.Helper()

	require.Equal(t, a.SystemID, b.SystemID)

	require.Equal(t, a.IPv4Mode, b.IPv4Mode)
	if a.IPv4Addr != nil || b.IPv4Addr != nil {
		require.NotNil(t, a.IPv4Addr)
		require.NotNil(t, b.IPv4Addr)
		require.Equal(t, a.IPv4Addr.String(), b.IPv4Addr.String())
	}

	require.Equal(t, a.IPv6Mode, b.IPv6Mode)
	if a.IPv6Addr != nil || b.IPv6Addr != nil {
		require.NotNil(t, a.IPv6Addr)
		require.NotNil(t, b.IPv6Addr)
		require.Equal(t, a.IPv6Addr.String(), b.IPv6Addr.String())
	}
}

func compareVnBindings(t testing.TB, a, b datacenter.VNBinding, strict bool) {
	t.Helper()

	if len(a.AccessSwitchNodeIDs) != 0 || len(b.AccessSwitchNodeIDs) != 0 { // nil and [] slices are equal for our purpose
		compareSlices(t, a.AccessSwitchNodeIDs, b.AccessSwitchNodeIDs, "VnBindings.AccessSwitchNodeIDs")
	}

	require.Equal(t, a.SystemID, b.SystemID)

	if a.VLAN != nil || // the caller specified a VLAN, so we check it
		((a.VLAN != nil || b.VLAN != nil) && strict) { // strict mode means we always check
		require.NotNil(t, a.VLAN)
		require.NotNil(t, b.VLAN)
		require.Equal(t, a.VLAN, b.VLAN)
	}
}

func TestSetVirtualNetworkLeafBindings(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	bindingSliceToMap := func(t testing.TB, in []datacenter.VNBinding) map[ObjectId]*datacenter.VNBinding {
		t.Helper()

		result := make(map[ObjectId]*datacenter.VNBinding, len(in))
		for _, binding := range in {
			result[ObjectId(binding.SystemID)] = &binding
		}

		if len(in) != len(result) {
			t.Fatalf("after converting binding slice to map, had %d bindings, expected %d", len(result), len(in))
		}

		return result
	}

	sviSliceToMap := func(t testing.TB, in []datacenter.SVIAddressing) map[ObjectId]*datacenter.SVIAddressing {
		t.Helper()

		result := make(map[ObjectId]*datacenter.SVIAddressing, len(in))
		for _, sviIp := range in {
			result[ObjectId(sviIp.SystemID)] = &sviIp
		}

		if len(in) != len(result) {
			t.Fatalf("after converting svi slice to map, had %d svis, expected %d", len(result), len(in))
		}

		return result
	}

	compareBindingsMaps := func(t testing.TB, a, b map[ObjectId]*datacenter.VNBinding) {
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

	compareSviMaps := func(t testing.TB, a, b map[ObjectId]*datacenter.SVIAddressing) {
		t.Helper()

		if len(a) != len(b) {
			t.Fatalf("svi count mismatch: %d vs %d", len(a), len(b))
		}

		for id, sviIpA := range a {
			sviIpB, ok := b[id]
			if !ok {
				t.Fatalf("svi %q from 'a' map not found in 'b' map", id)
			}

			compareSVIAddressing(t, *sviIpA, *sviIpB)
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
			rzId, err := bp.CreateSecurityZone(ctx, datacenter.SecurityZone{
				Label:   rzLabel,
				Type:    enum.SecurityZoneTypeEVPN,
				VRFName: rzLabel,
			})
			require.NoError(t, err)

			vnPrefix := randomPrefix(t, "10.0.0.0/8", 24)
			vnId, err := bp.CreateVirtualNetwork(ctx, datacenter.VirtualNetwork{
				IPv4Enabled:    true,
				IPv4Subnet:     &vnPrefix,
				Label:          randString(6, "hex"),
				SecurityZoneID: rzId,
				Type:           enum.VnTypeVxlan,
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
				bindings := make(map[ObjectId]*datacenter.VNBinding, count)
				for j := range count {
					bindings[leafIds[j]] = &datacenter.VNBinding{
						AccessSwitchNodeIDs: nil,
						SystemID:            string(leafIds[j]),
						VLAN:                pointer.To(uint16(100*(count) + rand.Intn(100))),
					}
				}

				sviIps := make(map[ObjectId]*datacenter.SVIAddressing, count)
				lastOctets := make([]int, count)
				randomIntsN(lastOctets, 253) // 253 possible values while avoiding .0, .1 and .255
				for j := range count {
					_, ipNet, err := net.ParseCIDR(vnPrefix.String())
					require.NoError(t, err)
					ipNet.IP[3] = byte(lastOctets[j] + 2) // add 2 to to avoid .0 and .1
					sviIps[leafIds[j]] = &datacenter.SVIAddressing{
						SystemID: string(leafIds[j]),
						IPv4Addr: ipNet,
						IPv4Mode: oneOf(enum.IPv4SVIModeForced, enum.IPv4SVIModeEnabled),
						IPv6Mode: enum.IPv6SVIModeDisabled,
					}
				}

				log.Printf("testing SetVirtualNetworkLeafBindings() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				err = bp.SetVirtualNetworkLeafBindings(ctx, VirtualNetworkBindingsRequest{
					VnId:               ObjectId(vnId),
					VnBindings:         bindings,
					SviIps:             sviIps,
					DhcpServiceEnabled: toPtr(datacenter.DHCPServiceEnabled(false)),
				})
				require.NoError(t, err)

				vn, err := bp.GetVirtualNetwork(ctx, vnId)
				require.NoError(t, err)
				compareBindingsMaps(t, bindings, bindingSliceToMap(t, vn.Bindings))
				compareSviMaps(t, sviIps, sviSliceToMap(t, vn.SVIIPs))
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

	bindingSliceToMap := func(t testing.TB, in []datacenter.VNBinding) map[ObjectId]*datacenter.VNBinding {
		t.Helper()

		result := make(map[ObjectId]*datacenter.VNBinding, len(in))
		for _, binding := range in {
			result[ObjectId(binding.SystemID)] = &binding
		}

		if len(in) != len(result) {
			t.Fatalf("after converting binding slice to map, had %d bindings, expected %d", len(result), len(in))
		}

		return result
	}

	sviSliceToMap := func(t testing.TB, in []datacenter.SVIAddressing) map[ObjectId]*datacenter.SVIAddressing {
		t.Helper()

		result := make(map[ObjectId]*datacenter.SVIAddressing, len(in))
		for _, sviIp := range in {
			result[ObjectId(sviIp.SystemID)] = &sviIp
		}

		if len(in) != len(result) {
			t.Fatalf("after converting svi slice to map, had %d svis, expected %d", len(result), len(in))
		}

		return result
	}

	compareBindingsMaps := func(t testing.TB, a, b map[ObjectId]*datacenter.VNBinding) {
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

	compareSviMaps := func(t testing.TB, a, b map[ObjectId]*datacenter.SVIAddressing) {
		t.Helper()

		if len(a) != len(b) {
			t.Fatalf("svi count mismatch: %d vs %d", len(a), len(b))
		}

		for id, sviIpA := range a {
			sviIpB, ok := b[id]
			if !ok {
				t.Fatalf("svi %q from 'a' map not found in 'b' map", id)
			}

			compareSVIAddressing(t, *sviIpA, *sviIpB)
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
			rzId, err := bp.CreateSecurityZone(ctx, datacenter.SecurityZone{
				Label:   rzLabel,
				Type:    enum.SecurityZoneTypeEVPN,
				VRFName: rzLabel,
			})
			require.NoError(t, err)

			vnPrefix := randomPrefix(t, "10.0.0.0/8", 24)
			vnBinding := datacenter.VNBinding{
				SystemID: string(fixedLeafId),
				VLAN:     pointer.To(uint16(rand.Intn(89) + 11)),
			}
			sviIp := datacenter.SVIAddressing{
				SystemID: string(fixedLeafId),
				IPv4Addr: &net.IPNet{
					IP:   []byte{vnPrefix.IP[0], vnPrefix.IP[1], vnPrefix.IP[2], 2},
					Mask: vnPrefix.Mask,
				},
				IPv4Mode: oneOf(enum.IPv4SVIModeEnabled, enum.IPv4SVIModeForced),
				IPv6Mode: enum.IPv6SVIModeDisabled,
			}
			vnId, err := bp.CreateVirtualNetwork(ctx, datacenter.VirtualNetwork{
				IPv4Enabled:    true,
				IPv4Subnet:     &vnPrefix,
				Label:          randString(6, "hex"),
				SecurityZoneID: rzId,
				Type:           enum.VnTypeVxlan,
				Bindings:       []datacenter.VNBinding{vnBinding},
				SVIIPs:         []datacenter.SVIAddressing{sviIp},
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
				requestBindings := make(map[ObjectId]*datacenter.VNBinding, len(leafIds))
				for j, leafId := range leafIds {
					if j < count {
						requestBindings[leafId] = &datacenter.VNBinding{
							AccessSwitchNodeIDs: nil,
							SystemID:            string(leafIds[j]),
							VLAN:                pointer.To(uint16(100*(count) + rand.Intn(100))),
						}
					} else {
						requestBindings[leafId] = nil
					}
				}

				requestSviIps := make(map[ObjectId]*datacenter.SVIAddressing, count)
				lastOctets := make([]int, count)
				randomIntsN(lastOctets, 253) // 253 possible values while avoiding .0, .1 and .255
				for j, leafId := range leafIds {
					if j < count {
						_, ipNet, err := net.ParseCIDR(vnPrefix.String())
						require.NoError(t, err)
						ipNet.IP[3] = byte(lastOctets[j] + 2) // add 2 to to avoid .0 and .1
						requestSviIps[leafId] = &datacenter.SVIAddressing{
							SystemID: string(leafIds[j]),
							IPv4Addr: ipNet,
							IPv4Mode: oneOf(enum.IPv4SVIModeForced, enum.IPv4SVIModeEnabled),
							IPv6Mode: enum.IPv6SVIModeDisabled,
						}
					} else {
						requestSviIps[leafId] = nil
					}
				}

				log.Printf("testing UpdateVirtualNetworkLeafBindings() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				err = bp.UpdateVirtualNetworkLeafBindings(ctx, VirtualNetworkBindingsRequest{
					VnId:               ObjectId(vnId),
					VnBindings:         requestBindings,
					SviIps:             requestSviIps,
					DhcpServiceEnabled: toPtr(datacenter.DHCPServiceEnabled(false)),
				})
				require.NoError(t, err)

				vn, err := bp.GetVirtualNetwork(ctx, vnId)
				require.NoError(t, err)

				actualBindings := bindingSliceToMap(t, vn.Bindings)
				delete(actualBindings, fixedLeafId)
				for k, v := range requestBindings {
					if v == nil {
						delete(requestBindings, k)
					}
				}

				actualSviIps := sviSliceToMap(t, vn.SVIIPs)
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
