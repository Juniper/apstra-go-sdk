// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
)

const (
	apiUrlBlueprintAntiAffinityPolicy = apiUrlBlueprintByIdPrefix + "anti-affinity-policy"
)

// getAntiAffinityPolicy is for Apstra 4.2.0 and earlier (not available in 4.2.1)
func (o *TwoStageL3ClosClient) getAntiAffinityPolicy(ctx context.Context) (*rawAntiAffinityPolicy, error) {
	if !compatibility.EqApstra420.Check(o.client.apiVersion) {
		return nil, fmt.Errorf("apstra %s does not support %q", o.client.apiVersion, apiUrlBlueprintAntiAffinityPolicy)
	}

	var result rawAntiAffinityPolicy
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintAntiAffinityPolicy, o.blueprintId),
		apiResponse: &result,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &result, nil
}

// setAntiAffinityPolicy is for Apstra 4.2.0 and earlier (not available in 4.2.1)
func (o *TwoStageL3ClosClient) setAntiAffinityPolicy(ctx context.Context, in *rawAntiAffinityPolicy) error {
	if in == nil {
		return nil
	}

	if !compatibility.EqApstra420.Check(o.client.apiVersion) {
		return fmt.Errorf("apstra %s does not support %q", o.client.apiVersion, apiUrlBlueprintAntiAffinityPolicy)
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintAntiAffinityPolicy, o.blueprintId),
		apiInput: &in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
