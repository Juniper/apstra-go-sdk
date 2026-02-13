// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/str"
)

const (
	apiUrlBlueprintSecurityZones               = apiUrlBlueprintById + apiUrlPathDelim + "security-zones"
	apiUrlBlueprintSecurityZonesPrefix         = apiUrlBlueprintSecurityZones + apiUrlPathDelim
	apiUrlBlueprintSecurityZoneById            = apiUrlBlueprintSecurityZonesPrefix + "%s"
	apiUrlBlueprintSecurityZoneByIdDhcpServers = apiUrlBlueprintSecurityZoneById + apiUrlPathDelim + "dhcp-servers"
)

var (
	_ internal.IDer    = (*SecurityZone)(nil)
	_ json.Unmarshaler = (*SecurityZone)(nil)
)

type SecurityZone struct {
	Label             string                 `json:"label"`
	Description       *string                `json:"vrf_description,omitempty"`
	Type              enum.SecurityZoneType  `json:"sz_type"`
	VRFName           string                 `json:"vrf_name"`
	RoutingPolicyID   string                 `json:"routing_policy_id,omitempty"`   // automatically assigned if not set
	RouteTarget       *string                `json:"route_target,omitempty"`        // calculated only value
	RTPolicy          *RTPolicy              `json:"rt_policy,omitempty"`           // can be null
	VLAN              *VLAN                  `json:"vlan_id,omitempty"`             // can be null
	VNI               *int                   `json:"vni_id,omitempty"`              // can be null
	JunosEVPNIRBMode  *enum.JunosEVPNIRBMode `json:"junos_evpn_irb_mode,omitempty"` // can be null in POST, required in PUT AOS-58916
	AddressingSupport *enum.AddressingScheme `json:"addressing_support,omitempty"`  // Apstra 6.1+ only
	DisableIPv4       *bool                  `json:"disable_ipv4,omitempty"`        // Apstra 6.1+ only

	id string
}

func (o SecurityZone) ID() *string {
	if o.id == "" {
		return nil
	}
	return &o.id
}

func (o *SecurityZone) SetID(id string) {
	o.id = id
}

