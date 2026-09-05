// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	compareapstra "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/apstra"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestEvpnInterconnectGroup(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	securityZoneCount := 3
	routingPolicyCount := 3

	// in order to "wipe out" the association between a security zone and a
	// routing policy the map entry for that security zone must be present.
	populateMissingSZIDs := func(data apstra.EVPNInterconnectGroup, ids []string) {
		if data.InterconnectSecurityZones == nil {
			data.InterconnectSecurityZones = make(map[string]apstra.InterconnectSecurityZone)
		}
		for _, id := range ids {
			if _, ok := data.InterconnectSecurityZones[id]; ok {
				continue
			}

			data.InterconnectSecurityZones[id] = apstra.InterconnectSecurityZone{}
		}
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			t.Logf("Creating blueprint")
			bpClient := dctestobj.TestBlueprintA(t, ctx, client.Client)
			fs, err := bpClient.GetFabricSettings(ctx)
			require.NoError(t, err)

			fs.EsiMacMsb = pointer.To(uint8((rand.Int() & 254) | 2))
			t.Logf("Setting blueprint ESI MAC MSB to %d", *fs.EsiMacMsb)
			err = bpClient.SetFabricSettings(ctx, fs)
			require.NoError(t, err)

			securityZoneIDs := make([]string, securityZoneCount)
			for i := range securityZoneIDs {
				vrfName := testutils.RandString(6, "hex")
				id, err := bpClient.CreateSecurityZone(ctx, datacenter.SecurityZone{
					Label:   vrfName,
					Type:    enum.SecurityZoneTypeEVPN,
					VRFName: vrfName,
				})
				require.NoError(t, err)
				securityZoneIDs[i] = id
			}

			routingPolicyIDs := make([]string, routingPolicyCount)
			importPolicies := []apstra.DcRoutingPolicyImportPolicy{
				apstra.DcRoutingPolicyImportPolicyAll,
				apstra.DcRoutingPolicyImportPolicyDefaultOnly,
				apstra.DcRoutingPolicyImportPolicyExtraOnly,
			}
			for i := range routingPolicyIDs {
				importPolicy := importPolicies[i%len(importPolicies)]
				id, err := bpClient.CreateRoutingPolicy(ctx, &apstra.DcRoutingPolicyData{
					Label:        testutils.RandString(6, "hex"),
					PolicyType:   apstra.DcRoutingPolicyTypeUser,
					ImportPolicy: importPolicy,
					ExportPolicy: apstra.DcRoutingExportPolicy{},
				})
				require.NoError(t, err)
				routingPolicyIDs[i] = string(id)
			}

			type testStep struct {
				config apstra.EVPNInterconnectGroup
			}

			type testCase struct {
				steps []testStep
			}

			testCases := map[string]testCase{
				"start_minimal": {
					steps: []testStep{
						{
							config: apstra.EVPNInterconnectGroup{
								Label:       pointer.To("a" + testutils.RandString(6, "hex")),
								RouteTarget: pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
							},
						},
						{
							config: apstra.EVPNInterconnectGroup{
								Label:       pointer.To("a" + testutils.RandString(6, "hex")),
								RouteTarget: pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
								ESIMAC:      testutils.RandomHardwareAddr([]byte{*fs.EsiMacMsb}, []byte{^*fs.EsiMacMsb}), // match policy MAC MSB
								InterconnectSecurityZones: map[string]apstra.InterconnectSecurityZone{
									securityZoneIDs[0]: {
										RoutingPolicyID: &routingPolicyIDs[0],
										RouteTarget:     pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       true,
									},
									securityZoneIDs[1]: {
										RoutingPolicyID: &routingPolicyIDs[1],
										RouteTarget:     pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       false,
									},
								},
							},
						},
						{
							config: apstra.EVPNInterconnectGroup{
								Label:       pointer.To("a" + testutils.RandString(6, "hex")),
								RouteTarget: pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
								ESIMAC:      testutils.RandomHardwareAddr([]byte{*fs.EsiMacMsb}, []byte{^*fs.EsiMacMsb}), // match policy MAC MSB
								InterconnectSecurityZones: map[string]apstra.InterconnectSecurityZone{
									securityZoneIDs[2]: {
										RoutingPolicyID: &routingPolicyIDs[2],
										RouteTarget:     pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       false,
									},
									securityZoneIDs[1]: {
										RoutingPolicyID: &routingPolicyIDs[1],
										RouteTarget:     pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       true,
									},
								},
							},
						},
						{
							config: apstra.EVPNInterconnectGroup{
								Label:       pointer.To("a" + testutils.RandString(6, "hex")),
								RouteTarget: pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
							},
						},
					},
				},
				"start_maximal": {
					steps: []testStep{
						{
							config: apstra.EVPNInterconnectGroup{
								Label:       pointer.To("a" + testutils.RandString(6, "hex")),
								RouteTarget: pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
								ESIMAC:      testutils.RandomHardwareAddr([]byte{*fs.EsiMacMsb}, []byte{^*fs.EsiMacMsb}), // match policy MAC MSB
								InterconnectSecurityZones: map[string]apstra.InterconnectSecurityZone{
									securityZoneIDs[0]: {
										RoutingPolicyID: &routingPolicyIDs[0],
										RouteTarget:     pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       true,
									},
									securityZoneIDs[1]: {
										RoutingPolicyID: &routingPolicyIDs[1],
										RouteTarget:     pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       false,
									},
								},
							},
						},
						{
							config: apstra.EVPNInterconnectGroup{
								Label:       pointer.To("a" + testutils.RandString(6, "hex")),
								RouteTarget: pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
							},
						},
						{
							config: apstra.EVPNInterconnectGroup{
								Label:       pointer.To("a" + testutils.RandString(6, "hex")),
								RouteTarget: pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
								ESIMAC:      testutils.RandomHardwareAddr([]byte{*fs.EsiMacMsb}, []byte{^*fs.EsiMacMsb}), // match policy MAC MSB
								InterconnectSecurityZones: map[string]apstra.InterconnectSecurityZone{
									securityZoneIDs[0]: {
										RoutingPolicyID: &routingPolicyIDs[0],
										RouteTarget:     pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       true,
									},
									securityZoneIDs[1]: {
										RoutingPolicyID: &routingPolicyIDs[1],
										RouteTarget:     pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)),
										L3Enabled:       false,
									},
								},
							},
						},
					},
				},
			}

			var createdIds []string
			idMutex := new(sync.Mutex)

			wg := sync.WaitGroup{}
			wg.Add(len(testCases))

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					t.Cleanup(func() { wg.Done() })
					t.Parallel()

					// use the first config/step for initial creation
					config := tCase.steps[0].config
					populateMissingSZIDs(config, securityZoneIDs)

					require.Greater(t, len(tCase.steps), 0)
					if config.ESIMAC == nil {
						t.Log("creating EVPN Interconnect Group with unspecified ESI MAC")
					} else {
						t.Logf("creating EVPN Interconnect Group with ESI MAC %s", config.ESIMAC)
					}

					id, err := bpClient.CreateEVPNInterconnectGroup(ctx, config)
					require.NoError(t, err)
					idMutex.Lock()
					createdIds = append(createdIds, id)
					idMutex.Unlock()

					get, err := bpClient.GetEVPNInterconnectGroup(ctx, id)
					require.NoError(t, err)
					require.NotNil(t, get.ID())
					require.Equal(t, id, *get.ID())
					compareapstra.EVPNInterconnectGroup(t, config, get)

					for i, step := range tCase.steps {
						require.NoError(t, step.config.SetID(id))
						populateMissingSZIDs(step.config, securityZoneIDs)
						t.Logf("%s update step %d", tName, i)
						if step.config.ESIMAC == nil {
							t.Log("updating EVPN Interconnect Group with unspecified ESI MAC")
						} else {
							t.Logf("updating EVPN Interconnect Group with ESI MAC %s", step.config.ESIMAC)
						}
						err = bpClient.UpdateEVPNInterconnectGroup(ctx, step.config)
						require.NoError(t, err)

						get, err = bpClient.GetEVPNInterconnectGroup(ctx, id)
						require.NoError(t, err)
						require.NotNil(t, get.ID())
						require.Equal(t, id, *get.ID())
						compareapstra.EVPNInterconnectGroup(t, step.config, get)

						require.NotNil(t, step.config.Label)
						get, err = bpClient.GetEVPNInterconnectGroupByLabel(ctx, *step.config.Label)
						require.NoError(t, err)
						require.NotNil(t, get.ID())
						require.Equal(t, id, *get.ID())
						compareapstra.EVPNInterconnectGroup(t, step.config, get)
					}
				})
			}

			t.Run("get_all", func(t *testing.T) {
				t.Parallel()

				wg.Wait()

				all, err := bpClient.GetEVPNInterconnectGroups(ctx)
				require.NoError(t, err)
				require.Equal(t, len(testCases), len(all))

				retrievedIds := make([]string, len(all))
				for i, o := range all {
					require.NotNil(t, o.ID())
					retrievedIds[i] = *o.ID()
				}

				compare.SlicesAsSets(t, createdIds, retrievedIds, "created and retrieved IDs do not match")

				for _, evpnInterconnectGroup := range all {
					t.Run("delete_"+*evpnInterconnectGroup.ID(), func(t *testing.T) {
						t.Parallel()

						err = bpClient.DeleteEVPNInterconnectGroup(ctx, *evpnInterconnectGroup.ID())
						require.NoError(t, err)

						var ace apstra.ClientErr

						err = bpClient.DeleteEVPNInterconnectGroup(ctx, *evpnInterconnectGroup.ID())
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, apstra.ErrNotfound, ace.Type())

						_, err = bpClient.GetEVPNInterconnectGroup(ctx, *evpnInterconnectGroup.ID())
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, apstra.ErrNotfound, ace.Type())

						require.NotNil(t, evpnInterconnectGroup.Label)
						_, err = bpClient.GetEVPNInterconnectGroupByLabel(ctx, *evpnInterconnectGroup.Label)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, apstra.ErrNotfound, ace.Type())
					})
				}
			})
		})
	}
}

