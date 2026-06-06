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
	_ internal.IDer     = (*SecurityZone)(nil)
	_ internal.IDSetter = (*SecurityZone)(nil)
	_ json.Unmarshaler  = (*SecurityZone)(nil)
)

type SecurityZone struct {
	Label             string                 `json:"label"`
	Description       *string                `json:"vrf_description,omitempty"`
	Type              enum.SecurityZoneType  `json:"sz_type"`
	VRFName           string                 `json:"vrf_name"`
	RoutingPolicyID   string                 `json:"routing_policy_id,omitempty"`   // automatically assigned if not set
	RouteTarget       *string                `json:"route_target,omitempty"`        // calculated only value
	RTPolicy          *RTPolicy              `json:"rt_policy,omitempty"`           // can be null
	VLAN              *uint16                `json:"vlan_id,omitempty"`             // can be null
	VNI               *int                   `json:"vni_id,omitempty"`              // can be null
	JunosEVPNIRBMode  *enum.JunosEVPNIRBMode `json:"junos_evpn_irb_mode,omitempty"` // can be null in POST, required in PUT AOS-58916
	AddressingSupport *enum.AddressingScheme `json:"addressing_support,omitempty"`  // Apstra 6.1+ only
	DisableIPv4       *bool                  `json:"disable_ipv4,omitempty"`        // Apstra 6.1+ only
	VTEPAddressing    *enum.AddressingScheme `json:"vtep_addressing,omitempty"`     // Apstra 6.1+ only
	Tags              []string               // Apstra 5.0.0+ and read-only - JSON struct tag is omitted to prevent marshaling this attribute

	id string
}

// ID returns a pointer to a copy of the object's ID, or nil when no ID is set.
func (o SecurityZone) ID() *string {
	if o.id == "" {
		return nil
	}
	id := o.id
	return &id
}

func (o *SecurityZone) SetID(id string) error {
	if o.id != "" {
		return errors.IDAlreadySet(fmt.Sprintf("id already has value %q", o.id))
	}

	o.id = id
	return nil
}

func (o *SecurityZone) UnmarshalJSON(bytes []byte) error {
	type securityZoneAlias SecurityZone // type alias prevents recursion

	var aux struct {
		securityZoneAlias
		ID   string   `json:"id"`   // the `id` struct element cannot be unmarshaled so we temporarily stash that value here
		Tags []string `json:"tags"` // the `Tags` struct element cannot be unmarshaled so we temporarily stash that value here
	}

	// unmarshal everything which can be handled by the `json` package.
	if err := json.Unmarshal(bytes, &aux); err != nil {
		return err
	}

	*o = SecurityZone(aux.securityZoneAlias)
	o.id = aux.ID
	o.Tags = aux.Tags

	return nil
}
