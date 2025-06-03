// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

const (
	apiUrlBlueprintEvpnInterconnectGroups       = apiUrlBlueprintById + apiUrlPathDelim + "evpn_interconnect_groups"
	apiUrlBlueprintEvpnInterconnectGroupsPrefix = apiUrlBlueprintEvpnInterconnectGroups + apiUrlPathDelim
	apiUrlBlueprintEvpnInterconnectGroupById    = apiUrlBlueprintEvpnInterconnectGroupsPrefix + "%s"
)

var _ json.Unmarshaler = (*EvpnInterconnectGroup)(nil)

type EvpnInterconnectGroup struct {
	Id   ObjectId
	Data *EvpnInterconnectGroupData
}

func (o *EvpnInterconnectGroup) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id                        ObjectId                                  `json:"id"`
		Label                     string                                    `json:"label"`
		RouteTarget               string                                    `json:"interconnect_route_target"`
		EsiMac                    *string                                   `json:"interconnect_esi_mac"`
		InterconnectSecurityZones map[ObjectId]InterconnectSecurityZoneData `json:"interconnect_security_zones"`
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(EvpnInterconnectGroupData)
	o.Data.Label = raw.Label
	o.Data.RouteTarget = raw.RouteTarget
	if raw.EsiMac != nil {
		esiMac, err := net.ParseMAC(*raw.EsiMac)
		if err != nil {
			return fmt.Errorf("parsing interconnect esi mac: %s", err)
		}
		o.Data.EsiMac = esiMac
	}
	o.Data.InterconnectSecurityZones = raw.InterconnectSecurityZones

	return nil
}

type InterconnectSecurityZoneData struct {
	RoutingPolicyId *ObjectId `json:"routing_policy_id"`
	RouteTarget     *string   `json:"interconnect_route_target"`
	L3Enabled       bool      `json:"enabled_for_l3"`
}

var _ json.Marshaler = (*EvpnInterconnectGroupData)(nil)

type EvpnInterconnectGroupData struct {
	Label                     string
	RouteTarget               string
	EsiMac                    net.HardwareAddr
	InterconnectSecurityZones map[ObjectId]InterconnectSecurityZoneData
}

func (o EvpnInterconnectGroupData) MarshalJSON() ([]byte, error) {
	var raw struct {
		Label                     string                                    `json:"label"`
		RouteTarget               string                                    `json:"interconnect_route_target"`
		EsiMac                    string                                    `json:"interconnect_esi_mac,omitempty"`
		InterconnectSecurityZones map[ObjectId]InterconnectSecurityZoneData `json:"interconnect_security_zones"`
	}

	raw.Label = o.Label
	raw.RouteTarget = o.RouteTarget
	if o.EsiMac != nil {
		raw.EsiMac = o.EsiMac.String()
	}
	raw.InterconnectSecurityZones = o.InterconnectSecurityZones

	return json.Marshal(raw)
}

func (o *TwoStageL3ClosClient) CreateEvpnInterconnectGroup(ctx context.Context, in *EvpnInterconnectGroupData) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintEvpnInterconnectGroups, o.Id()),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *TwoStageL3ClosClient) GetEvpnInterconnectGroup(ctx context.Context, id ObjectId) (*EvpnInterconnectGroup, error) {
	var response EvpnInterconnectGroup

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintEvpnInterconnectGroupById, o.Id(), id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *TwoStageL3ClosClient) GetAllEvpnInterconnectGroups(ctx context.Context) ([]EvpnInterconnectGroup, error) {
	var response struct {
		Items []EvpnInterconnectGroup `json:"evpn_interconnect_groups"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintEvpnInterconnectGroups, o.Id()),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *TwoStageL3ClosClient) GetEvpnInterconnectGroupByName(ctx context.Context, name string) (*EvpnInterconnectGroup, error) {
	items, err := o.GetAllEvpnInterconnectGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllEvpnInterconnectGroups: %w", err)
	}

	var evpnInterconnectGroup *EvpnInterconnectGroup
	for _, item := range items {
		if item.Data.Label == name {
			if evpnInterconnectGroup == nil {
				evpnInterconnectGroup = &item
			} else {
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("found multiple EVPN Interconnect Groups with label %q", name),
				}
			}
		}
	}

	if evpnInterconnectGroup == nil {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("EVPN Interconnect Group with label %q not found", name),
		}
	}

	return evpnInterconnectGroup, nil
}

func (o *TwoStageL3ClosClient) UpdateEvpnInterconnectGroup(ctx context.Context, id ObjectId, cfg *EvpnInterconnectGroupData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintEvpnInterconnectGroupById, o.Id(), id),
		apiInput: cfg,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) DeleteEvpnInterconnectGroup(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintEvpnInterconnectGroupById, o.Id(), id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
