// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package slice

import (
	"fmt"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/internal"
)

// FindByID searches the given slice for elements where the element's ID()
// method returns a non-nil pointer to a string which matches the passed id.
//
// If no match, both the returned pointer and the error are nil.
// If exactly one match is found, it returns a pointer to the matching element.
// If more than one match is found, an error is returned.
func FindByID[T internal.IDer](elements []T, id string) (*T, error) {
	var result *T
	for i, element := range elements {
		idPtr := element.ID()
		if idPtr != nil && id == *idPtr {
			if result == nil {
				result = &elements[i]
			} else {
				return nil, sdk.ErrMultipleMatch(fmt.Sprintf("found multiple elements with ID: %s", id))
			}
		}
	}
	return result, nil
}

// MustFindByID searches the given slice for elements where the element's ID()
// method returns a non-nil pointer to a string which matches the passed id.
//
// If no match, both the returned pointer and the error are nil.
// If exactly one match is found, it returns a pointer to the matching element.
// If more than one match is found, it panics.
func MustFindByID[T internal.IDer](elements []T, id string) *T {
	result, err := FindByID[T](elements, id)
	if err != nil {
		panic(err)
	}
	return result
}

// Remove removes the item at index i from slice s and returns the modified
// slice. Order is not preserved. If i is invalid (out of range) it panics.
func Remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
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
