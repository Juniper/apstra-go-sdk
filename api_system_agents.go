package goapstra

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	apiUrlSystemAgents          = "/api/system-agents"
	apiUrlSystemAgentsPrefix    = "/api/system-agents" + apiUrlPathDelim
	apiUrlSystemAgentsById      = apiUrlSystemAgentsPrefix + "%s"
	apiUrlSystemAgentCheck      = apiUrlSystemAgentsPrefix + "%s" + "/check"
	apiUrlSystemAgentInstall    = apiUrlSystemAgentsPrefix + "%s" + "/install-agent"
	apiUrlSystemAgentUninstall  = apiUrlSystemAgentsPrefix + "%s" + "/uninstall-agent"
	apiUrlSystemAgentJobHistory = apiUrlSystemAgentsPrefix + "%s" + "/job-history"
)

const ( // new block resets iota to 0
	AgentCxnStateNone = AgentCxnState(iota)
	AgentCxnStateConnected
	AgentCxnStateDisconnected
	AgentCxnStateAuthFail
	AgentCxnStateUnknown

	agentCxnStateNone         = rawAgentCxnState("")
	agentCxnStateConnected    = rawAgentCxnState("connected")
	agentCxnStateDisconnected = rawAgentCxnState("disconnected")
	agentCxnStateAuthFail     = rawAgentCxnState("auth_failed")
	agentCxnStateUnknown      = "system agent connection state %d unknown"
)

const ( // new block resets iota to 0
	AgentTypeDefault = AgentType(iota) // default type 0
	AgentTypeOffbox
	AgentTypeOnbox
	AgentTypeUnknown

	agentTypeDefault = agentTypeOffbox
	agentTypeOffbox  = rawAgentType("offbox")
	agentTypeOnbox   = rawAgentType("onbox")
	agentTypeUnknown = "system agent type %d unknown"
)

const ( // new block resets iota to 0
	AgentModeFull = AgentMode(iota) // default type 0
	AgentModeTelemetry
	AgentModeUnknown

	agentModeFull      = rawAgentMode("full_control")
	agentModeTelemetry = rawAgentMode("telemetry_only")
	agentModeUnknown   = "system agent mode %d unknown"
)

const ( // new block resets iota to 0
	AgentJobTypeNull = AgentJobType(iota) // default type 0
	AgentJobTypeNone
	AgentJobTypeInstall
	AgentJobTypeCheck
	AgentJobTypeRevertToPristine
	AgentJobTypeUninstall
	AgentJobTypeUpgrade
	AgentJobTypeUnknown

	agentJobTypeNull             = rawAgentJobType("")
	agentJobTypeNone             = rawAgentJobType("none")
	agentJobTypeInstall          = rawAgentJobType("install")
	agentJobTypeCheck            = rawAgentJobType("check")
	agentJobTypeRevertToPristine = rawAgentJobType("revertToPristine")
	agentJobTypeUninstall        = rawAgentJobType("uninstall")
	agentJobTypeUpgrade          = rawAgentJobType("upgrade")
	agentJobTypeUnknown          = "system agent job type %d unknown"
)

const ( // new block resets iota to 0
	AgentJobStateNull = AgentJobState(iota)
	AgentJobStateInit
	AgentJobStateInProgress
	AgentJobStateSuccess
	AgentJobStateUnknown

	agentJobStateNull       = rawAgentJobState("")
	agentJobStateInit       = rawAgentJobState("init")
	agentJobStateInProgress = rawAgentJobState("inprogress")
	agentJobStateSuccess    = rawAgentJobState("success")
	agentJobStateUnknown    = "system agent job state %d unknown"
)

const ( // new block resets iota to 0
	AgentPlatformNull = AgentPlatform(iota) // default type 0
	AgentPlatformJunos
	AgentPlatformEOS
	AgentPlatformNXOS
	AgentPlatformUnknown

	agentPlatformNull    = rawAgentPlatform("")
	agentPlatformJunos   = rawAgentPlatform("junos")
	agentPlatformEOS     = rawAgentPlatform("eos")
	agentPlatformNXOS    = rawAgentPlatform("nxos")
	agentPlatformUnknown = "system agent platform %d unknown"
)

type JobId int

