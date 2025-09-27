// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/internal/str"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
	"net/http"
)

func (c Client) CreateInterfaceMap2(ctx context.Context, v design.InterfaceMap) (string, error) {
	if v.ID() != nil {
		return "", fmt.Errorf("id must be nil in %s", str.FuncName())
	}

	var response struct {
		ID string `json:"id"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      design.InterfaceMapsURL,
		apiInput:    v,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.ID, nil
}

func (c Client) GetInterfaceMap2(ctx context.Context, id string) (design.InterfaceMap, error) {
	var response design.InterfaceMap
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(design.InterfaceMapURLByID, id),
		apiResponse: &response,
	})
	if err != nil {
		return response, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (c Client) UpdateInterfaceMap2(ctx context.Context, v design.InterfaceMap) error {
	if v.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(design.InterfaceMapURLByID, *v.ID()),
		apiInput: &v,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) DeleteInterfaceMap2(ctx context.Context, id string) error {
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(design.InterfaceMapURLByID, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) ListInterfaceMaps2(ctx context.Context) ([]string, error) {
	var response struct {
		Items []string `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      design.InterfaceMapsURL,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (c Client) GetInterfaceMaps2(ctx context.Context) ([]design.InterfaceMap, error) {
	var response struct {
		Items []design.InterfaceMap `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      design.InterfaceMapsURL,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (c Client) GetInterfaceMapByLabel2(ctx context.Context, label string) (design.InterfaceMap, error) {
	var result []design.InterfaceMap

	all, err := c.GetInterfaceMaps2(ctx)
	if err != nil {
		return zero.SliceItem(result), fmt.Errorf("%s failed getting all %T candidates: %w", str.FuncName(), result, err)
	}

	for _, item := range all {
		if item.Label == label {
			result = append(result, item)
		}
	}

	switch len(result) {
	case 0:
		return zero.SliceItem(result), ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("%T with label %s not found", zero.SliceItem(result), label),
		}
	case 1:
		return result[0], nil
	default: // len(result) > 1
		return zero.SliceItem(result), ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("found multiple candidate %T with label %s", result, label),
		}
	}
}
