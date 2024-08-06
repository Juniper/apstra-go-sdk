//go:build integration
// +build integration

package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDResourceGenerators(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, expected, actual *FreeformResourceGeneratorData) {
		t.Helper()

		require.NotNil(t, expected)
		require.NotNil(t, actual)

		if expected.Label != "" {
			require.Equal(t, expected.Label, actual.Label)
		}
		if expected.Scope != "" {
			require.Equal(t, expected.Scope, actual.Scope)
		}
		if expected.ResourceType.String() != "" {
			require.Equal(t, expected.ResourceType, actual.ResourceType)
		}
		if expected.ContainerId != "" {
			require.Equal(t, expected.ContainerId, actual.ContainerId)
		}
		if expected.AllocatedFrom.String() != "" {
			require.Equal(t, expected.AllocatedFrom, actual.AllocatedFrom)
		}
		if expected.SubnetPrefixLen != nil {
			require.Equal(t, expected.SubnetPrefixLen, actual.SubnetPrefixLen)
		}
		if expected.ScopeNodePoolLabel != nil {
			require.Equal(t, *expected.ScopeNodePoolLabel, *actual.ScopeNodePoolLabel)
		}
	}

	type testCase struct {
		steps []FreeformResourceGeneratorData
	}

	for _, client := range clients {

		// create a blueprint with 4 systems
		ffc, _ := testFFBlueprintB(ctx, t, client.client, 4)

		testCases := map[string]testCase{
			"asn": {
				steps: []FreeformResourceGeneratorData{
					{
						Label:         randString(6, "hex"),
						Scope:         "node('system', name='target')",
						ContainerId:   testResourceGroup(ctx, t, ffc),
						AllocatedFrom: toPtr(testResourceGroupAsnData(ctx, t, ffc)),
						ResourceType:  FFResourceTypeAsn,
					},
					{
						Label:         randString(6, "hex"),
						Scope:         "node('system', name='target')",
						AllocatedFrom: toPtr(testResourceGroupAsnData(ctx, t, ffc)),
						ResourceType:  FFResourceTypeAsn,
					},
				},
			},
			"int": {
				steps: []FreeformResourceGeneratorData{
					{
						Label:         randString(6, "hex"),
						Scope:         "node('system', name='target')",
						ContainerId:   testResourceGroup(ctx, t, ffc),
						AllocatedFrom: toPtr(testResourceGroupIntData(ctx, t, ffc)),
						ResourceType:  FFResourceTypeInt,
					},
					{
						Label:         randString(6, "hex"),
						Scope:         "node('system', name='target')",
						AllocatedFrom: toPtr(testResourceGroupIntData(ctx, t, ffc)),
					},
				},
			},
			"ipv4": {
				steps: []FreeformResourceGeneratorData{
					{
						Label:           randString(6, "hex"),
						Scope:           "node('system', name='target')",
						ContainerId:     testResourceGroup(ctx, t, ffc),
						AllocatedFrom:   toPtr(testResourceGroupIpv4Data(ctx, t, ffc)),
						ResourceType:    FFResourceTypeIpv4,
						SubnetPrefixLen: toPtr(30),
					},
					{
						Label:           randString(6, "hex"),
						Scope:           "node('system', name='target')",
						AllocatedFrom:   toPtr(testResourceGroupIpv4Data(ctx, t, ffc)),
						SubnetPrefixLen: toPtr(29),
					},
				},
			},
			"ipv6": {
				steps: []FreeformResourceGeneratorData{
					{
						Label:           randString(6, "hex"),
						Scope:           "node('system', name='target')",
						ContainerId:     testResourceGroup(ctx, t, ffc),
						AllocatedFrom:   toPtr(testResourceGroupIpv6Data(ctx, t, ffc)),
						ResourceType:    FFResourceTypeIpv6,
						SubnetPrefixLen: toPtr(125),
					},
					{
						Label:           randString(6, "hex"),
						Scope:           "node('system', name='target')",
						AllocatedFrom:   toPtr(testResourceGroupIpv6Data(ctx, t, ffc)),
						SubnetPrefixLen: toPtr(126),
					},
				},
			},
		}

		for tName, tCase := range testCases {
			t.Run(tName, func(t *testing.T) {
				var id ObjectId
				for i, step := range tCase.steps {
					if i == 0 {
						id, err = ffc.CreateResourceGenerator(ctx, &step)
						require.NoError(t, err)
					} else {
						err = ffc.UpdateResourceGenerator(ctx, id, &step)
						require.NoError(t, err)
					}
					fetched, err := ffc.GetResourceGenerator(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, fetched.Id)
					compare(t, &step, fetched.Data)

					fetched, err = ffc.GetResourceGeneratorByName(ctx, step.Label)
					require.NoError(t, err)
					require.Equal(t, id, fetched.Id)
					compare(t, &step, fetched.Data)
				}

				raGroups, err := ffc.GetAllResourceGenerators(ctx)
				require.NoError(t, err)
				ids := make([]ObjectId, len(raGroups))
				for i, template := range raGroups {
					ids[i] = template.Id
				}
				require.Contains(t, ids, id)

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
			})
		}

	}
}
