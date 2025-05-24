// Copyright (c) Juniper Networks, Inc., 2024-2025.
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

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/stretchr/testify/require"
)

func TestCRUDRemoteGateway(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	compare := func(t *testing.T, a, b *TwoStageL3ClosRemoteGatewayData) {
		require.NotNil(t, a)
		require.NotNil(t, b)

		require.Equal(t, a.Label, b.Label, "label does not match")
		require.Equal(t, a.GwIp.String(), b.GwIp.String(), fmt.Sprintf("gateway ip addresses do not match"))
		require.Equal(t, a.GwAsn, b.GwAsn, "gateway ASN does not match")

		require.NotNil(t, b.RouteTypes)
		if a.RouteTypes == nil {
			require.Equal(t, enum.RemoteGatewayRouteTypeAll.Value, b.RouteTypes.Value, "expected default route type not found")
		} else {
			require.Equal(t, a.RouteTypes.Value, b.RouteTypes.Value, "route type does not match")
		}

		require.NotNil(t, b.Ttl)
		if a.Ttl == nil {
			require.Equal(t, uint8(30), *b.Ttl, "expected default ttl not found")
		} else {
			require.Equal(t, *a.Ttl, *b.Ttl, "ttl does not match")
		}

		require.NotNil(t, b.KeepaliveTimer)
		if a.KeepaliveTimer == nil {
			require.Equal(t, uint16(10), *b.KeepaliveTimer, "expected default keepalive timer not found")
		} else {
			require.Equal(t, *a.KeepaliveTimer, *b.KeepaliveTimer, "keepalive timer does not match")
		}

		require.NotNil(t, b.HoldtimeTimer)
		if a.HoldtimeTimer == nil {
			require.Equal(t, uint16(30), *b.HoldtimeTimer, "expected default holdtime timer not found")
		} else {
			require.Equal(t, *a.HoldtimeTimer, *b.HoldtimeTimer, "holdtime timer does not match")
		}

		if a.EvpnInterconnectGroupId == nil {
			require.Nil(t, b.EvpnInterconnectGroupId)
		} else {
			require.Equal(t, *a.EvpnInterconnectGroupId, *b.EvpnInterconnectGroupId, fmt.Sprintf("evpn interconnect group id does not match"))
		}

		compareSlicesAsSets(t, a.LocalGwNodes, b.LocalGwNodes, "mismatched local gateway node sets")
	}

	for _, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			t.Logf("testing remote gateway CRUD methods against Apstra %s", client.client.apiVersion)

			bp := testBlueprintA(ctx, t, client.client)

			leafIds, err := getSystemIdsByRole(ctx, bp, "leaf")
			require.NoError(t, err)

			evpnInterConnectGroupId, err := bp.CreateEvpnInterconnectGroup(ctx, &EvpnInterconnectGroupData{
				Label:       randString(6, "hex"),
				RouteTarget: fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1),
			})
			require.NoError(t, err)
			_ = evpnInterConnectGroupId

			type testStep struct {
				config TwoStageL3ClosRemoteGatewayData
			}

			type testCase struct {
				steps []testStep
			}

			testCases := map[string]testCase{
				"mandatory_fields_only": {
					steps: []testStep{
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:        randString(6, "hex"),
								GwIp:         netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:        rand.Uint32(),
								LocalGwNodes: []ObjectId{leafIds[rand.Intn(len(leafIds))]},
							},
						},
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:        randString(6, "hex"),
								GwIp:         netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:        rand.Uint32(),
								LocalGwNodes: []ObjectId{leafIds[rand.Intn(len(leafIds))]},
							},
						},
					},
				},
				"start_minimal": {
					steps: []testStep{
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:        randString(6, "hex"),
								GwIp:         netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:        rand.Uint32(),
								LocalGwNodes: []ObjectId{leafIds[rand.Intn(len(leafIds))]},
							},
						},
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:                   randString(6, "hex"),
								GwIp:                    netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:                   rand.Uint32(),
								RouteTypes:              &enum.RemoteGatewayRouteTypeAll,
								Ttl:                     toPtr(uint8(rand.Intn(254) + 2)),               // 2-255
								KeepaliveTimer:          toPtr(uint16(rand.Intn(math.MaxUint16) + 1)),   // 1-65535
								HoldtimeTimer:           toPtr(uint16(rand.Intn(math.MaxUint16-2) + 3)), // 3-65535
								Password:                toPtr(randString(6, "hex")),
								EvpnInterconnectGroupId: &evpnInterConnectGroupId,
								LocalGwNodes:            leafIds,
							},
						},
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:                   randString(6, "hex"),
								GwIp:                    netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:                   rand.Uint32(),
								RouteTypes:              &enum.RemoteGatewayRouteTypeAll,
								Ttl:                     toPtr(uint8(rand.Intn(254) + 2)),               // 2-255
								KeepaliveTimer:          toPtr(uint16(rand.Intn(math.MaxUint16) + 1)),   // 1-65535
								HoldtimeTimer:           toPtr(uint16(rand.Intn(math.MaxUint16-2) + 3)), // 3-65535
								Password:                toPtr(randString(6, "hex")),
								EvpnInterconnectGroupId: nil,
								LocalGwNodes:            []ObjectId{leafIds[rand.Intn(len(leafIds))]},
							},
						},
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:                   randString(6, "hex"),
								GwIp:                    netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:                   rand.Uint32(),
								RouteTypes:              &enum.RemoteGatewayRouteTypeAll,
								Ttl:                     toPtr(uint8(rand.Intn(254) + 2)),               // 2-255
								KeepaliveTimer:          toPtr(uint16(rand.Intn(math.MaxUint16) + 1)),   // 1-65535
								HoldtimeTimer:           toPtr(uint16(rand.Intn(math.MaxUint16-2) + 3)), // 3-65535
								Password:                toPtr(randString(6, "hex")),
								EvpnInterconnectGroupId: &evpnInterConnectGroupId,
								LocalGwNodes:            leafIds,
							},
						},
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:        randString(6, "hex"),
								GwIp:         netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:        rand.Uint32(),
								LocalGwNodes: []ObjectId{leafIds[rand.Intn(len(leafIds))]},
							},
						},
					},
				},
				"start_maximal": {
					steps: []testStep{
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:                   randString(6, "hex"),
								GwIp:                    netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:                   rand.Uint32(),
								RouteTypes:              &enum.RemoteGatewayRouteTypeAll,
								Ttl:                     toPtr(uint8(rand.Intn(254) + 2)),               // 2-255
								KeepaliveTimer:          toPtr(uint16(rand.Intn(math.MaxUint16) + 1)),   // 1-65535
								HoldtimeTimer:           toPtr(uint16(rand.Intn(math.MaxUint16-2) + 3)), // 3-65535
								Password:                toPtr(randString(6, "hex")),
								EvpnInterconnectGroupId: &evpnInterConnectGroupId,
								LocalGwNodes:            leafIds,
							},
						},
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:        randString(6, "hex"),
								GwIp:         netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:        rand.Uint32(),
								LocalGwNodes: []ObjectId{leafIds[rand.Intn(len(leafIds))]},
							},
						},
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:                   randString(6, "hex"),
								GwIp:                    netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:                   rand.Uint32(),
								EvpnInterconnectGroupId: &evpnInterConnectGroupId,
								LocalGwNodes:            []ObjectId{leafIds[rand.Intn(len(leafIds))]},
							},
						},
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:        randString(6, "hex"),
								GwIp:         netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:        rand.Uint32(),
								LocalGwNodes: []ObjectId{leafIds[rand.Intn(len(leafIds))]},
							},
						},
						{
							config: TwoStageL3ClosRemoteGatewayData{
								Label:                   randString(6, "hex"),
								GwIp:                    netIpToNetIpAddr(t, randomIpv4()),
								GwAsn:                   rand.Uint32(),
								RouteTypes:              &enum.RemoteGatewayRouteTypeFiveOnly,
								Ttl:                     toPtr(uint8(rand.Intn(254) + 2)),               // 2-255
								KeepaliveTimer:          toPtr(uint16(rand.Intn(math.MaxUint16) + 1)),   // 1-65535
								HoldtimeTimer:           toPtr(uint16(rand.Intn(math.MaxUint16-2) + 3)), // 3-65535
								Password:                toPtr(randString(6, "hex")),
								EvpnInterconnectGroupId: nil,
								LocalGwNodes:            leafIds,
							},
						},
					},
				},
			}

			wg := new(sync.WaitGroup)
			wg.Add(len(testCases))

			var createdIds []ObjectId
			idMutex := new(sync.Mutex)

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					t.Cleanup(func() { wg.Done() })
					// t.Parallel() - do not parallelize - use leaf nodes one-at-a-time

					require.Greater(t, len(tCase.steps), 0)

					t.Log("creating remote gateay")
					id, err := bp.CreateRemoteGateway(ctx, &tCase.steps[0].config)
					require.NoError(t, err)
					idMutex.Lock()
					createdIds = append(createdIds, id)
					idMutex.Unlock()

					t.Log("getting remote gateay")
					gw, err := bp.GetRemoteGateway(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, gw.Id)
					compare(t, &tCase.steps[0].config, gw.Data)

					t.Log("getting remote gateay by name")
					gw, err = bp.GetRemoteGatewayByName(ctx, tCase.steps[0].config.Label)
					require.NoError(t, err)
					require.Equal(t, id, gw.Id)
					compare(t, &tCase.steps[0].config, gw.Data)

					for i, step := range tCase.steps {
						t.Logf("updating remote gateay (step %d)", i)
						err = bp.UpdateRemoteGateway(ctx, gw.Id, &step.config)
						require.NoError(t, err)

						t.Logf("getting remote gateay (step %d)", i)
						gw, err := bp.GetRemoteGateway(ctx, id)
						require.NoError(t, err)
						require.Equal(t, id, gw.Id)
						compare(t, &step.config, gw.Data)
					}
				})
			}

			t.Run("get_all", func(t *testing.T) {
				t.Parallel()

				wg.Wait()

				all, err := bp.GetAllRemoteGateways(ctx)
				require.NoError(t, err)
				require.Equal(t, len(testCases), len(all))

				retrievedIds := make([]ObjectId, len(all))
				for i, o := range all {
					retrievedIds[i] = o.Id
				}

				compareSlicesAsSets(t, createdIds, retrievedIds, "created and retrieved IDs do not match")

				for _, each := range all {
					t.Run("delete_"+each.Id.String(), func(t *testing.T) {
						t.Parallel()

						err = bp.DeleteRemoteGateway(ctx, each.Id)
						require.NoError(t, err)

						var ace ClientErr

						err = bp.DeleteRemoteGateway(ctx, each.Id)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, ErrNotfound, ace.errType)

						_, err = bp.GetRemoteGateway(ctx, each.Id)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, ErrNotfound, ace.errType)

						_, err = bp.GetRemoteGatewayByName(ctx, each.Data.Label)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, ErrNotfound, ace.errType)
					})
				}
			})
		})
	}
}
