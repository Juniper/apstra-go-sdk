// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package slice

import "github.com/Juniper/apstra-go-sdk/internal"

// ObjectWithID searches the given slice for the first element whose ID method
// returns a non-nil pointer equal to the specified id string.
//
// It returns a pointer to the matching element, or nil if no match is found.
func ObjectWithID[T internal.IDer](elements []T, id string) *T {
	for i, element := range elements {
		idPtr := element.ID()
		if idPtr != nil && id == *idPtr {
			return &elements[i]
		}
	}
	return nil
}

// Remove removes the item at index i from slice s and returns the modified
// slice. Order is not preserved. If i is invalid (out of range) it panics.
func Remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
