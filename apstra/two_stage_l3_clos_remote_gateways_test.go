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

	"golang.org/x/exp/constraints"
)

func checkRemoteGatewayDataAreEqual(t *testing.T, a, b *RemoteGatewayData, skipNilValues bool) {
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

func possiblyNilValuesMatch[A constraints.Integer](a, b *A, skipNilValues bool) bool {
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
	for _, routeType := range RemoteGatewayRouteTypesEnum.Members() {
		remoteGwCfgs = append(remoteGwCfgs, RemoteGatewayData{
			RouteTypes:     routeType,
			LocalGwNodes:   nil, // blueprint-specific details set in client loop below
			GwAsn:          uint32(rand.Int()),
			GwIp:           net.IPv4(uint8(rand.Int()), uint8(rand.Int()), uint8(rand.Int()), uint8(rand.Int())),
			GwName:         randString(5, "hex"),
			Ttl:            nil,
			KeepaliveTimer: nil,
			HoldtimeTimer:  nil,
		})
	}

	for clientName, client := range clients {
		bp, bpDel := testBlueprintA(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		localGwNodes, err := getSystemIdsByRole(ctx, bp, "leaf")
		if err != nil {
			t.Fatal(err)
		}

		ids := make([]ObjectId, len(remoteGwCfgs))
		for i, cfg := range remoteGwCfgs {
			cfg.LocalGwNodes = localGwNodes

			log.Printf("testing createRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			ids[i], err = bp.createRemoteGateway(ctx, cfg.raw())
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing getRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			raw, err := bp.getRemoteGateway(ctx, ids[i])
			if err != nil {
				t.Fatal(err)
			}

			polishedById, err := raw.polish()
			if err != nil {
				t.Fatal(err)
			}

			if ids[i] != polishedById.Id {
				t.Fatalf("expected ID %q, got %q", ids[i], polishedById.Id)
			}

			if cfg.Ttl == nil || cfg.KeepaliveTimer == nil || cfg.HoldtimeTimer == nil {
				checkRemoteGatewayDataAreEqual(t, &cfg, polishedById.Data, true)
			} else {
				checkRemoteGatewayDataAreEqual(t, &cfg, polishedById.Data, false)
			}

			log.Printf("testing getRemoteGatewayByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			raw, err = bp.getRemoteGatewayByName(ctx, remoteGwCfgs[i].GwName)
			if err != nil {
				t.Fatal(err)
			}

			polishedByName, err := raw.polish()
			if err != nil {
				t.Fatal(err)
			}

			if polishedById.Id != polishedByName.Id {
				t.Fatalf("id fetched by ID doesn't match id fetched by name: %q vs. %q", polishedById.Id, polishedByName.Id)
			}

			checkRemoteGatewayDataAreEqual(t, polishedById.Data, polishedByName.Data, false)
		}

		log.Printf("testing getAllRemoteGateways() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		remoteGws, err := bp.getAllRemoteGateways(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(remoteGwCfgs) != len(remoteGws) {
			t.Fatalf("expected %d remote gateways, got %d", len(remoteGwCfgs), len(remoteGws))
		}

		for _, id := range ids {
			log.Printf("testing deleteRemoteGateway() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bp.deleteRemoteGateway(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
		}

		log.Printf("testing getAllRemoteGateways() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		remoteGws, err = bp.getAllRemoteGateways(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if 0 != len(remoteGws) {
			t.Fatalf("expected 0 remote gateways, got %d", len(remoteGws))
		}
	}
}
