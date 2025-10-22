// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/internal/str"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
)

func (c Client) CreateRackType2(ctx context.Context, v design.RackType) (string, error) {
	if v.ID() != nil {
		return "", fmt.Errorf("id must be nil in %s", str.FuncName())
	}

	var response struct {
		ID string `json:"id"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      design.RackTypesURL,
		apiInput:    v,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.ID, nil
}

func (c Client) GetRackType2(ctx context.Context, id string) (design.RackType, error) {
	var response design.RackType
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(design.RackTypeURLByID, id),
		apiResponse: &response,
	})
	if err != nil {
		return response, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (c Client) UpdateRackType2(ctx context.Context, v design.RackType) error {
	if v.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(design.RackTypeURLByID, *v.ID()),
		apiInput: &v,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) DeleteRackType2(ctx context.Context, id string) error {
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(design.RackTypeURLByID, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) ListRackTypes2(ctx context.Context) ([]string, error) {
	var response struct {
		Items []string `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      design.RackTypesURL,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (c Client) GetRackTypes2(ctx context.Context) ([]design.RackType, error) {
	var response struct {
		Items []design.RackType `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      design.RackTypesURL,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (c Client) GetRackTypeByLabel2(ctx context.Context, label string) (design.RackType, error) {
	var result []design.RackType

	all, err := c.GetRackTypes2(ctx)
	if err != nil {
		return zero.SliceItem(result), fmt.Errorf("%s failed getting all %T candidates: %w", str.FuncName(), zero.SliceItem(result), err)
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
			err:     fmt.Errorf("found multiple candidate %T with label %s", zero.SliceItem(result), label),
		}
	}
}
