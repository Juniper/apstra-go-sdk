//go:build integration
// +build integration

package apstra

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"testing"
)

func compareRtPolicy(t *testing.T, a, b *RtPolicy) {
	if (a != nil) != (b != nil) { // XOR
		t.Fatalf("RtPolicy exists mismatch: %t vs %t", a != nil, b != nil)
	}

	if a != nil && b != nil {
		compareSlices(t, a.ImportRTs, b.ImportRTs, "RtPolicy ImportRTs elements")
		compareSlices(t, a.ExportRTs, b.ExportRTs, "RtPolicy ExportRTs elements")
	}
}

func comapareSviIps(t *testing.T, a, b SviIp) {
	if a.SystemId != b.SystemId {
		t.Fatalf("SystemId mismatch: %q vs. %q", a.SystemId, b.SystemId)
	}

	if !a.Ipv4Addr.Equal(b.Ipv4Addr) {
		t.Fatalf("Ipv4Addr mismatch: %q vs. %q", a.Ipv4Addr.String(), b.Ipv4Addr.String())
	}

	if a.Ipv4Mode != b.Ipv4Mode {
		t.Fatalf("Ipv4Mode mismatch: %q vs. %q", a.Ipv4Mode.String(), b.Ipv4Mode.String())
	}

	if a.Ipv4Requirement != b.Ipv4Requirement {
		t.Fatalf("Ipv4Requirement mismatch: %q vs. %q", a.Ipv4Requirement.String(), b.Ipv4Requirement.String())
	}

	if !a.Ipv6Addr.Equal(b.Ipv6Addr) {
		t.Fatalf("Ipv6Addr mismatch: %q vs. %q", a.Ipv6Addr.String(), b.Ipv6Addr.String())
	}

	if a.Ipv6Mode != b.Ipv6Mode {
		t.Fatalf("Ipv6Mode mismatch: %q vs. %q", a.Ipv6Mode.String(), b.Ipv6Mode.String())
	}

	if a.Ipv6Requirement != b.Ipv6Requirement {
		t.Fatalf("Ipv6Requirement mismatch: %q vs. %q", a.Ipv6Requirement.String(), b.Ipv6Requirement.String())
	}
}

func compareSviIpSlices(t *testing.T, a, b []SviIp) {
	if len(a) != len(b) {
		t.Fatalf("SviIps length mismatch: %d vs %d", len(a), len(b))
	}
	for i := range a {
		log.Printf("comparing SviIps at index %d", i)
		comapareSviIps(t, a[i], b[i])
	}
}

func compareVnBindings(t *testing.T, a, b VnBinding, strict bool) {
	compareSlices(t, a.AccessSwitchNodeIds, b.AccessSwitchNodeIds, "VnBindings.AccessSwitchNodeIds")

	if a.SystemId != b.SystemId {
		t.Fatalf("SystemId mismatch: %q vs. %q", a.SystemId, b.SystemId)
	}

	if strict && (a.VlanId != nil) != (b.VlanId != nil) {
		t.Fatalf("VlanId exists mismatch: %t vs. %t", a.VlanId != nil, b.VlanId != nil)
	}

	if a.VlanId != nil && b.VlanId != nil && *a.VlanId != *b.VlanId {
		t.Fatalf("VlanId mismatch: %d vs. %d", *a.VlanId, *b.VlanId)
	}
}

func compareVnBindingSlices(t *testing.T, a, b []VnBinding, strict bool) {
	if len(a) != len(b) {
		t.Fatalf("VnBindings length mismatch: %d vs %d", len(a), len(b))
	}
	for i := range a {
		log.Printf("comparing VnBindings at index %d", i)
		compareVnBindings(t, a[i], b[i], strict)
	}
}

