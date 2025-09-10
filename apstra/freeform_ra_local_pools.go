// Copyright (c) Juniper Networks, Inc., 2024-2024.
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
	apiUrlFfRaLocalPools    = apiUrlBlueprintById + apiUrlPathDelim + "ra-local-pools"
	apiUrlFfRaLocalPoolById = apiUrlFfRaLocalPools + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformRaLocalIntPool)

type FreeformRaLocalIntPool struct {
	Id   ObjectId
	Data *FreeformRaLocalIntPoolData
}

func (o *FreeformRaLocalIntPool) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id    ObjectId `json:"id"`
		Label string   `json:"label"`
		// PoolType     string    `json:"pool_type"` // not used because always (?) "integer"
		ResourceType string    `json:"resource_type"`
		OwnerId      ObjectId  `json:"owner_id"`
		GeneratorId  *ObjectId `json:"generator_id"`
		Definition   struct {
			Chunks []FFLocalIntPoolChunk `json:"chunks"`
		} `json:"definition"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(FreeformRaLocalIntPoolData)
	o.Data.Label = raw.Label
	err = o.Data.ResourceType.FromString(raw.ResourceType)
	if err != nil {
		return err
	}
	o.Data.OwnerId = raw.OwnerId
	o.Data.GeneratorId = raw.GeneratorId
	o.Data.Chunks = raw.Definition.Chunks

	return err
}

var _ json.Marshaler = new(FreeformRaLocalIntPoolData)

type FreeformRaLocalIntPoolData struct {
	ResourceType enum.FFResourceType
	Label        string
	OwnerId      ObjectId
	GeneratorId  *ObjectId
	Chunks       []FFLocalIntPoolChunk
}

func (o FreeformRaLocalIntPoolData) MarshalJSON() ([]byte, error) {
	var raw struct {
		Label        string    `json:"label"`
		PoolType     string    `json:"pool_type"`
		OwnerId      ObjectId  `json:"owner_id"`
		GeneratorId  *ObjectId `json:"generator_id"`
		ResourceType string    `json:"resource_type"`
		Definition   struct {
			Chunks []FFLocalIntPoolChunk `json:"chunks"`
		} `json:"definition"`
	}

	raw.Label = o.Label
	raw.OwnerId = o.OwnerId
	raw.PoolType = "integer"
	raw.GeneratorId = o.GeneratorId
	raw.ResourceType = o.ResourceType.String()
	raw.Definition.Chunks = o.Chunks

	return json.Marshal(&raw)
}

type FFLocalIntPoolChunk struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func (o *FreeformClient) CreateRaLocalIntPool(ctx context.Context, in *FreeformRaLocalIntPoolData) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFfRaLocalPools, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *FreeformClient) GetAllRaLocalIntPools(ctx context.Context) ([]FreeformRaLocalIntPool, error) {
	var response struct {
		Items []FreeformRaLocalIntPool `json:"items"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfRaLocalPools, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) GetRaLocalIntPool(ctx context.Context, id ObjectId) (*FreeformRaLocalIntPool, error) {
	var response FreeformRaLocalIntPool

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfRaLocalPoolById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) UpdateRaLocalIntPool(ctx context.Context, id ObjectId, in *FreeformRaLocalIntPoolData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFfRaLocalPoolById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) DeleteRaLocalIntPool(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfRaLocalPoolById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
