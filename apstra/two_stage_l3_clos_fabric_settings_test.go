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

	compareAntiAffinityPolicy := func(set, get AntiAffinityPolicy) error {
		if set.Algorithm != get.Algorithm {
			return fmt.Errorf("set AntiAffinityPolicy Algorithm %s got %s", set.Algorithm, get.Algorithm)
		}

		if set.MaxLinksPerPort != get.MaxLinksPerPort {
			return fmt.Errorf("set AntiAffinityPolicy MaxLinksPerPort %d got %d", set.MaxLinksPerPort, get.MaxLinksPerPort)
		}

		if set.MaxLinksPerSlot != get.MaxLinksPerSlot {
			return fmt.Errorf("set AntiAffinityPolicy MaxLinksPerSlot %d got %d", set.MaxLinksPerSlot, get.MaxLinksPerSlot)
		}

		if set.MaxPerSystemLinksPerPort != get.MaxPerSystemLinksPerPort {
			return fmt.Errorf("set AntiAffinityPolicy MaxPerSystemLinksPerPort %d got %d", set.MaxPerSystemLinksPerPort, get.MaxPerSystemLinksPerPort)
		}

		if set.MaxPerSystemLinksPerSlot != get.MaxPerSystemLinksPerSlot {
			return fmt.Errorf("set AntiAffinityPolicy MaxPerSystemLinksPerSlot %d got %d", set.MaxPerSystemLinksPerSlot, get.MaxPerSystemLinksPerSlot)
		}

		if set.Mode != get.Mode {
			return fmt.Errorf("set AntiAffinityPolicy Mode %s got %s", set.Mode, get.Mode)
		}

		return nil
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

		if set.JunosGracefulRestart != get.JunosGracefulRestart {
			return fmt.Errorf("set JunosGracefulRestart %t got %t", set.JunosGracefulRestart, get.JunosGracefulRestart)
		}

		if set.OptimiseSzFootprint != get.OptimiseSzFootprint {
			return fmt.Errorf("set OptimiseSzFootprint %t got %t", set.OptimiseSzFootprint, get.OptimiseSzFootprint)
		}

		if set.JunosEvpnRoutingInstanceVlanAware != get.JunosEvpnRoutingInstanceVlanAware {
			return fmt.Errorf("set JunosEvpnRoutingInstanceVlanAware %t got %t", set.JunosEvpnRoutingInstanceVlanAware, get.JunosEvpnRoutingInstanceVlanAware)
		}

		if set.EvpnGenerateType5HostRoutes != get.EvpnGenerateType5HostRoutes {
			return fmt.Errorf("set EvpnGenerateType5HostRoutes %t got %t", set.EvpnGenerateType5HostRoutes, get.EvpnGenerateType5HostRoutes)
		}

		if set.MaxFabricRoutes != nil && *set.MaxFabricRoutes != *get.MaxFabricRoutes {
			return fmt.Errorf("set MaxFabricRoutes %d got %d", *set.MaxFabricRoutes, *get.MaxFabricRoutes)
		}

		if set.MaxMlagRoutes != nil && *set.MaxMlagRoutes != *get.MaxMlagRoutes {
			return fmt.Errorf("set MaxMlagRoutes %d got %d", *set.MaxMlagRoutes, *get.MaxMlagRoutes)
		}

		if set.JunosExOverlayEcmpDisabled != get.JunosEvpnRoutingInstanceVlanAware {
			return fmt.Errorf("set JunosEvpnRoutingInstanceVlanAware %t got %t", set.JunosEvpnRoutingInstanceVlanAware, get.JunosEvpnRoutingInstanceVlanAware)
		}

		if set.DefaultSviL3Mtu != nil && *set.DefaultSviL3Mtu != *get.DefaultSviL3Mtu {
			return fmt.Errorf("set DefaultSviL3Mtu  %d got %d", *set.DefaultSviL3Mtu, *get.DefaultSviL3Mtu)
		}

		if set.JunosEvpnMaxNexthopAndInterfaceNumberDisabled != get.JunosEvpnMaxNexthopAndInterfaceNumberDisabled {
			return fmt.Errorf("set JunosEvpnMaxNexthopAndInterfaceNumberDisabled %t got %t", set.JunosEvpnMaxNexthopAndInterfaceNumberDisabled, get.JunosEvpnMaxNexthopAndInterfaceNumberDisabled)
		}

		if set.FabricL3Mtu != nil && *set.FabricL3Mtu != *get.FabricL3Mtu {
			return fmt.Errorf("set FabricL3Mtu  %d got %d", *set.FabricL3Mtu, *get.FabricL3Mtu)
		}

		if set.Ipv6Enabled != get.Ipv6Enabled {
			return fmt.Errorf("set Ipv6Enabled %t got %t", set.Ipv6Enabled, get.Ipv6Enabled)
		}

		// don't check overlay control protocol - it's an immutable value. attempts to set it have no effect.
		//if set.OverlayControlProtocol != get.OverlayControlProtocol {
		//	return fmt.Errorf("set OverlayControlProtocol %s got %s", set.OverlayControlProtocol, get.OverlayControlProtocol)
		//}

		if set.ExternalRouterMtu != nil && *set.ExternalRouterMtu != *get.ExternalRouterMtu {
			return fmt.Errorf("set ExternalRouterMtu %d got %d", *set.ExternalRouterMtu, *get.ExternalRouterMtu)
		}

		if set.MaxEvpnRoutes != nil && *set.MaxEvpnRoutes != *get.MaxEvpnRoutes {
			return fmt.Errorf("set MaxEvpnRoutes %d got %d", *set.MaxEvpnRoutes, *get.MaxEvpnRoutes)
		}

		if set.AntiAffinityPolicy != nil {
			err = compareAntiAffinityPolicy(*get.AntiAffinityPolicy, *set.AntiAffinityPolicy)
			if err != nil {
				return err
			}
		}

		return nil
	}

	type testCase struct {
		fabricSettings    FabricSettings
		versionConstraint *version.Constraint
	}

	testCases := map[string]testCase{
		"zerovalues": {
			fabricSettings: FabricSettings{},
		},
		"lotsofvalues": {
			fabricSettings: FabricSettings{
				JunosEvpnDuplicateMacRecoveryTime:             toPtr(uint16(16)),
				MaxExternalRoutes:                             toPtr(uint32(239832)),
				EsiMacMsb:                                     toPtr(uint8(32)),
				JunosGracefulRestart:                          false,
				OptimiseSzFootprint:                           true,
				JunosEvpnRoutingInstanceVlanAware:             true,
				EvpnGenerateType5HostRoutes:                   true,
				MaxFabricRoutes:                               toPtr(uint32(84231)),
				MaxMlagRoutes:                                 toPtr(uint32(76112)),
				JunosExOverlayEcmpDisabled:                    false,
				DefaultSviL3Mtu:                               toPtr(uint16(9100)),
				JunosEvpnMaxNexthopAndInterfaceNumberDisabled: false,
				FabricL3Mtu:                                   toPtr(uint16(9178)),
				Ipv6Enabled:                                   true,
				OverlayControlProtocol:                        OverlayControlProtocolEvpn,
				ExternalRouterMtu:                             toPtr(uint16(9100)),
				MaxEvpnRoutes:                                 toPtr(uint32(92342)),
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
			fs, err := bpClient.GetFabricSettings(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if fs.OverlayControlProtocol != OverlayControlProtocolEvpn {
				t.Fatalf("expected OverlayControlProtocol %q, got %q", OverlayControlProtocolEvpn, fs.OverlayControlProtocol)
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
