// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"net/netip"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/stretchr/testify/require"
)

func TestTwoStageL3ClosClient_SetSecurityZoneLoopbacks(t *testing.T) {
	ctx := context.Background()

	mustParsePrefixPtr := func(s string) *netip.Prefix {
		p := netip.MustParsePrefix(s)
		return &p
	}

	ipV4PoolId := ObjectId("Private-10_0_0_0-8")
	ipv4PoolPrefix := netip.MustParsePrefix("10.0.0.0/8")

	ipV6PoolId := ObjectId("Private-fc01-a05-fab-48")
	ipv6PoolPrefix := netip.MustParsePrefix("fc01:a05:fab::/48")

	type testCase struct {
		ipv4 *netip.Prefix
		ipv6 *netip.Prefix
	}

	persistIpv4 := mustParsePrefixPtr("192.0.2.8/32")
	persistIpv6 := mustParsePrefixPtr("3fff::8/128")
	testCases := []testCase{
		{
			ipv4: mustParsePrefixPtr("192.0.2.7/32"),
			ipv6: mustParsePrefixPtr("3fff::7/128"),
		},
		{
			ipv4: &netip.Prefix{},
			ipv6: &netip.Prefix{},
		},
		{
			ipv4: persistIpv4,
			ipv6: persistIpv6,
		},
		{
			ipv4: nil,
			ipv6: nil,
		},
	}

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	getLoopbackNodeId := func(t *testing.T, ctx context.Context, bpClient *TwoStageL3ClosClient, systemId, securityZoenId ObjectId) ObjectId {
		query := new(PathQuery).
			SetBlueprintId(bpClient.Id()).
			SetClient(bpClient.client).
			Node([]QEEAttribute{
				NodeTypeSystem.QEEAttribute(),
				{"id", QEStringVal(systemId)},
			}).
			Out([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeInterface.QEEAttribute(),
				{Key: "if_type", Value: QEStringVal("loopback")},
				{Key: "name", Value: QEStringVal("n_interface")},
			}).
			In([]QEEAttribute{RelationshipTypeMemberInterfaces.QEEAttribute()}).
			Node([]QEEAttribute{NodeTypeSecurityZoneInstance.QEEAttribute()}).
			In([]QEEAttribute{RelationshipTypeInstantiatedBy.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeSecurityZone.QEEAttribute(),
				{"id", QEStringVal(securityZoenId)},
			})

		var response struct {
			Items []struct {
				Interface struct {
					Id ObjectId `json:"id"`
				} `json:"n_interface"`
			} `json:"items"`
		}

		err := query.Do(ctx, &response)
		require.NoError(t, err)
		require.Equalf(t, 1, len(response.Items), "expected 1 loopback, found %d", len(response.Items))

		return response.Items[0].Interface.Id
	}

	for _, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			if !compatibility.SecurityZoneLoopbackApiSupported.Check(client.client.apiVersion) {
				t.Skipf("skipping due to version constraint %q", compatibility.SecurityZoneLoopbackApiSupported)
			}

			// create a blueprint with IPv6 enabled
			bpClient := testBlueprintH(ctx, t, client.client)

			// fetch all security zones
			szs, err := bpClient.GetAllSecurityZones(ctx)
			require.NoError(t, err)
			require.Greater(t, len(szs), 0)

			// find the default security zone ID
			var szId ObjectId
			for _, sz := range szs {
				if sz.Data.Label == "Default routing zone" {
					szId = sz.Id
				}
			}
			require.NotEmpty(t, szId)

			// assign an IPv4 pool to leaf loopbacks so that we can "remove" (cause it to revert to a pool address) a loopback IPv4 address
			err = bpClient.SetResourceAllocation(ctx, &ResourceGroupAllocation{
				ResourceGroup: ResourceGroup{
					Type: ResourceTypeIp4Pool,
					Name: ResourceGroupNameLeafIp4,
				},
				PoolIds: []ObjectId{ipV4PoolId},
			})
			require.NoError(t, err)

			// assign an IPv6 pool to leaf loopbacks so that we can "remove" (cause it to revert to a pool) address a loopback IPv6 address
			err = bpClient.SetResourceAllocation(ctx, &ResourceGroupAllocation{
				ResourceGroup: ResourceGroup{
					Type: ResourceTypeIp6Pool,
					Name: ResourceGroupNameLeafIp6,
				},
				PoolIds: []ObjectId{ipV6PoolId},
			})
			require.NoError(t, err)

			leafIds, err := getSystemIdsByRole(ctx, bpClient, "leaf")
			require.NoError(t, err)
			require.Greater(t, len(leafIds), 0)
			leafId := leafIds[0]

			loopbackNodeId := getLoopbackNodeId(t, ctx, bpClient, leafId, szId)

			for i, tCase := range testCases {
				t.Run(fmt.Sprintf("test_case_%d", i), func(t *testing.T) {
					err := bpClient.SetSecurityZoneLoopbacks(ctx, szId, map[ObjectId]SecurityZoneLoopback{
						loopbackNodeId: {
							IPv4Addr: tCase.ipv4,
							IPv6Addr: tCase.ipv6,
						},
					})
					require.NoError(t, err)

					// fetch by interface id
					actualByIfId, err := bpClient.GetSecurityZoneLoopbackByInterfaceId(ctx, loopbackNodeId)
					require.NoError(t, err)
					require.Equal(t, loopbackNodeId, actualByIfId.Id)
					require.Equal(t, szId, actualByIfId.SecurityZoneId)
					require.Equal(t, 0, actualByIfId.LoopbackId)

					// fetch by system id
					actualsBySysId, err := bpClient.GetSecurityZoneLoopbacksBySystemId(ctx, leafId)
					require.NoError(t, err)
					var actualBySysId *SecurityZoneLoopback
					for _, szl := range actualsBySysId {
						if szl.Id == loopbackNodeId {
							actualBySysId = &szl
							break
						}
					}
					require.NotNil(t, actualBySysId)
					require.Equal(t, loopbackNodeId, actualBySysId.Id)
					require.Equal(t, szId, actualBySysId.SecurityZoneId)
					require.Equal(t, 0, actualBySysId.LoopbackId)

					// fetch all
					actualMap, err := bpClient.GetSecurityZoneLoopbacks(ctx, szId)
					require.NoError(t, err)
					require.Contains(t, actualMap, loopbackNodeId)
					require.Equal(t, loopbackNodeId, actualMap[loopbackNodeId].Id)
					require.Equal(t, szId, actualMap[loopbackNodeId].SecurityZoneId)
					require.Equal(t, 0, actualMap[loopbackNodeId].LoopbackId)

					switch {
					case tCase.ipv4 == nil:
						require.Equalf(t, persistIpv4.String(), actualByIfId.IPv4Addr.String(),
							"we sent <nil>, so actual ipv4 address (by if id) should use the old value %s, got %s",
							persistIpv4.String(), actualByIfId.IPv4Addr.String())
						require.Equalf(t, persistIpv4.String(), actualMap[loopbackNodeId].IPv4Addr.String(),
							"we sent <nil>, so actual ipv4 address (by map) should use the old value %s, got %s",
							persistIpv4.String(), actualMap[loopbackNodeId].IPv4Addr.String())
						require.Equalf(t, persistIpv4.String(), actualBySysId.IPv4Addr.String(),
							"we sent <nil>, so actual ipv4 address (by sys id) should use the old value %s, got %s",
							persistIpv4.String(), actualBySysId.IPv4Addr.String())
					case !tCase.ipv4.IsValid():
						require.Truef(t, ipv4PoolPrefix.Contains(actualByIfId.IPv4Addr.Addr()),
							"we sent <invalid>, so actual (by if id) ipv4 address should fall within the pool prefix %s, got %s",
							ipv4PoolPrefix, actualByIfId.IPv4Addr)
						require.Truef(t, ipv4PoolPrefix.Contains(actualMap[loopbackNodeId].IPv4Addr.Addr()),
							"we sent <invalid>, so actual (by map) ipv4 address should fall within the pool prefix %s, got %s",
							ipv4PoolPrefix, actualMap[loopbackNodeId])
						require.Truef(t, ipv4PoolPrefix.Contains(actualBySysId.IPv4Addr.Addr()),
							"we sent <invalid>, so actual (by sys id) ipv4 address should fall within the pool prefix %s, got %s",
							ipv4PoolPrefix, actualBySysId.IPv4Addr)
					default:
						require.Equalf(t, tCase.ipv4.String(), actualByIfId.IPv4Addr.String(), "expected: %s actual (by if id) %s",
							tCase.ipv4.String(), actualByIfId.IPv4Addr.String())
						require.Equalf(t, tCase.ipv4.String(), actualByIfId.IPv4Addr.String(), "expected: %s actual (by map) %s",
							tCase.ipv4.String(), actualMap[loopbackNodeId].IPv4Addr.String())
						require.Equalf(t, tCase.ipv4.String(), actualBySysId.IPv4Addr.String(), "expected: %s actual (by sys id) %s",
							tCase.ipv4.String(), actualBySysId.IPv4Addr.String())
					}

					switch {
					case tCase.ipv6 == nil:
						require.Equalf(t, persistIpv6.String(), actualByIfId.IPv6Addr.String(),
							"we sent <nil>, so actual ipv6 address (by if id) should use the old value %s, got %s",
							persistIpv6.String(), actualByIfId.IPv6Addr.String())
						require.Equalf(t, persistIpv6.String(), actualMap[loopbackNodeId].IPv6Addr.String(),
							"we sent <nil>, so actual ipv6 address (by map) should use the old value %s, got %s",
							persistIpv6.String(), actualMap[loopbackNodeId].IPv6Addr.String())
						require.Equalf(t, persistIpv6.String(), actualBySysId.IPv6Addr.String(),
							"we sent <nil>, so actual ipv6 address (by sys id) should use the old value %s, got %s",
							persistIpv6.String(), actualBySysId.IPv6Addr.String())
					case !tCase.ipv6.IsValid():
						require.Truef(t, ipv6PoolPrefix.Contains(actualByIfId.IPv6Addr.Addr()),
							"we sent <invalid>, so actual (by if id) ipv6 address should fall within the pool prefix %s, got %s",
							ipv6PoolPrefix, actualByIfId.IPv6Addr)
						require.Truef(t, ipv6PoolPrefix.Contains(actualMap[loopbackNodeId].IPv6Addr.Addr()),
							"we sent <invalid>, so actual (by map) ipv6 address should fall within the pool prefix %s, got %s",
							ipv6PoolPrefix, actualMap[loopbackNodeId])
						require.Truef(t, ipv6PoolPrefix.Contains(actualBySysId.IPv6Addr.Addr()),
							"we sent <invalid>, so actual (by sys d) ipv6 address should fall within the pool prefix %s, got %s",
							ipv6PoolPrefix, actualBySysId.IPv6Addr)
					default:
						require.Equalf(t, tCase.ipv6.String(), actualByIfId.IPv6Addr.String(), "expected: %s actual (by if id) %s",
							tCase.ipv6.String(), actualByIfId.IPv6Addr.String())
						require.Equalf(t, tCase.ipv6.String(), actualByIfId.IPv6Addr.String(), "expected: %s actual (by map) %s",
							tCase.ipv6.String(), actualMap[loopbackNodeId].IPv6Addr.String())
						require.Equalf(t, tCase.ipv6.String(), actualBySysId.IPv6Addr.String(), "expected: %s actual (by sys id) %s",
							tCase.ipv6.String(), actualBySysId.IPv6Addr.String())
					}
				})
			}
		})
	}
}
