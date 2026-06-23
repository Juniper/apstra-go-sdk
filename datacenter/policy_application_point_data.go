// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import "github.com/Juniper/apstra-go-sdk/enum"

type PolicyApplicationPointData struct {
	Id    string                           `json:"id"`
	Label string                           `json:"label"`
	Type  *enum.PolicyApplicationPointType `json:"type"`
}
