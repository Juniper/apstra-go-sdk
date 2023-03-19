package goapstra

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const (
	apiUrlBlueprintRoutingPolicies       = apiUrlBlueprintById + apiUrlPathDelim + "routing-policies"
	apiUrlBlueprintRoutingPoliciesPrefix = apiUrlBlueprintRoutingPolicies + apiUrlPathDelim
	apiUrlBlueprintRoutingPolicyById     = apiUrlBlueprintRoutingPoliciesPrefix + "%s"
)

type DcRoutingPolicyType int
type dcRoutingPolicyType string

const (
	DcRoutingPolicyTypeNone = DcRoutingPolicyType(iota)
	DcRoutingPolicyTypeDefault
	DcRoutingPolicyTypeUser
	DcRoutingPolicyTypeUnknown = "unknown datacenter routing policy type %q"

	dcRoutingPolicyTypeNone    = dcRoutingPolicyType("")
	dcRoutingPolicyTypeDefault = dcRoutingPolicyType("default_immutable")
	dcRoutingPolicyTypeUser    = dcRoutingPolicyType("user_defined")
	dcRoutingPolicyTypeUnknown = "unknown datacenter routing policy type %d"
)

func (o DcRoutingPolicyType) Int() int {
	return int(o)
}

func (o DcRoutingPolicyType) String() string {
	switch o {
	case DcRoutingPolicyTypeNone:
		return string(dcRoutingPolicyTypeNone)
	case DcRoutingPolicyTypeDefault:
		return string(dcRoutingPolicyTypeDefault)
	case DcRoutingPolicyTypeUser:
		return string(dcRoutingPolicyTypeUser)
	default:
		return fmt.Sprintf(dcRoutingPolicyTypeUnknown, o)
	}
}

func (o *DcRoutingPolicyType) FromString(in string) error {
	i, err := dcRoutingPolicyType(in).parse()
	if err != nil {
		return err
	}
	*o = DcRoutingPolicyType(i)
	return nil
}

func (o DcRoutingPolicyType) raw() dcRoutingPolicyType {
	return dcRoutingPolicyType(o.String())
}

func (o dcRoutingPolicyType) string() string {
	return string(o)
}

func (o dcRoutingPolicyType) parse() (int, error) {
	switch o {
	case dcRoutingPolicyTypeNone:
		return int(DcRoutingPolicyTypeNone), nil
	case dcRoutingPolicyTypeDefault:
		return int(DcRoutingPolicyTypeDefault), nil
	case dcRoutingPolicyTypeUser:
		return int(DcRoutingPolicyTypeUser), nil
	default:
		return 0, fmt.Errorf(DcRoutingPolicyTypeUnknown, o)
	}
}

type PrefixFilterAction int
type prefixFilterAction string

const (
	PrefixFilterActionNone = PrefixFilterAction(iota)
	PrefixFilterActionPermit
	PrefixFilterActionDeny
	PrefixFilterActionUnknown = "unknown prefix filter action %q"

	prefixFilterActionNone    = prefixFilterAction("")
	prefixFilterActionPermit  = prefixFilterAction("permit")
	prefixFilterActionDeny    = prefixFilterAction("deny")
	prefixFilterActionUnknown = "unknown prefix filter action %d"
)

func (o PrefixFilterAction) Int() int {
	return int(o)
}

func (o PrefixFilterAction) String() string {
	switch o {
	case PrefixFilterActionNone:
		return string(prefixFilterActionNone)
	case PrefixFilterActionPermit:
		return string(prefixFilterActionPermit)
	case PrefixFilterActionDeny:
		return string(prefixFilterActionDeny)
	default:
		return fmt.Sprintf(prefixFilterActionUnknown, o)
	}
}

func (o *PrefixFilterAction) FromString(in string) error {
	i, err := prefixFilterAction(in).parse()
	if err != nil {
		return err
	}
	*o = PrefixFilterAction(i)
	return nil
}

func (o PrefixFilterAction) raw() prefixFilterAction {
	return prefixFilterAction(o.String())
}

func (o prefixFilterAction) string() string {
	return string(o)
}

