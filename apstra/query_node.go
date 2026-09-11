// Copyright (c) Juniper Networks, Inc., 2023-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"encoding"
	"errors"
	"fmt"
	"strings"
)

const (
	NodeTypeNone = NodeType(iota)
	NodeTypeAntiAffinityPolicy
	NodeTypeEpApplicationInstance
	NodeTypeDeviceProfile
	NodeTypeDomain
	NodeTypeEpEndpointPolicy
	NodeTypeEpGroup
	NodeTypeEvpnInterconnectGroup
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
	NodeTypeRoutingZoneConstraint
	NodeTypeSecurityZone
	NodeTypeSwitchingZone
	NodeTypeSecurityZoneInstance
	NodeTypeSecurityZonePolicy
	NodeTypeSystem
	NodeTypeTag
	NodeTypeVirtualNetwork
	NodeTypeVirtualNetworkInstance
	NodeTypeVirtualNetworkPolicy

	nodeTypeNone                   = nodeType("")
	nodeTypeAntiAffinityPolicy     = nodeType("anti_affinity_policy")
	nodeTypeEpApplicationInstance  = nodeType("ep_application_instance")
	nodeTypeDeviceProfile          = nodeType("device_profile")
	nodeTypeDomain                 = nodeType("domain")
	nodeTypeEpEndpointPolicy       = nodeType("ep_endpoint_policy")
	nodeTypeEpGroup                = nodeType("ep_group")
	nodeTypeEvpnInterconnectGroup  = nodeType("evpn_interconnect_group")
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
	nodeTypeRoutingZoneConstraint  = nodeType("routing_zone_constraint")
	nodeTypeSecurityZone           = nodeType("security_zone")
	nodeTypeSwitchingZone          = nodeType("switching_zone")
	nodeTypeSecurityZoneInstance   = nodeType("sz_instance")
	nodeTypeSecurityZonePolicy     = nodeType("security_zone_policy")
	nodeTypeSystem                 = nodeType("system")
	nodeTypeTag                    = nodeType("tag")
	nodeTypeVirtualNetwork         = nodeType("virtual_network")
	nodeTypeVirtualNetworkInstance = nodeType("vn_instance")
	nodeTypeVirtualNetworkPolicy   = nodeType("virtual_network_policy")
	nodeTypeUnknown                = "unknown node type"
)

var (
	_ encoding.TextMarshaler   = (*NodeType)(nil)
	_ encoding.TextUnmarshaler = (*NodeType)(nil)
)

type (
	NodeType int
	nodeType string
)

func (o NodeType) MarshalText() (text []byte, err error) {
	s := o.String()
	if strings.HasPrefix(s, nodeTypeUnknown) {
		return nil, errors.New("cannot marshal: " + s)
	}
	return []byte(o.String()), nil
}

func (o *NodeType) UnmarshalText(text []byte) error {
	return o.FromString(string(text))
}

func (o *NodeType) FromString(in string) error {
	// Loop over every type starting at zero.
	t := NodeType(0)
	for {
		s := t.String()

		// Check if we've run off the end of the list of known types.
		if strings.HasPrefix(s, nodeTypeUnknown) {
			return fmt.Errorf(nodeTypeUnknown+" %q", in) // The input string is not valid.
		}

		// Check the input string against t's string representation.
		if s == in {
			*o = t // We found our guy.
			return nil
		}

		t++ // Increment t and try again.
	}
}

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
	case NodeTypeEvpnInterconnectGroup:
		return string(nodeTypeEvpnInterconnectGroup)
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
	case NodeTypeRoutingZoneConstraint:
		return string(nodeTypeRoutingZoneConstraint)
	case NodeTypeSecurityZone:
		return string(nodeTypeSecurityZone)
	case NodeTypeSwitchingZone:
		return string(nodeTypeSwitchingZone)
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
		return fmt.Sprintf(nodeTypeUnknown+" %d", o)
	}
}

func (o NodeType) QEEAttribute() QEEAttribute {
	return QEEAttribute{
		Key:   "type",
		Value: QEStringVal(o.String()),
	}
}
