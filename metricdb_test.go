package goapstra

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

const (
	metricDbBlueprintPayloadTodo = `{"items":[
        {
          "application": "iba",
          "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665",
          "name": "imbalanced_system_count_out_of_range"
        },
        {
          "application": "iba",
          "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/4eb11184-4b32-4106-8e90-edb312042683",
          "name": "System Interface Counters"
        },
        {
          "application": "iba",
          "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/4eb11184-4b32-4106-8e90-edb312042683",
          "name": "Average Interface Counters"
        },
        {
          "application": "iba",
          "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665",
          "name": "leaf_fab_int_tx_avg"
        },
        {
          "application": "iba",
          "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665",
          "name": "std_dev_percentage"
        },
        {
          "application": "iba",
          "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665",
          "name": "system_imbalance"
        }
    ]}`
	mockGetMetricdbPayload = `{"items":[
        {
          "application": "cluster_health_info",
          "namespace": "agent",
          "name": "health_aggr_3600"
        },
        {
          "application": "cluster_health_info",
          "namespace": "file_registry",
          "name": "file_aggr_3600"
        },
        {
          "application": "cluster_health_info",
          "namespace": "file_registry",
          "name": "directory_aggr_3600"
        },
        {
          "application": "cluster_health_info",
          "namespace": "agent",
          "name": "health"
        },
        {
          "application": "cluster_health_info",
          "namespace": "agent",
          "name": "utilization"
        },
        {
          "application": "cluster_health_info",
          "namespace": "node",
          "name": "utilization_aggr_3600"
        },
        {
          "application": "cluster_health_info",
          "namespace": "node",
          "name": "disk_utilization_aggr_3600"
        },
        {
          "application": "cluster_health_info",
          "namespace": "container",
          "name": "utilization"
        },
        {
          "application": "cluster_health_info",
          "namespace": "node",
          "name": "utilization"
        },
        {
          "application": "cluster_health_info",
          "namespace": "file_registry",
          "name": "directory"
        },
        {
          "application": "cluster_health_info",
          "namespace": "node",
          "name": "disk_utilization"
        },
        {
          "application": "cluster_health_info",
          "namespace": "file_registry",
          "name": "file"
        },
        {
          "application": "cluster_health_info",
          "namespace": "container",
          "name": "utilization_aggr_3600"
        },
        {
          "application": "cluster_health_info",
          "namespace": "agent",
          "name": "utilization_aggr_3600"
        },
        {
          "application": "cluster_health_info",
          "namespace": "container",
          "name": "file_usage_aggr_3600"
        },
        {
          "application": "cluster_health_info",
          "namespace": "container",
          "name": "file_usage"
        }
    ]}`
)

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
