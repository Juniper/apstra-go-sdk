// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func compareAntiAffinityPolicy(t testing.TB, set, get AntiAffinityPolicy) {
	t.Helper()

	if set.Algorithm != get.Algorithm {
		t.Errorf("set AntiAffinityPolicy Algorithm %s got %s", set.Algorithm, get.Algorithm)
	}

	if set.MaxLinksPerPort != get.MaxLinksPerPort {
		t.Errorf("set AntiAffinityPolicy MaxLinksPerPort %d got %d", set.MaxLinksPerPort, get.MaxLinksPerPort)
	}

	if set.MaxLinksPerSlot != get.MaxLinksPerSlot {
		t.Errorf("set AntiAffinityPolicy MaxLinksPerSlot %d got %d", set.MaxLinksPerSlot, get.MaxLinksPerSlot)
	}

	if set.MaxPerSystemLinksPerPort != get.MaxPerSystemLinksPerPort {
		t.Errorf("set AntiAffinityPolicy MaxPerSystemLinksPerPort %d got %d", set.MaxPerSystemLinksPerPort, get.MaxPerSystemLinksPerPort)
	}

	if set.MaxPerSystemLinksPerSlot != get.MaxPerSystemLinksPerSlot {
		t.Errorf("set AntiAffinityPolicy MaxPerSystemLinksPerSlot %d got %d", set.MaxPerSystemLinksPerSlot, get.MaxPerSystemLinksPerSlot)
	}

	if set.Mode != get.Mode {
		t.Errorf("set AntiAffinityPolicy Mode %s got %s", set.Mode, get.Mode)
	}
}

