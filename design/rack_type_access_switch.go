// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
	"github.com/Juniper/apstra-go-sdk/speed"
)

var (
	_ logicalDeviceIDer = (*AccessSwitch)(nil)
	_ json.Marshaler    = (*AccessSwitch)(nil)
	_ json.Unmarshaler  = (*AccessSwitch)(nil)
)

type AccessSwitch struct {
	Count         int
	ESILAGInfo    *RackTypeAccessSwitchESILAGInfo
	Label         string
	Links         []RackTypeLink
	LogicalDevice LogicalDevice
	Tags          []Tag
}

// logicalDeviceID returns *string representing the ID of the embedded logical
// device. If the LD ID is unset, nil is returned
func (a AccessSwitch) logicalDeviceID() *string {
	if a.LogicalDevice.id == "" {
		return nil
	}
	return pointer.To(a.LogicalDevice.id)
}

// Replicate returns a copy of itself with zero values for metadata fields
func (a AccessSwitch) Replicate() AccessSwitch {
	result := AccessSwitch{
		Count:         a.Count,
		Label:         a.Label,
		Links:         make([]RackTypeLink, len(a.Links)),
		LogicalDevice: a.LogicalDevice.Replicate(),
		Tags:          make([]Tag, len(a.Tags)),
		// ESILAGInfo: nil,
	}

	if a.ESILAGInfo != nil {
		result.ESILAGInfo = pointer.To(a.ESILAGInfo.Replicate())
	}

	for i, link := range a.Links {
		result.Links[i] = link.Replicate()
	}

	for i, tag := range a.Tags {
		result.Tags[i] = tag.Replicate()
	}

	return result
}

func (a AccessSwitch) MarshalJSON() ([]byte, error) {
	result := rawAccessSwitch{
		Count:           zero.PreferDefault(a.Count, 1),
		Label:           a.Label,
		Links:           a.Links,
		LogicalDeviceID: fmt.Sprintf("%x", mustHashForComparison(a.LogicalDevice, md5.New())),
		TagLabels:       make([]string, len(a.Tags)),

		// ESILinkCount:        0,  // set by a.ESILagInfo below
		// ESILinkSpeed:        "", // set by a.ESILagInfo below
		// ESIPortChannelIDMax: 0,  // set by a.ESILagInfo below
		// ESIPortChannelIDMin: 0,  // set by a.ESILagInfo below
		// RedundancyProtocol:  "", // set by a.ESILAGInfo below
	}

	for i, tag := range a.Tags {
		result.TagLabels[i] = tag.Label
	}
	slices.Sort(result.TagLabels)

	if a.ESILAGInfo != nil {
		result.ESILinkCount = a.ESILAGInfo.LinkCount
		result.ESILinkSpeed = a.ESILAGInfo.LinkSpeed
		result.ESIPortChannelIDMax = a.ESILAGInfo.PortChannelIdMax
		result.ESIPortChannelIDMin = a.ESILAGInfo.PortChannelIdMin
		result.RedundancyProtocol = "esi"
	}

	return json.Marshal(result)
}

func (a *AccessSwitch) UnmarshalJSON(bytes []byte) error {
	var raw rawAccessSwitch
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling access switch: %w", err)
	}

	a.Count = raw.Count
	a.Label = raw.Label
	a.LogicalDevice = NewLogicalDevice(raw.LogicalDeviceID)
	a.Links = raw.Links
	a.LogicalDevice = NewLogicalDevice(raw.LogicalDeviceID)
	a.Tags = make([]Tag, len(raw.TagLabels))
	for i, rawTagLabel := range raw.TagLabels {
		a.Tags[i].Label = rawTagLabel // tag description must be filled by the caller
	}

	// we're done unless "the switch" is actually an ESI LAG pair
	if raw.RedundancyProtocol == "" {
		return nil
	}

	// having failed to find a reason to return early, save the ESI LAG info
	a.ESILAGInfo = &RackTypeAccessSwitchESILAGInfo{
		LinkCount:        raw.ESILinkCount,
		LinkSpeed:        raw.ESILinkSpeed,
		PortChannelIdMax: raw.ESIPortChannelIDMax,
		PortChannelIdMin: raw.ESIPortChannelIDMin,
	}

	return nil
}

// it is safe and reasonable to have a "raw" type for objects which:
// 1) are marshaled and unmarshaled symmetrically (have no metadata to suppress)
// 2) have JSON layout which doesn't align with their public struct layout
type rawAccessSwitch struct {
	Count              int            `json:"instance_count"`
	Label              string         `json:"label"`
	Links              []RackTypeLink `json:"links"`
	LogicalDeviceID    string         `json:"logical_device"`
	RedundancyProtocol string         `json:"redundancy_protocol,omitempty"`
	TagLabels          []string       `json:"tags"`

	ESILinkCount        int         `json:"access_access_link_count"`
	ESILinkSpeed        speed.Speed `json:"access_access_link_speed"`
	ESIPortChannelIDMax int         `json:"access_access_link_port_channel_id_max"`
	ESIPortChannelIDMin int         `json:"access_access_link_port_channel_id_min"`
}
