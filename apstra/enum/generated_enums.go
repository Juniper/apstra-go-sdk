// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// contents of this file are auto-generated by ./generator/generator.go - DO NOT EDIT

package enum

import oenum "github.com/orsinium-labs/enum"

var _ enum = (*ApiFeature)(nil)

func (o ApiFeature) String() string {
	return o.Value
}

func (o *ApiFeature) FromString(s string) error {
	if ApiFeatures.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*DeployMode)(nil)

func (o DeployMode) String() string {
	return o.Value
}

func (o *DeployMode) FromString(s string) error {
	if DeployModes.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*DeviceProfileType)(nil)

func (o DeviceProfileType) String() string {
	return o.Value
}

func (o *DeviceProfileType) FromString(s string) error {
	if DeviceProfileTypes.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*FFResourceType)(nil)

func (o FFResourceType) String() string {
	return o.Value
}

func (o *FFResourceType) FromString(s string) error {
	if FFResourceTypes.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*FeatureSwitch)(nil)

func (o FeatureSwitch) String() string {
	return o.Value
}

func (o *FeatureSwitch) FromString(s string) error {
	if FeatureSwitches.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*IbaWidgetType)(nil)

func (o IbaWidgetType) String() string {
	return o.Value
}

func (o *IbaWidgetType) FromString(s string) error {
	if IbaWidgetTypes.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*InterfaceNumberingIpv4Type)(nil)

func (o InterfaceNumberingIpv4Type) String() string {
	return o.Value
}

func (o *InterfaceNumberingIpv4Type) FromString(s string) error {
	if InterfaceNumberingIpv4Types.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*InterfaceNumberingIpv6Type)(nil)

func (o InterfaceNumberingIpv6Type) String() string {
	return o.Value
}

func (o *InterfaceNumberingIpv6Type) FromString(s string) error {
	if InterfaceNumberingIpv6Types.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*JunosEvpnIrbMode)(nil)

func (o JunosEvpnIrbMode) String() string {
	return o.Value
}

func (o *JunosEvpnIrbMode) FromString(s string) error {
	if JunosEvpnIrbModes.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*PolicyApplicationPointType)(nil)

func (o PolicyApplicationPointType) String() string {
	return o.Value
}

func (o *PolicyApplicationPointType) FromString(s string) error {
	if PolicyApplicationPointTypes.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*PolicyRuleAction)(nil)

func (o PolicyRuleAction) String() string {
	return o.Value
}

func (o *PolicyRuleAction) FromString(s string) error {
	if PolicyRuleActions.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*PolicyRuleProtocol)(nil)

func (o PolicyRuleProtocol) String() string {
	return o.Value
}

func (o *PolicyRuleProtocol) FromString(s string) error {
	if PolicyRuleProtocols.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*PortRole)(nil)

func (o PortRole) String() string {
	return o.Value
}

func (o *PortRole) FromString(s string) error {
	if PortRoles.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*RemoteGatewayRouteType)(nil)

func (o RemoteGatewayRouteType) String() string {
	return o.Value
}

func (o *RemoteGatewayRouteType) FromString(s string) error {
	if RemoteGatewayRouteTypes.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*RenderedConfigType)(nil)

func (o RenderedConfigType) String() string {
	return o.Value
}

func (o *RenderedConfigType) FromString(s string) error {
	if RenderedConfigTypes.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*ResourcePoolType)(nil)

func (o ResourcePoolType) String() string {
	return o.Value
}

func (o *ResourcePoolType) FromString(s string) error {
	if ResourcePoolTypes.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*RoutingZoneConstraintMode)(nil)

func (o RoutingZoneConstraintMode) String() string {
	return o.Value
}

func (o *RoutingZoneConstraintMode) FromString(s string) error {
	if RoutingZoneConstraintModes.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*StorageSchemaPath)(nil)

func (o StorageSchemaPath) String() string {
	return o.Value
}

