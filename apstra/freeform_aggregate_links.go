// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"
)

const (
	apiUrlFfAggLinks    = apiUrlBlueprintById + apiUrlPathDelim + "aggregate-links"
	apiUrlFfAggLinkById = apiUrlFfAggLinks + apiUrlPathDelim + "%s"
)

type FreeformAggregateLinkMemberEndpoint struct {
	AggIntfId     ObjectId // not used in create
	SystemId      ObjectId
	PortChannelId int
	LagMode       RackLinkLagMode
	Ipv4Address   netip.Prefix
	Ipv6Address   netip.Prefix
}

var _ json.Marshaler = (*FreeformAggregateLinkData)(nil)

type FreeformAggregateLinkData struct {
	Label         string
	Endpoints     [2][]FreeformAggregateLinkMemberEndpoint
	MemberLinkIds []ObjectId
	Tags          []string
}

func (o FreeformAggregateLinkData) MarshalJSON() ([]byte, error) {
	type Endpoint struct {
		System struct {
			Id ObjectId `json:"id"`
		} `json:"system"`
		Interface struct {
			Id            ObjectId    `json:"id,omitempty"`
			PortChannelId int         `json:"port_channel_id"`
			LagMode       string      `json:"lag_mode"`
			Ipv4Address   interface{} `json:"ipv4_addr"`
			Ipv6Address   interface{} `json:"ipv6_addr"`
		} `json:"interface"`
		EndpointGroup int `json:"endpoint_group"`
	}

	var raw struct {
		Label         string     `json:"label"`
		Endpoints     []Endpoint `json:"endpoints"`
		MemberLinkIds []ObjectId `json:"member_link_ids"`
		Tags          []string   `json:"tags"`
	}

	raw.Label = o.Label
	raw.MemberLinkIds = o.MemberLinkIds

	raw.Endpoints = make([]Endpoint, len(o.Endpoints[0])+len(o.Endpoints[1]))

	for i, endpoint := range o.Endpoints[0] {
		var Ipv4Address interface{}
		var Ipv6Address interface{}
		if o.Endpoints[0][i].Ipv4Address.IsValid() {
			Ipv4Address = o.Endpoints[0][i].Ipv4Address.String()
		} else {
			Ipv4Address = nil
		}
		if o.Endpoints[0][i].Ipv6Address.IsValid() {
			Ipv6Address = o.Endpoints[0][i].Ipv6Address.String()
		} else {
			Ipv6Address = nil
		}
		raw.Endpoints[i] = Endpoint{
			System: struct {
				Id ObjectId `json:"id"`
			}{
				Id: endpoint.SystemId,
			},
			Interface: struct {
				Id            ObjectId    `json:"id,omitempty"`
				PortChannelId int         `json:"port_channel_id"`
				LagMode       string      `json:"lag_mode"`
				Ipv4Address   interface{} `json:"ipv4_addr"`
				Ipv6Address   interface{} `json:"ipv6_addr"`
			}{
				Id:            endpoint.AggIntfId,
				PortChannelId: endpoint.PortChannelId,
				LagMode:       endpoint.LagMode.String(),
				Ipv4Address:   Ipv4Address,
				Ipv6Address:   Ipv6Address,
			},
			EndpointGroup: 0,
		}
	}

	skip := len(o.Endpoints[0])
	for i, endpoint := range o.Endpoints[1] {
		var Ipv4Address interface{}
		var Ipv6Address interface{}
		if o.Endpoints[1][i].Ipv4Address.IsValid() {
			Ipv4Address = o.Endpoints[1][i].Ipv4Address.String()
		} else {
			Ipv4Address = nil
		}
		if o.Endpoints[1][i].Ipv6Address.IsValid() {
			Ipv6Address = o.Endpoints[1][i].Ipv6Address.String()
		} else {
			Ipv6Address = nil
		}
		raw.Endpoints[i+skip] = Endpoint{
			System: struct {
				Id ObjectId `json:"id"`
			}{
				Id: endpoint.SystemId,
			},
			Interface: struct {
				Id            ObjectId    `json:"id,omitempty"`
				PortChannelId int         `json:"port_channel_id"`
				LagMode       string      `json:"lag_mode"`
				Ipv4Address   interface{} `json:"ipv4_addr"`
				Ipv6Address   interface{} `json:"ipv6_addr"`
			}{
				Id:            endpoint.AggIntfId,
				PortChannelId: endpoint.PortChannelId,
				LagMode:       endpoint.LagMode.String(),
				Ipv4Address:   Ipv4Address,
				Ipv6Address:   Ipv6Address,
			},
			EndpointGroup: 1,
		}
	}
	raw.Tags = o.Tags
	return json.Marshal(raw)
}

var _ json.Unmarshaler = new(FreeformAggregateLink)

type FreeformAggregateLink struct {
	Id   ObjectId
	Data *FreeformAggregateLinkData
}

