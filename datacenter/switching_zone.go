// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import (
	"fmt"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/errors"
)

type SwitchingZone struct {
	ImplementationType *enum.SwitchingZoneImplementationType `json:"impl_type,omitempty"`
	Label              *string                               `json:"label,omitempty"`
	MACVRFDescription  *string                               `json:"mac_vrf_description,omitempty"`
	MACVRFName         *string                               `json:"mac_vrf_name,omitempty"`
	MACVRFServiceType  *enum.SwitchingZoneMACVRFServiceType  `json:"mac_vrf_service_type,omitempty"`
	RouteTarget        *string                               `json:"route_target,omitempty"`
	Tags               []string                              `json:"tags,omitempty"`
	id                 string
}

// ID returns a pointer to a copy of the object's ID, or nil when no ID is set.
func (z SwitchingZone) ID() *string {
	if z.id == "" {
		return nil
	}
	id := z.id
	return &id
}

func (z *SwitchingZone) SetID(id string) error {
	if z.id != "" {
		return errors.IDAlreadySet(fmt.Sprintf("id already has value %q", z.id))
	}

	z.id = id
	return nil
}