func TestEvpnInterconnectGroupErrors(t *testing.T) {
	ctx := context.Background()

	bogusRZID := "bogus_rz"
	bogusRPID := "bogus_rp"
	bogusIGID := "bogus_ig"
	bogusVNID := "bogus_vn"

	clients := testclient.GetTestClients(t, ctx)
	for _, c := range clients {
		t.Run(c.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bp := dctestobj.TestBlueprintA(t, ctx, c.Client)
			rzID := dctestobj.TestSecurityZoneA(t, ctx, bp)
			rpID := dctestobj.TestRoutingPolicyA(t, ctx, bp)
			defaultRZ, err := bp.GetDefaultSecurityZone(ctx)
			require.NoError(t, err)
			require.NotNil(t, defaultRZ.ID())

			var ace apstra.ClientErr

			// Test Get with bogus group ID.
			t.Run("get_with_bogus_group_id", func(t *testing.T) {
				ctx := testutils.ContextWithTestID(ctx, t)

				_, err := bp.GetEVPNInterconnectGroup(ctx, bogusIGID)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, apstra.ErrNotfound, ace.Type())
				require.NotNil(t, ace.Detail())
				require.IsType(t, apstra.ErrNotFoundDetail{}, ace.Detail())
				require.Equal(t, apstra.NodeTypeEvpnInterconnectGroup, ace.Detail().(apstra.ErrNotFoundDetail).Type)
				require.Equal(t, bogusIGID, ace.Detail().(apstra.ErrNotFoundDetail).ID)
			})

			// Test Update with bogus group ID.
			t.Run("update_with_bogus_group_id", func(t *testing.T) {
				ctx := testutils.ContextWithTestID(ctx, t)

				var u apstra.EVPNInterconnectGroup
				_ = u.SetID(bogusIGID)
				err := bp.UpdateEVPNInterconnectGroup(ctx, u)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, apstra.ErrNotfound, ace.Type())
				require.NotNil(t, ace.Detail())
				require.IsType(t, apstra.ErrNotFoundDetail{}, ace.Detail())
				require.Equal(t, apstra.NodeTypeEvpnInterconnectGroup, ace.Detail().(apstra.ErrNotFoundDetail).Type)
				require.Equal(t, bogusIGID, ace.Detail().(apstra.ErrNotFoundDetail).ID)
			})

			// Test Delete with bogus group ID.
			t.Run("delete_with_bogus_group_id", func(t *testing.T) {
				ctx := testutils.ContextWithTestID(ctx, t)

				err := bp.DeleteEVPNInterconnectGroup(ctx, bogusIGID)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, apstra.ErrNotfound, ace.Type())
				require.NotNil(t, ace.Detail())
				require.IsType(t, apstra.ErrNotFoundDetail{}, ace.Detail())
				require.Equal(t, apstra.NodeTypeEvpnInterconnectGroup, ace.Detail().(apstra.ErrNotFoundDetail).Type)
				require.Equal(t, bogusIGID, ace.Detail().(apstra.ErrNotFoundDetail).ID)
			})

			// Test Create with bogus Virtual Network ID.
			t.Run("create_with_bogus_vn_id", func(t *testing.T) {
				ctx := testutils.ContextWithTestID(ctx, t)

				_, err := bp.CreateEVPNInterconnectGroup(ctx, apstra.EVPNInterconnectGroup{
					Label:       pointer.To("a" + testutils.RandString(6, "hex")),
					RouteTarget: pointer.To(testutils.RandomRouteTarget(t)),
					InterconnectVirtualNetworks: map[string]apstra.InterconnectVirtualNetwork{
						bogusVNID: {
							L2Enabled:      true,
							L3Enabled:      true,
							TranslationVNI: pointer.To(uint32(rand.Intn(4000) + 1)),
						},
					},
				})
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ace.Type(), apstra.ErrNotfound)
				require.NotNil(t, ace.Detail())
				require.IsType(t, apstra.ErrNotFoundDetail{}, ace.Detail())
				require.Equal(t, apstra.NodeTypeVirtualNetwork, ace.Detail().(apstra.ErrNotFoundDetail).Type)
				require.Equal(t, bogusVNID, ace.Detail().(apstra.ErrNotFoundDetail).ID)
			})

			// Test Create with bogus Routing Zone ID.
			t.Run("create_with_bogus_routing_zone", func(t *testing.T) {
				ctx := testutils.ContextWithTestID(ctx, t)

				_, err = bp.CreateEVPNInterconnectGroup(ctx, apstra.EVPNInterconnectGroup{
					Label:       pointer.To("a" + testutils.RandString(6, "hex")),
					RouteTarget: pointer.To(testutils.RandomRouteTarget(t)),
					InterconnectSecurityZones: map[string]apstra.InterconnectSecurityZone{
						bogusRZID: {
							RouteTarget:     pointer.To(testutils.RandomRouteTarget(t)),
							RoutingPolicyID: pointer.To(rpID),
						},
					},
				})
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ace.Type(), apstra.ErrNotfound)
				require.NotNil(t, ace.Detail())
				require.IsType(t, apstra.ErrNotFoundDetail{}, ace.Detail())
				require.Equal(t, apstra.NodeTypeSecurityZone, ace.Detail().(apstra.ErrNotFoundDetail).Type)
				require.Equal(t, bogusRZID, ace.Detail().(apstra.ErrNotFoundDetail).ID)
			})

			// Test Create with default Routing Zone ID.
			t.Run("create_with_default_routing_zone", func(t *testing.T) {
				ctx := testutils.ContextWithTestID(ctx, t)

				_, err = bp.CreateEVPNInterconnectGroup(ctx, apstra.EVPNInterconnectGroup{
					Label:       pointer.To("a" + testutils.RandString(6, "hex")),
					RouteTarget: pointer.To(testutils.RandomRouteTarget(t)),
					InterconnectSecurityZones: map[string]apstra.InterconnectSecurityZone{
						*defaultRZ.ID(): {
							RouteTarget:     pointer.To(testutils.RandomRouteTarget(t)),
							RoutingPolicyID: pointer.To(rpID),
						},
					},
				})
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ace.Type(), apstra.ErrUnknown)
				require.Nil(t, ace.Detail())
				require.Contains(t, err.Error(), "not supported")
			})

			// Test Create with bogus Routing Policy ID.
			t.Run("create_with_bogus_routing_policy", func(t *testing.T) {
				ctx := testutils.ContextWithTestID(ctx, t)

				_, err = bp.CreateEVPNInterconnectGroup(ctx, apstra.EVPNInterconnectGroup{
					Label:       pointer.To("a" + testutils.RandString(6, "hex")),
					RouteTarget: pointer.To(testutils.RandomRouteTarget(t)),
					InterconnectSecurityZones: map[string]apstra.InterconnectSecurityZone{
						rzID: {
							RouteTarget:     pointer.To(testutils.RandomRouteTarget(t)),
							RoutingPolicyID: pointer.To(bogusRPID),
						},
					},
				})
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ace.Type(), apstra.ErrNotfound)
				require.NotNil(t, ace.Detail())
				require.IsType(t, apstra.ErrNotFoundDetail{}, ace.Detail())
				require.Equal(t, apstra.NodeTypeRoutingPolicy, ace.Detail().(apstra.ErrNotFoundDetail).Type)
				require.Equal(t, bogusRPID, ace.Detail().(apstra.ErrNotFoundDetail).ID)
			})

			// Create a valid Interconnect Group for use in Update tests
			dciID, err := bp.CreateEVPNInterconnectGroup(ctx, apstra.EVPNInterconnectGroup{
				Label:       pointer.To("a" + testutils.RandString(6, "hex")),
				RouteTarget: pointer.To(testutils.RandomRouteTarget(t)),
			})
			require.NoError(t, err)
			//dci, err := bp.GetEVPNInterconnectGroup(ctx, dciID)
			//require.NoError(t, err)

			// Test Update with bogus Virtual Network ID.
			t.Run("update_with_bogus_virtual_network", func(t *testing.T) {
				ctx := testutils.ContextWithTestID(ctx, t)

				var u apstra.EVPNInterconnectGroup
				_ = u.SetID(dciID)
				u.InterconnectVirtualNetworks = map[string]apstra.InterconnectVirtualNetwork{
					bogusVNID: {
						L2Enabled:      true,
						L3Enabled:      true,
						TranslationVNI: pointer.To(uint32(rand.Intn(4000) + 1)),
					},
				}
				err = bp.UpdateEVPNInterconnectGroup(ctx, u)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ace.Type(), apstra.ErrNotfound)
				require.NotNil(t, ace.Detail())
				require.IsType(t, apstra.ErrNotFoundDetail{}, ace.Detail())
				require.Equal(t, apstra.NodeTypeVirtualNetwork, ace.Detail().(apstra.ErrNotFoundDetail).Type)
				require.Equal(t, bogusVNID, ace.Detail().(apstra.ErrNotFoundDetail).ID)
			})

			// Test Update with bogus Routing Zone ID.
			t.Run("update_with_bogus_routing_zone", func(t *testing.T) {
				ctx := testutils.ContextWithTestID(ctx, t)

				var u apstra.EVPNInterconnectGroup
				_ = u.SetID(dciID)
				u.InterconnectSecurityZones = map[string]apstra.InterconnectSecurityZone{
					bogusRZID: {
						RouteTarget:     pointer.To(testutils.RandomRouteTarget(t)),
						RoutingPolicyID: pointer.To(rpID),
					},
				}
				err = bp.UpdateEVPNInterconnectGroup(ctx, u)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ace.Type(), apstra.ErrNotfound)
				require.NotNil(t, ace.Detail())
				require.IsType(t, apstra.ErrNotFoundDetail{}, ace.Detail())
				require.Equal(t, apstra.NodeTypeSecurityZone, ace.Detail().(apstra.ErrNotFoundDetail).Type)
				require.Equal(t, bogusRZID, ace.Detail().(apstra.ErrNotFoundDetail).ID)
			})

			// Test Update with default Routing Zone ID.
			t.Run("update_with_default_routing_zone", func(t *testing.T) {
				ctx := testutils.ContextWithTestID(ctx, t)

				var u apstra.EVPNInterconnectGroup
				_ = u.SetID(dciID)
				u.InterconnectSecurityZones = map[string]apstra.InterconnectSecurityZone{
					*defaultRZ.ID(): {
						RouteTarget:     pointer.To(testutils.RandomRouteTarget(t)),
						RoutingPolicyID: pointer.To(rpID),
					},
				}

				err = bp.UpdateEVPNInterconnectGroup(ctx, u)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ace.Type(), apstra.ErrUnknown)
				require.Nil(t, ace.Detail())
				require.Contains(t, err.Error(), "not supported")
			})

			// Test Update with bogus Routing Policy ID.
			t.Run("update_with_bogus_routing_policy", func(t *testing.T) {
				ctx := testutils.ContextWithTestID(ctx, t)

				var u apstra.EVPNInterconnectGroup
				_ = u.SetID(dciID)
				u.InterconnectSecurityZones = map[string]apstra.InterconnectSecurityZone{
					rzID: {
						RouteTarget:     pointer.To(testutils.RandomRouteTarget(t)),
						RoutingPolicyID: pointer.To(bogusRPID),
					},
				}
				err = bp.UpdateEVPNInterconnectGroup(ctx, u)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ace.Type(), apstra.ErrNotfound)
				require.NotNil(t, ace.Detail())
				require.IsType(t, apstra.ErrNotFoundDetail{}, ace.Detail())
				require.Equal(t, apstra.NodeTypeRoutingPolicy, ace.Detail().(apstra.ErrNotFoundDetail).Type)
				require.Equal(t, bogusRPID, ace.Detail().(apstra.ErrNotFoundDetail).ID)
			})
		})
	}
}

