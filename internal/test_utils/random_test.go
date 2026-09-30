// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package testutils_test

import (
	"fmt"
	"net/netip"
	"strings"
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

func TestRandomHardwareAddr(t *testing.T) {
	type testCase struct {
		set   []byte
		unset []byte
	}

	testCases := map[string]testCase{
		"laa": {
			set: []byte{2},
		},
		"group": {
			set: []byte{1},
		},
		"laa_and_not_group": {
			set:   []byte{2},
			unset: []byte{1},
		},
		"group_and_not_laa": {
			set:   []byte{1},
			unset: []byte{2},
		},
		"laa_and_group": {
			set: []byte{3},
		},
		"last_byte_128": {
			set:   []byte{0, 0, 0, 0, 0, 128},
			unset: []byte{0, 0, 0, 0, 0, 127},
		},
		"last_byte_high": {
			set: []byte{0, 0, 0, 0, 0, 128},
		},
		"last_byte_low": {
			unset: []byte{0, 0, 0, 0, 0, 128},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result := testutils.RandomHardwareAddr(tCase.set, tCase.unset)

			for i, setByte := range tCase.set {
				require.Equal(t, setByte, result[i]&setByte)
			}

			for i, unsetByte := range tCase.unset {
				require.Equal(t, ^unsetByte, result[i]|^unsetByte)
			}
		})
	}
}

func TestRandomPrefixes(t *testing.T) {
	type testCase struct {
		cidr  string
		count int
		bits  int
	}

	testCases := map[string]testCase{
		"10.0.0.0/8_24_5": {
			cidr:  "10.0.0.0/8",
			bits:  24,
			count: 5,
		},
		"192.0.2.0/24_26_3": {
			cidr:  "192.0.2.0/24",
			bits:  26,
			count: 3,
		},
		"192.168.0.0/16_24_256": {
			cidr:  "192.168.0.0/16",
			bits:  24,
			count: 256,
		},
		"3fff::/20_64_100": {
			cidr:  "3fff::/20",
			bits:  64,
			count: 100,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			// parse the container block - we'll use it during validation
			cidr, err := netip.ParsePrefix(tCase.cidr)
			require.NoError(t, err)

			got := testutils.RandomPrefixes(t, tCase.cidr, tCase.bits, tCase.count)
			require.Equal(t, tCase.count, len(got))

			sb := new(strings.Builder)
			sb.WriteString(fmt.Sprintf("%s -> /%d, %d samples\n", tCase.cidr, tCase.bits, tCase.count))
			for i, p := range got {
				require.True(t, cidr.Contains(p.Addr()))
				require.Equal(t, tCase.bits, p.Bits())
				sb.WriteString(fmt.Sprintf("\t%d\t %s\n", i+1, p.String()))
			}
			t.Log(sb.String())
		})
	}
}
