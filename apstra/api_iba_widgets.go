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
	IbaWidgetTypeStage          = IbaWidgetType{Value: "stage"}
	IbaWidgetTypeAnomalyHeatmap = IbaWidgetType{Value: "anomaly_heatmap"}
	IbaWidgetTypes              = enum.New(IbaWidgetTypeStage, IbaWidgetTypeAnomalyHeatmap)
)

type IbaWidget struct {
	Id        ObjectId
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      *IbaWidgetData
}

type IbaWidgetData struct {
	AggregationPeriod  *time.Duration
	OrderBy            string
	StageName          string
	ShowContext        bool
	Description        string
	AnomalousOnly      bool
	SpotlightMode      bool
	ProbeId            ObjectId
	Label              string
	Filter             string
	TimeSeriesDuration *time.Duration
	DataSource         string
	MaxItems           *int
	CombineGraphs      string
	VisibleColumns     []string
	Type               IbaWidgetType
	UpdatedBy          string
}

type rawIbaWidget struct {
	AggregationPeriod  *int     `json:"aggregation_period,omitempty"`
	OrderBy            string   `json:"orderby,omitempty"`
	StageName          string   `json:"stage_name,omitempty"`
	ShowContext        bool     `json:"show_context,omitempty"`
	Description        string   `json:"description,omitempty"`
	AnomalousOnly      bool     `json:"anomalous_only,omitempty"`
	CreatedAt          *string  `json:"created_at,omitempty"`
	SpotlightMode      bool     `json:"spotlight_mode,omitempty"`
	UpdatedAt          *string  `json:"updated_at,omitempty"`
	ProbeId            string   `json:"probe_id"`
	Label              string   `json:"label"`
	Filter             string   `json:"filter,omitempty"`
	TimeSeriesDuration *int     `json:"time_series_duration,omitempty"`
	DataSource         string   `json:"data_source,omitempty"`
	MaxItems           *int     `json:"max_items,omitempty"`
	CombineGraphs      string   `json:"combine_graphs,omitempty"`
	VisibleColumns     []string `json:"visible_columns,omitempty"`
	Type               string   `json:"type"`
	Id                 ObjectId `json:"id,omitempty"`
	UpdatedBy          string   `json:"updated_by,omitempty"`
}

func (o *rawIbaWidget) polish() (*IbaWidget, error) {
	var created, updated time.Time

	if o.CreatedAt != nil {
		t, err := time.Parse("2006-01-02T15:04:05.000000+0000", *o.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failure parsing create time %s - %w", *o.CreatedAt, err)
		}
		created = t
	}

	if o.UpdatedAt != nil {
		t, err := time.Parse("2006-01-02T15:04:05.000000+0000", *o.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failure parsing create time %s - %w", *o.CreatedAt, err)
		}
		created = t
	}

	widgetType := IbaWidgetTypes.Parse(o.Type)
	if widgetType == nil {
		return nil, fmt.Errorf("failure to parse returned Iba Widget type %s", o.Type)
	}

	var aggregationPeriod *time.Duration
	if o.AggregationPeriod != nil {
		td := time.Duration(float64(*o.AggregationPeriod) * float64(time.Second))
		aggregationPeriod = &td
	}

	var timeSeriesDuration *time.Duration
	if o.TimeSeriesDuration != nil {
		td := time.Duration(float64(*o.AggregationPeriod) * float64(time.Second))
		timeSeriesDuration = &td
	}

	return &IbaWidget{
		Id:        o.Id,
		CreatedAt: created,
		UpdatedAt: updated,
		Data: &IbaWidgetData{
			AggregationPeriod:  aggregationPeriod,
			OrderBy:            o.OrderBy,
			StageName:          o.StageName,
			ShowContext:        o.ShowContext,
			Description:        o.Description,
			AnomalousOnly:      o.AnomalousOnly,
			SpotlightMode:      o.SpotlightMode,
			ProbeId:            ObjectId(o.ProbeId),
			Label:              o.Label,
			Filter:             o.Filter,
			TimeSeriesDuration: timeSeriesDuration,
			DataSource:         o.DataSource,
			MaxItems:           o.MaxItems,
			CombineGraphs:      o.CombineGraphs,
			VisibleColumns:     o.VisibleColumns,
			Type:               *widgetType,
			UpdatedBy:          o.UpdatedBy,
		},
	}, nil
}

func (o *IbaWidgetData) raw() *rawIbaWidget {
	var aggPeriod, timeSeriesDuration int
	if o.AggregationPeriod != nil {
		aggPeriod = int(o.AggregationPeriod.Seconds())
	} else {
		aggPeriod = 1
	}
	if o.TimeSeriesDuration != nil {
		timeSeriesDuration = int(o.TimeSeriesDuration.Seconds())
	} else {
		timeSeriesDuration = 1
	}
	return &rawIbaWidget{
		AggregationPeriod:  &aggPeriod,
		OrderBy:            o.OrderBy,
		StageName:          o.StageName,
		ShowContext:        o.ShowContext,
		Description:        o.Description,
		AnomalousOnly:      o.AnomalousOnly,
		CreatedAt:          nil,
		SpotlightMode:      o.SpotlightMode,
		UpdatedAt:          nil,
		ProbeId:            o.ProbeId.String(),
		Label:              o.Label,
		Filter:             o.Filter,
		TimeSeriesDuration: &timeSeriesDuration,
		DataSource:         o.DataSource,
		MaxItems:           o.MaxItems,
		CombineGraphs:      o.CombineGraphs,
		VisibleColumns:     o.VisibleColumns,
		Type:               o.Type.Value,
	}
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

func (o *Client) getIbaWidget(ctx context.Context, bpId ObjectId, id ObjectId) (*rawIbaWidget, error) {
	response := &rawIbaWidget{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaWidgetsById, bpId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getIbaWidgetsByLabel(ctx context.Context, bpId ObjectId, label string) ([]rawIbaWidget, error) {
	allIbaWidgets, err := o.getAllIbaWidgets(ctx, bpId)
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

func (o *Client) getIbaWidgetByLabel(ctx context.Context, bpId ObjectId, label string) (*rawIbaWidget, error) {
	rawWidgets, err := o.getIbaWidgetsByLabel(ctx, bpId, label)
	if err != nil {
		return nil, err
	}

	switch len(rawWidgets) {
	case 0:
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("IBA widget with label %q not found in blueprint %q", label, bpId),
		}
	case 1:
		return &rawWidgets[0], nil
	}

	return nil, ClientErr{
		errType: ErrMultipleMatch,
		err:     fmt.Errorf("multiple IBA widget with label %q found in blueprint %q", label, bpId),
	}
}

func (o *Client) createIbaWidget(ctx context.Context, bpId ObjectId, widget *rawIbaWidget) (ObjectId, error) {
	response := &objectIdResponse{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlIbaWidgets, bpId),
		apiInput:    &widget,
		apiResponse: &response,
	})
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

func (o *Client) deleteIbaWidget(ctx context.Context, bpId ObjectId, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlIbaWidgetsById, bpId, id),
	})
}
