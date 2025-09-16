// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package dctestobj

import (
	"context"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestBlueprintA(t testing.TB, ctx context.Context, client *apstra.Client) *apstra.TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &apstra.CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      testutils.RandString(5, "hex"),
		TemplateId: "L3_Collapsed_ESI",
	})
	require.NoError(t, err)

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	})

	return bpClient
}

func TestBlueprintB(t testing.TB, ctx context.Context, client *apstra.Client) *apstra.TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &apstra.CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      testutils.RandString(5, "hex"),
		TemplateId: "L2_Virtual",
	})
	require.NoError(t, err)

	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	})

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	return bpClient
}

func TestBlueprintC(t testing.TB, ctx context.Context, client *apstra.Client) *apstra.TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &apstra.CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      testutils.RandString(5, "hex"),
		TemplateId: "L2_Virtual_EVPN",
	})
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	})

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	return bpClient
}

func TestBlueprintD(t testing.TB, ctx context.Context, client *apstra.Client) *apstra.TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &apstra.CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      testutils.RandString(5, "hex"),
		TemplateId: "L2_Virtual_ESI_2x_Links",
	})
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	})

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	query := new(apstra.PathQuery).
		SetBlueprintId(bpId).
		SetBlueprintType(apstra.BlueprintTypeStaging).
		SetClient(client).
		Node([]apstra.QEEAttribute{
			apstra.NodeTypeSystem.QEEAttribute(),
			{Key: "system_type", Value: apstra.QEStringVal("switch")},
			{Key: "role", Value: apstra.QEStringVal("leaf")},
			{Key: "name", Value: apstra.QEStringVal("n_leaf")},
		})
	var response struct {
		Items []struct {
			Leaf struct {
				ID string `json:"id"`
			} `json:"n_leaf"`
		} `json:"items"`
	}
	require.NoError(t, query.Do(ctx, &response))

	assignments := make(apstra.SystemIdToInterfaceMapAssignment)
	for _, item := range response.Items {
		assignments[item.Leaf.ID] = "Juniper_vQFX__AOS-7x10-Leaf"
	}

	require.NoError(t, bpClient.SetInterfaceMapAssignments(ctx, assignments))

	return bpClient
}

func TestBlueprintE(t testing.TB, ctx context.Context, client *apstra.Client) *apstra.TwoStageL3ClosClient {
	bpId, err := client.CreateBlueprintFromTemplate(ctx, &apstra.CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      testutils.RandString(5, "hex"),
		TemplateId: "L2_ESI_Access",
	})
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	})

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	leafQuery := new(apstra.PathQuery).
		SetBlueprintId(bpId).
		SetBlueprintType(apstra.BlueprintTypeStaging).
		SetClient(client).
		Node([]apstra.QEEAttribute{
			apstra.NodeTypeSystem.QEEAttribute(),
			{Key: "system_type", Value: apstra.QEStringVal("switch")},
			{Key: "role", Value: apstra.QEStringVal("leaf")},
			{Key: "name", Value: apstra.QEStringVal("n_leaf")},
		})
	var leafResponse struct {
		Items []struct {
			Leaf struct {
				ID string `json:"id"`
			} `json:"n_leaf"`
		} `json:"items"`
	}
	err = leafQuery.Do(ctx, &leafResponse)
	require.NoError(t, err)

	leafAssignements := make(apstra.SystemIdToInterfaceMapAssignment)
	for _, item := range leafResponse.Items {
		leafAssignements[item.Leaf.ID] = "Juniper_vQFX__AOS-7x10-Leaf"
	}
	err = bpClient.SetInterfaceMapAssignments(ctx, leafAssignements)
	require.NoError(t, err)

	accessQuery := new(apstra.PathQuery).
		SetBlueprintId(bpId).
		SetBlueprintType(apstra.BlueprintTypeStaging).
		SetClient(client).
		Node([]apstra.QEEAttribute{
			apstra.NodeTypeSystem.QEEAttribute(),
			{Key: "system_type", Value: apstra.QEStringVal("switch")},
			{Key: "role", Value: apstra.QEStringVal("access")},
			{Key: "name", Value: apstra.QEStringVal("n_access")},
		})
	var accessResponse struct {
		Items []struct {
			Leaf struct {
				ID string `json:"id"`
			} `json:"n_access"`
		} `json:"items"`
	}
	require.NoError(t, accessQuery.Do(ctx, &accessResponse))

	accessAssignements := make(apstra.SystemIdToInterfaceMapAssignment)
	for _, item := range accessResponse.Items {
		accessAssignements[item.Leaf.ID] = "Juniper_vQFX__AOS-8x10-1"
	}
	require.NoError(t, bpClient.SetInterfaceMapAssignments(ctx, accessAssignements))

	return bpClient
}

func TestBlueprintF(t testing.TB, ctx context.Context, client *apstra.Client) *apstra.TwoStageL3ClosClient {
	templateId := TestTemplateA(t, ctx, client)

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &apstra.CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      testutils.RandString(5, "hex"),
		TemplateId: templateId,
	})
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	})

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	return bpClient
}

