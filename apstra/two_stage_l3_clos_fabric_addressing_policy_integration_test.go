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
		bpClient, bpDel := testBlueprintA(ctx, t, client.client)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err := bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		log.Printf("testing GetFabricAddressingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		fap, err := bpClient.GetFabricAddressingPolicy(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if fap.EsiMacMsb != defaultEsiMacMsb {
			t.Fatalf("expected mac msb to be %d (apstra default?) got %d", defaultEsiMacMsb, fap.EsiMacMsb)
		}

		if fap.Ipv6Enabled {
			t.Fatal("expected ipv6 to be disabled")
		}

		newMsb := uint8(rand.Intn(100) + 100) // value 100 - 199
		newMsb = newMsb + newMsb%2            // make it even

		log.Printf("testing SetFabricAddressingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.SetFabricAddressingPolicy(ctx, &TwoStageL3ClosFabricAddressingPolicy{
			EsiMacMsb:   newMsb,
			Ipv6Enabled: true,
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetFabricAddressingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		fap, err = bpClient.GetFabricAddressingPolicy(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if newMsb != fap.EsiMacMsb {
			t.Fatalf("new fabric addressing policy mac msb: expected %d got %d", newMsb, fap.EsiMacMsb)
		}

		if !fap.Ipv6Enabled {
			t.Fatal("enabling ipv6 in the fabric addressing policy failed")
		}
	}
}
