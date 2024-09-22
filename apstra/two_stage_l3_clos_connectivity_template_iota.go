// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import "fmt"

type (
	CtPrimitivePolicyTypeName int
	ctPrimitivePolicyTypeName string
)

const (
	CtPrimitivePolicyTypeNameNone = CtPrimitivePolicyTypeName(iota)
	CtPrimitivePolicyTypeNameBatch
	CtPrimitivePolicyTypeNamePipeline
	CtPrimitivePolicyTypeNameAttachSingleVlan
	CtPrimitivePolicyTypeNameAttachMultipleVlan
	CtPrimitivePolicyTypeNameAttachLogicalLink
	CtPrimitivePolicyTypeNameAttachStaticRoute
	CtPrimitivePolicyTypeNameAttachCustomStaticRoute
	CtPrimitivePolicyTypeNameAttachIpEndpointWithBgpNsxt
	CtPrimitivePolicyTypeNameAttachBgpOverSubinterfacesOrSvi
	CtPrimitivePolicyTypeNameAttachBgpWithPrefixPeeringForSviOrSubinterface
	CtPrimitivePolicyTypeNameAttachExistingRoutingPolicy
	CtPrimitivePolicyTypeNameAttachRoutingZoneConstraint
	CtPrimitivePolicyTypeNameUnknown = "unknown CT primitive policy name %q"

	ctPrimitivePolicyTypeNameNone                                           = ctPrimitivePolicyTypeName("")
	ctPrimitivePolicyTypeNameBatch                                          = ctPrimitivePolicyTypeName("batch")
	ctPrimitivePolicyTypeNamePipeline                                       = ctPrimitivePolicyTypeName("pipeline")
	ctPrimitivePolicyTypeNameAttachSingleVlan                               = ctPrimitivePolicyTypeName("AttachSingleVLAN")
	ctPrimitivePolicyTypeNameAttachMultipleVlan                             = ctPrimitivePolicyTypeName("AttachMultipleVLAN")
	ctPrimitivePolicyTypeNameAttachLogicalLink                              = ctPrimitivePolicyTypeName("AttachLogicalLink")
	ctPrimitivePolicyTypeNameAttachStaticRoute                              = ctPrimitivePolicyTypeName("AttachStaticRoute")
	ctPrimitivePolicyTypeNameAttachCustomStaticRoute                        = ctPrimitivePolicyTypeName("AttachCustomStaticRoute")
	ctPrimitivePolicyTypeNameAttachIpEndpointWithBgpNsxt                    = ctPrimitivePolicyTypeName("AttachIpEndpointWithBgpNsxt")
	ctPrimitivePolicyTypeNameAttachBgpOverSubinterfacesOrSvi                = ctPrimitivePolicyTypeName("AttachBgpOverSubinterfacesOrSvi")
	ctPrimitivePolicyTypeNameAttachBgpWithPrefixPeeringForSviOrSubinterface = ctPrimitivePolicyTypeName("AttachBgpWithPrefixPeeringForSviOrSubinterface")
	ctPrimitivePolicyTypeNameAttachExistingRoutingPolicy                    = ctPrimitivePolicyTypeName("AttachExistingRoutingPolicy")
	ctPrimitivePolicyTypeNameAttachRoutingZoneConstraint                    = ctPrimitivePolicyTypeName("AttachRoutingZoneConstraint")
	ctPrimitivePolicyTypeNameUnknown                                        = "unknown CT primitive policy name %d"
)

func (o CtPrimitivePolicyTypeName) String() string {
	return string(o.raw())
}

