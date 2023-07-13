package apstra

import (
	"encoding/json"
	"fmt"
	"net"
)

// connectivityTemplateAttributes are the data structures which make the various
// CT primitives (single VLAN, multiple VLAN, static route, etc...) different
// from each other. In Apstra 4.1.2 there are 10 CT primitives, so there are 10
// implementations of the connectivityTemplateAttributes interface.
type connectivityTemplateAttributes interface {
	raw() (json.RawMessage, error)
	policyTypeName() CtPrimitivePolicyTypeName
	label() string
	description() string
	fromRawJson(json.RawMessage) error
}

// AttachSingleVlan
var _ connectivityTemplateAttributes = &ConnectivityTemplatePrimitiveAttributesAttachSingleVlan{}

type ConnectivityTemplatePrimitiveAttributesAttachSingleVlan struct {
	Tagged   bool
	VnNodeId *ObjectId
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) fromRawJson(in json.RawMessage) error {
	var raw rawConnectivityTemplatePrimitiveAttributesAttachSingleVlan
	err := json.Unmarshal(in, &raw)
	if err != nil {
		return err
	}

	return raw.polish(o)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) raw() (json.RawMessage, error) {
	var tagType string
	if o.Tagged {
		tagType = "vlan_tagged"
	} else {
		tagType = "untagged"
	}

	raw := rawConnectivityTemplatePrimitiveAttributesAttachSingleVlan{
		TagType:  tagType,
		VnNodeId: o.VnNodeId,
	}

	return json.Marshal(&raw)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) policyTypeName() CtPrimitivePolicyTypeName {
	return CtPrimitivePolicyTypeNameAttachSingleVlan
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) label() string {
	return "Virtual Network (Single)"
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) description() string {
	return "Add a single VLAN to interfaces, as tagged or untagged."
}

// AttachMultipleVLAN
var _ connectivityTemplateAttributes = &ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan{}

type ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan struct {
	UntaggedVnNodeId *ObjectId
	TaggedVnNodeIds  []ObjectId
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan) fromRawJson(in json.RawMessage) error {
	var raw rawConnectivityTemplatePrimitiveAttributesAttachMultipleVlan
	err := json.Unmarshal(in, &raw)
	if err != nil {
		return err
	}

	return raw.polish(o)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan) raw() (json.RawMessage, error) {
	raw := rawConnectivityTemplatePrimitiveAttributesAttachMultipleVlan{
		UntaggedVnNodeId: o.UntaggedVnNodeId,
		TaggedVnNodeIds:  o.TaggedVnNodeIds,
	}

	return json.Marshal(&raw)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan) policyTypeName() CtPrimitivePolicyTypeName {
	return CtPrimitivePolicyTypeNameAttachMultipleVLAN
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan) label() string {
	return "Virtual Network (Multiple)"
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan) description() string {
	return "Add a list of VLANs to interfaces, as tagged or untagged."
}

// AttachLogicalLink
var _ connectivityTemplateAttributes = &ConnectivityTemplatePrimitiveAttributesAttachLogicalLink{}

