// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparefreeform

import (
	"regexp"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

// matches AOS-created labels like: server_name<->leaf_a,leaf_b[1]
var aggregateLinkLabelRegex = regexp.MustCompile(`^.+<->.+_\[\d+\]$`)

func AggregateLink(t testing.TB, req, resp apstra.FreeformAggregateLink, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Aggregate Link")

	if req.ID() != nil {
		require.NotNil(t, resp.ID(), msg)
		require.Equal(t, *req.ID(), *resp.ID(), msg)
	}

	switch {
	case req.Label == nil:
		// expect default label
		require.NotNil(t, resp.Label, msg)
		require.Regexp(t, aggregateLinkLabelRegex, *resp.Label, msg)
	case req.Label != nil && *req.Label == "":
		// expect no label
		require.Nil(t, resp.Label, msg)
	case req.Label != nil && *req.Label != "":
		// expect provided label
		require.NotNil(t, resp.Label, msg)
		require.Equal(t, *req.Label, *resp.Label, msg)
	}

	require.ElementsMatch(t, req.MemberLinkIds, resp.MemberLinkIds, msg)
	AggregateLinkEndpointGroup(t, req.EndpointGroups[0], resp.EndpointGroups[0], testmessage.Add(msg, "Comparing Aggregate Link Endpoint Group 0")...)
	AggregateLinkEndpointGroup(t, req.EndpointGroups[1], resp.EndpointGroups[1], testmessage.Add(msg, "Comparing Aggregate Link Endpoint Group 1")...)

	require.ElementsMatch(t, req.Tags, resp.Tags, msg)
}
