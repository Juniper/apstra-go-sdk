//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestGenericSystemLoopbacks(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		t.Logf("creating test blueprint in %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpClient, bpDelete := testBlueprintH(ctx, t, client.client)
		defer func() {
			err := bpDelete(context.Background())
			if err != nil {
				t.Fatal(err)
			}
		}()

		//szLabelA := randString(5, "hex")
		//securityZoneIdA, err := bpClient.CreateSecurityZone(ctx, &SecurityZoneData{
		//	Label:   szLabelA,
		//	SzType:  SecurityZoneTypeEVPN,
		//	VrfName: szLabelA,
		//})

		//szLabelB := randString(5, "hex")
		//securityZoneIdB, err := bpClient.CreateSecurityZone(ctx, &SecurityZoneData{
		//	Label:   szLabelB,
		//	SzType:  SecurityZoneTypeEVPN,
		//	VrfName: szLabelB,
		//})

		query := new(PathQuery).
			SetBlueprintId(bpClient.Id()).
			SetClient(bpClient.Client()).
			Node([]QEEAttribute{
				NodeTypeSystem.QEEAttribute(),
				{Key: "role", Value: QEStringVal("generic")},
				{Key: "name", Value: QEStringVal("n_system")},
			})

		var queryResult struct {
			Items []struct {
				System struct {
					Id string `json:"id"`
				} `json:"n_system"`
			} `json:"items"`
		}

		t.Logf("determining generic system node IDs in %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = query.Do(ctx, &queryResult)
		if err != nil {
			t.Fatal(err)
		}

		systemIds := make([]string, len(queryResult.Items))
		for i, item := range queryResult.Items {
			systemIds[i] = item.System.Id
		}
		t.Logf(`[ "%s" ]`, strings.Join(systemIds, `", "`))

		compareLoopbacks := func(a, b *GenericSystemLoopback) error {
			aIP4 := a.Ipv4Addr != nil
			bIP4 := b.Ipv4Addr != nil
			if (aIP4 || bIP4) && !(aIP4 && bIP4) { // xor
				return fmt.Errorf("generic system loopbacks do not match: a has ipv4: %t, b has IPv4: %t", aIP4, bIP4)
			}
			if aIP4 && bIP4 && a.Ipv4Addr.String() != b.Ipv4Addr.String() {
				return fmt.Errorf("generic system loopbacks do not match: a has ipv4: %s, b has IPv4: %s", a.Ipv4Addr.String(), b.Ipv4Addr.String())
			}

			aIP6 := a.Ipv6Addr != nil
			bIP6 := b.Ipv6Addr != nil
			if (aIP6 || bIP6) && !(aIP6 && bIP6) { // xor
				return fmt.Errorf("generic system loopbacks do not match: a has ipv6: %t, b has IPv6: %t", aIP6, bIP6)
			}
			if aIP6 && bIP6 && a.Ipv6Addr.String() != b.Ipv6Addr.String() {
				return fmt.Errorf("generic system loopbacks do not match: a has ipv6: %s, b has IPv6: %s", a.Ipv6Addr.String(), b.Ipv6Addr.String())
			}

			if a.Ipv6Enabled != b.Ipv6Enabled {
				return fmt.Errorf("generic system loopbacks do not match: a has ipv6 enabled: %t, b has IPv6 enabled: %t", a.Ipv6Enabled, b.Ipv6Enabled)
			}

			if a.LoopbackNodeId != b.LoopbackNodeId {
				return fmt.Errorf("generic system loopbacks do not match: a has loopback node id: %q, b has loopback node id: %q", a.LoopbackNodeId, b.LoopbackNodeId)
			}

			//if a.SecurityZoneId != b.SecurityZoneId {
			//	return fmt.Errorf("generic system loopbacks do not match: a has security zone: %q, b has security zone: %q", a.SecurityZoneId, b.SecurityZoneId)
			//}

			return nil
		}

		//randomIpv4Ptr := func() *net.IPNet {
		//	ip := randomIpv4()
		//	return &net.IPNet{
		//		IP:   ip,
		//		Mask: net.IPMask{255, 255, 255, 255},
		//	}
		//}

		randomIpv6Ptr := func() *net.IPNet {
			ip := randomIpv6()
			return &net.IPNet{
				IP:   ip,
				Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			}
		}

		type testCase struct {
			loopback GenericSystemLoopback
		}

		testCases := map[string]testCase{
			//"v4_only": {
			//	loopback: GenericSystemLoopback{
			//		Ipv4Addr: randomIpv4Ptr(),
			//	},
			//},
			"v6_only": {
				loopback: GenericSystemLoopback{
					Ipv6Addr: randomIpv6Ptr(),
				},
			},
		}

		for _, systemId := range systemIds {
			loopbacks, err := bpClient.GetGenericSystemLoopbacks(ctx, ObjectId(systemId))
			if err != nil {
				t.Fatal(err)
			}
			if len(loopbacks) != 0 {
				t.Fatalf("expected no loopbacks, got %d loopbacks", len(loopbacks))
			}

			for tName, tCase := range testCases {
				tName, tCase := tName, tCase
				if tCase.loopback.Ipv6Addr != nil && bpClient.client.apiVersion == "4.1.0" {
					t.Log("not testing IPv6 scenarios with apstra 4.1.0")
					continue
				}

				t.Run(tName, func(t *testing.T) {
					// set ipv6 flag (read only) to the expected value based on whether we have a v6 address
					if tCase.loopback.Ipv6Addr != nil {
						tCase.loopback.Ipv6Enabled = true
					}

					err = bpClient.SetGenericSystemLoopback(ctx, ObjectId(systemId), 0, &tCase.loopback)
					if err != nil {
						t.Fatal(err)
					}

					loopback, err := bpClient.GetGenericSystemLoopback(ctx, ObjectId(systemId), 0)
					if err != nil {
						t.Fatal(err)
					}

					if loopback.LoopbackNodeId == "" {
						t.Fatal("lopoback node id should not be empty after read")
					}
					tCase.loopback.LoopbackNodeId = loopback.LoopbackNodeId

					err = compareLoopbacks(&tCase.loopback, loopback)
					if err != nil {
						t.Fatal(err)
					}

					loopbacks, err = bpClient.GetGenericSystemLoopbacks(ctx, ObjectId(systemId))
					if err != nil {
						t.Fatal(err)
					}
					if len(loopbacks) != 1 {
						t.Fatalf("expected 1 loopback, got %d", len(loopbacks))
					}

					loopbackFromMap, ok := loopbacks[0]
					if !ok {
						t.Fatal("loopback 0 not found in map")
					}

					err = compareLoopbacks(&tCase.loopback, &loopbackFromMap)
					if err != nil {
						t.Fatal(err)
					}
				})
			}
		}
	}
}
