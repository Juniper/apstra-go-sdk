// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestCreateGetDeleteRackBasedTemplate(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	dn := testutils.RandString(5, "hex")
	req := apstra.CreateRackBasedTemplateRequest{
		DisplayName: dn,
		Spine: &apstra.TemplateElementSpineRequest{
			Count:                  2,
			LinkPerSuperspineSpeed: "10G",
			LogicalDevice:          "AOS-7x10-Spine",
			LinkPerSuperspineCount: 1,
			Tags:                   []apstra.ObjectId{"firewall", "hypervisor"},
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
		id, err := client.Client.CreateRackBasedTemplate(context.TODO(), &req)
		if err != nil {
			t.Fatal(err)
		}

		rbt, err := client.Client.GetRackBasedTemplate(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
		if rbt.Data.DisplayName != dn {
			t.Fatalf("new template displayname mismatch: '%s' vs. '%s'", dn, rbt.Data.DisplayName)
		}

		err = client.Client.DeleteTemplate(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetRackBasedTemplateByName(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	name := "L2 Pod"

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			rbt, err := client.Client.GetRackBasedTemplateByName(ctx, name)
			require.NoError(t, err)
			require.Equal(t, name, rbt.Data.DisplayName)
		})
	}
}

func TestRackBasedTemplateMethods(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	compareSpine := func(t testing.TB, req apstra.TemplateElementSpineRequest, rbt apstra.Spine) {
		t.Helper()

		require.Equal(t, req.Count, rbt.Count)
		require.Equal(t, req.LinkPerSuperspineCount, rbt.LinkPerSuperspineCount)
		require.Equal(t, req.LinkPerSuperspineSpeed, rbt.LinkPerSuperspineSpeed)
		require.Equal(t, len(req.Tags), len(rbt.Tags))
	}

	compareRackInfo := func(t testing.TB, req, rbt apstra.TemplateRackBasedRackInfo) {
		t.Helper()

		require.Equal(t, req.Count, rbt.Count)
	}

	compareRackInfos := func(t testing.TB, req, rbt map[apstra.ObjectId]apstra.TemplateRackBasedRackInfo) {
		t.Helper()

		require.Equal(t, len(req), len(rbt))

		for k, reqRI := range req {
			if rbtRI, ok := rbt[k]; ok {
				compareRackInfo(t, reqRI, rbtRI)
			} else {
				t.Fatalf("rack type mismatch expected rack based info %q not found", k)
			}
		}
	}

	compareRequestToTemplate := func(t testing.TB, req apstra.CreateRackBasedTemplateRequest, rbt apstra.TemplateRackBasedData) {
		t.Helper()

		require.Equal(t, req.DisplayName, rbt.DisplayName)
		compareSpine(t, *req.Spine, rbt.Spine)
		compareRackInfos(t, req.RackInfos, rbt.RackInfo)
		require.Equal(t, req.DhcpServiceIntent.Active, rbt.DhcpServiceIntent.Active)
		if req.AntiAffinityPolicy != nil {
			require.NotNilf(t, rbt.AntiAffinityPolicy, "rbt.AntiAffinityPolicy is nil")
			compare.AntiAffinityPolicy(t, *req.AntiAffinityPolicy, *rbt.AntiAffinityPolicy)
		}

		require.Equal(t, req.AsnAllocationPolicy.SpineAsnScheme, rbt.AsnAllocationPolicy.SpineAsnScheme)
		require.Equal(t, req.VirtualNetworkPolicy.OverlayControlProtocol, rbt.VirtualNetworkPolicy.OverlayControlProtocol)
	}

	type testCase struct {
		request            apstra.CreateRackBasedTemplateRequest
		versionConstraints version.Constraints
	}

	spines := []apstra.TemplateElementSpineRequest{
		{
			Count:                  2,
			LinkPerSuperspineSpeed: "10G",
			LogicalDevice:          "AOS-7x10-Spine",
			LinkPerSuperspineCount: 1,
			Tags:                   []apstra.ObjectId{"firewall", "hypervisor"},
		},
		{
			Count:                  1,
			LinkPerSuperspineSpeed: "10G",
			LogicalDevice:          "AOS-7x10-Spine",
			LinkPerSuperspineCount: 2,
		},
	}

	rackInfos := []map[apstra.ObjectId]apstra.TemplateRackBasedRackInfo{
		{"access_switch": {Count: 1}},
		{"access_switch": {Count: 2}},
	}

	testCases := []testCase{
		{
			versionConstraints: compatibility.EqApstra420,
			request: apstra.CreateRackBasedTemplateRequest{
				DisplayName:          testutils.RandString(5, "hex"),
				Spine:                &spines[0],
				RackInfos:            rackInfos[0],
				DhcpServiceIntent:    &apstra.DhcpServiceIntent{Active: true},
				AntiAffinityPolicy:   &apstra.AntiAffinityPolicy{Algorithm: apstra.AlgorithmHeuristic}, // 4.2.0 only?
				AsnAllocationPolicy:  &apstra.AsnAllocationPolicy{SpineAsnScheme: apstra.AsnAllocationSchemeSingle},
				VirtualNetworkPolicy: &apstra.VirtualNetworkPolicy{},
			},
		},
		{
			versionConstraints: compatibility.EqApstra420,
			request: apstra.CreateRackBasedTemplateRequest{
				DisplayName:          testutils.RandString(5, "hex"),
				Spine:                &spines[1],
				RackInfos:            rackInfos[1],
				DhcpServiceIntent:    &apstra.DhcpServiceIntent{Active: false},
				AntiAffinityPolicy:   &apstra.AntiAffinityPolicy{Algorithm: apstra.AlgorithmHeuristic},
				AsnAllocationPolicy:  &apstra.AsnAllocationPolicy{SpineAsnScheme: apstra.AsnAllocationSchemeSingle},
				VirtualNetworkPolicy: &apstra.VirtualNetworkPolicy{},
			},
		},
		{
			versionConstraints: compatibility.GeApstra421,
			request: apstra.CreateRackBasedTemplateRequest{
				DisplayName:          testutils.RandString(5, "hex"),
				Spine:                &spines[0],
				RackInfos:            rackInfos[0],
				DhcpServiceIntent:    &apstra.DhcpServiceIntent{Active: true},
				AsnAllocationPolicy:  &apstra.AsnAllocationPolicy{SpineAsnScheme: apstra.AsnAllocationSchemeSingle},
				VirtualNetworkPolicy: &apstra.VirtualNetworkPolicy{},
			},
		},
		{
			request: apstra.CreateRackBasedTemplateRequest{
				DisplayName:          testutils.RandString(5, "hex"),
				Spine:                &spines[1],
				RackInfos:            rackInfos[1],
				DhcpServiceIntent:    &apstra.DhcpServiceIntent{Active: false},
				AsnAllocationPolicy:  &apstra.AsnAllocationPolicy{SpineAsnScheme: apstra.AsnAllocationSchemeSingle},
				VirtualNetworkPolicy: &apstra.VirtualNetworkPolicy{},
			},
			versionConstraints: compatibility.GeApstra421,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test-%d", i), func(t *testing.T) {
			ctx := testutils.ContextWithTestID(ctx, t)

			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					if !tc.versionConstraints.Check(client.APIVersion()) {
						t.Skipf("skipping testcase %d because of versionConstraint %s, version %s", i, tc.versionConstraints, client.APIVersion())
					}

					id, err := client.Client.CreateRackBasedTemplate(ctx, &tc.request)
					require.NoError(t, err)

					rbt, err := client.Client.GetRackBasedTemplate(ctx, id)
					require.NoError(t, err)

					require.Equal(t, id, rbt.ID())
					compareRequestToTemplate(t, tc.request, *rbt.Data)

					for j := i; j < i+len(testCases); j++ { // j counts up from i
						k := j % len(testCases) // k counts up from i, but loops back to zero

						if !testCases[k].versionConstraints.Check(client.APIVersion()) {
							continue
						}

						req := testCases[k].request
						err = client.Client.UpdateRackBasedTemplate(ctx, id, &req)
						require.NoError(t, err)

						rbt, err = client.Client.GetRackBasedTemplate(ctx, id)
						require.NoError(t, err)

						require.Equal(t, id, rbt.ID())
						compareRequestToTemplate(t, req, *rbt.Data)
					}

					err = client.Client.DeleteTemplate(ctx, id)
					require.NoError(t, err)
				})
			}
		})
	}
}
