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

func Spine(t testing.TB, req, resp design.Spine, msg ...string) {
	msg = addMsg(msg, "Comparing Spine")

	require.Equal(t, req.Count, resp.Count, msg)
	require.Equal(t, req.LinkPerSuperspineCount, resp.LinkPerSuperspineCount, msg)
	require.Equal(t, req.LinkPerSuperspineSpeed, resp.LinkPerSuperspineSpeed, msg)
	LogicalDevice(t, req.LogicalDevice, req.LogicalDevice, msg...)
	require.Equal(t, len(req.Tags), len(resp.Tags), msg)
	for i := range len(req.Tags) {
		Tag(t, req.Tags[i], resp.Tags[i], addMsg(msg, "Comparing Tag %d", i)...)
	}
}
