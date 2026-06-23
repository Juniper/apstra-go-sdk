// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter_test

import (
	"encoding/json"
	"testing"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/stretchr/testify/require"
)

func TestPortRange_MarshalText(t *testing.T) {
	type testCase struct {
		in     datacenter.PortRange
		exp    string
		expErr bool
	}

	testCases := map[string]testCase{
		"single_port": {
			in:  datacenter.PortRange{First: 80, Last: 80},
			exp: "80",
		},
		"range": {
			in:  datacenter.PortRange{First: 20, Last: 21},
			exp: "20-21",
		},
		"invalid_reversed": {
			in:     datacenter.PortRange{First: 21, Last: 20},
			expErr: true,
		},
		"invalid_zero_start": {
			in:     datacenter.PortRange{First: 0, Last: 10},
			expErr: true,
		},
		"invalid_zero_end": {
			in:     datacenter.PortRange{First: 65535, Last: 0},
			expErr: true,
		},
		"invalid_zero_both": {
			in:     datacenter.PortRange{First: 0, Last: 0},
			expErr: true,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			result, err := tCase.in.MarshalText()
			if tCase.expErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tCase.exp, string(result))
		})
	}
}

func TestPortRange_UnmarshalText(t *testing.T) {
	type testCase struct {
		in     string
		exp    datacenter.PortRange
		expErr bool
	}

	testCases := map[string]testCase{
		"single_port": {
			in:  "80",
			exp: datacenter.PortRange{First: 80, Last: 80},
		},
		"single_port_with_whitespace": {
			in:  " 80 ",
			exp: datacenter.PortRange{First: 80, Last: 80},
		},
		"range": {
			in:  "20-21",
			exp: datacenter.PortRange{First: 20, Last: 21},
		},
		"range_with_whitespace": {
			in:  " 20 - 21 ",
			exp: datacenter.PortRange{First: 20, Last: 21},
		},
		"invalid_empty": {
			in:     "",
			expErr: true,
		},
		"invalid_space": {
			in:     " ",
			expErr: true,
		},
		"invalid_too_many_parts": {
			in:     "1-2-3",
			expErr: true,
		},
		"invalid_non_numeric": {
			in:     "abc",
			expErr: true,
		},
		"invalid_zero_single": {
			in:     "0",
			expErr: true,
		},
		"invalid_zero_range": {
			in:     "0-10",
			expErr: true,
		},
		"invalid_reversed": {
			in:     "21-20",
			expErr: true,
		},
		"invalid_overflow": {
			in:     "65536",
			expErr: true,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			var result datacenter.PortRange
			err := result.UnmarshalText([]byte(tCase.in))
			if tCase.expErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tCase.exp, result)
		})
	}
}

func TestPortRange_UnmarshalText_NilReceiver(t *testing.T) {
	var p *datacenter.PortRange
	err := p.UnmarshalText([]byte("80"))
	require.Error(t, err)
	require.ErrorContains(t, err, "nil *PortRange")
}

func TestPortRange_UnmarshalText_DoesNotMutateOnError(t *testing.T) {
	original := datacenter.PortRange{First: 80, Last: 80}
	got := original

	err := got.UnmarshalText([]byte("0"))
	require.Error(t, err)
	require.Equal(t, original, got)
}

func TestPortRange_JSONUsesTextMarshaler(t *testing.T) {
	type payload struct {
		Ports datacenter.PortRange `json:"ports"`
	}

	t.Run("marshal_as_json_string", func(t *testing.T) {
		in := payload{
			Ports: datacenter.PortRange{First: 20, Last: 21},
		}

		b, err := json.Marshal(in)
		require.NoError(t, err)
		require.JSONEq(t, `{"ports":"20-21"}`, string(b))
	})

	t.Run("unmarshal_from_json_string", func(t *testing.T) {
		var out payload
		err := json.Unmarshal([]byte(`{"ports":"80"}`), &out)
		require.NoError(t, err)
		require.Equal(t, datacenter.PortRange{First: 80, Last: 80}, out.Ports)
	})
}

func TestPortRange_canonicalize(t *testing.T) {
	type testCase struct {
		in  datacenter.PortRange
		exp datacenter.PortRange
	}

	testCases := map[string]testCase{
		"single_port": {
			in:  datacenter.PortRange{First: 80, Last: 80},
			exp: datacenter.PortRange{First: 80, Last: 80},
		},
		"ordered_range": {
			in:  datacenter.PortRange{First: 80, Last: 81},
			exp: datacenter.PortRange{First: 80, Last: 81},
		},
		"unordered_range": {
			in:  datacenter.PortRange{First: 81, Last: 80},
			exp: datacenter.PortRange{First: 80, Last: 81},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			tCase.in.Canonicalize()
			require.Equal(t, tCase.exp, tCase.in)
		})
	}
}
