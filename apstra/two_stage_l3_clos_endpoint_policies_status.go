// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/enum"
)

const (
	apiUrlBlueprintExperienceWebEndpointPolicies = apiUrlBlueprintById + apiUrlPathDelim + "experience/web/endpoint-policies"
)

var _ json.Unmarshaler = (*EndpointPolicyStatus)(nil)

type EndpointPolicyStatus struct {
	Id             ObjectId
	Label          string
	Description    string
	Status         enum.EndpointPolicyStatus
	AppPointsCount int
	Tags           []string
	topLevel       bool
}

func (o *EndpointPolicyStatus) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id                ObjectId                  `json:"id"`
		Label             string                    `json:"label"`
		Description       string                    `json:"description"`
		Status            enum.EndpointPolicyStatus `json:"status"`
		AppPointsCount    int                       `json:"app_points_count"`
		Tags              []string                  `json:"tags"`
		TopLevelPolicyIds []ObjectId                `json:"top_level_policy_ids"`
		Visible           bool                      `json:"visible"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshal raw EndpointPolicyStatus: %w", err)
	}
	if raw.AppPointsCount < 0 {
		return fmt.Errorf("app_points_count must be greater than or equal to zero, got %d", raw.AppPointsCount)
	}

	o.Id = raw.Id
	o.Label = raw.Label
	o.Description = raw.Description
	o.Tags = raw.Tags
	o.Status = raw.Status
	o.AppPointsCount = raw.AppPointsCount
	o.Tags = raw.Tags

	if raw.Visible && len(raw.TopLevelPolicyIds) == 0 {
		o.topLevel = true
	} else {
		o.topLevel = false
	}

	return nil
}

// GetAllConnectivityTemplateStatus returns map[ObjectId]EndpointPolicyStatus representing only Connectivity Template
// policy elements. Specifically, those are policy elements with no associated "top level" policy IDs and with
// visible == true.
func (o TwoStageL3ClosClient) GetAllConnectivityTemplateStatus(ctx context.Context) (map[ObjectId]EndpointPolicyStatus, error) {
	var apiResponse struct {
		EndpointPolicies []EndpointPolicyStatus `json:"endpoint_policies"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintExperienceWebEndpointPolicies, o.Id()),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make(map[ObjectId]EndpointPolicyStatus)
	for _, endpointPolicy := range apiResponse.EndpointPolicies {
		if endpointPolicy.topLevel {
			result[endpointPolicy.Id] = endpointPolicy
		}
	}

	return result, nil
}
