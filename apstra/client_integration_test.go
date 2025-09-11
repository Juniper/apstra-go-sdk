// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestClientLog(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()

			client.Client.Logf(1, "log test - client %q", client.Name())
		})
	}
}

func TestLoginEmptyPassword(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			if client.Type() == testclient.ClientTypeAPIOps {
				t.Skipf("skipping test - api-ops type clients do not log in or out")
			}

			c := *client.Client // don't use iterator variable because it points to the shared client object
			c.SetPassword("")
			err := c.Login(ctx)
			require.Error(t, err)
		})
	}
}

func TestLoginBadPassword(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			if client.Type() == testclient.ClientTypeAPIOps {
				t.Skipf("skipping test - api-ops type clients do not log in or out")
			}

			c := *client.Client
			c.SetPassword(testutils.RandString(10, "hex"))

			err := c.Login(ctx)
			require.Error(t, err)
		})
	}
}

func TestLogoutAuthFail(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			if client.Type() == testclient.ClientTypeAPIOps {
				t.Skipf("skipping test - api-ops type clients do not log in or out")
			}

			c := *client.Client

			err := c.Login(ctx)
			require.NoError(t, err)

			client.Client.SetAuthtoken(testutils.RandJWT())
			err = c.Logout(ctx)
			require.Error(t, err)
		})
	}
}

func TestGetBlueprintOverlayControlProtocol(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		bpFunc      func(testing.TB, context.Context, *apstra.Client) *apstra.TwoStageL3ClosClient
		expectedOcp apstra.OverlayControlProtocol
	}

	testCases := []testCase{
		{bpFunc: dctestobj.TestBlueprintA, expectedOcp: apstra.OverlayControlProtocolEvpn},
		{bpFunc: dctestobj.TestBlueprintB, expectedOcp: apstra.OverlayControlProtocolNone},
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			for i := range testCases {
				i := i
				t.Run(fmt.Sprintf("test_case_%d", i), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					bpClient := testCases[i].bpFunc(t, ctx, client.Client)

					ocp, err := client.Client.BlueprintOverlayControlProtocol(ctx, bpClient.Id())
					require.NoError(t, err)

					if ocp != testCases[i].expectedOcp {
						t.Fatalf("expected overlay control protocol %q, got %q", testCases[i].expectedOcp.String(), ocp.String())
					}
					log.Printf("blueprint %q has overlay control protocol %q", bpClient.Id(), ocp.String())
				})
			}
		})
	}
}

