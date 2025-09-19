// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/hash"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	"github.com/Juniper/apstra-go-sdk/speed"
)

type LeafMlagInfo struct {
	LeafLeafL3LinkCount         int
	LeafLeafL3LinkSpeed         speed.Speed
	LeafLeafL3LinkPortChannelId int
	LeafLeafLinkCount           int
	LeafLeafLinkSpeed           speed.Speed
	LeafLeafLinkPortChannelId   int
	MlagVlanId                  int
}

var (
	_ json.Marshaler   = (*LeafSwitch)(nil)
	_ json.Unmarshaler = (*LeafSwitch)(nil)
)

type LeafSwitch struct {
	Label              string
	LinkPerSpineCount  *int
	LinkPerSpineSpeed  *speed.Speed
	LogicalDevice      LogicalDevice
	RedundancyProtocol *enum.LeafRedundancyProtocol
	Tags               []Tag
	MlagInfo           *LeafMlagInfo
}

func (l LeafSwitch) MarshalJSON() ([]byte, error) {
	raw := rawLeafSwitch{
		Label:              l.Label,
		LinkPerSpineCount:  l.LinkPerSpineCount,
		LinkPerSpineSpeed:  l.LinkPerSpineSpeed,
		LogicalDeviceID:    fmt.Sprintf("%x", hash.StructMust(l.LogicalDevice, sha256.New())),
		RedundancyProtocol: l.RedundancyProtocol,
		TagLabels:          slice.WithCapacityOrNil("", len(l.Tags)),
	}

	for _, tag := range l.Tags {
		raw.TagLabels = append(raw.TagLabels, tag.Label)
	}

	if l.MlagInfo != nil {
		raw.LeafLeafL3LinkCount = nil
		raw.LeafLeafL3LinkSpeed = nil
		raw.LeafLeafL3LinkPortChannelId = nil
		raw.LeafLeafLinkCount = nil
		raw.LeafLeafLinkSpeed = nil
		raw.LeafLeafLinkPortChannelId = nil
		raw.MlagVlanId = nil
	}

	return json.Marshal(raw)
}

func (l *LeafSwitch) UnmarshalJSON(bytes []byte) error {
	var raw rawLeafSwitch
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling rawLeafSwitch: %w", err)
	}

	l.Label = raw.Label
	l.LinkPerSpineCount = raw.LinkPerSpineCount
	l.LinkPerSpineSpeed = raw.LinkPerSpineSpeed
	l.LogicalDevice = NewLogicalDevice(raw.LogicalDeviceID)
	l.RedundancyProtocol = raw.RedundancyProtocol
	l.Tags = make([]Tag, len(raw.TagLabels))
	for i, rawTagLabel := range raw.TagLabels {
		l.Tags[i].Label = rawTagLabel // tag description must be filled by the caller
	}

	// look for reasons to NOT save MLAG info
	if raw.RedundancyProtocol == nil ||
		*raw.RedundancyProtocol != enum.LeafRedundancyProtocolMLAG ||
		raw.LeafLeafL3LinkCount == nil ||
		raw.LeafLeafL3LinkSpeed == nil ||
		raw.LeafLeafL3LinkPortChannelId == nil ||
		raw.LeafLeafLinkCount == nil ||
		raw.LeafLeafLinkSpeed == nil ||
		raw.LeafLeafLinkPortChannelId == nil {
		return nil
	}

	// having failed to find an excuse, save the MLAG info
	l.MlagInfo = &LeafMlagInfo{
		LeafLeafL3LinkCount:         *raw.LeafLeafL3LinkCount,
		LeafLeafL3LinkSpeed:         *raw.LeafLeafL3LinkSpeed,
		LeafLeafL3LinkPortChannelId: *raw.LeafLeafL3LinkPortChannelId,
		LeafLeafLinkCount:           *raw.LeafLeafLinkCount,
		LeafLeafLinkSpeed:           *raw.LeafLeafLinkSpeed,
		LeafLeafLinkPortChannelId:   *raw.LeafLeafLinkPortChannelId,
		MlagVlanId:                  *raw.MlagVlanId,
	}

	return nil
}

type rawLeafSwitch struct {
	Label              string                       `json:"label"`
	LinkPerSpineCount  *int                         `json:"link_per_spine_count,omitempty"`
	LinkPerSpineSpeed  *speed.Speed                 `json:"link_per_spine_speed,omitempty"`
	LogicalDeviceID    string                       `json:"logical_device"`
	RedundancyProtocol *enum.LeafRedundancyProtocol `json:"redundancy_protocol,omitempty"`
	TagLabels          []string                     `json:"taglabels"`

	LeafLeafL3LinkCount         *int         `json:"leaf_leaf_l3_link_count,omitempty"`
	LeafLeafL3LinkSpeed         *speed.Speed `json:"leaf_leaf_l3_link_speed,omitempty"`
	LeafLeafL3LinkPortChannelId *int         `json:"leaf_leaf_l3_link_port_channel_id,omitempty"`
	LeafLeafLinkCount           *int         `json:"leaf_leaf_link_count,omitempty"`
	LeafLeafLinkSpeed           *speed.Speed `json:"leaf_leaf_link_speed,omitempty"`
	LeafLeafLinkPortChannelId   *int         `json:"leaf_leaf_link_port_channel_id,omitempty"`
	MlagVlanId                  *int         `json:"mlag_vlan_id,omitempty"`
}
