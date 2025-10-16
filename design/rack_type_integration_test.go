// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package design_test

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/design"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/Juniper/apstra-go-sdk/speed"
	"github.com/stretchr/testify/require"
)

var testRackTypes = map[string]design.RackType{
	"collapsed_1xleaf": {
		Label:                    testutils.RandString(6, "hex"),
		Description:              testutils.RandString(6, "hex"),
		FabricConnectivityDesign: enum.FabricConnectivityDesignL3Collapsed,
		LeafSwitches: []design.LeafSwitch{
			{
				Label:         testutils.RandString(6, "hex"),
				Tags:          []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
				LogicalDevice: testLogicalDevices["leaf_48x25_4x400"],
			},
		},
	},
	"leaf_esi_access_esi_servers": {
		Label:                    testutils.RandString(6, "hex"),
		Description:              testutils.RandString(6, "hex"),
		FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
		LeafSwitches: []design.LeafSwitch{
			{
				Label:              "leaf",
				LinkPerSpineCount:  pointer.To(1),
				LinkPerSpineSpeed:  pointer.To(speed.Speed("400G")),
				LogicalDevice:      testLogicalDevices["leaf_48x25_4x400"],
				RedundancyProtocol: enum.LeafRedundancyProtocolESI,
				Tags:               []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
			},
		},
		AccessSwitches: []design.AccessSwitch{
			{
				Count: 1,
				ESILAGInfo: pointer.To(design.RackTypeAccessSwitchESILAGInfo{
					LinkCount:        2,
					LinkSpeed:        "400G",
					PortChannelIdMax: 19,
					PortChannelIdMin: 10,
				}),
				Label:         "access",
				LogicalDevice: testLogicalDevices["leaf_48x25_4x400"],
				Tags:          []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
				Links: []design.RackTypeLink{
					{
						Label:              testutils.RandString(6, "hex"),
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 1,
						Speed:              "400G",
						AttachmentType:     enum.LinkAttachmentTypeDual,
						LAGMode:            enum.LAGModeActiveLACP,
						Tags:               []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
					},
				},
			},
		},
		GenericSystems: []design.GenericSystem{
			{
				AsnDomain: nil,
				Count:     1,
				Label:     testutils.RandString(6, "hex"),
				Links: []design.RackTypeLink{
					{
						Label:              testutils.RandString(6, "hex"),
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 1,
						Speed:              "25G",
						AttachmentType:     enum.LinkAttachmentTypeSingle,
						LAGMode:            enum.LAGModeNone,
						SwitchPeer:         enum.LinkSwitchPeerFirst,
						Tags:               []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
					},
				},
				LogicalDevice:    testLogicalDevices["generic_4x25"],
				Loopback:         nil,
				ManagementLevel:  enum.SystemManagementLevelUnmanaged,
				PortChannelIDMax: 19,
				PortChannelIDMin: 10,
				Tags:             []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
			},
			{
				AsnDomain: pointer.To(enum.FeatureSwitchEnabled),
				Count:     2,
				Label:     testutils.RandString(6, "hex"),
				Links: []design.RackTypeLink{
					{
						Label:              testutils.RandString(6, "hex"),
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 2,
						Speed:              "25G",
						AttachmentType:     enum.LinkAttachmentTypeSingle,
						LAGMode:            enum.LAGModePassiveLACP,
						SwitchPeer:         enum.LinkSwitchPeerSecond,
						Tags:               []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
					},
				},
				LogicalDevice:    testLogicalDevices["generic_4x25"],
				Loopback:         pointer.To(enum.FeatureSwitchEnabled),
				ManagementLevel:  enum.SystemManagementLevelUnmanaged,
				PortChannelIDMax: 29,
				PortChannelIDMin: 20,
				Tags:             []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
			},
			{
				AsnDomain: nil,
				Count:     1,
				Label:     testutils.RandString(6, "hex"),
				Links: []design.RackTypeLink{
					{
						Label:              testutils.RandString(6, "hex"),
						TargetSwitchLabel:  "access",
						LinkPerSwitchCount: 1,
						Speed:              "25G",
						AttachmentType:     enum.LinkAttachmentTypeSingle,
						LAGMode:            enum.LAGModeStatic,
						SwitchPeer:         enum.LinkSwitchPeerFirst,
						Tags:               []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
					},
				},
				LogicalDevice:    testLogicalDevices["generic_4x25"],
				Loopback:         nil,
				ManagementLevel:  enum.SystemManagementLevelUnmanaged,
				PortChannelIDMax: 39,
				PortChannelIDMin: 30,
				Tags:             []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
			},
			{
				AsnDomain: pointer.To(enum.FeatureSwitchEnabled),
				Count:     2,
				Label:     testutils.RandString(6, "hex"),
				Links: []design.RackTypeLink{
					{
						Label:              testutils.RandString(6, "hex"),
						TargetSwitchLabel:  "access",
						LinkPerSwitchCount: 2,
						Speed:              "25G",
						AttachmentType:     enum.LinkAttachmentTypeSingle,
						LAGMode:            enum.LAGModeActiveLACP,
						SwitchPeer:         enum.LinkSwitchPeerSecond,
						Tags:               []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
					},
				},
				LogicalDevice:    testLogicalDevices["generic_4x25"],
				Loopback:         pointer.To(enum.FeatureSwitchEnabled),
				ManagementLevel:  enum.SystemManagementLevelUnmanaged,
				PortChannelIDMax: 49,
				PortChannelIDMin: 40,
				Tags:             []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
			},
			{
				AsnDomain: pointer.To(enum.FeatureSwitchEnabled),
				Count:     2,
				Label:     testutils.RandString(6, "hex"),
				Links: []design.RackTypeLink{
					{
						Label:              testutils.RandString(6, "hex"),
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 2,
						Speed:              "25G",
						AttachmentType:     enum.LinkAttachmentTypeDual,
						LAGMode:            enum.LAGModeActiveLACP,
						Tags:               []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
					},
				},
				LogicalDevice:    testLogicalDevices["generic_4x25"],
				Loopback:         pointer.To(enum.FeatureSwitchEnabled),
				ManagementLevel:  enum.SystemManagementLevelUnmanaged,
				PortChannelIDMax: 59,
				PortChannelIDMin: 50,
				Tags:             []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
			},
			{
				AsnDomain: nil,
				Count:     1,
				Label:     testutils.RandString(6, "hex"),
				Links: []design.RackTypeLink{
					{
						Label:              testutils.RandString(6, "hex"),
						TargetSwitchLabel:  "access",
						LinkPerSwitchCount: 1,
						Speed:              "25G",
						AttachmentType:     enum.LinkAttachmentTypeDual,
						LAGMode:            enum.LAGModeStatic,
						Tags:               []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
					},
				},
				LogicalDevice:    testLogicalDevices["generic_4x25"],
				Loopback:         nil,
				ManagementLevel:  enum.SystemManagementLevelUnmanaged,
				PortChannelIDMax: 69,
				PortChannelIDMin: 60,
				Tags:             []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
			},
		},
	},
	"leaf_mlag_2xaccess": {
		Label:                    testutils.RandString(6, "hex"),
		Description:              testutils.RandString(6, "hex"),
		FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
		LeafSwitches: []design.LeafSwitch{
			{
				LinkPerSpineCount: pointer.To(2),
				LinkPerSpineSpeed: pointer.To(speed.Speed("400G")),
				Label:             "leaf",
				Tags: []design.Tag{
					{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")},
					{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")},
				},
				MLAGInfo: pointer.To(design.RackTypeLeafSwitchMLAGInfo{
					LeafLeafL3LinkCount:         1,
					LeafLeafL3LinkSpeed:         "400G",
					LeafLeafL3LinkPortChannelId: 3,
					LeafLeafLinkCount:           1,
					LeafLeafLinkSpeed:           "400G",
					LeafLeafLinkPortChannelId:   2,
					MLAGVLAN:                    100,
				}),
				RedundancyProtocol: enum.LeafRedundancyProtocolMLAG,
				LogicalDevice:      testLogicalDevices["leaf_48x25_4x400"],
			},
		},
		AccessSwitches: []design.AccessSwitch{
			{
				Count:         1,
				Label:         "dual-homed-access",
				LogicalDevice: testLogicalDevices["leaf_48x25_4x400"],
				Tags:          []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
				Links: []design.RackTypeLink{
					{
						Label:              testutils.RandString(6, "hex"),
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 1,
						Speed:              "400G",
						AttachmentType:     enum.LinkAttachmentTypeDual,
						LAGMode:            enum.LAGModeActiveLACP,
						Tags:               []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
					},
				},
			},
			{
				Count:         1,
				Label:         "lefty",
				LogicalDevice: testLogicalDevices["leaf_48x25_4x400"],
				Tags:          []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
				Links: []design.RackTypeLink{
					{
						Label:              testutils.RandString(6, "hex"),
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 1,
						Speed:              "25G",
						AttachmentType:     enum.LinkAttachmentTypeSingle,
						SwitchPeer:         enum.LinkSwitchPeerFirst,
						LAGMode:            enum.LAGModeActiveLACP,
						Tags:               []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
					},
				},
			},
			{
				Count:         1,
				Label:         "righty",
				LogicalDevice: testLogicalDevices["leaf_48x25_4x400"],
				Tags:          []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
				Links: []design.RackTypeLink{
					{
						Label:              testutils.RandString(6, "hex"),
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 1,
						Speed:              "25G",
						AttachmentType:     enum.LinkAttachmentTypeSingle,
						SwitchPeer:         enum.LinkSwitchPeerSecond,
						LAGMode:            enum.LAGModeActiveLACP,
						Tags:               []design.Tag{{Label: testutils.RandString(6, "hex"), Description: testutils.RandString(6, "hex")}},
					},
				},
			},
		},
	},
}

func TestRackType_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create design.RackType
		update design.RackType
	}

	testCases := map[string]testCase{
		"collapsed_vs_esi": {
			create: testRackTypes["collapsed_1xleaf"],
			update: testRackTypes["leaf_esi_access_esi_servers"],
		},
		"mlag_vs_esi": {
			create: testRackTypes["leaf_mlag_2xaccess"],
			update: testRackTypes["leaf_esi_access_esi_servers"],
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			ctx := testutils.ContextWithTestID(ctx, t)
			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					var id string
					var err error
					var obj design.RackType

					// create the object
					id, err = client.Client.CreateRackType2(ctx, tCase.create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, 10, func(ctx context.Context) error {
						_ = client.Client.DeleteRackType2(ctx, id)
						return nil
					})

					// retrieve the object by ID and validate
					obj, err = client.Client.GetRackType2(ctx, id)
					require.NoError(t, err)
					idPtr := obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.RackType(t, tCase.create, obj)

					// retrieve the object by label and validate
					obj, err = client.Client.GetRackTypeByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.RackType(t, tCase.create, obj)

					// retrieve the list of IDs (ours must be in there)
					ids, err := client.Client.ListRackTypes2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := client.Client.GetRackTypes2(ctx)
					require.NoError(t, err)
					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.RackType(t, tCase.create, obj)

					// update the object and validate
					require.NoError(t, tCase.update.SetID(id))
					require.NotNil(t, tCase.update.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateRackType2(ctx, tCase.update)
					require.NoError(t, err)

					// retrieve the updated object by ID and validate
					update, err := client.Client.GetRackType2(ctx, id)
					require.NoError(t, err)
					idPtr = update.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.RackType(t, tCase.update, update)

					// restore the object to the original state
					require.NoError(t, tCase.create.SetID(id))
					require.NotNil(t, tCase.create.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateRackType2(ctx, tCase.create)
					require.NoError(t, err)

					// retrieve the object by ID and validate
					obj, err = client.Client.GetRackType2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.RackType(t, tCase.create, obj)

					// delete the object
					err = client.Client.DeleteRackType2(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = client.Client.GetRackType2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetRackTypeByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListRackTypes2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = client.Client.GetRackTypes2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateRackType2(ctx, tCase.update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = client.Client.DeleteRackType2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}
