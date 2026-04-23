// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import (
	"encoding/json"
	"fmt"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/errors"
	"github.com/Juniper/apstra-go-sdk/internal"
)

var (
	_ internal.IDer     = (*SwitchingZone)(nil)
	_ internal.IDSetter = (*SwitchingZone)(nil)
	_ json.Marshaler    = (*SwitchingZone)(nil)
	_ json.Unmarshaler  = (*SwitchingZone)(nil)
)

type SwitchingZone struct {
	Label             *string                              `json:"label,omitempty"`
	MACVRFDescription *string                              `json:"mac_vrf_description,omitempty"`
	MACVRFName        *string                              `json:"mac_vrf_name,omitempty"`
	MACVRFServiceType *enum.SwitchingZoneMACVRFServiceType `json:"mac_vrf_service_type,omitempty"`
	RouteTarget       *string                              `json:"route_target,omitempty"`
	Tags              []string                             `json:"tags,omitempty"`
	id                string
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

// MarshalJSON is implemented because the API is expecting an extra (mandatory) field
// `impl_type` which must always be "mac_vrf". Rather than exposing this nonsense to
// the caller, we just sneak it in here at the last minute.
func (z SwitchingZone) MarshalJSON() ([]byte, error) {
	type switchingZoneAlias SwitchingZone
	return json.Marshal(struct {
		switchingZoneAlias
		ImplementationType string `json:"impl_type"`
	}{
		switchingZoneAlias: switchingZoneAlias(z),
		ImplementationType: "mac_vrf",
	})
}

// UnmarshalJSON is implemented because we want to unmarshal the private
// struct element `id`.
func (z *SwitchingZone) UnmarshalJSON(bytes []byte) error {
	type alias SwitchingZone

	aux := struct {
		ID string `json:"id"`
		*alias
	}{
		alias: (*alias)(z),
	}

	if err := json.Unmarshal(bytes, &aux); err != nil {
		return err
	}

	z.id = aux.ID

	return nil
}
