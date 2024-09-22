// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sort"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

func ensureRemoteGatewayDataEqual(t *testing.T, a, b *RemoteGatewayData, skipNilValues bool) {
	sort.Slice(a.LocalGwNodes, func(i, j int) bool {
		if a.LocalGwNodes[i] > a.LocalGwNodes[j] {
			return true
		}
		return false
	})

	sort.Slice(b.LocalGwNodes, func(i, j int) bool {
		if b.LocalGwNodes[i] > b.LocalGwNodes[j] {
			return true
		}
		return false
	})

	compareSlices(t, a.LocalGwNodes, b.LocalGwNodes, "local gateway nodes don't match")

	if a.RouteTypes.Value != b.RouteTypes.Value {
		t.Fatalf("remote gateway route types don't match: %q vs %q", a.RouteTypes.Value, b.RouteTypes.Value)
	}

	if a.GwName != b.GwName {
		t.Fatalf("remote gateway names don't match: %q vs %q", a.GwName, b.GwName)
	}

	if !a.GwIp.Equal(b.GwIp) {
		t.Fatalf("remote gateway IPs don't match: %q vs %q", a.GwIp.String(), b.GwIp.String())
	}

	if a.GwAsn != b.GwAsn {
		t.Fatalf("remote gateway ASNs don't match: %q vs %q", a.GwAsn, b.GwAsn)
	}

	if !possiblyNilValuesMatch(a.Ttl, b.Ttl, skipNilValues) {
		if a == nil || b == nil {
			t.Fatalf("remote gateway TTLs don't match: %v vs. %v", a.Ttl, b.Ttl)
		}
		t.Fatalf("remote gateway TTLs don't match: %d vs. %d", *a.Ttl, *b.Ttl)
	}

	if !possiblyNilValuesMatch(a.KeepaliveTimer, b.KeepaliveTimer, skipNilValues) {
		if a == nil || b == nil {
			t.Fatalf("remote gateway Keepalive timers don't match: %v vs. %v", a.KeepaliveTimer, b.KeepaliveTimer)
		}
		t.Fatalf("remote gateway Keepalive timers don't match: %v vs. %v", a.KeepaliveTimer, b.KeepaliveTimer)
	}

	if !possiblyNilValuesMatch(a.HoldtimeTimer, b.HoldtimeTimer, skipNilValues) {
		if a == nil || b == nil {
			t.Fatalf("remote gateway Holdtime timers don't match: %v vs. %v", a.HoldtimeTimer, b.HoldtimeTimer)
		}
		t.Fatalf("remote gateway Holdtime timers don't match: %v vs. %v", a.HoldtimeTimer, b.HoldtimeTimer)
	}
}

func possiblyNilValuesMatch[A comparable](a, b *A, skipNilValues bool) bool {
	if a == nil && b == nil { // two nil values match
		return true
	}

	if a == nil || b == nil { // only one value is nil
		if skipNilValues {
			return true // but that's okay!
		}
		return false // that's not okay
	}

	// neither value is nil
	if *a == *b {
		return true // they match!
	}

	return false // they don't match
}