type ConnectivityTemplatePrimitiveAttributesAttachLogicalLink struct {
	SecurityZone            *ObjectId
	Tagged                  bool
	Vlan                    *Vlan
	IPv4AddressingNumbered  bool
	IPv6AddressingLinkLocal bool
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachLogicalLink) fromRawJson(in json.RawMessage) error {
	var raw rawConnectivityTemplatePrimitiveAttributesAttachLogicalLink
	err := json.Unmarshal(in, &raw)
	if err != nil {
		return err
	}

	return raw.polish(o)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachLogicalLink) raw() (json.RawMessage, error) {
	var intfType string
	switch o.Tagged {
	case true:
		intfType = "tagged"
	case false:
		intfType = "untagged"
	}

	if o.Vlan != nil {
		err := o.Vlan.validate()
		if err != nil {
			return nil, err
		}
	}

	var iPv4AddressingType string
	switch o.IPv4AddressingNumbered {
	case true:
		iPv4AddressingType = "numbered"
	case false:
		iPv4AddressingType = "none"
	}

	var iPv6AddressingType string
	switch o.IPv6AddressingLinkLocal {
	case true:
		iPv6AddressingType = "link_local"
	case false:
		iPv6AddressingType = "none"
	}

	raw := rawConnectivityTemplatePrimitiveAttributesAttachLogicalLink{
		InterfaceType:      intfType,
		Vlan:               o.Vlan,
		Ipv4AddressingType: iPv4AddressingType,
		Ipv6AddressingType: iPv6AddressingType,
		SecurityZone:       o.SecurityZone,
	}

	return json.Marshal(&raw)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachLogicalLink) policyTypeName() CtPrimitivePolicyTypeName {
	return CtPrimitivePolicyTypeNameAttachLogicalLink
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachLogicalLink) label() string {
	return "IP Link"
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachLogicalLink) description() string {
	return "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table."
}

// AttachStaticRoute
var _ connectivityTemplateAttributes = &ConnectivityTemplatePrimitiveAttributesAttachStaticRoute{}

type ConnectivityTemplatePrimitiveAttributesAttachStaticRoute struct {
	ShareIpEndpoint bool
	Network         *net.IPNet
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachStaticRoute) raw() (json.RawMessage, error) {
	var network *string
	if o.Network != nil {
		s := o.Network.String()
		network = &s
	}

	raw := rawConnectivityTemplatePrimitiveAttributesAttachStaticRoute{
		ShareIpEndpoint: o.ShareIpEndpoint,
		Network:         network,
	}

	return json.Marshal(&raw)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachStaticRoute) policyTypeName() CtPrimitivePolicyTypeName {
	return CtPrimitivePolicyTypeNameAttachStaticRoute
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachStaticRoute) label() string {
	return "Static Route"
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachStaticRoute) description() string {
	return "Create a static route to user defined subnet via next hop derived from either IP link or VN endpoint."
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachStaticRoute) fromRawJson(in json.RawMessage) error {
	var raw rawConnectivityTemplatePrimitiveAttributesAttachStaticRoute
	err := json.Unmarshal(in, &raw)
	if err != nil {
		return err
	}

	return raw.polish(o)
}

// AttachCustomStaticRoute
var _ connectivityTemplateAttributes = &ConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute{}

type ConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute struct {
	Network      *net.IPNet
	NextHop      net.IP
	SecurityZone *ObjectId
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute) raw() (json.RawMessage, error) {
	var network, nexthop *string

	if o.Network != nil {
		s := o.Network.String()
		network = &s
	}

	if o.NextHop != nil {
		s := o.NextHop.String()
		nexthop = &s
	}

	raw := rawConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute{
		Network:      network,
		NextHop:      nexthop,
		SecurityZone: o.SecurityZone,
	}

	return json.Marshal(&raw)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute) policyTypeName() CtPrimitivePolicyTypeName {
	return CtPrimitivePolicyTypeNameAttachCustomStaticRoute
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute) label() string {
	return "Custom Static Route"
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute) description() string {
	return "Create a static route with user defined next hop and destination network."
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute) fromRawJson(in json.RawMessage) error {
	var raw rawConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute
	err := json.Unmarshal(in, &raw)
	if err != nil {
		return err
	}

	return raw.polish(o)
}

// AttachIpEndpointWithBgpNsxt
var _ connectivityTemplateAttributes = &ConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt{}

type ConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt struct {
	Asn                *uint32
	Bfd                bool
	Holdtime           *uint16
	Ipv4Addr           net.IP
	Ipv6Addr           net.IP
	Ipv4Safi           bool
	Ipv6Safi           bool
	Keepalive          *uint16
	LocalAsn           *uint32
	NeighborAsnDynamic bool
	Password           *string
	Ttl                uint8
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt) fromRawJson(in json.RawMessage) error {
	var raw rawConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt
	err := json.Unmarshal(in, &raw)
	if err != nil {
		return err
	}

	return raw.polish(o)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt) raw() (json.RawMessage, error) {
	var ipv4Addr, ipv6Addr *string

	if len(o.Ipv4Addr) != 0 {
		s := o.Ipv4Addr.String()
		ipv4Addr = &s
	}

	if len(o.Ipv6Addr) != 0 {
		s := o.Ipv6Addr.String()
		ipv6Addr = &s
	}

	var neighborAsnType string
	if o.NeighborAsnDynamic {
		neighborAsnType = "dynamic"
	} else {
		neighborAsnType = "static"
	}

	raw := rawConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt{
		Asn:             o.Asn,
		Bfd:             o.Bfd,
		Holdtime:        o.Holdtime,
		Ipv4Addr:        ipv4Addr,
		Ipv4Safi:        o.Ipv4Safi,
		Ipv6Addr:        ipv6Addr,
		Ipv6Safi:        o.Ipv6Safi,
		Keepalive:       o.Keepalive,
		LocalAsn:        o.LocalAsn,
		NeighborAsnType: neighborAsnType,
		Password:        o.Password,
		Ttl:             o.Ttl,
	}

	return json.Marshal(&raw)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt) policyTypeName() CtPrimitivePolicyTypeName {
	return CtPrimitivePolicyTypeNameAttachIpEndpointWithBgpNsxt
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt) label() string {
	return "BGP Peering (IP Endpoint)"
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt) description() string {
	return "Create a BGP peering session with a user-specified BGP neighbor addressed peer."
}

// AttachBgpOverSubinterfacesOrSvi
var _ connectivityTemplateAttributes = &ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi{}

type ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi struct {
	Bfd                   bool
	Holdtime              *uint16
	Ipv4Safi              bool
	Ipv6Safi              bool
	Keepalive             *uint16
	LocalAsn              *uint32
	NeighborAsnDynamic    bool // 'static', 'dynamic'
	Password              *string
	PeerFromLoopback      bool
	PeerTo                CtPrimitiveBgpPeerTo
	SessionAddressingIpv4 CtPrimitiveIPv4ProtocolSessionAddressing
	SessionAddressingIpv6 CtPrimitiveIPv6ProtocolSessionAddressing
	Ttl                   uint8
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi) fromRawJson(in json.RawMessage) error {
	var raw rawConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi
	err := json.Unmarshal(in, &raw)
	if err != nil {
		return err
	}

	return json.Unmarshal(in, o)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi) raw() (json.RawMessage, error) {
	var neighborAsnType string
	switch o.NeighborAsnDynamic {
	case true:
		neighborAsnType = "dynamic"
	case false:
		neighborAsnType = "static"
	}

	var peerFrom string
	switch o.PeerFromLoopback {
	case true:
		peerFrom = "loopback"
	case false:
		peerFrom = "interface"
	}

	raw := rawConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi{
		Bfd:                   o.Bfd,
		Holdtime:              o.Holdtime,
		Ipv4Safi:              o.Ipv4Safi,
		Ipv6Safi:              o.Ipv6Safi,
		Keepalive:             o.Keepalive,
		LocalAsn:              o.LocalAsn,
		NeighborAsnType:       neighborAsnType,
		Password:              o.Password,
		PeerFrom:              peerFrom,
		PeerTo:                o.PeerTo.raw(),
		SessionAddressingIpv4: o.SessionAddressingIpv4.raw(),
		SessionAddressingIpv6: o.SessionAddressingIpv6.raw(),
		Ttl:                   o.Ttl,
	}

	return json.Marshal(&raw)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi) policyTypeName() CtPrimitivePolicyTypeName {
	return CtPrimitivePolicyTypeNameAttachBgpOverSubinterfacesOrSvi
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi) label() string {
	return "BGP Peering (Generic System)"
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi) description() string {
	return "Create a BGP peering session with Generic Systems inherited from AOS Generic System properties such as loopback and ASN (addressed, or link-local peer)."
}

