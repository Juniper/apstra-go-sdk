package apstra

import (
	"fmt"

	oenum "github.com/orsinium-labs/enum"
)

type enum interface {
	String() string
	FromString(string) error
}

var (
	_                  enum = new(DeployMode)
	DeployModeDeploy        = DeployMode{Value: "deploy"}
	DeployModeDrain         = DeployMode{Value: "drain"}
	DeployModeNone          = DeployMode{Value: ""}
	DeployModeReady         = DeployMode{Value: "ready"}
	DeployModeUndeploy      = DeployMode{Value: "undeploy"}
	DeployModes             = oenum.New(
		DeployModeDeploy,
		DeployModeDrain,
		DeployModeNone,
		DeployModeReady,
		DeployModeUndeploy,
	)

	_                           enum = new(DeviceProfileType)
	DeviceProfileTypeModular         = DeviceProfileType{Value: "modular"}
	DeviceProfileTypeMonolithic      = DeviceProfileType{Value: "monolithic"}
	DeviceProfileTypes               = oenum.New(DeviceProfileTypeModular, DeviceProfileTypeMonolithic)

	_                         enum = new(FeatureSwitchEnum)
	FeatureSwitchEnumEnabled       = FeatureSwitchEnum{Value: "enabled"}
	FeatureSwitchEnumDisabled      = FeatureSwitchEnum{Value: "disabled"}
	FeatureSwitchEnums             = oenum.New(
		FeatureSwitchEnumEnabled,
		FeatureSwitchEnumDisabled,
	)

	_                           enum = new(IbaWidgetType)
	IbaWidgetTypeStage               = IbaWidgetType{Value: "stage"}
	IbaWidgetTypeAnomalyHeatmap      = IbaWidgetType{Value: "anomaly_heatmap"}
	IbaWidgetTypes                   = oenum.New(
		IbaWidgetTypeStage,
		IbaWidgetTypeAnomalyHeatmap,
	)

	_                          enum = new(JunosEvpnIrbMode)
	JunosEvpnIrbModeSymmetric       = JunosEvpnIrbMode{Value: "symmetric"}
	JunosEvpnIrbModeAsymmetric      = JunosEvpnIrbMode{Value: "asymmetric"}
	JunosEvpnIrbModes               = oenum.New(
		JunosEvpnIrbModeSymmetric,
		JunosEvpnIrbModeAsymmetric,
	)

	_                                        enum = new(PolicyApplicationPointType)
	PolicyApplicationPointTypeGroup               = PolicyApplicationPointType{Value: "group"}
	PolicyApplicationPointTypeInternal            = PolicyApplicationPointType{Value: "internal"}
	PolicyApplicationPointTypeExternal            = PolicyApplicationPointType{Value: "external"}
	PolicyApplicationPointTypeSecurityZone        = PolicyApplicationPointType{Value: "security_zone"}
	PolicyApplicationPointTypeVirtualNetwork      = PolicyApplicationPointType{Value: "virtual_network"}
	PolicyApplicationPointTypes                   = oenum.New(
		PolicyApplicationPointTypeGroup,
		PolicyApplicationPointTypeInternal,
		PolicyApplicationPointTypeExternal,
		PolicyApplicationPointTypeSecurityZone,
		PolicyApplicationPointTypeVirtualNetwork,
	)

	_                         enum = new(PolicyRuleAction)
	PolicyRuleActionDeny           = PolicyRuleAction{Value: "deny"}
	PolicyRuleActionDenyLog        = PolicyRuleAction{Value: "deny_log"}
	PolicyRuleActionPermit         = PolicyRuleAction{Value: "permit"}
	PolicyRuleActionPermitLog      = PolicyRuleAction{Value: "permit_log"}
	PolicyRuleActions              = oenum.New(
		PolicyRuleActionDeny,
		PolicyRuleActionDenyLog,
		PolicyRuleActionPermit,
		PolicyRuleActionPermitLog,
	)

	_                      enum = new(PolicyRuleProtocol)
	PolicyRuleProtocolIcmp      = PolicyRuleProtocol{Value: "ICMP"}
	PolicyRuleProtocolIp        = PolicyRuleProtocol{Value: "IP"}
	PolicyRuleProtocolTcp       = PolicyRuleProtocol{Value: "TCP"}
	PolicyRuleProtocolUdp       = PolicyRuleProtocol{Value: "UDP"}
	PolicyRuleProtocols         = oenum.New(
		PolicyRuleProtocolIcmp,
		PolicyRuleProtocolIp,
		PolicyRuleProtocolTcp,
		PolicyRuleProtocolUdp,
	)

	_                               enum = new(RemoteGatewayRouteTypes)
	RemoteGatewayRouteTypesAll           = RemoteGatewayRouteTypes{Value: "all"}
	RemoteGatewayRouteTypesFiveOnly      = RemoteGatewayRouteTypes{Value: "type5_only"}
	RemoteGatewayRouteTypesEnum          = oenum.New(
		RemoteGatewayRouteTypesAll,
		RemoteGatewayRouteTypesFiveOnly,
	)

	_                            enum = new(TcpStateQualifier)
	TcpStateQualifierEstablished      = TcpStateQualifier{Value: "established"}
	TcpStateQualifiers                = oenum.New(
		TcpStateQualifierEstablished,
	)

	_                      enum = new(FFResourceType)
	FFResourceTypeAsn           = FFResourceType{Value: "asn"}
	FFResourceTypeHostIpv4      = FFResourceType{Value: "host_ip"}
	FFResourceTypeHostIpv6      = FFResourceType{Value: "host_ipv6"}
	FFResourceTypeInt           = FFResourceType{Value: "integer"}
	FFResourceTypeIpv4          = FFResourceType{Value: "ip"}
	FFResourceTypeIpv6          = FFResourceType{Value: "ipv6"}
	FFResourceTypeVlan          = FFResourceType{Value: "vlan"}
	FFResourceTypeVni           = FFResourceType{Value: "vni"}
	FFResourceTypes             = oenum.New(
		FFResourceTypeAsn,
		FFResourceTypeHostIpv4,
		FFResourceTypeHostIpv6,
		FFResourceTypeInt,
		FFResourceTypeIpv4,
		FFResourceTypeIpv6,
		FFResourceTypeVlan,
		FFResourceTypeVni,
	)

	_                                   enum = new(StorageSchemaPath)
	StorageSchemaPathXCVR                    = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.xcvr"}
	StorageSchemaPathGRAPH                   = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.graph"}
	StorageSchemaPathROUTE                   = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.route"}
	StorageSchemaPathMAC                     = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.mac"}
	StorageSchemaPathOPTICAL_XCVR            = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.optical_xcvr"}
	StorageSchemaPathHOSTNAME                = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.hostname"}
	StorageSchemaPathGENERIC                 = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.generic"}
	StorageSchemaPathLAG                     = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.lag"}
	StorageSchemaPathBGP                     = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.bgp"}
	StorageSchemaPathINTERFACE               = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.interface"}
	StorageSchemaPathMLAG                    = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.mlag"}
	StorageSchemaPathIBA_STRING_DATA         = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.iba_string_data"}
	StorageSchemaPathIBA_INTEGER_DATA        = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.iba_integer_data"}
	StorageSchemaPathROUTE_LOOKUP            = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.route_lookup"}
	StorageSchemaPathINTERFACE_COUNTERS      = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.interface_counters"}
	StorageSchemaPathARP                     = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.arp"}
	StorageSchemaPathCPP_GRAPH               = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.cpp_graph"}
	StorageSchemaPathNSXT                    = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.nsxt"}
	StorageSchemaPathENVIRONMENT             = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.environment"}
	StorageSchemaPathLLDP                    = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.lldp"}
	StorageSchemaPaths                       = oenum.New(StorageSchemaPathXCVR,
		StorageSchemaPathGRAPH,
		StorageSchemaPathROUTE,
		StorageSchemaPathMAC,
		StorageSchemaPathOPTICAL_XCVR,
		StorageSchemaPathHOSTNAME,
		StorageSchemaPathGENERIC,
		StorageSchemaPathLAG,
		StorageSchemaPathBGP,
		StorageSchemaPathINTERFACE,
		StorageSchemaPathMLAG,
		StorageSchemaPathIBA_STRING_DATA,
		StorageSchemaPathIBA_INTEGER_DATA,
		StorageSchemaPathROUTE_LOOKUP,
		StorageSchemaPathINTERFACE_COUNTERS,
		StorageSchemaPathARP,
		StorageSchemaPathCPP_GRAPH,
		StorageSchemaPathNSXT,
		StorageSchemaPathENVIRONMENT,
		StorageSchemaPathLLDP,
	)

	_                   enum = new(FFLinkType)
	FFLinkTypeEthernet       = FFLinkType{Value: "ethernet"}
	FFLinkTypeAggregate      = FFLinkType{Value: "aggregate_link"}
	FFLinkTypes              = oenum.New(
		FFLinkTypeEthernet,
		FFLinkTypeAggregate,
	)

	_                                  enum = new(InterfaceNumberingIpv4Type)
	InterfaceNumberingIpv4TypeNone          = InterfaceNumberingIpv4Type{Value: ""}
	InterfaceNumberingIpv4TypeNumbered      = InterfaceNumberingIpv4Type{Value: "numbered"}
	InterfaceNumberingIpv4Types             = oenum.New(
		InterfaceNumberingIpv4TypeNone,
		InterfaceNumberingIpv4TypeNumbered,
	)

	_                                   enum = new(InterfaceNumberingIpv6Type)
	InterfaceNumberingIpv6TypeNone           = InterfaceNumberingIpv6Type{Value: ""}
	InterfaceNumberingIpv6TypeNumbered       = InterfaceNumberingIpv6Type{Value: "numbered"}
	InterfaceNumberingIpv6TypeLinkLocal      = InterfaceNumberingIpv6Type{Value: "link_local"}
	InterfaceNumberingIpv6Types              = oenum.New(
		InterfaceNumberingIpv6TypeNone,
		InterfaceNumberingIpv6TypeNumbered,
		InterfaceNumberingIpv6TypeLinkLocal,
	)

	_                    enum = new(ResourcePoolType)
	ResourcePoolTypeAsn       = ResourcePoolType{Value: "asn"}
	ResourcePoolTypeInt       = ResourcePoolType{Value: "integer"}
	ResourcePoolTypeIpv4      = ResourcePoolType{Value: "ip"}
	ResourcePoolTypeIpv6      = ResourcePoolType{Value: "ipv6"}
	ResourcePoolTypeVlan      = ResourcePoolType{Value: "vlan"}
	ResourcePoolTypeVni       = ResourcePoolType{Value: "vni"}
	ResourcePoolTypes         = oenum.New(
		ResourcePoolTypeAsn,
		ResourcePoolTypeInt,
		ResourcePoolTypeIpv4,
		ResourcePoolTypeIpv6,
		ResourcePoolTypeVlan,
		ResourcePoolTypeVni,
	)
	_                       enum = new(CollectorOSType)
	CollectorOSTypeJunos         = CollectorOSType{Value: "junos"}
	CollectorOSTypeJunosEvo      = CollectorOSType{Value: "junos_evo"}
	CollectorOSTypes             = oenum.New(
		CollectorOSTypeJunos,
		CollectorOSTypeJunosEvo,
	)

	_                             enum = new(CollectorOSFamily)
	CollectorOSFamilyACX               = CollectorOSFamily{Value: "acx"}
	CollectorOSFamilyACX_F             = CollectorOSFamily{Value: "acx-f"}
	CollectorOSFamilyACX_QFX_7K        = CollectorOSFamily{Value: "acx-qfx-7k"}
	CollectorOSFamilyPTX               = CollectorOSFamily{Value: "ptx"}
	CollectorOSFamilyPTX1K             = CollectorOSFamily{Value: "ptx1k"}
	CollectorOSFamilyQFX_MS_FIXED      = CollectorOSFamily{Value: "qfx-ms-fixed"}
	CollectorOSFamilyJunos             = CollectorOSFamily{Value: "junos"}
	CollectorOSFamilyJunos_EX          = CollectorOSFamily{Value: "junos-ex"}
	CollectorOSFamilyJunos_QFX         = CollectorOSFamily{Value: "junos-qfx"}
	CollectorOSFamilies                = oenum.New(
		CollectorOSFamilyACX,
		CollectorOSFamilyACX_F,
		CollectorOSFamilyACX_QFX_7K,
		CollectorOSFamilyPTX,
		CollectorOSFamilyPTX1K,
		CollectorOSFamilyQFX_MS_FIXED,
		CollectorOSFamilyJunos,
		CollectorOSFamilyJunos_EX,
		CollectorOSFamilyJunos_QFX,
	)

	CollectorSourceTypeCLI = CollectorSourceType{Value: "cli"}
	CollectorSourceTypes   = oenum.New(
		CollectorSourceTypeCLI,
	)

	_                              enum = new(RoutingZoneConstraintMode)
	RoutingZoneConstraintModeNone       = RoutingZoneConstraintMode{Value: "none"}
	RoutingZoneConstraintModeAllow      = RoutingZoneConstraintMode{Value: "allow"}
	RoutingZoneConstraintModeDeny       = RoutingZoneConstraintMode{Value: "deny"}
	RoutingZoneConstraintModes          = oenum.New(
		RoutingZoneConstraintModeNone,
		RoutingZoneConstraintModeAllow,
		RoutingZoneConstraintModeDeny,
	)
)

