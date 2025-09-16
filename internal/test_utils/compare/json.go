// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package compare

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func JSON(t testing.TB, m1, m2 json.RawMessage) {
	t.Helper()

	var map1 interface{}
	require.NoError(t, json.Unmarshal(m1, &map1))

	var map2 interface{}
	require.NoError(t, json.Unmarshal(m2, &map2))

	require.True(t, reflect.DeepEqual(map1, map2))
}
