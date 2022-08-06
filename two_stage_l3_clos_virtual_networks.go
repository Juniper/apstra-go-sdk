package goapstra

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

const (
	apiUrlVirtualNetworks    = apiUrlBlueprintById + apiUrlPathDelim + "virtual-networks"
	apiUrlVirtualNetworkById = apiUrlVirtualNetworks + apiUrlPathDelim + "%s"

	dhcpServiceDisabled = dhcpServiceMode("dhcpServiceDisabled")
	dhcpServiceEnabled  = dhcpServiceMode("dhcpServiceEnabled")
)

type DhcpServiceEnabled bool
type dhcpServiceMode string

func (o DhcpServiceEnabled) raw() dhcpServiceMode {
	if o {
		return dhcpServiceEnabled
	}
	return dhcpServiceDisabled
}

func (o dhcpServiceMode) parse() DhcpServiceEnabled {
	if o == dhcpServiceEnabled {
		return true
	}
	return false
}

const (
	l3ConnectivityEnabled  = l3ConnectivityMode("l3Enabled")
	l3ConnectivityDisabled = l3ConnectivityMode("l3Disabled")
)

type L3ConnectivityEnabled bool
type l3ConnectivityMode string

func (o L3ConnectivityEnabled) raw() l3ConnectivityMode {
	if o {
		return l3ConnectivityEnabled
	}
	return l3ConnectivityDisabled
}

func (o l3ConnectivityMode) parse() L3ConnectivityEnabled {
	if o == l3ConnectivityEnabled {
		return true
	}
	return false
}

type SviIpRequirement int
type sviIpRequirement string

const (
	SviIpRequirementNone = SviIpRequirement(iota)
	SviIpRequirementOptional
	SviIpRequirementForbidden
	SviIpRequirementMandatory
	SviIpRequirementIntentionConflict
	SviIpRequirementUnknown = "SVI IP requirement mode '%s' unknown"

	sviIpRequirementNone              = sviIpRequirement("")
	sviIpRequirementOptional          = sviIpRequirement("optional")
	sviIpRequirementForbidden         = sviIpRequirement("forbidden")
	sviIpRequirementMandatory         = sviIpRequirement("mandatory")
	sviIpRequirementIntentionConflict = sviIpRequirement("intention_conflict")
	sviIpRequirementUnknown           = "SVI IP requirement mode %d unknown"
)

func (o SviIpRequirement) String() string {
	return string(o.raw())
}
func (o SviIpRequirement) int() int {
	return int(o)
}
func (o SviIpRequirement) raw() sviIpRequirement {
	switch o {
	case SviIpRequirementNone:
		return sviIpRequirementNone
	case SviIpRequirementOptional:
		return sviIpRequirementOptional
	case SviIpRequirementForbidden:
		return sviIpRequirementForbidden
	case SviIpRequirementMandatory:
		return sviIpRequirementMandatory
	case SviIpRequirementIntentionConflict:
		return sviIpRequirementIntentionConflict
	default:
		return sviIpRequirement(fmt.Sprintf(sviIpRequirementUnknown, o))
	}
}

func (o sviIpRequirement) string() string {
	return string(o)
}

func (o sviIpRequirement) parse() (SviIpRequirement, error) {
	switch o {
	case sviIpRequirementNone:
		return SviIpRequirementNone, nil
	case sviIpRequirementOptional:
		return SviIpRequirementOptional, nil
	case sviIpRequirementForbidden:
		return SviIpRequirementForbidden, nil
	case sviIpRequirementMandatory:
		return SviIpRequirementMandatory, nil
	case sviIpRequirementIntentionConflict:
		return SviIpRequirementIntentionConflict, nil
	default:
		return 0, fmt.Errorf(SviIpRequirementUnknown, o)
	}
}

type Ipv4Mode int
type ipv4Mode string