type AgentType int

func (o AgentType) Int() int {
	return int(o)
}

func (o AgentType) String() string {
	switch o {
	case AgentTypeDefault:
		return string(agentTypeDefault)
	case AgentTypeOffbox:
		return string(agentTypeOffbox)
	case AgentTypeOnbox:
		return string(agentTypeOnbox)
	default:
		return fmt.Sprintf(agentTypeUnknown, o)
	}
}

func (o AgentType) raw() rawAgentType {
	return rawAgentType(o.String())
}

type rawAgentType string

func (o rawAgentType) string() string {
	return string(o)
}

func (o rawAgentType) parse() int {
	switch o {
	case "":
		return int(AgentTypeDefault)
	case agentTypeOffbox:
		return int(AgentTypeOffbox)
	case agentTypeOnbox:
		return int(AgentTypeOnbox)
	default:
		return int(AgentTypeUnknown)
	}
}

type AgentMode int

func (o AgentMode) Int() int {
	return int(o)
}

func (o AgentMode) String() string {
	switch o {
	case AgentModeFull:
		return string(agentModeFull)
	case AgentModeTelemetry:
		return string(agentModeTelemetry)
	default:
		return fmt.Sprintf(agentModeUnknown, o)
	}
}

func (o AgentMode) raw() rawAgentMode {
	return rawAgentMode(o.String())
}

type rawAgentMode string

func (o rawAgentMode) string() string {
	return string(o)
}

func (o rawAgentMode) parse() int {
	switch o {
	case agentModeFull:
		return int(AgentModeFull)
	case agentModeTelemetry:
		return int(AgentModeTelemetry)
	default:
		return int(AgentModeUnknown)
	}
}

type AgentJobState int

func (o AgentJobState) Int() int {
	return int(o)
}

func (o AgentJobState) String() string {
	switch o {
	case AgentJobStateNull:
		return string(agentJobStateNull)
	case AgentJobStateInit:
		return string(agentJobStateInit)
	case AgentJobStateInProgress:
		return string(agentJobStateInProgress)
	case AgentJobStateSuccess:
		return string(agentJobStateSuccess)
	default:
		return fmt.Sprintf(agentJobStateUnknown, o)
	}
}

func (o AgentJobState) raw() rawAgentJobState {
	return rawAgentJobState(o.String())
}

func (o AgentJobState) HasExited() bool {
	switch o {
	// todo: more states which look like "exited" ?
	case AgentJobStateSuccess:
		return true
	}
	return false
}

type rawAgentJobState string

func (o rawAgentJobState) string() string {
	return string(o)
}

func (o rawAgentJobState) parse() int {
	switch o {
	case agentJobStateNull:
		return int(AgentJobStateNull)
	case agentJobStateInit:
		return int(AgentJobStateInit)
	case agentJobStateInProgress:
		return int(AgentJobStateInProgress)
	case agentJobStateSuccess:
		return int(AgentJobStateSuccess)
	default:
		return int(AgentJobStateUnknown)
	}

}

type AgentJobType int

func (o AgentJobType) Int() int {
	return int(o)
}

func (o AgentJobType) String() string {
	switch o {
	case AgentJobTypeNull:
		return string(agentJobTypeNull)
	case AgentJobTypeInstall:
		return string(agentJobTypeInstall)
	case AgentJobTypeCheck:
		return string(agentJobTypeCheck)
	case AgentJobTypeRevertToPristine:
		return string(agentJobTypeRevertToPristine)
	case AgentJobTypeUninstall:
		return string(agentJobTypeUninstall)
	case AgentJobTypeUpgrade:
		return string(agentJobTypeUpgrade)
	case AgentJobTypeNone:
		return string(agentJobTypeNone)
	default:
		return fmt.Sprintf(agentJobTypeUnknown, o)
	}
}

func (o AgentJobType) raw() rawAgentJobType {
	return rawAgentJobType(o.String())
}

type rawAgentJobType string

func (o rawAgentJobType) string() string {
	return string(o)
}

