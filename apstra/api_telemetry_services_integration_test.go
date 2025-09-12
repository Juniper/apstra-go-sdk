// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"bytes"
	"log"
	"testing"

	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetTelemetryServicesDeviceMapping(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			result, err := client.Client.GetTelemetryServicesDeviceMapping(ctx)
			require.NoError(t, err)

			buf := bytes.NewBuffer([]byte{})
			err = testutils.PrettyPrint(result, buf)
			require.NoError(t, err)

			log.Print(buf.String())
		})
	}
}
