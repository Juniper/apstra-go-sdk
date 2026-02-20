// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"encoding/json"
	"net/netip"

	"github.com/Juniper/apstra-go-sdk/enum"
)

var (
	_ json.Marshaler   = (*FreeformAggregateLinkEndpoint)(nil)
	_ json.Unmarshaler = (*FreeformAggregateLinkEndpoint)(nil)
)

// FreeformAggregateLinkEndpoint represents a member of a FreeformAggregateLinkEndpointGroup.
// FreeformAggregateLinkEndpoint is the logical LAG interface (ae1, bond0) on a device. Because
// a LAG may be terminated by a multi-chassis scheme (MLAG, ESI LAG), one endpoint is not enough
// to describe one end of a LAG. Each LAG has exactly two FreeformAggregateLinkEndpointGroup,
// and each FreeformAggregateLinkEndpointGroup has at least one FreeformAggregateLinkEndpoint.
// Because of L3 MLAG schemes (e.g. HSRP routing over Cisco VPC) it is possible for each
// FreeformAggregateLinkEndpoint to have its own IPv4 and IPv6 address.`
type FreeformAggregateLinkEndpoint struct {
	SystemID      string
	IfName        string
	IPv4Addr      *netip.Prefix
	IPv6Addr      *netip.Prefix
	PortChannelID int
	Tags          []string
	LAGMode       enum.LAGMode

	endpointGroup int
	id            string
}

func (o FreeformAggregateLinkEndpoint) ID() *string {
	if o.id == "" {
		return nil
	}
	return &o.id
}

func (o FreeformAggregateLinkEndpoint) MarshalJSON() ([]byte, error) {
	type rawSystem struct {
		ID string `json:"id"`
	}

	type rawInterface struct {
		IfName        string        `json:"if_name,omitempty"`
		PortChannelId int           `json:"port_channel_id"`
		LagMode       enum.LAGMode  `json:"lag_mode"`
		IPv4Addr      *netip.Prefix `json:"ipv4_addr,omitempty"`
		IPv6Addr      *netip.Prefix `json:"ipv6_addr,omitempty"`
		Tags          []string      `json:"tags"`
	}

	raw := struct {
		System        rawSystem    `json:"system"`
		Interface     rawInterface `json:"interface"`
		EndpointGroup int          `json:"endpoint_group"`
	}{
		System: rawSystem{ID: o.SystemID},
		Interface: rawInterface{
			IfName:        o.IfName,
			PortChannelId: o.PortChannelID,
			LagMode:       o.LAGMode,
			IPv4Addr:      o.IPv4Addr,
			IPv6Addr:      o.IPv6Addr,
			Tags:          o.Tags,
		},
		EndpointGroup: o.endpointGroup,
	}

	return json.Marshal(raw)
}

func (o *FreeformAggregateLinkEndpoint) UnmarshalJSON(bytes []byte) error {
	type rawSystem struct {
		ID string `json:"id"`
	}

	type rawInterface struct {
		IfName        string        `json:"if_name,omitempty"`
		PortChannelId int           `json:"port_channel_id"`
		LagMode       enum.LAGMode  `json:"lag_mode"`
		IPv4Addr      *netip.Prefix `json:"ipv4_addr,omitempty"`
		IPv6Addr      *netip.Prefix `json:"ipv6_addr,omitempty"`
		Tags          []string      `json:"tags"`
	}

	var raw struct {
		ID            string       `json:"id"`
		System        rawSystem    `json:"system"`
		Interface     rawInterface `json:"interface"`
		EndpointGroup int          `json:"endpoint_group"`
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	o.id = raw.ID
	o.endpointGroup = raw.EndpointGroup
	o.SystemID = raw.System.ID
	o.IfName = raw.Interface.IfName
	o.IPv4Addr = raw.Interface.IPv4Addr
	o.IPv6Addr = raw.Interface.IPv6Addr
	o.PortChannelID = raw.Interface.PortChannelId
	o.Tags = raw.Interface.Tags
	o.LAGMode = raw.Interface.LagMode

	return nil
}
