// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package speed_test

import (
	"encoding/json"
	"testing"

	"github.com/Juniper/apstra-go-sdk/speed"
	"github.com/stretchr/testify/require"
)

func TestSpeed_Bps(t *testing.T) {
	type testCase struct {
		ss []speed.Speed
		e  int64
	}

	testCases := map[string]testCase{
		"10M": {
			ss: []speed.Speed{"10000000", "10m", " 10 m ", "10mbps", " 10 mbps ", "10mb/s", " 10 mb/s ", "10m", " 10 M ", "10Mbps", " 10 Mbps ", "10Mb/s", " 10 Mb/s "},
			e:  10_000_000,
		},
		"10G": {
			ss: []speed.Speed{"10000000000", "10g", " 10 g ", "10gbps", " 10 gbps ", "10gb/s", " 10 gb/s ", "10g", " 10 G ", "10Gbps", " 10 Gbps ", "10Gb/s", " 10 Gb/s "},
			e:  10_000_000_000,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			for _, s := range tCase.ss {
				t.Run(string(s), func(t *testing.T) {
					t.Parallel()

					require.Equal(t, tCase.e, s.Bps())
				})
			}
		})
	}
}

func TestSpeed_Equal(t *testing.T) {
	type testCase struct {
		a        speed.Speed
		b        speed.Speed
		notEqual bool
	}

	testCases := map[string]testCase{
		"t_same_suffix": {
			a: "10M",
			b: "10M",
		},
		"t_same_numeric": {
			a: "10000000",
			b: "10000000",
		},
		"t_different_case": {
			a: "10M",
			b: "10m",
		},
		"t_different_suffix": {
			a: "10M",
			b: "10mbps",
		},
		"t_different_whitespace": {
			a: "10M",
			b: " 10 m bps ",
		},
		"t_different_units": {
			a: "10G",
			b: "10000M",
		},
		"f_same_units": {
			a:        "1G",
			b:        "10G",
			notEqual: true,
		},
		"f_different_units": {
			a:        "10G",
			b:        "1000M",
			notEqual: true,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result := tCase.a.Equal(tCase.b)
			if tCase.notEqual {
				require.False(t, result)
			} else {
				require.True(t, result)
			}
		})
	}
}

