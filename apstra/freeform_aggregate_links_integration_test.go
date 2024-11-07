// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDFFAggLink(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t testing.TB, expected, actual FreeformAggregateLinkData) {
		t.Helper()

		require.NotNil(t, actual, "actual is nil")
		require.NotNil(t, expected, "expected is nil")
		require.Equal(t, expected.Label, actual.Label, "Label")
		require.Equal(t, expected.Endpoints, actual.Endpoints, "Endpoints")
		require.Equal(t, expected.MemberLinkIds, actual.MemberLinkIds, "MemberLinkIds")
	}

	type testCase struct {
		steps []FreeformAggregateLinkData
	}

	for _, client := range clients {
		ffc, intSysIds, _ := testFFBlueprintB(ctx, t, client.client, 4, 0)
		// todo build links now
		link1 := []FreeformLinkRequest{
			{
				Label: randString(2, "hex"),
				Endpoints: [2]FreeformEndpoint{
					{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
						IfName:           toPtr("ge-0/0/0"),
						TransformationId: toPtr(1),
					}}},
					{SystemId: intSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
						IfName:           toPtr("ge-0/0/0"),
						TransformationId: toPtr(1),
					}}},
				},
			},
		}
		link2 := []FreeformLinkRequest{
			{
				Label: randString(2, "hex"),
				Endpoints: [2]FreeformEndpoint{
					{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
						IfName:           toPtr("ge-0/0/1"),
						TransformationId: toPtr(1),
					}}},
					{SystemId: intSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
						IfName:           toPtr("ge-0/0/1"),
						TransformationId: toPtr(1),
					}}},
				},
			},
		}

		var link1id ObjectId
		var link2id ObjectId
		link1id, err = ffc.CreateLink(ctx, &link1[0])
		require.NoError(t, err)
		link2id, err = ffc.CreateLink(ctx, &link2[0])
		require.NoError(t, err)

		testCases := map[string]testCase{
			"start_with_minimal_config": {
				steps: []FreeformAggregateLinkData{
					{
						Label: randString(2, "hex"),
						Endpoints: [2][]FreeformAggregateLinkMemberEndpoint{
							{
								{
									SystemId:      intSysIds[0],
									PortChannelId: 1,
									LagMode:       RackLinkLagModeActive,
								},
							},
							{
								{
									SystemId:      intSysIds[1],
									PortChannelId: 1,
									LagMode:       RackLinkLagModeActive,
								},
							},
						},
						MemberLinkIds: []ObjectId{link1id, link2id},
					},
					{
						Label: randString(2, "hex"),
						Endpoints: [2][]FreeformAggregateLinkMemberEndpoint{
							{
								{
									SystemId:      intSysIds[0],
									PortChannelId: 1,
									LagMode:       RackLinkLagModeActive,
									Ipv4Address:   netip.MustParsePrefix("192.168.2.1/31"),
								},
							},
							{
								{
									SystemId:      intSysIds[1],
									PortChannelId: 1,
									LagMode:       RackLinkLagModeActive,
									Ipv4Address:   netip.MustParsePrefix("192.168.2.2/31"),
								},
							},
						},
						MemberLinkIds: []ObjectId{link1id, link2id},
					},
				},
			},
			"start_with_second_config": {
				steps: []FreeformAggregateLinkData{
					{
						Label: randString(2, "hex"),
						Endpoints: [2][]FreeformAggregateLinkMemberEndpoint{
							{
								{
									SystemId:      intSysIds[2],
									PortChannelId: 1,
									LagMode:       RackLinkLagModeActive,
								},
							},
							{
								{
									SystemId:      intSysIds[3],
									PortChannelId: 1,
									LagMode:       RackLinkLagModeActive,
								},
							},
						},
						MemberLinkIds: []ObjectId{link1id, link2id},
					},
					{
						Label: randString(2, "hex"),
						Endpoints: [2][]FreeformAggregateLinkMemberEndpoint{
							{
								{
									SystemId:      intSysIds[2],
									PortChannelId: 1,
									LagMode:       RackLinkLagModeActive,
									Ipv4Address:   netip.MustParsePrefix("192.168.2.0/31"),
								},
							},
							{
								{
									SystemId:      intSysIds[3],
									PortChannelId: 1,
									LagMode:       RackLinkLagModeActive,
									Ipv4Address:   netip.MustParsePrefix("192.168.2.1/31"),
								},
							},
						},
						MemberLinkIds: []ObjectId{link1id, link2id},
					},
				},
			},
		}

		copyAcutalIntIds := func(t testing.TB, src, dst []FreeformAggregateLinkMemberEndpoint) {
			t.Helper()

			srcMap := make(map[ObjectId]FreeformAggregateLinkMemberEndpoint, len(src))
			for _, srcEp := range src {
				srcMap[srcEp.SystemId] = srcEp
			}

			for i, dstEp := range dst {
				srcEp, ok := srcMap[dstEp.SystemId]
				if !ok {
					log.Fatalf("src does not contain endpoint %s", dstEp.SystemId)
				}
				// copy actual intfId from source to destination here
				dst[i].AggIntfId = srcEp.AggIntfId
			}
		}

		for tName, tCase := range testCases {
			tName, tCase := tName, tCase

			t.Run(tName, func(t *testing.T) {
				// t.Parallel()

				// create the link
				id, err := ffc.CreateAggregateLink(ctx, &tCase.steps[0])
				require.NoError(t, err)
				require.NotEmpty(t, id)

				// read the link
				link, err := ffc.GetAggregateLink(ctx, id)
				require.NoError(t, err)
				require.Equal(t, id, link.Id, "link Id After Create")
				//				compare(t, *link.Data, tCase.steps[0])
				// todo copy agginterfaceids into tcase.steps[0]
				copyAcutalIntIds(t, link.Data.Endpoints[0], tCase.steps[0].Endpoints[0])
				copyAcutalIntIds(t, link.Data.Endpoints[1], tCase.steps[0].Endpoints[1])

				// todo compare tcase.steps[0], link.Data.Endpoints[0]
				require.Equal(t, tCase.steps[0].Endpoints[0], link.Data.Endpoints[0])
				require.Equal(t, tCase.steps[0].Endpoints[1], link.Data.Endpoints[1])

				// update the link once for each "step", including the first step (values used at creation)
				for i, step := range tCase.steps {
					t.Run(fmt.Sprintf("update_step_%d", i), func(t *testing.T) {
						// copy the aggInterfaceIds into the update request
						copyAcutalIntIds(t, link.Data.Endpoints[0], step.Endpoints[0])
						copyAcutalIntIds(t, link.Data.Endpoints[1], step.Endpoints[1])
						// clear the aggregate interface ID

						// update the link
						require.NoError(t, ffc.UpdateAggregateLink(ctx, id, &step))

						// read the link
						link, err = ffc.GetAggregateLink(ctx, id)
						require.NoError(t, err)
						require.Equal(t, id, link.Id, fmt.Sprintf("linkId after update iteration %d", i))
						compare(t, step, *link.Data)
					})
				}

				// delete the link
				err = ffc.DeleteAggregateLink(ctx, id)
				require.NoError(t, err)

				var ace ClientErr

				// fetching a previously deleted link should fail
				_, err = ffc.GetAggregateLink(ctx, id)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ErrNotfound, ace.Type())

				// deleting a previously deleted link should fail
				err = ffc.DeleteAggregateLink(ctx, id)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ErrNotfound, ace.Type())
			})
		}
	}
}
