package goapstra

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

const (
	metricDbBlueprintPayloadTodo = `{"items":[
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
	mockGetMetricdbPayload = `{"items":[
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

//      "application": "cluster_health_info",
//      "namespace": "agent",
//      "name": "health_aggr_3600",
// Empty

//
//      "application": "cluster_health_info",
//      "namespace": "file_registry",
//      "name": "file_aggr_3600"
//  {   "node": "AosController",
//      "timestamp": "2022-05-13T00:44:09.723419Z",
//      "file_path": "/var/lib/aos/metricdb/cluster_health_info/node/disk_utilization/disk_utilization-189-2022-05-12--16-45-55.616443.tel",
//      "size": 34535 }

//
//      "application": "iba",
//      "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665",
//      "name": "imbalanced_system_count_out_of_range",
// Empty

//
//      "application": "cluster_health_info",
//      "namespace": "file_registry",
//      "name": "directory_aggr_3600",
//  {   "node": "AosController",
//      "timestamp": "2022-05-13T00:44:09.730437Z",
//      "size": 78081287,
//      "directory_path": "/var/lib/aos/metricdb/cluster_health_info" }

//
//      "application": "cluster_health_info",
//      "namespace": "agent",
//      "name": "health",
// Empty

//
//      "application": "cluster_health_info",
//      "namespace": "agent",
//      "name": "utilization",
//  {   "node": "AosController",
//      "container": "aos_sysdb_1",
//      "timestamp": "2022-05-13T00:00:00.186722Z",
//      "agent": "CognacSysdb",
//      "memory": 151576576,
//      "cpu": 0 }

//
//      "application": "iba",
//      "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/4eb11184-4b32-4106-8e90-edb312042683",
//      "name": "System Interface Counters"
//  {   "timestamp": "2022-05-14T00:00:01.367573Z",
//      "aggregate_rx_bps": 1376,
//      "max_ifc_rx_utilization": 0,
//      "max_ifc_tx_utilization": 0,
//      "system_id": "WS3119350041",
//      "aggregate_rx_utilization": 0,
//      "aggregate_tx_bps": 1500,
//      "aggregate_tx_utilization": 0 }

func (o *mockApstraApi) createMetricdb() error {
	return json.Unmarshal([]byte(mockGetMetricdbPayload), &o.metricDb)
}

func (o *mockApstraApi) createMetricdbIbaData(blueprintId ObjectId, ibaName string, ibaId ObjectId) error {
	o.metricDb.Items = append(o.metricDb.Items, metricdbResponseElem{
		Application: "iba",
		Namespace:   string(blueprintId) + apiUrlPathDelim + string(ibaId),
		Name:        ibaName,
	})
	return nil
}

func TestUnmarshalMockMetricdbData(t *testing.T) {
	result := &metricdbResponse{}
	err := json.Unmarshal([]byte(mockGetMetricdbPayload), result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMetricdb(t *testing.T) {
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
		_, err := client.getMetricdb(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
	}
}
