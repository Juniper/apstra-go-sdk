// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestModularDeviceProfile(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			mdp1 := &apstra.ModularDeviceProfile{
				Label:            testutils.RandString(5, "hex"),
				ChassisProfileId: "Juniper_QFX10016",
				SlotConfigurations: map[uint64]apstra.ModularDeviceSlotConfiguration{
					0: {LinecardProfileId: "Juniper_QFX10000_30C_M"},
					2: {LinecardProfileId: "Juniper_QFX10000_30C_M"},
					4: {LinecardProfileId: "Juniper_QFX10000_30C_M"},
				},
			}

			id, err := client.Client.CreateModularDeviceProfile(ctx, mdp1)
			require.NoError(t, err)

			mdp2, err := client.Client.GetModularDeviceProfile(ctx, id)
			require.NoError(t, err)

			require.Equal(t, *mdp1, *mdp2)

			mdp1.Label = testutils.RandString(5, "hex")
			mdp1.ChassisProfileId = "Juniper_QFX10016"
			mdp1.SlotConfigurations = map[uint64]apstra.ModularDeviceSlotConfiguration{
				1: {LinecardProfileId: "Juniper_QFX10000_30C"},
				3: {LinecardProfileId: "Juniper_QFX10000_30C"},
				5: {LinecardProfileId: "Juniper_QFX10000_30C"},
			}
			err = client.Client.UpdateModularDeviceProfile(ctx, id, mdp1)
			require.NoError(t, err)

			mdp2, err = client.Client.GetModularDeviceProfile(ctx, id)
			require.NoError(t, err)
			require.Equal(t, *mdp1, *mdp2)

			err = client.Client.DeleteModularDeviceProfile(ctx, id)
			require.NoError(t, err)
		})
	}
}
