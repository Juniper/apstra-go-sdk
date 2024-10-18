// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	apiUrlBlueprintNodeConfigRender   = apiUrlBlueprintNodeByIdPrefix + "config-rendering?type=%s"
	apiUrlBlueprintSystemConfigRender = apiUrlBlueprintSystemByIdPrefix + "config-rendering?type=%s"
)

func (o *Client) GetNodeRenderedConfig(ctx context.Context, bpId, nodeId ObjectId, rcType enum.RenderedConfigType) (string, error) {
	var apiResponse struct {
		Config string `json:"config"`
	}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintNodeConfigRender, bpId, nodeId, rcType.String()),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return apiResponse.Config, nil
}

func (o *Client) GetSystemRenderedConfig(ctx context.Context, bpId, sysId ObjectId, rcType enum.RenderedConfigType) (string, error) {
	var apiResponse struct {
		Config string `json:"config"`
	}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSystemConfigRender, bpId, sysId, rcType.String()),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return apiResponse.Config, nil
}
