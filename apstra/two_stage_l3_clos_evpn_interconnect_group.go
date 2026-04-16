// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/internal/str"
)

const (
	apiUrlBlueprintEvpnInterconnectGroups       = apiUrlBlueprintById + apiUrlPathDelim + "evpn_interconnect_groups"
	apiUrlBlueprintEvpnInterconnectGroupsPrefix = apiUrlBlueprintEvpnInterconnectGroups + apiUrlPathDelim
	apiUrlBlueprintEvpnInterconnectGroupById    = apiUrlBlueprintEvpnInterconnectGroupsPrefix + "%s"
)

type InterconnectVirtualNetwork struct {
	L2Enabled      bool    `json:"l2"`
	L3Enabled      bool    `json:"l3"`
	TranslationVNI *uint32 `json:"translation_vni,omitempty"`
}

type InterconnectSecurityZone struct {
	L3Enabled       bool    `json:"enabled_for_l3"`
	RouteTarget     *string `json:"interconnect_route_target"`
	RoutingPolicyId *string `json:"routing_policy_id"`
}

var (
	_ internal.IDer    = (*EVPNInterconnectGroup)(nil)
	_ json.Marshaler   = (*EVPNInterconnectGroup)(nil)
	_ json.Unmarshaler = (*EVPNInterconnectGroup)(nil)
)

type EVPNInterconnectGroup struct {
	Label                       *string                               `json:"label,omitempty"`
	RouteTarget                 *string                               `json:"interconnect_route_target,omitempty"`
	ESIMAC                      net.HardwareAddr                      `json:"interconnect_esi_mac,omitempty"`
	InterconnectSecurityZones   map[string]InterconnectSecurityZone   `json:"interconnect_security_zones,omitempty"`
	InterconnectVirtualNetworks map[string]InterconnectVirtualNetwork `json:"interconnect_virtual_networks,omitempty"`

	id string
}

func (e EVPNInterconnectGroup) ID() *string {
	if e.id == "" {
		return nil
	}
	return &e.id
}

func (e *EVPNInterconnectGroup) SetID(id string) error {
	if e.id != "" {
		return sdk.ErrIDAlreadySet(fmt.Sprintf("id already has value %q", e.id))
	}

	e.id = id
	return nil
}

func (e EVPNInterconnectGroup) MarshalJSON() ([]byte, error) {
	type Alias EVPNInterconnectGroup

	return json.Marshal(&struct {
		EsiMac string `json:"interconnect_esi_mac,omitempty"`
		*Alias
	}{
		EsiMac: func() string {
			if e.ESIMAC == nil {
				return ""
			}
			return e.ESIMAC.String()
		}(),
		Alias: (*Alias)(&e),
	})
}

func (e *EVPNInterconnectGroup) UnmarshalJSON(bytes []byte) error {
	type Alias EVPNInterconnectGroup

	aux := &struct {
		ID     string `json:"id"`
		ESIMAC string `json:"interconnect_esi_mac"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(bytes, aux); err != nil {
		return err
	}

	e.id = aux.ID

	if aux.ESIMAC != "" {
		hw, err := net.ParseMAC(aux.ESIMAC)
		if err != nil {
			return err
		}
		e.ESIMAC = hw
	} else {
		e.ESIMAC = nil
	}

	return nil
}

func (o *TwoStageL3ClosClient) CreateEVPNInterconnectGroup(ctx context.Context, in EVPNInterconnectGroup) (string, error) {
	var response struct {
		ID string `json:"id"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintEvpnInterconnectGroups, o.Id()),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.ID, nil
}

func (o *TwoStageL3ClosClient) GetEVPNInterconnectGroup(ctx context.Context, id string) (EVPNInterconnectGroup, error) {
	var response EVPNInterconnectGroup

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintEvpnInterconnectGroupById, o.Id(), id),
		apiResponse: &response,
	})
	if err != nil {
		return EVPNInterconnectGroup{}, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *TwoStageL3ClosClient) GetAllEVPNInterconnectGroups(ctx context.Context) ([]EVPNInterconnectGroup, error) {
	var response struct {
		Items []EVPNInterconnectGroup `json:"evpn_interconnect_groups"`
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

func (o *TwoStageL3ClosClient) GetEVPNInterconnectGroupByName(ctx context.Context, name string) (EVPNInterconnectGroup, error) {
	items, err := o.GetAllEVPNInterconnectGroups(ctx)
	if err != nil {
		return EVPNInterconnectGroup{}, fmt.Errorf("GetAllEVPNInterconnectGroups: %w", err)
	}

	var evpnInterconnectGroup *EVPNInterconnectGroup
	for _, item := range items {
		if item.Label != nil && *item.Label == name {
			if evpnInterconnectGroup == nil {
				evpnInterconnectGroup = &item
			} else {
				return EVPNInterconnectGroup{}, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("found multiple EVPN Interconnect Groups with label %q", name),
				}
			}
		}
	}

	if evpnInterconnectGroup == nil {
		return EVPNInterconnectGroup{}, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("EVPN Interconnect Group with label %q not found", name),
		}
	}

	return *evpnInterconnectGroup, nil
}

func (o *TwoStageL3ClosClient) UpdateEVPNInterconnectGroup(ctx context.Context, v EVPNInterconnectGroup) error {
	if v.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintEvpnInterconnectGroupById, o.Id(), *v.ID()),
		apiInput: &v,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) DeleteEVPNInterconnectGroup(ctx context.Context, id string) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintEvpnInterconnectGroupById, o.Id(), id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
