package apstra

import "encoding/json"

type xConnectivityTemplateAttributes interface {
	raw() (json.RawMessage, error)
	policyTypeName() string
	label() string
	description() string
}

// AttachSingleVlan
var _ xConnectivityTemplateAttributes = ConnectivityTemplatePrimitiveAttributesAttachSingleVlan{}

type ConnectivityTemplatePrimitiveAttributesAttachSingleVlan struct {
	Tagged   bool
	VnNodeId *ObjectId
}

func (o ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) raw() (json.RawMessage, error) {
	var tagType string
	if o.Tagged {
		tagType = "vlan_tagged"
	} else {
		tagType = "untagged"
	}

	raw := struct {
		TagType  string    `json:"tag_type"`
		VnNodeId *ObjectId `json:"vn_node_id"`
	}{
		TagType:  tagType,
		VnNodeId: o.VnNodeId,
	}

	return json.Marshal(&raw)
}

func (o ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) policyTypeName() string {
	return "AttachSingleVlan"
}

func (o ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) label() string {
	return "Virtual Network (Single)"
}

func (o ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) description() string {
	return "Add a single VLAN to interfaces, as tagged or untagged."
}

// AttachMultipleVLAN
var _ xConnectivityTemplateAttributes = ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan{}

type ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan struct {
	UntaggedVnNodeId *ObjectId
	TaggedVnNodeIds  []ObjectId
}

func (o ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan) raw() (json.RawMessage, error) {
	raw := struct {
		UntaggedVnNodeId *ObjectId  `json:"untagged_vn_node_id"`
		TaggedVnNodeIds  []ObjectId `json:"tagged_vn_node_ids"`
	}{
		UntaggedVnNodeId: o.UntaggedVnNodeId,
		TaggedVnNodeIds:  o.TaggedVnNodeIds,
	}

	return json.Marshal(&raw)
}

func (o ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan) policyTypeName() string {
	return "AttachMultipleVLAN"
}

func (o ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan) label() string {
	return "Virtual Network (Multiple)"
}

func (o ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan) description() string {
	return "Add a list of VLANs to interfaces, as tagged or untagged."
}

// AttachLogicalLink
var _ xConnectivityTemplateAttributes = ConnectivityTemplatePrimitiveAttributesAttachLogicalLink{}

type ConnectivityTemplatePrimitiveAttributesAttachLogicalLink struct {
	SecurityZone            *ObjectId
	Tagged                  bool
	Vlan                    *Vlan
	IPv4AddressingNumbered  bool
	IPv6AddressingLinkLocal bool
}

func (o ConnectivityTemplatePrimitiveAttributesAttachLogicalLink) raw() (json.RawMessage, error) {
	var interfaceType string
	switch o.Tagged {
	case true:
		interfaceType = "tagged"
	case false:
		interfaceType = "untagged"
	}

	if o.Vlan != nil {
		err := o.Vlan.validate()
		if err != nil {
			return nil, err
		}
	}

	var IPv4AddressingType string
	switch o.IPv4AddressingNumbered {
	case true:
		IPv4AddressingType = "numbered"
	case false:
		IPv4AddressingType = "none"
	}

	var IPv6AddressingType string
	switch o.IPv6AddressingLinkLocal {
	case true:
		IPv6AddressingType = "link_local"
	case false:
		IPv6AddressingType = "none"
	}

	raw := struct {
		SecurityZone       *ObjectId `json:"security_zone"`
		InterfaceType      string    `json:"interface_type"`
		VlanId             *Vlan     `json:"vlan_id"`
		IPv4AddressingType string    `json:"ipv4_addressing_type"`
		IPv6AddressingType string    `json:"ipv6_addressing_type"`
	}{
		SecurityZone:       o.SecurityZone,
		InterfaceType:      interfaceType,
		VlanId:             o.Vlan,
		IPv4AddressingType: IPv4AddressingType,
		IPv6AddressingType: IPv6AddressingType,
	}

	return json.Marshal(&raw)
}

func (o ConnectivityTemplatePrimitiveAttributesAttachLogicalLink) policyTypeName() string {
	return "AttachLogicalLink"
}

func (o ConnectivityTemplatePrimitiveAttributesAttachLogicalLink) label() string {
	return "IP Link"
}

