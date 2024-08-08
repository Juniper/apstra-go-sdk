package apstra

import (
	"context"
	"math/rand"
	"sort"
	"testing"

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
		ResourceType: FFResourceTypeIpv4,
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
		ResourceType: FFResourceTypeIpv6,
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
		ResourceType: FFResourceTypeVlan,
		Label:        label,
		OwnerId:      ownerSystemId,
		Chunks:       chunks,
	})
	require.NoError(t, err)

	return id
}
