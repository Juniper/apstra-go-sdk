// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiUrlBlueprintPropertySets = apiUrlBlueprintById + apiUrlPathDelim + "property-sets"
	apiUrlBlueprintPropertySet  = apiUrlBlueprintById + apiUrlPathDelim + "property-sets" + apiUrlPathDelim + "%s"
)

type TwoStageL3ClosPropertySet struct {
	Id         ObjectId        `json:"id"`
	Label      string          `json:"label"`
	Stale      bool            `json:"stale"`
	Values     json.RawMessage `json:"values"`
	ValuesYaml string          `json:"values_yaml"`
}

type importPropertySetRequest struct {
	Id   ObjectId `json:"id"`
	Keys []string `json:"keys,omitempty"`
}

type importPropertySetResponse struct {
	Id ObjectId `json:"id"`
}

func (o *TwoStageL3ClosClient) importPropertySet(ctx context.Context, psid ObjectId, keys ...string) (ObjectId, error) {
	importPropertySetUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintPropertySets, o.blueprintId.String()))
	if err != nil {
		return "", err
	}

	response := &importPropertySetResponse{}
	err = o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodPost,
		url:         importPropertySetUrl,
		apiInput:    importPropertySetRequest{Id: psid, Keys: keys},
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) getAllPropertySets(ctx context.Context) ([]TwoStageL3ClosPropertySet, error) {
	result := &struct {
		Items []TwoStageL3ClosPropertySet `json:"items"`
	}{}
	allImportedPropertySetUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintPropertySets, o.blueprintId.String()))
	if err != nil {
		return nil, err
	}
	err = o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		url:         allImportedPropertySetUrl,
		apiResponse: result,
	})
	return result.Items, convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) getPropertySetByName(ctx context.Context, name string) (*TwoStageL3ClosPropertySet, error) {
	allps, err := o.getAllPropertySets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	for _, t := range allps {
		if t.Label == name {
			return &t, nil
		}
	}
	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("property Set with name '%s' not found", name),
	}
}

func (o *TwoStageL3ClosClient) getPropertySet(ctx context.Context, psid ObjectId) (*TwoStageL3ClosPropertySet, error) {
	result := &TwoStageL3ClosPropertySet{}

	getImportedPropertySetUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintPropertySet, o.blueprintId.String(), psid.String()))
	err = o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		url:         getImportedPropertySetUrl,
		apiResponse: result,
	})
	return result, convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) updatePropertySet(ctx context.Context, psid ObjectId, keys ...string) error {
	updateImportedPropertySetUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintPropertySet, o.blueprintId.String(), psid.String()))
	if err != nil {
		return err
	}

	err = o.client.talkToApstra(ctx, talkToApstraIn{
		method: http.MethodPut,
		url:    updateImportedPropertySetUrl,
		apiInput: importPropertySetRequest{
			Id:   psid,
			Keys: keys,
		},
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *TwoStageL3ClosClient) deletePropertySet(ctx context.Context, pid ObjectId) error {
	deleteImportedPropertySetUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintPropertySet, o.blueprintId.String(), pid.String()))
	if err != nil {
		return err
	}
	err = o.client.talkToApstra(ctx, talkToApstraIn{
		method: http.MethodDelete,
		url:    deleteImportedPropertySetUrl,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
