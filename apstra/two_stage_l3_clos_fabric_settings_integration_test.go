// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"sync"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestSetGetFabricSettings(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		fabricSettings    []*apstra.FabricSettings
		versionConstraint *compatibility.Constraint
	}

	testCases := map[string]testCase{
		"nil_value": {
			fabricSettings: []*apstra.FabricSettings{nil},
		},
		"zero_value": {
			fabricSettings: []*apstra.FabricSettings{{}},
		},
		"lots_of_values_including_ipv6_enable": {
			versionConstraint: &compatibility.FabricSettingsIPv6EnabledOK,
			fabricSettings: []*apstra.FabricSettings{
				{
					JunosEvpnDuplicateMacRecoveryTime:     pointer.To(uint16(16)),
					MaxExternalRoutes:                     pointer.To(uint32(239832)),
					EsiMacMsb:                             pointer.To(uint8(32)),
					JunosGracefulRestart:                  &enum.FeatureSwitchDisabled,
					OptimiseSzFootprint:                   &enum.FeatureSwitchEnabled,
					JunosEvpnRoutingInstanceVlanAware:     &enum.FeatureSwitchEnabled,
					EvpnGenerateType5HostRoutes:           &enum.FeatureSwitchEnabled,
					MaxFabricRoutes:                       pointer.To(uint32(84231)),
					MaxMlagRoutes:                         pointer.To(uint32(76112)),
					JunosExOverlayEcmp:                    &enum.FeatureSwitchDisabled,
					DefaultSviL3Mtu:                       pointer.To(uint16(9100)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &enum.FeatureSwitchDisabled,
					FabricL3Mtu:                           pointer.To(uint16(9178)),
					Ipv6Enabled:                           pointer.To(false), // do not enable because it's a one-way trip
					ExternalRouterMtu:                     pointer.To(uint16(9100)),
					MaxEvpnRoutes:                         pointer.To(uint32(92342)),
					AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
						Algorithm:                apstra.AlgorithmHeuristic,
						MaxLinksPerPort:          2,
						MaxLinksPerSlot:          2,
						MaxPerSystemLinksPerPort: 2,
						MaxPerSystemLinksPerSlot: 2,
						Mode:                     apstra.AntiAffinityModeEnabledLoose,
					},
				},
				{
					JunosEvpnDuplicateMacRecoveryTime:     pointer.To(uint16(15)),
					MaxExternalRoutes:                     pointer.To(uint32(239732)),
					EsiMacMsb:                             pointer.To(uint8(30)),
					JunosGracefulRestart:                  &enum.FeatureSwitchEnabled,
					OptimiseSzFootprint:                   &enum.FeatureSwitchDisabled,
					JunosEvpnRoutingInstanceVlanAware:     &enum.FeatureSwitchDisabled,
					EvpnGenerateType5HostRoutes:           &enum.FeatureSwitchEnabled,
					MaxFabricRoutes:                       pointer.To(uint32(84230)),
					MaxMlagRoutes:                         pointer.To(uint32(76110)),
					JunosExOverlayEcmp:                    &enum.FeatureSwitchEnabled,
					DefaultSviL3Mtu:                       pointer.To(uint16(9050)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &enum.FeatureSwitchEnabled,
					FabricL3Mtu:                           pointer.To(uint16(9176)),
					Ipv6Enabled:                           pointer.To(false), // do not enable because it's a one-way trip
					ExternalRouterMtu:                     pointer.To(uint16(9050)),
					MaxEvpnRoutes:                         pointer.To(uint32(92332)),
					AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
						Algorithm:                apstra.AlgorithmHeuristic,
						MaxLinksPerPort:          4,
						MaxLinksPerSlot:          4,
						MaxPerSystemLinksPerPort: 4,
						MaxPerSystemLinksPerSlot: 4,
						Mode:                     apstra.AntiAffinityModeEnabledStrict,
					},
				},
			},
		},
		"lots_of_values_including_anycast_gw_mac": {
			versionConstraint: &compatibility.FabricSettingsDefaultAnycastGWMacOK,
			fabricSettings: []*apstra.FabricSettings{
				{
					DefaultAnycastGWMAC:                   testutils.RandomMAC(),
					JunosEvpnDuplicateMacRecoveryTime:     pointer.To(uint16(17)),
					MaxExternalRoutes:                     pointer.To(uint32(239833)),
					EsiMacMsb:                             pointer.To(uint8(34)),
					JunosGracefulRestart:                  &enum.FeatureSwitchEnabled,
					OptimiseSzFootprint:                   &enum.FeatureSwitchDisabled,
					JunosEvpnRoutingInstanceVlanAware:     &enum.FeatureSwitchDisabled,
					EvpnGenerateType5HostRoutes:           &enum.FeatureSwitchDisabled,
					MaxFabricRoutes:                       pointer.To(uint32(84232)),
					MaxMlagRoutes:                         pointer.To(uint32(76113)),
					JunosExOverlayEcmp:                    &enum.FeatureSwitchEnabled,
					DefaultSviL3Mtu:                       pointer.To(uint16(9102)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &enum.FeatureSwitchEnabled,
					FabricL3Mtu:                           pointer.To(uint16(9176)),
					ExternalRouterMtu:                     pointer.To(uint16(9102)),
					MaxEvpnRoutes:                         pointer.To(uint32(92343)),
					AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
						Algorithm:                apstra.AlgorithmHeuristic,
						MaxLinksPerPort:          3,
						MaxLinksPerSlot:          3,
						MaxPerSystemLinksPerPort: 3,
						MaxPerSystemLinksPerSlot: 3,
						Mode:                     apstra.AntiAffinityModeEnabledStrict,
					},
				},
				{
					DefaultAnycastGWMAC:                   testutils.RandomMAC(),
					JunosEvpnDuplicateMacRecoveryTime:     pointer.To(uint16(14)),
					MaxExternalRoutes:                     pointer.To(uint32(239731)),
					EsiMacMsb:                             pointer.To(uint8(28)),
					JunosGracefulRestart:                  &enum.FeatureSwitchDisabled,
					OptimiseSzFootprint:                   &enum.FeatureSwitchEnabled,
					JunosEvpnRoutingInstanceVlanAware:     &enum.FeatureSwitchEnabled,
					EvpnGenerateType5HostRoutes:           &enum.FeatureSwitchDisabled,
					MaxFabricRoutes:                       pointer.To(uint32(84229)),
					MaxMlagRoutes:                         pointer.To(uint32(76109)),
					JunosExOverlayEcmp:                    &enum.FeatureSwitchDisabled,
					DefaultSviL3Mtu:                       pointer.To(uint16(9050)),
					JunosEvpnMaxNexthopAndInterfaceNumber: &enum.FeatureSwitchDisabled,
					FabricL3Mtu:                           pointer.To(uint16(9174)),
					ExternalRouterMtu:                     pointer.To(uint16(9050)),
					MaxEvpnRoutes:                         pointer.To(uint32(92332)),
					AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
						Algorithm:                apstra.AlgorithmHeuristic,
						MaxLinksPerPort:          4,
						MaxLinksPerSlot:          4,
						MaxPerSystemLinksPerPort: 4,
						MaxPerSystemLinksPerSlot: 4,
						Mode:                     apstra.AntiAffinityModeEnabledLoose,
					},
				},
			},
		},
	}

	// create a test blueprint in each test instance
	bpwg := new(sync.WaitGroup)
	bpwg.Add(len(clients))
	bpClients := make([]*apstra.TwoStageL3ClosClient, len(clients))
	for i, client := range clients {
		go func() {
			defer bpwg.Done()
			bpClients[i] = dctestobj.TestBlueprintC(t, ctx, client.Client)
		}()
	}
	bpwg.Wait()

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			// do not use t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			for _, bpClient := range bpClients {
				apiVersion := version.Must(version.NewVersion(bpClient.Client().ApiVersion()))
				t.Run(apiVersion.String(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					if tCase.versionConstraint != nil && !tCase.versionConstraint.Check(apiVersion) {
						t.Skipf("skipping test %q due to version constraints: %q. API version %q",
							tName, tCase.versionConstraint, apiVersion)
					}

					for _, set := range tCase.fabricSettings {
						err := bpClient.SetFabricSettings(ctx, set)
						require.NoError(t, err)

						get, err := bpClient.GetFabricSettings(ctx)
						require.NoError(t, err)
						require.NotNil(t, get)
						if set != nil {
							compare.FabricSettings(t, *set, *get)
						}
					}
				})
			}
		})
	}
}

