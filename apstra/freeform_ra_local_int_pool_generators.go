// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	apiUrlFfRaLocalPoolGenerators    = apiUrlBlueprintById + apiUrlPathDelim + "ra-local-pool-generators"
	apiUrlFfRaLocalPoolGeneratorById = apiUrlFfRaLocalPoolGenerators + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformRaLocalIntPoolGenerator)

type FreeformRaLocalIntPoolGenerator struct {
	Id   ObjectId
	Data *FreeformRaLocalIntPoolGeneratorData
}

func (o *FreeformRaLocalIntPoolGenerator) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id           ObjectId `json:"id"`
		Label        string   `json:"label"`
		Scope        string   `json:"scope"`
		PoolType     string   `json:"pool_type"`
		ResourceType string   `json:"resource_type"`
		Definition   struct {
			Chunks []FFLocalIntPoolChunk `json:"chunks"`
		} `json:"definition"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	if raw.PoolType != "integer" {
		return errors.New("pool type mismatch")
	}

	o.Id = raw.Id
	o.Data = new(FreeformRaLocalIntPoolGeneratorData)
	o.Data.Label = raw.Label
	o.Data.Scope = raw.Scope
	o.Data.Chunks = raw.Definition.Chunks
	err = o.Data.ResourceType.FromString(raw.ResourceType)
	if err != nil {
		return err
	}

	return nil
}

type FreeformRaLocalIntPoolGeneratorData struct {
	ResourceType enum.FFResourceType
	Label        string
	Scope        string
	Chunks       []FFLocalIntPoolChunk
}

var _ json.Marshaler = new(FreeformRaLocalIntPoolGeneratorData)

func (o FreeformRaLocalIntPoolGeneratorData) MarshalJSON() ([]byte, error) {
	var raw struct {
		Label        string `json:"label"`
		Scope        string `json:"scope"`
		PoolType     string `json:"pool_type"`
		ResourceType string `json:"resource_type"`
		Definition   struct {
			Chunks []FFLocalIntPoolChunk `json:"chunks,omitempty"`
		} `json:"definition"`
	}

	raw.Label = o.Label
	raw.Scope = o.Scope
	raw.PoolType = "integer"
	raw.ResourceType = o.ResourceType.String()
	raw.Definition.Chunks = o.Chunks

	return json.Marshal(&raw)
}

func (o *FreeformClient) CreateRaLocalIntPoolGenerator(ctx context.Context, in *FreeformRaLocalIntPoolGeneratorData) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFfRaLocalPoolGenerators, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *FreeformClient) GetAllLocalIntPoolGenerators(ctx context.Context) ([]FreeformRaLocalIntPoolGenerator, error) {
	var response struct {
		Items []FreeformRaLocalIntPoolGenerator `json:"items"`
	}

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfRaLocalPoolGenerators, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) GetRaLocalIntPoolGenerator(ctx context.Context, id ObjectId) (*FreeformRaLocalIntPoolGenerator, error) {
	var response FreeformRaLocalIntPoolGenerator

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfRaLocalPoolGeneratorById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) UpdateRaLocalIntPoolGenerator(ctx context.Context, id ObjectId, in *FreeformRaLocalIntPoolGeneratorData) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFfRaLocalPoolGeneratorById, o.blueprintId, id),
		apiInput: &in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) DeleteRaLocalPoolGenerator(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfRaLocalPoolGeneratorById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
