// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedatacenter

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func InterconnectVirtualNetwork(t testing.TB, req, resp apstra.InterconnectVirtualNetwork, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Interconnect Virtual Network")

	require.Equal(t, req.L2Enabled, resp.L2Enabled, msg)
	require.Equal(t, req.L3Enabled, resp.L3Enabled, msg)
	if req.TranslationVNI == nil {
		require.Nil(t, resp.TranslationVNI, msg)
	} else {
		require.NotNil(t, resp.TranslationVNI, msg)
		require.Equal(t, *req.TranslationVNI, *resp.TranslationVNI, msg)
	}
}
