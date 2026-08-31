// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	aerrors "github.com/Juniper/apstra-go-sdk/errors"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/internal/str"
	"github.com/Juniper/apstra-go-sdk/internal/urls"
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

// ID returns a pointer to a copy of the object's ID, or nil when no ID is set.
func (e EVPNInterconnectGroup) ID() *string {
	if e.id == "" {
		return nil
	}
	id := e.id
	return &id
}

func (e *EVPNInterconnectGroup) SetID(id string) error {
	if e.id != "" {
		return aerrors.IDAlreadySet(fmt.Sprintf("id already has value %q", e.id))
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

// parseErr handles error text like this:
//
//	{
//	  "errors": {
//	    "interconnect_security_zones": {
//	      "ajksdfalkj": "EVPN Interconnect routing zone node ID \"ajksdfalkj\" does not exist",
//	      "newxSuZde_5WSjOO_A": {
//	        "routing_policy_id": "EVPN Interconnect routing policy \"ioeoruwpq\" does not exist"
//	      }
//	    }
//	  }
//	}
//
// It was written specifically to handle errors provoked by PUT or PATCH to the evpn_interconnect_groups
// API endpoint, and specifically errors resulting from use of an invalid routing zone ID.
func (e *EVPNInterconnectGroup) parseErr(err error) error {
	var ttae TalkToApstraErr // We only handle TalkToApstraErr errors.
	if !errors.As(err, &ttae) {
		return err
	}

	// We only handle errors from the evpn_interconnect_groups API endpoint.
	if !urls.DatacenterEvpnInterconnectGroupByIDRegex.MatchString(ttae.Request.URL.Path) {
		return convertTtaeToAceWherePossible(err)
	}

	var target struct {
		Errors struct {
			InterconnectSecurityZones map[string]jsontext.Value `json:"interconnect_security_zones"`
			Extra                     map[string]jsontext.Value `json:",embed"`
		} `json:"errors"`
	}

	if fail := json.Unmarshal([]byte(ttae.Msg), &target); fail != nil {
		return convertTtaeToAceWherePossible(err)
	}

	var result error

	for zone, val := range target.Errors.InterconnectSecurityZones {
		// Missing Routing Zone is communicated via a string value.
		noSuchRZMsg, _ := json.Marshal(fmt.Sprintf("EVPN Interconnect routing zone node ID %q does not exist", zone))
		if bytes.Equal(noSuchRZMsg, val) {
			newErr := ClientErr{
				errType: ErrNotfound,
				err:     errors.New(val.String()),
				detail:  NodeTypeSecurityZone,
			}
			if result == nil {
				result = newErr
			} else {
				result = errors.Join(result, newErr)
			}
			continue
		}

		// Missing Routing Policy is communicated via a struct.
		var target2 struct {
			RoutingPolicyID *string `json:"routing_policy_id"`
		}
		if fail := json.Unmarshal(val, &target2); fail != nil {
			// We have failed to unmarshal the error message, so we don't know what it is. Wrap it in an UnhandledApstraErr and continue.
			newErr := aerrors.UnhandledApstraErr(fmt.Sprintf("interconnect_security_zones: %q: %q", zone, val.String()))
			if result == nil {
				result = newErr
			} else {
				result = errors.Join(result, aerrors.UnhandledApstraErr(fmt.Sprintf("interconnect_security_zones: %q: %q", zone, val.String())))
			}
			continue
		}

		if target2.RoutingPolicyID != nil {
			// We have found a routing policy ID error.
			if strings.HasSuffix(*target2.RoutingPolicyID, " does not exist") {
				newErr := ClientErr{
					errType: ErrNotfound,
					err:     errors.New(val.String()),
					detail:  NodeTypeRoutingPolicy,
				}
				if result == nil {
					result = newErr
				} else {
					result = errors.Join(result, newErr)
				}
				continue
			}

			result = errors.Join(result, aerrors.UnhandledApstraErr(fmt.Sprintf("interconnect_security_zones: %q: %q", zone, val.String())))
		}
	}

	// Handle any other errors that may have been returned by the API.
	// We don't know what they are, so we just wrap them in an UnhandledApstraErr.
	for k, v := range target.Errors.Extra {
		result = errors.Join(result, aerrors.UnhandledApstraErr(fmt.Sprintf("%s: %s", k, v.String())))
	}

	return result
}

func (o *TwoStageL3ClosClient) CreateEVPNInterconnectGroup(ctx context.Context, in EVPNInterconnectGroup) (string, error) {
	var response struct {
		ID string `json:"id"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(urls.DatacenterEvpnInterconnectGroups, o.Id()),
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
		urlStr:      fmt.Sprintf(urls.DatacenterEvpnInterconnectGroupByID, o.Id(), id),
		apiResponse: &response,
	})
	if err != nil {
		return EVPNInterconnectGroup{}, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *TwoStageL3ClosClient) GetEVPNInterconnectGroups(ctx context.Context) ([]EVPNInterconnectGroup, error) {
	var response struct {
		Items []EVPNInterconnectGroup `json:"evpn_interconnect_groups"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(urls.DatacenterEvpnInterconnectGroups, o.Id()),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *TwoStageL3ClosClient) GetEVPNInterconnectGroupByLabel(ctx context.Context, name string) (EVPNInterconnectGroup, error) {
	items, err := o.GetEVPNInterconnectGroups(ctx)
	if err != nil {
		return EVPNInterconnectGroup{}, fmt.Errorf("GetEVPNInterconnectGroups: %w", err)
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
		urlStr:   fmt.Sprintf(urls.DatacenterEvpnInterconnectGroupByID, o.Id(), *v.ID()),
		apiInput: &v,
	})
	if err != nil {
		return v.parseErr(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) DeleteEVPNInterconnectGroup(ctx context.Context, id string) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(urls.DatacenterEvpnInterconnectGroupByID, o.Id(), id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
