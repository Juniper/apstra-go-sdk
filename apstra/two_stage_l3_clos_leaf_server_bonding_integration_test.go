package apstra

import (
	"context"
	"fmt"
	"log"
	"sort"
	"testing"
)

func TestSetGenericServerBonding(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintE(ctx, t, client.client)
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
				{"label", QEStringVal("l2_esi_acs_single_001")},
			}).
			In([]QEEAttribute{RelationshipTypePartOfRack.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeSystem.QEEAttribute(),
				{"system_type", QEStringVal("switch")},
				{"role", QEStringVal("leaf")},
				{"name", QEStringVal("n_system")},
			})

		accessQuery := new(PathQuery).
			SetBlueprintId(bpClient.Id()).
			SetBlueprintType(BlueprintTypeStaging).
			SetClient(bpClient.client).
			Node([]QEEAttribute{
				NodeTypeRack.QEEAttribute(),
				{"label", QEStringVal("l2_esi_acs_single_001")},
			}).
			In([]QEEAttribute{RelationshipTypePartOfRack.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeSystem.QEEAttribute(),
				{"system_type", QEStringVal("switch")},
				{"role", QEStringVal("access")},
				{"name", QEStringVal("n_system")},
			})

		var switchQueryResult struct {
			Items []struct {
				System struct {
					Id    ObjectId `json:"id"`
					Label string   `json:"label"`
				} `json:"n_system"`
			} `json:"items"`
		}

		// get a slice of leaf switch IDs
		err = leafQuery.Do(ctx, &switchQueryResult)
		if err != nil {
			t.Fatal(err)
		}
		sort.Slice(switchQueryResult.Items, func(i, j int) bool {
			return switchQueryResult.Items[i].System.Label < switchQueryResult.Items[j].System.Label
		})
		leafIds := make([]ObjectId, len(switchQueryResult.Items))
		for i := range switchQueryResult.Items {
			leafIds[i] = switchQueryResult.Items[i].System.Id
		}

		// get a slice of access switch IDs
		err = accessQuery.Do(ctx, &switchQueryResult)
		if err != nil {
			t.Fatal(err)
		}
		sort.Slice(switchQueryResult.Items, func(i, j int) bool {
			return switchQueryResult.Items[i].System.Label < switchQueryResult.Items[j].System.Label
		})
		accessIds := make([]ObjectId, len(switchQueryResult.Items))
		for i := range switchQueryResult.Items {
			accessIds[i] = switchQueryResult.Items[i].System.Id
		}

		// used for querying links to learn new server ID
		linkQueryResult := struct {
			Items []struct {
				Generic struct {
					Id ObjectId `json:"id"`
				} `json:"n_generic"`
			} `json:"items"`
		}{}

		server1LinkCount := 4
		server1FirstPort := 5
		server1Links := make([]SwitchLink, server1LinkCount)
		for i := 0; i < server1LinkCount; i++ {
			server1Links[i] = SwitchLink{
				SystemEndpoint: SwitchLinkEndpoint{},
				SwitchEndpoint: SwitchLinkEndpoint{
					TransformationId: 1,
					SystemId:         leafIds[0],
					IfName:           fmt.Sprintf("xe-0/0/%d", server1FirstPort+i),
				},
			}
		}

		server2Links := make([]SwitchLink, len(leafIds))
		for i, id := range leafIds {
			server2Links[i] = SwitchLink{
				SystemEndpoint: SwitchLinkEndpoint{},
				SwitchEndpoint: SwitchLinkEndpoint{
					TransformationId: 1,
					SystemId:         id,
					IfName:           "xe-0/0/2",
				},
			}
		}

		server3LinkCount := 4
		server3FirstPort := 5
		server3Links := make([]SwitchLink, server3LinkCount)
		for i := 0; i < server3LinkCount; i++ {
			server3Links[i] = SwitchLink{
				SystemEndpoint: SwitchLinkEndpoint{},
				SwitchEndpoint: SwitchLinkEndpoint{
					TransformationId: 1,
					SystemId:         accessIds[0],
					IfName:           fmt.Sprintf("xe-0/0/%d", server3FirstPort+i),
				},
			}
		}

		server4Links := make([]SwitchLink, len(accessIds))
		for i, id := range accessIds {
			server4Links[i] = SwitchLink{
				SystemEndpoint: SwitchLinkEndpoint{},
				SwitchEndpoint: SwitchLinkEndpoint{
					TransformationId: 1,
					SystemId:         id,
					IfName:           "xe-0/0/9",
				},
			}
		}

		serverLinks := [][]SwitchLink{server1Links, server2Links, server3Links, server4Links}
		serverIdToLinkIDs := make(map[ObjectId][]ObjectId, len(serverLinks))
		for _, links := range serverLinks {
			// create generic system
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

			err = linkQuery.Do(ctx, &linkQueryResult)
			if err != nil {
				t.Fatal(err)
			}

			if len(linkQueryResult.Items) != 1 {
				t.Fatalf("expected 1 item, got %d items", len(linkQueryResult.Items))
			}

			serverIdToLinkIDs[linkQueryResult.Items[0].Generic.Id] = linkIds

		}

		lagCount := func(ctx context.Context, serverId ObjectId) int {
			links, err := bpClient.GetCablingMapLinksBySystem(ctx, serverId)
			if err != nil {
				t.Fatal(err)
			}
			var lagCount int
			for _, link := range links {
				if link.Type == LinkTypeAggregateLink {
					lagCount++
				}
			}
			return lagCount
		}

		// one big LAG
		for serverId, linkIds := range serverIdToLinkIDs {
			groupLabel := "bond0"
			request := make(SetLinkLagParamsRequest)
			for _, linkId := range linkIds {
				params := LinkLagParams{
					GroupLabel: &groupLabel,
					LagMode:    RackLinkLagModeActive,
					Tags:       []string{"a", "b"},
				}
				request[linkId] = params
			}

			err = bpClient.SetLinkLagParams(ctx, &request)
			if err != nil {
				t.Fatal(err)
			}

			lc := lagCount(ctx, serverId)
			if lc != 1 {
				t.Fatalf("expected 1 LAG, got %d LAGs", lc)
			}
		}

		// no LAG
		for serverId, linkIds := range serverIdToLinkIDs {
			request := make(SetLinkLagParamsRequest)
			for _, linkId := range linkIds {
				params := LinkLagParams{
					LagMode: RackLinkLagModeNone,
					Tags:    []string{"no lag", "still no lag"},
				}
				request[linkId] = params
			}

			err = bpClient.SetLinkLagParams(ctx, &request)
			if err != nil {
				t.Fatal(err)
			}

			lc := lagCount(ctx, serverId)
			if lc != 0 {
				t.Fatalf("expected 0 LAG, got %d LAGs", lc)
			}
		}

		for id := range serverIdToLinkIDs {
			log.Printf("testing DeleteGenericSystem() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.DeleteGenericSystem(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
