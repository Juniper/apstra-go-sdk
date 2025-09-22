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
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/speed"
)

var (
	_ logicalDeviceIDer      = (*LeafSwitch)(nil)
	_ replicator[LeafSwitch] = (*LeafSwitch)(nil)
	_ json.Marshaler         = (*LeafSwitch)(nil)
	_ json.Unmarshaler       = (*LeafSwitch)(nil)
)

type LeafSwitch struct {
	Label              string
	LinkPerSpineCount  *int
	LinkPerSpineSpeed  *speed.Speed
	LogicalDevice      LogicalDevice
	RedundancyProtocol enum.LeafRedundancyProtocol
	Tags               []Tag
	MLAGInfo           *RackTypeLeafSwitchMLAGInfo
}

// logicalDeviceID returns *string representing the ID of the embedded logical
// device. If the LD ID is unset, nil is returned
func (l LeafSwitch) logicalDeviceID() *string {
	if l.LogicalDevice.id == "" {
		return nil
	}
	return pointer.To(l.LogicalDevice.id)
}

// replicate returns a copy of itself with zero values for metadata fields
func (l LeafSwitch) replicate() LeafSwitch {
	result := LeafSwitch{
		Label:              l.Label,
		LogicalDevice:      l.LogicalDevice.replicate(),
		RedundancyProtocol: l.RedundancyProtocol,
		Tags:               make([]Tag, len(l.Tags)),

		// LinkPerSpineCount:  nil,
		// LinkPerSpineSpeed:  nil,
		// MLAGInfo:           nil,
	}

	for i, tag := range l.Tags {
		result.Tags[i] = tag.replicate()
	}

	if l.LinkPerSpineCount != nil {
		result.LinkPerSpineCount = pointer.To(*l.LinkPerSpineCount)
	}

	if l.LinkPerSpineSpeed != nil {
		result.LinkPerSpineSpeed = pointer.To(*l.LinkPerSpineSpeed)
	}

	if l.MLAGInfo != nil {
		result.MLAGInfo = pointer.To(l.MLAGInfo.replicate())
	}

	return result
}

func (l LeafSwitch) MarshalJSON() ([]byte, error) {
	result := rawLeafSwitch{
		Label:              l.Label,
		LinkPerSpineCount:  l.LinkPerSpineCount,
		LinkPerSpineSpeed:  l.LinkPerSpineSpeed,
		LogicalDeviceID:    fmt.Sprintf("%x", mustHashForComparison(l.LogicalDevice, md5.New())),
		RedundancyProtocol: l.RedundancyProtocol.String(),
		TagLabels:          make([]string, len(l.Tags)),

		// LeafLeafL3LinkCount:         0,  // set by l.MLAGInfo below
		// LeafLeafL3LinkSpeed:         "", // set by l.MLAGInfo below
		// LeafLeafL3LinkPortChannelId: 0,  // set by l.MLAGInfo below
		// LeafLeafLinkCount:           0,  // set by l.MLAGInfo below
		// LeafLeafLinkSpeed:           "", // set by l.MLAGInfo below
		// LeafLeafLinkPortChannelId:   0,  // set by l.MLAGInfo below
		// MLAGVLAN:                    0,  // set by l.MLAGInfo below

	}

	for i, tag := range l.Tags {
		result.TagLabels[i] = tag.Label
	}
	slices.Sort(result.TagLabels)

	if l.MLAGInfo != nil {
		result.LeafLeafL3LinkCount = l.MLAGInfo.LeafLeafL3LinkCount
		result.LeafLeafL3LinkSpeed = l.MLAGInfo.LeafLeafL3LinkSpeed
		result.LeafLeafL3LinkPortChannelId = l.MLAGInfo.LeafLeafL3LinkPortChannelId
		result.LeafLeafLinkCount = l.MLAGInfo.LeafLeafLinkCount
		result.LeafLeafLinkSpeed = l.MLAGInfo.LeafLeafLinkSpeed
		result.LeafLeafLinkPortChannelId = l.MLAGInfo.LeafLeafLinkPortChannelId
		result.MLAGVLAN = l.MLAGInfo.MLAGVLAN
	}

	return json.Marshal(result)
}

func (l *LeafSwitch) UnmarshalJSON(bytes []byte) error {
	var raw rawLeafSwitch
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling leaf switch: %w", err)
	}

	l.Label = raw.Label
	l.LinkPerSpineCount = raw.LinkPerSpineCount
	l.LinkPerSpineSpeed = raw.LinkPerSpineSpeed
	l.LogicalDevice = NewLogicalDevice(raw.LogicalDeviceID)
	err = l.RedundancyProtocol.FromString(raw.RedundancyProtocol)
	if err != nil {
		return fmt.Errorf("parsing field 'RedundancyProtocol': %w", err)
	}
	l.Tags = make([]Tag, len(raw.TagLabels))
	for i, rawTagLabel := range raw.TagLabels {
		l.Tags[i].Label = rawTagLabel // tag description must be filled by the caller
	}

	// look for reasons to return before adding MLAG info
	if l.RedundancyProtocol == enum.LeafRedundancyProtocolNone {
		return nil
	}

	// having failed to find a reason to return early, save the MLAG info
	l.MLAGInfo = &RackTypeLeafSwitchMLAGInfo{
		LeafLeafL3LinkCount:         raw.LeafLeafL3LinkCount,
		LeafLeafL3LinkSpeed:         raw.LeafLeafL3LinkSpeed,
		LeafLeafL3LinkPortChannelId: raw.LeafLeafL3LinkPortChannelId,
		LeafLeafLinkCount:           raw.LeafLeafLinkCount,
		LeafLeafLinkSpeed:           raw.LeafLeafLinkSpeed,
		LeafLeafLinkPortChannelId:   raw.LeafLeafLinkPortChannelId,
		MLAGVLAN:                    raw.MLAGVLAN,
	}

	return nil
}

// it is safe and reasonable to have a "raw" type for objects which:
// 1) are marshaled and unmarshaled symmetrically (have no metadata to suppress)
// 2) have JSON layout which doesn't align with their public struct layout
type rawLeafSwitch struct {
	Label              string       `json:"label"`
	LinkPerSpineCount  *int         `json:"link_per_spine_count,omitempty"`
	LinkPerSpineSpeed  *speed.Speed `json:"link_per_spine_speed,omitempty"`
	LogicalDeviceID    string       `json:"logical_device"`
	RedundancyProtocol string       `json:"redundancy_protocol,omitempty"`
	TagLabels          []string     `json:"tags"`

	LeafLeafL3LinkCount         int         `json:"leaf_leaf_l3_link_count,omitempty"`
	LeafLeafL3LinkSpeed         speed.Speed `json:"leaf_leaf_l3_link_speed,omitempty"`
	LeafLeafL3LinkPortChannelId int         `json:"leaf_leaf_l3_link_port_channel_id,omitempty"`
	LeafLeafLinkCount           int         `json:"leaf_leaf_link_count,omitempty"`
	LeafLeafLinkSpeed           speed.Speed `json:"leaf_leaf_link_speed,omitempty"`
	LeafLeafLinkPortChannelId   int         `json:"leaf_leaf_link_port_channel_id,omitempty"`
	MLAGVLAN                    int         `json:"mlag_vlan_id,omitempty"`
}
