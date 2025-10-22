// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparepolicy

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/policy"
	"github.com/stretchr/testify/require"
)

func VirtualNetwork(t testing.TB, req, resp policy.VirtualNetwork, msg ...string) {
	msg = addMsg(msg, "Comparing Virtual Network Policy")

	require.Equal(t, req.OverlayControlProtocol, resp.OverlayControlProtocol, msg)
}
