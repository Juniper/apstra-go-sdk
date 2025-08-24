// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package testutils

func Range(n int) []int {
	r := make([]int, n)
	for i := range n {
		r[i] = i
	}
	return r
}
