// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import "fmt"

const (
	NodeTypeNone = NodeType(iota)
	NodeTypeAntiAffinityPolicy
	NodeTypeEpApplicationInstance
	NodeTypeDeviceProfile
	NodeTypeDomain
	NodeTypeEpEndpointPolicy
	NodeTypeEpGroup
	NodeTypeFabricAddressingPolicy
	NodeTypeFabricPolicy
	NodeTypeInterface
	NodeTypeInterfaceMap
	NodeTypeLink
	NodeTypeLogicalDevice
	NodeTypeMetadata
	NodeTypePolicy
	NodeTypeProtocol
	NodeTypeRack
	NodeTypeRedundancyGroup
	NodeTypeRouteTargetPolicy
	NodeTypeRoutingPolicy
	NodeTypeSecurityZone
	NodeTypeSecurityZoneInstance
	NodeTypeSecurityZonePolicy
	NodeTypeSystem
	NodeTypeTag
	NodeTypeVirtualNetwork
	NodeTypeVirtualNetworkInstance
	NodeTypeVirtualNetworkPolicy
	NodeTypeUnknown = "unknown node type %s"

	nodeTypeNone                   = nodeType("")
	nodeTypeAntiAffinityPolicy     = nodeType("anti_affinity_policy")
	nodeTypeEpApplicationInstance  = nodeType("ep_application_instance")
	nodeTypeDeviceProfile          = nodeType("device_profile")
	nodeTypeDomain                 = nodeType("domain")
	nodeTypeEpEndpointPolicy       = nodeType("ep_endpoint_policy")
	nodeTypeEpGroup                = nodeType("ep_group")
	nodeTypeFabricAddressingPolicy = nodeType("fabric_addressing_policy")
	nodeTypeFabricPolicy           = nodeType("fabric_policy")
	nodeTypeInterface              = nodeType("interface")
	nodeTypeInterfaceMap           = nodeType("interface_map")
	nodeTypeLink                   = nodeType("link")
	nodeTypeLogicalDevice          = nodeType("logical_device")
	nodeTypeMetadata               = nodeType("metadata")
	nodeTypePolicy                 = nodeType("policy")
	nodeTypeProtocol               = nodeType("protocol")
	nodeTypeRack                   = nodeType("rack")
	nodeTypeRedundancyGroup        = nodeType("redundancy_group")
	nodeTypeRouteTargetPolicy      = nodeType("route_target_policy")
	nodeTypeRoutingPolicy          = nodeType("routing_policy")
	nodeTypeSecurityZone           = nodeType("security_zone")
	nodeTypeSecurityZoneInstance   = nodeType("sz_instance")
	nodeTypeSecurityZonePolicy     = nodeType("security_zone_policy")
	nodeTypeSystem                 = nodeType("system")
	nodeTypeTag                    = nodeType("tag")
	nodeTypeVirtualNetwork         = nodeType("virtual_network")
	nodeTypeVirtualNetworkInstance = nodeType("vn_instance")
	nodeTypeVirtualNetworkPolicy   = nodeType("virtual_network_policy")
	nodeTypeUnknown                = "unknown node type %d"
)

type (
	NodeType int
	nodeType string
)

func (o NodeType) String() string {
	switch o {
	case NodeTypeNone:
		return string(nodeTypeNone)
	case NodeTypeAntiAffinityPolicy:
		return string(nodeTypeAntiAffinityPolicy)
	case NodeTypeEpApplicationInstance:
		return string(nodeTypeEpApplicationInstance)
	case NodeTypeDeviceProfile:
		return string(nodeTypeDeviceProfile)
	case NodeTypeDomain:
		return string(nodeTypeDomain)
	case NodeTypeEpEndpointPolicy:
		return string(nodeTypeEpEndpointPolicy)
	case NodeTypeEpGroup:
		return string(nodeTypeEpGroup)
	case NodeTypeFabricAddressingPolicy:
		return string(nodeTypeFabricAddressingPolicy)
	case NodeTypeFabricPolicy:
		return string(nodeTypeFabricPolicy)
	case NodeTypeInterface:
		return string(nodeTypeInterface)
	case NodeTypeInterfaceMap:
		return string(nodeTypeInterfaceMap)
	case NodeTypeLink:
		return string(nodeTypeLink)
	case NodeTypeLogicalDevice:
		return string(nodeTypeLogicalDevice)
	case NodeTypeMetadata:
		return string(nodeTypeMetadata)
	case NodeTypePolicy:
		return string(nodeTypePolicy)
	case NodeTypeProtocol:
		return string(nodeTypeProtocol)
	case NodeTypeRack:
		return string(nodeTypeRack)
	case NodeTypeRedundancyGroup:
		return string(nodeTypeRedundancyGroup)
	case NodeTypeRouteTargetPolicy:
		return string(nodeTypeRouteTargetPolicy)
	case NodeTypeRoutingPolicy:
		return string(nodeTypeRoutingPolicy)
	case NodeTypeSecurityZone:
		return string(nodeTypeSecurityZone)
	case NodeTypeSecurityZoneInstance:
		return string(nodeTypeSecurityZoneInstance)
	case NodeTypeSecurityZonePolicy:
		return string(nodeTypeSecurityZonePolicy)
	case NodeTypeSystem:
		return string(nodeTypeSystem)
	case NodeTypeTag:
		return string(nodeTypeTag)
	case NodeTypeVirtualNetwork:
		return string(nodeTypeVirtualNetwork)
	case NodeTypeVirtualNetworkInstance:
		return string(nodeTypeVirtualNetworkInstance)
	case NodeTypeVirtualNetworkPolicy:
		return string(nodeTypeVirtualNetworkPolicy)
	default:
		return fmt.Sprintf(nodeTypeUnknown, o)
	}
}

func (o NodeType) QEEAttribute() QEEAttribute {
	return QEEAttribute{
		Key:   "type",
		Value: QEStringVal(o.String()),
	}
}
