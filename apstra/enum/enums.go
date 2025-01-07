// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:generate go run generator/generator.go
//go:generate go run mvdan.cc/gofumpt -w generated_enums.go

package enum

import oenum "github.com/orsinium-labs/enum"

// Attention! This file must contain only `type` and `var()` declarations of the
// sort included below. The `generated_enums.go` file is auto-generated based on
// the contents of this file. After editing this file, run `go generate ./...`
// or `go generate apstra/enum/enums.go` from the repository root directory to
// refresh `generated_enums.go`.

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
