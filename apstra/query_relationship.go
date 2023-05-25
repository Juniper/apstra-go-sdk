package apstra

import "fmt"

const (
	RelationshipTypeNone = iota
	RelationshipTypeComposedOfSystems
	RelationshipTypeDeviceProfile
	RelationshipTypeHostedInterfaces
	RelationshipTypeInterfaceMap
	RelationshipTypeLink
	RelationshipTypeLogicalDevice
	RelationshipTypeTag
	RelationshipTypeUnknown = "unknown node type %s"

	relationshipTypeNone              = relationshipType("")
	relationshipTypeComposedOfSystems = relationshipType("composed_of_systems")
	relationshipTypeDeviceProfile     = relationshipType("device_profile")
	relationshipTypeHostedInterfaces  = relationshipType("hosted_interfaces")
	relationshipTypeInterfaceMap      = relationshipType("interface_map")
	relationshipTypeLink              = relationshipType("link")
	relationshipTypeLogicalDevice     = relationshipType("logical_device")
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
	case RelationshipTypeLink:
		return string(relationshipTypeLink)
	case RelationshipTypeLogicalDevice:
		return string(relationshipTypeLogicalDevice)
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
