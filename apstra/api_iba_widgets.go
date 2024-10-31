// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"encoding/json"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"time"
)

type IbaWidget struct {
	Id        ObjectId
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      *IbaWidgetData
}

type IbaWidgetData struct {
	AggregationPeriod  *time.Duration
	WidgetType         enum.IbaWidgetType
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
	UpdatedBy          string
}

var _ json.Marshaler = new(IbaWidgetData)
var _ json.Unmarshaler = new(IbaWidgetData)

func (i *IbaWidgetData) UnmarshalJSON(bytes []byte) error {
	//TODO implement me
	var raw struct {
		AggregationPeriod  *int     `json:"aggregation_period,omitempty"`
		OrderBy            string   `json:"orderby,omitempty"`
		StageName          string   `json:"stage_name,omitempty"`
		ShowContext        bool     `json:"show_context,omitempty"`
		Description        string   `json:"description,omitempty"`
		AnomalousOnly      bool     `json:"anomalous_only,omitempty"`
		SpotlightMode      bool     `json:"spotlight_mode,omitempty"`
		ProbeId            string   `json:"probe_id"`
		Label              string   `json:"label"`
		Filter             string   `json:"filter,omitempty"`
		TimeSeriesDuration *int     `json:"time_series_duration,omitempty"`
		DataSource         string   `json:"data_source,omitempty"`
		MaxItems           *int     `json:"max_items,omitempty"`
		CombineGraphs      string   `json:"combine_graphs,omitempty"`
		VisibleColumns     []string `json:"visible_columns,omitempty"`
		Id                 ObjectId `json:"id,omitempty"`
		UpdatedBy          string   `json:"updated_by,omitempty"`
		Type               string   `json:"type,omitempty"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	var t enum.IbaWidgetType
	err = t.FromString(raw.Type)
	if err != nil {
		return err
	}

	var aggregationPeriod *time.Duration
	if raw.AggregationPeriod != nil {
		td := time.Duration(float64(*raw.AggregationPeriod) * float64(time.Second))
		aggregationPeriod = &td
	}

	var timeSeriesDuration *time.Duration
	if raw.TimeSeriesDuration != nil {
		td := time.Duration(float64(*raw.TimeSeriesDuration) * float64(time.Second))
		timeSeriesDuration = &td
	}

	*i = IbaWidgetData{
		AggregationPeriod:  aggregationPeriod,
		OrderBy:            raw.OrderBy,
		StageName:          raw.StageName,
		ShowContext:        raw.ShowContext,
		Description:        raw.Description,
		AnomalousOnly:      raw.AnomalousOnly,
		SpotlightMode:      raw.SpotlightMode,
		ProbeId:            ObjectId(raw.ProbeId),
		Label:              raw.Label,
		Filter:             raw.Filter,
		TimeSeriesDuration: timeSeriesDuration,
		DataSource:         raw.DataSource,
		MaxItems:           raw.MaxItems,
		CombineGraphs:      raw.CombineGraphs,
		VisibleColumns:     raw.VisibleColumns,
		UpdatedBy:          raw.UpdatedBy,
		WidgetType:         t,
	}

	return err
}

func (i *IbaWidgetData) MarshalJSON() ([]byte, error) {
	var aggPeriod, timeSeriesDuration int
	if i.AggregationPeriod != nil {
		aggPeriod = int(i.AggregationPeriod.Seconds())
	} else {
		aggPeriod = 1
	}
	if i.TimeSeriesDuration != nil {
		timeSeriesDuration = int(i.TimeSeriesDuration.Seconds())
	} else {
		timeSeriesDuration = 1
	}

	var raw = struct {
		AggregationPeriod  *int     `json:"aggregation_period,omitempty"`
		OrderBy            string   `json:"orderby,omitempty"`
		StageName          string   `json:"stage_name,omitempty"`
		ShowContext        bool     `json:"show_context,omitempty"`
		Description        string   `json:"description,omitempty"`
		AnomalousOnly      bool     `json:"anomalous_only,omitempty"`
		SpotlightMode      bool     `json:"spotlight_mode,omitempty"`
		ProbeId            string   `json:"probe_id"`
		Label              string   `json:"label"`
		Filter             string   `json:"filter,omitempty"`
		TimeSeriesDuration *int     `json:"time_series_duration,omitempty"`
		DataSource         string   `json:"data_source,omitempty"`
		MaxItems           *int     `json:"max_items,omitempty"`
		CombineGraphs      string   `json:"combine_graphs,omitempty"`
		VisibleColumns     []string `json:"visible_columns,omitempty"`
		Id                 ObjectId `json:"id,omitempty"`
		UpdatedBy          string   `json:"updated_by,omitempty"`
		Type               string   `json:"type,omitempty"`
	}{
		AggregationPeriod:  &aggPeriod,
		OrderBy:            i.OrderBy,
		StageName:          i.StageName,
		ShowContext:        i.ShowContext,
		Description:        i.Description,
		AnomalousOnly:      i.AnomalousOnly,
		SpotlightMode:      i.SpotlightMode,
		ProbeId:            i.ProbeId.String(),
		Label:              i.Label,
		Filter:             i.Filter,
		TimeSeriesDuration: &timeSeriesDuration,
		DataSource:         i.DataSource,
		MaxItems:           i.MaxItems,
		CombineGraphs:      i.CombineGraphs,
		VisibleColumns:     i.VisibleColumns,
		Type:               i.WidgetType.String(),
	}
	return json.Marshal(raw)
}
