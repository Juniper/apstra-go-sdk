package goapstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	apiUrlAnomalies = "/api/anomalies"
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

// unpackAnomaly is clunky. It extracts instances of Anomaly as returned by
// apiUrlAnomalies, and attempts to gracefully handle the unpredictable
// Anomaly.Identity.Properties list as returned by Apstra. Sometimes it
// sends strings, sometimes it sends integers.
func unpackAnomaly(in []byte) (*Anomaly, error) {
	rawAnomaly := &rawAnomaly{}
	err := json.Unmarshal(in, rawAnomaly)
	if err != nil {
		return nil, err
	}

	actual, err := rawAnomaly.Actual.Parse()
	if err != nil {
		return nil, fmt.Errorf("error unpacking raw anomaly 'actual' - %w", err)
	}

	anomalous, err := rawAnomaly.Anomalous.Parse()
	if err != nil {
		return nil, fmt.Errorf("error unpacking raw anomaly 'anomalous' - %w", err)
	}

	anomaly := &Anomaly{}

	anomaly.Actual = *actual
	anomaly.Anomalous = *anomalous
	anomaly.AnomalyType = rawAnomaly.AnomalyType
	anomaly.Id = rawAnomaly.Id
	anomaly.LastModifiedAt = rawAnomaly.LastModifiedAt
	anomaly.Severity = rawAnomaly.Severity
	anomaly.Identity.AnomalyType = rawAnomaly.Identity.AnomalyType
	anomaly.Identity.ItemId = rawAnomaly.Identity.ItemId
	anomaly.Identity.ProbeId = rawAnomaly.Identity.ProbeId
	anomaly.Identity.ProbeLabel = rawAnomaly.Identity.ProbeLabel
	anomaly.Identity.StageName = rawAnomaly.Identity.StageName

	for _, raw := range rawAnomaly.Identity.Properties {
		property, err := raw.Parse()
		if err != nil {
			return nil, fmt.Errorf("error parsing rawAnomalyProperty - %w", err)
		}
		anomaly.Identity.Properties = append(anomaly.Identity.Properties, *property)
	}

	return anomaly, nil
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

func (o *Client) getAnomalies(ctx context.Context) ([]Anomaly, error) {
	apstraUrl, err := url.Parse(apiUrlAnomalies)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlAnomalies, err)
	}
	response := &struct {
		Items []json.RawMessage `json:"items"`
	}{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodGet,
		url:            apstraUrl,
		apiResponse:    response,
		unsynchronized: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting anomalies - %w", err)
	}
	var result []Anomaly
	for _, ra := range response.Items {
		a, err := unpackAnomaly(ra)
		if err != nil {
			return nil, fmt.Errorf("error unpacking anomaly - %w", err)
		}
		result = append(result, *a)
	}
	return result, nil
}
