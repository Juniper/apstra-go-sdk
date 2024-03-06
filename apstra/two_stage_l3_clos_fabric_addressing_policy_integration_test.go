//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
)

const (
	defaultEsiMacMsb = 2
)

func TestGetSetGetFAP(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintA(ctx, t, client.client)

		log.Printf("testing GetFabricAddressingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		fap, err := bpClient.GetFabricAddressingPolicy(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if *fap.EsiMacMsb != defaultEsiMacMsb {
			t.Fatalf("expected mac msb to be %d (apstra default?) got %d", defaultEsiMacMsb, fap.EsiMacMsb)
		}

		if *fap.Ipv6Enabled {
			t.Fatal("expected ipv6 to be disabled")
		}

		newMsb := uint8(rand.Intn(100) + 100) // value 100 - 199
		newMsb = newMsb + newMsb%2            // make it even

		ipv6Enabled := true

		var fabricL3Mtu *uint16
		if !fabricL3MtuForbidden.Check(client.client.apiVersion) {
			flm := uint16(rand.Intn(550)*2 + 8000) // even number 8000 - 9100
			fabricL3Mtu = &flm
		}

		log.Printf("testing SetFabricAddressingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.SetFabricAddressingPolicy(ctx, &TwoStageL3ClosFabricAddressingPolicy{
			EsiMacMsb:   &newMsb,
			Ipv6Enabled: &ipv6Enabled,
			FabricL3Mtu: fabricL3Mtu,
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetFabricAddressingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		fap, err = bpClient.GetFabricAddressingPolicy(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if newMsb != *fap.EsiMacMsb {
			t.Fatalf("new fabric addressing policy mac msb: expected %d got %d", newMsb, fap.EsiMacMsb)
		}

		if *fap.Ipv6Enabled != ipv6Enabled {
			t.Fatal("enabling ipv6 in the fabric addressing policy failed")
		}

		if fap.FabricL3Mtu != nil {
			if *fabricL3Mtu != *fap.FabricL3Mtu {
				t.Fatalf("expected fabric MTU %d, got %d", *fabricL3Mtu, *fap.FabricL3Mtu)
			}
		}
	}
}
