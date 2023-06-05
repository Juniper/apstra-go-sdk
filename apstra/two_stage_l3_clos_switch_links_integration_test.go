//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
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
		links := make([]SwitchLink, len(leafResult.Items))
		for i, item := range leafResult.Items {
			links[i] = SwitchLink{
				Tags:           []string{"link blah", "link also blah"},
				SystemEndpoint: SwitchLinkEndpoint{},
				SwitchEndpoint: SwitchLinkEndpoint{
					TransformationId: 1,
					SystemId:         item.Leaf.Id,
					IfName:           "xe-0/0/3",
				},
			}
		}

		log.Printf("testing CreateLinksWithNewServer() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		linkIds, err := bpClient.CreateLinksWithNewServer(ctx, &CreateLinksWithNewServerRequest{
			Server: System{
				Hostname:         randString(5, "hex"),
				Label:            randString(5, "hex"),
				LogicalDeviceId:  "AOS-2x10-1",
				PortChannelIdMin: 0,
				PortChannelIdMax: 0,
				Tags:             []string{"blah", "also blah"},
			},
			Links: links,
		})
		if err != nil {
			t.Fatal(err)
		}

		linkQuery := new(MatchQuery).
			SetBlueprintId(bpClient.blueprintId).
			SetBlueprintType(BlueprintTypeStaging).
			SetClient(bpClient.Client())

		for _, linkId := range linkIds {
			linkQuery.Match(
				new(PathQuery).
					Node([]QEEAttribute{NodeTypeLink.QEEAttribute(),
						{"id", QEStringVal(linkId.String())},
					}).
					In([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
					Node([]QEEAttribute{NodeTypeInterface.QEEAttribute()}).
					In([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
					Node([]QEEAttribute{NodeTypeSystem.QEEAttribute(),
						{"role", QEStringVal("generic")},
						{"name", QEStringVal("n_generic")},
					}),
			)
		}

		linkResult := struct {
			Items []struct {
				Generic struct {
					Id ObjectId `json:"id"`
				} `json:"n_generic"`
			} `json:"items"`
		}{}

		err = linkQuery.Do(ctx, &linkResult)
		if err != nil {
			t.Fatal(err)
		}

		if len(linkResult.Items) != 1 {
			t.Fatalf("expected 1 item, got %d items", len(linkResult.Items))
		}

		log.Printf("testing DeleteGenericSystem() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeleteGenericSystem(ctx, linkResult.Items[0].Generic.Id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
