// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/sha256"
	"encoding/json"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/speed"
	"github.com/stretchr/testify/require"
)

func TestRackTypeID(t *testing.T) {
	var obj RackType
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
		obj = testutils.ZeroOf(obj)
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

func TestRackType_replicate(t *testing.T) {
	original := RackType{
		Label:                    testutils.RandString(6, "hex"),
		Description:              testutils.RandString(6, "hex"),
		FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
		LeafSwitches: []LeafSwitch{
			{
				Label:             testutils.RandString(6, "hex"),
				LinkPerSpineCount: pointer.To(3),
				LinkPerSpineSpeed: pointer.To(speed.Speed("10G")),
				LogicalDevice: LogicalDevice{
					Label: testutils.RandString(6, "hex"),
					Panels: []LogicalDevicePanel{
						{
							PanelLayout: LogicalDevicePanelLayout{RowCount: 8, ColumnCount: 8},
							PortGroups: []LogicalDevicePanelPortGroup{
								{
									Count: 64,
									Speed: "10G",
									Roles: LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleGeneric},
								},
							},
							PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
						},
					},
				},
				Tags: []Tag{
					{
						Label:       testutils.RandString(6, "hex"),
						Description: testutils.RandString(10, "hex"),
					},
				},
			},
		},
		id:             testutils.RandString(6, "hex"),
		createdAt:      pointer.To(testutils.RandTime(time.Now().Add(-2 * time.Minute))),
		lastModifiedAt: pointer.To(testutils.RandTime(time.Now().Add(-1 * time.Minute))),
	}

	replicant := original.replicate()

	require.Equal(t, mustHashForComparison(original, sha256.New()), mustHashForComparison(replicant, sha256.New()))

	// wipe out values which cannot be replicated before comparing values
	original.id = ""
	original.createdAt = nil
	original.lastModifiedAt = nil

	require.Equal(t, original, replicant)
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
		"collapsed_one_of_each": {
			v: rackTypeTestCollapsedOneOfEach,
			e: rackTypeTestCollapsedOneOfEachJSON,
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
		})
	}
}

func TestRackType_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		e RackType
		v string
	}

	testCases := map[string]testCase{
		"collapsed_one_of_each": {
			v: rackTypeTestCollapsedOneOfEachJSON,
			e: rackTypeTestCollapsedOneOfEach,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			var r RackType
			err := json.Unmarshal([]byte(tCase.v), &r)
			require.NoError(t, err)

			//// clear test data which doesn't marshal
			//tCase.e.createdAt = nil
			//tCase.e.lastModifiedAt = nil

			require.Equal(t, tCase.e, r)
		})
	}
}
