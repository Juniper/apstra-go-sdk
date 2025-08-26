// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build tools

package tools

import (
	// license compliance
	_ "github.com/google/go-licenses/v2"

	// opinionated code formatting
	_ "mvdan.cc/gofumpt"
)
