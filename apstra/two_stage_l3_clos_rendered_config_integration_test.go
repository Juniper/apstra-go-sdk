// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestGetNodeRenderedConfig(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bp := testBlueprintI(ctx, t, client.client)
			leafIds, err := getSystemIdsByRole(ctx, bp, "leaf")
			require.NoError(t, err)

			for _, leafId := range leafIds {
				t.Run("leaf_"+leafId.String()+"_staging", func(t *testing.T) {
					t.Parallel()

					stagingConfig, err := bp.GetNodeRenderedConfig(ctx, leafId, enum.RenderedConfigTypeStaging)
					require.NoError(t, err)
					lineCount := 0
					scanner := bufio.NewScanner(strings.NewReader(stagingConfig))
					for scanner.Scan() {
						lineCount++
					}
					require.Greaterf(t, lineCount, 100, "staging config less than 100 lines is sus")
				})

				t.Run("leaf_"+leafId.String()+"_deployed", func(t *testing.T) {
					t.Parallel()

					deployedConfig, err := bp.GetNodeRenderedConfig(ctx, leafId, enum.RenderedConfigTypeDeployed)
					require.NoError(t, err)
					lineCount := 0
					scanner := bufio.NewScanner(strings.NewReader(deployedConfig))
					for scanner.Scan() {
						lineCount++
					}
					require.Greaterf(t, lineCount, 100, "deployed config less than 100 lines is sus")
				})
			}

		})
	}
}
