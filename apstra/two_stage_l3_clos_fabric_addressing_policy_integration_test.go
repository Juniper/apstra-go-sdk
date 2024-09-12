//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/stretchr/testify/require"
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
		if compatibility.GeApstra421.Check(client.client.apiVersion) {
			continue
		}

		bpClient := testBlueprintA(ctx, t, client.client)

		log.Printf("testing GetFabricAddressingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		fap, err := bpClient.GetFabricAddressingPolicy(ctx)
		require.NoError(t, err)

		require.NotNil(t, fap.EsiMacMsb)
		require.Equal(t, *fap.EsiMacMsb, uint8(defaultEsiMacMsb))

		require.NotNil(t, fap.Ipv6Enabled)
		require.False(t, *fap.Ipv6Enabled)

		newMsb := uint8(rand.Intn(100) + 100) // value 100 - 199
		newMsb = newMsb + newMsb%2            // make it even

		fabricL3Mtu := uint16(rand.Intn(550)*2 + 8000) // even number 8000 - 9100

		log.Printf("testing SetFabricAddressingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		require.NoError(t, bpClient.SetFabricAddressingPolicy(ctx, &TwoStageL3ClosFabricAddressingPolicy{
			EsiMacMsb:   &newMsb,
			Ipv6Enabled: toPtr(true),
			FabricL3Mtu: &fabricL3Mtu,
		}))

		log.Printf("testing GetFabricAddressingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		fap, err = bpClient.GetFabricAddressingPolicy(ctx)
		require.NoError(t, err)

		require.NotNil(t, *fap.EsiMacMsb)
		require.Equal(t, newMsb, *fap.EsiMacMsb)

		require.NotNil(t, *fap.Ipv6Enabled)
		require.True(t, *fap.Ipv6Enabled)

		require.NotNil(t, fap.FabricL3Mtu)
		require.Equal(t, fabricL3Mtu, *fap.FabricL3Mtu)
	}
}
