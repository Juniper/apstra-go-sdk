// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	apiUrlBlueprintRemoteGateways       = apiUrlBlueprintById + apiUrlPathDelim + "remote_gateways"
	apiUrlBlueprintRemoteGatewaysPrefix = apiUrlBlueprintRemoteGateways + apiUrlPathDelim
	apiUrlBlueprintRemoteGatewayById    = apiUrlBlueprintRemoteGatewaysPrefix + "%s"
)

var _ json.Unmarshaler = (*TwoStageL3ClosRemoteGateway)(nil)

type TwoStageL3ClosRemoteGateway struct {
	Id   ObjectId
	Data *TwoStageL3ClosRemoteGatewayData
}

func (o *TwoStageL3ClosRemoteGateway) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id                      ObjectId                     `json:"id"`
		Label                   string                       `json:"gw_name"`
		GwIp                    string                       `json:"gw_ip"`
		GwAsn                   uint32                       `json:"gw_asn"`
		RouteTypes              *enum.RemoteGatewayRouteType `json:"evpn_route_types"`
		Ttl                     *uint8                       `json:"ttl"`
		KeepaliveTimer          *uint16                      `json:"keepalive_timer"`
		HoldtimeTimer           *uint16                      `json:"holdtime_timer"`
		EvpnInterconnectGroupId *ObjectId                    `json:"evpn_interconnect_group_id"`
		LocalGwNodes            []struct {
			NodeId ObjectId `json:"node_id"`
			// Label              string        `json:"label"`
			// Role               enum.NodeRole `json:"role"`
			// EvpnInternalRd     interface{}   `json:"evpn_internal_rd"`
			// EvpnInterconnectRd interface{}   `json:"evpn_interconnect_rd"`
		} `json:"local_gw_nodes"`
	}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(TwoStageL3ClosRemoteGatewayData)
	o.Data.Label = raw.Label
	gwIp, err := netip.ParseAddr(raw.GwIp)
	if err != nil {
		return fmt.Errorf("parse gw_ip address: %q", err)
	}
	o.Data.GwIp = gwIp
	o.Data.GwAsn = raw.GwAsn
	o.Data.RouteTypes = raw.RouteTypes
	o.Data.Ttl = raw.Ttl
	o.Data.KeepaliveTimer = raw.KeepaliveTimer
	o.Data.HoldtimeTimer = raw.HoldtimeTimer
	o.Data.EvpnInterconnectGroupId = raw.EvpnInterconnectGroupId
	o.Data.LocalGwNodes = make([]ObjectId, len(raw.LocalGwNodes))
	for i, localGwNode := range raw.LocalGwNodes {
		o.Data.LocalGwNodes[i] = localGwNode.NodeId
	}

	return nil
}

var _ json.Marshaler = (*TwoStageL3ClosRemoteGatewayData)(nil)

type TwoStageL3ClosRemoteGatewayData struct {
	Label                   string
	GwIp                    netip.Addr
	GwAsn                   uint32
	RouteTypes              *enum.RemoteGatewayRouteType
	Ttl                     *uint8
	KeepaliveTimer          *uint16
	HoldtimeTimer           *uint16
	Password                *string
	EvpnInterconnectGroupId *ObjectId
	LocalGwNodes            []ObjectId
}

func (o TwoStageL3ClosRemoteGatewayData) MarshalJSON() ([]byte, error) {
	var raw struct {
		Label                   string                       `json:"gw_name"`
		GwIp                    string                       `json:"gw_ip"`
		GwAsn                   uint32                       `json:"gw_asn"`
		LocalGwNodes            []ObjectId                   `json:"local_gw_nodes"`
		RouteTypes              *enum.RemoteGatewayRouteType `json:"evpn_route_types,omitempty"`
		Ttl                     *uint8                       `json:"ttl,omitempty"`
		KeepaliveTimer          *uint16                      `json:"keepalive_timer,omitempty"`
		HoldtimeTimer           *uint16                      `json:"holdtime_timer,omitempty"`
		Password                *string                      `json:"password"`
		EvpnInterconnectGroupId *ObjectId                    `json:"evpn_interconnect_group_id"`
	}

	raw.Label = o.Label
	raw.GwIp = o.GwIp.String()
	raw.GwAsn = o.GwAsn
	raw.LocalGwNodes = o.LocalGwNodes
	raw.RouteTypes = o.RouteTypes
	raw.Ttl = o.Ttl
	raw.KeepaliveTimer = o.KeepaliveTimer
	raw.HoldtimeTimer = o.HoldtimeTimer
	raw.Password = o.Password
	raw.EvpnInterconnectGroupId = o.EvpnInterconnectGroupId

	return json.Marshal(&raw)
}

func (o *TwoStageL3ClosClient) CreateRemoteGateway(ctx context.Context, in *TwoStageL3ClosRemoteGatewayData) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintRemoteGateways, o.Id()),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *TwoStageL3ClosClient) GetRemoteGateway(ctx context.Context, id ObjectId) (*TwoStageL3ClosRemoteGateway, error) {
	var response TwoStageL3ClosRemoteGateway

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintRemoteGatewayById, o.Id(), id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *TwoStageL3ClosClient) GetAllRemoteGateways(ctx context.Context) ([]TwoStageL3ClosRemoteGateway, error) {
	var response struct {
		RemoteGateways []TwoStageL3ClosRemoteGateway `json:"remote_gateways"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintRemoteGateways, o.Id()),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.RemoteGateways, nil
}

func (o *TwoStageL3ClosClient) GetRemoteGatewayByName(ctx context.Context, name string) (*TwoStageL3ClosRemoteGateway, error) {
	all, err := o.GetAllRemoteGateways(ctx)
	if err != nil {
		return nil, err
	}

	var result *TwoStageL3ClosRemoteGateway

	for _, each := range all {
		if each.Data.Label == name {
			if result != nil {
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("found multiple remote gateways with name %q", name),
				}
			}
			result = &each
		}
	}

	if result != nil {
		return result, nil
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no remote gateway with name %q", name),
	}
}

func (o *TwoStageL3ClosClient) UpdateRemoteGateway(ctx context.Context, id ObjectId, in *TwoStageL3ClosRemoteGatewayData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintRemoteGatewayById, o.Id(), id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) DeleteRemoteGateway(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintRemoteGatewayById, o.Id(), id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