const (
	Ipv4ModeDisabled = Ipv4Mode(iota)
	Ipv4ModeEnabled
	Ipv4ModeForced
	Ipv4ModeUnknown = "unknown IPv4 mode '%s'"

	ipv4ModeDisabled = ipv4Mode("disabled")
	ipv4ModeEnabled  = ipv4Mode("enabled")
	ipv4ModeForced   = ipv4Mode("forced")
	ipv4ModeUnknown  = "unknown IPv4 mode %d"
)

func (o Ipv4Mode) string() string {
	return string(o.raw())
}

func (o Ipv4Mode) int() int {
	return int(o)
}

func (o Ipv4Mode) raw() ipv4Mode {
	switch o {
	case Ipv4ModeDisabled:
		return ipv4ModeDisabled
	case Ipv4ModeEnabled:
		return ipv4ModeEnabled
	case Ipv4ModeForced:
		return ipv4ModeForced
	default:
		return ipv4Mode(fmt.Sprintf(ipv4ModeUnknown, o))
	}
}

func (o ipv4Mode) string() string {
	return string(o)
}

func (o ipv4Mode) parse() (Ipv4Mode, error) {
	switch o {
	case ipv4ModeDisabled:
		return Ipv4ModeDisabled, nil
	case ipv4ModeEnabled:
		return Ipv4ModeEnabled, nil
	case ipv4ModeForced:
		return Ipv4ModeForced, nil
	default:
		return 0, fmt.Errorf(Ipv4ModeUnknown, o)
	}
}

type Ipv6Mode int
type ipv6Mode string

const (
	Ipv6ModeDisabled = Ipv6Mode(iota)
	Ipv6ModeEnabled
	Ipv6ModeForced
	Ipv6ModeLinkLocal
	Ipv6ModeUnknown = "unknown IPv6 mode '%s'"

	ipv6ModeDisabled  = ipv6Mode("disabled")
	ipv6ModeEnabled   = ipv6Mode("enabled")
	ipv6ModeForced    = ipv6Mode("forced")
	ipv6ModeLinkLocal = ipv6Mode("link_local")
	ipv6ModeUnknown   = "unknown IPv6 mode %d"
)

func (o Ipv6Mode) string() string {
	return string(o.raw())
}

func (o Ipv6Mode) int() int {
	return int(o)
}

func (o Ipv6Mode) raw() ipv6Mode {
	switch o {
	case Ipv6ModeDisabled:
		return ipv6ModeDisabled
	case Ipv6ModeEnabled:
		return ipv6ModeEnabled
	case Ipv6ModeLinkLocal:
		return ipv6ModeLinkLocal
	case Ipv6ModeForced:
		return ipv6ModeForced
	default:
		return ipv6Mode(fmt.Sprintf(ipv6ModeUnknown, o))
	}
}

func (o ipv6Mode) string() string {
	return string(o)
}

func (o ipv6Mode) parse() (Ipv6Mode, error) {
	switch o {
	case ipv6ModeDisabled:
		return Ipv6ModeDisabled, nil
	case ipv6ModeEnabled:
		return Ipv6ModeEnabled, nil
	case ipv6ModeLinkLocal:
		return Ipv6ModeLinkLocal, nil
	case ipv6ModeForced:
		return Ipv6ModeForced, nil
	default:
		return 0, fmt.Errorf(Ipv6ModeUnknown, o)
	}
}

type VnType int
type vnType string

const (
	VnTypeVlan = VnType(iota)
	VnTypeVxlan
	VnTypeOverlay
	VnTypeUnknown = "unknown VN type '%s'"

	vnTypeVlan    = vnType("vlan")
	vnTypeVxlan   = vnType("vxlan")
	vnTypeOverlay = vnType("overlay")
	vnTypeUnknown = "unknown VN type '%d'"
)

