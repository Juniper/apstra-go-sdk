// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

type VNBinding struct {
	AccessSwitchNodeIDs []string `json:"access_switch_node_ids"`
	SystemID            string   `json:"system_id"` // graph node id of a leaf switch
	VLAN                *uint16  `json:"vlan_id"`   // optional (auto-assign)
}
