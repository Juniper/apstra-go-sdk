// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package enum

type enum interface {
	String() string
	FromString(string) error
}
