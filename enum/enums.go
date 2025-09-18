// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:generate go run generator/generator.go
//go:generate go run mvdan.cc/gofumpt -w generated_enums.go

package enum

import oenum "github.com/orsinium-labs/enum"

// Attention! This file must contain only `type` and `var()` declarations of the
// sort included below. The `generated_enums.go` file is auto-generated based on
// the contents of this file. After editing this file, run `go generate ./...`
// or `go generate enum/enums.go` from the repository root directory to refresh
// `generated_enums.go`.

type DeployMode oenum.Member[string]

var (
	DeployModeDeploy   = DeployMode{Value: "deploy"}
	DeployModeDrain    = DeployMode{Value: "drain"}
	DeployModeNone     = DeployMode{Value: ""}
	DeployModeReady    = DeployMode{Value: "ready"}
	DeployModeUndeploy = DeployMode{Value: "undeploy"}
)

type DeviceProfileType oenum.Member[string]

var (
	DeviceProfileTypeModular    = DeviceProfileType{Value: "modular"}
	DeviceProfileTypeMonolithic = DeviceProfileType{Value: "monolithic"}
)

type FeatureSwitch oenum.Member[string]

var (
	FeatureSwitchEnabled  = FeatureSwitch{Value: "enabled"}
	FeatureSwitchDisabled = FeatureSwitch{Value: "disabled"}
)

type IbaWidgetType oenum.Member[string]

var (
	IbaWidgetTypeStage          = IbaWidgetType{Value: "stage"}
	IbaWidgetTypeAnomalyHeatmap = IbaWidgetType{Value: "anomaly_heatmap"}
)

type JunosEvpnIrbMode oenum.Member[string]

var (
	JunosEvpnIrbModeSymmetric  = JunosEvpnIrbMode{Value: "symmetric"}
	JunosEvpnIrbModeAsymmetric = JunosEvpnIrbMode{Value: "asymmetric"}
)

type PolicyApplicationPointType oenum.Member[string]

var (
	PolicyApplicationPointTypeGroup          = PolicyApplicationPointType{Value: "group"}
	PolicyApplicationPointTypeInternal       = PolicyApplicationPointType{Value: "internal"}
	PolicyApplicationPointTypeExternal       = PolicyApplicationPointType{Value: "external"}
	PolicyApplicationPointTypeSecurityZone   = PolicyApplicationPointType{Value: "security_zone"}
	PolicyApplicationPointTypeVirtualNetwork = PolicyApplicationPointType{Value: "virtual_network"}
)

type PolicyRuleAction oenum.Member[string]

var (
	PolicyRuleActionDeny      = PolicyRuleAction{Value: "deny"}
	PolicyRuleActionDenyLog   = PolicyRuleAction{Value: "deny_log"}
	PolicyRuleActionPermit    = PolicyRuleAction{Value: "permit"}
	PolicyRuleActionPermitLog = PolicyRuleAction{Value: "permit_log"}
)

type PolicyRuleProtocol oenum.Member[string]

var (
	PolicyRuleProtocolIcmp = PolicyRuleProtocol{Value: "ICMP"}
	PolicyRuleProtocolIp   = PolicyRuleProtocol{Value: "IP"}
	PolicyRuleProtocolTcp  = PolicyRuleProtocol{Value: "TCP"}
	PolicyRuleProtocolUdp  = PolicyRuleProtocol{Value: "UDP"}
)

type RemoteGatewayRouteType oenum.Member[string]

var (
	RemoteGatewayRouteTypeAll      = RemoteGatewayRouteType{Value: "all"}
	RemoteGatewayRouteTypeFiveOnly = RemoteGatewayRouteType{Value: "type5_only"}
)

type TcpStateQualifier oenum.Member[string]

var TcpStateQualifierEstablished = TcpStateQualifier{Value: "established"}

type FFResourceType oenum.Member[string]

var (
	FFResourceTypeAsn      = FFResourceType{Value: "asn"}
	FFResourceTypeHostIpv4 = FFResourceType{Value: "host_ip"}
	FFResourceTypeHostIpv6 = FFResourceType{Value: "host_ipv6"}
	FFResourceTypeInt      = FFResourceType{Value: "integer"}
	FFResourceTypeIpv4     = FFResourceType{Value: "ip"}
	FFResourceTypeIpv6     = FFResourceType{Value: "ipv6"}
	FFResourceTypeVlan     = FFResourceType{Value: "vlan"}
	FFResourceTypeVni      = FFResourceType{Value: "vni"}
)

type StorageSchemaPath oenum.Member[string]

