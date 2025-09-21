// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import "github.com/Juniper/apstra-go-sdk/speed"

var _ replicator[RackTypeLeafSwitchMlagInfo] = (*RackTypeLeafSwitchMlagInfo)(nil)

type RackTypeLeafSwitchMlagInfo struct {
	LeafLeafL3LinkCount         int
	LeafLeafL3LinkSpeed         speed.Speed
	LeafLeafL3LinkPortChannelId int
	LeafLeafLinkCount           int
	LeafLeafLinkSpeed           speed.Speed
	LeafLeafLinkPortChannelId   int
	MlagVlanId                  int
}

func (r RackTypeLeafSwitchMlagInfo) replicate() RackTypeLeafSwitchMlagInfo {
	return r
}
