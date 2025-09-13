// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetNodeRenderedConfig(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bp := dctestobj.TestBlueprintI(t, ctx, client.Client)
			leafIds, err := testutils.GetSystemIdsByRole(ctx, bp, "leaf")
			require.NoError(t, err)

			for _, leafId := range leafIds {
				t.Run("leaf_"+leafId.String()+"_staging", func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					config, err := bp.Client().GetNodeRenderedConfig(ctx, bp.Id(), leafId, enum.RenderedConfigTypeStaging)
					require.NoError(t, err)
					lineCount := 0
					scanner := bufio.NewScanner(strings.NewReader(config))
					for scanner.Scan() {
						lineCount++
					}
					require.Greaterf(t, lineCount, 100, "staging config less than 100 lines is sus")
				})

				t.Run("leaf_"+leafId.String()+"_deployed", func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					config, err := bp.Client().GetNodeRenderedConfig(ctx, bp.Id(), leafId, enum.RenderedConfigTypeDeployed)
					require.NoError(t, err)
					lineCount := 0
					scanner := bufio.NewScanner(strings.NewReader(config))
					for scanner.Scan() {
						lineCount++
					}
					require.Greaterf(t, lineCount, 100, "deployed config less than 100 lines is sus")
				})
			}
		})
	}
}