type DeployMode oenum.Member[string]

func (o DeployMode) String() string {
	return o.Value
}

func (o *DeployMode) FromString(s string) error {
	if DeployModes.Parse(s) == nil {
		return fmt.Errorf("failed to parse DeployMode value %q", s)
	}

	o.Value = s
	return nil
}

type DeviceProfileType oenum.Member[string]

func (o DeviceProfileType) String() string {
	return o.Value
}

func (o *DeviceProfileType) FromString(s string) error {
	t := DeviceProfileTypes.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse DeviceProfileType %q", s)
	}
	o.Value = t.Value
	return nil
}

type FeatureSwitchEnum oenum.Member[string]

func (o FeatureSwitchEnum) String() string {
	return o.Value
}

func (o *FeatureSwitchEnum) FromString(s string) error {
	t := FeatureSwitchEnums.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse FeatureSwitchEnum %q", s)
	}
	o.Value = t.Value
	return nil
}

type IbaWidgetType oenum.Member[string]

func (o IbaWidgetType) String() string {
	return o.Value
}

func (o *IbaWidgetType) FromString(s string) error {
	t := IbaWidgetTypes.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse IbaWidgetTypes %q", s)
	}
	o.Value = t.Value
	return nil
}

type JunosEvpnIrbMode oenum.Member[string]

