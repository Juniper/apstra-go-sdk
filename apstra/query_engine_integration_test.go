// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestParsingQueryInfo(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			bpClient := dctestobj.TestBlueprintA(t, ctx, client.Client)

			// the type of info we expect the query to return (a slice of these)
			var qResponse struct {
				Count int `json:"count"`
				Items []struct {
					LogicalDevice struct {
						Id    string `json:"id"`
						Label string `json:"label"`
					} `json:"n_logical_device"`
					System struct {
						Id    string `json:"id"`
						Label string `json:"label"`
					} `json:"n_system"`
				} `json:"items"`
			}

			err := new(apstra.PathQuery).
				SetClient(bpClient.Client()).
				SetBlueprintId(bpClient.Id()).
				Node([]apstra.QEEAttribute{
					{"type", apstra.QEStringVal("system")},
					{"name", apstra.QEStringVal("n_system")},
					{"role", apstra.QEStringValIsIn{"superspine", "spine", "leaf"}},
					{"external", apstra.QEBoolVal(false)},
				}).
				Out([]apstra.QEEAttribute{
					{"type", apstra.QEStringVal("logical_device")},
				}).
				Node([]apstra.QEEAttribute{
					{"type", apstra.QEStringVal("logical_device")},
					{"name", apstra.QEStringVal("n_logical_device")},
				}).
				Do(ctx, &qResponse)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("query produced %d results", qResponse.Count)
			for i, item := range qResponse.Items {
				log.Printf("  %d id: '%s', label: '%s', logical_device: '%s'", i, item.System.Id, item.System.Label, item.LogicalDevice.Label)
			}
		})
	}
}

func TestRawQueryWithBlueprint(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			bpClient := dctestobj.TestBlueprintA(t, ctx, client.Client)

			query := new(apstra.RawQuery).
				SetBlueprintType(apstra.BlueprintTypeStaging).
				SetClient(client.Client).
				SetBlueprintId(bpClient.Id()).
				SetQuery("node(type='system', role='leaf', name='n_system')")

			var queryResponse struct {
				Count int `json:"count"`
				Items []struct {
					System struct {
						Id    string `json:"id"`
						Label string `json:"label"`
					} `json:"n_system"`
				} `json:"items"`
			}

			err := query.Do(ctx, &queryResponse)
			require.NoError(t, err)

			qr1 := queryResponse

			err = json.Unmarshal(query.RawResult(), &queryResponse)
			require.NoError(t, err)

			qr2 := queryResponse

			require.Equal(t, qr1, qr2)
		})
	}
}
