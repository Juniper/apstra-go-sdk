package apstra

import (
	"context"
	"math/rand"
	"os"
	"sort"
	"sync"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/stretchr/testify/require"
)

func TestNextInterface(t *testing.T) {
	type testCase struct {
		t string
		e string
	}

	testCases := []testCase{
		{t: "xe-0/0/0", e: "xe-0/0/1"},
		{t: "xe-0/0/9", e: "xe-0/0/10"},
	}

	for i, tc := range testCases {
		r := nextInterface(tc.t)
		if tc.e != r {
			t.Fatalf("test case %d: expected %s got %s", i, tc.e, r)
		}
	}
}

func testRaResourceIpv4(ctx context.Context, t testing.TB, cidrBlock string, bits int, client *FreeformClient) ObjectId {
	prefix := randomPrefix(t, cidrBlock, bits)
	id, err := client.CreateRaResource(ctx, &FreeformRaResourceData{
		ResourceType: enum.FFResourceTypeIpv4,
		Label:        randString(6, "hex"),
		Value:        toPtr(prefix.String()),
		GroupId:      testResourceGroup(ctx, t, client),
	})
	require.NoError(t, err)

	return id
}

func testRaResourceIpv6(ctx context.Context, t testing.TB, cidrBlock string, bits int, client *FreeformClient) ObjectId {
	prefix := randomPrefix(t, cidrBlock, bits)
	id, err := client.CreateRaResource(ctx, &FreeformRaResourceData{
		ResourceType: enum.FFResourceTypeIpv6,
		Label:        randString(6, "hex"),
		Value:        toPtr(prefix.String()),
		GroupId:      testResourceGroup(ctx, t, client),
	})
	require.NoError(t, err)

	return id
}

func testRaLocalVlanPool(ctx context.Context, t testing.TB, client *FreeformClient, ownerSystemId ObjectId, label string) ObjectId {
	ranges := rand.Intn(4) + 1
	ints, err := getRandInts(vlanMin+10, vlanMax, ranges*2)
	require.NoError(t, err)
	sort.Ints(ints)

	chunks := make([]FFLocalIntPoolChunk, ranges)
	for i := range ranges {
		chunks[i] = FFLocalIntPoolChunk{
			Start: ints[i*2],
			End:   ints[(i*2)+1],
		}
	}

	id, err := client.CreateRaLocalIntPool(ctx, &FreeformRaLocalIntPoolData{
		ResourceType: enum.FFResourceTypeVlan,
		Label:        label,
		OwnerId:      ownerSystemId,
		Chunks:       chunks,
	})
	require.NoError(t, err)

	return id
}

func TestSamples(t *testing.T) {
	type testCase struct {
		env      *string
		count    *int
		length   int
		expected int
	}

	initialEnvVal, ok := os.LookupEnv(envSampleSize)
	if ok {
		require.NoError(t, os.Unsetenv(envSampleSize))
	}

	testCases := map[string]testCase{
		"simple": {
			length:   5,
			expected: 5,
		},
		"env_valid": {
			length:   5,
			expected: 2,
			env:      toPtr("2"),
		},
		"count_wins": {
			length:   5,
			expected: 2,
			env:      toPtr("1"),
			count:    toPtr(2),
		},
		"env_over": {
			length:   5,
			expected: 5,
			env:      toPtr("10"),
		},
		"count_wins_over": {
			length:   5,
			expected: 5,
			env:      toPtr("1"),
			count:    toPtr(10),
		},
		"count_wins_both_over": {
			length:   5,
			expected: 5,
			env:      toPtr("9"),
			count:    toPtr(10),
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(testCases))
	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			// don't use t.Parallel due to risk of screwing up the environment

			if tCase.env != nil {
				require.NoError(t, os.Setenv(envSampleSize, *tCase.env))
			}

			var count []int
			if tCase.count != nil {
				count = []int{*tCase.count}
			}

			result := samples(t, tCase.length, count...)
			wg.Done()

			require.Equalf(t, tCase.expected, len(result), "expected %d samples, got %d", tCase.expected, len(result))
			for _, sample := range result {
				require.GreaterOrEqual(t, sample, 0)
				require.LessOrEqual(t, sample, tCase.length)
			}
		})
	}

	wg.Wait()
	if ok {
		// reset the environment
		require.NoError(t, os.Setenv(envSampleSize, initialEnvVal))
	}
}
