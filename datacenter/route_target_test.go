// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRouteTarget_MarshalText(t *testing.T) {
	type testCase struct {
		in     RouteTarget
		exp    string
		expErr bool
	}

	mkAS2 := func(asn uint16, local uint32) RouteTarget {
		var rt RouteTarget
		rt.e = rtEncodingAS2Local4
		binary.BigEndian.PutUint16(rt.v[0:2], asn)
		binary.BigEndian.PutUint32(rt.v[2:6], local)
		return rt
	}

	mkAS4 := func(asn uint32, local uint16) RouteTarget {
		var rt RouteTarget
		rt.e = rtEncodingAS4Local2
		binary.BigEndian.PutUint32(rt.v[0:4], asn)
		binary.BigEndian.PutUint16(rt.v[4:6], local)
		return rt
	}

	mkIPv4 := func(a, b, c, d byte, local uint16) RouteTarget {
		var rt RouteTarget
		rt.e = rtEncodingIPv4Local2
		rt.v[0], rt.v[1], rt.v[2], rt.v[3] = a, b, c, d
		binary.BigEndian.PutUint16(rt.v[4:6], local)
		return rt
	}

	testCases := map[string]testCase{
		"as2_local4": {
			in:  mkAS2(65000, 123456),
			exp: "65000:123456",
		},
		"as4_local2": {
			in:  mkAS4(70000, 321),
			exp: "70000:321",
		},
		"ipv4_local2": {
			in:  mkIPv4(192, 0, 2, 1, 456),
			exp: "192.0.2.1:456",
		},
		"unknown": {
			in:     RouteTarget{},
			expErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := tc.in.MarshalText()
			if tc.expErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.exp, string(got))
		})
	}
}

func TestRouteTarget_UnmarshalText(t *testing.T) {
	type testCase struct {
		in     string
		exp    RouteTarget
		expErr bool
	}

	mkAS2 := func(asn uint16, local uint32) [6]byte {
		var v [6]byte
		binary.BigEndian.PutUint16(v[0:2], asn)
		binary.BigEndian.PutUint32(v[2:6], local)
		return v
	}

	mkAS4 := func(asn uint32, local uint16) [6]byte {
		var v [6]byte
		binary.BigEndian.PutUint32(v[0:4], asn)
		binary.BigEndian.PutUint16(v[4:6], local)
		return v
	}

	mkIPv4 := func(a, b, c, d byte, local uint16) [6]byte {
		var v [6]byte
		v[0], v[1], v[2], v[3] = a, b, c, d
		binary.BigEndian.PutUint16(v[4:6], local)
		return v
	}

	testCases := map[string]testCase{
		"as2_local4": {
			in:  "65000:123456",
			exp: RouteTarget{e: rtEncodingAS2Local4, v: mkAS2(65000, 123456)},
		},
		"as4_local2": {
			in:  "70000:321",
			exp: RouteTarget{e: rtEncodingAS4Local2, v: mkAS4(70000, 321)},
		},
		"ipv4_local2": {
			in:  "192.0.2.1:456",
			exp: RouteTarget{e: rtEncodingIPv4Local2, v: mkIPv4(192, 0, 2, 1, 456)},
		},
		"ambiguous_prefers_as2": {
			in:  "1:1",
			exp: RouteTarget{e: rtEncodingAS2Local4, v: mkAS2(1, 1)},
		},
		"bad_format": {
			in:     "1",
			expErr: true,
		},
		"bad_part0": {
			in:     "abc:1",
			expErr: true,
		},
		"bad_part1": {
			in:     "1:abc",
			expErr: true,
		},
		"bad_ipv6": {
			in:     "3fff::/64:65536",
			expErr: true,
		},
		"as4_part0_overflow": {
			in:     "4294967296:1",
			expErr: true,
		},
		"as2_part1_overflow": {
			in:     "1:4294967296",
			expErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var got RouteTarget
			err := got.UnmarshalText([]byte(tc.in))
			if tc.expErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.exp, got)
		})
	}
}

func TestRouteTarget_RoundTrip(t *testing.T) {
	inputs := []string{
		"65000:123456",
		"70000:321",
		"192.0.2.1:456",
		"1:1",
	}

	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			rt, err := NewRouteTarget(in)
			require.NoError(t, err)

			out, err := rt.MarshalText()
			require.NoError(t, err)
			require.Equal(t, in, string(out))
		})
	}
}
