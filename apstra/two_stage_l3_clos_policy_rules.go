package apstra

import (
	"context"
	"fmt"
	"github.com/orsinium-labs/enum"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	portAny       = "any"
	portRangeSep  = "-"
	portRangesSep = ","
)

type PolicyRuleAction enum.Member[string]

var (
	PolicyRuleActionDeny      = PolicyRuleAction{Value: "deny"}
	PolicyRuleActionDenyLog   = PolicyRuleAction{Value: "deny_log"}
	PolicyRuleActionPermit    = PolicyRuleAction{Value: "permit"}
	PolicyRuleActionPermitLog = PolicyRuleAction{Value: "permit_log"}
	PolicyRuleActions         = enum.New(
		PolicyRuleActionDeny,
		PolicyRuleActionDenyLog,
		PolicyRuleActionPermit,
		PolicyRuleActionPermitLog,
	)
)

type PolicyRuleProtocol enum.Member[string]

var (
	PolicyRuleProtocolIcmp = PolicyRuleProtocol{Value: "ICMP"}
	PolicyRuleProtocolIp   = PolicyRuleProtocol{Value: "IP"}
	PolicyRuleProtocolTcp  = PolicyRuleProtocol{Value: "TCP"}
	PolicyRuleProtocolUdp  = PolicyRuleProtocol{Value: "UDP"}
	PolicyRuleProtocols    = enum.New(
		PolicyRuleProtocolIcmp,
		PolicyRuleProtocolIp,
		PolicyRuleProtocolTcp,
		PolicyRuleProtocolUdp,
	)
)

type TcpStateQualifier enum.Member[string]

var (
	TcpStateQualifierEstablished = TcpStateQualifier{Value: "established"}
	TcpStateQualifiers           = enum.New(TcpStateQualifierEstablished)
)

type PortRange struct {
	First uint16
	Last  uint16
}

func (o PortRange) string() string {
	switch {
	case o.First == o.Last:
		return strconv.Itoa(int(o.First))
	case o.First < o.Last:
		return strconv.Itoa(int(o.First)) + portRangeSep + strconv.Itoa(int(o.Last))
	default:
		return strconv.Itoa(int(o.Last)) + portRangeSep + strconv.Itoa(int(o.First))
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
			First: uint16(first),
			Last:  uint16(last),
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
	Id   ObjectId
	Data *PolicyRuleData
}

type PolicyRuleData struct {
	Label             string
	Description       string
	Protocol          PolicyRuleProtocol
	Action            PolicyRuleAction
	SrcPort           PortRanges
	DstPort           PortRanges
	TcpStateQualifier *TcpStateQualifier
}

func (o PolicyRuleData) raw() *rawPolicyRule {
	var tcpStateQualifier *string
	if o.TcpStateQualifier != nil {
		s := o.TcpStateQualifier.Value
		tcpStateQualifier = &s
	}

	return &rawPolicyRule{
		Label:             o.Label,
		Description:       o.Description,
		Protocol:          o.Protocol.Value,
		Action:            o.Action.Value,
		SrcPort:           rawPortRanges(o.SrcPort.string()),
		DstPort:           rawPortRanges(o.DstPort.string()),
		TcpStateQualifier: tcpStateQualifier,
	}
}

type rawPolicyRule struct {
	Id                ObjectId      `json:"id,omitempty"`
	Label             string        `json:"label"`
	Description       string        `json:"description"`
	Protocol          string        `json:"protocol"`
	Action            string        `json:"action"`
	SrcPort           rawPortRanges `json:"src_port"`
	DstPort           rawPortRanges `json:"dst_port"`
	TcpStateQualifier *string       `json:"tcp_state_qualifier,omitempty"`
}

func (o rawPolicyRule) polish() (*PolicyRule, error) {
	action := PolicyRuleActions.Parse(o.Action)
	if action == nil {
		return nil, fmt.Errorf("policy rule %q has unknown action %q", o.Id, o.Action)
	}

	protocol := PolicyRuleProtocols.Parse(o.Protocol)
	if protocol == nil {
		return nil, fmt.Errorf("policy rule %q has unknown protocol %q", o.Id, o.Protocol)
	}

	var tcpStateQualifier *TcpStateQualifier
	if o.TcpStateQualifier != nil {
		tcpStateQualifier = TcpStateQualifiers.Parse(*o.TcpStateQualifier)
		if tcpStateQualifier == nil {
			return nil, fmt.Errorf("cannot parse policy rule %q tcp state qualifier: %q", o.Id, *o.TcpStateQualifier)
		}
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
		Id: o.Id,
		Data: &PolicyRuleData{
			Label:             o.Label,
			Description:       o.Description,
			Protocol:          *protocol,
			Action:            *action,
			SrcPort:           srcPort,
			DstPort:           dstPort,
			TcpStateQualifier: tcpStateQualifier,
		},
	}, nil
}

func (o *TwoStageL3ClosClient) getPolicyRuleIdByLabel(ctx context.Context, policyId ObjectId, label string) (ObjectId, error) {
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

func (o *TwoStageL3ClosClient) addPolicyRule(ctx context.Context, rule *rawPolicyRule, position int, policyId ObjectId) (ObjectId, error) {
	// ensure exclusive access to the policy while we recalculate the rules
	lockId := o.lockId(policyId)
	o.client.lock(lockId)
	defer o.client.unlock(lockId)

	policy, err := o.getPolicy(ctx, policyId)
	if err != nil {
		return "", err
	}

	currentRuleCount := len(policy.Rules)

	if position < 0 {
		position = currentRuleCount
	}

	switch {
	case currentRuleCount == 0:
		// empty rule set is an easy case
		policy.Rules = []rawPolicyRule{*rule}
	case position == 0:
		// insert at the beginning
		policy.Rules = append([]rawPolicyRule{*rule}, policy.Rules...)
	case position >= currentRuleCount:
		// insert at the end
		policy.Rules = append(policy.Rules, *rule)
	default:
		// insert somewhere in the middle
		policy.Rules = append(policy.Rules[:position+1], policy.Rules[position:]...)
		policy.Rules[position] = *rule
	}

	// push the new policy
	err = o.updatePolicy(ctx, policyId, policy.request())
	if err != nil {
		return "", err
	}

	return o.getPolicyRuleIdByLabel(ctx, policyId, rule.Label)
}

func (o *TwoStageL3ClosClient) deletePolicyRuleById(ctx context.Context, policyId ObjectId, ruleId ObjectId) error {
	// ensure exclusive access to the policy while we recalculate the rules
	lockId := o.lockId(policyId)
	o.client.lock(lockId)
	defer o.client.unlock(lockId)

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
		return ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("rule id '%s' not found in policy '%s'", ruleId, policyId),
		}
	}

	policy.Rules = append(policy.Rules[:ruleIdx], policy.Rules[ruleIdx+1:]...)
	return o.updatePolicy(ctx, policyId, policy.request())
}
