// Copyright (c) Juniper Networks, Inc., 2023-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"log"
	"net/netip"
	"sort"
	"strings"
	"testing"

	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestGetCablingMapLinks(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	linksAreEqual := func(a, b CablingMapLink) bool {
		if a.ID != b.ID {
			return false
		}
		if (a.Label == nil) != (b.Label == nil) {
			return false
		}
		if (a.Label != nil && b.Label != nil) && *a.Label != *b.Label {
			return false
		}
		return true
	}

	for clientName, client := range clients {
		bpClient := testBlueprintB(ctx, t, client.client)

		log.Printf("testing GetCablingMapLinks() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		links, err := bpClient.GetCablingMapLinks(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if len(links) != 16 {
			t.Fatalf("expectex 16 links, got %d", len(links))
		}

		systemIdToLinks := make(map[string][]CablingMapLink)
		for i, link := range links {
			for _, endpoint := range link.Endpoints {
				systemIdToLinks[endpoint.System.ID] = append(systemIdToLinks[endpoint.System.ID], links[i])
			}
		}

		for systemId, linksA := range systemIdToLinks {
			sort.Slice(linksA, func(i, j int) bool {
				return linksA[i].ID < linksA[j].ID
			})

			linksB, err := bpClient.GetCablingMapLinksBySystem(ctx, systemId)
			if err != nil {
				t.Fatal(err)
			}

			if len(linksA) != len(linksB) {
				t.Fatalf("length of linksA (%d) doesn't match length of linksB(%d)", len(linksA), len(linksB))
			}

			sort.Slice(linksB, func(i, j int) bool {
				return linksB[i].ID < linksB[j].ID
			})

			for i := range linksA {
				if !linksAreEqual(linksA[i], linksB[i]) {
					t.Fatalf("linksA[%d] doesn't match linksB[%d]: %v vs. %v", i, i, linksA[i], linksB[i])
				}
			}
		}
	}
}

