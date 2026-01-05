// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package dctestobj

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestSecurityZoneA(t testing.TB, ctx context.Context, bp *apstra.TwoStageL3ClosClient) string {
	t.Helper()

	rs := testutils.RandString(6, "hex")

	id, err := bp.CreateSecurityZone(ctx, apstra.SecurityZone{
		Label:            rs,
		Type:             enum.SecurityZoneTypeEVPN,
		VRFName:          rs,
		RoutingPolicyID:  "",
		RouteTarget:      nil,
		RTPolicy:         nil,
		VLAN:             nil,
		VNI:              nil,
		JunosEVPNIRBMode: nil,
	})
	require.NoError(t, err)

	return id
}