func compareFabricSettings(t testing.TB, set, get FabricSettings) {
	t.Helper()

	if set.AntiAffinityPolicy != nil {
		require.NotNil(t, get.AntiAffinityPolicy)
		compareAntiAffinityPolicy(t, *get.AntiAffinityPolicy, *set.AntiAffinityPolicy)
	}

	if set.DefaultSviL3Mtu != nil {
		require.NotNil(t, get.DefaultSviL3Mtu)
		require.Equalf(t, *set.DefaultSviL3Mtu, *get.DefaultSviL3Mtu, "DefaultSviL3Mtu: set %d get %d", *set.DefaultSviL3Mtu, *get.DefaultSviL3Mtu)
	}

	if set.EsiMacMsb != nil {
		require.NotNil(t, get.EsiMacMsb)
		require.Equalf(t, *set.EsiMacMsb, *get.EsiMacMsb, "EsiMacMsb: set %d get %d", *set.EsiMacMsb, *get.EsiMacMsb)
	}

	if set.EvpnGenerateType5HostRoutes != nil && *set.EvpnGenerateType5HostRoutes != *get.EvpnGenerateType5HostRoutes {
		require.NotNil(t, get.EvpnGenerateType5HostRoutes)
		require.Equalf(t, *set.EvpnGenerateType5HostRoutes, *get.EvpnGenerateType5HostRoutes, "EvpnGenerateType5HostRoutes: set %d get %d", *set.EvpnGenerateType5HostRoutes, *get.EvpnGenerateType5HostRoutes)
	}

	if set.ExternalRouterMtu != nil {
		require.NotNil(t, get.ExternalRouterMtu)
		require.Equalf(t, *set.ExternalRouterMtu, *get.ExternalRouterMtu, "ExternalRouterMtu: set %d get %d", *set.ExternalRouterMtu, *get.ExternalRouterMtu)
	}

	if set.FabricL3Mtu != nil && *set.FabricL3Mtu != *get.FabricL3Mtu {
		require.NotNil(t, get.FabricL3Mtu)
		require.Equalf(t, *set.FabricL3Mtu, *get.FabricL3Mtu, "FabricL3Mtu: set %d get %d", *set.FabricL3Mtu, *get.FabricL3Mtu)
	}

	if set.Ipv6Enabled != nil {
		require.NotNil(t, get.Ipv6Enabled)
		require.Equalf(t, *set.Ipv6Enabled, *get.Ipv6Enabled, "Ipv6Enabled: set %d get %d", *set.Ipv6Enabled, *get.Ipv6Enabled)
	}

	if set.JunosEvpnDuplicateMacRecoveryTime != nil {
		require.NotNil(t, get.JunosEvpnDuplicateMacRecoveryTime)
		require.Equalf(t, *set.JunosEvpnDuplicateMacRecoveryTime, *get.JunosEvpnDuplicateMacRecoveryTime, "JunosEvpnDuplicateMacRecoveryTime: set %d get %d", *set.JunosEvpnDuplicateMacRecoveryTime, *get.JunosEvpnDuplicateMacRecoveryTime)
	}

	if set.JunosEvpnRoutingInstanceVlanAware != nil {
		require.NotNil(t, get.JunosEvpnRoutingInstanceVlanAware)
		require.Equalf(t, *set.JunosEvpnRoutingInstanceVlanAware, *get.JunosEvpnRoutingInstanceVlanAware, "JunosEvpnRoutingInstanceVlanAware: set %d get %d", *set.JunosEvpnRoutingInstanceVlanAware, *get.JunosEvpnRoutingInstanceVlanAware)
	}

	if set.JunosEvpnMaxNexthopAndInterfaceNumber != nil {
		require.NotNil(t, get.JunosEvpnMaxNexthopAndInterfaceNumber)
		require.Equalf(t, *set.JunosEvpnMaxNexthopAndInterfaceNumber, *get.JunosEvpnMaxNexthopAndInterfaceNumber, "JunosEvpnMaxNexthopAndInterfaceNumber: set %d get %d", *set.JunosEvpnMaxNexthopAndInterfaceNumber, *get.JunosEvpnMaxNexthopAndInterfaceNumber)
	}

	if set.JunosExOverlayEcmp != nil {
		require.NotNil(t, get.JunosExOverlayEcmp)
		require.Equalf(t, *set.JunosExOverlayEcmp, *get.JunosExOverlayEcmp, "JunosExOverlayEcmp: set %d get %d", *set.JunosExOverlayEcmp, *get.JunosExOverlayEcmp)
	}

	if set.JunosGracefulRestart != nil {
		require.NotNil(t, get.JunosGracefulRestart)
		require.Equalf(t, *set.JunosGracefulRestart, *get.JunosGracefulRestart, "JunosGracefulRestart: set %d get %d", *set.JunosGracefulRestart, *get.JunosGracefulRestart)
	}

	if set.MaxEvpnRoutes != nil {
		require.NotNil(t, get.MaxEvpnRoutes)
		require.Equalf(t, *set.MaxEvpnRoutes, *get.MaxEvpnRoutes, "MaxEvpnRoutes: set %d get %d", *set.MaxEvpnRoutes, *get.MaxEvpnRoutes)
	}

	if set.MaxExternalRoutes != nil {
		require.NotNil(t, get.MaxExternalRoutes)
		require.Equalf(t, *set.MaxExternalRoutes, *get.MaxExternalRoutes, "MaxExternalRoutes: set %d get %d", *set.MaxExternalRoutes, *get.MaxExternalRoutes)
	}

	if set.MaxFabricRoutes != nil {
		require.NotNil(t, get.MaxFabricRoutes)
		require.Equalf(t, *set.MaxFabricRoutes, *get.MaxFabricRoutes, "MaxFabricRoutes: set %d get %d", *set.MaxFabricRoutes, *get.MaxFabricRoutes)
	}

	if set.MaxMlagRoutes != nil {
		require.NotNil(t, get.MaxMlagRoutes)
		require.Equalf(t, *set.MaxMlagRoutes, *get.MaxMlagRoutes, "MaxMlagRoutes: set %d get %d", *set.MaxMlagRoutes, *get.MaxMlagRoutes)
	}

	if set.OptimiseSzFootprint != nil && *set.OptimiseSzFootprint != *get.OptimiseSzFootprint {
		t.Errorf("set OptimiseSzFootprint %s got %s", *set.OptimiseSzFootprint, *get.OptimiseSzFootprint)
	}

	// don't check overlay control protocol - it's an immutable value. attempts to set it have no effect.
	//if set.OverlayControlProtocol != get.OverlayControlProtocol {
	//	t.Errorf("set OverlayControlProtocol %s got %s", set.OverlayControlProtocol, get.OverlayControlProtocol)
	//}
}

