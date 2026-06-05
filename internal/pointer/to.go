// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pointer

func To[A any](a A) *A {
	return &a
}

func ToCopyOf[A any](a A) *A {
	_copy := a
	return &_copy
}
