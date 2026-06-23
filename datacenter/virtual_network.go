// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/errors"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/internal/parse"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
)

var (
	_ internal.IDer     = (*VirtualNetwork)(nil)
	_ internal.IDSetter = (*VirtualNetwork)(nil)
	_ json.Unmarshaler  = (*VirtualNetwork)(nil)
	_ json.Marshaler    = (*VirtualNetwork)(nil)
)

type VirtualNetwork struct {
	Bindings                  []VNBinding        `json:"bound_to"`
	Description               string             `json:"description"`
	DHCPService               DHCPServiceEnabled `json:"dhcp_service"`
	IPv4Enabled               bool               `json:"ipv4_enabled"`
	IPv4Subnet                *net.IPNet         `json:"-"`
	IPv6Enabled               bool               `json:"ipv6_enabled"`
	IPv6Subnet                *net.IPNet         `json:"-"`
	L3MTU                     *int               `json:"l3_mtu,omitempty"`
	Label                     string             `json:"label"`
	ReservedVLAN              *uint16            `json:"reserved_vlan_id,omitempty"`
	RTPolicy                  *RTPolicy          `json:"rt_policy"`
	SecurityZoneID            string             `json:"security_zone_id,omitempty"`
	SVIIPs                    []SVIAddressing    `json:"svi_ips"`
	Tags                      []string           `json:"tags"`
	Type                      enum.VnType        `json:"vn_type"`
	VirtualGatewayIPv4        net.IP             `json:"-"`
	VirtualGatewayIPv6        net.IP             `json:"-"`
	VirtualGatewayIPv4Enabled bool               `json:"virtual_gateway_ipv4_enabled"`
	VirtualGatewayIPv6Enabled bool               `json:"virtual_gateway_ipv6_enabled"`
	VNI                       *uint32            `json:"vn_id,omitempty,string"` // VNI as a string, null when unset
	VirtualMAC                net.HardwareAddr   `json:"-"`

	id string
}

// ID returns a pointer to a copy of the object's ID, or nil when no ID is set.
func (vn VirtualNetwork) ID() *string {
	if vn.id == "" {
		return nil
	}
	return pointer.ToCopyOf(vn.id)
}

func (vn *VirtualNetwork) SetID(id string) error {
	if vn.id != "" {
		return errors.IDAlreadySet(fmt.Sprintf("id already has value %q", vn.id))
	}

	vn.id = id
	return nil
}

func (vn VirtualNetwork) MarshalJSON() ([]byte, error) {
	var ipv4Subnet, ipv6Subnet, virtualGatewayIPv4, virtualGatewayIPv6, virtualMAC string
	if vn.IPv4Subnet != nil && vn.IPv4Subnet.IP != nil { // todo: handle zero/removal of existing IP address value
		ipv4Subnet = vn.IPv4Subnet.String()
	}
	if vn.IPv6Subnet != nil && vn.IPv6Subnet.IP != nil { // todo: handle zero/removal of existing IP address value
		ipv6Subnet = vn.IPv6Subnet.String()
	}
	if len(vn.VirtualGatewayIPv4.To4()) == net.IPv4len { // todo: handle zero/removal o fexisting IP address value
		virtualGatewayIPv4 = vn.VirtualGatewayIPv4.String()
	}
	if len(vn.VirtualGatewayIPv6) == net.IPv6len { // todo: handle zero/removal o fexisting IP address value
		virtualGatewayIPv6 = vn.VirtualGatewayIPv6.String()
	}
	virtualMAC = vn.VirtualMAC.String()

	type virtualNetwork VirtualNetwork
	return json.Marshal(struct {
		virtualNetwork
		IPv4Subnet         string `json:"ipv4_subnet,omitempty"`
		IPv6Subnet         string `json:"ipv6_subnet,omitempty"`
		VirtualGatewayIPv4 string `json:"virtual_gateway_ipv4,omitempty"`
		VirtualGatewayIPv6 string `json:"virtual_gateway_ipv6,omitempty"`
		VirtualMAC         string `json:"virtual_mac,omitempty"`
	}{
		virtualNetwork:     virtualNetwork(vn),
		IPv4Subnet:         ipv4Subnet,
		IPv6Subnet:         ipv6Subnet,
		VirtualGatewayIPv4: virtualGatewayIPv4,
		VirtualGatewayIPv6: virtualGatewayIPv6,
		VirtualMAC:         virtualMAC,
	})
}

