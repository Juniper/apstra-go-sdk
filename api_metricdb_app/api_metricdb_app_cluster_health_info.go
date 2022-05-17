package api_metricdb_app

import (
	"time"
)

const (
	CHINamespaceAgent     = "agent"
	CHINamespaceContainer = "container"
	CHINamespaceFileReg   = "file_registry"
	CHINamespaceNode      = "node"

	CHINameHealth    = "health"
	CHINameUtil      = "utilization"
	CHINameFileUsage = "file_usage"
	CHINameDir       = "directory"
	CHINameDisk      = "file"
	CHINameDiskUtil  = "disk_utilization"
)

// AgentHealth is returned by Apstra in response to POSTs
// to ApiUrlMetricdbQuery with metricDbQuery like one of:
//   { "application": "cluster_health_info", "namespace": "agent", "name": "health" },
//   { "application": "cluster_health_info", "namespace": "agent", "name": "health_aggr_3600" },
type AgentHealth struct{} // never seen one of these yet

// AgentUtilization is returned by Apstra in response to POSTs
// to ApiUrlMetricdbQuery with metricDbQuery like one of:
//   { "application": "cluster_health_info", "namespace": "agent", "name": "utilization" },
//   { "application": "cluster_health_info", "namespace": "agent", "name": "utilization_aggr_3600" },
type AgentUtilization struct {
	Node      string    `json:"node"`
	Container string    `json:"container"`
	Timestamp time.Time `json:"timestamp"`
	Agent     string    `json:"agent"`
	Memory    int       `json:"memory"`
	Cpu       int       `json:"cpu"`
}

// ContainerFileUsage is returned by Apstra in response to POSTs
// to ApiUrlMetricdbQuery with metricDbQuery like one of:
//   { "application": "cluster_health_info", "namespace": "container", "name": "file_usage" }
//   { "application": "cluster_health_info", "namespace": "container", "name": "file_usage_aggr_3600" },
type ContainerFileUsage struct {
	Node      string    `json:"node"`
	FileName  string    `json:"file_name"`
	Container string    `json:"container"`
	Timestamp time.Time `json:"timestamp"`
	Size      int       `json:"size"`
}

// ContainerUtilization is returned by Apstra in response to POSTs
// to ApiUrlMetricdbQuery with metricDbQuery like one of:
//   { "application": "cluster_health_info", "namespace": "container", "name": "utilization" },
//   { "application": "cluster_health_info", "namespace": "container", "name": "utilization_aggr_3600" },
type ContainerUtilization struct {
	Node               string    `json:"node"`
	Container          string    `json:"container"`
	Timestamp          time.Time `json:"timestamp"`
	CumulativeFileSize int       `json:"cumulative_file_size"`
	Memory             int       `json:"memory"`
	Cpu                int       `json:"cpu"`
}

// FileRegistryDirectory is returned by Apstra in response to POSTs
// to ApiUrlMetricdbQuery with metricDbQuery like one of:
//   { "application": "cluster_health_info", "namespace": "file_registry", "name": "directory" },
//   { "application": "cluster_health_info", "namespace": "file_registry", "name": "directory_aggr_3600" },
type FileRegistryDirectory struct {
	Node          string    `json:"node"`
	Timestamp     time.Time `json:"timestamp"`
	Size          int       `json:"size"`
	DirectoryPath string    `json:"directory_path"`
}

// FileRegistryFile is returned by Apstra in response to POSTs
// to ApiUrlMetricdbQuery with metricDbQuery like one of
//   { "application": "cluster_health_info", "namespace": "file_registry", "name": "file" },
//   { "application": "cluster_health_info", "namespace": "file_registry", "name": "file_aggr_3600" },
type FileRegistryFile struct {
	Node      string    `json:"node"`
	Timestamp time.Time `json:"timestamp"`
	FilePath  string    `json:"file_path"`
	Size      int       `json:"size"`
}

// NodeDiskUtilization is returned by Apstra in response to POSTs
// to ApiUrlMetricdbQuery with metricDbQuery like one of:
//   { "application": "cluster_health_info", "namespace": "node", "name": "disk_utilization" },
//   { "application": "cluster_health_info", "namespace": "node", "name": "disk_utilization_aggr_3600" },
type NodeDiskUtilization struct {
	Node      string    `json:"node"`
	Timestamp time.Time `json:"timestamp"`
	Partition string    `json:"partition"`
	Size      int       `json:"size"`
}