func (o *StorageSchemaPath) FromString(s string) error {
	if StorageSchemaPaths.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var _ enum = (*TcpStateQualifier)(nil)

func (o TcpStateQualifier) String() string {
	return o.Value
}

func (o *TcpStateQualifier) FromString(s string) error {
	if TcpStateQualifiers.Parse(s) == nil {
		return newEnumParseError(o, s)
	}
	o.Value = s
	return nil
}

var (
	_           enum = new(ApiFeature)
	ApiFeatures      = oenum.New(
		ApiFeatureAiFabric,
		ApiFeatureCentral,
		ApiFeatureEnterprise,
		ApiFeatureFreeform,
		ApiFeatureFullAccess,
		ApiFeatureTaskApi,
	)

	_           enum = new(DeployMode)
	DeployModes      = oenum.New(
		DeployModeDeploy,
		DeployModeDrain,
		DeployModeNone,
		DeployModeReady,
		DeployModeUndeploy,
	)

	_                  enum = new(DeviceProfileType)
	DeviceProfileTypes      = oenum.New(
		DeviceProfileTypeModular,
		DeviceProfileTypeMonolithic,
	)

	_               enum = new(FFResourceType)
	FFResourceTypes      = oenum.New(
		FFResourceTypeAsn,
		FFResourceTypeHostIpv4,
		FFResourceTypeHostIpv6,
		FFResourceTypeInt,
		FFResourceTypeIpv4,
		FFResourceTypeIpv6,
		FFResourceTypeVlan,
		FFResourceTypeVni,
	)

	_               enum = new(FeatureSwitch)
	FeatureSwitches      = oenum.New(
		FeatureSwitchEnabled,
		FeatureSwitchDisabled,
	)

	_              enum = new(IbaWidgetType)
	IbaWidgetTypes      = oenum.New(
		IbaWidgetTypeStage,
		IbaWidgetTypeAnomalyHeatmap,
	)

	_                           enum = new(InterfaceNumberingIpv4Type)
	InterfaceNumberingIpv4Types      = oenum.New(
		InterfaceNumberingIpv4TypeNone,
		InterfaceNumberingIpv4TypeNumbered,
	)

	_                           enum = new(InterfaceNumberingIpv6Type)
	InterfaceNumberingIpv6Types      = oenum.New(
		InterfaceNumberingIpv6TypeNone,
		InterfaceNumberingIpv6TypeNumbered,
		InterfaceNumberingIpv6TypeLinkLocal,
	)

	_                 enum = new(JunosEvpnIrbMode)
	JunosEvpnIrbModes      = oenum.New(
		JunosEvpnIrbModeSymmetric,
		JunosEvpnIrbModeAsymmetric,
	)

	_                           enum = new(PolicyApplicationPointType)
	PolicyApplicationPointTypes      = oenum.New(
		PolicyApplicationPointTypeGroup,
		PolicyApplicationPointTypeInternal,
		PolicyApplicationPointTypeExternal,
		PolicyApplicationPointTypeSecurityZone,
		PolicyApplicationPointTypeVirtualNetwork,
	)

	_                 enum = new(PolicyRuleAction)
	PolicyRuleActions      = oenum.New(
		PolicyRuleActionDeny,
		PolicyRuleActionDenyLog,
		PolicyRuleActionPermit,
		PolicyRuleActionPermitLog,
	)

	_                   enum = new(PolicyRuleProtocol)
	PolicyRuleProtocols      = oenum.New(
		PolicyRuleProtocolIcmp,
		PolicyRuleProtocolIp,
		PolicyRuleProtocolTcp,
		PolicyRuleProtocolUdp,
	)

	_         enum = new(PortRole)
	PortRoles      = oenum.New(
		PortRoleAccess,
		PortRoleGeneric,
		PortRoleL3Server,
		PortRoleLeaf,
		PortRolePeer,
		PortRoleSpine,
		PortRoleSuperspine,
		PortRoleUnused,
	)

	_                       enum = new(RemoteGatewayRouteType)
	RemoteGatewayRouteTypes      = oenum.New(
		RemoteGatewayRouteTypeAll,
		RemoteGatewayRouteTypeFiveOnly,
	)

	_                   enum = new(RenderedConfigType)
	RenderedConfigTypes      = oenum.New(
		RenderedConfigTypeStaging,
		RenderedConfigTypeDeployed,
	)

	_                 enum = new(ResourcePoolType)
	ResourcePoolTypes      = oenum.New(
		ResourcePoolTypeAsn,
		ResourcePoolTypeInt,
		ResourcePoolTypeIpv4,
		ResourcePoolTypeIpv6,
		ResourcePoolTypeVlan,
		ResourcePoolTypeVni,
	)

	_                          enum = new(RoutingZoneConstraintMode)
	RoutingZoneConstraintModes      = oenum.New(
		RoutingZoneConstraintModeNone,
		RoutingZoneConstraintModeAllow,
		RoutingZoneConstraintModeDeny,
	)

	_                  enum = new(StorageSchemaPath)
	StorageSchemaPaths      = oenum.New(
		StorageSchemaPathARP,
		StorageSchemaPathBGP,
		StorageSchemaPathCppGraph,
		StorageSchemaPathEnvironment,
		StorageSchemaPathGeneric,
		StorageSchemaPathGraph,
		StorageSchemaPathHostname,
		StorageSchemaPathIbaData,
		StorageSchemaPathIbaIntegerData,
		StorageSchemaPathIbaStringData,
		StorageSchemaPathInterface,
		StorageSchemaPathInterfaceCounters,
		StorageSchemaPathLAG,
		StorageSchemaPathLLDP,
		StorageSchemaPathMAC,
		StorageSchemaPathMLAG,
		StorageSchemaPathNSXT,
		StorageSchemaPathOpticalXcvr,
		StorageSchemaPathRoute,
		StorageSchemaPathRouteLookup,
		StorageSchemaPathXcvr,
	)

	_                  enum = new(TcpStateQualifier)
	TcpStateQualifiers      = oenum.New(
		TcpStateQualifierEstablished,
	)
)
