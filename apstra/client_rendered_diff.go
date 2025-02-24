// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintNodeConfigRenderDiff   = apiUrlBlueprintNodeByIdPrefix + "config-incremental"
	apiUrlBlueprintSystemConfigRenderDiff = apiUrlBlueprintSystemByIdPrefix + "config-incremental"
)

type RenderDiff struct {
	Config             string          `json:"config"`
	PristineConfig     json.RawMessage `json:"pristine_config"`
	Context            json.RawMessage `json:"context"`
	SupportsDiffConfig bool            `json:"supports_diff_config"`
}

func (o *Client) GetNodeRenderedConfigDiff(ctx context.Context, bpId, nodeId ObjectId) (*RenderDiff, error) {
	var apiResponse RenderDiff

	err := o.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintNodeConfigRenderDiff, bpId, nodeId),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &apiResponse, nil
}

func (o *Client) GetSystemRenderedConfigDiff(ctx context.Context, bpId, sysId ObjectId) (*RenderDiff, error) {
	var apiResponse RenderDiff

	err := o.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSystemConfigRenderDiff, bpId, sysId),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &apiResponse, nil
}
