//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"
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
		if err != nil {
			t.Fatal(err)
		}

		var bpClient *TwoStageL3ClosClient
		if len(bpIds) > 0 {
			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bpClient, err = client.client.NewTwoStageL3ClosClient(ctx, bpIds[0])
			if err != nil {
				t.Fatal(err)
			}
		} else {
			var deleteFunc func(ctx2 context.Context) error
			bpClient, deleteFunc = testBlueprintA(ctx, t, client.client)
			defer deleteFunc(ctx)
		}

		log.Printf("testing GetAllSystemNodeInfos() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		nodeInfos, err := bpClient.GetAllSystemNodeInfos(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for nodeId := range nodeInfos {
			log.Printf("testing GetSystemNodeInfo() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			nodeInfo, err := bpClient.GetSystemNodeInfo(ctx, nodeId)
			if err != nil {
				t.Fatal(err)
			}
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
		bpClient, deleteFunc := testBlueprintB(ctx, t, client.client)
		defer deleteFunc(ctx)

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

func randomIpv4() net.IP {
	return []byte{
		byte(rand.Intn(222) + 1),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
	}
}

func randomIpv6() net.IP {
	return []byte{
		0x20, 0x01,
		0x0d, 0xb8,
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		byte(rand.Intn(256)), byte(rand.Intn(256)),
	}
}

func TestSetSystemLoopbackIpv4v6(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient, deleteFunc := testBlueprintG(ctx, t, client.client)
		defer deleteFunc(ctx)

		ipv6Enabled := true
		err = bpClient.SetFabricAddressingPolicy(ctx, &TwoStageL3ClosFabricAddressingPolicy{Ipv6Enabled: &ipv6Enabled})
		if err != nil {
			t.Fatal(err)
		}

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

		ipv4Map := make(map[ObjectId]net.IPNet, len(systemIds))
		ipv6Map := make(map[ObjectId]net.IPNet, len(systemIds))
		for _, id := range systemIds {
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

		for _, id := range systemIds {
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
		bpClient, deleteFunc := testBlueprintG(ctx, t, client.client)
		defer deleteFunc(ctx)

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
			log.Printf("testing SetSystemPortChannelMinMax against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			log.Printf("testing SetSystemPortChannelMinMax for %s as %d, %d", id, channelIdMin, channelIdMax)

			err = bpClient.SetSystemPortChannelMinMax(ctx, id, channelIdMin, channelIdMax)
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
