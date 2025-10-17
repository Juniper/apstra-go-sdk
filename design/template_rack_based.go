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
	type rawRackTypeCount struct {
		Count      int    `json:"count"`
		RackTypeId string `json:"rack_type_id"`
	}

	raw := struct {
		ID                   string                   `json:"id,omitempty"` // ID must be marshaled for pod-based template embedding
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
		ID:                   t.id,
		DisplayName:          t.Label,
		Type:                 t.TemplateType(),
		RackTypes:            make([]RackType, 0, len(t.Racks)),
		RackTypeCounts:       make([]rawRackTypeCount, 0, len(t.Racks)),
		AntiAffinityPolicy:   t.AntiAffinityPolicy,
		AsnAllocationPolicy:  t.AsnAllocationPolicy,
		Capability:           t.Capability,
		DHCPServiceIntent:    t.DHCPServiceIntent,
		Spine:                t.Spine,
		VirtualNetworkPolicy: t.VirtualNetworkPolicy,
	}

	// used to generate IDs within the template
	hasher := md5.New()

	// set the spine logical device ID if necessary
	if raw.Spine.LogicalDevice.ID() == nil {
		raw.Spine.LogicalDevice.mustSetHashID(hasher)
	}

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
		raw.RackTypeCounts = append(raw.RackTypeCounts, rawRackTypeCount{Count: rackTypeWithCount.Count, RackTypeId: rackType.id})

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
		Count      int    `json:"count"`
		RackTypeId string `json:"rack_type_id"`
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
			rackType.id = "" // we don't want the ID of the embedded rack type
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

func (t TemplateRackBased) digest(h hash.Hash) []byte {
	h.Reset()
	return mustHashForComparison(t, h)
}

func (t *TemplateRackBased) setHashID(h hash.Hash) error {
	return t.SetID(fmt.Sprintf("%x", t.digest(h)))
}

func (t *TemplateRackBased) mustSetHashID(h hash.Hash) {
	err := t.SetID(fmt.Sprintf("%x", t.digest(h)))
	if err != nil {
		panic(err)
	}
}

func NewTemplateRackBased(id string) TemplateRackBased {
	return TemplateRackBased{id: id}
}
