// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pointer

import (
	"reflect"
)

func ZeroOf[A any](a *A) A {
	var zero A
	return zero
}

func TypeOf[A any](_ *A) string {
	return reflect.TypeFor[A]().String()
}