func (o prefixFilterAction) parse() (int, error) {
	switch o {
	case prefixFilterActionNone:
		return int(PrefixFilterActionNone), nil
	case prefixFilterActionPermit:
		return int(PrefixFilterActionPermit), nil
	case prefixFilterActionDeny:
		return int(PrefixFilterActionDeny), nil
	default:
		return 0, fmt.Errorf(PrefixFilterActionUnknown, o)
	}
}

type DcRoutingPolicyImportPolicy int
type dcRoutingPolicyImportPolicy string

const (
	DcRoutingPolicyImportPolicyNone = DcRoutingPolicyImportPolicy(iota)
	DcRoutingPolicyImportPolicyDefaultOnly
	DcRoutingPolicyImportPolicyAll
	DcRoutingPolicyImportPolicyExtraOnly
	DcRoutingPolicyImportPolicyUnknown = "unknown import policy %s"

	dcRoutingPolicyImportPolicyNone        = dcRoutingPolicyImportPolicy("")
	dcRoutingPolicyImportPolicyDefaultOnly = dcRoutingPolicyImportPolicy("default_only")
	dcRoutingPolicyImportPolicyAll         = dcRoutingPolicyImportPolicy("all")
	dcRoutingPolicyImportPolicyExtraOnly   = dcRoutingPolicyImportPolicy("extra_only")
	dcRoutingPolicyImportPolicyUnknown     = "unknown import policy %d"
)

func (o DcRoutingPolicyImportPolicy) Int() int {
	return int(o)
}

func (o DcRoutingPolicyImportPolicy) String() string {
	switch o {
	case DcRoutingPolicyImportPolicyNone:
		return string(dcRoutingPolicyImportPolicyNone)
	case DcRoutingPolicyImportPolicyDefaultOnly:
		return string(dcRoutingPolicyImportPolicyDefaultOnly)
	case DcRoutingPolicyImportPolicyAll:
		return string(dcRoutingPolicyImportPolicyAll)
	case DcRoutingPolicyImportPolicyExtraOnly:
		return string(dcRoutingPolicyImportPolicyExtraOnly)
	default:
		return fmt.Sprintf(dcRoutingPolicyImportPolicyUnknown, o)
	}
}

func (o *DcRoutingPolicyImportPolicy) FromString(in string) error {
	i, err := dcRoutingPolicyImportPolicy(in).parse()
	if err != nil {
		return err
	}
	*o = DcRoutingPolicyImportPolicy(i)
	return nil
}

func (o DcRoutingPolicyImportPolicy) raw() dcRoutingPolicyImportPolicy {
	return dcRoutingPolicyImportPolicy(o.String())
}

func (o dcRoutingPolicyImportPolicy) string() string {
	return string(o)
}

func (o dcRoutingPolicyImportPolicy) parse() (int, error) {
	switch o {
	case dcRoutingPolicyImportPolicyNone:
		return int(DcRoutingPolicyImportPolicyNone), nil
	case dcRoutingPolicyImportPolicyDefaultOnly:
		return int(DcRoutingPolicyImportPolicyDefaultOnly), nil
	case dcRoutingPolicyImportPolicyAll:
		return int(DcRoutingPolicyImportPolicyAll), nil
	case dcRoutingPolicyImportPolicyExtraOnly:
		return int(DcRoutingPolicyImportPolicyExtraOnly), nil
	default:
		return 0, fmt.Errorf(DcRoutingPolicyImportPolicyUnknown, o)
	}
}

type PrefixFilter struct {
	Action PrefixFilterAction `json:"action"`
	Prefix net.IPNet          `json:"prefix"`
	GeMask int                `json:"ge_mask"`
	LeMask int                `json:"le_mask"`
}

func (o *PrefixFilter) raw() *rawPrefixFilter {
	return &rawPrefixFilter{
		Action: o.Action.raw(),
		Prefix: o.Prefix.String(),
		GeMask: o.GeMask,
		LeMask: o.LeMask,
	}
}

type rawPrefixFilter struct {
	Action prefixFilterAction `json:"action"`
	Prefix string             `json:"prefix"`
	GeMask int                `json:"ge_mask"`
	LeMask int                `json:"le_mask"`
}

