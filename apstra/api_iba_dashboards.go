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
	apiUrlIbaDashboards                 = "/api/blueprints/%s/iba/dashboards"
	apiUrlIbaDashboardsPrefix           = apiUrlIbaDashboards + apiUrlPathDelim
	apiUrlIbaDashboardsById             = apiUrlIbaDashboardsPrefix + "%s"
	apiUrlIbaPredefinedDashboards       = "/api/blueprints/%s/iba/predefined-dashboards"
	apiUrlIbaPredefinedDashboardsPrefix = apiUrlIbaPredefinedDashboards + apiUrlPathDelim
	apiUrlIbaPredefinedDashboardsById   = apiUrlIbaPredefinedDashboardsPrefix + "%s"
)

var _ json.Unmarshaler = (*IbaDashboard)(nil)

type IbaDashboard struct {
	Id   ObjectId
	Data *IbaDashboardData
}

func (i *IbaDashboard) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id                  string            `json:"id"`
		Label               string            `json:"label"`
		Description         string            `json:"description"`
		Default             bool              `json:"default"`
		CreatedAt           *string           `json:"created_at"`
		UpdatedAt           *string           `json:"updated_at"`
		IbaWidgetGrid       [][]IbaWidgetData `json:"grid"`
		PredefinedDashboard string            `json:"predefined_dashboard"`
		UpdatedBy           string            `json:"updated_by"`
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

var _ json.Marshaler = (*IbaDashboardData)(nil)

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
	})
}

func (o *Client) listAllIbaPredefinedDashboardIds(ctx context.Context, blueprintId ObjectId) ([]ObjectId, error) {
	var response struct {
		Items []struct {
			Name ObjectId `json:"name"`
		} `json:"items"`
	}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaPredefinedDashboards, blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	ids := make([]ObjectId, len(response.Items))
	for i, r := range response.Items {
		ids[i] = r.Name
	}

	return ids, nil
}

func (o *Client) instantiateIbaPredefinedDashboard(ctx context.Context, blueprintId ObjectId, dashboardId ObjectId, label string) (ObjectId, error) {
	var response objectIdResponse
	var in struct {
		Label string `json:"label"`
	}
	in.Label = label
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlIbaPredefinedDashboardsById, blueprintId, dashboardId),
		apiInput:    &in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *Client) getAllIbaDashboards(ctx context.Context, BlueprintId ObjectId) ([]IbaDashboard, error) {
	var response struct {
		Items []IbaDashboard `json:"items"`
	}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaDashboards, BlueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getIbaDashboard(ctx context.Context, blueprintId ObjectId, id ObjectId) (*IbaDashboard, error) {
	var response IbaDashboard
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaDashboardsById, blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return &response, nil
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
	if in.UpdatedBy != "" {
		return "", errors.New("attempt to create dashboard with non-empty updated_by value - this value can be set only by the server")
	}
	if in.PredefinedDashboard != "" {
		return "", errors.New("attempt to create dashboard with non-empty predefined_dashboard value - this value can " +
			"be set only by the server, and only when a dashboard is instantiated from a predefined template")
	}
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
	if in.UpdatedBy != "" {
		return errors.New("attempt to update dashboard with non-empty updated_by value - this value can be set only by the server")
	}
	if in.PredefinedDashboard != "" {
		return errors.New("attempt to update dashboard with non-empty predefined_dashboard value - this value can " +
			"be set only by the server, and only when a dashboard is instantiated from a predefined template")
	}

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