func (o rawAgentJobType) parse() int {
	switch o {
	case agentJobTypeNull:
		return int(AgentJobTypeNull)
	case agentJobTypeInstall:
		return int(AgentJobTypeInstall)
	case agentJobTypeCheck:
		return int(AgentJobTypeCheck)
	case agentJobTypeRevertToPristine:
		return int(AgentJobTypeRevertToPristine)
	case agentJobTypeUninstall:
		return int(AgentJobTypeUninstall)
	case agentJobTypeUpgrade:
		return int(AgentJobTypeUpgrade)
	case agentJobTypeNone:
		return int(AgentJobTypeNone)
	default:
		return int(AgentJobTypeUnknown)
	}

}

type AgentPlatform int

func (o AgentPlatform) Int() int {
	return int(o)
}

func (o AgentPlatform) String() string {
	switch o {
	case AgentPlatformNull:
		return string(agentPlatformNull)
	case AgentPlatformJunos:
		return string(agentPlatformJunos)
	case AgentPlatformEOS:
		return string(agentPlatformEOS)
	case AgentPlatformNXOS:
		return string(agentPlatformNXOS)
	default:
		return fmt.Sprintf(agentPlatformUnknown, o)
	}
}

func (o AgentPlatform) raw() rawAgentPlatform {
	return rawAgentPlatform(o.String())
}

type rawAgentPlatform string

func (o rawAgentPlatform) string() string {
	return string(o)
}

func (o rawAgentPlatform) parse() int {
	switch o {
	case agentPlatformNull:
		return int(AgentPlatformNull)
	case agentPlatformEOS:
		return int(AgentPlatformEOS)
	case agentPlatformJunos:
		return int(AgentPlatformJunos)
	case agentPlatformNXOS:
		return int(AgentPlatformNXOS)
	default:
		return int(AgentPlatformUnknown)
	}
}

type AgentCxnState int

func (o AgentCxnState) Int() int {
	return int(o)
}

func (o AgentCxnState) String() string {
	switch o {
	case AgentCxnStateNone:
		return string(agentCxnStateNone)
	case AgentCxnStateConnected:
		return string(agentCxnStateConnected)
	case AgentCxnStateDisconnected:
		return string(agentCxnStateDisconnected)
	case AgentCxnStateAuthFail:
		return string(agentCxnStateAuthFail)
	default:
		return fmt.Sprintf(agentCxnStateUnknown, o)
	}

}

func (o AgentCxnState) raw() rawAgentCxnState {
	return rawAgentCxnState(o.String())
}

type rawAgentCxnState string

func (o rawAgentCxnState) string() string {
	return string(o)
}

func (o rawAgentCxnState) parse() int {
	switch o {
	case agentCxnStateNone:
		return int(AgentCxnStateNone)
	case agentCxnStateConnected:
		return int(AgentCxnStateConnected)
	case agentCxnStateDisconnected:
		return int(AgentCxnStateDisconnected)
	case agentCxnStateAuthFail:
		return int(AgentCxnStateAuthFail)
	default:
		return int(AgentCxnStateUnknown)
	}
}

type AgentPackages map[string]string

// todo: method for inventory
//   to be checked by agent profile and agent create/update methods,
//   refuse to reference packages which do not exist.

func (o *AgentPackages) raw() rawAgentPackages {
	// todo: one of these lines causes 'null' in JSON output, while the other causes '[]' ... which one is correct for the API?
	//raw := rawAgentPackages{}
	var raw rawAgentPackages
	for k, v := range *o {
		raw = append(raw, k+apstraSystemAgentPlatformStringSep+v)
	}
	return raw
}

type rawAgentPackages []string

func (o rawAgentPackages) polish() AgentPackages {
	var polish AgentPackages
	if len(o) > 0 {
		polish = make(map[string]string)
	}
	for _, s := range o {
		kv := strings.SplitN(s, apstraSystemAgentPlatformStringSep, 2)
		switch len(kv) {
		case 2:
			polish[kv[0]] = kv[1]
		case 1:
			polish[kv[0]] = ""
		}
	}
	return polish
}

type rawAgentJobTypes []rawAgentJobType

func (o rawAgentJobTypes) polish() []AgentJobType {
	result := make([]AgentJobType, len(o))
	for i, t := range o {
		result[i] = AgentJobType(t.parse())
	}
	return result
}

type agentJobHistoryResponse struct {
	Items []rawAgentJobStatus `json:"items"`
}

type optionsAgentsResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type getAgentsResponse struct {
	Items []rawSystemAgent `json:"items"`
}

type AgentStatus struct {
	ConnectionState   AgentCxnState
	PackagesInstalled AgentPackages
	CurrentTask       string
	PendingType       string
	JobId             JobId
	OperationMode     AgentMode
	PlatformVersion   string
	Platform          AgentPlatform
	State             string
	SystemId          SystemId
	HasCredential     bool
	Error             string
	StatusMessage     string
	AosVersion        string
}

func (o *AgentStatus) raw() *rawAgentStatus {
	return &rawAgentStatus{
		ConnectionState:   rawAgentCxnState(o.ConnectionState.String()),
		PackagesInstalled: o.PackagesInstalled.raw(),
		CurrentTask:       o.CurrentTask,
		PendingType:       o.PendingType,
		JobId:             o.JobId,
		OperationMode:     rawAgentMode(o.OperationMode.String()),
		PlatformVersion:   o.PlatformVersion,
		Platform:          rawAgentPlatform(o.Platform.String()),
		State:             o.State,
		SystemId:          o.SystemId,
		HasCredential:     o.HasCredential,
		Error:             o.Error,
		StatusMessage:     o.StatusMessage,
		AosVersion:        o.AosVersion,
	}
}

type rawAgentStatus struct {
	ConnectionState   rawAgentCxnState `json:"connection_state"`
	PackagesInstalled rawAgentPackages `json:"packages_installed"`
	CurrentTask       string           `json:"current_task"`
	PendingType       string           `json:"pending_type"`
	JobId             JobId            `json:"job_id"`
	OperationMode     rawAgentMode     `json:"operation_mode"`
	PlatformVersion   string           `json:"platform_version"`
	Platform          rawAgentPlatform `json:"platform"`
	State             string           `json:"state"`
	SystemId          SystemId         `json:"system_id"`
	HasCredential     bool             `json:"has_credential"`
	Error             string           `json:"error"`
	StatusMessage     string           `json:"status_message"`
	AosVersion        string           `json:"aos_version"`
}

func (o *rawAgentStatus) polish() *AgentStatus {
	return &AgentStatus{
		ConnectionState:   AgentCxnState(o.ConnectionState.parse()),
		PackagesInstalled: o.PackagesInstalled.polish(),
		CurrentTask:       o.CurrentTask,
		PendingType:       o.PendingType,
		JobId:             o.JobId,
		OperationMode:     AgentMode(o.OperationMode.parse()),
		PlatformVersion:   o.PlatformVersion,
		Platform:          AgentPlatform(o.Platform.parse()),
		State:             o.State,
		SystemId:          o.SystemId,
		HasCredential:     o.HasCredential,
		Error:             o.Error,
		StatusMessage:     o.StatusMessage,
		AosVersion:        o.AosVersion,
	}
}

type AgentJobStatus struct {
	Started        time.Time     `json:"started"`
	JobType        AgentJobType  `json:"job_type"`
	Finished       time.Time     `json:"finished"`
	HostId         string        `json:"host_id"` // todo: device s/n? own type?
	CurrentTask    string        `json:"current_task"`
	IsLogAvailable bool          `json:"is_log_available"`
	JobId          JobId         `json:"job_id"`
	Created        time.Time     `json:"created"`
	State          AgentJobState `json:"state"`
	AgentType      AgentType     `json:"agent_type"`
	Error          string        `json:"error"`
}

func (o *AgentJobStatus) raw() *rawAgentJobStatus {
	return &rawAgentJobStatus{
		Started:        o.Started,
		JobType:        o.JobType.raw(),
		Finished:       o.Finished,
		HostId:         o.HostId,
		CurrentTask:    o.CurrentTask,
		IsLogAvailable: o.IsLogAvailable,
		JobId:          o.JobId,
		Created:        o.Created,
		State:          o.State.raw(),
		AgentType:      o.AgentType.raw(),
		Error:          o.Error,
	}
}

