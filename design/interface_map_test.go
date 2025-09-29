// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
	"github.com/stretchr/testify/require"
)

func TestInterfaceMap_ID(t *testing.T) {
	var obj InterfaceMap
	var id *string
	desiredId := testutils.RandString(6, "hex")

	t.Run("nil_ID_when_unset", func(t *testing.T) {
		id = obj.ID()
		require.Nil(t, id)
	})

	t.Run("set_id", func(t *testing.T) {
		err := obj.SetID(desiredId)
		require.NoError(t, err)
	})

	t.Run("check_id_after_set", func(t *testing.T) {
		id = obj.ID()
		require.NotNil(t, id)
		require.Equal(t, desiredId, *id)
	})

	t.Run("must_set_id", func(t *testing.T) {
		obj = zero.Of(obj)
		obj.MustSetID(desiredId)
	})

	t.Run("check_id_after_must_set", func(t *testing.T) {
		id = obj.ID()
		require.NotNil(t, id)
		require.Equal(t, desiredId, *id)
	})

	t.Run("must_set_id_panic", func(t *testing.T) {
		require.Panics(t, func() { obj.MustSetID(desiredId) })
	})
}

func TestInterfaceMap_timestamps(t *testing.T) {
	testCases := map[string]InterfaceMap{
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

func TestInterfaceMap_MarshalJSON(t *testing.T) {
	type testCase struct {
		v InterfaceMap
		e string
	}

	testCases := map[string]testCase{
		"test-48x10+4x100": {
			v: interfaceMapTestX,
			e: interfaceMapTestXJSON,
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
		})
	}
}

func TestInterfaceMap_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		v string
		e InterfaceMap
	}

	testCases := map[string]testCase{
		"test-48x10+4x100": {
			v: interfaceMapTestXJSON,
			e: interfaceMapTestX,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			var r InterfaceMap
			err := json.Unmarshal([]byte(tCase.v), &r)
			require.NoError(t, err)
			require.Equal(t, tCase.e, r)
		})
	}
}
