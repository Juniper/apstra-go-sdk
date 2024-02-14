package apstra

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/orsinium-labs/enum"
	"net/http"
)

const (
	apiUrlBlueprintFabricSettings = apiUrlBlueprintById + apiUrlPathDelim + "fabric-settings"
)

type junosEvpnRoutingInstanceType enum.Member[string]

var (
	junosEvpnRoutingInstanceTypeVlanAware = junosEvpnRoutingInstanceType{Value: "vlan_aware"}
	junosEvpnRoutingInstanceTypeDefault   = junosEvpnRoutingInstanceType{Value: "default"}
	junosEvpnRoutingInstanceTypes         = enum.New(junosEvpnRoutingInstanceTypeVlanAware, junosEvpnRoutingInstanceTypeDefault)
)

type FabricSettings struct {
	JunosEvpnDuplicateMacRecoveryTime     *uint16                 // supported in 4.2.0 (not in GUI)
	MaxExternalRoutes                     *uint32                 // supported in 4.2.0 ( virtual-network-policy)
	EsiMacMsb                             *uint8                  // supported in 4.2.0 (fabric-addressing-policy)
	JunosGracefulRestart                  *FeatureSwitchEnum      // supported in 4.2.0 ( virtual-network-policy)
	OptimiseSzFootprint                   *FeatureSwitchEnum      // supported in 4.2.0 (patch fabric-settings node)
	JunosEvpnRoutingInstanceVlanAware     *FeatureSwitchEnum      // supported in 4.2.0 ( virtual-network-policy)
	EvpnGenerateType5HostRoutes           *FeatureSwitchEnum      // supported in 4.2.0 ( virtual-network-policy)
	MaxFabricRoutes                       *uint32                 // supported in 4.2.0 ( virtual-network-policy)
	MaxMlagRoutes                         *uint32                 // supported in 4.2.0 ( virtual-network-policy)
	JunosExOverlayEcmp                    *FeatureSwitchEnum      // supported in 4.2.0 ( virtual-network-policy)
	DefaultSviL3Mtu                       *uint16                 // supported in 4.2.0 ( virtual-network-policy)
	JunosEvpnMaxNexthopAndInterfaceNumber *FeatureSwitchEnum      // supported in 4.2.0 ( virtual-network-policy)
	FabricL3Mtu                           *uint16                 // supported in 4.2.0 (fabric-addressing-policy)
	Ipv6Enabled                           *bool                   // supported in 4.2.0 (fabric-addressing-policy)
	OverlayControlProtocol                *OverlayControlProtocol // supported in 4.2.0 ( virtual-network-policy)
	ExternalRouterMtu                     *uint16                 // supported in 4.2.0 ( virtual-network-policy)
	MaxEvpnRoutes                         *uint32                 // supported in 4.2.0 (virtual-network-policy)
	AntiAffinityPolicy                    *AntiAffinityPolicy     // supported in 4.2.0 (anti-affinity-policy)
	//DefaultFabricEviRouteTarget                 string   					Not Exposed via WebUI
	//FrrRdVlanOffset                             string					Not Exposed via WebUI
}

