package apstra

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-version"
	"net/http"
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
	MaxEvpnRoutes                         *uint32 `json:"max_evpn_routes,omitempty"`
	MaxExternalRoutes                     *uint32 `json:"max_external_routes,omitempty"`
	MaxFabricRoutes                       *uint32 `json:"max_fabric_routes,omitempty"`
	MaxMlagRoutes                         *uint32 `json:"max_mlag_routes,omitempty"`
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
	if in.DefaultSviL3Mtu == nil &&
		in.EvpnGenerateType5HostRoutes == nil &&
		in.ExternalRouterMtu == nil &&
		in.JunosEvpnDuplicateMacRecoveryTime == nil &&
		in.JunosEvpnMaxNexthopAndInterfaceNumber == nil &&
		in.JunosEvpnRoutingInstanceType == nil &&
		in.JunosExOverlayEcmp == nil &&
		in.JunosGracefulRestart == nil &&
		in.MaxEvpnRoutes == nil &&
		in.MaxExternalRoutes == nil &&
		in.MaxFabricRoutes == nil &&
		in.MaxMlagRoutes == nil {
		return nil // nothing to do if all relevant input fields are nil
	}

	if !o.client.apiVersion.Equal(version.Must(version.NewVersion(apstra420))) {
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

type rawVirtualNetworkPolicy41x struct {
	EvpnGenerateType5HostRoutes *string `json:"evpn_generate_type5_host_routes"`
	ExternalRouterMtu           *uint16 `json:"external_router_mtu"`
	MaxFabricRoutes             *uint32 `json:"max_fabric_routes"`
	MaxMlagRoutes               *uint32 `json:"max_mlag_routes"`
	MaxEvpnRoutes               *uint32 `json:"max_evpn_routes"`
	MaxExternalRoutes           *uint32 `json:"max_external_routes"`
	OverlayControlProtocol      *string `json:"overlay_control_protocol"`
}

func (o *TwoStageL3ClosClient) getVirtualNetworkPolicy41x(ctx context.Context) (*rawVirtualNetworkPolicy41x, error) {
	vnpNodeIds, err := o.NodeIdsByType(ctx, NodeTypeVirtualNetworkPolicy)
	if err != nil {
		return nil, err
	}
	if len(vnpNodeIds) != 1 {
		return nil, fmt.Errorf("expected 1 %s node, got %d", NodeTypeVirtualNetworkPolicy.String(), len(vnpNodeIds))
	}

	var result rawVirtualNetworkPolicy41x

	err = o.client.GetNode(ctx, o.blueprintId, vnpNodeIds[0], &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (o *TwoStageL3ClosClient) setVirtualNetworkPolicy41x(ctx context.Context, in *rawFabricSettings) error {
	if in.EvpnGenerateType5HostRoutes == nil &&
		in.ExternalRouterMtu == nil &&
		in.MaxEvpnRoutes == nil &&
		in.MaxExternalRoutes == nil &&
		in.MaxFabricRoutes == nil &&
		in.MaxMlagRoutes == nil {
		return nil // nothing to do if all relevant input fields are nil
	}

	if !o.client.apiVersion.Equal(version.Must(version.NewVersion(apstra412))) {
		return fmt.Errorf("setRawVirtualNetworkPolicy412() must not be invoked with apstra %s", o.client.apiVersion)
	}

	apiInput := rawVirtualNetworkPolicy41x{
		EvpnGenerateType5HostRoutes: in.EvpnGenerateType5HostRoutes,
		ExternalRouterMtu:           in.ExternalRouterMtu,
		MaxEvpnRoutes:               in.MaxEvpnRoutes,
		MaxExternalRoutes:           in.MaxExternalRoutes,
		MaxFabricRoutes:             in.MaxFabricRoutes,
		MaxMlagRoutes:               in.MaxMlagRoutes,
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