func (o JunosEvpnIrbMode) String() string {
	return o.Value
}

func (o *JunosEvpnIrbMode) FromString(s string) error {
	t := JunosEvpnIrbModes.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse JunosEvpnIrbMode %q", s)
	}
	o.Value = t.Value
	return nil
}

type PolicyApplicationPointType oenum.Member[string]

func (o PolicyApplicationPointType) String() string {
	return o.Value
}

func (o *PolicyApplicationPointType) FromString(s string) error {
	t := PolicyApplicationPointTypes.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse PolicyApplicationPointType %q", s)
	}
	o.Value = t.Value
	return nil
}

type PolicyRuleAction oenum.Member[string]

func (o PolicyRuleAction) String() string {
	return o.Value
}

func (o *PolicyRuleAction) FromString(s string) error {
	t := PolicyRuleActions.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse PolicyRuleAction %q", s)
	}
	o.Value = t.Value
	return nil
}

type PolicyRuleProtocol oenum.Member[string]

func (o PolicyRuleProtocol) String() string {
	return o.Value
}

func (o *PolicyRuleProtocol) FromString(s string) error {
	t := PolicyRuleProtocols.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse PolicyRuleProtocol %q", s)
	}
	o.Value = t.Value
	return nil
}

type RemoteGatewayRouteTypes oenum.Member[string]