func (o VnType) String() string {
	return string(o.raw())
}
func (o VnType) int() int {
	return int(o)
}
func (o VnType) raw() vnType {
	switch o {
	case VnTypeOverlay:
		return vnTypeOverlay
	case VnTypeVlan:
		return vnTypeVlan
	case VnTypeVxlan:
		return vnTypeVxlan
	default:
		return vnType(fmt.Sprintf(vnTypeUnknown, o))
	}
}
func (o vnType) string() string {
	return string(o)
}
func (o vnType) parse() (VnType, error) {
	switch o {
	case vnTypeOverlay:
		return VnTypeOverlay, nil
	case vnTypeVlan:
		return VnTypeVlan, nil
	case vnTypeVxlan:
		return VnTypeVxlan, nil
	default:
		return 0, fmt.Errorf(VnTypeUnknown, o)
	}
}

type SviIps struct {
	SystemId        ObjectId         `json:"system_id"`
	Ipv4Addr        net.IP           `json:"ipv4_addr"`
	Ipv4Mode        Ipv4Mode         `json:"ipv4_mode"`
	Ipv4Requirement SviIpRequirement `json:"ipv4_requirement"`
	Ipv6Addr        net.IP           `json:"ipv6_addr"`
	Ipv6Mode        Ipv6Mode         `json:"ipv6_mode"`
	Ipv6Requirement SviIpRequirement `json:"ipv6_requirement"`
}

func (o *SviIps) raw() *rawSviIps {
	return &rawSviIps{
		SystemId:        o.SystemId,
		Ipv4Addr:        o.Ipv4Addr.String(),
		Ipv4Mode:        o.Ipv4Mode.raw(),
		Ipv4Requirement: o.Ipv4Requirement.raw(),
		Ipv6Addr:        o.Ipv6Addr.String(),
		Ipv6Mode:        o.Ipv6Mode.raw(),
		Ipv6Requirement: o.Ipv6Requirement.raw(),
	}
}

type rawSviIps struct {
	SystemId        ObjectId         `json:"system_id"`
	Ipv4Addr        string           `json:"ipv4_addr"`
	Ipv4Mode        ipv4Mode         `json:"ipv4_mode"`
	Ipv4Requirement sviIpRequirement `json:"ipv4_requirement"`
	Ipv6Addr        string           `json:"ipv6_addr"`
	Ipv6Mode        ipv6Mode         `json:"ipv6_mode"`
	Ipv6Requirement sviIpRequirement `json:"ipv6_requirement"`
}

func (o *rawSviIps) parse() (*SviIps, error) {
	var ipv4Addr, ipv6Addr net.IP
	var err error

	if o.Ipv4Addr != "" {
		ipv4Addr, _, err = net.ParseCIDR(o.Ipv4Addr)
		if err != nil {
			return nil, err
		}
	}

	if o.Ipv6Addr != "" {
		ipv6Addr, _, err = net.ParseCIDR(o.Ipv6Addr)
		if err != nil {
			return nil, err
		}
	}

	ipv4mode, err := o.Ipv4Mode.parse()
	if err != nil {
		return nil, err
	}
	ipv6mode, err := o.Ipv6Mode.parse()
	if err != nil {
		return nil, err
	}

	ipv4Requirement, err := o.Ipv4Requirement.parse()
	if err != nil {
		return nil, err
	}
	ipv6Requirement, err := o.Ipv6Requirement.parse()
	if err != nil {
		return nil, err
	}

	return &SviIps{
		SystemId:        o.SystemId,
		Ipv4Addr:        ipv4Addr,
		Ipv4Mode:        ipv4mode,
		Ipv4Requirement: ipv4Requirement,
		Ipv6Addr:        ipv6Addr,
		Ipv6Mode:        ipv6mode,
		Ipv6Requirement: ipv6Requirement,
	}, nil
}

type VNBoundTo struct {
	AccessSwitchNodeIds []ObjectId `json:"access_switch_node_ids"`
	SystemId            ObjectId   `json:"system_id"`
	VlanId              uint16     `json:"vlan_id"`
}

