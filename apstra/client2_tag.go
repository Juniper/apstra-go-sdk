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

func (c Client) CreateTag2(ctx context.Context, v design.Tag) (string, error) {
	if v.ID() != nil {
		return "", fmt.Errorf("id must be nil in %s", str.FuncName())
	}

	var response struct {
		ID string `json:"id"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      design.TagsURL,
		apiInput:    v,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.ID, nil
}

func (c Client) GetTag2(ctx context.Context, id string) (design.Tag, error) {
	var response design.Tag
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(design.TagURLByID, id),
		apiResponse: &response,
	})
	if err != nil {
		return response, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (c Client) UpdateTag2(ctx context.Context, v design.Tag) error {
	if v.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(design.TagURLByID, *v.ID()),
		apiInput: &v,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) DeleteTag2(ctx context.Context, id string) error {
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(design.TagURLByID, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) ListTags2(ctx context.Context) ([]string, error) {
	var response struct {
		Items []string `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      design.TagsURL,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (c Client) GetTags2(ctx context.Context) ([]design.Tag, error) {
	var response struct {
		Items []design.Tag `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      design.TagsURL,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (c Client) GetTagByLabel2(ctx context.Context, label string) (design.Tag, error) {
	var result []design.Tag

	all, err := c.GetTags2(ctx)
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
