// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/Juniper/apstra-go-sdk/enum"
)

var (
	_ json.Marshaler   = (*SVIAddressing)(nil)
	_ json.Unmarshaler = (*SVIAddressing)(nil)
)

type SVIAddressing struct {
	SystemID string
	IPv4Addr *net.IPNet
	IPv4Mode enum.IPv4SVIMode
	IPv6Addr *net.IPNet
	IPv6Mode enum.IPv6SVIMode
}

func (o *SVIAddressing) MarshalJSON() ([]byte, error) {
	var raw struct {
		IPv4Addr string `json:"ipv4_addr,omitempty"`
		IPv4Mode string `json:"ipv4_mode,omitempty"`
		IPv6Addr string `json:"ipv6_addr,omitempty"`
		IPv6Mode string `json:"ipv6_mode,omitempty"`
		SystemId string `json:"system_id"`
	}

	if o.IPv4Addr != nil {
		raw.IPv4Addr = o.IPv4Addr.String()
	}
	raw.IPv4Mode = o.IPv4Mode.String()

	if o.IPv6Addr != nil {
		raw.IPv6Addr = o.IPv6Addr.String()
	}
	raw.IPv6Mode = o.IPv6Mode.String()

	raw.SystemId = o.SystemID

	return json.Marshal(raw)
}

func (o *SVIAddressing) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		IPv4Addr string `json:"ipv4_addr"`
		IPv4Mode string `json:"ipv4_mode"`
		IPv6Addr string `json:"ipv6_addr"`
		IPv6Mode string `json:"ipv6_mode"`
		SystemId string `json:"system_id"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("while unmarshaling an SVIAddressing - %w", err)
	}

	o.IPv4Addr = nil
	if raw.IPv4Addr != "" {
		var ip net.IP
		ip, o.IPv4Addr, err = net.ParseCIDR(raw.IPv4Addr)
		if err != nil {
			return fmt.Errorf("while parsing SVIAddressing.IPv4Addr - %w", err)
		}
		o.IPv4Addr.IP = ip
	}

	err = o.IPv4Mode.FromString(raw.IPv4Mode)
	if err != nil {
		return fmt.Errorf("while parsing SVIAddressing.IPv4Mode - %w", err)
	}

	o.IPv6Addr = nil
	if raw.IPv6Addr != "" {
		var ip net.IP
		ip, o.IPv6Addr, err = net.ParseCIDR(raw.IPv6Addr)
		if err != nil {
			return fmt.Errorf("while parsing SVIAddressing.IPv6Addr - %w", err)
		}
		o.IPv6Addr.IP = ip
	}

	err = o.IPv6Mode.FromString(raw.IPv6Mode)
	if err != nil {
		return fmt.Errorf("while parsing SVIAddressing.IPv6Mode - %w", err)
	}

	o.SystemID = raw.SystemId

	return nil
}