func TestCreateDeleteRemoteGateway(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	var remoteGwCfgs []RemoteGatewayData
	for _, routeType := range enum.RemoteGatewayRouteTypesEnum.Members() {
		remoteGwCfgs = append(remoteGwCfgs, RemoteGatewayData{
			RouteTypes:     routeType,
			LocalGwNodes:   nil, // blueprint-specific details set in client loop below
			GwAsn:          rand.Uint32(),
			GwIp:           net.IPv4(uint8(rand.Int()), uint8(rand.Int()), uint8(rand.Int()), uint8(rand.Int())),
			GwName:         randString(5, "hex"),
			Ttl:            nil,
			KeepaliveTimer: nil,
			HoldtimeTimer:  nil,
			Password:       nil,
		})
	}

	randomGwCfgWithNil := RemoteGatewayData{
		RouteTypes:   enum.RemoteGatewayRouteTypesEnum.Members()[rand.Intn(len(enum.RemoteGatewayRouteTypesEnum.Members()))],
		LocalGwNodes: nil,
		GwAsn:        rand.Uint32(),
	}

	randTtl := uint8(rand.Int())
	randKeepaliveTimer := uint16(rand.Int())
	randHoldtimeTimer := uint16(rand.Int())
	randomPassword := randString(5, "hex")
	randomGwCfg := RemoteGatewayData{
		RouteTypes:     enum.RemoteGatewayRouteTypesEnum.Members()[rand.Intn(len(enum.RemoteGatewayRouteTypesEnum.Members()))],
		LocalGwNodes:   nil,
		GwAsn:          rand.Uint32(),
		Ttl:            &randTtl,
		KeepaliveTimer: &randKeepaliveTimer,
		HoldtimeTimer:  &randHoldtimeTimer,
		Password:       &randomPassword,
	}

	for clientName, client := range clients {
		bp := testBlueprintA(ctx, t, client.client)

		// populate blueprint-specific facts into config slice
		localGwNodes, err := getSystemIdsByRole(ctx, bp, "leaf")
		if err != nil {
			t.Fatal(err)
		}
		for i := range remoteGwCfgs {
			remoteGwCfgs[i].LocalGwNodes = localGwNodes
		}

		// populate blueprint-specific facts into random config slice
		randomGwCfg.LocalGwNodes = localGwNodes[rand.Intn(len(localGwNodes)):]
		randomGwCfgWithNil.LocalGwNodes = localGwNodes[rand.Intn(len(localGwNodes)):]

		ids := make([]ObjectId, len(remoteGwCfgs))
		for i, cfg := range remoteGwCfgs {
			log.Printf("testing CreateRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			ids[i], err = bp.CreateRemoteGateway(ctx, &cfg)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			gatewayById, err := bp.GetRemoteGateway(ctx, ids[i])
			if err != nil {
				t.Fatal(err)
			}

			if ids[i] != gatewayById.Id {
				t.Fatalf("expected ID %q, got %q", ids[i], gatewayById.Id)
			}

			ensureRemoteGatewayDataEqual(t, &cfg, gatewayById.Data, true)

			log.Printf("testing GetRemoteGatewayByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			gatewayByName, err := bp.GetRemoteGatewayByName(ctx, remoteGwCfgs[i].GwName)
			if err != nil {
				t.Fatal(err)
			}

			if gatewayById.Id != gatewayByName.Id {
				t.Fatalf("id fetched by ID doesn't match id fetched by name: %q vs. %q", gatewayById.Id, gatewayByName.Id)
			}

			ensureRemoteGatewayDataEqual(t, gatewayById.Data, gatewayByName.Data, false)

			log.Printf("testing UpdateRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bp.UpdateRemoteGateway(ctx, ids[i], gatewayById.Data)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			gatewayById, err = bp.GetRemoteGateway(ctx, ids[i])
			if err != nil {
				t.Fatal(err)
			}

			ensureRemoteGatewayDataEqual(t, gatewayById.Data, gatewayByName.Data, false)

			randomGwCfgWithNil.GwName = randString(5, "hex")
			randomGwCfgWithNil.GwIp = net.IPv4(uint8(rand.Int()), uint8(rand.Int()), uint8(rand.Int()), uint8(rand.Int()))
			log.Printf("testing UpdateRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bp.UpdateRemoteGateway(ctx, ids[i], &randomGwCfgWithNil)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			gatewayById, err = bp.GetRemoteGateway(ctx, ids[i])
			if err != nil {
				t.Fatal(err)
			}

			ensureRemoteGatewayDataEqual(t, gatewayById.Data, &randomGwCfgWithNil, true)

			randomGwCfg.GwName = randString(5, "hex")
			randomGwCfg.GwIp = net.IPv4(uint8(rand.Int()), uint8(rand.Int()), uint8(rand.Int()), uint8(rand.Int()))
			log.Printf("testing UpdateRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bp.UpdateRemoteGateway(ctx, ids[i], &randomGwCfg)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			gatewayById, err = bp.GetRemoteGateway(ctx, ids[i])
			if err != nil {
				t.Fatal(err)
			}

			ensureRemoteGatewayDataEqual(t, gatewayById.Data, &randomGwCfg, false)
		}

		log.Printf("testing GetAllRemoteGateways() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		remoteGws, err := bp.GetAllRemoteGateways(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(remoteGwCfgs) != len(remoteGws) {
			t.Fatalf("expected %d remote gateways, got %d", len(remoteGwCfgs), len(remoteGws))
		}

		for _, id := range ids {
			log.Printf("testing DeleteRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bp.DeleteRemoteGateway(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
		}

		log.Printf("testing GetAllRemoteGateways() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		remoteGws, err = bp.GetAllRemoteGateways(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if 0 != len(remoteGws) {
			t.Fatalf("expected 0 remote gateways, got %d", len(remoteGws))
		}
	}
}
