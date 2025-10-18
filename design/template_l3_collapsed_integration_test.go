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
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/design"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
	"github.com/Juniper/apstra-go-sdk/policy"
	"github.com/stretchr/testify/require"
)

var testTemplatesL3Collapsed = map[string]design.TemplateL3Collapsed{
	"L3_Collapsed_ACS": {
		Label: testutils.RandString(6, "hex"),
		Racks: []design.RackTypeWithCount{
			{
				RackType: design.RackType{
					Label:                    "Collapsed 1xleaf",
					FabricConnectivityDesign: enum.FabricConnectivityDesignL3Collapsed,
					LeafSwitches: []design.LeafSwitch{
						{
							Label: "leaf",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-7x10-Leaf",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 7},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleSpine}},
											{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRolePeer}},
											{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}},
											{Count: 1, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric}},
										},
									},
								},
							},
						},
					},
					AccessSwitches: []design.AccessSwitch{
						{
							Count: 1,
							Label: "access",
							Links: []design.RackTypeLink{
								{
									Label:              "leaf_link",
									TargetSwitchLabel:  "leaf",
									LinkPerSwitchCount: 1,
									Speed:              "10G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									LAGMode:            enum.LAGModeActiveLACP,
								},
							},
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-8x10-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 4},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{Count: 8, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess}},
										},
									},
								},
							},
						},
					},
					GenericSystems: []design.GenericSystem{
						{
							Count: 2,
							Label: "generic",
							Links: []design.RackTypeLink{
								{
									Label:              "link",
									TargetSwitchLabel:  "access",
									LinkPerSwitchCount: 1,
									Speed:              "10G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
								},
							},
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-1x10-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 1},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{Count: 1, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess}},
										},
									},
								},
							},
							ManagementLevel: enum.SystemManagementLevelUnmanaged,
						},
					},
				},
				Count: 1,
			},
		},
		MeshLinkCount:        1,
		MeshLinkSpeed:        "10G",
		DHCPServiceIntent:    policy.DHCPServiceIntent{Active: true},
		VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolEVPN},
	},
	"L3_Collapsed_ESI": {
		Label: testutils.RandString(6, "hex"),
		Racks: []design.RackTypeWithCount{
			{
				Count: 1,
				RackType: design.RackType{
					Label:                    "Collapsed 2xleafs",
					FabricConnectivityDesign: enum.FabricConnectivityDesignL3Collapsed,
					LeafSwitches: []design.LeafSwitch{
						{
							Label: "esi_leaf",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-7x10-Leaf",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 7},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleSpine}},
											{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRolePeer}},
											{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}},
											{Count: 1, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric}},
										},
									},
								},
							},
							RedundancyProtocol: enum.LeafRedundancyProtocolESI,
						},
					},
					GenericSystems: []design.GenericSystem{
						{
							Count: 2,
							Label: "generic",
							Links: []design.RackTypeLink{
								{
									Label:              "link",
									TargetSwitchLabel:  "esi_leaf",
									LinkPerSwitchCount: 1,
									Speed:              "10G",
									AttachmentType:     enum.LinkAttachmentTypeDual,
									LAGMode:            enum.LAGModeActiveLACP,
								},
							},
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-2x10-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 2},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess}},
										},
									},
								},
							},
							ManagementLevel: enum.SystemManagementLevelUnmanaged,
						},
					},
				},
			},
		},
		MeshLinkCount:        1,
		MeshLinkSpeed:        "10G",
		DHCPServiceIntent:    policy.DHCPServiceIntent{Active: true},
		VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolEVPN},
	},
}

