package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/hashicorp/go-version"
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

type FabricSettings struct { //										 4.2.0							4.1.2							4.1.1							4.1.0
	AntiAffinityPolicy                    *AntiAffinityPolicy     // /anti-affinity-policy			/anti-affinity-policy			/anti-affinity-policy			/anti-affinity-policy
	DefaultSviL3Mtu                       *uint16                 // virtual_network_policy node	not supported					not supported					not supported.
	EsiMacMsb                             *uint8                  // /fabric-addressing-policy		/fabric-addressing-policy		/fabric-addressing-policy		/fabric-addressing-policy
	EvpnGenerateType5HostRoutes           *enum.FeatureSwitchEnum // virtual_network_policy node	virtual_network_policy node		virtual_network_policy node		virtual_network_policy node
	ExternalRouterMtu                     *uint16                 // virtual_network_policy node	virtual_network_policy node		virtual_network_policy node		virtual_network_policy node
	FabricL3Mtu                           *uint16                 // /fabric-addressing-policy		not supported					not supported					not supported
	Ipv6Enabled                           *bool                   // /fabric-addressing-policy		/fabric-addressing-policy		/fabric-addressing-policy		/fabric-addressing-policy
	JunosEvpnDuplicateMacRecoveryTime     *uint16                 // virtual_network_policy node	not supported					not supported					not supported
	JunosEvpnMaxNexthopAndInterfaceNumber *enum.FeatureSwitchEnum // virtual_network_policy node	not supported					not supported					not supported
	JunosEvpnRoutingInstanceVlanAware     *enum.FeatureSwitchEnum // virtual_network_policy node	not supported					not supported					not supported
	JunosExOverlayEcmp                    *enum.FeatureSwitchEnum // virtual_network_policy node	not supported					not supported					not supported
	JunosGracefulRestart                  *enum.FeatureSwitchEnum // virtual_network_policy node	not supported					not supported					not supported
	MaxEvpnRoutes                         *uint32                 // virtual_network_policy node	virtual_network_policy node		virtual_network_policy node		virtual_network_policy node
	MaxExternalRoutes                     *uint32                 // virtual_network_policy node	virtual_network_policy node		virtual_network_policy node		virtual_network_policy node
	MaxFabricRoutes                       *uint32                 // virtual_network_policy node	virtual_network_policy node		virtual_network_policy node		virtual_network_policy node
	MaxMlagRoutes                         *uint32                 // virtual_network_policy node	virtual_network_policy node		virtual_network_policy node		virtual_network_policy node
	OptimiseSzFootprint                   *enum.FeatureSwitchEnum // security_zone_policy node		not supported					not supported					not supported
	OverlayControlProtocol                *OverlayControlProtocol // virtual_network_policy node	virtual_network_policy node		virtual_network_policy node		virtual_network_policy node
	SpineLeafLinks                        *AddressingScheme       // blueprint creation only		blueprint creation only			blueprint creation only			blueprint creation only
	SpineSuperspineLinks                  *AddressingScheme       // blueprint creation only		blueprint creation only			blueprint creation only			blueprint creation only
}

