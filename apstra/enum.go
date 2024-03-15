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