func (o CtPrimitivePolicyTypeName) raw() ctPrimitivePolicyTypeName {
	switch o {
	case CtPrimitivePolicyTypeNameNone:
		return ctPrimitivePolicyTypeNameNone
	case CtPrimitivePolicyTypeNameBatch:
		return ctPrimitivePolicyTypeNameBatch
	case CtPrimitivePolicyTypeNamePipeline:
		return ctPrimitivePolicyTypeNamePipeline
	case CtPrimitivePolicyTypeNameAttachSingleVlan:
		return ctPrimitivePolicyTypeNameAttachSingleVlan
	case CtPrimitivePolicyTypeNameAttachMultipleVlan:
		return ctPrimitivePolicyTypeNameAttachMultipleVlan
	case CtPrimitivePolicyTypeNameAttachLogicalLink:
		return ctPrimitivePolicyTypeNameAttachLogicalLink
	case CtPrimitivePolicyTypeNameAttachStaticRoute:
		return ctPrimitivePolicyTypeNameAttachStaticRoute
	case CtPrimitivePolicyTypeNameAttachCustomStaticRoute:
		return ctPrimitivePolicyTypeNameAttachCustomStaticRoute
	case CtPrimitivePolicyTypeNameAttachIpEndpointWithBgpNsxt:
		return ctPrimitivePolicyTypeNameAttachIpEndpointWithBgpNsxt
	case CtPrimitivePolicyTypeNameAttachBgpOverSubinterfacesOrSvi:
		return ctPrimitivePolicyTypeNameAttachBgpOverSubinterfacesOrSvi
	case CtPrimitivePolicyTypeNameAttachBgpWithPrefixPeeringForSviOrSubinterface:
		return ctPrimitivePolicyTypeNameAttachBgpWithPrefixPeeringForSviOrSubinterface
	case CtPrimitivePolicyTypeNameAttachExistingRoutingPolicy:
		return ctPrimitivePolicyTypeNameAttachExistingRoutingPolicy
	case CtPrimitivePolicyTypeNameAttachRoutingZoneConstraint:
		return ctPrimitivePolicyTypeNameAttachRoutingZoneConstraint
	default:
		return ctPrimitivePolicyTypeName(fmt.Sprintf(ctPrimitivePolicyTypeNameUnknown, o))
	}
}

func (o *CtPrimitivePolicyTypeName) FromString(in string) error {
	i, err := ctPrimitivePolicyTypeName(in).parse()
	if err != nil {
		return err
	}
	*o = CtPrimitivePolicyTypeName(i)
	return nil
}

func (o ctPrimitivePolicyTypeName) parse() (int, error) {
	switch o {
	case ctPrimitivePolicyTypeNameNone:
		return int(CtPrimitivePolicyTypeNameNone), nil
	case ctPrimitivePolicyTypeNameBatch:
		return int(CtPrimitivePolicyTypeNameBatch), nil
	case ctPrimitivePolicyTypeNamePipeline:
		return int(CtPrimitivePolicyTypeNamePipeline), nil
	case ctPrimitivePolicyTypeNameAttachSingleVlan:
		return int(CtPrimitivePolicyTypeNameAttachSingleVlan), nil
	case ctPrimitivePolicyTypeNameAttachMultipleVlan:
		return int(CtPrimitivePolicyTypeNameAttachMultipleVlan), nil
	case ctPrimitivePolicyTypeNameAttachLogicalLink:
		return int(CtPrimitivePolicyTypeNameAttachLogicalLink), nil
	case ctPrimitivePolicyTypeNameAttachStaticRoute:
		return int(CtPrimitivePolicyTypeNameAttachStaticRoute), nil
	case ctPrimitivePolicyTypeNameAttachCustomStaticRoute:
		return int(CtPrimitivePolicyTypeNameAttachCustomStaticRoute), nil
	case ctPrimitivePolicyTypeNameAttachIpEndpointWithBgpNsxt:
		return int(CtPrimitivePolicyTypeNameAttachIpEndpointWithBgpNsxt), nil
	case ctPrimitivePolicyTypeNameAttachBgpOverSubinterfacesOrSvi:
		return int(CtPrimitivePolicyTypeNameAttachBgpOverSubinterfacesOrSvi), nil
	case ctPrimitivePolicyTypeNameAttachBgpWithPrefixPeeringForSviOrSubinterface:
		return int(CtPrimitivePolicyTypeNameAttachBgpWithPrefixPeeringForSviOrSubinterface), nil
	case ctPrimitivePolicyTypeNameAttachExistingRoutingPolicy:
		return int(CtPrimitivePolicyTypeNameAttachExistingRoutingPolicy), nil
	case ctPrimitivePolicyTypeNameAttachRoutingZoneConstraint:
		return int(CtPrimitivePolicyTypeNameAttachRoutingZoneConstraint), nil
	default:
		return 0, fmt.Errorf(CtPrimitivePolicyTypeNameUnknown, o)
	}
}

