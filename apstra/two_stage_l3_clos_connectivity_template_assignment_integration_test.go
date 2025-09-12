// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func compareConnectivityTemplateAssignments(a, b map[ObjectId]bool, applicationPointId ObjectId, t *testing.T) {
	if len(a) != len(b) {
		t.Fatalf("Connectivity template assignment maps for interface %q have different length: %d vs. %d", applicationPointId, len(a), len(b))
	}

	for ctId, aUsed := range a {
		var ok bool
		var bUsed bool
		if bUsed, ok = b[ctId]; !ok {
			t.Fatalf("Connectivity template assignment maps for interface %q don't both have connectivty template %q", applicationPointId, ctId)
		}

		if aUsed != bUsed {
			t.Fatalf("Connectivity template assignment maps for interface %q don't agree about connectivty template %q: a: %t b: %t",
				applicationPointId, ctId, aUsed, bUsed)
		}
	}
}

func compareInterfacesConnectivityTemplateAssignments(a, b map[ObjectId]map[ObjectId]bool, t *testing.T) {
	if len(a) != len(b) {
		t.Fatalf("Connectivity template assignment maps have different length: %d vs. %d", len(a), len(b))
	}

	for applicationPointId, aCTs := range a {
		// aCTs and bCTs are map[CT ID]bool indicating whether the CT is applied to applicationPointId
		var ok bool
		var bCTs map[ObjectId]bool
		if bCTs, ok = b[applicationPointId]; !ok {
			t.Fatalf("Connectivity template assignment map key missing: %q", applicationPointId)
		}

		compareConnectivityTemplateAssignments(aCTs, bCTs, applicationPointId, t)
	}
}

