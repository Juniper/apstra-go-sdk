// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package testutils

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
