// Copyright (c) Juniper Networks, Inc., 2024-2024.
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

	getLoopback := func(ctx context.Context, t *testing.T, bpClient *TwoStageL3ClosClient, id ObjectId) *SecurityZoneLoopback {
		t.Helper()
		var response struct {
			Nodes map[ObjectId]struct {
				IPv4Addr *string `json:"ipv4_addr"`
				IPv6Addr *string `json:"ipv6_addr"`
			} `json:"nodes"`
		}
		err := bpClient.GetNodes(ctx, NodeTypeInterface, &response)
		require.NoError(t, err)

		node, ok := response.Nodes[id]
		require.Truef(t, ok, "Didn't find interface node %s", id)
		require.NotNil(t, node.IPv4Addr)
		require.NotNil(t, node.IPv6Addr)

		ipv4Addr, err := netip.ParsePrefix(*node.IPv4Addr)
		require.NoError(t, err)
		ipv6Addr, err := netip.ParsePrefix(*node.IPv6Addr)
		require.NoError(t, err)

		return &SecurityZoneLoopback{
			IPv4Addr: &ipv4Addr,
			IPv6Addr: &ipv6Addr,
		}
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

			loopbackNodeId := getLoopbackNodeId(t, ctx, bpClient, leafIds[0], szId)

			for i, tCase := range testCases {
				t.Run(fmt.Sprintf("test_case_%d", i), func(t *testing.T) {
					err := bpClient.SetSecurityZoneLoopbacks(ctx, szId, map[ObjectId]SecurityZoneLoopback{
						loopbackNodeId: {
							IPv4Addr: tCase.ipv4,
							IPv6Addr: tCase.ipv6,
						},
					})
					require.NoError(t, err)

					actual := getLoopback(ctx, t, bpClient, loopbackNodeId)
					require.NotNil(t, actual.IPv4Addr)
					require.NotNil(t, actual.IPv6Addr)

					switch {
					case tCase.ipv4 == nil:
						require.Equalf(t, persistIpv4.String(), actual.IPv4Addr.String(),
							"we sent <nil>, so actual ipv4 address should use the old value %s, got %s",
							persistIpv4.String(), actual.IPv4Addr.String())
					case !tCase.ipv4.IsValid():
						require.Truef(t, ipv4PoolPrefix.Contains(actual.IPv4Addr.Addr()),
							"we sent <invalid>, so actual ipv4 address should fall within the pool prefix %s, got %s",
							ipv4PoolPrefix, tCase.ipv4)
					default:
						require.Equalf(t, tCase.ipv4.String(), actual.IPv4Addr.String(), "expected: %s actual %s",
							tCase.ipv4.String(), actual.IPv4Addr.String())
					}

					switch {
					case tCase.ipv6 == nil:
						require.Equalf(t, persistIpv6.String(), actual.IPv6Addr.String(),
							"we sent <nil>, so actual ipv6 address should use the old value %s, got %s",
							persistIpv6.String(), actual.IPv6Addr.String())
					case !tCase.ipv6.IsValid():
						require.Truef(t, ipv6PoolPrefix.Contains(actual.IPv6Addr.Addr()),
							"we sent <invalid>, so actual ipv6 address should fall within the pool prefix %s, got %s",
							ipv6PoolPrefix, tCase.ipv6)
					default:
						require.Equalf(t, tCase.ipv6.String(), actual.IPv6Addr.String(), "expected: %s actual %s",
							tCase.ipv6.String(), actual.IPv6Addr.String())
					}
				})
			}
		})
	}
}
