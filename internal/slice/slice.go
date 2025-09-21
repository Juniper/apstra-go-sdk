// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package slice

type IDer interface {
	ID() *string
}

func ObjectWithID[T IDer](elements []T, id string) *T {
	for _, element := range elements {
		idPtr := element.ID()
		if idPtr != nil && id == *idPtr {
			return &element
		}
	}
	return nil
}
