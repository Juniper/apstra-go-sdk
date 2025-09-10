// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func TestCRUDResourceGenerators(t *testing.T) {
	ctx := wrapCtxWithTestId(t, context.Background())
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, expected, actual *FreeformResourceGeneratorData) {
		t.Helper()

		require.NotNil(t, expected)
		require.NotNil(t, actual)

		if expected.Label != "" {
			require.Equal(t, expected.Label, actual.Label, "label mismatch")
		}
		if expected.Scope != "" {
			require.Equal(t, expected.Scope, actual.Scope, "scope mismatch")
		}
		if expected.ResourceType.String() != "" {
			require.Equal(t, expected.ResourceType, actual.ResourceType, "resource type mismatch")
		}
		if expected.ContainerId != "" {
			require.Equal(t, expected.ContainerId, actual.ContainerId, "container id mismatch")
		}
		if expected.AllocatedFrom != nil {
			require.NotNil(t, actual.AllocatedFrom)
			require.Equal(t, expected.AllocatedFrom, actual.AllocatedFrom, "allocated from mismatch")
		}
		if expected.SubnetPrefixLen != nil {
			require.NotNil(t, actual.SubnetPrefixLen)
			require.Equal(t, expected.SubnetPrefixLen, actual.SubnetPrefixLen, "subnet prefix len mismatch")
		}
		if expected.ScopeNodePoolLabel != nil {
			require.NotNil(t, actual.ScopeNodePoolLabel)
			require.Equal(t, *expected.ScopeNodePoolLabel, *actual.ScopeNodePoolLabel, "scope node pool label mismatch")
		}
	}

	type testCase struct {
		steps []FreeformResourceGeneratorData
	}

	for clientName, client := range clients {
		t.Run(client.client.apiVersion.String()+"_"+client.clientType+"_"+clientName, func(t *testing.T) {
			t.Parallel()
			ctx := wrapCtxWithTestId(t, ctx)

			// create a blueprint with 4 systems
			systemCount := 4
			ffc, systemIds, _ := testFFBlueprintB(ctx, t, client.client, systemCount, 0)
			require.Equal(t, systemCount, len(systemIds))

			// attach a VLAN pool to each system
			type systemIdAndPoolLabel struct {
				systemId  ObjectId
				poolLabel string
			}

			systemsAndVlanPools := make([]systemIdAndPoolLabel, systemCount)
			for i, systemId := range systemIds {
				label := "a" + randString(6, "hex") // all numeric values cause problem
				_ = testRaLocalVlanPool(ctx, t, ffc, systemId, label)
				systemsAndVlanPools[i] = systemIdAndPoolLabel{
					systemId:  systemId,
					poolLabel: label,
				}
			}

			testCases := map[string]testCase{
				"asn": {
					steps: []FreeformResourceGeneratorData{
						{
							ResourceType:  enum.FFResourceTypeAsn,
							ContainerId:   testResourceGroup(ctx, t, ffc),
							AllocatedFrom: toPtr(testResourceGroupAsn(ctx, t, ffc)),
							Label:         randString(6, "hex"),
							Scope:         "node('system', name='target')",
						},
						{
							AllocatedFrom: toPtr(testResourceGroupAsn(ctx, t, ffc)),
							Label:         randString(6, "hex"),
							Scope:         fmt.Sprintf("node('system', name='target', id='%s')", systemIds[rand.IntN(systemCount)]),
						},
					},
				},
				"int": {
					steps: []FreeformResourceGeneratorData{
						{
							ResourceType:  enum.FFResourceTypeInt,
							ContainerId:   testResourceGroup(ctx, t, ffc),
							AllocatedFrom: toPtr(testResourceGroupInt(ctx, t, ffc)),
							Label:         randString(6, "hex"),
							Scope:         "node('system', name='target')",
						},
						{
							AllocatedFrom: toPtr(testResourceGroupInt(ctx, t, ffc)),
							Label:         randString(6, "hex"),
							Scope:         fmt.Sprintf("node('system', name='target', id='%s')", systemIds[rand.IntN(systemCount)]),
						},
					},
				},
				"ipv4": {
					steps: []FreeformResourceGeneratorData{
						{
							ResourceType:    enum.FFResourceTypeIpv4,
							ContainerId:     testResourceGroup(ctx, t, ffc),
							AllocatedFrom:   toPtr(testResourceGroupIpv4(ctx, t, ffc)),
							Label:           randString(6, "hex"),
							Scope:           "node('system', name='target')",
							SubnetPrefixLen: toPtr(30),
						},
						{
							AllocatedFrom:   toPtr(testResourceGroupIpv4(ctx, t, ffc)),
							Label:           randString(6, "hex"),
							Scope:           fmt.Sprintf("node('system', name='target', id='%s')", systemIds[rand.IntN(systemCount)]),
							SubnetPrefixLen: toPtr(29),
						},
					},
				},
				"ipv6": {
					steps: []FreeformResourceGeneratorData{
						{
							ResourceType:    enum.FFResourceTypeIpv6,
							ContainerId:     testResourceGroup(ctx, t, ffc),
							AllocatedFrom:   toPtr(testResourceGroupIpv6(ctx, t, ffc)),
							Label:           randString(6, "hex"),
							Scope:           "node('system', name='target')",
							SubnetPrefixLen: toPtr(125),
						},
						{
							AllocatedFrom:   toPtr(testResourceGroupIpv6(ctx, t, ffc)),
							Label:           randString(6, "hex"),
							Scope:           fmt.Sprintf("node('system', name='target', id='%s')", systemIds[rand.IntN(systemCount)]),
							SubnetPrefixLen: toPtr(126),
						},
					},
				},
				"vlan": {
					steps: []FreeformResourceGeneratorData{
						{
							ResourceType:       enum.FFResourceTypeVlan,
							ContainerId:        testResourceGroup(ctx, t, ffc),
							ScopeNodePoolLabel: toPtr(systemsAndVlanPools[0].poolLabel),
							Label:              randString(6, "hex"),
							Scope:              fmt.Sprintf("node(id='%s', name='target')", systemsAndVlanPools[0].systemId),
						},
						{
							Label:              randString(6, "hex"),
							Scope:              fmt.Sprintf("node(id='%s', name='target')", systemsAndVlanPools[0].systemId),
							ScopeNodePoolLabel: toPtr(systemsAndVlanPools[0].poolLabel),
						},
					},
				},
				"ipv4_host": {
					steps: []FreeformResourceGeneratorData{
						{
							ResourceType:  enum.FFResourceTypeHostIpv4,
							Label:         randString(6, "hex"),
							Scope:         "node('system', name='target')",
							AllocatedFrom: toPtr(testRaResourceIpv4(ctx, t, "10.1.0.0/16", 24, ffc)),
							ContainerId:   testResourceGroup(ctx, t, ffc),
						},
						{
							Label:         randString(6, "hex"),
							Scope:         fmt.Sprintf("node(id='%s', name='target')", systemIds[0]),
							AllocatedFrom: toPtr(testRaResourceIpv4(ctx, t, "10.2.0.0/16", 23, ffc)),
						},
					},
				},
				"ipv6_host": {
					steps: []FreeformResourceGeneratorData{
						{
							ResourceType:  enum.FFResourceTypeHostIpv6,
							Label:         randString(6, "hex"),
							Scope:         "node('system', name='target')",
							AllocatedFrom: toPtr(testRaResourceIpv6(ctx, t, "3fff:0:0::/48", 64, ffc)),
							ContainerId:   testResourceGroup(ctx, t, ffc),
						},
						{
							Label:         randString(6, "hex"),
							Scope:         fmt.Sprintf("node(id='%s', name='target')", systemIds[0]),
							AllocatedFrom: toPtr(testRaResourceIpv6(ctx, t, "3fff:0:1::/48", 64, ffc)),
						},
					},
				},
			}

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					t.Parallel()

					var err error
					var id ObjectId
					var resourceGenerator *FreeformResourceGenerator
					for i, step := range tCase.steps {
						if i == 0 {
							id, err = ffc.CreateResourceGenerator(ctx, &step)
							require.NoError(t, err)
						} else {
							err = ffc.UpdateResourceGenerator(ctx, id, &step)
							require.NoError(t, err)
						}

						resourceGenerator, err = ffc.GetResourceGenerator(ctx, id)
						require.NoError(t, err)
						require.Equal(t, id, resourceGenerator.Id)
						compare(t, &step, resourceGenerator.Data)

						resourceGenerator, err = ffc.GetResourceGeneratorByName(ctx, step.Label)
						require.NoError(t, err)
						require.Equal(t, id, resourceGenerator.Id)
						compare(t, &step, resourceGenerator.Data)

						raGroups, err := ffc.GetAllResourceGenerators(ctx)
						require.NoError(t, err)
						ids := make([]ObjectId, len(raGroups))
						for i, template := range raGroups {
							ids[i] = template.Id
						}
						require.Contains(t, ids, id)
					}

					err = ffc.DeleteResourceGenerator(ctx, id)
					require.NoError(t, err)

					var ace ClientErr

					_, err = ffc.GetResourceGenerator(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ErrNotfound, ace.Type())

					err = ffc.DeleteResourceGenerator(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ErrNotfound, ace.Type())

					_, err = ffc.GetResourceGeneratorByName(ctx, resourceGenerator.Data.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ErrNotfound, ace.Type())

					raGroups, err := ffc.GetAllResourceGenerators(ctx)
					require.NoError(t, err)
					ids := make([]ObjectId, len(raGroups))
					for i, template := range raGroups {
						ids[i] = template.Id
					}
					require.NotContains(t, ids, id)
				})
			}
		})
	}
}
