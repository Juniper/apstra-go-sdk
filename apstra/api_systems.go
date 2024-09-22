// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

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

	systemCommsOn  = "on"
	systemCommsOff = "off"

	ErrAgentNotConnect = iota
)

const ( // new block resets iota to 0
	SystemAdminStateNone = SystemAdminState(iota) // default type 0
	SystemAdminStateNormal
	SystemAdminStateDecomm
	SystemAdminStateMaint
	SystemAdminStateUnknown = "unknown system admin state '%s'"

	systemAdminStateNone    = rawSystemAdminState("")
	systemAdminStateNormal  = rawSystemAdminState("normal")
	systemAdminStateDecomm  = rawSystemAdminState("decomm")
	systemAdminStateMaint   = rawSystemAdminState("maint")
	systemAdminStateUnknown = "unknown system admin state '%d'"
)

type SystemAdminState int

func (o SystemAdminState) Int() int {
	return int(o)
}

func (o SystemAdminState) String() string {
	switch o {
	case SystemAdminStateNone:
		return string(systemAdminStateNone)
	case SystemAdminStateNormal:
		return string(systemAdminStateNormal)
	case SystemAdminStateDecomm:
		return string(systemAdminStateDecomm)
	case SystemAdminStateMaint:
		return string(systemAdminStateMaint)
	default:
		return fmt.Sprintf(systemAdminStateUnknown, o)
	}
}

func (o SystemAdminState) raw() rawSystemAdminState {
	return rawSystemAdminState(o.String())
}

type rawSystemAdminState string

func (o rawSystemAdminState) string() string {
	return string(o)
}

func (o rawSystemAdminState) parse() (int, error) {
	switch o {
	case systemAdminStateNone:
		return int(SystemAdminStateNone), nil
	case systemAdminStateNormal:
		return int(SystemAdminStateNormal), nil
	case systemAdminStateDecomm:
		return int(SystemAdminStateDecomm), nil
	case systemAdminStateMaint:
		return int(SystemAdminStateMaint), nil
	default:
		return 0, fmt.Errorf(SystemAdminStateUnknown, o)
	}
}

type SystemId string

type optionsSystemsResponse struct {
	Items   []SystemId `json:"items"`
	Methods []string   `json:"methods"`
}

type ManagedSystemInfo struct {
	ContainerStatus SystemContainerStatus
	DeviceKey       string
	Facts           SystemFacts
	Id              SystemId
	Services        []string
	Status          SystemStatus
	UserConfig      SystemUserConfig
}

type rawManagedSystemInfo struct {
	ContainerStatus SystemContainerStatus `json:"container_status"`
	DeviceKey       string                `json:"device_key"`
	Facts           SystemFacts           `json:"facts"`
	Id              SystemId              `json:"id"`
	Services        []string              `json:"services"`
	Status          rawSystemStatus       `json:"status"`
	UserConfig      rawSystemUserConfig   `json:"user_config"`
}

func (o *rawManagedSystemInfo) polish() (*ManagedSystemInfo, error) {
	userConfig, err := o.UserConfig.polish()
	if err != nil {
		return nil, err
	}

	status, err := o.Status.polish()
	if err != nil {
		return nil, err
	}

	return &ManagedSystemInfo{
		ContainerStatus: o.ContainerStatus,
		DeviceKey:       o.DeviceKey,
		Facts:           o.Facts,
		Id:              o.Id,
		Services:        o.Services,
		Status:          *status,
		UserConfig:      *userConfig,
	}, nil
}

