package goapstra

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	apiUrlSystemAgents         = "/api/system-agents"
	apiUrlSystemAgentsPrefix   = "/api/system-agents" + apiUrlPathDelim
	apiUrlSystemAgentsById     = apiUrlSystemAgentsPrefix + "%s"
	apiUrlSystemAgentInstall   = apiUrlSystemAgentsPrefix + "%s" + "/install-agent"
	apiUrlSystemAgentUninstall = apiUrlSystemAgentsPrefix + "%s" + "/uninstall-agent"

	SystemAgentJobInstall             = "install"
	SystemAgentOperationModeFull      = "full_control"
	SystemAgentOperationModeTelemetry = "telemetry_only"
	SystemAgentTypeOffbox             = "offbox"
	SystemAgentTypeOnbox              = "onbox"
	systemAgentConnectStateConnected  = "connected"
)

type optionsSystemAgentsResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type getSystemAgentsResponse struct {
	Items []SystemAgentInfo `json:"items"`
}

type SystemAgentInfo struct {
	Status struct {
		ConnectionState   string   `json:"connection_state"`
		PackagesInstalled []string `json:"packages_installed"`
		CurrentTask       string   `json:"current_task"`
		PendingType       string   `json:"pending_type"`
		JobId             int      `json:"job_id"`
		OperationMode     string   `json:"operation_mode"`
		PlatformVersion   string   `json:"platform_version"`
		Platform          string   `json:"platform"`
		State             string   `json:"state"`
		SystemId          ObjectId `json:"system_id"`
		HasCredential     bool     `json:"has_credential"`
		Error             string   `json:"error"`
		StatusMessage     string   `json:"status_message"`
		AosVersion        string   `json:"aos_version"`
	} `json:"status"`
	PlatformConfig struct {
		ContainerEnable bool `json:"container_enable"`
	} `json:"platform_config"`
	RunningConfig ConfigInfo `json:"running_config"`
	LastJobStatus struct {
		Started        time.Time `json:"started"`
		JobType        string    `json:"job_type"`
		Finished       time.Time `json:"finished"`
		HostId         string    `json:"host_id"`
		CurrentTask    string    `json:"current_task"`
		IsLogAvailable bool      `json:"is_log_available"`
		JobId          int       `json:"job_id"`
		Created        time.Time `json:"created"`
		State          string    `json:"state"`
		AgentType      string    `json:"agent_type"`
		Error          string    `json:"error"`
	} `json:"last_job_status"`
	DeviceFacts struct {
		DeviceOsFamily  string `json:"device_os_family"`
		Hostname        string `json:"hostname"`
		DeviceState     string `json:"device_state"`
		DeviceOsVersion string `json:"device_os_version"`
	} `json:"device_facts"`
	TelemetryExtStatus struct {
		PackagesInstalled []string `json:"packages_installed"`
		StatusMessage     string   `json:"status_message"`
	} `json:"telemetry_ext_status"`
	Config          ConfigInfo `json:"config"`
	Id              ObjectId   `json:"id"`
	ContainerStatus struct {
		Status     string `json:"status"`
		Name       string `json:"name"`
		TaskId     string `json:"task_id"`
		LastUpdate string `json:"last_update"`
		Host       string `json:"host"`
		Error      string `json:"error"`
		ServiceId  string `json:"service_id"`
	} `json:"container_status"`
	PlatformStatus struct {
		JobId           int    `json:"job_id"`
		PlatformVersion string `json:"platform_version"`
		Platform        string `json:"platform"`
		State           string `json:"state"`
		HasCredential   bool   `json:"has_credential"`
		Error           string `json:"error"`
		CurrentTask     string `json:"current_task"`
	} `json:"platform_status"`
}

type ConfigInfo struct {
	Profile             string   `json:"profile"`
	ForcePackageInstall bool     `json:"force_package_install"`
	InstallRequirements bool     `json:"install_requirements"`
	Packages            []string `json:"packages"`
	OpenOptions         struct {
		AdditionalProp1 string `json:"additionalProp1"`
		AdditionalProp2 string `json:"additionalProp2"`
		AdditionalProp3 string `json:"additionalProp3"`
	} `json:"open_options"`
	Label           string   `json:"label"`
	Platform        string   `json:"platform"`
	ManagementIp    string   `json:"management_ip"`
	AgentType       string   `json:"agent_type"`
	AllowedJobTypes []string `json:"allowed_job_types"`
	OperationMode   string   `json:"operation_mode"`
	Id              ObjectId `json:"id"`
}