func TestCRUDIntegerPools(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	validate := func(req *apstra.IntPoolRequest, resp *apstra.IntPool) {
		require.Equal(t, req.DisplayName, resp.DisplayName)
		require.Equal(t, len(req.Ranges), len(resp.Ranges))

		for i := range req.Ranges {
			reqFirst := req.Ranges[i].(*apstra.IntRangeRequest).First
			reqLast := req.Ranges[i].(*apstra.IntRangeRequest).Last
			respFirst := resp.Ranges[i].First
			respLast := resp.Ranges[i].Last

			require.Equal(t, reqFirst, respFirst)
			require.Equal(t, reqLast, respLast)
		}

		require.Equal(t, len(req.Tags), len(resp.Tags))
		compare.SlicesAsSets(t, req.Tags, resp.Tags, "tags mismatch")
	}

	randomTags := func(min, max int) []string {
		var result []string
		for i := 0; i < rand.Intn(max-min)+min; i++ {
			result = append(result, testutils.RandString(5, "hex"))
		}
		return result
	}

	randomRanges := func(minRanges, maxRanges int, minVal, maxVal uint32) []apstra.IntfIntRange {
		rangeCount := rand.Intn(maxRanges-minRanges) + minRanges
		valMap := make(map[int]struct{})
		for len(valMap) < rangeCount*2 {
			valMap[rand.Intn(int(maxVal-minVal))+int(minVal)] = struct{}{}
		}
		valSlice := make([]int, len(valMap))
		var i int
		for k := range valMap {
			valSlice[i] = k
			i++
		}
		sort.Ints(valSlice)

		result := make([]apstra.IntfIntRange, rangeCount)
		for i = 0; i < rangeCount; i++ {
			result[i] = &apstra.IntRangeRequest{
				First: uint32(valSlice[(i * 2)]),
				Last:  uint32(valSlice[(i*2)+1]),
			}
		}
		return result
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			pools, err := client.Client.GetIntegerPools(ctx)
			require.NoError(t, err)

			beforePoolCount := len(pools)
			request := apstra.IntPoolRequest{
				DisplayName: testutils.RandString(5, "hex"),
				Ranges:      randomRanges(2, 5, 10000, 10999),
				Tags:        randomTags(2, 5),
			}

			id, err := client.Client.CreateIntegerPool(ctx, &request)
			require.NoError(t, err)

			pool, err := client.Client.GetIntegerPool(ctx, id)
			require.NoError(t, err)

			validate(&request, pool)

			pools, err = client.Client.GetIntegerPools(ctx)
			require.NoError(t, err)
			require.Equal(t, beforePoolCount+1, len(pools))

			poolIdx := -1
			for i, p := range pools {
				if p.Id == id {
					poolIdx = i
					break
				}
			}
			require.GreaterOrEqual(t, poolIdx, 0, "just-created pool id not found among pools")

			validate(&request, &pools[poolIdx])

			poolIds, err := client.Client.ListIntegerPoolIds(ctx)
			require.NoError(t, err)
			require.Equal(t, beforePoolCount+1, len(poolIds))

			var found bool
			for _, poolId := range poolIds {
				if poolId == id {
					found = true
					break
				}
			}
			require.True(t, found, "newly created pool ID not found among pool ID list")

			request = apstra.IntPoolRequest{
				DisplayName: testutils.RandString(5, "hex"),
				Ranges:      randomRanges(2, 5, 1, math.MaxUint32),
				Tags:        randomTags(2, 5),
			}

			err = client.Client.UpdateIntegerPool(ctx, id, &request)
			require.NoError(t, err)

			pool, err = client.Client.GetIntegerPool(ctx, id)
			require.NoError(t, err)

			validate(&request, pool)

			err = client.Client.DeleteIntegerPool(ctx, id)
			require.NoError(t, err)

			pools, err = client.Client.GetIntegerPools(ctx)
			require.NoError(t, err)

			for i := len(pools) - 1; i >= 0; i-- {
				if pools[i].Status == apstra.PoolStatusDeleting {
					log.Printf("dropping pool %s from fetched pool list because it has status %s", pools[i].Id, pools[i].Status.String())
					pools[i] = pools[len(pools)-1]
					pools = pools[:len(pools)-1]
				}
			}

			if len(pools) != beforePoolCount {
				t.Fatalf("pools before creation: %d; after creation: %d", beforePoolCount, len(pools))
			}
		})
	}
}

func TestBlueprintOverlayControlProtocol(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		templateId apstra.ObjectId
		expected   apstra.OverlayControlProtocol
	}

	testCases := map[string]testCase{
		"L2_Virtual_EVPN": {
			templateId: "L2_Virtual_EVPN",
			expected:   apstra.OverlayControlProtocolEvpn,
		},
		"L2_Virtual": {
			templateId: "L2_Virtual",
			expected:   apstra.OverlayControlProtocolNone,
		},
	}

	createBlueprint := func(t testing.TB, templateId apstra.ObjectId, client *apstra.Client) apstra.ObjectId {
		t.Helper()

		id, err := client.CreateBlueprintFromTemplate(ctx, &apstra.CreateBlueprintFromTemplateRequest{
			RefDesign:  enum.RefDesignDatacenter,
			Label:      testutils.RandString(5, "hex"),
			TemplateId: templateId,
		})
		require.NoError(t, err)

		t.Cleanup(func() { require.NoError(t, client.DeleteBlueprint(ctx, id)) })

		return id
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					ctx := testutils.ContextWithTestID(ctx, t)

					bpId := createBlueprint(t, tCase.templateId, client.Client)

					ocp, err := client.Client.BlueprintOverlayControlProtocol(ctx, bpId)
					require.NoError(t, err)
					require.Equal(t, tCase.expected, ocp)
				})
			}
		})
	}
}
