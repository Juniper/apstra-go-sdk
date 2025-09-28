// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/speed"
)

var _ internal.Replicator[RackTypeLeafSwitchMLAGInfo] = (*RackTypeLeafSwitchMLAGInfo)(nil)

type RackTypeLeafSwitchMLAGInfo struct {
	LeafLeafL3LinkCount         int         `json:"leaf_leaf_l3_link_count"`
	LeafLeafL3LinkSpeed         speed.Speed `json:"leaf_leaf_l3_link_speed"`
	LeafLeafL3LinkPortChannelId int         `json:"leaf_leaf_l3_link_port_channel_id"`
	LeafLeafLinkCount           int         `json:"leaf_leaf_link_count"`
	LeafLeafLinkSpeed           speed.Speed `json:"leaf_leaf_link_speed"`
	LeafLeafLinkPortChannelId   int         `json:"leaf_leaf_link_port_channel_id"`
	MLAGVLAN                    int         `json:"mlagvlan"`
}

// Replicate returns a copy of itself with zero values for metadata fields
func (r RackTypeLeafSwitchMLAGInfo) Replicate() RackTypeLeafSwitchMLAGInfo {
	return r
}
