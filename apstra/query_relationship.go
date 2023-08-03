package apstra

import "fmt"

const (
	RelationshipTypeNone = RelationshipType(iota)
	RelationshipTypeComposedOfSystems
	RelationshipTypeDeviceProfile
	RelationshipTypeHostedInterfaces
	RelationshipTypeInterfaceMap
	RelationshipTypeInstantiates
	RelationshipTypeLink
	RelationshipTypeLogicalDevice
	RelationshipTypeMemberVNs
	RelationshipTypePartOfRack
	RelationshipTypeTag
	RelationshipTypeUnknown = "unknown node type %s"

	relationshipTypeNone              = relationshipType("")
	relationshipTypeComposedOfSystems = relationshipType("composed_of_systems")
	relationshipTypeDeviceProfile     = relationshipType("device_profile")
	relationshipTypeHostedInterfaces  = relationshipType("hosted_interfaces")
	relationshipTypeInterfaceMap      = relationshipType("interface_map")
	relationshipTypeInstantiates      = relationshipType("instantiates")
	relationshipTypeLink              = relationshipType("link")
	relationshipTypeLogicalDevice     = relationshipType("logical_device")
	relationshipTypeMemberVNs         = relationshipType("member_vns")
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
	case RelationshipTypeHostedInterfaces:
		return string(relationshipTypeHostedInterfaces)
	case RelationshipTypeInterfaceMap:
		return string(relationshipTypeInterfaceMap)
	case RelationshipTypeInstantiates:
		return string(relationshipTypeInstantiates)
	case RelationshipTypeLink:
		return string(relationshipTypeLink)
	case RelationshipTypeLogicalDevice:
		return string(relationshipTypeLogicalDevice)
	case RelationshipTypeMemberVNs:
		return string(relationshipTypeMemberVNs)
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
