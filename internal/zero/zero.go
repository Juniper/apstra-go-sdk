// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package zero

func PreferDefault[T comparable](value, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}

func Of[T any](_ T) T {
	var zero T
	return zero
}

func SliceItem[T any](_ []T) T {
	var zero T
	return zero
}
