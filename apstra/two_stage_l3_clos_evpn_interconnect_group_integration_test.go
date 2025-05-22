// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvpnInterconnectGroup(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	compare := func(t *testing.T, a, b *EvpnInterconnectGroupData) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)

		require.Equal(t, a.Label, b.Label)

		require.NotNil(t, b.EsiMac)
		if a.EsiMac != nil {
			require.Equal(t, a.EsiMac, b.EsiMac)
		}

		require.Equal(t, a.RouteTarget, b.RouteTarget)
	}

	for clientName, client := range clients {
		t.Run(clientName, func(t *testing.T) {
			t.Parallel()

			bpClient := testBlueprintA(ctx, t, client.client)

			type testStep struct {
				config EvpnInterconnectGroupData
			}

			type testCase struct {
				steps []testStep
			}

			testCases := map[string]testCase{
				"start_minimal": {
					steps: []testStep{
						{
							config: EvpnInterconnectGroupData{
								Label:       "a" + randString(6, "hex"),
								RouteTarget: fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1),
							},
						},
						{
							config: EvpnInterconnectGroupData{
								Label:       "a" + randString(6, "hex"),
								RouteTarget: fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1),
								EsiMac:      randomHardwareAddr(), // unicast mac
							},
						},
						{
							config: EvpnInterconnectGroupData{
								Label:       "a" + randString(6, "hex"),
								RouteTarget: fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1),
							},
						},
					},
				},
				"start_maximal": {
					steps: []testStep{
						{
							config: EvpnInterconnectGroupData{
								Label:       "a" + randString(6, "hex"),
								RouteTarget: fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1),
								EsiMac:      randomHardwareAddr(),
							},
						},
						{
							config: EvpnInterconnectGroupData{
								Label:       "a" + randString(6, "hex"),
								RouteTarget: fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1),
							},
						},
						{
							config: EvpnInterconnectGroupData{
								Label:       "a" + randString(6, "hex"),
								RouteTarget: fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1),
								EsiMac:      randomHardwareAddr(),
							},
						},
					},
				},
			}

			var createdIds []ObjectId
			idMutex := new(sync.Mutex)

			wg := sync.WaitGroup{}
			wg.Add(len(testCases))

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					t.Cleanup(func() { wg.Done() })
					t.Parallel()

					require.Greater(t, len(tCase.steps), 0)
					id, err := bpClient.CreateEvpnInterconnectGroup(ctx, &tCase.steps[0].config)
					require.NoError(t, err)
					idMutex.Lock()
					createdIds = append(createdIds, id)
					idMutex.Unlock()

					get, err := bpClient.GetEvpnInterconnectGroup(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, get.Id)
					require.NotNil(t, get.Data)
					compare(t, &tCase.steps[0].config, get.Data)

					for i, step := range tCase.steps {
						t.Logf("%s update step %d", tName, i)
						err := bpClient.UpdateEvpnInterconnectGroup(ctx, id, &step.config)
						require.NoError(t, err)

						get, err := bpClient.GetEvpnInterconnectGroup(ctx, id)
						require.NoError(t, err)
						require.Equal(t, id, get.Id)
						require.NotNil(t, get.Data)
						compare(t, &step.config, get.Data)

						get, err = bpClient.GetEvpnInterconnectGroupByName(ctx, step.config.Label)
						require.NoError(t, err)
						require.Equal(t, id, get.Id)
						require.NotNil(t, get.Data)
						compare(t, &step.config, get.Data)
					}
				})
			}

			t.Run("get_all", func(t *testing.T) {
				t.Parallel()

				wg.Wait()

				all, err := bpClient.GetAllEvpnInterconnectGroups(ctx)
				require.NoError(t, err)
				require.Equal(t, len(testCases), len(all))

				retrievedIds := make([]ObjectId, len(all))
				for i, o := range all {
					retrievedIds[i] = o.Id
				}

				compareSlicesAsSets(t, createdIds, retrievedIds, "created and retrieved IDs do not match")

				for _, evpnInterconnectGroup := range all {
					t.Run("delete_"+evpnInterconnectGroup.Id.String(), func(t *testing.T) {

						t.Parallel()

						err = bpClient.DeleteEvpnInterconnectGroup(ctx, evpnInterconnectGroup.Id)
						require.NoError(t, err)

						var ace ClientErr

						err = bpClient.DeleteEvpnInterconnectGroup(ctx, evpnInterconnectGroup.Id)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, ErrNotfound, ace.errType)

						_, err = bpClient.GetEvpnInterconnectGroup(ctx, evpnInterconnectGroup.Id)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, ErrNotfound, ace.errType)

						_, err = bpClient.GetEvpnInterconnectGroupByName(ctx, evpnInterconnectGroup.Data.Label)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, ErrNotfound, ace.errType)
					})
				}
			})
		})

	}
}
