package goapstra

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"testing"
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
	o.metricdbMetric.Items = append(o.metricdbMetric.Items, metricdbClusterHealthInfo.Items...)
	err = json.Unmarshal([]byte(mockGetMetricdbClusterHealthAppPayload), &metricdbIbaInfo)
	if err != nil {
		return err
	}
	o.metricdbMetric.Items = append(o.metricdbMetric.Items, metricdbIbaInfo.Items...)

	return nil
}

func (o *mockApstraApi) createMetricdbIbaData(blueprintId ObjectId, ibaName string, ibaId ObjectId) error {
	o.metricdbMetric.Items = append(o.metricdbMetric.Items, metricdbResponseElem{
		Application: "iba",
		Namespace:   string(blueprintId) + apiUrlPathDelim + string(ibaId),
		Name:        ibaName,
	})
	return nil
}

func TestUnmarshalMockMetricdbData(t *testing.T) {
	result := &metricdbMetricResponse{}
	err := json.Unmarshal([]byte(mockGetMetricdbClusterHealthAppPayload), result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMetricdb(t *testing.T) {
	DebugLevel = 2
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
		log.Printf("testing getMetricdb() with %s client", clientName)
		err := client.Login(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		metrics, err := client.getMetricdb(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		buf := bytes.NewBuffer([]byte{})
		err = pp(metrics, buf)
		if err != nil {
			t.Fatal(err)
		}

		debugStr(2, buf.String())
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
	td = append(td, testData{testString: "foo"})
	td = append(td, testData{testString: "foo_aggr_"})
	td = append(td, testData{testString: "_aggr_3600"})
	td = append(td, testData{testString: "foo_aggr_-3600"})
	td = append(td, testData{testString: "foo_aggr_3600", expectBool: true, expectName: "foo", expectSecs: 3600})

	for i := range td {
		useAgg, name, secs, err := useAggregation(td[i].testString)
		if err != nil {
			t.Fatal(err)
		}
		if useAgg != td[i].expectBool {
			t.Fatalf("'%s' expected '%t', got '%t'", td[i].testString, td[i].expectBool, useAgg)
		}
		if secs != td[i].expectSecs {
			t.Fatalf("'%s' expected '%d', got '%d'", td[i].testString, td[i].expectSecs, secs)
		}
		if name != td[i].expectName {
			t.Fatalf("'%s' expected '%s', got '%s'", td[i].testString, td[i].expectName, name)
		}
	}
}
