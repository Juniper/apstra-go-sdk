// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestTemplateRailCollapsed_ID(t *testing.T) {
	t.Parallel()

	var obj TemplateRailCollapsed
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

	t.Run("set_id_panic", func(t *testing.T) {
		require.Panics(t, func() { obj.SetID(desiredId) })
	})
}

func TestTemplateRailCollapsed_timestamps(t *testing.T) {
	testCases := map[string]TemplateRailCollapsed{
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

func TestTemplateRailCollapsed_MarshalJSON(t *testing.T) {
	type testCase struct {
		v TemplateRailCollapsed
		e string
	}

	testCases := map[string]testCase{
		"collapsed_fabric_128gpu": {
			v: railCollapsedSmall,
			e: railCollapsedSmallJSON,
		},
		"collapsed_fabric_512gpu": {
			v: railCollapsedMedium,
			e: railCollapsedMediumJSON,
		},
		"collapsed_fabric_1024gpu": {
			v: railCollapsedLarge,
			e: railCollapsedLargeJSON,
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
			delete(eMap, "id")
			delete(eMap, "created_at")
			delete(eMap, "last_modified_at")
			e, err := json.Marshal(eMap)
			require.NoError(t, err)

			require.JSONEq(t, string(e), string(r))
		})
	}
}

func TestTemplateRailCollapsed_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		v string
		e TemplateRailCollapsed
	}

	testCases := map[string]testCase{
		"collapsed_fabric_128gpu": {
			v: railCollapsedSmallJSON,
			e: railCollapsedSmall,
		},
		"collapsed_fabric_512gpu": {
			v: railCollapsedMediumJSON,
			e: railCollapsedMedium,
		},
		"collapsed_fabric_1024gpu": {
			v: railCollapsedLargeJSON,
			e: railCollapsedLarge,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var r TemplateRailCollapsed
			err := json.Unmarshal([]byte(tCase.v), &r)
			require.NoError(t, err)
			require.Equal(t, tCase.e, r)
		})
	}
}

func TestTemplateRailCollapsed_TemplateType(t *testing.T) {
	t.Parallel()

	r := TemplateRailCollapsed{}.TemplateType()
	require.Equal(t, enum.TemplateTypeRailCollapsed, r)
}
