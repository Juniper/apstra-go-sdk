// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetMetricdbMetrics(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)
			_, err := client.Client.GetMetricdbMetrics(ctx)
			require.NoError(t, err)
		})
	}
}

func TestQueryMetricdb(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()

			metrics, err := client.Client.GetMetricdbMetrics(ctx)
			require.NoError(t, err)

			for i := len(metrics) - 1; i >= 0; i-- {
				if metrics[i].Application == "audit" {
					// we cannot use these for some reason...
					metrics[i] = metrics[len(metrics)-1]
					metrics = metrics[:len(metrics)-1]
				}
			}

			var result *apstra.MetricDbQueryResponse
			if len(metrics) > 0 { // do not call rand.Intn() with '0'
				i := rand.Intn(len(metrics))
				log.Printf("randomly requesting metric %q (%d) of %d available", metrics[i], i, len(metrics))
				var q apstra.MetricDbQueryRequest
				q.SetMetric(metrics[i])
				q.SetBegin(time.Now().Add(-time.Hour))
				q.SetEnd(time.Now())

				result, err = client.Client.QueryMetricdb(ctx, &q)
				require.NoError(t, err)
				m := q.Metric()
				log.Printf("got %d results for the last hour of %s/%s/%s", len(result.Items), m.Application, m.Namespace, m.Name)
			}
		})
	}
}
