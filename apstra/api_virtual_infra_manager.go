package goapstra

import (
	"context"
	"net/http"
	"time"
)

const (
	apiUrlVirtualInfraManagers = "/api/virtual-infra-managers"
)

type virtualInfraMgrsResponse struct {
	Items []VirtualInfraMgrInfo
}

type VirtualInfraMgrInfo struct {
	ConnectionState              string    `json:"connection_state"`
	LastSuccessfulCollectionTime time.Time `json:"last_successful_collection_time"`
	ServiceEnabled               bool      `json:"service_enabled"`
	ManagementIp                 string    `json:"management_ip"`
	SystemId                     string    `json:"system_id"`
	AgentId                      string    `json:"agent_id"`
	VirtualInfraType             string    `json:"virtual_infra_type"`
}

func (o *Client) getVirtualInfraMgrs(ctx context.Context) ([]VirtualInfraMgrInfo, error) {
	response := &virtualInfraMgrsResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlVirtualInfraManagers,
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}

	return response.Items, err
}
