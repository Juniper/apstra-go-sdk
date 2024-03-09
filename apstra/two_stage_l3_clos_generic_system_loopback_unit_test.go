package apstra

import (
	"encoding/json"
	"net"
	"strings"
	"testing"
)

func TestGenericSystemLoopback_MarshalJSON(t *testing.T) {
	type testCase struct {
		val GenericSystemLoopback
		exp string
		err string
	}

	testCases := map[string]testCase{
		"both_nil": {
			val: GenericSystemLoopback{},
			exp: `{"ipv4_addr":null,"ipv6_addr":null}`,
		},
		"v4_only": {
			val: GenericSystemLoopback{
				Ipv4Addr: &net.IPNet{
					IP:   []byte{1, 1, 1, 1},
					Mask: net.CIDRMask(32, 32),
				},
			},
			exp: `{"ipv4_addr":"1.1.1.1/32","ipv6_addr":null}`,
		},
		"v6_only": {
			val: GenericSystemLoopback{
				Ipv6Addr: &net.IPNet{
					IP:   net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					Mask: net.CIDRMask(128, 128),
				},
			},
			exp: `{"ipv4_addr":null,"ipv6_addr":"2001:db8::1/128"}`,
		},
		"both": {
			val: GenericSystemLoopback{
				Ipv4Addr: &net.IPNet{
					IP:   net.IP{1, 1, 1, 1},
					Mask: net.CIDRMask(32, 32),
				},
				Ipv6Addr: &net.IPNet{
					IP:   net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					Mask: net.CIDRMask(128, 128),
				},
			},
			exp: `{"ipv4_addr":"1.1.1.1/32","ipv6_addr":"2001:db8::1/128"}`,
		},
		"v4_discontiguous": {
			val: GenericSystemLoopback{
				Ipv4Addr: &net.IPNet{
					IP:   net.IP{1, 1, 1, 1},
					Mask: net.IPv4Mask(255, 255, 255, 127),
				},
			},
			err: "IPNet object may have discontiguous bitmask: 1.1.1.1/ffffff7f",
		},
		"v6_discontiguous": {
			val: GenericSystemLoopback{
				Ipv6Addr: &net.IPNet{
					IP:   net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 127},
				},
			},
			err: "IPNet object may have discontiguous bitmask: 2001:db8::1/ffffffffffffffffffffffffffffff7f",
		},
		"v4_short_addr": {
			val: GenericSystemLoopback{
				Ipv4Addr: &net.IPNet{
					IP:   net.IP{1, 1, 1},
					Mask: net.CIDRMask(32, 32),
				},
			},
			err: "IPNet object does not render to string: <nil>",
		},
		"v6_short_addr": {
			val: GenericSystemLoopback{
				Ipv6Addr: &net.IPNet{
					IP:   net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, // 15 bytes
					Mask: net.CIDRMask(128, 128),
				},
			},
			err: "IPNet object does not render to string: <nil>",
		},
		"v4_long_addr": {
			val: GenericSystemLoopback{
				Ipv4Addr: &net.IPNet{
					IP:   net.IP{1, 1, 1, 1, 1},
					Mask: net.IPMask{255, 255, 255, 255},
				},
			},
			err: "IPNet object does not render to string: <nil>",
		},
		"v6_long_addr": {
			val: GenericSystemLoopback{
				Ipv6Addr: &net.IPNet{
					IP:   net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, // 17 bytes
					Mask: net.CIDRMask(128, 128),
				},
			},
			err: "IPNet object does not render to string: <nil>",
		},
		"v4_short_mask": {
			val: GenericSystemLoopback{
				Ipv4Addr: &net.IPNet{
					IP:   net.IP{1, 1, 1, 1},
					Mask: net.IPMask{255, 255, 255},
				},
			},
			err: "IPNet object does not render to string: <nil>",
		},
		"v6_short_mask": {
			val: GenericSystemLoopback{
				Ipv6Addr: &net.IPNet{
					IP:   net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, // 15 bytes
				},
			},
			err: "IPNet object does not render to string: <nil>",
		},
		"v4_long_mask": {
			val: GenericSystemLoopback{
				Ipv4Addr: &net.IPNet{
					IP:   net.IP{1, 1, 1, 1},
					Mask: net.IPMask{255, 255, 255, 255, 255},
				},
			},
			err: "IPNet object does not render to string: <nil>",
		},
		"v6_long_mask": {
			val: GenericSystemLoopback{
				Ipv6Addr: &net.IPNet{
					IP:   net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					Mask: net.CIDRMask(128, 128),
				},
			},
			err: "IPNet object does not render to string: <nil>",
		},
	}

	for tName, tCase := range testCases {
		tName, tCase := tName, tCase
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			result, err := json.Marshal(tCase.val)
			if len(tCase.err) == 0 && err != nil {
				t.Fatal(err)
			}
			if len(tCase.err) > 0 {
				if err == nil {
					t.Fatalf("test case %s expected error %s, got none", tName, tCase.err)
				}
				if !strings.Contains(err.Error(), tCase.err) {
					t.Fatalf("expected error: %q, got: %q", tCase.err, err.Error())
				}
			} else {
				if tCase.exp != string(result) {
					t.Fatalf("expected:\n%s\n\ngot:\n%s\n\n", tCase.exp, string(result))
				}
				//if tCase.val.Ipv4Addr != nil {
				//	t.Log(tCase.val.Ipv4Addr.String())
				//}
				//if tCase.val.Ipv6Addr != nil {
				//	t.Log(tCase.val.Ipv6Addr.String())
				//}
				//t.Log(string(result))
			}
		})
	}
}
