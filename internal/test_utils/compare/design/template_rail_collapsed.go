// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedesign

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/design"
	comparepolicy "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/policy"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func TemplateRailCollapsed(t testing.TB, req, resp design.TemplateRailCollapsed, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Rail Collapsed Template")

	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, len(req.Racks), len(resp.Racks), msg)
	for i := range len(req.Racks) {
		RackTypeWithCount(t, req.Racks[i], resp.Racks[i], testmessage.Add(msg, "Comparing Rack %d", i)...)
	}

	comparepolicy.DHCPServiceIntent(t, req.DHCPServiceIntent, resp.DHCPServiceIntent, msg...)
	if req.VirtualNetworkPolicy != nil {
		require.NotNil(t, resp.VirtualNetworkPolicy)
		comparepolicy.VirtualNetwork(t, *req.VirtualNetworkPolicy, *resp.VirtualNetworkPolicy, msg...)
	}
	if req.ID() != nil && resp.ID() != nil {
		require.Equal(t, req.ID(), resp.ID(), msg)
	}
	if req.CreatedAt() != nil && resp.CreatedAt() != nil {
		require.Equal(t, req.CreatedAt(), resp.CreatedAt(), msg)
	}
	if req.LastModifiedAt() != nil && resp.LastModifiedAt() != nil {
		require.Equal(t, req.LastModifiedAt(), resp.LastModifiedAt(), msg)
	}
}
