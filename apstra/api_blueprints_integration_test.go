// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestListAllBlueprintIds(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			blueprints, err := client.Client.ListAllBlueprintIds(ctx)
			require.NoError(t, err)

			result, err := json.Marshal(blueprints)
			require.NoError(t, err)

			log.Println(string(result))
		})
	}
}

func TestGetAllBlueprintStatus(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bps, err := client.Client.GetAllBlueprintStatus(ctx)
			require.NoError(t, err)

			log.Println(len(bps))
		})
	}
}

func TestCreateDeleteBlueprint(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			req := apstra.CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      testutils.RandString(10, "hex"),
				TemplateId: "L2_Virtual_EVPN",
				FabricSettings: &apstra.FabricSettings{
					FabricL3Mtu:          testutils.ToPtr(uint16(rand.Intn(50)*2 + 9100)),
					SpineLeafLinks:       testutils.ToPtr(apstra.AddressingSchemeIp46),
					SpineSuperspineLinks: testutils.ToPtr(apstra.AddressingSchemeIp46),
				},
			}

			id, err := client.Client.CreateBlueprintFromTemplate(ctx, &req)
			require.NoError(t, err)

			bp, err := client.Client.GetBlueprint(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id, bp.Id)
			require.Equal(t, req.Label, bp.Label)

			bpClient, err := client.Client.NewTwoStageL3ClosClient(ctx, id)
			require.NoError(t, err)

			if req.FabricSettings != nil && req.FabricSettings.FabricL3Mtu != nil {
				fap, err := bpClient.GetFabricSettings(ctx)
				require.NoError(t, err)
				require.Equal(t, *req.FabricSettings.FabricL3Mtu, *fap.FabricL3Mtu)
			}

			log.Printf("got id '%s', deleting blueprint...\n", id)
			err = client.Client.DeleteBlueprint(ctx, id)
			require.NoError(t, err)
		})
	}
}

func TestGetPatchGetPatchNode(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bpClient := dctestobj.TestBlueprintA(t, ctx, client.Client)

			type metadataNode struct {
				Tags         interface{}     `json:"tags,omitempty"`
				PropertySet  interface{}     `json:"property_set,omitempty"`
				Label        string          `json:"label,omitempty"`
				UserIp       interface{}     `json:"user_ip,omitempty"`
				TemplateJson interface{}     `json:"template_json,omitempty"`
				Design       string          `json:"design,omitempty"`
				User         interface{}     `json:"user,omitempty"`
				Type         string          `json:"type,omitempty"`
				Id           apstra.ObjectId `json:"id,omitempty"`
			}

			type nodes struct {
				Nodes map[string]metadataNode `json:"nodes"`
			}
			var nodesA, nodesB nodes

			// fetch all metadata nodes into nodesA
			require.NoError(t, bpClient.GetNodes(ctx, apstra.NodeTypeMetadata, &nodesA))

			// sanity check
			require.Equal(t, 1, len(nodesA.Nodes))

			newName := testutils.RandString(10, "hex")
			// loop should run just once (len check above)
			for idA, nodeA := range nodesA.Nodes {
				log.Printf("node id: %s ; label: %s\n", idA, nodeA.Label)

				// change name to newName
				req := metadataNode{Label: newName}
				var resp metadataNode
				if compatibility.PatchNodeSupportsUnsafeArg.Check(client.APIVersion()) {
					var ace apstra.ClientErr
					err := bpClient.PatchNode(ctx, nodeA.Id, req, &resp)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), apstra.ErrUnsafePatchProhibited)

					log.Printf("Apstra %s complained that this patch attempt is unsafe. Good!", client.Client.ApiVersion())

					require.NoError(t, bpClient.PatchNodeUnsafe(ctx, nodeA.Id, req, &resp))
				} else {
					require.NoError(t, bpClient.PatchNode(ctx, nodeA.Id, req, &resp))
				}
				if resp.Label != newName {
					t.Fatalf("expected new blueprint name %q, got %q", newName, resp.Label)
				}
				log.Printf("response indicates name changed '%s' -> '%s'", nodeA.Label, resp.Label)

				// fetch changed node(s) (still expecting one) into nodesB
				require.NoError(t, bpClient.GetNodes(ctx, apstra.NodeTypeMetadata, &nodesB))
				for idB, nodeB := range nodesB.Nodes {
					log.Printf("node id: %s ; label: %s\n", idB, nodeB.Label)
					require.Equalf(t, nodeB.Label, newName, "expected new blueprint name %q, got %q", newName, nodeB.Label)
				}
			}
		})
	}
}

