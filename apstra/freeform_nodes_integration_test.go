// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFfNodes(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			intSystemCount := rand.IntN(5) + 2
			extSystemCount := rand.IntN(5) + 2

			bp, intSysIds, extSysIds := testFFBlueprintB(ctx, t, client.client, intSystemCount, extSystemCount)
			require.Equal(t, intSystemCount, len(intSysIds))
			require.Equal(t, extSystemCount, len(extSysIds))

			var target struct {
				Nodes map[ObjectId]struct {
					Id         ObjectId `json:"id"`
					Type       string   `json:"type"`
					SystemType string   `json:"system_type"`
				} `json:"nodes"`
			}
			err := bp.GetNodes(ctx, NodeTypeSystem, &target)
			require.NoError(t, err)
			require.Equal(t, len(target.Nodes), intSystemCount+extSystemCount)

			for id, node := range target.Nodes {
				require.Equal(t, node.Type, NodeTypeSystem.String())
				require.Contains(t, []string{SystemTypeInternal.String(), systemTypeExternal.string()}, node.SystemType)
				require.Equal(t, id, node.Id)
				switch node.SystemType {
				case SystemTypeInternal.String():
					require.Contains(t, intSysIds, node.Id)
				case SystemTypeExternal.String():
					require.Contains(t, extSysIds, node.Id)
				default:
					t.Fatalf("it should have been impossible to get here")
				}
			}
		})
	}
}
