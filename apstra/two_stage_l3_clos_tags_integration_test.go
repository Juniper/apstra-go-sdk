// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNodeTagging(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bpClient := testBlueprintD(ctx, t, client.client)

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
				} `json:"items"`
			}

			require.NoError(t, query.Do(ctx, &queryResponse))
			require.Equal(t, 1, len(queryResponse.Items))

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
				require.NoError(t, bpClient.SetNodeTags(ctx, spine1, tc.tags))

				log.Printf("testing GetNodeTags() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				tags, err := bpClient.GetNodeTags(ctx, spine1)
				require.NoError(t, err)

				sort.Strings(tc.tags)
				sort.Strings(tags)
				compareSlices(t, tc.tags, tags, fmt.Sprintf("test case %d spine 1 tags expected vs. reality", i))
			}
		})
	}
}
