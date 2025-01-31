// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	// Do not use apiUrlVirtualNetworks directly. The rawVirtualNetwork objects
	// in the returned map do not match the objects when retrieved using
	// apiUrlVirtualNetworkById
	apiUrlVirtualNetworks    = apiUrlBlueprintById + apiUrlPathDelim + "virtual-networks"
	apiUrlVirtualNetworkById = apiUrlVirtualNetworks + apiUrlPathDelim + "%s"
)

var _ json.Marshaler = (*DhcpServiceEnabled)(nil)

type DhcpServiceEnabled bool

func (o *DhcpServiceEnabled) FromString(s string) error {
	var dsm enum.DhcpServiceMode

	err := dsm.FromString(s)
	if err != nil {
		return fmt.Errorf("while parsing dhcp service mode - %w", err)
	}

	*o = dsm == enum.DhcpServiceModeEnabled

	return nil
}

func (o DhcpServiceEnabled) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o DhcpServiceEnabled) String() string {
	if o {
		return enum.DhcpServiceModeEnabled.String()
	}

	return enum.DhcpServiceModeDisabled.String()
}

type (
	SystemRole int
	systemRole string
)

const (
	SystemRoleNone = SystemRole(iota)
	SystemRoleAccess
	SystemRoleGeneric
	SystemRoleLeaf
	SystemRoleSpine
	SystemRoleSuperSpine
	SystemRoleRedundancyGroup
	SystemRoleUnknown = "unknown System Role '%s'"

	systemRoleNone            = systemRole("")
	systemRoleAccess          = systemRole("access")
	systemRoleGeneric         = systemRole("generic")
	systemRoleLeaf            = systemRole("leaf")
	systemRoleSpine           = systemRole("spine")
	systemRoleSuperSpine      = systemRole("superspine")
	systemRoleRedundancyGroup = systemRole("redundancy_group")
	systemRoleUnknown         = "unknown System Role '%d'"
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
	case SystemRoleGeneric:
		return systemRoleGeneric
	case SystemRoleLeaf:
		return systemRoleLeaf
	case SystemRoleSpine:
		return systemRoleSpine
	case SystemRoleSuperSpine:
		return systemRoleSuperSpine
	case SystemRoleRedundancyGroup:
		return systemRoleRedundancyGroup
	default:
		return systemRole(fmt.Sprintf(systemRoleUnknown, o))
	}
}

func (o *SystemRole) FromString(in string) error {
	i, err := systemRole(in).parse()
	if err != nil {
		return err
	}
	*o = SystemRole(i)
	return nil
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
	case systemRoleGeneric:
		return int(SystemRoleGeneric), nil
	case systemRoleLeaf:
		return int(SystemRoleLeaf), nil
	case systemRoleSpine:
		return int(SystemRoleSpine), nil
	case systemRoleSuperSpine:
		return int(SystemRoleSuperSpine), nil
	case systemRoleRedundancyGroup:
		return int(SystemRoleRedundancyGroup), nil
	default:
		return 0, fmt.Errorf(SystemRoleUnknown, o)
	}
}

var (
	_ json.Marshaler   = (*SviIp)(nil)
	_ json.Unmarshaler = (*SviIp)(nil)
)

type SviIp struct {
	SystemId ObjectId
	Ipv4Addr *net.IPNet
	Ipv4Mode enum.SviIpv4Mode
	Ipv6Addr *net.IPNet
	Ipv6Mode enum.SviIpv6Mode
}

func (o *SviIp) MarshalJSON() ([]byte, error) {
	var raw struct {
		Ipv4Addr string   `json:"ipv4_addr,omitempty"`
		Ipv4Mode string   `json:"ipv4_mode,omitempty"`
		Ipv6Addr string   `json:"ipv6_addr,omitempty"`
		Ipv6Mode string   `json:"ipv6_mode,omitempty"`
		SystemId ObjectId `json:"system_id"`
	}

	if o.Ipv4Addr != nil {
		raw.Ipv4Addr = o.Ipv4Addr.String()
	}
	raw.Ipv4Mode = o.Ipv4Mode.String()

	if o.Ipv6Addr != nil {
		raw.Ipv6Addr = o.Ipv6Addr.String()
	}
	raw.Ipv6Mode = o.Ipv6Mode.String()

	raw.SystemId = o.SystemId

	return json.Marshal(raw)
}

