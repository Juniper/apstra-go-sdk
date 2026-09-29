// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package enum

import (
	"strings"
	"testing"
)

func TestResourceGroup_Type(t *testing.T) {
	rgTypeToGroups := make(map[string][]string)
	for _, rg := range ResourceGroups.Members() {
		rgt := rg.Type()
		if rgt == nil {
			t.Fatalf("ResourceGroup.Type() returned nil for %T with value %q", rg, rg.Value)
		}
		rgTypeToGroups[rgt.String()] = append(rgTypeToGroups[rgt.String()], rg.String())
	}
	for k, v := range rgTypeToGroups {
		sb := new(strings.Builder)
		for _, s := range v {
			sb.WriteString("\t" + s + "\n")
		}
		t.Logf("%q type resource groups:\n%s", k, sb.String())
	}
}
