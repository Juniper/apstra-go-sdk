package apstra

import (
	"context"
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

type rawIbaDashboard struct {
	Id                  string       `json:"id,omitempty"`
	Label               string       `json:"label"`
	Description         string       `json:"description"`
	Default             bool         `json:"default,omitempty"`
	CreatedAt           string       `json:"created_at,omitempty"`
	UpdatedAt           string       `json:"updated_at,omitempty"`
	IbaWidgetGrid       [][]ObjectId `json:"grid"`
	PredefinedDashboard string       `json:"predefined_dashboard,omitempty"`
	UpdatedBy           string       `json:"updated_by,omitempty"`
}

type IbaDashboard struct {
	Id        ObjectId
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      *IbaDashboardData
}

type IbaDashboardData struct {
	Description         string
	Default             bool
	Label               string
	IbaWidgetGrid       [][]ObjectId
	PredefinedDashboard string
	UpdatedBy           string
}

func (o *IbaDashboardData) raw() *rawIbaDashboard {
	IbaWidgetGrid := make([][]ObjectId, len(o.IbaWidgetGrid))

	for i, g := range o.IbaWidgetGrid {
		for _, v := range g {
			IbaWidgetGrid[i] = append(IbaWidgetGrid[i], v)
		}
	}

	return &rawIbaDashboard{
		Label:               o.Label,
		Description:         o.Description,
		Default:             o.Default,
		IbaWidgetGrid:       IbaWidgetGrid,
		PredefinedDashboard: o.PredefinedDashboard,
		UpdatedBy:           o.UpdatedBy,
	}
}

func (o *rawIbaDashboard) polish() (*IbaDashboard, error) {
	var err error

	IbaWidgetGrid := make([][]ObjectId, len(o.IbaWidgetGrid))
	for i, g := range o.IbaWidgetGrid {
		for _, v := range g {
			IbaWidgetGrid[i] = append(IbaWidgetGrid[i], ObjectId(v))
		}
	}

	created, err := time.Parse("2006-01-02T15:04:05.000000+0000", o.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failure parsing create time %s - %w", o.CreatedAt, err)
	}
	updated, err := time.Parse("2006-01-02T15:04:05.000000+0000", o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failure parsing update time %s - %w", o.UpdatedAt, err)
	}
	return &IbaDashboard{
		Id:        ObjectId(o.Id),
		CreatedAt: created,
		UpdatedAt: updated,
		Data: &IbaDashboardData{
			Description:         o.Description,
			Default:             o.Default,
			Label:               o.Label,
			IbaWidgetGrid:       IbaWidgetGrid,
			PredefinedDashboard: o.PredefinedDashboard,
			UpdatedBy:           o.UpdatedBy,
		},
	}, nil
}

func (o *Client) getAllIbaDashboards(ctx context.Context, blueprint_id ObjectId) ([]rawIbaDashboard, error) {
	response := &struct {
		Items []rawIbaDashboard `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodGet, urlStr: fmt.Sprintf(apiUrlIbaDashboards, blueprint_id.String()),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getIbaDashboard(ctx context.Context, blueprintId ObjectId, id ObjectId) (*rawIbaDashboard, error) {
	response := &rawIbaDashboard{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodGet, urlStr: fmt.Sprintf(apiUrlIbaDashboardsById, blueprintId, id), apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getIbaDashboardByLabel(ctx context.Context, blueprintId ObjectId, label string) (*rawIbaDashboard, error) {
	dashes, err := o.getAllIbaDashboards(ctx, blueprintId)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	var result []rawIbaDashboard
	for _, w := range dashes {
		if w.Label == label {
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

func (o *Client) createIbaDashboard(ctx context.Context, blueprintId ObjectId, in *rawIbaDashboard) (ObjectId, error) {
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

func (o *Client) updateIbaDashboard(ctx context.Context, blueprintId ObjectId, id ObjectId, in *rawIbaDashboard) error {
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