func (o ConnectivityTemplatePrimitiveAttributesAttachLogicalLink) description() string {
	return "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table. Can be assigned to physical interfaces or single-chassis LAGs (not applicable to ESI LAG or MLAG interfaces)."
}

//	"AttachStaticRoute"
// "Static Route"
// "Create a static route to user defined subnet via next hop derived from either IP link or VN endpoint."

//	"AttachCustomStaticRoute"
// "Custom Static Route"
// "Create a static route with user defined next hop and destination network."

//	"AttachIpEndpointWithBgpNsxt"
// "BGP Peering (IP Endpoint)"
// "Create a BGP peering session with a user-specified BGP neighbor addressed peer."

// AttachBgpOverSubinterfacesOrSvi
var _ xConnectivityTemplateAttributes = ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi{}

type ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi struct {
	Ipv4Safi              bool
	Ipv6Safi              bool
	TTL                   uint8
	BFD                   bool
	Password              *string
	Keepalive             *uint16
	Holdtime              *uint16
	SessionAddressingIpv4 CtPrimitiveIPv4ProtocolSessionAddressing
	SessionAddressingIpv6 CtPrimitiveIPv6ProtocolSessionAddressing
	LocalAsn              *uint32
	PeerFromLoopback      bool
	PeerTo                CtPrimitiveBgpPeerTo
	NeighborAsnDynamic    bool // 'static', 'dynamic'
}

func (o ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi) raw() (json.RawMessage, error) {
	var peerFrom string
	switch o.PeerFromLoopback {
	case true:
		peerFrom = "loopback"
	case false:
		peerFrom = "interface"
	}

	var neighborAsnType string
	switch o.NeighborAsnDynamic {
	case true:
		neighborAsnType = "dynamic"
	case false:
		neighborAsnType = "static"
	}

	raw := struct {
		Ipv4Safi              bool                                     `json:"ipv4_safi"`
		Ipv6Safi              bool                                     `json:"ipv6_safi"`
		TTL                   uint8                                    `json:"ttl"`
		BFD                   bool                                     `json:"bfd"`
		Password              *string                                  `json:"password"`
		Keepalive             *uint16                                  `json:"keepalive_timer"`
		Holdtime              *uint16                                  `json:"holdtime_timer"`
		LocalAsn              *uint32                                  `json:"local_asn"`
		NeighborAsnType       string                                   `json:"neighbor_asn_type"`
		PeerFrom              string                                   `json:"peer_from"`
		PeerTo                ctPrimitiveBgpPeerTo                     `json:"peer_to"`
		SessionAddressingIpv4 ctPrimitiveIPv4ProtocolSessionAddressing `json:"session_addressing_ipv4"`
		SessionAddressingIpv6 ctPrimitiveIPv6ProtocolSessionAddressing `json:"session_addressing_ipv6"`
	}{
		Ipv4Safi:              o.Ipv4Safi,
		Ipv6Safi:              o.Ipv6Safi,
		TTL:                   o.TTL,
		BFD:                   o.BFD,
		Password:              o.Password,
		Keepalive:             o.Keepalive,
		Holdtime:              o.Holdtime,
		LocalAsn:              o.LocalAsn,
		NeighborAsnType:       neighborAsnType,
		PeerFrom:              peerFrom,
		PeerTo:                o.PeerTo.raw(),
		SessionAddressingIpv4: o.SessionAddressingIpv4.raw(),
		SessionAddressingIpv6: o.SessionAddressingIpv6.raw(),
	}

	return json.Marshal(&raw)
}

func (o ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi) policyTypeName() string {
	return "AttachBgpOverSubinterfacesOrSvi"
}

func (o ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi) label() string {
	return "BGP Peering (Generic System)"
}

func (o ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi) description() string {
	return "Create a BGP peering session with Generic Systems inherited from AOS Generic System properties such as loopback and ASN (addressed, or link-local peer). Static route is automatically created when selecting loopback peering."
}

//	"AttachBgpWithPrefixPeeringForSviOrSubinterface"
// "Dynamic BGP Peering"
// "Configure dynamic BGP peering with IP prefix specified."

//	"AttachExistingRoutingPolicy"
// "Routing Policy"
// "Allocate routing policy to specific BGP sessions."

// "AttachRoutingZoneConstraint"
// "Routing Zone Constraint"
// "Assign a Routing Zone Constraint"
