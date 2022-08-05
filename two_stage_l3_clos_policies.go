package goapstra

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	apiUrlPolicies   = apiUrlBlueprintById + apiUrlPathDelim + "policies"
	apiUrlPolicyById = apiUrlPolicies + apiUrlPathDelim + "%s"

	portAny       = "any"
	portRangeSep  = "-"
	portRangesSep = ","
)

// RULE_SCHEMA = {
//    'id': s.Optional(s.NodeId(description='ID of the rule node')),
//    'label': s.GenericName(description='Unique user-friendly name of the rule'),
//    'protocol': s.SecurityRuleProtocol(),
//    'src_port': s.PortSetOrAny(),
//    'dst_port': s.PortSetOrAny(),
//    'description': s.Optional(s.Description(), load_default=''),
//    'action': s.SecurityRuleAction() //             ['deny', 'deny_log', 'permit', 'permit_log'],
//}

type PolicyRuleAction int
type policyRuleAction string

const (
	PolicyRuleActionDeny = iota
	PolicyRuleActionDenyLog
	PolicyRuleActionPermit
	PolicyRuleActionPermitLog
	PolicyRuleActionUnknown = "unknown policy action '%s'"

	policyRuleActionDeny      = policyRuleAction("deny")
	policyRuleActionDenyLog   = policyRuleAction("deny_log")
	policyRuleActionPermit    = policyRuleAction("permit")
	policyRuleActionPermitLog = policyRuleAction("permit_log")
	policyRuleActionUnknown   = "unknown policy action %d"
)

func (o PolicyRuleAction) Int() int {
	return int(o)
}

func (o PolicyRuleAction) String() string {
	return string(o.raw())
}

func (o PolicyRuleAction) raw() policyRuleAction {
	switch o {
	case PolicyRuleActionDeny:
		return policyRuleActionDeny
	case PolicyRuleActionDenyLog:
		return policyRuleActionDenyLog
	case PolicyRuleActionPermit:
		return policyRuleActionPermit
	case PolicyRuleActionPermitLog:
		return policyRuleActionPermitLog
	default:
		return policyRuleAction(fmt.Sprintf(policyRuleActionUnknown, o))
	}
}

func (o policyRuleAction) string() string {
	return string(o)
}

func (o policyRuleAction) parse() (int, error) {
	switch o {
	case policyRuleActionDeny:
		return PolicyRuleActionDeny, nil
	case policyRuleActionDenyLog:
		return PolicyRuleActionDenyLog, nil
	case policyRuleActionPermit:
		return PolicyRuleActionPermit, nil
	case policyRuleActionPermitLog:
		return PolicyRuleActionPermitLog, nil
	default:
		return 0, fmt.Errorf(PolicyRuleActionUnknown, o)
	}
}

type PortRange struct {
	first uint16
	last  uint16
}

func (o PortRange) string() string {
	switch {
	case o.first == o.last:
		return strconv.Itoa(int(o.first))
	case o.first < o.last:
		return strconv.Itoa(int(o.first)) + portRangeSep + strconv.Itoa(int(o.last))
	default:
		return strconv.Itoa(int(o.last)) + portRangeSep + strconv.Itoa(int(o.first))
	}
}

type rawPortRanges string

func (o rawPortRanges) parse() (PortRanges, error) {
	if o == portAny {
		return []PortRange{}, nil
	}

	rawRangeSlice := strings.Split(string(o), portRangesSep)
	result := make([]PortRange, len(rawRangeSlice))
	for i, raw := range rawRangeSlice {
		var first, last uint64
		var err error
		portStrs := strings.Split(raw, portRangeSep)
		switch len(portStrs) {
		case 1:
			first, err = strconv.ParseUint(raw, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("error parsing port range '%s' - %w", raw, err)
			}
			last = first
		case 2:
			first, err = strconv.ParseUint(portStrs[0], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("error parsing first element of port range '%s' - %w", raw, err)
			}
			last, err = strconv.ParseUint(portStrs[1], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("error parsing last element of port range '%s' - %w", raw, err)
			}
		default:
			return nil, fmt.Errorf("cannot parse port range '%s'", raw)
		}
		if first > math.MaxUint16 || last > math.MaxUint16 {
			return nil, fmt.Errorf("port spec '%s' falls outside of range %d-%d", raw, 0, math.MaxUint16)
		}
		result[i] = PortRange{
			first: uint16(first),
			last:  uint16(last),
		}
	}
	return result, nil
}

type PortRanges []PortRange

