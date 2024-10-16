// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"net/http"
)

const (
	apiUrlBlueprintNodeConfigRender   = apiUrlBlueprintNodeByIdPrefix + "config-rendering?type=%s"
	apiUrlBlueprintSystemConfigRender = apiUrlBlueprintSystemByIdPrefix + "config-rendering?type=%s"
)

func (o *TwoStageL3ClosClient) GetNodeRenderedConfig(ctx context.Context, id ObjectId, rcType enum.RenderedConfigType) (string, error) {
	var apiResponse struct {
		Config string `json:"config"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintNodeConfigRender, o.blueprintId, id, rcType.String()),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return apiResponse.Config, nil
}

func (o *TwoStageL3ClosClient) GetSystemRenderedConfig(ctx context.Context, id ObjectId, rcType enum.RenderedConfigType) (string, error) {
	var apiResponse struct {
		Config string `json:"config"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSystemConfigRender, o.blueprintId, id, rcType.String()),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return apiResponse.Config, nil
}
