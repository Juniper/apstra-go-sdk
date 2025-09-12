// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func compareRtPolicy(t testing.TB, a, b *RtPolicy) {
	t.Helper()

	if (a != nil) != (b != nil) { // XOR
		t.Fatalf("RtPolicy exists mismatch: %t vs %t", a != nil, b != nil)
	}

	if a != nil && b != nil {
		compareSlices(t, a.ImportRTs, b.ImportRTs, "RtPolicy ImportRTs elements")
		compareSlices(t, a.ExportRTs, b.ExportRTs, "RtPolicy ExportRTs elements")
	}
}

func comapareSviIps(t testing.TB, a, b SviIp) {
	t.Helper()

	require.Equal(t, a.SystemId, b.SystemId)

	require.Equal(t, a.Ipv4Mode, b.Ipv4Mode)
	if a.Ipv4Addr != nil || b.Ipv4Addr != nil {
		require.NotNil(t, a.Ipv4Addr)
		require.NotNil(t, b.Ipv4Addr)
		require.Equal(t, a.Ipv4Addr.String(), b.Ipv4Addr.String())
	}

	require.Equal(t, a.Ipv6Mode, b.Ipv6Mode)
	if a.Ipv6Addr != nil || b.Ipv6Addr != nil {
		require.NotNil(t, a.Ipv6Addr)
		require.NotNil(t, b.Ipv6Addr)
		require.Equal(t, a.Ipv6Addr.String(), b.Ipv6Addr.String())
	}
}

func compareSviIpSlices(t testing.TB, a, b []SviIp) {
	t.Helper()

	require.Equal(t, len(a), len(b))
	for i := range a {
		log.Printf("comparing SviIps at index %d", i)
		comapareSviIps(t, a[i], b[i])
	}
}

func compareVnBindings(t testing.TB, a, b VnBinding, strict bool) {
	t.Helper()

	if len(a.AccessSwitchNodeIds) != 0 || len(b.AccessSwitchNodeIds) != 0 { // nil and [] slices are equal for our purpose
		compareSlices(t, a.AccessSwitchNodeIds, b.AccessSwitchNodeIds, "VnBindings.AccessSwitchNodeIds")
	}

	require.Equal(t, a.SystemId, b.SystemId)

	if a.VlanId != nil || // the caller specified a VLAN, so we check it
		((a.VlanId != nil || b.VlanId != nil) && strict) { // strict mode means we always check
		require.NotNil(t, a.VlanId)
		require.NotNil(t, b.VlanId)
		require.Equal(t, a.VlanId, b.VlanId)
	}
}

func compareVnBindingSlices(t testing.TB, a, b []VnBinding, strict bool) {
	t.Helper()

	require.Equal(t, len(a), len(b))
	for i := range a {
		log.Printf("comparing VnBindings at index %d", i)
		compareVnBindings(t, a[i], b[i], strict)
	}
}

func compareVirtualNetworkData(t testing.TB, a, b *VirtualNetworkData, strict bool) {
	t.Helper()

	require.Equal(t, a.DhcpService, b.DhcpService)
	require.Equal(t, a.Ipv4Enabled, b.Ipv4Enabled)
	require.Equal(t, a.Ipv4Subnet, b.Ipv4Subnet)
	require.Equal(t, a.Ipv6Enabled, b.Ipv6Enabled)
	require.Equal(t, a.Ipv6Subnet, b.Ipv6Subnet)
	require.Equal(t, a.Label, b.Label)
	require.Equal(t, a.RouteTarget, b.RouteTarget)
	require.Equal(t, a.SecurityZoneId, b.SecurityZoneId)
	require.Equal(t, a.VirtualGatewayIpv4, b.VirtualGatewayIpv4)
	require.Equal(t, a.VirtualGatewayIpv6, b.VirtualGatewayIpv6)
	require.Equal(t, a.VirtualGatewayIpv4Enabled, b.VirtualGatewayIpv4Enabled)
	require.Equal(t, a.VirtualGatewayIpv6Enabled, b.VirtualGatewayIpv6Enabled)
	require.Equal(t, a.VnType, b.VnType)
	require.Equal(t, a.VirtualMac, b.VirtualMac)

	if a.L3Mtu != nil || b.L3Mtu != nil {
		require.NotNil(t, a.L3Mtu)
		require.NotNil(t, b.L3Mtu)
		require.Equal(t, a.L3Mtu, b.L3Mtu)
	}

	if a.ReservedVlanId != nil || b.ReservedVlanId != nil {
		require.NotNil(t, a.ReservedVlanId)
		require.NotNil(t, b.ReservedVlanId)
		require.Equal(t, a.ReservedVlanId, b.ReservedVlanId)
	}

	if a.VnId != nil || b.VnId != nil {
		require.NotNil(t, a.VnId)
		require.NotNil(t, b.VnId)
		require.Equal(t, a.VnId, b.VnId)
	}

	compareRtPolicy(t, a.RtPolicy, b.RtPolicy)
	compareSviIpSlices(t, a.SviIps, b.SviIps)
	compareVnBindingSlices(t, a.VnBindings, b.VnBindings, strict)
}

