// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package testutils

func ZeroOf[T any](example T) T {
	var zero T
	return zero
}
