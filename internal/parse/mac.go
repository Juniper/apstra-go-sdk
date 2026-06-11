// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package parse

import (
	"fmt"
	"net"
)

// MACFromString is an improvement on calling net.ParseMAC() directly because it
// handles empty strings gracefully (returns nil net.HardwareAddr)
func MACFromString(s string) (net.HardwareAddr, error) {
	if s == "" {
		return nil, nil
	}

	mac, err := net.ParseMAC(s)
	if err != nil {
		return nil, fmt.Errorf("cannot parse hardware address %q", s)
	}

	return mac, nil
}
