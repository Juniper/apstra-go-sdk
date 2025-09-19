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

func Tag(t testing.TB, a, b design.Tag) {
	t.Helper()

	if a.ID() != nil && b.ID() != nil {
		require.Equal(t, *a.ID(), *b.ID(), "IDs do not match")
	}

	require.Equal(t, a.Label, b.Label, "Labels do not match")
	require.Equal(t, a.Description, b.Description, "Descriptions do not match")
}
