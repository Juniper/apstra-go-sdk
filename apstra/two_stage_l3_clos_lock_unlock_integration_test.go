package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTwoStageL3ClosClient_Lock(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for _, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			bp := testBlueprintA(ctx, t, client.client)
			err := bp.Lock(ctx)
			require.Error(t, err)
			var ace ClientErr
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.Type(), ErrAlreadyLocked)
			require.IsType(t, ErrAlreadyLockedDetail{}, ace.Detail())
			require.NotNil(t, ace.detail.(ErrAlreadyLockedDetail).UserId)
			require.Equal(t, bp.client.ID(), *ace.detail.(ErrAlreadyLockedDetail).UserId)

			err = bp.Unlock(ctx)
			require.NoError(t, err)

			err = bp.Lock(ctx)
			require.NoError(t, err)
		})
	}
}
