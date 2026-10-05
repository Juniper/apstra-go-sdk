// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package resourcetestobj

import (
	"context"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func RandomIPv4Pool(t testing.TB, ctx context.Context, client *apstra.Client) string {
	t.Helper()

	subnetCount := 3

	subnets := testutils.RandomPrefixes(t, "10.0.0.0/8", 27, 3)

	request := apstra.NewIpPoolRequest{
		DisplayName: testutils.RandString(6, "hex"),
		Subnets:     make([]apstra.NewIpSubnet, subnetCount),
	}

	for i := range subnets {
		request.Subnets[i] = apstra.NewIpSubnet{Network: subnets[i].String()}
	}

	id, err := client.CreateIp4Pool(ctx, &request)
	require.NoError(t, err)

	testutils.CleanupWithFreshContext(t, 10*time.Second, func(ctx context.Context) error {
		return client.DeleteIp4Pool(ctx, id)
	})

	return string(id)
}
