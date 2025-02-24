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
	apiUrlFfGroupGenerators    = apiUrlBlueprintById + apiUrlPathDelim + "ra-group-generators"
	apiUrlFfGroupGeneratorById = apiUrlFfGroupGenerators + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformGroupGenerator)

type FreeformGroupGenerator struct {
	Id   ObjectId
	Data *FreeformGroupGeneratorData
}

func (o *FreeformGroupGenerator) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id       ObjectId  `json:"id"`
		ParentId *ObjectId `json:"parent_id"`
		Label    string    `json:"label"`
		Scope    string    `json:"scope"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(FreeformGroupGeneratorData)
	o.Data.ParentId = raw.ParentId
	o.Data.Label = raw.Label
	o.Data.Scope = raw.Scope

	return nil
}

var _ json.Marshaler = new(FreeformGroupGeneratorData)

type FreeformGroupGeneratorData struct {
	ParentId *ObjectId
	Label    string
	Scope    string
}

func (o FreeformGroupGeneratorData) MarshalJSON() ([]byte, error) {
	var raw struct {
		ParentId *ObjectId `json:"parent_id"`
		Label    string    `json:"label"`
		Scope    string    `json:"scope"`
	}

	raw.ParentId = o.ParentId
	raw.Label = o.Label
	raw.Scope = o.Scope

	return json.Marshal(&raw)
}

func (o *FreeformClient) CreateGroupGenerator(ctx context.Context, in *FreeformGroupGeneratorData) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFfGroupGenerators, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *FreeformClient) GetAllGroupGenerators(ctx context.Context) ([]FreeformGroupGenerator, error) {
	var response struct {
		Items []FreeformGroupGenerator `json:"items"`
	}

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfGroupGenerators, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) GetGroupGenerator(ctx context.Context, id ObjectId) (*FreeformGroupGenerator, error) {
	var response FreeformGroupGenerator

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfGroupGeneratorById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) GetGroupGeneratorByName(ctx context.Context, name string) (*FreeformGroupGenerator, error) {
	all, err := o.GetAllGroupGenerators(ctx)
	if err != nil {
		return nil, err
	}

	var result *FreeformGroupGenerator
	for _, ffGrpGen := range all {
		ffGrpGen := ffGrpGen
		if ffGrpGen.Data.Label == name {
			if result != nil {
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple freeform group generators in blueprint %q have name %q", o.client.id, name),
				}
			}

			result = &ffGrpGen
		}
	}

	if result == nil {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no freeform group generator in blueprint %q has name %q", o.client.id, name),
		}
	}

	return result, nil
}

func (o *FreeformClient) UpdateGroupGenerator(ctx context.Context, id ObjectId, in *FreeformGroupGeneratorData) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFfGroupGeneratorById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) DeleteGroupGenerator(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfGroupGeneratorById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
