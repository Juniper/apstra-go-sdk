package apstra

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

const (
	// Do not use apiUrlVirtualNetworks directly. The rawVirtualNetwork objects
	// in the returned map do not match the objects when retrieved using
	// apiUrlVirtualNetworkById
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

func (o dhcpServiceMode) polish() DhcpServiceEnabled {
	return o == dhcpServiceEnabled
}

//const (
//	l3ConnectivityEnabled  = l3ConnectivityMode("l3Enabled")
//	l3ConnectivityDisabled = l3ConnectivityMode("l3Disabled")
//)
//
//type L3ConnectivityMode bool
//type l3ConnectivityMode string
//
//func (o L3ConnectivityMode) raw() l3ConnectivityMode {
//	if o {
//		return l3ConnectivityEnabled
//	}
//	return l3ConnectivityDisabled
//}
//
//func (o l3ConnectivityMode) polish() L3ConnectivityMode {
//	return o == l3ConnectivityEnabled
//}

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

func (o sviIpRequirement) parse() (int, error) {
	switch o {
	case sviIpRequirementNone:
		return int(SviIpRequirementNone), nil
	case sviIpRequirementOptional:
		return int(SviIpRequirementOptional), nil
	case sviIpRequirementForbidden:
		return int(SviIpRequirementForbidden), nil
	case sviIpRequirementMandatory:
		return int(SviIpRequirementMandatory), nil
	case sviIpRequirementIntentionConflict:
		return int(SviIpRequirementIntentionConflict), nil
	default:
		return 0, fmt.Errorf(SviIpRequirementUnknown, o)
	}
}

type Ipv4Mode int
type ipv4Mode string

const (
	Ipv4ModeNone = Ipv4Mode(iota)
	Ipv4ModeDisabled
	Ipv4ModeEnabled
	Ipv4ModeForced
	Ipv4ModeUnknown = "unknown IPv4 mode '%s'"

	ipv4ModeNone     = ipv4Mode("")
	ipv4ModeDisabled = ipv4Mode("disabled")
	ipv4ModeEnabled  = ipv4Mode("enabled")
	ipv4ModeForced   = ipv4Mode("forced")
	ipv4ModeUnknown  = "unknown IPv4 mode %d"
)

func (o Ipv4Mode) String() string {
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

func (o ipv4Mode) parse() (int, error) {
	switch o {
	case ipv4ModeDisabled:
		return int(Ipv4ModeDisabled), nil
	case ipv4ModeEnabled:
		return int(Ipv4ModeEnabled), nil
	case ipv4ModeForced:
		return int(Ipv4ModeForced), nil
	default:
		return 0, fmt.Errorf(Ipv4ModeUnknown, o)
	}
}

type Ipv6Mode int
type ipv6Mode string

const (
	Ipv6ModeNone = Ipv6Mode(iota)
	Ipv6ModeDisabled
	Ipv6ModeEnabled
	Ipv6ModeForced
	Ipv6ModeLinkLocal
	Ipv6ModeUnknown = "unknown IPv6 mode '%s'"

	ipv6ModeNone      = ipv6Mode("")
	ipv6ModeDisabled  = ipv6Mode("disabled")
	ipv6ModeEnabled   = ipv6Mode("enabled")
	ipv6ModeForced    = ipv6Mode("forced")
	ipv6ModeLinkLocal = ipv6Mode("link_local")
	ipv6ModeUnknown   = "unknown IPv6 mode %d"
)

func (o Ipv6Mode) String() string {
	return string(o.raw())
}

func (o Ipv6Mode) int() int {
	return int(o)
}

func (o Ipv6Mode) raw() ipv6Mode {
	switch o {
	case Ipv6ModeNone:
		return ipv6ModeNone
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

func (o ipv6Mode) parse() (int, error) {
	switch o {
	case ipv6ModeNone:
		return int(Ipv6ModeNone), nil
	case ipv6ModeDisabled:
		return int(Ipv6ModeDisabled), nil
	case ipv6ModeEnabled:
		return int(Ipv6ModeEnabled), nil
	case ipv6ModeLinkLocal:
		return int(Ipv6ModeLinkLocal), nil
	case ipv6ModeForced:
		return int(Ipv6ModeForced), nil
	default:
		return 0, fmt.Errorf(Ipv6ModeUnknown, o)
	}
}

type VnType int
type vnType string

const (
	VnTypeNone = VnType(iota)
	VnTypeVlan
	VnTypeVxlan
	VnTypeOverlay
	VnTypeUnknown = "unknown VN type '%s'"

	vnTypeNone    = vnType("")
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
	case VnTypeNone:
		return vnTypeNone
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
func (o vnType) parse() (int, error) {
	switch o {
	case vnTypeNone:
		return int(VnTypeNone), nil
	case vnTypeOverlay:
		return int(VnTypeOverlay), nil
	case vnTypeVlan:
		return int(VnTypeVlan), nil
	case vnTypeVxlan:
		return int(VnTypeVxlan), nil
	default:
		return 0, fmt.Errorf(VnTypeUnknown, o)
	}
}

type SystemRole int
type systemRole string

const (
	SystemRoleNone = SystemRole(iota)
	SystemRoleAccess
	SystemRoleLeaf
	SystemRoleUnknown = "unknown System Role '%s'"

	systemRoleNone    = systemRole("")
	systemRoleAccess  = systemRole("access")
	systemRoleLeaf    = systemRole("leaf")
	systemRoleUnknown = "unknown System Role '%d'"
)

func (o SystemRole) String() string {
	return string(o.raw())
}
func (o SystemRole) int() int {
	return int(o)
}
func (o SystemRole) raw() systemRole {
	switch o {
	case SystemRoleNone:
		return systemRoleNone
	case SystemRoleAccess:
		return systemRoleAccess
	case SystemRoleLeaf:
		return systemRoleLeaf
	default:
		return systemRole(fmt.Sprintf(systemRoleUnknown, o))
	}
}
func (o systemRole) string() string {
	return string(o)
}
func (o systemRole) parse() (int, error) {
	switch o {
	case systemRoleNone:
		return int(SystemRoleNone), nil
	case systemRoleAccess:
		return int(SystemRoleAccess), nil
	case systemRoleLeaf:
		return int(SystemRoleLeaf), nil
	default:
		return 0, fmt.Errorf(SystemRoleUnknown, o)
	}
}

type SviIp struct {
	SystemId        ObjectId         `json:"system_id"`
	Ipv4Addr        net.IP           `json:"ipv4_addr"`
	Ipv4Mode        Ipv4Mode         `json:"ipv4_mode"`
	Ipv4Requirement SviIpRequirement `json:"ipv4_requirement"`
	Ipv6Addr        net.IP           `json:"ipv6_addr"`
	Ipv6Mode        Ipv6Mode         `json:"ipv6_mode"`
	Ipv6Requirement SviIpRequirement `json:"ipv6_requirement"`
}

func (o *SviIp) raw() *rawSviIp {
	var ipv4Addr, ipv6Addr string
	if len(o.Ipv4Addr) != 0 {
		ipv4Addr = o.Ipv4Addr.String()
	}
	if len(o.Ipv6Addr) != 0 {
		ipv6Addr = o.Ipv6Addr.String()
	}
	return &rawSviIp{
		SystemId:        o.SystemId,
		Ipv4Addr:        ipv4Addr,
		Ipv4Mode:        o.Ipv4Mode.raw(),
		Ipv4Requirement: o.Ipv4Requirement.raw(),
		Ipv6Addr:        ipv6Addr,
		Ipv6Mode:        o.Ipv6Mode.raw(),
		Ipv6Requirement: o.Ipv6Requirement.raw(),
	}
}

type rawSviIp struct {
	Ipv4Addr        string           `json:"ipv4_addr,omitempty"`
	Ipv4Mode        ipv4Mode         `json:"ipv4_mode,omitempty"`
	Ipv4Requirement sviIpRequirement `json:"ipv4_requirement,omitempty"` // not present in swagger example, not present in GET
	Ipv6Addr        string           `json:"ipv6_addr,omitempty"`
	Ipv6Mode        ipv6Mode         `json:"ipv6_mode,omitempty"`
	Ipv6Requirement sviIpRequirement `json:"ipv6_requirement,omitempty"` // not present in swagger example, not present in GET
	SystemId        ObjectId         `json:"system_id"`
}

func (o *rawSviIp) parse() (*SviIp, error) {
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

	return &SviIp{
		SystemId:        o.SystemId,
		Ipv4Addr:        ipv4Addr,
		Ipv4Mode:        Ipv4Mode(ipv4mode),
		Ipv4Requirement: SviIpRequirement(ipv4Requirement),
		Ipv6Addr:        ipv6Addr,
		Ipv6Mode:        Ipv6Mode(ipv6mode),
		Ipv6Requirement: SviIpRequirement(ipv6Requirement),
	}, nil
}

type VnBinding struct {
	//AccessSwitches []interface `json:"access_switches"`
	AccessSwitchNodeIds []ObjectId `json:"access_switch_node_ids"`
	//Role                string     `json:"role"`      // so far: "leaf", possibly graphdb "role" element
	SystemId ObjectId `json:"system_id"` // graphdb node id of a leaf (so far) switch
	VlanId   *Vlan    `json:"vlan_id"`   // optional (auto-assign)
	//Selected            bool       `json:"selected?"`
	// Tags []interface `json:"tags"` //sent as empty string by 4.1.2 web UI, not seen in 4.1.0 or 4.1.1

	//PodData struct {
	//	Description     interface{} `json:"description"`
	//	GlobalCatalogId interface{} `json:"global_catalog_id"`
	//	Label           string      `json:"label"`
	//	Position        int         `json:"position"`
	//	Type            string      `json:"type"`
	//	Id              string      `json:"id"`
	//} `json:"pod-data"`
}

type VirtualNetwork struct {
	Id   ObjectId
	Data *VirtualNetworkData
}

type Endpoint struct {
	InterfaceId ObjectId `json:"interface_id"`
	TagType     string   `json:"tag_type"`
	Label       string   `json:"label"`
}

type VirtualNetworkData struct {
	DhcpService               DhcpServiceEnabled
	Ipv4Enabled               bool
	Ipv4Subnet                *net.IPNet
	Ipv6Enabled               bool
	Ipv6Subnet                *net.IPNet
	Label                     string
	ReservedVlanId            *Vlan
	RouteTarget               string
	RtPolicy                  *RtPolicy
	SecurityZoneId            ObjectId
	SviIps                    []SviIp
	VirtualGatewayIpv4        net.IP
	VirtualGatewayIpv6        net.IP
	VirtualGatewayIpv4Enabled bool
	VirtualGatewayIpv6Enabled bool
	VnBindings                []VnBinding
	VnId                      *VNI
	VnType                    VnType
	VirtualMac                net.HardwareAddr
}

func (o *VirtualNetworkData) raw() *rawVirtualNetwork {
	var ipv4Subnet, ipv6Subnet string
	if o.Ipv4Subnet != nil {
		ipv4Subnet = o.Ipv4Subnet.String()
	}
	if o.Ipv6Subnet != nil {
		ipv6Subnet = o.Ipv6Subnet.String()
	}

	sviIps := make([]rawSviIp, len(o.SviIps))
	for i := range o.SviIps {
		sviIps[i] = *o.SviIps[i].raw()
	}

	var virtualGatewayIpv4, virtualGatewayIpv6 string
	if len(o.VirtualGatewayIpv4) == 4 {
		virtualGatewayIpv4 = o.VirtualGatewayIpv4.String()
	}
	if len(o.VirtualGatewayIpv6) == 16 {
		virtualGatewayIpv6 = o.VirtualGatewayIpv6.String()
	}

	var vnId string
	if o.VnId != nil {
		vnId = strconv.Itoa(int(*o.VnId))
	}

	return &rawVirtualNetwork{
		DhcpService:               o.DhcpService.raw(),
		Ipv4Enabled:               o.Ipv4Enabled,
		Ipv4Subnet:                ipv4Subnet,
		Ipv6Enabled:               o.Ipv6Enabled,
		Ipv6Subnet:                ipv6Subnet,
		Label:                     o.Label,
		ReservedVlanId:            o.ReservedVlanId,
		RouteTarget:               o.RouteTarget,
		RtPolicy:                  o.RtPolicy,
		SecurityZoneId:            o.SecurityZoneId,
		SviIps:                    sviIps,
		VirtualGatewayIpv4:        virtualGatewayIpv4,
		VirtualGatewayIpv6:        virtualGatewayIpv6,
		VirtualGatewayIpv4Enabled: o.VirtualGatewayIpv4Enabled,
		VirtualGatewayIpv6Enabled: o.VirtualGatewayIpv6Enabled,
		VnBindings:                o.VnBindings,
		VnId:                      vnId,
		VnType:                    o.VnType.raw(),
		VirtualMac:                o.VirtualMac.String(),
	}
}

type rawVirtualNetwork struct {
	Id                        ObjectId        `json:"id,omitempty"`
	DhcpService               dhcpServiceMode `json:"dhcp_service"`
	Ipv4Enabled               bool            `json:"ipv4_enabled"`
	Ipv4Subnet                string          `json:"ipv4_subnet,omitempty"`
	Ipv6Enabled               bool            `json:"ipv6_enabled"`
	Ipv6Subnet                string          `json:"ipv6_subnet,omitempty"`
	Label                     string          `json:"label"`
	ReservedVlanId            *Vlan           `json:"reserved_vlan_id,omitempty"`
	RouteTarget               string          `json:"route_target,omitempty"` // not mentioned in swagger, seen in 4.1.1: "10000:1"
	RtPolicy                  *RtPolicy       `json:"rt_policy"`
	SecurityZoneId            ObjectId        `json:"security_zone_id,omitempty"`
	SviIps                    []rawSviIp      `json:"svi_ips"`
	VirtualGatewayIpv4        string          `json:"virtual_gateway_ipv4,omitempty"`
	VirtualGatewayIpv6        string          `json:"virtual_gateway_ipv6,omitempty"`
	VirtualGatewayIpv4Enabled bool            `json:"virtual_gateway_ipv4_enabled"`
	VirtualGatewayIpv6Enabled bool            `json:"virtual_gateway_ipv6_enabled"`
	VnBindings                []VnBinding     `json:"bound_to"`
	VnId                      string          `json:"vn_id,omitempty"` // VNI as a string, null when unset
	VnType                    vnType          `json:"vn_type"`
	VirtualMac                string          `json:"virtual_mac,omitempty"`
	//CreatePolicyTagged      bool            `json:"create_policy_tagged"`
	//CreatePolicyUntagged    bool            `json:"create_policy_untagged"`
	//DefaultEndpointTagTypes interface{}     `json:"default_endpoint_tag_types"`    // what is this? not present in 4.1.1 api response
	//Description             string          `json:"description"`                   // not used in the web UI
	//FloatingIps             []interface{}   `json:"floating_ips"`                  // seen in 4.1.1 api response
	//ForceMoveUntaggedEndpoints bool         `json:"force_move_untagged_endpoints"` // not used in post/get with web UI
	//L3Connectivity          *l3ConnectivityMode `json:"l3_connectivity,omitempty"` // does not appear in 4.1.2 swagger
	//VniIds                  []interface{}   `json:"vni_ids,omitempty"`             // unknown, sent by web UI as empty list
	//Endpoints               []interface{}   `json:"endpoints"`                     // unknown, maybe relates to servers, etc?
}

func (o rawVirtualNetwork) polish() (*VirtualNetwork, error) {
	var err error

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

	sviIps := make([]SviIp, len(o.SviIps))
	for i, sviIp := range o.SviIps {
		SviIp, err := sviIp.parse()
		if err != nil {
			return nil, err
		}
		sviIps[i] = *SviIp
	}

	var vnId *VNI
	if o.VnId != "" {
		vniUint64, err := strconv.ParseUint(o.VnId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing VNID from string %q - %w", o.VnId, err)
		}
		vni := VNI(uint32(vniUint64))
		vnId = &vni
	}

	vntype, err := o.VnType.parse()
	if err != nil {
		return nil, err
	}

	var virtualMac net.HardwareAddr
	if o.VirtualMac != "" {
		virtualMac, err = net.ParseMAC(o.VirtualMac)
		if err != nil {
			return nil, fmt.Errorf("error parsing mac address %q - %w", o.VirtualMac, err)
		}
	}

	return &VirtualNetwork{
		Id: o.Id,
		Data: &VirtualNetworkData{
			DhcpService:               o.DhcpService.polish(),
			Ipv4Enabled:               o.Ipv4Enabled,
			Ipv4Subnet:                ipv4Subnet,
			Ipv6Enabled:               o.Ipv6Enabled,
			Ipv6Subnet:                ipv6Subnet,
			Label:                     o.Label,
			ReservedVlanId:            o.ReservedVlanId,
			RouteTarget:               o.RouteTarget,
			RtPolicy:                  o.RtPolicy,
			SecurityZoneId:            o.SecurityZoneId,
			SviIps:                    sviIps,
			VirtualGatewayIpv4:        net.ParseIP(o.VirtualGatewayIpv4),
			VirtualGatewayIpv6:        net.ParseIP(o.VirtualGatewayIpv6),
			VirtualGatewayIpv4Enabled: o.VirtualGatewayIpv4Enabled,
			VirtualGatewayIpv6Enabled: o.VirtualGatewayIpv6Enabled,
			VnBindings:                o.VnBindings,
			VnId:                      vnId,
			VnType:                    VnType(vntype),
			VirtualMac:                virtualMac,
		},
	}, nil
}

func (o *TwoStageL3ClosClient) listAllVirtualNetworkIds(ctx context.Context) ([]ObjectId, error) {
	apstraUrl, err := o.urlWithParam(fmt.Sprintf(apiUrlVirtualNetworks, o.blueprintId))
	if err != nil {
		return nil, err
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
	for id := range response.VirtualNetworks {
		result[i] = id
		i++
	}
	return result, nil
}

func (o *TwoStageL3ClosClient) getVirtualNetwork(ctx context.Context, vnId ObjectId) (*rawVirtualNetwork, error) {
	apstraUrl, err := o.urlWithParam(fmt.Sprintf(apiUrlVirtualNetworkById, o.blueprintId, vnId))
	if err != nil {
		return nil, err
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

	return response, nil
}

func (o *TwoStageL3ClosClient) createVirtualNetwork(ctx context.Context, cfg *rawVirtualNetwork) (ObjectId, error) {
	if cfg.Id != "" {
		return "", fmt.Errorf("refusing to create virtual network using input data with a populated ID field")
	}
	response := &objectIdResponse{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlVirtualNetworks, o.blueprintId),
		apiInput:    cfg,
		apiResponse: response,
	})

	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *TwoStageL3ClosClient) updateVirtualNetwork(ctx context.Context, id ObjectId, cfg *rawVirtualNetwork) error {
	if cfg.Id != "" {
		return fmt.Errorf("refusing to update virtual network using input data with a populated ID field")
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlVirtualNetworkById, o.blueprintId, id),
		apiInput: cfg,
	})

	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) deleteVirtualNetwork(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlVirtualNetworkById, o.blueprintId, id),
	})

	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