func TestBlueprintG(t testing.TB, ctx context.Context, client *apstra.Client) *apstra.TwoStageL3ClosClient {
	t.Helper()

	templateId := TestTemplateB(t, ctx, client)

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &apstra.CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      testutils.RandString(5, "hex"),
		TemplateId: templateId,
		FabricSettings: &apstra.FabricSettings{
			SpineSuperspineLinks: testutils.ToPtr(apstra.AddressingSchemeIp46),
			SpineLeafLinks:       testutils.ToPtr(apstra.AddressingSchemeIp46),
		},
	})
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	})

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	require.NoError(t, bpClient.SetFabricSettings(ctx, &apstra.FabricSettings{Ipv6Enabled: testutils.ToPtr(true)}))

	return bpClient
}

// TestBlueprintH creates a test blueprint using client and returns a *TwoStageL3ClosClient.
// The blueprint will use a dual-stack fabric and have ipv6 enabled.
func TestBlueprintH(t testing.TB, ctx context.Context, client *apstra.Client) *apstra.TwoStageL3ClosClient {
	t.Helper()

	bpRequest := apstra.CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      testutils.RandString(5, "hex"),
		TemplateId: "L2_Virtual_EVPN",
		FabricSettings: &apstra.FabricSettings{
			SpineSuperspineLinks: testutils.ToPtr(apstra.AddressingSchemeIp46),
			SpineLeafLinks:       testutils.ToPtr(apstra.AddressingSchemeIp46),
		},
	}

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &bpRequest)
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	})

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	// enable IPv6
	require.NoError(t, bpClient.SetFabricSettings(ctx, &apstra.FabricSettings{Ipv6Enabled: testutils.ToPtr(true)}))

	return bpClient
}

// TestBlueprintI returns a collapsed fabric which has been committed and has no build errors
func TestBlueprintI(t testing.TB, ctx context.Context, client *apstra.Client) *apstra.TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &apstra.CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      testutils.RandString(5, "hex"),
		TemplateId: "L3_Collapsed_ESI",
	})
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	})

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	// assign leaf interface maps
	leafIds, err := testutils.GetSystemIdsByRole(ctx, bpClient, "leaf")
	require.NoError(t, err)
	mappings := make(apstra.SystemIdToInterfaceMapAssignment, len(leafIds))
	for _, leafId := range leafIds {
		mappings[leafId.String()] = "Juniper_vQFX__AOS-7x10-Leaf"
	}
	err = bpClient.SetInterfaceMapAssignments(ctx, mappings)
	require.NoError(t, err)

	// set leaf loopback pool
	err = bpClient.SetResourceAllocation(ctx, &apstra.ResourceGroupAllocation{
		ResourceGroup: apstra.ResourceGroup{
			Type: apstra.ResourceTypeIp4Pool,
			Name: apstra.ResourceGroupNameLeafIp4,
		},
		PoolIds: []apstra.ObjectId{"Private-10_0_0_0-8"},
	})
	require.NoError(t, err)

	// set leaf-leaf pool
	err = bpClient.SetResourceAllocation(ctx, &apstra.ResourceGroupAllocation{
		ResourceGroup: apstra.ResourceGroup{
			Type: apstra.ResourceTypeIp4Pool,
			Name: apstra.ResourceGroupNameLeafLeafIp4,
		},
		PoolIds: []apstra.ObjectId{"Private-10_0_0_0-8"},
	})
	require.NoError(t, err)

	// set leaf ASN pool
	err = bpClient.SetResourceAllocation(ctx, &apstra.ResourceGroupAllocation{
		ResourceGroup: apstra.ResourceGroup{
			Type: apstra.ResourceTypeAsnPool,
			Name: apstra.ResourceGroupNameLeafAsn,
		},
		PoolIds: []apstra.ObjectId{"Private-64512-65534"},
	})
	require.NoError(t, err)

	// set VN VNI pool
	err = bpClient.SetResourceAllocation(ctx, &apstra.ResourceGroupAllocation{
		ResourceGroup: apstra.ResourceGroup{
			Type: apstra.ResourceTypeVniPool,
			Name: apstra.ResourceGroupNameEvpnL3Vni,
		},
		PoolIds: []apstra.ObjectId{"Default-10000-20000"},
	})
	require.NoError(t, err)

	// set VN VNI pool
	err = bpClient.SetResourceAllocation(ctx, &apstra.ResourceGroupAllocation{
		ResourceGroup: apstra.ResourceGroup{
			Type: apstra.ResourceTypeVniPool,
			Name: apstra.ResourceGroupNameVxlanVnIds,
		},
		PoolIds: []apstra.ObjectId{"Default-10000-20000"},
	})
	require.NoError(t, err)

	// commit
	bpStatus, err := client.GetBlueprintStatus(ctx, bpClient.Id())
	require.NoError(t, err)
	_, err = client.DeployBlueprint(ctx, &apstra.BlueprintDeployRequest{
		Id:          bpClient.Id(),
		Description: "initial commit in test: " + t.Name(),
		Version:     bpStatus.Version,
	})
	require.NoError(t, err)

	return bpClient
}