type VirtualNetwork struct {
	Id                        ObjectId              `json:"id"`
	Label                     string                `json:"label"`
	VnType                    VnType                `json:"vn_type"`
	BoundTo                   []VNBoundTo           `json:"bound_to"`
	SviIps                    []SviIps              `json:"svi_ips"`
	DhcpService               DhcpServiceEnabled    `json:"dhcp_service"`
	L3Connectivity            L3ConnectivityEnabled `json:"l3_connectivity"`
	Ipv4Enabled               bool                  `json:"ipv4_enabled"`
	Ipv6Enabled               bool                  `json:"ipv6_enabled"`
	Ipv4Subnet                *net.IPNet            `json:"ipv4_subnet"`
	Ipv6Subnet                *net.IPNet            `json:"ipv6_subnet"`
	VirtualGatewayIpv4Enabled bool                  `json:"virtual_gateway_ipv4_enabled"`
	VirtualGatewayIpv6Enabled bool                  `json:"virtual_gateway_ipv6_enabled"`
	SecurityZoneId            ObjectId              `json:"security_zone_id"`
}

type rawVirtualNetwork struct {
	Id                        ObjectId           `json:"id"`
	Label                     string             `json:"label"`
	VnType                    vnType             `json:"vn_type"`
	BoundTo                   []VNBoundTo        `json:"bound_to"`
	SviIps                    []rawSviIps        `json:"svi_ips"`
	DhcpService               dhcpServiceMode    `json:"dhcp_service"`
	L3Connectivity            l3ConnectivityMode `json:"l3_connectivity"`
	Ipv4Enabled               bool               `json:"ipv4_enabled"`
	Ipv6Enabled               bool               `json:"ipv6_enabled"`
	Ipv4Subnet                string             `json:"ipv4_subnet"`
	Ipv6Subnet                string             `json:"ipv6_subnet"`
	VirtualGatewayIpv4Enabled bool               `json:"virtual_gateway_ipv4_enabled"`
	VirtualGatewayIpv6Enabled bool               `json:"virtual_gateway_ipv6_enabled"`
	SecurityZoneId            ObjectId           `json:"security_zone_id"`
	//VnId                      string           `json:"vn_id"`                // todo
	//RtPolicy                  interface{}      `json:"rt_policy"`            // todo
	//VirtualGatewayIpv4        interface{}      `json:"virtual_gateway_ipv4"` // todo
	//VirtualGatewayIpv6        interface{}      `json:"virtual_gateway_ipv6"` // todo
	//FloatingIps               []interface{}    `json:"floating_ips"`         // todo
	//ReservedVlanId            interface{}      `json:"reserved_vlan_id"`     // todo
	//Description               interface{}      `json:"description"`          // todo
	//VirtualMac                interface{}      `json:"virtual_mac"`          // todo
	//RouteTarget               interface{}      `json:"route_target"`         // todo
	//Endpoints                 []interface{}    `json:"endpoints"`            // todo
}

func (o rawVirtualNetwork) parse() (*VirtualNetwork, error) {
	vntype, err := o.VnType.parse()
	if err != nil {
		return nil, err
	}

	sviips := make([]SviIps, len(o.SviIps))
	for i, sviIp := range o.SviIps {
		SviIp, err := sviIp.parse()
		if err != nil {
			return nil, err
		}
		sviips[i] = *SviIp
	}

	var ipv4Subnet *net.IPNet
	if o.Ipv4Subnet != "" {
		_, ipv4Subnet, err = net.ParseCIDR(o.Ipv4Subnet)
		if err != nil {
			return nil, err
		}
	}

	var ipv6Subnet *net.IPNet
	if o.Ipv6Subnet != "" {
		_, ipv6Subnet, err = net.ParseCIDR(o.Ipv6Subnet)
		if err != nil {
			return nil, err
		}
	}

	return &VirtualNetwork{
		Id:                        o.Id,
		Label:                     o.Label,
		VnType:                    vntype,
		BoundTo:                   o.BoundTo,
		SviIps:                    sviips,
		DhcpService:               o.DhcpService.parse(),
		L3Connectivity:            o.L3Connectivity.parse(),
		Ipv4Enabled:               o.Ipv4Enabled,
		Ipv6Enabled:               o.Ipv6Enabled,
		Ipv4Subnet:                ipv4Subnet,
		Ipv6Subnet:                ipv6Subnet,
		VirtualGatewayIpv4Enabled: o.VirtualGatewayIpv4Enabled,
		VirtualGatewayIpv6Enabled: o.VirtualGatewayIpv6Enabled,
		SecurityZoneId:            o.SecurityZoneId,
	}, nil
}

