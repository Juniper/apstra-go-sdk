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

	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestLogicalDevice_ID(t *testing.T) {
	t.Parallel()

	var obj LogicalDevice
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

func TestLogicalDevice_Replicate(t *testing.T) {
	t.Parallel()

	testCases := []LogicalDevice{
		logicalDeviceTest1x1,
		logicalDeviceTest48x10plus4x100,
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

func TestLogicalDevice_timestamps(t *testing.T) {
	testCases := map[string]LogicalDevice{
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

func TestLogicalDevice_MarshalJSON(t *testing.T) {
	type testCase struct {
		v LogicalDevice
		e string
	}

	testCases := map[string]testCase{
		"test-1x1": {
			v: logicalDeviceTest1x1,
			e: logicalDeviceTest1x1JSON,
		},
		"test-48x10+4x100": {
			v: logicalDeviceTest48x10plus4x100,
			e: logicalDeviceTest48x10plus4x100JSON,
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

func TestLogicalDevice_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		v string
		e LogicalDevice
	}

	testCases := map[string]testCase{
		"test-1x1": {
			v: logicalDeviceTest1x1JSON,
			e: logicalDeviceTest1x1,
		},
		"test-48x10+4x100": {
			v: logicalDeviceTest48x10plus4x100JSON,
			e: logicalDeviceTest48x10plus4x100,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var r LogicalDevice
			err := json.Unmarshal([]byte(tCase.v), &r)
			require.NoError(t, err)
			require.Equal(t, tCase.e, r)
		})
	}
}
