// Copyright (c) Juniper Networks, Inc., 2025-2026.
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
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/policy"
	"github.com/stretchr/testify/require"
)

func TemplateA(t testing.TB, ctx context.Context, client *apstra.Client) apstra.ObjectId {
	t.Helper()

	rackId := RackTypeA(t, ctx, client)

	request := apstra.CreateRackBasedTemplateRequest{
		DisplayName: testutils.RandString(5, "hex"),
		Spine: &apstra.TemplateElementSpineRequest{
			Count:         1,
			LogicalDevice: "AOS-16x40-1",
		},
		RackInfos: map[apstra.ObjectId]apstra.TemplateRackBasedRackInfo{
			rackId: {Count: 1},
		},
		AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
			Algorithm:                apstra.AlgorithmHeuristic,
			MaxLinksPerPort:          1,
			MaxLinksPerSlot:          1,
			MaxPerSystemLinksPerPort: 1,
			MaxPerSystemLinksPerSlot: 1,
			Mode:                     apstra.AntiAffinityModeDisabled,
		},
		AsnAllocationPolicy:  &apstra.AsnAllocationPolicy{SpineAsnScheme: apstra.AsnAllocationSchemeDistinct},
		VirtualNetworkPolicy: &apstra.VirtualNetworkPolicy{OverlayControlProtocol: apstra.OverlayControlProtocolEvpn},
	}

	id, err := client.CreateRackBasedTemplate(ctx, &request)
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteTemplate(ctx, id)
	})

	return id
}

func TemplateB(t testing.TB, ctx context.Context, client *apstra.Client) apstra.ObjectId {
	t.Helper()

	rbt, err := client.GetRackBasedTemplate(ctx, "L2_Virtual")
	require.NoError(t, err)

	rbt.Data.DisplayName = testutils.RandString(5, "hex")
	for k, v := range rbt.Data.RackInfo {
		v.RackTypeData = nil
		rbt.Data.RackInfo[k] = v
	}

	id, err := client.CreateRackBasedTemplate(ctx, &apstra.CreateRackBasedTemplateRequest{
		DisplayName: rbt.Data.DisplayName,
		Spine: &apstra.TemplateElementSpineRequest{
			Count:                  rbt.Data.Spine.Count,
			LinkPerSuperspineSpeed: rbt.Data.Spine.LinkPerSuperspineSpeed,
			LogicalDevice:          "AOS-7x10-Spine",
			LinkPerSuperspineCount: rbt.Data.Spine.LinkPerSuperspineCount,
		},
		RackInfos:            rbt.Data.RackInfo,
		DhcpServiceIntent:    &rbt.Data.DhcpServiceIntent,
		AntiAffinityPolicy:   rbt.Data.AntiAffinityPolicy,
		AsnAllocationPolicy:  &rbt.Data.AsnAllocationPolicy,
		VirtualNetworkPolicy: &rbt.Data.VirtualNetworkPolicy,
	})
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteTemplate(ctx, id)
	})

	return id
}

// TemplateC returns the ID of a rack-based template with a single rack produced by
// the RackTypeB() function.
func TemplateC(t testing.TB, ctx context.Context, client *apstra.Client) string {
	t.Helper()

	rackTypeB, err := client.GetRackType2(ctx, RackTypeB(t, ctx, client))
	require.NoError(t, err)

	request := design.TemplateRackBased{
		Label: testutils.RandString(6, "hex"),
		Racks: []design.RackTypeWithCount{
			{Count: 1, RackType: rackTypeB},
		},
		ASNAllocationPolicy: &policy.ASNAllocation{SpineASNScheme: enum.ASNAllocationSchemeDistinct},
		Spine: design.Spine{
			Count:         1,
			LogicalDevice: rackTypeB.LeafSwitches[0].LogicalDevice,
		},
		VirtualNetworkPolicy: &policy.VirtualNetwork{OverlayControlProtocol: enum.OverlayControlProtocolEVPN},
	}

	id, err := client.CreateTemplate2(ctx, &request)
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, testutils.DefaultCleanupTimeout, func(ctx context.Context) error {
		return client.DeleteTemplate2(ctx, id)
	})

	return id
}
