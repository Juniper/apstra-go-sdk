// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/sha256"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRackTypeLink_replicate(t *testing.T) {
	type testCase struct {
		v RackTypeLink
	}

	testCases := map[string]testCase{
		"simple": {v: linkSimple},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			r := tCase.v.Replicate()
			require.Equal(t, mustHashForComparison(tCase.v, sha256.New()), mustHashForComparison(r, sha256.New()))
			require.Equal(t, tCase.v, r)
		})
	}
}

func TestRackTypeLink_MarshalJSON(t *testing.T) {
	type testCase struct {
		v RackTypeLink
		e string
	}

	testCases := map[string]testCase{
		"simple": {
			v: linkSimple,
			e: linkSimpleJSON,
		},
		"complicated": {
			v: linkComplicated,
			e: linkComplicatedJSON,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			r, err := json.Marshal(tCase.v)
			require.NoError(t, err)

			require.JSONEq(t, tCase.e, string(r))
		})
	}
}

//func TestRackType_UnmarshalJSON(t *testing.T) {
//	type testCase struct {
//		e RackType
//		v string
//	}
//
//	testCases := map[string]testCase{
//		"collapsed_simple": {
//			v: rackTypeTestCollapsedSimpleJSON,
//			e: rackTypeTestCollapsedSimple,
//		},
//		"collapsed_simple_with_access": {
//			v: rackTypeTestCollapsedSimpleWithAccessJSON,
//			e: rackTypeTestCollapsedSimpleWithAccess,
//		},
//		"collapsed_esi": {
//			v: rackTypeTestCollapsedESIJSON,
//			e: rackTypeTestCollapsedESI,
//		},
//	}
//
//	for tName, tCase := range testCases {
//		t.Run(tName, func(t *testing.T) {
//			t.Parallel()
//			var r RackType
//			err := json.Unmarshal([]byte(tCase.v), &r)
//			require.NoError(t, err)
//
//			require.Equal(t, tCase.e, r)
//		})
//	}
//}