func (o FabricSettings) raw() *rawFabricSettings {
	var antiAffinityPolicy *rawAntiAffinityPolicy
	if o.AntiAffinityPolicy != nil {
		antiAffinityPolicy = o.AntiAffinityPolicy.raw()
	}

	var jeriType *string
	if o.JunosEvpnRoutingInstanceVlanAware != nil {
		if o.JunosEvpnRoutingInstanceVlanAware.String() == FeatureSwitchEnumEnabled.String() {
			jeriType = toPtr(junosEvpnRoutingInstanceTypeVlanAware.Value)
		} else {
			jeriType = toPtr(junosEvpnRoutingInstanceTypeDefault.Value)
		}
	}
	return &rawFabricSettings{
		JunosEvpnDuplicateMacRecoveryTime:     o.JunosEvpnDuplicateMacRecoveryTime,
		MaxExternalRoutes:                     o.MaxExternalRoutes,
		EsiMacMsb:                             o.EsiMacMsb,
		JunosGracefulRestart:                  stringerPtrToStringPtr(o.JunosGracefulRestart),
		OptimiseSzFootprint:                   stringerPtrToStringPtr(o.OptimiseSzFootprint),
		JunosEvpnRoutingInstanceType:          jeriType,
		EvpnGenerateType5HostRoutes:           stringerPtrToStringPtr(o.EvpnGenerateType5HostRoutes),
		MaxFabricRoutes:                       o.MaxFabricRoutes,
		MaxMlagRoutes:                         o.MaxMlagRoutes,
		JunosExOverlayEcmp:                    stringerPtrToStringPtr(o.JunosExOverlayEcmp),
		DefaultSviL3Mtu:                       o.DefaultSviL3Mtu,
		JunosEvpnMaxNexthopAndInterfaceNumber: stringerPtrToStringPtr(o.JunosEvpnMaxNexthopAndInterfaceNumber),
		FabricL3Mtu:                           o.FabricL3Mtu,
		Ipv6Enabled:                           o.Ipv6Enabled,
		OverlayControlProtocol:                stringerPtrToStringPtr(o.OverlayControlProtocol),
		ExternalRouterMtu:                     o.ExternalRouterMtu,
		MaxEvpnRoutes:                         o.MaxEvpnRoutes,
		AntiAffinity:                          antiAffinityPolicy,
	}
}

type rawFabricSettings struct {
	JunosEvpnDuplicateMacRecoveryTime     *uint16                `json:"junos_evpn_duplicate_mac_recovery_time,omitempty"`
	MaxExternalRoutes                     *uint32                `json:"max_external_routes,omitempty"`
	EsiMacMsb                             *uint8                 `json:"esi_mac_msb,omitempty"`
	JunosGracefulRestart                  *string                `json:"junos_graceful_restart,omitempty"`
	OptimiseSzFootprint                   *string                `json:"optimise_sz_footprint,omitempty"`
	JunosEvpnRoutingInstanceType          *string                `json:"junos_evpn_routing_instance_type,omitempty"`
	EvpnGenerateType5HostRoutes           *string                `json:"evpn_generate_type5_host_routes,omitempty"`
	MaxFabricRoutes                       *uint32                `json:"max_fabric_routes,omitempty"`
	MaxMlagRoutes                         *uint32                `json:"max_mlag_routes,omitempty"`
	JunosExOverlayEcmp                    *string                `json:"junos_ex_overlay_ecmp,omitempty"`
	DefaultSviL3Mtu                       *uint16                `json:"default_svi_l3_mtu,omitempty"`
	JunosEvpnMaxNexthopAndInterfaceNumber *string                `json:"junos_evpn_max_nexthop_and_interface_number,omitempty"`
	FabricL3Mtu                           *uint16                `json:"fabric_l3_mtu,omitempty"`
	Ipv6Enabled                           *bool                  `json:"ipv6_enabled,omitempty"`
	OverlayControlProtocol                *string                `json:"overlay_control_protocol,omitempty"`
	ExternalRouterMtu                     *uint16                `json:"external_router_mtu,omitempty"`
	MaxEvpnRoutes                         *uint32                `json:"max_evpn_routes,omitempty"`
	AntiAffinity                          *rawAntiAffinityPolicy `json:"anti_affinity,omitempty"`
	//FrrRdVlanOffset                       string                `json:"frr_rd_vlan_offset"`
	//DefaultFabricEviRouteTarget           string                `json:"default_fabric_evi_route_target"`
}

