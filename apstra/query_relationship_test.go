// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra_test

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/stretchr/testify/require"
)

func TestRelationshipType_FromString(t *testing.T) {
	type testCase struct {
		data   string
		exp    apstra.RelationshipType
		expErr string
	}

	testCases := map[string]testCase{
		"valid_relationship_type_hosted_interfaces": {
			data: "hosted_interfaces",
			exp:  apstra.RelationshipTypeHostedInterfaces,
		},
		"invalid_relationship_type_bogus": {
			data:   "bogus",
			expErr: `unknown relationship type "bogus"`,
		},
		"empty_string": {
			data: "",
			exp:  apstra.RelationshipTypeNone,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var rt apstra.RelationshipType
			err := rt.FromString(tCase.data)
			if tCase.expErr != "" {
				require.Error(t, err)
				require.EqualError(t, err, tCase.expErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tCase.exp, rt)
		})
	}
}

func TestRelationshipType_MarshalText(t *testing.T) {
	type testCase struct {
		data   apstra.RelationshipType
		exp    string
		expErr string
	}

	testCases := map[string]testCase{
		"valid_relationship_type_hosted_interfaces": {
			data: apstra.RelationshipTypeHostedInterfaces,
			exp:  "hosted_interfaces",
		},
		"invalid_relationship_type": {
			data:   apstra.RelationshipType(-1),
			expErr: "cannot marshal: unknown relationship type -1",
		},
		"empty_string": {
			data: apstra.RelationshipTypeNone,
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

func TestRelationshipType_UnmarshalText(t *testing.T) {
	type testCase struct {
		data   string
		exp    apstra.RelationshipType
		expErr string
	}

	testCases := map[string]testCase{
		"valid_relationship_type_hosted_interfaces": {
			data: "hosted_interfaces",
			exp:  apstra.RelationshipTypeHostedInterfaces,
		},
		"invalid_relationship_type_bogus": {
			data:   "bogus",
			expErr: `unknown relationship type "bogus"`,
		},
		"empty_string": {
			data: "",
			exp:  apstra.RelationshipTypeNone,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var rt apstra.RelationshipType
			err := rt.UnmarshalText([]byte(tCase.data))
			if tCase.expErr != "" {
				require.Error(t, err)
				require.EqualError(t, err, tCase.expErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tCase.exp, rt)
		})
	}
}
