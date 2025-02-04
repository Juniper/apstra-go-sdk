// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_DoRawJsonTransaction(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	parseUrl := func(t *testing.T, urlStr string) *url.URL {
		t.Helper()

		u, err := url.Parse(urlStr)
		require.NoError(t, err)

		return u
	}

	for _, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			var idResponse objectIdResponse

			// create an IP pool
			err = client.client.DoRawJsonTransaction(ctx, RawJsonRequest{
				Method: http.MethodPost,
				Url:    parseUrl(t, "/api/resources/ip-pools"),
				Payload: NewIpPoolRequest{
					DisplayName: randString(6, "hex"),
					Subnets:     []NewIpSubnet{{Network: "10.0.0.0/24"}},
				},
			}, &idResponse)
			require.NoError(t, err)

			var itemsResponse struct {
				Items []ObjectId `json:"items"`
			}

			// ensure the pool ID exists
			err = client.client.DoRawJsonTransaction(ctx, RawJsonRequest{
				Method:  http.MethodOptions,
				Url:     parseUrl(t, "/api/resources/ip-pools"),
				Payload: nil,
			}, &itemsResponse)
			require.NoError(t, err)
			require.Contains(t, itemsResponse.Items, idResponse.Id)

			// delete the pool
			err = client.client.DoRawJsonTransaction(ctx, RawJsonRequest{
				Method: http.MethodDelete,
				Url:    parseUrl(t, fmt.Sprintf("/api/resources/ip-pools/%s", idResponse.Id)),
			}, nil)

			// ensure the pool ID does not exist
			err = client.client.DoRawJsonTransaction(ctx, RawJsonRequest{
				Method:  http.MethodOptions,
				Url:     parseUrl(t, "/api/resources/ip-pools"),
				Payload: nil,
			}, &itemsResponse)
			require.NoError(t, err)
			require.NotContains(t, itemsResponse.Items, idResponse.Id)
		})
	}
}
