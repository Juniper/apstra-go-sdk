// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package dctestobj

import (
	"maps"
	"math/rand"
	"slices"

	"github.com/Juniper/apstra-go-sdk/datacenter"
)

func RandomPortRanges(n int) datacenter.PortRanges {
	randomIntMap := make(map[uint16]struct{}, 2*n)
	for len(randomIntMap) < n*2 {
		var v uint16
		for v == 0 { // loop until v != 0 - we throw out the value if zero is chosen
			v = uint16(rand.Int())
		}
		randomIntMap[v] = struct{}{}
	}

	beginEnds := slices.Collect(maps.Keys(randomIntMap))
	slices.Sort(beginEnds)

	result := make(datacenter.PortRanges, 0, n)
	for i := range n {
		result = append(result, datacenter.PortRange{
			First: beginEnds[i*2],
			Last:  beginEnds[i*2+1],
		})
	}

	return result
}
