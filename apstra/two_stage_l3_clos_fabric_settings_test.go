//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-version"
	"log"
	"testing"
)

func TestSetGetFabricSettings(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	compareFabricSettings := func(set, get FabricSettings) error {
		if set.JunosEvpnDuplicateMacRecoveryTime != nil &&
			*set.JunosEvpnDuplicateMacRecoveryTime != *get.JunosEvpnDuplicateMacRecoveryTime {
			return fmt.Errorf("set junosEvpnDuplicateMacRecoveryTime %d got %d", *set.JunosEvpnDuplicateMacRecoveryTime, *get.JunosEvpnDuplicateMacRecoveryTime)
		}

		if set.MaxExternalRoutes != nil && *set.MaxExternalRoutes != *get.MaxExternalRoutes {
			return fmt.Errorf("set MaxExternalRoutes %d got %d", *set.MaxExternalRoutes, *get.MaxExternalRoutes)
		}

		if set.EsiMacMsb != nil && *set.EsiMacMsb != *get.EsiMacMsb {
			return fmt.Errorf("set EsiMacMsb %d got %d", *set.EsiMacMsb, *get.EsiMacMsb)
		}

		if set.MaxFabricRoutes != nil && *set.MaxFabricRoutes != *get.MaxFabricRoutes {
			return fmt.Errorf("set MaxFabricRoutes %d got %d", *set.MaxFabricRoutes, *get.MaxFabricRoutes)
		}

		if set.MaxMlagRoutes != nil && *set.MaxMlagRoutes != *get.MaxMlagRoutes {
			return fmt.Errorf("set MaxMlagRoutes %d got %d", *set.MaxMlagRoutes, *get.MaxMlagRoutes)
		}

		if set.DefaultSviL3Mtu != nil && *set.DefaultSviL3Mtu != *get.DefaultSviL3Mtu {
			return fmt.Errorf("set DefaultSviL3Mtu  %d got %d", *set.DefaultSviL3Mtu, *get.DefaultSviL3Mtu)
		}

		if set.FabricL3Mtu != nil && *set.FabricL3Mtu != *get.FabricL3Mtu {
			return fmt.Errorf("set FabricL3Mtu  %d got %d", *set.FabricL3Mtu, *get.FabricL3Mtu)
		}

		if set.ExternalRouterMtu != nil && *set.ExternalRouterMtu != *get.ExternalRouterMtu {
			return fmt.Errorf("set ExternalRouterMtu %d got %d", *set.ExternalRouterMtu, *get.ExternalRouterMtu)
		}

		if set.MaxEvpnRoutes != nil && *set.MaxEvpnRoutes != *get.MaxEvpnRoutes {
			return fmt.Errorf("set MaxEvpnRoutes %d got %d", *set.MaxEvpnRoutes, *get.MaxEvpnRoutes)
		}

		return nil
	}

	type testCase struct {
		fabricSettings    FabricSettings
		versionConstraint *version.Constraint
	}

	uint32ptr := func(i uint32) *uint32 {
		return &i
	}

	uint16ptr := func(i uint16) *uint16 {
		return &i
	}

	uint8ptr := func(i uint8) *uint8 {
		return &i
	}

	testCases := map[string]testCase{
		"zerovalues": {
			fabricSettings: FabricSettings{},
		},
		"lotsofvalues": {
			fabricSettings: FabricSettings{
				JunosEvpnDuplicateMacRecoveryTime:             uint16ptr(16),
				MaxExternalRoutes:                             uint32ptr(239832),
				EsiMacMsb:                                     uint8ptr(32),
				JunosGracefulRestart:                          false,
				OptimiseSzFootprint:                           true,
				JunosEvpnRoutingInstanceVlanAware:             true,
				EvpnGenerateType5HostRoutes:                   true,
				MaxFabricRoutes:                               uint32ptr(84231),
				MaxMlagRoutes:                                 uint32ptr(76112),
				JunosExOverlayEcmpDisabled:                    false,
				DefaultSviL3Mtu:                               uint16ptr(9100),
				JunosEvpnMaxNexthopAndInterfaceNumberDisabled: false,
				FabricL3Mtu:                                   uint16ptr(9178),
				Ipv6Enabled:                                   true,
				OverlayControlProtocol:                        OverlayControlProtocolEvpn,
				ExternalRouterMtu:                             uint16ptr(9100),
				MaxEvpnRoutes:                                 uint32ptr(92342),
				AntiAffinityPolicy:                            nil,
			},
		},
	}

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintC(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		apiVersion := version.Must(version.NewVersion(client.client.apiVersion))
		t.Run("initial fetch", func(t *testing.T) {
			if !version.MustConstraints(version.NewConstraint(">= 4.2.1")).Check(apiVersion) {
				t.Skipf("skipping test %q due to mismatch version %q", "initial fetch", apiVersion)
			}

			log.Printf("testing GetFabricSettings() against %s %s (%s)", client.clientType, clientName, apiVersion)
			_, err := bpClient.GetFabricSettings(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		for tName, tCase := range testCases {
			tName, tCase := tName, tCase
			t.Run(tName, func(t *testing.T) {
				if tCase.versionConstraint != nil && !tCase.versionConstraint.Check(apiVersion) {
					t.Skipf("skipping test %q due to mismatch version %q", tName, apiVersion)
				}

				log.Printf("testing SetFabricSettings() against %s %s (%s)", client.clientType, clientName, apiVersion)
				err = bpClient.SetFabricSettings(ctx, &tCase.fabricSettings)
				if err != nil {
					t.Fatal(err)
				}

				log.Printf("testing GetFabricSettings() against %s %s (%s)", client.clientType, clientName, apiVersion)
				fs, err := bpClient.GetFabricSettings(ctx)
				if err != nil {
					t.Fatal(err)
				}
				err = compareFabricSettings(tCase.fabricSettings, *fs)
				if err != nil {
					t.Fatal(err)
				}
			})
		}

	}
}