// AttachBgpWithPrefixPeeringForSviOrSubinterface
var _ connectivityTemplateAttributes = &ConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface{}

type ConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface struct {
	Bfd                   bool
	Holdtime              *uint16
	Ipv4Safi              bool
	Ipv6Safi              bool
	Keepalive             *uint16
	LocalAsn              *uint32
	Password              *string
	PrefixNeighborIpv4    *net.IPNet
	PrefixNeighborIpv6    *net.IPNet
	SessionAddressingIpv4 bool
	SessionAddressingIpv6 bool
	Ttl                   int
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface) raw() (json.RawMessage, error) {
	var prefixNeighborIpv4, prefixNeighborIpv6 *string

	if o.PrefixNeighborIpv4 != nil {
		s := o.PrefixNeighborIpv4.String()
		prefixNeighborIpv4 = &s
	}

	if o.PrefixNeighborIpv6 != nil {
		s := o.PrefixNeighborIpv6.String()
		prefixNeighborIpv6 = &s
	}

	raw := rawConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface{
		Bfd:                   o.Bfd,
		Holdtime:              o.Holdtime,
		Ipv4Safi:              o.Ipv4Safi,
		Ipv6Safi:              o.Ipv6Safi,
		Keepalive:             o.Keepalive,
		LocalAsn:              o.LocalAsn,
		Password:              o.Password,
		PrefixNeighborIpv4:    prefixNeighborIpv4,
		PrefixNeighborIpv6:    prefixNeighborIpv6,
		SessionAddressingIpv4: o.SessionAddressingIpv4,
		SessionAddressingIpv6: o.SessionAddressingIpv6,
		Ttl:                   o.Ttl,
	}

	return json.Marshal(&raw)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface) policyTypeName() CtPrimitivePolicyTypeName {
	return CtPrimitivePolicyTypeNameAttachBgpWithPrefixPeeringForSviOrSubinterface
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface) label() string {
	return "Dynamic BGP Peering"
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface) description() string {
	return "Configure dynamic BGP peering with IP prefix specified."
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface) fromRawJson(in json.RawMessage) error {
	var raw rawConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface
	err := json.Unmarshal(in, &raw)
	if err != nil {
		return err
	}

	return raw.polish(o)
}

// AttachExistingRoutingPolicy
var _ connectivityTemplateAttributes = &ConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy{}

type ConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy struct {
	RpToAttach *string `json:"rp_to_attach"`
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy) fromRawJson(in json.RawMessage) error {
	var raw rawConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy
	err := json.Unmarshal(in, &raw)
	if err != nil {
		return err
	}

	return raw.polish(o)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy) raw() (json.RawMessage, error) {
	raw := rawConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy{
		RpToAttach: o.RpToAttach,
	}
	return json.Marshal(&raw)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy) policyTypeName() CtPrimitivePolicyTypeName {
	return CtPrimitivePolicyTypeNameAttachExistingRoutingPolicy
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy) label() string {
	return "Routing Policy"
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy) description() string {
	return "Allocate routing policy to specific BGP sessions."
}

// AttachRoutingZoneConstraint
var _ connectivityTemplateAttributes = &ConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint{}

type ConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint struct {
	RoutingZoneConstraint *string
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint) raw() (json.RawMessage, error) {
	raw := rawConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint{
		RoutingZoneConstraint: o.RoutingZoneConstraint,
	}

	return json.Marshal(&raw)
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint) policyTypeName() CtPrimitivePolicyTypeName {
	return CtPrimitivePolicyTypeNameAttachRoutingZoneConstraint
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint) label() string {
	return "Routing Zone Constraint"
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint) description() string {
	return "Assign a Routing Zone Constraint"
}

func (o *ConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint) fromRawJson(in json.RawMessage) error {
	var raw rawConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint
	err := json.Unmarshal(in, &raw)
	if err != nil {
		return err
	}

	return raw.polish(o)
}

