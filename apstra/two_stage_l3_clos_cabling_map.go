// Copyright (c) Juniper Networks, Inc., 2023-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Juniper/apstra-go-sdk/enum"
)

const (
	apiUrlBlueprintExperienceWeb = apiUrlBlueprintById + apiUrlPathDelim + "experience/web/"
	apiUrlBlueprintCablingMap    = apiUrlBlueprintExperienceWeb + "cabling-map"

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

type CablingMapLinkEndpointInterface struct {
	Description    *string                       `json:"description,omitempty"`
	TagLabels      []string                      `json:"tags,omitempty"`
	IfType         *enum.InterfaceType           `json:"if_type,omitempty"`
	OperationState *enum.InterfaceOperationState `json:"operation_state,omitempty"`
	IfName         *string                       `json:"if_name,omitempty"`
	PortChannelID  *int                          `json:"port_channel_id,omitempty"`
	ID             string                        `json:"id"`
	LAGMode        *enum.LAGMode                 `json:"lag_mode,omitempty"`
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
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintCablingMap, o.blueprintId))
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
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintCablingMap, o.blueprintId))
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
