// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/device"
	"github.com/Juniper/apstra-go-sdk/internal/str"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
)

func (c Client) CreateDeviceProfile(ctx context.Context, v device.Profile) (string, error) {
	if v.ID() != nil {
		return "", fmt.Errorf("id must be nil in %s", str.FuncName())
	}
	if v.Predefined {
		return "", fmt.Errorf("predefined must be false in %s", str.FuncName())
	}

	var response struct {
		ID string `json:"id"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      device.ProfilesURL,
		apiInput:    v,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.ID, nil
}

func (c Client) GetDeviceProfile(ctx context.Context, id string) (device.Profile, error) {
	var response device.Profile
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(device.ProfileURLByID, id),
		apiResponse: &response,
	})
	if err != nil {
		return response, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (c Client) UpdateDeviceProfile(ctx context.Context, v device.Profile) error {
	if v.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(device.ProfileURLByID, *v.ID()),
		apiInput: &v,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) DeleteDeviceProfile(ctx context.Context, id string) error {
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(device.ProfileURLByID, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) ListDeviceProfiles(ctx context.Context) ([]string, error) {
	var response struct {
		Items []string `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      device.ProfilesURL,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (c Client) GetDeviceProfiles(ctx context.Context) ([]device.Profile, error) {
	var response struct {
		Items []device.Profile `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      device.ProfilesURL,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (c Client) GetDeviceProfileByLabel(ctx context.Context, label string) (device.Profile, error) {
	result, err := c.GetDeviceProfiles(ctx)
	if err != nil {
		return zero.SliceItem(result), fmt.Errorf("%s failed getting all %T candidates: %w", str.FuncName(), zero.SliceItem(result), err)
	}

	for i := len(result) - 1; i >= 0; i-- {
		if result[i].Label != label {
			result = slice.Remove(result, i)
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
	default:
		return zero.SliceItem(result), ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("found multiple candidate %T with label %s", zero.SliceItem(result), label),
		}
	}
}
