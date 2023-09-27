package apstra

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlIbaDashboards       = "/api/blueprints/%s/iba/dashboards"
	apiUrlIbaDashboardsPrefix = apiUrlIbaDashboards + apiUrlPathDelim
	apiUrlIbaDashboardsById   = apiUrlIbaDashboardsPrefix + "%s"
)

type rawIbaDashboard struct {
	Id                  string     `json:"id"`
	Label               string     `json:"label"`
	Description         string     `json:"description"`
	Default             bool       `json:"default,omitempty"`
	CreatedAt           string     `json:"created_at"`
	UpdatedAt           string     `json:"updated_at"`
	IbaWidgetGrid       [][]string `json:"grid"`
	PredefinedDashboard string     `json:"predefined_dashboard,omitempty"`
	UpdatedBy           string     `json:"updated_by,omitempty"`
}

type IbaDashboard struct {
	Id             ObjectId
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Data           *IbaDashboardData
}

type IbaDashboardData struct {
	Description         string
	Default             bool
	Label               string
	IbaWidgetGrid       [][]ObjectId
	PredefinedDashboard string
	UpdatedBy           string
}

type rawIbaDashboardData struct {
	Label               string     `json:"label"`
	Description         string     `json:"description"`
	Default             bool       `json:"default,omitempty"`
	IbaWidgetGrid       [][]string `json:"grid"`
	PredefinedDashboard string     `json:"predefined_dashboard,omitempty"`
	UpdatedBy           string     `json:"updated_by,omitempty"`
}

func (o *IbaDashboardData) raw() *rawIbaDashboardData {
	IbaWidgetGrid := make([][]string, len(o.IbaWidgetGrid))

	for i, g := range o.IbaWidgetGrid {
		for _, v := range g {
			IbaWidgetGrid[i] = append(IbaWidgetGrid[i], v.String())
		}
	}

	return &rawIbaDashboardData{
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
		Id:             ObjectId(o.Id),
		CreatedAt:      created,
		LastModifiedAt: updated,
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
	l := len(dashes)

	if l == 0 {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no Iba Dashboards with label '%s' found", label),
		}
	}

	if l > 1 {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("%d Iba Dashboards with label '%s' found, expected 1", l, label),
		}
	}

	return &dashes[0], nil
}

func (o *Client) createIbaDashboard(ctx context.Context, blueprintId ObjectId, in *rawIbaDashboardData) (ObjectId, error) {
	response := &objectIdResponse{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPost, urlStr: fmt.Sprintf(apiUrlIbaDashboards, blueprintId),
		apiInput: in, apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *Client) updateIbaDashboard(ctx context.Context, blueprintId ObjectId, id ObjectId, in *rawIbaDashboardData) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPut, urlStr: fmt.Sprintf(apiUrlIbaDashboardsById, blueprintId, id), apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *Client) deleteIbaDashboard(ctx context.Context, blueprintId ObjectId, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete, urlStr: fmt.Sprintf(apiUrlIbaDashboardsById, blueprintId, id),
	})
}
