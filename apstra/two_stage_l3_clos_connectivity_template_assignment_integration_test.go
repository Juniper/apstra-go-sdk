//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
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
		bpClient, bpDel := testBlueprintC(ctx, t, client.client)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err := bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

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
				VnType:         VnTypeVxlan,
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
			//Node([]QEEAttribute{{"id", QEStringVal(leaf1Id.String())}}).
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
	}
}