var (
	StorageSchemaPathARP               = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.arp"}
	StorageSchemaPathBGP               = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.bgp"}
	StorageSchemaPathCppGraph          = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.cpp_graph"}
	StorageSchemaPathEnvironment       = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.environment"}
	StorageSchemaPathGeneric           = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.generic"}
	StorageSchemaPathGraph             = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.graph"}
	StorageSchemaPathHostname          = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.hostname"}
	StorageSchemaPathIbaData           = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.iba_data"}
	StorageSchemaPathIbaIntegerData    = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.iba_integer_data"}
	StorageSchemaPathIbaStringData     = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.iba_string_data"}
	StorageSchemaPathInterface         = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.interface"}
	StorageSchemaPathInterfaceCounters = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.interface_counters"}
	StorageSchemaPathLAG               = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.lag"}
	StorageSchemaPathLLDP              = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.lldp"}
	StorageSchemaPathMAC               = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.mac"}
	StorageSchemaPathMLAG              = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.mlag"}
	StorageSchemaPathNSXT              = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.nsxt"}
	StorageSchemaPathOpticalXcvr       = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.optical_xcvr"}
	StorageSchemaPathRoute             = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.route"}
	StorageSchemaPathRouteLookup       = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.route_lookup"}
	StorageSchemaPathXcvr              = StorageSchemaPath{Value: "aos.sdk.telemetry.schemas.xcvr"}
)

type InterfaceNumberingIpv4Type oenum.Member[string]

var (
	InterfaceNumberingIpv4TypeNone     = InterfaceNumberingIpv4Type{Value: ""}
	InterfaceNumberingIpv4TypeNumbered = InterfaceNumberingIpv4Type{Value: "numbered"}
)

type InterfaceNumberingIpv6Type oenum.Member[string]

var (
	InterfaceNumberingIpv6TypeNone      = InterfaceNumberingIpv6Type{Value: ""}
	InterfaceNumberingIpv6TypeNumbered  = InterfaceNumberingIpv6Type{Value: "numbered"}
	InterfaceNumberingIpv6TypeLinkLocal = InterfaceNumberingIpv6Type{Value: "link_local"}
)

type ResourcePoolType oenum.Member[string]

var (
	ResourcePoolTypeAsn  = ResourcePoolType{Value: "asn"}
	ResourcePoolTypeInt  = ResourcePoolType{Value: "integer"}
	ResourcePoolTypeIpv4 = ResourcePoolType{Value: "ip"}
	ResourcePoolTypeIpv6 = ResourcePoolType{Value: "ipv6"}
	ResourcePoolTypeVlan = ResourcePoolType{Value: "vlan"}
	ResourcePoolTypeVni  = ResourcePoolType{Value: "vni"}
)

type RoutingZoneConstraintMode oenum.Member[string]

var (
	RoutingZoneConstraintModeNone  = RoutingZoneConstraintMode{Value: "none"}
	RoutingZoneConstraintModeAllow = RoutingZoneConstraintMode{Value: "allow"}
	RoutingZoneConstraintModeDeny  = RoutingZoneConstraintMode{Value: "deny"}
)

type ApiFeature oenum.Member[string]

var (
	ApiFeatureAiFabric   = ApiFeature{Value: "ai_fabric"}
	ApiFeatureCentral    = ApiFeature{Value: "central"}
	ApiFeatureEnterprise = ApiFeature{Value: "enterprise"}
	ApiFeatureFreeform   = ApiFeature{Value: "freeform"}
	ApiFeatureFullAccess = ApiFeature{Value: "full_access"}
	ApiFeatureTaskApi    = ApiFeature{Value: "task_api"}
)

type PortRole oenum.Member[string]

var (
	PortRoleAccess     = PortRole{Value: "access"}
	PortRoleGeneric    = PortRole{Value: "generic"}
	PortRoleL3Server   = PortRole{Value: "l3_server"} // todo: remove this, LogicalDevicePortRoles.Validate and simplify LogicalDevicePortRoles.IncludeAllUses
	PortRoleLeaf       = PortRole{Value: "leaf"}
	PortRolePeer       = PortRole{Value: "peer"}
	PortRoleSpine      = PortRole{Value: "spine"}
	PortRoleSuperspine = PortRole{Value: "superspine"}
	PortRoleUnused     = PortRole{Value: "unused"}
)

type RenderedConfigType oenum.Member[string]

var (
	RenderedConfigTypeStaging  = RenderedConfigType{Value: "staging"}
	RenderedConfigTypeDeployed = RenderedConfigType{Value: "deployed"}
)

type SviIpv4Mode oenum.Member[string]

var (
	SviIpv4ModeDisabled = SviIpv4Mode{Value: "disabled"}
	SviIpv4ModeEnabled  = SviIpv4Mode{Value: "enabled"}
	SviIpv4ModeForced   = SviIpv4Mode{Value: "forced"}
)

