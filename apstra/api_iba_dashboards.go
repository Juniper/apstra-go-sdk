// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	apiUrlIbaDashboards       = "/api/blueprints/%s/iba/dashboards"
	apiUrlIbaDashboardsPrefix = apiUrlIbaDashboards + apiUrlPathDelim
	apiUrlIbaDashboardsById   = apiUrlIbaDashboardsPrefix + "%s"
)

var _ json.Marshaler = new(IbaDashboard)
var _ json.Unmarshaler = new(IbaDashboard)

var _ json.Marshaler = new(IbaDashboardData)

type IbaDashboard struct {
	Id   ObjectId
	Data *IbaDashboardData
}

func (i *IbaDashboard) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id                  string            `json:"id,omitempty"`
		Label               string            `json:"label"`
		Description         string            `json:"description"`
		Default             bool              `json:"default,omitempty"`
		CreatedAt           *string           `json:"created_at,omitempty"`
		UpdatedAt           *string           `json:"updated_at,omitempty"`
		IbaWidgetGrid       [][]IbaWidgetData `json:"grid"`
		PredefinedDashboard string            `json:"predefined_dashboard,omitempty"`
		UpdatedBy           string            `json:"updated_by,omitempty"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	i.Id = ObjectId(raw.Id)
	i.Data = &IbaDashboardData{
		Description:         raw.Description,
		Default:             raw.Default,
		Label:               raw.Label,
		IbaWidgetGrid:       raw.IbaWidgetGrid,
		PredefinedDashboard: raw.PredefinedDashboard,
		UpdatedBy:           raw.UpdatedBy,
	}

	return nil
}

func (i *IbaDashboard) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id                  string            `json:"id,omitempty"`
		Label               string            `json:"label"`
		Description         string            `json:"description"`
		Default             bool              `json:"default,omitempty"`
		IbaWidgetGrid       [][]IbaWidgetData `json:"grid"`
		PredefinedDashboard string            `json:"predefined_dashboard,omitempty"`
		UpdatedBy           string            `json:"updated_by,omitempty"`
	}{
		Id:                  i.Id.String(),
		Label:               i.Data.Label,
		Description:         i.Data.Description,
		Default:             i.Data.Default,
		IbaWidgetGrid:       i.Data.IbaWidgetGrid,
		PredefinedDashboard: i.Data.PredefinedDashboard,
		UpdatedBy:           i.Data.UpdatedBy,
	},
	)
}

type IbaDashboardData struct {
	Description         string
	Default             bool
	Label               string
	IbaWidgetGrid       [][]IbaWidgetData
	PredefinedDashboard string
	UpdatedBy           string
}

func (i *IbaDashboardData) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id                  string            `json:"id,omitempty"`
		Label               string            `json:"label"`
		Description         string            `json:"description"`
		Default             bool              `json:"default,omitempty"`
		IbaWidgetGrid       [][]IbaWidgetData `json:"grid"`
		PredefinedDashboard string            `json:"predefined_dashboard,omitempty"`
		UpdatedBy           string            `json:"updated_by,omitempty"`
	}{
		Label:               i.Label,
		Description:         i.Description,
		Default:             i.Default,
		IbaWidgetGrid:       i.IbaWidgetGrid,
		PredefinedDashboard: i.PredefinedDashboard,
		UpdatedBy:           i.UpdatedBy,
	},
	)
}

func (o *Client) getAllIbaDashboards(ctx context.Context, BlueprintId ObjectId) ([]IbaDashboard, error) {
	response := &struct {
		Items []IbaDashboard `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodGet, urlStr: fmt.Sprintf(apiUrlIbaDashboards, BlueprintId.String()),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getIbaDashboard(ctx context.Context, blueprintId ObjectId, id ObjectId) (*IbaDashboard, error) {
	response := &IbaDashboard{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodGet, urlStr: fmt.Sprintf(apiUrlIbaDashboardsById, blueprintId, id), apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getIbaDashboardByLabel(ctx context.Context, blueprintId ObjectId, label string) (*IbaDashboard, error) {
	dashes, err := o.getAllIbaDashboards(ctx, blueprintId)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	var result []IbaDashboard
	for _, w := range dashes {
		if w.Data.Label == label {
			result = append(result, w)
		}
	}
	l := len(result)

	if l == 0 {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no Iba Dashboards with label '%s' found", label),
		}
	}

	if l > 1 {
		return nil, ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("%d Iba Dashboards with label '%s' found, expected 1", l, label),
		}
	}

	return &result[0], nil
}

func (o *Client) createIbaDashboard(ctx context.Context, blueprintId ObjectId, in *IbaDashboardData) (ObjectId, error) {
	var response objectIdResponse
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlIbaDashboards, blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err == nil {
		return response.Id, nil
	}

	err = convertTtaeToAceWherePossible(err)

	var ace ClientErr
	if !(errors.As(err, &ace) && ace.IsRetryable()) {
		return "", err // fatal error
	}

	retryMax := o.GetTuningParam("ibaDashboardMaxRetries")
	retryInterval := time.Duration(o.GetTuningParam("ibaDashboardRetryIntervalMs")) * time.Millisecond

	for i := 0; i < retryMax; i++ {
		// Make a random wait, in case multiple threads are running
		if rand.Int()%2 == 0 {
			time.Sleep(retryInterval)
		}

		time.Sleep(retryInterval * time.Duration(i))

		e := o.talkToApstra(ctx, &talkToApstraIn{
			method:      http.MethodPost,
			urlStr:      fmt.Sprintf(apiUrlIbaDashboards, blueprintId),
			apiInput:    in,
			apiResponse: &response,
		})
		if e == nil {
			return response.Id, nil // success!
		}

		e = convertTtaeToAceWherePossible(e)
		if !(errors.As(e, &ace) && ace.IsRetryable()) {
			return "", e // return the fatal error
		}

		err = errors.Join(err, e) // the error is retryable; stack it with the rest
	}

	return "", errors.Join(err, fmt.Errorf("reached retry limit %d", retryMax))
}

func (o *Client) updateIbaDashboard(ctx context.Context, blueprintId ObjectId, id ObjectId, in *IbaDashboardData) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPut, urlStr: fmt.Sprintf(apiUrlIbaDashboardsById, blueprintId, id), apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *Client) deleteIbaDashboard(ctx context.Context, blueprintId ObjectId, id ObjectId) error {
	return convertTtaeToAceWherePossible(o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete, urlStr: fmt.Sprintf(apiUrlIbaDashboardsById, blueprintId, id),
	}))
}
