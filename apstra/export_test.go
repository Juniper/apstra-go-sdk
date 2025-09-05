// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"time"
)

const (
	VlanMin = vlanMin
	VlanMax = vlanMax
	VniMin  = vniMin
	VniMax  = vniMax
)

type (
	VersionsAosdiResponse  versionsAosdiResponse
	VersionsApiResponse    versionsApiResponse
	VersionsBuildResponse  versionsBuildResponse
	VersionsServerResponse versionsServerResponse
)

// GetVersionsAosdi is in the export_test file because this private function needs to be exposed only for test code
func (o *Client) GetVersionsAosdi(ctx context.Context) (*VersionsAosdiResponse, error) {
	result, err := o.getVersionsAosdi(ctx)
	return (*VersionsAosdiResponse)(result), err
}

// GetVersionsApi is in the export_test file because this private function needs to be exposed only for test code
func (o *Client) GetVersionsApi(ctx context.Context) (*VersionsApiResponse, error) {
	result, err := o.getVersionsApi(ctx)
	return (*VersionsApiResponse)(result), err
}

// GetVersionsBuild is in the export_test file because this private function needs to be exposed only for test code
func (o *Client) GetVersionsBuild(ctx context.Context) (*VersionsBuildResponse, error) {
	result, err := o.getVersionsBuild(ctx)
	return (*VersionsBuildResponse)(result), err
}

// GetVersionsServer is in the export_test file because this private function needs to be exposed only for test code
func (o *Client) GetVersionsServer(ctx context.Context) (*VersionsServerResponse, error) {
	result, err := o.getVersionsServer(ctx)
	return (*VersionsServerResponse)(result), err
}

// GetFeatures is in the export_test file because this private function needs to be exposed only for test code
func (o *Client) GetFeatures(ctx context.Context) error {
	return o.getFeatures(ctx)
}

// SetAuthtoken is in the export_test file because this private function needs to be exposed only for test code
func (o *Client) SetAuthtoken(t string) {
	o.httpHeaders[apstraAuthHeader] = t
}

// SetPassword is in the export_test file because this private function needs to be exposed only for test code
func (o *Client) SetPassword(p string) {
	o.cfg.Pass = p
}

// Metric is in the export_test file because this private struct element needed to be exposed for test code
func (o *MetricDbQueryRequest) Metric() MetricdbMetric {
	return o.metric
}

// SetMetric is in the export_test file because this private struct element needed to be exposed for test code
func (o *MetricDbQueryRequest) SetMetric(m MetricdbMetric) {
	o.metric = m
}

// SetBegin is in the export_test file because this private struct element needed to be exposed for test code
func (o *MetricDbQueryRequest) SetBegin(t time.Time) {
	o.begin = t
}

// SetEnd is in the export_test file because this private struct element needed to be exposed for test code
func (o *MetricDbQueryRequest) SetEnd(t time.Time) {
	o.end = t
}