func (o RemoteGatewayRouteTypes) String() string {
	return o.Value
}

func (o *RemoteGatewayRouteTypes) FromString(s string) error {
	t := RemoteGatewayRouteTypesEnum.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse RemoteGatewayRouteTypes %q", s)
	}
	o.Value = t.Value
	return nil
}

type TcpStateQualifier oenum.Member[string]

func (o TcpStateQualifier) String() string {
	return o.Value
}

func (o *TcpStateQualifier) FromString(s string) error {
	t := TcpStateQualifiers.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse TcpStateQualifier %q", s)
	}
	o.Value = t.Value
	return nil
}

type FFResourceType oenum.Member[string]

func (o FFResourceType) String() string {
	return o.Value
}

func (o *FFResourceType) FromString(s string) error {
	t := FFResourceTypes.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse FFResourceType %q", s)
	}
	o.Value = t.Value
	return nil
}

type StorageSchemaPath oenum.Member[string]

func (o StorageSchemaPath) String() string {
	return o.Value
}

func (o *StorageSchemaPath) FromString(s string) error {
	t := StorageSchemaPaths.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse StorageSchemaPath %q", s)
	}
	o.Value = t.Value
	return nil
}

type FFLinkType oenum.Member[string]

func (o FFLinkType) String() string {
	return o.Value
}

