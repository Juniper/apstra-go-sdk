// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCommitCheck(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for _, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			si, err := client.client.GetAllSystemsInfo(ctx)
			require.NoError(t, err)

			var vExIds []ObjectId
			for _, v := range si {
				if v.Facts.AosHclModel != "Juniper_vEX" {
					continue // wrong platform
				}

				if !v.Status.IsAcknowledged {
					continue // needs ack
				}

				if v.Status.BlueprintActive {
					continue // in use
				}

				vExIds = append(vExIds, ObjectId(v.Id))
			}

			if len(vExIds) == 0 {
				t.Skip("no Juniper_vEX switches available")
			}
			sysId := vExIds[0]

			bp := testBlueprintJ(ctx, t, client.client, sysId)

			err = bp.RunCommitCheck(ctx, nil, nil)
			require.NoError(t, err)

			var ar *TwoStageL3ClosCommitCheckResult
		ccloop:
			for {
				ar, err = bp.CommitCheckResults(ctx)
				require.NoError(t, err)
				require.GreaterOrEqual(t, 1, len(ar.Systems))

				if !ar.Blueprint.State.IsFinal() {
					continue
				}

				for _, v := range ar.Systems {
					if !v.State.IsFinal() {
						continue ccloop
					}
				}

				break
			}

			for sysId := range ar.Systems {
				err = bp.RunCommitCheck(ctx, &sysId, nil)
				require.NoError(t, err)

				_, err := bp.CommitCheckResult(ctx, sysId)
				require.NoError(t, err)
			}
		})
	}
}
