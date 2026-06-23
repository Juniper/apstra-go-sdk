// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package dctestobj

import (
	"context"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/query"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestVirtualNetworkA(t testing.TB, ctx context.Context, bp *apstra.TwoStageL3ClosClient, szId string) string {
	t.Helper()

	leafIds, err := query.SystemIdsByRole(ctx, bp, "leaf")
	require.NoError(t, err)

	vnBindings := make([]datacenter.VNBinding, len(leafIds))
	for i, leafId := range leafIds {
		vnBindings[i] = datacenter.VNBinding{SystemID: string(leafId)}
	}

	id, err := bp.CreateVirtualNetwork(ctx, datacenter.VirtualNetwork{
		IPv4Enabled:               true,
		Label:                     testutils.RandString(6, "hex"),
		SecurityZoneID:            szId,
		VirtualGatewayIPv4Enabled: true,
		Bindings:                  vnBindings,
		Type:                      enum.VnTypeVxlan,
	})
	require.NoError(t, err)
	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return bp.DeleteVirtualNetwork(ctx, id)
	})

	return id
}
