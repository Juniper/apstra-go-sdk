//go:build integration

package apstra

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestGetFeatures(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	type testCase struct {
		versionConstraing *version.Constraint
		feature           enum.ApiFeature
		expExists         bool
		expEnabled        bool
		expDisabled       bool
	}

	testCases := map[string]testCase{
		"a": {
			versionConstraing: nil,
			feature:           enum.ApiFeature{},
			expExists:         false,
			expEnabled:        false,
			expDisabled:       false,
		},
	}

	for clientName, client := range clients {
		t.Run(clientName, func(t *testing.T) {
			t.Parallel()

			require.NoError(t, client.client.getFeatures(ctx))
		})
	}
}
