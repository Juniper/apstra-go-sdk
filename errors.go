// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package sdk

type ErrAPIResponseInvalid string

func (e ErrAPIResponseInvalid) Error() string {
	return string(e)
}

type ErrInternal string

func (e ErrInternal) Error() string {
	return string(e)
}

type ErrMultipleMatch string

func (e ErrMultipleMatch) Error() string {
	return string(e)
}

type ErrNotFound string

func (e ErrNotFound) Error() string {
	return string(e)
}

type ErrWrongType string

func (e ErrWrongType) Error() string {
	return string(e)
}
