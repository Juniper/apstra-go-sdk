// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"hash"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
	"github.com/stretchr/testify/require"
)

func TestDigest(t *testing.T) {
	type testCase struct {
		v      ider
		h      hash.Hash
		expHex string
	}

	testCases := map[string]testCase{
		"tag_zero_value_md5": {
			v:      &Tag{},
			h:      md5.New(),
			expHex: "ee6d9c5b3212b57cbbf2ab1e2ad58343",
		},
		"tag_with_id_one_md5": {
			v:      &Tag{id: "one"},
			h:      md5.New(),
			expHex: "ee6d9c5b3212b57cbbf2ab1e2ad58343",
		},
		"tag_with_id_two_md5": {
			v:      &Tag{id: "two"},
			h:      md5.New(),
			expHex: "ee6d9c5b3212b57cbbf2ab1e2ad58343",
		},
		"tag_with_everything_md5": {
			v: &Tag{
				Label:          "label",
				Description:    "description",
				id:             "id",
				createdAt:      pointer.To(timeutils.TimeParseMust(time.RFC3339, "2006-01-02T15:04:05Z")),
				lastModifiedAt: pointer.To(timeutils.TimeParseMust(time.RFC3339, "2006-01-02T15:04:05Z")),
			},
			h:      md5.New(),
			expHex: "4aaea4ea9adf028373ae2a617b0c5f4e",
		},
		"tag_zero_value_sha256": {
			v:      &Tag{},
			h:      sha256.New(),
			expHex: "76446a5e9d7bbe112130e83f369d1a9c00a9c258f1f9f5d5c5ebf9c655c8677a",
		},
		"tag_with_id_one_sha256": {
			v:      &Tag{id: "one"},
			h:      sha256.New(),
			expHex: "76446a5e9d7bbe112130e83f369d1a9c00a9c258f1f9f5d5c5ebf9c655c8677a",
		},
		"tag_with_id_two_sha256": {
			v:      &Tag{id: "two"},
			h:      sha256.New(),
			expHex: "76446a5e9d7bbe112130e83f369d1a9c00a9c258f1f9f5d5c5ebf9c655c8677a",
		},
		"tag_with_everything_sha256": {
			v: &Tag{
				Label:          "label",
				Description:    "description",
				id:             "id",
				createdAt:      pointer.To(timeutils.TimeParseMust(time.RFC3339, "2006-01-02T15:04:05Z")),
				lastModifiedAt: pointer.To(timeutils.TimeParseMust(time.RFC3339, "2006-01-02T15:04:05Z")),
			},
			h:      sha256.New(),
			expHex: "5b4c3446bccb55d89ad7ef5e2fb04e8da5be0e509e72b528ab4fcae31031663d",
		},
		"logical_device_zero_value_md5": {
			v:      &LogicalDevice{},
			h:      md5.New(),
			expHex: "e66be34a23f2eea874b75957113e2f86",
		},
		"logical_device_with_id_one_md5": {
			v:      &LogicalDevice{id: "one"},
			h:      md5.New(),
			expHex: "e66be34a23f2eea874b75957113e2f86",
		},
		"logical_device_with_id_two_md5": {
			v:      &LogicalDevice{id: "two"},
			h:      md5.New(),
			expHex: "e66be34a23f2eea874b75957113e2f86",
		},
		"logical_device_with_everything_md5": {
			v: &LogicalDevice{
				Label:          "label",
				Panels:         []LogicalDevicePanel{},
				id:             "id",
				createdAt:      pointer.To(timeutils.TimeParseMust(time.RFC3339, "2006-01-02T15:04:05Z")),
				lastModifiedAt: pointer.To(timeutils.TimeParseMust(time.RFC3339, "2006-01-02T15:04:05Z")),
			},
			h:      md5.New(),
			expHex: "3c1180d1042f7f74ade356728269296a",
		},
		"logical_device_zero_value_sha256": {
			v:      &LogicalDevice{},
			h:      sha256.New(),
			expHex: "94588a302eb7b4439b74996bcf3b97c552ecb7eff5c23d9944f5d579af5adee7",
		},
		"logical_device_with_id_one_sha256": {
			v:      &LogicalDevice{id: "one"},
			h:      sha256.New(),
			expHex: "94588a302eb7b4439b74996bcf3b97c552ecb7eff5c23d9944f5d579af5adee7",
		},
		"logical_device_with_id_two_sha256": {
			v:      &LogicalDevice{id: "two"},
			h:      sha256.New(),
			expHex: "94588a302eb7b4439b74996bcf3b97c552ecb7eff5c23d9944f5d579af5adee7",
		},
		"logical_device_with_everything_sha256": {
			v: &LogicalDevice{
				Label:          "label",
				Panels:         []LogicalDevicePanel{},
				id:             "id",
				createdAt:      pointer.To(timeutils.TimeParseMust(time.RFC3339, "2006-01-02T15:04:05Z")),
				lastModifiedAt: pointer.To(timeutils.TimeParseMust(time.RFC3339, "2006-01-02T15:04:05Z")),
			},
			h:      sha256.New(),
			expHex: "2f86576c0b403b2a4ff77445ed469cad73b7a7f171ac4132d5552c7c58f5545c",
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			idBefore := tCase.v.ID()

			b, err := hashForComparison(tCase.v, tCase.h)
			require.NoError(t, err)
			require.Equal(t, tCase.expHex, fmt.Sprintf("%x", b))

			b = mustHashForComparison(tCase.v, tCase.h)
			require.Equal(t, tCase.expHex, fmt.Sprintf("%x", b))

			// ensure we didn't change the ID
			if idBefore == nil {
				require.Nil(t, tCase.v.ID())
			} else {
				idAfter := tCase.v.ID()
				require.NotNil(t, idAfter)
				require.Equal(t, *idBefore, *idAfter)
			}
		})
	}
}

