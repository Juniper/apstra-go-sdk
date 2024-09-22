// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetMetricdbMetrics(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	require.NoError(t, err)

	for clientName, client := range clients {
		log.Printf("testing getMetricdbMetrics() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, err = client.client.getMetricdbMetrics(context.TODO())
		require.NoError(t, err)
	}
}

func TestQueryMetricdb(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			log.Printf("testing getMetricdbMetrics() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			metrics, err := client.client.GetMetricdbMetrics(ctx)
			require.NoError(t, err)

			for i := len(metrics) - 1; i >= 0; i-- {
				if metrics[i].Application == "audit" {
					// we cannot use these for some reason...
					metrics[i] = metrics[len(metrics)-1]
					metrics = metrics[:len(metrics)-1]
				}
			}

			var result *MetricDbQueryResponse
			if len(metrics) > 0 { // do not call rand.Intn() with '0'
				i := rand.Intn(len(metrics))
				log.Printf("randomly requesting metric %q (%d) of %d available", metrics[i], i, len(metrics))
				q := MetricDbQueryRequest{
					metric: metrics[i],
					begin:  time.Now().Add(-time.Hour),
					end:    time.Now(),
				}

				log.Printf("testing QueryMetricdb() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				result, err = client.client.QueryMetricdb(ctx, &q)
				require.NoError(t, err)
				log.Printf("got %d results for the last hour of %s/%s/%s",
					len(result.Items), q.metric.Application, q.metric.Namespace, q.metric.Name)
			}
		})
	}
}
