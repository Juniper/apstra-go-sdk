// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInterfaceMapDigest_ID(t *testing.T) {
	var obj InterfaceMapDigest
	var id *string
	desiredId := testutils.RandString(6, "hex")

	t.Run("nil_ID_when_unset", func(t *testing.T) {
		id = obj.ID()
		require.Nil(t, id)
	})

	t.Run("check_id_after_set", func(t *testing.T) {
		obj.id = desiredId
		id = obj.ID()
		require.NotNil(t, id)
		require.Equal(t, desiredId, *id)
	})
}

func TestInterfaceMapDigest_timestamps(t *testing.T) {
	testCases := map[string]InterfaceMapDigest{
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

func TestInterfaceMapDigest_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		v string
		e InterfaceMapDigest
	}

	testCases := map[string]testCase{
		"Juniper_QFX5120_32C__AOS_32x40_4": {
			v: interfaceMapDigestJuniper_QFX5120_32C__AOS_32x40_4JSON,
			e: interfaceMapDigestJuniper_QFX5120_32C__AOS_32x40_4,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			var r InterfaceMapDigest
			err := json.Unmarshal([]byte(tCase.v), &r)
			require.NoError(t, err)
			require.Equal(t, tCase.e, r)
		})
	}
}