type (
	CtPrimitiveBgpPeerTo int
	ctPrimitiveBgpPeerTo string
)

const (
	CtPrimitiveBgpPeerToLoopback = CtPrimitiveBgpPeerTo(iota)
	CtPrimitiveBgpPeerToInterfaceOrIpEndpoint
	CtPrimitiveBgpPeerToInterfaceOrSharedIpEndpoint
	CtPrimitiveBgpPeerToUnknown = "unknown CtPrimitiveBgpPeerTo value %q"

	ctPrimitiveBgpPeerToLoopback                    = ctPrimitiveBgpPeerTo("loopback")
	ctPrimitiveBgpPeerToInterfaceOrIpEndpoint       = ctPrimitiveBgpPeerTo("interface_or_ip_endpoint")
	ctPrimitiveBgpPeerToInterfaceOrSharedIpEndpoint = ctPrimitiveBgpPeerTo("interface_or_shared_ip_endpoint")
	ctPrimitiveBgpPeerToUnknown                     = "unknown ctPrimitiveBgpPeerTo value %d"
)

func (o CtPrimitiveBgpPeerTo) String() string {
	switch o {
	case CtPrimitiveBgpPeerToLoopback:
		return string(ctPrimitiveBgpPeerToLoopback)
	case CtPrimitiveBgpPeerToInterfaceOrIpEndpoint:
		return string(ctPrimitiveBgpPeerToInterfaceOrIpEndpoint)
	case CtPrimitiveBgpPeerToInterfaceOrSharedIpEndpoint:
		return string(ctPrimitiveBgpPeerToInterfaceOrSharedIpEndpoint)
	default:
		return fmt.Sprintf(ctPrimitiveBgpPeerToUnknown, o)
	}
}

func (o *CtPrimitiveBgpPeerTo) FromString(in string) error {
	i, err := ctPrimitiveBgpPeerTo(in).parse()
	if err != nil {
		return err
	}
	*o = CtPrimitiveBgpPeerTo(i)
	return nil
}

func (o CtPrimitiveBgpPeerTo) raw() ctPrimitiveBgpPeerTo {
	return ctPrimitiveBgpPeerTo(o.String())
}

func (o ctPrimitiveBgpPeerTo) parse() (int, error) {
	switch o {
	case ctPrimitiveBgpPeerToLoopback:
		return int(CtPrimitiveBgpPeerToLoopback), nil
	case ctPrimitiveBgpPeerToInterfaceOrIpEndpoint:
		return int(CtPrimitiveBgpPeerToInterfaceOrIpEndpoint), nil
	case ctPrimitiveBgpPeerToInterfaceOrSharedIpEndpoint:
		return int(CtPrimitiveBgpPeerToInterfaceOrSharedIpEndpoint), nil
	default:
		return 0, fmt.Errorf(CtPrimitiveBgpPeerToUnknown, o)
	}
}

type (
	CtPrimitiveIPv4ProtocolSessionAddressing int
	ctPrimitiveIPv4ProtocolSessionAddressing string
)

const (
	CtPrimitiveIPv4ProtocolSessionAddressingNone = CtPrimitiveIPv4ProtocolSessionAddressing(iota)
	CtPrimitiveIPv4ProtocolSessionAddressingAddressed
	CtPrimitiveIPv4ProtocolSessionAddressingUnknown = "unknown CtPrimitiveIPv4ProtocolSessionAddressing value %q"

	ctPrimitiveIPv4ProtocolSessionAddressingNone      = ctPrimitiveIPv4ProtocolSessionAddressing("none")
	ctPrimitiveIPv4ProtocolSessionAddressingAddressed = ctPrimitiveIPv4ProtocolSessionAddressing("addressed")
	ctPrimitiveIPv4ProtocolSessionAddressingUnknown   = "unknown ctPrimitiveIPv4ProtocolSessionAddressing value %d"
)

