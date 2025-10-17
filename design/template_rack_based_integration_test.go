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
	"github.com/Juniper/apstra-go-sdk/policy"
	"github.com/Juniper/apstra-go-sdk/speed"
	"github.com/stretchr/testify/require"
)

var testTemplatesRackBased = map[string]design.TemplateRackBased{
	"L2_Virtual_EVPN": {
		Label: testutils.RandString(6, "hex"),
		Racks: []design.RackTypeWithCount{
			{
				Count: 4,
				RackType: design.RackType{
					Label:                    "L2 Virtual",
					FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
					LeafSwitches: []design.LeafSwitch{
						{
							Label:             "leaf",
							LinkPerSpineCount: pointer.To(1),
							LinkPerSpineSpeed: pointer.To(speed.Speed("10G")),
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-7x10-Leaf",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 7},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 2,
												Speed: "10G",
												Roles: []enum.PortRole{enum.PortRoleLeaf, enum.PortRoleSpine},
											},
											{
												Count: 2,
												Speed: "10G",
												Roles: []enum.PortRole{enum.PortRolePeer},
											},
											{
												Count: 2,
												Speed: "10G",
												Roles: []enum.PortRole{enum.PortRoleGeneric, enum.PortRoleAccess},
											},
											{
												Count: 1,
												Speed: "10G",
												Roles: []enum.PortRole{enum.PortRoleGeneric},
											},
										},
									},
								},
							},
						},
					},
					GenericSystems: []design.GenericSystem{
						{
							Count:     2,
							Label:     "generic",
							ASNDomain: &enum.FeatureSwitchDisabled,
							Links: []design.RackTypeLink{
								{
									Label:              "link",
									TargetSwitchLabel:  "leaf",
									LinkPerSwitchCount: 1,
									Speed:              "10G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									LAGMode:            enum.LAGModeNone,
								},
							},
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-1x10-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 1},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 1,
												Speed: "10G",
												Roles: []enum.PortRole{enum.PortRoleLeaf, enum.PortRoleAccess},
											},
										},
									},
								},
							},
							Loopback:        &enum.FeatureSwitchDisabled,
							ManagementLevel: enum.SystemManagementLevelUnmanaged,
						},
					},
				},
			},
		},
		AsnAllocationPolicy: &design.AsnAllocationPolicy{SpineAsnScheme: enum.AsnAllocationSchemeDistinct},
		Capability:          &enum.TemplateCapabilityBlueprint,
		DHCPServiceIntent:   policy.DHCPServiceIntent{Active: true},
		Spine: design.Spine{
			Count: 2,
			LogicalDevice: design.LogicalDevice{
				Label: "AOS-7x10-Spine",
				Panels: []design.LogicalDevicePanel{
					{
						PanelLayout: design.LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 7},
						PortGroups: []design.LogicalDevicePanelPortGroup{
							{
								Count: 5,
								Speed: "10G",
								Roles: []enum.PortRole{enum.PortRoleSuperspine, enum.PortRoleLeaf},
							},
							{
								Count: 2,
								Speed: "10G",
								Roles: []enum.PortRole{enum.PortRoleGeneric},
							},
						},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
					},
				},
			},
		},
		VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolEVPN},
	},
	"pod_mlag": {
		Label: testutils.RandString(6, "hex"),
		Racks: []design.RackTypeWithCount{
			{
				Count: 1,
				RackType: design.RackType{
					Label:                    "L2 MLAG 1x access",
					FabricConnectivityDesign: enum.FabricConnectivityDesign{Value: "l3clos"},
					LeafSwitches: []design.LeafSwitch{
						{
							Label:             "leaf",
							LinkPerSpineCount: pointer.To(1),
							LinkPerSpineSpeed: pointer.To(speed.Speed("10G")),
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-7x10-Leaf",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 7},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexing{Value: "T-B, L-R"},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRole{Value: "leaf"}, enum.PortRole{Value: "spine"}}},
											{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRole{Value: "peer"}}},
											{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRole{Value: "generic"}, enum.PortRole{Value: "access"}}},
											{Count: 1, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRole{Value: "generic"}}},
										},
									},
								},
							},
							RedundancyProtocol: enum.LeafRedundancyProtocol{Value: "mlag"},
							Tags:               []design.Tag{},
							MLAGInfo: &design.RackTypeLeafSwitchMLAGInfo{
								LeafLeafLinkCount: 2,
								LeafLeafLinkSpeed: "10G",
								MLAGVLAN:          2999,
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
									LinkPerSwitchCount: 2,
									Speed:              "10G",
									AttachmentType:     enum.LinkAttachmentType{Value: "dualAttached"},
									LAGMode:            enum.LAGMode{Value: "lacp_active"},
								},
							},
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-8x10-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 4},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexing{Value: "T-B, L-R"},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{Count: 8, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRole{Value: "leaf"}, enum.PortRole{Value: "generic"}, enum.PortRole{Value: "peer"}, enum.PortRole{Value: "access"}}},
										},
									},
								},
							},
						},
					},
					GenericSystems: []design.GenericSystem{
						{
							ASNDomain: &enum.FeatureSwitchDisabled,
							Count:     2,
							Label:     "generic",
							Links: []design.RackTypeLink{
								{
									Label:              "link",
									TargetSwitchLabel:  "access",
									LinkPerSwitchCount: 1,
									Speed:              "10G",
									AttachmentType:     enum.LinkAttachmentType{Value: "singleAttached"},
								},
							},
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-2x10-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 2},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexing{Value: "T-B, L-R"},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRole{Value: "leaf"}, enum.PortRole{Value: "access"}}},
										},
									},
								},
							},
							Loopback:        &enum.FeatureSwitchDisabled,
							ManagementLevel: enum.SystemManagementLevel{Value: "unmanaged"},
						},
					},
				},
			},
		},
		AsnAllocationPolicy: &design.AsnAllocationPolicy{SpineAsnScheme: enum.AsnAllocationSchemeSingle},
		Capability:          &enum.TemplateCapabilityPod,
		DHCPServiceIntent:   policy.DHCPServiceIntent{Active: true},
		Spine: design.Spine{
			Count:                  2,
			LinkPerSuperspineCount: 1,
			LinkPerSuperspineSpeed: "10G",
			LogicalDevice: design.LogicalDevice{
				Label: "AOS-32x10-Spine",
				Panels: []design.LogicalDevicePanel{
					{
						PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 16},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexing{Value: "T-B, L-R"},
						PortGroups: []design.LogicalDevicePanelPortGroup{
							{Count: 24, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRole{Value: "superspine"}, enum.PortRole{Value: "leaf"}}},
							{Count: 8, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRole{Value: "generic"}}},
						},
					},
				},
			},
		},
		VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
	},
}