type SviIpv6Mode oenum.Member[string]

var (
	SviIpv6ModeDisabled  = SviIpv6Mode{Value: "disabled"}
	SviIpv6ModeEnabled   = SviIpv6Mode{Value: "enabled"}
	SviIpv6ModeForced    = SviIpv6Mode{Value: "forced"}
	SviIpv6ModeLinkLocal = SviIpv6Mode{Value: "link_local"}
)

type DhcpServiceMode oenum.Member[string]

var (
	DhcpServiceModeDisabled = DhcpServiceMode{Value: "dhcpServiceDisabled"}
	DhcpServiceModeEnabled  = DhcpServiceMode{Value: "dhcpServiceEnabled"}
)

type VnType oenum.Member[string]

var (
	VnTypeExternal = VnType{Value: "disabled"}
	VnTypeVlan     = VnType{Value: "vlan"}
	VnTypeVxlan    = VnType{Value: "vxlan"}
)

type RedundancyGroupType oenum.Member[string]

var (
	RedundancyGroupTypeEsi  = RedundancyGroupType{Value: "esi"}
	RedundancyGroupTypeMlag = RedundancyGroupType{Value: "mlag"}
)

type SystemType oenum.Member[string]

var (
	SystemTypeServer = SystemType{Value: "server"}
	SystemTypeSwitch = SystemType{Value: "switch"}
)

type NodeRole oenum.Member[string]

var (
	NodeRoleAccess        = NodeRole{Value: "access"}
	NodeRoleGeneric       = NodeRole{Value: "generic"}
	NodeRoleLeaf          = NodeRole{Value: "leaf"}
	NodeRoleRemoteGateway = NodeRole{Value: "remote_gateway"}
	NodeRoleSpine         = NodeRole{Value: "spine"}
	NodeRoleSuperspine    = NodeRole{Value: "superspine"}
)

type RefDesign oenum.Member[string]

var (
	RefDesignDatacenter    = RefDesign{Value: "two_stage_l3clos"}
	RefDesignFreeform      = RefDesign{Value: "freeform"}
	RefDesignRailCollapsed = RefDesign{Value: "rail_collapsed"}
)

type ConfigletStyle oenum.Member[string]

var (
	ConfigletStyleCumulus = ConfigletStyle{Value: "cumulus"}
	ConfigletStyleEos     = ConfigletStyle{Value: "eos"}
	ConfigletStyleJunos   = ConfigletStyle{Value: "junos"}
	ConfigletStyleNxos    = ConfigletStyle{Value: "nxos"}
	ConfigletStyleSonic   = ConfigletStyle{Value: "sonic"}
)

type ConfigletSection oenum.Member[string]

var (
	ConfigletSectionDeleteBasedInterface = ConfigletSection{Value: "delete_based_interface"}
	ConfigletSectionFile                 = ConfigletSection{Value: "file"}
	ConfigletSectionFrr                  = ConfigletSection{Value: "frr"}
	ConfigletSectionInterface            = ConfigletSection{Value: "interface"}
	ConfigletSectionOspf                 = ConfigletSection{Value: "ospf"}
	ConfigletSectionSetBasedInterface    = ConfigletSection{Value: "set_based_interface"}
	ConfigletSectionSetBasedSystem       = ConfigletSection{Value: "set_based_system"}
	ConfigletSectionSystem               = ConfigletSection{Value: "system"}
	ConfigletSectionSystemTop            = ConfigletSection{Value: "system_top"}
)

type SecurityZoneType oenum.Member[string]

var (
	SecurityZoneTypeEvpn            = SecurityZoneType{Value: "evpn"}
	SecurityZoneTypeL3Fabric        = SecurityZoneType{Value: "l3_fabric"}
	SecurityZoneTypeVirtualL3Fabric = SecurityZoneType{Value: "virtual_l3_fabric"}
)

type LockStatus oenum.Member[string]

var (
	LockStatusLocked                 = LockStatus{Value: "locked"}
	LockStatusLockedByAdmin          = LockStatus{Value: "locked_by_admin"}
	LockStatusLockedByDeletedUser    = LockStatus{Value: "locked_by_deleted_user"}
	LockStatusLockedByRestrictedUser = LockStatus{Value: "locked_by_restricted_user"}
	LockStatusUnlocked               = LockStatus{Value: "unlocked"}
)

type LockType oenum.Member[string]

var (
	LockTypeLockedByChanges = LockType{Value: "lock_by_changes"}
	LockTypeLockedByUser    = LockType{Value: "lock_by_user"}
	LockTypeUnlocked        = LockType{Value: "unlocked"}
)

