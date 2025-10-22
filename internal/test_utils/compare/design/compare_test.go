// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedesign

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddMsg(t *testing.T) {
	type testCase struct {
		old      []string
		msg      string
		args     []any
		expected string
	}

	testCases := map[string]testCase{
		"empty": {},
		"add_string_to_empty": {
			msg:      "test",
			expected: "test",
		},
		"add_string_with_args_to_empty": {
			msg:      "test %d %s",
			expected: "test 1 two",
			args:     []any{1, "two"},
		},
		"add string with args to string": {
			old:      []string{"test one two"},
			msg:      "%s",
			args:     []any{"three"},
			expected: "test one two: three",
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result := addMsg(tCase.old, tCase.msg, tCase.args...)
			require.Equal(t, 1, len(result))
			require.Equal(t, tCase.expected, result[0])
		})
	}
}
