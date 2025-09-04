// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"bytes"
	"context"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"log"
	"testing"

	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestGetTelemetryServicesDeviceMapping(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			result, err := client.Client.GetTelemetryServicesDeviceMapping(ctx)
			require.NoError(t, err)

			buf := bytes.NewBuffer([]byte{})
			err = testutils.PrettyPrint(result, buf)
			require.NoError(t, err)

			log.Print(buf.String())
		})
	}
}