func TestFabricSettingsRoutesMaxDefaultVsZero(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		name              string
		fabricSettings    apstra.FabricSettings
		versionConstraint version.Constraints
	}

	// testCases are a slice here to ensure run order
	testCases := []testCase{
		{
			name:           "defaults",
			fabricSettings: apstra.FabricSettings{},
		},
		{
			name: "values",
			fabricSettings: apstra.FabricSettings{
				MaxEvpnRoutes:     pointer.To(uint32(10000)),
				MaxExternalRoutes: pointer.To(uint32(11000)),
				MaxFabricRoutes:   pointer.To(uint32(12000)),
				MaxMlagRoutes:     pointer.To(uint32(13000)),
			},
		},
		{
			name: "zeros",
			fabricSettings: apstra.FabricSettings{
				MaxEvpnRoutes:     pointer.To(uint32(0)),
				MaxExternalRoutes: pointer.To(uint32(0)),
				MaxFabricRoutes:   pointer.To(uint32(0)),
				MaxMlagRoutes:     pointer.To(uint32(0)),
			},
		},
		{
			name:           "restore_defaults",
			fabricSettings: apstra.FabricSettings{},
		},
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			ctx := testutils.ContextWithTestID(ctx, t)
			t.Parallel()

			bpClient := dctestobj.TestBlueprintC(t, ctx, client.Client)

			for _, tCase := range testCases {
				t.Run(tCase.name, func(t *testing.T) {
					ctx := testutils.ContextWithTestID(ctx, t)

					apiVersion := version.Must(version.NewVersion(bpClient.Client().ApiVersion()))
					if tCase.versionConstraint != nil && !tCase.versionConstraint.Check(apiVersion) {
						t.Skipf("skipping test %q due to version constraints: %q. API version %q",
							tCase.name, tCase.versionConstraint, apiVersion)
					}

					err := bpClient.SetFabricSettings(ctx, &tCase.fabricSettings)
					if err != nil {
						t.Fatal(err)
					}

					fs, err := bpClient.GetFabricSettings(ctx)
					if err != nil {
						t.Fatal(err)
					}
					compare.FabricSettings(t, tCase.fabricSettings, *fs)
				})
			}
		})
	}
}

