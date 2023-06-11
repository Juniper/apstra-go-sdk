//go:build integration
// +build integration

package apstra

import (
	"context"
	"errors"
	"log"
	"testing"
)

func TestGetSetTransformationId(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintF(ctx, t, client.client)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err := bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		leafQuery := new(PathQuery).
			SetBlueprintId(bpClient.blueprintId).
			SetBlueprintType(BlueprintTypeStaging).
			SetClient(bpClient.client).
			Node([]QEEAttribute{
				NodeTypeSystem.QEEAttribute(),
				{"system_type", QEStringVal("switch")},
				{"role", QEStringVal("leaf")},
				{"name", QEStringVal("n_leaf")},
			})
		var leafResponse struct {
			Items []struct {
				Leaf struct {
					ID ObjectId `json:"id"`
				} `json:"n_leaf"`
			} `json:"items"`
		}
		err = leafQuery.Do(ctx, &leafResponse)
		if err != nil {
			t.Fatal(err)
		}

		if len(leafResponse.Items) != 1 {
			t.Fatalf("expected 1 leaf, got %d", len(leafResponse.Items))
		}

		leafId := leafResponse.Items[0].Leaf.ID
		leafAssignements := SystemIdToInterfaceMapAssignment{leafId.String(): "Cisco_3172PQ_NXOS__AOS-48x10_6x40-1"}
		err = bpClient.SetInterfaceMapAssignments(ctx, leafAssignements)
		if err != nil {
			t.Fatal(err)
		}

		ifName := "Ethernet1/1"
		linkIds, err := bpClient.CreateLinksWithNewServer(ctx, &CreateLinksWithNewServerRequest{
			Links: []CreateLinkRequest{
				{
					SwitchEndpoint: SwitchLinkEndpoint{
						TransformationId: 1,
						SystemId:         leafId,
						IfName:           ifName,
					},
				},
			},
			Server: CreateLinksWithNewServerRequestServer{
				LogicalDeviceId: "AOS-1x10-1",
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(linkIds) != 1 {
			t.Fatalf("expected 1 link ID, got %d", len(linkIds))
		}
		linkId := linkIds[0]

		interfaceQuery := new(PathQuery).
			SetBlueprintType(BlueprintTypeStaging).
			SetBlueprintId(bpClient.blueprintId).
			SetClient(bpClient.client).
			Node([]QEEAttribute{
				NodeTypeLink.QEEAttribute(),
				{"id", QEStringVal(linkId)},
			}).
			In([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeInterface.QEEAttribute(),
				{"name", QEStringVal("n_interface")},
			}).
			In([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeSystem.QEEAttribute(),
				{"system_type", QEStringVal("switch")},
			})

		var interfaceQueryResponse struct {
			Items []struct {
				Interface struct {
					Id ObjectId `json:"id"`
				} `json:"n_interface"`
			} `json:"items"`
		}

		err = interfaceQuery.Do(ctx, &interfaceQueryResponse)
		if err != nil {
			t.Fatal(err)
		}
		if len(interfaceQueryResponse.Items) != 1 {
			t.Fatalf("expected 1 result, got %d", len(interfaceQueryResponse.Items))
		}

		type testCase struct {
			transformId int
			expSetIdErr bool
		}

		testCases := []testCase{
			{2, false},
			{1, false},
			{2, false},
			{3, true},
		}

		for i, tc := range testCases {
			log.Printf("testing SetTransformationId() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			request := SetTransformationRequest{
				Force: false,
				Interface: struct {
					TransformationId int
					SystemId         ObjectId
					IfName           string
				}{
					TransformationId: tc.transformId,
					SystemId:         leafId,
					IfName:           ifName,
				},
			}
			err = bpClient.SetTransformationId(ctx, &request)
			if tc.expSetIdErr && err == nil {
				t.Fatalf("test case %d should have produced an error", i)
			}
			if !tc.expSetIdErr && err != nil {
				t.Fatalf("test case %d: %s", i, err.Error())
			}

			if tc.expSetIdErr {
				var ace ApstraClientErr
				if errors.As(err, &ace) && ace.Type() == ErrCannotChangeTransform {
					log.Println("got the error we wanted:", err)
					continue
				}
				t.Fatalf("test case %d: these aren't the errors you're looking for %s", i, err)
			}

			log.Printf("testing GetTransformationId() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			gtidResult, err := bpClient.GetTransformationId(ctx, interfaceQueryResponse.Items[0].Interface.Id)
			if err != nil {
				t.Fatalf("test case %d: %s", i, err.Error())
			}
			if tc.transformId != gtidResult {
				t.Fatalf("test case %d: expected transform id %d, got %d", i, tc.transformId, gtidResult)
			}

			log.Printf("testing GetTransformationIdByIfName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			gtidResult, err = bpClient.GetTransformationIdByIfName(ctx, leafId, ifName)
			if err != nil {
				t.Fatalf("test case %d: %s", i, err.Error())
			}
			if tc.transformId != gtidResult {
				t.Fatalf("test case %d: expected transform id %d, got %d", i, tc.transformId, gtidResult)
			}
		}
	}
}
