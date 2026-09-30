// Copyright (c) Juniper Networks, Inc., 2022-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"strings"
	"testing"

	"github.com/Juniper/apstra-go-sdk/enum"
)

func TestResourceGroupType(t *testing.T) {
	rgTypeToGroups := make(map[string][]string)
	for _, rg := range enum.ResourceGroups.Members() {
		rgt := resourceGroupType(rg)
		if rgt == nil {
			t.Fatalf("resourceGroupType() returned nil for %T with value %q", rg, rg.Value)
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

	if len(rgTypeToGroups[enum.ResourceTypeIPv4.Value]) != len(rgTypeToGroups[enum.ResourceTypeIPv6.Value]) {
		t.Fatalf("Mismatch in quantity of IPv4 and IPv6 resource group types: %d IPv4 vs %d IPv6",
			len(rgTypeToGroups[enum.ResourceTypeIPv4.Value]), len(rgTypeToGroups[enum.ResourceTypeIPv6.Value]))
	}
}