func (o FabricSettings) raw() *rawFabricSettings {
	var antiAffinityPolicy *rawAntiAffinityPolicy
	if o.AntiAffinityPolicy != nil {
		antiAffinityPolicy = o.AntiAffinityPolicy.raw()
	}

	var jeriType *string
	if o.JunosEvpnRoutingInstanceVlanAware != nil {
		if o.JunosEvpnRoutingInstanceVlanAware.String() == enum.FeatureSwitchEnumEnabled.String() {
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

	return &rawFabricSettings{
		AntiAffinity:                          antiAffinityPolicy,
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

func (o rawFabricSettings) polish() (*FabricSettings, error) {
	var ocp *OverlayControlProtocol
	if o.OverlayControlProtocol != nil {
		ocp = new(OverlayControlProtocol)
		err := ocp.FromString(*o.OverlayControlProtocol)
		if err != nil {
			return nil, fmt.Errorf("failed to parse overlay_control_protocol value %q", *o.OverlayControlProtocol)
		}
	}

	var junosEvpnRoutingInstanceVlanAware *enum.FeatureSwitchEnum
	if o.JunosEvpnRoutingInstanceType != nil {
		x := junosEvpnRoutingInstanceTypes.Parse(*o.JunosEvpnRoutingInstanceType)
		if x == nil {
			return nil, fmt.Errorf("failed to parse junos_evpn_routing_instance_type value %q", *o.JunosEvpnRoutingInstanceType)
		}
		if *x == junosEvpnRoutingInstanceTypeDefault {
			junosEvpnRoutingInstanceVlanAware = &enum.FeatureSwitchEnumDisabled
		} else {
			junosEvpnRoutingInstanceVlanAware = &enum.FeatureSwitchEnumEnabled
		}
	}

	antiAffinityPolicy, err := o.AntiAffinity.polish()
	if err != nil {
		return nil, err
	}

	var spineLeafLinks *AddressingScheme
	if o.SpineLeafLinks != nil {
		i, err := o.SpineLeafLinks.parse()
		if err != nil {
			return nil, err
		}
		spineLeafLinks = toPtr(AddressingScheme(i))
	}

	var spineSuperspineLinks *AddressingScheme
	if o.SpineSuperspineLinks != nil {
		i, err := o.SpineSuperspineLinks.parse()
		if err != nil {
			return nil, err
		}
		spineSuperspineLinks = toPtr(AddressingScheme(i))
	}

	return &FabricSettings{
		AntiAffinityPolicy:                    antiAffinityPolicy,
		DefaultSviL3Mtu:                       o.DefaultSviL3Mtu,
		EsiMacMsb:                             o.EsiMacMsb,
		EvpnGenerateType5HostRoutes:           featureSwitchEnumFromStringPtr(o.EvpnGenerateType5HostRoutes),
		ExternalRouterMtu:                     o.ExternalRouterMtu,
		FabricL3Mtu:                           o.FabricL3Mtu,
		Ipv6Enabled:                           o.Ipv6Enabled,
		JunosEvpnDuplicateMacRecoveryTime:     o.JunosEvpnDuplicateMacRecoveryTime,
		JunosEvpnMaxNexthopAndInterfaceNumber: featureSwitchEnumFromStringPtr(o.JunosEvpnMaxNexthopAndInterfaceNumber),
		JunosEvpnRoutingInstanceVlanAware:     junosEvpnRoutingInstanceVlanAware,
		JunosExOverlayEcmp:                    featureSwitchEnumFromStringPtr(o.JunosExOverlayEcmp),
		JunosGracefulRestart:                  featureSwitchEnumFromStringPtr(o.JunosGracefulRestart),
		MaxEvpnRoutes:                         o.MaxEvpnRoutes,
		MaxExternalRoutes:                     o.MaxExternalRoutes,
		MaxFabricRoutes:                       o.MaxFabricRoutes,
		MaxMlagRoutes:                         o.MaxMlagRoutes,
		OptimiseSzFootprint:                   featureSwitchEnumFromStringPtr(o.OptimiseSzFootprint),
		OverlayControlProtocol:                ocp,
		SpineLeafLinks:                        spineLeafLinks,
		SpineSuperspineLinks:                  spineSuperspineLinks,
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

// getFabricSettings420 does the same job as setFabricSettings, but for Apstra 4.2.0, which collects
// the parameters in rawFabricSettings from 3 different places
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

	antiAffinityPolicy, err := o.getAntiAffinityPolicy(ctx)
	if err != nil {
		return nil, err
	}

	return &rawFabricSettings{
		AntiAffinity:                          antiAffinityPolicy,
		DefaultSviL3Mtu:                       virtualNetworkPolicy.DefaultSviL3Mtu,
		EsiMacMsb:                             fabricAddressingPolicy.EsiMacMsb,
		EvpnGenerateType5HostRoutes:           virtualNetworkPolicy.EvpnGenerateType5HostRoutes,
		ExternalRouterMtu:                     virtualNetworkPolicy.ExternalRouterMtu,
		FabricL3Mtu:                           fabricAddressingPolicy.FabricL3Mtu,
		Ipv6Enabled:                           fabricAddressingPolicy.Ipv6Enabled,
		JunosEvpnDuplicateMacRecoveryTime:     virtualNetworkPolicy.JunosEvpnDuplicateMacRecoveryTime,
		JunosEvpnMaxNexthopAndInterfaceNumber: virtualNetworkPolicy.JunosEvpnMaxNexthopAndInterfaceNumber,
		JunosEvpnRoutingInstanceType:          virtualNetworkPolicy.JunosEvpnRoutingInstanceType,
		JunosExOverlayEcmp:                    virtualNetworkPolicy.JunosExOverlayEcmp,
		JunosGracefulRestart:                  virtualNetworkPolicy.JunosGracefulRestart,
		MaxEvpnRoutes:                         virtualNetworkPolicy.MaxEvpnRoutes,
		MaxExternalRoutes:                     virtualNetworkPolicy.MaxExternalRoutes,
		MaxFabricRoutes:                       virtualNetworkPolicy.MaxFabricRoutes,
		MaxMlagRoutes:                         virtualNetworkPolicy.MaxMlagRoutes,
		OptimiseSzFootprint:                   &optimiseFootprint,
		OverlayControlProtocol:                virtualNetworkPolicy.OverlayControlProtocol,
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
	if err != nil {
		return err
	}

	err = o.setVirtualNetworkPolicy420(ctx, in)
	if err != nil {
		return err
	}

	err = o.setSzFootprintOptimization420(ctx, in.OptimiseSzFootprint)
	if err != nil {
		return err
	}

	err = o.setAntiAffinityPolicy(ctx, in.AntiAffinity)
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
