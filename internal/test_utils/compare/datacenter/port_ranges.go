// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedatacenter

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func PortRanges(t testing.TB, req, resp datacenter.PortRanges, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Port Ranges")

	require.Equal(t, len(req), len(resp), msg)
	for i := range req {
		PortRange(t, req[i], resp[i], testmessage.Add(msg, "Comparing Port Range at index %d", i)...)
	}
}
