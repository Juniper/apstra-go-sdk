// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package designtestobj

import (
	"context"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/speed"
	"github.com/stretchr/testify/require"
)

func RackTypeA(t testing.TB, ctx context.Context, client *apstra.Client) apstra.ObjectId {
	t.Helper()

	request := apstra.RackTypeRequest{
		DisplayName:              testutils.RandString(5, "hex"),
		FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
		LeafSwitches: []apstra.RackElementLeafSwitchRequest{
			{
				Label:             testutils.RandString(5, "hex"),
				LinkPerSpineCount: 1,
				LinkPerSpineSpeed: "40G",
				LogicalDeviceId:   "AOS-48x10_6x40-1",
			},
		},
	}

	id, err := client.CreateRackType(ctx, &request)
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteRackType(ctx, id)
	})

	return id
}

// RackTypeB creates a rack with an ESI-enabled leaf pair and seven access switches:
// - one access switch homed only to the first leaf
// - one access switch homed only to the second leaf
// - one access switch dual-homed to both leafs
// - an ESI-enabled single-homed access switch pair
// - an ESI-enabled dual-homed access switch pair
func RackTypeB(t testing.TB, ctx context.Context, client *apstra.Client) string {
	t.Helper()

	ld, err := client.GetLogicalDevice2(ctx, "AOS-48x10_4x100-1")
	require.NoError(t, err)

	request := design.RackType{
		Label:                    testutils.RandString(5, "hex"),
		FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
		LeafSwitches: []design.RackTypeLeafSwitch{
			{
				Label:              "leaf",
				LinkPerSpineCount:  pointer.To(1),
				LinkPerSpineSpeed:  pointer.To(speed.Speed("10G")),
				LogicalDevice:      ld,
				RedundancyProtocol: enum.LeafRedundancyProtocolESI,
			},
		},
		AccessSwitches: []design.RackTypeAccessSwitch{
			{
				Count: 1,
				Label: "access_1",
				Links: []design.RackTypeLink{
					{
						Label:              "link",
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 1,
						Speed:              "10G",
						AttachmentType:     enum.LinkAttachmentTypeSingle,
						SwitchPeer:         enum.LinkSwitchPeerFirst,
						LAGMode:            enum.LAGModeActiveLACP,
					},
				},
				LogicalDevice: ld,
			},
			{
				Count: 1,
				Label: "access_2",
				Links: []design.RackTypeLink{
					{
						Label:              "link",
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 1,
						Speed:              "10G",
						AttachmentType:     enum.LinkAttachmentTypeSingle,
						SwitchPeer:         enum.LinkSwitchPeerSecond,
						LAGMode:            enum.LAGModeActiveLACP,
					},
				},
				LogicalDevice: ld,
			},
			{
				Count: 1,
				Label: "access_3",
				Links: []design.RackTypeLink{
					{
						Label:              "link",
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 1,
						Speed:              "10G",
						AttachmentType:     enum.LinkAttachmentTypeDual,
						LAGMode:            enum.LAGModeActiveLACP,
					},
				},
				LogicalDevice: ld,
			},
			{
				Count: 1,
				Label: "access_4",
				Links: []design.RackTypeLink{
					{
						Label:              "link",
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 1,
						Speed:              "10G",
						AttachmentType:     enum.LinkAttachmentTypeSingle,
						LAGMode:            enum.LAGModeActiveLACP,
					},
				},
				LogicalDevice: ld,
				ESILAGInfo: &design.RackTypeAccessSwitchESILAGInfo{
					LinkCount: 1,
					LinkSpeed: "10G",
				},
			},
			{
				Count: 1,
				Label: "access_5",
				Links: []design.RackTypeLink{
					{
						Label:              "link",
						TargetSwitchLabel:  "leaf",
						LinkPerSwitchCount: 1,
						Speed:              "10G",
						AttachmentType:     enum.LinkAttachmentTypeDual,
						LAGMode:            enum.LAGModeActiveLACP,
					},
				},
				LogicalDevice: ld,
				ESILAGInfo: &design.RackTypeAccessSwitchESILAGInfo{
					LinkCount: 1,
					LinkSpeed: "10G",
				},
			},
		},
	}

	id, err := client.CreateRackType2(ctx, request)
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, testutils.DefaultCleanupTimeout, func(ctx context.Context) error {
		return client.DeleteRackType2(ctx, id)
	})

	return id
}
