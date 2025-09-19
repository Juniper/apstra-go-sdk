// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package compare

import "testing"

func SlicesAsSets[A comparable](t testing.TB, a, b []A, info string) {
	t.Helper()

	if len(a) != len(b) {
		t.Fatalf("%s slice length mismatch: %d vs %d", info, len(a), len(b))
	}

	mapA := make(map[A]struct{}, len(a))
	for _, v := range a {
		mapA[v] = struct{}{}
	}

	mapB := make(map[A]struct{}, len(b))
	for _, v := range b {
		mapB[v] = struct{}{}
	}

	for k := range mapA {
		if _, ok := mapB[k]; !ok {
			t.Fatalf("%s slice contents mismatch: element %v found only in one slice", info, k)
		}
	}
}
