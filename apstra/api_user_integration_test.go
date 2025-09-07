// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"testing"

	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestUserLogin(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			if client.Type() == testclient.ClientTypeAPIOps {
				t.Skipf("skipping test - api-ops type clients do not log in or out")
			}

			var err error

			err = client.Client.Login(ctx)
			require.NoError(t, err)

			err = client.Client.Logout(ctx)
			require.NoError(t, err)

			err = client.Client.Logout(ctx)
			require.NoError(t, err)

			err = client.Client.Login(ctx)
			require.NoError(t, err)
		})
	}
}
