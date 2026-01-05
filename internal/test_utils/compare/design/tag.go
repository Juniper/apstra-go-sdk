// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedesign

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func Tag(t testing.TB, req, resp design.Tag, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Tag")

	if req.ID() != nil && resp.ID() != nil {
		require.Equal(t, *req.ID(), *resp.ID(), msg)
	}

	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, req.Description, resp.Description, msg)
}
