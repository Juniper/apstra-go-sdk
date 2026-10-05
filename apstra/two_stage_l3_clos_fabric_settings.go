// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/enum"
	oenum "github.com/orsinium-labs/enum"
)

const (
	apiUrlBlueprintFabricSettings = apiUrlBlueprintById + apiUrlPathDelim + "fabric-settings"
)

type junosEvpnRoutingInstanceType oenum.Member[string]

var (
	junosEvpnRoutingInstanceTypeVlanAware = junosEvpnRoutingInstanceType{Value: "vlan_aware"}
	junosEvpnRoutingInstanceTypeDefault   = junosEvpnRoutingInstanceType{Value: "default"}
	junosEvpnRoutingInstanceTypes         = oenum.New(junosEvpnRoutingInstanceTypeVlanAware, junosEvpnRoutingInstanceTypeDefault)
)

var (
	_ json.Marshaler   = (*FabricSettings)(nil)
	_ json.Unmarshaler = (*FabricSettings)(nil)
)

type FabricSettings struct { //										 4.2.0                          6.1.0
	AntiAffinityPolicy                    *AntiAffinityPolicy     // anti-affinity-policy           fabric_policy node
	DefaultAnycastGWMAC                   net.HardwareAddr        // n/a                            fabric_policy node
	DefaultSviL3Mtu                       *uint16                 // virtual_network_policy node    fabric_policy node
	EsiMacMsb                             *uint8                  // /fabric-addressing-policy      fabric_policy node
	EvpnGenerateType5HostRoutes           *enum.FeatureSwitch     // virtual_network_policy node    fabric_policy node
	ExternalRouterMtu                     *uint16                 // virtual_network_policy node    fabric_policy node
	FabricL3Mtu                           *uint16                 // /fabric-addressing-policy      fabric_policy node
	Ipv6Enabled                           *bool                   // /fabric-addressing-policy      n/a
	JunosEvpnDuplicateMacRecoveryTime     *uint16                 // virtual_network_policy node    fabric_policy node
	JunosEvpnMaxNexthopAndInterfaceNumber *enum.FeatureSwitch     // virtual_network_policy node    fabric_policy node
	JunosEvpnRoutingInstanceVlanAware     *enum.FeatureSwitch     // virtual_network_policy node    fabric_policy node
	JunosExOverlayEcmp                    *enum.FeatureSwitch     // virtual_network_policy node    fabric_policy node
	JunosGracefulRestart                  *enum.FeatureSwitch     // virtual_network_policy node    fabric_policy node
	MaxEvpnRoutes                         *uint32                 // virtual_network_policy node    fabric_policy node
	MaxExternalRoutes                     *uint32                 // virtual_network_policy node    fabric_policy node
	MaxFabricRoutes                       *uint32                 // virtual_network_policy node    fabric_policy node
	MaxMlagRoutes                         *uint32                 // virtual_network_policy node    fabric_policy node
	OptimiseSzFootprint                   *enum.FeatureSwitch     // security_zone_policy node      fabric_policy node
	OverlayControlProtocol                *OverlayControlProtocol // virtual_network_policy node    fabric_policy node
	SpineLeafLinks                        *AddressingScheme       // blueprint creation only        n/a
	SpineSuperspineLinks                  *AddressingScheme       // blueprint creation only        n/a
}

