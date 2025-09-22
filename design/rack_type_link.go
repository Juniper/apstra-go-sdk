// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
	"github.com/Juniper/apstra-go-sdk/speed"
)

var (
	_ replicator[RackTypeLink] = (*RackTypeLink)(nil)
	_ json.Marshaler           = (*RackTypeLink)(nil)
	_ json.Unmarshaler         = (*RackTypeLink)(nil)
)

type RackTypeLink struct {
	Label              string
	TargetSwitchLabel  string
	LinkPerSwitchCount int
	Speed              speed.Speed
	AttachmentType     enum.LinkAttachmentType
	LAGMode            enum.LAGMode
	SwitchPeer         enum.LinkSwitchPeer
	RailIndex          *int
	Tags               []Tag
}

func (r RackTypeLink) MarshalJSON() ([]byte, error) {
	result := rawRackTypeLink{
		Label:              r.Label,
		TargetSwitchLabel:  r.TargetSwitchLabel,
		LinkPerSwitchCount: zero.PreferDefault(r.LinkPerSwitchCount, 1),
		Speed:              r.Speed,
		AttachmentType:     zero.PreferDefault(r.AttachmentType, enum.LinkAttachmentTypeSingle),
		LAGMode:            r.LAGMode,
		SwitchPeer:         r.SwitchPeer,
		RailIndex:          r.RailIndex,
		Tags:               make([]string, len(r.Tags)),
	}

	for i, tag := range r.Tags {
		result.Tags[i] = tag.Label
	}
	sort.Strings(result.Tags)

	return json.Marshal(result)
}

func (r *RackTypeLink) UnmarshalJSON(bytes []byte) error {
	var raw rawRackTypeLink
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshalling link: %w", err)
	}

	r.Label = raw.Label
	r.TargetSwitchLabel = raw.TargetSwitchLabel
	r.LinkPerSwitchCount = raw.LinkPerSwitchCount
	r.Speed = raw.Speed
	r.AttachmentType = raw.AttachmentType
	r.LAGMode = raw.LAGMode
	r.SwitchPeer = raw.SwitchPeer
	r.RailIndex = raw.RailIndex
	r.Tags = make([]Tag, len(raw.Tags))
	for i, tag := range raw.Tags {
		r.Tags[i] = Tag{Label: tag}
	}

	return nil
}

func (r RackTypeLink) replicate() RackTypeLink {
	r.RailIndex = pointer.To(*r.RailIndex)
	return r
}

type rawRackTypeLink struct {
	Label              string                  `json:"label"`
	TargetSwitchLabel  string                  `json:"target_switch_label"`
	LinkPerSwitchCount int                     `json:"link_per_switch_count"`
	Speed              speed.Speed             `json:"speed"`
	AttachmentType     enum.LinkAttachmentType `json:"attachment_type"`
	LAGMode            enum.LAGMode            `json:"lag_mode"`
	SwitchPeer         enum.LinkSwitchPeer     `json:"switch_peer,omitempty"`
	RailIndex          *int                    `json:"rail_index"`
	Tags               []string                `json:"tags"`
}
