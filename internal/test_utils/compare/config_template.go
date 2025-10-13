// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package compare

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/stretchr/testify/require"
)

func ConfigTemplate(t testing.TB, req, resp design.ConfigTemplate, msg ...string) {
	msg = addMsg(msg, "Comparing Config Template")

	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, req.Predefined, resp.Predefined, msg)
	require.Equal(t, req.Text, resp.Text, msg)
}
