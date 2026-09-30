// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package iresources

import (
	"strings"

	"github.com/Juniper/apstra-go-sdk/enum"
)

func GroupType(rg enum.ResourceGroup) *enum.ResourceType {
	switch { // Order matters: IPv6 must be checked before IPv4
	case strings.HasSuffix(rg.Value, "_asns"):
		return &enum.ResourceTypeASN
	case strings.HasPrefix(rg.Value, "ipv6_") || strings.HasSuffix(rg.Value, "_ipv6"):
		return &enum.ResourceTypeIPv6
	case strings.HasSuffix(rg.Value, "_ips") || strings.HasSuffix(rg.Value, "_subnets"):
		return &enum.ResourceTypeIPv4
	case strings.HasSuffix(rg.Value, "_vnis") || strings.HasSuffix(rg.Value, "_vn_ids"):
		return &enum.ResourceTypeVNI
	}

	return nil
}
