package apstra

import (
	"context"
	"fmt"
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
	JunosEvpnDuplicateMacRecoveryTime             *uint16
	MaxExternalRoutes                             *uint32
	EsiMacMsb                                     *uint8
	JunosGracefulRestart                          bool
	OptimiseSzFootprint                           bool
	JunosEvpnRoutingInstanceVlanAware             bool
	EvpnGenerateType5HostRoutes                   bool
	MaxFabricRoutes                               *uint32
	MaxMlagRoutes                                 *uint32
	JunosExOverlayEcmpDisabled                    bool
	DefaultSviL3Mtu                               *uint16
	JunosEvpnMaxNexthopAndInterfaceNumberDisabled bool
	FabricL3Mtu                                   *uint16
	Ipv6Enabled                                   bool
	OverlayControlProtocol                        OverlayControlProtocol
	ExternalRouterMtu                             *uint16
	MaxEvpnRoutes                                 *uint32
	AntiAffinityPolicy                            *AntiAffinityPolicy
	//DefaultFabricEviRouteTarget                 string
	//FrrRdVlanOffset                             string
}

func (o FabricSettings) raw() *rawFabricSettings {
	junosGracefulRestart := FeatureSwitchEnumDisabled.String()
	if o.JunosGracefulRestart {
		junosGracefulRestart = FeatureSwitchEnumEnabled.String()
	}

	optimiseSzFootprint := FeatureSwitchEnumDisabled.String()
	if o.OptimiseSzFootprint {
		optimiseSzFootprint = FeatureSwitchEnumEnabled.String()
	}

	junosRoutingInstanceType := junosEvpnRoutingInstanceTypeDefault.Value
	if o.JunosEvpnRoutingInstanceVlanAware {
		junosRoutingInstanceType = junosEvpnRoutingInstanceTypeVlanAware.Value
	}

	evpnGenerateType5HostRoutes := FeatureSwitchEnumDisabled.String()
	if o.EvpnGenerateType5HostRoutes {
		evpnGenerateType5HostRoutes = FeatureSwitchEnumEnabled.String()
	}

	junosExOverlayEcmp := FeatureSwitchEnumDisabled.String()
	if !o.JunosExOverlayEcmpDisabled {
		junosExOverlayEcmp = FeatureSwitchEnumEnabled.String()
	}

	junosEvpnMaxNexthopAndInterfaceNumber := FeatureSwitchEnumDisabled.String()
	if !o.JunosEvpnMaxNexthopAndInterfaceNumberDisabled {
		junosEvpnMaxNexthopAndInterfaceNumber = FeatureSwitchEnumEnabled.String()
	}

	var antiAffinityPolicy *rawAntiAffinityPolicy
	if o.AntiAffinityPolicy != nil {
		antiAffinityPolicy = o.AntiAffinityPolicy.raw()
	}

	return &rawFabricSettings{
		JunosEvpnDuplicateMacRecoveryTime:     o.JunosEvpnDuplicateMacRecoveryTime,
		MaxExternalRoutes:                     o.MaxExternalRoutes,
		EsiMacMsb:                             o.EsiMacMsb,
		JunosGracefulRestart:                  junosGracefulRestart,
		OptimiseSzFootprint:                   optimiseSzFootprint,
		JunosEvpnRoutingInstanceType:          junosRoutingInstanceType,
		EvpnGenerateType5HostRoutes:           evpnGenerateType5HostRoutes,
		MaxFabricRoutes:                       o.MaxFabricRoutes,
		MaxMlagRoutes:                         o.MaxMlagRoutes,
		JunosExOverlayEcmp:                    junosExOverlayEcmp,
		DefaultSviL3Mtu:                       o.DefaultSviL3Mtu,
		JunosEvpnMaxNexthopAndInterfaceNumber: junosEvpnMaxNexthopAndInterfaceNumber,
		FabricL3Mtu:                           o.FabricL3Mtu,
		Ipv6Enabled:                           o.Ipv6Enabled,
		OverlayControlProtocol:                o.OverlayControlProtocol.String(),
		ExternalRouterMtu:                     o.ExternalRouterMtu,
		MaxEvpnRoutes:                         o.MaxEvpnRoutes,
		AntiAffinity:                          antiAffinityPolicy,
	}
}

