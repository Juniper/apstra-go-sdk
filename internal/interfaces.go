// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

type IDer interface {
	ID() *string
}

type IDSetter interface {
	IDer
	SetID(string) error
	MustSetID(string)
}

type Replicator[T any] interface {
	Replicate() T
}
