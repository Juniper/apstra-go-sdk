package apstra

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRaLocalPoolGenA(t *testing.T) {
	var x *FreeformRaLocalPoolGenerator
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)
	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)
		x.Id = "foo"
		x.Data.ResourceType = FFResourceTypeVlan
		x.Data.PoolType = "integer"
		x.Data.Scope = "match(node('interface', tag=has_any(['tag1', 'tag2', 'tag3']), name='target').out('link').node('link', role='internal'))"
		x.Data.Chunks[0].Start = 10
		x.Data.Chunks[0].End = 20
		_, err := ffc.CreateLocalPoolGenerator(ctx, x)
		if err != nil {
			return
		}
		require.NoError(t, err)
		t.Log(x.Id)
	}
}
