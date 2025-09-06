// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"log"
	"math/rand/v2"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetSetSystemAgentManagerConfiguration(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			// initial fetch
			mgrCfg, err := client.Client.GetSystemAgentManagerConfig(ctx)
			require.NoError(t, err)

			// new config with opposite values
			testCfg := &apstra.SystemAgentManagerConfig{
				SkipRevertToPristineOnUninstall: !mgrCfg.SkipRevertToPristineOnUninstall,
				SkipPristineValidation:          !mgrCfg.SkipPristineValidation,
				SkipInterfaceShutdownOnUpgrade:  !mgrCfg.SkipInterfaceShutdownOnUpgrade,
			}
			if compatibility.HasDeviceOsImageDownloadTimeout.Check(client.APIVersion()) {
				testCfg.DeviceOsImageDownloadTimeout = testutils.ToPtr(rand.IntN(2700) + 1)
			}

			// set new config
			err = client.Client.SetSystemAgentManagerConfig(ctx, testCfg)
			if !compatibility.SystemManagerHasSkipInterfaceShutdownOnUpgrade.Check(client.APIVersion()) {
				require.Error(t, err)
				log.Printf("apstra %s refused to run with SkipInterfaceShutdownOnUpgrade set to %t", client.Client.ApiVersion(), testCfg.SkipInterfaceShutdownOnUpgrade)

				testCfg.SkipInterfaceShutdownOnUpgrade = false
				err = client.Client.SetSystemAgentManagerConfig(ctx, testCfg)
			}
			require.NoError(t, err)

			// fetch new config
			testCfgRetrieved, err := client.Client.GetSystemAgentManagerConfig(ctx)
			require.NoError(t, err)

			// validate field as expected
			require.Equal(t, testCfg.SkipPristineValidation, testCfgRetrieved.SkipPristineValidation)
			require.Equal(t, testCfg.SkipInterfaceShutdownOnUpgrade, testCfgRetrieved.SkipInterfaceShutdownOnUpgrade)
			require.Equal(t, testCfg.SkipRevertToPristineOnUninstall, testCfgRetrieved.SkipRevertToPristineOnUninstall)

			// reset to original configuration
			err = client.Client.SetSystemAgentManagerConfig(ctx, mgrCfg)
			require.NoError(t, err)
		})
	}
}