type rawAgentJobStatus struct {
	Started        time.Time        `json:"started"`
	JobType        rawAgentJobType  `json:"job_type"`
	Finished       time.Time        `json:"finished"`
	HostId         string           `json:"host_id"` // todo: device s/n? own type?
	CurrentTask    string           `json:"current_task"`
	IsLogAvailable bool             `json:"is_log_available"`
	JobId          JobId            `json:"job_id"`
	Created        time.Time        `json:"created"`
	State          rawAgentJobState `json:"state"`
	AgentType      rawAgentType     `json:"agent_type"`
	Error          string           `json:"error"`
}

func (o *rawAgentJobStatus) polish() *AgentJobStatus {
	return &AgentJobStatus{
		Started:        o.Started,
		JobType:        AgentJobType(o.JobType.parse()),
		Finished:       o.Finished,
		HostId:         o.HostId,
		CurrentTask:    o.CurrentTask,
		IsLogAvailable: o.IsLogAvailable,
		JobId:          o.JobId,
		Created:        o.Created,
		State:          AgentJobState(o.State.parse()),
		AgentType:      AgentType(o.AgentType.parse()),
		Error:          o.Error,
	}
}

type TelemetryExtStatus struct {
	PackagesInstalled AgentPackages `json:"packages_installed"`
	StatusMessage     string        `json:"status_message"`
}

func (o *TelemetryExtStatus) raw() *rawTelemetryExtStatus {
	return &rawTelemetryExtStatus{
		PackagesInstalled: o.PackagesInstalled.raw(),
		StatusMessage:     o.StatusMessage,
	}
}

type rawTelemetryExtStatus struct {
	PackagesInstalled rawAgentPackages `json:"packages_installed"`
	StatusMessage     string           `json:"status_message"`
}

func (o *rawTelemetryExtStatus) polish() *TelemetryExtStatus {
	return &TelemetryExtStatus{
		PackagesInstalled: o.PackagesInstalled.polish(),
		StatusMessage:     o.StatusMessage,
	}
}

type PlatformStatus struct {
	JobId           JobId  `json:"job_id"`
	PlatformVersion string `json:"platform_version"`
	Platform        string `json:"platform"`
	State           string `json:"state"`
	HasCredential   bool   `json:"has_credential"`
	Error           string `json:"error"`
	CurrentTask     string `json:"current_task"`
}

type ContainerStatus struct {
	Status     string `json:"status"`
	Name       string `json:"name"`
	TaskId     string `json:"task_id"`
	LastUpdate string `json:"last_update"`
	Host       string `json:"host"`
	Error      string `json:"error"`
	ServiceId  string `json:"service_id"`
}

type DeviceFacts struct {
	DeviceOsFamily  string `json:"device_os_family"`
	Hostname        string `json:"hostname"`
	DeviceState     string `json:"device_state"`
	DeviceOsVersion string `json:"device_os_version"`
}

type PlatformConfig struct {
	ContainerEnable bool `json:"container_enable"`
}

type SystemAgent struct {
	Id                 ObjectId           `json:"id"`
	Config             SystemAgentConfig  `json:"config"`
	RunningConfig      SystemAgentConfig  `json:"running_config"`
	LastJobStatus      AgentJobStatus     `json:"last_job_status"`
	Status             AgentStatus        `json:"status"`
	TelemetryExtStatus TelemetryExtStatus `json:"telemetry_ext_status"`
	ContainerStatus    ContainerStatus    `json:"container_status"`
	DeviceFacts        DeviceFacts        `json:"device_facts"`
	PlatformConfig     PlatformConfig     `json:"platform_config"`
	PlatformStatus     PlatformStatus     `json:"platform_status"`
}

func (o SystemAgent) raw() *rawSystemAgent {
	return &rawSystemAgent{
		Config:             *o.Config.raw(),
		Id:                 o.Id,
		LastJobStatus:      *o.LastJobStatus.raw(),
		RunningConfig:      *o.RunningConfig.raw(),
		Status:             *o.Status.raw(),
		TelemetryExtStatus: *o.TelemetryExtStatus.raw(),
		ContainerStatus:    o.ContainerStatus,
		DeviceFacts:        o.DeviceFacts,
		PlatformConfig:     o.PlatformConfig,
		PlatformStatus:     o.PlatformStatus,
	}
}

