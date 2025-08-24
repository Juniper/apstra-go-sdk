// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra_test

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/stretchr/testify/require"
)

func TestInterfaceSettingParam(t *testing.T) {
	type testCase struct {
		expected string
		val      apstra.InterfaceSettingParam
	}

	testCases := map[string]testCase{
		"a": {
			expected: `{\"global\":{\"breakout\":false,\"fpc\":0,\"pic\":0,\"port\":0,\"speed\":\"100g\"},\"interface\":{\"speed\":\"\"}}`,
			val: apstra.InterfaceSettingParam{
				Global: struct {
					Breakout bool   `json:"breakout"`
					Fpc      int    `json:"fpc"`
					Pic      int    `json:"pic"`
					Port     int    `json:"port"`
					Speed    string `json:"speed"`
				}{
					Breakout: false,
					Fpc:      0,
					Pic:      0,
					Port:     0,
					Speed:    "100g",
				},
				Interface: struct {
					Speed string `json:"speed"`
				}{},
			},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result := tCase.val.String()
			require.Equal(t, tCase.expected, result)
		})
	}
}