func TestSetGetFabricSettings(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		fabricSettings    FabricSettings
		versionConstraint version.Constraints
	}

	testCases := map[string]testCase{
		"no_fabric_settings": {
			fabricSettings: FabricSettings{},
		},
		"lots_of_values": {
			fabricSettings: FabricSettings{
				JunosEvpnDuplicateMacRecoveryTime:     toPtr(uint16(16)),
				MaxExternalRoutes:                     toPtr(uint32(239832)),
				EsiMacMsb:                             toPtr(uint8(32)),
				JunosGracefulRestart:                  &enum.FeatureSwitchDisabled,
				OptimiseSzFootprint:                   &enum.FeatureSwitchEnabled,
				JunosEvpnRoutingInstanceVlanAware:     &enum.FeatureSwitchEnabled,
				EvpnGenerateType5HostRoutes:           &enum.FeatureSwitchEnabled,
				MaxFabricRoutes:                       toPtr(uint32(84231)),
				MaxMlagRoutes:                         toPtr(uint32(76112)),
				JunosExOverlayEcmp:                    &enum.FeatureSwitchDisabled,
				DefaultSviL3Mtu:                       toPtr(uint16(9100)),
				JunosEvpnMaxNexthopAndInterfaceNumber: &enum.FeatureSwitchDisabled,
				FabricL3Mtu:                           toPtr(uint16(9178)),
				Ipv6Enabled:                           toPtr(false), // do not enable because it's a one-way trip
				ExternalRouterMtu:                     toPtr(uint16(9100)),
				MaxEvpnRoutes:                         toPtr(uint32(92342)),
				AntiAffinityPolicy: &AntiAffinityPolicy{
					Algorithm:                enum.AntiAffinityAlgorithmHeuristic,
					MaxLinksPerPort:          2,
					MaxLinksPerSlot:          2,
					MaxPerSystemLinksPerPort: 2,
					MaxPerSystemLinksPerSlot: 2,
					Mode:                     enum.AntiAffinityModeEnabledLoose,
				},
			},
		},
		"different_values": {
			fabricSettings: FabricSettings{
				JunosEvpnDuplicateMacRecoveryTime:     toPtr(uint16(15)),
				MaxExternalRoutes:                     toPtr(uint32(239732)),
				EsiMacMsb:                             toPtr(uint8(30)),
				JunosGracefulRestart:                  &enum.FeatureSwitchEnabled,
				OptimiseSzFootprint:                   &enum.FeatureSwitchDisabled,
				JunosEvpnRoutingInstanceVlanAware:     &enum.FeatureSwitchDisabled,
				EvpnGenerateType5HostRoutes:           &enum.FeatureSwitchEnabled,
				MaxFabricRoutes:                       toPtr(uint32(84230)),
				MaxMlagRoutes:                         toPtr(uint32(76110)),
				JunosExOverlayEcmp:                    &enum.FeatureSwitchEnabled,
				DefaultSviL3Mtu:                       toPtr(uint16(9050)),
				JunosEvpnMaxNexthopAndInterfaceNumber: &enum.FeatureSwitchEnabled,
				FabricL3Mtu:                           toPtr(uint16(9176)),
				Ipv6Enabled:                           toPtr(false), // do not enable because it's a one-way trip
				ExternalRouterMtu:                     toPtr(uint16(9050)),
				MaxEvpnRoutes:                         toPtr(uint32(92332)),
				AntiAffinityPolicy: &AntiAffinityPolicy{
					Algorithm:                enum.AntiAffinityAlgorithmHeuristic,
					MaxLinksPerPort:          4,
					MaxLinksPerSlot:          4,
					MaxPerSystemLinksPerPort: 4,
					MaxPerSystemLinksPerSlot: 4,
					Mode:                     enum.AntiAffinityModeEnabledStrict,
				},
			},
		},
	}

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(clientName, func(t *testing.T) {
			bpClient := testBlueprintC(ctx, t, client.client)

			for tName, tCase := range testCases {
				tName, tCase := tName, tCase
				t.Run(tName, func(t *testing.T) {
					if tCase.versionConstraint != nil && !tCase.versionConstraint.Check(bpClient.client.apiVersion) {
						t.Skipf("skipping test %q due to version constraints: %q. API version %q",
							tName, tCase.versionConstraint, bpClient.client.apiVersion)
					}

					log.Printf("testing SetFabricSettings() against %s %s (%s)", client.clientType, clientName, bpClient.client.apiVersion)
					err = bpClient.SetFabricSettings(ctx, &tCase.fabricSettings)
					if err != nil {
						t.Fatal(err)
					}

					log.Printf("testing GetFabricSettings() against %s %s (%s)", client.clientType, clientName, bpClient.client.apiVersion)
					fs, err := bpClient.GetFabricSettings(ctx)
					if err != nil {
						t.Fatal(err)
					}
					compareFabricSettings(t, tCase.fabricSettings, *fs)
				})
			}
		})
	}
}