type rawSystemAgent struct {
	Id                 ObjectId              `json:"id"`
	Status             rawAgentStatus        `json:"status"`
	Config             rawSystemAgentConfig  `json:"config"`
	RunningConfig      rawSystemAgentConfig  `json:"running_config"`
	LastJobStatus      rawAgentJobStatus     `json:"last_job_status"`
	TelemetryExtStatus rawTelemetryExtStatus `json:"telemetry_ext_status"`
	DeviceFacts        DeviceFacts           `json:"device_facts"`
	ContainerStatus    ContainerStatus       `json:"container_status"`
	PlatformConfig     PlatformConfig        `json:"platform_config"`
	PlatformStatus     PlatformStatus        `json:"platform_status"`
}

func (o *rawSystemAgent) polish() *SystemAgent {
	return &SystemAgent{
		Id:                 o.Id,                      //
		Config:             *o.Config.polish(),        //
		LastJobStatus:      *o.LastJobStatus.polish(), //
		RunningConfig:      *o.RunningConfig.polish(), //
		Status:             *o.Status.polish(),        //
		TelemetryExtStatus: *o.TelemetryExtStatus.polish(),
		ContainerStatus:    o.ContainerStatus, //
		DeviceFacts:        o.DeviceFacts,     //
		PlatformConfig:     o.PlatformConfig,  //
		PlatformStatus:     o.PlatformStatus,  //
	}
}

type SystemAgentConfig struct {
	Id                  ObjectId
	Label               string
	Profile             ObjectId
	ForcePackageInstall bool
	InstallRequirements bool
	Packages            AgentPackages
	OpenOptions         map[string]string
	Platform            string
	ManagementIp        string
	AgentType           AgentType
	OperationMode       AgentMode
	AllowedJobTypes     []AgentJobType
}

func (o *SystemAgentConfig) raw() *rawSystemAgentConfig {
	return &rawSystemAgentConfig{
		Profile:             o.Profile,
		ForcePackageInstall: o.ForcePackageInstall,
		InstallRequirements: o.InstallRequirements,
		Packages:            o.Packages.raw(),
		OpenOptions:         o.OpenOptions,
		Label:               o.Label,
		Platform:            o.Platform,
		ManagementIp:        o.ManagementIp,
		AgentType:           rawAgentType(o.AgentType.String()),
		OperationMode:       rawAgentMode(o.OperationMode.String()),
		Id:                  o.Id,
	}
}

type rawSystemAgentConfig struct {
	Profile             ObjectId          `json:"profile"`
	ForcePackageInstall bool              `json:"force_package_install"`
	InstallRequirements bool              `json:"install_requirements"`
	Packages            rawAgentPackages  `json:"packages"`
	OpenOptions         map[string]string `json:"open_options"`
	Label               string            `json:"label"`
	Platform            string            `json:"platform"`
	ManagementIp        string            `json:"management_ip"`
	AgentType           rawAgentType      `json:"agent_type"`
	OperationMode       rawAgentMode      `json:"operation_mode"`
	Id                  ObjectId          `json:"id"`
	AllowedJobTypes     rawAgentJobTypes  `json:"allowed_job_types"`
}

func (o *rawSystemAgentConfig) polish() *SystemAgentConfig {
	return &SystemAgentConfig{
		Profile:             o.Profile,
		ForcePackageInstall: o.ForcePackageInstall,
		InstallRequirements: o.InstallRequirements,
		Packages:            o.Packages.polish(),
		OpenOptions:         o.OpenOptions,
		Label:               o.Label,
		Platform:            o.Platform,
		ManagementIp:        o.ManagementIp,
		AgentType:           AgentType(o.AgentType.parse()),
		OperationMode:       AgentMode(o.OperationMode.parse()),
		Id:                  o.Id,
		AllowedJobTypes:     o.AllowedJobTypes.polish(),
	}
}

// SystemAgentRequest is used when creating/updating system agents
type SystemAgentRequest struct {
	AgentType           AgentType
	ManagementIp        string
	OperationMode       AgentMode
	JobOnCreate         AgentJobType
	Profile             ObjectId
	Username            string
	Password            string
	Packages            AgentPackages
	ForcePackageInstall bool
	InstallRequirements bool
	EnableMonitor       bool
	Label               string
	Platform            AgentPlatform
}

