// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import (
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

type routeTargetEncoding uint8

const (
	rtEncodingUnknown    routeTargetEncoding = iota
	rtEncodingAS2Local4                      // RFC 4360 section 3, selected when parsing text where both parts are <= 65535
	rtEncodingAS4Local2                      // RFC 5668 (4-octet AS specific extended community)
	rtEncodingIPv4Local2                     // RFC 4360 section 4
)

var (
	_ encoding.TextMarshaler   = (*RouteTarget)(nil)
	_ encoding.TextUnmarshaler = (*RouteTarget)(nil)
)

// RouteTarget represents the value outlined in RFC 4360 section 3 and section 4, and RFC 5668. The value is stored in
// a 6-byte array along with a parameter which indicates how those bytes are to be interpreted. One of:
// - 2-byte ASN followed by 4-byte numerical value
// - 4-byte ASN followed by 2-byte numerical value
// - 4-byte IPv4 address followed by 2-byte numerical value
type RouteTarget struct {
	e routeTargetEncoding
	v [6]byte
}

func (r RouteTarget) MarshalText() ([]byte, error) {
	switch r.e {
	case rtEncodingAS2Local4:
		v2 := uint64(binary.BigEndian.Uint16(r.v[0:2]))
		v4 := uint64(binary.BigEndian.Uint32(r.v[2:6]))
		return []byte(strconv.FormatUint(v2, 10) + ":" + strconv.FormatUint(v4, 10)), nil
	case rtEncodingAS4Local2:
		v4 := uint64(binary.BigEndian.Uint32(r.v[0:4]))
		v2 := uint64(binary.BigEndian.Uint16(r.v[4:6]))
		return []byte(strconv.FormatUint(v4, 10) + ":" + strconv.FormatUint(v2, 10)), nil
	case rtEncodingIPv4Local2:
		ip := net.IP(r.v[0:4])
		v2 := uint64(binary.BigEndian.Uint16(r.v[4:6]))
		return []byte(ip.String() + ":" + strconv.FormatUint(v2, 10)), nil
	case rtEncodingUnknown:
		return nil, errors.New("cannot marshal route target with unknown encoding")
	default:
		return nil, fmt.Errorf("cannot marshal route target with unhandled encoding type %d", r.e)
	}
}

// UnmarshalText parses text strings ([]byte) like "1:1", "65535:4294967295", "4294967295:65535" or "192.0.2.1:65535"
// into the RouteTarget, storing their value and encoding type. Note that in the case of RT strings where both parts
// are numerical values <= 65535, the rtEncodingAS2Local4 encoding is preferred.
func (r *RouteTarget) UnmarshalText(text []byte) error {
	if r == nil {
		return fmt.Errorf("cannot unmarshal into nil %T", r)
	}

	parts := strings.Split(string(text), ":")
	if len(parts) != 2 {
		return fmt.Errorf("route target %q has invalid format", string(text))
	}

	var e routeTargetEncoding
	var v [6]byte

	if n, err := strconv.ParseUint(parts[0], 10, 32); err == nil {
		if n <= math.MaxUint16 {
			e = rtEncodingAS2Local4
			binary.BigEndian.PutUint16(v[0:2], uint16(n))
		} else {
			e = rtEncodingAS4Local2
			binary.BigEndian.PutUint32(v[0:4], uint32(n))
		}
	} else {
		if a, err := netip.ParseAddr(parts[0]); err == nil {
			if !a.Is4() {
				return fmt.Errorf("route target %q uses unsupported address", string(text))
			}
			e = rtEncodingIPv4Local2
			*(*[4]byte)(v[:4]) = a.As4() // store IP address in first four bytes of v
		} else {
			// cannot parse part 0 as uint nor as ip address
			return fmt.Errorf("cannot parse 1st part of route target %q", string(text))
		}
	}

	switch e {
	case rtEncodingAS4Local2, rtEncodingIPv4Local2: // part 1 must be 16 bit
		n, err := strconv.ParseUint(parts[1], 10, 16)
		if err != nil {
			return fmt.Errorf("parsing 2nd part of route target %q: %w", string(text), err)
		}
		binary.BigEndian.PutUint16(v[4:6], uint16(n))
	case rtEncodingAS2Local4: // part 1 must be 32 bit
		n, err := strconv.ParseUint(parts[1], 10, 32)
		if err != nil {
			return fmt.Errorf("parsing 2nd part of route target %q: %w", string(text), err)
		}
		binary.BigEndian.PutUint32(v[2:6], uint32(n))
	default:
		return fmt.Errorf("internal error: failed parsing route target %q", string(text))
	}

	r.e = e
	r.v = v
	return nil
}

func NewRouteTarget(in string) (RouteTarget, error) {
	var result RouteTarget
	return result, result.UnmarshalText([]byte(in))
}
