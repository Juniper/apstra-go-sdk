// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

func defaultIfZero[T comparable](val, def T) T {
	var zero T
	if val == zero {
		return def
	}
	return val
}