func (o FabricSettings) MarshalJSON() ([]byte, error) {
	var antiAffinityPolicy *rawAntiAffinityPolicy
	if o.AntiAffinityPolicy != nil {
		antiAffinityPolicy = o.AntiAffinityPolicy.raw()
	}

	var jeriType *string
	if o.JunosEvpnRoutingInstanceVlanAware != nil {
		if o.JunosEvpnRoutingInstanceVlanAware.String() == enum.FeatureSwitchEnabled.String() {
			jeriType = toPtr(junosEvpnRoutingInstanceTypeVlanAware.Value)
		} else {
			jeriType = toPtr(junosEvpnRoutingInstanceTypeDefault.Value)
		}
	}

	var spineLeafLinks *addressingScheme
	if o.SpineLeafLinks != nil {
		spineLeafLinks = toPtr(addressingScheme(o.SpineLeafLinks.String()))
	}

	var spineSuperspineLinks *addressingScheme
	if o.SpineSuperspineLinks != nil {
		spineSuperspineLinks = toPtr(addressingScheme(o.SpineSuperspineLinks.String()))
	}

	var DefaultAnycastGWMAC *string
	if len(o.DefaultAnycastGWMAC.String()) > 0 {
		DefaultAnycastGWMAC = toPtr(o.DefaultAnycastGWMAC.String())
	}

	raw := rawFabricSettings{
		AntiAffinity:                          antiAffinityPolicy,
		DefaultAnycastGWMAC:                   DefaultAnycastGWMAC,
		DefaultSviL3Mtu:                       o.DefaultSviL3Mtu,
		EsiMacMsb:                             o.EsiMacMsb,
		EvpnGenerateType5HostRoutes:           stringerPtrToStringPtr(o.EvpnGenerateType5HostRoutes),
		ExternalRouterMtu:                     o.ExternalRouterMtu,
		FabricL3Mtu:                           o.FabricL3Mtu,
		Ipv6Enabled:                           o.Ipv6Enabled,
		JunosEvpnDuplicateMacRecoveryTime:     o.JunosEvpnDuplicateMacRecoveryTime,
		JunosEvpnMaxNexthopAndInterfaceNumber: stringerPtrToStringPtr(o.JunosEvpnMaxNexthopAndInterfaceNumber),
		JunosEvpnRoutingInstanceType:          jeriType,
		JunosExOverlayEcmp:                    stringerPtrToStringPtr(o.JunosExOverlayEcmp),
		JunosGracefulRestart:                  stringerPtrToStringPtr(o.JunosGracefulRestart),
		MaxEvpnRoutes:                         o.MaxEvpnRoutes,
		MaxExternalRoutes:                     o.MaxExternalRoutes,
		MaxFabricRoutes:                       o.MaxFabricRoutes,
		MaxMlagRoutes:                         o.MaxMlagRoutes,
		OptimiseSzFootprint:                   stringerPtrToStringPtr(o.OptimiseSzFootprint),
		OverlayControlProtocol:                stringerPtrToStringPtr(o.OverlayControlProtocol),
		SpineLeafLinks:                        spineLeafLinks,
		SpineSuperspineLinks:                  spineSuperspineLinks,
	}

	return json.Marshal(&raw)
}

func (o *FabricSettings) UnmarshalJSON(bytes []byte) error {
	var raw rawFabricSettings
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return fmt.Errorf("unmarshaling rawFabricSettings: %w", err)
	}

	var err error

	if o.AntiAffinityPolicy, err = raw.AntiAffinity.polish(); err != nil {
		return fmt.Errorf("parsing AntiAffinityPolicy: %w", err)
	}

	if raw.DefaultAnycastGWMAC != nil && *raw.DefaultAnycastGWMAC != "" {
		if o.DefaultAnycastGWMAC, err = net.ParseMAC(*raw.DefaultAnycastGWMAC); err != nil {
			return fmt.Errorf("parsing DefaultAnycastGWMAC: %w", err)
		}
	} else {
		o.DefaultAnycastGWMAC = nil
	}

	o.DefaultSviL3Mtu = raw.DefaultSviL3Mtu
	o.EsiMacMsb = raw.EsiMacMsb
	o.EvpnGenerateType5HostRoutes = featureSwitchEnumFromStringPtr(raw.EvpnGenerateType5HostRoutes)
	o.ExternalRouterMtu = raw.ExternalRouterMtu
	o.FabricL3Mtu = raw.FabricL3Mtu
	o.Ipv6Enabled = raw.Ipv6Enabled
	o.JunosEvpnDuplicateMacRecoveryTime = raw.JunosEvpnDuplicateMacRecoveryTime
	o.JunosEvpnMaxNexthopAndInterfaceNumber = featureSwitchEnumFromStringPtr(raw.JunosEvpnMaxNexthopAndInterfaceNumber)

	if raw.JunosEvpnRoutingInstanceType != nil {
		parsed := junosEvpnRoutingInstanceTypes.Parse(*raw.JunosEvpnRoutingInstanceType)
		if parsed == nil {
			return fmt.Errorf("cannot parse junos_evpn_routing_instance_type value %q", *raw.JunosEvpnRoutingInstanceType)
		}

		if *parsed == junosEvpnRoutingInstanceTypeDefault {
			o.JunosEvpnRoutingInstanceVlanAware = &enum.FeatureSwitchDisabled
		} else {
			o.JunosEvpnRoutingInstanceVlanAware = &enum.FeatureSwitchEnabled
		}
	} else {
		o.JunosEvpnRoutingInstanceVlanAware = nil
	}

	o.JunosExOverlayEcmp = featureSwitchEnumFromStringPtr(raw.JunosExOverlayEcmp)
	o.JunosGracefulRestart = featureSwitchEnumFromStringPtr(raw.JunosGracefulRestart)
	o.MaxEvpnRoutes = raw.MaxEvpnRoutes
	o.MaxExternalRoutes = raw.MaxExternalRoutes
	o.MaxFabricRoutes = raw.MaxFabricRoutes
	o.MaxMlagRoutes = raw.MaxMlagRoutes
	o.OptimiseSzFootprint = featureSwitchEnumFromStringPtr(raw.OptimiseSzFootprint)

	if raw.OverlayControlProtocol != nil {
		ocp := new(OverlayControlProtocol)
		if err := ocp.FromString(*raw.OverlayControlProtocol); err != nil {
			return fmt.Errorf("parsing overlay_control_protocol %q: %w", *raw.OverlayControlProtocol, err)
		}
	} else {
		o.OverlayControlProtocol = nil
	}

	if raw.SpineLeafLinks != nil {
		x := new(addressingScheme)
		i, err := x.parse()
		if err != nil {
			return fmt.Errorf("parsing spine_leaf_links: %w", err)
		}
		o.SpineLeafLinks = toPtr(AddressingScheme(i))
	} else {
		o.SpineLeafLinks = nil
	}

	if raw.SpineSuperspineLinks != nil {
		x := new(addressingScheme)
		i, err := x.parse()
		if err != nil {
			return fmt.Errorf("parsing spine_superspine_links: %w", err)
		}
		o.SpineSuperspineLinks = toPtr(AddressingScheme(i))
	} else {
		o.SpineSuperspineLinks = nil
	}

	return nil
}

