// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package dctestobj

import (
	"context"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestVirtualNetworkA(t testing.TB, ctx context.Context, bp *apstra.TwoStageL3ClosClient, szId apstra.ObjectId) apstra.ObjectId {
	t.Helper()

	leafIds, err := testutils.GetSystemIdsByRole(ctx, bp, "leaf")
	require.NoError(t, err)

	vnBindings := make([]apstra.VnBinding, len(leafIds))
	for i, leafId := range leafIds {
		vnBindings[i] = apstra.VnBinding{SystemId: leafId}
	}

	id, err := bp.CreateVirtualNetwork(ctx, &apstra.VirtualNetworkData{
		Ipv4Enabled:               true,
		Label:                     testutils.RandString(6, "hex"),
		SecurityZoneId:            szId,
		VirtualGatewayIpv4Enabled: true,
		VnBindings:                vnBindings,
		VnType:                    enum.VnTypeVxlan,
	})
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return bp.DeleteVirtualNetwork(ctx, id)
	})

	return id
}