// Each implementation of connectivityTemplateAttributes needs a "raw" struct
// with JSON tags wire-style elements. The 10 "raw" structs follow, each with a
// `polish()` method. Note that rather than returning a polished struct (or
// pointer), these methods polish into an existing struct referenced by a caller
// supplied pointer.
type rawConnectivityTemplatePrimitiveAttributesAttachSingleVlan struct {
	VnNodeId *ObjectId `json:"vn_node_id"`
	TagType  string    `json:"tag_type"`
}

func (o rawConnectivityTemplatePrimitiveAttributesAttachSingleVlan) polish(t *ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) error {
	var tagged bool
	switch o.TagType {
	case "vlan_tagged":
		tagged = true
	case "untagged":
		tagged = false
	default:
		return fmt.Errorf("unexpected tag_type %q", o.TagType)
	}

	t.Tagged = tagged
	t.VnNodeId = o.VnNodeId

	return nil
}

type rawConnectivityTemplatePrimitiveAttributesAttachMultipleVlan struct {
	UntaggedVnNodeId *ObjectId  `json:"untagged_vn_node_id"`
	TaggedVnNodeIds  []ObjectId `json:"tagged_vn_node_ids"`
}

func (o rawConnectivityTemplatePrimitiveAttributesAttachMultipleVlan) polish(t *ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan) error {
	t.UntaggedVnNodeId = o.UntaggedVnNodeId
	t.TaggedVnNodeIds = o.TaggedVnNodeIds

	return nil
}

type rawConnectivityTemplatePrimitiveAttributesAttachLogicalLink struct {
	InterfaceType      string    `json:"interface_type"`
	Vlan               *Vlan     `json:"vlan_id"`
	Ipv6AddressingType string    `json:"ipv6_addressing_type"`
	Ipv4AddressingType string    `json:"ipv4_addressing_type"`
	SecurityZone       *ObjectId `json:"security_zone"`
}

func (o rawConnectivityTemplatePrimitiveAttributesAttachLogicalLink) polish(t *ConnectivityTemplatePrimitiveAttributesAttachLogicalLink) error {
	var tagged bool
	switch o.InterfaceType {
	case "tagged":
		tagged = true
	case "untagged":
		tagged = false
	case "":
		tagged = false
	default:
		return fmt.Errorf("unexpected interfaceType %q", o.InterfaceType)
	}

	var ipv4Numbered bool
	switch o.Ipv4AddressingType {
	case "numbered":
		ipv4Numbered = true
	case "none":
		ipv4Numbered = false
	default:
		return fmt.Errorf("unexpected ipv4_addressing_type %q", o.Ipv4AddressingType)
	}

	var ipv6LinkLocal bool
	switch o.Ipv6AddressingType {
	case "link_local":
		ipv6LinkLocal = true
	case "none":
		ipv6LinkLocal = false
	default:
		return fmt.Errorf("unexpected ipv6_addressing_type %q", o.Ipv6AddressingType)
	}

	t.SecurityZone = o.SecurityZone
	t.Tagged = tagged
	t.Vlan = o.Vlan
	t.IPv4AddressingNumbered = ipv4Numbered
	t.IPv6AddressingLinkLocal = ipv6LinkLocal

	return nil
}

type rawConnectivityTemplatePrimitiveAttributesAttachStaticRoute struct {
	ShareIpEndpoint bool    `json:"share_ip_endpoint"`
	Network         *string `json:"network"`
}

func (o rawConnectivityTemplatePrimitiveAttributesAttachStaticRoute) polish(t *ConnectivityTemplatePrimitiveAttributesAttachStaticRoute) error {
	var network *net.IPNet

	if o.Network != nil {
		var err error
		_, network, err = net.ParseCIDR(*o.Network)
		if err != nil {
			return err
		}
	}

	t.ShareIpEndpoint = o.ShareIpEndpoint
	t.Network = network

	return nil
}

type rawConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute struct {
	Network      *string   `json:"network"`
	NextHop      *string   `json:"next_hop"`
	SecurityZone *ObjectId `json:"security_zone"`
}

