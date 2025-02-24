// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"net/http"
)

const (
	apiUrlConfigAudit = apiUrlConfigPrefix + "audit"
)

type AuditConfig struct {
	Syslogs []string `json:"syslogs"`
}

func (o *Client) getAuditConfig(ctx context.Context) (*AuditConfig, error) {
	response := &AuditConfig{}
	err := o.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlConfigAudit,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *Client) putAuditConfig(ctx context.Context, config *AuditConfig) error {
	err := o.talkToApstra(ctx, talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   apiUrlConfigAudit,
		apiInput: config,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