func (o *TwoStageLThreeClosClient) listAllVirtualNetworkIds(ctx context.Context, bpType BlueprintType) ([]ObjectId, error) {
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlVirtualNetworks, o.blueprintId))
	if err != nil {
		return nil, err
	}

	if bpType != BlueprintTypeNone {
		params := apstraUrl.Query()
		params.Set(blueprintTypeParam, bpType.string())
		apstraUrl.RawQuery = params.Encode()
	}

	response := &struct {
		VirtualNetworks map[ObjectId]rawVirtualNetwork `json:"virtual_networks"`
	}{}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodGet,
		url:            apstraUrl,
		apiResponse:    response,
		unsynchronized: true,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	result := make([]ObjectId, len(response.VirtualNetworks))
	i := 0
	for id, _ := range response.VirtualNetworks {
		result[i] = id
		i++
	}
	return result, nil
}

func (o *TwoStageLThreeClosClient) getVirtualNetwork(ctx context.Context, vnId ObjectId, bpType BlueprintType) (*VirtualNetwork, error) {
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlVirtualNetworkById, o.blueprintId, vnId))
	if err != nil {
		return nil, err
	}

	if bpType != BlueprintTypeNone {
		params := apstraUrl.Query()
		params.Set(blueprintTypeParam, bpType.string())
		apstraUrl.RawQuery = params.Encode()
	}

	response := &rawVirtualNetwork{}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodGet,
		url:            apstraUrl,
		apiResponse:    response,
		unsynchronized: true,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.parse()
}

func (o *TwoStageLThreeClosClient) getVirtualNetworkBySubnet(ctx context.Context, desiredNet *net.IPNet, vrf ObjectId, bpType BlueprintType) (*VirtualNetwork, error) {
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlVirtualNetworks, o.blueprintId))
	if err != nil {
		return nil, err
	}

	if bpType != BlueprintTypeNone {
		params := apstraUrl.Query()
		params.Set(blueprintTypeParam, bpType.string())
		apstraUrl.RawQuery = params.Encode()
	}

	response := &struct {
		VirtualNetworks map[ObjectId]rawVirtualNetwork `json:"virtual_networks"`
	}{}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodGet,
		url:            apstraUrl,
		apiResponse:    response,
		unsynchronized: true,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	target := desiredNet.String()
	for _, rawVn := range response.VirtualNetworks {
		if rawVn.SecurityZoneId == vrf {
			_, ipv4net, err := net.ParseCIDR(rawVn.Ipv4Subnet)
			if err != nil {
				return nil, err
			}
			if ipv4net.String() == target {
				return o.getVirtualNetwork(ctx, rawVn.Id, bpType)
			}

			_, ipv6net, err := net.ParseCIDR(rawVn.Ipv4Subnet)
			if err != nil {
				return nil, err
			}
			if ipv6net.String() == target {
				return o.getVirtualNetwork(ctx, rawVn.Id, bpType)
			}
		}

	}

	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("virtual network for subnet '%s' in vrf '%s' not found", desiredNet.String(), vrf),
	}
}
