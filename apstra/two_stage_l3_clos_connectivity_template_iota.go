package apstra

import "fmt"

type CtPrimitivePolicyTypeName int
type ctPrimitivePolicyTypeName string

const (
	CtPrimitivePolicyTypeNameNone = CtPrimitivePolicyTypeName(iota)
	CtPrimitivePolicyTypeNameBatch
	CtPrimitivePolicyTypeNamePipeline
	CtPrimitivePolicyTypeNameAttachSingleVlan
	CtPrimitivePolicyTypeNameAttachMultipleVLAN
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
	ctPrimitivePolicyTypeNameAttachSingleVlan                               = ctPrimitivePolicyTypeName("AttachSingleVlan")
	ctPrimitivePolicyTypeNameAttachMultipleVLAN                             = ctPrimitivePolicyTypeName("AttachMultipleVLAN")
	ctPrimitivePolicyTypeNameAttachLogicalLink                              = ctPrimitivePolicyTypeName("AttachLogicalLink")
	ctPrimitivePolicyTypeNameAttachStaticRoute                              = ctPrimitivePolicyTypeName("AttachStaticRoute")
	ctPrimitivePolicyTypeNameAttachCustomStaticRoute                        = ctPrimitivePolicyTypeName("AttachCustomStaticRoute")
	ctPrimitivePolicyTypeNameAttachIpEndpointWithBgpNsxt                    = ctPrimitivePolicyTypeName("AttachIpEndpointWithBgpNsxt")
	ctPrimitivePolicyTypeNameAttachBgpOverSubinterfacesOrSvi                = ctPrimitivePolicyTypeName("AttachBgpOverSubinterfacesOrSvi")
	ctPrimitivePolicyTypeNameAttachBgpWithPrefixPeeringForSviOrSubinterface = ctPrimitivePolicyTypeName("BgpWithPrefixPeeringForSviOrSubinterface")
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
	case CtPrimitivePolicyTypeNameAttachMultipleVLAN:
		return ctPrimitivePolicyTypeNameAttachMultipleVLAN
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
	case ctPrimitivePolicyTypeNameAttachMultipleVLAN:
		return int(CtPrimitivePolicyTypeNameAttachMultipleVLAN), nil
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

type CtPrimitiveBgpPeerTo int
type ctPrimitiveBgpPeerTo string

const (
	CtPrimitiveBgpPeerToLoopback = CtPrimitiveBgpPeerTo(iota)
	CtPrimitiveBgpPeerToInterfaceOrIpEndpoint
	CtPrimitiveBgpPeerToInterfaceOrSharedIpEndpoint
	CtPrimitiveBgpPeerToInterfaceUnknown = "unknown CtPrimitiveBgpPeerTo value %q"

	ctPrimitiveBgpPeerToLoopback                    = ctPrimitiveBgpPeerTo("loopback")
	ctPrimitiveBgpPeerToInterfaceOrIpEndpoint       = ctPrimitiveBgpPeerTo("interface_or_ip_endpoint")
	ctPrimitiveBgpPeerToInterfaceOrSharedIpEndpoint = ctPrimitiveBgpPeerTo("interface_or_shared_ip_endpoint")
	ctPrimitiveBgpPeerToInterfaceUnknown            = "unknown ctPrimitiveBgpPeerTo value %d"
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
		return fmt.Sprintf(ctPrimitiveBgpPeerToInterfaceUnknown, o)
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
		return 0, fmt.Errorf(CtPrimitiveBgpPeerToInterfaceUnknown, o)
	}
}

type CtPrimitiveIPv4ProtocolSessionAddressing int
type ctPrimitiveIPv4ProtocolSessionAddressing string

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

type CtPrimitiveIPv6ProtocolSessionAddressing int
type ctPrimitiveIPv6ProtocolSessionAddressing string

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
	default:
		return 0, fmt.Errorf(CtPrimitiveIPv6ProtocolSessionAddressingUnknown, o)
	}
}