func TestTemplateL3Collapsed_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create design.TemplateL3Collapsed
		update design.TemplateL3Collapsed
	}

	testCases := map[string]testCase{
		"L3_Collapsed_ACS_to_L3_Collapsed_ESI": {
			create: testTemplatesL3Collapsed["L3_Collapsed_ACS"],
			update: testTemplatesL3Collapsed["L3_Collapsed_ESI"],
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			ctx := testutils.ContextWithTestID(ctx, t)
			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					require.NotEqual(t, tCase.create, zero.Of(tCase.create)) // make sure we didn't use a bogus map key
					require.NotEqual(t, tCase.update, zero.Of(tCase.update)) // make sure we didn't use a bogus map key

					var id string
					var err error
					var obj design.TemplateL3Collapsed

					// create the object (by type)
					id, err = client.Client.CreateTemplateL3Collapsed2(ctx, tCase.create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, 10, func(ctx context.Context) error {
						_ = client.Client.DeleteTemplate2(ctx, id)
						return nil
					})

					// retrieve the object by ID then validate
					template, err := client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok := template.(*design.TemplateL3Collapsed)
					require.True(t, ok)
					idPtr := objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateL3Collapsed(t, tCase.create, *objPtr)

					// retrieve the object by ID (by type) then validate
					obj, err = client.Client.GetTemplateL3Collapsed2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateL3Collapsed(t, tCase.create, obj)

					// retrieve the object by label then validate
					template, err = client.Client.GetTemplateByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplateL3Collapsed)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateL3Collapsed(t, tCase.create, *objPtr)

					// retrieve the object by label (by type) then validate
					obj, err = client.Client.GetTemplateL3CollapsedByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateL3Collapsed(t, tCase.create, obj)

					// retrieve the list of IDs (ours must be in there)
					ids, err := client.Client.ListTemplates2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of IDs (by type) (ours must be in there)
					ids, err = client.Client.ListTemplatesL3Collapsed2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) then validate
					templates, err := client.Client.GetTemplates2(ctx)
					require.NoError(t, err)
					templatePtr := slice.MustFindByID(templates, id)
					require.NotNil(t, templatePtr)
					objPtr, ok = (*templatePtr).(*design.TemplateL3Collapsed)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateL3Collapsed(t, tCase.create, *objPtr)

					// retrieve the list of objects (by type) (ours must be in there) then validate
					objs, err := client.Client.GetTemplatesL3Collapsed2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateL3Collapsed(t, tCase.create, obj)

					// update the object then validate
					require.NoError(t, tCase.update.SetID(id))
					require.NotNil(t, tCase.update.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateTemplate2(ctx, &tCase.update)
					require.NoError(t, err)

					// retrieve the updated object by ID then validate
					template, err = client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplateL3Collapsed)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateL3Collapsed(t, tCase.update, *objPtr)

					// retrieve the updated object by ID (by type) type then validate
					update, err := client.Client.GetTemplateL3Collapsed2(ctx, id)
					require.NoError(t, err)
					idPtr = update.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateL3Collapsed(t, tCase.update, update)

					// restore the object (by type)
					require.NoError(t, tCase.create.SetID(id))
					require.NotNil(t, tCase.create.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateTemplateL3Collapsed2(ctx, tCase.create)
					require.NoError(t, err)

					// retrieve the restored object by ID then validate
					template, err = client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplateL3Collapsed)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateL3Collapsed(t, tCase.create, *objPtr)

					// retrieve the restored object by ID (by type) then validate
					obj, err = client.Client.GetTemplateL3Collapsed2(ctx, id)
					require.NoError(t, err)
					idPtr = update.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateL3Collapsed(t, tCase.create, obj)

					// delete the object
					err = client.Client.DeleteTemplate2(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = client.Client.GetTemplate2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by ID (by type)
					_, err = client.Client.GetTemplateL3Collapsed2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetTemplateByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label (by type)
					_, err = client.Client.GetTemplateL3CollapsedByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListTemplates2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of IDs (by type) (ours must *not* be in there)
					ids, err = client.Client.ListTemplatesL3Collapsed2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					templates, err = client.Client.GetTemplates2(ctx)
					require.NoError(t, err)
					templatePtr = slice.MustFindByID(templates, id)
					require.Nil(t, templatePtr)

					// retrieve the list of objects (by type) (ours must *not* be in there)
					objs, err = client.Client.GetTemplatesL3Collapsed2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateTemplate2(ctx, &tCase.update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// update the object (by type)
					err = client.Client.UpdateTemplateL3Collapsed2(ctx, tCase.update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = client.Client.DeleteTemplate2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}
