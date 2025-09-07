// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"bytes"
	"context"
	"log"
	"testing"

	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetVirtualInfraMgrs(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			vim, err := client.Client.GetVirtualInfraMgrs(ctx)
			require.NoError(t, err)

			buf := bytes.NewBuffer([]byte{})
			err = testutils.PrettyPrint(vim, buf)
			require.NoError(t, err)

			log.Println(buf.String())
		})
	}
}
