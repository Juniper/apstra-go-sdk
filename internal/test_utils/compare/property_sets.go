// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package compare

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/stretchr/testify/require"
)

func PropertySetData(t testing.TB, a, b apstra.PropertySetData) {
	t.Helper()

	require.Equal(t, a.Label, b.Label)
	require.Equal(t, len(a.Blueprints), len(b.Blueprints))
	if len(a.Blueprints) > 0 {
		require.Equal(t, a.Blueprints, b.Blueprints)
	}
	JSON(t, a.Values, b.Values)
}
