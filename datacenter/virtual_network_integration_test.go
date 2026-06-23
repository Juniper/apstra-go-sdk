// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package datacenter_test

import (
	"context"
	"math/rand"
	"net"
	"sync"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/query"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	comparedatacenter "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/datacenter"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/deepcopy"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestVirtualNetwork_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		constraints []compatibility.Constraint
		create      datacenter.VirtualNetwork
		update      *datacenter.VirtualNetwork
	}

	// Create maps keyed by clent name: blueprints, slices of leaf id, non-default RZ and SZ
	bpMap := make(map[string]*apstra.TwoStageL3ClosClient, len(clients))
	leafsMap := make(map[string][]string, len(clients))
	rzMap := make(map[string]string, len(clients))
	szMap := make(map[string]string, len(clients))
	wg := new(sync.WaitGroup)
	wg.Add(len(clients))
	mu := new(sync.Mutex)
	for _, client := range clients {
		go func() {
			mu.Lock()
			defer mu.Unlock()

			bpMap[client.Name()] = dctestobj.TestBlueprintH(t, ctx, client.Client)
			leafIds, err := query.SystemIdsByRole(ctx, bpMap[client.Name()], "leaf")
			require.NoError(t, err)

			leafsMap[client.Name()] = make([]string, len(leafIds))
			for i, leafId := range leafIds {
				leafsMap[client.Name()][i] = string(leafId)
			}

			rzName := testutils.RandString(6, "hex")
			rzID, err := bpMap[client.Name()].CreateSecurityZone(ctx, datacenter.SecurityZone{
				Label:   rzName,
				Type:    enum.SecurityZoneTypeEVPN,
				VRFName: rzName,
			})
			require.NoError(t, err)
			rzMap[client.Name()] = rzID

			if compatibility.SwitchingZoneSupported.Check(version.Must(version.NewVersion(bpMap[client.Name()].Client().ApiVersion()))) {
				szMap[client.Name()], err = bpMap[client.Name()].CreateSwitchingZone(ctx, datacenter.SwitchingZone{
					MACVRFName:        pointer.To(testutils.RandString(6, "hex")),
					MACVRFServiceType: pointer.To(enum.SwitchingZoneMACVRFServiceTypeVLANBundle),
				})
				require.NoError(t, err)
			}

			wg.Done()
		}()
	}
	wg.Wait()

	testCases := map[string]testCase{
		"simple": {
			create: datacenter.VirtualNetwork{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.VnTypeVlan,
				Bindings: []datacenter.VNBinding{
					{
						AccessSwitchNodeIDs: nil,
						SystemID:            "",
						VLAN:                nil,
					},
				},
			},
		},
		"remove_binding_empty_slice": {
			constraints: []compatibility.Constraint{compatibility.EmptyVnBindingsOk},
			create: datacenter.VirtualNetwork{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.VnTypeVlan,
				Bindings: []datacenter.VNBinding{
					{
						AccessSwitchNodeIDs: nil,
						SystemID:            "",
						VLAN:                nil,
					},
				},
			},
			update: &datacenter.VirtualNetwork{
				Label:    testutils.RandString(6, "hex"),
				Type:     enum.VnTypeVlan,
				Bindings: []datacenter.VNBinding{},
			},
		},
		"remove_binding_nil_slice": {
			constraints: []compatibility.Constraint{compatibility.EmptyVnBindingsOk},
			create: datacenter.VirtualNetwork{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.VnTypeVlan,
				Bindings: []datacenter.VNBinding{
					{
						AccessSwitchNodeIDs: nil,
						SystemID:            "",
						VLAN:                nil,
					},
				},
			},
			update: &datacenter.VirtualNetwork{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.VnTypeVlan,
			},
		},
		"start_minimal_vlan": {
			constraints: []compatibility.Constraint{compatibility.EmptyVnBindingsOk},
			create: datacenter.VirtualNetwork{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.VnTypeVlan,
			},
			update: &datacenter.VirtualNetwork{
				Bindings: []datacenter.VNBinding{
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
				},
				Description: testutils.RandString(6, "hex"),
				DHCPService: true,
				IPv4Enabled: true,
				IPv4Subnet:  pointer.To(testutils.RandomPrefix(t, "10.0.0.0/8", 23)),
				IPv6Enabled: true,
				IPv6Subnet:  pointer.To(testutils.RandomPrefix(t, "3fff::/20", 64)),
				L3MTU:       pointer.To(9002 + (rand.Intn(50) * 2)),
				Label:       testutils.RandString(6, "hex"),
				//ReservedVLAN: pointer.To(uint16(90 + rand.Intn(4000))),
				//RTPolicy: pointer.To(datacenter.RTPolicy{
				//	ImportRTs: testutils.RandomRouteTargets(t, 1, 3),
				//	ExportRTs: testutils.RandomRouteTargets(t, 1, 3),
				//}),
				SecurityZoneID: "nondefault",
				SVIIPs: []datacenter.SVIAddressing{
					{
						IPv4Mode: enum.IPv4SVIModeEnabled,
						IPv6Mode: enum.IPv6SVIModeForced,
					},
				},
				// Tags:                      testutils.RandomStrings(3, 5, 6, "hex"),
				Type: enum.VnTypeVlan,
				// VirtualGatewayIPv4:        make(net.IP, 0), // auto filled based on subnet
				// VirtualGatewayIPv6:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv4Enabled: true,
				VirtualGatewayIPv6Enabled: true,
				VNI:                       nil,
				// VirtualMAC:                testutils.RandomHardwareAddr([]byte{2}, []byte{1}),
			},
		},
		"start_maximal_vlan": {
			constraints: []compatibility.Constraint{compatibility.EmptyVnBindingsOk},
			create: datacenter.VirtualNetwork{
				Bindings: []datacenter.VNBinding{
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
				},
				Description: testutils.RandString(6, "hex"),
				DHCPService: true,
				IPv4Enabled: true,
				IPv4Subnet:  pointer.To(testutils.RandomPrefix(t, "10.0.0.0/8", 23)),
				IPv6Enabled: true,
				IPv6Subnet:  pointer.To(testutils.RandomPrefix(t, "3fff::/20", 64)),
				L3MTU:       pointer.To(9002 + (rand.Intn(50) * 2)),
				Label:       testutils.RandString(6, "hex"),
				//ReservedVLAN: pointer.To(uint16(90 + rand.Intn(4000))),
				//RTPolicy: pointer.To(datacenter.RTPolicy{
				//	ImportRTs: testutils.RandomRouteTargets(t, 1, 3),
				//	ExportRTs: testutils.RandomRouteTargets(t, 1, 3),
				//}),
				SecurityZoneID: "nondefault",
				SVIIPs: []datacenter.SVIAddressing{
					{
						IPv4Mode: enum.IPv4SVIModeEnabled,
						IPv6Mode: enum.IPv6SVIModeForced,
					},
				},
				// Tags:                      testutils.RandomStrings(3, 5, 6, "hex"),
				Type: enum.VnTypeVlan,
				// VirtualGatewayIPv4:        make(net.IP, 0), // auto filled based on subnet
				// VirtualGatewayIPv6:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv4Enabled: true,
				VirtualGatewayIPv6Enabled: true,
				VNI:                       nil,
				// VirtualMAC:                testutils.RandomHardwareAddr([]byte{2}, []byte{1}),
			},
			update: &datacenter.VirtualNetwork{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.VnTypeVlan,
			},
		},
		"start_minimal_vxlan": {
			constraints: []compatibility.Constraint{compatibility.EmptyVnBindingsOk},
			create: datacenter.VirtualNetwork{
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVxlan,
				SecurityZoneID: "nondefault",
			},
			update: &datacenter.VirtualNetwork{
				Bindings: []datacenter.VNBinding{
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
				},
				Description: testutils.RandString(6, "hex"),
				DHCPService: true,
				IPv4Enabled: true,
				IPv4Subnet:  pointer.To(testutils.RandomPrefix(t, "192.0.2.0/24", 28)),
				IPv6Enabled: true,
				IPv6Subnet:  pointer.To(testutils.RandomPrefix(t, "3fff::/20", 64)),
				L3MTU:       pointer.To(9002 + (rand.Intn(50) * 2)),
				Label:       testutils.RandString(6, "hex"),
				// ReservedVLAN: pointer.To(uint16(90 + rand.Intn(4000))),
				RTPolicy: pointer.To(datacenter.RTPolicy{
					ImportRTs: testutils.RandomRouteTargets(t, 1, 3),
					ExportRTs: testutils.RandomRouteTargets(t, 1, 3),
				}),
				SecurityZoneID: "nondefault",
				SVIIPs: []datacenter.SVIAddressing{
					{
						IPv4Mode: enum.IPv4SVIModeEnabled,
						IPv6Mode: enum.IPv6SVIModeForced,
					},
				},
				// Tags:                      testutils.RandomStrings(3, 5, 6, "hex"),
				Type:                      enum.VnTypeVxlan,
				VirtualGatewayIPv4:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv6:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv4Enabled: true,
				VirtualGatewayIPv6Enabled: true,
				VNI:                       pointer.To(uint32(10000 + rand.Intn(1000))),
				VirtualMAC:                testutils.RandomHardwareAddr([]byte{2}, []byte{1}),
			},
		},
		"start_maximal_vxlan": {
			constraints: []compatibility.Constraint{compatibility.EmptyVnBindingsOk},
			create: datacenter.VirtualNetwork{
				Bindings:     []datacenter.VNBinding{{}, {}},
				Description:  testutils.RandString(6, "hex"),
				DHCPService:  true,
				IPv4Enabled:  true,
				IPv4Subnet:   pointer.To(testutils.RandomPrefix(t, "10.0.0.0/8", 23)),
				IPv6Enabled:  true,
				IPv6Subnet:   pointer.To(testutils.RandomPrefix(t, "3fff::/20", 64)),
				L3MTU:        pointer.To(9002 + (rand.Intn(50) * 2)),
				Label:        testutils.RandString(6, "hex"),
				ReservedVLAN: pointer.To(uint16(90 + rand.Intn(4000))),
				RTPolicy: pointer.To(datacenter.RTPolicy{
					ImportRTs: testutils.RandomRouteTargets(t, 1, 3),
					ExportRTs: testutils.RandomRouteTargets(t, 1, 3),
				}),
				SecurityZoneID: "nondefault",
				SVIIPs: []datacenter.SVIAddressing{
					{
						IPv4Mode: enum.IPv4SVIModeEnabled,
						IPv6Mode: enum.IPv6SVIModeForced,
					},
				},
				// Tags:                      testutils.RandomStrings(3, 5, 6, "hex"),
				Type:                      enum.VnTypeVxlan,
				VirtualGatewayIPv4:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv6:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv4Enabled: true,
				VirtualGatewayIPv6Enabled: true,
				VNI:                       pointer.To(uint32(10000 + rand.Intn(1000))),
				VirtualMAC:                testutils.RandomHardwareAddr([]byte{2}, []byte{1}),
			},
			update: &datacenter.VirtualNetwork{
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVxlan,
				SecurityZoneID: "nondefault",
			},
		},
		"start_minimal_vlan_with_binding_same_security_zone_for_4x": {
			create: datacenter.VirtualNetwork{
				Bindings: []datacenter.VNBinding{{}},
				// Description: testutils.RandString(6, "hex"), does not update with 4x
				Label: testutils.RandString(6, "hex"),
				Type:  enum.VnTypeVlan,
			},
			update: &datacenter.VirtualNetwork{
				Bindings: []datacenter.VNBinding{
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
				},
				// Description: testutils.RandString(6, "hex"), does not update with 4x
				DHCPService: true,
				IPv4Enabled: true,
				IPv4Subnet:  pointer.To(testutils.RandomPrefix(t, "10.0.0.0/8", 23)),
				IPv6Enabled: true,
				IPv6Subnet:  pointer.To(testutils.RandomPrefix(t, "3fff::/20", 64)),
				L3MTU:       pointer.To(9002 + (rand.Intn(50) * 2)),
				Label:       testutils.RandString(6, "hex"),
				//ReservedVLAN: pointer.To(uint16(90 + rand.Intn(4000))),
				//RTPolicy: pointer.To(datacenter.RTPolicy{
				//	ImportRTs: testutils.RandomRouteTargets(t, 1, 3),
				//	ExportRTs: testutils.RandomRouteTargets(t, 1, 3),
				//}),
				SVIIPs: []datacenter.SVIAddressing{
					{
						IPv4Mode: enum.IPv4SVIModeEnabled,
						IPv6Mode: enum.IPv6SVIModeForced,
					},
				},
				// Tags:                      testutils.RandomStrings(3, 5, 6, "hex"),
				Type: enum.VnTypeVlan,
				// VirtualGatewayIPv4:        make(net.IP, 0), // auto filled based on subnet
				// VirtualGatewayIPv6:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv4Enabled: true,
				VirtualGatewayIPv6Enabled: true,
				VNI:                       nil,
				// VirtualMAC:                testutils.RandomHardwareAddr([]byte{2}, []byte{1}),
			},
		},
		"start_maximal_vlan_with_binding_same_security_zone_for_4x": {
			create: datacenter.VirtualNetwork{
				Bindings: []datacenter.VNBinding{
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
				},
				// Description: testutils.RandString(6, "hex"), does not update with 4x
				DHCPService: true,
				IPv4Enabled: true,
				IPv4Subnet:  pointer.To(testutils.RandomPrefix(t, "10.0.0.0/8", 23)),
				IPv6Enabled: true,
				IPv6Subnet:  pointer.To(testutils.RandomPrefix(t, "3fff::/20", 64)),
				L3MTU:       pointer.To(9002 + (rand.Intn(50) * 2)),
				Label:       testutils.RandString(6, "hex"),
				//ReservedVLAN: pointer.To(uint16(90 + rand.Intn(4000))),
				//RTPolicy: pointer.To(datacenter.RTPolicy{
				//	ImportRTs: testutils.RandomRouteTargets(t, 1, 3),
				//	ExportRTs: testutils.RandomRouteTargets(t, 1, 3),
				//}),
				SecurityZoneID: "nondefault",
				SVIIPs: []datacenter.SVIAddressing{
					{
						IPv4Mode: enum.IPv4SVIModeEnabled,
						IPv6Mode: enum.IPv6SVIModeForced,
					},
				},
				// Tags:                      testutils.RandomStrings(3, 5, 6, "hex"),
				Type: enum.VnTypeVlan,
				// VirtualGatewayIPv4:        make(net.IP, 0), // auto filled based on subnet
				// VirtualGatewayIPv6:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv4Enabled: true,
				VirtualGatewayIPv6Enabled: true,
				VNI:                       nil,
				// VirtualMAC:                testutils.RandomHardwareAddr([]byte{2}, []byte{1}),
			},
			update: &datacenter.VirtualNetwork{
				Bindings: []datacenter.VNBinding{{}},
				// Description: testutils.RandString(6, "hex"), does not update with 4x
				Label:          testutils.RandString(6, "hex"),
				SecurityZoneID: "nondefault",
				Type:           enum.VnTypeVlan,
			},
		},
		"start_minimal_vxlan_with_binding_same_security_zone_for_4x": {
			create: datacenter.VirtualNetwork{
				Bindings: []datacenter.VNBinding{{}},
				// Description: testutils.RandString(6, "hex"), does not update with 4x
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVxlan,
				SecurityZoneID: "nondefault",
			},
			update: &datacenter.VirtualNetwork{
				Bindings: []datacenter.VNBinding{
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
				},
				// Description: testutils.RandString(6, "hex"), does not update with 4x
				DHCPService: true,
				IPv4Enabled: true,
				IPv4Subnet:  pointer.To(testutils.RandomPrefix(t, "192.0.2.0/24", 28)),
				IPv6Enabled: true,
				IPv6Subnet:  pointer.To(testutils.RandomPrefix(t, "3fff::/20", 64)),
				L3MTU:       pointer.To(9002 + (rand.Intn(50) * 2)),
				Label:       testutils.RandString(6, "hex"),
				// ReservedVLAN: pointer.To(uint16(90 + rand.Intn(4000))),
				RTPolicy: pointer.To(datacenter.RTPolicy{
					ImportRTs: testutils.RandomRouteTargets(t, 1, 3),
					ExportRTs: testutils.RandomRouteTargets(t, 1, 3),
				}),
				SecurityZoneID: "nondefault",
				SVIIPs: []datacenter.SVIAddressing{
					{
						IPv4Mode: enum.IPv4SVIModeEnabled,
						IPv6Mode: enum.IPv6SVIModeForced,
					},
				},
				// Tags:                      testutils.RandomStrings(3, 5, 6, "hex"),
				Type:                      enum.VnTypeVxlan,
				VirtualGatewayIPv4:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv6:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv4Enabled: true,
				VirtualGatewayIPv6Enabled: true,
				VNI:                       pointer.To(uint32(10000 + rand.Intn(1000))),
				VirtualMAC:                testutils.RandomHardwareAddr([]byte{2}, []byte{1}),
			},
		},
		"start_maximal_vxlan_with_binding_same_security_zone_for_4x": {
			create: datacenter.VirtualNetwork{
				Bindings: []datacenter.VNBinding{
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
				},
				// Description: testutils.RandString(6, "hex"), does not update with 4x
				DHCPService:  true,
				IPv4Enabled:  true,
				IPv4Subnet:   pointer.To(testutils.RandomPrefix(t, "10.0.0.0/8", 23)),
				IPv6Enabled:  true,
				IPv6Subnet:   pointer.To(testutils.RandomPrefix(t, "3fff::/20", 64)),
				L3MTU:        pointer.To(9002 + (rand.Intn(50) * 2)),
				Label:        testutils.RandString(6, "hex"),
				ReservedVLAN: pointer.To(uint16(90 + rand.Intn(4000))),
				RTPolicy: pointer.To(datacenter.RTPolicy{
					ImportRTs: testutils.RandomRouteTargets(t, 1, 3),
					ExportRTs: testutils.RandomRouteTargets(t, 1, 3),
				}),
				SecurityZoneID: "nondefault",
				SVIIPs: []datacenter.SVIAddressing{
					{
						IPv4Mode: enum.IPv4SVIModeEnabled,
						IPv6Mode: enum.IPv6SVIModeForced,
					},
				},
				// Tags:                      testutils.RandomStrings(3, 5, 6, "hex"),
				Type:                      enum.VnTypeVxlan,
				VirtualGatewayIPv4:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv6:        make(net.IP, 0), // auto filled based on subnet
				VirtualGatewayIPv4Enabled: true,
				VirtualGatewayIPv6Enabled: true,
				VNI:                       pointer.To(uint32(10000 + rand.Intn(1000))),
				VirtualMAC:                testutils.RandomHardwareAddr([]byte{2}, []byte{1}),
			},
			update: &datacenter.VirtualNetwork{
				Bindings: []datacenter.VNBinding{{}},
				// Description: testutils.RandString(6, "hex"), does not update with 4x
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVxlan,
				SecurityZoneID: "nondefault",
			},
		},
		"clear_bindings_vlan": {
			constraints: []compatibility.Constraint{compatibility.EmptyVnBindingsOk},
			create: datacenter.VirtualNetwork{
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVlan,
				SecurityZoneID: "nondefault",
				Bindings: []datacenter.VNBinding{
					{},
				},
			},
			update: &datacenter.VirtualNetwork{
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVlan,
				SecurityZoneID: "nondefault",
			},
		},
		"clear_bindings_vxlan": {
			constraints: []compatibility.Constraint{compatibility.EmptyVnBindingsOk},
			create: datacenter.VirtualNetwork{
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVxlan,
				SecurityZoneID: "nondefault",
				Bindings: []datacenter.VNBinding{
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
				},
			},
			update: &datacenter.VirtualNetwork{
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVxlan,
				SecurityZoneID: "nondefault",
				Bindings: []datacenter.VNBinding{
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
				},
			},
		},
		"unreserve_vlan": {
			create: datacenter.VirtualNetwork{
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVxlan,
				SecurityZoneID: "nondefault",
				ReservedVLAN:   pointer.To(uint16(90 + rand.Intn(1000))),
				Bindings: []datacenter.VNBinding{
					{}, {}, {},
				},
			},
			update: &datacenter.VirtualNetwork{
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVxlan,
				SecurityZoneID: "nondefault",
				Bindings: []datacenter.VNBinding{
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
				},
			},
		},
		"reserve_vlan": {
			create: datacenter.VirtualNetwork{
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVxlan,
				SecurityZoneID: "nondefault",
				Bindings: []datacenter.VNBinding{
					{VLAN: pointer.To(uint16(90 + rand.Intn(4000)))},
				},
			},
			update: &datacenter.VirtualNetwork{
				Label:          testutils.RandString(6, "hex"),
				Type:           enum.VnTypeVxlan,
				SecurityZoneID: "nondefault",
				ReservedVLAN:   pointer.To(uint16(90 + rand.Intn(1000))),
				Bindings: []datacenter.VNBinding{
					{}, {}, {},
				},
			},
		},
		"add_tags": {
			constraints: []compatibility.Constraint{compatibility.VirtualNetworkAPITags},
			create: datacenter.VirtualNetwork{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.VnTypeVlan,
			},
			update: &datacenter.VirtualNetwork{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.VnTypeVlan,
				Tags:  testutils.RandomStrings(3, 5, 6, "hex"),
			},
		},
		"clear_tags": {
			constraints: []compatibility.Constraint{compatibility.VirtualNetworkAPITags},
			create: datacenter.VirtualNetwork{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.VnTypeVlan,
				Tags:  testutils.RandomStrings(3, 5, 6, "hex"),
			},
			update: &datacenter.VirtualNetwork{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.VnTypeVlan,
			},
		},
		"change_from_default_switching_zone": {
			constraints: []compatibility.Constraint{compatibility.SwitchingZoneSupported},
			create: datacenter.VirtualNetwork{
				Label:           testutils.RandString(6, "hex"),
				SecurityZoneID:  "nondefault",
				SwitchingZoneID: "",
				Type:            enum.VnTypeVxlan,
			},
			update: &datacenter.VirtualNetwork{
				Label:           testutils.RandString(6, "hex"),
				SecurityZoneID:  "nondefault",
				SwitchingZoneID: "nondefault",
				Type:            enum.VnTypeVxlan,
			},
		},
		"change_to_default_switching_zone": {
			constraints: []compatibility.Constraint{compatibility.SwitchingZoneSupported},
			create: datacenter.VirtualNetwork{
				Label:           testutils.RandString(6, "hex"),
				SwitchingZoneID: "nondefault",
				Type:            enum.VnTypeVlan,
			},
			update: &datacenter.VirtualNetwork{
				Label:           testutils.RandString(6, "hex"),
				SwitchingZoneID: "",
				Type:            enum.VnTypeVlan,
			},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			ctx := testutils.ContextWithTestID(ctx, t)
			// t.Parallel() -- don't run test cases in parallel to avoid collisions between VLAN numbers, etc.

			for clientName := range bpMap {
				t.Run(clientName, func(t *testing.T) {
					ctx := testutils.ContextWithTestID(ctx, t)
					t.Parallel()

					for _, constraint := range tCase.constraints {
						if !constraint.Check(version.Must(version.NewVersion(bpMap[clientName].Client().ApiVersion()))) {
							t.Skipf("skipping test case because constraints not satisfied: %v", tCase.constraints)
						}
					}

					// deep copy because we modify these values below
					var update *datacenter.VirtualNetwork
					create := deepcopy.VirtualNetwork(tCase.create)
					if tCase.update != nil {
						update = pointer.To(deepcopy.VirtualNetwork(*tCase.update))
					}

					bp := bpMap[clientName]
					var leafs []string

					// fill in create bindings with random leaf IDs; clone VLAN ID to/from binding as necessary
					leafs = make([]string, len(leafsMap[clientName]))
					copy(leafs, leafsMap[clientName])
					rand.Shuffle(len(leafs), func(i, j int) { leafs[i], leafs[j] = leafs[j], leafs[i] })
					for i := range len(create.Bindings) {
						leaf, ok := slice.Pop(&leafs)
						require.True(t, ok, "test case requires more leafs than available in blueprint")
						create.Bindings[i].SystemID = leaf
						if create.Type == enum.VnTypeVlan {
							if len(create.Bindings) > 0 {
								switch {
								case create.VNI != nil && create.Bindings[0].VLAN == nil:
									create.Bindings[0].VLAN = pointer.ToCopyOf(uint16(*create.VNI))
								case create.VNI == nil && create.Bindings[0].VLAN != nil:
									create.VNI = pointer.ToCopyOf(uint32(*create.Bindings[0].VLAN))
								case create.VNI != nil && create.Bindings[0].VLAN != nil:
									require.Equal(t, *create.VNI, *create.Bindings[0].VLAN, "create test case VNI and binding VLAN must match when type is VLAN")
								}
							}
						}
					}
					// fill in SVI IP info, if any
					if len(create.SVIIPs) > len(create.Bindings) {
						t.Fatalf("test case assigns more SVI IPs (%d) than switches (%d) in create operation", len(create.SVIIPs), len(create.Bindings))
					}
					for i := range len(create.SVIIPs) {
						create.SVIIPs[i].SystemID = create.Bindings[i].SystemID
						if create.SVIIPs[i].IPv4Mode == enum.IPv4SVIModeForced || create.SVIIPs[i].IPv4Mode == enum.IPv4SVIModeEnabled {
							ip := make(net.IP, 4)
							copy(ip, create.IPv4Subnet.IP)
							ip[3] += 3 + uint8(i)
							create.SVIIPs[i].IPv4Addr = &net.IPNet{IP: ip, Mask: create.IPv4Subnet.Mask}
						}
						if create.SVIIPs[i].IPv6Mode == enum.IPv6SVIModeForced || create.SVIIPs[i].IPv6Mode == enum.IPv6SVIModeEnabled {
							ip := make(net.IP, 16)
							copy(ip, create.IPv6Subnet.IP)
							ip[15] += 3 + uint8(i)
							create.SVIIPs[i].IPv6Addr = &net.IPNet{IP: ip, Mask: create.IPv6Subnet.Mask}
						}
					}
					// set gateway addresses
					if create.VirtualGatewayIPv4Enabled && create.VirtualGatewayIPv4 != nil && len(create.VirtualGatewayIPv4) == 0 {
						ip := make(net.IP, 4)
						copy(ip, create.IPv4Subnet.IP)
						ip[3] += uint8(2)
						create.VirtualGatewayIPv4 = ip
					}
					if create.VirtualGatewayIPv6Enabled && create.VirtualGatewayIPv6 != nil && len(create.VirtualGatewayIPv6) == 0 {
						ip := make(net.IP, 16)
						copy(ip, create.IPv6Subnet.IP)
						ip[15] += uint8(2)
						create.VirtualGatewayIPv6 = ip
					}
					// use the reserved vlan for each binding if one has been selected
					if create.ReservedVLAN != nil {
						for i := range len(create.Bindings) {
							create.Bindings[i].VLAN = pointer.ToCopyOf(*create.ReservedVLAN)
						}
					}
					// set the routing zone ID if necessary
					switch create.SecurityZoneID {
					case "nondefault":
						create.SecurityZoneID = rzMap[clientName]
					case "default":
						zoneID, err := bpMap[clientName].DefaultSecurityZoneID(ctx)
						require.NoError(t, err)
						require.NotNil(t, zoneID)
						create.SecurityZoneID = *zoneID
					}
					// set switching zone ID if necessary
					switch create.SwitchingZoneID {
					case "nondefault":
						create.SwitchingZoneID = szMap[clientName]
					case "default":
						zoneID, err := bpMap[clientName].DefaultSwitchingZoneID(ctx)
						require.NoError(t, err)
						require.NotNil(t, zoneID)
						create.SwitchingZoneID = *zoneID
					}

					if update != nil {
						// fill in update bindings with random leaf IDs; clone VLAN ID to/from binding as necessary
						leafs = make([]string, len(leafsMap[clientName]))
						copy(leafs, leafsMap[clientName])
						rand.Shuffle(len(leafs), func(i, j int) { leafs[i], leafs[j] = leafs[j], leafs[i] })
						for i := range len(update.Bindings) {
							leaf, ok := slice.Pop(&leafs)
							require.True(t, ok, "test case requires more leafs than available in blueprint")
							update.Bindings[i].SystemID = leaf
							if update.Type == enum.VnTypeVlan {
								if len(update.Bindings) > 0 {
									switch {
									case update.VNI != nil && update.Bindings[0].VLAN == nil:
										update.Bindings[0].VLAN = pointer.ToCopyOf(uint16(*update.VNI))
									case update.VNI == nil && update.Bindings[0].VLAN != nil:
										update.VNI = pointer.ToCopyOf(uint32(*update.Bindings[0].VLAN))
									case update.VNI != nil && update.Bindings[0].VLAN != nil:
										require.Equal(t, *update.VNI, *update.Bindings[0].VLAN, "update test case VNI and binding VLAN must match when type is VLAN")
									}
								}
							}
						}
						// fill in SVI IP info, if any
						if len(update.SVIIPs) > len(update.Bindings) {
							t.Fatalf("test case assigns more SVI IPs (%d) than switches (%d) in update operation", len(update.SVIIPs), len(update.Bindings))
						}
						for i := range len(update.SVIIPs) {
							update.SVIIPs[i].SystemID = update.Bindings[i].SystemID
							if update.SVIIPs[i].IPv4Mode == enum.IPv4SVIModeForced || update.SVIIPs[i].IPv4Mode == enum.IPv4SVIModeEnabled {
								ip, ipNet, _ := net.ParseCIDR(testutils.RandomHostIP(t, update.IPv4Subnet.String()).String())
								ipNet.IP = ip
								update.SVIIPs[i].IPv4Addr = ipNet
							}
							if update.SVIIPs[i].IPv6Mode == enum.IPv6SVIModeForced || update.SVIIPs[i].IPv6Mode == enum.IPv6SVIModeEnabled {
								ip, ipNet, _ := net.ParseCIDR(testutils.RandomHostIP(t, update.IPv6Subnet.String()).String())
								ipNet.IP = ip
								update.SVIIPs[i].IPv6Addr = ipNet
							}
						}
						// fill in SVI IP info, if any
						if len(update.SVIIPs) > len(update.Bindings) {
							t.Fatalf("test case assigns more SVI IPs (%d) than switches (%d) in update operation", len(update.SVIIPs), len(update.Bindings))
						}
						for i := range len(update.SVIIPs) {
							update.SVIIPs[i].SystemID = update.Bindings[i].SystemID
							if update.SVIIPs[i].IPv4Mode == enum.IPv4SVIModeForced || update.SVIIPs[i].IPv4Mode == enum.IPv4SVIModeEnabled {
								ip := make(net.IP, 4)
								copy(ip, update.IPv4Subnet.IP)
								ip[3] += 3 + uint8(i)
								update.SVIIPs[i].IPv4Addr = &net.IPNet{IP: ip, Mask: update.IPv4Subnet.Mask}
							}
							if update.SVIIPs[i].IPv6Mode == enum.IPv6SVIModeForced || update.SVIIPs[i].IPv6Mode == enum.IPv6SVIModeEnabled {
								ip := make(net.IP, 16)
								copy(ip, update.IPv6Subnet.IP)
								ip[15] += 3 + uint8(i)
								update.SVIIPs[i].IPv6Addr = &net.IPNet{IP: ip, Mask: update.IPv6Subnet.Mask}
							}
						}
						// set gateway addresses
						if update.VirtualGatewayIPv4Enabled && update.VirtualGatewayIPv4 != nil && len(update.VirtualGatewayIPv4) == 0 {
							ip := make(net.IP, 4)
							copy(ip, update.IPv4Subnet.IP)
							ip[3] += uint8(2)
							update.VirtualGatewayIPv4 = ip
						}
						if update.VirtualGatewayIPv6Enabled && update.VirtualGatewayIPv6 != nil && len(update.VirtualGatewayIPv6) == 0 {
							ip := make(net.IP, 16)
							copy(ip, update.IPv6Subnet.IP)
							ip[15] += uint8(2)
							update.VirtualGatewayIPv6 = ip
						}
						// use the reserved vlan for each binding if one has been selected
						if update.ReservedVLAN != nil {
							for i := range len(update.Bindings) {
								update.Bindings[i].VLAN = pointer.ToCopyOf(*update.ReservedVLAN)
							}
						}
						// set the routing zone ID if necessary
						switch update.SecurityZoneID {
						case "nondefault":
							update.SecurityZoneID = rzMap[clientName]
						case "default":
							zoneID, err := bpMap[clientName].DefaultSecurityZoneID(ctx)
							require.NoError(t, err)
							require.NotNil(t, zoneID)
							update.SecurityZoneID = *zoneID
						}
						// set switching zone ID if necessary
						switch update.SwitchingZoneID {
						case "nondefault":
							update.SwitchingZoneID = szMap[clientName]
						case "default":
							szID, err := bpMap[clientName].DefaultSwitchingZoneID(ctx)
							require.NoError(t, err)
							require.NotNil(t, szID)
							update.SwitchingZoneID = *szID
						}
					}

					// create the object
					id, err := bp.CreateVirtualNetwork(ctx, create)
					require.NoError(t, err)
					require.NotEmpty(t, id)

					// save the ID into the object
					require.NoError(t, create.SetID(id))

					// retrieve the object by ID and validate
					obj, err := bp.GetVirtualNetwork(ctx, id)
					require.NoError(t, err)
					idPtr := obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedatacenter.VirtualNetwork(t, create, obj)

					// retrieve the object by label and validate
					obj, err = bp.GetVirtualNetworkByLabel(ctx, create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedatacenter.VirtualNetwork(t, create, obj)

					// retrieve the list of IDs - ours must be in there
					ids, err := bp.ListVirtualNetworks(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := bp.GetVirtualNetworks(ctx)
					require.NoError(t, err)
					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedatacenter.VirtualNetwork(t, create, obj)

					if update != nil {
						// update the object and validate
						require.NoError(t, update.SetID(id))
						require.NotNil(t, update.ID())
						require.Equal(t, id, *update.ID())
						err = bp.UpdateVirtualNetwork(ctx, *update)
						require.NoError(t, err)

						// retrieve the updated object by ID and validate
						obj, err = bp.GetVirtualNetwork(ctx, id)
						require.NoError(t, err)
						idPtr = obj.ID()
						require.NotNil(t, idPtr)
						require.Equal(t, id, *idPtr)
						comparedatacenter.VirtualNetwork(t, *update, obj)
					}

					// delete the object
					err = bp.DeleteVirtualNetwork(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = bp.GetVirtualNetwork(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = bp.GetVirtualNetworkByLabel(ctx, create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = bp.ListVirtualNetworks(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = bp.GetVirtualNetworks(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = bp.UpdateVirtualNetwork(ctx, create)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = bp.DeleteVirtualNetwork(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}