func (o FabricSettings) rawBlueprintRequestFabricAddressingPolicy() *rawBlueprintRequestFabricAddressingPolicy {
	var spineSuperspineLinks addressingScheme
	if o.SpineSuperspineLinks != nil {
		spineSuperspineLinks = addressingScheme(o.SpineSuperspineLinks.String())
	}

	var spineLeafLinks addressingScheme
	if o.SpineLeafLinks != nil {
		spineLeafLinks = addressingScheme(o.SpineLeafLinks.String())
	}

	return &rawBlueprintRequestFabricAddressingPolicy{
		SpineSuperspineLinks: spineSuperspineLinks,
		SpineLeafLinks:       spineLeafLinks,
		FabricL3Mtu:          o.FabricL3Mtu,
	}
}

type rawFabricSettings struct {
	AntiAffinity                          *rawAntiAffinityPolicy `json:"anti_affinity,omitempty"`
	DefaultAnycastGWMAC                   *string                `json:"default_anycast_gw_mac,omitempty"`
	DefaultSviL3Mtu                       *uint16                `json:"default_svi_l3_mtu,omitempty"`
	EsiMacMsb                             *uint8                 `json:"esi_mac_msb,omitempty"`
	EvpnGenerateType5HostRoutes           *string                `json:"evpn_generate_type5_host_routes,omitempty"`
	ExternalRouterMtu                     *uint16                `json:"external_router_mtu,omitempty"`
	FabricL3Mtu                           *uint16                `json:"fabric_l3_mtu,omitempty"`
	Ipv6Enabled                           *bool                  `json:"ipv6_enabled,omitempty"`
	JunosEvpnDuplicateMacRecoveryTime     *uint16                `json:"junos_evpn_duplicate_mac_recovery_time,omitempty"`
	JunosEvpnMaxNexthopAndInterfaceNumber *string                `json:"junos_evpn_max_nexthop_and_interface_number,omitempty"`
	JunosEvpnRoutingInstanceType          *string                `json:"junos_evpn_routing_instance_type,omitempty"`
	JunosExOverlayEcmp                    *string                `json:"junos_ex_overlay_ecmp,omitempty"`
	JunosGracefulRestart                  *string                `json:"junos_graceful_restart,omitempty"`
	MaxEvpnRoutes                         *uint32                `json:"max_evpn_routes"`
	MaxExternalRoutes                     *uint32                `json:"max_external_routes"`
	MaxFabricRoutes                       *uint32                `json:"max_fabric_routes"`
	MaxMlagRoutes                         *uint32                `json:"max_mlag_routes"`
	OptimiseSzFootprint                   *string                `json:"optimise_sz_footprint,omitempty"`
	OverlayControlProtocol                *string                `json:"overlay_control_protocol,omitempty"`
	SpineLeafLinks                        *addressingScheme      `json:"spine_leaf_links,omitempty"`       // ['ipv4', 'ipv6', 'ipv4_ipv6'],
	SpineSuperspineLinks                  *addressingScheme      `json:"spine_superspine_links,omitempty"` // ['ipv4', 'ipv6', 'ipv4_ipv6']
	// leaf_loopbacks ['ipv4', 'ipv4_ipv6']
	// spine_loopbacks ['ipv4', 'ipv4_ipv6']
	// mlag_svi_subnets ['ipv4', 'ipv4_ipv6']
	// leaf_l3_peer_links ['ipv4', 'ipv4_ipv6']
	// default_fabric_evi_route_target
}

func (o *TwoStageL3ClosClient) getFabricSettings(ctx context.Context) (*FabricSettings, error) {
	var response FabricSettings
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintFabricSettings, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *TwoStageL3ClosClient) setFabricSettings(ctx context.Context, in *FabricSettings) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintFabricSettings, o.blueprintId),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