func TestTemplateRackBased_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create design.TemplateRackBased
		update design.TemplateRackBased
	}

	testCases := map[string]testCase{
		"L2_Virtual_EVPN_to_pod_mlag": {
			create: testTemplatesRackBased["L2_Virtual_EVPN"],
			update: testTemplatesRackBased["pod_mlag"],
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
					var obj design.TemplateRackBased

					// create the object (by type)
					id, err = client.Client.CreateTemplateRackBased2(ctx, tCase.create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, 10, func(ctx context.Context) error {
						_ = client.Client.DeleteTemplate2(ctx, id)
						return nil
					})

					// retrieve the object by ID then validate
					template, err := client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok := template.(*design.TemplateRackBased)
					require.True(t, ok)
					idPtr := objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRackBased(t, tCase.create, *objPtr)

					// retrieve the object by ID (by type) then validate
					obj, err = client.Client.GetTemplateRackBased2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRackBased(t, tCase.create, obj)

					// retrieve the object by label then validate
					template, err = client.Client.GetTemplateByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplateRackBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRackBased(t, tCase.create, *objPtr)

					// retrieve the object by label (by type) then validate
					obj, err = client.Client.GetTemplateRackBasedByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRackBased(t, tCase.create, obj)

					// retrieve the list of IDs (ours must be in there)
					ids, err := client.Client.ListTemplates2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of IDs (by type) (ours must be in there)
					ids, err = client.Client.ListTemplatesRackBased2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) then validate
					templates, err := client.Client.GetTemplates2(ctx)
					require.NoError(t, err)
					templatePtr := slice.MustFindByID(templates, id)
					require.NotNil(t, templatePtr)
					objPtr, ok = (*templatePtr).(*design.TemplateRackBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRackBased(t, tCase.create, *objPtr)

					// retrieve the list of objects (by type) (ours must be in there) then validate
					objs, err := client.Client.GetTemplatesRackBased2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRackBased(t, tCase.create, obj)

					// update the object then validate
					require.NoError(t, tCase.update.SetID(id))
					require.NotNil(t, tCase.update.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateTemplate2(ctx, &tCase.update)
					require.NoError(t, err)

					// retrieve the updated object by ID then validate
					template, err = client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplateRackBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRackBased(t, tCase.update, *objPtr)

					// retrieve the updated object by ID (by type) type then validate
					update, err := client.Client.GetTemplateRackBased2(ctx, id)
					require.NoError(t, err)
					idPtr = update.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRackBased(t, tCase.update, update)

					// restore the object (by type)
					require.NoError(t, tCase.create.SetID(id))
					require.NotNil(t, tCase.create.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateTemplateRackBased2(ctx, tCase.create)
					require.NoError(t, err)

					// retrieve the restored object by ID then validate
					template, err = client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplateRackBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRackBased(t, tCase.create, *objPtr)

					// retrieve the restored object by ID (by type) then validate
					obj, err = client.Client.GetTemplateRackBased2(ctx, id)
					require.NoError(t, err)
					idPtr = update.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRackBased(t, tCase.create, obj)

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
					_, err = client.Client.GetTemplateRackBased2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetTemplateByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label (by type)
					_, err = client.Client.GetTemplateRackBasedByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListTemplates2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of IDs (by type) (ours must *not* be in there)
					ids, err = client.Client.ListTemplatesRackBased2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					templates, err = client.Client.GetTemplates2(ctx)
					require.NoError(t, err)
					templatePtr = slice.MustFindByID(templates, id)
					require.Nil(t, templatePtr)

					// retrieve the list of objects (by type) (ours must *not* be in there)
					objs, err = client.Client.GetTemplatesRackBased2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateTemplate2(ctx, &tCase.update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// update the object (by type)
					err = client.Client.UpdateTemplateRackBased2(ctx, tCase.update)
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
