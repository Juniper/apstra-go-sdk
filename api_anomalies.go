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
	Actual struct {
		Value int `json:"value"`
	} `json:"actual"`
	Anomalous struct {
		ValueMax int `json:"value_max"`
	} `json:"anomalous"`
	AnomalyType string   `json:"anomaly_type"`
	Id          ObjectId `json:"id"`
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

type rawAnomalyProperty struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

type rawAnomaly struct {
	Actual struct {
		Value int `json:"value"`
	} `json:"actual"`
	Anomalous struct {
		ValueMax int `json:"value_max"`
	} `json:"anomalous"`
	AnomalyType string   `json:"anomaly_type"`
	Id          ObjectId `json:"id"`
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

type getAnomaliesResponse struct {
	Items []json.RawMessage `json:items`
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
	anomaly := &Anomaly{}

	anomaly.Actual = rawAnomaly.Actual
	anomaly.Anomalous = rawAnomaly.Anomalous
	anomaly.AnomalyType = rawAnomaly.AnomalyType
	anomaly.Id = rawAnomaly.Id
	anomaly.LastModifiedAt = rawAnomaly.LastModifiedAt
	anomaly.Severity = rawAnomaly.Severity
	anomaly.Identity.AnomalyType = rawAnomaly.Identity.AnomalyType
	anomaly.Identity.ItemId = rawAnomaly.Identity.ItemId
	anomaly.Identity.ProbeId = rawAnomaly.Identity.ProbeId
	anomaly.Identity.ProbeLabel = rawAnomaly.Identity.ProbeLabel
	anomaly.Identity.StageName = rawAnomaly.Identity.StageName

	var intVal int
	var stringVal string
	var unmarshalTypeError *json.UnmarshalTypeError

	for _, p := range rawAnomaly.Identity.Properties {
		// try unmarshaling as a string
		err = json.Unmarshal(p.Value, &stringVal)
		switch {
		case err == nil:
			anomaly.Identity.Properties = append(anomaly.Identity.Properties, AnomalyProperty{
				Key:   p.Key,
				Value: stringVal,
			})
		case errors.As(err, &unmarshalTypeError) && unmarshalTypeError.Value == "number":
			// oops, this is probably a number
			err = json.Unmarshal(p.Value, &intVal)
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling data detected as number '%c' - %w", p.Value, err)
			}
			anomaly.Identity.Properties = append(anomaly.Identity.Properties, AnomalyProperty{
				Key:   p.Key,
				Value: strconv.Itoa(intVal),
			})
		default:
			return nil, fmt.Errorf("unhandled error in unmarshaling anomaly payload - %w", err)
		}
		if err != nil && errors.Is(err, new(json.UnmarshalTypeError)) {
			return nil, err
		}
	}

	//anomaly := Anomaly(*rawAnomaly)
	//anomaly = rawAnomaly
	//return anomaly, nil
	return anomaly, nil
}

func (o *Client) getAnomalies(ctx context.Context) ([]*Anomaly, error) {
	apstraUrl, err := url.Parse(apiUrlAnomalies)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlAnomalies, err)
	}
	response := &getAnomaliesResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodGet,
		url:            apstraUrl,
		apiResponse:    response,
		unsynchronized: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting anomalies - %w", err)
	}
	var result []*Anomaly
	for _, ra := range response.Items {
		a, err := unpackAnomaly(ra)
		if err != nil {
			return nil, fmt.Errorf("error unpacking anomaly - %w", err)
		}
		result = append(result, a)
	}
	return result, nil
}
