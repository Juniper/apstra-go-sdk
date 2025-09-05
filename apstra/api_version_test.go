// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"encoding/json"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestGetVersion(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, t.Context())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, t.Context())

			ver, err := client.Client.GetVersion(ctx)
			require.NoError(t, err)

			result, err := json.Marshal(ver)
			require.NoError(t, err)

			log.Println(string(result))
		})
	}
}