func (o *SviIp) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Ipv4Addr string   `json:"ipv4_addr"`
		Ipv4Mode string   `json:"ipv4_mode"`
		Ipv6Addr string   `json:"ipv6_addr"`
		Ipv6Mode string   `json:"ipv6_mode"`
		SystemId ObjectId `json:"system_id"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("while unmarshaling an SviIp - %w", err)
	}

	o.Ipv4Addr = nil
	if raw.Ipv4Addr != "" {
		var ip net.IP
		ip, o.Ipv4Addr, err = net.ParseCIDR(raw.Ipv4Addr)
		if err != nil {
			return fmt.Errorf("while parsing SviIp.Ipv4Addr - %w", err)
		}
		o.Ipv4Addr.IP = ip
	}

	err = o.Ipv4Mode.FromString(raw.Ipv4Mode)
	if err != nil {
		return fmt.Errorf("while parsing SviIp.Ipv4Mode - %w", err)
	}

	o.Ipv6Addr = nil
	if raw.Ipv6Addr != "" {
		var ip net.IP
		ip, o.Ipv6Addr, err = net.ParseCIDR(raw.Ipv6Addr)
		if err != nil {
			return fmt.Errorf("while parsing SviIp.Ipv6Addr - %w", err)
		}
		o.Ipv6Addr.IP = ip
	}

	err = o.Ipv6Mode.FromString(raw.Ipv6Mode)
	if err != nil {
		return fmt.Errorf("while parsing SviIp.Ipv6Mode - %w", err)
	}

	o.SystemId = raw.SystemId

	return nil
}

type VnBinding struct {
	AccessSwitchNodeIds []ObjectId `json:"access_switch_node_ids"`
	SystemId            ObjectId   `json:"system_id"` // graphdb node id of a leaf (so far) switch
	VlanId              *Vlan      `json:"vlan_id"`   // optional (auto-assign)
}

var _ json.Unmarshaler = (*VirtualNetwork)(nil)

type VirtualNetwork struct {
	Id   ObjectId
	Data *VirtualNetworkData
}

func (o *VirtualNetwork) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id                        ObjectId    `json:"id"`
		Description               string      `json:"description"`
		DhcpService               string      `json:"dhcp_service"`
		Ipv4Enabled               bool        `json:"ipv4_enabled"`
		Ipv4Subnet                string      `json:"ipv4_subnet"`
		Ipv6Enabled               bool        `json:"ipv6_enabled"`
		Ipv6Subnet                string      `json:"ipv6_subnet"`
		L3Mtu                     *int        `json:"l3_mtu"`
		Label                     string      `json:"label"`
		ReservedVlanId            *Vlan       `json:"reserved_vlan_id"`
		RouteTarget               string      `json:"route_target"`
		RtPolicy                  *RtPolicy   `json:"rt_policy"`
		SecurityZoneId            ObjectId    `json:"security_zone_id"`
		SviIps                    []SviIp     `json:"svi_ips"`
		VirtualGatewayIpv4        string      `json:"virtual_gateway_ipv4"`
		VirtualGatewayIpv6        string      `json:"virtual_gateway_ipv6"`
		VirtualGatewayIpv4Enabled bool        `json:"virtual_gateway_ipv4_enabled"`
		VirtualGatewayIpv6Enabled bool        `json:"virtual_gateway_ipv6_enabled"`
		VnBindings                []VnBinding `json:"bound_to"`
		VnId                      string      `json:"vn_id"`
		VnType                    string      `json:"vn_type"`
		VirtualMac                string      `json:"virtual_mac"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("while unmarshaling raw API response - %w", err)
	}

	o.Id = raw.Id
	o.Data = new(VirtualNetworkData)

	o.Data.Description = raw.Description

	err = o.Data.DhcpService.FromString(raw.DhcpService)
	if err != nil {
		return fmt.Errorf("while unmarshaling dhcp_service %q - %w", raw.DhcpService, err)
	}

	o.Data.Ipv4Enabled = raw.Ipv4Enabled
	o.Data.Ipv4Subnet, err = ipNetFromString(raw.Ipv4Subnet)
	if err != nil {
		return fmt.Errorf("while parsing virtual network data ipv4_subnet %q - %w", raw.Ipv4Subnet, err)
	}

	o.Data.Ipv6Enabled = raw.Ipv6Enabled
	o.Data.Ipv6Subnet, err = ipNetFromString(raw.Ipv6Subnet)
	if err != nil {
		return fmt.Errorf("while parsing virtual network data ipv6_subnet %q - %w", raw.Ipv6Subnet, err)
	}

	o.Data.L3Mtu = raw.L3Mtu
	o.Data.Label = raw.Label
	o.Data.ReservedVlanId = raw.ReservedVlanId
	o.Data.RouteTarget = raw.RouteTarget
	o.Data.RtPolicy = raw.RtPolicy
	o.Data.SecurityZoneId = raw.SecurityZoneId
	o.Data.SviIps = raw.SviIps

	o.Data.VirtualGatewayIpv4, err = ipFromString(raw.VirtualGatewayIpv4)
	if err != nil {
		return fmt.Errorf("while parsing virtual network data virtual_gateway_ipv4 %q - %w", raw.VirtualGatewayIpv4, err)
	}

	o.Data.VirtualGatewayIpv6, err = ipFromString(raw.VirtualGatewayIpv6)
	if err != nil {
		return fmt.Errorf("while parsing virtual network data virtual_gateway_ipv6 %q - %w", raw.VirtualGatewayIpv6, err)
	}

	o.Data.VirtualGatewayIpv4Enabled = raw.VirtualGatewayIpv4Enabled
	o.Data.VirtualGatewayIpv6Enabled = raw.VirtualGatewayIpv6Enabled
	o.Data.VnBindings = raw.VnBindings

	if raw.VnId != "" {
		vnId, err := strconv.Atoi(raw.VnId)
		if err != nil {
			return fmt.Errorf("while parsing virtual network data vn_id %q - %w", raw.VnId, err)
		}
		o.Data.VnId = (*VNI)(toPtr(uint32(vnId)))
	}

	err = o.Data.VnType.FromString(raw.VnType)
	if err != nil {
		return fmt.Errorf("while parsing virtual network data vn_type %q - %w", raw.VnType, err)
	}

	o.Data.VirtualMac, err = macFromString(raw.VirtualMac)
	if err != nil {
		return fmt.Errorf("while parsing virtual network data virtual_mac %q - %w", raw.VirtualMac, err)
	}

	return nil
}