func (o *SecurityZone) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID                string                 `json:"id"`
		Label             string                 `json:"label"`
		Description       *string                `json:"vrf_description"`
		Type              enum.SecurityZoneType  `json:"sz_type"`
		VRFName           string                 `json:"vrf_name"`
		RoutingPolicyID   string                 `json:"routing_policy_id"`
		RouteTarget       *string                `json:"route_target"`
		RTPolicy          *RTPolicy              `json:"rt_policy"`
		VLAN              *VLAN                  `json:"vlan_id"`
		VNI               *int                   `json:"vni_id"`
		JunosEVPNIRBMode  *enum.JunosEVPNIRBMode `json:"junos_evpn_irb_mode"`
		AddressingSupport *enum.AddressingScheme `json:"addressing_support"`
		DisableIPv4       *bool                  `json:"disable_ipv4"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.id = raw.ID
	o.Label = raw.Label
	o.Description = raw.Description
	o.Type = raw.Type
	o.VRFName = raw.VRFName
	o.RoutingPolicyID = raw.RoutingPolicyID
	o.RouteTarget = raw.RouteTarget
	o.RTPolicy = raw.RTPolicy
	o.VLAN = raw.VLAN
	o.VNI = raw.VNI
	o.JunosEVPNIRBMode = raw.JunosEVPNIRBMode
	o.AddressingSupport = raw.AddressingSupport
	o.DisableIPv4 = raw.DisableIPv4

	return nil
}

// CreateSecurityZone creates an Apstra Routing Zone / Security Zone / VRF.
// If cfg.JunosEVPNIRBMode is omitted, but the API's version-dependent behavior
// requires that field, it will be set to JunosEVPNIRBModeAsymmetric in the
// request sent to the API.
func (o TwoStageL3ClosClient) CreateSecurityZone(ctx context.Context, cfg SecurityZone) (string, error) {
	szAddressingSupportOK := compatibility.SecurityZoneAddressingSupported.Check(o.client.apiVersion)

	if (cfg.AddressingSupport != nil || cfg.DisableIPv4 != nil) && !szAddressingSupportOK {
		return "", fmt.Errorf("AddressingSupport and DisableIPv4 must be nil with Apstra %s", o.client.apiVersion)
	}

	if cfg.AddressingSupport != nil &&
		*cfg.AddressingSupport != enum.AddressingSchemeIPv6 &&
		cfg.DisableIPv4 != nil &&
		*cfg.DisableIPv4 {
		return "", fmt.Errorf("disabling IPv4 not permitted with addressing scheme %s", *cfg.AddressingSupport)
	}

	response := &objectIdResponse{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSecurityZones, o.blueprintId),
		apiInput:    cfg,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return string(response.Id), nil
}

// GetSecurityZone fetches the Security Zone / Routing Zone / VRF with the given id
func (o TwoStageL3ClosClient) GetSecurityZone(ctx context.Context, id string) (SecurityZone, error) {
	var response SecurityZone
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSecurityZoneById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return SecurityZone{}, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

// GetSecurityZoneByLabel fetches the Security Zone / Routing Zone / VRF with
// the given label.
func (o TwoStageL3ClosClient) GetSecurityZoneByLabel(ctx context.Context, label string) (SecurityZone, error) {
	zones, err := o.GetSecurityZones(ctx)
	if err != nil {
		return SecurityZone{}, err
	}

	for _, zone := range zones {
		if zone.Label == label {
			return zone, nil
		}
	}

	return SecurityZone{}, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("security zone with label %q not found in blueprint %q", label, o.blueprintId),
	}
}

// GetSecurityZoneByVRFName fetches the Security Zone / Routing Zone / VRF with
// the given label.
func (o TwoStageL3ClosClient) GetSecurityZoneByVRFName(ctx context.Context, vrfName string) (SecurityZone, error) {
	zones, err := o.GetSecurityZones(ctx)
	if err != nil {
		return SecurityZone{}, err
	}

	for _, zone := range zones {
		if zone.VRFName == vrfName {
			return zone, nil
		}
	}

	return SecurityZone{}, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("security zone with vrf name %q not found in blueprint %q", vrfName, o.blueprintId),
	}
}

// GetSecurityZones returns []SecurityZone representing all Security Zones /
// Routing Zones / VRFs on the system.
func (o TwoStageL3ClosClient) GetSecurityZones(ctx context.Context) ([]SecurityZone, error) {
	response := &struct {
		Items map[string]SecurityZone `json:"items"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSecurityZones, o.blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}

	// This API endpoint returns a map. Convert to list for consistency with other 'GetAll' functions.
	result := make([]SecurityZone, len(response.Items))
	var i int
	for _, v := range response.Items {
		result[i] = v
		i++
	}

	return result, nil
}

// UpdateSecurityZone replaces the configuration of zone zoneId with the supplied CreateSecurityZoneCfg
func (o TwoStageL3ClosClient) UpdateSecurityZone(ctx context.Context, v SecurityZone) error {
	if v.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	szAddressingSupportOK := compatibility.SecurityZoneAddressingSupported.Check(o.client.apiVersion)

	if (v.AddressingSupport != nil || v.DisableIPv4 != nil) && !szAddressingSupportOK {
		return fmt.Errorf("AddressingSupport and DisableIPv4 must be nil with Apstra %s", o.client.apiVersion)
	}

	if v.AddressingSupport != nil &&
		*v.AddressingSupport != enum.AddressingSchemeIPv6 &&
		v.DisableIPv4 != nil &&
		*v.DisableIPv4 {
		return fmt.Errorf("disabling IPv4 not permitted with addressing scheme %s", *v.AddressingSupport)
	}

	// workaround for error:
	//
	// {
	//  "error_code": 422,
	//  "errors": {
	//    "disable_ipv4": "IPv4 support can only be disabled when addressing_support=\"ipv6\""
	//  }
	// }
	//
	// JP says the API behavior is deliberate: A shim layer sets an omitted `disable_ipv4` attribute
	// to the current value in PUT requests /even when doing so produces an unsupported combination/.
	//
	// Because of this API behavior, when setting addressing_support to something other than IPv6,
	// we will explicitly enable IPv4 support (disable = false)
	if v.AddressingSupport != nil &&
		*v.AddressingSupport != enum.AddressingSchemeIPv6 &&
		v.DisableIPv4 == nil &&
		szAddressingSupportOK {
		v.DisableIPv4 = pointer.To(false)
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintSecurityZoneById, o.blueprintId, v.id),
		apiInput: v,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o TwoStageL3ClosClient) DeleteSecurityZone(ctx context.Context, id string) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintSecurityZoneById, o.blueprintId, id),
	})
	return convertTtaeToAceWherePossible(err)
}
