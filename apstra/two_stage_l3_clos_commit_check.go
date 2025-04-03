// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	apiUrlBlueprintCommitCheckAll    = "/api/blueprints/%s/commit-check"
	apiUrlBlueprintCommitCheckSystem = "/api/blueprints/%s/systems/%s/commit-check"

	apiUrlBlueprintCommitCheckAllResults    = "/api/blueprints/%s/commit-check-result"
	apiUrlBlueprintCommitCheckSystemResults = "/api/blueprints/%s/systems/%s/commit-check-result"
)

// RunCommitCheck kicks off a configuration validation event. If id is non-nil, then the
// validation will be requested only for the system with the given id. If configString is
// non-nil, then the validation will use the supplied configuration rather than the
// configuration automatically rendered for the given id. configString must not be
// supplied without supplying id.
func (o *TwoStageL3ClosClient) RunCommitCheck(ctx context.Context, id *ObjectId, configString *string) error {
	if configString != nil && id == nil {
		return fmt.Errorf("RunCommitCheck: cfg requires id")
	}

	var urlStr string
	if id == nil {
		urlStr = fmt.Sprintf(apiUrlBlueprintCommitCheckAll, o.blueprintId)
	} else {
		urlStr = fmt.Sprintf(apiUrlBlueprintCommitCheckSystem, o.blueprintId, *id)
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPost,
		urlStr: urlStr,
		apiInput: struct {
			ConfigString *string `json:"config_string,omitempty"`
		}{
			ConfigString: configString,
		},
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

type TwoStageL3ClosCommitCheckResult struct {
	Blueprint struct {
		State            enum.CommitCheckState     `json:"state"`
		Validity         *enum.CommitCheckValidity `json:"validity"`
		BlueprintVersion *uint                     `json:"blueprint_version"`
	} `json:"blueprint"`
	Systems map[ObjectId]SystemCommitCheckResult `json:"systems"`
}

var _ json.Unmarshaler = &SystemCommitCheckResult{}

type SystemCommitCheckResult struct {
	State            enum.CommitCheckState
	Error            *string
	Source           enum.CommitCheckSource
	FinishedAt       *time.Time
	BlueprintId      *ObjectId
	BlueprintVersion *int64
	Validity         *enum.CommitCheckValidity
}

func (o *SystemCommitCheckResult) UnmarshalJSON(b []byte) error {
	var raw struct {
		State            enum.CommitCheckState     `json:"state"`
		Source           enum.CommitCheckSource    `json:"source"`
		FinishedAt       *time.Time                `json:"finished_at"`
		BlueprintId      *ObjectId                 `json:"blueprint_id"`
		BlueprintVersion *int64                    `json:"blueprint_version"`
		Validity         *enum.CommitCheckValidity `json:"validity"`
		Error            string                    `json:"error"`
	}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	o.State = raw.State
	o.Source = raw.Source
	o.FinishedAt = raw.FinishedAt
	o.BlueprintId = raw.BlueprintId
	o.BlueprintVersion = raw.BlueprintVersion
	o.Validity = raw.Validity

	if raw.Error != "" {
		o.Error = &raw.Error
	}

	return nil
}

func (o *TwoStageL3ClosClient) CommitCheckResults(ctx context.Context) (*TwoStageL3ClosCommitCheckResult, error) {
	var result TwoStageL3ClosCommitCheckResult

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintCommitCheckAllResults, o.blueprintId),
		apiResponse: &result,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &result, err
}

func (o *TwoStageL3ClosClient) CommitCheckResult(ctx context.Context, id ObjectId) (*SystemCommitCheckResult, error) {
	var result SystemCommitCheckResult

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintCommitCheckSystemResults, o.blueprintId, id),
		apiResponse: &result,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &result, err
}
