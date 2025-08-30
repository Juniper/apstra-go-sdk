package compare

import (
	"fmt"
	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/stretchr/testify/require"
	"testing"
)

func IntPool(t testing.TB, req apstra.IntPoolRequest, data apstra.IntPool) {
	t.Helper()

	require.NotNil(t, req.DisplayName, data.DisplayName)
	SlicesAsSets(t, req.Tags, data.Tags, "tags mismatch")
	require.Equal(t, len(req.Ranges), len(data.Ranges))

	requestedRanges := make([]string, len(req.Ranges))
	for i, r := range req.Ranges {
		switch rt := r.(type) {
		case apstra.IntRangeRequest:
			requestedRanges[i] = fmt.Sprintf("%d-%d", rt.First, rt.Last)
		default:
			t.Fatalf("unhandled type %T", r)
		}
	}

	for _, r := range data.Ranges {
		require.Contains(t, requestedRanges, fmt.Sprintf("%d-%d", r.First, r.Last))
	}
}

func IpPool(t testing.TB, req apstra.NewIpPoolRequest, data apstra.IpPool) {
	t.Helper()

	require.Equal(t, req.DisplayName, data.DisplayName)
	SlicesAsSets(t, req.Tags, data.Tags, "tags mismatch")
	require.Equal(t, len(req.Subnets), len(data.Subnets))
	for i := range len(req.Subnets) {
		require.Contains(t, req.Subnets, data.Subnets[i].Network.String())
	}
}
