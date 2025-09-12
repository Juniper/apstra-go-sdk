// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetVersionsAll(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			aosdi, err := client.Client.GetVersionsAosdi(ctx)
			require.NoError(t, err)
			if err != nil {
				t.Fatal(err)
			}
			api, err := client.Client.GetVersionsApi(ctx)
			require.NoError(t, err)

			build, err := client.Client.GetVersionsBuild(ctx)
			require.NoError(t, err)

			server, err := client.Client.GetVersionsServer(ctx)
			require.NoError(t, err)

			body, err := json.Marshal(&struct {
				Aosdi  *apstra.VersionsAosdiResponse  `json:"aosdi"`
				Api    *apstra.VersionsApiResponse    `json:"api"`
				Build  *apstra.VersionsBuildResponse  `json:"build"`
				Server *apstra.VersionsServerResponse `json:"server"`
			}{
				Aosdi:  aosdi,
				Api:    api,
				Build:  build,
				Server: server,
			})
			require.NoError(t, err)

			log.Println(string(body))
		})
	}
}
