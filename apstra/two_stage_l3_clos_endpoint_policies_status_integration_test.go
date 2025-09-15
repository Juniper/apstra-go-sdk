// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func TestGetAllConnectivityTemplateStatus(t *testing.T) {
	ctx := context.Background()

	connectivityTemplateCount := 3

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	getApplicationPointIds := func(t *testing.T, ctx context.Context, bp *TwoStageL3ClosClient) []ObjectId {
		t.Helper()

		query := new(PathQuery).
			SetBlueprintId(bp.Id()).
			SetClient(bp.Client()).
			Node([]QEEAttribute{
				NodeTypeSystem.QEEAttribute(),
				{Key: "system_type", Value: QEStringVal("switch")},
			}).
			Out([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeInterface.QEEAttribute(),
				{Key: "name", Value: QEStringVal("server_interface")},
			}).
			Out([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
			Node([]QEEAttribute{NodeTypeLink.QEEAttribute()}).
			In([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
			Node([]QEEAttribute{NodeTypeInterface.QEEAttribute()}).
			In([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeSystem.QEEAttribute(),
				{Key: "system_type", Value: QEStringVal("server")},
			})

		var result struct {
			Items []struct {
				Interface struct {
					Id ObjectId `json:"id"`
				} `json:"server_interface"`
			} `json:"items"`
		}

		err = query.Do(ctx, &result)
		require.NoError(t, err)

		applicationPointIds := make([]ObjectId, len(result.Items))
		for i, item := range result.Items {
			applicationPointIds[i] = item.Interface.Id
		}

		return applicationPointIds
	}

	createCt := func(t *testing.T, ctx context.Context, bp *TwoStageL3ClosClient, vlan Vlan) EndpointPolicyStatus {
		t.Helper()

		result := EndpointPolicyStatus{
			Label:       randString(6, "hex"),
			Description: randString(6, "hex"),
			Status:      enum.EndpointPolicyStatusReady,
			Tags:        randStrings(3, 6),
			topLevel:    true,
		}

		szLabel := randString(6, "hex")
		szId, err := bp.CreateSecurityZone(ctx, &SecurityZoneData{
			Label:   szLabel,
			SzType:  SecurityZoneTypeEVPN,
			VrfName: szLabel,
		})
		require.NoError(t, err)

		vlanPtr := &vlan
		if vlan == 0 { // this produces an "incomplete" CT
			vlanPtr = nil
			result.Status = enum.EndpointPolicyStatusIncomplete
		}

		ct := ConnectivityTemplate{
			Label:       result.Label,
			Description: result.Description,
			Tags:        result.Tags,
			Subpolicies: []*ConnectivityTemplatePrimitive{
				{
					Attributes: &ConnectivityTemplatePrimitiveAttributesAttachLogicalLink{
						SecurityZone:       &szId,
						Tagged:             true,
						Vlan:               vlanPtr,
						IPv4AddressingType: CtPrimitiveIPv4AddressingTypeNumbered,
					},
				},
			},
		}
		require.NoError(t, ct.SetIds())
		require.NoError(t, ct.SetUserData())
		require.NoError(t, bp.CreateConnectivityTemplate(ctx, &ct))

		result.Id = *ct.Id

		return result
	}

	compare := func(t *testing.T, a, b EndpointPolicyStatus) {
		t.Helper()

		require.Equal(t, a.Id, b.Id)
		require.Equal(t, a.Label, b.Label)
		require.Equal(t, a.Description, b.Description)
		require.Equal(t, a.Status, b.Status)
		require.Equal(t, a.AppPointsCount, b.AppPointsCount)
		compareSlicesAsSets(t, a.Tags, b.Tags, "tag mismatch")
		require.Equal(t, a.topLevel, b.topLevel)
	}

	assignCt := func(t *testing.T, ctx context.Context, bp *TwoStageL3ClosClient, apIds []ObjectId, epStatus EndpointPolicyStatus) {
		t.Helper()

		selectedAPIds := samples(t, apIds, epStatus.AppPointsCount)
		assignments := make(map[ObjectId]map[ObjectId]bool, len(selectedAPIds))
		for _, apId := range selectedAPIds {
			assignments[apId] = map[ObjectId]bool{epStatus.Id: true}
		}

		err := bp.SetApplicationPointsConnectivityTemplates(ctx, assignments)
		require.NoError(t, err)
	}

	var vlan Vlan = 10

	for clientName, client := range clients {
		t.Run(clientName, func(t *testing.T) {
			t.Parallel()

			t.Log("Creating Blueprint")
			bpClient := testBlueprintC(ctx, t, client.client)

			t.Logf("Creating %d Connectivity Templates", connectivityTemplateCount)
			expectedConnectivityTemplateStatuses := make([]EndpointPolicyStatus, connectivityTemplateCount)
			for i := range expectedConnectivityTemplateStatuses {
				vlan++
				expectedConnectivityTemplateStatuses[i] = createCt(t, ctx, bpClient, vlan)
			}

			t.Logf("Checking %d Connectivity Templates", connectivityTemplateCount)
			actualConnectivityTemplateStatuses, err := bpClient.GetAllConnectivityTemplateStatus(ctx)
			require.NoError(t, err)
			require.Equal(t, connectivityTemplateCount, len(actualConnectivityTemplateStatuses))
			for _, expectedConnectivityTemplateStatus := range expectedConnectivityTemplateStatuses {
				require.Contains(t, actualConnectivityTemplateStatuses, expectedConnectivityTemplateStatus.Id)
				compare(t, expectedConnectivityTemplateStatus, actualConnectivityTemplateStatuses[expectedConnectivityTemplateStatus.Id])
			}

			t.Log("Discovering Application Points")
			applicationPointIds := getApplicationPointIds(t, ctx, bpClient)

			t.Logf("Distribuging Connectivity Templates across %d Application Points", len(applicationPointIds))
			for i := range expectedConnectivityTemplateStatuses {
				expectedConnectivityTemplateStatuses[i].Status = enum.EndpointPolicyStatusAssigned
				expectedConnectivityTemplateStatuses[i].AppPointsCount = rand.Intn(connectivityTemplateCount) + 1
				assignCt(t, ctx, bpClient, applicationPointIds, expectedConnectivityTemplateStatuses[i])
			}

			t.Logf("Checking %d Connectivity Templates", connectivityTemplateCount)
			actualConnectivityTemplateStatuses, err = bpClient.GetAllConnectivityTemplateStatus(ctx)
			require.NoError(t, err)
			require.Equal(t, connectivityTemplateCount, len(actualConnectivityTemplateStatuses))
			for _, expectedConnectivityTemplateStatus := range expectedConnectivityTemplateStatuses {
				require.Contains(t, actualConnectivityTemplateStatuses, expectedConnectivityTemplateStatus.Id)
				compare(t, expectedConnectivityTemplateStatus, actualConnectivityTemplateStatuses[expectedConnectivityTemplateStatus.Id])
			}

			t.Log("Creating incomplete Connectivity Template")
			expected := createCt(t, ctx, bpClient, 0)

			t.Log("Checking incomplete Connectivity Template")
			actualConnectivityTemplateStatuses, err = bpClient.GetAllConnectivityTemplateStatus(ctx)
			require.NoError(t, err)
			require.Contains(t, actualConnectivityTemplateStatuses, expected.Id)
			compare(t, expected, actualConnectivityTemplateStatuses[expected.Id])

			log.Println(applicationPointIds)
		})
	}
}
