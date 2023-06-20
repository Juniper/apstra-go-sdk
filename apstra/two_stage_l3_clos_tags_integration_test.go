//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"sort"
	"testing"
)

func TestNodeTagging(t *testing.T) {
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

		query := new(PathQuery).
			SetBlueprintId(bpClient.blueprintId).
			SetBlueprintType(BlueprintTypeStaging).
			SetClient(bpClient.client).
			Node([]QEEAttribute{
				NodeTypeSystem.QEEAttribute(),
				{"label", QEStringVal("spine1")},
				{"name", QEStringVal("n_system")},
			})

		var queryResponse struct {
			Items []struct {
				System struct {
					Id ObjectId `json:"id"`
				} `json:"n_system"`
			} `json:"items""`
		}

		err = query.Do(ctx, &queryResponse)
		if err != nil {
			t.Fatal(err)
		}

		if len(queryResponse.Items) != 1 {
			t.Fatalf("expected 1 query match for 'spine 1', got %d", len(queryResponse.Items))
		}

		spine1 := queryResponse.Items[0].System.Id

		type testCase struct {
			tags []string
		}

		testCases := []testCase{
			{tags: []string{}},
			{tags: []string{"a"}},
			{tags: []string{"a", "b"}},
			{tags: []string{"a"}},
			{tags: []string{}},
			{tags: []string{"c", "d"}},
			{tags: []string{"c", "d"}},
			{tags: []string{"a"}},
		}

		for i, tc := range testCases {
			log.Printf("testing SetNodeTags() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetNodeTags(ctx, spine1, tc.tags)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetNodeTags() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			tags, err := bpClient.GetNodeTags(ctx, spine1)
			if err != nil {
				t.Fatal(err)
			}

			sort.Strings(tc.tags)
			sort.Strings(tags)
			compareSlices(t, tc.tags, tags, fmt.Sprintf("test case %d spine 1 tags expected vs. reality", i))
		}
	}
}
