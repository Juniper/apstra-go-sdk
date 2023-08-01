package apstra

import "fmt"

const (
	NodeTypeNone = NodeType(iota)
	NodeTypeEpApplicationInstance
	NodeTypeEpEndpointPolicy
	NodeTypeEpGroup
	NodeTypeInterface
	NodeTypeInterfaceMap
	NodeTypeLink
	NodeTypeLogicalDevice
	NodeTypeMetadata
	NodeTypePolicy
	NodeTypeRack
	NodeTypeRedundancyGroup
	NodeTypeRoutingPolicy
	NodeTypeSecurityZone
	NodeTypeSystem
	NodeTypeTag
	NodeTypeVirtualNetwork
	NodeTypeUnknown = "unknown node type %s"

	nodeTypeNone                  = nodeType("")
	nodeTypeEpApplicationInstance = nodeType("ep_application_instance")
	nodeTypeEpEndpointPolicy      = nodeType("ep_endpoint_policy")
	nodeTypeEpGroup               = nodeType("ep_group")
	nodeTypeInterface             = nodeType("interface")
	nodeTypeInterfaceMap          = nodeType("interface_map")
	nodeTypeLink                  = nodeType("link")
	nodeTypeLogicalDevice         = nodeType("logical_device")
	nodeTypeMetadata              = nodeType("metadata")
	nodeTypePolicy                = nodeType("policy")
	nodeTypeRack                  = nodeType("rack")
	nodeTypeRedundancyGroup       = nodeType("redundancy_group")
	nodeTypeRoutingPolicy         = nodeType("routing_policy")
	nodeTypeSecurityZone          = nodeType("security_zone")
	nodeTypeSystem                = nodeType("system")
	nodeTypeTag                   = nodeType("tag")
	nodeTypeVirtualNetwork        = nodeType("virtual_network")
	nodeTypeUnknown               = "unknown node type %d"
)

type NodeType int
type nodeType string

func (o NodeType) String() string {
	switch o {
	case NodeTypeNone:
		return string(nodeTypeNone)
	case NodeTypeEpApplicationInstance:
		return string(nodeTypeEpApplicationInstance)
	case NodeTypeEpEndpointPolicy:
		return string(nodeTypeEpEndpointPolicy)
	case NodeTypeEpGroup:
		return string(nodeTypeEpGroup)
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
	case NodeTypeRack:
		return string(nodeTypeRack)
	case NodeTypeRedundancyGroup:
		return string(nodeTypeRedundancyGroup)
	case NodeTypeRoutingPolicy:
		return string(nodeTypeRoutingPolicy)
	case NodeTypeSecurityZone:
		return string(nodeTypeSecurityZone)
	case NodeTypeSystem:
		return string(nodeTypeSystem)
	case NodeTypeTag:
		return string(nodeTypeTag)
	case NodeTypeVirtualNetwork:
		return string(nodeTypeVirtualNetwork)
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
