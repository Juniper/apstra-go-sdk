//go:build integration
// +build integration

package apstra

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestListAllBlueprintIds(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		blueprints, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		result, err := json.Marshal(blueprints)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(string(result))
	}
}

func TestGetAllBlueprintStatus(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getAllBlueprintStatus() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bps, err := client.client.getAllBlueprintStatus(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Println(len(bps))
	}
}

func TestCreateDeleteBlueprint(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		var fabricSettings *FabricSettings
		if rackBasedTemplateFabricAddressingPolicyForbidden().Includes(client.client.apiVersion.String()) {
			// forbidden in the template means we can use this feature in the blueprint
			fabricSettings = &FabricSettings{}

			if !fabricL3MtuForbidden.Check(client.client.apiVersion) {
				fabricL3Mtu := uint16(rand.Intn(550)*2 + 8000) // even number 8000 - 9100
				fabricSettings.FabricL3Mtu = &fabricL3Mtu
				fabricSettings.SpineLeafLinks = toPtr(AddressingSchemeIp46)
				fabricSettings.SpineSuperspineLinks = toPtr(AddressingSchemeIp46)
			}
		}

		req := CreateBlueprintFromTemplateRequest{
			RefDesign:      RefDesignTwoStageL3Clos,
			Label:          randString(10, "hex"),
			TemplateId:     "L2_Virtual_EVPN",
			FabricSettings: fabricSettings,
		}

		log.Printf("testing createBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.CreateBlueprintFromTemplate(ctx, &req)
		if err != nil {
			t.Fatal(err)
		}

		bp, err := client.client.GetBlueprint(ctx, id)
		if err != nil {
			t.Fatal(err)
		}

		if id != bp.Id {
			t.Fatalf("expected id %q, got %q", id, bp.Id)
		}

		if req.Label != bp.Label {
			t.Fatalf("expected label %q, got %q", req.Label, bp.Label)
		}

		bpClient, err := client.client.NewTwoStageL3ClosClient(ctx, id)
		if err != nil {
			t.Fatal(err)
		}

		if req.FabricSettings != nil && req.FabricSettings.FabricL3Mtu != nil {
			fap, err := bpClient.GetFabricAddressingPolicy(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if *req.FabricSettings.FabricL3Mtu != *fap.FabricL3Mtu {
				t.Fatalf("expected fabric MTU %d, got %d", *req.FabricSettings.FabricL3Mtu, *fap.FabricL3Mtu)
			}
		}

		log.Printf("got id '%s', deleting blueprint...\n", id)
		log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteBlueprint(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetPatchGetPatchNode(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintA(ctx, t, client.client)

		type metadataNode struct {
			Tags         interface{} `json:"tags,omitempty"`
			PropertySet  interface{} `json:"property_set,omitempty"`
			Label        string      `json:"label,omitempty"`
			UserIp       interface{} `json:"user_ip,omitempty"`
			TemplateJson interface{} `json:"template_json,omitempty"`
			Design       string      `json:"design,omitempty"`
			User         interface{} `json:"user,omitempty"`
			Type         string      `json:"type,omitempty"`
			Id           ObjectId    `json:"id,omitempty"`
		}

		type nodes struct {
			Nodes map[string]metadataNode `json:"nodes"`
		}
		var nodesA, nodesB nodes

		// fetch all metadata nodes into nodesA
		log.Printf("testing getNodes() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.GetNodes(ctx, NodeTypeMetadata, &nodesA)
		if err != nil {
			t.Fatal()
		}

		// sanity check
		if len(nodesA.Nodes) != 1 {
			t.Fatalf("not expecting %d '%s' nodes", len(nodesA.Nodes), NodeTypeMetadata)
		}

		newName := randString(10, "hex")
		// loop should run just once (len check above)
		for idA, nodeA := range nodesA.Nodes {
			log.Printf("node id: %s ; label: %s\n", idA, nodeA.Label)

			// change name to newName
			req := metadataNode{Label: newName}
			resp := &metadataNode{}
			log.Printf("testing patchNode(%s) against %s %s (%s)", bpClient.Id(), client.clientType, clientName, client.client.ApiVersion())
			err := bpClient.PatchNode(ctx, nodeA.Id, req, resp)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Label != newName {
				t.Fatalf("expected new blueprint name %q, got %q", newName, resp.Label)
			}
			log.Printf("response indicates name changed '%s' -> '%s'", nodeA.Label, resp.Label)

			// fetch changed node(s) (still expecting one) into nodesB
			log.Printf("testing getNodes(%s) against %s %s (%s)", bpClient.Id(), client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.GetNodes(ctx, NodeTypeMetadata, &nodesB)
			if err != nil {
				t.Fatal()
			}
			for idB, nodeB := range nodesB.Nodes {
				log.Printf("node id: %s ; label: %s\n", idB, nodeB.Label)
				if nodeB.Label != newName {
					t.Fatalf("expected new blueprint name %q, got %q", newName, nodeB.Label)
				}
			}
		}
	}
}

func TestGetNodes(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintB(ctx, t, client.client)

		type node struct {
			Id         ObjectId `json:"id"`
			Label      string   `json:"label"`
			SystemType string   `json:"system_type"`
		}
		equal := func(a, b node) bool {
			return a.Id == b.Id &&
				a.Label == b.Label &&
				a.SystemType == b.SystemType
		}

		var response struct {
			Nodes map[ObjectId]node `json:"nodes"`
		}
		log.Printf("testing GetNodes() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.Client().GetNodes(ctx, bpClient.Id(), NodeTypeSystem, &response)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("got %d nodes. Fetch each one...", len(response.Nodes))
		var nodeB node
		for id, nodeA := range response.Nodes {
			log.Printf("testing GetNode() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.Client().GetNode(ctx, bpClient.Id(), id, &nodeB)
			if err != nil {
				t.Fatal()
			}
			if !equal(nodeA, nodeB) {
				t.Fatalf("nodes don't match:\n%v\n%v", nodeA, nodeB)
			}
		}
	}
}

func TestPatchNodes(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintB(ctx, t, client.client)

		type node struct {
			Id         ObjectId `json:"id"`
			Label      string   `json:"label"`
			SystemType string   `json:"system_type,omitempty"`
		}

		var getResponse struct {
			Nodes map[ObjectId]node `json:"nodes"`
		}
		log.Printf("testing GetNodes() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.Client().GetNodes(ctx, bpClient.Id(), NodeTypeSystem, &getResponse)
		if err != nil {
			t.Fatal(err)
		}

		var patch []interface{}
		for k, v := range getResponse.Nodes {
			if v.SystemType == "server" {
				patch = append(patch, node{
					Id:    k,
					Label: randString(5, "hex"),
				})
			}
		}

		err = client.client.PatchNodes(ctx, bpClient.Id(), patch)
		if err != nil {
			t.Fatal(err)
		}

		for _, n := range patch {
			var result node
			err = client.client.GetNode(ctx, bpClient.Id(), n.(node).Id, &result)
			if err != nil {
				t.Fatal(err)
			}

			if n.(node).Label != result.Label {
				t.Fatalf("patch expected label %s, got label %s", n.(node).Label, result.Label)
			}
		}
	}
}

func TestCreateDeleteEvpnBlueprint(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	type testCase struct {
		req                CreateBlueprintFromTemplateRequest
		versionConstraints version.Constraints
	}

	testCases := map[string]testCase{
		"simple": {
			req: CreateBlueprintFromTemplateRequest{
				RefDesign:  RefDesignTwoStageL3Clos,
				Label:      randString(5, "hex"),
				TemplateId: "L2_Virtual_EVPN",
			},
		},
		"4.1.1_and_later": {
			versionConstraints: geApstra411,
			req: CreateBlueprintFromTemplateRequest{
				RefDesign:  RefDesignTwoStageL3Clos,
				Label:      randString(5, "hex"),
				TemplateId: "L2_Virtual_EVPN",
				FabricSettings: &FabricSettings{
					SpineLeafLinks:       toPtr(AddressingSchemeIp46),
					SpineSuperspineLinks: toPtr(AddressingSchemeIp46),
				},
			},
		},
		"4.2.0_specific_test": {
			versionConstraints: geApstra420,
			req: CreateBlueprintFromTemplateRequest{
				RefDesign:  RefDesignTwoStageL3Clos,
				Label:      randString(5, "hex"),
				TemplateId: "L2_Virtual_EVPN",
				FabricSettings: &FabricSettings{
					SpineLeafLinks:       toPtr(AddressingSchemeIp46),
					SpineSuperspineLinks: toPtr(AddressingSchemeIp46),
					FabricL3Mtu:          toPtr(uint16(9178)),
				},
			},
		},
		"lots_of_values": {
			versionConstraints: geApstra420,
			req: CreateBlueprintFromTemplateRequest{
				RefDesign:  RefDesignTwoStageL3Clos,
				Label:      randString(5, "hex"),
				TemplateId: "L2_Virtual_EVPN",
				FabricSettings: &FabricSettings{
					JunosEvpnDuplicateMacRecoveryTime:     toPtr(uint16(16)),
					MaxExternalRoutes:                     toPtr(uint32(239832)),
					EsiMacMsb:                             toPtr(uint8(32)),
					JunosGracefulRestart:                  &FeatureSwitchEnumDisabled,
					OptimiseSzFootprint:                   &FeatureSwitchEnumDisabled,
					JunosEvpnRoutingInstanceVlanAware:     &FeatureSwitchEnumDisabled,
					EvpnGenerateType5HostRoutes:           &FeatureSwitchEnumDisabled,
					MaxFabricRoutes:                       toPtr(uint32(84231)),
					MaxMlagRoutes:                         toPtr(uint32(76112)),
					JunosExOverlayEcmp:                    &FeatureSwitchEnumDisabled,
					DefaultSviL3Mtu:                       toPtr(uint16(9100)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &FeatureSwitchEnumDisabled,
					FabricL3Mtu:                           toPtr(uint16(9178)),
					Ipv6Enabled:                           toPtr(false), // do not enable because it's a one-way trip
					ExternalRouterMtu:                     toPtr(uint16(9100)),
					MaxEvpnRoutes:                         toPtr(uint32(92342)),
					AntiAffinityPolicy: &AntiAffinityPolicy{
						Algorithm:                AlgorithmHeuristic,
						MaxLinksPerPort:          2,
						MaxLinksPerSlot:          2,
						MaxPerSystemLinksPerPort: 2,
						MaxPerSystemLinksPerSlot: 2,
						Mode:                     AntiAffinityModeEnabledLoose,
					},
					SpineLeafLinks:       toPtr(AddressingSchemeIp4),
					SpineSuperspineLinks: toPtr(AddressingSchemeIp4),
				},
			},
		},
		"different_values": {
			versionConstraints: geApstra420,
			req: CreateBlueprintFromTemplateRequest{
				RefDesign:  RefDesignTwoStageL3Clos,
				Label:      randString(5, "hex"),
				TemplateId: "L2_Virtual_EVPN",
				FabricSettings: &FabricSettings{
					JunosEvpnDuplicateMacRecoveryTime:     toPtr(uint16(14)),
					MaxExternalRoutes:                     toPtr(uint32(233832)),
					EsiMacMsb:                             toPtr(uint8(50)),
					JunosGracefulRestart:                  &FeatureSwitchEnumEnabled,
					OptimiseSzFootprint:                   &FeatureSwitchEnumEnabled,
					JunosEvpnRoutingInstanceVlanAware:     &FeatureSwitchEnumEnabled,
					EvpnGenerateType5HostRoutes:           &FeatureSwitchEnumEnabled,
					MaxFabricRoutes:                       toPtr(uint32(82231)),
					MaxMlagRoutes:                         toPtr(uint32(74112)),
					JunosExOverlayEcmp:                    &FeatureSwitchEnumEnabled,
					DefaultSviL3Mtu:                       toPtr(uint16(9070)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &FeatureSwitchEnumEnabled,
					FabricL3Mtu:                           toPtr(uint16(9172)),
					Ipv6Enabled:                           toPtr(false), // do not enable because it's a one-way trip
					ExternalRouterMtu:                     toPtr(uint16(9080)),
					MaxEvpnRoutes:                         toPtr(uint32(91342)),
					AntiAffinityPolicy: &AntiAffinityPolicy{
						Algorithm:                AlgorithmHeuristic,
						MaxLinksPerPort:          4,
						MaxLinksPerSlot:          4,
						MaxPerSystemLinksPerPort: 4,
						MaxPerSystemLinksPerSlot: 4,
						Mode:                     AntiAffinityModeEnabledLoose,
					},
					SpineLeafLinks:       toPtr(AddressingSchemeIp46),
					SpineSuperspineLinks: toPtr(AddressingSchemeIp46),
				},
			},
		},
	}

	fetchFabricAddressingScheme := func(t testing.TB, client *TwoStageL3ClosClient) (AddressingScheme, AddressingScheme) {
		t.Helper()

		query := new(PathQuery).
			SetClient(client.client).
			SetBlueprintId(client.blueprintId)
		if version.MustConstraints(version.NewConstraint(">=" + apstra421)).Check(client.client.apiVersion) {
			query.Node([]QEEAttribute{
				NodeTypeFabricPolicy.QEEAttribute(),
				{Key: "name", Value: QEStringVal("node")},
			})
		} else {
			query.Node([]QEEAttribute{
				NodeTypeFabricAddressingPolicy.QEEAttribute(),
				{Key: "name", Value: QEStringVal("node")},
			})
		}

		var queryResponse struct {
			Items []struct {
				Node struct {
					SpineLeafLinks       addressingScheme `json:"spine_leaf_links"`
					SpineSuperspineLinks addressingScheme `json:"spine_superspine_links"`
				} `json:"node"`
			} `json:"items"`
		}

		err := query.Do(ctx, &queryResponse)
		require.NoError(t, err)

		if len(queryResponse.Items) != 1 {
			t.Fatalf("got %d responses when querying for fabric addressing", len(queryResponse.Items))
		}

		spineLeaf, err := queryResponse.Items[0].Node.SpineLeafLinks.parse()
		require.NoError(t, err)
		spineSuperspine, err := queryResponse.Items[0].Node.SpineLeafLinks.parse()
		require.NoError(t, err)

		return AddressingScheme(spineLeaf), AddressingScheme(spineSuperspine)
	}

	for tName, tCase := range testCases {
		tName, tCase := tName, tCase

		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			for clientName, client := range clients {
				if tCase.versionConstraints != nil && !tCase.versionConstraints.Check(client.client.apiVersion) {
					t.Skipf("skipping test case %q with Apstra %s due to version constraint %q", tName, client.client.apiVersion, tCase.versionConstraints)
				}

				t.Logf("testing CreateBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				id, err := client.client.CreateBlueprintFromTemplate(ctx, &tCase.req)
				require.NoError(t, err)

				t.Logf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				bpClient, err := client.client.NewTwoStageL3ClosClient(ctx, id)
				require.NoError(t, err)

				if tCase.req.FabricSettings != nil {
					t.Logf("testing GetFabricSettings() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					fabricSettings, err := bpClient.GetFabricSettings(ctx)
					require.NoError(t, err)

					t.Logf("comparing create-time vs. fetched blueprint fabric settings against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					compareFabricSettings(t, *tCase.req.FabricSettings, *fabricSettings)

					if tCase.req.FabricSettings.SpineLeafLinks != nil || tCase.req.FabricSettings.SpineSuperspineLinks != nil {
						spineLeaf, spineSuperspine := fetchFabricAddressingScheme(t, bpClient)

						if tCase.req.FabricSettings.SpineLeafLinks != nil && *tCase.req.FabricSettings.SpineLeafLinks != spineLeaf {
							t.Fatalf("expected spine leaf addressing: %q, got %q", *tCase.req.FabricSettings.SpineLeafLinks, spineLeaf)
						}

						if tCase.req.FabricSettings.SpineLeafLinks != nil && *tCase.req.FabricSettings.SpineLeafLinks != spineLeaf {
							t.Fatalf("expected spine superspine addressing: %q, got %q", *tCase.req.FabricSettings.SpineSuperspineLinks, spineSuperspine)
						}
					}
				}

				t.Logf("testing DeleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				err = client.client.DeleteBlueprint(ctx, id)
				require.NoError(t, err)
			}
		})
	}
}

func TestCreateDeleteIpFabricBlueprint(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	type testCase struct {
		req                CreateBlueprintFromTemplateRequest
		versionConstraints version.Constraints
	}

	testCases := map[string]testCase{
		"simple": {
			req: CreateBlueprintFromTemplateRequest{
				RefDesign:  RefDesignTwoStageL3Clos,
				Label:      randString(5, "hex"),
				TemplateId: "L2_Virtual",
			},
		},
		"4.1.1_and_later": {
			versionConstraints: geApstra411,
			req: CreateBlueprintFromTemplateRequest{
				RefDesign:  RefDesignTwoStageL3Clos,
				Label:      randString(5, "hex"),
				TemplateId: "L2_Virtual",
				FabricSettings: &FabricSettings{
					SpineLeafLinks:       toPtr(AddressingSchemeIp46),
					SpineSuperspineLinks: toPtr(AddressingSchemeIp46),
				},
			},
		},
		"4.2.0_specific_test": {
			versionConstraints: geApstra420,
			req: CreateBlueprintFromTemplateRequest{
				RefDesign:  RefDesignTwoStageL3Clos,
				Label:      randString(5, "hex"),
				TemplateId: "L2_Virtual",
				FabricSettings: &FabricSettings{
					SpineLeafLinks:       toPtr(AddressingSchemeIp46),
					SpineSuperspineLinks: toPtr(AddressingSchemeIp46),
					FabricL3Mtu:          toPtr(uint16(9178)),
				},
			},
		},
		"lots_of_values": {
			versionConstraints: geApstra420,
			req: CreateBlueprintFromTemplateRequest{
				RefDesign:  RefDesignTwoStageL3Clos,
				Label:      randString(5, "hex"),
				TemplateId: "L2_Virtual",
				FabricSettings: &FabricSettings{
					JunosEvpnDuplicateMacRecoveryTime:     toPtr(uint16(16)),
					MaxExternalRoutes:                     toPtr(uint32(239832)),
					EsiMacMsb:                             toPtr(uint8(32)),
					JunosGracefulRestart:                  &FeatureSwitchEnumDisabled,
					OptimiseSzFootprint:                   &FeatureSwitchEnumDisabled,
					JunosEvpnRoutingInstanceVlanAware:     &FeatureSwitchEnumDisabled,
					EvpnGenerateType5HostRoutes:           &FeatureSwitchEnumDisabled,
					MaxFabricRoutes:                       toPtr(uint32(84231)),
					MaxMlagRoutes:                         toPtr(uint32(76112)),
					JunosExOverlayEcmp:                    &FeatureSwitchEnumDisabled,
					DefaultSviL3Mtu:                       toPtr(uint16(9100)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &FeatureSwitchEnumDisabled,
					FabricL3Mtu:                           toPtr(uint16(9178)),
					Ipv6Enabled:                           toPtr(false), // do not enable because it's a one-way trip
					ExternalRouterMtu:                     toPtr(uint16(9100)),
					MaxEvpnRoutes:                         toPtr(uint32(92342)),
					AntiAffinityPolicy: &AntiAffinityPolicy{
						Algorithm:                AlgorithmHeuristic,
						MaxLinksPerPort:          2,
						MaxLinksPerSlot:          2,
						MaxPerSystemLinksPerPort: 2,
						MaxPerSystemLinksPerSlot: 2,
						Mode:                     AntiAffinityModeEnabledLoose,
					},
					SpineLeafLinks:       toPtr(AddressingSchemeIp4),
					SpineSuperspineLinks: toPtr(AddressingSchemeIp4),
				},
			},
		},
		"different_values": {
			versionConstraints: geApstra420,
			req: CreateBlueprintFromTemplateRequest{
				RefDesign:  RefDesignTwoStageL3Clos,
				Label:      randString(5, "hex"),
				TemplateId: "L2_Virtual",
				FabricSettings: &FabricSettings{
					JunosEvpnDuplicateMacRecoveryTime:     toPtr(uint16(14)),
					MaxExternalRoutes:                     toPtr(uint32(233832)),
					EsiMacMsb:                             toPtr(uint8(50)),
					JunosGracefulRestart:                  &FeatureSwitchEnumEnabled,
					OptimiseSzFootprint:                   &FeatureSwitchEnumEnabled,
					JunosEvpnRoutingInstanceVlanAware:     &FeatureSwitchEnumEnabled,
					EvpnGenerateType5HostRoutes:           &FeatureSwitchEnumEnabled,
					MaxFabricRoutes:                       toPtr(uint32(82231)),
					MaxMlagRoutes:                         toPtr(uint32(74112)),
					JunosExOverlayEcmp:                    &FeatureSwitchEnumEnabled,
					DefaultSviL3Mtu:                       toPtr(uint16(9070)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &FeatureSwitchEnumEnabled,
					FabricL3Mtu:                           toPtr(uint16(9172)),
					Ipv6Enabled:                           toPtr(false), // do not enable because it's a one-way trip
					ExternalRouterMtu:                     toPtr(uint16(9080)),
					MaxEvpnRoutes:                         toPtr(uint32(91342)),
					AntiAffinityPolicy: &AntiAffinityPolicy{
						Algorithm:                AlgorithmHeuristic,
						MaxLinksPerPort:          4,
						MaxLinksPerSlot:          4,
						MaxPerSystemLinksPerPort: 4,
						MaxPerSystemLinksPerSlot: 4,
						Mode:                     AntiAffinityModeEnabledLoose,
					},
					SpineLeafLinks:       toPtr(AddressingSchemeIp46),
					SpineSuperspineLinks: toPtr(AddressingSchemeIp46),
				},
			},
		},
	}

	fetchFabricAddressingScheme := func(t testing.TB, client *TwoStageL3ClosClient) (AddressingScheme, AddressingScheme) {
		t.Helper()

		query := new(PathQuery).
			SetClient(client.client).
			SetBlueprintId(client.blueprintId)
		if version.MustConstraints(version.NewConstraint(">=" + apstra421)).Check(client.client.apiVersion) {
			query.Node([]QEEAttribute{
				NodeTypeFabricPolicy.QEEAttribute(),
				{Key: "name", Value: QEStringVal("node")},
			})
		} else {
			query.Node([]QEEAttribute{
				NodeTypeFabricAddressingPolicy.QEEAttribute(),
				{Key: "name", Value: QEStringVal("node")},
			})
		}

		var queryResponse struct {
			Items []struct {
				Node struct {
					SpineLeafLinks       addressingScheme `json:"spine_leaf_links"`
					SpineSuperspineLinks addressingScheme `json:"spine_superspine_links"`
				} `json:"node"`
			} `json:"items"`
		}

		err := query.Do(ctx, &queryResponse)
		require.NoError(t, err)

		if len(queryResponse.Items) != 1 {
			t.Fatalf("got %d responses when querying for fabric addressing", len(queryResponse.Items))
		}

		spineLeaf, err := queryResponse.Items[0].Node.SpineLeafLinks.parse()
		require.NoError(t, err)
		spineSuperspine, err := queryResponse.Items[0].Node.SpineLeafLinks.parse()
		require.NoError(t, err)

		return AddressingScheme(spineLeaf), AddressingScheme(spineSuperspine)
	}

	for tName, tCase := range testCases {
		tName, tCase := tName, tCase

		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			for clientName, client := range clients {
				if tCase.versionConstraints != nil && !tCase.versionConstraints.Check(client.client.apiVersion) {
					t.Skipf("skipping test case %q with Apstra %s due to version constraint %q", tName, client.client.apiVersion, tCase.versionConstraints)
				}

				t.Logf("testing CreateBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				id, err := client.client.CreateBlueprintFromTemplate(ctx, &tCase.req)
				require.NoError(t, err)

				t.Logf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				bpClient, err := client.client.NewTwoStageL3ClosClient(ctx, id)
				require.NoError(t, err)

				if tCase.req.FabricSettings != nil {
					t.Logf("testing GetFabricSettings() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					fabricSettings, err := bpClient.GetFabricSettings(ctx)
					require.NoError(t, err)

					t.Logf("comparing create-time vs. fetched blueprint fabric settings against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					compareFabricSettings(t, *tCase.req.FabricSettings, *fabricSettings)

					if tCase.req.FabricSettings.SpineLeafLinks != nil || tCase.req.FabricSettings.SpineSuperspineLinks != nil {
						spineLeaf, spineSuperspine := fetchFabricAddressingScheme(t, bpClient)

						if tCase.req.FabricSettings.SpineLeafLinks != nil && *tCase.req.FabricSettings.SpineLeafLinks != spineLeaf {
							t.Fatalf("expected spine leaf addressing: %q, got %q", *tCase.req.FabricSettings.SpineLeafLinks, spineLeaf)
						}

						if tCase.req.FabricSettings.SpineLeafLinks != nil && *tCase.req.FabricSettings.SpineLeafLinks != spineLeaf {
							t.Fatalf("expected spine superspine addressing: %q, got %q", *tCase.req.FabricSettings.SpineSuperspineLinks, spineSuperspine)
						}
					}
				}

				t.Logf("testing DeleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				err = client.client.DeleteBlueprint(ctx, id)
				require.NoError(t, err)
			}
		})
	}
}

func TestCreateDeleteBlueprintWithRoutingLimits(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	blueprintRequest := CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignTwoStageL3Clos,
		Label:      randString(5, "hex"),
		TemplateId: "L2_Virtual",
	}

	type testCase struct {
		name string
		//versionConstraints   version.Constraints
		fabricSettings FabricSettings
	}

	testCases := []testCase{
		{
			name:           "create_with_defaults",
			fabricSettings: FabricSettings{},
		},
		{
			name: "create_with_zeros",
			fabricSettings: FabricSettings{
				MaxEvpnRoutes:     toPtr(uint32(0)),
				MaxExternalRoutes: toPtr(uint32(0)),
				MaxFabricRoutes:   toPtr(uint32(0)),
				MaxMlagRoutes:     toPtr(uint32(0)),
			},
		},
		{
			name: "create_with_values",
			fabricSettings: FabricSettings{
				MaxEvpnRoutes:     toPtr(uint32(20001)),
				MaxExternalRoutes: toPtr(uint32(20002)),
				MaxFabricRoutes:   toPtr(uint32(20003)),
				MaxMlagRoutes:     toPtr(uint32(20004)),
			},
		},
	}

	for clientName, client := range clients {
		t.Run(client.client.apiVersion.String(), func(t *testing.T) {
			if !geApstra421.Check(client.client.apiVersion) {
				t.Skipf("skipping Apstra %s client due to version constraint", client.client.apiVersion)
			}

			for _, tCase := range testCases {
				clientName, client := clientName, client
				tCase := tCase

				t.Run(tCase.name, func(t *testing.T) {

					bpr := blueprintRequest
					bpr.FabricSettings = &tCase.fabricSettings
					t.Logf("testing CreateBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					id, err := client.client.CreateBlueprintFromTemplate(ctx, &bpr)
					require.NoError(t, err)

					t.Logf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					bpClient, err := client.client.NewTwoStageL3ClosClient(ctx, id)
					require.NoError(t, err)

					t.Logf("testing GetFabricSettings() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					fabricSettings, err := bpClient.GetFabricSettings(ctx)
					require.NoError(t, err)

					t.Logf("comparing create-time vs. fetched blueprint fabric settings against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					compareFabricSettings(t, tCase.fabricSettings, *fabricSettings)

					t.Logf("testing DeleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
					err = client.client.DeleteBlueprint(ctx, id)
					require.NoError(t, err)
				})
			}
		})
	}
}