func TestGetDcNodes(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bpClient := dctestobj.TestBlueprintB(t, ctx, client.Client)

			type node struct {
				Id         apstra.ObjectId `json:"id"`
				Label      string          `json:"label"`
				SystemType string          `json:"system_type"`
			}

			var response struct {
				Nodes map[apstra.ObjectId]node `json:"nodes"`
			}
			err := bpClient.Client().GetNodes(ctx, bpClient.Id(), apstra.NodeTypeSystem, &response)
			require.NoError(t, err)

			log.Printf("got %d nodes. Fetch each one...", len(response.Nodes))
			var nodeB node
			for id, nodeA := range response.Nodes {
				err = bpClient.Client().GetNode(ctx, bpClient.Id(), id, &nodeB)
				require.NoError(t, err)
				require.Equal(t, nodeB, nodeA)
			}
		})
	}
}

func TestPatchNodes(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bpClient := dctestobj.TestBlueprintB(t, ctx, client.Client)

			type node struct {
				Id         apstra.ObjectId `json:"id"`
				Label      string          `json:"label"`
				SystemType string          `json:"system_type,omitempty"`
			}

			var getResponse struct {
				Nodes map[apstra.ObjectId]node `json:"nodes"`
			}
			err := bpClient.Client().GetNodes(ctx, bpClient.Id(), apstra.NodeTypeSystem, &getResponse)
			require.NoError(t, err)

			var patch []interface{}
			for k, v := range getResponse.Nodes {
				if v.SystemType == "server" {
					patch = append(patch, node{
						Id:    k,
						Label: testutils.RandString(5, "hex"),
					})
				}
			}

			require.NoError(t, client.Client.PatchNodes(ctx, bpClient.Id(), patch))

			for _, n := range patch {
				var result node
				require.NoError(t, client.Client.GetNode(ctx, bpClient.Id(), n.(node).Id, &result))
				require.Equal(t, n.(node).Label, result.Label)
			}
		})
	}
}

