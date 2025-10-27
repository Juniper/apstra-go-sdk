// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package policy

import "github.com/Juniper/apstra-go-sdk/enum"

type ASNAllocation struct {
	SpineASNScheme enum.ASNAllocationScheme `json:"spine_asn_scheme"`
}
