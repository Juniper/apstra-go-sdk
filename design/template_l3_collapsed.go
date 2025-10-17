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
	"github.com/Juniper/apstra-go-sdk/speed"
)

var (
	_ Template          = (*TemplateL3Collapsed)(nil)
	_ internal.IDSetter = (*TemplateL3Collapsed)(nil)
	_ json.Marshaler    = (*TemplateL3Collapsed)(nil)
	_ json.Unmarshaler  = (*TemplateL3Collapsed)(nil)
)

type TemplateL3Collapsed struct {
	Label                string
	Racks                []RackTypeWithCount
	MeshLinkCount        int
	MeshLinkSpeed        speed.Speed
	AntiAffinityPolicy   *policy.AntiAffinity // required for 4.2.0 only
	DHCPServiceIntent    policy.DHCPServiceIntent
	VirtualNetworkPolicy *policy.VirtualNetwork

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (t TemplateL3Collapsed) TemplateType() enum.TemplateType {
	return enum.TemplateTypeL3Collapsed
}

func (t TemplateL3Collapsed) ID() *string {
	if t.id == "" {
		return nil
	}
	return &t.id
}

// SetID sets a previously un-set id attribute. If the id attribute is found to
// have an existing value, an error is returned. Presence of an existing value
// is the only reason SetID will return an error. If the id attribute is known
// to be empty, use MustSetID.
func (t *TemplateL3Collapsed) SetID(id string) error {
	if t.id != "" {
		return sdk.ErrIDIsSet(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (t *TemplateL3Collapsed) MustSetID(id string) {
	err := t.SetID(id)
	if err != nil {
		panic(err)
	}
}

func (t TemplateL3Collapsed) MarshalJSON() ([]byte, error) {
	type rawRackTypeCount struct {
		RackTypeId string `json:"rack_type_id"`
		Count      int    `json:"count"`
	}

	raw := struct {
		Type                 enum.TemplateType        `json:"type"`
		Capability           enum.TemplateCapability  `json:"capability"`
		DisplayName          string                   `json:"display_name"`
		RackTypes            []RackType               `json:"rack_types"`
		RackTypeCounts       []rawRackTypeCount       `json:"rack_type_counts"`
		MeshLinkCount        int                      `json:"mesh_link_count"`
		MeshLinkSpeed        speed.Speed              `json:"mesh_link_speed"`
		AntiAffinityPolicy   *policy.AntiAffinity     `json:"anti_affinity_policy,omitempty"`
		DHCPServiceIntent    policy.DHCPServiceIntent `json:"dhcp_service_intent"`
		VirtualNetworkPolicy *policy.VirtualNetwork   `json:"virtual_network_policy,omitempty"`
	}{
		Type:                 t.TemplateType(),
		Capability:           enum.TemplateCapabilityBlueprint,
		DisplayName:          t.Label,
		MeshLinkCount:        t.MeshLinkCount,
		MeshLinkSpeed:        t.MeshLinkSpeed,
		AntiAffinityPolicy:   t.AntiAffinityPolicy,
		DHCPServiceIntent:    t.DHCPServiceIntent,
		VirtualNetworkPolicy: t.VirtualNetworkPolicy,
	}

	// used to generate rack type IDs within the template
	hash := md5.New()

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

func (t *TemplateL3Collapsed) UnmarshalJSON(bytes []byte) error {
	type rawRackTypeCount struct {
		RackTypeId string `json:"rack_type_id"`
		Count      int    `json:"count"`
	}

	var raw struct {
		Type                 enum.TemplateType        `json:"type"`
		DisplayName          string                   `json:"display_name"`
		RackTypes            []RackType               `json:"rack_types"`
		RackTypeCounts       []rawRackTypeCount       `json:"rack_type_counts"`
		MeshLinkCount        int                      `json:"mesh_link_count"`
		MeshLinkSpeed        speed.Speed              `json:"mesh_link_speed"`
		AntiAffinityPolicy   *policy.AntiAffinity     `json:"anti_affinity_policy,omitempty"`
		DHCPServiceIntent    policy.DHCPServiceIntent `json:"dhcp_service_intent"`
		VirtualNetworkPolicy *policy.VirtualNetwork   `json:"virtual_network_policy"`

		ID             string     `json:"id"`
		CreatedAt      *time.Time `json:"created_at"`
		LastModifiedAt *time.Time `json:"last_modified_at"`
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return fmt.Errorf("unmarshaling l3 collapsed template: %w", err)
	}

	if raw.Type != enum.TemplateTypeL3Collapsed {
		return fmt.Errorf("l3 collapsed templates must have type %q, got %q", enum.TemplateTypeL3Collapsed, raw.Type)
	}

	t.Label = raw.DisplayName
	t.MeshLinkCount = raw.MeshLinkCount
	t.MeshLinkSpeed = raw.MeshLinkSpeed
	t.AntiAffinityPolicy = raw.AntiAffinityPolicy
	t.DHCPServiceIntent = raw.DHCPServiceIntent
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

func (t TemplateL3Collapsed) CreatedAt() *time.Time {
	return t.createdAt
}

func (t TemplateL3Collapsed) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}

func NewTemplateL3Collapsed(id string) TemplateL3Collapsed {
	return TemplateL3Collapsed{id: id}
}
