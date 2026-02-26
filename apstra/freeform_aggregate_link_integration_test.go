// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package apstra_test

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	comparefreeform "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/freeform"
	fftestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/freeform_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestCRUDFreeformAggregateLink(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create apstra.FreeformAggregateLink
		update apstra.FreeformAggregateLink
	}

	// create2i1e creates two internal systems and 1 external system:
	//         +-----------+
	//         |  external |
	//         +-----------+
	//            /     \
	//           /       \
	//       link1     link2
	//         /           \
	// +------------+   +------------+
	// | internal1  |   | internal2  |
	// +------------+   +------------+
	//
	// five IDs are returned: external, internal1, internal2, link1, link2
	create2i1e := func(t *testing.T, ctx context.Context, bp apstra.FreeformClient) (string, string, string, string, string) {
		t.Helper()

		dpID, err := bp.ImportDeviceProfile(ctx, "Juniper_vEX")
		require.NoError(t, err)

		eSysID, err := bp.CreateSystem(ctx, &apstra.FreeformSystemData{
			Type:  apstra.SystemTypeExternal,
			Label: testutils.RandString(6, "hex"),
		})
		require.NoError(t, err)

		iSysID1, err := bp.CreateSystem(ctx, &apstra.FreeformSystemData{
			Type:            apstra.SystemTypeInternal,
			Label:           testutils.RandString(6, "hex"),
			DeviceProfileId: &dpID,
		})
		require.NoError(t, err)

		link1ID, err := bp.CreateLink(ctx, &apstra.FreeformLinkRequest{
			Label: testutils.RandString(6, "hex"),
			Endpoints: [2]apstra.FreeformEthernetEndpoint{
				{SystemId: eSysID},
				{
					SystemId: iSysID1,
					Interface: apstra.FreeformInterface{
						Data: &apstra.FreeformInterfaceData{
							IfName:           pointer.To("ge-0/0/0"),
							TransformationId: pointer.To(1),
						},
					},
				},
			},
		})
		require.NoError(t, err)

		iSysID2, err := bp.CreateSystem(ctx, &apstra.FreeformSystemData{
			Type:            apstra.SystemTypeInternal,
			Label:           testutils.RandString(6, "hex"),
			DeviceProfileId: &dpID,
		})
		require.NoError(t, err)

		link2ID, err := bp.CreateLink(ctx, &apstra.FreeformLinkRequest{
			Label: testutils.RandString(6, "hex"),
			Endpoints: [2]apstra.FreeformEthernetEndpoint{
				{SystemId: eSysID},
				{
					SystemId: iSysID2,
					Interface: apstra.FreeformInterface{
						Data: &apstra.FreeformInterfaceData{
							IfName:           pointer.To("ge-0/0/0"),
							TransformationId: pointer.To(1),
						},
					},
				},
			},
		})
		require.NoError(t, err)

		return string(eSysID), string(iSysID1), string(iSysID2), string(link1ID), string(link2ID)
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(context.Background(), t)

			bp := fftestobj.TestBlueprintA(t, ctx, client.Client)
			require.NotNil(t, bp)

			e1, i1, i2, l1, l2 := create2i1e(t, ctx, *bp)

			testCases := map[string]testCase{
				"with_labels": {
					create: apstra.FreeformAggregateLink{
						Label:         pointer.To(testutils.RandString(6, "hex")),
						MemberLinkIds: []string{l1, l2},
						EndpointGroups: [2]apstra.FreeformAggregateLinkEndpointGroup{
							{
								Label: pointer.To(testutils.RandString(6, "hex")),
								Tags:  []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      e1,
										IfName:        "bond0",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 11,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModePassiveLACP,
									},
								},
							},
							{
								Label: pointer.To(testutils.RandString(6, "hex")),
								Tags:  []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      i1,
										IfName:        "ae21",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 221,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeActiveLACP,
									},
									{
										SystemID:      i2,
										IfName:        "ae22",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 222,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeActiveLACP,
									},
								},
							},
						},
						Tags: []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
					},
					update: apstra.FreeformAggregateLink{
						Label:         pointer.To(testutils.RandString(6, "hex")),
						MemberLinkIds: []string{l1, l2},
						EndpointGroups: [2]apstra.FreeformAggregateLinkEndpointGroup{
							{
								Label: pointer.To(testutils.RandString(6, "hex")),
								Tags:  []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      e1,
										IfName:        "bond1",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 12,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeStatic,
									},
								},
							},
							{
								Label: pointer.To(testutils.RandString(6, "hex")),
								Tags:  []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      i1,
										IfName:        "ae31",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 231,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeStatic,
									},
									{
										SystemID:      i2,
										IfName:        "ae32",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 232,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeStatic,
									},
								},
							},
						},
						Tags: []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
					},
				},
				"labels_are_nil": {
					create: apstra.FreeformAggregateLink{
						MemberLinkIds: []string{l1, l2},
						EndpointGroups: [2]apstra.FreeformAggregateLinkEndpointGroup{
							{
								Tags: []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      e1,
										IfName:        "bond0",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 11,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModePassiveLACP,
									},
								},
							},
							{
								Tags: []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      i1,
										IfName:        "ae21",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 221,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeActiveLACP,
									},
									{
										SystemID:      i2,
										IfName:        "ae22",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 222,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeActiveLACP,
									},
								},
							},
						},
						Tags: []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
					},
					update: apstra.FreeformAggregateLink{
						MemberLinkIds: []string{l1, l2},
						EndpointGroups: [2]apstra.FreeformAggregateLinkEndpointGroup{
							{
								Tags: []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      e1,
										IfName:        "bond1",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 12,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeStatic,
									},
								},
							},
							{
								Tags: []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      i1,
										IfName:        "ae31",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 231,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeStatic,
									},
									{
										SystemID:      i2,
										IfName:        "ae32",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 232,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeStatic,
									},
								},
							},
						},
						Tags: []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
					},
				},
				"clear_labels": {
					create: apstra.FreeformAggregateLink{
						Label:         pointer.To(testutils.RandString(6, "hex")),
						MemberLinkIds: []string{l1, l2},
						EndpointGroups: [2]apstra.FreeformAggregateLinkEndpointGroup{
							{
								Label: pointer.To(testutils.RandString(6, "hex")),
								Tags:  []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      e1,
										IfName:        "bond0",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 11,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModePassiveLACP,
									},
								},
							},
							{
								Label: pointer.To(testutils.RandString(6, "hex")),
								Tags:  []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      i1,
										IfName:        "ae21",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 221,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeActiveLACP,
									},
									{
										SystemID:      i2,
										IfName:        "ae22",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 222,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeActiveLACP,
									},
								},
							},
						},
						Tags: []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
					},
					update: apstra.FreeformAggregateLink{
						Label:         pointer.To(""),
						MemberLinkIds: []string{l1, l2},
						EndpointGroups: [2]apstra.FreeformAggregateLinkEndpointGroup{
							{
								Label: pointer.To(""),
								Tags:  []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      e1,
										IfName:        "bond1",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 12,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeStatic,
									},
								},
							},
							{
								Label: pointer.To(""),
								Tags:  []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
								Endpoints: []apstra.FreeformAggregateLinkEndpoint{
									{
										SystemID:      i1,
										IfName:        "ae31",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 231,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeStatic,
									},
									{
										SystemID:      i2,
										IfName:        "ae32",
										IPv4Addr:      testutils.RandomHostIP(t, "192.0.2.0/24"),
										IPv6Addr:      testutils.RandomHostIP(t, "3fff::/64"),
										PortChannelID: 232,
										Tags:          []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
										LAGMode:       enum.LAGModeStatic,
									},
								},
							},
						},
						Tags: []string{testutils.RandString(6, "hex"), testutils.RandString(6, "hex")},
					},
				},
			}

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					ctx := testutils.ContextWithTestID(context.Background(), t)

					create, update := tCase.create, tCase.update // because we modify these values below

					// create the object
					id, err := bp.CreateAggregateLink(ctx, tCase.create)
					require.NoError(t, err)

					// retrieve the object by ID and validate
					obj, err := bp.GetAggregateLink(ctx, id)
					require.NoError(t, err)
					idPtr := obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparefreeform.AggregateLink(t, create, obj)

					// retrieve the object by label and validate
					if create.Label != nil {
						obj, err = bp.GetAggregateLinkByLabel(ctx, *create.Label)
						require.NoError(t, err)
						idPtr = obj.ID()
						require.NotNil(t, idPtr)
						require.Equal(t, id, *idPtr)
						comparefreeform.AggregateLink(t, create, obj)
					}

					// retrieve the list of IDs - ours must be in there
					ids, err := bp.ListAggregateLinks(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := bp.GetAggregateLinks(ctx)
					require.NoError(t, err)
					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparefreeform.AggregateLink(t, create, obj)

					// update the object and validate
					update.SetID(id)
					require.NotNil(t, update.ID())
					require.Equal(t, id, *update.ID())
					err = bp.UpdateAggregateLink(ctx, update)
					require.NoError(t, err)

					// retrieve the updated object by ID and validate
					obj, err = bp.GetAggregateLink(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparefreeform.AggregateLink(t, update, obj)

					// delete the object
					err = bp.DeleteAggregateLink(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = bp.GetAggregateLink(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					if create.Label != nil {
						_, err = bp.GetAggregateLinkByLabel(ctx, *create.Label)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, apstra.ErrNotfound, ace.Type())
					}

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = bp.ListAggregateLinks(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = bp.GetAggregateLinks(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = bp.UpdateAggregateLink(ctx, update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = bp.DeleteAggregateLink(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}
