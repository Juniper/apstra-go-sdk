package apstra

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlIbaWidgets       = "/api/blueprints/%s/iba/widgets"
	apiUrlIbaWidgetsPrefix = apiUrlIbaWidgets + apiUrlPathDelim
	apiUrlIbaWidgetsById   = apiUrlIbaWidgetsPrefix + "%s"
)

type IbaWidgetType int
type ibaWidgetType string

const (
	IbaWidgetTypeStage = IbaWidgetType(iota)
	IbaWidgetTypeAnomalyHeatmap

	ibaWidgetTypeStage          = ibaWidgetType("stage")
	ibaWidgetTypeAnomalyHeatmap = ibaWidgetType("anomaly_heatmap")
	ibaWidgetTypeUnknown        = "unknown widget type %s"

	IbaWidgetTypeUnknown = "Unknown Widget Type %d"
)

func (o IbaWidgetType) Int() int {
	return int(o)
}

func (o IbaWidgetType) String() string {
	switch o {
	case IbaWidgetTypeStage:
		return string(ibaWidgetTypeStage)
	case IbaWidgetTypeAnomalyHeatmap:
		return string(ibaWidgetTypeAnomalyHeatmap)
	default:
		return fmt.Sprintf(IbaWidgetTypeUnknown, o)
	}
}

func (o IbaWidgetType) string() string {
	return o.String()
}

func (o *IbaWidgetType) FromString(s string) error {
	i, err := ibaWidgetType(s).parse()
	if err != nil {
		return err
	}
	*o = IbaWidgetType(i)
	return nil
}

func (o IbaWidgetType) raw() ibaWidgetType {
	return ibaWidgetType(o.String())
}

func (o ibaWidgetType) parse() (int, error) {
	switch o {
	case ibaWidgetTypeStage:
		return int(IbaWidgetTypeStage), nil
	case ibaWidgetTypeAnomalyHeatmap:
		return int(IbaWidgetTypeAnomalyHeatmap), nil
	default:
		return 0, fmt.Errorf(ibaWidgetTypeUnknown, o)
	}
}

func (o ibaWidgetType) string() string {
	return string(o)
}

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
	var wtype IbaWidgetType
	err = wtype.FromString(o.Type)
	if err != nil {
		return nil, fmt.Errorf("failure parsing widget type %w", err)
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
			Type:               wtype,
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