func (o CtPrimitiveIPv4ProtocolSessionAddressing) String() string {
	switch o {
	case CtPrimitiveIPv4ProtocolSessionAddressingNone:
		return string(ctPrimitiveIPv4ProtocolSessionAddressingNone)
	case CtPrimitiveIPv4ProtocolSessionAddressingAddressed:
		return string(ctPrimitiveIPv4ProtocolSessionAddressingAddressed)
	default:
		return fmt.Sprintf(ctPrimitiveIPv4ProtocolSessionAddressingUnknown, o)
	}
}

func (o *CtPrimitiveIPv4ProtocolSessionAddressing) FromString(in string) error {
	i, err := ctPrimitiveIPv4ProtocolSessionAddressing(in).parse()
	if err != nil {
		return err
	}
	*o = CtPrimitiveIPv4ProtocolSessionAddressing(i)
	return nil
}

func (o CtPrimitiveIPv4ProtocolSessionAddressing) raw() ctPrimitiveIPv4ProtocolSessionAddressing {
	return ctPrimitiveIPv4ProtocolSessionAddressing(o.String())
}

func (o ctPrimitiveIPv4ProtocolSessionAddressing) parse() (int, error) {
	switch o {
	case ctPrimitiveIPv4ProtocolSessionAddressingNone:
		return int(CtPrimitiveIPv4ProtocolSessionAddressingNone), nil
	case ctPrimitiveIPv4ProtocolSessionAddressingAddressed:
		return int(CtPrimitiveIPv4ProtocolSessionAddressingAddressed), nil
	default:
		return 0, fmt.Errorf(CtPrimitiveIPv4ProtocolSessionAddressingUnknown, o)
	}
}

type (
	CtPrimitiveIPv6ProtocolSessionAddressing int
	ctPrimitiveIPv6ProtocolSessionAddressing string
)

const (
	CtPrimitiveIPv6ProtocolSessionAddressingNone = CtPrimitiveIPv6ProtocolSessionAddressing(iota)
	CtPrimitiveIPv6ProtocolSessionAddressingAddressed
	CtPrimitiveIPv6ProtocolSessionAddressingLinkLocal
	CtPrimitiveIPv6ProtocolSessionAddressingUnknown = "unknown CtPrimitiveIPv6ProtocolSessionAddressing value %q"

	ctPrimitiveIPv6ProtocolSessionAddressingNone      = ctPrimitiveIPv6ProtocolSessionAddressing("none")
	ctPrimitiveIPv6ProtocolSessionAddressingAddressed = ctPrimitiveIPv6ProtocolSessionAddressing("addressed")
	ctPrimitiveIPv6ProtocolSessionAddressingLinkLocal = ctPrimitiveIPv6ProtocolSessionAddressing("link_local")
	ctPrimitiveIPv6ProtocolSessionAddressingUnknown   = "unknown ctPrimitiveIPv6ProtocolSessionAddressing value %d"
)

func (o CtPrimitiveIPv6ProtocolSessionAddressing) String() string {
	switch o {
	case CtPrimitiveIPv6ProtocolSessionAddressingNone:
		return string(ctPrimitiveIPv6ProtocolSessionAddressingNone)
	case CtPrimitiveIPv6ProtocolSessionAddressingAddressed:
		return string(ctPrimitiveIPv6ProtocolSessionAddressingAddressed)
	case CtPrimitiveIPv6ProtocolSessionAddressingLinkLocal:
		return string(ctPrimitiveIPv6ProtocolSessionAddressingLinkLocal)
	default:
		return fmt.Sprintf(ctPrimitiveIPv6ProtocolSessionAddressingUnknown, o)
	}
}

func (o *CtPrimitiveIPv6ProtocolSessionAddressing) FromString(in string) error {
	i, err := ctPrimitiveIPv6ProtocolSessionAddressing(in).parse()
	if err != nil {
		return err
	}
	*o = CtPrimitiveIPv6ProtocolSessionAddressing(i)
	return nil
}

