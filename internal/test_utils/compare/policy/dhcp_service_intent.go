// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparepolicy

import (
	"testing"

	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/Juniper/apstra-go-sdk/policy"
	"github.com/stretchr/testify/require"
)

func DHCPServiceIntent(t testing.TB, req, resp policy.DHCPServiceIntent, msg ...string) {
	msg = testmessage.Add(msg, "Comparing DHCP Service Intent Policy")

	require.Equal(t, req.Active, resp.Active, msg)
}
