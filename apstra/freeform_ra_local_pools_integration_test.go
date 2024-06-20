package apstra

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRaLpA(t *testing.T) {
	var x *FreeformRaLocalIntPool
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)
	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)
		x.Id = "foo"
		x.Data.ResourceType = FFResourceTypeVlan
		x.Data.Chunks[0].Start = 10
		x.Data.Chunks[0].End = 20
		_, err := ffc.CreateLocalIntPool(ctx, x)
		if err != nil {
			return
		}
		require.NoError(t, err)
		t.Log(x.Id)
	}
}