func (vn *VirtualNetwork) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID                        string             `json:"id"`
		Description               string             `json:"description"`
		DHCPService               DHCPServiceEnabled `json:"dhcp_service"`
		IPv4Enabled               bool               `json:"ipv4_enabled"`
		IPv4Subnet                string             `json:"ipv4_subnet"`
		IPv6Enabled               bool               `json:"ipv6_enabled"`
		IPv6Subnet                string             `json:"ipv6_subnet"`
		L3MTU                     *int               `json:"l3_mtu"`
		Label                     string             `json:"label"`
		ReservedVLAN              *uint16            `json:"reserved_vlan_id"`
		RTPolicy                  *RTPolicy          `json:"rt_policy"`
		SecurityZoneID            string             `json:"security_zone_id"`
		SVIIPs                    []SVIAddressing    `json:"svi_ips"`
		Tags                      []string           `json:"tags"`
		VirtualGatewayIPv4        string             `json:"virtual_gateway_ipv4"`
		VirtualGatewayIPv6        string             `json:"virtual_gateway_ipv6"`
		VirtualGatewayIPv4Enabled bool               `json:"virtual_gateway_ipv4_enabled"`
		VirtualGatewayIPv6Enabled bool               `json:"virtual_gateway_ipv6_enabled"`
		VNBindings                []VNBinding        `json:"bound_to"`
		VNI                       string             `json:"vn_id"`
		VNType                    string             `json:"vn_type"`
		VirtualMAC                string             `json:"virtual_mac"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("while unmarshaling raw API response - %w", err)
	}

	vn.id = raw.ID

	vn.Description = raw.Description
	vn.DHCPService = raw.DHCPService
	vn.IPv4Enabled = raw.IPv4Enabled
	vn.IPv4Subnet, err = parse.IPNetFromString(raw.IPv4Subnet)
	if err != nil {
		return fmt.Errorf("while parsing virtual network data ipv4_subnet %q - %w", raw.IPv4Subnet, err)
	}

	vn.IPv6Enabled = raw.IPv6Enabled
	vn.IPv6Subnet, err = parse.IPNetFromString(raw.IPv6Subnet)
	if err != nil {
		return fmt.Errorf("while parsing virtual network ipv6_subnet %q - %w", raw.IPv6Subnet, err)
	}

	vn.L3MTU = raw.L3MTU
	vn.Label = raw.Label
	vn.ReservedVLAN = raw.ReservedVLAN
	vn.RTPolicy = raw.RTPolicy
	vn.SecurityZoneID = raw.SecurityZoneID
	vn.SVIIPs = raw.SVIIPs
	vn.Tags = raw.Tags

	vn.VirtualGatewayIPv4, err = parse.IPFromString(raw.VirtualGatewayIPv4)
	if err != nil {
		return fmt.Errorf("while parsing virtual network virtual_gateway_ipv4 %q - %w", raw.VirtualGatewayIPv4, err)
	}

	vn.VirtualGatewayIPv6, err = parse.IPFromString(raw.VirtualGatewayIPv6)
	if err != nil {
		return fmt.Errorf("while parsing virtual network virtual_gateway_ipv6 %q - %w", raw.VirtualGatewayIPv6, err)
	}

	vn.VirtualGatewayIPv4Enabled = raw.VirtualGatewayIPv4Enabled
	vn.VirtualGatewayIPv6Enabled = raw.VirtualGatewayIPv6Enabled
	vn.Bindings = raw.VNBindings

	if raw.VNI != "" {
		vnId, err := strconv.Atoi(raw.VNI)
		if err != nil {
			return fmt.Errorf("while parsing virtual network vn_id %q - %w", raw.VNI, err)
		}
		vn.VNI = pointer.To(uint32(vnId))
	}

	err = vn.Type.FromString(raw.VNType)
	if err != nil {
		return fmt.Errorf("while parsing virtual network vn_type %q - %w", raw.VNType, err)
	}

	vn.VirtualMAC, err = parse.MACFromString(raw.VirtualMAC)
	if err != nil {
		return fmt.Errorf("while parsing virtual network virtual_mac %q - %w", raw.VirtualMAC, err)
	}

	return nil
}
