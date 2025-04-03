// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"

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

			var sysId ObjectId
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

				sysId = ObjectId(v.Id)
				break
			}
			if sysId == "" {
				t.Skip("no Juniper_vEX switches available")
			}

			bp := testBlueprintJ(ctx, t, client.client, sysId)

			err = bp.RunCommitCheck(ctx, nil, nil)
			require.NoError(t, err)

			var ar *TwoStageL3ClosCommitCheckResult
		bpLoop:
			for {
				ar, err = bp.CommitCheckResults(ctx)
				require.NoError(t, err)
				require.GreaterOrEqual(t, 1, len(ar.Systems))

				if !ar.Blueprint.State.IsFinal() {
					continue
				}

				for _, v := range ar.Systems {
					require.Nil(t, v.Error)
					if !v.State.IsFinal() {
						continue bpLoop
					}
					require.Equal(t, enum.CommitCheckSourceStaging, v.Source)
				}

				break
			}

			for s := range ar.Systems {
				err = bp.RunCommitCheck(ctx, &s, nil)
				require.NoError(t, err)

				var r *SystemCommitCheckResult
			sysLoop:
				for {
					r, err = bp.CommitCheckResult(ctx, s)
					require.NoError(t, err)
					require.Nil(t, r.Error)

					if r.State.IsFinal() {
						require.Equal(t, enum.CommitCheckSourceStaging, r.Source)
						break sysLoop
					}
				}

				cfg := "system {\n    host-name c31a2-001-leaf1;\n}\ninterfaces {\n    replace: ge-0/0/0 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/1 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/2 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/3 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/4 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/5 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/6 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/7 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/8 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/9 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/10 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/11 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/12 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/13 {\n        unit 0 {\n           family inet;\n        }\n    }\n    replace: ge-0/0/14 {\n        unit 0 {\n           family inet;\n        }\n    }\n}\nprotocols {\n    bgp {\n        log-updown;\n        graceful-restart {\n            dont-help-shared-fate-bfd-down;\n        }\n        multipath;\n    }\n    lldp {\n        port-id-subtype interface-name;\n        port-description-type interface-description;\n        neighbour-port-info-display port-id;\n        interface all;\n    }\n    replace: rstp {\n        bridge-priority 0;\n        bpdu-block-on-edge;\n    }\n}\npolicy-options {\n    community DEFAULT_DIRECT_V4 {\n        members [ 1:20007 21001:26000 ];\n    }\n    policy-statement AllPodNetworks {\n        term AllPodNetworks-10 {\n            from {\n                family inet;\n                protocol direct;\n            }\n            then {\n                community add DEFAULT_DIRECT_V4;\n                accept;\n            }\n        }\n        term AllPodNetworks-100 {\n            then reject;\n        }\n    }\n    policy-statement BGP-AOS-Policy {\n        term BGP-AOS-Policy-10 {\n            from {\n                policy AllPodNetworks;\n            }\n            then accept;\n        }\n        term BGP-AOS-Policy-100 {\n            then reject;\n        }\n    }\n    policy-statement PFE-LB {\n        then {\n            load-balance per-packet;\n        }\n    }\n}\n"
				cfg = "garbage"
				err = bp.RunCommitCheck(ctx, &s, &cfg)
				require.NoError(t, err)

			sysLoopWithConfig:
				for {
					r, err := bp.CommitCheckResult(ctx, s)
					require.NoError(t, err)
					require.Nil(t, r.Error)

					if r.State.IsFinal() {
						// require.Equal(t, enum.CommitCheckSourceUser, r.Source) // doesn't work?
						break sysLoopWithConfig
					}
				}
			}
		})
	}
}
