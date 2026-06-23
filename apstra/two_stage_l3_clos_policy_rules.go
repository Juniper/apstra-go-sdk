// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"time"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
)

type PolicyRule struct {
	Id   ObjectId
	Data *PolicyRuleData
}

type PolicyRuleData struct {
	Label             string
	Description       string
	Protocol          enum.PolicyRuleProtocol
	Action            enum.PolicyRuleAction
	SrcPort           datacenter.PortRanges
	DstPort           datacenter.PortRanges
	TcpStateQualifier *enum.TcpStateQualifier
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
		SrcPort:           o.SrcPort,
		DstPort:           o.DstPort,
		TcpStateQualifier: tcpStateQualifier,
	}
}

type rawPolicyRule struct {
	Id                ObjectId              `json:"id,omitempty"`
	Label             string                `json:"label"`
	Description       string                `json:"description"`
	Protocol          string                `json:"protocol"`
	Action            string                `json:"action"`
	SrcPort           datacenter.PortRanges `json:"src_port"`
	DstPort           datacenter.PortRanges `json:"dst_port"`
	TcpStateQualifier *string               `json:"tcp_state_qualifier,omitempty"`
}

func (o rawPolicyRule) polish() (*PolicyRule, error) {
	action := enum.PolicyRuleActions.Parse(o.Action)
	if action == nil {
		return nil, fmt.Errorf("policy rule %q has unknown action %q", o.Id, o.Action)
	}

	protocol := enum.PolicyRuleProtocols.Parse(o.Protocol)
	if protocol == nil {
		return nil, fmt.Errorf("policy rule %q has unknown protocol %q", o.Id, o.Protocol)
	}

	var tcpStateQualifier *enum.TcpStateQualifier
	if o.TcpStateQualifier != nil {
		tcpStateQualifier = enum.TcpStateQualifiers.Parse(*o.TcpStateQualifier)
		if tcpStateQualifier == nil {
			return nil, fmt.Errorf("cannot parse policy rule %q tcp state qualifier: %q", o.Id, *o.TcpStateQualifier)
		}
	}

	return &PolicyRule{
		Id: o.Id,
		Data: &PolicyRuleData{
			Label:             o.Label,
			Description:       o.Description,
			Protocol:          *protocol,
			Action:            *action,
			SrcPort:           o.SrcPort,
			DstPort:           o.DstPort,
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