// NodeUtilization is returned by Apstra in response to POSTs
// to ApiUrlMetricdbQuery with metricDbQuery like one of:
//   { "application": "cluster_health_info", "namespace": "node", "name": "utilization" },
//   { "application": "cluster_health_info", "namespace": "node", "name": "utilization_aggr_3600" },
type NodeUtilization struct {
	Node      string    `json:"node"`
	Timestamp time.Time `json:"timestamp"`
	Cpu       int       `json:"cpu"`
	Memory    int64     `json:"memory"`
}

/*




    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/4eb11184-4b32-4106-8e90-edb312042683", "name": "System Interface Counters" },
    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/4eb11184-4b32-4106-8e90-edb312042683", "name": "Average Interface Counters" },
    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665", "name": "imbalanced_system_count_out_of_range" },
    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665", "name": "leaf_fab_int_tx_avg" },
    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665", "name": "std_dev_percentage" },
    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665", "name": "system_imbalance" },
   cluster_health_info/agent
   agent_metric_path = os.path.join(
       CONSTANTS.metric_application_cluster_health,
       CONSTANTS.metric_namespace_cluster_health_agent)

   cluster_health_info/container
   container_metric_path = os.path.join(
       CONSTANTS.metric_application_cluster_health,
       CONSTANTS.metric_namespace_cluster_health_container)

   cluster_health_info/file_registry
   registry_metric_path = os.path.join(
       CONSTANTS.metric_application_cluster_health,
       CONSTANTS.metric_namespace_cluster_health_file_registry)

   cluster_health_info/node
   node_metric_path = os.path.join(
       CONSTANTS.metric_application_cluster_health,
       CONSTANTS.metric_namespace_cluster_health_node)

    { "application": "cluster_health_info", "namespace": "agent", "name": "health" },
      # empty

    { "application": "cluster_health_info", "namespace": "agent", "name": "health_aggr_3600" },
      # empty

    { "application": "cluster_health_info", "namespace": "agent", "name": "utilization" },
      "node": "AosController",
      "container": "iba4207e3df",
      "timestamp": "2022-05-14T20:10:00.724435Z",
      "agent": "tacspawner",
      "memory": 65826816,
      "cpu": 0

    { "application": "cluster_health_info", "namespace": "agent", "name": "utilization_aggr_3600" },
      "node": "AosController",
      "container": "iba4207e3df",
      "timestamp": "2022-05-14T20:44:00.536622Z",
      "agent": "tacspawner",
      "memory": 65826816,
      "cpu": 0

    { "application": "cluster_health_info", "namespace": "container", "name": "file_usage" }
      "node": "AosController",
      "file_name": "PipelineAgentiba,4207e3df-ab32-4aa0-9cb3-03c4b17ed5c7,iba4207e3df_2022-03-18--22-09-51_55-2022-03-18--22-09-51.674532.tel",
      "container": "iba4207e3df",
      "timestamp": "2022-05-14T20:10:00.724944Z",
      "size": 2483335

    { "application": "cluster_health_info", "namespace": "container", "name": "file_usage_aggr_3600" },
      "node": "AosController",
      "file_name": "PipelineAgentiba,4207e3df-ab32-4aa0-9cb3-03c4b17ed5c7,iba4207e3df_2022-03-18--22-09-51_55-2022-03-18--22-09-51.674532.tel",
      "container": "iba4207e3df",
      "timestamp": "2022-05-14T20:43:59.627640Z",
      "size": 2483655

    { "application": "cluster_health_info", "namespace": "container", "name": "utilization" },
      "node": "AosController",
      "container": "iba4207e3df",
      "timestamp": "2022-05-14T20:10:00.724852Z",
      "cumulative_file_size": 9260744,
      "memory": 268009472,
      "cpu": 0

    { "application": "cluster_health_info", "namespace": "container", "name": "utilization_aggr_3600" },
      "node": "AosController",
      "container": "iba4207e3df",
      "timestamp": "2022-05-14T20:43:59.242497Z",
      "cumulative_file_size": 9262664,
      "memory": 268009472,
      "cpu": 0

    { "application": "cluster_health_info", "namespace": "file_registry", "name": "directory" },
      "node": "AosController",
      "timestamp": "2022-05-14T20:12:07.378641Z",
      "size": 79455693,
      "directory_path": "/var/lib/aos/metricdb_apps/cluster_health_info"

    { "application": "cluster_health_info", "namespace": "file_registry", "name": "directory_aggr_3600" },
      "node": "AosController",
      "timestamp": "2022-05-14T20:44:10.302570Z",
      "size": 79606735,
      "directory_path": "/var/lib/aos/metricdb_apps/cluster_health_info"

    { "application": "cluster_health_info", "namespace": "file_registry", "name": "file" },
      "node": "AosController",
      "timestamp": "2022-05-14T20:12:07.378885Z",
      "file_path": "/var/lib/aos/metricdb_apps/cluster_health_info/node/disk_utilization/disk_utilization-189-2022-05-12--16-45-55.616443.tel",
      "size": 50893

    { "application": "cluster_health_info", "namespace": "file_registry", "name": "file_aggr_3600" },
      "node": "AosController",
      "timestamp": "2022-05-14T20:44:10.273072Z",
      "file_path": "/var/lib/aos/metricdb_apps/cluster_health_info/node/disk_utilization/disk_utilization-189-2022-05-12--16-45-55.616443.tel",
      "size": 50893

    { "application": "cluster_health_info", "namespace": "node", "name": "disk_utilization" },
      "node": "AosController",
      "timestamp": "2022-05-14T20:10:10.721289Z",
      "partition": "aos--server--vg-var+log",
      "size": 1222959104

    { "application": "cluster_health_info", "namespace": "node", "name": "disk_utilization_aggr_3600" },
      "node": "AosController",
      "timestamp": "2022-05-14T20:43:59.242444Z",
      "partition": "aos--server--vg-var+log",
      "size": 1235595264

    { "application": "cluster_health_info", "namespace": "node", "name": "utilization" },
      "node": "AosController",
      "timestamp": "2022-05-14T20:10:10.721232Z",
      "cpu": 1,
      "memory": 13455171584

    { "application": "cluster_health_info", "namespace": "node", "name": "utilization_aggr_3600" },
      "node": "AosController",
      "timestamp": "2022-05-14T20:43:59.242391Z",
      "cpu": 1,
      "memory": 13454567690

    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/4eb11184-4b32-4106-8e90-edb312042683", "name": "System Interface Counters" },
      "timestamp": "2022-05-14T20:10:01.473914Z",
      "aggregate_rx_bps": 2274,
      "max_ifc_rx_utilization": 0,
      "max_ifc_tx_utilization": 0,
      "system_id": "EP226",
      "aggregate_rx_utilization": 0,
      "aggregate_tx_bps": 2133,
      "aggregate_tx_utilization": 0

    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/4eb11184-4b32-4106-8e90-edb312042683", "name": "Average Interface Counters" },
      "symbol_errors_per_second_average": 0,
      "tx_error_pps_average": 0,
      "runts_per_second_average": 0,
      "interface": "xe-0/0/1",
      "rx_utilization_average": 0,
      "speed": 10000000000,
      "tx_discard_pps_average": 0,
      "tx_unicast_pps_average": 0,
      "rx_error_pps_average": 0,
      "fcs_errors_per_second_average": 0,
      "system_id": "WS3119350041",
      "alignment_errors_per_second_average": 0,
      "tx_broadcast_pps_average": 0,
      "timestamp": "2022-05-14T20:10:00.812870Z",
      "tx_multicast_pps_average": 0,
      "rx_broadcast_pps_average": 0,
      "rx_discard_pps_average": 0,
      "giants_per_second_average": 0,
      "tx_bps_average": 88,
      "rx_bps_average": 78,
      "tx_utilization_average": 0,
      "rx_multicast_pps_average": 0,
      "rx_unicast_pps_average": 0

    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665", "name": "imbalanced_system_count_out_of_range" },
      # empty
    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665", "name": "leaf_fab_int_tx_avg" },
      # empty
    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665", "name": "std_dev_percentage" },
      # empty
    { "application": "iba", "namespace": "db10754a-610e-475b-9baa-4c85f82282e8/d6ffca4d-ba91-4833-bf43-714bc0c5b665", "name": "system_imbalance" },
      # empty


*/