func (o rawFabricSettings) polish() (*FabricSettings, error) {
	var ocp *OverlayControlProtocol
	if o.OverlayControlProtocol != nil {
		ocp = new(OverlayControlProtocol)
		err := ocp.FromString(*o.OverlayControlProtocol)
		if err != nil {
			return nil, fmt.Errorf("failed to parse overlay_control_protocol value %q", *o.OverlayControlProtocol)
		}
	}

	var junosEvpnRoutingInstanceVlanAware *FeatureSwitchEnum
	if o.JunosEvpnRoutingInstanceType != nil {
		x := junosEvpnRoutingInstanceTypes.Parse(*o.JunosEvpnRoutingInstanceType)
		if x == nil {
			return nil, fmt.Errorf("failed to parse junos_evpn_routing_instance_type value %q", *o.JunosEvpnRoutingInstanceType)
		}
		if *x == junosEvpnRoutingInstanceTypeDefault {
			junosEvpnRoutingInstanceVlanAware = &FeatureSwitchEnumDisabled
		} else {
			junosEvpnRoutingInstanceVlanAware = &FeatureSwitchEnumEnabled
		}
	}

	antiAffinityPolicy, err := o.AntiAffinity.polish()
	if err != nil {
		return nil, err
	}
	return &FabricSettings{
		JunosEvpnDuplicateMacRecoveryTime:     o.JunosEvpnDuplicateMacRecoveryTime,
		MaxExternalRoutes:                     o.MaxExternalRoutes,
		EsiMacMsb:                             o.EsiMacMsb,
		JunosGracefulRestart:                  featureSwitchEnumFromStringPtr(o.JunosGracefulRestart),
		OptimiseSzFootprint:                   featureSwitchEnumFromStringPtr(o.OptimiseSzFootprint),
		JunosEvpnRoutingInstanceVlanAware:     junosEvpnRoutingInstanceVlanAware,
		EvpnGenerateType5HostRoutes:           featureSwitchEnumFromStringPtr(o.EvpnGenerateType5HostRoutes),
		MaxFabricRoutes:                       o.MaxFabricRoutes,
		MaxMlagRoutes:                         o.MaxMlagRoutes,
		JunosExOverlayEcmp:                    featureSwitchEnumFromStringPtr(o.JunosExOverlayEcmp),
		DefaultSviL3Mtu:                       o.DefaultSviL3Mtu,
		JunosEvpnMaxNexthopAndInterfaceNumber: featureSwitchEnumFromStringPtr(o.JunosEvpnMaxNexthopAndInterfaceNumber),
		FabricL3Mtu:                           o.FabricL3Mtu,
		Ipv6Enabled:                           o.Ipv6Enabled,
		OverlayControlProtocol:                ocp,
		ExternalRouterMtu:                     o.ExternalRouterMtu,
		MaxEvpnRoutes:                         o.MaxEvpnRoutes,
		AntiAffinityPolicy:                    antiAffinityPolicy,
	}, nil
}

