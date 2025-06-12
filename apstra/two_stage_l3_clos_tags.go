// Copyright (c) Juniper Networks, Inc., 2023-2025.
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
	apiUrlBlueprintTags    = apiUrlBlueprintByIdPrefix + "tags"
	apiUrlBlueprintTagById = apiUrlBlueprintTags + apiUrlPathDelim + "%s"
)

type TwoStageL3ClosTagData struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}

var _ json.Unmarshaler = (*TwoStageL3ClosTag)(nil)

type TwoStageL3ClosTag struct {
	Id   ObjectId               `json:"id"`
	Data *TwoStageL3ClosTagData `json:"data"`
}

func (o *TwoStageL3ClosTag) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id          ObjectId `json:"id"`
		Label       string   `json:"label"`
		Description string   `json:"description"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshal tag data: %w", err)
	}

	o.Id = raw.Id
	o.Data = new(TwoStageL3ClosTagData)
	o.Data.Label = raw.Label
	o.Data.Description = raw.Description

	return nil
}

func (o TwoStageL3ClosClient) CreateTag(ctx context.Context, in TwoStageL3ClosTagData) (ObjectId, error) {
	var apiResponse objectIdResponse

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintTags, o.Id()),
		apiInput:    &in,
		apiResponse: &apiResponse,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return apiResponse.Id, nil
}

func (o TwoStageL3ClosClient) GetTag(ctx context.Context, id ObjectId) (TwoStageL3ClosTag, error) {
	var apiResponse TwoStageL3ClosTag

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintTagById, o.Id(), id),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return TwoStageL3ClosTag{}, convertTtaeToAceWherePossible(err)
	}

	return apiResponse, nil
}

func (o TwoStageL3ClosClient) GetAllTags(ctx context.Context) (map[ObjectId]TwoStageL3ClosTag, error) {
	var apiResponse struct {
		Items []TwoStageL3ClosTag `json:"items"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintTags, o.Id()),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make(map[ObjectId]TwoStageL3ClosTag, len(apiResponse.Items))
	for _, item := range apiResponse.Items {
		result[item.Id] = item
	}

	return result, nil
}

func (o TwoStageL3ClosClient) UpdateTag(ctx context.Context, id ObjectId, in TwoStageL3ClosTagData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintTagById, o.Id(), id),
		apiInput: &in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o TwoStageL3ClosClient) DeleteTag(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintTagById, o.Id(), id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