func TestPatchCablingMapLinks(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	// It's not possible to clear IPv4 or IPv6 addresses from a link:
	// - If a pool is assigned, a pool value takes over.
	// - If no pool is assigned, but a manual value was previously set, the manual value persists until a pool is assigned.
	// This test assumes that specific IPv4 and IPv6 pools have been assigned, checks for those prefixes.
	compareLinks := func(t *testing.T, req, resp CablingMapLink) {
		require.Equal(t, req.Endpoints[0].Interface.ID, resp.Endpoints[0].Interface.ID)
		if req.Endpoints[0].Interface.IfName != nil { // IF Name must match or be nil
			if *req.Endpoints[0].Interface.IfName == "" {
				require.Nil(t, resp.Endpoints[0].Interface.IfName)
			} else {
				require.NotNil(t, resp.Endpoints[0].Interface.IfName)
				require.Equal(t, *req.Endpoints[0].Interface.IfName, *resp.Endpoints[0].Interface.IfName)
			}
		}
		if req.Endpoints[0].Interface.IPv4Addr != nil { // IPv4 Addr must match or come from pool
			require.NotNil(t, resp.Endpoints[0].Interface.IPv4Addr)
			if req.Endpoints[0].Interface.IPv4Addr.IsValid() {
				require.Equal(t, *req.Endpoints[0].Interface.IPv4Addr, *resp.Endpoints[0].Interface.IPv4Addr)
			} else {
				require.True(t, strings.HasPrefix(resp.Endpoints[0].Interface.IPv4Addr.String(), "10."))
			}
		}
		if req.Endpoints[0].Interface.IPv6Addr != nil { // IPv6 Addr must match or come from pool
			require.NotNil(t, resp.Endpoints[0].Interface.IPv6Addr)
			if req.Endpoints[0].Interface.IPv6Addr.IsValid() {
				require.Equal(t, *req.Endpoints[0].Interface.IPv6Addr, *resp.Endpoints[0].Interface.IPv6Addr)
			} else {
				require.True(t, strings.HasPrefix(resp.Endpoints[0].Interface.IPv6Addr.String(), "fc01:a05:fab"))
			}
		}

		require.Equal(t, req.Endpoints[1].Interface.ID, resp.Endpoints[1].Interface.ID)
		if req.Endpoints[1].Interface.IfName != nil { // IF Name must match or be nil
			if *req.Endpoints[1].Interface.IfName == "" {
				require.Nil(t, resp.Endpoints[1].Interface.IfName)
			} else {
				require.NotNil(t, resp.Endpoints[1].Interface.IfName)
				require.Equal(t, *req.Endpoints[1].Interface.IfName, *resp.Endpoints[1].Interface.IfName)
			}
		}
		if req.Endpoints[1].Interface.IPv4Addr != nil { // IPv4 Addr must match or come from pool
			require.NotNil(t, resp.Endpoints[1].Interface.IPv4Addr)
			if req.Endpoints[1].Interface.IPv4Addr.IsValid() {
				require.Equal(t, *req.Endpoints[1].Interface.IPv4Addr, *resp.Endpoints[1].Interface.IPv4Addr)
			} else {
				require.True(t, strings.HasPrefix(resp.Endpoints[1].Interface.IPv4Addr.String(), "10."))
			}
		}
		if req.Endpoints[1].Interface.IPv6Addr != nil { // IPv6 Addr must match or come from pool
			require.NotNil(t, resp.Endpoints[1].Interface.IPv6Addr)
			if req.Endpoints[1].Interface.IPv6Addr.IsValid() {
				require.Equal(t, *req.Endpoints[1].Interface.IPv6Addr, *resp.Endpoints[1].Interface.IPv6Addr)
			} else {
				require.True(t, strings.HasPrefix(resp.Endpoints[1].Interface.IPv6Addr.String(), "fc01:a05:fab"))
			}
		}
	}

	for clientName, client := range clients {
		t.Run(clientName, func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(context.Background(), t)

			bpClient := testBlueprintH(ctx, t, client.client)
			links, err := bpClient.GetCablingMapLinks(ctx)
			require.NoError(t, err)

			// Enable Apstra 6.1+ dual-stack
			if compatibility.SecurityZoneAddressingSupported.Check(client.client.apiVersion) {
				rz, err := bpClient.GetSecurityZoneByVRFName(ctx, "default")
				require.NoError(t, err)
				rz.AddressingSupport = &enum.AddressingSchemeIPv46
				err = bpClient.UpdateSecurityZone(ctx, rz)
				require.NoError(t, err)
			}

			// Assign an IPv4 pool for spine-leaf links
			err = bpClient.SetResourceAllocation(ctx, &ResourceGroupAllocation{
				ResourceGroup: ResourceGroup{Type: ResourceTypeIp4Pool, Name: ResourceGroupNameSpineLeafIp4},
				PoolIds:       []ObjectId{"Private-10_0_0_0-8"},
			})
			require.NoError(t, err)

			// Assign an IPv6 pool for spine-leaf links
			err = bpClient.SetResourceAllocation(ctx, &ResourceGroupAllocation{
				ResourceGroup: ResourceGroup{Type: ResourceTypeIp6Pool, Name: ResourceGroupNameSpineLeafIp6},
				PoolIds:       []ObjectId{"Private-fc01-a05-fab-48"},
			})
			require.NoError(t, err)

			// get a single spine-leaf link for testing
			var req, resp *CablingMapLink
			for _, v := range links {
				require.NotNil(t, v.Role)
				if *v.Role != enum.LinkRoleSpineLeaf {
					continue
				}
				req = &v
				break
			}
			require.NotNil(t, req)
			linkID := req.ID

			// assign random values to patchable fields in the request.
			ipv4Addr := netip.MustParsePrefix(pointer.To(randomSlash31(t)).String())
			ipv6Addr := netip.MustParsePrefix(pointer.To(randomSlash127(t)).String())
			ipv4Base := ipv4Addr.Masked()
			ipv6Base := ipv6Addr.Masked()
			ipv4Next := netip.PrefixFrom(ipv4Base.Addr().Next(), ipv4Base.Bits())
			ipv6Next := netip.PrefixFrom(ipv6Base.Addr().Next(), ipv6Base.Bits())
			req.Endpoints[0].Interface.IfName = pointer.To(randString(6, "hex"))
			req.Endpoints[0].Interface.IPv4Addr = pointer.To(ipv4Base)
			req.Endpoints[0].Interface.IPv6Addr = pointer.To(ipv6Base)
			req.Endpoints[1].Interface.IfName = pointer.To(randString(6, "hex"))
			req.Endpoints[1].Interface.IPv4Addr = pointer.To(ipv4Next)
			req.Endpoints[1].Interface.IPv6Addr = pointer.To(ipv6Next)

			// patch the link
			err = bpClient.PatchCablingMapLinks(ctx, []CablingMapLink{*req})
			require.NoError(t, err)

			// retrieve and check the patched link
			links, err = bpClient.GetCablingMapLinks(ctx)
			require.NoError(t, err)
			for _, v := range links {
				if v.ID == linkID {
					resp = &v
					break
				}
			}
			require.NotNil(t, resp)
			compareLinks(t, *req, *resp)

			// patch the link with a minimal (no-op) patch to ensure that nothing changes
			req.Endpoints[0].Interface.IfName = nil   //    do not modify this field
			req.Endpoints[0].Interface.IPv4Addr = nil //  do not modify this field
			req.Endpoints[0].Interface.IPv6Addr = nil //  do not modify this field
			req.Endpoints[1].Interface.IfName = nil   //    do not modify this field
			req.Endpoints[1].Interface.IPv4Addr = nil //  do not modify this field
			req.Endpoints[1].Interface.IPv6Addr = nil //  do not modify this field
			err = bpClient.PatchCablingMapLinks(ctx, []CablingMapLink{*req})
			require.NoError(t, err)

			// retrieve and check the no-op patched link
			links, err = bpClient.GetCablingMapLinks(ctx)
			require.NoError(t, err)
			for _, v := range links {
				if v.ID == linkID {
					resp = &v
					break
				}
			}
			require.NotNil(t, resp)
			compareLinks(t, *req, *resp)

			// clear the values from the patchable fields using our non-nil empty value signals.
			req.Endpoints[0].Interface.IfName = pointer.To("")
			req.Endpoints[0].Interface.IPv4Addr = new(netip.Prefix)
			req.Endpoints[0].Interface.IPv6Addr = new(netip.Prefix)
			req.Endpoints[1].Interface.IfName = pointer.To("")
			req.Endpoints[1].Interface.IPv4Addr = new(netip.Prefix)
			req.Endpoints[1].Interface.IPv6Addr = new(netip.Prefix)

			// patch the link
			err = bpClient.PatchCablingMapLinks(ctx, []CablingMapLink{*req})
			require.NoError(t, err)

			// retrieve and check the patched link
			links, err = bpClient.GetCablingMapLinks(ctx)
			require.NoError(t, err)
			for _, v := range links {
				if v.ID == linkID {
					resp = &v
					break
				}
			}
			require.NotNil(t, resp)
			compareLinks(t, *req, *resp)
		})
	}
}
