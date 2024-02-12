//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-version"
	"log"
	"math/rand"
	"sort"
	"testing"
)

func TestSetGenericServerBonding(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	addLinksToEsiAccessMinVersion, err := version.NewVersion("4.1.2")
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		clientVersion, err := version.NewVersion(client.client.apiVersion.String())
		if err != nil {
			t.Fatal(err)
		}

		bpClient, bpDel := testBlueprintE(ctx, t, client.client)
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

		type testCase struct {
			switchIds       []ObjectId
			linkCount       int
			firstInterface  string
			bondStrategy    string
			logicalDeviceId string
			systemType      SystemType
		}

		testCases := []testCase{
			{
				switchIds:       []ObjectId{leafIds[0]},
				linkCount:       1,
				firstInterface:  "xe-0/0/5",
				logicalDeviceId: "AOS-4x10-1",
				systemType:      SystemTypeServer,
			},
			{
				switchIds:       []ObjectId{leafIds[0]},
				linkCount:       4,
				firstInterface:  "xe-0/0/5",
				logicalDeviceId: "AOS-4x10-1",
				systemType:      SystemTypeServer,
			},
			{
				switchIds:       []ObjectId{leafIds[0]},
				linkCount:       1,
				firstInterface:  "xe-0/0/5",
				logicalDeviceId: "AOS-4x10-1",
				systemType:      SystemTypeExternal,
			},
			{
				switchIds:       []ObjectId{leafIds[0]},
				linkCount:       4,
				firstInterface:  "xe-0/0/5",
				logicalDeviceId: "AOS-4x10-1",
				systemType:      SystemTypeExternal,
			},
			{
				switchIds:       []ObjectId{accessIds[0]},
				linkCount:       1,
				firstInterface:  "xe-0/0/5",
				logicalDeviceId: "AOS-4x10-1",
				systemType:      SystemTypeServer,
			},
			{
				switchIds:       []ObjectId{accessIds[0]},
				linkCount:       4,
				firstInterface:  "xe-0/0/5",
				logicalDeviceId: "AOS-4x10-1",
				systemType:      SystemTypeServer,
			},
			{
				switchIds:       leafIds,
				linkCount:       4,
				firstInterface:  "xe-0/0/5",
				logicalDeviceId: "AOS-4x10-1",
				systemType:      SystemTypeServer,
			},
			{
				switchIds:       leafIds,
				linkCount:       4,
				firstInterface:  "xe-0/0/5",
				logicalDeviceId: "AOS-4x10-1",
				systemType:      SystemTypeExternal,
			},
			{
				switchIds:       accessIds,
				linkCount:       4,
				firstInterface:  "xe-0/0/5",
				logicalDeviceId: "AOS-4x10-1",
				systemType:      SystemTypeServer,
			},
		}

	TESTCASE:
		for i, tc := range testCases {
			var hostTags []string
			for j := 0; j < rand.Intn(5)+2; j++ {
				hostTags = append(hostTags, randString(5, "hex"))
			}

			if sliceContainsAnyOf(accessIds, tc.switchIds) && clientVersion.LessThan(addLinksToEsiAccessMinVersion) {
				continue TESTCASE // cannot use this test case against early apstra
			}

			var links []CreateLinkRequest
			var ifName string
			for j := 0; j < tc.linkCount; j++ {
				var linkTags []string
				for k := 0; k < rand.Intn(5)+2; k++ {
					linkTags = append(linkTags, randString(5, "hex"))
				}

				switchModulo := j % len(tc.switchIds)

				switchId := tc.switchIds[switchModulo]

				if j == 0 {
					ifName = tc.firstInterface
				} else if switchModulo == 0 {
					ifName = nextInterface(ifName)
				}

				links = append(links, CreateLinkRequest{
					Tags: linkTags,
					SwitchEndpoint: SwitchLinkEndpoint{
						TransformationId: 1,
						SystemId:         switchId,
						IfName:           ifName,
					},
				})
			}

			// poId min/max can only be set for "internal" generic systems
			var portChannelIdMin, portChannelIdMax int
			if tc.systemType == SystemTypeServer {
				portChannelIdMin = rand.Intn(100) + 100
				portChannelIdMax = rand.Intn(100) + 200
			}

			request := CreateLinksWithNewSystemRequest{
				Links: links,
				System: CreateLinksWithNewSystemRequestSystem{
					Hostname:         randString(5, "hex"),
					Label:            randString(5, "hex"),
					LogicalDeviceId:  ObjectId(tc.logicalDeviceId),
					PortChannelIdMin: portChannelIdMin,
					PortChannelIdMax: portChannelIdMax,
					Tags:             hostTags,
					Type:             tc.systemType,
				},
			}
			log.Printf("testing CreateLinksWithNewSystem() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			linkIds, err := bpClient.CreateLinksWithNewSystem(ctx, &request)
			if err != nil {
				t.Fatalf("test case %d - %s", i, err)
			}

			genericId, err := bpClient.SystemNodeFromLinkIds(ctx, linkIds, SystemNodeRoleGeneric)
			if err != nil {
				t.Fatal(err)
			}

			systemTags, err := bpClient.GetNodeTags(ctx, genericId)
			if err != nil {
				t.Fatal(err)
			}

			sort.Strings(systemTags)
			sort.Strings(request.System.Tags)
			compareSlices(t, systemTags, request.System.Tags, fmt.Sprintf("test case %d system tags", i))

			var node struct {
				PoMax int `json:"port_channel_id_max"`
				PoMin int `json:"port_channel_id_min"`
			}
			err = bpClient.client.GetNode(ctx, bpClient.blueprintId, genericId, &node)
			if err != nil {
				t.Fatal(err)
			}
			if request.System.PortChannelIdMin != node.PoMin {
				t.Fatalf("expected port channel id min: %d got %d", request.System.PortChannelIdMin, node.PoMin)
			}
			if request.System.PortChannelIdMax != node.PoMax {
				t.Fatalf("expected port channel id max: %d got %d", request.System.PortChannelIdMax, node.PoMax)
			}

			originalLinkTagDigests := make([]string, len(request.Links))
			for j, link := range request.Links {
				sort.Strings(link.Tags)
				originalLinkTagDigests[j] = fmt.Sprintf("%v", link.Tags)
			}
			sort.Strings(originalLinkTagDigests)

			observedLinkTagDigests := make([]string, len(linkIds))
			for j, linkId := range linkIds {
				tags, err := bpClient.GetNodeTags(ctx, linkId)
				if err != nil {
					t.Fatal(err)
				}
				sort.Strings(tags)
				observedLinkTagDigests[j] = fmt.Sprintf("%v", tags)
			}
			sort.Strings(observedLinkTagDigests)

			compareSlices(t, originalLinkTagDigests, observedLinkTagDigests, "link tag digests")

			// set individual port channels
			lagRequest := make(SetLinkLagParamsRequest)
			for _, linkId := range linkIds {
				lagRequest[linkId] = LinkLagParams{
					LagMode: RackLinkLagModeActive,
				}
			}
			err = bpClient.SetLinkLagParams(ctx, &lagRequest)
			if err != nil {
				t.Fatal(err)
			}

			typeCount, lagMemberCount, err := countSystemLinkTypes(ctx, genericId, bpClient)
			if err != nil {
				t.Fatal(err)
			}
			if typeCount[LinkTypeAggregateLink] != typeCount[LinkTypeEthernet] {
				t.Fatalf("expected count of aggregate to match count of ethernet, got %d and %d",
					typeCount[LinkTypeAggregateLink], typeCount[LinkTypeEthernet])
			}
			if len(linkIds) != lagMemberCount {
				t.Fatalf("expected %d lag member links, got %d", len(linkIds), lagMemberCount)
			}

			observedLinkTagDigests = make([]string, len(linkIds))
			for j, linkId := range linkIds {
				tags, err := bpClient.GetNodeTags(ctx, linkId)
				if err != nil {
					t.Fatal(err)
				}
				sort.Strings(tags)
				observedLinkTagDigests[j] = fmt.Sprintf("%v", tags)
			}
			sort.Strings(observedLinkTagDigests)

			compareSlices(t, originalLinkTagDigests, observedLinkTagDigests, "link tag digests")

			// create one big port channel and wipe out link tags
			lagRequest = make(SetLinkLagParamsRequest)
			for _, linkId := range linkIds {
				lagRequest[linkId] = LinkLagParams{
					GroupLabel: "one big lag",
					LagMode:    RackLinkLagModeActive,
					Tags:       []string{},
				}
			}
			err = bpClient.SetLinkLagParams(ctx, &lagRequest)
			if err != nil {
				t.Fatal(err)
			}

			typeCount, lagMemberCount, err = countSystemLinkTypes(ctx, genericId, bpClient)
			if err != nil {
				t.Fatal(err)
			}
			if typeCount[LinkTypeAggregateLink] != 1 {
				t.Fatal("expected one big LAG")
			}
			if len(linkIds) != lagMemberCount {
				t.Fatalf("expected %d lag member links, got %d", len(linkIds), lagMemberCount)
			}

			for _, linkId := range linkIds {
				tags, err := bpClient.GetNodeTags(ctx, linkId)
				if err != nil {
					t.Fatal(err)
				}
				if len(tags) != 0 {
					t.Fatalf("expected no link tags, got %v", tags)
				}
			}

			// create port channels with no more than two links
			if len(linkIds) < 2 {
				err = bpClient.DeleteGenericSystem(ctx, genericId)
				if err != nil {
					t.Fatal(err)
				}
				continue TESTCASE
			}

			lagRequest = make(SetLinkLagParamsRequest)
			for i, linkId := range linkIds {
				lagRequest[linkId] = LinkLagParams{
					GroupLabel: fmt.Sprintf("paired links %d", i/2),
					LagMode:    RackLinkLagModeActive,
					Tags:       []string{"paired links", fmt.Sprintf("link %d of the pair", i%2)},
				}
			}
			err = bpClient.SetLinkLagParams(ctx, &lagRequest)
			if err != nil {
				t.Fatal(err)
			}

			typeCount, lagMemberCount, err = countSystemLinkTypes(ctx, genericId, bpClient)
			if err != nil {
				t.Fatal(err)
			}
			if typeCount[LinkTypeAggregateLink] != typeCount[LinkTypeEthernet]/2 {
				t.Fatalf("expected half as many aggregate links as ethernet, got %d and %d",
					typeCount[LinkTypeAggregateLink], typeCount[LinkTypeEthernet])
			}
			if len(linkIds) != lagMemberCount {
				t.Fatalf("expected %d lag member links, got %d", len(linkIds), lagMemberCount)
			}

			observedLinkTagDigests = make([]string, len(linkIds))
			for j, linkId := range linkIds {
				tags, err := bpClient.GetNodeTags(ctx, linkId)
				if err != nil {
					t.Fatal(err)
				}
				sort.Strings(tags)
				observedLinkTagDigests[j] = fmt.Sprintf("%v", tags)
			}
			sort.Strings(observedLinkTagDigests)

			expectedLinkTagDigests := make([]string, len(lagRequest))
			var j int
			for _, params := range lagRequest {
				sort.Strings(params.Tags)
				expectedLinkTagDigests[j] = fmt.Sprintf("%v", params.Tags)
				j++
			}
			sort.Strings(expectedLinkTagDigests)

			compareSlices(t, expectedLinkTagDigests, observedLinkTagDigests, "link tag digests")

			// disable LAG
			lagRequest = make(SetLinkLagParamsRequest)
			for _, linkId := range linkIds {
				lagRequest[linkId] = LinkLagParams{
					LagMode: RackLinkLagModeNone,
				}
			}
			err = bpClient.SetLinkLagParams(ctx, &lagRequest)
			if err != nil {
				t.Fatal(err)
			}

			typeCount, lagMemberCount, err = countSystemLinkTypes(ctx, genericId, bpClient)
			if err != nil {
				t.Fatal(err)
			}
			if typeCount[LinkTypeAggregateLink] != 0 {
				t.Fatalf("expected 0 LAGs got %d", typeCount[LinkTypeAggregateLink])
			}
			if 0 != lagMemberCount {
				t.Fatalf("expected 0 lag member links, got %d", lagMemberCount)
			}

			// delete the server
			err = bpClient.DeleteGenericSystem(ctx, genericId)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
