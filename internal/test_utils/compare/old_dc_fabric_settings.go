// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package compare

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/stretchr/testify/require"
)

func AntiAffinityPolicy(t testing.TB, set, get apstra.AntiAffinityPolicy) {
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

func FabricSettings(t testing.TB, set, get apstra.FabricSettings) {
	t.Helper()

	if set.AntiAffinityPolicy != nil {
		require.NotNil(t, get.AntiAffinityPolicy)
		AntiAffinityPolicy(t, *get.AntiAffinityPolicy, *set.AntiAffinityPolicy)
	}

	if set.DefaultAnycastGWMAC != nil {
		require.NotNil(t, get.DefaultAnycastGWMAC)
		require.Equalf(t, set.DefaultAnycastGWMAC, get.DefaultAnycastGWMAC, "DefaultAnycastGWMAC: set %s get %s", set.DefaultAnycastGWMAC, get.DefaultAnycastGWMAC)
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
