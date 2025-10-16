// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparepolicy

import "fmt"

func addMsg(old []string, new string, a ...any) []string {
	var prefix string
	if len(old) > 0 && old[0] != "" {
		prefix = old[0] + ": "
	}

	return []string{prefix + fmt.Sprintf(new, a...)}
}