type IbaWidgetDataSource oenum.Member[string]

var (
	IbaWidgetDataSourceRealTime   = IbaWidgetDataSource{Value: "real_time"}
	IbaWidgetDataSourceTimeSeries = IbaWidgetDataSource{Value: "time_series"}
)

type IbaWidgetAggregationType oenum.Member[string]

var (
	IbaWidgetAggregationTypeUnset   = IbaWidgetAggregationType{Value: "unset"}
	IbaWidgetAggregationTypeMin     = IbaWidgetAggregationType{Value: "min"}
	IbaWidgetAggregationTypeAverage = IbaWidgetAggregationType{Value: "average"}
	IbaWidgetAggregationTypeNone    = IbaWidgetAggregationType{Value: "none"}
	IbaWidgetAggregationTypeAnyOf   = IbaWidgetAggregationType{Value: "any_of"}
	IbaWidgetAggregationTypeLast    = IbaWidgetAggregationType{Value: "last"}
	IbaWidgetAggregationTypeAllOf   = IbaWidgetAggregationType{Value: "all_of"}
	IbaWidgetAggregationTypeMax     = IbaWidgetAggregationType{Value: "max"}
)

type IbaWidgetCombineGraph oenum.Member[string]

var (
	IbaWidgetCombineGraphNone    = IbaWidgetCombineGraph{Value: "none"}
	IbaWidgetCombineGraphLinear  = IbaWidgetCombineGraph{Value: "linear"}
	IbaWidgetCombineGraphStacked = IbaWidgetCombineGraph{Value: "stacked"}
)

type FabricConnectivityDesign oenum.Member[string]

var (
	FabricConnectivityDesignL3Clos        = FabricConnectivityDesign{Value: "l3clos"}
	FabricConnectivityDesignL3Collapsed   = FabricConnectivityDesign{Value: "l3collapsed"}
	FabricConnectivityDesignRailCollapsed = FabricConnectivityDesign{Value: "rail_collapsed"}
)

type EndpointPolicyStatus oenum.Member[string]

var (
	EndpointPolicyStatusAssigned   = EndpointPolicyStatus{Value: "assigned"}
	EndpointPolicyStatusIncomplete = EndpointPolicyStatus{Value: "incomplete"}
	EndpointPolicyStatusReady      = EndpointPolicyStatus{Value: "ready"}
)

type SpeedUnit oenum.Member[string]

var (
	SpeedUnitM = SpeedUnit{Value: "M"}
	SpeedUnitG = SpeedUnit{Value: "G"}
	// SpeedUnitT = SpeedUnit{Value: "T"}
)

type LinkSpeed oenum.Member[string]

var (
	LinkSpeed10M   = LinkSpeed{Value: "10M"}
	LinkSpeed10m   = LinkSpeed{Value: "10m"}
	LinkSpeed100M  = LinkSpeed{Value: "100M"}
	LinkSpeed100m  = LinkSpeed{Value: "100m"}
	LinkSpeed1G    = LinkSpeed{Value: "1G"}
	LinkSpeed1g    = LinkSpeed{Value: "1g"}
	LinkSpeed2500M = LinkSpeed{Value: "2500M"}
	LinkSpeed2500m = LinkSpeed{Value: "2500m"}
	LinkSpeed5G    = LinkSpeed{Value: "5G"}
	LinkSpeed5g    = LinkSpeed{Value: "5g"}
	LinkSpeed10G   = LinkSpeed{Value: "10G"}
	LinkSpeed10g   = LinkSpeed{Value: "10g"}
	LinkSpeed25G   = LinkSpeed{Value: "25G"}
	LinkSpeed25g   = LinkSpeed{Value: "25g"}
	LinkSpeed40G   = LinkSpeed{Value: "40G"}
	LinkSpeed40g   = LinkSpeed{Value: "40g"}
	LinkSpeed50G   = LinkSpeed{Value: "50G"}
	LinkSpeed50g   = LinkSpeed{Value: "50g"}
	LinkSpeed100G  = LinkSpeed{Value: "100G"}
	LinkSpeed100g  = LinkSpeed{Value: "100g"}
	LinkSpeed150G  = LinkSpeed{Value: "150G"}
	LinkSpeed150g  = LinkSpeed{Value: "150g"}
	LinkSpeed200G  = LinkSpeed{Value: "200G"}
	LinkSpeed200g  = LinkSpeed{Value: "200g"}
	LinkSpeed400G  = LinkSpeed{Value: "400G"}
	LinkSpeed400g  = LinkSpeed{Value: "400g"}
	LinkSpeed800G  = LinkSpeed{Value: "800G"}
	LinkSpeed800g  = LinkSpeed{Value: "800g"}
)