func TestAssignClearCtToInterface(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	vnCount := 2

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			bpClient := testBlueprintC(ctx, t, client.client)

			leafIds, err := getSystemIdsByRole(ctx, bpClient, "leaf")
			if err != nil {
				t.Fatal(err)
			}

			// create assignments for the leaf switch nodes
			assignments := make(SystemIdToInterfaceMapAssignment, len(leafIds))
			bindings := make([]VnBinding, len(leafIds))
			for i, leafId := range leafIds {
				assignments[leafId.String()] = "Juniper_vQFX__AOS-7x10-Leaf"
				bindings[i] = VnBinding{SystemId: leafId}
			}

			log.Printf("testing SetInterfaceMapAssignments() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetInterfaceMapAssignments(ctx, assignments)
			if err != nil {
				t.Fatal(err)
			}

			vrf := randString(5, "hex")
			szId, err := bpClient.CreateSecurityZone(ctx, &SecurityZoneData{
				Label:   vrf,
				SzType:  SecurityZoneTypeEVPN,
				VrfName: vrf,
			})
			if err != nil {
				t.Fatal(err)
			}

			vnIds := make([]ObjectId, vnCount)
			cts := make([]ConnectivityTemplate, vnCount)
			ctIds := make([]ObjectId, vnCount)
			for i := 0; i < vnCount; i++ {
				log.Printf("testing CreateVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				vnIds[i], err = bpClient.CreateVirtualNetwork(ctx, &VirtualNetworkData{
					Label:          randString(6, "hex"),
					SecurityZoneId: szId,
					VnBindings:     bindings,
					VnType:         enum.VnTypeVxlan,
				})
				if err != nil {
					t.Fatal(err)
				}

				cts[i] = ConnectivityTemplate{
					Label: randString(5, "hex"),
					Subpolicies: []*ConnectivityTemplatePrimitive{{
						Attributes: &ConnectivityTemplatePrimitiveAttributesAttachSingleVlan{
							Tagged:   true,
							VnNodeId: &vnIds[i],
						},
					}},
				}

				err = cts[i].SetIds()
				if err != nil {
					t.Fatal(err)
				}

				err = cts[i].SetUserData()
				if err != nil {
					t.Fatal(err)
				}

				log.Printf("testing CreateConnectivityTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				err = bpClient.CreateConnectivityTemplate(ctx, &cts[i])
				if err != nil {
					t.Fatal(err)
				}

				ctIds[i] = *cts[i].Id
			}

			// Graph query which picks out generic-facing interfaces on leaf switches
			query := new(PathQuery).
				SetBlueprintType(BlueprintTypeStaging).
				SetBlueprintId(bpClient.blueprintId).
				SetClient(bpClient.client).
				Node([]QEEAttribute{
					NodeTypeSystem.QEEAttribute(),
					{"role", QEStringVal("leaf")},
				}).
				Out([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
				Node([]QEEAttribute{
					NodeTypeInterface.QEEAttribute(),
					{"name", QEStringVal("switch_port")},
				}).
				Out([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
				Node([]QEEAttribute{NodeTypeLink.QEEAttribute()}).
				In([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
				Node([]QEEAttribute{NodeTypeInterface.QEEAttribute()}).
				In([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
				Node([]QEEAttribute{
					NodeTypeSystem.QEEAttribute(),
					{"role", QEStringVal("generic")},
				})

			var queryResponse struct {
				Items []struct {
					Interface struct {
						Id ObjectId `json:"id"`
					} `json:"switch_port"`
				} `json:"items"`
			}

			err = query.Do(ctx, &queryResponse)
			if err != nil {
				t.Fatal(err)
			}
			if len(queryResponse.Items) == 0 {
				t.Fatal("graph query found no generic-system-facing leaf switch ports")
			}

			// collect the server-facing interfaces of leaf switches
			leafInterfaceIds := make([]ObjectId, len(queryResponse.Items))
			for i, item := range queryResponse.Items {
				leafInterfaceIds[i] = item.Interface.Id
			}

			// assign a CT to a lone interface
			log.Printf("testing SetApplicationPointConnectivityTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetApplicationPointConnectivityTemplates(ctx, leafInterfaceIds[0], ctIds)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetInterfaceConnectivityTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			assignedCts, err := bpClient.GetInterfaceConnectivityTemplates(ctx, leafInterfaceIds[0])
			if err != nil {
				t.Fatal(err)
			}

			compareSlicesAsSets(t, ctIds, assignedCts, "assigned slices do not match intent")

			err = bpClient.DelApplicationPointConnectivityTemplates(ctx, leafInterfaceIds[0], ctIds)
			if err != nil {
				t.Fatal(err)
			}

			assignedCts, err = bpClient.GetInterfaceConnectivityTemplates(ctx, leafInterfaceIds[0])
			if err != nil {
				t.Fatal(err)
			}

			if len(assignedCts) != 0 {
				t.Fatalf("expected 0 connectivity templates assigned to interface, got %d", len(assignedCts))
			}

			// create the outer map (keyed by application point ID)
			ctAssignments := make(map[ObjectId]map[ObjectId]bool, len(leafInterfaceIds))
			for _, interfaceId := range leafInterfaceIds {
				// create the inner map (keyed by connectivity template ID)
				ctAssignments[interfaceId] = make(map[ObjectId]bool, len(ctIds))
				for _, ctId := range ctIds {
					ctAssignments[interfaceId][ctId] = randBool()
				}
			}

			// set the assignments selected above
			log.Printf("testing SetApplicationPointsConnectivityTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetApplicationPointsConnectivityTemplates(ctx, ctAssignments)
			if err != nil {
				t.Fatal(err)
			}

			// retrieve the assignments
			log.Printf("testing GetConnectivityTemplatesByApplicationPoints() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			apToPolicyInfo, err := bpClient.GetConnectivityTemplatesByApplicationPoints(ctx, leafInterfaceIds)
			if err != nil {
				t.Fatal(err)
			}

			// check our work
			compareInterfacesConnectivityTemplateAssignments(ctAssignments, apToPolicyInfo, t)

			// loop over individual interfaces, checking each
			for interfaceId, expected := range ctAssignments {
				log.Printf("testing GetApplicationPointConnectivityTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				result, err := bpClient.GetApplicationPointConnectivityTemplates(ctx, interfaceId)
				if err != nil {
					t.Fatal(err)
				}

				compareConnectivityTemplateAssignments(expected, result, interfaceId, t)
			}

			log.Printf("testing GetAllApplicationPointsConnectivityTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			all, err := bpClient.GetAllApplicationPointsConnectivityTemplates(ctx)
			if err != nil {
				t.Fatal(err)
			}

			for applicationPointId, expectedCtInfo := range ctAssignments {
				actualCtInfo, ok := all[applicationPointId]
				if !ok {
					t.Fatalf("GetAllApplicationPointsConnectivityTemplates() didn't find a record for %q", applicationPointId)
				}

				compareConnectivityTemplateAssignments(expectedCtInfo, actualCtInfo, applicationPointId, t)
			}

			for _, ctId := range ctIds {
				log.Printf("testing GetApplicationPointsConnectivityTemplatesByCt() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				applicationPoints, err := bpClient.GetApplicationPointsConnectivityTemplatesByCt(ctx, ctId)
				if err != nil {
					t.Fatal(err)
				}

				for applicationPointId, applicationPointCtMap := range applicationPoints {
					if applicationPointCtMap[ctId] != apToPolicyInfo[applicationPointId][ctId] {
						t.Fatalf("application point %s, connectivity template %s, expected: %t actual: %t",
							applicationPointId, ctId, applicationPointCtMap[ctId], apToPolicyInfo[applicationPointId][ctId])
					}
				}
			}
		})
	}
}

func TestSetDelApplicationPointConnectivityTemplates_Errors(t *testing.T) {
	ctx := context.Background()

	ctCount := 5

	type testCase struct {
		apIdx  int   // index of application point ID in our slice of AP IDs. Negative value indicates "use a bogus AP ID"
		ctIdxs []int // indexes of connectivity template IDs in our slice of CT IDs. Negative value indicates "use a bogus CT ID"
	}

	testCases := map[string]testCase{
		"not_bogus_one_CT": {
			apIdx:  1,
			ctIdxs: []int{1},
		},
		"not_bogus_two_CTs": {
			apIdx:  1,
			ctIdxs: []int{3, 2},
		},
		"bogus_AP_one_CT": {
			apIdx:  -1,
			ctIdxs: []int{3},
		},
		"bogus_AP_two_CTs": {
			apIdx:  -1,
			ctIdxs: []int{3, 2},
		},
		"good_AP_one_bogus_CT": {
			apIdx:  3,
			ctIdxs: []int{-1},
		},
		"good_AP_two_bogus_CTs": {
			apIdx:  3,
			ctIdxs: []int{-1, -1},
		},
		"good_AP_blended_CTs_a": {
			apIdx:  2,
			ctIdxs: []int{2, -1, 0},
		},
		"good_AP_blended_CTs_b": {
			apIdx:  2,
			ctIdxs: []int{-1, 2, -1},
		},
		"good_AP_blended_CTs_c": {
			apIdx:  1,
			ctIdxs: []int{2, -1, -1, 1},
		},
		"good_AP_blended_CTs_d": {
			apIdx:  3,
			ctIdxs: []int{-1, 2, 1, -1},
		},
		"good_AP_blended_CTs_e": {
			apIdx:  3,
			ctIdxs: []int{2, -1, 0, -1},
		},
		"good_AP_blended_CTs_f": {
			apIdx:  1,
			ctIdxs: []int{-1, 3, -1, 1},
		},
	}

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			bpClient := testBlueprintC(ctx, t, client.client)

			zones, err := bpClient.GetAllSecurityZones(ctx)
			require.NoError(t, err)
			require.Equal(t, 1, len(zones))

			cts := make([]*ConnectivityTemplate, ctCount)
			ctIds := make([]ObjectId, ctCount)
			ctLabel := randString(5, "hex")
			for i := range ctCount {
				cts[i] = &ConnectivityTemplate{
					Label: ctLabel + fmt.Sprintf("-%d", i),
					Subpolicies: []*ConnectivityTemplatePrimitive{
						{
							Attributes: &ConnectivityTemplatePrimitiveAttributesAttachLogicalLink{
								SecurityZone:       &zones[0].Id,
								Tagged:             true,
								Vlan:               toPtr(Vlan(i + 101)),
								IPv4AddressingType: CtPrimitiveIPv4AddressingTypeNumbered,
							},
						},
					},
				}

				require.NoError(t, cts[i].SetIds())
				require.NoError(t, cts[i].SetUserData())
				require.NoError(t, bpClient.CreateConnectivityTemplate(ctx, cts[i]))
				ctIds[i] = *cts[i].Id
			}

			// Graph query which picks out generic-facing interfaces on leaf switches
			query := new(PathQuery).
				SetBlueprintType(BlueprintTypeStaging).
				SetBlueprintId(bpClient.blueprintId).
				SetClient(bpClient.client).
				Node([]QEEAttribute{
					NodeTypeSystem.QEEAttribute(),
					{"role", QEStringVal("leaf")},
				}).
				Out([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
				Node([]QEEAttribute{
					NodeTypeInterface.QEEAttribute(),
					{"name", QEStringVal("switch_port")},
				}).
				Out([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
				Node([]QEEAttribute{NodeTypeLink.QEEAttribute()}).
				In([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
				Node([]QEEAttribute{NodeTypeInterface.QEEAttribute()}).
				In([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
				Node([]QEEAttribute{
					NodeTypeSystem.QEEAttribute(),
					{"role", QEStringVal("generic")},
				})

			var queryResponse struct {
				Items []struct {
					Interface struct {
						Id ObjectId `json:"id"`
					} `json:"switch_port"`
				} `json:"items"`
			}

			err = query.Do(ctx, &queryResponse)
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(queryResponse.Items), 2)

			// collect the server-facing interfaces of leaf switches
			leafInterfaceIds := make([]ObjectId, len(queryResponse.Items))
			for i, item := range queryResponse.Items {
				leafInterfaceIds[i] = item.Interface.Id
			}

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					// t.Parallel() -- do not use -- all tests use the same interfaces

					var testApId ObjectId
					var errorExpected bool
					if tCase.apIdx >= 0 {
						testApId = leafInterfaceIds[tCase.apIdx]
					} else {
						testApId = ObjectId(randString(6, "hex"))
						errorExpected = true
					}

					testCtIds := make([]ObjectId, len(tCase.ctIdxs))
					for i, idx := range tCase.ctIdxs {
						if idx >= 0 {
							testCtIds[i] = ctIds[idx]
						} else {
							testCtIds[i] = ObjectId(randString(6, "hex"))
							errorExpected = true
						}
					}

					log.Printf("testing SetApplicationPointConnectivityTemplates() error handling against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					setErr := bpClient.SetApplicationPointConnectivityTemplates(ctx, testApId, testCtIds)
					if !errorExpected {
						require.NoError(t, setErr)
					}

					log.Printf("testing DelApplicationPointConnectivityTemplates() error handling against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					delErr := bpClient.DelApplicationPointConnectivityTemplates(ctx, testApId, testCtIds)
					if !errorExpected {
						require.NoError(t, delErr)
						return
					}

					// if we got here, there should be errors to inspect
					for _, err := range []error{setErr, delErr} {
						require.Error(t, err)

						var ace ClientErr
						require.ErrorAs(t, setErr, &ace)
						require.Equal(t, ErrCtAssignmentFailed, ace.Type())
						detail := ace.detail.(*ErrCtAssignmentFailedDetail)

						// collect the bad data we used
						var bogusApIds []ObjectId
						if tCase.apIdx < 0 {
							bogusApIds = []ObjectId{testApId}
						}
						var bogusCtIds []ObjectId
						for i, idx := range tCase.ctIdxs {
							if idx < 0 {
								bogusCtIds = append(bogusCtIds, testCtIds[i])
							}
						}

						require.Equal(t, len(bogusCtIds), len(detail.InvalidConnectivityTemplateIds))
						for _, bogusCtId := range bogusCtIds {
							require.Contains(t, detail.InvalidConnectivityTemplateIds, bogusCtId)
						}

						// bogus CT IDs take precedence over bogus AP IDs, so only check for
						// expected AP IDs when no bogus CT IDs were used in the request.
						if len(bogusCtIds) == 0 {
							require.Equal(t, len(bogusApIds), len(detail.InvalidApplicationPointIds))
							for _, bogusApId := range bogusApIds {
								require.Contains(t, detail.InvalidApplicationPointIds, bogusApId)
							}
						}
					}
				})
			}
		})
	}
}

func TestSetApplicationPointsConnectivityTemplates_Errors(t *testing.T) {
	ctx := context.Background()

	ctCount := 8

	type testCase struct {
		apIdxToCtIdxs map[int][]int
	}

	testCases := map[string]testCase{
		"one_bogus_CT_of_many": {
			apIdxToCtIdxs: map[int][]int{
				0: {0},
				1: {0, 1},
				2: {0, 1, 2},
				3: {0, 1, 2, 3},
				4: {0, 1, 2, 3, 4},
				5: {0, 1, 2, -1, 3, 4, 5},
				6: {0, 1, 2, 3, 4, 5, 6},
				7: {0, 1, 2, 3, 4, 5, 6, 7},
			},
		},
		"several_bogus_CTs_of_many": {
			apIdxToCtIdxs: map[int][]int{
				0: {0},
				1: {0, 1},
				2: {0, 1, 2},
				3: {0, 1, 2, 3},
				4: {0, 1, 2, 3, 4},
				5: {0, 1, -1, -1, 3, 4, 5},
				6: {0, 1, 2, 3, 4, 5, 6},
				7: {0, 1, 2, -1, -1, 3, 4, 5, 6, 7, -1},
			},
		},
		"one_bogus_AP_of_many": {
			apIdxToCtIdxs: map[int][]int{
				0:  {0},
				1:  {0, 1},
				2:  {0, 1, 2},
				3:  {0, 1, 2, 3},
				-1: {0, 1, 2, 3, 4},
				4:  {0, 1, 2, 3, 4},
				5:  {0, 1, 2, 3, 4, 5},
				6:  {0, 1, 2, 3, 4, 5, 6},
				7:  {0, 1, 2, 3, 4, 5, 6, 7},
			},
		},
		"several_bogus_APs_of_many": {
			apIdxToCtIdxs: map[int][]int{
				0:  {0},
				1:  {0, 1},
				2:  {0, 1, 2},
				3:  {0, 1, 2, 3},
				-1: {0, 1, 2, 3, 4},
				4:  {0, 1, 2, 3, 4},
				5:  {0, 1, 2, 3, 4, 5},
				-2: {0, 1, 2, 3, 4, 5},
				6:  {0, 1, 2, 3, 4, 5, 6},
				-3: {0, 1, 2, 3, 4, 5, 6, 7},
				7:  {0, 1, 2, 3, 4, 5, 6, 7},
			},
		},
		"several_bogus_APs_and_CTs": {
			apIdxToCtIdxs: map[int][]int{
				0:  {0},
				1:  {0, 1},
				2:  {0, 1, 2},
				3:  {0, 1, 2, -1, 3},
				-1: {0, 1, 2, 3, 4},
				4:  {0, 1, 2, 3, 4},
				5:  {0, -1, 1, 2, -1, 3, 4, 5},
				-2: {0, 1, -1, 2, 3, 4, 5},
				6:  {0, 1, 2, 3, 4, 5, 6},
				-3: {0, 1, 2, 3, 4, 5, 6, 7},
				7:  {0, 1, 2, 3, 4, 5, 6, 7},
			},
		},
	}

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			bpClient := testBlueprintC(ctx, t, client.client)

			zones, err := bpClient.GetAllSecurityZones(ctx)
			require.NoError(t, err)
			require.Equal(t, 1, len(zones))

			cts := make([]*ConnectivityTemplate, ctCount)
			ctIds := make([]ObjectId, ctCount)
			ctLabel := randString(5, "hex")
			for i := range ctCount {
				cts[i] = &ConnectivityTemplate{
					Label: ctLabel + fmt.Sprintf("-%d", i),
					Subpolicies: []*ConnectivityTemplatePrimitive{
						{
							Attributes: &ConnectivityTemplatePrimitiveAttributesAttachLogicalLink{
								SecurityZone:       &zones[0].Id,
								Tagged:             true,
								Vlan:               toPtr(Vlan(i + 101)),
								IPv4AddressingType: CtPrimitiveIPv4AddressingTypeNumbered,
							},
						},
					},
				}

				require.NoError(t, cts[i].SetIds())
				require.NoError(t, cts[i].SetUserData())
				require.NoError(t, bpClient.CreateConnectivityTemplate(ctx, cts[i]))
				ctIds[i] = *cts[i].Id
			}

			// Graph query which picks out generic-facing interfaces on leaf switches
			query := new(PathQuery).
				SetBlueprintType(BlueprintTypeStaging).
				SetBlueprintId(bpClient.blueprintId).
				SetClient(bpClient.client).
				Node([]QEEAttribute{
					NodeTypeSystem.QEEAttribute(),
					{"role", QEStringVal("leaf")},
				}).
				Out([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
				Node([]QEEAttribute{
					NodeTypeInterface.QEEAttribute(),
					{"name", QEStringVal("switch_port")},
				}).
				Out([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
				Node([]QEEAttribute{NodeTypeLink.QEEAttribute()}).
				In([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
				Node([]QEEAttribute{NodeTypeInterface.QEEAttribute()}).
				In([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
				Node([]QEEAttribute{
					NodeTypeSystem.QEEAttribute(),
					{"role", QEStringVal("generic")},
				})

			var queryResponse struct {
				Items []struct {
					Interface struct {
						Id ObjectId `json:"id"`
					} `json:"switch_port"`
				} `json:"items"`
			}

			err = query.Do(ctx, &queryResponse)
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(queryResponse.Items), 2)

			// collect the server-facing interfaces of leaf switches
			leafInterfaceIds := make([]ObjectId, len(queryResponse.Items))
			for i, item := range queryResponse.Items {
				leafInterfaceIds[i] = item.Interface.Id
			}

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					// t.Parallel() -- do not use -- all tests use the same interfaces

					var bogusApIds []ObjectId
					var bogusCtIds []ObjectId

					setRequest := make(map[ObjectId]map[ObjectId]bool)
					delRequest := make(map[ObjectId]map[ObjectId]bool)
					for apIdx, ctIdxs := range tCase.apIdxToCtIdxs {
						var apId ObjectId
						if apIdx >= 0 {
							apId = leafInterfaceIds[apIdx]
						} else {
							apId = ObjectId("bogus-AP-" + randString(6, "hex"))
							bogusApIds = append(bogusApIds, apId)
						}

						setRequest[apId] = map[ObjectId]bool{}
						delRequest[apId] = map[ObjectId]bool{}
						for _, ctIdx := range ctIdxs {
							if ctIdx >= 0 {
								ctId := ctIds[ctIdx]
								setRequest[apId][ctId] = true
								delRequest[apId][ctId] = false
							} else {
								ctId := ObjectId("bogus-CT-" + randString(6, "hex"))
								setRequest[apId][ctId] = true
								delRequest[apId][ctId] = false
								bogusCtIds = append(bogusCtIds, ctId)
							}
						}

					}

					log.Printf("testing SetApplicationPointsConnectivityTemplates() error handling when setting with bogus values against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					setCtx := context.WithValue(ctx, CtxKeyTestID, tName+"(set)")
					setErr := bpClient.SetApplicationPointsConnectivityTemplates(setCtx, setRequest)
					if len(bogusApIds) == 0 && len(bogusCtIds) == 0 {
						require.NoError(t, setErr)
					}

					log.Printf("testing SetApplicationPointsConnectivityTemplates() error handling when clearing with bogus values against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					delCtx := context.WithValue(ctx, CtxKeyTestID, tName+"(del)")
					delErr := bpClient.SetApplicationPointsConnectivityTemplates(delCtx, delRequest)
					if len(bogusApIds) == 0 && len(bogusCtIds) == 0 {
						require.NoError(t, setErr)
						return
					}

					// if we got here, there should be errors to inspect
					for _, err := range []error{setErr, delErr} {
						require.Error(t, err)

						var ace ClientErr
						require.ErrorAs(t, setErr, &ace)
						require.Equal(t, ErrCtAssignmentFailed, ace.Type())
						detail := ace.detail.(*ErrCtAssignmentFailedDetail)

						require.Equal(t, len(bogusCtIds), len(detail.InvalidConnectivityTemplateIds))
						for _, bogusCtId := range bogusCtIds {
							require.Contains(t, detail.InvalidConnectivityTemplateIds, bogusCtId)
						}

						// bogus CT IDs take precedence over bogus AP IDs, so only check for
						// expected AP IDs when no bogus CT IDs were used in the request.
						if len(bogusCtIds) == 0 {
							require.Equal(t, len(bogusApIds), len(detail.InvalidApplicationPointIds))
							for _, id := range bogusApIds {
								require.Contains(t, detail.InvalidApplicationPointIds, id)
							}
						}
					}
				})
			}
		})
	}
}