func (o CtPrimitiveIPv6ProtocolSessionAddressing) raw() ctPrimitiveIPv6ProtocolSessionAddressing {
	return ctPrimitiveIPv6ProtocolSessionAddressing(o.String())
}

func (o ctPrimitiveIPv6ProtocolSessionAddressing) parse() (int, error) {
	switch o {
	case ctPrimitiveIPv6ProtocolSessionAddressingNone:
		return int(CtPrimitiveIPv6ProtocolSessionAddressingNone), nil
	case ctPrimitiveIPv6ProtocolSessionAddressingAddressed:
		return int(CtPrimitiveIPv6ProtocolSessionAddressingAddressed), nil
	case ctPrimitiveIPv6ProtocolSessionAddressingLinkLocal:
		return int(CtPrimitiveIPv6ProtocolSessionAddressingLinkLocal), nil
	default:
		return 0, fmt.Errorf(CtPrimitiveIPv6ProtocolSessionAddressingUnknown, o)
	}
}

type (
	CtPrimitiveIPv4AddressingType int
	ctPrimitiveIPv4AddressingType string
)

const (
	CtPrimitiveIPv4AddressingTypeNone     = CtPrimitiveIPv4AddressingType(iota)
	CtPrimitiveIPv4AddressingTypeNumbered = CtPrimitiveIPv4AddressingType(iota)
	CtPrimitiveIPv4AddressingTypeUnknown  = "unknown CtPrimitiveIPv4AddressingType value %q"

	ctPrimitiveIPv4AddressingTypeNone     = ctPrimitiveIPv4AddressingType("none")
	ctPrimitiveIPv4AddressingTypeNumbered = ctPrimitiveIPv4AddressingType("numbered")
	ctPrimitiveIPv4AddressingTypeUnknown  = "unknown ctPrimitiveIPv4AddressingType value %d"
)

func (o CtPrimitiveIPv4AddressingType) String() string {
	switch o {
	case CtPrimitiveIPv4AddressingTypeNone:
		return string(ctPrimitiveIPv4AddressingTypeNone)
	case CtPrimitiveIPv4AddressingTypeNumbered:
		return string(ctPrimitiveIPv4AddressingTypeNumbered)
	default:
		return fmt.Sprintf(ctPrimitiveIPv4AddressingTypeUnknown, o)
	}
}

func (o *CtPrimitiveIPv4AddressingType) FromString(in string) error {
	i, err := ctPrimitiveIPv4AddressingType(in).parse()
	if err != nil {
		return err
	}
	*o = CtPrimitiveIPv4AddressingType(i)
	return nil
}

func (o CtPrimitiveIPv4AddressingType) raw() ctPrimitiveIPv4AddressingType {
	return ctPrimitiveIPv4AddressingType(o.String())
}

func (o ctPrimitiveIPv4AddressingType) parse() (int, error) {
	switch o {
	case ctPrimitiveIPv4AddressingTypeNone:
		return int(CtPrimitiveIPv4AddressingTypeNone), nil
	case ctPrimitiveIPv4AddressingTypeNumbered:
		return int(CtPrimitiveIPv4AddressingTypeNumbered), nil
	default:
		return 0, fmt.Errorf(CtPrimitiveIPv4AddressingTypeUnknown, o)
	}
}

type (
	CtPrimitiveIPv6AddressingType int
	ctPrimitiveIPv6AddressingType string
)

const (
	CtPrimitiveIPv6AddressingTypeLinkLocal = CtPrimitiveIPv6AddressingType(iota)
	CtPrimitiveIPv6AddressingTypeNone
	CtPrimitiveIPv6AddressingTypeNumbered
	CtPrimitiveIPv6AddressingTypeUnknown = "unknown CtPrimitiveIPv6AddressingType value %q"

	ctPrimitiveIPv6AddressingTypeLinkLocal = ctPrimitiveIPv6AddressingType("link_local")
	ctPrimitiveIPv6AddressingTypeNone      = ctPrimitiveIPv6AddressingType("none")
	ctPrimitiveIPv6AddressingTypeNumbered  = ctPrimitiveIPv6AddressingType("numbered")
	ctPrimitiveIPv6AddressingTypeUnknown   = "unknown ctPrimitiveIPv6AddressingType value %d"
)

