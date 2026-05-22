// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

type IDer interface {
	ID() *string
}

type IDSetter interface {
	SetID(string) error
}
