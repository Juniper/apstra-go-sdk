// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/internal/str"
	"github.com/Juniper/apstra-go-sdk/internal/urls"
)

func (c *TwoStageL3ClosClient) CreateSwitchingZone(ctx context.Context, v datacenter.SwitchingZone) (string, error) {
	var response struct {
		ID string `json:"id"`
	}
	err := c.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(urls.DatacenterSwitchingZones, c.blueprintId),
		apiInput:    v,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.ID, nil
}

func (c *TwoStageL3ClosClient) GetSwitchingZone(ctx context.Context, id string) (datacenter.SwitchingZone, error) {
	var result datacenter.SwitchingZone
	err := c.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(urls.DatacenterSwitchingZoneByID, c.blueprintId, id),
		apiResponse: &result,
	})
	if err != nil {
		return result, convertTtaeToAceWherePossible(err)
	}

	return result, nil
}

func (c *TwoStageL3ClosClient) GetSwitchingZones(ctx context.Context) ([]datacenter.SwitchingZone, error) {
	var result struct {
		Items map[string]datacenter.SwitchingZone `json:"items"`
	}
	err := c.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(urls.DatacenterSwitchingZones, c.blueprintId),
		apiResponse: &result,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return slices.Collect(maps.Values(result.Items)), nil
}

func (c *TwoStageL3ClosClient) GetSwitchingZoneByLabel(ctx context.Context, label string) (datacenter.SwitchingZone, error) {
	items, err := c.GetSwitchingZones(ctx)
	if err != nil {
		return datacenter.SwitchingZone{}, err
	}

	var result *datacenter.SwitchingZone
	for _, item := range items {
		if item.Label != nil && *item.Label == label {
			if result == nil {
				result = &item
			} else {
				return datacenter.SwitchingZone{}, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("found multiple Switching Zones with label %q", label),
				}
			}
		}
	}

	if result == nil {
		return datacenter.SwitchingZone{}, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no Switching Zone with label %q", label),
		}
	}

	return *result, nil
}

func (c *TwoStageL3ClosClient) GetDefaultSwitchingZone(ctx context.Context) (datacenter.SwitchingZone, error) {
	// collect all Switching Zones here without unpacking them
	var response struct {
		Items map[string]json.RawMessage `json:"items"`
	}

	err := c.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(urls.DatacenterSwitchingZones, c.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return datacenter.SwitchingZone{}, convertTtaeToAceWherePossible(err)
	}

	// impl_type is omitted from the SwitchingZone type. This struct exists as a filter used here only.
	var target struct {
		ImplType string `json:"impl_type"`
	}

	// loop over switching zones looking for the one (we expect) with impl_type == default
	var result *datacenter.SwitchingZone
	for i, item := range response.Items {
		if err = json.Unmarshal(item, &target); err != nil {
			return datacenter.SwitchingZone{}, fmt.Errorf("error unmarshalling Switching Zone %q: %w", i, err)
		}
		if target.ImplType == "default" {
			if result != nil {
				return datacenter.SwitchingZone{}, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("found multiple Switching Zones with impl_type == default"),
				}
			}
			if err = json.Unmarshal(item, &result); err != nil {
				return datacenter.SwitchingZone{}, fmt.Errorf("error unmarshalling default Switching Zone: %w", err)
			}
		}
	}

	if result == nil {
		return datacenter.SwitchingZone{}, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no Switching Zone with impl_type == default"),
		}
	}

	return *result, nil
}

func (c *TwoStageL3ClosClient) ListSwitchingZones(ctx context.Context) ([]string, error) {
	items, err := c.GetSwitchingZones(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(items))
	for i, item := range items {
		idPtr := item.ID()
		if idPtr == nil {
			return nil, ClientErr{
				errType: ErrInvalidId,
				err:     fmt.Errorf("the Switching Zone at index %d has nil ID", i),
			}
		}
		result[i] = *idPtr
	}

	return result, nil
}

func (c *TwoStageL3ClosClient) UpdateSwitchingZone(ctx context.Context, v datacenter.SwitchingZone) error {
	if v.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	err := c.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(urls.DatacenterSwitchingZoneByID, c.blueprintId, *v.ID()),
		apiInput: v,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c *TwoStageL3ClosClient) DeleteSwitchingZone(ctx context.Context, id string) error {
	err := c.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(urls.DatacenterSwitchingZoneByID, c.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