func compareVirtualNetworkData(t *testing.T, a, b *VirtualNetworkData, strict bool) {
	if a.DhcpService != b.DhcpService {
		t.Fatalf("DhcpService mismatch: %q vs. %q", a.DhcpService.raw(), b.DhcpService.raw())
	}

	if a.Ipv4Enabled != b.Ipv4Enabled {
		t.Fatalf("Ipv4Enabled mismatch: %t vs. %t", a.Ipv4Enabled, b.Ipv4Enabled)
	}

	if a.Ipv4Subnet.String() != b.Ipv4Subnet.String() {
		t.Fatalf("Ipv4Subnet mismatch: %q vs. %q", a.Ipv4Subnet.String(), b.Ipv4Subnet.String())
	}

	if a.Ipv6Enabled != b.Ipv6Enabled {
		t.Fatalf("Ipv6Enabled mismatch: %t vs. %t", a.Ipv6Enabled, b.Ipv6Enabled)
	}

	if a.Ipv6Subnet.String() != b.Ipv6Subnet.String() {
		t.Fatalf("Ipv6Subnet mismatch: %q vs. %q", a.Ipv6Subnet.String(), b.Ipv6Subnet.String())
	}

	aL3Mtu := a.L3Mtu != nil
	bL3Mtu := b.L3Mtu != nil
	if (aL3Mtu || bL3Mtu) && !(aL3Mtu && bL3Mtu) { // xor
		t.Fatalf("L3 MTU setting mismatch: set %t vs. set %t", aL3Mtu, bL3Mtu)
	}

	if aL3Mtu && bL3Mtu && (*a.L3Mtu != *b.L3Mtu) {
		t.Fatalf("L3 MTU setting mismatch: %d vs. %d", *a.L3Mtu, *b.L3Mtu)
	}

	if a.Label != b.Label {
		t.Fatalf("Label mismatch: %q vs. %q", a.Label, b.Label)
	}

	if (a.ReservedVlanId != nil) != (b.ReservedVlanId != nil) { // XOR
		t.Fatalf("ReservedVlanId exists mismatch: %t vs %t", a.ReservedVlanId != nil, b.ReservedVlanId != nil)
	}

	if a.ReservedVlanId != nil && b.ReservedVlanId != nil && *a.ReservedVlanId != *b.ReservedVlanId {
		t.Fatalf("ReservedVlanId mismatch: %d vs %d", *a.ReservedVlanId, *b.ReservedVlanId)
	}

	if a.RouteTarget != b.RouteTarget {
		t.Fatalf("RouteTarget mismatch: %q vs. %q", a.RouteTarget, b.RouteTarget)
	}

	compareRtPolicy(t, a.RtPolicy, b.RtPolicy)

	if a.SecurityZoneId != b.SecurityZoneId {
		t.Fatalf("SecurityZoneId mismatch: %q vs %q", a.SecurityZoneId, b.SecurityZoneId)
	}

	compareSviIpSlices(t, a.SviIps, b.SviIps)

	if !a.VirtualGatewayIpv4.Equal(b.VirtualGatewayIpv4) {
		t.Fatalf("VirtualGatwayIpv4 mismatch: %q vs. %q", a.VirtualGatewayIpv4.String(), b.VirtualGatewayIpv4.String())
	}

	if !a.VirtualGatewayIpv6.Equal(b.VirtualGatewayIpv6) {
		t.Fatalf("VirtualGatwayIpv6 mismatch: %q vs. %q", a.VirtualGatewayIpv6.String(), b.VirtualGatewayIpv6.String())
	}

	if a.VirtualGatewayIpv4Enabled != b.VirtualGatewayIpv4Enabled {
		t.Fatalf("VirtualGatewayIpv4Enabled mismatch: %t vs %t", a.VirtualGatewayIpv4Enabled, b.VirtualGatewayIpv4Enabled)
	}

	if a.VirtualGatewayIpv6Enabled != b.VirtualGatewayIpv6Enabled {
		t.Fatalf("VirtualGatewayIpv6Enabled mismatch: %t vs %t", a.VirtualGatewayIpv6Enabled, b.VirtualGatewayIpv6Enabled)
	}

	compareVnBindingSlices(t, a.VnBindings, b.VnBindings, strict)

	if (a.VnId != nil) != (b.VnId != nil) {
		t.Fatalf("VnId exists mismatch: %t vs. %t", a.VnId != nil, b.VnId != nil)
	}

	if a.VnId != nil && b.VnId != nil && *a.VnId != *b.VnId {
		t.Fatalf("VnId mismatch: %d vs. %d", *a.VnId, *b.VnId)
	}

	if a.VnType != b.VnType {
		t.Fatalf("VnType mismatch: %q vs. %q", a.VnType.String(), b.VnType.String())
	}

	if a.VirtualMac.String() != b.VirtualMac.String() {
		t.Fatalf("VirtualMac mismatch: %q vs. %q", a.VirtualMac.String(), b.VirtualMac.String())
	}
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
				Ipv4Mode: Ipv4ModeEnabled,
				Ipv6Mode: Ipv6ModeDisabled,
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
			VnType:                    VnTypeVxlan,
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

		//log.Printf("testing SetResourceAllocation() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		//err = bpClient.SetResourceAllocation(ctx, rga)
		//if err != nil {
		//	t.Fatal()
		//}
		//
		//log.Printf("testing GetResourceAllocation() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		//rga, err = bpClient.GetResourceAllocation(ctx, &rga.ResourceGroup)
		//if err != nil {
		//	t.Fatal()
		//}
		//
		//if len(rga.PoolIds) != 0 {
		//	t.Fatalf("expected 0 pool ids, got %d", len(rga.PoolIds))
		//}
		//
		//log.Printf("testing DeleteSecurityZone() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		//err = bpClient.DeleteSecurityZone(ctx, zoneId)
		//if err != nil {
		//	t.Fatal(err)
		//}
	}
}