func (o *SystemAgentRequest) raw() *rawSystemAgentRequest {
	return &rawSystemAgentRequest{
		AgentType:           rawAgentType(o.AgentType.String()),
		ManagementIp:        o.ManagementIp,
		Profile:             o.Profile,
		OperationMode:       rawAgentMode(o.OperationMode.String()),
		JobOnCreate:         rawAgentJobType(o.JobOnCreate.String()),
		Username:            o.Username,
		ForcePackageInstall: o.ForcePackageInstall,
		InstallRequirements: o.InstallRequirements,
		EnableMonitor:       o.EnableMonitor,
		Password:            o.Password,
		Packages:            o.Packages.raw(),
		Label:               o.Label,
		Platform:            rawAgentPlatform(o.Platform.String()),
	}
}

type rawSystemAgentRequest struct {
	AgentType           rawAgentType     `json:"agent_type,omitempty"`
	ManagementIp        string           `json:"management_ip,omitempty"`
	Profile             ObjectId         `json:"profile,omitempty"`
	OperationMode       rawAgentMode     `json:"operation_mode,omitempty"`
	JobOnCreate         rawAgentJobType  `json:"job_on_create,omitempty"`
	Username            string           `json:"username,omitempty"`
	ForcePackageInstall bool             `json:"force_package_install,omitempty"`
	InstallRequirements bool             `json:"install_requirements,omitempty"`
	EnableMonitor       bool             `json:"enable_monitor,omitempty"`
	Password            string           `json:"password,omitempty"`
	Packages            rawAgentPackages `json:"packages,omitempty"`
	Label               string           `json:"label,omitempty"`
	Platform            rawAgentPlatform `json:"platform,omitempty"`
}

type jobIdResponse struct {
	Id JobId `json:"id"`
}

func (o *Client) listAgents(ctx context.Context) ([]ObjectId, error) {
	response := &optionsAgentsResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlSystemAgents,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getSystemAgent(ctx context.Context, id ObjectId) (*SystemAgent, error) {
	response := &rawSystemAgent{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlSystemAgentsById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.polish(), nil
}

func (o *Client) getAllAgentsInfo(ctx context.Context) ([]SystemAgent, error) {
	response := &getAgentsResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlSystemAgents,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	var result []SystemAgent
	for _, i := range response.Items {
		result = append(result, *i.polish())
	}
	return result, nil
}

func (o *Client) getSystemAgentByManagementIp(ctx context.Context, ip string) (*SystemAgent, error) {
	asa, err := o.getAllAgentsInfo(ctx)
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

func (o *Client) createAgent(ctx context.Context, request *SystemAgentRequest) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlSystemAgents,
		apiInput:    request.raw(),
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *Client) updateSystemAgent(ctx context.Context, id ObjectId, cfg *SystemAgentRequest) error {
	rawCfg := cfg.raw()
	rawCfg.AgentType = "" // cannot change agent type
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPatch,
		urlStr:      fmt.Sprintf(apiUrlSystemAgentsById, id),
		apiInput:    rawCfg,
		apiResponse: response,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) deleteSystemAgent(ctx context.Context, id ObjectId) error {
	agentInfo, err := o.getSystemAgent(ctx, id)
	if err != nil {
		return fmt.Errorf("error fetching agent info prior to deletion - %w", convertTtaeToAceWherePossible(err))
	}

	_, err = o.SystemAgentRunJob(ctx, id, AgentJobTypeUninstall)
	if err != nil {
		return fmt.Errorf("error running agent uninstall job prior to deletion - %w", convertTtaeToAceWherePossible(err))
	}

	err = o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlSystemAgentsById, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	// wait for agent's system comms status to go down before returning from "deleteSystemAgent" because
	// a) deleteSystem is probably next in line
	// b) apstra complains:
	//    	Can't delete the device in neither STOCKED nor DECOMM state. Device is in OOS-READY state.
	if agentInfo.Status.SystemId != "" {
		minuteCountdown, _ := context.WithTimeout(ctx, 1*time.Minute)
		for {
			systemInfo, err := o.getSystemInfo(minuteCountdown, agentInfo.Status.SystemId)
			if err != nil {
				return fmt.Errorf("error checking system state after agent deletion - %w", convertTtaeToAceWherePossible(err))
			}
			if systemInfo.Status.CommState == systemCommsOn {
				continue
			}
			break
		}
	}
	return nil
}

func (o *Client) systemAgentStartJob(ctx context.Context, id ObjectId, job AgentJobType) (JobId, error) {
	var urlStr string
	switch job {
	case AgentJobTypeCheck:
		urlStr = fmt.Sprintf(apiUrlSystemAgentCheck, id)
	case AgentJobTypeInstall:
		urlStr = fmt.Sprintf(apiUrlSystemAgentInstall, id)
	case AgentJobTypeUninstall:
		urlStr = fmt.Sprintf(apiUrlSystemAgentUninstall, id)
	default:
		return 0, fmt.Errorf("don't know how to run job '%s' (type %d)", job.String(), job)
	}
	response := &jobIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      urlStr,
		apiResponse: response,
	})
	if err != nil {
		return 0, convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *Client) getSystemAgentJobHistory(ctx context.Context, id ObjectId) ([]AgentJobStatus, error) {
	response := &agentJobHistoryResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlSystemAgentJobHistory, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	var result []AgentJobStatus
	for _, js := range response.Items {
		result = append(result, *js.polish())
	}
	return result, nil
}

func (o *Client) getSystemAgentJobStatus(ctx context.Context, agentId ObjectId, jobId JobId) (*AgentJobStatus, error) {
	jobs, err := o.getSystemAgentJobHistory(ctx, agentId)
	if err != nil {
		return nil, fmt.Errorf("error getting agent job history - %w", err)
	}

	// pick out and return the requested jobId
	for _, j := range jobs {
		if j.JobId == jobId {
			return &j, nil
		}
	}

	// jobId not found - return error
	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("agent '%s' job history does not include job '%d'", agentId, jobId),
	}
}

