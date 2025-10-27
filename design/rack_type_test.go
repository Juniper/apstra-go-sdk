// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestRackType_ID(t *testing.T) {
	t.Parallel()

	var obj RackType
	var id *string
	desiredId := testutils.RandString(6, "hex")

	t.Run("nil_ID_when_unset", func(t *testing.T) {
		id = obj.ID()
		require.Nil(t, id)
	})

	t.Run("set_id", func(t *testing.T) {
		obj.SetID(desiredId)
	})

	t.Run("check_id_after_set", func(t *testing.T) {
		id = obj.ID()
		require.NotNil(t, id)
		require.Equal(t, desiredId, *id)
	})

	t.Run("check_id_after_must_set", func(t *testing.T) {
		id = obj.ID()
		require.NotNil(t, id)
		require.Equal(t, desiredId, *id)
	})
}

func TestRackType_Replicate(t *testing.T) {
	t.Parallel()

	testCases := []RackType{
		rackTypeTestCollapsedSimple,
		rackTypeTestCollapsedESI,
		rackTypeTestCollapsedSimpleWithAccess,
		rackTypeTestRackBasedESIWithAccessESI,
		rackTypeTestRackBasedMLAGWithAccessPair,
		rackTypeTestCollapsedESIWithGenericSystems,
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			t.Parallel()

			result := tc.Replicate()

			require.Equal(t, mustHashForComparison(tc, sha256.New()), mustHashForComparison(result, sha256.New()))

			// wipe out values which cannot be replicated before comparing values
			tc.id = ""
			tc.createdAt = nil
			tc.lastModifiedAt = nil

			require.Equal(t, tc, result)
		})
	}
}

func TestRackType_timestamps(t *testing.T) {
	testCases := map[string]RackType{
		"a": {
			createdAt:      pointer.To(testutils.RandTime(time.Now().Add(-2 * time.Minute))),
			lastModifiedAt: pointer.To(testutils.RandTime(time.Now().Add(-1 * time.Minute))),
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			t.Run("created_at", func(t *testing.T) {
				t.Parallel()
				createdAt := tCase.CreatedAt()
				require.NotNil(t, createdAt)
				require.Equal(t, *tCase.createdAt, *createdAt)
			})
			t.Run("last_modified_at", func(t *testing.T) {
				t.Parallel()
				lastModifiedAt := tCase.LastModifiedAt()
				require.NotNil(t, lastModifiedAt)
				require.Equal(t, *tCase.lastModifiedAt, *lastModifiedAt)
			})
		})
	}
}

func TestRackType_MarshalJSON(t *testing.T) {
	type testCase struct {
		v RackType
		e string
	}

	testCases := map[string]testCase{
		"collapsed_simple": {
			v: rackTypeTestCollapsedSimple,
			e: rackTypeTestCollapsedSimpleJSON,
		},
		"collapsed_simple_with_access": {
			v: rackTypeTestCollapsedSimpleWithAccess,
			e: rackTypeTestCollapsedSimpleWithAccessJSON,
		},
		"collapsed_esi": {
			v: rackTypeTestCollapsedESI,
			e: rackTypeTestCollapsedESIJSON,
		},
		"rack_based_esi_with_access_esi": {
			v: rackTypeTestRackBasedESIWithAccessESI,
			e: rackTypeTestRackBasedESIWithAccessESIJSON,
		},
		"rack_based_mlag_with_access_pair": {
			v: rackTypeTestRackBasedMLAGWithAccessPair,
			e: rackTypeTestRackBasedMLAGWithAccessPairJSON,
		},
		"collapsed_esi_with_generic_systems": {
			v: rackTypeTestCollapsedESIWithGenericSystems,
			e: rackTypeTestCollapsedESIWithGenericSystemsJSON,
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
			delete(eMap, "status")
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

func TestRackType_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		e RackType
		v string
	}

	testCases := map[string]testCase{
		"collapsed_simple": {
			v: rackTypeTestCollapsedSimpleJSON,
			e: rackTypeTestCollapsedSimple,
		},
		"collapsed_simple_with_access": {
			v: rackTypeTestCollapsedSimpleWithAccessJSON,
			e: rackTypeTestCollapsedSimpleWithAccess,
		},
		"collapsed_esi": {
			v: rackTypeTestCollapsedESIJSON,
			e: rackTypeTestCollapsedESI,
		},
		"collapsed_esi_with_generic_systems": {
			e: rackTypeTestCollapsedESIWithGenericSystems,
			v: rackTypeTestCollapsedESIWithGenericSystemsJSON,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var r RackType
			err := json.Unmarshal([]byte(tCase.v), &r)
			require.NoError(t, err)

			// set attributes known to marshal differently than the origin struct
			for i, accessSwitch := range tCase.e.AccessSwitches {
				for j, link := range accessSwitch.Links {
					if link.AttachmentType.String() == "" {
						tCase.e.AccessSwitches[i].Links[j].AttachmentType = enum.LinkAttachmentTypeSingle
					}
				}
			}

			require.Equal(t, tCase.e, r)
		})
	}
}
