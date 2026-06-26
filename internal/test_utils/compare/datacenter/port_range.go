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

func PortRange(t testing.TB, req, resp datacenter.PortRange, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Port Range")

	require.Equal(t, req.First, resp.First, msg)
	require.Equal(t, req.Last, resp.Last, msg)
}
