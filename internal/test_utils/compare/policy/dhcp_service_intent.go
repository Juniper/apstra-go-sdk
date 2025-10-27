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

func DHCPServiceIntent(t testing.TB, req, resp policy.DHCPServiceIntent, msg ...string) {
	msg = addMsg(msg, "Comparing DHCP Service Intent Policy")

	require.Equal(t, req.Active, resp.Active, msg)
}
