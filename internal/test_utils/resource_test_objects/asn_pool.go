// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package resourcetestobj

import (
	"context"
	"math"
	"slices"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func RandomASNPool(t testing.TB, ctx context.Context, client *apstra.Client) string {
	t.Helper()

	rangeCount := 3

	ints, err := testutils.GetRandInts(1, math.MaxUint32, 2*rangeCount)
	require.NoError(t, err)

	slices.Sort(ints)

	request := apstra.AsnPoolRequest{
		DisplayName: testutils.RandString(6, "hex"),
		Ranges:      make([]apstra.IntfIntRange, rangeCount),
	}

	for i := range rangeCount {
		request.Ranges[i] = apstra.IntRange{
			First: uint32(ints[i*2]),
			Last:  uint32(ints[i*2+1]),
		}
	}

	id, err := client.CreateAsnPool(ctx, &request)
	require.NoError(t, err)

	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteAsnPool(ctx, id)
	})

	return string(id)
}
