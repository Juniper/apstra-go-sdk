// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
)

const (
	apiUrlBlueprintVirtualNetworkPolicy = apiUrlBlueprintByIdPrefix + "virtual-network-policy"
)

type rawVirtualNetworkPolicy420 struct {
	DefaultSviL3Mtu                       *uint16 `json:"default_svi_l3_mtu,omitempty"`
	EvpnGenerateType5HostRoutes           *string `json:"evpn_generate_type5_host_routes,omitempty"`
	ExternalRouterMtu                     *uint16 `json:"external_router_mtu,omitempty"`
	JunosEvpnDuplicateMacRecoveryTime     *uint16 `json:"junos_evpn_duplicate_mac_recovery_time,omitempty"`
	JunosEvpnMaxNexthopAndInterfaceNumber *string `json:"junos_evpn_max_nexthop_and_interface_number,omitempty"`
	JunosEvpnRoutingInstanceType          *string `json:"junos_evpn_routing_instance_type,omitempty"`
	JunosExOverlayEcmp                    *string `json:"junos_ex_overlay_ecmp,omitempty"`
	JunosGracefulRestart                  *string `json:"junos_graceful_restart,omitempty"`
	MaxEvpnRoutes                         *uint32 `json:"max_evpn_routes"`
	MaxExternalRoutes                     *uint32 `json:"max_external_routes"`
	MaxFabricRoutes                       *uint32 `json:"max_fabric_routes"`
	MaxMlagRoutes                         *uint32 `json:"max_mlag_routes"`
	OverlayControlProtocol                *string `json:"overlay_control_protocol,omitempty"`
}

func (o *TwoStageL3ClosClient) getVirtualNetworkPolicy420(ctx context.Context) (*rawVirtualNetworkPolicy420, error) {
	vnpNodeIds, err := o.NodeIdsByType(ctx, NodeTypeVirtualNetworkPolicy)
	if err != nil {
		return nil, err
	}
	if len(vnpNodeIds) != 1 {
		return nil, fmt.Errorf("expected 1 %s node, got %d", NodeTypeVirtualNetworkPolicy.String(), len(vnpNodeIds))
	}

	var result rawVirtualNetworkPolicy420

	err = o.client.GetNode(ctx, o.blueprintId, vnpNodeIds[0], &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (o *TwoStageL3ClosClient) setVirtualNetworkPolicy420(ctx context.Context, in *rawFabricSettings) error {
	if !compatibility.EqApstra420.Check(o.client.apiVersion) {
		return fmt.Errorf("setRawVirtualNetworkPolicy420() must not be invoked with apstra %s", o.client.apiVersion)
	}

	apiInput := rawVirtualNetworkPolicy420{
		DefaultSviL3Mtu:                       in.DefaultSviL3Mtu,
		EvpnGenerateType5HostRoutes:           in.EvpnGenerateType5HostRoutes,
		ExternalRouterMtu:                     in.ExternalRouterMtu,
		JunosEvpnDuplicateMacRecoveryTime:     in.JunosEvpnDuplicateMacRecoveryTime,
		JunosEvpnMaxNexthopAndInterfaceNumber: in.JunosEvpnMaxNexthopAndInterfaceNumber,
		JunosEvpnRoutingInstanceType:          in.JunosEvpnRoutingInstanceType,
		JunosExOverlayEcmp:                    in.JunosExOverlayEcmp,
		JunosGracefulRestart:                  in.JunosGracefulRestart,
		MaxEvpnRoutes:                         in.MaxEvpnRoutes,
		MaxExternalRoutes:                     in.MaxExternalRoutes,
		MaxFabricRoutes:                       in.MaxFabricRoutes,
		MaxMlagRoutes:                         in.MaxMlagRoutes,
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintVirtualNetworkPolicy, o.blueprintId),
		apiInput: &apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