func (o *TwoStageL3ClosClient) getFabricSettings(ctx context.Context) (*rawFabricSettings, error) {
	var response rawFabricSettings
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

func (o *TwoStageL3ClosClient) getFabricSettings420(ctx context.Context) (*rawFabricSettings, error) {
	fabricAddressingPolicy, err := o.GetFabricAddressingPolicy(ctx)
	if err != nil {
		return nil, err
	}

	virtualNetworkPolicy, err := o.getVirtualNetworkPolicy420(ctx)
	if err != nil {
		return nil, err
	}

	optimiseFootprint, err := o.getSzFootprintOptimization420(ctx)
	if err != nil {
		return nil, err
	}

	antiAffinityPolicy, err := o.getAntiAffinityPolicy420(ctx)
	if err != nil {
		return nil, err
	}

	return &rawFabricSettings{
		EsiMacMsb:   fabricAddressingPolicy.EsiMacMsb,
		FabricL3Mtu: fabricAddressingPolicy.FabricL3Mtu,
		Ipv6Enabled: fabricAddressingPolicy.Ipv6Enabled,

		JunosEvpnDuplicateMacRecoveryTime:     virtualNetworkPolicy.JunosEvpnDuplicateMacRecoveryTime,
		MaxExternalRoutes:                     virtualNetworkPolicy.MaxExternalRoutes,
		JunosGracefulRestart:                  virtualNetworkPolicy.JunosGracefulRestart,
		JunosEvpnRoutingInstanceType:          virtualNetworkPolicy.JunosEvpnRoutingInstanceType,
		EvpnGenerateType5HostRoutes:           virtualNetworkPolicy.EvpnGenerateType5HostRoutes,
		MaxFabricRoutes:                       virtualNetworkPolicy.MaxFabricRoutes,
		MaxMlagRoutes:                         virtualNetworkPolicy.MaxMlagRoutes,
		JunosExOverlayEcmp:                    virtualNetworkPolicy.JunosExOverlayEcmp,
		DefaultSviL3Mtu:                       virtualNetworkPolicy.DefaultSviL3Mtu,
		JunosEvpnMaxNexthopAndInterfaceNumber: virtualNetworkPolicy.JunosEvpnMaxNexthopAndInterfaceNumber,
		ExternalRouterMtu:                     virtualNetworkPolicy.ExternalRouterMtu,
		MaxEvpnRoutes:                         virtualNetworkPolicy.MaxEvpnRoutes,
		OverlayControlProtocol:                virtualNetworkPolicy.OverlayControlProtocol,

		OptimiseSzFootprint: &optimiseFootprint,

		AntiAffinity: antiAffinityPolicy,
	}, nil
}

func (o *TwoStageL3ClosClient) setFabricSettings(ctx context.Context, in *rawFabricSettings) error {
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

// setFabricSettings420 does the same job as setFabricSettings, but for Apstra 4.2.0, which controls
// the parameters in rawFabricSettings in 3 different places
func (o *TwoStageL3ClosClient) setFabricSettings420(ctx context.Context, in *rawFabricSettings) error {
	err := o.SetFabricAddressingPolicy(ctx, &TwoStageL3ClosFabricAddressingPolicy{
		Ipv6Enabled: in.Ipv6Enabled,
		EsiMacMsb:   in.EsiMacMsb,
		FabricL3Mtu: in.FabricL3Mtu,
	})

	err = o.setVirtualNetworkPolicy420(ctx, in)
	if err != nil {
		return err
	}

	err = o.setSzFootprintOptimization420(ctx, in.OptimiseSzFootprint)
	if err != nil {
		return err
	}

	err = o.setAntiAffinityPolicy420(ctx, in.AntiAffinity)
	if err != nil {
		return err
	}

	return nil
}

func (o *TwoStageL3ClosClient) getSzFootprintOptimization420(ctx context.Context) (string, error) {
	if !o.client.apiVersion.Equal(version.Must(version.NewVersion(apstra420))) {
		return "", fmt.Errorf("getSzFootprintOptimization420() must not be invoked with apstra %s", o.client.apiVersion)
	}

	securityZonePolicyNodeIds, err := o.NodeIdsByType(ctx, NodeTypeSecurityZonePolicy)
	if err != nil {
		return "", err
	}
	if len(securityZonePolicyNodeIds) != 1 {
		return "", fmt.Errorf("expected 1 %s, got %d", NodeTypeSecurityZonePolicy.String(), len(securityZonePolicyNodeIds))
	}

	var node struct {
		FootprintOptimise string `json:"footprint_optimise"`
	}

	err = o.client.GetNode(ctx, o.blueprintId, securityZonePolicyNodeIds[0], &node)
	if err != nil {
		return "", err
	}

	return node.FootprintOptimise, nil
}

func (o *TwoStageL3ClosClient) setSzFootprintOptimization420(ctx context.Context, in *string) error {
	if in == nil {
		return nil
	}

	if !o.client.apiVersion.Equal(version.Must(version.NewVersion(apstra420))) {
		return fmt.Errorf("setSzFootprintOptimization420() must not be invoked with apstra %s", o.client.apiVersion)
	}

	securityZonePolicyNodeIds, err := o.NodeIdsByType(ctx, NodeTypeSecurityZonePolicy)
	if err != nil {
		return err
	}
	if len(securityZonePolicyNodeIds) != 1 {
		return fmt.Errorf("expected 1 %s, got %d", NodeTypeSecurityZonePolicy.String(), len(securityZonePolicyNodeIds))
	}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPatch,
		urlStr: fmt.Sprintf(apiUrlBlueprintNodeById, o.blueprintId, securityZonePolicyNodeIds[0]),
		apiInput: &struct {
			FootprintOptimise string `json:"footprint_optimise"`
		}{
			FootprintOptimise: *in,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to patch %s node - %w", NodeTypeSecurityZonePolicy.String(), convertTtaeToAceWherePossible(err))
	}

	return nil
}
