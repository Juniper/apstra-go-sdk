package apstra

import "fmt"

const (
	RelationshipTypeNone = RelationshipType(iota)
	RelationshipTypeComposedOf
	RelationshipTypeComposedOfSystems
	RelationshipTypeDeviceProfile
	RelationshipTypeEpAffectedBy
	RelationshipTypeEpMemberOf
	RelationshipTypeEpNested
	RelationshipTypeEpTopLevel
	RelationshipTypeHostedInterfaces
	RelationshipTypeHostedVnInstances
	RelationshipTypeInterfaceMap
	RelationshipTypeInstantiatedBy
	RelationshipTypeInstantiates
	RelationshipTypeLink
	RelationshipTypeLogicalDevice
	RelationshipTypeMemberInterfaces
	RelationshipTypeMemberVNs
	RelationshipTypePartOfRack
	RelationshipTypePolicy
	RelationshipTypeProtocol
	RelationshipTypeRouteTargetPolicy
	RelationshipTypeSecurityPolicy
	RelationshipTypeTag
	RelationshipTypeUnknown = "unknown node type %s"

	relationshipTypeNone              = relationshipType("")
	relationshipTypeComposedOf        = relationshipType("composed_of")
	relationshipTypeComposedOfSystems = relationshipType("composed_of_systems")
	relationshipTypeDeviceProfile     = relationshipType("device_profile")
	relationshipTypeEpAffectedBy      = relationshipType("ep_affected_by")
	relationshipTypeEpMemberOf        = relationshipType("ep_member_of")
	relationshipTypeEpNested          = relationshipType("ep_nested")
	relationshipTypeEpTopLevel        = relationshipType("ep_top_level")
	relationshipTypeHostedInterfaces  = relationshipType("hosted_interfaces")
	relationshipTypeHostedVnInstances = relationshipType("hosted_vn_instances")
	relationshipTypeInterfaceMap      = relationshipType("interface_map")
	relationshipTypeInstantiatedBy    = relationshipType("instantiated_by")
	relationshipTypeInstantiates      = relationshipType("instantiates")
	relationshipTypeLink              = relationshipType("link")
	relationshipTypeLogicalDevice     = relationshipType("logical_device")
	relationshipTypeMemberInterfaces  = relationshipType("member_interfaces")
	relationshipTypeMemberVNs         = relationshipType("member_vns")
	relationshipTypePartOfRack        = relationshipType("part_of_rack")
	relationshipTypePolicy            = relationshipType("policy")
	relationshipTypeProtocol          = relationshipType("protocol")
	relationshipTypeRouteTargetPolicy = relationshipType("route_target_policy")
	relationshipTypeSecurityPolicy    = relationshipType("security_policy")
	relationshipTypeTag               = relationshipType("tag")
	relationshipTypeUnknown           = "unknown node type %d"
)

type RelationshipType int
type relationshipType string

func (o RelationshipType) String() string {
	switch o {
	case RelationshipTypeNone:
		return string(relationshipTypeNone)
	case RelationshipTypeComposedOf:
		return string(relationshipTypeComposedOf)
	case RelationshipTypeComposedOfSystems:
		return string(relationshipTypeComposedOfSystems)
	case RelationshipTypeDeviceProfile:
		return string(relationshipTypeDeviceProfile)
	case RelationshipTypeEpAffectedBy:
		return string(relationshipTypeEpAffectedBy)
	case RelationshipTypeEpMemberOf:
		return string(relationshipTypeEpMemberOf)
	case RelationshipTypeEpNested:
		return string(relationshipTypeEpNested)
	case RelationshipTypeEpTopLevel:
		return string(relationshipTypeEpTopLevel)
	case RelationshipTypeHostedInterfaces:
		return string(relationshipTypeHostedInterfaces)
	case RelationshipTypeHostedVnInstances:
		return string(relationshipTypeHostedVnInstances)
	case RelationshipTypeInterfaceMap:
		return string(relationshipTypeInterfaceMap)
	case RelationshipTypeInstantiatedBy:
		return string(relationshipTypeInstantiatedBy)
	case RelationshipTypeInstantiates:
		return string(relationshipTypeInstantiates)
	case RelationshipTypeLink:
		return string(relationshipTypeLink)
	case RelationshipTypeLogicalDevice:
		return string(relationshipTypeLogicalDevice)
	case RelationshipTypeMemberInterfaces:
		return string(relationshipTypeMemberInterfaces)
	case RelationshipTypeMemberVNs:
		return string(relationshipTypeMemberVNs)
	case RelationshipTypePartOfRack:
		return string(relationshipTypePartOfRack)
	case RelationshipTypePolicy:
		return string(relationshipTypePolicy)
	case RelationshipTypeProtocol:
		return string(relationshipTypeProtocol)
	case RelationshipTypeRouteTargetPolicy:
		return string(relationshipTypeRouteTargetPolicy)
	case RelationshipTypeSecurityPolicy:
		return string(relationshipTypeSecurityPolicy)
	case RelationshipTypeTag:
		return string(relationshipTypeTag)
	default:
		return fmt.Sprintf(relationshipTypeUnknown, o)
	}
}

func (o RelationshipType) QEEAttribute() QEEAttribute {
	return QEEAttribute{
		Key:   "type",
		Value: QEStringVal(o.String()),
	}
}
