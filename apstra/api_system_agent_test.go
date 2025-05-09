// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"math/rand/v2"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/stretchr/testify/require"
)

func TestGetSetSystemAgentManagerConfiguration(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	require.NoError(t, err)

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			// initial fetch
			log.Printf("testing GetSystemAgentManagerConfig() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			mgrCfg, err := client.client.GetSystemAgentManagerConfig(context.Background())
			require.NoError(t, err)

			// new config with opposite values
			testCfg := &SystemAgentManagerConfig{
				SkipRevertToPristineOnUninstall: !mgrCfg.SkipRevertToPristineOnUninstall,
				SkipPristineValidation:          !mgrCfg.SkipPristineValidation,
				SkipInterfaceShutdownOnUpgrade:  !mgrCfg.SkipInterfaceShutdownOnUpgrade,
			}
			if compatibility.HasDeviceOsImageDownloadTimeout.Check(client.client.apiVersion) {
				testCfg.DeviceOsImageDownloadTimeout = toPtr(rand.IntN(2700) + 1)
			}

			// set new config
			log.Printf("testing SetSystemAgentManagerConfig() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.SetSystemAgentManagerConfig(context.Background(), testCfg)
			if !compatibility.SystemManagerHasSkipInterfaceShutdownOnUpgrade.Check(client.client.apiVersion) {
				require.Error(t, err)
				log.Printf("apstra %s refused to run with SkipInterfaceShutdownOnUpgrade set to %t", client.client.apiVersion, testCfg.SkipInterfaceShutdownOnUpgrade)

				testCfg.SkipInterfaceShutdownOnUpgrade = false
				err = client.client.SetSystemAgentManagerConfig(context.Background(), testCfg)
			}
			require.NoError(t, err)

			// fetch new config
			log.Printf("testing GetSystemAgentManagerConfig() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			testCfgRetrieved, err := client.client.GetSystemAgentManagerConfig(context.Background())
			require.NoError(t, err)

			// validate field as expected
			require.Equal(t, testCfg.SkipPristineValidation, testCfgRetrieved.SkipPristineValidation)
			require.Equal(t, testCfg.SkipInterfaceShutdownOnUpgrade, testCfgRetrieved.SkipInterfaceShutdownOnUpgrade)
			require.Equal(t, testCfg.SkipRevertToPristineOnUninstall, testCfgRetrieved.SkipRevertToPristineOnUninstall)

			// reset to original configuration
			log.Printf("testing SetSystemAgentManagerConfig() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.SetSystemAgentManagerConfig(context.Background(), mgrCfg)
			require.NoError(t, err)
		})
	}
}
