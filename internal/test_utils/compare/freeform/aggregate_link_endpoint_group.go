// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparefreeform

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func AggregateLinkEndpointGroup(t testing.TB, req, resp apstra.FreeformAggregateLinkEndpointGroup, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Aggregate Link Endpoint Group")

	if req.ID() != nil {
		require.NotNil(t, resp.ID(), msg)
		require.Equal(t, *req.ID(), *resp.ID(), msg)
	}

	if req.Label != nil {
		if *req.Label == "" {
			require.Nil(t, resp.Label, msg) // an empty string in the request will clear API value if any
		} else {
			require.NotNil(t, resp.Label)
			require.Equal(t, *req.Label, *resp.Label, msg)
		}
	}

	require.ElementsMatch(t, req.Tags, resp.Tags, msg)

	require.Equal(t, len(req.Endpoints), len(resp.Endpoints), msg)

	reqEndpoints := make(map[string]apstra.FreeformAggregateLinkEndpoint, len(req.Endpoints))
	for _, ep := range req.Endpoints {
		require.NotContains(t, reqEndpoints, ep.SystemID, msg)
		reqEndpoints[ep.SystemID] = ep
	}

	respEndpoints := make(map[string]apstra.FreeformAggregateLinkEndpoint, len(resp.Endpoints))
	for _, ep := range resp.Endpoints {
		require.NotContains(t, respEndpoints, ep.SystemID, msg)
		respEndpoints[ep.SystemID] = ep
	}

	for _, reqEp := range req.Endpoints {
		respEp, ok := respEndpoints[reqEp.SystemID]
		require.True(t, ok, msg)

		AggregateLinkEndpoint(t, reqEp, respEp, msg...)
	}
}
