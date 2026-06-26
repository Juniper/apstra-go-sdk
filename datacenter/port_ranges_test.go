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

func TestPortRanges_MarshalText(t *testing.T) {
	type testCase struct {
		in     datacenter.PortRanges
		exp    string
		expErr bool
	}

	testCases := map[string]testCase{
		"nil_slice_encodes_as_any": {
			in:  nil,
			exp: "any",
		},
		"empty_slice_encodes_as_any": {
			in:  datacenter.PortRanges{},
			exp: "any",
		},
		"single_port": {
			in: datacenter.PortRanges{
				{First: 80, Last: 80},
			},
			exp: "80",
		},
		"single_range": {
			in: datacenter.PortRanges{
				{First: 20, Last: 21},
			},
			exp: "20-21",
		},
		"many_ranges": {
			in: datacenter.PortRanges{
				{First: 20, Last: 21},
				{First: 80, Last: 80},
				{First: 443, Last: 443},
			},
			exp: "20-21,80,443",
		},
		"with_invalid_range": {
			in: datacenter.PortRanges{
				{First: 20, Last: 21},
				{First: 0, Last: 10},
			},
			expErr: true,
		},
		"invalid_order": {
			in: datacenter.PortRanges{
				{First: 81, Last: 81},
				{First: 80, Last: 80},
			},
			expErr: true,
		},
		"invalid_overlap": {
			in: datacenter.PortRanges{
				{First: 100, Last: 200},
				{First: 150, Last: 250},
			},
			expErr: true,
		},
		"invalid_a_contains_b": {
			in: datacenter.PortRanges{
				{First: 100, Last: 200},
				{First: 150, Last: 150},
			},
			expErr: true,
		},
		"invalid_b_contains_a": {
			in: datacenter.PortRanges{
				{First: 150, Last: 150},
				{First: 100, Last: 200},
			},
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

func TestPortRanges_UnmarshalText(t *testing.T) {
	type testCase struct {
		in     string
		exp    datacenter.PortRanges
		expErr bool
	}

	testCases := map[string]testCase{
		"any": {
			in:  "any",
			exp: nil, // current implementation sets *prs = nil
		},
		"single_port": {
			in: "80",
			exp: datacenter.PortRanges{
				{First: 80, Last: 80},
			},
		},
		"single_range": {
			in: "20-21",
			exp: datacenter.PortRanges{
				{First: 20, Last: 21},
			},
		},
		"many_ranges": {
			in: "20-21,80,443",
			exp: datacenter.PortRanges{
				{First: 20, Last: 21},
				{First: 80, Last: 80},
				{First: 443, Last: 443},
			},
		},
		"whitespace": {
			in: " 20 - 21 , 80 ",
			exp: datacenter.PortRanges{
				{First: 20, Last: 21},
				{First: 80, Last: 80},
			},
		},
		"invalid_empty": {
			in:     "",
			expErr: true,
		},
		"invalid_trailing_separator": {
			in:     "80,",
			expErr: true,
		},
		"invalid_leading_separator": {
			in:     ",80",
			expErr: true,
		},
		"invalid_double_separator": {
			in:     "80,,443",
			expErr: true,
		},
		"invalid_nested_range": {
			in:     "20-21,21-20",
			expErr: true,
		},
		"invalid_non_numeric": {
			in:     "20-21,abc",
			expErr: true,
		},
		"invalid_zero": {
			in:     "0,80",
			expErr: true,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			var result datacenter.PortRanges
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

func TestPortRanges_UnmarshalText_NilReceiver(t *testing.T) {
	var prs *datacenter.PortRanges
	err := prs.UnmarshalText([]byte("80"))
	require.Error(t, err)
	require.ErrorContains(t, err, "nil *PortRanges")
}

func TestPortRanges_UnmarshalText_DoesNotMutateOnError(t *testing.T) {
	original := datacenter.PortRanges{{First: 80, Last: 80}}
	result := make(datacenter.PortRanges, len(original))
	copy(result, original)

	err := result.UnmarshalText([]byte("20-21,0"))
	require.Error(t, err)
	require.Equal(t, original, result)
}

func TestPortRanges_JSONUsesTextMarshaler(t *testing.T) {
	type payload struct {
		Ports datacenter.PortRanges `json:"ports"`
	}

	t.Run("marshal_any_as_json_string", func(t *testing.T) {
		in := payload{
			Ports: nil,
		}

		b, err := json.Marshal(in)
		require.NoError(t, err)
		require.JSONEq(t, `{"ports":"any"}`, string(b))
	})

	t.Run("marshal_ranges_as_json_string", func(t *testing.T) {
		in := payload{
			Ports: datacenter.PortRanges{
				{First: 20, Last: 21},
				{First: 80, Last: 80},
			},
		}

		b, err := json.Marshal(in)
		require.NoError(t, err)
		require.JSONEq(t, `{"ports":"20-21,80"}`, string(b))
	})

	t.Run("unmarshal_any_from_json_string", func(t *testing.T) {
		var out payload
		err := json.Unmarshal([]byte(`{"ports":"any"}`), &out)
		require.NoError(t, err)
		require.Nil(t, out.Ports)
	})

	t.Run("unmarshal_ranges_from_json_string", func(t *testing.T) {
		var out payload
		err := json.Unmarshal([]byte(`{"ports":"20-21,80"}`), &out)
		require.NoError(t, err)
		require.Equal(t, datacenter.PortRanges{
			{First: 20, Last: 21},
			{First: 80, Last: 80},
		}, out.Ports)
	})
}

func TestPortRanges_canonicalize(t *testing.T) {
	type testCase struct {
		in  datacenter.PortRanges
		exp datacenter.PortRanges
	}

	testCases := map[string]testCase{
		"single_port": {
			in:  datacenter.PortRanges{{First: 80, Last: 80}},
			exp: datacenter.PortRanges{{First: 80, Last: 80}},
		},
		"single_range": {
			in:  datacenter.PortRanges{{First: 100, Last: 200}},
			exp: datacenter.PortRanges{{First: 100, Last: 200}},
		},
		"individual_ports_ordered": {
			in:  datacenter.PortRanges{{First: 22, Last: 22}, {First: 80, Last: 80}, {First: 443, Last: 443}},
			exp: datacenter.PortRanges{{First: 22, Last: 22}, {First: 80, Last: 80}, {First: 443, Last: 443}},
		},
		"port_ranges_ordered": {
			in:  datacenter.PortRanges{{First: 100, Last: 200}, {First: 300, Last: 400}, {First: 500, Last: 600}},
			exp: datacenter.PortRanges{{First: 100, Last: 200}, {First: 300, Last: 400}, {First: 500, Last: 600}},
		},
		"blend_ordered": {
			in:  datacenter.PortRanges{{First: 100, Last: 200}, {First: 250, Last: 250}, {First: 300, Last: 400}},
			exp: datacenter.PortRanges{{First: 100, Last: 200}, {First: 250, Last: 250}, {First: 300, Last: 400}},
		},
		"individual_ports_unordered": {
			in:  datacenter.PortRanges{{First: 443, Last: 443}, {First: 80, Last: 80}, {First: 22, Last: 22}},
			exp: datacenter.PortRanges{{First: 22, Last: 22}, {First: 80, Last: 80}, {First: 443, Last: 443}},
		},
		"port_ranges_unordered": {
			in:  datacenter.PortRanges{{First: 500, Last: 600}, {First: 100, Last: 200}, {First: 300, Last: 400}},
			exp: datacenter.PortRanges{{First: 100, Last: 200}, {First: 300, Last: 400}, {First: 500, Last: 600}},
		},
		"blend_unordered": {
			in:  datacenter.PortRanges{{First: 250, Last: 250}, {First: 300, Last: 400}, {First: 100, Last: 200}},
			exp: datacenter.PortRanges{{First: 100, Last: 200}, {First: 250, Last: 250}, {First: 300, Last: 400}},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			tCase.in.Canonicalize()
			require.Equal(t, tCase.exp, tCase.in)
		})
	}
}