type rawFabricSettings struct {
	JunosEvpnDuplicateMacRecoveryTime     *uint16                `json:"junos_evpn_duplicate_mac_recovery_time,omitempty"`
	MaxExternalRoutes                     *uint32                `json:"max_external_routes,omitempty"`
	EsiMacMsb                             *uint8                 `json:"esi_mac_msb,omitempty"`
	JunosGracefulRestart                  string                 `json:"junos_graceful_restart,omitempty"`
	OptimiseSzFootprint                   string                 `json:"optimise_sz_footprint,omitempty"`
	JunosEvpnRoutingInstanceType          string                 `json:"junos_evpn_routing_instance_type,omitempty"`
	EvpnGenerateType5HostRoutes           string                 `json:"evpn_generate_type5_host_routes,omitempty"`
	MaxFabricRoutes                       *uint32                `json:"max_fabric_routes,omitempty"`
	MaxMlagRoutes                         *uint32                `json:"max_mlag_routes,omitempty"`
	JunosExOverlayEcmp                    string                 `json:"junos_ex_overlay_ecmp,omitempty"`
	DefaultSviL3Mtu                       *uint16                `json:"default_svi_l3_mtu,omitempty"`
	JunosEvpnMaxNexthopAndInterfaceNumber string                 `json:"junos_evpn_max_nexthop_and_interface_number,omitempty"`
	FabricL3Mtu                           *uint16                `json:"fabric_l3_mtu,omitempty"`
	Ipv6Enabled                           bool                   `json:"ipv6_enabled"`
	OverlayControlProtocol                string                 `json:"overlay_control_protocol,omitempty"`
	ExternalRouterMtu                     *uint16                `json:"external_router_mtu,omitempty"`
	MaxEvpnRoutes                         *uint32                `json:"max_evpn_routes,omitempty"`
	AntiAffinity                          *rawAntiAffinityPolicy `json:"anti_affinity,omitempty"`
	//FrrRdVlanOffset                       string                `json:"frr_rd_vlan_offset"`
	//DefaultFabricEviRouteTarget           string                `json:"default_fabric_evi_route_target"`
}

func (o rawFabricSettings) polish() (*FabricSettings, error) {

	junosGracefulRestart := FeatureSwitchEnums.Parse(o.JunosGracefulRestart)
	if junosGracefulRestart == nil {
		return nil, fmt.Errorf("failed to parse junos_graceful_restart value %q", o.JunosGracefulRestart)
	}

	optimizeSzFootprint := FeatureSwitchEnums.Parse(o.OptimiseSzFootprint)
	if optimizeSzFootprint == nil {
		return nil, fmt.Errorf("failed to parse optimise_sz_footprint value %q", o.OptimiseSzFootprint)
	}

	junosRoutingInstanceType := junosEvpnRoutingInstanceTypes.Parse(o.JunosEvpnRoutingInstanceType)
	if junosRoutingInstanceType == nil {
		return nil, fmt.Errorf("failed to parse junos_evpn_routing_instance_type value %q", o.JunosEvpnRoutingInstanceType)
	}

	evpnGenerateType5HostRoutes := FeatureSwitchEnums.Parse(o.EvpnGenerateType5HostRoutes)
	if evpnGenerateType5HostRoutes == nil {
		return nil, fmt.Errorf("failed to parse evpn_generate_type5_host_routes value %q", o.EvpnGenerateType5HostRoutes)
	}

	junosExOverlayEcmp := FeatureSwitchEnums.Parse(o.JunosExOverlayEcmp)
	if junosExOverlayEcmp == nil {
		return nil, fmt.Errorf("failed to parse junos_ex_overlay_ecmp value %q", o.JunosExOverlayEcmp)
	}

	junosEvpnMaxNexthopAndInterfaceNumber := FeatureSwitchEnums.Parse(o.JunosEvpnMaxNexthopAndInterfaceNumber)
	if junosEvpnMaxNexthopAndInterfaceNumber == nil {
		return nil, fmt.Errorf("failed to parse junos_evpn_max_nexthop_and_interface_number value %q", o.JunosEvpnMaxNexthopAndInterfaceNumber)
	}

	var ocp OverlayControlProtocol
	err := ocp.FromString(o.OverlayControlProtocol)
	if err != nil {
		return nil, fmt.Errorf("failed to parse overlay_control_protocol value %q", o.OverlayControlProtocol)
	}

	antiAffinityPolicy, err := o.AntiAffinity.polish()
	if err != nil {
		return nil, err
	}
	return &FabricSettings{
		JunosEvpnDuplicateMacRecoveryTime:             o.JunosEvpnDuplicateMacRecoveryTime,
		MaxExternalRoutes:                             o.MaxExternalRoutes,
		EsiMacMsb:                                     o.EsiMacMsb,
		JunosGracefulRestart:                          junosGracefulRestart.Value == FeatureSwitchEnumEnabled.Value,
		OptimiseSzFootprint:                           optimizeSzFootprint.Value == FeatureSwitchEnumEnabled.Value,
		JunosEvpnRoutingInstanceVlanAware:             junosRoutingInstanceType.Value == junosEvpnRoutingInstanceTypeVlanAware.Value,
		EvpnGenerateType5HostRoutes:                   evpnGenerateType5HostRoutes.Value == FeatureSwitchEnumEnabled.Value,
		MaxFabricRoutes:                               o.MaxFabricRoutes,
		MaxMlagRoutes:                                 o.MaxMlagRoutes,
		JunosExOverlayEcmpDisabled:                    junosExOverlayEcmp.Value == FeatureSwitchEnumDisabled.Value,
		DefaultSviL3Mtu:                               o.DefaultSviL3Mtu,
		JunosEvpnMaxNexthopAndInterfaceNumberDisabled: junosEvpnMaxNexthopAndInterfaceNumber.Value == FeatureSwitchEnumEnabled.Value,
		FabricL3Mtu:                                   o.FabricL3Mtu,
		Ipv6Enabled:                                   o.Ipv6Enabled,
		OverlayControlProtocol:                        ocp,
		ExternalRouterMtu:                             o.ExternalRouterMtu,
		MaxEvpnRoutes:                                 o.MaxEvpnRoutes,
		AntiAffinityPolicy:                            antiAffinityPolicy,
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
