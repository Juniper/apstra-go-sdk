// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package enum

type enum interface {
	String() string
	FromString(string) error
}