func TestSpeed_MarshalJSON(t *testing.T) {
	type testCase struct {
		values      []speed.Speed
		expected    string
		expectError bool
	}

	testCases := map[string]testCase{
		"empty": {
			values:   []speed.Speed{""},
			expected: "null",
		},
		"10M": {
			values: []speed.Speed{
				"10000000",
				" 10000000 ",
				"10M",
				"10m",
				" 10 M ",
				" 10 m ",
				"10Mbps",
				"10mbps",
				" 10 Mbps ",
				" 10 mbps ",
				"10Mb/s",
				"10mb/s",
				" 10 Mb/s ",
				" 10 mb/s ",
			},
			expected: `{"unit":"M","value":10}`,
		},
		"10M_error": {
			values: []speed.Speed{
				"10000001",
				"9999999",
			},
			expectError: true,
		},
		"100M": {
			values: []speed.Speed{
				"100000000",
				" 100000000 ",
				"100M",
				"100m",
				" 100 M ",
				" 100 m ",
				"100Mbps",
				"100mbps",
				" 100 Mbps ",
				" 100 mbps ",
				"100Mb/s",
				"100mb/s",
				" 100 Mb/s ",
				" 100 mb/s ",
			},
			expected: `{"unit":"M","value":100}`,
		},
		"100M_error": {
			values: []speed.Speed{
				"100000001",
				"99999999",
			},
			expectError: true,
		},
		"1G": {
			values: []speed.Speed{
				"1000000000",
				" 1000000000 ",
				"1G",
				"1g",
				"1000M",
				"1000m",
				" 1000 M ",
				" 1000 m ",
				" 1 G ",
				" 1 g ",
				"1000Mbps",
				"1000mbps",
				"1Gbps",
				"1gbps",
				" 1000 Mbps ",
				" 1000 mbps ",
				" 1 Gbps ",
				" 1 gbps ",
				"1000Mb/s",
				"1000mb/s",
				"1Gb/s",
				"1gb/s",
				" 1000 Mb/s ",
				" 1000 mb/s ",
				" 1 Gb/s ",
				" 1 gb/s ",
			},
			expected: `{"unit":"G","value":1}`,
		},
		"1G_error": {
			values: []speed.Speed{
				"1000000001",
				"999999999",
				"1001M",
				"999M",
			},
			expectError: true,
		},
		"2G": {
			values: []speed.Speed{
				"2000000000",
				" 2000000000 ",
				"2G",
				"2g",
				"2000M",
				"2000m",
				" 2000 M ",
				" 2000 m ",
				" 2 G ",
				" 2 g ",
				"2000Mbps",
				"2000mbps",
				"2Gbps",
				"2gbps",
				" 2000 Mbps ",
				" 2000 mbps ",
				" 2 Gbps ",
				" 2 gbps ",
				"2000Mb/s",
				"2000mb/s",
				"2Gb/s",
				"2gb/s",
				" 2000 Mb/s ",
				" 2000 mb/s ",
				" 2 Gb/s ",
				" 2 gb/s ",
			},
			expected: `{"unit":"G","value":2}`,
		},
		"2G_error": {
			values: []speed.Speed{
				"2000000001",
				"1999999999",
				"2001M",
				"1999M",
			},
			expectError: true,
		},
		"5G": {
			values: []speed.Speed{
				"5000000000",
				" 5000000000 ",
				"5G",
				"5g",
				"5000M",
				"5000m",
				" 5000 M ",
				" 5000 m ",
				" 5 G ",
				" 5 g ",
				"5000Mbps",
				"5000mbps",
				"5Gbps",
				"5gbps",
				" 5000 Mbps ",
				" 5000 mbps ",
				" 5 Gbps ",
				" 5 gbps ",
				"5000Mb/s",
				"5000mb/s",
				"5Gb/s",
				"5gb/s",
				" 5000 Mb/s ",
				" 5000 mb/s ",
				" 5 Gb/s ",
				" 5 gb/s ",
			},
			expected: `{"unit":"G","value":5}`,
		},
		"5G_error": {
			values: []speed.Speed{
				"5000000001",
				"4999999999",
				"5001M",
				"4999M",
			},
			expectError: true,
		},
		"10G": {
			values: []speed.Speed{
				"10000000000",
				" 10000000000 ",
				"10G",
				"10g",
				"10000M",
				"10000m",
				" 10000 M ",
				" 10000 m ",
				" 10 G ",
				" 10 g ",
				"10000Mbps",
				"10000mbps",
				"10Gbps",
				"10gbps",
				" 10000 Mbps ",
				" 10000 mbps ",
				" 10 Gbps ",
				" 10 gbps ",
				"10000Mb/s",
				"10000mb/s",
				"10Gb/s",
				"10gb/s",
				" 10000 Mb/s ",
				" 10000 mb/s ",
				" 10 Gb/s ",
				" 10 gb/s ",
			},
			expected: `{"unit":"G","value":10}`,
		},
		"10G_error": {
			values: []speed.Speed{
				"10000000001",
				"9999999999",
				"10001M",
				"9999M",
			},
			expectError: true,
		},
		"25G": {
			values: []speed.Speed{
				"25000000000",
				" 25000000000 ",
				"25G",
				"25g",
				"25000M",
				"25000m",
				" 25000 M ",
				" 25000 m ",
				" 25 G ",
				" 25 g ",
				"25000Mbps",
				"25000mbps",
				"25Gbps",
				"25gbps",
				" 25000 Mbps ",
				" 25000 mbps ",
				" 25 Gbps ",
				" 25 gbps ",
				"25000Mb/s",
				"25000mb/s",
				"25Gb/s",
				"25gb/s",
				" 25000 Mb/s ",
				" 25000 mb/s ",
				" 25 Gb/s ",
				" 25 gb/s ",
			},
			expected: `{"unit":"G","value":25}`,
		},
		"25G_error": {
			values: []speed.Speed{
				"25000000001",
				"24999999999",
				"25001M",
				"24999M",
			},
			expectError: true,
		},
		"50G": {
			values: []speed.Speed{
				"50000000000",
				" 50000000000 ",
				"50G",
				"50g",
				"50000M",
				"50000m",
				" 50000 M ",
				" 50000 m ",
				" 50 G ",
				" 50 g ",
				"50000Mbps",
				"50000mbps",
				"50Gbps",
				"50gbps",
				" 50000 Mbps ",
				" 50000 mbps ",
				" 50 Gbps ",
				" 50 gbps ",
				"50000Mb/s",
				"50000mb/s",
				"50Gb/s",
				"50gb/s",
				" 50000 Mb/s ",
				" 50000 mb/s ",
				" 50 Gb/s ",
				" 50 gb/s ",
			},
			expected: `{"unit":"G","value":50}`,
		},
		"50G_error": {
			values: []speed.Speed{
				"50000000001",
				"49999999999",
				"50001M",
				"49999M",
			},
			expectError: true,
		},
		"100G": {
			values: []speed.Speed{
				"100000000000",
				" 100000000000 ",
				"100G",
				"100g",
				"100000M",
				"100000m",
				" 100000 M ",
				" 100000 m ",
				" 100 G ",
				" 100 g ",
				"100000Mbps",
				"100000mbps",
				"100Gbps",
				"100gbps",
				" 100000 Mbps ",
				" 100000 mbps ",
				" 100 Gbps ",
				" 100 gbps ",
				"100000Mb/s",
				"100000mb/s",
				"100Gb/s",
				"100gb/s",
				" 100000 Mb/s ",
				" 100000 mb/s ",
				" 100 Gb/s ",
				" 100 gb/s ",
			},
			expected: `{"unit":"G","value":100}`,
		},
		"100G_error": {
			values: []speed.Speed{
				"100000000001",
				"99999999999",
				"100001M",
				"99999M",
			},
			expectError: true,
		},
		"150G": {
			values: []speed.Speed{
				"150000000000",
				" 150000000000 ",
				"150G",
				"150g",
				"150000M",
				"150000m",
				" 150000 M ",
				" 150000 m ",
				" 150 G ",
				" 150 g ",
				"150000Mbps",
				"150000mbps",
				"150Gbps",
				"150gbps",
				" 150000 Mbps ",
				" 150000 mbps ",
				" 150 Gbps ",
				" 150 gbps ",
				"150000Mb/s",
				"150000mb/s",
				"150Gb/s",
				"150gb/s",
				" 150000 Mb/s ",
				" 150000 mb/s ",
				" 150 Gb/s ",
				" 150 gb/s ",
			},
			expected: `{"unit":"G","value":150}`,
		},
		"150G_error": {
			values: []speed.Speed{
				"150000000001",
				"149999999999",
				"150001M",
				"149999M",
			},
			expectError: true,
		},
		"200G": {
			values: []speed.Speed{
				"200000000000",
				" 200000000000 ",
				"200G",
				"200g",
				"200000M",
				"200000m",
				" 200000 M ",
				" 200000 m ",
				" 200 G ",
				" 200 g ",
				"200000Mbps",
				"200000mbps",
				"200Gbps",
				"200gbps",
				" 200000 Mbps ",
				" 200000 mbps ",
				" 200 Gbps ",
				" 200 gbps ",
				"200000Mb/s",
				"200000mb/s",
				"200Gb/s",
				"200gb/s",
				" 200000 Mb/s ",
				" 200000 mb/s ",
				" 200 Gb/s ",
				" 200 gb/s ",
			},
			expected: `{"unit":"G","value":200}`,
		},
		"200G_error": {
			values: []speed.Speed{
				"200000000001",
				"199999999999",
				"200001M",
				"199999M",
			},
			expectError: true,
		},
		"400G": {
			values: []speed.Speed{
				"400000000000",
				" 400000000000 ",
				"400G",
				"400g",
				"400000M",
				"400000m",
				" 400000 M ",
				" 400000 m ",
				" 400 G ",
				" 400 g ",
				"400000Mbps",
				"400000mbps",
				"400Gbps",
				"400gbps",
				" 400000 Mbps ",
				" 400000 mbps ",
				" 400 Gbps ",
				" 400 gbps ",
				"400000Mb/s",
				"400000mb/s",
				"400Gb/s",
				"400gb/s",
				" 400000 Mb/s ",
				" 400000 mb/s ",
				" 400 Gb/s ",
				" 400 gb/s ",
			},
			expected: `{"unit":"G","value":400}`,
		},
		"400G_error": {
			values: []speed.Speed{
				"400000000001",
				"399999999999",
				"400001M",
				"399999M",
			},
			expectError: true,
		},
		"800G": {
			values: []speed.Speed{
				"800000000000",
				" 800000000000 ",
				"800G",
				"800g",
				"800000M",
				"800000m",
				" 800000 M ",
				" 800000 m ",
				" 800 G ",
				" 800 g ",
				"800000Mbps",
				"800000mbps",
				"800Gbps",
				"800gbps",
				" 800000 Mbps ",
				" 800000 mbps ",
				" 800 Gbps ",
				" 800 gbps ",
				"800000Mb/s",
				"800000mb/s",
				"800Gb/s",
				"800gb/s",
				" 800000 Mb/s ",
				" 800000 mb/s ",
				" 800 Gb/s ",
				" 800 gb/s ",
			},
			expected: `{"unit":"G","value":800}`,
		},
		"800G_error": {
			values: []speed.Speed{
				"800000000001",
				"799999999999",
				"800001M",
				"799999M",
			},
			expectError: true,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			for _, value := range tCase.values {
				t.Run(string(value), func(t *testing.T) {
					t.Parallel()

					result, err := json.Marshal(value)
					if tCase.expectError {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
					}

					require.Equal(t, tCase.expected, string(result))
				})
			}
		})
	}
}