func (o *FreeformAggregateLink) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id        ObjectId `json:"id"`
		Label     string   `json:"label"`
		Endpoints []struct {
			System struct {
				Id ObjectId `json:"id"`
			} `json:"system"`
			Interface struct {
				Id ObjectId `json:"id"`
				// IfName     string   `json:"if_name"`
				PortChannelId int    `json:"port_channel_id"`
				LagMode       string `json:"lag_mode"`
				Ipv4Addr      string `json:"ipv4_addr"`
				Ipv6Addr      string `json:"ipv6_addr"`
				// Tags       []string `json:"tags"`
			} `json:"interface"`
			EndpointGroup int `json:"endpoint_group"`
		} `json:"endpoints"`
		MemberLinkIds []ObjectId `json:"member_link_ids"`
		Tags          []string   `json:"tags"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(FreeformAggregateLinkData)
	o.Data.Label = raw.Label
	o.Data.MemberLinkIds = raw.MemberLinkIds
	o.Data.Tags = raw.Tags
	for _, endpoint := range raw.Endpoints {
		if endpoint.EndpointGroup < 0 || endpoint.EndpointGroup > 2 {
			return fmt.Errorf("only endpoints 0 and 1 are valid, got endpoint group %d", endpoint.EndpointGroup)
		}

		var lagmode RackLinkLagMode
		err := lagmode.FromString(endpoint.Interface.LagMode)
		if err != nil {
			return fmt.Errorf("while parsing lag mode %q - %w", endpoint.Interface.LagMode, err)
		}

		var ipv4address netip.Prefix
		if endpoint.Interface.Ipv4Addr != "" {
			ipv4address, err = netip.ParsePrefix(endpoint.Interface.Ipv4Addr)
			if err != nil {
				return fmt.Errorf("while parsing ipv4 address %q - %w", endpoint.Interface.Ipv4Addr, err)
			}
		}

		var ipv6address netip.Prefix
		if endpoint.Interface.Ipv6Addr != "" {
			ipv6address, err = netip.ParsePrefix(endpoint.Interface.Ipv6Addr)
			if err != nil {
				return fmt.Errorf("while parsing ipv6 address %q - %w", endpoint.Interface.Ipv6Addr, err)
			}
		}

		o.Data.Endpoints[endpoint.EndpointGroup] = append(o.Data.Endpoints[endpoint.EndpointGroup], FreeformAggregateLinkMemberEndpoint{
			AggIntfId:     endpoint.Interface.Id,
			SystemId:      endpoint.System.Id,
			PortChannelId: endpoint.Interface.PortChannelId,
			LagMode:       lagmode,
			Ipv4Address:   ipv4address,
			Ipv6Address:   ipv6address,
		})
	}

	if len(o.Data.Endpoints[0]) == 0 || len(o.Data.Endpoints[1]) == 0 {
		return fmt.Errorf("each endpoint group must have atleast 1 member got %d and %d", len(o.Data.Endpoints[0]), len(o.Data.Endpoints[1]))
	}

	return nil
}

type FreeformEndpointGroup struct {
	Id   ObjectId `json:"id"`
	Name string   `json:"label"`
	Tags []string `json:"tags"`
}

func (o *FreeformClient) CreateAggregateLink(ctx context.Context, in *FreeformAggregateLinkData) (ObjectId, error) {
	var response objectIdResponse

	for _, epSlice := range in.Endpoints {
		for _, ep := range epSlice {
			if ep.AggIntfId != "" {
				return "", fmt.Errorf("aggregate interface id must be empty when creating an aggregate link")
			}
		}
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFfAggLinks, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *FreeformClient) GetAggregateLink(ctx context.Context, id ObjectId) (*FreeformAggregateLink, error) {
	var response FreeformAggregateLink

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfAggLinkById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) DeleteAggregateLink(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfAggLinkById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

var _ json.Marshaler = new(FreeformAggInterfaceData)

type FreeformAggInterfaceData struct {
	IfName        *string
	PortChannelId *int
	LagMode       *string
	Ipv4Address   *net.IPNet
	Ipv6Address   *net.IPNet
	Tags          []string
	EndpointGroup *int
}

func (o FreeformAggInterfaceData) MarshalJSON() ([]byte, error) {
	var raw struct {
		IfName        *string  `json:"if_name,omitempty"`
		PortChannelId *int     `json:"port_channel_id,omitempty"`
		Ipv4Addr      string   `json:"ipv4_addr,omitempty"`
		Ipv6Addr      string   `json:"ipv6_addr,omitempty"`
		Tags          []string `json:"tags"`
		EndpointGroup *int     `json:"endpoint_group,omitempty"`
	}

	raw.IfName = o.IfName
	raw.PortChannelId = o.PortChannelId
	if o.Ipv4Address != nil {
		raw.Ipv4Addr = o.Ipv4Address.String()
		if strings.Contains(raw.Ipv4Addr, "<nil>") {
			return nil, fmt.Errorf("cannot marshall ipv4 address - %s", raw.Ipv4Addr)
		}
	}
	if o.Ipv6Address != nil {
		raw.Ipv6Addr = o.Ipv6Address.String()
		if strings.Contains(raw.Ipv6Addr, "<nil>") {
			return nil, fmt.Errorf("cannot marshall ipv6 address - %s", raw.Ipv6Addr)
		}
	}
	raw.Tags = o.Tags
	raw.EndpointGroup = o.EndpointGroup
	return json.Marshal(&raw)
}

func (o *FreeformClient) GetAllAggregateLinks(ctx context.Context) ([]FreeformAggregateLink, error) {
	var response struct {
		Items []FreeformAggregateLink `json:"items"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfAggLinks, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) UpdateAggregateLink(ctx context.Context, id ObjectId, in *FreeformAggregateLinkData) error {
	for _, epSlice := range in.Endpoints {
		for _, ep := range epSlice {
			if ep.AggIntfId == "" {
				return fmt.Errorf("aggregate interface id must be non empty when updating an aggregate link")
			}
		}
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFfAggLinkById, o.blueprintId, id),
		apiInput: &in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