func TestEvpnInterconnectGroupClearTranslationVNI(t *testing.T) {
	ctx := context.Background()

	clients := testclient.GetTestClients(t, ctx)
	for _, c := range clients {
		t.Run(c.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bp := dctestobj.TestBlueprintA(t, ctx, c.Client)
			rzID := dctestobj.TestSecurityZoneA(t, ctx, bp)
			vnID := dctestobj.TestVirtualNetworkA(t, ctx, bp, rzID)
			rpID := dctestobj.TestRoutingPolicyA(t, ctx, bp)

			// Create the DCI with a translation VNI.
			dciID, err := bp.CreateEVPNInterconnectGroup(ctx, apstra.EVPNInterconnectGroup{
				Label:       pointer.To(testutils.RandString(6, "hex")),
				RouteTarget: pointer.To(testutils.RandomRouteTarget(t)),
				InterconnectSecurityZones: map[string]apstra.InterconnectSecurityZone{
					rzID: {
						RouteTarget:     pointer.To(testutils.RandomRouteTarget(t)),
						RoutingPolicyID: pointer.To(rpID),
					},
				},
				InterconnectVirtualNetworks: map[string]apstra.InterconnectVirtualNetwork{
					vnID: {
						L2Enabled:      true,
						L3Enabled:      true,
						TranslationVNI: pointer.To(uint32(rand.Intn(4000) + 1)),
					},
				},
			})
			require.NoError(t, err)

			// Clear the translation VNI
			var update apstra.EVPNInterconnectGroup
			_ = update.SetID(dciID)
			update.InterconnectVirtualNetworks = map[string]apstra.InterconnectVirtualNetwork{
				vnID: {
					L2Enabled:      true,
					L3Enabled:      true,
					TranslationVNI: pointer.To(uint32(0)), // zero should clear the translation VNI
				},
			}
			err = bp.UpdateEVPNInterconnectGroup(ctx, update)
			require.NoError(t, err)

			// Read back the DCI and verify that the translation VNI is cleared
			dci, err := bp.GetEVPNInterconnectGroup(ctx, dciID)
			require.NoError(t, err)
			require.Contains(t, dci.InterconnectVirtualNetworks, vnID)
			require.Nil(t, dci.InterconnectVirtualNetworks[vnID].TranslationVNI)
		})
	}
}
