//go:build integration
// +build integration

package apstra

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"math/rand"
	"net"
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

		//query := new(PathQuery).
		//	SetBlueprintId(bpClient.Id()).
		//	SetClient(bpClient.Client()).
		//	Node([]QEEAttribute{
		//		NodeTypeSystem.QEEAttribute(),
		//		{Key: "role", Value: QEStringVal("generic")},
		//		{Key: "name", Value: QEStringVal("n_system")},
		//	})
		//
		//var queryResult struct {
		//	Items []struct {
		//		System struct {
		//			Id string `json:"id"`
		//		} `json:"n_system"`
		//	} `json:"items"`
		//}

		t.Logf("determining generic system node IDs in %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		systemIds, err := getSystemIdsByRole(ctx, bpClient, "generic")
		if err != nil {
			t.Fatal(err)
		}

		//err = query.Do(ctx, &queryResult)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//
		//
		//systemIds := make([]string, len(queryResult.Items))
		//for i, item := range queryResult.Items {
		//	systemIds[i] = item.System.Id
		//}
		//t.Logf(`[ "%s" ]`, strings.Join(systemIds, `", "`))

		compareLoopbacks := func(set, get *GenericSystemLoopback) error {
			setIP4 := set.Ipv4Addr != nil
			getIP4 := get.Ipv4Addr != nil
			if (setIP4 || getIP4) && !(setIP4 && getIP4) { // xor
				return fmt.Errorf("generic system loopbacks do not match: a has ipv4: %t, b has IPv4: %t", setIP4, getIP4)
			}
			if setIP4 && getIP4 && set.Ipv4Addr.String() != get.Ipv4Addr.String() {
				return fmt.Errorf("generic system loopbacks do not match: a has ipv4: %s, b has IPv4: %s", set.Ipv4Addr.String(), get.Ipv4Addr.String())
			}

			setIP6 := set.Ipv6Addr != nil
			getIP6 := get.Ipv6Addr != nil
			if (setIP6 || getIP6) && !(setIP6 && getIP6) { // xor
				return fmt.Errorf("generic system loopbacks do not match: a has ipv6: %t, b has IPv6: %t", setIP6, getIP6)
			}
			if setIP6 && getIP6 && set.Ipv6Addr.String() != get.Ipv6Addr.String() {
				return fmt.Errorf("generic system loopbacks do not match: a has ipv6: %s, b has IPv6: %s", set.Ipv6Addr.String(), get.Ipv6Addr.String())
			}

			if set.Ipv6Enabled != get.Ipv6Enabled {
				return fmt.Errorf("generic system loopbacks do not match: a has ipv6 enabled: %t, b has IPv6 enabled: %t", set.Ipv6Enabled, get.Ipv6Enabled)
			}

			return nil
		}

		type testCase struct {
			loopback       GenericSystemLoopback
			apiConstraints version.Constraints
		}

		v4HostMask := net.IPMask{255, 255, 255, 255}
		v6HostMask := net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
		_ = v4HostMask
		_ = v6HostMask

		testCases := map[string]testCase{
			//"v4_only": {
			//	loopback: GenericSystemLoopback{
			//		Ipv4Addr: &net.IPNet{IP: randomIpv4(), Mask: v4HostMask},
			//	},
			//},
			"v6_only": {
				apiConstraints: version.MustConstraints(version.NewConstraint(">=4.1.1")),
				loopback: GenericSystemLoopback{
					Ipv6Addr: &net.IPNet{IP: randomIpv6(), Mask: v6HostMask},
				},
			},
		}

		apiVersion := version.Must(version.NewVersion(client.client.apiVersion))
		systemId := systemIds[rand.Intn(len(systemIds))]

		loopbacks, err := bpClient.GetGenericSystemLoopbacks(ctx, systemId)
		if err != nil {
			t.Fatal(err)
		}
		if len(loopbacks) != 0 {
			t.Fatalf("expected no loopbacks, got %d loopbacks", len(loopbacks))
		}

		t.Run("expect_404", func(t *testing.T) {
			_, err = bpClient.GetGenericSystemLoopback(ctx, systemId, 0)
			if err != nil {
				var ace ClientErr
				if !(errors.As(err, &ace) && ace.Type() == ErrNotfound) {
					t.Fatalf("got an error, but not the expected 404: " + err.Error())
				}
			}
		})

		for tName, tCase := range testCases {
			tName, tCase, bpClient, apiVersion := tName, tCase, *bpClient, *apiVersion
			t.Run(tName, func(t *testing.T) {
				if !tCase.apiConstraints.Check(&apiVersion) {
					t.Skipf("skipping ipv6 test with apstra %s blueprint", &apiVersion)
				}

				err = bpClient.SetGenericSystemLoopback(ctx, systemId, 0, &tCase.loopback)
				if err != nil {
					t.Fatal(err)
				}

				loopback, err := bpClient.GetGenericSystemLoopback(ctx, systemId, 0)
				if err != nil {
					t.Fatal(err)
				}

				if loopback.LoopbackNodeId == "" {
					t.Fatal("loopback node id should not be empty after read")
				}

				// set ipv6 flag (read only) to the expected value based
				// on whether we requested a loopback IPv6 address
				if tCase.loopback.Ipv6Addr != nil {
					tCase.loopback.Ipv6Enabled = true
				}

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
