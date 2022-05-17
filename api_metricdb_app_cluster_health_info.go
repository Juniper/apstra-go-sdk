package goapstra

import (
	"context"
	"fmt"
	"time"
)

const (
	MetricdbCHINamespaceAgent     = "agent"
	MetricdbCHINamespaceContainer = "container"
	MetricdbCHINamespaceFileReg   = "file_registry"
	MetricdbCHINamespaceNode      = "node"

	MetricdbCHINameHealth    = "health"
	MetricdbCHINameUtil      = "utilization"
	MetricdbCHINameFileUsage = "file_usage"
	MetricdbCHINameDir       = "directory"
	MetricdbCHINameFile      = "file"
	MetricdbCHINameDiskUtil  = "disk_utilization"
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

//
func QueryMetricdbClusterHealthInfo(ctx context.Context, in *MetricDbQuery) (interface{}, error) {
	_, baseMetricName, _, err := useAggregation(in.Name)
	if err != nil {
		return nil, err
	}
	//	regexp.Match(k)
	//	if strings.HasSuffix(in.Name, MetricdbCHINameAggrSuffix) {
	//		aggr = aggr
	//	}
	switch in.Namespace + baseMetricName {
	case MetricdbCHINamespaceAgent + MetricdbCHINameHealth:
	case MetricdbCHINamespaceAgent + MetricdbCHINameUtil:
	case MetricdbCHINamespaceContainer + MetricdbCHINameFileUsage:
	case MetricdbCHINamespaceContainer + MetricdbCHINameUtil:
	case MetricdbCHINamespaceFileReg + MetricdbCHINameDir:
	case MetricdbCHINamespaceFileReg + MetricdbCHINameFile:
	case MetricdbCHINamespaceNode + MetricdbCHINameDiskUtil:
	case MetricdbCHINamespaceNode + MetricdbCHINameUtil:
	default:
		return nil, fmt.Errorf("unknown metricdb combination: '%s/%s/%s'", in.Application, in.Namespace, in.Name)
	}
	return nil, nil
}
