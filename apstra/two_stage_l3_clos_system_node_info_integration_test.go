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
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemNodeInfo(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing ListAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpIds, err := client.client.ListAllBlueprintIds(ctx)
		require.NoError(t, err)

		var bpClient *TwoStageL3ClosClient
		if len(bpIds) > 0 {
			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bpClient, err = client.client.NewTwoStageL3ClosClient(ctx, bpIds[0])
			require.NoError(t, err)
		} else {
			bpClient = testBlueprintA(ctx, t, client.client)
		}

		log.Printf("testing GetAllSystemNodeInfos() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		nodeInfos, err := bpClient.GetAllSystemNodeInfos(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for nodeId := range nodeInfos {
			log.Printf("testing GetSystemNodeInfo() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			nodeInfo, err := bpClient.GetSystemNodeInfo(ctx, nodeId)
			require.NoError(t, err)
			log.Println(nodeInfo.Id)
		}
	}
}

func TestSetSystemAsn(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		if client.client.ApiVersion() == "4.1.0" {
			continue
		}
		bpClient := testBlueprintB(ctx, t, client.client)

		time.Sleep(2 * time.Second) // todo: fix this terrible workaround for 404s from
		//  /api/blueprints/<id>/experience/web/system-info
		//  shortly after blueprint creation:
		//  {"errors": "Cache for <id> blueprint staging not found"}
		//  see https://apstrktr.atlassian.net/browse/AOS-44024
		//  and https://apstra-eng.slack.com/archives/C2DFCFHJR/p1703621403168039

		log.Printf("testing GetAllSystemNodeInfos() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		nodeInfos, err := bpClient.GetAllSystemNodeInfos(ctx)
		if err != nil {
			t.Fatal(err)
		}

		var systemIds []ObjectId
		for id, info := range nodeInfos {
			if info.Role == SystemRoleGeneric {
				systemIds = append(systemIds, id)
			}
		}

		asnMap := make(map[ObjectId]uint32, len(systemIds))
		for _, id := range systemIds {
			for asnMap[id] == 0 {
				asnMap[id] = rand.Uint32()
			}

			asn := asnMap[id]
			log.Printf("testing SetGenericSystemAsn() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetGenericSystemAsn(ctx, id, &asn)
			if err != nil {
				t.Fatal(err)
			}
		}

		log.Printf("testing GetAllSystemNodeInfos() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		nodeInfos, err = bpClient.GetAllSystemNodeInfos(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for nodeId, asn := range asnMap {
			if nodeInfos[nodeId].Asn == nil {
				t.Fatalf("expected node %q to have asn %d, got nil", nodeId, asn)
			}
			if *nodeInfos[nodeId].Asn != asn {
				t.Fatalf("expected node %q to have asn %d, got %d", nodeId, asn, nodeInfos[nodeId].Asn)
			}

			log.Printf("testing SetGenericSystemAsn() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetGenericSystemAsn(ctx, nodeId, nil)
			if err != nil {
				t.Fatal(err)
			}
		}

		log.Printf("testing GetAllSystemNodeInfos() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		nodeInfos, err = bpClient.GetAllSystemNodeInfos(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for _, id := range systemIds {
			if nodeInfos[id].Asn != nil {
				t.Fatalf("expected node %q to have no ASN, got %d", id, nodeInfos[id].Asn)
			}
		}
	}
}

func TestSetSystemLoopbackIpv4v6(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintG(ctx, t, client.client)

		log.Printf("testing GetAllSystemNodeInfos() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		nodeInfos, err := bpClient.GetAllSystemNodeInfos(ctx)
		if err != nil {
			t.Fatal(err)
		}

		var genericIDs []ObjectId
		for id, info := range nodeInfos {
			if info.Role == SystemRoleGeneric {
				genericIDs = append(genericIDs, id)
			}
		}

		ipv4Map := make(map[ObjectId]net.IPNet, len(genericIDs))
		ipv6Map := make(map[ObjectId]net.IPNet, len(genericIDs))
		for _, id := range genericIDs {
			ipv4Map[id] = net.IPNet{
				IP:   randomIpv4(),
				Mask: net.CIDRMask(rand.Intn(9)+24, 32),
			}

			var v6Mask net.IPMask
			if rand.Int()%2 == 0 {
				v6Mask = net.CIDRMask(64, 128)
			} else {
				v6Mask = net.CIDRMask(128, 128)
			}
			ipv6Map[id] = net.IPNet{
				IP:   randomIpv6(),
				Mask: v6Mask,
			}

			ipv4Net := ipv4Map[id]
			log.Printf("testing SetGenericSystemLoopbackIpv4() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetGenericSystemLoopbackIpv4(ctx, id, &ipv4Net, 0)
			if err != nil {
				t.Fatal(err)
			}

			ipv6Net := ipv6Map[id]
			log.Printf("testing SetGenericSystemLoopbackIpv6() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetGenericSystemLoopbackIpv6(ctx, id, &ipv6Net, 0)
			if err != nil {
				t.Fatal(err)
			}
		}

		log.Printf("testing GetAllSystemNodeInfos() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		nodeInfos, err = bpClient.GetAllSystemNodeInfos(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for nodeId, ip := range ipv4Map {
			if nodeInfos[nodeId].LoopbackIpv4 == nil {
				t.Fatalf("expected node %q to have loopback %s, got nil", nodeId, ip.String())
			}
			if !nodeInfos[nodeId].LoopbackIpv4.IP.Equal(ip.IP) {
				t.Fatalf("expected node %q to have loopback IP %s, got %s", nodeId, ip.IP, nodeInfos[nodeId].LoopbackIpv4.IP)
			}
			if nodeInfos[nodeId].LoopbackIpv4.Mask.String() != ip.Mask.String() {
				t.Fatalf("expected node %q to have loopback IP %s, got %s", nodeId, ip.IP, nodeInfos[nodeId].LoopbackIpv4.IP)
			}

			log.Printf("testing SetGenericSystemLoopbackIpv4() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetGenericSystemLoopbackIpv4(ctx, nodeId, nil, 0)
			if err != nil {
				t.Fatal(err)
			}
		}

		for nodeId, ip := range ipv6Map {
			if nodeInfos[nodeId].LoopbackIpv6 == nil {
				t.Fatalf("expected node %q to have loopback %s, got nil", nodeId, ip.String())
			}
			if !nodeInfos[nodeId].LoopbackIpv6.IP.Equal(ip.IP) {
				t.Fatalf("expected node %q to have loopback IP %s, got %s", nodeId, ip.IP, nodeInfos[nodeId].LoopbackIpv6.IP)
			}
			if nodeInfos[nodeId].LoopbackIpv6.Mask.String() != ip.Mask.String() {
				t.Fatalf("expected node %q to have loopback IP %s, got %s", nodeId, ip.IP, nodeInfos[nodeId].LoopbackIpv6.IP)
			}

			log.Printf("testing SetGenericSystemLoopbackIpv6() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetGenericSystemLoopbackIpv6(ctx, nodeId, nil, 0)
			if err != nil {
				t.Fatal(err)
			}
		}

		log.Printf("testing GetAllSystemNodeInfos() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		nodeInfos, err = bpClient.GetAllSystemNodeInfos(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for _, id := range genericIDs {
			if nodeInfos[id].LoopbackIpv4 != nil {
				t.Fatalf("expected node %q to have no Loopback IPv4, got %s", id, nodeInfos[id].LoopbackIpv4)
			}

			if nodeInfos[id].LoopbackIpv6 != nil {
				t.Fatalf("expected node %q to have no Loopback IPv6, got %s", id, nodeInfos[id].LoopbackIpv6)
			}
		}
	}
}

func TestSetSystemPortChannelIdMinMax(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintG(ctx, t, client.client)

		log.Printf("testing GetAllSystemNodeInfos() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		nodeInfos, err := bpClient.GetAllSystemNodeInfos(ctx)
		if err != nil {
			t.Fatal(err)
		}

		var systemIds []ObjectId
		for id, info := range nodeInfos {
			if info.Role == SystemRoleGeneric {
				systemIds = append(systemIds, id)
			}
		}

		if len(systemIds) == 0 {
			t.Fatal("cannot test - no generic systems found")
		}

		channelIdMinMap := make(map[ObjectId]int, len(systemIds))
		channelIdMaxMap := make(map[ObjectId]int, len(systemIds))
		channelIdMin := 1
		channelIdIncrement := 5
		for _, id := range systemIds {
			channelIdMinMap[id] = channelIdMin
			channelIdMax := channelIdMin + channelIdIncrement
			channelIdMaxMap[id] = channelIdMax
			log.Printf("testing SetGenericSystemPortChannelMinMax against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			log.Printf("testing SetGenericSystemPortChannelMinMax for %s as %d, %d", id, channelIdMin, channelIdMax)

			err = bpClient.SetGenericSystemPortChannelMinMax(ctx, id, channelIdMin, channelIdMax)
			if err != nil {
				t.Fatal(err)
			}
			channelIdMin = channelIdMax + 1
		}

		log.Printf("testing GetAllSystemNodeInfos() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		nodeInfos, err = bpClient.GetAllSystemNodeInfos(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for _, nodeId := range systemIds {
			log.Printf("generic id %s port channel id min = %d, port channel id max = %d", nodeId,
				nodeInfos[nodeId].PortChannelIdMin, nodeInfos[nodeId].PortChannelIdMax)

			if nodeInfos[nodeId].PortChannelIdMin != channelIdMinMap[nodeId] {
				t.Fatalf("Expected Port Channel Id Min %d for Generic System %s Got %d",
					nodeInfos[nodeId].PortChannelIdMin, nodeId, channelIdMinMap[nodeId])
			}

			if nodeInfos[nodeId].PortChannelIdMax != channelIdMaxMap[nodeId] {
				t.Fatalf("Expected Port Channel Id Max %d for Generic System %s Got %d",
					nodeInfos[nodeId].PortChannelIdMax, nodeId, channelIdMaxMap[nodeId])
			}
		}
	}
}

func TestSetGenericSystemLoopbackIPs(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	// this struct represents each test environment we'll run test cases against
	type testBlueprint struct {
		name           string
		testClient     testClient
		bpClient       *TwoStageL3ClosClient
		genericSystems []ObjectId
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(clients))

	testBlueprints := make([]testBlueprint, len(clients))

	var i int
	for name, client := range clients {
		name, client := name, client
		go func(i int) {
			defer wg.Done()

			testBlueprints[i].name = name
			testBlueprints[i].testClient = client
			testBlueprints[i].bpClient = testBlueprintG(ctx, t, client.client)

			nodeInfos, err := testBlueprints[i].bpClient.GetAllSystemNodeInfos(ctx)
			require.NoError(t, err)
			for id, info := range nodeInfos {
				if info.Role == SystemRoleGeneric {
					testBlueprints[i].genericSystems = append(testBlueprints[i].genericSystems, id)
				}
			}
			if len(testBlueprints[i].genericSystems) == 0 {
				t.Error("no generic systems found in blueprint")
			}
		}(i)
		i++
	}

	wg.Wait()

	type testCase struct {
		ip4          *net.IPNet
		ip6          *net.IPNet
		errGetRegexp string
		errSetRegexp string
	}

	testCases := map[string]testCase{
		"empty": {
			errGetRegexp: "Loopback with ID [0-9]+ for system .* is not found",
		},
		"v4_only": {
			ip4: &net.IPNet{
				IP:   randomIpv4(),
				Mask: net.CIDRMask(32, 32),
			},
		},
		"v6_only": {
			ip6: &net.IPNet{
				IP:   randomIpv6(),
				Mask: net.CIDRMask(128, 128),
			},
		},
		"v4_and_v6": {
			ip4: &net.IPNet{
				IP:   randomIpv4(),
				Mask: net.CIDRMask(32, 32),
			},
			ip6: &net.IPNet{
				IP:   randomIpv6(),
				Mask: net.CIDRMask(128, 128),
			},
		},
		"v4_bogus_mask": {
			ip4: &net.IPNet{
				IP:   randomIpv4(),
				Mask: net.CIDRMask(24, 32),
			},
			errSetRegexp: "ip4 value does not contain a valid mask for loopback interfaces",
		},
		"v6_bogus_mask": {
			ip6: &net.IPNet{
				IP:   randomIpv6(),
				Mask: net.CIDRMask(64, 128),
			},
			errSetRegexp: "ip6 value does not contain a valid mask for loopback interfaces",
		},
	}

	for tName, tCase := range testCases {
		tName, tCase := tName, tCase
		t.Run(tName, func(t *testing.T) {
			for _, testBlueprint := range testBlueprints {
				tCase := tCase
				testBlueprint := testBlueprint
				t.Run(testBlueprint.name, func(t *testing.T) {
					t.Parallel()

					gsNodeId := testBlueprint.genericSystems[0]
					bp := testBlueprint.bpClient

					err := bp.SetGenericSystemLoopbackIPs(ctx, gsNodeId, GenericSystemLoopback{
						Ipv4Addr: tCase.ip4,
						Ipv6Addr: tCase.ip6,
					})
					if len(tCase.errSetRegexp) == 0 {
						require.NoError(t, err)
					} else {
						if assert.Error(t, err) {
							assert.Regexp(t, regexp.MustCompile(tCase.errSetRegexp), err.Error())
						}
						return
					}

					lo0, err := bp.GetGenericSystemLoopback(ctx, gsNodeId, 0)
					if len(tCase.errGetRegexp) == 0 {
						require.NoError(t, err)
					} else {
						if assert.Error(t, err) {
							assert.Regexp(t, regexp.MustCompile(tCase.errGetRegexp), err.Error())
						}
						return
					}

					if tCase.ip4 == nil {
						tCase.ip4 = new(net.IPNet)
					}
					if !(lo0.Ipv4Addr.String() == tCase.ip4.String()) {
						t.Errorf("expected %s / got %s", tCase.ip4, lo0.Ipv4Addr)
					}
					if tCase.ip6 == nil {
						tCase.ip6 = new(net.IPNet)
					}
					if !(lo0.Ipv6Addr.String() == tCase.ip6.String()) {
						t.Errorf("expected %s / got %s", tCase.ip6, lo0.Ipv6Addr)
					}
				})
			}
		})
	}
}
