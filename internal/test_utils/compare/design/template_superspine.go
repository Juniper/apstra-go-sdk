// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedesign

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/stretchr/testify/require"
)

func Superspine(t testing.TB, req, resp design.Superspine, msg ...string) {
	msg = addMsg(msg, "Comparing Superspine")

	require.Equal(t, req.PlaneCount, resp.PlaneCount, msg)
	require.Equal(t, req.SuperspinePerPlane, resp.SuperspinePerPlane, msg)
	LogicalDevice(t, req.LogicalDevice, req.LogicalDevice, msg...)
	require.Equal(t, len(req.Tags), len(resp.Tags), msg)
	for i := range len(req.Tags) {
		Tag(t, req.Tags[i], resp.Tags[i], addMsg(msg, "Comparing Tag %d", i)...)
	}
}