func (o *FFLinkType) FromString(s string) error {
	t := FFLinkTypes.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse FFLinkType %q", s)
	}
	o.Value = t.Value
	return nil
}

type InterfaceNumberingIpv4Type oenum.Member[string]

func (o InterfaceNumberingIpv4Type) String() string {
	return o.Value
}

func (o *InterfaceNumberingIpv4Type) FromString(s string) error {
	t := InterfaceNumberingIpv4Types.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse InterfaceNumberingIpv4Type %q", s)
	}
	o.Value = t.Value
	return nil
}

type InterfaceNumberingIpv6Type oenum.Member[string]

func (o InterfaceNumberingIpv6Type) String() string {
	return o.Value
}

func (o *InterfaceNumberingIpv6Type) FromString(s string) error {
	t := InterfaceNumberingIpv6Types.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse InterfaceNumberingIpv6Type %q", s)
	}
	o.Value = t.Value
	return nil
}

type ResourcePoolType oenum.Member[string]

func (o ResourcePoolType) String() string {
	return o.Value
}

func (o *ResourcePoolType) FromString(s string) error {
	t := ResourcePoolTypes.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse ResourcePoolType %q", s)
	}
	o.Value = t.Value
	return nil
}

type CollectorOSType oenum.Member[string]

func (o CollectorOSType) String() string { return o.Value }

func (o *CollectorOSType) FromString(s string) error {
	t := CollectorOSTypes.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse CollectorOSType %q", s)
	}
	o.Value = t.Value
	return nil
}

type CollectorOSFamily oenum.Member[string]

func (o CollectorOSFamily) String() string { return o.Value }

func (o *CollectorOSFamily) FromString(s string) error {
	t := CollectorOSFamilies.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse CollectorOSFamily %q", s)
	}
	o.Value = t.Value
	return nil
}

type CollectorSourceType oenum.Member[string]

func (o CollectorSourceType) String() string { return o.Value }

func (o *CollectorSourceType) FromString(s string) error {
	t := CollectorSourceTypes.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse CollectorOSFamily %q", s)
	}
	o.Value = t.Value
	return nil
}

type RoutingZoneConstraintMode oenum.Member[string]

func (o RoutingZoneConstraintMode) String() string {
	return o.Value
}

func (o *RoutingZoneConstraintMode) FromString(s string) error {
	t := RoutingZoneConstraintModes.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse RoutingZoneConstraintMode %q", s)
	}
	return nil
}
