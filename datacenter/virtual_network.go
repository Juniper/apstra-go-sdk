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
	"github.com/Juniper/apstra-go-sdk/internal/str"
)

var (
	_ internal.IDer     = (*VirtualNetwork)(nil)
	_ internal.IDSetter = (*VirtualNetwork)(nil)
	_ json.Unmarshaler  = (*VirtualNetwork)(nil)
	_ json.Marshaler    = (*VirtualNetwork)(nil)
)

type VirtualNetwork struct {
	Bindings                  []VNBinding         `json:"bound_to,omitzero"`
	Description               *string             `json:"description,omitempty"` // pointer to empty string sends null to the API
	DHCPService               *DHCPServiceEnabled `json:"dhcp_service,omitempty"`
	EncapsulateInnerVLAN      *bool               `json:"encapsulate_inner_vlan,omitempty"`
	IPv4Enabled               *bool               `json:"ipv4_enabled,omitempty"`
	IPv4Subnet                *net.IPNet          `json:"-"` // pointer to zero value sends null to the API
	IPv6Enabled               *bool               `json:"ipv6_enabled,omitempty"`
	IPv6Subnet                *net.IPNet          `json:"-"` // pointer to zero value sends null to the API
	L3MTU                     *int                `json:"-"` // pointer to zero value sends null to the API
	Label                     *string             `json:"label,omitempty"`
	ReservedVLAN              *uint16             `json:"-"` // pointer to zero value sends null to the API
	RTPolicy                  *RTPolicy           `json:"-"` // pointer to a zero-value RTPolicy will send null to the API
	SVIIPs                    []SVIAddressing     `json:"svi_ips,omitzero"`
	SecurityZoneID            string              `json:"security_zone_id,omitempty"`
	SwitchingZoneID           string              `json:"switching_zone_id,omitempty"`
	Tags                      []string            `json:"tags,omitzero"`
	Type                      enum.VnType         `json:"vn_type,omitzero"`
	VirtualGatewayIPv4        net.IP              `json:"-"` // empty slice sends null to the API
	VirtualGatewayIPv6        net.IP              `json:"-"` // empty slice sends null to the API
	VirtualGatewayIPv4Enabled *bool               `json:"virtual_gateway_ipv4_enabled,omitempty"`
	VirtualGatewayIPv6Enabled *bool               `json:"virtual_gateway_ipv6_enabled,omitempty"`
	VNI                       *uint32             `json:"vn_id,omitempty,string"` // VNI is a string
	VirtualMAC                net.HardwareAddr    `json:"-"`                      // empty slice sends null to the API

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
	if vn.IPv4Enabled != nil && !*vn.IPv4Enabled && vn.IPv4Subnet != nil {
		return nil, fmt.Errorf("marshaling virtual network: IPv4 is disabled but IPv4 subnet is set to %q", vn.IPv4Subnet.String())
	}
	if vn.IPv4Enabled != nil && !*vn.IPv4Enabled && vn.VirtualGatewayIPv4 != nil {
		return nil, fmt.Errorf("marshaling virtual network: IPv4 is disabled but virtual gateway IPv4 is set to %q", vn.VirtualGatewayIPv4.String())
	}
	if vn.IPv6Enabled != nil && !*vn.IPv6Enabled && vn.IPv6Subnet != nil {
		return nil, fmt.Errorf("marshaling virtual network: IPv6 is disabled but IPv6 subnet is set to %q", vn.IPv6Subnet.String())
	}
	if vn.IPv6Enabled != nil && !*vn.IPv6Enabled && vn.VirtualGatewayIPv6 != nil {
		return nil, fmt.Errorf("marshaling virtual network: IPv6 is disabled but virtual gateway IPv6 is set to %q", vn.VirtualGatewayIPv6.String())
	}

	var description, ipv4Subnet, ipv6Subnet, l3MTU, reservedVLAN, rtPolicy, virtualGatewayIPv4, virtualGatewayIPv6, virtualMAC json.RawMessage
	if vn.Description != nil {
		if *vn.Description == "" {
			description = json.RawMessage("null") // explicit removal via 'null' when description exists but is empty
		} else {
			description = str.QuoteJSONString(*vn.Description) // send the provided value
		}
	}
	if vn.IPv4Subnet != nil { // a value was provided, so we'll send *something*
		if vn.IPv4Subnet.IP.IsUnspecified() || len(vn.IPv4Subnet.IP) == 0 {
			ipv4Subnet = json.RawMessage("null") // explicit removal via nil, empty or unspecified subnet IP; send null
		} else {
			ipv4Subnet = str.QuoteJSONString(vn.IPv4Subnet.String()) // send the provided value
		}
	}
	if vn.L3MTU != nil { // a value was provided, so we'll send *something*
		if *vn.L3MTU == 0 {
			l3MTU = json.RawMessage("null") // explicit removal via pointer to zero value; send null
		} else {
			l3MTU = json.RawMessage(strconv.Itoa(*vn.L3MTU))
		}
	}
	if vn.RTPolicy != nil { // a value was provided, so we'll send *something*
		if vn.RTPolicy.ExportRTs == nil || vn.RTPolicy.ImportRTs == nil {
			rtPolicy = json.RawMessage("null") // zero value RTPolicy; send null to the API
		} else {
			var err error
			rtPolicy, err = json.Marshal(vn.RTPolicy) // marshal the RTPolicy to JSON and send the provided value
			if err != nil {
				return nil, fmt.Errorf("marshaling virtual network: marshaling RTPolicy: %w", err)
			}
		}
	}
	if vn.VirtualGatewayIPv4 != nil { // a value was provided, so we'll send *something*
		if vn.VirtualGatewayIPv4.IsUnspecified() || len(vn.VirtualGatewayIPv4) == 0 {
			virtualGatewayIPv4 = json.RawMessage("null") // explicit removal via nil, empty or unspecified virtual gateway IP; send null
		} else {
			virtualGatewayIPv4 = str.QuoteJSONString(vn.VirtualGatewayIPv4.String()) // send the provided value
		}
	}
	if vn.IPv6Subnet != nil { // a value was provided, so we'll send *something*
		if vn.IPv6Subnet.IP.IsUnspecified() || len(vn.IPv6Subnet.IP) == 0 {
			ipv6Subnet = json.RawMessage("null") // explicit removal via nil, empty or unspecified subnet IP; send null
		} else {
			ipv6Subnet = str.QuoteJSONString(vn.IPv6Subnet.String()) // send the provided value
		}
	}
	if vn.VirtualGatewayIPv6 != nil { // a value was provided, so we'll send *something*
		if vn.VirtualGatewayIPv6.IsUnspecified() || len(vn.VirtualGatewayIPv6) == 0 {
			virtualGatewayIPv6 = json.RawMessage("null") // explicit removal via nil, empty or unspecified virtual gateway IP; send null
		} else {
			virtualGatewayIPv6 = str.QuoteJSONString(vn.VirtualGatewayIPv6.String()) // send the provided value
		}
	}
	if vn.ReservedVLAN != nil { // a value was provided, so we'll send *something*
		if *vn.ReservedVLAN == 0 {
			reservedVLAN = json.RawMessage("null") // explicit removal via 'null'
		} else {
			reservedVLAN = json.RawMessage(strconv.Itoa(int(*vn.ReservedVLAN))) // send the provided value
		}
	}
	if vn.VirtualMAC != nil { // a value was provided, so we'll send *something*
		if len(vn.VirtualMAC) == 0 {
			virtualMAC = json.RawMessage("null") // explicit removal via 'null' when MAC address exists but is empty
		} else {
			virtualMAC = str.QuoteJSONString(vn.VirtualMAC.String()) // send the provided value
		}
	}
	if vn.IPv4Enabled != nil {
		if !*vn.IPv4Enabled {
			ipv4Subnet = json.RawMessage("null")         // ipv4Subnet must be cleared when IPv4 is disabled
			virtualGatewayIPv4 = json.RawMessage("null") // virtualGatewayIPv4 must be cleared when IPv4 is disabled
		}
	}
	if vn.IPv6Enabled != nil {
		if !*vn.IPv6Enabled {
			ipv6Subnet = json.RawMessage("null")         // ipv6Subnet must be cleared when IPv6 is disabled
			virtualGatewayIPv6 = json.RawMessage("null") // virtualGatewayIPv6 must be cleared when IPv6 is disabled
		}
	}

	type virtualNetwork VirtualNetwork
	return json.Marshal(struct {
		virtualNetwork
		Description        json.RawMessage `json:"description,omitempty"`
		IPv4Subnet         json.RawMessage `json:"ipv4_subnet,omitempty"`
		IPv6Subnet         json.RawMessage `json:"ipv6_subnet,omitempty"`
		L3MTU              json.RawMessage `json:"l3_mtu,omitempty"`
		ReservedVLAN       json.RawMessage `json:"reserved_vlan_id,omitempty"`
		RTPolicy           json.RawMessage `json:"rt_policy,omitempty"`
		VirtualGatewayIPv4 json.RawMessage `json:"virtual_gateway_ipv4,omitempty"`
		VirtualGatewayIPv6 json.RawMessage `json:"virtual_gateway_ipv6,omitempty"`
		VirtualMAC         json.RawMessage `json:"virtual_mac,omitempty"`
	}{
		Description:        description,
		virtualNetwork:     virtualNetwork(vn),
		IPv4Subnet:         ipv4Subnet,
		IPv6Subnet:         ipv6Subnet,
		L3MTU:              l3MTU,
		ReservedVLAN:       reservedVLAN,
		RTPolicy:           rtPolicy,
		VirtualGatewayIPv4: virtualGatewayIPv4,
		VirtualGatewayIPv6: virtualGatewayIPv6,
		VirtualMAC:         virtualMAC,
	})
}

