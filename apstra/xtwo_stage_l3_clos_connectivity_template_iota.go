package apstra

import "fmt"

type CtPrimitiveBgpPeerTo int
type ctPrimitiveBgpPeerTo string

const (
	CtPrimitiveBgpPeerToLoopback = CtPrimitiveBgpPeerTo(iota)
	CtPrimitiveBgpPeerToInterfaceOrIpEndpoint
	CtPrimitiveBgpPeerToInterfaceOrSharedIpEndpoint
	CtPrimitiveBgpPeerToInterfaceUnknown = "unknown CtPrimitiveBgpPeerTo value %d"

	ctPrimitiveBgpPeerToLoopback                    = ctPrimitiveBgpPeerTo("loopback")
	ctPrimitiveBgpPeerToInterfaceOrIpEndpoint       = ctPrimitiveBgpPeerTo("interface_or_ip_endpoint")
	ctPrimitiveBgpPeerToInterfaceOrSharedIpEndpoint = ctPrimitiveBgpPeerTo("interface_or_shared_ip_endpoint")
	ctPrimitiveBgpPeerToInterfaceUnknown            = "unknown ctPrimitiveBgpPeerTo value %q"
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
		return fmt.Sprintf(CtPrimitiveBgpPeerToInterfaceUnknown, o)
	}
}

func (o CtPrimitiveBgpPeerTo) raw() ctPrimitiveBgpPeerTo {
	return ctPrimitiveBgpPeerTo(o.String())
}

type CtPrimitiveIPv4ProtocolSessionAddressing int
type ctPrimitiveIPv4ProtocolSessionAddressing string

const (
	CtPrimitiveIPv4ProtocolSessionAddressingNone = CtPrimitiveIPv4ProtocolSessionAddressing(iota)
	CtPrimitiveIPv4ProtocolSessionAddressingAddressed
	CtPrimitiveIPv4ProtocolSessionAddressingUnknown = "unknown CtPrimitiveIPv4ProtocolSessionAddressing value %d"

	ctPrimitiveIPv4ProtocolSessionAddressingNone      = ctPrimitiveIPv4ProtocolSessionAddressing("none")
	ctPrimitiveIPv4ProtocolSessionAddressingAddressed = ctPrimitiveIPv4ProtocolSessionAddressing("addressed")
	ctPrimitiveIPv4ProtocolSessionAddressingUnknown   = "unknown ctPrimitiveIPv4ProtocolSessionAddressing value %q"
)

func (o CtPrimitiveIPv4ProtocolSessionAddressing) String() string {
	switch o {
	case CtPrimitiveIPv4ProtocolSessionAddressingNone:
		return string(ctPrimitiveIPv4ProtocolSessionAddressingNone)
	case CtPrimitiveIPv4ProtocolSessionAddressingAddressed:
		return string(ctPrimitiveIPv4ProtocolSessionAddressingAddressed)
	default:
		return fmt.Sprintf(CtPrimitiveIPv4ProtocolSessionAddressingUnknown, o)
	}
}

func (o CtPrimitiveIPv4ProtocolSessionAddressing) raw() ctPrimitiveIPv4ProtocolSessionAddressing {
	return ctPrimitiveIPv4ProtocolSessionAddressing(o.String())
}

type CtPrimitiveIPv6ProtocolSessionAddressing int
type ctPrimitiveIPv6ProtocolSessionAddressing string

const (
	CtPrimitiveIPv6ProtocolSessionAddressingNone = CtPrimitiveIPv6ProtocolSessionAddressing(iota)
	CtPrimitiveIPv6ProtocolSessionAddressingAddressed
	CtPrimitiveIPv6ProtocolSessionAddressingLinkLocal
	CtPrimitiveIPv6ProtocolSessionAddressingUnknown = "unknown CtPrimitiveIPv6ProtocolSessionAddressing value %d"

	ctPrimitiveIPv6ProtocolSessionAddressingNone      = ctPrimitiveIPv6ProtocolSessionAddressing("none")
	ctPrimitiveIPv6ProtocolSessionAddressingAddressed = ctPrimitiveIPv6ProtocolSessionAddressing("addressed")
	ctPrimitiveIPv6ProtocolSessionAddressingLinkLocal = ctPrimitiveIPv6ProtocolSessionAddressing("link_local")
	ctPrimitiveIPv6ProtocolSessionAddressingUnknown   = "unknown ctPrimitiveIPv6ProtocolSessionAddressing value %q"
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
		return fmt.Sprintf(CtPrimitiveIPv6ProtocolSessionAddressingUnknown, o)
	}
}

func (o CtPrimitiveIPv6ProtocolSessionAddressing) raw() ctPrimitiveIPv6ProtocolSessionAddressing {
	return ctPrimitiveIPv6ProtocolSessionAddressing(o.String())
}
