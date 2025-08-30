// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"time"
)

const (
	VniMin = vniMin
	VniMax = vniMax
)

// GetFeatures is in the export_test file because this private function needs to be exposed only for test code
func (o *Client) GetFeatures(ctx context.Context) error {
	return o.getFeatures(ctx)
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