func TestCreateUpdateDeleteVirtualNetwork(t *testing.T) {
	var ace ClientErr
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	randStr := randString(5, "hex")
	label := "test-" + randStr
	vrfName := "test-" + randStr

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			bpClient := testBlueprintC(ctx, t, client.client)

			bpClient.SetType(BlueprintTypeStaging)

			log.Printf("testing CreateSecurityZone() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			zoneId, err := bpClient.CreateSecurityZone(ctx, &SecurityZoneData{
				SzType:  SecurityZoneTypeEVPN,
				VrfName: vrfName,
				Label:   label,
			})
			if err != nil {
				t.Fatal(err)
			}

			var result struct {
				Items []struct {
					System struct {
						SystemId string `json:"id"`
					} `json:"system"`
				} `json:"items"`
			}

			query := new(PathQuery).
				SetClient(client.client).
				SetBlueprintId(bpClient.Id()).
				Node([]QEEAttribute{
					{"type", QEStringVal("system")},
					{"system_type", QEStringVal("switch")},
					{"role", QEStringVal("leaf")},
					{"name", QEStringVal("system")},
				})

			err = query.Do(ctx, &result)
			if err != nil {
				t.Fatal(err)
			}

			sviIps := make([]SviIp, len(result.Items))
			vnBindings := make([]VnBinding, len(result.Items))
			for i := range result.Items {
				leafId := ObjectId(result.Items[i].System.SystemId)
				sviIps[i] = SviIp{
					SystemId: leafId,
					Ipv4Mode: enum.SviIpv4ModeEnabled,
					Ipv6Mode: enum.SviIpv6ModeDisabled,
				}
				vnBindings[i] = VnBinding{
					SystemId: leafId,
				}
			}

			l3Mtu := toPtr(1280 + (2 * rand.Intn(3969))) // 1280 - 9216 even numbers only

			createData := VirtualNetworkData{
				Ipv4Enabled:               true,
				L3Mtu:                     l3Mtu,
				Label:                     label,
				SecurityZoneId:            zoneId,
				SviIps:                    sviIps[:1],
				VirtualGatewayIpv4Enabled: true,
				VnBindings:                vnBindings[:1],
				VnType:                    enum.VnTypeVxlan,
			}

			log.Printf("testing CreateVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			vnId, err := bpClient.CreateVirtualNetwork(ctx, &createData)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("created virtual network - id:'%s', name: '%s', label:'%s'", vnId, vrfName, label)

			log.Printf("testing CreateVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			shouldFail, err := bpClient.CreateVirtualNetwork(ctx, &createData)
			if err == nil {
				t.Fatalf("Creating two virtual networks with name %q should have failed, but %q and %q seem to coexist",
					label, vnId, shouldFail)
			}
			if !errors.As(err, &ace) || ace.Type() != ErrExists {
				t.Fatalf("creating two VNs with same name should fail, but not for this reason: %q", err.Error())
			}

			log.Printf("testing GetVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			getById, err := bpClient.GetVirtualNetwork(ctx, vnId)
			if err != nil {
				t.Fatal(err)
			}
			if vnId != getById.Id {
				t.Fatalf("Virtual Network ID mismatch: %q vs. %q", vnId, getById.Id)
			}
			compareVirtualNetworkData(t, &createData, getById.Data, false)

			getByName, err := bpClient.GetVirtualNetworkByName(ctx, getById.Data.Label)
			if err != nil {
				t.Fatal(err)
			}
			if vnId != getByName.Id {
				t.Fatalf("Virtual Network ID mismatch: %q vs. %q", vnId, getByName.Id)
			}
			compareVirtualNetworkData(t, &createData, getByName.Data, false)

			newVlan := Vlan(100)
			createData.ReservedVlanId = &newVlan
			createData.Label = randString(10, "hex")
			createData.L3Mtu = toPtr(1280 + (2 * rand.Intn(3969))) // 1280 - 9216 even numbers only

			for i := range createData.VnBindings {
				createData.VnBindings[i].VlanId = &newVlan
			}

			log.Printf("testing UpdateVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.UpdateVirtualNetwork(ctx, vnId, &createData)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			getById, err = bpClient.GetVirtualNetwork(ctx, vnId)
			if err != nil {
				t.Fatal(err)
			}
			if vnId != getById.Id {
				t.Fatalf("Virtual Network ID mismatch: %q vs. %q", vnId, getById.Id)
			}
			compareVirtualNetworkData(t, &createData, getById.Data, true)

			vnMap, err := bpClient.GetAllVirtualNetworks(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if len(vnMap) != 1 {
				t.Fatalf("expected one VN got %d", len(vnMap))
			}
			if _, ok := vnMap[vnId]; !ok {
				t.Fatalf("map does not contain virtual network %q", vnId)
			}
			batchData := createData
			batchData.SviIps = nil // the "get all" API call omits SVI info. for. some. reason.
			compareVirtualNetworkData(t, &batchData, vnMap[vnId].Data, true)

			log.Printf("testing DeleteVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.DeleteVirtualNetwork(ctx, vnId)
			if err != nil {
				t.Fatal(err)
			}

			// get the deleted VN, expect 404
			log.Printf("testing GetVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			_, err = bpClient.GetVirtualNetwork(ctx, vnId)
			if err == nil {
				t.Fatal("GetVirtualNetwork after DeleteVirtualNetwork should have produced an error")
			}
			if !errors.As(err, &ace) || ace.Type() != ErrNotfound {
				t.Fatalf("expected a 404/NotFound error after deletion")
			}

			// delete the deleted VN, expect 404
			log.Printf("testing DeleteVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.DeleteVirtualNetwork(ctx, vnId)
			if err == nil {
				t.Fatal("DeleteVirtualNetwork after DeleteVirtualNetwork should have produced an error")
			}
			if !errors.As(err, &ace) || ace.Type() != ErrNotfound {
				t.Fatalf("expected a 404/NotFound error after deletion")
			}
		})
	}
}

func TestVirtualNetworkTags(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	randStr := randString(5, "hex")
	label := "test-" + randStr
	vrfName := "test-" + randStr

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			if !compatibility.VirtualNetworkTags.Check(client.client.apiVersion) {
				t.Skipf("skipping virtual network tag test with version %s", client.client.apiVersion)
			}

			t.Parallel()

			bpClient := testBlueprintC(ctx, t, client.client)

			log.Printf("testing CreateSecurityZone() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			zoneId, err := bpClient.CreateSecurityZone(ctx, &SecurityZoneData{
				SzType:  SecurityZoneTypeEVPN,
				VrfName: vrfName,
				Label:   label,
			})
			require.NoError(t, err)

			var result struct {
				Items []struct {
					System struct {
						SystemId string `json:"id"`
					} `json:"system"`
				} `json:"items"`
			}

			query := new(PathQuery).
				SetClient(client.client).
				SetBlueprintId(bpClient.Id()).
				Node([]QEEAttribute{
					{"type", QEStringVal("system")},
					{"system_type", QEStringVal("switch")},
					{"role", QEStringVal("leaf")},
					{"name", QEStringVal("system")},
				})

			err = query.Do(ctx, &result)
			require.NoError(t, err)

			vnBindings := make([]VnBinding, len(result.Items))
			for i := range result.Items {
				leafId := ObjectId(result.Items[i].System.SystemId)
				vnBindings[i] = VnBinding{
					SystemId: leafId,
				}
			}

			createData := VirtualNetworkData{
				Ipv4Enabled:               true,
				Label:                     label,
				SecurityZoneId:            zoneId,
				VirtualGatewayIpv4Enabled: true,
				VnBindings:                vnBindings[:1],
				VnType:                    enum.VnTypeVxlan,
			}

			log.Printf("testing CreateVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			vnId, err := bpClient.CreateVirtualNetwork(ctx, &createData)
			require.NoError(t, err)

			log.Printf("testing GetVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			vn, err := bpClient.GetVirtualNetwork(ctx, vnId)
			require.NoError(t, err)
			require.Equal(t, 0, len(vn.Data.Tags))

			tags := make([]string, rand.Intn(10)+1)
			for i := range tags {
				tags[i] = randString(6, "hex")
			}

			log.Printf("setting tags on virtual network against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetNodeTags(ctx, vn.Id, tags)
			require.NoError(t, err)

			log.Printf("testing GetVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			vn, err = bpClient.GetVirtualNetwork(ctx, vnId)
			require.NoError(t, err)
			require.Equal(t, len(tags), len(vn.Data.Tags))
			compareSlicesAsSets(t, tags, vn.Data.Tags, "tag set mismatch")

			log.Printf("clearing tags on virtual network against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetNodeTags(ctx, vn.Id, nil)
			require.NoError(t, err)

			log.Printf("testing GetVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			vn, err = bpClient.GetVirtualNetwork(ctx, vnId)
			require.NoError(t, err)
			require.Equal(t, 0, len(vn.Data.Tags))

			createData.Tags = tags
			var ace ClientErr

			log.Printf("ensuring that CreateVirtualNetwork() with tags fails against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			_, err = bpClient.CreateVirtualNetwork(ctx, &createData)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ErrNotSupported, ace.Type())

			log.Printf("ensuring that UpdateVirtualNetwork() with tags fails against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.UpdateVirtualNetwork(ctx, vnId, &createData)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ErrNotSupported, ace.Type())
		})
	}
}
