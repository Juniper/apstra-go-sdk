// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package dctestobj

import (
	"context"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestTemplateA(t testing.TB, ctx context.Context, client *apstra.Client) apstra.ObjectId {
	t.Helper()

	rackId := TestRackA(t, ctx, client)

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

func TestTemplateB(t testing.TB, ctx context.Context, client *apstra.Client) apstra.ObjectId {
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
