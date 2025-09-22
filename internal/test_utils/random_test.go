// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package testutils_test

import (
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestRandTime(t *testing.T) {
	type testCase struct {
		bounds []time.Time
	}

	// baseStart is testutils.RandTime's unbounded start time
	baseStart := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	testCases := map[string]testCase{
		"no_bounds": {},
		"one_bound": {
			bounds: []time.Time{time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		"two_bounds": {
			bounds: []time.Time{
				time.Date(2006, 10, 12, 0, 0, 0, 0, time.FixedZone("boston", -4*60*60)),
				time.Date(2006, 10, 12, 23, 59, 59, 1e9-1, time.FixedZone("boston", -4*60*60)),
			},
		},
		"reversed_bounds": {
			bounds: []time.Time{
				time.Date(2013, time.October, 12, 0, 0, 0, 0, time.FixedZone("nashua", -4*60*60)),
				time.Date(2013, time.October, 12, 23, 59, 59, 1e9-1, time.FixedZone("nashua", -4*60*60)),
			},
		},
		"range_too_small": {
			bounds: []time.Time{
				time.Date(1969, time.July, 20, 15, 17, 40, 0, time.FixedZone("houston", -5*60*60)),
				time.Date(1969, time.July, 20, 15, 17, 40, 1e9-1, time.FixedZone("houston", -5*60*60)),
			},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result := testutils.RandTime(tCase.bounds...)
			maxExpected := pointer.To(time.Now()) // collect "now" after collecting result to avoid race condition

			var minExpected *time.Time
			switch len(tCase.bounds) {
			case 0:
				minExpected = pointer.To(baseStart.Truncate(time.Second))
			case 1:
				minExpected = &tCase.bounds[0]
			default:
				minExpected = &tCase.bounds[0]
				maxExpected = &tCase.bounds[1]
			}

			if minExpected.After(*maxExpected) {
				minExpected, maxExpected = maxExpected, minExpected
			}

			if maxExpected.Sub(*minExpected) < time.Second {
				require.True(t, minExpected.Equal(result))
				return
			}

			require.True(t, minExpected.Before(result))
			require.True(t, maxExpected.After(result))
		})
	}
}
