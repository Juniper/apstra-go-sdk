// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package enum

import "fmt"

type ErrorType int

const (
	errorTypeUnknown = ErrorType(iota)
	ErrorTypeParsingFailed
)

type Error struct {
	errType     ErrorType
	stringVal   string
	parentError error
	typeName    string
}

func (o Error) Error() string {
	switch o.errType {
	case ErrorTypeParsingFailed:
		return fmt.Sprintf("failed to parse %s %q", o.typeName, o.stringVal)
	default:
		return o.parentError.Error()
	}
}

func (o Error) Type() ErrorType {
	return o.errType
}

func newEnumParseError(e enum, s string) Error {
	return Error{
		errType:   ErrorTypeParsingFailed,
		stringVal: s,
		typeName:  fmt.Sprintf("%T", e),
	}
}
