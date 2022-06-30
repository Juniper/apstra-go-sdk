package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	apiUrlSystems       = "/api/systems"
	apiUrlSystemsPrefix = "/api/systems" + apiUrlPathDelim
	apiUrlSystemsById   = apiUrlSystemsPrefix + "%s"

	systemAdminStateNormal = "normal"
	systemAdminStateDecomm = "decomm"

	systemCommsOn  = "on"
	systemCommsOff = "off"

	ErrAgentNotConnect = iota
)

type SystemId string

type optionsSystemsResponse struct {
	Items   []SystemId `json:"items"`
	Methods []string   `json:"methods"`
}

type ManagedSystemInfo struct {
	ContainerStatus SystemContainerStatus `json:"container_status"`
	DeviceKey       SystemId              `json:"device_key"`
	Facts           SystemFacts           `json:"facts"`
	Id              SystemId              `json:"id"`
	Services        []string              `json:"services"`
	Status          SystemStatus          `json:"status"`
	UserConfig      SystemUserConfig      `json:"user_config"`
}

type rawManagedSystemInfo struct {
	ContainerStatus SystemContainerStatus `json:"container_status"`
	DeviceKey       SystemId              `json:"device_key"`
	Facts           SystemFacts           `json:"facts"`
	Id              SystemId              `json:"id"`
	Services        []string              `json:"services"`
	Status          rawSystemStatus       `json:"status"`
	UserConfig      SystemUserConfig      `json:"user_config"`
}

func (o *rawManagedSystemInfo) polish() *ManagedSystemInfo {
	return &ManagedSystemInfo{
		ContainerStatus: o.ContainerStatus,
		DeviceKey:       o.DeviceKey,
		Facts:           o.Facts,
		Id:              o.Id,
		Services:        o.Services,
		Status:          *o.Status.polish(),
		UserConfig:      o.UserConfig,
	}
}

type SystemContainerStatus struct {
	Error  string `json:"error"`
	Host   string `json:"host"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SystemFacts struct {
	AosHclModel      string `json:"aos_hcl_model"`
	AosServer        string `json:"aos_server"`
	AosVersion       string `json:"aos_version"`
	ChassisMacRanges string `json:"chassis_mac_ranges"`
	HwModel          string `json:"hw_model"`
	HwVersion        string `json:"hw_version"`
	MgmtIfname       string `json:"mgmt_ifname"`
	MgmtIpaddr       string `json:"mgmt_ipaddr"`
	MgmtMacaddr      string `json:"mgmt_macaddr"`
	OsArch           string `json:"os_arch"`
	OsFamily         string `json:"os_family"`
	OsVersion        string `json:"os_version"`
	OsVersionInfo    struct {
		Build string `json:"build"`
		Major string `json:"major"`
		Minor string `json:"minor"`
	} `json:"os_version_info"`
	SerialNumber string `json:"serial_number"`
	Vendor       string `json:"vendor"`
}

type SystemStatus struct {
	AgentStartTime  time.Time `json:"agent_start_time"`
	CommState       string    `json:"comm_state"`
	DeviceStartTime time.Time `json:"device_start_time"`
	ErrorMessage    string    `json:"error_message"`
	IsAcknowledged  bool      `json:"is_acknowledged"`
	OperationMode   AgentMode `json:"operation_mode"`
	State           string    `json:"state"`
}

type rawSystemStatus struct {
	AgentStartTime  time.Time        `json:"agent_start_time"`
	CommState       string           `json:"comm_state"`
	DeviceStartTime time.Time        `json:"device_start_time"`
	ErrorMessage    string           `json:"error_message"`
	IsAcknowledged  bool             `json:"is_acknowledged"`
	OperationMode   rawAgentMode     `json:"operation_mode"`
	State           string           `json:"state"`
	UserConfig      SystemUserConfig `json:"user_config"`
}

func (o *rawSystemStatus) polish() *SystemStatus {
	return &SystemStatus{
		AgentStartTime:  o.AgentStartTime,
		CommState:       o.CommState,
		DeviceStartTime: o.DeviceStartTime,
		ErrorMessage:    o.ErrorMessage,
		IsAcknowledged:  o.IsAcknowledged,
		OperationMode:   AgentMode(o.OperationMode.parse()),
		State:           o.State,
	}
}

type systemUpdate struct {
	UserConfig SystemUserConfig `json:"user_config"`
}

type SystemUserConfig struct {
	AdminState  string `json:"admin_state,omitempty"`
	AosHclModel string `json:"aos_hcl_model,omitempty"`
	Location    string `json:"location,omitempty"`
}

func (o *Client) listSystems(ctx context.Context) ([]SystemId, error) {
	method := http.MethodOptions
	urlStr := apiUrlSystems
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &optionsSystemsResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error listing systems - %w", err)
	}

	return response.Items, nil
}

func (o *Client) getSystemInfo(ctx context.Context, id SystemId) (*ManagedSystemInfo, error) {
	method := http.MethodGet
	urlStr := fmt.Sprintf(apiUrlSystemsById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &rawManagedSystemInfo{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}
	return response.polish(), nil
}

func (o *Client) acknowledgeSystemByAgentId(ctx context.Context, agentId ObjectId, location string) error {
	agent, err := o.getAgentInfo(ctx, agentId)
	if err != nil {
		return fmt.Errorf("cannot get info for agent '%s' - %w", agentId, err)
	}

	// todo: maybe this test isn't needed?
	if agent.Status.ConnectionState != AgentCxnStateConnected {
		return ApstraClientErr{
			errType: ErrAgentNotConnect,
			err: fmt.Errorf("cannot acknowledge system with connection state '%s', must be '%s'",
				agent.Status.ConnectionState.String(), AgentCxnStateConnected.String()),
		}
	}

	if agent.Status.SystemId == "" {
		return fmt.Errorf("cannot acknowledge system from agent '%s' - system ID is empty", agentId)
	}

	return o.acknowledgeSystem(ctx, agent.Status.SystemId, location)
}

func (o *Client) acknowledgeSystem(ctx context.Context, id SystemId, location string) error {
	systemInfo, err := o.getSystemInfo(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting system info - %w", err)
	}

	systemInfo.UserConfig.AdminState = systemAdminStateNormal

	return o.updateSystem(ctx, id, &SystemUserConfig{
		AdminState:  systemAdminStateNormal,
		AosHclModel: systemInfo.Facts.AosHclModel,
		Location:    location,
	})
}

func (o *Client) updateSystem(ctx context.Context, id SystemId, cfg *SystemUserConfig) error {
	method := http.MethodPut
	urlStr := fmt.Sprintf(apiUrlSystemsById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}

	update := &systemUpdate{UserConfig: *cfg}
	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   method,
		url:      apstraUrl,
		apiInput: update,
	})
}

func (o *Client) deleteSystem(ctx context.Context, id SystemId) error {
	method := http.MethodDelete
	urlStr := fmt.Sprintf(apiUrlSystemsById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}

	return o.talkToApstra(ctx, &talkToApstraIn{
		method: method,
		url:    apstraUrl,
	})
}
