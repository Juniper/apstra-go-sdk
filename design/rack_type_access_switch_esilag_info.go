// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"github.com/Juniper/apstra-go-sdk/speed"
)

type RackTypeAccessSwitchESILAGInfo struct {
	LinkCount        int         `json:"access_access_link_count"`
	LinkSpeed        speed.Speed `json:"access_access_link_speed"`
	PortChannelIdMax int         `json:"access_access_link_port_channel_id_max"`
	PortChannelIdMin int         `json:"access_access_link_port_channel_id_min"`
}

func (r RackTypeAccessSwitchESILAGInfo) Replicate() RackTypeAccessSwitchESILAGInfo {
	return r
}