func (o PortRanges) string() string {
	if len(o) == 0 {
		return portAny
	}
	sb := strings.Builder{}
	sb.WriteString(o[0].string())
	for _, pr := range o[1:] {
		sb.WriteString(portRangesSep + pr.string())
	}
	return sb.String()
}

type PolicyRule struct {
	Id          ObjectId
	Label       string
	Description string
	Protocol    string
	Action      PolicyRuleAction
	SrcPort     PortRanges
	DstPort     PortRanges
}

func (o PolicyRule) raw() *rawPolicyRule {
	return &rawPolicyRule{
		Id:          o.Id,
		Label:       o.Label,
		Description: o.Description,
		Protocol:    o.Protocol,
		Action:      o.Action.raw(),
		SrcPort:     rawPortRanges(o.SrcPort.string()),
		DstPort:     rawPortRanges(o.DstPort.string()),
	}
}

type rawPolicyRule struct {
	Id          ObjectId         `json:"id,omitempty"`
	Label       string           `json:"label"`
	Description string           `json:"description"`
	Protocol    string           `json:"protocol"`
	Action      policyRuleAction `json:"action"`
	SrcPort     rawPortRanges    `json:"src_port"`
	DstPort     rawPortRanges    `json:"dst_port"`
}

func (o rawPolicyRule) polish() (*PolicyRule, error) {
	action, err := o.Action.parse()
	if err != nil {
		return nil, err
	}
	srcPort, err := o.SrcPort.parse()
	if err != nil {
		return nil, err
	}
	dstPort, err := o.DstPort.parse()
	if err != nil {
		return nil, err
	}
	return &PolicyRule{
		Id:          o.Id,
		Label:       o.Label,
		Description: o.Description,
		Protocol:    o.Protocol,
		Action:      PolicyRuleAction(action),
		SrcPort:     srcPort,
		DstPort:     dstPort,
	}, nil
}

type Policy struct {
	Enabled             bool                   `json:"enabled"`
	Label               string                 `json:"label"`
	Description         string                 `json:"description"`
	SrcApplicationPoint PolicyApplicationPoint `json:"src_application_point"`
	DstApplicationPoint PolicyApplicationPoint `json:"dst_application_point"`
	Rules               []PolicyRule           `json:"rules"`
	Tags                []TagLabel             `json:"tags"`
	Id                  ObjectId               `json:"object_id,omitempty"`
}

func (o Policy) request() *policyRequest {
	rules := make([]rawPolicyRule, len(o.Rules))
	for i, rule := range o.Rules {
		rules[i] = *rule.raw()
	}
	return &policyRequest{
		Enabled:             o.Enabled,
		Label:               o.Label,
		Description:         o.Description,
		SrcApplicationPoint: o.SrcApplicationPoint.objectId(),
		DstApplicationPoint: o.DstApplicationPoint.objectId(),
		Rules:               rules,
		Tags:                o.Tags,
		Id:                  o.Id,
	}
}

type policyRequest struct {
	Enabled             bool            `json:"enabled"`
	Label               string          `json:"label"`
	Description         string          `json:"description"`
	SrcApplicationPoint ObjectId        `json:"src_application_point"`
	DstApplicationPoint ObjectId        `json:"dst_application_point"`
	Rules               []rawPolicyRule `json:"rules"`
	Tags                []TagLabel      `json:"tags"`
	Id                  ObjectId        `json:"object_id,omitempty"`
}

type policyResponse struct {
	Enabled             bool                         `json:"enabled"`
	Label               string                       `json:"label"`
	Description         string                       `json:"description"`
	SrcApplicationPoint PolicyApplicationPointDigest `json:"src_application_point"`
	DstApplicationPoint PolicyApplicationPointDigest `json:"dst_application_point"`
	Rules               []rawPolicyRule              `json:"rules"`
	Tags                []TagLabel                   `json:"tags"`
	Id                  ObjectId                     `json:"id,omitempty"`
}

func (o policyResponse) polish() (*Policy, error) {
	rules := make([]PolicyRule, len(o.Rules))
	for i, raw := range o.Rules {
		polished, err := raw.polish()
		if err != nil {
			return nil, err
		}
		rules[i] = *polished
	}
	return &Policy{
		Enabled:             o.Enabled,
		Label:               o.Label,
		Description:         o.Description,
		SrcApplicationPoint: o.SrcApplicationPoint,
		DstApplicationPoint: o.DstApplicationPoint,
		Rules:               rules,
		Tags:                o.Tags,
		Id:                  o.Id,
	}, nil
}

