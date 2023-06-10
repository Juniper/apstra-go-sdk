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
)

func TestCreateDeleteServer(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintD(ctx, t, client.client)
		defer func() {
			err := bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		//bpClient, err := client.client.NewTwoStageL3ClosClient(ctx, "c09e3975-6799-41a3-ab1a-d96f93cd5d3e")
		//if err != nil {
		//	t.Fatal(err)
		//}

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

		log.Printf("testing CreateLinksWithNewServer() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		linkIds, err := bpClient.CreateLinksWithNewServer(ctx, &CreateLinksWithNewServerRequest{
			Server: CreateLinksWithNewServerRequestServer{
				Hostname:         randString(5, "hex"),
				Label:            randString(5, "hex"),
				LogicalDeviceId:  "AOS-2x10-1",
				PortChannelIdMin: 0,
				PortChannelIdMax: 0,
				Tags:             desiredTags,
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

		log.Printf("testing DeleteGenericSystem() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeleteGenericSystem(ctx, systemId)
		if err != nil {
			t.Fatal(err)
		}
	}
}
