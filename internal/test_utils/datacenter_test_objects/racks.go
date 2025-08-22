// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package dctestobj

import (
	"context"
	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRackA(t testing.TB, ctx context.Context, client *apstra.Client) apstra.ObjectId {
	t.Helper()

	request := apstra.RackTypeRequest{
		DisplayName:              testutils.RandString(5, "hex"),
		FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
		LeafSwitches: []apstra.RackElementLeafSwitchRequest{
			{
				Label:             testutils.RandString(5, "hex"),
				LinkPerSpineCount: 1,
				LinkPerSpineSpeed: "40G",
				LogicalDeviceId:   "AOS-48x10_6x40-1",
			},
		},
	}

	id, err := client.CreateRackType(ctx, &request)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, client.DeleteRackType(ctx, id))
	})

	return id
}
