//go:build integration
// +build integration

package apstra

import (
	"context"
	"github.com/hashicorp/go-version"
	"log"
	"testing"
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

	if set.JunosEvpnDuplicateMacRecoveryTime != nil &&
		*set.JunosEvpnDuplicateMacRecoveryTime != *get.JunosEvpnDuplicateMacRecoveryTime {
		t.Errorf("set junosEvpnDuplicateMacRecoveryTime %d got %d", *set.JunosEvpnDuplicateMacRecoveryTime, *get.JunosEvpnDuplicateMacRecoveryTime)
	}

	if set.MaxExternalRoutes != nil && *set.MaxExternalRoutes != *get.MaxExternalRoutes {
		t.Errorf("set MaxExternalRoutes %d got %d", *set.MaxExternalRoutes, *get.MaxExternalRoutes)
	}

	if set.EsiMacMsb != nil && *set.EsiMacMsb != *get.EsiMacMsb {
		t.Errorf("set EsiMacMsb %d got %d", *set.EsiMacMsb, *get.EsiMacMsb)
	}

	if set.JunosGracefulRestart != nil && *set.JunosGracefulRestart != *get.JunosGracefulRestart {
		t.Errorf("set JunosGracefulRestart %s got %s", *set.JunosGracefulRestart, *get.JunosGracefulRestart)
	}

	if set.OptimiseSzFootprint != nil && *set.OptimiseSzFootprint != *get.OptimiseSzFootprint {
		t.Errorf("set OptimiseSzFootprint %s got %s", *set.OptimiseSzFootprint, *get.OptimiseSzFootprint)
	}

	if set.JunosEvpnRoutingInstanceVlanAware != nil && *set.JunosEvpnRoutingInstanceVlanAware != *get.JunosEvpnRoutingInstanceVlanAware {
		t.Errorf("set JunosEvpnRoutingInstanceVlanAware %s got %s", *set.JunosEvpnRoutingInstanceVlanAware, *get.JunosEvpnRoutingInstanceVlanAware)
	}

	if set.EvpnGenerateType5HostRoutes != nil && *set.EvpnGenerateType5HostRoutes != *get.EvpnGenerateType5HostRoutes {
		t.Errorf("set EvpnGenerateType5HostRoutes %s got %s", *set.EvpnGenerateType5HostRoutes, *get.EvpnGenerateType5HostRoutes)
	}

	if set.MaxFabricRoutes != nil && *set.MaxFabricRoutes != *get.MaxFabricRoutes {
		t.Errorf("set MaxFabricRoutes %d got %d", *set.MaxFabricRoutes, *get.MaxFabricRoutes)
	}

	if set.MaxMlagRoutes != nil && *set.MaxMlagRoutes != *get.MaxMlagRoutes {
		t.Errorf("set MaxMlagRoutes %d got %d", *set.MaxMlagRoutes, *get.MaxMlagRoutes)
	}

	if set.JunosEvpnRoutingInstanceVlanAware != nil && *set.JunosEvpnRoutingInstanceVlanAware != *get.JunosEvpnRoutingInstanceVlanAware {
		t.Errorf("set JunosEvpnRoutingInstanceVlanAware %s got %s", *set.JunosEvpnRoutingInstanceVlanAware, *get.JunosEvpnRoutingInstanceVlanAware)
	}

	if set.DefaultSviL3Mtu != nil && *set.DefaultSviL3Mtu != *get.DefaultSviL3Mtu {
		t.Errorf("set DefaultSviL3Mtu  %d got %d", *set.DefaultSviL3Mtu, *get.DefaultSviL3Mtu)
	}

	if set.JunosEvpnMaxNexthopAndInterfaceNumber != nil && *set.JunosEvpnMaxNexthopAndInterfaceNumber != *get.JunosEvpnMaxNexthopAndInterfaceNumber {
		t.Errorf("set JunosEvpnMaxNexthopAndInterfaceNumber %s got %s", *set.JunosEvpnMaxNexthopAndInterfaceNumber, *get.JunosEvpnMaxNexthopAndInterfaceNumber)
	}

	if set.FabricL3Mtu != nil && *set.FabricL3Mtu != *get.FabricL3Mtu {
		t.Errorf("set FabricL3Mtu  %d got %d", *set.FabricL3Mtu, *get.FabricL3Mtu)
	}

	if set.Ipv6Enabled != nil && *set.Ipv6Enabled != *get.Ipv6Enabled {
		t.Errorf("set Ipv6Enabled %t got %t", *set.Ipv6Enabled, *get.Ipv6Enabled)
	}

	// don't check overlay control protocol - it's an immutable value. attempts to set it have no effect.
	//if set.OverlayControlProtocol != get.OverlayControlProtocol {
	//	t.Errorf("set OverlayControlProtocol %s got %s", set.OverlayControlProtocol, get.OverlayControlProtocol)
	//}

	if set.ExternalRouterMtu != nil && *set.ExternalRouterMtu != *get.ExternalRouterMtu {
		t.Errorf("set ExternalRouterMtu %d got %d", *set.ExternalRouterMtu, *get.ExternalRouterMtu)
	}

	if set.MaxEvpnRoutes != nil && *set.MaxEvpnRoutes != *get.MaxEvpnRoutes {
		t.Errorf("set MaxEvpnRoutes %d got %d", *set.MaxEvpnRoutes, *get.MaxEvpnRoutes)
	}

	if set.AntiAffinityPolicy != nil {
		compareAntiAffinityPolicy(t, *get.AntiAffinityPolicy, *set.AntiAffinityPolicy)
	}
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
		"41x_compatible_1": {
			fabricSettings: FabricSettings{
				AntiAffinityPolicy: &AntiAffinityPolicy{
					Algorithm:                AlgorithmHeuristic,
					MaxLinksPerPort:          2,
					MaxLinksPerSlot:          2,
					MaxPerSystemLinksPerPort: 2,
					MaxPerSystemLinksPerSlot: 2,
					Mode:                     AntiAffinityModeEnabledStrict,
				},
				EsiMacMsb:                   toPtr(uint8(4)),
				EvpnGenerateType5HostRoutes: &FeatureSwitchEnumEnabled,
				ExternalRouterMtu:           toPtr(uint16(9002)),
				Ipv6Enabled:                 toPtr(false), // do not enable because it's a one-way trip
				MaxEvpnRoutes:               toPtr(uint32(10000)),
				MaxExternalRoutes:           toPtr(uint32(11000)),
				MaxFabricRoutes:             toPtr(uint32(12000)),
				MaxMlagRoutes:               toPtr(uint32(13000)),
			},
		},
		"412_compatible_2": {
			fabricSettings: FabricSettings{
				AntiAffinityPolicy: &AntiAffinityPolicy{
					Algorithm:                AlgorithmHeuristic,
					MaxLinksPerPort:          3,
					MaxLinksPerSlot:          3,
					MaxPerSystemLinksPerPort: 3,
					MaxPerSystemLinksPerSlot: 3,
					Mode:                     AntiAffinityModeEnabledLoose,
				},
				EsiMacMsb:                   toPtr(uint8(6)),
				EvpnGenerateType5HostRoutes: &FeatureSwitchEnumDisabled,
				ExternalRouterMtu:           toPtr(uint16(9004)),
				Ipv6Enabled:                 toPtr(false), // do not enable because it's a one-way trip
				MaxEvpnRoutes:               toPtr(uint32(20000)),
				MaxExternalRoutes:           toPtr(uint32(21000)),
				MaxFabricRoutes:             toPtr(uint32(22000)),
				MaxMlagRoutes:               toPtr(uint32(23000)),
			},
		},
		"lots_of_values": {
			versionConstraint: version.MustConstraints(version.NewConstraint(">" + apstra412)),
			fabricSettings: FabricSettings{
				JunosEvpnDuplicateMacRecoveryTime:     toPtr(uint16(16)),
				MaxExternalRoutes:                     toPtr(uint32(239832)),
				EsiMacMsb:                             toPtr(uint8(32)),
				JunosGracefulRestart:                  &FeatureSwitchEnumDisabled,
				OptimiseSzFootprint:                   &FeatureSwitchEnumEnabled,
				JunosEvpnRoutingInstanceVlanAware:     &FeatureSwitchEnumEnabled,
				EvpnGenerateType5HostRoutes:           &FeatureSwitchEnumEnabled,
				MaxFabricRoutes:                       toPtr(uint32(84231)),
				MaxMlagRoutes:                         toPtr(uint32(76112)),
				JunosExOverlayEcmp:                    &FeatureSwitchEnumDisabled,
				DefaultSviL3Mtu:                       toPtr(uint16(9100)),
				JunosEvpnMaxNexthopAndInterfaceNumber: &FeatureSwitchEnumDisabled,
				FabricL3Mtu:                           toPtr(uint16(9178)),
				Ipv6Enabled:                           toPtr(false), // do not enable because it's a one-way trip
				ExternalRouterMtu:                     toPtr(uint16(9100)),
				MaxEvpnRoutes:                         toPtr(uint32(92342)),
				AntiAffinityPolicy: &AntiAffinityPolicy{
					Algorithm:                AlgorithmHeuristic,
					MaxLinksPerPort:          2,
					MaxLinksPerSlot:          2,
					MaxPerSystemLinksPerPort: 2,
					MaxPerSystemLinksPerSlot: 2,
					Mode:                     AntiAffinityModeEnabledLoose,
				},
			},
		},
		"different_values": {
			versionConstraint: version.MustConstraints(version.NewConstraint(">" + apstra412)),
			fabricSettings: FabricSettings{
				JunosEvpnDuplicateMacRecoveryTime:     toPtr(uint16(15)),
				MaxExternalRoutes:                     toPtr(uint32(239732)),
				EsiMacMsb:                             toPtr(uint8(30)),
				JunosGracefulRestart:                  &FeatureSwitchEnumEnabled,
				OptimiseSzFootprint:                   &FeatureSwitchEnumDisabled,
				JunosEvpnRoutingInstanceVlanAware:     &FeatureSwitchEnumDisabled,
				EvpnGenerateType5HostRoutes:           &FeatureSwitchEnumEnabled,
				MaxFabricRoutes:                       toPtr(uint32(84230)),
				MaxMlagRoutes:                         toPtr(uint32(76110)),
				JunosExOverlayEcmp:                    &FeatureSwitchEnumEnabled,
				DefaultSviL3Mtu:                       toPtr(uint16(9050)),
				JunosEvpnMaxNexthopAndInterfaceNumber: &FeatureSwitchEnumEnabled,
				FabricL3Mtu:                           toPtr(uint16(9176)),
				Ipv6Enabled:                           toPtr(false), // do not enable because it's a one-way trip
				ExternalRouterMtu:                     toPtr(uint16(9050)),
				MaxEvpnRoutes:                         toPtr(uint32(92332)),
				AntiAffinityPolicy: &AntiAffinityPolicy{
					Algorithm:                AlgorithmHeuristic,
					MaxLinksPerPort:          4,
					MaxLinksPerSlot:          4,
					MaxPerSystemLinksPerPort: 4,
					MaxPerSystemLinksPerSlot: 4,
					Mode:                     AntiAffinityModeEnabledStrict,
				},
			},
		},
	}

	for clientName, client := range clients {
		client := client
		bpClient, bpDel := testBlueprintC(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

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
	}
}

func TestSetGetFabricSettingsV6(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		client := client

		bpClient, bpDel := testBlueprintC(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		t.Run("enable_and_check_ipv6", func(t *testing.T) {
			fsSet := &FabricSettings{
				AntiAffinityPolicy: &AntiAffinityPolicy{
					Algorithm:                AlgorithmHeuristic,
					MaxLinksPerPort:          2,
					MaxLinksPerSlot:          2,
					MaxPerSystemLinksPerPort: 2,
					MaxPerSystemLinksPerSlot: 2,
					Mode:                     AntiAffinityModeEnabledStrict,
				},
				EsiMacMsb:                   toPtr(uint8(4)),
				EvpnGenerateType5HostRoutes: &FeatureSwitchEnumEnabled,
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
	}
}