func TestCreateDeleteEvpnBlueprint(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		req        apstra.CreateBlueprintFromTemplateRequest
		constraint *compatibility.Constraint
	}

	testCases := map[string]testCase{
		"simple": {
			req: apstra.CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      testutils.RandString(5, "hex"),
				TemplateId: "L2_Virtual_EVPN",
			},
		},
		"4.1.1_and_later": {
			req: apstra.CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      testutils.RandString(5, "hex"),
				TemplateId: "L2_Virtual_EVPN",
				FabricSettings: &apstra.FabricSettings{
					SpineLeafLinks:       testutils.ToPtr(apstra.AddressingSchemeIp46),
					SpineSuperspineLinks: testutils.ToPtr(apstra.AddressingSchemeIp46),
				},
			},
		},
		"4.2.0_specific_test": {
			req: apstra.CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      testutils.RandString(5, "hex"),
				TemplateId: "L2_Virtual_EVPN",
				FabricSettings: &apstra.FabricSettings{
					SpineLeafLinks:       testutils.ToPtr(apstra.AddressingSchemeIp46),
					SpineSuperspineLinks: testutils.ToPtr(apstra.AddressingSchemeIp46),
					FabricL3Mtu:          testutils.ToPtr(uint16(9178)),
				},
			},
		},
		"lots_of_values": {
			req: apstra.CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      testutils.RandString(5, "hex"),
				TemplateId: "L2_Virtual_EVPN",
				FabricSettings: &apstra.FabricSettings{
					JunosEvpnDuplicateMacRecoveryTime:     testutils.ToPtr(uint16(16)),
					MaxExternalRoutes:                     testutils.ToPtr(uint32(239832)),
					EsiMacMsb:                             testutils.ToPtr(uint8(32)),
					JunosGracefulRestart:                  &enum.FeatureSwitchDisabled,
					OptimiseSzFootprint:                   &enum.FeatureSwitchDisabled,
					JunosEvpnRoutingInstanceVlanAware:     &enum.FeatureSwitchDisabled,
					EvpnGenerateType5HostRoutes:           &enum.FeatureSwitchDisabled,
					MaxFabricRoutes:                       testutils.ToPtr(uint32(84231)),
					MaxMlagRoutes:                         testutils.ToPtr(uint32(76112)),
					JunosExOverlayEcmp:                    &enum.FeatureSwitchDisabled,
					DefaultSviL3Mtu:                       testutils.ToPtr(uint16(9100)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &enum.FeatureSwitchDisabled,
					FabricL3Mtu:                           testutils.ToPtr(uint16(9178)),
					Ipv6Enabled:                           testutils.ToPtr(false), // do not enable because it's a one-way trip
					ExternalRouterMtu:                     testutils.ToPtr(uint16(9100)),
					MaxEvpnRoutes:                         testutils.ToPtr(uint32(92342)),
					AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
						Algorithm:                apstra.AlgorithmHeuristic,
						MaxLinksPerPort:          2,
						MaxLinksPerSlot:          2,
						MaxPerSystemLinksPerPort: 2,
						MaxPerSystemLinksPerSlot: 2,
						Mode:                     apstra.AntiAffinityModeEnabledLoose,
					},
					SpineLeafLinks:       testutils.ToPtr(apstra.AddressingSchemeIp4),
					SpineSuperspineLinks: testutils.ToPtr(apstra.AddressingSchemeIp4),
				},
			},
		},
		"different_values": {
			req: apstra.CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      testutils.RandString(5, "hex"),
				TemplateId: "L2_Virtual_EVPN",
				FabricSettings: &apstra.FabricSettings{
					JunosEvpnDuplicateMacRecoveryTime:     testutils.ToPtr(uint16(14)),
					MaxExternalRoutes:                     testutils.ToPtr(uint32(233832)),
					EsiMacMsb:                             testutils.ToPtr(uint8(50)),
					JunosGracefulRestart:                  &enum.FeatureSwitchEnabled,
					OptimiseSzFootprint:                   &enum.FeatureSwitchEnabled,
					JunosEvpnRoutingInstanceVlanAware:     &enum.FeatureSwitchEnabled,
					EvpnGenerateType5HostRoutes:           &enum.FeatureSwitchEnabled,
					MaxFabricRoutes:                       testutils.ToPtr(uint32(82231)),
					MaxMlagRoutes:                         testutils.ToPtr(uint32(74112)),
					JunosExOverlayEcmp:                    &enum.FeatureSwitchEnabled,
					DefaultSviL3Mtu:                       testutils.ToPtr(uint16(9070)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &enum.FeatureSwitchEnabled,
					FabricL3Mtu:                           testutils.ToPtr(uint16(9172)),
					Ipv6Enabled:                           testutils.ToPtr(false), // do not enable because it's a one-way trip
					ExternalRouterMtu:                     testutils.ToPtr(uint16(9080)),
					MaxEvpnRoutes:                         testutils.ToPtr(uint32(91342)),
					AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
						Algorithm:                apstra.AlgorithmHeuristic,
						MaxLinksPerPort:          4,
						MaxLinksPerSlot:          4,
						MaxPerSystemLinksPerPort: 4,
						MaxPerSystemLinksPerSlot: 4,
						Mode:                     apstra.AntiAffinityModeEnabledLoose,
					},
					SpineLeafLinks:       testutils.ToPtr(apstra.AddressingSchemeIp46),
					SpineSuperspineLinks: testutils.ToPtr(apstra.AddressingSchemeIp46),
				},
			},
		},
	}

	fetchFabricAddressingScheme := func(t testing.TB, client *apstra.TwoStageL3ClosClient) (apstra.AddressingScheme, apstra.AddressingScheme) {
		t.Helper()

		query := new(apstra.PathQuery).
			SetClient(client.Client()).
			SetBlueprintId(client.Id())
		if compatibility.BpHasFabricAddressingPolicyNode.Check(version.Must(version.NewVersion(client.Client().ApiVersion()))) {
			query.Node([]apstra.QEEAttribute{
				apstra.NodeTypeFabricAddressingPolicy.QEEAttribute(),
				{Key: "name", Value: apstra.QEStringVal("node")},
			})
		} else {
			query.Node([]apstra.QEEAttribute{
				apstra.NodeTypeFabricPolicy.QEEAttribute(),
				{Key: "name", Value: apstra.QEStringVal("node")},
			})
		}

		var queryResponse struct {
			Items []struct {
				Node struct {
					SpineLeafLinks       string `json:"spine_leaf_links"`
					SpineSuperspineLinks string `json:"spine_superspine_links"`
				} `json:"node"`
			} `json:"items"`
		}

		err := query.Do(ctx, &queryResponse)
		require.NoError(t, err)
		require.Equal(t, 1, len(queryResponse.Items))

		var spineLeaf, spineSuperspine apstra.AddressingScheme
		require.NoError(t, spineLeaf.FromString(queryResponse.Items[0].Node.SpineLeafLinks))
		require.NoError(t, spineSuperspine.FromString(queryResponse.Items[0].Node.SpineLeafLinks))

		return spineLeaf, spineSuperspine
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			ctx := testutils.ContextWithTestID(ctx, t)

			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					if tCase.constraint != nil && !tCase.constraint.Check(client.APIVersion()) {
						t.Skipf("skipping test case %q with Apstra %s due to version constraint %q", tName, client.Client.ApiVersion(), tCase.constraint)
					}

					id, err := client.Client.CreateBlueprintFromTemplate(ctx, &tCase.req)
					require.NoError(t, err)

					if !compatibility.FabricSettingsApiOk.Check(client.APIVersion()) && tCase.req.FabricSettings != nil {
						// 4.2.0 cannot set fabric settings when creating blueprint, so we have to do it afterward
						bp, err := client.Client.NewTwoStageL3ClosClient(ctx, id)
						require.NoError(t, err)

						// spine/leaf and spine/superspine addressing cannot be set, so we clear these
						fs := tCase.req.FabricSettings
						fs.SpineLeafLinks = nil
						fs.SpineSuperspineLinks = nil
						require.NoError(t, bp.SetFabricSettings(ctx, fs))
					}

					bpClient, err := client.Client.NewTwoStageL3ClosClient(ctx, id)
					require.NoError(t, err)

					if tCase.req.FabricSettings != nil {
						fabricSettings, err := bpClient.GetFabricSettings(ctx)
						require.NoError(t, err)

						compare.FabricSettings(t, *tCase.req.FabricSettings, *fabricSettings)

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

					require.NoError(t, client.Client.DeleteBlueprint(ctx, id))
				})
			}
		})
	}
}

