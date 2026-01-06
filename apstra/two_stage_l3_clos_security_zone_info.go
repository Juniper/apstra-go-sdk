// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"

	"github.com/Juniper/apstra-go-sdk/enum"
)

const (
	apiUrlWebExpSecurityZones = apiUrlBlueprintByIdPrefix + "experience/web/security-zones"
)

var _ json.Unmarshaler = (*TwoStageL3ClosSecurityZoneInfoInterface)(nil)

type TwoStageL3ClosSecurityZoneInfoInterface struct {
	Id       ObjectId
	Ipv4Addr *netip.Prefix
	Ipv6Addr *netip.Prefix
}

func (o *TwoStageL3ClosSecurityZoneInfoInterface) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id       ObjectId `json:"id"`
		Ipv4Addr *string  `json:"ipv4_addr"`
		Ipv6Addr *string  `json:"ipv6_addr"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id

	if raw.Ipv4Addr != nil {
		ipv4Addr, err := netip.ParsePrefix(*raw.Ipv4Addr)
		if err != nil {
			return fmt.Errorf("failed parsing ipv4 address %q - %w", *raw.Ipv4Addr, err)
		}
		if !ipv4Addr.Addr().Is4() {
			return fmt.Errorf("invalid ipv4 address %q", *raw.Ipv4Addr)
		}
		o.Ipv4Addr = &ipv4Addr
	} else {
		o.Ipv4Addr = nil
	}

	if raw.Ipv6Addr != nil {
		ipv6Addr, err := netip.ParsePrefix(*raw.Ipv6Addr)
		if err != nil {
			return fmt.Errorf("failed parsing ipv6 address %q - %w", *raw.Ipv6Addr, err)
		}
		if !ipv6Addr.Addr().Is6() {
			return fmt.Errorf("invalid ipv6 address %q", *raw.Ipv6Addr)
		}
		o.Ipv6Addr = &ipv6Addr
	} else {
		o.Ipv6Addr = nil
	}

	return nil
}

type TwoStageL3ClosSecurityZoneInfo struct {
	Id               ObjectId               `json:"id"`
	Label            string                 `json:"label"`
	VrfName          string                 `json:"vrf_name"`
	SecurityZoneType *enum.SecurityZoneType `json:"sz_type"`
	VlanId           *VLAN                  `json:"vlan_id"`
	RoutingPolicyId  *ObjectId              `json:"routing_policy_id"`
	JunosEvpnIrbMode *enum.JunosEVPNIRBMode `json:"junos_evpn_irb_mode"`
	RouteTarget      *string                `json:"route_target"`
	VniId            *int                   `json:"vni_id"`
	RtPolicy         *RTPolicy              `json:"rt_policy"`
	Tags             []string               `json:"tags"`
	MemberInterfaces []struct {
		Loopbacks     []TwoStageL3ClosSecurityZoneInfoInterface `json:"loopbacks"`
		Subinterfaces []TwoStageL3ClosSecurityZoneInfoInterface `json:"subinterfaces"`
		HostingSystem struct {
			Id       ObjectId `json:"id"`
			Label    *string  `json:"label"`
			Hostname *string  `json:"hostname"`
		} `json:"hosting_system"`
	} `json:"member_interfaces"`
}

// GetAllSecurityZoneInfos returns map[ObjectId]TwoStageL3CloSecurityZoneInfo
// keyed by Security Zone ID.
func (o *TwoStageL3ClosClient) GetAllSecurityZoneInfos(ctx context.Context) (map[ObjectId]TwoStageL3ClosSecurityZoneInfo, error) {
	var apiResponse struct {
		SecurityZones []TwoStageL3ClosSecurityZoneInfo `json:"security_zones"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlWebExpSecurityZones, o.blueprintId),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make(map[ObjectId]TwoStageL3ClosSecurityZoneInfo, len(apiResponse.SecurityZones))
	for _, zone := range apiResponse.SecurityZones {
		result[zone.Id] = zone
	}

	return result, nil
}

// GetSecurityZoneInfo returns TwoStageL3ClosSecurityZoneInfo describing
// the Security Zone represented by ID.
func (o *TwoStageL3ClosClient) GetSecurityZoneInfo(ctx context.Context, id ObjectId) (*TwoStageL3ClosSecurityZoneInfo, error) {
	allSZs, err := o.GetAllSecurityZoneInfos(ctx)
	if err != nil {
		return nil, err
	}

	var result *TwoStageL3ClosSecurityZoneInfo
	for _, v := range allSZs {
		if v.Id == id {
			result = &v
			break
		}
	}

	if result == nil {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("security zone with id %v not found", id),
		}
	}

	return result, nil
}
