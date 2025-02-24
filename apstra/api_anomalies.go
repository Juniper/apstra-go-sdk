// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	apiUrlBlueprintAnomalies          = apiUrlBlueprintById + apiUrlPathDelim + "anomalies"
	apiUrlBlueprintAnomaliesByNode    = apiUrlBlueprintById + apiUrlPathDelim + "anomalies_nodes_count"
	apiUrlBlueprintAnomaliesByService = apiUrlBlueprintById + apiUrlPathDelim + "anomalies_services_count"
)

type AnomalyProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Anomaly struct {
	Actual      Actual    `json:"actual"`
	Anomalous   Anomalous `json:"anomalous"`
	AnomalyType string    `json:"anomaly_type"`
	Id          ObjectId  `json:"id"`
	Identity    struct {
		AnomalyType string            `json:"anomaly_type"`
		ItemId      ObjectId          `json:"item_id"`
		ProbeId     ObjectId          `json:"probe_id"`
		ProbeLabel  string            `json:"probe_label"`
		Properties  []AnomalyProperty `json:"properties"`
		StageName   string            `json:"stage_name"`
	} `json:"identity"`
	LastModifiedAt time.Time `json:"last_modified_at"`
	Severity       string    `json:"severity"`
}

type Actual struct {
	Value string `json:"value"`
}

type Anomalous struct {
	Value    string `json:"value"`
	ValueMin string `json:"value_min"`
	ValueMax string `json:"value_max"`
}

type rawAnomalyProperty struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

type rawAnomaly struct {
	Actual      rawActual    `json:"actual"`
	Anomalous   rawAnomalous `json:"anomalous"`
	AnomalyType string       `json:"anomaly_type"`
	Id          ObjectId     `json:"id"`
	Identity    struct {
		AnomalyType string               `json:"anomaly_type"`
		ItemId      ObjectId             `json:"item_id"`
		ProbeId     ObjectId             `json:"probe_id"`
		ProbeLabel  string               `json:"probe_label"`
		Properties  []rawAnomalyProperty `json:"properties"`
		StageName   string               `json:"stage_name"`
	} `json:"identity"`
	LastModifiedAt time.Time `json:"last_modified_at"`
	Severity       string    `json:"severity"`
}

type rawActual struct {
	Value json.RawMessage `json:"value"`
}

func (o rawActual) Parse() (*Actual, error) {
	var val string
	var err error

	if o.Value != nil {
		val, err = unpackIntOrStringAsString(o.Value)
		if err != nil {
			return nil, fmt.Errorf("error parsing rawActual payload - %w", err)
		}
	}

	return &Actual{Value: val}, nil
}

type rawAnomalous struct {
	Value    json.RawMessage `json:"value"`
	ValueMin json.RawMessage `json:"value_min"`
	ValueMax json.RawMessage `json:"value_max"`
}

func (o rawAnomalous) Parse() (*Anomalous, error) {
	var val, min, max string
	var err error

	if o.Value != nil {
		val, err = unpackIntOrStringAsString(o.Value)
		if err != nil {
			return nil, fmt.Errorf("error parsing rawAnomalous 'value' payload - %w", err)
		}
	}

	if o.ValueMin != nil {
		min, err = unpackIntOrStringAsString(o.ValueMin)
		if err != nil {
			return nil, fmt.Errorf("error parsing rawAnomalous 'value_min' payload - %w", err)
		}
	}

	if o.ValueMax != nil {
		max, err = unpackIntOrStringAsString(o.ValueMax)
		if err != nil {
			return nil, fmt.Errorf("error parsing rawAnomalous 'value_max' payload - %w", err)
		}
	}

	return &Anomalous{
		Value:    val,
		ValueMin: min,
		ValueMax: max,
	}, nil
}

func (o rawAnomalyProperty) Parse() (*AnomalyProperty, error) {
	val, err := unpackIntOrStringAsString(o.Value)
	if err != nil {
		return nil, fmt.Errorf("error unpacking unpredictable JSON payload - %w", err)
	}
	return &AnomalyProperty{
		Key:   o.Key,
		Value: val,
	}, nil
}

func unpackIntOrStringAsString(raw json.RawMessage) (string, error) {
	var intVal int
	var stringVal string
	var unmarshalTypeError *json.UnmarshalTypeError
	err := json.Unmarshal(raw, &stringVal)
	switch {
	case err == nil:
		return stringVal, nil
	case errors.As(err, &unmarshalTypeError) && unmarshalTypeError.Value == "number":
		// oops, this is probably a number
		err = json.Unmarshal(raw, &intVal)
		if err != nil {
			return "", fmt.Errorf("error unmarshaling data detected as number '%c' - %w", raw, err)
		}
		return strconv.Itoa(intVal), nil
	default:
		return "", fmt.Errorf("unhandled error in unmarshaling anomaly payload - %w", err)
	}
}