func (o *rawPrefixFilter) polish() (*PrefixFilter, error) {
	action, err := o.Action.parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing prefix filter action %q - %w", o.Action, err)
	}

	_, ipNet, err := net.ParseCIDR(o.Prefix)
	if err != nil {
		return nil, fmt.Errorf("error parsing prefix %q - %w", o.Prefix, err)
	}
	if ipNet == nil {
		return nil, fmt.Errorf("result of parsing prefix %q is nil", o.Prefix)
	}

	return &PrefixFilter{
		Action: PrefixFilterAction(action),
		Prefix: *ipNet,
		GeMask: o.GeMask,
		LeMask: o.LeMask,
	}, nil
}

type DcRoutingExportPolicy struct {
	StaticRoutes         bool `json:"static_routes"`
	Loopbacks            bool `json:"loopbacks"`
	SpineSuperspineLinks bool `json:"spine_superspine_links"`
	L3EdgeServerLinks    bool `json:"l3edge_server_links"`
	SpineLeafLinks       bool `json:"spine_leaf_links"`
	L2EdgeSubnets        bool `json:"l2edge_subnets"`
}

type rawDcRoutingPolicy struct {
	Id                     ObjectId                    `json:"id,omitempty"`
	Label                  string                      `json:"label"`
	Description            string                      `json:"description"`
	PolicyType             dcRoutingPolicyType         `json:"policy_type"`
	ImportPolicy           dcRoutingPolicyImportPolicy `json:"import_policy"`
	ExportPolicy           DcRoutingExportPolicy       `json:"export_policy"`
	ExpectDefaultIpv4Route bool                        `json:"expect_default_ipv4_route"`
	ExpectDefaultIpv6Route bool                        `json:"expect_default_ipv6_route"`
	AggregatePrefixes      []string                    `json:"aggregate_prefixes"`
	ExtraImportRoutes      []rawPrefixFilter           `json:"extra_import_routes"`
	ExtraExportRoutes      []rawPrefixFilter           `json:"extra_export_routes"`
}

func (o rawDcRoutingPolicy) polish() (*DcRoutingPolicy, error) {
	policyType, err := o.PolicyType.parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing datacenter routing policy type %q - %w", o.PolicyType, err)
	}

	importPolicy, err := o.ImportPolicy.parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing datacenter import policy %q - %w", o.ImportPolicy, err)
	}

	aggregatePrefixes := make([]net.IPNet, len(o.AggregatePrefixes))
	for i := range o.AggregatePrefixes {
		_, p, err := net.ParseCIDR(o.AggregatePrefixes[i])
		if err != nil {
			return nil, fmt.Errorf("error parsing aggregate prefix %q - %w", o.AggregatePrefixes[i], err)
		}
		if p == nil {
			return nil, fmt.Errorf("result of parsing prefix %q is nil", o.AggregatePrefixes[i])
		}
		aggregatePrefixes[i] = *p
	}

	extraImportRoutes := make([]PrefixFilter, len(o.ExtraImportRoutes))
	for i := range o.ExtraImportRoutes {
		r, err := o.ExtraImportRoutes[i].polish()
		if err != nil {
			return nil, fmt.Errorf("error parsing extra import route - %w", err)
		}
		extraImportRoutes[i] = *r
	}

	extraExportRoutes := make([]PrefixFilter, len(o.ExtraExportRoutes))
	for i := range o.ExtraExportRoutes {
		r, err := o.ExtraExportRoutes[i].polish()
		if err != nil {
			return nil, fmt.Errorf("error parsing extra export route - %w", err)
		}
		extraExportRoutes[i] = *r
	}

	return &DcRoutingPolicy{
		Id: o.Id,
		Data: &DcRoutingPolicyData{
			Label:                  o.Label,
			Description:            o.Description,
			PolicyType:             DcRoutingPolicyType(policyType),
			ImportPolicy:           DcRoutingPolicyImportPolicy(importPolicy),
			ExportPolicy:           o.ExportPolicy,
			ExpectDefaultIpv4Route: o.ExpectDefaultIpv4Route,
			ExpectDefaultIpv6Route: o.ExpectDefaultIpv6Route,
			AggregatePrefixes:      aggregatePrefixes,
			ExtraImportRoutes:      extraImportRoutes,
			ExtraExportRoutes:      extraExportRoutes,
		},
	}, nil
}

type DcRoutingPolicy struct {
	Id   ObjectId
	Data *DcRoutingPolicyData
}

