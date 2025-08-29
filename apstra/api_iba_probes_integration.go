// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	apiUrlIbaProbes       = apiUrlBlueprintById + apiUrlPathDelim + "probes"
	apiUrlIbaProbesPrefix = apiUrlIbaProbes + apiUrlPathDelim
	apiUrlIbaProbesById   = apiUrlIbaProbesPrefix + "%s"
)

type IbaProbe struct {
	Id              ObjectId                 `json:"id"`
	Label           string                   `json:"label"`
	Tags            []string                 `json:"tags"`
	Stages          []map[string]interface{} `json:"stages"`
	Processors      []map[string]interface{} `json:"processors"`
	PredefinedProbe string                   `json:"predefined_probe"`
	Description     string                   `json:"description"`
}

type IbaProbeState struct {
	Id              ObjectId                 `json:"id"`
	Label           string                   `json:"label"`
	TaskError       string                   `json:"task_error"`
	Stages          []map[string]interface{} `json:"stages"`
	Processors      []map[string]interface{} `json:"processors"`
	AnomalyCount    int                      `json:"anomaly_count"`
	Tags            []string                 `json:"tags"`
	Disabled        bool                     `json:"disabled"`
	State           string                   `json:"state"`
	Version         int                      `json:"version"`
	TaskState       string                   `json:"task_state"`
	IbaUnit         string                   `json:"iba_unit"`
	PredefinedProbe string                   `json:"predefined_probe"`
	Description     string                   `json:"description"`
}

func (o *IbaProbeState) IbaProbe() IbaProbe {
	return IbaProbe{
		Id:              o.Id,
		Label:           o.Label,
		Tags:            o.Tags,
		Stages:          o.Stages,
		Processors:      o.Processors,
		PredefinedProbe: o.PredefinedProbe,
		Description:     o.Description,
	}
}

func (o *Client) getAllIbaProbes(ctx context.Context, bpId ObjectId) ([]IbaProbe, error) {
	response := &struct {
		Items []IbaProbe `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaProbes, bpId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getIbaProbeByLabel(ctx context.Context, bpId ObjectId, label string) (*IbaProbe, error) {
	pps, err := o.getAllIbaProbes(ctx, bpId)
	if err != nil {
		return nil, err
	}
	var probe IbaProbe
	i := 0
	for _, p := range pps {
		if p.Label == label {
			probe = p
			i = i + 1
		}
	}
	if i == 0 {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no Predefined Probe with label '%s' found", label),
		}
	}
	if i > 1 {
		return nil, ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("too many probes with label %s found, expected 1 got %d", label, i),
		}
	}
	return &probe, nil
}

func (o *Client) getIbaProbe(ctx context.Context, bpId ObjectId, id ObjectId) (*IbaProbe, error) {
	response := &IbaProbe{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaProbesById, bpId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getIbaProbeState(ctx context.Context, bpId ObjectId, id ObjectId) (*IbaProbeState, error) {
	response := &IbaProbeState{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaProbesById, bpId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) deleteIbaProbe(ctx context.Context, bpId ObjectId, id ObjectId) error {
	return convertTtaeToAceWherePossible(o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlIbaProbesById, bpId, id),
	}))
}

func (o *Client) createIbaProbeFromJson(ctx context.Context, bpId ObjectId, probeJson json.RawMessage) (ObjectId, error) {
	var response objectIdResponse
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlIbaProbes, bpId),
		apiInput:    probeJson,
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

	retryMax := o.GetTuningParam("createProbeMaxRetries")
	retryInterval := time.Duration(o.GetTuningParam("createProbeRetryIntervalMs")) * time.Millisecond

	for i := 0; i < retryMax; i++ {
		// Make a random wait, in case multiple threads are running
		if rand.Int()%2 == 0 {
			time.Sleep(retryInterval)
		}

		time.Sleep(retryInterval * time.Duration(i))

		e := o.talkToApstra(ctx, &talkToApstraIn{
			method:      http.MethodPost,
			urlStr:      fmt.Sprintf(apiUrlIbaProbes, bpId),
			apiInput:    probeJson,
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
