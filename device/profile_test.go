// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProfile_MarshalJSON(t *testing.T) {
	type testCase struct {
		v Profile
		e string
	}

	testCases := map[string]testCase{
		"generic_1x10": {
			v: testProfileGeneric1x10,
			e: testProfileGeneric1x10JSON,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			r, err := json.Marshal(tCase.v)
			require.NoError(t, err)

			// get rid of extraneous fields in the expected string value
			eMap := map[string]json.RawMessage{}
			require.NoError(t, json.Unmarshal([]byte(tCase.e), &eMap))
			delete(eMap, "created_at")
			delete(eMap, "last_modified_at")
			e, err := json.Marshal(eMap)
			require.NoError(t, err)

			require.JSONEq(t, string(e), string(r))

			/*
			   inspect raw json with this
			   pbpaste | jq 'walk(if type == "object" then . | to_entries | sort_by(.key) | from_entries else . end)'
			*/
		})
	}
}

func TestProfile_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		e Profile
		v string
	}

	testCases := map[string]testCase{
		"generic_1x10": {
			v: testProfileGeneric1x10JSON,
			e: testProfileGeneric1x10,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			var r Profile
			err := json.Unmarshal([]byte(tCase.v), &r)
			require.NoError(t, err)

			require.Equal(t, tCase.e, r)
		})
	}
}