type Endpoint struct {
	InterfaceId ObjectId `json:"interface_id"`
	TagType     string   `json:"tag_type"`
	Label       string   `json:"label"`
}

var _ json.Marshaler = (*VirtualNetworkData)(nil)

type VirtualNetworkData struct {
	Description               string
	DhcpService               DhcpServiceEnabled `json:"dhcp_service"`
	Ipv4Enabled               bool               `json:"ipv4_enabled"`
	Ipv4Subnet                *net.IPNet         `json:"ipv4_subnet,omitempty"`
	Ipv6Enabled               bool               `json:"ipv6_enabled"`
	Ipv6Subnet                *net.IPNet         `json:"ipv6_subnet,omitempty"`
	L3Mtu                     *int               `json:"l3_mtu,omitempty"`
	Label                     string             `json:"label"`
	ReservedVlanId            *Vlan              `json:"reserved_vlan_id,omitempty"`
	RouteTarget               string             `json:"route_target,omitempty"`
	RtPolicy                  *RtPolicy          `json:"rt_policy"`
	SecurityZoneId            ObjectId           `json:"security_zone_id,omitempty"`
	SviIps                    []SviIp            `json:"svi_ips"`
	VirtualGatewayIpv4        net.IP             `json:"virtual_gateway_ipv4,omitempty"`
	VirtualGatewayIpv6        net.IP             `json:"virtual_gateway_ipv6,omitempty"`
	VirtualGatewayIpv4Enabled bool               `json:"virtual_gateway_ipv4_enabled"`
	VirtualGatewayIpv6Enabled bool               `json:"virtual_gateway_ipv6_enabled"`
	VnBindings                []VnBinding        `json:"bound_to"`
	VnId                      *VNI               `json:"vn_id,omitempty"` // VNI as a string, null when unset
	VnType                    enum.VnType        `json:"vn_type"`
	VirtualMac                net.HardwareAddr   `json:"virtual_mac,omitempty"`
}

