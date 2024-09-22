// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	apiUrlLeafServerLinkLabels = apiUrlBlueprintByIdPrefix + "leaf-server-link-labels"
)

type LinkLagParams struct {
	GroupLabel string
	LagMode    RackLinkLagMode
	Tags       []string
}

func (o *LinkLagParams) raw() (*rawLinkLagParams, error) {
	groupLabel := o.GroupLabel
	if groupLabel == "" {
		initUUID()
		uuid1, err := uuid.NewUUID()
		if err != nil {
			return nil, fmt.Errorf("error generating type 1 uuid - %w", err)
		}
		groupLabel = uuid1.String()
	}

	return &rawLinkLagParams{
		GroupLabel: groupLabel,
		LagMode:    rackLinkLagMode(o.LagMode.String()),
		Tags:       o.Tags,
	}, nil
}

type rawLinkLagParams struct {
	GroupLabel string          `json:"group_label"`
	LagMode    rackLinkLagMode `json:"lag_mode,omitempty"`
	Tags       []string        `json:"tags"`
}

// SetLinkLagParamsRequest is a map of LAG parameters keyed by link node ID
type SetLinkLagParamsRequest map[ObjectId]LinkLagParams

// SetLinkLagParams configures the links identified in the request
// Links with no supplied GroupLabel will be given a unique random label making
// them the only members of their own group.
func (o *TwoStageL3ClosClient) SetLinkLagParams(ctx context.Context, req *SetLinkLagParamsRequest) error {
	var apiInput struct {
		Requests map[ObjectId]rawLinkLagParams `json:"links"`
	}
	apiInput.Requests = make(map[ObjectId]rawLinkLagParams)
	for k, v := range *req {
		raw, err := v.raw()
		if err != nil {
			return err
		}
		apiInput.Requests[k] = *raw
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlLeafServerLinkLabels, o.blueprintId),
		apiInput: &apiInput,
	})
	if err == nil {
		return nil // success!
	}

	// if we got here, then we have an error
	err = convertTtaeToAceWherePossible(err)
	var ace ClientErr
	if !errors.As(err, &ace) {
		return err // cannot handle
	}

	if ace.Type() != ErrLagHasAssignedStructrues {
		return err // cannot handle
	}

	var ds detailedStatus
	if json.Unmarshal([]byte(ace.Error()), &ds) != nil {
		return err // unmarshal fail - surface the original error
	}

	// unpack the error
	var e struct {
		Raw json.RawMessage `json:"links"`
	}
	if json.Unmarshal(ds.Errors, &e) != nil {
		return err // unmarshal fail - surface the original error
	}

	var errors []string
	if bytes.HasPrefix(bytes.TrimSpace(e.Raw), []byte{'['}) {
		// raw value is an array
		if json.Unmarshal(e.Raw, &errors) != nil {
			return err // unmarshal fail - surface the original error
		}
	} else {
		// raw value is a singleton
		var s string
		if json.Unmarshal(e.Raw, &s) != nil {
			return err // unmarshal fail - surface the original error
		}
		errors = []string{s}
	}
	if len(errors) == 0 {
		return err // didn't find embedded errors - surface the original error
	}

	aceDetail := ErrLagHasAssignedStructuresDetail{GroupLabels: make([]string, len(errors))}
	for i, s := range errors {
		m := regexpLagHasAssignedStructures.FindStringSubmatch(s)
		if len(m) != 2 {
			return fmt.Errorf("cannot handle lag with assigned structures error %q - %w", s, err)
		}

		aceDetail.GroupLabels[i] = m[1]
	}

	ace.detail = aceDetail
	return ace
}
