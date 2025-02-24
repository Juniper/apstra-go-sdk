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
	apiUrlBlueprintConfiglets       = apiUrlBlueprintById + apiUrlPathDelim + "configlets"
	apiUrlBlueprintConfigletsPrefix = apiUrlBlueprintConfiglets + apiUrlPathDelim
	apiUrlBlueprintConfigletsById   = apiUrlBlueprintConfigletsPrefix + "%s"
)

type TwoStageL3ClosConfigletData struct {
	Label     string         `json:"label"`
	Condition string         `json:"condition"`
	Data      *ConfigletData `json:"configlet"`
}

var _ json.Unmarshaler = (*TwoStageL3ClosConfiglet)(nil)

type TwoStageL3ClosConfiglet struct {
	Id   ObjectId
	Data *TwoStageL3ClosConfigletData
}

func (o *TwoStageL3ClosConfiglet) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id        ObjectId       `json:"id"`
		Condition string         `json:"condition"`
		Label     string         `json:"label"`
		Data      *ConfigletData `json:"configlet"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = &TwoStageL3ClosConfigletData{
		Label:     raw.Label,
		Condition: raw.Condition,
		Data:      raw.Data,
	}

	return nil
}

func (o *TwoStageL3ClosClient) getAllConfiglets(ctx context.Context) ([]TwoStageL3ClosConfiglet, error) {
	var response struct {
		Items []TwoStageL3ClosConfiglet `json:"items"`
	}

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintConfiglets, o.blueprintId.String()),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *TwoStageL3ClosClient) getAllConfigletIds(ctx context.Context) ([]ObjectId, error) {
	configlets, err := o.getAllConfiglets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	ids := make([]ObjectId, len(configlets))
	for i, c := range configlets {
		ids[i] = c.Id
	}

	return ids, nil
}

func (o *TwoStageL3ClosClient) getConfiglet(ctx context.Context, id ObjectId) (*TwoStageL3ClosConfiglet, error) {
	var response TwoStageL3ClosConfiglet
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintConfigletsById, o.blueprintId.String(), id.String()),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *TwoStageL3ClosClient) getConfigletByName(ctx context.Context, name string) (*TwoStageL3ClosConfiglet, error) {
	cgs, err := o.getAllConfiglets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	idx := -1
	for i, t := range cgs {
		if t.Data.Label == name {
			if idx == -1 {
				idx = i
			} else { // This is clearly the second occurrence
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("name '%s' does not uniquely identify a configlet", name),
				}
			}
		}
	}
	if idx != -1 {
		return &cgs[idx], nil
	}
	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no Configlet with name '%s' found", name),
	}
}

func (o *TwoStageL3ClosClient) createConfiglet(ctx context.Context, in *TwoStageL3ClosConfigletData) (ObjectId, error) {
	response := &objectIdResponse{}
	if len(in.Label) == 0 {
		in.Label = in.Data.DisplayName
	}

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintConfiglets, o.blueprintId.String()),
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *TwoStageL3ClosClient) updateConfiglet(ctx context.Context, id ObjectId, in *TwoStageL3ClosConfigletData) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintConfigletsById, o.blueprintId.String(), id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) deleteConfiglet(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintConfigletsById, o.blueprintId.String(), id.String()),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
