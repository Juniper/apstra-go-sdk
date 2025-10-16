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

func RackTypeWithCount(t testing.TB, req, resp design.RackTypeWithCount, msg ...string) {
	msg = addMsg(msg, "Comparing Rack Type With Count")

	require.Equal(t, req.Count, resp.Count, msg)
	RackType(t, req.RackType, resp.RackType, msg...)
}