type SystemContainerStatus struct {
	Error  string `json:"error"`
	Host   string `json:"host"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SystemFacts struct {
	AosHclModel      ObjectId `json:"aos_hcl_model"`
	AosServer        string   `json:"aos_server"`
	AosVersion       string   `json:"aos_version"`
	ChassisMacRanges string   `json:"chassis_mac_ranges"`
	HwModel          string   `json:"hw_model"`
	HwVersion        string   `json:"hw_version"`
	MgmtIfname       string   `json:"mgmt_ifname"`
	MgmtIpaddr       string   `json:"mgmt_ipaddr"`
	MgmtMacaddr      string   `json:"mgmt_macaddr"`
	OsArch           string   `json:"os_arch"`
	OsFamily         string   `json:"os_family"`
	OsVersion        string   `json:"os_version"`
	OsVersionInfo    struct {
		Build string `json:"build"`
		Major string `json:"major"`
		Minor string `json:"minor"`
	} `json:"os_version_info"`
	SerialNumber string `json:"serial_number"`
	Vendor       string `json:"vendor"`
}

type SystemStatus struct {
	AgentStartTime  time.Time             `json:"agent_start_time"`
	CommState       string                `json:"comm_state"`
	DeviceStartTime time.Time             `json:"device_start_time"`
	ErrorMessage    string                `json:"error_message"`
	IsAcknowledged  bool                  `json:"is_acknowledged"`
	OperationMode   SystemManagementLevel `json:"operation_mode"`
	State           string                `json:"state"`
}

type rawSystemStatus struct {
	AgentStartTime  time.Time             `json:"agent_start_time"`
	CommState       string                `json:"comm_state"`
	DeviceStartTime time.Time             `json:"device_start_time"`
	ErrorMessage    string                `json:"error_message"`
	IsAcknowledged  bool                  `json:"is_acknowledged"`
	OperationMode   systemManagementLevel `json:"operation_mode"`
	State           string                `json:"state"`
	UserConfig      SystemUserConfig      `json:"user_config"`
}

func (o *rawSystemStatus) polish() (*SystemStatus, error) {
	operationMode, err := o.OperationMode.parse()
	if err != nil {
		return nil, err
	}

	return &SystemStatus{
		AgentStartTime:  o.AgentStartTime,
		CommState:       o.CommState,
		DeviceStartTime: o.DeviceStartTime,
		ErrorMessage:    o.ErrorMessage,
		IsAcknowledged:  o.IsAcknowledged,
		OperationMode:   SystemManagementLevel(operationMode),
		State:           o.State,
	}, nil
}

type systemUpdate struct {
	UserConfig rawSystemUserConfig `json:"user_config"`
}

type SystemUserConfig struct {
	AdminState  SystemAdminState `json:"admin_state,omitempty"`
	AosHclModel ObjectId         `json:"aos_hcl_model,omitempty"`
	Location    string           `json:"location,omitempty"`
}

func (o *SystemUserConfig) raw() *rawSystemUserConfig {
	return &rawSystemUserConfig{
		AdminState:  o.AdminState.raw(),
		AosHclModel: o.AosHclModel,
		Location:    o.Location,
	}
}

type rawSystemUserConfig struct {
	AdminState  rawSystemAdminState `json:"admin_state,omitempty"`
	AosHclModel ObjectId            `json:"aos_hcl_model,omitempty"`
	Location    string              `json:"location,omitempty"`
}

func (o *rawSystemUserConfig) polish() (*SystemUserConfig, error) {
	adminState, err := o.AdminState.parse()
	if err != nil {
		return nil, err
	}

	return &SystemUserConfig{
		AdminState:  SystemAdminState(adminState),
		AosHclModel: o.AosHclModel,
		Location:    o.Location,
	}, nil
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

func (o *Client) getSystemInfo(ctx context.Context, id SystemId) (*rawManagedSystemInfo, error) {
	response := &rawManagedSystemInfo{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlSystemsById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *Client) getAllSystemsInfo(ctx context.Context) ([]rawManagedSystemInfo, error) {
	response := &struct{ Items []rawManagedSystemInfo }{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlSystems,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *Client) updateSystemByAgentId(ctx context.Context, agentId ObjectId, cfg *SystemUserConfig) error {
	agent, err := o.getSystemAgent(ctx, agentId)
	if err != nil {
		return fmt.Errorf("cannot get info for agent '%s' - %w", agentId, err)
	}

	if agent.Status.SystemId == "" {
		return fmt.Errorf("cannot acknowledge system from agent '%s' - system ID is empty", agentId)
	}

	return o.updateSystem(ctx, agent.Status.SystemId, cfg)
}

func (o *Client) updateSystem(ctx context.Context, id SystemId, cfg *SystemUserConfig) error {
	method := http.MethodPut
	urlStr := fmt.Sprintf(apiUrlSystemsById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}

	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   method,
		url:      apstraUrl,
		apiInput: &systemUpdate{UserConfig: *cfg.raw()},
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