func (o rawConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute) polish(t *ConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute) error {
	var network *net.IPNet
	var nextHop net.IP

	if o.Network != nil {
		var err error
		_, network, err = net.ParseCIDR(*o.Network)
		if err != nil {
			return err
		}
	}

	if o.NextHop != nil {
		nextHop = net.ParseIP(*o.NextHop)
	}

	t.Network = network
	t.NextHop = nextHop
	t.SecurityZone = o.SecurityZone

	return nil
}

type rawConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt struct {
	Asn             *uint32 `json:"asn"`
	Bfd             bool    `json:"bfd"`
	Holdtime        *uint16 `json:"holdtime_timer"`
	Ipv4Addr        *string `json:"ipv4_addr"`
	Ipv4Safi        bool    `json:"ipv4_safi"`
	Ipv6Addr        *string `json:"ipv6_addr"`
	Ipv6Safi        bool    `json:"ipv6_safi"`
	Keepalive       *uint16 `json:"keepalive_timer"`
	LocalAsn        *uint32 `json:"local_asn"`
	NeighborAsnType string  `json:"neighbor_asn_type"`
	Password        *string `json:"password"`
	Ttl             uint8   `json:"ttl"`
}

func (o rawConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt) polish(t *ConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt) error {
	var neighborAsnDynamic bool
	switch o.NeighborAsnType {
	case "static":
		neighborAsnDynamic = false
	case "dynamic":
		neighborAsnDynamic = true
	default:
		return fmt.Errorf("unhandled neighbor asn type %q", o.NeighborAsnType)
	}

	var ipv4Addr, ipv6Addr string
	if o.Ipv4Addr != nil {
		ipv4Addr = *o.Ipv4Addr
	}
	if o.Ipv6Addr != nil {
		ipv6Addr = *o.Ipv6Addr
	}

	t.Asn = o.Asn
	t.Bfd = o.Bfd
	t.Holdtime = o.Holdtime
	t.Ipv4Addr = net.ParseIP(ipv4Addr)
	t.Ipv4Safi = o.Ipv4Safi
	t.Ipv6Addr = net.ParseIP(ipv6Addr)
	t.Ipv6Safi = o.Ipv6Safi
	t.Keepalive = o.Keepalive
	t.LocalAsn = o.LocalAsn
	t.NeighborAsnDynamic = neighborAsnDynamic
	t.Password = o.Password
	t.Ttl = o.Ttl

	return nil
}

type rawConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi struct {
	Bfd                   bool                                     `json:"bfd"`
	Holdtime              *uint16                                  `json:"holdtime_timer"`
	Ipv4Safi              bool                                     `json:"ipv4_safi"`
	Ipv6Safi              bool                                     `json:"ipv6_safi"`
	Keepalive             *uint16                                  `json:"keepalive_timer"`
	LocalAsn              *uint32                                  `json:"local_asn"`
	NeighborAsnType       string                                   `json:"neighbor_asn_type"` // static / dynamic
	Password              *string                                  `json:"password"`
	PeerFrom              string                                   `json:"peer_from"`
	PeerTo                ctPrimitiveBgpPeerTo                     `json:"peer_to"`
	SessionAddressingIpv4 ctPrimitiveIPv4ProtocolSessionAddressing `json:"session_addressing_ipv4"`
	SessionAddressingIpv6 ctPrimitiveIPv6ProtocolSessionAddressing `json:"session_addressing_ipv6"`
	Ttl                   uint8                                    `json:"ttl"`
}