func TestCreateDeleteIpFabricBlueprint(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		req           apstra.CreateBlueprintFromTemplateRequest
		compatibility *compatibility.Constraint
	}

	testCases := map[string]testCase{
		"simple": {
			req: apstra.CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      testutils.RandString(5, "hex"),
				TemplateId: "L2_Virtual",
			},
		},
		"4.1.1_and_later": {
			req: apstra.CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      testutils.RandString(5, "hex"),
				TemplateId: "L2_Virtual",
				FabricSettings: &apstra.FabricSettings{
					SpineLeafLinks:       testutils.ToPtr(apstra.AddressingSchemeIp46),
					SpineSuperspineLinks: testutils.ToPtr(apstra.AddressingSchemeIp46),
				},
			},
		},
		"4.2.0_specific_test": {
			req: apstra.CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      testutils.RandString(5, "hex"),
				TemplateId: "L2_Virtual",
				FabricSettings: &apstra.FabricSettings{
					SpineLeafLinks:       testutils.ToPtr(apstra.AddressingSchemeIp46),
					SpineSuperspineLinks: testutils.ToPtr(apstra.AddressingSchemeIp46),
					FabricL3Mtu:          testutils.ToPtr(uint16(9178)),
				},
			},
		},
		"lots_of_values": {
			req: apstra.CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      testutils.RandString(5, "hex"),
				TemplateId: "L2_Virtual",
				FabricSettings: &apstra.FabricSettings{
					JunosEvpnDuplicateMacRecoveryTime:     testutils.ToPtr(uint16(16)),
					MaxExternalRoutes:                     testutils.ToPtr(uint32(239832)),
					EsiMacMsb:                             testutils.ToPtr(uint8(32)),
					JunosGracefulRestart:                  &enum.FeatureSwitchDisabled,
					OptimiseSzFootprint:                   &enum.FeatureSwitchDisabled,
					JunosEvpnRoutingInstanceVlanAware:     &enum.FeatureSwitchDisabled,
					EvpnGenerateType5HostRoutes:           &enum.FeatureSwitchDisabled,
					MaxFabricRoutes:                       testutils.ToPtr(uint32(84231)),
					MaxMlagRoutes:                         testutils.ToPtr(uint32(76112)),
					JunosExOverlayEcmp:                    &enum.FeatureSwitchDisabled,
					DefaultSviL3Mtu:                       testutils.ToPtr(uint16(9100)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &enum.FeatureSwitchDisabled,
					FabricL3Mtu:                           testutils.ToPtr(uint16(9178)),
					Ipv6Enabled:                           testutils.ToPtr(false), // do not enable because it's a one-way trip
					ExternalRouterMtu:                     testutils.ToPtr(uint16(9100)),
					MaxEvpnRoutes:                         testutils.ToPtr(uint32(92342)),
					AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
						Algorithm:                apstra.AlgorithmHeuristic,
						MaxLinksPerPort:          2,
						MaxLinksPerSlot:          2,
						MaxPerSystemLinksPerPort: 2,
						MaxPerSystemLinksPerSlot: 2,
						Mode:                     apstra.AntiAffinityModeEnabledLoose,
					},
					SpineLeafLinks:       testutils.ToPtr(apstra.AddressingSchemeIp4),
					SpineSuperspineLinks: testutils.ToPtr(apstra.AddressingSchemeIp4),
				},
			},
		},
		"different_values": {
			req: apstra.CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      testutils.RandString(5, "hex"),
				TemplateId: "L2_Virtual",
				FabricSettings: &apstra.FabricSettings{
					JunosEvpnDuplicateMacRecoveryTime:     testutils.ToPtr(uint16(14)),
					MaxExternalRoutes:                     testutils.ToPtr(uint32(233832)),
					EsiMacMsb:                             testutils.ToPtr(uint8(50)),
					JunosGracefulRestart:                  &enum.FeatureSwitchEnabled,
					OptimiseSzFootprint:                   &enum.FeatureSwitchEnabled,
					JunosEvpnRoutingInstanceVlanAware:     &enum.FeatureSwitchEnabled,
					EvpnGenerateType5HostRoutes:           &enum.FeatureSwitchEnabled,
					MaxFabricRoutes:                       testutils.ToPtr(uint32(82231)),
					MaxMlagRoutes:                         testutils.ToPtr(uint32(74112)),
					JunosExOverlayEcmp:                    &enum.FeatureSwitchEnabled,
					DefaultSviL3Mtu:                       testutils.ToPtr(uint16(9070)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &enum.FeatureSwitchEnabled,
					FabricL3Mtu:                           testutils.ToPtr(uint16(9172)),
					Ipv6Enabled:                           testutils.ToPtr(false), // do not enable because it's a one-way trip
					ExternalRouterMtu:                     testutils.ToPtr(uint16(9080)),
					MaxEvpnRoutes:                         testutils.ToPtr(uint32(91342)),
					AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
						Algorithm:                apstra.AlgorithmHeuristic,
						MaxLinksPerPort:          4,
						MaxLinksPerSlot:          4,
						MaxPerSystemLinksPerPort: 4,
						MaxPerSystemLinksPerSlot: 4,
						Mode:                     apstra.AntiAffinityModeEnabledLoose,
					},
					SpineLeafLinks:       testutils.ToPtr(apstra.AddressingSchemeIp46),
					SpineSuperspineLinks: testutils.ToPtr(apstra.AddressingSchemeIp46),
				},
			},
		},
	}

	fetchFabricAddressingScheme := func(t testing.TB, client *apstra.TwoStageL3ClosClient) (apstra.AddressingScheme, apstra.AddressingScheme) {
		t.Helper()

		query := new(apstra.PathQuery).
			SetClient(client.Client()).
			SetBlueprintId(client.Id())
		if compatibility.BpHasFabricAddressingPolicyNode.Check(version.Must(version.NewVersion(client.Client().ApiVersion()))) {
			query.Node([]apstra.QEEAttribute{
				apstra.NodeTypeFabricAddressingPolicy.QEEAttribute(),
				{Key: "name", Value: apstra.QEStringVal("node")},
			})
		} else {
			query.Node([]apstra.QEEAttribute{
				apstra.NodeTypeFabricPolicy.QEEAttribute(),
				{Key: "name", Value: apstra.QEStringVal("node")},
			})
		}

		var queryResponse struct {
			Items []struct {
				Node struct {
					SpineLeafLinks       string `json:"spine_leaf_links"`
					SpineSuperspineLinks string `json:"spine_superspine_links"`
				} `json:"node"`
			} `json:"items"`
		}

		err := query.Do(ctx, &queryResponse)
		require.NoError(t, err)
		require.Equal(t, 1, len(queryResponse.Items))

		var spineLeaf, spineSuperspine apstra.AddressingScheme
		require.NoError(t, spineLeaf.FromString(queryResponse.Items[0].Node.SpineLeafLinks))
		require.NoError(t, spineSuperspine.FromString(queryResponse.Items[0].Node.SpineSuperspineLinks))

		return spineLeaf, spineSuperspine
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			ctx := testutils.ContextWithTestID(ctx, t)

			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					if tCase.compatibility != nil && !tCase.compatibility.Check(client.APIVersion()) {
						t.Skipf("skipping test case %q with Apstra %s due to version constraint %q", tName, client.Client.ApiVersion(), tCase.compatibility)
					}

					id, err := client.Client.CreateBlueprintFromTemplate(ctx, &tCase.req)
					require.NoError(t, err)

					if !compatibility.FabricSettingsApiOk.Check(client.APIVersion()) && tCase.req.FabricSettings != nil {
						// 4.2.0 cannot set fabric settings when creating blueprint, so we have to do it afterward
						bp, err := client.Client.NewTwoStageL3ClosClient(ctx, id)
						require.NoError(t, err)

						// spine/leaf and spine/superspine addressing cannot be set, so we clear these
						fs := tCase.req.FabricSettings
						fs.SpineLeafLinks = nil
						fs.SpineSuperspineLinks = nil
						require.NoError(t, bp.SetFabricSettings(ctx, fs))
					}

					bpClient, err := client.Client.NewTwoStageL3ClosClient(ctx, id)
					require.NoError(t, err)

					if tCase.req.FabricSettings != nil {
						fabricSettings, err := bpClient.GetFabricSettings(ctx)
						require.NoError(t, err)

						compare.FabricSettings(t, *tCase.req.FabricSettings, *fabricSettings)

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

					require.NoError(t, client.Client.DeleteBlueprint(ctx, id))
				})
			}
		})
	}
}

