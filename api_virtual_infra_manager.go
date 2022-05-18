package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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
	apstraUrl, err := url.Parse(apiUrlVirtualInfraManagers)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVirtualInfraManagers, err)
	}
	response := &virtualInfraMgrsResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}

	return response.Items, err
}
