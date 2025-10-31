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
	"github.com/Juniper/apstra-go-sdk/internal/zero"
	"github.com/Juniper/apstra-go-sdk/policy"
	"github.com/Juniper/apstra-go-sdk/speed"
	"github.com/stretchr/testify/require"
)

var testTemplatesPodBased = map[string]design.TemplatePodBased{
	"L2_superspine_multi_plane": {
		Label: testutils.RandString(6, "hex"),
		Superspine: design.Superspine{
			PlaneCount:         4,
			SuperspinePerPlane: 4,
			Tags:               []design.Tag{},
			LogicalDevice: design.LogicalDevice{
				Label: "AOS-32x40-3",
				Panels: []design.LogicalDevicePanel{
					{
						PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 16},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
						PortGroups: []design.LogicalDevicePanelPortGroup{
							{Count: 32, Speed: "40G", Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRoleSpine}},
						},
					},
				},
			},
		},
		Pods: []design.PodWithCount{
			{
				Count: 2,
				Pod: design.TemplateRackBased{
					Label: "L2 Pod",
					Racks: []design.RackTypeWithCount{
						{
							Count: 1,
							RackType: design.RackType{
								Label:                    "L2 One Leaf",
								FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
								LeafSwitches: []design.RackTypeLeafSwitch{
									{
										Label:             "leaf",
										LinkPerSpineCount: pointer.To(1),
										LinkPerSpineSpeed: pointer.To(speed.Speed("40G")),
										LogicalDevice: design.LogicalDevice{
											Label: "AOS-64x10+16x40-2",
											Panels: []design.LogicalDevicePanel{
												{
													PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 32},
													PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
													PortGroups: []design.LogicalDevicePanelPortGroup{
														{Count: 64, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}},
													},
												},
												{
													PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
													PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
													PortGroups: []design.LogicalDevicePanelPortGroup{
														{Count: 16, Speed: "40G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleSpine}},
													},
												},
											},
										},
										Tags: []design.Tag{},
									},
								},
								GenericSystems: []design.RackTypeGenericSystem{
									{
										ASNDomain: &enum.FeatureSwitchDisabled,
										Count:     48,
										Label:     "generic",
										Links: []design.RackTypeLink{
											{
												Label:              "link",
												TargetSwitchLabel:  "leaf",
												LinkPerSwitchCount: 1,
												Speed:              "10G",
												AttachmentType:     enum.LinkAttachmentTypeSingle,
												Tags:               []design.Tag{},
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
										Loopback:        &enum.FeatureSwitchDisabled,
										ManagementLevel: enum.SystemManagementLevelUnmanaged,
										Tags:            []design.Tag{},
									},
								},
							},
						},
						{
							Count: 1,
							RackType: design.RackType{
								Label:                    "L2 Mlag Leaf",
								FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
								LeafSwitches: []design.RackTypeLeafSwitch{
									{
										Label:             "leaf",
										LinkPerSpineCount: pointer.To(1),
										LinkPerSpineSpeed: pointer.To(speed.Speed("40G")),
										LogicalDevice: design.LogicalDevice{
											Label: "AOS-64x10+16x40-2",
											Panels: []design.LogicalDevicePanel{
												{
													PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 32},
													PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
													PortGroups: []design.LogicalDevicePanelPortGroup{
														{Count: 64, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}},
													},
												},
												{
													PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
													PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
													PortGroups: []design.LogicalDevicePanelPortGroup{
														{Count: 16, Speed: "40G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleSpine}},
													},
												},
											},
										},
										RedundancyProtocol: enum.LeafRedundancyProtocolMLAG,
										Tags:               []design.Tag{},
										MLAGInfo: &design.RackTypeLeafSwitchMLAGInfo{
											LeafLeafLinkCount: 4,
											LeafLeafLinkSpeed: "40G",
											MLAGVLAN:          2999,
										},
									},
								},
								GenericSystems: []design.RackTypeGenericSystem{
									{
										ASNDomain: &enum.FeatureSwitchDisabled,
										Count:     24,
										Label:     "generic(1)",
										Links: []design.RackTypeLink{
											{
												Label:              "link",
												TargetSwitchLabel:  "leaf",
												LinkPerSwitchCount: 1,
												Speed:              "10G",
												AttachmentType:     enum.LinkAttachmentTypeSingle,
												SwitchPeer:         enum.LinkSwitchPeerFirst,
												Tags:               []design.Tag{},
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
										Loopback:        &enum.FeatureSwitchDisabled,
										ManagementLevel: enum.SystemManagementLevelUnmanaged,
										Tags:            []design.Tag{},
									},
									{
										ASNDomain: &enum.FeatureSwitchDisabled,
										Count:     24, Label: "generic(2)",
										Links: []design.RackTypeLink{
											{
												Label:              "link",
												TargetSwitchLabel:  "leaf",
												LinkPerSwitchCount: 1,
												Speed:              "10G",
												AttachmentType:     enum.LinkAttachmentTypeSingle,
												SwitchPeer:         enum.LinkSwitchPeerSecond,
												Tags:               []design.Tag{},
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
										Loopback:        &enum.FeatureSwitchDisabled,
										ManagementLevel: enum.SystemManagementLevelUnmanaged,
										Tags:            []design.Tag{},
									},
								},
							},
						},
					},
					ASNAllocationPolicy: &policy.ASNAllocation{SpineASNScheme: enum.ASNAllocationSchemeSingle},
					DHCPServiceIntent:   policy.DHCPServiceIntent{Active: true},
					Spine: design.Spine{
						Count:                  4,
						LinkPerSuperspineCount: 1,
						LinkPerSuperspineSpeed: "40G",
						LogicalDevice: design.LogicalDevice{
							Label: "AOS-32x40-3",
							Panels: []design.LogicalDevicePanel{
								{
									PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 16},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									PortGroups: []design.LogicalDevicePanelPortGroup{
										{Count: 32, Speed: "40G", Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRoleSpine}},
									},
								},
							},
						},
						Tags: []design.Tag{},
					},
					VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
				},
			},
		},
	},
	"L2_superspine_single_plane_with_acs": {
		Label: testutils.RandString(6, "hex"),
		Superspine: design.Superspine{
			PlaneCount:         1,
			SuperspinePerPlane: 4,
			LogicalDevice: design.LogicalDevice{
				Label: "AOS-32x10-3",
				Panels: []design.LogicalDevicePanel{
					{
						PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 16},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
						PortGroups: []design.LogicalDevicePanelPortGroup{
							{Count: 32, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRoleSpine}},
						},
					},
				},
			},
			Tags: []design.Tag{},
		},
		Pods: []design.PodWithCount{
			{
				Count: 1,
				Pod: design.TemplateRackBased{
					Label: "L2 Pod Single",
					Racks: []design.RackTypeWithCount{
						{
							Count: 1,
							RackType: design.RackType{
								Label:                    "L2 One Access",
								FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
								LeafSwitches: []design.RackTypeLeafSwitch{
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
														{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleSpine}},
														{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRolePeer}},
														{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}},
														{Count: 1, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric}},
													},
												},
											},
										},
										Tags: []design.Tag{},
									},
								},
								AccessSwitches: []design.RackTypeAccessSwitch{
									{
										Count: 1,
										Label: "access",
										Links: []design.RackTypeLink{
											{
												Label:              "leaf_link",
												TargetSwitchLabel:  "leaf",
												LinkPerSwitchCount: 2,
												Speed:              "10G",
												AttachmentType:     enum.LinkAttachmentTypeSingle,
												LAGMode:            enum.LAGModeActiveLACP,
												Tags:               []design.Tag{},
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
										Tags: []design.Tag{},
									},
								},
								GenericSystems: []design.RackTypeGenericSystem{
									{
										ASNDomain: &enum.FeatureSwitchDisabled,
										Count:     4,
										Label:     "generic",
										Links: []design.RackTypeLink{
											{
												Label:              "link",
												TargetSwitchLabel:  "access",
												LinkPerSwitchCount: 1,
												Speed:              "10G",
												AttachmentType:     enum.LinkAttachmentTypeSingle,
												Tags:               []design.Tag{},
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
										Loopback:        &enum.FeatureSwitchDisabled,
										ManagementLevel: enum.SystemManagementLevelUnmanaged,
										Tags:            []design.Tag{},
									},
								},
							},
						},
					},
					ASNAllocationPolicy: &policy.ASNAllocation{SpineASNScheme: enum.ASNAllocationSchemeSingle},
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
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									PortGroups: []design.LogicalDevicePanelPortGroup{
										{Count: 24, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleLeaf}},
										{Count: 8, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric}},
									},
								},
							},
						},
						Tags: []design.Tag{},
					},
					VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
				},
			},
			{
				Count: 1,
				Pod: design.TemplateRackBased{
					Label: "L2 Pod Mlag",
					Racks: []design.RackTypeWithCount{
						{
							Count: 1,
							RackType: design.RackType{
								Label:                    "L2 MLAG 1x access",
								FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
								LeafSwitches: []design.RackTypeLeafSwitch{
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
														{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleSpine}},
														{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRolePeer}},
														{Count: 2, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}},
														{Count: 1, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric}},
													},
												},
											},
										},
										RedundancyProtocol: enum.LeafRedundancyProtocolMLAG,
										Tags:               []design.Tag{},
										MLAGInfo: &design.RackTypeLeafSwitchMLAGInfo{
											LeafLeafLinkCount: 2,
											LeafLeafLinkSpeed: "10G",
											MLAGVLAN:          2999,
										},
									},
								},
								AccessSwitches: []design.RackTypeAccessSwitch{
									{
										Count: 1,
										Label: "access",
										Links: []design.RackTypeLink{
											{
												Label:              "leaf_link",
												TargetSwitchLabel:  "leaf",
												LinkPerSwitchCount: 2,
												Speed:              "10G",
												AttachmentType:     enum.LinkAttachmentTypeDual,
												LAGMode:            enum.LAGModeActiveLACP,
												Tags:               []design.Tag{},
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
										Tags: []design.Tag{},
									},
								},
								GenericSystems: []design.RackTypeGenericSystem{
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
												AttachmentType:     enum.LinkAttachmentTypeSingle,
												Tags:               []design.Tag{},
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
										Loopback:        &enum.FeatureSwitchDisabled,
										ManagementLevel: enum.SystemManagementLevelUnmanaged,
										Tags:            []design.Tag{},
									},
								},
							},
						},
					},
					ASNAllocationPolicy: &policy.ASNAllocation{SpineASNScheme: enum.ASNAllocationSchemeSingle},
					Capability:          (*enum.TemplateCapability)(nil),
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
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									PortGroups: []design.LogicalDevicePanelPortGroup{
										{Count: 24, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleLeaf}},
										{Count: 8, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric}},
									},
								},
							},
						},
						Tags: []design.Tag{},
					},
					VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
				},
			},
		},
	},
	"L2_superspine_single_plane": {
		Label: testutils.RandString(6, "hex"),
		Superspine: design.Superspine{
			PlaneCount:         1,
			SuperspinePerPlane: 4,
			LogicalDevice: design.LogicalDevice{
				Label: "AOS-32x40-3",
				Panels: []design.LogicalDevicePanel{
					{
						PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 16},
						PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
						PortGroups: []design.LogicalDevicePanelPortGroup{
							{Count: 32, Speed: "40G", Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRoleSpine}},
						},
					},
				},
			},
			Tags: []design.Tag{},
		},
		Pods: []design.PodWithCount{
			{
				Count: 2,
				Pod: design.TemplateRackBased{
					Label: "L2 Pod",
					Racks: []design.RackTypeWithCount{
						{
							Count: 1,
							RackType: design.RackType{
								Label:                    "L2 One Leaf",
								FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
								LeafSwitches: []design.RackTypeLeafSwitch{
									{
										Label:             "leaf",
										LinkPerSpineCount: pointer.To(1),
										LinkPerSpineSpeed: pointer.To(speed.Speed("40G")),
										LogicalDevice: design.LogicalDevice{
											Label: "AOS-64x10+16x40-2",
											Panels: []design.LogicalDevicePanel{
												{
													PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 32},
													PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
													PortGroups: []design.LogicalDevicePanelPortGroup{
														{Count: 64, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}},
													},
												},
												{
													PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
													PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
													PortGroups: []design.LogicalDevicePanelPortGroup{
														{Count: 16, Speed: "40G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleSpine}},
													},
												},
											},
										},
										Tags: []design.Tag{},
									},
								},
								GenericSystems: []design.RackTypeGenericSystem{
									{
										ASNDomain: &enum.FeatureSwitchDisabled,
										Count:     48,
										Label:     "generic",
										Links: []design.RackTypeLink{
											{
												Label:              "link",
												TargetSwitchLabel:  "leaf",
												LinkPerSwitchCount: 1,
												Speed:              "10G",
												AttachmentType:     enum.LinkAttachmentTypeSingle,
												Tags:               []design.Tag{},
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
										Loopback:        &enum.FeatureSwitchDisabled,
										ManagementLevel: enum.SystemManagementLevelUnmanaged,
										Tags:            []design.Tag{},
									},
								},
							},
						},
						{
							Count: 1, RackType: design.RackType{
							Label:                    "L2 Mlag Leaf",
							FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
							LeafSwitches: []design.RackTypeLeafSwitch{
								{
									Label:             "leaf",
									LinkPerSpineCount: pointer.To(1),
									LinkPerSpineSpeed: pointer.To(speed.Speed("40G")),
									LogicalDevice: design.LogicalDevice{
										Label: "AOS-64x10+16x40-2",
										Panels: []design.LogicalDevicePanel{
											{
												PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 32},
												PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
												PortGroups: []design.LogicalDevicePanelPortGroup{
													{Count: 64, Speed: "10G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}},
												},
											},
											{
												PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 8},
												PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
												PortGroups: []design.LogicalDevicePanelPortGroup{
													{Count: 16, Speed: "40G", Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleSpine}},
												},
											},
										},
									},
									RedundancyProtocol: enum.LeafRedundancyProtocolMLAG,
									Tags:               []design.Tag{},
									MLAGInfo: &design.RackTypeLeafSwitchMLAGInfo{
										LeafLeafLinkCount: 4,
										LeafLeafLinkSpeed: "40G",
										MLAGVLAN:          2999,
									},
								},
							},
							GenericSystems: []design.RackTypeGenericSystem{
								{
									ASNDomain: &enum.FeatureSwitchDisabled,
									Count:     24,
									Label:     "generic(1)",
									Links: []design.RackTypeLink{
										{
											Label:              "link",
											TargetSwitchLabel:  "leaf",
											LinkPerSwitchCount: 1,
											Speed:              "10G",
											AttachmentType:     enum.LinkAttachmentTypeSingle,
											SwitchPeer:         enum.LinkSwitchPeerFirst,
											Tags:               []design.Tag{},
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
									Loopback:        &enum.FeatureSwitchDisabled,
									ManagementLevel: enum.SystemManagementLevelUnmanaged,
									Tags:            []design.Tag{},
								},
								{
									ASNDomain: &enum.FeatureSwitchDisabled,
									Count:     24,
									Label:     "generic(2)",
									Links: []design.RackTypeLink{
										{
											Label:              "link",
											TargetSwitchLabel:  "leaf",
											LinkPerSwitchCount: 1,
											Speed:              "10G",
											AttachmentType:     enum.LinkAttachmentTypeSingle,
											SwitchPeer:         enum.LinkSwitchPeerSecond,
											Tags:               []design.Tag{},
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
									Loopback:        &enum.FeatureSwitchDisabled,
									ManagementLevel: enum.SystemManagementLevelUnmanaged,
									Tags:            []design.Tag{},
								},
							},
						},
						},
					},
					ASNAllocationPolicy: &policy.ASNAllocation{
						SpineASNScheme: enum.ASNAllocationSchemeSingle,
					},
					DHCPServiceIntent: policy.DHCPServiceIntent{Active: true},
					Spine: design.Spine{
						Count:                  4,
						LinkPerSuperspineCount: 1,
						LinkPerSuperspineSpeed: "40G",
						LogicalDevice: design.LogicalDevice{
							Label: "AOS-32x40-3",
							Panels: []design.LogicalDevicePanel{
								{
									PanelLayout:  design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 16},
									PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									PortGroups: []design.LogicalDevicePanelPortGroup{
										{Count: 32, Speed: "40G", Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRoleSpine}},
									},
								},
							},
						},
						Tags: []design.Tag{},
					},
					VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
				},
			},
		},
	},
}