func (o CtPrimitiveIPv6AddressingType) String() string {
	switch o {
	case CtPrimitiveIPv6AddressingTypeLinkLocal:
		return string(ctPrimitiveIPv6AddressingTypeLinkLocal)
	case CtPrimitiveIPv6AddressingTypeNone:
		return string(ctPrimitiveIPv6AddressingTypeNone)
	case CtPrimitiveIPv6AddressingTypeNumbered:
		return string(ctPrimitiveIPv6AddressingTypeNumbered)
	default:
		return fmt.Sprintf(ctPrimitiveIPv6AddressingTypeUnknown, o)
	}
}

func (o *CtPrimitiveIPv6AddressingType) FromString(in string) error {
	i, err := ctPrimitiveIPv6AddressingType(in).parse()
	if err != nil {
		return err
	}
	*o = CtPrimitiveIPv6AddressingType(i)
	return nil
}

func (o CtPrimitiveIPv6AddressingType) raw() ctPrimitiveIPv6AddressingType {
	return ctPrimitiveIPv6AddressingType(o.String())
}

func (o ctPrimitiveIPv6AddressingType) parse() (int, error) {
	switch o {
	case ctPrimitiveIPv6AddressingTypeLinkLocal:
		return int(CtPrimitiveIPv6AddressingTypeLinkLocal), nil
	case ctPrimitiveIPv6AddressingTypeNone:
		return int(CtPrimitiveIPv6AddressingTypeNone), nil
	case ctPrimitiveIPv6AddressingTypeNumbered:
		return int(CtPrimitiveIPv6AddressingTypeNumbered), nil
	default:
		return 0, fmt.Errorf(CtPrimitiveIPv6AddressingTypeUnknown, o)
	}
}

type (
	CtPrimitiveStatus int
	ctPrimitiveStatus string
)

const (
	CtPrimitiveStatusAssigned = CtPrimitiveStatus(iota)
	CtPrimitiveStatusIncomplete
	CtPrimitiveStatusReady
	CtPrimitiveStatusUnknown = "unknown CtPrimitiveStatus value %q"

	ctPrimitiveStatusAssigned   = ctPrimitiveStatus("assigned")
	ctPrimitiveStatusIncomplete = ctPrimitiveStatus("incomplete")
	ctPrimitiveStatusReady      = ctPrimitiveStatus("ready")
	ctPrimitiveStatusUnknown    = "unknown ctPrimitiveStatus value %d"
)

func (o CtPrimitiveStatus) String() string {
	switch o {
	case CtPrimitiveStatusAssigned:
		return string(ctPrimitiveStatusAssigned)
	case CtPrimitiveStatusIncomplete:
		return string(ctPrimitiveStatusIncomplete)
	case CtPrimitiveStatusReady:
		return string(ctPrimitiveStatusReady)
	default:
		return fmt.Sprintf(ctPrimitiveStatusUnknown, o)
	}
}

func (o *CtPrimitiveStatus) FromString(in string) error {
	i, err := ctPrimitiveStatus(in).parse()
	if err != nil {
		return err
	}
	*o = CtPrimitiveStatus(i)
	return nil
}

func (o CtPrimitiveStatus) raw() ctPrimitiveStatus {
	return ctPrimitiveStatus(o.String())
}

func (o ctPrimitiveStatus) parse() (int, error) {
	switch o {
	case ctPrimitiveStatusAssigned:
		return int(CtPrimitiveStatusAssigned), nil
	case ctPrimitiveStatusIncomplete:
		return int(CtPrimitiveStatusIncomplete), nil
	case ctPrimitiveStatusReady:
		return int(CtPrimitiveStatusReady), nil
	default:
		return 0, fmt.Errorf(CtPrimitiveStatusUnknown, o)
	}
}
