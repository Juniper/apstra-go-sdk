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

func AggregateLink(t testing.TB, req, resp apstra.FreeformAggregateLink, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Aggregate Link")

	if req.ID() != nil {
		require.NotNil(t, resp.ID(), msg)
		require.Equal(t, *req.ID(), *resp.ID(), msg)
	}

	require.Equal(t, req.Label, resp.Label, msg)
	require.ElementsMatch(t, req.MemberLinkIds, resp.MemberLinkIds, msg)
	AggregateLinkEndpointGroup(t, req.EndpointGroups[0], resp.EndpointGroups[0], testmessage.Add(msg, "Comparing Aggregate Link Endpoint Group 0")...)
	AggregateLinkEndpointGroup(t, req.EndpointGroups[1], resp.EndpointGroups[1], testmessage.Add(msg, "Comparing Aggregate Link Endpoint Group 1")...)

	require.ElementsMatch(t, req.Tags, resp.Tags, msg)
}
