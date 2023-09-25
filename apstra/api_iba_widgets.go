package apstra

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlIBAWidgets       = "/api/blueprints/%s/iba/widgets"
	apiUrlIBAWidgetsPrefix = apiUrlIBAWidgets + apiUrlPathDelim
	apiUrlIBAWidgetsById   = apiUrlIBAWidgetsPrefix + "%s"
)

type IBAWidget struct {
	Id        ObjectId
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      *IBAWidgetData
}

type IBAWidgetData struct {
	AggregationPeriod  int
	Orderby            string
	StageName          string
	ShowContext        bool
	Description        string
	AnomalousOnly      bool
	SpotlightMode      bool
	ProbeId            string
	Label              string
	Filter             string
	TimeSeriesDuration int
	DataSource         string
	MaxItems           int
	CombineGraphs      string
	VisibleColumns     []string
	Type               string
	UpdatedBy          string
}

type rawIBAWidget struct {
	AggregationPeriod  int      `json:"aggregation_period"`
	Orderby            string   `json:"orderby"`
	StageName          string   `json:"stage_name"`
	ShowContext        bool     `json:"show_context"`
	Description        string   `json:"description"`
	AnomalousOnly      bool     `json:"anomalous_only"`
	CreatedAt          string   `json:"created_at"`
	SpotlightMode      bool     `json:"spotlight_mode"`
	UpdatedAt          string   `json:"updated_at"`
	ProbeId            string   `json:"probe_id"`
	Label              string   `json:"label"`
	Filter             string   `json:"filter"`
	TimeSeriesDuration int      `json:"time_series_duration"`
	DataSource         string   `json:"data_source"`
	MaxItems           int      `json:"max_items"`
	CombineGraphs      string   `json:"combine_graphs"`
	VisibleColumns     []string `json:"visible_columns"`
	Type               string   `json:"type"`
	Id                 string   `json:"id"`
	UpdatedBy          string   `json:"updated_by"`
}

func (o *rawIBAWidget) polish() (*IBAWidget, error) {
	created, err := time.Parse("2006-01-02T15:04:05.000000+0000", o.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing create time %s - %w", o.CreatedAt, err)
	}
	updated, err := time.Parse("2006-01-02T15:04:05.000000+0000", o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing update time %s - %w", o.UpdatedAt, err)
	}
	return &IBAWidget{
		Id:        ObjectId(o.Id),
		CreatedAt: created,
		UpdatedAt: updated,
		Data: &IBAWidgetData{
			AggregationPeriod:  o.AggregationPeriod,
			Orderby:            o.Orderby,
			StageName:          o.StageName,
			ShowContext:        o.ShowContext,
			Description:        o.Description,
			AnomalousOnly:      o.AnomalousOnly,
			SpotlightMode:      o.SpotlightMode,
			ProbeId:            o.ProbeId,
			Label:              o.Label,
			Filter:             o.Filter,
			TimeSeriesDuration: o.TimeSeriesDuration,
			DataSource:         o.DataSource,
			MaxItems:           o.MaxItems,
			CombineGraphs:      o.CombineGraphs,
			VisibleColumns:     o.VisibleColumns,
			Type:               o.Type,
			UpdatedBy:          o.UpdatedBy,
		},
	}, nil
}

func (o *Client) getAllIBAWidgets(ctx context.Context, bp_id ObjectId) ([]rawIBAWidget, error) {
	response := &struct {
		Items []rawIBAWidget `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIBAWidgets, bp_id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getIBAWidget(ctx context.Context, bp_id ObjectId, id ObjectId) (*rawIBAWidget, error) {
	response := &rawIBAWidget{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIBAWidgetsById, bp_id, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getIBAWidgetsByLabel(ctx context.Context, bp_id ObjectId, label string) ([]rawIBAWidget, error) {
	allIBAWidgets, err := o.getAllIBAWidgets(ctx, bp_id)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	var result []rawIBAWidget
	for _, w := range allIBAWidgets {
		if w.Label == label {
			result = append(result, w)
		}
	}

	if len(result) == 0 {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("property set with label '%s' not found", label),
		}
	}
	return result, nil
}
