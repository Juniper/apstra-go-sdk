// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestClient_DoRawJsonTransaction(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	parseUrl := func(t *testing.T, urlStr string) *url.URL {
		t.Helper()

		u, err := url.Parse(urlStr)
		require.NoError(t, err)

		return u
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			var idResponse struct {
				Id string `json:"id"`
			}

			// create an IP pool
			err := client.Client.DoRawJsonTransaction(ctx, apstra.RawJsonRequest{
				Method: http.MethodPost,
				Url:    parseUrl(t, "/api/resources/ip-pools"),
				Payload: apstra.NewIpPoolRequest{
					DisplayName: testutils.RandString(6, "hex"),
					Subnets:     []apstra.NewIpSubnet{{Network: "10.0.0.0/24"}},
				},
			}, &idResponse)
			require.NoError(t, err)

			var itemsResponse struct {
				Items []string `json:"items"`
			}

			// ensure the pool ID exists
			err = client.Client.DoRawJsonTransaction(ctx, apstra.RawJsonRequest{
				Method:  http.MethodOptions,
				Url:     parseUrl(t, "/api/resources/ip-pools"),
				Payload: nil,
			}, &itemsResponse)
			require.NoError(t, err)
			require.Contains(t, itemsResponse.Items, idResponse.Id)

			// delete the pool
			err = client.Client.DoRawJsonTransaction(ctx, apstra.RawJsonRequest{
				Method: http.MethodDelete,
				Url:    parseUrl(t, fmt.Sprintf("/api/resources/ip-pools/%s", idResponse.Id)),
			}, nil)
			require.NoError(t, err)

			// pools don't disappear immediately
			time.Sleep(time.Second)

			// ensure the pool ID does not exist
			err = client.Client.DoRawJsonTransaction(ctx, apstra.RawJsonRequest{
				Method:  http.MethodOptions,
				Url:     parseUrl(t, "/api/resources/ip-pools"),
				Payload: nil,
			}, &itemsResponse)
			require.NoError(t, err)
			require.NotContains(t, itemsResponse.Items, idResponse.Id)
		})
	}
}
