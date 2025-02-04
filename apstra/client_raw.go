// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"errors"
	"net/url"
)

type RawJsonRequest struct {
	Method  string
	Url     *url.URL
	Payload any
}

func (o *Client) DoRawJsonTransaction(ctx context.Context, req RawJsonRequest, resp any) error {
	if req.Url == nil {
		return errors.New("no url provided")
	}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      req.Method,
		url:         req.Url,
		apiInput:    req.Payload,
		apiResponse: resp,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
