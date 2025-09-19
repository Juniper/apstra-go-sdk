// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package dctestobj

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestSecurityZoneA(t testing.TB, ctx context.Context, bp *apstra.TwoStageL3ClosClient) apstra.ObjectId {
	t.Helper()

	rs := testutils.RandString(6, "hex")

	id, err := bp.CreateSecurityZone(ctx, &apstra.SecurityZoneData{
		Label:            rs,
		SzType:           apstra.SecurityZoneTypeEVPN,
		VrfName:          rs,
		RoutingPolicyId:  "",
		RouteTarget:      nil,
		RtPolicy:         nil,
		VlanId:           nil,
		VniId:            nil,
		JunosEvpnIrbMode: nil,
	})
	require.NoError(t, err)

	return id
}
