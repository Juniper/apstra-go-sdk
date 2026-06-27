// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"time"

	"github.com/Juniper/apstra-go-sdk/datacenter"
)

func (o *TwoStageL3ClosClient) getPolicyRuleIDByLabel(ctx context.Context, policyID string, label string) (string, error) {
	start := time.Now()
	for i := 0; i <= dcClientMaxRetries; i++ {
		time.Sleep(dcClientRetryBackoff * time.Duration(i))
		policy, err := o.GetPolicy(ctx, policyID)
		if err != nil {
			return "", err
		}
		for i, rule := range policy.Rules {
			if rule.Label == label {
				idPtr := rule.ID()
				if idPtr == nil {
					return "", fmt.Errorf("policy rule at index %d has label %q, but id is nil", i, label)
				}
				return *idPtr, nil
			}
		}
	}
	return "", fmt.Errorf("rule '%s' didn't appear in policy '%s' after %s", label, policyID, time.Since(start))
}

// AddPolicyRule adds a policy rule at 'position' (bumping all other rules
// down). Position 0 makes the new policy first on the list, 1 makes it second
// on the list, etc... Use -1 for last on the list. The returned string is
// the ID of the new rule within the policy.
func (o *TwoStageL3ClosClient) AddPolicyRule(ctx context.Context, rule datacenter.PolicyRule, position int, policyID string) (string, error) {
	// ensure exclusive access to the policy while we recalculate the rules
	lockId := o.lockId(ObjectId(policyID))
	o.client.lock(lockId)
	defer o.client.unlock(lockId)

	policy, err := o.GetPolicy(ctx, policyID)
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
		policy.Rules = []datacenter.PolicyRule{rule}
	case position == 0:
		// insert at the beginning
		policy.Rules = append([]datacenter.PolicyRule{rule}, policy.Rules...)
	case position >= currentRuleCount:
		// insert at the end
		policy.Rules = append(policy.Rules, rule)
	default:
		// insert somewhere in the middle
		policy.Rules = append(policy.Rules[:position+1], policy.Rules[position:]...)
		policy.Rules[position] = rule
	}

	// push the new policy
	err = o.UpdatePolicy(ctx, policy)
	if err != nil {
		return "", err
	}

	return o.getPolicyRuleIDByLabel(ctx, policyID, rule.Label)
}

// DeletePolicyRuleById deletes the given rule. If the rule doesn't exist, a
// ClientErr with ErrNotFound is returned.
func (o *TwoStageL3ClosClient) DeletePolicyRuleByID(ctx context.Context, policyID string, ruleID string) error {
	// ensure exclusive access to the policy while we recalculate the rules
	lockId := o.lockId(ObjectId(policyID))
	o.client.lock(lockId)
	defer o.client.unlock(lockId)

	policy, err := o.GetPolicy(ctx, policyID)
	if err != nil {
		return err
	}

	ruleIdx := -1
	for i, rule := range policy.Rules {
		idPtr := rule.ID()
		if idPtr == nil {
			continue
		}
		if *idPtr == ruleID {
			ruleIdx = i
			break
		}
	}

	if ruleIdx < 0 {
		return ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("rule id '%s' not found in policy '%s'", ruleID, policyID),
		}
	}

	policy.Rules = append(policy.Rules[:ruleIdx], policy.Rules[ruleIdx+1:]...)
	return o.UpdatePolicy(ctx, policy)
}