func TestCreateDeleteBlueprintWithRoutingLimits(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	blueprintRequest := func() apstra.CreateBlueprintFromTemplateRequest {
		return apstra.CreateBlueprintFromTemplateRequest{
			RefDesign:  enum.RefDesignDatacenter,
			Label:      testutils.RandString(5, "hex"),
			TemplateId: "L2_Virtual",
		}
	}

	type testCase struct {
		name           string
		fabricSettings apstra.FabricSettings
	}

	testCases := []testCase{
		{
			name:           "create_with_defaults",
			fabricSettings: apstra.FabricSettings{},
		},
		{
			name: "create_with_zeros",
			fabricSettings: apstra.FabricSettings{
				MaxEvpnRoutes:     testutils.ToPtr(uint32(0)),
				MaxExternalRoutes: testutils.ToPtr(uint32(0)),
				MaxFabricRoutes:   testutils.ToPtr(uint32(0)),
				MaxMlagRoutes:     testutils.ToPtr(uint32(0)),
			},
		},
		{
			name: "create_with_values",
			fabricSettings: apstra.FabricSettings{
				MaxEvpnRoutes:     testutils.ToPtr(uint32(20001)),
				MaxExternalRoutes: testutils.ToPtr(uint32(20002)),
				MaxFabricRoutes:   testutils.ToPtr(uint32(20003)),
				MaxMlagRoutes:     testutils.ToPtr(uint32(20004)),
			},
		},
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			if !compatibility.GeApstra421.Check(client.APIVersion()) {
				t.Skipf("skipping Apstra %s client due to version constraint", client.Client.ApiVersion())
			}

			for _, tCase := range testCases {
				client := client
				tCase := tCase

				t.Run(tCase.name, func(t *testing.T) {
					ctx := testutils.ContextWithTestID(ctx, t)

					bpr := blueprintRequest()
					bpr.FabricSettings = &tCase.fabricSettings
					id, err := client.Client.CreateBlueprintFromTemplate(ctx, &bpr)
					require.NoError(t, err)

					bpClient, err := client.Client.NewTwoStageL3ClosClient(ctx, id)
					require.NoError(t, err)

					fabricSettings, err := bpClient.GetFabricSettings(ctx)
					require.NoError(t, err)

					compare.FabricSettings(t, tCase.fabricSettings, *fabricSettings)

					require.NoError(t, client.Client.DeleteBlueprint(ctx, id))
				})
			}
		})
	}
}

// This test deletes all blueprints, so is likely disruptive to other tests
// It can be run from the command line to quickly clean-up Apstra servers
// with lots of left-behind blueprints:
/*
DELETE_ALL_BLUEPRINTS_WITH_A_TEST=1 go test -v -run=TestDeleteAllBlueprints -tags=integration $(git rev-parse --show-toplevel)/apstra
*/
func TestDeleteAllBlueprints(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	if _, ok := os.LookupEnv("DELETE_ALL_BLUEPRINTS_WITH_A_TEST"); !ok {
		t.Skip("refusing to run without DELETE_ALL_BLUEPRINTS_WITH_A_TEST in environment")
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			ids, err := client.Client.ListAllBlueprintIds(ctx)
			require.NoError(t, err)

			for _, id := range ids {
				require.NoError(t, client.Client.DeleteBlueprint(ctx, id))
			}
		})
	}
}
