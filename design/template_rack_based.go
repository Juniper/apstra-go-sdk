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
	_ Template                               = (*TemplateRackBased)(nil)
	_ internal.IDSetter                      = (*TemplateRackBased)(nil)
	_ internal.Replicator[TemplateRackBased] = (*TemplateRackBased)(nil)
	_ json.Marshaler                         = (*TemplateRackBased)(nil)
	_ json.Unmarshaler                       = (*TemplateRackBased)(nil)
)

type TemplateRackBased struct {
	Label                string
	Racks                []RackTypeWithCount
	AntiAffinityPolicy   *policy.AntiAffinity // required for 4.2.0 only
	AsnAllocationPolicy  *AsnAllocationPolicy
	Capability           *enum.TemplateCapability
	DHCPServiceIntent    policy.DHCPServiceIntent
	Spine                Spine
	VirtualNetworkPolicy *policy.VirtualNetwork

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (t TemplateRackBased) TemplateType() enum.TemplateType {
	return enum.TemplateTypeRackBased
}

func (t TemplateRackBased) ID() *string {
	if t.id == "" {
		return nil
	}
	return &t.id
}

// SetID sets a previously un-set id attribute. If the id attribute is found to
// have an existing value, an error is returned. Presence of an existing value
// is the only reason SetID will return an error. If the id attribute is known
// to be empty, use MustSetID.
func (t *TemplateRackBased) SetID(id string) error {
	if t.id != "" {
		return sdk.ErrIDIsSet(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (t *TemplateRackBased) MustSetID(id string) {
	err := t.SetID(id)
	if err != nil {
		panic(err)
	}
}

// Replicate returns a copy of itself with zero values for metadata fields
func (t TemplateRackBased) Replicate() TemplateRackBased {
	result := TemplateRackBased{
		Label:                t.Label,
		Racks:                make([]RackTypeWithCount, len(t.Racks)),
		AntiAffinityPolicy:   t.AntiAffinityPolicy,
		AsnAllocationPolicy:  t.AsnAllocationPolicy,
		Capability:           t.Capability,
		DHCPServiceIntent:    t.DHCPServiceIntent,
		Spine:                t.Spine.Replicate(),
		VirtualNetworkPolicy: t.VirtualNetworkPolicy,
	}

	for i, rack := range t.Racks {
		result.Racks[i] = RackTypeWithCount{
			Count:    rack.Count,
			RackType: rack.RackType.Replicate(),
		}
	}

	return result
}

func (t TemplateRackBased) MarshalJSON() ([]byte, error) {
	// todo: set hash id for spine logical device
	type rawRackTypeCount struct {
		RackTypeId string `json:"rack_type_id"`
		Count      int    `json:"count"`
	}

	raw := struct {
		DisplayName          string                   `json:"display_name"`
		Type                 enum.TemplateType        `json:"type"`
		RackTypes            []RackType               `json:"rack_types"`
		RackTypeCounts       []rawRackTypeCount       `json:"rack_type_counts"`
		AntiAffinityPolicy   *policy.AntiAffinity     `json:"anti_affinity_policy,omitempty"`
		AsnAllocationPolicy  *AsnAllocationPolicy     `json:"asn_allocation_policy,omitempty"`
		Capability           *enum.TemplateCapability `json:"capability,omitempty"`
		DHCPServiceIntent    policy.DHCPServiceIntent `json:"dhcp_service_intent"`
		Spine                Spine                    `json:"spine"`
		VirtualNetworkPolicy *policy.VirtualNetwork   `json:"virtual_network_policy,omitempty"`
	}{
		DisplayName:          t.Label,
		Type:                 t.TemplateType(),
		AntiAffinityPolicy:   t.AntiAffinityPolicy,
		AsnAllocationPolicy:  t.AsnAllocationPolicy,
		Capability:           t.Capability,
		DHCPServiceIntent:    t.DHCPServiceIntent,
		Spine:                t.Spine,
		VirtualNetworkPolicy: t.VirtualNetworkPolicy,
	}

	// used to generate IDs within the template
	hash := md5.New()

	// set the spine ID if necessary
	if raw.Spine.LogicalDevice.ID() == nil {
		raw.Spine.LogicalDevice.mustSetHashID(hash)
	}

	// used to keep track of rack type quantity by ID in case we have identical racks
	rackTypeIDToCount := make(map[string]int, len(t.Racks))

	// initialize raw.RackTypeCounts so we can append to it without shuffling memory around
	raw.RackTypeCounts = make([]rawRackTypeCount, 0, len(t.Racks))

	// loop over racks, calculate a fresh ID, count the type of each
	for _, rack := range t.Racks {
		rackType := rack.RackType.Replicate() // fresh copy without metadata
		rackType.mustSetHashID(hash)          // assign the ID
		if _, ok := rackTypeIDToCount[rackType.id]; !ok {
			raw.RackTypes = append(raw.RackTypes, rackType) // previously unseen rack type - append it to the slice
		}
		rackTypeIDToCount[rackType.id] += rack.Count // adjust the quantity for this rack type
	}

	// prepare raw.RackTypeCounts from rackTypeIDToCount
	raw.RackTypeCounts = make([]rawRackTypeCount, 0, len(rackTypeIDToCount))
	for id, count := range rackTypeIDToCount {
		raw.RackTypeCounts = append(raw.RackTypeCounts, rawRackTypeCount{RackTypeId: id, Count: count})
	}

	return json.Marshal(&raw)
}

func (t *TemplateRackBased) UnmarshalJSON(bytes []byte) error {
	type rawRackTypeCount struct {
		RackTypeId string `json:"rack_type_id"`
		Count      int    `json:"count"`
	}

	var raw struct {
		DisplayName          string                   `json:"display_name"`
		Type                 enum.TemplateType        `json:"type"`
		RackTypes            []RackType               `json:"rack_types"`
		RackTypeCounts       []rawRackTypeCount       `json:"rack_type_counts"`
		AntiAffinityPolicy   *policy.AntiAffinity     `json:"anti_affinity_policy"`
		AsnAllocationPolicy  *AsnAllocationPolicy     `json:"asn_allocation_policy"`
		Capability           *enum.TemplateCapability `json:"capability"`
		DHCPServiceIntent    policy.DHCPServiceIntent `json:"dhcp_service_intent"`
		Spine                Spine                    `json:"spine"`
		VirtualNetworkPolicy *policy.VirtualNetwork   `json:"virtual_network_policy"`

		ID             string     `json:"id"`
		CreatedAt      *time.Time `json:"created_at"`
		LastModifiedAt *time.Time `json:"last_modified_at"`
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return fmt.Errorf("unmarshaling rack based template: %w", err)
	}

	if raw.Type != enum.TemplateTypeRackBased {
		return fmt.Errorf("rack based templates must have type %q, got %q", enum.TemplateTypeRackBased, raw.Type)
	}

	t.Label = raw.DisplayName
	t.AntiAffinityPolicy = raw.AntiAffinityPolicy
	t.AsnAllocationPolicy = raw.AsnAllocationPolicy
	t.Capability = raw.Capability
	t.DHCPServiceIntent = raw.DHCPServiceIntent
	t.Spine = raw.Spine
	t.VirtualNetworkPolicy = raw.VirtualNetworkPolicy
	t.id = raw.ID
	t.createdAt = raw.CreatedAt
	t.lastModifiedAt = raw.LastModifiedAt

	idToRackType := make(map[string]RackType, len(raw.RackTypes))
	for _, rackType := range raw.RackTypes {
		idToRackType[rackType.id] = rackType
	}

	t.Racks = make([]RackTypeWithCount, len(raw.RackTypeCounts))
	for i, rackTypeCount := range raw.RackTypeCounts {
		if rackType, ok := idToRackType[rackTypeCount.RackTypeId]; ok {
			t.Racks[i] = RackTypeWithCount{RackType: rackType, Count: rackTypeCount.Count}
			continue
		}

		// we should not get here
		return sdk.ErrAPIResponseInvalid(fmt.Sprintf("payload specifies %d instances of rack type %q which does not exist", rackTypeCount.Count, rackTypeCount.RackTypeId))
	}

	return nil
}

func (t TemplateRackBased) CreatedAt() *time.Time {
	return t.createdAt
}

func (t TemplateRackBased) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}
