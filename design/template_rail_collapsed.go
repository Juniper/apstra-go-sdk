// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/policy"
)

var (
	_ Template          = (*TemplateRailCollapsed)(nil)
	_ internal.IDSetter = (*TemplateRailCollapsed)(nil)
	_ json.Marshaler    = (*TemplateRailCollapsed)(nil)
	_ json.Unmarshaler  = (*TemplateRailCollapsed)(nil)
)

type TemplateRailCollapsed struct {
	Label                string
	Racks                []RackTypeWithCount
	DHCPServiceIntent    policy.DHCPServiceIntent
	VirtualNetworkPolicy *policy.VirtualNetwork

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (t TemplateRailCollapsed) TemplateType() enum.TemplateType {
	return enum.TemplateTypeRailCollapsed
}

func (t TemplateRailCollapsed) ID() *string {
	if t.id == "" {
		return nil
	}
	return &t.id
}

// SetID sets a previously un-set id attribute. If the id attribute is found to
// have an existing value, an error is returned. Presence of an existing value
// is the only reason SetID will return an error. If the id attribute is known
// to be empty, use MustSetID.
func (t *TemplateRailCollapsed) SetID(id string) error {
	if t.id != "" {
		return sdk.ErrIDIsSet(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (t *TemplateRailCollapsed) MustSetID(id string) {
	err := t.SetID(id)
	if err != nil {
		panic(err)
	}
}

func (t TemplateRailCollapsed) MarshalJSON() ([]byte, error) {
	type rawRackTypeCount struct {
		Count int    `json:"count"`
		ID    string `json:"rack_type_id"`
	}

	raw := struct {
		Type                 enum.TemplateType        `json:"type"`
		Capability           enum.TemplateCapability  `json:"capability"`
		DisplayName          string                   `json:"display_name"`
		RackTypes            []RackType               `json:"rack_types"`
		RackTypeCounts       []rawRackTypeCount       `json:"rack_type_counts"`
		DHCPServiceIntent    policy.DHCPServiceIntent `json:"dhcp_service_intent"`
		VirtualNetworkPolicy *policy.VirtualNetwork   `json:"virtual_network_policy,omitempty"`
	}{
		Type:                 t.TemplateType(),
		Capability:           enum.TemplateCapabilityBlueprint,
		DisplayName:          t.Label,
		DHCPServiceIntent:    t.DHCPServiceIntent,
		VirtualNetworkPolicy: t.VirtualNetworkPolicy,
	}

	// used to generate rack type IDs within the template
	hasher := md5.New()

	// Note that the looping over and duplicate handling of raw.RackTypes and
	// raw.RackTypeCounts shouldn't be necessary for a collapsed template which
	// only allows one rack. This strategy is required for the other template
	// types and who knows, a future L3 collapsed template may support more than
	// one rack, so maybe it'll be helpful.

	// keep track of rack type IDs (hashes of rack data). if two rack types are
	// identical twins (have the same contents) we don't want to add them to
	// raw.RackTypes twice. we will add them to raw.RackTypeCounts twice, and
	// the Apstra API will amend the totals as needed.
	rackTypeIDs := make(map[string]struct{}, len(t.Racks))

	// loop over racks, calculate a fresh ID, count the type of each
	for _, rackTypeWithCount := range t.Racks {
		rackType := rackTypeWithCount.RackType.Replicate() // fresh copy without metadata
		rackType.mustSetHashID(hasher)                     // assign the ID

		// add an entry to raw.RackTypeCounts without regard to twins
		raw.RackTypeCounts = append(raw.RackTypeCounts, rawRackTypeCount{Count: rackTypeWithCount.Count, ID: rackType.id})

		// add an entry to raw.RackTypes only if it's not a twin
		if _, ok := rackTypeIDs[rackType.id]; !ok {
			rackTypeIDs[rackType.id] = struct{}{}
			raw.RackTypes = append(raw.RackTypes, rackType)
		}
	}

	return json.Marshal(&raw)
}

func (t *TemplateRailCollapsed) UnmarshalJSON(bytes []byte) error {
	type rawRackTypeCount struct {
		Count int    `json:"count"`
		ID    string `json:"rack_type_id"`
	}

	var raw struct {
		Type                 enum.TemplateType        `json:"type"`
		DisplayName          string                   `json:"display_name"`
		RackTypes            []RackType               `json:"rack_types"`
		RackTypeCounts       []rawRackTypeCount       `json:"rack_type_counts"`
		DHCPServiceIntent    policy.DHCPServiceIntent `json:"dhcp_service_intent"`
		VirtualNetworkPolicy *policy.VirtualNetwork   `json:"virtual_network_policy"`

		ID             string     `json:"id"`
		CreatedAt      *time.Time `json:"created_at"`
		LastModifiedAt *time.Time `json:"last_modified_at"`
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return fmt.Errorf("unmarshaling l3 collapsed template: %w", err)
	}

	t.Label = raw.DisplayName
	t.Racks = make([]RackTypeWithCount, 0, len(raw.RackTypes))
	t.DHCPServiceIntent = raw.DHCPServiceIntent
	t.VirtualNetworkPolicy = raw.VirtualNetworkPolicy
	t.id = raw.ID
	t.createdAt = raw.CreatedAt
	t.lastModifiedAt = raw.LastModifiedAt

	idToCount := make(map[string]int, len(raw.RackTypeCounts))
	for _, v := range raw.RackTypeCounts {
		idToCount[v.ID] = v.Count
	}

	for _, v := range raw.RackTypes {
		count, ok := idToCount[v.id]
		if !ok {
			return sdk.ErrAPIResponseInvalid(fmt.Sprintf("rack type id %q has no associated count", v.id))
		}

		t.Racks = append(t.Racks, RackTypeWithCount{Count: count, RackType: v})
	}

	return nil
}

func (t TemplateRailCollapsed) CreatedAt() *time.Time {
	return t.createdAt
}

func (t TemplateRailCollapsed) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}

func NewTemplateRailCollapsed(id string) TemplateRailCollapsed {
	return TemplateRailCollapsed{id: id}
}
