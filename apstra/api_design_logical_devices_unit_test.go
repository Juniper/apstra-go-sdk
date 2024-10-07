// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"fmt"
	"testing"
)

var testSpeedStrings = [][]string{
	{"10000000", "10M"},
	{"10M", "10M"},
	{"10Mbps", "10M"},
	{"10Mb/s", "10M"},
	{"100000000", "100M"},
	{"100M", "100M"},
	{"100Mbps", "100M"},
	{"100Mb/s", "100M"},
	{"1000000000", "1G"},
	{"1000M", "1G"},
	{"1000Mbps", "1G"},
	{"1000Mb/s", "1G"},
	{"1000000000", "1G"},
	{"10G", "10G"},
	{"10Gbps", "10G"},
	{"10Gb/s", "10G"},
	{"10000000000", "10G"},
	{"25G", "25G"},
	{"25Gbps", "25G"},
	{"25Gb/s", "25G"},
	{"25000000000", "25G"},
	{"40G", "40G"},
	{"40Gbps", "40G"},
	{"40Gb/s", "40G"},
	{"40000000000", "40G"},
	{"50G", "50G"},
	{"50Gbps", "50G"},
	{"50Gb/s", "50G"},
	{"50000000000", "50G"},
	{"100G", "100G"},
	{"100Gbps", "100G"},
	{"100Gb/s", "100G"},
	{"100000000000", "100G"},
	{"200G", "200G"},
	{"200Gbps", "200G"},
	{"200Gb/s", "200G"},
	{"200000000000", "200G"},
	{"400G", "400G"},
	{"400Gbps", "400G"},
	{"400Gb/s", "400G"},
	{"400000000000", "400G"},
}

func TestParseLogicalDeviceSpeed(t *testing.T) {
	for _, test := range testSpeedStrings {
		r := LogicalDevicePortSpeed(test[0]).raw()
		s1 := fmt.Sprintf("%d%s", r.Value, r.Unit)
		s2 := string(r.parse())
		if s1 != s2 {
			t.Fatalf("conversion problem: %s %s %s %s", test[0], test[1], s1, s2)
		}
	}
}

func TestLogicalDevicePortSpeed_IsEqual(t *testing.T) {
	for _, s := range testSpeedStrings {
		if !LogicalDevicePortSpeed(s[0]).IsEqual(LogicalDevicePortSpeed(s[1])) {
			t.Fatalf("speeds not equal %s %s", s[0], s[1])
		}
	}
}
