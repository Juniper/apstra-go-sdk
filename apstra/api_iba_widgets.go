package apstra

import (
	"context"
	"fmt"
	"github.com/orsinium-labs/enum"
	"net/http"
	"time"
)

const (
	apiUrlIbaWidgets       = "/api/blueprints/%s/iba/widgets"
	apiUrlIbaWidgetsPrefix = apiUrlIbaWidgets + apiUrlPathDelim
	apiUrlIbaWidgetsById   = apiUrlIbaWidgetsPrefix + "%s"
)

type IbaWidgetType enum.Member[string]

var (
	IbaWidgetTypeStage          = IbaWidgetType{"stage"}
	IbaWidgetTypeAnomalyHeatmap = IbaWidgetType{"anomaly_heatmap"}
	IbaWidgetTypes              = enum.New(IbaWidgetTypeStage, IbaWidgetTypeAnomalyHeatmap)
)

type IbaWidget struct {
	Id        ObjectId
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      *IbaWidgetData
}

type IbaWidgetData struct {
	AggregationPeriod  time.Duration
	Orderby            string
	StageName          string
	ShowContext        bool
	Description        string
	AnomalousOnly      bool
	SpotlightMode      bool
	ProbeId            ObjectId
	Label              string
	Filter             string
	TimeSeriesDuration time.Duration
	DataSource         string
	MaxItems           int
	CombineGraphs      string
	VisibleColumns     []string
	Type               IbaWidgetType
	UpdatedBy          string
}

type rawIbaWidget struct {
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

func (o *rawIbaWidget) polish() (*IbaWidget, error) {
	created, err := time.Parse("2006-01-02T15:04:05.000000+0000", o.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failure parsing create time %s - %w", o.CreatedAt, err)
	}
	updated, err := time.Parse("2006-01-02T15:04:05.000000+0000", o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failure parsing update time %s - %w", o.UpdatedAt, err)
	}
	return &IbaWidget{
		Id:        ObjectId(o.Id),
		CreatedAt: created,
		UpdatedAt: updated,
		Data: &IbaWidgetData{
			AggregationPeriod:  time.Duration(float64(o.AggregationPeriod) * float64(time.Second)),
			Orderby:            o.Orderby,
			StageName:          o.StageName,
			ShowContext:        o.ShowContext,
			Description:        o.Description,
			AnomalousOnly:      o.AnomalousOnly,
			SpotlightMode:      o.SpotlightMode,
			ProbeId:            ObjectId(o.ProbeId),
			Label:              o.Label,
			Filter:             o.Filter,
			TimeSeriesDuration: time.Duration(float64(o.TimeSeriesDuration) * float64(time.Second)),
			DataSource:         o.DataSource,
			MaxItems:           o.MaxItems,
			CombineGraphs:      o.CombineGraphs,
			VisibleColumns:     o.VisibleColumns,
			Type:               *IbaWidgetTypes.Parse(o.Type),
			UpdatedBy:          o.UpdatedBy,
		},
	}, nil
}

func (o *Client) getAllIbaWidgets(ctx context.Context, bp_id ObjectId) ([]rawIbaWidget, error) {
	response := &struct {
		Items []rawIbaWidget `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaWidgets, bp_id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getIbaWidget(ctx context.Context, bp_id ObjectId, id ObjectId) (*rawIbaWidget, error) {
	response := &rawIbaWidget{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaWidgetsById, bp_id, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getIbaWidgetsByLabel(ctx context.Context, bp_id ObjectId, label string) ([]rawIbaWidget, error) {
	allIbaWidgets, err := o.getAllIbaWidgets(ctx, bp_id)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	var result []rawIbaWidget
	for _, w := range allIbaWidgets {
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
