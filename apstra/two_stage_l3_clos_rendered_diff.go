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

func (o *TwoStageL3ClosClient) GetNodeRenderedConfigDiff(ctx context.Context, id ObjectId) (*RenderDiff, error) {
	var apiResponse RenderDiff

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintNodeConfigRenderDiff, o.blueprintId, id),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &apiResponse, nil
}

func (o *TwoStageL3ClosClient) GetSystemRenderedConfigDiff(ctx context.Context, id ObjectId) (*RenderDiff, error) {
	var apiResponse RenderDiff

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSystemConfigRenderDiff, o.blueprintId, id),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &apiResponse, nil
}
