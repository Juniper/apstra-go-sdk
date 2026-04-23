// Copyright (c) Juniper Networks, Inc., 2024-2026.
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

type ASNAllocationScheme oenum.Member[string]

var (
	ASNAllocationSchemeDistinct = ASNAllocationScheme{Value: "distinct"}
	ASNAllocationSchemeSingle   = ASNAllocationScheme{Value: "single"}
)

type AddressingScheme oenum.Member[string]

var (
	AddressingSchemeIPv4  = AddressingScheme{Value: "ipv4"}
	AddressingSchemeIPv46 = AddressingScheme{Value: "ipv4_ipv6"}
	AddressingSchemeIPv6  = AddressingScheme{Value: "ipv6"}
)

type AntiAffinityMode oenum.Member[string]

var (
	AntiAffinityModeDisabled = AntiAffinityMode{Value: "disabled"}
	AntiAffinityModeLoose    = AntiAffinityMode{Value: "enabled_loose"}
	AntiAffinityModeStrict   = AntiAffinityMode{Value: "enabled_strict"}
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

type ConfigletStyle oenum.Member[string]

var (
	ConfigletStyleCumulus = ConfigletStyle{Value: "cumulus"}
	ConfigletStyleEos     = ConfigletStyle{Value: "eos"}
	ConfigletStyleJunos   = ConfigletStyle{Value: "junos"}
	ConfigletStyleNxos    = ConfigletStyle{Value: "nxos"}
	ConfigletStyleSonic   = ConfigletStyle{Value: "sonic"}
)

type DeployMode oenum.Member[string]

var (
	DeployModeDeploy   = DeployMode{Value: "deploy"}
	DeployModeDrain    = DeployMode{Value: "drain"}
	DeployModeNone     = DeployMode{Value: ""}
	DeployModeReady    = DeployMode{Value: "ready"}
	DeployModeUndeploy = DeployMode{Value: "undeploy"}
)

type DesignLogicalDevicePanelPortIndexing oenum.Member[string]

var (
	DesignLogicalDevicePanelPortIndexingLRTB = DesignLogicalDevicePanelPortIndexing{Value: "L-R, T-B"}
	DesignLogicalDevicePanelPortIndexingTBLR = DesignLogicalDevicePanelPortIndexing{Value: "T-B, L-R"}
)

type DeviceProfileType oenum.Member[string]

var (
	DeviceProfileTypeModular    = DeviceProfileType{Value: "modular"}
	DeviceProfileTypeMonolithic = DeviceProfileType{Value: "monolithic"}
)

type DhcpServiceMode oenum.Member[string]

var (
	DhcpServiceModeDisabled = DhcpServiceMode{Value: "dhcpServiceDisabled"}
	DhcpServiceModeEnabled  = DhcpServiceMode{Value: "dhcpServiceEnabled"}
)

type EndpointPolicyStatus oenum.Member[string]

var (
	EndpointPolicyStatusAssigned   = EndpointPolicyStatus{Value: "assigned"}
	EndpointPolicyStatusIncomplete = EndpointPolicyStatus{Value: "incomplete"}
	EndpointPolicyStatusReady      = EndpointPolicyStatus{Value: "ready"}
)

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

type FabricConnectivityDesign oenum.Member[string]

var (
	FabricConnectivityDesignL3Clos        = FabricConnectivityDesign{Value: "l3clos"}
	FabricConnectivityDesignL3Collapsed   = FabricConnectivityDesign{Value: "l3collapsed"}
	FabricConnectivityDesignRailCollapsed = FabricConnectivityDesign{Value: "rail_collapsed"}
)

type FeatureSwitch oenum.Member[string]

var (
	FeatureSwitchDisabled = FeatureSwitch{Value: "disabled"}
	FeatureSwitchEnabled  = FeatureSwitch{Value: "enabled"}
)

type IbaWidgetAggregationType oenum.Member[string]

var (
	IbaWidgetAggregationTypeAllOf   = IbaWidgetAggregationType{Value: "all_of"}
	IbaWidgetAggregationTypeAnyOf   = IbaWidgetAggregationType{Value: "any_of"}
	IbaWidgetAggregationTypeAverage = IbaWidgetAggregationType{Value: "average"}
	IbaWidgetAggregationTypeLast    = IbaWidgetAggregationType{Value: "last"}
	IbaWidgetAggregationTypeMax     = IbaWidgetAggregationType{Value: "max"}
	IbaWidgetAggregationTypeMin     = IbaWidgetAggregationType{Value: "min"}
	IbaWidgetAggregationTypeNone    = IbaWidgetAggregationType{Value: "none"}
	IbaWidgetAggregationTypeUnset   = IbaWidgetAggregationType{Value: "unset"}
)

type IbaWidgetCombineGraph oenum.Member[string]

var (
	IbaWidgetCombineGraphLinear  = IbaWidgetCombineGraph{Value: "linear"}
	IbaWidgetCombineGraphNone    = IbaWidgetCombineGraph{Value: "none"}
	IbaWidgetCombineGraphStacked = IbaWidgetCombineGraph{Value: "stacked"}
)

type IbaWidgetDataSource oenum.Member[string]

var (
	IbaWidgetDataSourceRealTime   = IbaWidgetDataSource{Value: "real_time"}
	IbaWidgetDataSourceTimeSeries = IbaWidgetDataSource{Value: "time_series"}
)

type IbaWidgetType oenum.Member[string]

var (
	IbaWidgetTypeAnomalyHeatmap = IbaWidgetType{Value: "anomaly_heatmap"}
	IbaWidgetTypeStage          = IbaWidgetType{Value: "stage"}
)

type InterfaceMapInterfaceState oenum.Member[string]

var (
	InterfaceMapInterfaceStateActive   = InterfaceMapInterfaceState{Value: "active"}
	InterfaceMapInterfaceStateInactive = InterfaceMapInterfaceState{Value: "inactive"}
)

type InterfaceNumberingIpv4Type oenum.Member[string]

var (
	InterfaceNumberingIpv4TypeNone     = InterfaceNumberingIpv4Type{Value: ""}
	InterfaceNumberingIpv4TypeNumbered = InterfaceNumberingIpv4Type{Value: "numbered"}
)

type InterfaceNumberingIpv6Type oenum.Member[string]

var (
	InterfaceNumberingIpv6TypeLinkLocal = InterfaceNumberingIpv6Type{Value: "link_local"}
	InterfaceNumberingIpv6TypeNone      = InterfaceNumberingIpv6Type{Value: ""}
	InterfaceNumberingIpv6TypeNumbered  = InterfaceNumberingIpv6Type{Value: "numbered"}
)

type InterfaceOperationState oenum.Member[string]

var (
	InterfaceOperationStateAdminDown = InterfaceOperationState{Value: "admin_down"}
	InterfaceOperationStateDown      = InterfaceOperationState{Value: "deduced_down"}
	InterfaceOperationStateUp        = InterfaceOperationState{Value: "up"}
)

type InterfaceState oenum.Member[string]

var (
	InterfaceStateActive   = InterfaceState{Value: "active"}
	InterfaceStateInactive = InterfaceState{Value: "inactive"}
)

type InterfaceType oenum.Member[string]

var (
	InterfaceTypeAnycastVtep       = InterfaceType{Value: "anycast_vtep"}
	InterfaceTypeEthernet          = InterfaceType{Value: "ethernet"}
	InterfaceTypeGlobalAnycastVtep = InterfaceType{Value: "global_anycast_vtep"}
	InterfaceTypeIp                = InterfaceType{Value: "ip"}
	InterfaceTypeLogicalVtep       = InterfaceType{Value: "logical_vtep"}
	InterfaceTypeLoopback          = InterfaceType{Value: "loopback"}
	InterfaceTypePortChannel       = InterfaceType{Value: "port_channel"}
	InterfaceTypeSubinterface      = InterfaceType{Value: "subinterface"}
	InterfaceTypeSvi               = InterfaceType{Value: "svi"}
	InterfaceTypeUnicastVtep       = InterfaceType{Value: "unicast_vtep"}
)

type JunosEVPNIRBMode oenum.Member[string]

var (
	JunosEVPNIRBModeAsymmetric = JunosEVPNIRBMode{Value: "asymmetric"}
	JunosEVPNIRBModeSymmetric  = JunosEVPNIRBMode{Value: "symmetric"}
)

type LAGMode oenum.Member[string]

var (
	LAGModeActiveLACP  = LAGMode{Value: "lacp_active"}
	LAGModeNone        = LAGMode{Value: ""}
	LAGModePassiveLACP = LAGMode{Value: "lacp_passive"}
	LAGModeStatic      = LAGMode{Value: "static_lag"}
)

type LeafRedundancyProtocol oenum.Member[string]

var (
	LeafRedundancyProtocolESI  = LeafRedundancyProtocol{Value: "esi"}
	LeafRedundancyProtocolMLAG = LeafRedundancyProtocol{Value: "mlag"}
	LeafRedundancyProtocolNone = LeafRedundancyProtocol{Value: ""}
)

type LinkAttachmentType oenum.Member[string]

var (
	LinkAttachmentTypeDual   = LinkAttachmentType{Value: "dualAttached"}
	LinkAttachmentTypeSingle = LinkAttachmentType{Value: "singleAttached"}
)

type LinkRole oenum.Member[string]

var (
	LinkRoleAccessL3PeerLink   = LinkRole{Value: "access_l3_peer_link"}
	LinkRoleAccessServer       = LinkRole{Value: "access_server"}
	LinkRoleLeafAccess         = LinkRole{Value: "leaf_access"}
	LinkRoleLeafL2Server       = LinkRole{Value: "leaf_l2_server"}
	LinkRoleLeafL3PeerLink     = LinkRole{Value: "leaf_l3_peer_link"}
	LinkRoleLeafL3Server       = LinkRole{Value: "leaf_l3_server"}
	LinkRoleLeafLeaf           = LinkRole{Value: "leaf_leaf"}
	LinkRoleLeafPairAccess     = LinkRole{Value: "leaf_pair_access"}
	LinkRoleLeafPairAccessPair = LinkRole{Value: "leaf_pair_access_pair"}
	LinkRoleLeafPairL2Server   = LinkRole{Value: "leaf_pair_l2_server"}
	LinkRoleLeafPeerLink       = LinkRole{Value: "leaf_peer_link"}
	LinkRoleSpineLeaf          = LinkRole{Value: "spine_leaf"}
	LinkRoleSpineSuperspine    = LinkRole{Value: "spine_superspine"}
	LinkRoleToExternalRouter   = LinkRole{Value: "to_external_router"}
	LinkRoleToGeneric          = LinkRole{Value: "to_generic"}
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

type LinkSwitchPeer oenum.Member[string]

var (
	LinkSwitchPeerFirst       = LinkSwitchPeer{Value: "first"}
	LinkSwitchPeerSecond      = LinkSwitchPeer{Value: "second"}
	LinkSwitchPeerUnspecified = LinkSwitchPeer{Value: ""}
)

type LinkType oenum.Member[string]

var (
	LinkTypeAggregateLink = LinkType{Value: "aggregate_link"}
	LinkTypeEthernet      = LinkType{Value: "ethernet"}
	LinkTypeLogicalLink   = LinkType{Value: "logical_link"}
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

type NodeRole oenum.Member[string]

var (
	NodeRoleAccess        = NodeRole{Value: "access"}
	NodeRoleGeneric       = NodeRole{Value: "generic"}
	NodeRoleLeaf          = NodeRole{Value: "leaf"}
	NodeRoleRemoteGateway = NodeRole{Value: "remote_gateway"}
	NodeRoleSpine         = NodeRole{Value: "spine"}
	NodeRoleSuperspine    = NodeRole{Value: "superspine"}
)

type OverlayControlProtocol oenum.Member[string]

var (
	OverlayControlProtocolEVPN = OverlayControlProtocol{Value: "evpn"}
	OverlayControlProtocolNone = OverlayControlProtocol{Value: ""}
)

type PolicyApplicationPointType oenum.Member[string]

var (
	PolicyApplicationPointTypeExternal       = PolicyApplicationPointType{Value: "external"}
	PolicyApplicationPointTypeGroup          = PolicyApplicationPointType{Value: "group"}
	PolicyApplicationPointTypeInternal       = PolicyApplicationPointType{Value: "internal"}
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

type PortRole oenum.Member[string]

var (
	PortRoleAccess     = PortRole{Value: "access"}
	PortRoleGeneric    = PortRole{Value: "generic"}
	PortRoleL3Server   = PortRole{Value: "l3_server"} // todo: remove this, LogicalDevicePortRoles.Validate and simplify LogicalDevicePortRoles.enableAll
	PortRoleLeaf       = PortRole{Value: "leaf"}
	PortRolePeer       = PortRole{Value: "peer"}
	PortRoleSpine      = PortRole{Value: "spine"}
	PortRoleSuperspine = PortRole{Value: "superspine"}
	PortRoleUnused     = PortRole{Value: "unused"}
)

type RedundancyGroupType oenum.Member[string]

var (
	RedundancyGroupTypeEsi  = RedundancyGroupType{Value: "esi"}
	RedundancyGroupTypeMlag = RedundancyGroupType{Value: "mlag"}
)

type RefDesign oenum.Member[string]

var (
	RefDesignDatacenter    = RefDesign{Value: "two_stage_l3clos"}
	RefDesignFreeform      = RefDesign{Value: "freeform"}
	RefDesignRailCollapsed = RefDesign{Value: "rail_collapsed"}
)

type RefDesignCapability oenum.Member[string]

var (
	RefDesignCapabilityDisabled    = RefDesignCapability{Value: "disabled"}
	RefDesignCapabilityFullSupport = RefDesignCapability{Value: "full_support"}
)

type RemoteGatewayRouteType oenum.Member[string]

var (
	RemoteGatewayRouteTypeAll      = RemoteGatewayRouteType{Value: "all"}
	RemoteGatewayRouteTypeFiveOnly = RemoteGatewayRouteType{Value: "type5_only"}
)

type RenderedConfigType oenum.Member[string]

var (
	RenderedConfigTypeDeployed = RenderedConfigType{Value: "deployed"}
	RenderedConfigTypeStaging  = RenderedConfigType{Value: "staging"}
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
	RoutingZoneConstraintModeAllow = RoutingZoneConstraintMode{Value: "allow"}
	RoutingZoneConstraintModeDeny  = RoutingZoneConstraintMode{Value: "deny"}
	RoutingZoneConstraintModeNone  = RoutingZoneConstraintMode{Value: "none"}
)

type SecurityZoneType oenum.Member[string]

var (
	SecurityZoneTypeEVPN            = SecurityZoneType{Value: "evpn"}
	SecurityZoneTypeL3Fabric        = SecurityZoneType{Value: "l3_fabric"}
	SecurityZoneTypeVirtualL3Fabric = SecurityZoneType{Value: "virtual_l3_fabric"}
)

type SpeedUnit oenum.Member[string]

var (
	SpeedUnitG = SpeedUnit{Value: "G"}
	SpeedUnitM = SpeedUnit{Value: "M"}
	// SpeedUnitT = SpeedUnit{Value: "T"}
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

type SwitchingZoneMACVRFServiceType oenum.Member[string]

var (
	SwitchingZoneMACVRFServiceTypeVLANAware  = SwitchingZoneMACVRFServiceType{Value: "vlan_aware"}
	SwitchingZoneMACVRFServiceTypeVLANBundle = SwitchingZoneMACVRFServiceType{Value: "vlan_bundle"}
)

type SystemManagementLevel oenum.Member[string]

var ( // do not introduce "not_installed" - that's an agent parameter with similar enum values
	SystemManagementLevelFullControl   = SystemManagementLevel{Value: "full_control"}
	SystemManagementLevelTelemetryOnly = SystemManagementLevel{Value: "telemetry_only"}
	SystemManagementLevelUnmanaged     = SystemManagementLevel{Value: "unmanaged"}
)

type SystemNodeRole oenum.Member[string]

var (
	SystemNodeRoleAccess        = SystemNodeRole{Value: "access"}
	SystemNodeRoleGeneric       = SystemNodeRole{Value: "generic"}
	SystemNodeRoleL3Server      = SystemNodeRole{Value: "l3_server"}
	SystemNodeRoleLeaf          = SystemNodeRole{Value: "leaf"}
	SystemNodeRoleRemoteGateway = SystemNodeRole{Value: "remote_gateway"}
	SystemNodeRoleSpine         = SystemNodeRole{Value: "spine"}
	SystemNodeRoleSuperspine    = SystemNodeRole{Value: "superspine"}
)

type SystemType oenum.Member[string]

var (
	SystemTypeServer = SystemType{Value: "server"}
	SystemTypeSwitch = SystemType{Value: "switch"}
)

type TcpStateQualifier oenum.Member[string]

var TcpStateQualifierEstablished = TcpStateQualifier{Value: "established"}

type TemplateCapability oenum.Member[string]

var (
	TemplateCapabilityBlueprint = TemplateCapability{Value: "blueprint"}
	TemplateCapabilityPod       = TemplateCapability{Value: "pod"}
)

type TemplateType oenum.Member[string]

var (
	TemplateTypeL3Collapsed   = TemplateType{Value: "l3_collapsed"}
	TemplateTypePodBased      = TemplateType{Value: "pod_based"}
	TemplateTypeRackBased     = TemplateType{Value: "rack_based"}
	TemplateTypeRailCollapsed = TemplateType{Value: "rail_collapsed"}
)

type VnType oenum.Member[string]

var (
	VnTypeExternal = VnType{Value: "disabled"}
	VnTypeVlan     = VnType{Value: "vlan"}
	VnTypeVxlan    = VnType{Value: "vxlan"}
)