func (vn *VirtualNetwork) UnmarshalJSON(bytes []byte) error {
	// Most of these don't need to be pointers because we'll always get someting back from the API.
	var raw struct {
		ID                        string             `json:"id"`
		Description               string             `json:"description"`
		DHCPService               DHCPServiceEnabled `json:"dhcp_service"`
		EncapsulateInnerVLAN      bool               `json:"encapsulate_inner_vlan"`
		IPv4Enabled               bool               `json:"ipv4_enabled"`
		IPv4Subnet                string             `json:"ipv4_subnet"`
		IPv6Enabled               bool               `json:"ipv6_enabled"`
		IPv6Subnet                string             `json:"ipv6_subnet"`
		L3MTU                     *int               `json:"l3_mtu"` // not sure if API always has an L3 MTU, so this is a pointer to be safe
		Label                     string             `json:"label"`
		ReservedVLAN              *uint16            `json:"reserved_vlan_id"` // not required to exist
		RTPolicy                  *RTPolicy          `json:"rt_policy"`        // not sure if API always has an RT policy, so this is a pointer to be safe
		SVIIPs                    []SVIAddressing    `json:"svi_ips"`
		SecurityZoneID            string             `json:"security_zone_id"`
		SwitchingZoneID           string             `json:"switching_zone_id"`
		Tags                      []string           `json:"tags"`
		VirtualGatewayIPv4        string             `json:"virtual_gateway_ipv4"`
		VirtualGatewayIPv6        string             `json:"virtual_gateway_ipv6"`
		VirtualGatewayIPv4Enabled bool               `json:"virtual_gateway_ipv4_enabled"`
		VirtualGatewayIPv6Enabled bool               `json:"virtual_gateway_ipv6_enabled"`
		VNBindings                []VNBinding        `json:"bound_to"`
		VNI                       *uint32            `json:"vn_id,string"` // vni may not exist
		VNType                    enum.VnType        `json:"vn_type"`
		VirtualMAC                string             `json:"virtual_mac"` // may be null - conversion to empty string is okay for our parser
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("while unmarshaling raw API response - %w", err)
	}

	vn.id = raw.ID

	// Convert empty description to nil because an empty description is not useful and this
	// is consistent with our marshaling behavior (where a pointer to an empty string is
	// marshaled as null to signal removal).
	if raw.Description != "" {
		vn.Description = &raw.Description
	} else {
		vn.Description = nil
	}

	vn.DHCPService = &raw.DHCPService
	vn.EncapsulateInnerVLAN = &raw.EncapsulateInnerVLAN

	vn.IPv4Enabled = &raw.IPv4Enabled
	vn.IPv4Subnet, err = parse.IPNetFromString(raw.IPv4Subnet)
	if err != nil {
		return fmt.Errorf("while parsing virtual network data ipv4_subnet %q - %w", raw.IPv4Subnet, err)
	}

	vn.IPv6Enabled = &raw.IPv6Enabled
	vn.IPv6Subnet, err = parse.IPNetFromString(raw.IPv6Subnet)
	if err != nil {
		return fmt.Errorf("while parsing virtual network ipv6_subnet %q - %w", raw.IPv6Subnet, err)
	}

	vn.L3MTU = raw.L3MTU
	vn.Label = &raw.Label
	vn.ReservedVLAN = raw.ReservedVLAN
	vn.RTPolicy = raw.RTPolicy
	vn.SVIIPs = raw.SVIIPs
	vn.SecurityZoneID = raw.SecurityZoneID
	vn.SwitchingZoneID = raw.SwitchingZoneID
	vn.Tags = raw.Tags

	vn.VirtualGatewayIPv4, err = parse.IPFromString(raw.VirtualGatewayIPv4)
	if err != nil {
		return fmt.Errorf("while parsing virtual network virtual_gateway_ipv4 %q - %w", raw.VirtualGatewayIPv4, err)
	}

	vn.VirtualGatewayIPv6, err = parse.IPFromString(raw.VirtualGatewayIPv6)
	if err != nil {
		return fmt.Errorf("while parsing virtual network virtual_gateway_ipv6 %q - %w", raw.VirtualGatewayIPv6, err)
	}

	vn.VirtualGatewayIPv4Enabled = &raw.VirtualGatewayIPv4Enabled
	vn.VirtualGatewayIPv6Enabled = &raw.VirtualGatewayIPv6Enabled
	vn.Bindings = raw.VNBindings
	vn.VNI = raw.VNI
	vn.Type = raw.VNType

	vn.VirtualMAC, err = parse.MACFromString(raw.VirtualMAC)
	if err != nil {
		return fmt.Errorf("while parsing virtual network virtual_mac %q - %w", raw.VirtualMAC, err)
	}

	return nil
}
