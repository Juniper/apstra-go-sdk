// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"math/rand"
	"net/netip"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/stretchr/testify/require"
)

func TestTwoStageL3ClosClient_GetSecurityZoneInfo(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	securityZoneCount := rand.Intn(3) + 3

	compareSzDataToSzInfo := func(t *testing.T, a *SecurityZoneData, b *TwoStageL3ClosSecurityZoneInfo) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)

		require.Equal(t, a.Label, b.Label)
		require.Equal(t, a.SzType.String(), b.SecurityZoneType.String())
		require.Equal(t, a.VrfName, b.VrfName)
		require.Equal(t, a.VlanId, b.VlanId)
		require.NotNil(t, a.JunosEvpnIrbMode)
		require.NotNil(t, b.JunosEvpnIrbMode)
		require.Equal(t, *a.JunosEvpnIrbMode, *b.JunosEvpnIrbMode)
	}

	compareLoopbacksToSzInfo := func(t *testing.T, a map[ObjectId]SecurityZoneLoopback, b TwoStageL3ClosSecurityZoneInfo) {
		t.Helper()

		// reorganize member interfaces into a map keyed by interface ID
		loopbacks := make(map[ObjectId]TwoStageL3ClosSecurityZoneInfoInterface)
		for _, memberInterface := range b.MemberInterfaces {
			for _, loopback := range memberInterface.Loopbacks {
				loopbacks[loopback.Id] = loopback
			}
		}

		for ifId, loopbackInfo := range a {
			if loopback, ok := loopbacks[ifId]; ok {
				require.Equal(t, loopbackInfo.IPv4Addr.String(), loopback.Ipv4Addr.String())
				require.Equal(t, loopbackInfo.IPv6Addr.String(), loopback.Ipv6Addr.String())
			} else {
				t.Fatalf("loopback not found for %s", ifId)
			}
		}
	}

	// sub-test per client
	for _, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			if !compatibility.SecurityZoneLoopbackApiSupported.Check(client.client.apiVersion) {
				t.Skip("Security zone loopback api not supported")
			}

			bp := testBlueprintH(ctx, t, client.client)

			// get all systems
			var response struct {
				Nodes map[ObjectId]struct {
					Id       ObjectId `json:"id"`
					Label    string   `json:"label"`
					Hostname string   `json:"hostname"`
					Role     string   `json:"role"`
					Ipv4Addr netip.Prefix
					Ipv6Addr netip.Prefix
				} `json:"nodes"`
			}
			err = bp.Client().GetNodes(ctx, bp.Id(), NodeTypeSystem, &response)
			require.NoError(t, err)

			var vnBindings []VnBinding // we'll use this binding slice when creating VNs later

			// whittle down the systems to just leaf switches and add a record to to vnBindings
			for k, v := range response.Nodes {
				if v.Role != "leaf" {
					delete(response.Nodes, k)
					continue
				}
				vnBindings = append(vnBindings, VnBinding{SystemId: k})
			}

			szIds := make([]ObjectId, securityZoneCount)
			szDatas := make([]SecurityZoneData, securityZoneCount)
			VlanIds := make([]int, securityZoneCount) // dedicated vlan slice ensures no collisions
			randomIntsN(VlanIds, vlanMax-2)

			// create security zones
			for i := range securityZoneCount {
				szDatas[i] = SecurityZoneData{
					Label:            randString(6, "hex"),
					SzType:           SecurityZoneTypeEVPN,
					VrfName:          randString(6, "hex"),
					VlanId:           toPtr(Vlan(VlanIds[i])),
					JunosEvpnIrbMode: oneOf(&enum.JunosEvpnIrbModeAsymmetric, &enum.JunosEvpnIrbModeSymmetric),
				}

				// create security zone, record ID
				szIds[i], err = bp.CreateSecurityZone(ctx, &szDatas[i])
				require.NoError(t, err)

				// create a VN to activate the SZ is on leaf switches
				_, err := bp.CreateVirtualNetwork(ctx, &VirtualNetworkData{
					Ipv4Enabled:    true,
					Ipv6Enabled:    true,
					Label:          randString(6, "hex"),
					SecurityZoneId: szIds[i],
					VnBindings:     vnBindings,
					VnType:         enum.VnTypeVxlan,
				})
				require.NoError(t, err)
			}

			// set ipv4 and ipv6 addresses of each leaf in each security zone
			szidToIfidToLoopbackInfo := make(map[ObjectId]map[ObjectId]SecurityZoneLoopback, securityZoneCount)
			for _, szId := range szIds {
				szidToIfidToLoopbackInfo[szId], err = bp.GetSecurityZoneLoopbacks(ctx, szId)
				require.NoError(t, err)
				for k, v := range szidToIfidToLoopbackInfo[szId] {
					v.IPv4Addr = toPtr(netip.MustParsePrefix(randomIpv4().String() + "/32"))
					v.IPv6Addr = toPtr(netip.MustParsePrefix(randomIpv6().String() + "/128"))
					szidToIfidToLoopbackInfo[szId][k] = v
				}

				err := bp.SetSecurityZoneLoopbacks(ctx, szId, szidToIfidToLoopbackInfo[szId])
				require.NoError(t, err)
			}

			infos, err := bp.GetAllSecurityZoneInfos(ctx)
			require.NoError(t, err)
			require.Equal(t, securityZoneCount+1, len(infos)) // one extra for default SZ

			// drop the default SZ from infos since we didn't set addresses in there
			for k := range infos {
				if !sliceContains(szIds, k) {
					delete(infos, k)
				}
			}
			require.Equal(t, securityZoneCount, len(infos))

			for i, szId := range szIds {
				info, err := bp.GetSecurityZoneInfo(ctx, szId)
				require.NoError(t, err)

				// check security zone details
				compareSzDataToSzInfo(t, &szDatas[i], info)

				// check per-sz loopback addresses
				compareLoopbacksToSzInfo(t, szidToIfidToLoopbackInfo[szId], infos[szId])
			}
		})
	}
}

func TestTwoStageL3ClosClient_GetSecurityZoneInfoBogus(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for _, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			bp := testBlueprintA(ctx, t, client.client)

			_, err := bp.GetSecurityZoneInfo(ctx, "bogus")
			require.Error(t, err)
			var ace ClientErr
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.errType, ErrNotfound)
		})
	}
}