type DcRoutingPolicyData struct {
	Label                  string                      `json:"label"`
	Description            string                      `json:"description"`
	PolicyType             DcRoutingPolicyType         `json:"policy_type"`
	ImportPolicy           DcRoutingPolicyImportPolicy `json:"import_policy"`
	ExportPolicy           DcRoutingExportPolicy       `json:"export_policy"`
	ExpectDefaultIpv4Route bool                        `json:"expect_default_ipv4_route"`
	ExpectDefaultIpv6Route bool                        `json:"expect_default_ipv6_route"`
	AggregatePrefixes      []net.IPNet                 `json:"aggregate_prefixes"`
	ExtraImportRoutes      []PrefixFilter              `json:"extra_import_routes"`
	ExtraExportRoutes      []PrefixFilter              `json:"extra_export_routes"`
}

func (o *DcRoutingPolicyData) raw() *rawDcRoutingPolicy {
	extraImportRoutes := make([]rawPrefixFilter, len(o.ExtraImportRoutes))
	for i := range o.ExtraImportRoutes {
		f := o.ExtraImportRoutes[i].raw()
		extraImportRoutes[i] = *f
	}

	extraExportRoutes := make([]rawPrefixFilter, len(o.ExtraExportRoutes))
	for i := range o.ExtraExportRoutes {
		f := o.ExtraExportRoutes[i].raw()
		extraExportRoutes[i] = *f
	}

	aggregatePrefixes := make([]string, len(o.AggregatePrefixes))
	for i := range o.AggregatePrefixes {
		aggregatePrefixes[i] = o.AggregatePrefixes[i].String()
	}

	return &rawDcRoutingPolicy{
		Label:        o.Label,
		Description:  o.Description,
		PolicyType:   o.PolicyType.raw(),
		ImportPolicy: o.ImportPolicy.raw(),
		ExportPolicy: DcRoutingExportPolicy{
			StaticRoutes:         o.ExportPolicy.StaticRoutes,
			Loopbacks:            o.ExportPolicy.Loopbacks,
			SpineSuperspineLinks: o.ExportPolicy.SpineSuperspineLinks,
			L3EdgeServerLinks:    o.ExportPolicy.L3EdgeServerLinks,
			SpineLeafLinks:       o.ExportPolicy.SpineLeafLinks,
			L2EdgeSubnets:        o.ExportPolicy.L2EdgeSubnets,
		},
		ExpectDefaultIpv4Route: o.ExpectDefaultIpv4Route,
		ExpectDefaultIpv6Route: o.ExpectDefaultIpv6Route,
		AggregatePrefixes:      aggregatePrefixes,
		ExtraImportRoutes:      extraImportRoutes,
		ExtraExportRoutes:      extraExportRoutes,
	}
}

func (o *TwoStageL3ClosClient) getAllRoutingPolicies(ctx context.Context) ([]rawDcRoutingPolicy, error) {
	response := &struct {
		Items []rawDcRoutingPolicy `json:"items"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintRoutingPolicies, o.blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *TwoStageL3ClosClient) getRoutingPolicy(ctx context.Context, id ObjectId) (*rawDcRoutingPolicy, error) {
	response := &rawDcRoutingPolicy{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintRoutingPolicyById, o.blueprintId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *TwoStageL3ClosClient) getDefaultRoutingPolicy(ctx context.Context) (*rawDcRoutingPolicy, error) {
	policies, err := o.getAllRoutingPolicies(ctx)
	if err != nil {
		return nil, err
	}

	for i := range policies {
		if policies[i].PolicyType == dcRoutingPolicyTypeDefault {
			return &policies[i], nil
		}
	}

	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err: fmt.Errorf("blueprint %q has %d policies, but none have type %q",
			o.blueprintId, len(policies), dcRoutingPolicyTypeDefault),
	}
}

func (o *TwoStageL3ClosClient) createRoutingPolicy(ctx context.Context, cfg *rawDcRoutingPolicy) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintRoutingPolicies, o.blueprintId),
		apiInput:    cfg,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *TwoStageL3ClosClient) updateRoutingPolicy(ctx context.Context, id ObjectId, cfg *rawDcRoutingPolicy) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintRoutingPolicyById, o.blueprintId, id),
		apiInput: cfg,
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) deleteRoutingPolicy(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintRoutingPolicyById, o.blueprintId, id),
	})
	return convertTtaeToAceWherePossible(err)
}
