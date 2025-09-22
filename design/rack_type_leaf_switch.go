// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/speed"
	"slices"
)

var (
	_ ider                   = (*LeafSwitch)(nil)
	_ replicator[LeafSwitch] = (*LeafSwitch)(nil)
	_ json.Marshaler         = (*LeafSwitch)(nil)
	_ json.Unmarshaler       = (*LeafSwitch)(nil)
)

type LeafSwitch struct {
	Label              string
	LinkPerSpineCount  *int
	LinkPerSpineSpeed  *speed.Speed
	LogicalDevice      LogicalDevice
	RedundancyProtocol *enum.LeafRedundancyProtocol
	Tags               []Tag
	MlagInfo           *RackTypeLeafSwitchMlagInfo

	id string
}

func (l LeafSwitch) ID() *string {
	if l.id == "" {
		return nil
	}
	return &l.id
}

// SetID sets a previously un-set id attribute. If the id attribute is found to
// have an existing value, an error is returned. Presence of an existing value
// is the only reason SetID will return an error. If the id attribute is known
// to be empty, use MustSetID.
func (l LeafSwitch) SetID(id string) error {
	if l.id != "" {
		return IDIsSet(fmt.Errorf("id already has value %q", l.id))
	}

	l.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (l LeafSwitch) MustSetID(id string) {
	err := l.SetID(id)
	if err != nil {
		panic(err)
	}
}

func (l LeafSwitch) replicate() LeafSwitch {
	var linkPerSpineCount *int
	if l.LinkPerSpineCount != nil {
		linkPerSpineCount = pointer.To(*l.LinkPerSpineCount)
	}

	var linkPerSpineSpeed *speed.Speed
	if l.LinkPerSpineSpeed != nil {
		linkPerSpineSpeed = pointer.To(*l.LinkPerSpineSpeed)
	}

	var tags []Tag
	if l.Tags != nil {
		tags = make([]Tag, len(l.Tags))
	}
	for i, tag := range l.Tags {
		tags[i] = tag.replicate()
	}

	var mlagInfo *RackTypeLeafSwitchMlagInfo
	if l.MlagInfo != nil {
		mlagInfo = pointer.To(*l.MlagInfo)
	}

	return LeafSwitch{
		Label:              l.Label,
		LinkPerSpineCount:  linkPerSpineCount,
		LinkPerSpineSpeed:  linkPerSpineSpeed,
		LogicalDevice:      l.LogicalDevice.replicate(),
		RedundancyProtocol: l.RedundancyProtocol,
		Tags:               tags,
		MlagInfo:           mlagInfo,
	}
}

func (l LeafSwitch) MarshalJSON() ([]byte, error) {
	raw := rawLeafSwitch{
		Label:              l.Label,
		LinkPerSpineCount:  l.LinkPerSpineCount,
		LinkPerSpineSpeed:  l.LinkPerSpineSpeed,
		LogicalDeviceID:    fmt.Sprintf("%x", mustHashForComparison(l.LogicalDevice, md5.New())),
		RedundancyProtocol: l.RedundancyProtocol,
		TagLabels:          make([]string, len(l.Tags)),
	}

	for i, tag := range l.Tags {
		raw.TagLabels[i] = tag.Label
	}
	slices.Sort(raw.TagLabels)

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

	// look for reasons to return before adding MLAG info
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

	// having failed to find a reason to return early, save the MLAG info
	l.MlagInfo = &RackTypeLeafSwitchMlagInfo{
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

// it is safe and reasonable to have a "raw" type for objects which:
// 1) are marshaled and unmarshaled symmetrically (have no metadata to suppress)
// 2) have JSON layout which doesn't align with their struct layout
type rawLeafSwitch struct {
	Label              string                       `json:"label"`
	LinkPerSpineCount  *int                         `json:"link_per_spine_count,omitempty"`
	LinkPerSpineSpeed  *speed.Speed                 `json:"link_per_spine_speed,omitempty"`
	LogicalDeviceID    string                       `json:"logical_device"`
	RedundancyProtocol *enum.LeafRedundancyProtocol `json:"redundancy_protocol,omitempty"`
	TagLabels          []string                     `json:"tags"`

	LeafLeafL3LinkCount         *int         `json:"leaf_leaf_l3_link_count,omitempty"`
	LeafLeafL3LinkSpeed         *speed.Speed `json:"leaf_leaf_l3_link_speed,omitempty"`
	LeafLeafL3LinkPortChannelId *int         `json:"leaf_leaf_l3_link_port_channel_id,omitempty"`
	LeafLeafLinkCount           *int         `json:"leaf_leaf_link_count,omitempty"`
	LeafLeafLinkSpeed           *speed.Speed `json:"leaf_leaf_link_speed,omitempty"`
	LeafLeafLinkPortChannelId   *int         `json:"leaf_leaf_link_port_channel_id,omitempty"`
	MlagVlanId                  *int         `json:"mlag_vlan_id,omitempty"`
}
