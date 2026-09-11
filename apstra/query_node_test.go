// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra_test

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/stretchr/testify/require"
)

func TestNodeType_FromString(t *testing.T) {
	type testCase struct {
		data   string
		exp    apstra.NodeType
		expErr string
	}

	testCases := map[string]testCase{
		"valid_node_type_system": {
			data: "system",
			exp:  apstra.NodeTypeSystem,
		},
		"invalid_node_type_bogus": {
			data:   "bogus",
			expErr: `unknown node type "bogus"`,
		},
		"empty_string": {
			data: "",
			exp:  apstra.NodeTypeNone,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var nt apstra.NodeType
			err := nt.FromString(tCase.data)
			if tCase.expErr != "" {
				require.Error(t, err)
				require.EqualError(t, err, tCase.expErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tCase.exp, nt)
		})
	}
}

func TestNodeType_MarshalText(t *testing.T) {
	type testCase struct {
		data   apstra.NodeType
		exp    string
		expErr string
	}

	testCases := map[string]testCase{
		"valid_node_type_system": {
			data: apstra.NodeTypeSystem,
			exp:  "system",
		},
		"invalid_nodetype": {
			data:   apstra.NodeType(-1),
			expErr: "cannot marshal: unknown node type -1",
		},
		"empty_string": {
			data: apstra.NodeTypeNone,
			exp:  "",
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			text, err := tCase.data.MarshalText()
			if tCase.expErr != "" {
				require.Error(t, err)
				require.EqualError(t, err, tCase.expErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tCase.exp, string(text))
		})
	}
}

func TestNodeType_UnmarshalText(t *testing.T) {
	type testCase struct {
		data   string
		exp    apstra.NodeType
		expErr string
	}

	testCases := map[string]testCase{
		"valid_node_type_system": {
			data: "system",
			exp:  apstra.NodeTypeSystem,
		},
		"invalid_node_type_bogus": {
			data:   "bogus",
			expErr: `unknown node type "bogus"`,
		},
		"empty_string": {
			data: "",
			exp:  apstra.NodeTypeNone,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var nt apstra.NodeType
			err := nt.UnmarshalText([]byte(tCase.data))
			if tCase.expErr != "" {
				require.Error(t, err)
				require.EqualError(t, err, tCase.expErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tCase.exp, nt)
		})
	}
}
