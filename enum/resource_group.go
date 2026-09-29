// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package enum

import "strings"

func (o ResourceGroup) Type() *ResourceType {
	switch {
	case strings.HasSuffix(o.Value, "_asns"):
		return &ResourceTypeASN
	case strings.HasPrefix(o.Value, "ipv6_") || strings.HasSuffix(o.Value, "_ipv6"):
		return &ResourceTypeIPv6
	case strings.HasSuffix(o.Value, "_ips") || o.Value == ResourceGroupMLAGDomainSVIIPv4.Value || o.Value == ResourceGroupVirtualNetworkSviIPv4.Value:
		return &ResourceTypeIPv4
	case strings.HasSuffix(o.Value, "_vnis") || strings.HasSuffix(o.Value, "_vn_ids"):
		return &ResourceTypeVNI
	}

	return nil
}