func (o *Client) systemAgentWaitForJobToExist(ctx context.Context, agentId ObjectId, jobId JobId) error {
	// loop until we find a reason to return
	for {
		// bail out if our context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, err := o.getSystemAgentJobStatus(ctx, agentId, jobId)
		if err != nil {
			var ace ApstraClientErr
			if !(errors.As(err, &ace) && ace.Type() == ErrNotfound) {
				// error other than notfound - stop looking - return error
				return fmt.Errorf("error getting job status - %w", err)
			}
		} else {
			// no error - the job exists - clean return
			return nil
		}
		time.Sleep(clientPollingIntervalMs * time.Millisecond)
	}
}

func (o *Client) systemAgentWaitForJobTermination(ctx context.Context, agentId ObjectId, jobId JobId) error {
	// loop until we find a reason to return
	for {
		// bail out if our context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		jobStatus, err := o.getSystemAgentJobStatus(ctx, agentId, jobId)
		if err != nil {
			return fmt.Errorf("error getting job status - %w", err)
		}

		if jobStatus.State.HasExited() {
			return nil
		}

		time.Sleep(clientPollingIntervalMs * time.Millisecond)
	}
}

func (o *Client) systemAgentWaitForConnection(ctx context.Context, agentId ObjectId) error {
	// loop until we find a reason to return
	for {
		// bail out if our context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		agentInfo, err := o.getSystemAgent(ctx, agentId)
		if err != nil {
			return fmt.Errorf("error getting agent info - %w", err)
		}

		switch agentInfo.Status.ConnectionState {
		case AgentCxnStateNone: //onbox agents don't "connect"
			return nil
		case AgentCxnStateConnected:
			return nil
		case AgentCxnStateAuthFail:
			return ApstraClientErr{
				errType: ErrAuthFail,
				err: fmt.Errorf("agent %s connection failure: '%s'",
					agentId, agentInfo.Status.ConnectionState.String()),
			}
		case AgentCxnStateDisconnected:
			// go around again
		default:
			return ApstraClientErr{
				errType: ErrUnknown,
				err: fmt.Errorf("unknown agent %s connection failure: '%s'",
					agentId, agentInfo.Status.ConnectionState.String()),
			}
		}

		time.Sleep(clientPollingIntervalMs * time.Millisecond)
	}
}
