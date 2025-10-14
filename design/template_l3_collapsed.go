// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
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
	RackTypes            []RackType
	RackTypeCounts       []RackTypeCount
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
	raw := struct {
		Type                 enum.TemplateType        `json:"type"`
		Capability           enum.TemplateCapability  `json:"capability"`
		DisplayName          string                   `json:"display_name"`
		RackTypes            []RackType               `json:"rack_types"`
		RackTypeCounts       []RackTypeCount          `json:"rack_type_counts"`
		MeshLinkCount        int                      `json:"mesh_link_count"`
		MeshLinkSpeed        speed.Speed              `json:"mesh_link_speed"`
		AntiAffinityPolicy   *policy.AntiAffinity     `json:"anti_affinity_policy,omitempty"`
		DHCPServiceIntent    policy.DHCPServiceIntent `json:"dhcp_service_intent"`
		VirtualNetworkPolicy *policy.VirtualNetwork   `json:"virtual_network_policy,omitempty"`
	}{
		Type:                 t.TemplateType(),
		Capability:           enum.TemplateCapabilityBlueprint,
		DisplayName:          t.Label,
		RackTypes:            t.RackTypes,
		RackTypeCounts:       t.RackTypeCounts,
		MeshLinkCount:        t.MeshLinkCount,
		MeshLinkSpeed:        t.MeshLinkSpeed,
		AntiAffinityPolicy:   t.AntiAffinityPolicy,
		DHCPServiceIntent:    t.DHCPServiceIntent,
		VirtualNetworkPolicy: t.VirtualNetworkPolicy,
	}

	return json.Marshal(&raw)
}

func (t *TemplateL3Collapsed) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Type                 enum.TemplateType        `json:"type"`
		DisplayName          string                   `json:"display_name"`
		RackTypes            []RackType               `json:"rack_types"`
		RackTypeCounts       []RackTypeCount          `json:"rack_type_counts"`
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
	t.RackTypes = raw.RackTypes
	t.RackTypeCounts = raw.RackTypeCounts
	t.MeshLinkCount = raw.MeshLinkCount
	t.MeshLinkSpeed = raw.MeshLinkSpeed
	t.AntiAffinityPolicy = raw.AntiAffinityPolicy
	t.DHCPServiceIntent = raw.DHCPServiceIntent
	t.VirtualNetworkPolicy = raw.VirtualNetworkPolicy
	t.id = raw.ID
	t.createdAt = raw.CreatedAt
	t.lastModifiedAt = raw.LastModifiedAt

	return nil
}

func (t TemplateL3Collapsed) CreatedAt() *time.Time {
	return t.createdAt
}

func (t TemplateL3Collapsed) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}