// also available in scotch/schemas/alerts.py:
//   BGP_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {                 // node
//   BLUEPRINT_RENDERING_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, { // node
//   CABLING_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {             // node
//   CONFIG_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {              // node
//   CONFIG_MISMATCH_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {     // node
//   DEPLOYMENT_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {          // node
//   EXTENSIBLE_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {          // ?
//   HOSTNAME_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {            // node
//   INTERFACE_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {           // node
//   LAG_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {                 // ?
//   LIVENESS_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {            // node
//   MAC_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {                 // ?
//   MLAG_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {                // ?
//   ROUTE_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {               // node
//   STREAMING_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {           // ?
//   PROBE_ANOMALY_SCHEMA = t.Object(BASE_ANOMALY_SCHEMA, {               // node

type BlueprintAnomaly struct {
	Id             ObjectId   `json:"id"`               // part of base schema
	LastModifiedAt *time.Time `json:"last_modified_at"` // part of base schema
	Severity       string     `json:"severity"`         // part of base schema
	AnomalyType    string     `json:"anomaly_type"`     // part of base schema

	Actual    json.RawMessage `json:"actual"`    // universal (near universal?)
	Expected  json.RawMessage `json:"expected"`  // universal (near universal?)
	Identity  json.RawMessage `json:"identity"`  // universal (near universal?)
	Role      *string         `json:"role"`      // near universal
	Anomalous json.RawMessage `json:"anomalous"` // probe
}

type BlueprintServiceAnomalyCount struct {
	AnomalyType string `json:"type"`
	Role        string `json:"role"`
	Count       int    `json:"count"`
}

// per JP Senior: I think the a reliable list is aos.reference_design.two_stage_l3clos.__init__.py's alert_types list:
// aos/reference_design/two_stage_l3clos/__init__.py:
// alert_types = ['bgp', 'cabling', 'counter', 'interface', 'hostname', 'liveness',
//               'route', 'config', 'deployment', 'blueprint_rendering', 'probe',
//               'streaming', 'mac', 'arp', 'lag', 'mlag', 'series',
//               'all']

type BlueprintNodeAnomalyCounts struct {
	Node     string   `json:"node"`
	SystemId ObjectId `json:"system_id"`
	All      int      `json:"all"`

	Arp                int `json:"arp"`
	Bgp                int `json:"bgp"`
	BlueprintRendering int `json:"blueprint_rendering"`
	Cabling            int `json:"cabling"`
	Config             int `json:"config"`
	Counter            int `json:"counter"`
	Deployment         int `json:"deployment"`
	Hostname           int `json:"hostname"`
	Interface          int `json:"interface"`
	Lag                int `json:"lag"`
	Liveness           int `json:"liveness"`
	Mac                int `json:"mac"`
	Mlag               int `json:"mlag"`
	Probe              int `json:"probe"`
	Route              int `json:"route"`
	Series             int `json:"series"`
	Streaming          int `json:"streaming"`
}

func (o *Client) getBlueprintAnomalies(ctx context.Context, blueprintId ObjectId) ([]BlueprintAnomaly, error) {
	var apiResonse struct {
		Items []BlueprintAnomaly
	}

	err := o.talkToApstra(ctx, talkToApstraIn{
		method:         http.MethodGet,
		urlStr:         fmt.Sprintf(apiUrlBlueprintAnomalies, blueprintId),
		apiResponse:    &apiResonse,
		unsynchronized: true,
	})
	return apiResonse.Items, convertTtaeToAceWherePossible(err)
}

func (o *Client) getBlueprintNodeAnomalyCounts(ctx context.Context, blueprintId ObjectId) ([]BlueprintNodeAnomalyCounts, error) {
	var apiResonse struct {
		Items []BlueprintNodeAnomalyCounts
	}

	err := o.talkToApstra(ctx, talkToApstraIn{
		method:         http.MethodGet,
		urlStr:         fmt.Sprintf(apiUrlBlueprintAnomaliesByNode, blueprintId),
		apiResponse:    &apiResonse,
		unsynchronized: true,
	})
	return apiResonse.Items, convertTtaeToAceWherePossible(err)
}

func (o *Client) getBlueprintServiceAnomalyCounts(ctx context.Context, blueprintId ObjectId) ([]BlueprintServiceAnomalyCount, error) {
	var apiResonse struct {
		Items []BlueprintServiceAnomalyCount
	}

	err := o.talkToApstra(ctx, talkToApstraIn{
		method:         http.MethodGet,
		urlStr:         fmt.Sprintf(apiUrlBlueprintAnomaliesByService, blueprintId),
		apiResponse:    &apiResonse,
		unsynchronized: true,
	})
	return apiResonse.Items, convertTtaeToAceWherePossible(err)
}