type SystemAgentCfg struct {
	AgentType           string   `json:"agent_type,omitempty"`
	ManagementIp        string   `json:"management_ip,omitempty"`
	Profile             string   `json:"profile,omitempty"`
	OperationMode       string   `json:"operation_mode,omitempty"`
	JobOnCreate         string   `json:"job_on_create,omitempty"`
	Username            string   `json:"username,omitempty"`
	ForcePackageInstall bool     `json:"force_package_install,omitempty"`
	InstallRequirements bool     `json:"install_requirements,omitempty"`
	EnableMonitor       bool     `json:"enable_monitor,omitempty"`
	Id                  string   `json:"id,omitempty"`
	Password            string   `json:"password,omitempty"`
	Packages            []string `json:"packages,omitempty,omitempty"`
	Label               string   `json:"label,omitempty"`
	Platform            string   `json:"platform,omitempty"`
}

func (o *Client) listSystemAgents(ctx context.Context) ([]ObjectId, error) {
	method := http.MethodOptions
	urlStr := apiUrlSystemAgents
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &optionsSystemAgentsResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' at '%s'", http.MethodOptions, apstraUrl.String())
	}
	return response.Items, nil
}

func (o *Client) getSystemAgent(ctx context.Context, id ObjectId) (*SystemAgentInfo, error) {
	method := http.MethodGet
	urlStr := fmt.Sprintf(apiUrlSystemAgentsById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &SystemAgentInfo{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		var ttae TalkToApstraErr
		if errors.As(err, &ttae) {
			if ttae.Response.StatusCode == http.StatusNotFound {
				return nil, ApstraClientErr{
					errType: ErrNotfound,
					err:     err,
				}
			}
		}
		return nil, err
	}
	return response, nil
}

func (o *Client) getAllSystemAgents(ctx context.Context) ([]SystemAgentInfo, error) {
	method := http.MethodGet
	urlStr := apiUrlSystemAgents
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &getSystemAgentsResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' at '%s'", method, apstraUrl.String())
	}
	return response.Items, nil
}

func (o *Client) getSystemAgentByManagementIp(ctx context.Context, ip string) (*SystemAgentInfo, error) {
	asa, err := o.getAllSystemAgents(ctx)
	if err != nil {
		return nil, err
	}
	for _, a := range asa {
		if a.Config.ManagementIp == ip || a.RunningConfig.ManagementIp == ip { // what's the difference?
			return &a, nil
		}
	}
	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no System Agent with management IP '%s' found", ip),
	}
}

func (o *Client) createSystemAgent(ctx context.Context, request *SystemAgentCfg) (ObjectId, error) {
	method := http.MethodPost
	urlStr := apiUrlSystemAgents
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &objectIdResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiInput:    request,
		apiResponse: response,
	})
	if err != nil {
		var ttae TalkToApstraErr
		if errors.As(err, &ttae) {
			if ttae.Response.StatusCode == http.StatusConflict {
				buf := make([]byte, 512)
				_, _ = io.ReadFull(ttae.Response.Body, buf)
				return "", ApstraClientErr{
					errType: ErrExists,
					err:     fmt.Errorf("%w - %s", ttae, string(buf)),
				}
			}
		}
		return "", err
	}

	return response.Id, nil
}

func (o *Client) updateSystemAgent(ctx context.Context, id ObjectId, request *SystemAgentCfg) error {
	method := http.MethodPatch
	urlStr := fmt.Sprintf(apiUrlSystemAgentsById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &objectIdResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiInput:    request,
		apiResponse: response,
	})
	if err != nil {
		var ttae TalkToApstraErr
		if errors.As(err, &ttae) {
			if ttae.Response.StatusCode == http.StatusNotFound {
				return ApstraClientErr{
					errType: ErrNotfound,
					err:     err,
				}
			}
		}
		return err
	}
	return nil
}

func (o *Client) deleteSystemAgent(ctx context.Context, id ObjectId) error {
	method := http.MethodDelete
	urlStr := fmt.Sprintf(apiUrlSystemAgentsById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method: method,
		url:    apstraUrl,
	})
	if err != nil {
		var ttae TalkToApstraErr
		if errors.As(err, &ttae) {
			if ttae.Response.StatusCode == http.StatusNotFound {
				return ApstraClientErr{
					errType: ErrNotfound,
					err:     err,
				}
			}
		}
		return err
	}
	return nil
}

func (o *Client) uninstallSystemAgent(ctx context.Context, id ObjectId) error {
	method := http.MethodPost
	urlStr := fmt.Sprintf(apiUrlSystemAgentUninstall, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method: method,
		url:    apstraUrl,
	})
	if err != nil {
		var ttae TalkToApstraErr
		if errors.As(err, &ttae) {
			if ttae.Response.StatusCode == http.StatusNotFound {
				return ApstraClientErr{
					errType: ErrNotfound,
					err:     err,
				}
			}
		}
		return err
	}
	return nil
}
