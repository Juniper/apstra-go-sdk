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

func TestTemplateRackBased_ID(t *testing.T) {
	t.Parallel()

	var obj TemplateRackBased
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

func TestTemplateRackBased_Replicate(t *testing.T) {
	t.Parallel()

	testCases := []TemplateRackBased{
		templateRackBasedL2VirtualEVPN,
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
			tc.Spine.LogicalDevice.id = ""
			tc.Spine.LogicalDevice.createdAt = nil
			tc.Spine.LogicalDevice.lastModifiedAt = nil
			for i := range tc.Racks {
				tc.Racks[i].RackType.id = ""
				tc.Racks[i].RackType.createdAt = nil
				tc.Racks[i].RackType.lastModifiedAt = nil
				for j := range tc.Racks[i].RackType.LeafSwitches {
					tc.Racks[i].RackType.LeafSwitches[j].LogicalDevice.id = ""
					tc.Racks[i].RackType.LeafSwitches[j].LogicalDevice.createdAt = nil
					tc.Racks[i].RackType.LeafSwitches[j].LogicalDevice.lastModifiedAt = nil
				}
				for j := range tc.Racks[i].RackType.AccessSwitches {
					tc.Racks[i].RackType.AccessSwitches[j].LogicalDevice.id = ""
					tc.Racks[i].RackType.AccessSwitches[j].LogicalDevice.createdAt = nil
					tc.Racks[i].RackType.AccessSwitches[j].LogicalDevice.lastModifiedAt = nil
				}
				for j := range tc.Racks[i].RackType.GenericSystems {
					tc.Racks[i].RackType.GenericSystems[j].LogicalDevice.id = ""
					tc.Racks[i].RackType.GenericSystems[j].LogicalDevice.createdAt = nil
					tc.Racks[i].RackType.GenericSystems[j].LogicalDevice.lastModifiedAt = nil
				}
			}

			require.Equal(t, tc, result)
		})
	}
}

func TestTemplateRackBased_timestamps(t *testing.T) {
	testCases := map[string]TemplateRackBased{
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

func TestTemplateRackBased_MarshalJSON(t *testing.T) {
	type testCase struct {
		v TemplateRackBased
		e string
	}

	testCases := map[string]testCase{
		"L2_Virtual_EVPN": {
			v: templateRackBasedL2VirtualEVPN,
			e: templateRackBasedL2VirtualEVPNJSON,
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

func TestTemplateRackBased_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		v string
		e TemplateRackBased
	}

	testCases := map[string]testCase{
		"L2_Virtual_EVPN": {
			v: templateRackBasedL2VirtualEVPNJSON,
			e: templateRackBasedL2VirtualEVPN,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var r TemplateRackBased
			err := json.Unmarshal([]byte(tCase.v), &r)
			require.NoError(t, err)
			require.Equal(t, tCase.e, r)
		})
	}
}

func TestTemplateRackBased_TemplateType(t *testing.T) {
	t.Parallel()

	r := TemplateRackBased{}.TemplateType()
	require.Equal(t, enum.TemplateTypeRackBased, r)
}