func (o rawConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi) polish(t *ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi) error {
	var neighborAsnDynamic bool
	switch o.NeighborAsnType {
	case "dynamic":
		neighborAsnDynamic = true
	case "static":
		neighborAsnDynamic = false
	default:
		return fmt.Errorf("unhandled neighbor ASN type %q", o.NeighborAsnType)
	}

	var peerFromLoopback bool
	switch o.PeerFrom {
	case "loopback":
		peerFromLoopback = true
	case "interface":
		peerFromLoopback = false
	default:
		return fmt.Errorf("unhandled peer from value %q", o.PeerFrom)
	}

	peerTo, err := o.PeerTo.parse()
	if err != nil {
		return err
	}

	var ipv4Addressing CtPrimitiveIPv4ProtocolSessionAddressing
	err = ipv4Addressing.FromString(string(o.SessionAddressingIpv4))
	if err != nil {
		return err
	}

	var ipv6Addressing CtPrimitiveIPv6ProtocolSessionAddressing
	err = ipv6Addressing.FromString(string(o.SessionAddressingIpv6))
	if err != nil {
		return err
	}

	t.Bfd = o.Bfd
	t.Holdtime = o.Holdtime
	t.Ipv4Safi = o.Ipv4Safi
	t.Ipv6Safi = o.Ipv6Safi
	t.Keepalive = o.Keepalive
	t.LocalAsn = o.LocalAsn
	t.NeighborAsnDynamic = neighborAsnDynamic
	t.Password = o.Password
	t.PeerFromLoopback = peerFromLoopback
	t.PeerTo = CtPrimitiveBgpPeerTo(peerTo)
	t.SessionAddressingIpv4 = ipv4Addressing
	t.SessionAddressingIpv6 = ipv6Addressing
	t.Ttl = o.Ttl

	return nil
}

type rawConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface struct {
	Bfd                   bool    `json:"bfd"`
	Holdtime              *uint16 `json:"holdtime_timer"`
	Ipv4Safi              bool    `json:"ipv4_safi"`
	Ipv6Safi              bool    `json:"ipv6_safi"`
	Keepalive             *uint16 `json:"keepalive_timer"`
	LocalAsn              *uint32 `json:"local_asn"`
	Password              *string `json:"password"`
	PrefixNeighborIpv4    *string `json:"prefix_neighbor_ipv4"`
	PrefixNeighborIpv6    *string `json:"prefix_neighbor_ipv6"`
	SessionAddressingIpv4 bool    `json:"session_addressing_ipv4"`
	SessionAddressingIpv6 bool    `json:"session_addressing_ipv6"`
	Ttl                   int     `json:"ttl"`
}

func (o rawConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface) polish(t *ConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface) error {
	var prefixNeighborIpv4, prefixNeighborIpv6 *net.IPNet
	var err error

	if o.PrefixNeighborIpv4 != nil {
		_, prefixNeighborIpv4, err = net.ParseCIDR(*o.PrefixNeighborIpv4)
		if err != nil {
			return err
		}
	}

	if o.PrefixNeighborIpv6 != nil {
		_, prefixNeighborIpv6, err = net.ParseCIDR(*o.PrefixNeighborIpv6)
		if err != nil {
			return err
		}
	}

	t.Bfd = o.Bfd
	t.Holdtime = o.Holdtime
	t.Ipv4Safi = o.Ipv4Safi
	t.Ipv6Safi = o.Ipv6Safi
	t.Keepalive = o.Keepalive
	t.LocalAsn = o.LocalAsn
	t.Password = o.Password
	t.PrefixNeighborIpv4 = prefixNeighborIpv4
	t.PrefixNeighborIpv6 = prefixNeighborIpv6
	t.SessionAddressingIpv4 = o.SessionAddressingIpv4
	t.SessionAddressingIpv6 = o.SessionAddressingIpv6
	t.Ttl = o.Ttl

	return nil
}

type rawConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy struct {
	RpToAttach *string `json:"rp_to_attach"`
}

func (o rawConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy) polish(t *ConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy) error {
	t.RpToAttach = o.RpToAttach

	return nil
}

type rawConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint struct {
	RoutingZoneConstraint *string `json:"routing_zone_constraint"`
}

func (o rawConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint) polish(t *ConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint) error {
	t.RoutingZoneConstraint = o.RoutingZoneConstraint

	return nil
}
