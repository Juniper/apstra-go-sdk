// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package timeutils

import "time"

type Stamper interface {
	CreatedAt() *time.Time
	LastModifiedAt() *time.Time
}