func (o VirtualNetworkData) MarshalJSON() ([]byte, error) {
	raw := struct {
		Description               string      `json:"description,omitempty"` // 5.0 and later only
		DhcpService               string      `json:"dhcp_service"`
		Ipv4Enabled               bool        `json:"ipv4_enabled"`
		Ipv4Subnet                string      `json:"ipv4_subnet,omitempty"`
		Ipv6Enabled               bool        `json:"ipv6_enabled"`
		Ipv6Subnet                string      `json:"ipv6_subnet,omitempty"`
		L3Mtu                     *int        `json:"l3_mtu,omitempty"`
		Label                     string      `json:"label"`
		ReservedVlanId            *Vlan       `json:"reserved_vlan_id,omitempty"`
		RouteTarget               string      `json:"route_target,omitempty"`
		RtPolicy                  *RtPolicy   `json:"rt_policy"`
		SecurityZoneId            ObjectId    `json:"security_zone_id,omitempty"`
		SviIps                    []SviIp     `json:"svi_ips"`
		VirtualGatewayIpv4        string      `json:"virtual_gateway_ipv4,omitempty"`
		VirtualGatewayIpv6        string      `json:"virtual_gateway_ipv6,omitempty"`
		VirtualGatewayIpv4Enabled bool        `json:"virtual_gateway_ipv4_enabled"`
		VirtualGatewayIpv6Enabled bool        `json:"virtual_gateway_ipv6_enabled"`
		VnBindings                []VnBinding `json:"bound_to"`
		VnId                      string      `json:"vn_id,omitempty"`
		VnType                    string      `json:"vn_type"`
		VirtualMac                string      `json:"virtual_mac,omitempty"`
	}{
		Description:               o.Description,
		Ipv4Enabled:               o.Ipv4Enabled,
		Ipv6Enabled:               o.Ipv6Enabled,
		L3Mtu:                     o.L3Mtu,
		Label:                     o.Label,
		ReservedVlanId:            o.ReservedVlanId,
		RouteTarget:               o.RouteTarget,
		RtPolicy:                  o.RtPolicy,
		SecurityZoneId:            o.SecurityZoneId,
		SviIps:                    o.SviIps,
		VirtualGatewayIpv4Enabled: o.VirtualGatewayIpv4Enabled,
		VirtualGatewayIpv6Enabled: o.VirtualGatewayIpv6Enabled,
		VnBindings:                o.VnBindings,
	}

	raw.DhcpService = o.DhcpService.String()
	if o.Ipv4Subnet != nil && o.Ipv4Subnet.IP != nil { // todo: handle zero/removal of existing IP address value
		raw.Ipv4Subnet = o.Ipv4Subnet.String()
	}
	if o.Ipv6Subnet != nil && o.Ipv6Subnet.IP != nil { // todo: handle zero/removal of existing IP address value
		raw.Ipv6Subnet = o.Ipv6Subnet.String()
	}
	if len(o.VirtualGatewayIpv4.To4()) == net.IPv4len { // todo: handle zero/removal o fexisting IP address value
		raw.VirtualGatewayIpv4 = o.VirtualGatewayIpv4.String()
	}
	if len(o.VirtualGatewayIpv6) == net.IPv6len { // todo: handle zero/removal o fexisting IP address value
		raw.VirtualGatewayIpv6 = o.VirtualGatewayIpv6.String()
	}
	if o.VnId != nil {
		raw.VnId = strconv.Itoa(int(*o.VnId))
	}
	raw.VnType = o.VnType.String()
	raw.VirtualMac = o.VirtualMac.String()

	return json.Marshal(raw)
}

func (o *TwoStageL3ClosClient) listAllVirtualNetworkIds(ctx context.Context) ([]ObjectId, error) {
	apstraUrl, err := o.urlWithParam(fmt.Sprintf(apiUrlVirtualNetworks, o.blueprintId))
	if err != nil {
		return nil, err
	}

	response := &struct {
		VirtualNetworks map[ObjectId]VirtualNetwork `json:"virtual_networks"`
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

func (o *TwoStageL3ClosClient) getVirtualNetwork(ctx context.Context, vnId ObjectId) (*VirtualNetwork, error) {
	apstraUrl, err := o.urlWithParam(fmt.Sprintf(apiUrlVirtualNetworkById, o.blueprintId, vnId))
	if err != nil {
		return nil, err
	}

	response := VirtualNetwork{}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodGet,
		url:            apstraUrl,
		apiResponse:    &response,
		unsynchronized: true,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *TwoStageL3ClosClient) getVirtualNetworkByName(ctx context.Context, name string) (*VirtualNetwork, error) {
	rawVns, err := o.getAllVirtualNetworks(ctx)
	if err != nil {
		return nil, err
	}

	var found int
	var rawVn VirtualNetwork
	for i := range rawVns {
		if rawVns[i].Data.Label == name {
			found++
			rawVn = rawVns[i]
		}
	}

	switch found {
	case 0:
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("virtual network with label %q not found", name),
		}
	case 1:
		// re-fetch is required here because data is missing when we "get all".
		return o.getVirtualNetwork(ctx, rawVn.Id)
	default:
		return nil, ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("found %d virtual networks with label %q", found, name),
		}
	}
}

func (o *TwoStageL3ClosClient) getAllVirtualNetworks(ctx context.Context) (map[ObjectId]VirtualNetwork, error) {
	var response struct {
		VirtualNetworks map[ObjectId]VirtualNetwork `json:"virtual_networks"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlVirtualNetworks, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.VirtualNetworks, nil
}

func (o *TwoStageL3ClosClient) createVirtualNetwork(ctx context.Context, cfg *VirtualNetworkData) (ObjectId, error) {
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

func (o *TwoStageL3ClosClient) updateVirtualNetwork(ctx context.Context, id ObjectId, cfg *VirtualNetworkData) error {
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
