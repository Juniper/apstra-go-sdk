// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package timeutils

import (
	"time"
)

func TimeParseMust(layout string, value string) time.Time {
	result, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return result
}
