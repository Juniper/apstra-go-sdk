// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package antiaffinity_test

import (
	"encoding/json"
	sdk "github.com/Juniper/apstra-go-sdk"
	antiaffinity "github.com/Juniper/apstra-go-sdk/anti_affinity"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestAntiAffinityPolicy_MarshalJSON(t *testing.T) {
	type testCase struct {
		v antiaffinity.Policy
		e string
	}

	testCases := map[string]testCase{
		"disabled": {
			v: antiaffinity.Policy{
				MaxLinksPerPort:          1,
				MaxLinksPerSlot:          2,
				MaxPerSystemLinksPerPort: 3,
				MaxPerSystemLinksPerSlot: 4,
				Mode:                     enum.AntiAffinityModeDisabled,
			},
			e: `{"mode":"disabled","algorithm":"heuristic","max_links_per_port":1,"max_links_per_slot":2,"max_per_system_links_per_port":3,"max_per_system_links_per_slot":4}`,
		},
		"loose": {
			v: antiaffinity.Policy{
				MaxLinksPerPort:          5,
				MaxLinksPerSlot:          6,
				MaxPerSystemLinksPerPort: 7,
				MaxPerSystemLinksPerSlot: 8,
				Mode:                     enum.AntiAffinityModeLoose,
			},
			e: `{"mode":"enabled_loose","algorithm":"heuristic","max_links_per_port":5,"max_links_per_slot":6,"max_per_system_links_per_port":7,"max_per_system_links_per_slot":8}`,
		},
		"strict": {
			v: antiaffinity.Policy{
				MaxLinksPerPort:          9,
				MaxLinksPerSlot:          10,
				MaxPerSystemLinksPerPort: 11,
				MaxPerSystemLinksPerSlot: 12,
				Mode:                     enum.AntiAffinityModeStrict,
			},
			e: `{"mode":"enabled_strict","algorithm":"heuristic","max_links_per_port":9,"max_links_per_slot":10,"max_per_system_links_per_port":11,"max_per_system_links_per_slot":12}`,
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

func TestAntiAffinityPolicy_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		v      string
		e      antiaffinity.Policy
		expErr error
	}

	testCases := map[string]testCase{
		"disabled": {
			v: `{"mode":"disabled","algorithm":"heuristic","max_links_per_port":1,"max_links_per_slot":2,"max_per_system_links_per_port":3,"max_per_system_links_per_slot":4}`,
			e: antiaffinity.Policy{
				MaxLinksPerPort:          1,
				MaxLinksPerSlot:          2,
				MaxPerSystemLinksPerPort: 3,
				MaxPerSystemLinksPerSlot: 4,
				Mode:                     enum.AntiAffinityModeDisabled,
			},
		},
		"loose": {
			v: `{"mode":"enabled_loose","algorithm":"heuristic","max_links_per_port":5,"max_links_per_slot":6,"max_per_system_links_per_port":7,"max_per_system_links_per_slot":8}`,
			e: antiaffinity.Policy{
				MaxLinksPerPort:          5,
				MaxLinksPerSlot:          6,
				MaxPerSystemLinksPerPort: 7,
				MaxPerSystemLinksPerSlot: 8,
				Mode:                     enum.AntiAffinityModeLoose,
			},
		},
		"strict": {
			v: `{"mode":"enabled_strict","algorithm":"heuristic","max_links_per_port":9,"max_links_per_slot":10,"max_per_system_links_per_port":11,"max_per_system_links_per_slot":12}`,
			e: antiaffinity.Policy{
				MaxLinksPerPort:          9,
				MaxLinksPerSlot:          10,
				MaxPerSystemLinksPerPort: 11,
				MaxPerSystemLinksPerSlot: 12,
				Mode:                     enum.AntiAffinityModeStrict,
			},
		},
		"bogus_algorithm": {
			v:      `{"algorithm":"bogus"}`,
			expErr: new(sdk.ErrAPIResponseInvalid),
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var r antiaffinity.Policy
			err := json.Unmarshal([]byte(tCase.v), &r)
			if tCase.expErr != nil {
				target := reflect.New(reflect.TypeOf(tCase.expErr).Elem()).Interface()
				require.ErrorAs(t, err, target)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tCase.e, r)
		})
	}
}
