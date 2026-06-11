// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package deepcopy

import (
	"net"
	"slices"
)

func cloneIPNet(in *net.IPNet) *net.IPNet {
	if in == nil {
		return nil
	}

	var out net.IPNet
	out.IP = slices.Clone(in.IP)
	out.Mask = slices.Clone(in.Mask)
	return &out
}