func TestOrderedMarshalJSON(t *testing.T) {
	type testCase struct {
		keys     []string
		values   map[string]json.RawMessage
		expJSON  string
		expError bool
	}

	testCases := map[string]testCase{
		"simple_ordered_keys": {
			keys:    []string{"a", "b", "c"},
			values:  map[string]json.RawMessage{"a": json.RawMessage(`1`), "b": json.RawMessage(`2`), "c": json.RawMessage(`3`)},
			expJSON: `{"a":1,"b":2,"c":3}`,
		},
		"keys_values_length_mismatch": {
			keys:     []string{"a", "b"},
			values:   map[string]json.RawMessage{"a": json.RawMessage(`1`), "b": json.RawMessage(`2`), "c": json.RawMessage(`3`)},
			expError: true,
		},
		"duplicate_keys": {
			keys:     []string{"a", "b", "a"},
			values:   map[string]json.RawMessage{"a": json.RawMessage(`1`), "b": json.RawMessage(`2`)},
			expError: true,
		},
		"missing_value_for_key": {
			keys:     []string{"a", "b", "c"},
			values:   map[string]json.RawMessage{"a": json.RawMessage(`1`), "b": json.RawMessage(`2`), "d": json.RawMessage(`4`)},
			expError: true,
		},
		"empty_keys_and_values": {
			keys:    []string{},
			values:  map[string]json.RawMessage{},
			expJSON: `{}`,
		},
		"keys_with_escaped_chars": {
			keys:    []string{"a\"", "b\n"},
			values:  map[string]json.RawMessage{"a\"": json.RawMessage(`true`), "b\n": json.RawMessage(`false`)},
			expJSON: `{"a\"":true,"b\n":false}`,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result, err := orderedMarshalJSON(tCase.keys, tCase.values)
			if tCase.expError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tCase.expJSON, string(result))
		})
	}
}
