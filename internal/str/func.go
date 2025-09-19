// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package str

import "runtime"

func FuncName() string {
	pc, _, _, ok := runtime.Caller(1) // step back one layer to get the caller
	if !ok {
		return "unknown"
	}
	return runtime.FuncForPC(pc).Name()
}
