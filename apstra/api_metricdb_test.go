//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestGetMetricdbMetrics(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getMetricdbMetrics() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, err = client.client.getMetricdbMetrics(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestUseAggregation(t *testing.T) {
	type testData struct {
		testString string
		expectBool bool
		expectName string
		expectSecs int
	}

	var td []testData
	td = append(td, testData{testString: "foo", expectBool: false, expectName: "foo", expectSecs: 0})
	td = append(td, testData{testString: "foo_aggr_", expectBool: false, expectName: "foo_aggr_", expectSecs: 0})
	td = append(td, testData{testString: "_aggr_3600", expectBool: false, expectName: "_aggr_3600", expectSecs: 0})
	td = append(td, testData{testString: "foo_aggr_-3600", expectBool: false, expectName: "foo_aggr_-3600", expectSecs: 0})
	td = append(td, testData{testString: "foo_aggr_3600", expectBool: true, expectName: "foo", expectSecs: 3600})

	for i := range td {
		useAgg, name, secs, err := useAggregation(td[i].testString)
		if err != nil {
			t.Fatal(err)
		}
		if useAgg != td[i].expectBool {
			t.Fatalf("'%s' expected bool '%t', got '%t'", td[i].testString, td[i].expectBool, useAgg)
		}
		if secs != td[i].expectSecs {
			t.Fatalf("'%s' expected time '%d', got '%d'", td[i].testString, td[i].expectSecs, secs)
		}
		if name != td[i].expectName {
			t.Fatalf("'%s' expected name '%s', got '%s'", td[i].testString, td[i].expectName, name)
		}
	}
}

func TestQueryMetricdb(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getMetricdbMetrics() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		metrics, err := client.client.getMetricdbMetrics(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		var result *MetricDbQueryResponse
		if len(metrics.Items) > 0 { // do not call rand.Intn() with '0'
			i := rand.Intn(len(metrics.Items))
			log.Printf("randomly requesting metric %d of %d available", i, len(metrics.Items))
			q := MetricDbQueryRequest{
				metric: metrics.Items[i],
				begin:  time.Now().Add(-time.Hour),
				end:    time.Now(),
			}

			log.Printf("testing QueryMetricdb() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			result, err = client.client.QueryMetricdb(context.TODO(), &q)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("got %d results for the last hour of %s/%s/%s",
				len(result.Items), q.metric.Application, q.metric.Namespace, q.metric.Name)
		}
	}
}
