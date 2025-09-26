// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package slice

// IDer is an interface for types that expose a string ID via the ID method.
// A nil return is reserved for cases where an ID has yet to be established.
type IDer interface {
	ID() *string
}

// ObjectWithID searches the given slice for the first element whose ID method
// returns a non-nil pointer equal to the specified id string.
//
// It returns a pointer to the matching element, or nil if no match is found.
func ObjectWithID[T IDer](elements []T, id string) *T {
	for i, element := range elements {
		idPtr := element.ID()
		if idPtr != nil && id == *idPtr {
			return &elements[i]
		}
	}
	return nil
}

// Uniq returns a new slice containing only the unique elements of the input
// slice `in`, preserving the order of their first appearance.
//
// The input elements must be of a comparable type.
func Uniq[T comparable](in []T) []T {
	if len(in) == 0 {
		if in == nil {
			return ([]T)(nil)
		}
		return make([]T, 0)
	}

	seen := make(map[T]struct{}, len(in))
	var result []T
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

// IsUniq reports whether all elements in the input slice `in` are unique.
//
// The input elements must be of a comparable type.
func IsUniq[T comparable](in []T) bool {
	return len(Uniq(in)) == len(in)
}
