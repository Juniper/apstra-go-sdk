// Copyright (c) Juniper Networks, Inc., 2023-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"net/url"

	"github.com/Juniper/apstra-go-sdk/enum"
)

const (
	apiUrlBlueprintExperienceWeb   = apiUrlBlueprintById + apiUrlPathDelim + "experience/web/"
	apiUrlBlueprintGETCablingMap   = apiUrlBlueprintExperienceWeb + "cabling-map"
	apiUrlBlueprintPATCHCablingMap = apiUrlBlueprintById + apiUrlPathDelim + "cabling-map"

	includeLagParam    = "aggregate_links"
	linksBySystemParam = "system_node_id"
)

type CablingMapLink struct {
	ID              string                    `json:"id,omitempty"`
	TagLabels       []string                  `json:"tags,omitempty"`
	Speed           *enum.LinkSpeed           `json:"speed,omitempty"`
	AggregateLinkId *string                   `json:"aggregate_link_id,omitempty"`
	GroupLabel      *string                   `json:"group_label,omitempty"`
	Label           *string                   `json:"label,omitempty"`
	Role            *enum.LinkRole            `json:"role,omitempty"`
	Endpoints       [2]CablingMapLinkEndpoint `json:"endpoints"`
	Type            *enum.LinkType            `json:"type,omitempty"`
	RailID          *string                   `json:"rail_id,omitempty"`
}

type CablingMapLinkEndpoint struct {
	Interface CablingMapLinkEndpointInterface `json:"interface,omitempty"`
	System    *CablingMapLinkEndpointSystem   `json:"system,omitempty"`
}

var _ json.Marshaler = (*CablingMapLinkEndpointInterface)(nil)

type CablingMapLinkEndpointInterface struct {
	Description    *string                       `json:"description,omitempty"`
	TagLabels      []string                      `json:"tags,omitempty"`
	IfType         *enum.InterfaceType           `json:"if_type,omitempty"`
	OperationState *enum.InterfaceOperationState `json:"operation_state,omitempty"`
	IfName         *string                       `json:"if_name,omitempty"`
	PortChannelID  *int                          `json:"port_channel_id,omitempty"`
	ID             string                        `json:"id"`
	LAGMode        *enum.LAGMode                 `json:"lag_mode,omitempty"`
	IPv4Addr       *netip.Prefix                 `json:"ipv4_addr,omitempty"`
	IPv4Enabled    *bool                         `json:"ipv4_enabled,omitempty"`
	IPv6Addr       *netip.Prefix                 `json:"ipv6_addr,omitempty"`
	IPv6Enabled    *bool                         `json:"ipv6_enabled,omitempty"`
}

func (c CablingMapLinkEndpointInterface) MarshalJSON() ([]byte, error) {
	raw := struct {
		ID       string          `json:"id"`
		IfName   json.RawMessage `json:"if_name,omitempty"`
		IPv4Addr json.RawMessage `json:"ipv4_addr,omitempty"`
		IPv6Addr json.RawMessage `json:"ipv6_addr,omitempty"`
	}{ID: c.ID}

	// Clear value from API by sending `null` if value is a pointer to an empty string.
	if c.IfName != nil {
		if *c.IfName == "" {
			raw.IfName = []byte("null")
		} else {
			raw.IfName, _ = json.Marshal(c.IfName)
		}
	}

	// Clear IPv4 value from API by sending `null` if value is a pointer to an invalid prefix.
	if c.IPv4Addr != nil {
		if c.IPv4Addr.IsValid() {
			raw.IPv4Addr, _ = json.Marshal(c.IPv4Addr)
		} else {
			raw.IPv4Addr = []byte("null")
		}
	}

	// Clear IPv6 value from API by sending `null` if value is a pointer to an invalid prefix.
	if c.IPv6Addr != nil {
		if c.IPv6Addr.IsValid() {
			raw.IPv6Addr, _ = json.Marshal(c.IPv6Addr)
		} else {
			raw.IPv6Addr = []byte("null")
		}
	}

	return json.Marshal(raw)
}

type CablingMapLinkEndpointSystem struct {
	ID    string              `json:"id"`
	Role  enum.SystemNodeRole `json:"role"`
	Label *string             `json:"label,omitempty"`
}

// Digest returns a string which uniquely identifies the endpoint by the system
// and interface where it terminates. It returns *string like "abc123:xe-0/0/1"
// or "def456:eth0"
// If any of the required elements are nil, nil is returned.
func (o CablingMapLinkEndpoint) Digest() *string {
	if o.System == nil || o.Interface.IfName == nil {
		return nil
	}

	result := o.System.ID + ":" + *o.Interface.IfName
	return &result
}

