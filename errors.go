// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra_go_sdk

type ErrIDIsSet string

func (e ErrIDIsSet) Error() string {
	return string(e)
}

type ErrNotFound string

func (e ErrNotFound) Error() string {
	return string(e)
}

type ErrMultipleMatch string

func (e ErrMultipleMatch) Error() string {
	return string(e)
}
