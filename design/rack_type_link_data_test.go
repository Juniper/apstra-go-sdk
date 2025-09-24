// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
)

var linkSimple = RackTypeLink{
	Label:              "Simple Link",
	TargetSwitchLabel:  "leafy",
	LinkPerSwitchCount: 2,
	Speed:              "50G",
}

const linkSimpleJSON = `{
  "label": "Simple Link",
  "target_switch_label": "leafy",
  "link_per_switch_count": 2,
  "speed": { "unit": "G", "value": 50 },
  "attachment_type": "singleAttached",
  "tags": []
}`

var linkComplicated = RackTypeLink{
	Label:              "Complicated Link",
	TargetSwitchLabel:  "leafy",
	LinkPerSwitchCount: 2,
	Speed:              "50G",
	AttachmentType:     enum.LinkAttachmentTypeDual,
	LAGMode:            enum.LAGModePassiveLACP,
	SwitchPeer:         enum.LinkSwitchPeerSecond,
	RailIndex:          pointer.To(3),
	Tags:               []Tag{{Label: "b", Description: "B"}, {Label: "a", Description: "A"}},
}

const linkComplicatedJSON = `{
  "label": "Complicated Link",
  "target_switch_label": "leafy",
  "link_per_switch_count": 2,
  "speed": { "unit": "G", "value": 50 },
  "attachment_type": "dualAttached",
  "tags": ["a", "b"],
  "rail_index": 3,
  "lag_mode" : "lacp_passive",
  "switch_peer": "second"
}`