func TestTemplatePodBased_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create design.TemplatePodBased
		update design.TemplatePodBased
	}

	testCases := map[string]testCase{
		"L2_superspine_multi_plane_to_L2_superspine_single_plane_with_acs": {
			create: testTemplatesPodBased["L2_superspine_multi_plane"],
			update: testTemplatesPodBased["L2_superspine_single_plane_with_acs"],
		},
		"L2_superspine_single_plane_with_acs_to_L2_superspine_single_plane": {
			create: testTemplatesPodBased["L2_superspine_single_plane_with_acs"],
			update: testTemplatesPodBased["L2_superspine_single_plane"],
		},
		"L2_superspine_single_plane_to_L2_superspine_multi_plane": {
			create: testTemplatesPodBased["L2_superspine_single_plane"],
			update: testTemplatesPodBased["L2_superspine_multi_plane"],
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			require.NotEqual(t, tCase.create, zero.Of(tCase.create)) // make sure we didn't use a bogus map key
			require.NotEqual(t, tCase.update, zero.Of(tCase.update)) // make sure we didn't use a bogus map key

			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					create, update := tCase.create, tCase.update // because we modify these values below

					var id string
					var err error
					var obj design.TemplatePodBased

					// create the object (by type)
					id, err = client.Client.CreateTemplatePodBased2(ctx, create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, 10, func(ctx context.Context) error {
						_ = client.Client.DeleteTemplate2(ctx, id)
						return nil
					})

					// retrieve the object by ID then validate
					template, err := client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok := template.(*design.TemplatePodBased)
					require.True(t, ok)
					idPtr := objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, create, *objPtr)

					// retrieve the object by ID (by type) then validate
					obj, err = client.Client.GetTemplatePodBased2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, create, obj)

					// retrieve the object by label then validate
					template, err = client.Client.GetTemplateByLabel2(ctx, create.Label)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplatePodBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, create, *objPtr)

					// retrieve the object by label (by type) then validate
					obj, err = client.Client.GetTemplatePodBasedByLabel2(ctx, create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, create, obj)

					// retrieve the list of IDs (ours must be in there)
					ids, err := client.Client.ListTemplates2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of IDs (by type) (ours must be in there)
					ids, err = client.Client.ListTemplatesPodBased2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) then validate
					templates, err := client.Client.GetTemplates2(ctx)
					require.NoError(t, err)
					templatePtr := slice.MustFindByID(templates, id)
					require.NotNil(t, templatePtr)
					objPtr, ok = (*templatePtr).(*design.TemplatePodBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, create, *objPtr)

					// retrieve the list of objects (by type) (ours must be in there) then validate
					objs, err := client.Client.GetTemplatesPodBased2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, create, obj)

					// update the object then validate
					update.SetID(id)
					require.NotNil(t, update.ID())
					require.Equal(t, id, *update.ID())
					err = client.Client.UpdateTemplate2(ctx, &update)
					require.NoError(t, err)

					// retrieve the updated object by ID then validate
					template, err = client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplatePodBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, update, *objPtr)

					// retrieve the updated object by ID (by type) type then validate
					obj, err = client.Client.GetTemplatePodBased2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, update, obj)

					// restore the object (by type)
					create.SetID(id)
					require.NotNil(t, create.ID())
					require.Equal(t, id, *update.ID())
					err = client.Client.UpdateTemplatePodBased2(ctx, create)
					require.NoError(t, err)

					// retrieve the restored object by ID then validate
					template, err = client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplatePodBased)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, create, *objPtr)

					// retrieve the restored object by ID (by type) then validate
					obj, err = client.Client.GetTemplatePodBased2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplatePodBased(t, create, obj)

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
					_, err = client.Client.GetTemplatePodBased2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetTemplateByLabel2(ctx, create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label (by type)
					_, err = client.Client.GetTemplatePodBasedByLabel2(ctx, create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListTemplates2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of IDs (by type) (ours must *not* be in there)
					ids, err = client.Client.ListTemplatesPodBased2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					templates, err = client.Client.GetTemplates2(ctx)
					require.NoError(t, err)
					templatePtr = slice.MustFindByID(templates, id)
					require.Nil(t, templatePtr)

					// retrieve the list of objects (by type) (ours must *not* be in there)
					objs, err = client.Client.GetTemplatesPodBased2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateTemplate2(ctx, &update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// update the object (by type)
					err = client.Client.UpdateTemplatePodBased2(ctx, update)
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
