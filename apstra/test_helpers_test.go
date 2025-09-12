// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	crand "crypto/rand"
	"math"
	"math/rand"
	"net"
	"os"
	"sort"
	"sync"
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
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
			length:   50,
			expected: 50,
		},
		"env_valid": {
			length:   50,
			expected: 12,
			env:      toPtr("12"),
		},
		"count_wins": {
			length:   50,
			expected: 13,
			env:      toPtr("11"),
			count:    toPtr(13),
		},
		"env_over": {
			length:   50,
			expected: 50,
			env:      toPtr("100"),
		},
		"count_wins_over": {
			length:   50,
			expected: 50,
			env:      toPtr("1"),
			count:    toPtr(100),
		},
		"both_over": {
			length:   50,
			expected: 50,
			env:      toPtr("101"),
			count:    toPtr(102),
		},
		"count_zero": {
			length:   50,
			expected: 50,
			count:    toPtr(0),
		},
		"env_zero": {
			length:   50,
			expected: 50,
			env:      toPtr("0"),
		},
		"both_zero": {
			length:   50,
			expected: 50,
			count:    toPtr(0),
			env:      toPtr("0"),
		},
		"count_wins_zero": {
			length:   50,
			expected: 50,
			count:    toPtr(0),
			env:      toPtr("23"),
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(testCases))
	for tName, tCase := range testCases {
		tName, tCase := tName, tCase
		t.Run(tName, func(t *testing.T) {
			if tCase.env != nil {
				t.Setenv(envSampleSize, *tCase.env)
			}

			var count []int
			if tCase.count != nil {
				count = []int{*tCase.count}
			}

			result := sampleIndexes(t, tCase.length, count...)
			wg.Done()

			require.Equalf(t, tCase.expected, len(result), "expected %d samples, got %d", tCase.expected, len(result))
			for _, sample := range result {
				require.GreaterOrEqual(t, sample, 0)
				require.LessOrEqual(t, sample, tCase.length)
			}
		})
	}

	// reset the environment after tests complete
	wg.Wait()
	if ok {
		require.NoError(t, os.Setenv(envSampleSize, initialEnvVal))
	} else {
		require.NoError(t, os.Unsetenv(envSampleSize))
	}
}

// Deprecated: Use testutils.RandomPrefix()
func randomPrefix(t testing.TB, cidrBlock string, bits int) net.IPNet {
	t.Helper()

	ip, block, err := net.ParseCIDR(cidrBlock)
	if err != nil {
		t.Fatalf("randomPrefix cannot parse cidrBlock - %s", err)
	}
	if block.IP.String() != ip.String() {
		t.Fatal("invocation of randomPrefix doesn't use a base block address")
	}

	mOnes, mBits := block.Mask.Size()
	if mOnes >= bits {
		t.Fatalf("cannot select a random /%d from within %s", bits, cidrBlock)
	}

	// generate a completely random address
	randomIP := make(net.IP, mBits/8)
	_, err = crand.Read(randomIP)
	if err != nil {
		t.Fatalf("rand read failed")
	}

	// mask off the "network" bits
	for i, b := range randomIP {
		mBitsThisByte := min(mOnes, 8)
		mOnes -= mBitsThisByte
		block.IP[i] = block.IP[i] | (b & byte(math.MaxUint8>>mBitsThisByte))
	}

	block.Mask = net.CIDRMask(bits, mBits)

	_, result, err := net.ParseCIDR(block.String())
	if err != nil {
		t.Fatal("failed to parse own CIDR block")
	}

	return *result
}
