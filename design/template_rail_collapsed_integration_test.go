// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package design_test

import (
	"context"
	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	comparedesign "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/design"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/Juniper/apstra-go-sdk/policy"
	"github.com/stretchr/testify/require"
	"testing"
)

var testTemplatesRailCollapsed = map[string]design.TemplateRailCollapsed{
	"collapsed_fabric_128gpu": {
		Label: testutils.RandString(6, "hex"),
		Racks: []design.RackTypeWithCount{
			{
				Count: 1,
				RackType: design.RackType{
					Label:                    "Collapsed 128GPU",
					FabricConnectivityDesign: enum.FabricConnectivityDesignRailCollapsed,
					LeafSwitches: []design.LeafSwitch{
						{
							Label: "leaf_1",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_1", Description: "Rail 1"},
								{Label: "rail_2", Description: "Rail 2"},
								{Label: "rail_3", Description: "Rail 3"},
								{Label: "rail_4", Description: "Rail 4"},
								{Label: "rail_5", Description: "Rail 5"},
								{Label: "rail_6", Description: "Rail 6"},
								{Label: "rail_7", Description: "Rail 7"},
								{Label: "rail_8", Description: "Rail 8"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
					},
					GenericSystems: []design.GenericSystem{
						{
							ASNDomain: &enum.FeatureSwitchDisabled,
							Count:     16,
							Label:     "server",
							Links: []design.RackTypeLink{
								{
									Label:              "link_1",
									TargetSwitchLabel:  "leaf_1",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(1),
									Tags: []design.Tag{
										{Label: "rail_1", Description: "Rail 1"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_2",
									TargetSwitchLabel:  "leaf_1",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(2),
									Tags: []design.Tag{
										{Label: "rail_2", Description: "Rail 2"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_3",
									TargetSwitchLabel:  "leaf_1",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(3),
									Tags: []design.Tag{
										{Label: "rail_3", Description: "Rail 3"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_4",
									TargetSwitchLabel:  "leaf_1",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(4),
									Tags: []design.Tag{
										{Label: "rail_4", Description: "Rail 4"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_5",
									TargetSwitchLabel:  "leaf_1",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(5),
									Tags: []design.Tag{
										{Label: "rail_5", Description: "Rail 5"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_6",
									TargetSwitchLabel:  "leaf_1",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(6),
									Tags: []design.Tag{
										{Label: "rail_6", Description: "Rail 6"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_7",
									TargetSwitchLabel:  "leaf_1",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(7),
									Tags: []design.Tag{
										{Label: "rail_7", Description: "Rail 7"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_8",
									TargetSwitchLabel:  "leaf_1",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(8),
									Tags: []design.Tag{
										{Label: "rail_8", Description: "Rail 8"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
							},
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-8x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 4},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 8,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
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
		DHCPServiceIntent:    policy.DHCPServiceIntent{Active: true},
		VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
	},
	"collapsed_fabric_512gpu": {
		Label: testutils.RandString(6, "hex"),
		Racks: []design.RackTypeWithCount{
			{
				Count: 1,
				RackType: design.RackType{
					Label:                    "Collapsed 512GPU",
					FabricConnectivityDesign: enum.FabricConnectivityDesignRailCollapsed,
					LeafSwitches: []design.LeafSwitch{
						{
							Label: "leaf_1",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_1", Description: "Rail 1"},
								{Label: "rail_2", Description: "Rail 2"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
						{
							Label: "leaf_2",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_3", Description: "Rail 3"},
								{Label: "rail_4", Description: "Rail 4"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
						{
							Label: "leaf_3",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_5", Description: "Rail 5"},
								{Label: "rail_6", Description: "Rail 6"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
						{
							Label: "leaf_4",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_7", Description: "Rail 7"},
								{Label: "rail_8", Description: "Rail 8"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
					},
					GenericSystems: []design.GenericSystem{
						{
							ASNDomain: &enum.FeatureSwitchDisabled,
							Count:     64,
							Label:     "server",
							Links: []design.RackTypeLink{
								{
									Label:              "link_1",
									TargetSwitchLabel:  "leaf_1",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(1),
									Tags: []design.Tag{
										{Label: "rail_1", Description: "Rail 1"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_2",
									TargetSwitchLabel:  "leaf_1",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(2),
									Tags: []design.Tag{
										{Label: "rail_2", Description: "Rail 2"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_3",
									TargetSwitchLabel:  "leaf_2",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(3),
									Tags: []design.Tag{
										{Label: "rail_3", Description: "Rail 3"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_4",
									TargetSwitchLabel:  "leaf_2",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(4),
									Tags: []design.Tag{
										{Label: "rail_4", Description: "Rail 4"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_5",
									TargetSwitchLabel:  "leaf_3",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(5),
									Tags: []design.Tag{
										{Label: "rail_5", Description: "Rail 5"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_6",
									TargetSwitchLabel:  "leaf_3",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(6),
									Tags: []design.Tag{
										{Label: "rail_6", Description: "Rail 6"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_7",
									TargetSwitchLabel:  "leaf_4",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(7),
									Tags: []design.Tag{
										{Label: "rail_7", Description: "Rail 7"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_8",
									TargetSwitchLabel:  "leaf_4",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(8),
									Tags: []design.Tag{
										{Label: "rail_8", Description: "Rail 8"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
							},
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-8x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 4},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 8,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
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
		DHCPServiceIntent:    policy.DHCPServiceIntent{Active: true},
		VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
	},
	"collapsed_fabric_1024gpu": {
		Label: testutils.RandString(6, "hex"),
		Racks: []design.RackTypeWithCount{
			{
				Count: 1,
				RackType: design.RackType{
					Label:                    "Collapsed 512GPU",
					FabricConnectivityDesign: enum.FabricConnectivityDesignRailCollapsed,
					LeafSwitches: []design.LeafSwitch{
						{
							Label: "leaf_1",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_1", Description: "Rail 1"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
						{
							Label: "leaf_2",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_2", Description: "Rail 2"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
						{
							Label: "leaf_3",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_3", Description: "Rail 3"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
						{
							Label: "leaf_4",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_4", Description: "Rail 4"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
						{
							Label: "leaf_5",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_5", Description: "Rail 5"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
						{
							Label: "leaf_6",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_6", Description: "Rail 6"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
						{
							Label: "leaf_7",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_7", Description: "Rail 7"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
						{
							Label: "leaf_8",
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-128x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 4, ColumnCount: 32},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 128,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleSuperspine, enum.PortRoleUnused, enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess, enum.PortRoleSpine},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
									},
								},
							},
							Tags: []design.Tag{
								{Label: "rail_8", Description: "Rail 8"},
								{Label: "stripe_1", Description: "Stripe 1"},
							},
						},
					},
					GenericSystems: []design.GenericSystem{
						{
							ASNDomain: &enum.FeatureSwitchDisabled,
							Count:     128,
							Label:     "server",
							Links: []design.RackTypeLink{
								{
									Label:              "link_1",
									TargetSwitchLabel:  "leaf_1",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(1),
									Tags: []design.Tag{
										{Label: "rail_1", Description: "Rail 1"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_2",
									TargetSwitchLabel:  "leaf_2",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(2),
									Tags: []design.Tag{
										{Label: "rail_2", Description: "Rail 2"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_3",
									TargetSwitchLabel:  "leaf_3",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(3),
									Tags: []design.Tag{
										{Label: "rail_3", Description: "Rail 3"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_4",
									TargetSwitchLabel:  "leaf_4",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(4),
									Tags: []design.Tag{
										{Label: "rail_4", Description: "Rail 4"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_5",
									TargetSwitchLabel:  "leaf_5",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(5),
									Tags: []design.Tag{
										{Label: "rail_5", Description: "Rail 5"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_6",
									TargetSwitchLabel:  "leaf_6",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(6),
									Tags: []design.Tag{
										{Label: "rail_6", Description: "Rail 6"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_7",
									TargetSwitchLabel:  "leaf_7",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(7),
									Tags: []design.Tag{
										{Label: "rail_7", Description: "Rail 7"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
								{
									Label:              "link_8",
									TargetSwitchLabel:  "leaf_8",
									LinkPerSwitchCount: 1,
									Speed:              "400G",
									AttachmentType:     enum.LinkAttachmentTypeSingle,
									RailIndex:          pointer.To(8),
									Tags: []design.Tag{
										{Label: "rail_8", Description: "Rail 8"},
										{Label: "stripe_1", Description: "Stripe 1"},
									},
								},
							},
							LogicalDevice: design.LogicalDevice{
								Label: "AOS-8x400-1",
								Panels: []design.LogicalDevicePanel{
									{
										PanelLayout: design.LogicalDevicePanelLayout{RowCount: 2, ColumnCount: 4},
										PortGroups: []design.LogicalDevicePanelPortGroup{
											{
												Count: 8,
												Speed: "400G",
												Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleGeneric, enum.PortRolePeer, enum.PortRoleAccess},
											},
										},
										PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
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
		DHCPServiceIntent:    policy.DHCPServiceIntent{Active: true},
		VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolNone},
	},
}

func TestTemplateRailCollapsed_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create design.TemplateRailCollapsed
		update design.TemplateRailCollapsed
	}

	testCases := map[string]testCase{
		"collapsed_fabric_128gpu_to_collapsed_fabric_512gpu": {
			create: testTemplatesRailCollapsed["collapsed_fabric_128gpu"],
			update: testTemplatesRailCollapsed["collapsed_fabric_512gpu"],
		},
		"collapsed_fabric_512gpu_to_collapsed_fabric_1024gpu": {
			create: testTemplatesRailCollapsed["collapsed_fabric_512gpu"],
			update: testTemplatesRailCollapsed["collapsed_fabric_1024gpu"],
		},
		"collapsed_fabric_1024gpu_to_collapsed_fabric_128gpu": {
			create: testTemplatesRailCollapsed["collapsed_fabric_1024gpu"],
			update: testTemplatesRailCollapsed["collapsed_fabric_128gpu"],
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
					var obj design.TemplateRailCollapsed

					// create the object (by type)
					id, err = client.Client.CreateTemplateRailCollapsed2(ctx, tCase.create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, 10, func(ctx context.Context) error {
						_ = client.Client.DeleteTemplate2(ctx, id)
						return nil
					})

					// retrieve the object by ID then validate
					template, err := client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok := template.(*design.TemplateRailCollapsed)
					require.True(t, ok)
					idPtr := objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRailCollapsed(t, tCase.create, *objPtr)

					// retrieve the object by ID (by type) then validate
					obj, err = client.Client.GetTemplateRailCollapsed2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRailCollapsed(t, tCase.create, obj)

					// retrieve the object by label then validate
					template, err = client.Client.GetTemplateByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplateRailCollapsed)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRailCollapsed(t, tCase.create, *objPtr)

					// retrieve the object by label (by type) then validate
					obj, err = client.Client.GetTemplateRailCollapsedByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRailCollapsed(t, tCase.create, obj)

					// retrieve the list of IDs (ours must be in there)
					ids, err := client.Client.ListTemplates2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of IDs (by type) (ours must be in there)
					ids, err = client.Client.ListTemplatesRailCollapsed2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) then validate
					templates, err := client.Client.GetTemplates2(ctx)
					require.NoError(t, err)
					templatePtr := slice.MustFindByID(templates, id)
					require.NotNil(t, templatePtr)
					objPtr, ok = (*templatePtr).(*design.TemplateRailCollapsed)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRailCollapsed(t, tCase.create, *objPtr)

					// retrieve the list of objects (by type) (ours must be in there) then validate
					objs, err := client.Client.GetTemplatesRailCollapsed2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRailCollapsed(t, tCase.create, obj)

					// update the object then validate
					require.NoError(t, tCase.update.SetID(id))
					require.NotNil(t, tCase.update.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateTemplate2(ctx, &tCase.update)
					require.NoError(t, err)

					// retrieve the updated object by ID then validate
					template, err = client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplateRailCollapsed)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRailCollapsed(t, tCase.update, *objPtr)

					// retrieve the updated object by ID (by type) type then validate
					update, err := client.Client.GetTemplateRailCollapsed2(ctx, id)
					require.NoError(t, err)
					idPtr = update.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRailCollapsed(t, tCase.update, update)

					// restore the object (by type)
					require.NoError(t, tCase.create.SetID(id))
					require.NotNil(t, tCase.create.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateTemplateRailCollapsed2(ctx, tCase.create)
					require.NoError(t, err)

					// retrieve the restored object by ID then validate
					template, err = client.Client.GetTemplate2(ctx, id)
					require.NoError(t, err)
					objPtr, ok = template.(*design.TemplateRailCollapsed)
					require.True(t, ok)
					idPtr = objPtr.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRailCollapsed(t, tCase.create, *objPtr)

					// retrieve the restored object by ID (by type) then validate
					obj, err = client.Client.GetTemplateRailCollapsed2(ctx, id)
					require.NoError(t, err)
					idPtr = update.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.TemplateRailCollapsed(t, tCase.create, obj)

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
					_, err = client.Client.GetTemplateRailCollapsed2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetTemplateByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label (by type)
					_, err = client.Client.GetTemplateRailCollapsedByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListTemplates2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of IDs (by type) (ours must *not* be in there)
					ids, err = client.Client.ListTemplatesRailCollapsed2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					templates, err = client.Client.GetTemplates2(ctx)
					require.NoError(t, err)
					templatePtr = slice.MustFindByID(templates, id)
					require.Nil(t, templatePtr)

					// retrieve the list of objects (by type) (ours must *not* be in there)
					objs, err = client.Client.GetTemplatesRailCollapsed2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateTemplate2(ctx, &tCase.update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// update the object (by type)
					err = client.Client.UpdateTemplateRailCollapsed2(ctx, tCase.update)
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
