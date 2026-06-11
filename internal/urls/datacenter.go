// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package urls

const (
	DatacenterSecurityZones           = blueprintByID + pathDelim + "security-zones"
	DatacenterSecurityZoneById        = DatacenterSecurityZones + pathDelim + "%s"
	DatacenterSecurityZoneDHCPServers = DatacenterSecurityZoneById + pathDelim + "dhcp-servers"
	DatacenterSecurityZoneLoopbacks   = DatacenterSecurityZoneById + pathDelim + "loopbacks"

	DatacenterSwitchingZones    = blueprintByID + pathDelim + "switching-zones"
	DatacenterSwitchingZoneByID = DatacenterSwitchingZones + pathDelim + "%s"

	DatacenterVirtualNetworks    = blueprintByID + pathDelim + "virtual-networks"
	DatacenterVirtualNetworkByID = DatacenterVirtualNetworks + pathDelim + "%s"
)
