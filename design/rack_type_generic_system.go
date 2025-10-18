// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
)

var (
	_ logicalDeviceIDer                  = (*GenericSystem)(nil)
	_ internal.Replicator[GenericSystem] = (*GenericSystem)(nil)
	_ json.Marshaler                     = (*GenericSystem)(nil)
	_ json.Unmarshaler                   = (*GenericSystem)(nil)
)

type GenericSystem struct {
	ASNDomain        *enum.FeatureSwitch
	Count            int
	Label            string
	Links            []RackTypeLink
	LogicalDevice    LogicalDevice
	Loopback         *enum.FeatureSwitch
	ManagementLevel  enum.SystemManagementLevel
	PortChannelIDMax int
	PortChannelIDMin int
	Tags             []Tag
}

// logicalDeviceID returns *string representing the ID of the embedded logical
// device. If the LD ID is unset, nil is returned
func (g GenericSystem) logicalDeviceID() *string {
	if g.LogicalDevice.id == "" {
		return nil
	}
	return pointer.To(g.LogicalDevice.id)
}

// Replicate returns a copy of itself with zero values for metadata fields
func (g GenericSystem) Replicate() GenericSystem {
	result := GenericSystem{
		Count:            g.Count,
		Label:            g.Label,
		Links:            make([]RackTypeLink, len(g.Links)),
		LogicalDevice:    g.LogicalDevice.Replicate(),
		ManagementLevel:  g.ManagementLevel,
		PortChannelIDMax: g.PortChannelIDMax,
		PortChannelIDMin: g.PortChannelIDMin,
		Tags:             make([]Tag, len(g.Tags)),
		// ASNDomain:     nil,
		// Loopback:      nil,
	}

	for i, link := range g.Links {
		result.Links[i] = link.Replicate()
	}

	for i, tag := range g.Tags {
		result.Tags[i] = tag.Replicate()
	}

	if g.ASNDomain != nil {
		result.ASNDomain = pointer.To(*g.ASNDomain)
	}

	if g.Loopback != nil {
		result.Loopback = pointer.To(*g.Loopback)
	}

	return result
}

func (g GenericSystem) MarshalJSON() ([]byte, error) {
	result := rawGenericSystem{
		AsnDomain:        g.ASNDomain,
		Count:            zero.PreferDefault(g.Count, 1),
		Label:            g.Label,
		Links:            g.Links,
		LogicalDeviceID:  fmt.Sprintf("%x", mustHashForComparison(g.LogicalDevice, md5.New())),
		Loopback:         g.Loopback,
		ManagementLevel:  g.ManagementLevel,
		PortChannelIDMax: g.PortChannelIDMax,
		PortChannelIDMin: g.PortChannelIDMin,
		TagLabels:        make([]string, len(g.Tags)),
	}

	for i, tag := range g.Tags {
		result.TagLabels[i] = tag.Label
	}
	slices.Sort(result.TagLabels)

	return json.Marshal(result)
}

func (g *GenericSystem) UnmarshalJSON(bytes []byte) error {
	var raw rawGenericSystem
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling access switch: %w", err)
	}

	g.ASNDomain = raw.AsnDomain
	g.Count = raw.Count
	g.Label = raw.Label
	g.Links = raw.Links
	g.LogicalDevice = NewLogicalDevice(raw.LogicalDeviceID)
	g.Loopback = raw.Loopback
	g.ManagementLevel = raw.ManagementLevel
	g.PortChannelIDMax = raw.PortChannelIDMax
	g.PortChannelIDMin = raw.PortChannelIDMin
	g.Tags = make([]Tag, len(raw.TagLabels))
	for i, rawTagLabel := range raw.TagLabels {
		g.Tags[i].Label = rawTagLabel // tag description must be filled by the caller
	}

	return nil
}

// it is safe and reasonable to have a "raw" type for objects which:
// 1) are marshaled and unmarshaled symmetrically (have no metadata to suppress)
// 2) have JSON layout which doesn't align with their public struct layout
type rawGenericSystem struct {
	AsnDomain        *enum.FeatureSwitch        `json:"asn_domain,omitempty"`
	Count            int                        `json:"count"`
	Label            string                     `json:"label"`
	Links            []RackTypeLink             `json:"links"`
	LogicalDeviceID  string                     `json:"logical_device"`
	Loopback         *enum.FeatureSwitch        `json:"loopback,omitempty"`
	ManagementLevel  enum.SystemManagementLevel `json:"management_level"`
	PortChannelIDMax int                        `json:"port_channel_id_max"`
	PortChannelIDMin int                        `json:"port_channel_id_min"`
	TagLabels        []string                   `json:"tags"`
}
