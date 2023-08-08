//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"strings"
	"testing"
)

func TestAssignClearCtToInterface(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintA(ctx, t, client.client)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err := bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		var systems struct {
			Nodes map[ObjectId]struct {
				Label      string     `json:"label"`
				SystemType string     `json:"system_type"`
				Role       systemRole `json:"role"`
			} `json:"nodes"`
		}

		log.Printf("testing GetNodes() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.GetNodes(ctx, NodeTypeSystem, &systems)
		if err != nil {
			t.Fatal(err)
		}

		// create assignments for the leaf switch nodes
		assignments := make(SystemIdToInterfaceMapAssignment)
		var bindings []VnBinding
		for k, v := range systems.Nodes {
			if v.SystemType != "switch" || v.Role != systemRoleLeaf {
				continue
			}
			assignments[k.String()] = "Juniper_vQFX__AOS-7x10-Leaf"
			bindings = append(bindings, VnBinding{SystemId: k})
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

		log.Printf("testing CreateVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		vnId, err := bpClient.CreateVirtualNetwork(ctx, &VirtualNetworkData{
			DhcpService:               false,
			Ipv4Enabled:               false,
			Ipv4Subnet:                nil,
			Ipv6Enabled:               false,
			Ipv6Subnet:                nil,
			Label:                     randString(6, "hex"),
			ReservedVlanId:            nil,
			RouteTarget:               "",
			RtPolicy:                  nil,
			SecurityZoneId:            szId,
			SviIps:                    nil,
			VirtualGatewayIpv4:        nil,
			VirtualGatewayIpv6:        nil,
			VirtualGatewayIpv4Enabled: false,
			VirtualGatewayIpv6Enabled: false,
			VnBindings:                bindings,
			VnId:                      nil,
			VnType:                    VnTypeVxlan,
			VirtualMac:                nil,
		})
		if err != nil {
			t.Fatal(err)
		}

		ct := ConnectivityTemplate{
			Label: randString(5, "hex"),
			Subpolicies: []*ConnectivityTemplatePrimitive{{
				Attributes: &ConnectivityTemplatePrimitiveAttributesAttachSingleVlan{
					VnNodeId: &vnId,
				},
			}},
		}

		err = ct.SetIds()
		if err != nil {
			t.Fatal(err)
		}

		ct.SetUserData()

		log.Printf("testing CreateConnectivityTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.CreateConnectivityTemplate(ctx, &ct)
		if err != nil {
			t.Fatal(err)
		}

		var leaf1Id *ObjectId
		for id, system := range systems.Nodes {
			if system.SystemType != "switch" || system.Role != systemRoleLeaf {
				continue
			}
			if strings.HasSuffix(system.Label, "leaf1") {
				leaf1Id = &id
				break
			}
		}

		if leaf1Id == nil {
			t.Fatal("failed to find leaf 1")
		}

		query := new(PathQuery).
			SetBlueprintType(BlueprintTypeStaging).
			SetBlueprintId(bpClient.blueprintId).
			SetClient(bpClient.client).
			Node([]QEEAttribute{{"id", QEStringVal(leaf1Id.String())}}).
			Out([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeInterface.QEEAttribute(),
				{"if_name", QEStringVal("ae1")},
			}).
			In([]QEEAttribute{RelationshipTypeComposedOf.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeInterface.QEEAttribute(),
				{"name", QEStringVal("n_interface")},
			})

		var queryResponse struct {
			Items []struct {
				Interface struct {
					Id ObjectId `json:"id"`
				} `json:"n_interface"`
			} `json:"items"`
		}

		err = query.Do(ctx, &queryResponse)
		if err != nil {
			t.Fatal(err)
		}
		if len(queryResponse.Items) != 1 {
			t.Fatalf("expected 1 item, got %d items", len(queryResponse.Items))
		}
		interfaceId := queryResponse.Items[0].Interface.Id

		ctsToAssign := []ObjectId{*ct.Id}
		err = bpClient.SetInterfaceConnectivityTemplates(ctx, interfaceId, ctsToAssign)
		if err != nil {
			t.Fatal(err)
		}

		assignedCts, err := bpClient.GetInterfaceConnectivityTemplates(ctx, interfaceId)
		if err != nil {
			t.Fatal(err)
		}

		compareSlices(t, ctsToAssign, assignedCts, "assigned slices do not match intent")

		err = bpClient.DelInterfaceConnectivityTemplates(ctx, interfaceId, ctsToAssign)
		if err != nil {
			t.Fatal(err)
		}

		assignedCts, err = bpClient.GetInterfaceConnectivityTemplates(ctx, interfaceId)
		if err != nil {
			t.Fatal(err)
		}

		if len(assignedCts) != 0 {
			t.Fatalf("expected 0 interfaces assigned to interface, got %d", len(assignedCts))
		}
	}

}
