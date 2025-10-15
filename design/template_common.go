// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

type RackTypeCount struct {
	RackTypeId string `json:"rack_type_id"`
	Count      int    `json:"count"`
}

type RackTypeWithCount struct {
	RackType RackType
	Count    int
}
