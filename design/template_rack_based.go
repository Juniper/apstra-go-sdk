// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"hash"
	"time"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/policy"
)

var (
	_ Template         = (*TemplateRackBased)(nil)
	_ internal.IDer    = (*TemplateRackBased)(nil)
	_ json.Marshaler   = (*TemplateRackBased)(nil)
	_ json.Unmarshaler = (*TemplateRackBased)(nil)
)

type TemplateRackBased struct {
	Label                string
	Racks                []RackTypeWithCount
	AntiAffinityPolicy   *policy.AntiAffinity // required for 4.2.0 only
	ASNAllocationPolicy  *policy.ASNAllocation
	Capability           *enum.TemplateCapability
	DHCPServiceIntent    policy.DHCPServiceIntent
	Spine                Spine
	VirtualNetworkPolicy *policy.VirtualNetwork

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time

	skipTypeDuringMarshalJSON bool
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

func (t *TemplateRackBased) setID(id string) {
	if t.id != "" {
		panic(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return
}

// Replicate returns a copy of itself with zero values for metadata fields
func (t TemplateRackBased) Replicate() TemplateRackBased {
	result := TemplateRackBased{
		Label:                t.Label,
		Racks:                make([]RackTypeWithCount, len(t.Racks)),
		AntiAffinityPolicy:   t.AntiAffinityPolicy,
		ASNAllocationPolicy:  t.ASNAllocationPolicy,
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
	type rawRackTypeCount struct {
		ID    string `json:"rack_type_id"`
		Count int    `json:"count"`
	}

	raw := struct {
		ID                   string                   `json:"id,omitempty"` // ID must be marshaled for pod-based template embedding
		DisplayName          string                   `json:"display_name"`
		Type                 *enum.TemplateType       `json:"type,omitempty"`
		RackTypes            []RackType               `json:"rack_types"`
		RackTypeCounts       []rawRackTypeCount       `json:"rack_type_counts"`
		AntiAffinityPolicy   *policy.AntiAffinity     `json:"anti_affinity_policy,omitempty"`
		AsnAllocationPolicy  *policy.ASNAllocation    `json:"asn_allocation_policy,omitempty"`
		Capability           *enum.TemplateCapability `json:"capability,omitempty"`
		DHCPServiceIntent    policy.DHCPServiceIntent `json:"dhcp_service_intent"`
		Spine                Spine                    `json:"spine"`
		VirtualNetworkPolicy *policy.VirtualNetwork   `json:"virtual_network_policy,omitempty"`
	}{
		ID:                   t.id,
		DisplayName:          t.Label,
		RackTypes:            make([]RackType, 0, len(t.Racks)),
		RackTypeCounts:       make([]rawRackTypeCount, 0, len(t.Racks)),
		AntiAffinityPolicy:   t.AntiAffinityPolicy,
		AsnAllocationPolicy:  t.ASNAllocationPolicy,
		Capability:           t.Capability,
		DHCPServiceIntent:    t.DHCPServiceIntent,
		Spine:                t.Spine,
		VirtualNetworkPolicy: t.VirtualNetworkPolicy,
	}

	if !t.skipTypeDuringMarshalJSON {
		raw.Type = pointer.To(t.TemplateType())
	}

	// used to generate IDs within the template
	hasher := md5.New()

	// set the spine logical device ID if necessary
	if raw.Spine.LogicalDevice.ID() == nil {
		raw.Spine.LogicalDevice.setHashID(hasher)
	}

	// keep track of rack type IDs (hashes of rack data). if two rack types are
	// identical twins (have the same contents) we don't want to add them to
	// raw.RackTypes twice. we will add them to raw.RackTypeCounts twice, and
	// the Apstra API will amend the totals as needed.
	rackTypeIDs := make(map[string]struct{}, len(t.Racks))

	// loop over racks, calculate a fresh ID, count the type of each
	for _, rackTypeWithCount := range t.Racks {
		rackType := rackTypeWithCount.RackType.Replicate() // fresh copy without metadata
		rackType.setHashID(hasher)                         // assign the ID

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

func (t *TemplateRackBased) UnmarshalJSON(bytes []byte) error {
	type rawRackTypeCount struct {
		Count int    `json:"count"`
		ID    string `json:"rack_type_id"`
	}

	var raw struct {
		DisplayName          string                   `json:"display_name"`
		Type                 enum.TemplateType        `json:"type"`
		RackTypes            []RackType               `json:"rack_types"`
		RackTypeCounts       []rawRackTypeCount       `json:"rack_type_counts"`
		AntiAffinityPolicy   *policy.AntiAffinity     `json:"anti_affinity_policy"`
		AsnAllocationPolicy  *policy.ASNAllocation    `json:"asn_allocation_policy"`
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

	t.Label = raw.DisplayName
	t.Racks = make([]RackTypeWithCount, 0, len(raw.RackTypes))
	t.AntiAffinityPolicy = raw.AntiAffinityPolicy
	t.ASNAllocationPolicy = raw.AsnAllocationPolicy
	t.Capability = raw.Capability
	t.DHCPServiceIntent = raw.DHCPServiceIntent
	t.Spine = raw.Spine
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

func (t TemplateRackBased) CreatedAt() *time.Time {
	return t.createdAt
}

func (t TemplateRackBased) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}

func (t TemplateRackBased) digest(h hash.Hash) []byte {
	h.Reset()
	return mustHashForComparison(t, h)
}

func (t *TemplateRackBased) setHashID(h hash.Hash) {
	t.setID(fmt.Sprintf("%x", t.digest(h)))
}

func NewTemplateRackBased(id string) TemplateRackBased {
	return TemplateRackBased{id: id}
}
