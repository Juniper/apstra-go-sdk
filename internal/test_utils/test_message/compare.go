// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package testmessage

import "fmt"

func Add(old []string, new string, a ...any) []string {
	var prefix string
	if len(old) > 0 && old[0] != "" {
		prefix = old[0] + ": "
	}

	return []string{prefix + fmt.Sprintf(new, a...)}
}
