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

func TestRouringPolicyA(t testing.TB, ctx context.Context, bp *apstra.TwoStageL3ClosClient) string {
	id, err := bp.CreateRoutingPolicy(ctx, &apstra.DcRoutingPolicyData{
		Label:        testutils.RandString(6, "hex"),
		Description:  testutils.RandString(6, "hex"),
		PolicyType:   apstra.DcRoutingPolicyTypeUser,
		ImportPolicy: apstra.DcRoutingPolicyImportPolicyAll,
	})
	require.NoError(t, err)

	return string(id)
}
