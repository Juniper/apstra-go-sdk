// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	apiUrlBlueprintRemoteGateways       = apiUrlBlueprintById + apiUrlPathDelim + "remote_gateways"
	apiUrlBlueprintRemoteGatewaysPrefix = apiUrlBlueprintRemoteGateways + apiUrlPathDelim
	apiUrlBlueprintRemoteGatewayById    = apiUrlBlueprintRemoteGatewaysPrefix + "%s"
)

type rawRemoteGatewayRequest struct {
	RouteTypes     string     `json:"evpn_route_types"`
	LocalGwNodes   []ObjectId `json:"local_gw_nodes"`
	GwAsn          uint32     `json:"gw_asn"`
	GwIp           string     `json:"gw_ip"`
	GwName         string     `json:"gw_name"`
	Ttl            *uint8     `json:"ttl,omitempty"`
	KeepaliveTimer *uint16    `json:"keepalive_timer,omitempty"`
	HoldtimeTimer  *uint16    `json:"holdtime_timer,omitempty"`
	Password       *string    `json:"password"`
}

type rawRemoteGatewayResponse struct {
	Id           ObjectId `json:"id"`
	RouteTypes   string   `json:"evpn_route_types"`
	LocalGwNodes []struct {
		NodeId ObjectId `json:"node_id"`
	} `json:"local_gw_nodes"`
	GwAsn          uint32  `json:"gw_asn"`
	GwIp           string  `json:"gw_ip"`
	GwName         string  `json:"gw_name"`
	Ttl            *uint8  `json:"ttl"`
	KeepaliveTimer *uint16 `json:"keepalive_timer"`
	HoldtimeTimer  *uint16 `json:"holdtime_timer"`
}

func (o *rawRemoteGatewayResponse) polish() (*RemoteGateway, error) {
	routeTypes := enum.RemoteGatewayRouteTypes.Parse(o.RouteTypes)
	if routeTypes == nil {
		return nil, fmt.Errorf("failed parsing remote gateway route types: %q", o.RouteTypes)
	}

	localGwNodes := make([]ObjectId, len(o.LocalGwNodes))
	for i, localGwNode := range o.LocalGwNodes {
		localGwNodes[i] = localGwNode.NodeId
	}

	gwIp := net.ParseIP(o.GwIp)
	if gwIp == nil {
		return nil, fmt.Errorf("faileed parsing remote gateway IP: %q", o.GwIp)
	}

	return &RemoteGateway{
		Id: o.Id,
		Data: &RemoteGatewayData{
			RouteTypes:     *routeTypes,
			LocalGwNodes:   localGwNodes,
			GwAsn:          o.GwAsn,
			GwIp:           gwIp,
			GwName:         o.GwName,
			Ttl:            o.Ttl,
			KeepaliveTimer: o.KeepaliveTimer,
			HoldtimeTimer:  o.HoldtimeTimer,
		},
	}, nil
}

type RemoteGatewayData struct {
	RouteTypes     enum.RemoteGatewayRouteType
	LocalGwNodes   []ObjectId
	GwAsn          uint32
	GwIp           net.IP
	GwName         string
	Ttl            *uint8
	KeepaliveTimer *uint16
	HoldtimeTimer  *uint16
	Password       *string
}

type RemoteGateway struct {
	Id   ObjectId
	Data *RemoteGatewayData
}

func (o *RemoteGatewayData) raw() *rawRemoteGatewayRequest {
	return &rawRemoteGatewayRequest{
		RouteTypes:     o.RouteTypes.Value,
		LocalGwNodes:   o.LocalGwNodes,
		GwAsn:          o.GwAsn,
		GwIp:           o.GwIp.String(),
		GwName:         o.GwName,
		Ttl:            o.Ttl,
		KeepaliveTimer: o.KeepaliveTimer,
		HoldtimeTimer:  o.HoldtimeTimer,
		Password:       o.Password,
	}
}

func (o *TwoStageL3ClosClient) createRemoteGateway(ctx context.Context, in *rawRemoteGatewayRequest) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, talkToApstraIn{
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

func (o *TwoStageL3ClosClient) getRemoteGateway(ctx context.Context, id ObjectId) (*rawRemoteGatewayResponse, error) {
	var response rawRemoteGatewayResponse

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintRemoteGatewayById, o.Id(), id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *TwoStageL3ClosClient) getAllRemoteGateways(ctx context.Context) ([]rawRemoteGatewayResponse, error) {
	var response struct {
		RemoteGateways []rawRemoteGatewayResponse `json:"remote_gateways"`
	}

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintRemoteGateways, o.Id()),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.RemoteGateways, nil
}

func (o *TwoStageL3ClosClient) getRemoteGatewayByName(ctx context.Context, name string) (*rawRemoteGatewayResponse, error) {
	rawRemoteGateways, err := o.getAllRemoteGateways(ctx)
	if err != nil {
		return nil, err
	}

	var result rawRemoteGatewayResponse
	var found bool

	for _, rawRemoteGateway := range rawRemoteGateways {
		if rawRemoteGateway.GwName == name {
			if found {
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("found multiple remote gateways named %q found", name),
				}
			}
			result = rawRemoteGateway
			found = true
		}
	}

	if found {
		return &result, nil
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no remote gateway named %q found", name),
	}
}

func (o *TwoStageL3ClosClient) updateRemoteGateway(ctx context.Context, id ObjectId, in *rawRemoteGatewayRequest) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintRemoteGatewayById, o.Id(), id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) deleteRemoteGateway(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintRemoteGatewayById, o.Id(), id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
