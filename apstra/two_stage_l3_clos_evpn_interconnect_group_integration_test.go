// Copyright (c) Juniper Networks, Inc., 2025-2026.
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

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func TestEvpnInterconnectGroup(t *testing.T) {
	securityZoneCount := 3
	routingPolicyCount := 3

	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	// in order to "wipe out" the association between a security zone and a
	// routing policy the map entry for that security zone must be present.
	populateMissingSZIDs := func(data *EvpnInterconnectGroupData, ids []ObjectId) {
		if data.InterconnectSecurityZones == nil {
			data.InterconnectSecurityZones = make(map[ObjectId]InterconnectSecurityZoneData)
		}
		for _, id := range ids {
			if _, ok := data.InterconnectSecurityZones[id]; ok {
				continue
			}

			data.InterconnectSecurityZones[id] = InterconnectSecurityZoneData{}
		}
	}

	compareSecurityZoneData := func(t *testing.T, a, b InterconnectSecurityZoneData) {
		require.Equal(t, a.RoutingPolicyId, b.RoutingPolicyId)
		require.Equal(t, a.RouteTarget, b.RouteTarget)
		require.Equal(t, a.L3Enabled, b.L3Enabled)
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

		require.Equal(t, len(a.InterconnectSecurityZones), len(b.InterconnectSecurityZones))
		for k, va := range a.InterconnectSecurityZones {
			vb, ok := b.InterconnectSecurityZones[k]
			if !ok {
				t.Fatalf("a has InterconnectSecurityZone %q, but b does not", k)
			}
			compareSecurityZoneData(t, va, vb)
		}
	}

	for clientName, client := range clients {
		t.Run(clientName, func(t *testing.T) {
			t.Parallel()

			t.Logf("Creating blueprint")
			bpClient := testBlueprintA(ctx, t, client.client)
			fs, err := bpClient.GetFabricSettings(ctx)
			require.NoError(t, err)

			fs.EsiMacMsb = toPtr(uint8((rand.Int() & 254) | 2))
			t.Logf("Setting blueprint ESI MAC MSB to %d", *fs.EsiMacMsb)
			err = bpClient.SetFabricSettings(ctx, fs)
			require.NoError(t, err)

			securityZoneIDs := make([]ObjectId, securityZoneCount)
			for i := range securityZoneIDs {
				vrfName := randString(6, "hex")
				id, err := bpClient.CreateSecurityZone(ctx, SecurityZone{
					Label:   vrfName,
					Type:    enum.SecurityZoneTypeEVPN,
					VRFName: vrfName,
				})
				require.NoError(t, err)
				securityZoneIDs[i] = ObjectId(id)
			}

			routingPolicyIDs := make([]ObjectId, routingPolicyCount)
			importPolicies := []DcRoutingPolicyImportPolicy{
				DcRoutingPolicyImportPolicyAll,
				DcRoutingPolicyImportPolicyDefaultOnly,
				DcRoutingPolicyImportPolicyExtraOnly,
			}
			for i := range routingPolicyIDs {
				importPolicy := importPolicies[i%len(importPolicies)]
				id, err := bpClient.CreateRoutingPolicy(ctx, &DcRoutingPolicyData{
					Label:        randString(6, "hex"),
					PolicyType:   DcRoutingPolicyTypeUser,
					ImportPolicy: importPolicy,
					ExportPolicy: DcRoutingExportPolicy{},
				})
				require.NoError(t, err)
				routingPolicyIDs[i] = id
			}

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
								EsiMac:      randomHardwareAddr([]byte{*fs.EsiMacMsb}, []byte{^*fs.EsiMacMsb}), // match policy MAC MSB
								InterconnectSecurityZones: map[ObjectId]InterconnectSecurityZoneData{
									securityZoneIDs[0]: {
										RoutingPolicyId: &routingPolicyIDs[0],
										RouteTarget:     toPtr(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       true,
									},
									securityZoneIDs[1]: {
										RoutingPolicyId: &routingPolicyIDs[1],
										RouteTarget:     toPtr(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       false,
									},
								},
							},
						},
						{
							config: EvpnInterconnectGroupData{
								Label:       "a" + randString(6, "hex"),
								RouteTarget: fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1),
								EsiMac:      randomHardwareAddr([]byte{*fs.EsiMacMsb}, []byte{^*fs.EsiMacMsb}), // match policy MAC MSB
								InterconnectSecurityZones: map[ObjectId]InterconnectSecurityZoneData{
									securityZoneIDs[2]: {
										RoutingPolicyId: &routingPolicyIDs[2],
										RouteTarget:     toPtr(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       false,
									},
									securityZoneIDs[1]: {
										RoutingPolicyId: &routingPolicyIDs[1],
										RouteTarget:     toPtr(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       true,
									},
								},
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
								EsiMac:      randomHardwareAddr([]byte{*fs.EsiMacMsb}, []byte{^*fs.EsiMacMsb}), // match policy MAC MSB
								InterconnectSecurityZones: map[ObjectId]InterconnectSecurityZoneData{
									securityZoneIDs[0]: {
										RoutingPolicyId: &routingPolicyIDs[0],
										RouteTarget:     toPtr(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       true,
									},
									securityZoneIDs[1]: {
										RoutingPolicyId: &routingPolicyIDs[1],
										RouteTarget:     toPtr(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       false,
									},
								},
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
								EsiMac:      randomHardwareAddr([]byte{*fs.EsiMacMsb}, []byte{^*fs.EsiMacMsb}), // match policy MAC MSB
								InterconnectSecurityZones: map[ObjectId]InterconnectSecurityZoneData{
									securityZoneIDs[0]: {
										RoutingPolicyId: &routingPolicyIDs[0],
										RouteTarget:     toPtr(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       true,
									},
									securityZoneIDs[1]: {
										RoutingPolicyId: &routingPolicyIDs[1],
										RouteTarget:     toPtr(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       false,
									},
								},
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

					// use the first config/step for initial creation
					config := tCase.steps[0].config
					populateMissingSZIDs(&config, securityZoneIDs)

					require.Greater(t, len(tCase.steps), 0)
					if config.EsiMac == nil {
						t.Log("creating EVPN Interconnect Group with unspecified ESI MAC")
					} else {
						t.Logf("creating EVPN Interconnect Group with ESI MAC %s", config.EsiMac)
					}

					id, err := bpClient.CreateEvpnInterconnectGroup(ctx, &config)
					require.NoError(t, err)
					idMutex.Lock()
					createdIds = append(createdIds, id)
					idMutex.Unlock()

					get, err := bpClient.GetEvpnInterconnectGroup(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, get.Id)
					require.NotNil(t, get.Data)
					compare(t, &config, get.Data)

					for i, step := range tCase.steps {
						populateMissingSZIDs(&step.config, securityZoneIDs)
						t.Logf("%s update step %d", tName, i)
						if step.config.EsiMac == nil {
							t.Log("updating EVPN Interconnect Group with unspecified ESI MAC")
						} else {
							t.Logf("updating EVPN Interconnect Group with ESI MAC %s", step.config.EsiMac)
						}
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