func TestFabricSettingsRoutesMaxDefaultVsZero(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		name              string
		fabricSettings    FabricSettings
		versionConstraint version.Constraints
	}

	// testCases are a slice here to ensure run order
	testCases := []testCase{
		{
			name:           "defaults",
			fabricSettings: FabricSettings{},
		},
		{
			name: "values",
			fabricSettings: FabricSettings{
				MaxEvpnRoutes:     toPtr(uint32(10000)),
				MaxExternalRoutes: toPtr(uint32(11000)),
				MaxFabricRoutes:   toPtr(uint32(12000)),
				MaxMlagRoutes:     toPtr(uint32(13000)),
			},
		},
		{
			name: "zeros",
			fabricSettings: FabricSettings{
				MaxEvpnRoutes:     toPtr(uint32(0)),
				MaxExternalRoutes: toPtr(uint32(0)),
				MaxFabricRoutes:   toPtr(uint32(0)),
				MaxMlagRoutes:     toPtr(uint32(0)),
			},
		},
		{
			name:           "restore_defaults",
			fabricSettings: FabricSettings{},
		},
	}

	for clientName, client := range clients {
		client := client
		bpClient := testBlueprintC(ctx, t, client.client)

		for _, tCase := range testCases {
			tCase := tCase
			t.Run(tCase.name, func(t *testing.T) {
				if tCase.versionConstraint != nil && !tCase.versionConstraint.Check(bpClient.client.apiVersion) {
					t.Skipf("skipping test %q due to version constraints: %q. API version %q",
						tCase.name, tCase.versionConstraint, bpClient.client.apiVersion)
				}

				log.Printf("testing SetFabricSettings() against %s %s (%s)", client.clientType, clientName, bpClient.client.apiVersion)
				err = bpClient.SetFabricSettings(ctx, &tCase.fabricSettings)
				if err != nil {
					t.Fatal(err)
				}

				log.Printf("testing GetFabricSettings() against %s %s (%s)", client.clientType, clientName, bpClient.client.apiVersion)
				fs, err := bpClient.GetFabricSettings(ctx)
				if err != nil {
					t.Fatal(err)
				}
				compareFabricSettings(t, tCase.fabricSettings, *fs)
			})
		}
	}
}

func TestSetGetFabricSettingsV6(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		clientName, client := clientName, client

		t.Run(clientName, func(t *testing.T) {
			bpClient := testBlueprintC(ctx, t, client.client)

			t.Run("enable_and_check_ipv6", func(t *testing.T) {
				fsSet := &FabricSettings{
					AntiAffinityPolicy: &AntiAffinityPolicy{
						Algorithm:                enum.AntiAffinityAlgorithmHeuristic,
						MaxLinksPerPort:          2,
						MaxLinksPerSlot:          2,
						MaxPerSystemLinksPerPort: 2,
						MaxPerSystemLinksPerSlot: 2,
						Mode:                     enum.AntiAffinityModeEnabledStrict,
					},
					EsiMacMsb:                   toPtr(uint8(4)),
					EvpnGenerateType5HostRoutes: &enum.FeatureSwitchEnabled,
					ExternalRouterMtu:           toPtr(uint16(9002)),
					Ipv6Enabled:                 toPtr(true),
					MaxEvpnRoutes:               toPtr(uint32(10000)),
					MaxExternalRoutes:           toPtr(uint32(11000)),
					MaxFabricRoutes:             toPtr(uint32(12000)),
					MaxMlagRoutes:               toPtr(uint32(13000)),
				}
				log.Printf("testing SetFabricSettings() against %s %s (%s)", client.clientType, clientName, bpClient.client.apiVersion)
				err = bpClient.SetFabricSettings(ctx, fsSet)
				if err != nil {
					t.Fatal(err)
				}
				log.Printf("testing GetFabricSettings() against %s %s (%s)", client.clientType, clientName, bpClient.client.apiVersion)
				fsGet, err := bpClient.GetFabricSettings(ctx)
				if err != nil {
					t.Fatal(err)
				}

				compareFabricSettings(t, *fsSet, *fsGet)
			})
		})
	}
}