// EndpointBySystemID returns the first (likely only) *CablingMapLinkEndpoint
// connected to the specified system.
func (o CablingMapLink) EndpointBySystemID(systemId string) *CablingMapLinkEndpoint {
	for _, endpoint := range o.Endpoints {
		if endpoint.System == nil {
			continue
		}
		if endpoint.System.ID == systemId {
			return &endpoint
		}
	}
	return nil
}

// OppositeEndpointBySystemID does the opposite of EndpointBySystemID. Rather
// than returning the specified endpoint, it returns the other one. Returns
// nil if both ends of the link land on the specified system or if no ends of
// the link end on the specified system.
func (o CablingMapLink) OppositeEndpointBySystemID(systemId string) *CablingMapLinkEndpoint {
	if o.Endpoints[0].System != nil &&
		o.Endpoints[1].System != nil &&
		o.Endpoints[0].System.ID == o.Endpoints[1].System.ID {
		// can't find an 'opposite of' system if the same system is on both ends
		return nil
	}

	for i, endpoint := range o.Endpoints {
		if endpoint.System == nil {
			continue
		}
		if endpoint.System.ID == systemId {
			return &o.Endpoints[(i+1)%2]
		}
	}
	return nil
}

// GetCablingMapLinks returns []CablingMapLink representing every link in the blueprint
func (o *TwoStageL3ClosClient) GetCablingMapLinks(ctx context.Context) ([]CablingMapLink, error) {
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintGETCablingMap, o.blueprintId))
	if err != nil {
		return nil, err
	}

	params := apstraUrl.Query()
	params.Set(includeLagParam, "true")
	apstraUrl.RawQuery = params.Encode()

	response := struct {
		Links []CablingMapLink `json:"links"`
	}{}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Links, nil
}

// GetCablingMapLinksBySystem returns []CablingMapLink representing every link (including LAGs)
func (o *TwoStageL3ClosClient) GetCablingMapLinksBySystem(ctx context.Context, systemNodeId string) ([]CablingMapLink, error) {
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintGETCablingMap, o.blueprintId))
	if err != nil {
		return nil, err
	}

	params := apstraUrl.Query()
	params.Set(includeLagParam, "true")
	params.Set(linksBySystemParam, systemNodeId)
	apstraUrl.RawQuery = params.Encode()

	response := struct {
		Links []CablingMapLink `json:"links"`
	}{}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Links, nil
}

func (o *TwoStageL3ClosClient) PatchCablingMapLinks(ctx context.Context, in []CablingMapLink) error {
	// make a copy of the input slice which we'll send to the API
	payload := struct {
		Links []CablingMapLink `json:"links"`
	}{
		Links: make([]CablingMapLink, len(in)),
	}

	// Only populate the patch-able fields into the payload.
	for i, link := range in {
		payload.Links[i] = CablingMapLink{
			Endpoints: [2]CablingMapLinkEndpoint{
				{Interface: CablingMapLinkEndpointInterface{ID: link.Endpoints[0].Interface.ID}},
				{Interface: CablingMapLinkEndpointInterface{ID: link.Endpoints[1].Interface.ID}},
			},
		}
		if link.Endpoints[0].Interface.IfName != nil {
			payload.Links[i].Endpoints[0].Interface.IfName = link.Endpoints[0].Interface.IfName
		}
		if link.Endpoints[0].Interface.IPv4Addr != nil {
			payload.Links[i].Endpoints[0].Interface.IPv4Addr = link.Endpoints[0].Interface.IPv4Addr
		}
		if link.Endpoints[0].Interface.IPv6Addr != nil {
			payload.Links[i].Endpoints[0].Interface.IPv6Addr = link.Endpoints[0].Interface.IPv6Addr
		}
		if link.Endpoints[1].Interface.IfName != nil {
			payload.Links[i].Endpoints[1].Interface.IfName = link.Endpoints[1].Interface.IfName
		}
		if link.Endpoints[1].Interface.IPv4Addr != nil {
			payload.Links[i].Endpoints[1].Interface.IPv4Addr = link.Endpoints[1].Interface.IPv4Addr
		}
		if link.Endpoints[1].Interface.IPv6Addr != nil {
			payload.Links[i].Endpoints[1].Interface.IPv6Addr = link.Endpoints[1].Interface.IPv6Addr
		}
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintPATCHCablingMap, o.blueprintId),
		apiInput: payload,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
