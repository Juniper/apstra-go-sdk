// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestCreateGetDeletePodBasedTemplate(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	dn := testutils.RandString(5, "hex")

	rbtdn := "rbtr-" + dn
	rbtr := apstra.CreateRackBasedTemplateRequest{
		DisplayName: rbtdn,
		Spine: &apstra.TemplateElementSpineRequest{
			Count:                  2,
			LinkPerSuperspineSpeed: "10G",
			LogicalDevice:          "AOS-7x10-Spine",
			LinkPerSuperspineCount: 1,
			Tags:                   nil,
		},
		RackInfos: map[apstra.ObjectId]apstra.TemplateRackBasedRackInfo{
			"access_switch": {
				Count: 1,
			},
		},
		DhcpServiceIntent:    &apstra.DhcpServiceIntent{Active: true},
		AntiAffinityPolicy:   &apstra.AntiAffinityPolicy{Algorithm: apstra.AlgorithmHeuristic},
		AsnAllocationPolicy:  &apstra.AsnAllocationPolicy{SpineAsnScheme: apstra.AsnAllocationSchemeSingle},
		VirtualNetworkPolicy: &apstra.VirtualNetworkPolicy{},
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			rbtid, err := client.Client.CreateRackBasedTemplate(ctx, &rbtr)
			require.NoError(t, err)

			pbtdn := "pbtr-" + dn
			pbtr := apstra.CreatePodBasedTemplateRequest{
				DisplayName: pbtdn,
				Superspine: &apstra.TemplateElementSuperspineRequest{
					PlaneCount:         1,
					Tags:               nil,
					SuperspinePerPlane: 4,
					LogicalDeviceId:    "AOS-4x40_8x10-1",
				},
				PodInfos: map[apstra.ObjectId]apstra.TemplatePodBasedInfo{
					rbtid: {
						Count: 1,
					},
				},
				AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
					Algorithm:                apstra.AlgorithmHeuristic,
					MaxLinksPerPort:          1,
					MaxLinksPerSlot:          1,
					MaxPerSystemLinksPerPort: 1,
					MaxPerSystemLinksPerSlot: 1,
					Mode:                     apstra.AntiAffinityModeDisabled,
				},
			}

			pbtid, err := client.Client.CreatePodBasedTemplate(ctx, &pbtr)
			require.NoError(t, err)

			pbt, err := client.Client.GetPodBasedTemplate(ctx, pbtid)
			require.NoError(t, err)

			require.Equal(t, pbtdn, pbt.Data.DisplayName)

			err = client.Client.DeleteTemplate(ctx, pbtid)
			require.NoError(t, err)

			err = client.Client.DeleteTemplate(ctx, rbtid)
			require.NoError(t, err)
		})
	}
}

func TestGetPodBasedTemplateByName(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	name := "L2 superspine single plane"

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			pbt, err := client.Client.GetPodBasedTemplateByName(ctx, name)
			require.NoError(t, err)
			require.Equal(t, name, pbt.Data.DisplayName)
		})
	}
}
