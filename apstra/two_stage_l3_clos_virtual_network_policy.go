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
	JunosEvpnDuplicateMacRecoveryTime     *uint16 `json:"junos_evpn_duplicate_mac_recovery_time,omitempty"`
	MaxExternalRoutes                     *uint32 `json:"max_external_routes,omitempty"`
	JunosGracefulRestart                  *string `json:"junos_graceful_restart,omitempty"`           // enabled/disabled
	JunosEvpnRoutingInstanceType          *string `json:"junos_evpn_routing_instance_type,omitempty"` // default/vlan_aware
	EvpnGenerateType5HostRoutes           *string `json:"evpn_generate_type5_host_routes,omitempty"`  // enabled/disabled
	MaxFabricRoutes                       *uint32 `json:"max_fabric_routes,omitempty"`
	MaxMlagRoutes                         *uint32 `json:"max_mlag_routes,omitempty"`
	JunosExOverlayEcmp                    *string `json:"junos_ex_overlay_ecmp,omitempty"` // enabled/disabled
	DefaultSviL3Mtu                       *uint16 `json:"default_svi_l3_mtu,omitempty"`
	JunosEvpnMaxNexthopAndInterfaceNumber *string `json:"junos_evpn_max_nexthop_and_interface_number,omitempty"` // enabled/disabled
	ExternalRouterMtu                     *uint16 `json:"external_router_mtu,omitempty"`
	MaxEvpnRoutes                         *uint32 `json:"max_evpn_routes,omitempty"`
	//FrrRdVlanOffset                       string `json:"frr_rd_vlan_offset,omitempty"` // not exposed in web UI
	//CumulusBridgeMacDerivation            string `json:"cumulus_bridge_mac_derivation,omitempty"`   // skipping cumulus support
	//OverlayControlProtocol                string `json:"overlay_control_protocol,omitempty"`        // immutable
	//DefaultFabricEviRouteTarget           string `json:"default_fabric_evi_route_target,omitempty"` // undocumented
	//CumulusVxlanArpSuppression            string `json:"cumulus_vxlan_arp_suppression,omitempty"`   // skipping cumulus support
}

func (o *TwoStageL3ClosClient) setRawVirtualNetworkPolicy420(ctx context.Context, in *rawFabricSettings) error {
	if in.JunosEvpnDuplicateMacRecoveryTime == nil &&
		in.MaxExternalRoutes == nil &&
		in.JunosGracefulRestart == nil &&
		in.JunosEvpnRoutingInstanceType == nil &&
		in.EvpnGenerateType5HostRoutes == nil &&
		in.MaxFabricRoutes == nil &&
		in.MaxMlagRoutes == nil &&
		in.JunosExOverlayEcmp == nil &&
		in.DefaultSviL3Mtu == nil &&
		in.JunosEvpnMaxNexthopAndInterfaceNumber == nil &&
		in.ExternalRouterMtu == nil &&
		in.MaxEvpnRoutes == nil {
		return nil // nothing to do if all relevant input fields are nil
	}

	if !o.client.apiVersion.Equal(version.Must(version.NewVersion(apstra420))) {
		return fmt.Errorf("setRawVirtualNetworkPolicy420() must not be invoked with apstra %s", o.client.apiVersion)
	}

	apiInput := rawVirtualNetworkPolicy420{
		JunosEvpnDuplicateMacRecoveryTime:     in.JunosEvpnDuplicateMacRecoveryTime,
		MaxExternalRoutes:                     in.MaxExternalRoutes,
		JunosGracefulRestart:                  in.JunosGracefulRestart,
		JunosEvpnRoutingInstanceType:          in.JunosEvpnRoutingInstanceType,
		EvpnGenerateType5HostRoutes:           in.EvpnGenerateType5HostRoutes,
		MaxFabricRoutes:                       in.MaxFabricRoutes,
		MaxMlagRoutes:                         in.MaxMlagRoutes,
		JunosExOverlayEcmp:                    in.JunosExOverlayEcmp,
		DefaultSviL3Mtu:                       in.DefaultSviL3Mtu,
		JunosEvpnMaxNexthopAndInterfaceNumber: in.JunosEvpnMaxNexthopAndInterfaceNumber,
		ExternalRouterMtu:                     in.ExternalRouterMtu,
		MaxEvpnRoutes:                         in.MaxEvpnRoutes,
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
