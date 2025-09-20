// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

const urlPrefix = "/api/design/"

type IDIsSet error

type IDer interface {
	ID() *string
	SetID(string) error
	MustSetID(string)
}