type PolicyApplicationPoint interface {
	objectId() ObjectId
}

type PolicyApplicationPointDigest struct {
	Id    ObjectId `json:"id"`
	Label string   `json:"label"`
	Type  string   `json:"type"` // group, internal, external, security_zone, virtual_network
}

func (o PolicyApplicationPointDigest) objectId() ObjectId {
	return o.Id
}

func (o *TwoStageLThreeClosClient) getAllPolicies(ctx context.Context) ([]Policy, error) {
	response := &struct {
		Policies []policyResponse `json:"policies"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodGet,
		urlStr:         fmt.Sprintf(apiUrlPolicies, o.blueprintId),
		apiResponse:    response,
		unsynchronized: true,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make([]Policy, len(response.Policies))
	for i, policy := range response.Policies {
		polished, err := policy.polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}
	return result, nil
}

func (o *TwoStageLThreeClosClient) getPolicy(ctx context.Context, id ObjectId) (*Policy, error) {
	response := &policyResponse{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodGet,
		urlStr:         fmt.Sprintf(apiUrlPolicyById, o.blueprintId, id),
		apiResponse:    response,
		unsynchronized: true,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.polish()
}

func (o *TwoStageLThreeClosClient) createPolicy(ctx context.Context, policy *Policy) (ObjectId, error) {
	response := &struct {
		Id ObjectId `json:"id"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodPost,
		urlStr:         fmt.Sprintf(apiUrlPolicies, o.blueprintId),
		apiInput:       policy,
		apiResponse:    response,
		unsynchronized: true,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *TwoStageLThreeClosClient) deletePolicy(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodDelete,
		urlStr:         fmt.Sprintf(apiUrlPolicyById, o.blueprintId, id),
		unsynchronized: true,
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *TwoStageLThreeClosClient) updatePolicy(ctx context.Context, id ObjectId, policy *Policy) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodPut,
		urlStr:         fmt.Sprintf(apiUrlPolicyById, o.blueprintId, id),
		apiInput:       policy.request(),
		unsynchronized: true,
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *TwoStageLThreeClosClient) getPolicyRuleIdByLabel(ctx context.Context, policyId ObjectId, label string) (ObjectId, error) {
	start := time.Now()
	for i := 0; i <= dcClientMaxRetries; i++ {
		time.Sleep(dcClientRetryBackoff * time.Duration(i))
		policy, err := o.getPolicy(ctx, policyId)
		if err != nil {
			return "", err
		}
		for _, rule := range policy.Rules {
			if rule.Label == label {
				return rule.Id, nil
			}
		}
	}
	return "", fmt.Errorf("rule '%s' didn't appear in policy '%s' after %s", label, policyId, time.Since(start))
}

func (o *TwoStageLThreeClosClient) addPolicyRule(ctx context.Context, rule *PolicyRule, position int, policyId ObjectId) (ObjectId, error) {
	policy, err := o.getPolicy(ctx, policyId)
	if err != nil {
		return "", err
	}

	currentRuleCount := len(policy.Rules)

	if position < 0 {
		position = currentRuleCount
	}

	switch {
	// empty rule set is an easy case
	case currentRuleCount == 0:
		policy.Rules = []PolicyRule{*rule}
	// zero insertion point
	case position <= 0:
		policy.Rules = append([]PolicyRule{*rule}, policy.Rules...)
	// end insertion point
	case position >= currentRuleCount:
		policy.Rules = append(policy.Rules, *rule)
	// insert in the middle
	default:
		policy.Rules = append(policy.Rules[:position+1], policy.Rules[position:]...)
		policy.Rules[position] = *rule
	}

	// push the new policy
	err = o.updatePolicy(ctx, policyId, policy)
	if err != nil {
		return "", err
	}

	return o.getPolicyRuleIdByLabel(ctx, policyId, rule.Label)
}

func (o *TwoStageLThreeClosClient) deletePolicyRuleById(ctx context.Context, policyId ObjectId, ruleId ObjectId) error {
	policy, err := o.getPolicy(ctx, policyId)
	if err != nil {
		return err
	}

	ruleIdx := -1
	for i, rule := range policy.Rules {
		if rule.Id == ruleId {
			ruleIdx = i
			break
		}
	}

	if ruleIdx < 0 {
		return ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("rule id '%s' not found in policy '%s'", ruleId, policyId),
		}
	}

	policy.Rules = append(policy.Rules[:ruleIdx], policy.Rules[ruleIdx+1:]...)
	return o.updatePolicy(ctx, policyId, policy)
}
