// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package fftestobj

import (
	"context"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

// TestBlueprintA is an empty FreeForm blueprint
func TestBlueprintA(t testing.TB, ctx context.Context, client *apstra.Client) *apstra.FreeformClient {
	t.Helper()

	bpId, err := client.CreateFreeformBlueprint(ctx, testutils.RandString(6, "hex"))
	require.NoError(t, err)

	bpClient, err := client.NewFreeformClient(ctx, bpId)
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	})

	return bpClient
}
