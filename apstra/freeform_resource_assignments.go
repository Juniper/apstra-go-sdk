// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiURLFfResourceAssignments = apiUrlFfRaResources + apiUrlPathDelim + "%s" + apiUrlPathDelim + "assignments"
)

func (o *FreeformClient) ListResourceAssignments(ctx context.Context, resource ObjectId) ([]ObjectId, error) {
	var response struct {
		Targets []struct {
			Id ObjectId `json:"node_id"`
		} `json:"assigned_to"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiURLFfResourceAssignments, o.blueprintId, resource),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make([]ObjectId, len(response.Targets))
	for i, target := range response.Targets {
		result[i] = target.Id
	}

	return result, nil
}

func (o *FreeformClient) UpdateResourceAssignments(ctx context.Context, resource ObjectId, targets []ObjectId) error {
	var apiInput struct {
		Targets []ObjectId `json:"assigned_to"`
	}

	if targets == nil {
		apiInput.Targets = []ObjectId{}
	} else {
		apiInput.Targets = targets
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiURLFfResourceAssignments, o.blueprintId, resource),
		apiInput: apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
