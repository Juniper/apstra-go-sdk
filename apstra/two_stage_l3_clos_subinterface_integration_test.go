// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func TestUpdateTwoStageL3ClosSubinterface(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compareSubinterfaces := func(t testing.TB, a, b TwoStageL3ClosSubinterface) {
		t.Helper()
		require.Equal(t, a.Ipv4AddrType, b.Ipv4AddrType)
		require.Equal(t, a.Ipv6AddrType, b.Ipv6AddrType)
		require.Equal(t, a.Ipv4Addr.String(), b.Ipv4Addr.String())
		require.Equal(t, a.Ipv6Addr.String(), b.Ipv6Addr.String())
	}

	compareEndpoints := func(t testing.TB, a, b TwoStageL3ClosSubinterfaceLinkEndpoint) {
		t.Helper()
		require.Equal(t, a.InterfaceId, b.InterfaceId)
		require.Equal(t, a.SubinterfaceId, b.SubinterfaceId)
		require.Equal(t, a.System.Id, b.System.Id)
		require.Equal(t, a.System.Label, b.System.Label)
		require.Equal(t, a.System.Role, b.System.Role)
		compareSubinterfaces(t, a.Subinterface, b.Subinterface)
	}

	compareLinks := func(t testing.TB, a, b *TwoStageL3ClosSubinterfaceLink) {
		t.Helper()
		require.NotNil(t, a)
		require.NotNil(t, b)
		require.EqualValues(t, len(a.Endpoints), 2)
		require.EqualValues(t, len(b.Endpoints), 2)
		require.EqualValues(t, a.Id, b.Id)
		require.EqualValues(t, a.SzId, b.SzId)
		require.EqualValues(t, a.SzLabel, b.SzLabel)

		if a.VlanId == nil {
			require.Nil(t, b.VlanId)
		} else {
			require.NotNil(t, b.VlanId)
			require.EqualValues(t, *a.VlanId, *b.VlanId)
		}

		require.Contains(t, []ObjectId{a.Endpoints[0].SubinterfaceId, a.Endpoints[1].SubinterfaceId}, b.Endpoints[0].SubinterfaceId)
		require.Contains(t, []ObjectId{a.Endpoints[0].SubinterfaceId, a.Endpoints[1].SubinterfaceId}, b.Endpoints[1].SubinterfaceId)
		require.Contains(t, []ObjectId{b.Endpoints[0].SubinterfaceId, b.Endpoints[1].SubinterfaceId}, a.Endpoints[0].SubinterfaceId)
		require.Contains(t, []ObjectId{b.Endpoints[0].SubinterfaceId, b.Endpoints[1].SubinterfaceId}, a.Endpoints[1].SubinterfaceId)

		if a.Endpoints[0].SubinterfaceId == b.Endpoints[0].SubinterfaceId {
			compareEndpoints(t, a.Endpoints[0], b.Endpoints[0])
			compareEndpoints(t, a.Endpoints[1], b.Endpoints[1])
		} else {
			compareEndpoints(t, a.Endpoints[0], b.Endpoints[1])
			compareEndpoints(t, a.Endpoints[1], b.Endpoints[0])
		}
	}

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			// create a blueprint and enable IPv6
			bp := testBlueprintC(ctx, t, client.client)
			settings, err := bp.GetFabricSettings(ctx)
			require.NoError(t, err)
			settings.Ipv6Enabled = toPtr(true)
			require.NoError(t, bp.SetFabricSettings(ctx, settings))

			// create a security zone within the blueprint
			szName := randString(6, "hex")
			szId, err := bp.CreateSecurityZone(ctx, SecurityZone{
				Label:   szName,
				Type:    enum.SecurityZoneTypeEVPN,
				VRFName: szName,
				VNI:     toPtr(rand.Intn(1000) + 10000),
			})
			require.NoError(t, err)

			// create a simple "ip link" connectivity template which uses the security zone
			ct := ConnectivityTemplate{
				Label: randString(6, "hex"),
				Subpolicies: []*ConnectivityTemplatePrimitive{{
					Label: "IP Link",
					Attributes: &ConnectivityTemplatePrimitiveAttributesAttachLogicalLink{
						Label:              "",
						SecurityZone:       (*ObjectId)(&szId),
						IPv4AddressingType: CtPrimitiveIPv4AddressingTypeNone,
						IPv6AddressingType: CtPrimitiveIPv6AddressingTypeLinkLocal,
					},
				}},
			}
			require.NoError(t, ct.SetIds())
			require.NoError(t, ct.SetUserData())
			require.NoError(t, bp.CreateConnectivityTemplate(ctx, &ct))

			// prep a graph query which finds all server-facing switch interfaces
			query := new(PathQuery).SetClient(client.client).SetBlueprintId(bp.Id()).
				Node([]QEEAttribute{
					NodeTypeSystem.QEEAttribute(),
					{Key: "system_type", Value: QEStringVal(SystemTypeServer.String())},
				}).
				Out([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
				Node([]QEEAttribute{NodeTypeInterface.QEEAttribute()}).
				Out([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
				Node([]QEEAttribute{
					NodeTypeLink.QEEAttribute(),
					{Key: "link_type", Value: QEStringVal("ethernet")},
				}).
				In([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
				Node([]QEEAttribute{
					NodeTypeInterface.QEEAttribute(),
					{Key: "name", Value: QEStringVal("n_interface")},
				}).
				In([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
				Node([]QEEAttribute{
					NodeTypeSystem.QEEAttribute(),
					{Key: "system_type", Value: QEStringVal(SystemTypeSwitch.String())},
				})

			// prep the response object and run the query
			var response struct {
				Items []struct {
					Interface struct {
						Id ObjectId `json:"id"`
					} `json:"n_interface"`
				} `json:"items"`
			}
			require.NoError(t, query.Do(ctx, &response))
			require.Greater(t, len(response.Items), 0)

			// extract interface IDs from the query response
			intfIds := make([]ObjectId, len(response.Items))
			for i, item := range response.Items {
				intfIds[i] = item.Interface.Id
			}

			// collect a CT application map for our switch port interfaces
			apMap, err := bp.GetConnectivityTemplatesByApplicationPoints(ctx, intfIds)
			require.NoError(t, err)
			require.EqualValuesf(t, len(response.Items), len(apMap),
				"found %d server-facing ports, but %d application points - these should match",
				len(response.Items), len(apMap))

			// "check the box" so that our CT will be assigned to all eligible application points
			for k, v := range apMap {
				if _, ok := v[*ct.Id]; !ok {
					t.Fatalf("required CT ID %q not found in CT map for interface %q", *ct.Id, k)
				}
				apMap[k][*ct.Id] = true // indicates the CT should be attached to this port
			}

			// assign the CT to the application points
			err = bp.SetApplicationPointsConnectivityTemplates(ctx, apMap)
			require.NoError(t, err)

			// fetch the links created by assigning CTs
			links, err := bp.GetAllSubinterfaceLinks(ctx)
			require.NoError(t, err)

			// prep new addresses for those links (two endpoints per link)
			subinterfaceConfigs := make(map[ObjectId]TwoStageL3ClosSubinterface, len(links)*2)
			for _, link := range links {
				require.EqualValuesf(t, len(link.Endpoints), 2, "link should have two endpoints, got %d", len(link.Endpoints))

				// pair of /31s for each link
				v4prefixes := make([]net.IPNet, 2)
				v4prefixes[0] = randomSlash31(t)
				v4prefixes[1] = randomSlash31(t)
				copy(v4prefixes[1].IP, v4prefixes[0].IP)
				copy(v4prefixes[1].Mask, v4prefixes[0].Mask)
				v4prefixes[1].IP[3]++

				// pair of /127s for each link
				v6prefixes := make([]net.IPNet, 2)
				v6prefixes[0] = randomSlash127(t)
				v6prefixes[1] = randomSlash127(t)
				copy(v6prefixes[1].IP, v6prefixes[0].IP)
				copy(v6prefixes[1].Mask, v6prefixes[0].Mask)
				v6prefixes[1].IP[15]++

				// prep an API payload for each end of the link
				for i, endpoint := range link.Endpoints {
					subinterfaceConfigs[endpoint.SubinterfaceId] = TwoStageL3ClosSubinterface{
						Ipv4AddrType: enum.InterfaceNumberingIpv4TypeNumbered,
						Ipv6AddrType: enum.InterfaceNumberingIpv6TypeNumbered,
						Ipv4Addr:     &v4prefixes[i],
						Ipv6Addr:     &v6prefixes[i],
					}
				}
			}

			// update the link endpoints
			require.NoError(t, bp.UpdateSubinterfaces(ctx, subinterfaceConfigs))

			time.Sleep(time.Second)

			// fetch the result
			links, err = bp.GetAllSubinterfaceLinks(ctx)
			require.NoError(t, err)

			// validate tha the fetched result matches the values we sent
			var totalEndpoints int
			for _, link := range links {
				for _, ep := range link.Endpoints {
					totalEndpoints++

					expectedSubinterface, ok := subinterfaceConfigs[ep.SubinterfaceId]
					if !ok {
						t.Fatalf("endpoint with subinterface ID %q was not expected", ep.SubinterfaceId)
					}

					require.EqualValues(t, expectedSubinterface.Ipv4AddrType, ep.Subinterface.Ipv4AddrType)
					require.EqualValues(t, expectedSubinterface.Ipv6AddrType, ep.Subinterface.Ipv6AddrType)
					require.EqualValues(t, expectedSubinterface.Ipv4Addr.String(), ep.Subinterface.Ipv4Addr.String())
					require.EqualValues(t, expectedSubinterface.Ipv6Addr.String(), ep.Subinterface.Ipv6Addr.String())

					si, err := bp.GetSubinterface(ctx, ep.SubinterfaceId)
					require.NoError(t, err)

					require.EqualValues(t, ep.Subinterface.Ipv4AddrType, si.Ipv4AddrType)
					require.EqualValues(t, ep.Subinterface.Ipv4Addr, si.Ipv4Addr)
					require.EqualValues(t, ep.Subinterface.Ipv6AddrType, si.Ipv6AddrType)
					require.EqualValues(t, ep.Subinterface.Ipv6Addr, si.Ipv6Addr)
				}
			}

			// make sure we didn't miss anything or get any extras
			require.EqualValues(t, len(subinterfaceConfigs), totalEndpoints)

			// clear the addresses for those links (two endpoints per link)
			for _, link := range links {
				// prep an API payload for each end of the link
				for _, endpoint := range link.Endpoints {
					subinterfaceConfigs[endpoint.SubinterfaceId] = TwoStageL3ClosSubinterface{
						Ipv4AddrType: enum.InterfaceNumberingIpv4TypeNone,
						Ipv6AddrType: enum.InterfaceNumberingIpv6TypeNone,
					}
				}
			}

			// update the link endpoints
			require.NoError(t, bp.UpdateSubinterfaces(ctx, subinterfaceConfigs))

			// fetch the result
			links, err = bp.GetAllSubinterfaceLinks(ctx)
			require.NoError(t, err)

			// validate tha the fetched result matches the values we sent
			totalEndpoints = 0
			for _, link := range links {
				for _, ep := range link.Endpoints {
					totalEndpoints++

					expectedSubinterface, ok := subinterfaceConfigs[ep.SubinterfaceId]
					if !ok {
						t.Fatalf("endpoint with subinterface ID %q was not expected", ep.SubinterfaceId)
					}

					require.EqualValues(t, expectedSubinterface.Ipv4AddrType, ep.Subinterface.Ipv4AddrType)
					require.EqualValues(t, expectedSubinterface.Ipv6AddrType, ep.Subinterface.Ipv6AddrType)
					require.EqualValues(t, expectedSubinterface.Ipv4Addr.String(), ep.Subinterface.Ipv4Addr.String())
					require.EqualValues(t, expectedSubinterface.Ipv6Addr.String(), ep.Subinterface.Ipv6Addr.String())

					si, err := bp.GetSubinterface(ctx, ep.SubinterfaceId)
					require.NoError(t, err)

					require.EqualValues(t, ep.Subinterface.Ipv4AddrType, si.Ipv4AddrType)
					require.EqualValues(t, ep.Subinterface.Ipv4Addr, si.Ipv4Addr)
					require.EqualValues(t, ep.Subinterface.Ipv6AddrType, si.Ipv6AddrType)
					require.EqualValues(t, ep.Subinterface.Ipv6Addr, si.Ipv6Addr)
				}
			}

			// make sure we didn't miss anything or get any extras
			require.EqualValues(t, len(subinterfaceConfigs), totalEndpoints)

			// test getting links by ID
			links, err = bp.GetAllSubinterfaceLinks(ctx)
			require.NoError(t, err)

			for _, link := range links {
				t.Run(fmt.Sprintf("link_id_%s", link.Id), func(t *testing.T) {
					t.Parallel()

					l, err := bp.GetSubinterfaceLink(ctx, link.Id)
					require.NoError(t, err)
					compareLinks(t, &link, l)
				})
			}
		})
	}
}
