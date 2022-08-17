package goapstra

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"testing"
	"time"
)

const (
	mockGetMetricdbIbaAppPayload = `{"items":[
	   { "application": "iba",
	     "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665",
	     "name": "imbalanced_system_count_out_of_range" },
	   { "application": "iba",
	     "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/4eb11184-4b32-4106-8e90-edb312042683",
	     "name": "System Interface Counters" },
	   { "application": "iba",
	     "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/4eb11184-4b32-4106-8e90-edb312042683",
	     "name": "Average Interface Counters" },
	   { "application": "iba",
	     "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665",
	     "name": "leaf_fab_int_tx_avg" },
	   { "application": "iba",
	     "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665",
	     "name": "std_dev_percentage" },
	   { "application": "iba",
	     "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665",
	     "name": "system_imbalance" }
	]}`
	mockGetMetricdbClusterHealthAppPayload = `{"items":[
        { "application": "cluster_health_info",
          "namespace": "agent",
          "name": "health_aggr_3600" },
        { "application": "cluster_health_info",
          "namespace": "file_registry",
          "name": "file_aggr_3600" },
        { "application": "cluster_health_info",
          "namespace": "file_registry",
          "name": "directory_aggr_3600" },
        { "application": "cluster_health_info",
          "namespace": "agent",
          "name": "health" },
        { "application": "cluster_health_info",
          "namespace": "agent",
          "name": "utilization" },
        { "application": "cluster_health_info",
          "namespace": "node",
          "name": "utilization_aggr_3600" },
        { "application": "cluster_health_info",
          "namespace": "node",
          "name": "disk_utilization_aggr_3600" },
        { "application": "cluster_health_info",
          "namespace": "container",
          "name": "utilization" },
        { "application": "cluster_health_info",
          "namespace": "node",
          "name": "utilization" },
        { "application": "cluster_health_info",
          "namespace": "file_registry",
          "name": "directory" },
        { "application": "cluster_health_info",
          "namespace": "node",
          "name": "disk_utilization" },
        { "application": "cluster_health_info",
          "namespace": "file_registry",
          "name": "file" },
        { "application": "cluster_health_info",
          "namespace": "container",
          "name": "utilization_aggr_3600" },
        { "application": "cluster_health_info",
          "namespace": "agent",
          "name": "utilization_aggr_3600" },
        { "application": "cluster_health_info",
          "namespace": "container",
          "name": "file_usage_aggr_3600" },
        { "application": "cluster_health_info",
          "namespace": "container",
          "name": "file_usage" }
    ]}`
)

func (o *mockApstraApi) createMetricdb() error {
	var err error

	var metricdbClusterHealthInfo metricdbMetricResponse
	err = json.Unmarshal([]byte(mockGetMetricdbClusterHealthAppPayload), &metricdbClusterHealthInfo)
	if err != nil {
		return err
	}

	var metricdbIbaInfo metricdbMetricResponse
	o.metricdb.metrics.Items = append(o.metricdb.metrics.Items, metricdbClusterHealthInfo.Items...)
	err = json.Unmarshal([]byte(mockGetMetricdbIbaAppPayload), &metricdbIbaInfo)
	if err != nil {
		return err
	}
	o.metricdb.metrics.Items = append(o.metricdb.metrics.Items, metricdbIbaInfo.Items...)

	return nil
}

func (o *mockApstraApi) createMetricdbIbaData(blueprintId ObjectId, ibaName string, ibaId ObjectId) error {
	o.metricdb.metrics.Items = append(o.metricdb.metrics.Items, MetricdbMetric{
		Application: "iba",
		Namespace:   string(blueprintId) + apiUrlPathDelim + string(ibaId),
		Name:        ibaName,
	})
	return nil
}

func TestUnmarshalMockMetricdbData(t *testing.T) {
	result := &metricdbMetricResponse{}
	for i, s := range []string{
		mockGetMetricdbIbaAppPayload,
		mockGetMetricdbClusterHealthAppPayload,
	} {
		log.Printf("TestUnmarshalMockMetricdbData test %d", i)
		err := json.Unmarshal([]byte(s), result)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetMetricdbMetrics(t *testing.T) {
	clients, apis, err := getTestClientsAndMockAPIs()
	if err != nil {
		t.Fatal(err)
	}

	_, mockExists := apis["mock"]
	if mockExists {
		err = apis["mock"].createMetricdb()
		if err != nil {
			log.Fatal(err)
		}
	}

	for clientName, client := range clients {
		log.Printf("testing getMetricdbMetrics() with %s client", clientName)
		err := client.Login(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		_, err = client.getMetricdbMetrics(context.TODO())
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
	clients, apis, err := getTestClientsAndMockAPIs()
	if err != nil {
		t.Fatal(err)
	}

	// todo: implement mock metricdb query

	_, mockExists := apis["mock"]
	if mockExists {
		err = apis["mock"].createMetricdb()
		if err != nil {
			log.Fatal(err)
		}
	}

	rand.Seed(time.Now().UnixNano())

	for clientName, client := range clients {
		log.Printf("testing getMetricdbMetrics() with %s client", clientName)
		err := client.Login(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		metrics, err := client.getMetricdbMetrics(context.TODO())
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

			result, err = client.QueryMetricdb(context.TODO(), &q)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("got %d results for the last hour of %s/%s/%s",
				len(result.Items), q.metric.Application, q.metric.Namespace, q.metric.Name)
			for i := range result.Items {
				log.Printf(string(result.Items[i]))
			}
		}
	}
}
