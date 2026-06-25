// Copyright (c) Juniper Networks, Inc., 2022-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	"github.com/Juniper/apstra-go-sdk/internal/str"
	"github.com/Juniper/apstra-go-sdk/internal/urls"
)

// CreatePolicy creates a policy within the DC blueprint, returns its ID
func (o *TwoStageL3ClosClient) CreatePolicy(ctx context.Context, v datacenter.Policy) (string, error) {
	var response struct {
		ID string `json:"id"`
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(urls.DatacenterPolicies, o.blueprintId),
		apiInput:    v,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.ID, nil
}

// GetPolicy returns Policy representing policy 'id' within the DC blueprint
func (o *TwoStageL3ClosClient) GetPolicy(ctx context.Context, id string) (datacenter.Policy, error) {
	var response datacenter.Policy
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(urls.DatacenterPolicyByID, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return response, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

// GetPolicies returns []Policy representing all policies configured within the DC blueprint
func (o *TwoStageL3ClosClient) GetPolicies(ctx context.Context) ([]datacenter.Policy, error) {
	response := &struct {
		Policies []datacenter.Policy `json:"policies"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(urls.DatacenterPolicies, o.blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Policies, nil
}

// GetPolicyByLabel returns *Policy representing policy identified by 'label' within the DC blueprint
func (o *TwoStageL3ClosClient) GetPolicyByLabel(ctx context.Context, label string) (datacenter.Policy, error) {
	all, err := o.GetPolicies(ctx)
	if err != nil {
		return slice.ZeroOf(all), err
	}

	var result *datacenter.Policy
	for _, policy := range all {
		if policy.Label == label {
			if result == nil {
				result = &policy
			} else {
				return pointer.ZeroOf(result), ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple matches for %s with label %q", pointer.TypeOf(result), label),
				}
			}
		}
	}

	if result == nil {
		return pointer.ZeroOf(result), ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("%s with label %q not found", pointer.TypeOf(result), label),
		}
	}

	return *result, nil
}

func (o *TwoStageL3ClosClient) ListPolicies(ctx context.Context) ([]string, error) {
	all, err := o.GetPolicies(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(all))
	for i, policy := range all {
		idPtr := policy.ID()
		if idPtr == nil {
			return nil, fmt.Errorf("policy at index %d has nil ID", i)
		}
		result = append(result, *idPtr)
	}

	return result, nil
}

// UpdatePolicy calls PUT to replace the configuration of policy
func (o *TwoStageL3ClosClient) UpdatePolicy(ctx context.Context, v datacenter.Policy) error {
	if v.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	if v.Tags == nil {
		v.Tags = []string{} // sending "null" doesn't clear the tags for some reason
	}
	if v.Rules == nil {
		v.Rules = []datacenter.PolicyRule{} // sending "null" doesn't clear the rules for some reason
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(urls.DatacenterPolicyByID, o.blueprintId, *v.ID()),
		apiInput: v,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

// DeletePolicy deletes policy 'id' within the DC blueprint
func (o *TwoStageL3ClosClient) DeletePolicy(ctx context.Context, id string) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(urls.DatacenterPolicyByID, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) setPolicyEnabled(ctx context.Context, id string, enabled bool) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPatch,
		urlStr: fmt.Sprintf(urls.DatacenterPolicyByID, o.blueprintId, id),
		apiInput: struct {
			Enabled bool `json:"enabled"`
		}{
			Enabled: enabled,
		},
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) EnablePolicy(ctx context.Context, id string) error {
	return o.setPolicyEnabled(ctx, id, true)
}

func (o *TwoStageL3ClosClient) DisablePolicy(ctx context.Context, id string) error {
	return o.setPolicyEnabled(ctx, id, false)
}