func TestSetGetFabricSettingsV6(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			ctx := testutils.ContextWithTestID(ctx, t)
			t.Parallel()

			if !compatibility.FabricSettingsIPv6EnabledOK.Check(client.APIVersion()) {
				t.Skipf("skipping test with Apstra %s", client.APIVersion())
			}

			bpClient := dctestobj.TestBlueprintC(t, ctx, client.Client)

			t.Run("enable_and_check_ipv6", func(t *testing.T) {
				fsSet := &apstra.FabricSettings{
					AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
						Algorithm:                apstra.AlgorithmHeuristic,
						MaxLinksPerPort:          2,
						MaxLinksPerSlot:          2,
						MaxPerSystemLinksPerPort: 2,
						MaxPerSystemLinksPerSlot: 2,
						Mode:                     apstra.AntiAffinityModeEnabledStrict,
					},
					EsiMacMsb:                   pointer.To(uint8(4)),
					EvpnGenerateType5HostRoutes: &enum.FeatureSwitchEnabled,
					ExternalRouterMtu:           pointer.To(uint16(9002)),
					Ipv6Enabled:                 pointer.To(true),
					MaxEvpnRoutes:               pointer.To(uint32(10000)),
					MaxExternalRoutes:           pointer.To(uint32(11000)),
					MaxFabricRoutes:             pointer.To(uint32(12000)),
					MaxMlagRoutes:               pointer.To(uint32(13000)),
				}

				err := bpClient.SetFabricSettings(ctx, fsSet)
				require.NoError(t, err)

				fsGet, err := bpClient.GetFabricSettings(ctx)
				require.NoError(t, err)

				compare.FabricSettings(t, *fsSet, *fsGet)
			})
		})
	}
}
