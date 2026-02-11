// Copyright (c) Juniper Networks, Inc., 2025-2026.
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

func AntiAffinity(t testing.TB, req, resp policy.AntiAffinity, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Anti Affinity Policy")

	require.Equal(t, req.MaxLinksPerPort, resp.MaxLinksPerPort, msg)
	require.Equal(t, req.MaxLinksPerSlot, resp.MaxLinksPerSlot, msg)
	require.Equal(t, req.MaxPerSystemLinksPerPort, resp.MaxPerSystemLinksPerPort, msg)
	require.Equal(t, req.MaxPerSystemLinksPerSlot, resp.MaxPerSystemLinksPerSlot, msg)
	require.Equal(t, req.Mode, resp.Mode, msg)
}
