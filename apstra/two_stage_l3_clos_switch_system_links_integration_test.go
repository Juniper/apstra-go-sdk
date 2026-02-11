// Copyright (c) Juniper Networks, Inc., 2023-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDeleteGenericSystem(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintD(ctx, t, client.client)

		leafQuery := new(PathQuery).
			SetBlueprintId(bpClient.Id()).
			SetBlueprintType(BlueprintTypeStaging).
			SetClient(bpClient.client).
			Node([]QEEAttribute{
				NodeTypeRack.QEEAttribute(),
				{"label", QEStringVal("l2_esi_2x_links_001")},
			}).
			In([]QEEAttribute{RelationshipTypePartOfRack.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeSystem.QEEAttribute(),
				{"system_type", QEStringVal("switch")},
				{"role", QEStringVal("leaf")},
				{"name", QEStringVal("n_leaf")},
			})

		var leafResult struct {
			Items []struct {
				Leaf struct {
					Id ObjectId `json:"id"`
				} `json:"n_leaf"`
			} `json:"items"`
		}

		err = leafQuery.Do(ctx, &leafResult)
		if err != nil {
			t.Fatal(err)
		}
		links := make([]CreateLinkRequest, len(leafResult.Items))
		for i, item := range leafResult.Items {
			links[i] = CreateLinkRequest{
				LagMode:        RackLinkLagModeActive,
				GroupLabel:     "foo",
				Tags:           []string{"link blah", "link also blah"},
				SystemEndpoint: SwitchLinkEndpoint{},
				SwitchEndpoint: SwitchLinkEndpoint{
					TransformationId: 1,
					SystemId:         item.Leaf.Id,
					IfName:           "xe-0/0/3",
				},
			}
		}

		var desiredTags []string
		for i := 0; i < rand.Intn(3)+2; i++ {
			desiredTags = append(desiredTags, randString(5, "hex"))
		}

		log.Printf("testing CreateLinksWithNewSystem() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		linkIds, err := bpClient.CreateLinksWithNewSystem(ctx, &CreateLinksWithNewSystemRequest{
			System: CreateLinksWithNewSystemRequestSystem{
				Hostname:         randString(5, "hex"),
				Label:            randString(5, "hex"),
				LogicalDeviceId:  "AOS-2x10-1",
				PortChannelIdMin: 0,
				PortChannelIdMax: 0,
				Tags:             desiredTags,
				Type:             SystemTypeServer,
			},
			Links: links,
		})
		if err != nil {
			t.Fatal(err)
		}

		systemId, err := bpClient.SystemNodeFromLinkIds(ctx, linkIds, SystemNodeRoleGeneric)
		if err != nil {
			t.Fatal(err)
		}

		observedTags, err := bpClient.GetNodeTags(ctx, systemId)
		if err != nil {
			t.Fatal(err)
		}

		cmLinks, err := bpClient.GetCablingMapLinksBySystem(ctx, systemId)
		if err != nil {
			t.Fatal(err)
		}

		var aggregateCount int
		for i := range cmLinks {
			if cmLinks[i].Type == LinkTypeAggregateLink {
				aggregateCount++
			}
		}

		if aggregateCount != 1 {
			t.Fatalf("expected 1 aggregate link, got %d", aggregateCount)
		}

		sort.Strings(desiredTags)
		sort.Strings(observedTags)
		compareSlices(t, desiredTags, observedTags, fmt.Sprintf("generic system tags"))

		newLinks := make([]CreateLinkRequest, len(leafResult.Items))
		for i, item := range leafResult.Items {
			newLinks[i] = CreateLinkRequest{
				LagMode:    RackLinkLagModePassive,
				GroupLabel: "bar",
				Tags:       []string{"a", "b"},
				SystemEndpoint: SwitchLinkEndpoint{
					SystemId: systemId,
				},
				SwitchEndpoint: SwitchLinkEndpoint{
					TransformationId: 1,
					SystemId:         item.Leaf.Id,
					IfName:           "xe-0/0/2",
				},
			}
		}

		log.Printf("testing AddLinksToSystem() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		newLinkIds, err := bpClient.AddLinksToSystem(ctx, newLinks)
		if err != nil {
			t.Fatal(err)
		}
		if len(leafResult.Items) != len(newLinkIds) {
			t.Fatalf("expected %d additional link IDs, got %d", len(leafResult.Items), len(newLinks))
		}

		log.Printf("testing DeleteGenericSystem() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeleteGenericSystem(ctx, systemId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCreateDeleteExternalGenericSystem(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintD(ctx, t, client.client)

		leafQuery := new(PathQuery).
			SetBlueprintId(bpClient.Id()).
			SetBlueprintType(BlueprintTypeStaging).
			SetClient(bpClient.client).
			Node([]QEEAttribute{
				NodeTypeRack.QEEAttribute(),
				{"label", QEStringVal("l2_esi_2x_links_001")},
			}).
			In([]QEEAttribute{RelationshipTypePartOfRack.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeSystem.QEEAttribute(),
				{"system_type", QEStringVal("switch")},
				{"role", QEStringVal("leaf")},
				{"name", QEStringVal("n_leaf")},
			})

		var leafResult struct {
			Items []struct {
				Leaf struct {
					Id ObjectId `json:"id"`
				} `json:"n_leaf"`
			} `json:"items"`
		}

		err = leafQuery.Do(ctx, &leafResult)
		if err != nil {
			t.Fatal(err)
		}
		links := make([]CreateLinkRequest, len(leafResult.Items))
		for i, item := range leafResult.Items {
			links[i] = CreateLinkRequest{
				LagMode:        RackLinkLagModeActive,
				GroupLabel:     "foo",
				Tags:           []string{"link blah", "link also blah"},
				SystemEndpoint: SwitchLinkEndpoint{},
				SwitchEndpoint: SwitchLinkEndpoint{
					TransformationId: 1,
					SystemId:         item.Leaf.Id,
					IfName:           "xe-0/0/3",
				},
			}
		}

		var desiredTags []string
		for i := 0; i < rand.Intn(3)+2; i++ {
			desiredTags = append(desiredTags, randString(5, "hex"))
		}

		log.Printf("testing CreateLinksWithNewSystem() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		linkIds, err := bpClient.CreateLinksWithNewSystem(ctx, &CreateLinksWithNewSystemRequest{
			System: CreateLinksWithNewSystemRequestSystem{
				Hostname:         randString(5, "hex"),
				Label:            randString(5, "hex"),
				LogicalDeviceId:  "AOS-2x10-1",
				PortChannelIdMin: 0,
				PortChannelIdMax: 0,
				Tags:             desiredTags,
				Type:             SystemTypeExternal,
			},
			Links: links,
		})
		if err != nil {
			t.Fatal(err)
		}

		systemId, err := bpClient.SystemNodeFromLinkIds(ctx, linkIds, SystemNodeRoleGeneric)
		if err != nil {
			t.Fatal(err)
		}

		observedTags, err := bpClient.GetNodeTags(ctx, systemId)
		if err != nil {
			t.Fatal(err)
		}

		cmLinks, err := bpClient.GetCablingMapLinksBySystem(ctx, systemId)
		if err != nil {
			t.Fatal(err)
		}

		var aggregateCount int
		for i := range cmLinks {
			if cmLinks[i].Type == LinkTypeAggregateLink {
				aggregateCount++
			}
		}

		if aggregateCount != 1 {
			t.Fatalf("expected 1 aggregate link, got %d", aggregateCount)
		}

		sort.Strings(desiredTags)
		sort.Strings(observedTags)
		compareSlices(t, desiredTags, observedTags, fmt.Sprintf("generic system tags"))

		newLinks := make([]CreateLinkRequest, len(leafResult.Items))
		for i, item := range leafResult.Items {
			newLinks[i] = CreateLinkRequest{
				LagMode:    RackLinkLagModePassive,
				GroupLabel: "bar",
				Tags:       []string{"a", "b"},
				SystemEndpoint: SwitchLinkEndpoint{
					SystemId: systemId,
				},
				SwitchEndpoint: SwitchLinkEndpoint{
					TransformationId: 1,
					SystemId:         item.Leaf.Id,
					IfName:           "xe-0/0/2",
				},
			}
		}

		log.Printf("testing AddLinksToSystem() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		newLinkIds, err := bpClient.AddLinksToSystem(ctx, newLinks)
		if err != nil {
			t.Fatal(err)
		}
		if len(leafResult.Items) != len(newLinkIds) {
			t.Fatalf("expected %d additional link IDs, got %d", len(leafResult.Items), len(newLinks))
		}

		log.Printf("testing DeleteGenericSystem() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeleteGenericSystem(ctx, systemId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestDeleteSwitchSystemLinks_WithCtAssigned(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			// create a blueprint
			log.Printf("creating test blueprint against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bp := testBlueprintC(ctx, t, client.client)

			// collect leaf switch IDs
			log.Printf("determining leaf switch IDs in blueprint %q %s %s (%s)", bp.Id(), client.clientType, clientName, client.client.ApiVersion())
			nodeMap, err := bp.GetAllSystemNodeInfos(ctx)
			require.NoError(t, err)
			var leafIds []ObjectId
			for _, node := range nodeMap {
				if node.Role == SystemRoleLeaf {
					leafIds = append(leafIds, node.Id)
				}
			}

			// assign leaf switch interface map
			log.Printf("assigning leaf switch interface maps in blueprint %q %s %s (%s)", bp.Id(), client.clientType, clientName, client.client.ApiVersion())
			assignments := make(SystemIdToInterfaceMapAssignment)
			for _, leafId := range leafIds {
				assignments[leafId.String()] = "Juniper_vQFX__AOS-7x10-Leaf"
			}
			require.NoError(t, bp.SetInterfaceMapAssignments(ctx, assignments))

			// create security zone
			szId := testSecurityZone(t, ctx, bp)

			// create virtual networks
			log.Printf("creating virtual networks in blueprint %q %s %s (%s)", bp.Id(), client.clientType, clientName, client.client.ApiVersion())
			vnIds := make([]ObjectId, 3)
			for i := range vnIds {
				vnIds[i] = testVirtualNetwork(t, ctx, bp, ObjectId(szId))
			}

			// create connectivity templates
			log.Printf("creating connectivity templates in blueprint %q %s %s (%s)", bp.Id(), client.clientType, clientName, client.client.ApiVersion())
			ctIds := make([]ObjectId, len(vnIds))
			for i, vnId := range vnIds {
				ct := ConnectivityTemplate{
					Id:    nil,
					Label: randString(6, "hex"),
					Subpolicies: []*ConnectivityTemplatePrimitive{
						{
							Label: "",
							Attributes: &ConnectivityTemplatePrimitiveAttributesAttachSingleVlan{
								Tagged:   true,
								VnNodeId: &vnId,
							},
						},
					},
				}
				err = bp.CreateConnectivityTemplate(ctx, &ct)
				require.NoError(t, err)

				ctIds[i] = *ct.Id
			}

			// create generic system
			log.Printf("creating generic system in blueprint %q %s %s (%s)", bp.Id(), client.clientType, clientName, client.client.ApiVersion())
			linkIds, err := bp.CreateLinksWithNewSystem(ctx, &CreateLinksWithNewSystemRequest{
				Links: []CreateLinkRequest{
					{
						SwitchEndpoint: SwitchLinkEndpoint{
							TransformationId: 1,
							SystemId:         leafIds[0],
							IfName:           "xe-0/0/7",
						},
						GroupLabel: "",
						LagMode:    0,
					},
					{
						SwitchEndpoint: SwitchLinkEndpoint{
							TransformationId: 1,
							SystemId:         leafIds[0],
							IfName:           "xe-0/0/8",
						},
						GroupLabel: "",
						LagMode:    0,
					},
					{
						SwitchEndpoint: SwitchLinkEndpoint{
							TransformationId: 1,
							SystemId:         leafIds[0],
							IfName:           "xe-0/0/9",
						},
						GroupLabel: "a",
						LagMode:    RackLinkLagModeActive,
					},
					{
						SwitchEndpoint: SwitchLinkEndpoint{
							TransformationId: 1,
							SystemId:         leafIds[0],
							IfName:           "xe-0/0/10",
						},
						GroupLabel: "a",
						LagMode:    RackLinkLagModeActive,
					},
					{
						SwitchEndpoint: SwitchLinkEndpoint{
							TransformationId: 1,
							SystemId:         leafIds[0],
							IfName:           "xe-0/0/11",
						},
						GroupLabel: "b",
						LagMode:    RackLinkLagModeActive,
					},
				},
				System: CreateLinksWithNewSystemRequestSystem{
					LogicalDeviceId: "AOS-1x10-1",
					Type:            SystemTypeServer,
				},
			})
			require.NoError(t, err)
			require.EqualValuesf(t, 5, len(linkIds), "expected 1 generic system id, got %d", len(linkIds))

			// collect links we'll fiddle with
			log.Printf("determining GS LAG IDs in blueprint %q %s %s (%s)", bp.Id(), client.clientType, clientName, client.client.ApiVersion())
			testLinks := []ObjectId{linkIds[1]}
			lagIdA, err := bp.lagIdFromMemberIds(ctx, linkIds[2:4])
			require.NoError(t, err)
			lagIdB, err := bp.lagIdFromMemberIds(ctx, linkIds[4:5])
			require.NoError(t, err)
			testLinks = append(testLinks, lagIdA, lagIdB)

			// determine the application points (switch ports) associated with the test link IDs
			log.Printf("determining GS switch port IDs in blueprint %q %s %s (%s)", bp.Id(), client.clientType, clientName, client.client.ApiVersion())
			applicationPoints := make([]ObjectId, len(testLinks))
			for i, linkId := range testLinks {
				q := new(PathQuery).SetBlueprintId(bp.Id()).SetClient(bp.Client()).
					Node([]QEEAttribute{{Key: "id", Value: QEStringVal(linkId.String())}}).
					In([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
					Node([]QEEAttribute{
						NodeTypeInterface.QEEAttribute(),
						{Key: "name", Value: QEStringVal("n_interface")},
					}).
					In([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
					Node([]QEEAttribute{{Key: "id", Value: QEStringVal(leafIds[0].String())}})

				var result struct {
					Items []struct {
						Interface struct {
							Id ObjectId `json:"id"`
						} `json:"n_interface"`
					} `json:"items"`
				}

				require.NoError(t, q.Do(ctx, &result))
				require.EqualValuesf(t, 1, len(result.Items), "expected 1 result got %d", len(result.Items))

				applicationPoints[i] = result.Items[0].Interface.Id
			}

			// assign CTs to application points
			log.Printf("assigning connectivity templatee in blueprint %q %s %s (%s)", bp.Id(), client.clientType, clientName, client.client.ApiVersion())
			for i, apId := range applicationPoints {
				require.NoError(t, bp.SetApplicationPointConnectivityTemplates(ctx, apId, []ObjectId{ctIds[i%len(ctIds)]}))
			}

			err = bp.DeleteLinksFromSystem(ctx, linkIds[1:])
			require.Error(t, err)
			require.IsType(t, ClientErr{}, err)
			require.EqualValues(t, err.(ClientErr).Type(), ErrCtAssignedToLink)
			require.EqualValues(t, 3, len(err.(ClientErr).Detail().(ErrCtAssignedToLinkDetail).LinkIds))
		})
	}
}
