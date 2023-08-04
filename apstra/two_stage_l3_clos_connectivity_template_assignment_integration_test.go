//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
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
			SecurityZoneId:            "",
			SviIps:                    nil,
			VirtualGatewayIpv4:        nil,
			VirtualGatewayIpv6:        nil,
			VirtualGatewayIpv4Enabled: false,
			VirtualGatewayIpv6Enabled: false,
			VnBindings:                bindings,
			VnId:                      nil,
			VnType:                    VnTypeVlan,
			VirtualMac:                nil,
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Println(vnId)
	}

}
