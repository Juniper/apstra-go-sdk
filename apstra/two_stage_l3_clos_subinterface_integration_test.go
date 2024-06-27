//go:build integration
// +build integration

package apstra

import (
	"context"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net"
	"testing"
)

func TestUpdateTwoStageL3ClosSubinterface(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for _, client := range clients {
		// create a blueprint and enable IPv6
		bp := testBlueprintC(ctx, t, client.client)
		settings, err := bp.GetFabricSettings(ctx)
		require.NoError(t, err)
		settings.Ipv6Enabled = toPtr(true)
		require.NoError(t, bp.SetFabricSettings(ctx, settings))

		// create a security zone within the blueprint
		szName := randString(6, "hex")
		szId, err := bp.CreateSecurityZone(ctx, &SecurityZoneData{
			Label:   szName,
			SzType:  SecurityZoneTypeEVPN,
			VrfName: szName,
			VniId:   toPtr(rand.Intn(1000) + 10000),
		})
		require.NoError(t, err)

		// create a simple "ip link" connectivity template which uses the security zone
		ct := ConnectivityTemplate{
			Label: randString(6, "hex"),
			Subpolicies: []*ConnectivityTemplatePrimitive{{
				Label: "IP Link",
				Attributes: &ConnectivityTemplatePrimitiveAttributesAttachLogicalLink{
					Label:              "",
					SecurityZone:       &szId,
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
		var subinterfaceApiPayload []TwoStageL3ClosSubinterface                              // slice for use as Update() payload -- do not attempt to size!
		expectedSubinterfaces := make(map[ObjectId]TwoStageL3ClosSubinterface, len(links)*2) // map for easy lookup later
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
				expectedSubinterfaces[endpoint.Subinterface.Id] = TwoStageL3ClosSubinterface{
					Id:           endpoint.Subinterface.Id,
					Ipv4AddrType: toPtr(InterfaceNumberingIpv4TypeNumbered),
					Ipv6AddrType: toPtr(InterfaceNumberingIpv6TypeNumbered),
					Ipv4Addr:     &v4prefixes[i],
					Ipv6Addr:     &v6prefixes[i],
				}
				subinterfaceApiPayload = append(subinterfaceApiPayload, expectedSubinterfaces[endpoint.Subinterface.Id])
			}
		}

		// update the link endpoints
		require.NoError(t, bp.UpdateSubinterfaces(ctx, subinterfaceApiPayload))

		// fetch the result
		links, err = bp.GetAllSubinterfaceLinks(ctx)
		require.NoError(t, err)

		// validate tha the fetched result matches the values we sent
		var totalEndpoints int
		for _, link := range links {
			for _, ep := range link.Endpoints {
				totalEndpoints++

				expectedSubinterface, ok := expectedSubinterfaces[ep.Subinterface.Id]
				if !ok {
					t.Fatalf("endpoint with subinterface ID %q was not expected", ep.Subinterface.Id)
				}

				require.EqualValues(t, *expectedSubinterface.Ipv4AddrType, *ep.Subinterface.Ipv4AddrType)
				require.EqualValues(t, *expectedSubinterface.Ipv6AddrType, *ep.Subinterface.Ipv6AddrType)
				require.EqualValues(t, expectedSubinterface.Ipv4Addr.String(), ep.Subinterface.Ipv4Addr.String())
				require.EqualValues(t, expectedSubinterface.Ipv6Addr.String(), ep.Subinterface.Ipv6Addr.String())
			}
		}

		// make sure we didn't miss anything or get any extras
		require.EqualValues(t, len(subinterfaceApiPayload), totalEndpoints)
	}
}
