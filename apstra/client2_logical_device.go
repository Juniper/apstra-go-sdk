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
)

func (c Client) CreateLogicalDevice2(ctx context.Context, v design.LogicalDevice) (string, error) {
	if v.ID() != nil {
		return "", fmt.Errorf("id must be nil in %s", str.FuncName())
	}

	var response struct {
		ID string `json:"id"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      design.LogicalDevicesUrl,
		apiInput:    v,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.ID, nil
}

func (c Client) GetLogicalDevice2(ctx context.Context, id string) (design.LogicalDevice, error) {
	var response design.LogicalDevice
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(design.LogicalDeviceUrlByID, id),
		apiResponse: &response,
	})
	if err != nil {
		return response, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (c Client) UpdateLogicalDevice2(ctx context.Context, v design.LogicalDevice) error {
	if v.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(design.LogicalDeviceUrlByID, *v.ID()),
		apiInput: &v,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) DeleteLogicalDevice2(ctx context.Context, id string) error {
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(design.LogicalDeviceUrlByID, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) ListLogicalDevices2(ctx context.Context) ([]string, error) {
	var response struct {
		Items []string `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      design.LogicalDevicesUrl,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (c Client) GetLogicalDevices2(ctx context.Context) ([]design.LogicalDevice, error) {
	var response struct {
		Items []design.LogicalDevice `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      design.LogicalDevicesUrl,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (c Client) GetLogicalDeviceByLabel2(ctx context.Context, label string) (design.LogicalDevice, error) {
	var result design.LogicalDevice

	all, err := c.GetLogicalDevices2(ctx)
	if err != nil {
		return result, fmt.Errorf("%s failed getting all %T candidates: %w", str.FuncName(), result, err)
	}

	for _, result = range all {
		if result.Label == label {
			return result, nil
		}
	}

	return result, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("%T with label %s not found", result, label),
	}
}
