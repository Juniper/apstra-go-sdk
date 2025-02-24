// Copyright (c) Juniper Networks, Inc., 2024-2025.
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
	apiUrlFfRaGroups    = apiUrlBlueprintById + apiUrlPathDelim + "ra-groups"
	apiUrlFfRaGroupById = apiUrlFfRaGroups + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformRaGroup)

type FreeformRaGroup struct {
	Id   ObjectId
	Data *FreeformRaGroupData
}

func (o *FreeformRaGroup) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id          ObjectId        `json:"id"`
		ParentId    *ObjectId       `json:"parent_id"`
		Label       string          `json:"label"`
		Tags        []string        `json:"tags"`
		Data        json.RawMessage `json:"data"`
		GeneratorId *ObjectId       `json:"generator_id"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	if o.Data == nil {
		o.Data = new(FreeformRaGroupData)
	}

	o.Id = raw.Id
	o.Data.ParentId = raw.ParentId
	o.Data.Label = raw.Label
	o.Data.Tags = raw.Tags
	o.Data.Data = raw.Data
	o.Data.GeneratorId = raw.GeneratorId

	return err
}

type FreeformRaGroupData struct {
	ParentId    *ObjectId       `json:"parent_id"`
	Label       string          `json:"label"`
	Tags        []string        `json:"tags"`
	Data        json.RawMessage `json:"data"`
	GeneratorId *ObjectId       `json:"generator_id"`
}

func (o *FreeformClient) CreateRaGroup(ctx context.Context, in *FreeformRaGroupData) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFfRaGroups, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *FreeformClient) GetAllRaGroups(ctx context.Context) ([]FreeformRaGroup, error) {
	var response struct {
		Items []FreeformRaGroup `json:"items"`
	}

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfRaGroups, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) GetRaGroup(ctx context.Context, id ObjectId) (*FreeformRaGroup, error) {
	var response FreeformRaGroup

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfRaGroupById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) GetRaGroupByName(ctx context.Context, name string) (*FreeformRaGroup, error) {
	all, err := o.GetAllRaGroups(ctx)
	if err != nil {
		return nil, err
	}

	var result *FreeformRaGroup
	for _, ffrag := range all {
		ffrag := ffrag
		if ffrag.Data.Label == name {
			if result != nil {
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple resource allocation groups in blueprint %q have name %q", o.client.id, name),
				}
			}

			result = &ffrag
		}
	}

	if result == nil {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no freeform resource allocation group  in blueprint %q has name %q", o.client.id, name),
		}
	}

	return result, nil
}

func (o *FreeformClient) UpdateRaGroup(ctx context.Context, id ObjectId, in *FreeformRaGroupData) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFfRaGroupById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) DeleteRaGroup(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfRaGroupById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
