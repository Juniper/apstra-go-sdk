// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package parse

import (
	"fmt"
	"net"
)

// IPFromString is an improvement on calling net.ParseIP() directly because it
// handles empty strings gracefully (returns nil net.IP) and because it returns
// errors in case of un-parseable input strings.
func IPFromString(s string) (net.IP, error) {
	if s == "" {
		return nil, nil
	}

	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("cannot parse IP %q", s)
	}

	return ip, nil
}

// IPNetFromString is an improvement on calling net.ParseCIDR() directly because
// it handles empty strings gracefully (returns nil pointer) and because it
// returns a net.IPNet with the actual IP address, rather than the base address.
func IPNetFromString(s string) (*net.IPNet, error) {
	if s == "" {
		return nil, nil
	}

	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, fmt.Errorf("while parsing CIDR string %q - %w", s, err)
	}
	ipNet.IP = ip

	return ipNet, nil
}
