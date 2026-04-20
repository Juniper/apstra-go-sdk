// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package errors

type APIResponseInvalid string

func (e APIResponseInvalid) Error() string {
	return string(e)
}

type IDAlreadySet string

func (e IDAlreadySet) Error() string {
	return string(e)
}

type Internal string

func (e Internal) Error() string {
	return string(e)
}

type MultipleMatch string

func (e MultipleMatch) Error() string {
	return string(e)
}

type NotFound string

func (e NotFound) Error() string {
	return string(e)
}

type WrongType string

func (e WrongType) Error() string {
	return string(e)
}
