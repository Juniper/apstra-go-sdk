package apstra

import "fmt"

const (
	RelationshipTypeNone = RelationshipType(iota)
	RelationshipTypeComposedOfSystems
	RelationshipTypeDeviceProfile
	RelationshipTypeEpAffectedBy
	RelationshipTypeEpMemberOf
	RelationshipTypeEpNested
	RelationshipTypeEpTopLevel
	RelationshipTypeHostedInterfaces
	RelationshipTypeInterfaceMap
	RelationshipTypeLink
	RelationshipTypeLogicalDevice
	RelationshipTypePartOfRack
	RelationshipTypeTag
	RelationshipTypeUnknown = "unknown node type %s"

	relationshipTypeNone              = relationshipType("")
	relationshipTypeComposedOfSystems = relationshipType("composed_of_systems")
	relationshipTypeDeviceProfile     = relationshipType("device_profile")
	relationshipTypeEpAffectedBy      = relationshipType("ep_affected_by")
	relationshipTypeEpMemberOf        = relationshipType("ep_member_of")
	relationshipTypeEpNested          = relationshipType("ep_nested")
	relationshipTypeEpTopLevel        = relationshipType("ep_top_level")
	relationshipTypeHostedInterfaces  = relationshipType("hosted_interfaces")
	relationshipTypeInterfaceMap      = relationshipType("interface_map")
	relationshipTypeLink              = relationshipType("link")
	relationshipTypeLogicalDevice     = relationshipType("logical_device")
	relationshipTypePartOfRack        = relationshipType("part_of_rack")
	relationshipTypeTag               = relationshipType("tag")
	relationshipTypeUnknown           = "unknown node type %d"
)

type RelationshipType int
type relationshipType string

func (o RelationshipType) String() string {
	switch o {
	case RelationshipTypeNone:
		return string(relationshipTypeNone)
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
	case RelationshipTypeInterfaceMap:
		return string(relationshipTypeInterfaceMap)
	case RelationshipTypeLink:
		return string(relationshipTypeLink)
	case RelationshipTypeLogicalDevice:
		return string(relationshipTypeLogicalDevice)
	case RelationshipTypePartOfRack:
		return string(relationshipTypePartOfRack)
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
