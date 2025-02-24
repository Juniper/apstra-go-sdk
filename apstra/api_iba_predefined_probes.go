// Copyright (c) Juniper Networks, Inc., 2023-2025.
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
	apiUrlIbaPredefinedProbes       = "/api/blueprints/%s/iba/predefined-probes"
	apiUrlIbaPredefinedProbesPrefix = apiUrlIbaPredefinedProbes + apiUrlPathDelim
	apiUrlIbaPredefinedProbesByName = apiUrlIbaPredefinedProbesPrefix + "%s"
)

type IbaPredefinedProbe struct {
	Name         string          `json:"name"`
	Experimental bool            `json:"experimental"`
	Description  string          `json:"description"`
	Schema       json.RawMessage `json:"schema"`
}

type IbaPredefinedProbeRequest struct {
	Name string
	Data json.RawMessage
}

func (o *Client) getAllIbaPredefinedProbes(ctx context.Context, bpId ObjectId) ([]IbaPredefinedProbe, error) {
	response := &struct {
		Items []IbaPredefinedProbe `json:"items"`
	}{}

	err := o.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaPredefinedProbes, bpId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getIbaPredefinedProbeByName(ctx context.Context, bpId ObjectId, name string) (*IbaPredefinedProbe, error) {
	pps, err := o.getAllIbaPredefinedProbes(ctx, bpId)
	if err != nil {
		return nil, err
	}

	for _, p := range pps {
		if p.Name == name {
			return &p, nil
		}
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no Predefined Probe with name '%s' found", name),
	}
}

func (o *Client) instantiatePredefinedIbaProbe(ctx context.Context, bpId ObjectId, in *IbaPredefinedProbeRequest) (ObjectId, error) {
	var response objectIdResponse
	err := o.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlIbaPredefinedProbesByName, bpId, in.Name),
		apiInput:    in.Data,
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

	retryMax := o.GetTuningParam("ibaPredefinedProbeMaxRetries")
	retryInterval := time.Duration(o.GetTuningParam("ibaPredefinedProbeRetryIntervalMs")) * time.Millisecond

	for i := 0; i < retryMax; i++ {
		// Make a random wait, in case multiple threads are running
		if rand.Int()%2 == 0 {
			time.Sleep(retryInterval)
		}

		time.Sleep(retryInterval * time.Duration(i))

		e := o.talkToApstra(ctx, talkToApstraIn{
			method:      http.MethodPost,
			urlStr:      fmt.Sprintf(apiUrlIbaPredefinedProbesByName, bpId, in.Name),
			apiInput:    in.Data,
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
