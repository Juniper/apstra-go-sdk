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
	RelationshipTypeNone = RelationshipType(iota)
	RelationshipTypeAttachedTaggedVLANs
	RelationshipTypeAttachedUntaggedVLAN
	RelationshipTypeComposedOf
	RelationshipTypeComposedOfSystems
	RelationshipTypeConstraint
	RelationshipTypeDeviceProfile
	RelationshipTypeEpAffectedBy
	RelationshipTypeEpFirstSubpolicy
	RelationshipTypeEpMemberOf
	RelationshipTypeEpNested
	RelationshipTypeEpSubpolicy
	RelationshipTypeEpTopLevel
	RelationshipTypeEvpnInterconnectPeer
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
	RelationshipTypePartOfRedundancyGroup
	RelationshipTypePolicy
	RelationshipTypeProtocol
	RelationshipTypeRouteTargetPolicy
	RelationshipTypeSecurityPolicy
	RelationshipTypeTag

	relationshipTypeNone                  = relationshipType("")
	relationshipTypeAttachedTaggedVLANs   = relationshipType("attached_tagged_vlans")
	relationshipTypeAttachedUntaggedVLAN  = relationshipType("attached_untagged_vlan")
	relationshipTypeComposedOf            = relationshipType("composed_of")
	relationshipTypeComposedOfSystems     = relationshipType("composed_of_systems")
	relationshipTypeConstraint            = relationshipType("constraint")
	relationshipTypeDeviceProfile         = relationshipType("device_profile")
	relationshipTypeEpAffectedBy          = relationshipType("ep_affected_by")
	relationshipTypeEpFirstSubpolicy      = relationshipType("ep_first_subpolicy")
	relationshipTypeEpMemberOf            = relationshipType("ep_member_of")
	relationshipTypeEpNested              = relationshipType("ep_nested")
	relationshipTypeEpSubpolicy           = relationshipType("ep_subpolicy")
	relationshipTypeEpTopLevel            = relationshipType("ep_top_level")
	relationshipTypeEvpnInterconnectPeer  = relationshipType("evpn_interconnect_peer")
	relationshipTypeHostedInterfaces      = relationshipType("hosted_interfaces")
	relationshipTypeHostedVnInstances     = relationshipType("hosted_vn_instances")
	relationshipTypeInterfaceMap          = relationshipType("interface_map")
	relationshipTypeInstantiatedBy        = relationshipType("instantiated_by")
	relationshipTypeInstantiates          = relationshipType("instantiates")
	relationshipTypeLink                  = relationshipType("link")
	relationshipTypeLogicalDevice         = relationshipType("logical_device")
	relationshipTypeMemberInterfaces      = relationshipType("member_interfaces")
	relationshipTypeMemberVNs             = relationshipType("member_vns")
	relationshipTypePartOfRack            = relationshipType("part_of_rack")
	relationshipTypePartOfRedundancyGroup = relationshipType("part_of_redundancy_group")
	relationshipTypePolicy                = relationshipType("policy")
	relationshipTypeProtocol              = relationshipType("protocol")
	relationshipTypeRouteTargetPolicy     = relationshipType("route_target_policy")
	relationshipTypeSecurityPolicy        = relationshipType("security_policy")
	relationshipTypeTag                   = relationshipType("tag")
	relationshipTypeUnknown               = "unknown relationship type"
)

var (
	_ encoding.TextMarshaler   = (*RelationshipType)(nil)
	_ encoding.TextUnmarshaler = (*RelationshipType)(nil)
)

type (
	RelationshipType int
	relationshipType string
)

func (o RelationshipType) MarshalText() (text []byte, err error) {
	s := o.String()
	if strings.HasPrefix(s, relationshipTypeUnknown) {
		return nil, errors.New("cannot marshal: " + s)
	}
	return []byte(o.String()), nil
}

func (o *RelationshipType) UnmarshalText(text []byte) error {
	return o.FromString(string(text))
}

func (o *RelationshipType) FromString(in string) error {
	// Loop over every type starting at zero.
	t := RelationshipType(0)
	for {
		s := t.String()

		// Check if we've run off the end of the list of known types.
		if strings.HasPrefix(s, relationshipTypeUnknown) {
			return fmt.Errorf(relationshipTypeUnknown+" %q", in) // The input string is not valid.
		}

		// Check the input string against t's string representation.
		if s == in {
			*o = t // We found our guy.
			return nil
		}

		t++ // Increment t and try again.
	}
}

func (o RelationshipType) String() string {
	switch o {
	case RelationshipTypeNone:
		return string(relationshipTypeNone)
	case RelationshipTypeAttachedTaggedVLANs:
		return string(relationshipTypeAttachedTaggedVLANs)
	case RelationshipTypeAttachedUntaggedVLAN:
		return string(relationshipTypeAttachedUntaggedVLAN)
	case RelationshipTypeComposedOf:
		return string(relationshipTypeComposedOf)
	case RelationshipTypeComposedOfSystems:
		return string(relationshipTypeComposedOfSystems)
	case RelationshipTypeConstraint:
		return string(relationshipTypeConstraint)
	case RelationshipTypeDeviceProfile:
		return string(relationshipTypeDeviceProfile)
	case RelationshipTypeEpAffectedBy:
		return string(relationshipTypeEpAffectedBy)
	case RelationshipTypeEpFirstSubpolicy:
		return string(relationshipTypeEpFirstSubpolicy)
	case RelationshipTypeEpMemberOf:
		return string(relationshipTypeEpMemberOf)
	case RelationshipTypeEpNested:
		return string(relationshipTypeEpNested)
	case RelationshipTypeEpSubpolicy:
		return string(relationshipTypeEpSubpolicy)
	case RelationshipTypeEpTopLevel:
		return string(relationshipTypeEpTopLevel)
	case RelationshipTypeEvpnInterconnectPeer:
		return string(relationshipTypeEvpnInterconnectPeer)
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
	case RelationshipTypePartOfRedundancyGroup:
		return string(relationshipTypePartOfRedundancyGroup)
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
		return fmt.Sprintf(relationshipTypeUnknown+" %d", o)
	}
}

func (o RelationshipType) QEEAttribute() QEEAttribute {
	return QEEAttribute{
		Key:   "type",
		Value: QEStringVal(o.String()),
	}
}
