// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestGenericSystemLoopbacks(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			t.Logf("creating test blueprint in %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bpClient := testBlueprintH(ctx, t, client.client)

			t.Logf("determining generic system node IDs in %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
			systemIds, err := getSystemIdsByRole(ctx, bpClient, "generic")
			if err != nil {
				t.Fatal(err)
			}

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

				return nil
			}

			type testCase struct {
				name           string
				loopback       GenericSystemLoopback
				apiConstraints version.Constraints
			}

			v4HostMask := net.IPMask{255, 255, 255, 255}
			v6HostMask := net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
			_ = v4HostMask
			_ = v6HostMask

			testCases := []testCase{ // this is a slice to ensure tests are run in order
				{
					name: "v4_only_1",
					loopback: GenericSystemLoopback{
						Ipv4Addr: &net.IPNet{IP: randomIpv4(), Mask: v4HostMask},
					},
				},
				{
					name: "v4_only_2",
					loopback: GenericSystemLoopback{
						Ipv4Addr: &net.IPNet{IP: randomIpv4(), Mask: v4HostMask},
					},
				},
				{
					name: "v6_only_1",
					loopback: GenericSystemLoopback{
						Ipv6Addr: &net.IPNet{IP: randomIpv6(), Mask: v6HostMask},
					},
				},
				{
					name: "v6_only_2",
					loopback: GenericSystemLoopback{
						Ipv6Addr: &net.IPNet{IP: randomIpv6(), Mask: v6HostMask},
					},
				},
				{name: "none_1"},
				{
					name: "v4_and_v6_1",
					loopback: GenericSystemLoopback{
						Ipv4Addr: &net.IPNet{IP: randomIpv4(), Mask: v4HostMask},
						Ipv6Addr: &net.IPNet{IP: randomIpv6(), Mask: v6HostMask},
					},
				},
				{
					name: "v4_and_v6_2",
					loopback: GenericSystemLoopback{
						Ipv4Addr: &net.IPNet{IP: randomIpv4(), Mask: v4HostMask},
						Ipv6Addr: &net.IPNet{IP: randomIpv6(), Mask: v6HostMask},
					},
				},
				{name: "none_2"},
			}

			systemId := systemIds[rand.Intn(len(systemIds))]

			t.Run("expect_empty_loopback_map", func(t *testing.T) {
				loopbacks, err := bpClient.GetGenericSystemLoopbacks(ctx, systemId)
				if err != nil {
					t.Fatal(err)
				}
				if len(loopbacks) != 0 {
					t.Fatalf("expected no loopbacks, got %d loopbacks", len(loopbacks))
				}
			})

			t.Run("expect_404", func(t *testing.T) {
				_, err = bpClient.GetGenericSystemLoopback(ctx, systemId, 0)
				if err != nil {
					var ace ClientErr
					if !(errors.As(err, &ace) && ace.Type() == ErrNotfound) {
						t.Fatalf("got an error, but not the expected 404: " + err.Error())
					}
				}
			})

			for _, tCase := range testCases {
				tCase, bpClient := tCase, *bpClient
				// do not use t.Parallel()
				t.Run(tCase.name, func(t *testing.T) {
					if !tCase.apiConstraints.Check(bpClient.client.apiVersion) {
						t.Skipf("skipping ipv6 test with apstra %s blueprint", bpClient.client.apiVersion)
					}

					err = bpClient.SetGenericSystemLoopback(ctx, systemId, 0, &tCase.loopback)
					if err != nil {
						t.Fatal(err)
					}

					loopback, err := bpClient.GetGenericSystemLoopback(ctx, systemId, 0)
					if tCase.loopback.Ipv4Addr == nil && tCase.loopback.Ipv6Addr == nil {
						if err == nil {
							t.Fatal("fetching loopback with no addresses defined should have produced an error")
						}
						loopback = nil
					} else {
						if err != nil {
							t.Fatal(err)
						}
					}

					if loopback != nil && loopback.LoopbackNodeId == "" {
						t.Fatal("loopback node id should not be empty after read")
					}

					if loopback != nil {
						err = compareLoopbacks(&tCase.loopback, loopback)
						if err != nil {
							t.Fatal(err)
						}
					}

					loopbacks, err := bpClient.GetGenericSystemLoopbacks(ctx, systemId)
					if err != nil {
						t.Fatal(err)
					}
					if loopback != nil && len(loopbacks) != 1 {
						t.Fatalf("expected 1 loopback, got %d", len(loopbacks))
					}
					if loopback == nil && len(loopbacks) != 0 {
						t.Fatalf("expected 0 loopback, got %d", len(loopbacks))
					}

					if loopback != nil {
						loopbackFromMap, ok := loopbacks[0]
						if !ok {
							t.Fatal("loopback 0 not found in map")
						}

						err = compareLoopbacks(&tCase.loopback, &loopbackFromMap)
						if err != nil {
							t.Fatal(err)
						}
					}
				})
			}
		})
	}
}
