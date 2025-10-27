// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"github.com/Juniper/apstra-go-sdk/speed"
)

type RackTypeWithCount struct {
	Count    int
	RackType RackType
}

type PodWithCount struct {
	Count int
	Pod   TemplateRackBased
}

type Spine struct {
	Count                  int           `json:"count"`
	LinkPerSuperspineCount int           `json:"link_per_superspine_count"`
	LinkPerSuperspineSpeed speed.Speed   `json:"link_per_superspine_speed"`
	LogicalDevice          LogicalDevice `json:"logical_device"`
	Tags                   []Tag         `json:"tags"`
}

// Replicate returns a copy of itself with zero values for metadata fields
func (s Spine) Replicate() Spine {
	return Spine{
		Count:                  s.Count,
		LinkPerSuperspineCount: s.LinkPerSuperspineCount,
		LinkPerSuperspineSpeed: s.LinkPerSuperspineSpeed,
		LogicalDevice:          s.LogicalDevice.Replicate(),
		Tags:                   s.Tags,
	}
}

type Superspine struct {
	PlaneCount         int           `json:"plane_count"`
	SuperspinePerPlane int           `json:"superspine_per_plane"`
	LogicalDevice      LogicalDevice `json:"logical_device"`
	Tags               []Tag         `json:"tags"`
}
