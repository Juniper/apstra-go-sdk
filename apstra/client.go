// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	mutexmap "github.com/Juniper/apstra-go-sdk/mutex_map"
	"github.com/hashicorp/go-version"
)

const (
	DefaultTimeout   = 10 * time.Second
	apstraAuthHeader = "Authtoken"
	ErrUnknown       = iota
	ErrAsnOutOfRange
	ErrAsnRangeOverlap
	ErrCannotChangeTransform
	ErrRangeOverlap
	ErrAuthFail
	ErrCompatibility
	ErrConflict
	ErrExists
	ErrInUse
	ErrMultipleMatch
	ErrNotfound
	ErrNotSupported
	ErrUncommitted
	ErrWrongType
	ErrReadOnly
	ErrCtAssignedToLink
	ErrLagHasAssignedStructrues
	ErrTimeout
	ErrAgentProfilePlatformRequired
	ErrIbaCurrentMountConflictsWithExistingMount
	ErrInvalidId
	ErrUnsafePatchProhibited
	ErrCtAssignmentFailed

	clientPollingIntervalMs = 1000

	defaultTimerPollingIntervalMs = 1000
	defaultTimerRetryIntervalMs   = 100
	defaultTimerTimeoutSec        = 10
	defaultMaxRetries             = 5

	mutexKeySeparator   = ":"
	mutexKeyHttpHeaders = "http headers"
	tuningParamLock     = "tuning param lock"
	userAgentDefault    = "apstra-go-sdk"
)

type ErrCtAssignedToLinkDetail struct {
	LinkIds []ObjectId
}

type ErrLagHasAssignedStructuresDetail struct {
	GroupLabels []string
}

type ErrCtAssignmentFailedDetail struct {
	InvalidApplicationPointIds     []ObjectId
	InvalidConnectivityTemplateIds []ObjectId
}

type ClientErr struct {
	errType   int
	err       error
	detail    interface{}
	retryable bool
}

func (o ClientErr) Error() string {
	return o.err.Error()
}

func (o ClientErr) Type() int {
	return o.errType
}

func (o ClientErr) Detail() interface{} {
	return o.detail
}

func (o ClientErr) IsRetryable() bool {
	return o.retryable
}

type apstraHttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// ClientCfg is passed to NewClient() when instantiating a new apstra Client.
// Scheme, Host, Port, User(name) and Pass(word) describe the Apstra API. Each
// of these can be set by environment variable, the names of which are
// controlled by these constants: EnvApstraScheme, EnvApstraUser, EnvApstraPass,
// EnvApstraHost and EnvApstraPort.
// If Logger is nil, the Client will log to log.Default().
// LogLevel controls log verbosity. 0 is default logging level, higher values
// produce more detailed logs. Negative values disable logging.
// HttpClient is optional.
// Timeout is used to create a contextWithTimeout for any passed contexts which
// do not expire. negative values == infinite timeout, 0/default uses
// DefaultTimeout value, positive values are used directly.
// ErrChan, when not nil, is used by async operations to deliver any errors to
// the caller's code.
type ClientCfg struct {
	Url          string         // URL to access Apstra
	User         string         // Apstra API/UI username
	Pass         string         // Apstra API/UI password
	LogLevel     int            // set < 0 for no logging
	Logger       Logger         // optional caller-created logger sorted by increasing verbosity
	HttpClient   *http.Client   // optional
	Timeout      time.Duration  // <0 = infinite; 0 = DefaultTimeout; >0 = this value is used
	ErrChan      chan<- error   // async client errors (apstra task polling, etc) sent here
	Experimental bool           // used to enable experimental features
	UserAgent    string         // may used to set a custom user-agent
	tuningParams map[string]int // various tunable parameters keyed by name
	APIopsDCID   *string        // indicates that we should be talking to API-ops proxy using this DC ID
}

// TaskId represents outstanding tasks on an Apstra server
type TaskId string

// objectIdResponse is returned by various calls which create an Apstra object
type objectIdResponse struct {
	Id ObjectId `json:"id"`
}

// ObjectId known to Apstra for various objects/resources
type ObjectId string

func (o ObjectId) ObjectId() ObjectId {
	return o
}

func (o ObjectId) String() string {
	return string(o)
}

// Client interacts with an AOS API server
type Client struct {
	apiVersion  *version.Version         // as reported by apstra API
	baseUrl     *url.URL                 // everything up to the file path, generated based on env and cfg
	cfg         ClientCfg                // passed by the caller when creating Client
	id          ObjectId                 // Apstra user ID
	httpClient  apstraHttpClient         // used when talking to apstra
	httpHeaders map[string]string        // default set of http headers
	tmQuit      chan struct{}            // task monitor exit trigger
	taskMonChan chan *taskMonitorMonReq  // send tasks for monitoring here
	ctx         context.Context          // copied from ClientCfg, for async operations
	logger      Logger                   // logs sent here
	mutexMap    mutexmap.MutexMap        // some client operations are not concurrency safe. Their mutexes live here.
	features    map[enum.ApiFeature]bool // true/false indicate feature enabled/disabled status
	skipGzip    bool                     // prevents setting 'Accept-Encoding: gzip' - only implemented for api-ops proxy
}

// GetTuningParam returns a named timer value from the client configuration if one has been configured.
// For un-configured parameters containing "TimeoutSec", "PollingIntervalMs" or "MaxRetries" a default
// value is returned. Other un-configured values return zero.
func (o *Client) GetTuningParam(name string) int {
	o.lock(tuningParamLock)
	defer o.unlock(tuningParamLock)
	result, ok := o.cfg.tuningParams[name]
	if ok {
		return result
	}

	switch {
	case strings.Contains(name, "TimeoutSec"):
		return defaultTimerTimeoutSec
	case strings.Contains(name, "RetryIntervalMs"):
		return defaultTimerRetryIntervalMs
	case strings.Contains(name, "PollingIntervalMs"):
		return defaultTimerPollingIntervalMs
	case strings.Contains(name, "MaxRetries"):
		return defaultMaxRetries
	}

	return result
}

// SetTuningParam configures a value in the client configuration. It's a simple map, so any name may be used.
func (o *Client) SetTuningParam(name string, val int) {
	o.lock(tuningParamLock)
	defer o.unlock(tuningParamLock)
	if o.cfg.tuningParams == nil {
		o.cfg.tuningParams = make(map[string]int)
	}
	o.cfg.tuningParams[name] = val
}

// SetContext sets the internal context.Context used by background pollers. This
// context should not have a timeout/deadline, but can be used to cancel
// background tasks.
func (o *Client) SetContext(ctx context.Context) {
	o.ctx = ctx
}

// ID returns the Apstra User ID associated with the client or an empty string if not logged in.
func (o *Client) ID() ObjectId {
	return o.id
}

func (o *Client) NewTwoStageL3ClosClient(ctx context.Context, blueprintId ObjectId) (*TwoStageL3ClosClient, error) {
	bp, err := o.getBlueprintStatus(ctx, blueprintId)
	if err != nil {
		return nil, err
	}
	if bp.Design != enum.RefDesignDatacenter {
		return nil, fmt.Errorf("cannot create '%s' client for blueprint '%s' (type '%s')",
			enum.RefDesignDatacenter, blueprintId, bp.Design)
	}

	result := &TwoStageL3ClosClient{
		client:        o,
		blueprintId:   blueprintId,
		nodeIdsByType: make(map[NodeType][]ObjectId),
	}
	result.Mutex = &TwoStageL3ClosMutex{client: result}

	return result, nil
}

// CreateFreeformBlueprint returns the ID of the new freeform blueprint
func (o *Client) CreateFreeformBlueprint(ctx context.Context, label string) (ObjectId, error) {
	var request struct {
		Design string `json:"design"`
		Label  string `json:"label"`
	}

	request.Design = enum.RefDesignFreeform.String()
	request.Label = label

	var response postBlueprintsResponse

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlBlueprints,
		apiInput:    &request,
		apiResponse: &response,
	})
	if err != nil {
		return response.Id, convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *Client) NewFreeformClient(ctx context.Context, blueprintId ObjectId) (*FreeformClient, error) {
	bp, err := o.getBlueprintStatus(ctx, blueprintId)
	if err != nil {
		return nil, err
	}
	if bp.Design != enum.RefDesignFreeform {
		return nil, fmt.Errorf("cannot create '%s' client for blueprint '%s' (type '%s')",
			enum.RefDesignFreeform, blueprintId, bp.Design)
	}

	return &FreeformClient{
		client:      o,
		blueprintId: blueprintId,
	}, nil
}

func (o ClientCfg) validate() error {
	switch {
	case o.Url == "":
		return errors.New("error Url for Apstra Service cannot be empty")
	case o.User == "" && o.APIopsDCID == nil:
		return errors.New("error username for Apstra service cannot be empty")
	case o.Pass == "" && o.APIopsDCID == nil:
		return errors.New("error password for Apstra service cannot be empty")
	}

	return nil
}

// NewClient creates a Client object
func (o ClientCfg) NewClient(ctx context.Context) (*Client, error) {
	if proxyId, ok := os.LookupEnv(envAosOpsEdgeId); ok {
		if proxyId == "" {
			return nil, fmt.Errorf("environment variable %s must not be empty if set", envAosOpsEdgeId)
		}
		o.APIopsDCID = &proxyId
	}

	err := o.validate()
	if err != nil {
		return nil, err
	}

	var logger Logger
	if o.Logger == nil && o.LogLevel >= 0 {
		logger = log.Default()
	}

	baseUrl, err := url.Parse(o.Url)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", o.Url, err)
	}

	for strings.HasSuffix(baseUrl.Path, apiUrlPathDelim) {
		baseUrl.Path = strings.TrimSuffix(baseUrl.Path, apiUrlPathDelim)
	}

	httpClient := o.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	httpHeaders := map[string]string{"Accept": "application/json"}
	if o.UserAgent == "" {
		httpHeaders["User-Agent"] = userAgentDefault
	} else {
		httpHeaders["User-Agent"] = o.UserAgent
	}

	c := &Client{
		cfg:         o,
		baseUrl:     baseUrl,
		httpClient:  httpClient,
		httpHeaders: httpHeaders,
		logger:      logger,
		taskMonChan: make(chan *taskMonitorMonReq),
		mutexMap:    mutexmap.NewMutexMap(),
		ctx:         context.Background(),
	}

	if _, ok := os.LookupEnv(envAosOpsNoGzip); ok {
		c.skipGzip = true
	}

	if o.APIopsDCID != nil {
		c.startTaskMonitor() // because this client will never "log in"
	}

	err = c.getFeatures(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting features from new client - %w", err)
	}

	// must call getApiVersion() before apiVersionSupported()
	_, err = c.getApiVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting API version from new client - %w", err)
	}

	if !c.apiVersionSupported() {
		msg := fmt.Sprintf("unsupported API version: '%s'", c.apiVersion)
		c.logStr(0, msg)
		if !c.cfg.Experimental {
			return nil, ClientErr{
				errType: ErrCompatibility,
				err:     errors.New(msg),
			}
		}
	}

	c.logStr(1, fmt.Sprintf("Apstra client for %s created", c.baseUrl.String()))

	return c, nil
}

func (o *Client) getApiVersion(ctx context.Context) (*version.Version, error) {
	if o.apiVersion != nil {
		return o.apiVersion, nil
	}

	apiResponse, err := o.getVersionsApi(ctx)
	if err != nil {
		return nil, err
	}

	apiVersion, err := version.NewVersion(apiResponse.Version)
	if err != nil {
		return nil, fmt.Errorf("failed parsing API version string %q - %w", apiResponse.Version, err)
	}

	// ToDo - This can be removed once the results from /versions/api and /versions/build reconverge.
	//  Expected late 2024.
	if apiVersion.Equal(version.Must(version.NewVersion("5.1.0"))) {
		// This *might* be an AI-enabled 5.0.0 pre-release, which merely *claims* to be 5.1.0.

		// Fetch the Build Version to find out if we're talking to an AI-enabled pre-release.
		buildVerResp, err := o.getVersionsBuild(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed fetching AI prerelease workaround version - %w", err)
		}

		// Parse the Build Version.
		buildVersion, err := version.NewVersion(buildVerResp.Version)
		if err != nil {
			return nil, fmt.Errorf("failed parsing AI prerelease workaround version - %w", err)
		}

		if buildVersion.LessThan(apiVersion) && strings.HasPrefix(buildVersion.Prerelease(), "a") {
			// We're going to trust the "Build" version string rather than the "API" version string
			apiVersion = buildVersion
		}
	}

	// set the embedded API version
	o.apiVersion = apiVersion

	// return the result
	return o.apiVersion, nil
}

func (o *Client) apiVersionSupported() bool {
	if o.apiVersion == nil {
		panic("apiVersionSupported() invoked before o.apiVersion got populated")
	}

	return compatibility.ServerVersionSupported.Check(o.apiVersion)
}

// lock secures a sync.Mutex specified by id. The sync.Mutex will be created
// if it does not exist.
func (o *Client) lock(id string) {
	o.mutexMap.Lock(id)
}

// unlock releases the sync.Mutex specified by id. It is a run-time error if
// the specified sync.Mutex does not exist or is not locked
func (o *Client) unlock(id string) {
	o.mutexMap.Unlock(id)
}

// Login submits username and password from the ClientCfg (Client.cfg) to the
// Apstra API, retrieves an authorization token. It is optional. If the client
// is not already logged in, Apstra will send HTTP 401. The client will log
// itself in and resubmit the request.
func (o *Client) Login(ctx context.Context) error {
	if o.cfg.APIopsDCID != nil {
		return nil // we never log in or out when configured to use the api-ops proxy
	}

	err := o.login(ctx)
	if err != nil {
		return err
	}

	o.startTaskMonitor()

	return nil
}

// Logout invalidates the Apstra API token held by Client
func (o *Client) Logout(ctx context.Context) error {
	if o.cfg.APIopsDCID != nil {
		return nil // we never log in or out when configured to use the api-ops proxy
	}

	o.stopTaskMonitor()

	return o.logout(ctx)
}

// GetBlueprint returns *Blueprint detailing the requested blueprint
func (o *Client) GetBlueprint(ctx context.Context, in ObjectId) (*Blueprint, error) {
	return o.getBlueprint(ctx, in)
}

// GetStreamingConfig returns a slice of *StreamingConfigInfo representing
// the requested Apstra streaming configs / receivers
func (o *Client) GetStreamingConfig(ctx context.Context, id ObjectId) (*StreamingConfigInfo, error) {
	return o.getStreamingConfig(ctx, id)
}

// NewStreamingConfig creates a StreamingConfig (Streaming Receiver) on the
// Apstra server.
func (o *Client) NewStreamingConfig(ctx context.Context, cfg *StreamingConfigParams) (ObjectId, error) {
	response, err := o.newStreamingConfig(ctx, cfg)
	return response.Id, err
}

// DeleteStreamingConfig deletes the specified streaming config / receiver from
// the Apstra server configuration.
func (o *Client) DeleteStreamingConfig(ctx context.Context, id ObjectId) error {
	return o.deleteStreamingConfig(ctx, id)
}

// GetVersion calls apiUrlVersion, returns the Apstra server version as a
// VersionResponse
func (o *Client) GetVersion(ctx context.Context) (*VersionResponse, error) {
	return o.getVersion(ctx)
}

// GetVirtualInfraMgrs returns all Virtual Infrastructure Managers configured in Apstra
func (o *Client) GetVirtualInfraMgrs(ctx context.Context) ([]VirtualInfraMgrInfo, error) {
	return o.getVirtualInfraMgrs(ctx)
}

// GetMetricdbMetrics returns []MetricdbMetric representing the various metricdb
// application/namespace/name paths available to be queried from Apstra
func (o *Client) GetMetricdbMetrics(ctx context.Context) ([]MetricdbMetric, error) {
	response, err := o.getMetricdbMetrics(ctx)
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}

// QueryMetricdb returns a MetricDbQueryResponse including all available data
// for the metric and time range specified in the
func (o *Client) QueryMetricdb(ctx context.Context, q *MetricDbQueryRequest) (*MetricDbQueryResponse, error) {
	return o.queryMetricdb(ctx, q.begin, q.end, q.metric)
}

// GetBlueprintAnomalies returns []BlueprintAnomaly representing all anomalies in
// the blueprint.
func (o *Client) GetBlueprintAnomalies(ctx context.Context, blueprintId ObjectId) ([]BlueprintAnomaly, error) {
	return o.getBlueprintAnomalies(ctx, blueprintId)
}

// GetBlueprintNodeAnomalyCounts returns []BlueprintNodeAnomalyCounts
// which summarize current anomalies on a per-node basis in the blueprint.
// Nodes which are not currently experiencing an anomaly are not represented in
// the returned slice.
func (o *Client) GetBlueprintNodeAnomalyCounts(ctx context.Context, blueprintId ObjectId) ([]BlueprintNodeAnomalyCounts, error) {
	return o.getBlueprintNodeAnomalyCounts(ctx, blueprintId)
}

// GetBlueprintServiceAnomalyCounts returns []BlueprintServiceAnomalyCount
// which summarize current anomalies on a per-service basis in the blueprint.
// Services which are not currently experiencing an anomaly are not represented
// in the returned slice.
func (o *Client) GetBlueprintServiceAnomalyCounts(ctx context.Context, blueprintId ObjectId) ([]BlueprintServiceAnomalyCount, error) {
	return o.getBlueprintServiceAnomalyCounts(ctx, blueprintId)
}

// GetAsnPools returns ASN pools configured on Apstra
func (o *Client) GetAsnPools(ctx context.Context) ([]AsnPool, error) {
	return o.getAsnPools(ctx)
}

// GetAsnPoolByName returns ASN pools configured on Apstra
func (o *Client) GetAsnPoolByName(ctx context.Context, desired string) (*AsnPool, error) {
	return o.getAsnPoolByName(ctx, desired)
}

// ListAsnPoolIds returns ASN pools configured on Apstra
func (o *Client) ListAsnPoolIds(ctx context.Context) ([]ObjectId, error) {
	return o.listAsnPoolIds(ctx)
}

// CreateAsnPool adds an ASN pool to Apstra
func (o *Client) CreateAsnPool(ctx context.Context, in *AsnPoolRequest) (ObjectId, error) {
	response, err := o.createAsnPool(ctx, in)
	if err != nil {
		return "", fmt.Errorf("error creating ASN pool - %w", err)
	}
	return response, nil
}

// GetAsnPool returns, by ObjectId, a specific ASN pool
func (o *Client) GetAsnPool(ctx context.Context, in ObjectId) (*AsnPool, error) {
	return o.getAsnPool(ctx, in)
}

// DeleteAsnPool deletes an ASN pool, by ObjectId from Apstra
func (o *Client) DeleteAsnPool(ctx context.Context, in ObjectId) error {
	return o.deleteAsnPool(ctx, in)
}

// UpdateAsnPool updates an ASN pool by ObjectId with new ASN pool config
func (o *Client) UpdateAsnPool(ctx context.Context, id ObjectId, cfg *AsnPoolRequest) error {
	return o.updateAsnPool(ctx, id, cfg)
}

// GetIntegerPools returns Integer Pools configured on Apstra
func (o *Client) GetIntegerPools(ctx context.Context) ([]IntPool, error) {
	rawPools, err := o.getIntPools(ctx, apiUrlResourcesIntegerPools)
	if err != nil {
		return nil, err
	}

	result := make([]IntPool, len(rawPools))
	for i, rawPool := range rawPools {
		polished, err := rawPool.polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}

	return result, nil
}

// GetIntegerPoolByName returns Integer Pools configured on Apstra
func (o *Client) GetIntegerPoolByName(ctx context.Context, desired string) (*IntPool, error) {
	raw, err := o.getIntPoolByName(ctx, desired)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// ListIntegerPoolIds returns Integer Pools configured on Apstra
func (o *Client) ListIntegerPoolIds(ctx context.Context) ([]ObjectId, error) {
	return o.listIntPoolIds(ctx, apiUrlResourcesIntegerPools)
}

// CreateIntegerPool adds an Integer Pool to Apstra
func (o *Client) CreateIntegerPool(ctx context.Context, in *IntPoolRequest) (ObjectId, error) {
	response, err := o.createIntPool(ctx, in, apiUrlResourcesIntegerPools)
	if err != nil {
		return "", fmt.Errorf("error creating Integer Pool - %w", err)
	}
	return response, nil
}

// GetIntegerPool returns, by ObjectId, a specific Integer Pool
func (o *Client) GetIntegerPool(ctx context.Context, in ObjectId) (*IntPool, error) {
	rawPool, err := o.getIntPool(ctx, apiUrlResourcesIntegerPoolById, in)
	if err != nil {
		return nil, err
	}
	return rawPool.polish()
}

// DeleteIntegerPool deletes an Integer Pool, by ObjectId from Apstra
func (o *Client) DeleteIntegerPool(ctx context.Context, in ObjectId) error {
	return o.deleteIntPool(ctx, apiUrlResourcesIntegerPoolById, in)
}

// UpdateIntegerPool updates an Integer Pool by ObjectId with new Integer Pool config
func (o *Client) UpdateIntegerPool(ctx context.Context, id ObjectId, cfg *IntPoolRequest) error {
	return o.updateIntPool(ctx, apiUrlResourcesIntegerPoolById, id, cfg)
}

// GetVniPools returns Vni pools configured on Apstra
func (o *Client) GetVniPools(ctx context.Context) ([]VniPool, error) {
	return o.getVniPools(ctx)
}

// ListVniPoolIds returns Vni pools configured on Apstra
func (o *Client) ListVniPoolIds(ctx context.Context) ([]ObjectId, error) {
	return o.listVniPoolIds(ctx)
}

// CreateVniPool adds a VNI pool to Apstra
func (o *Client) CreateVniPool(ctx context.Context, in *VniPoolRequest) (ObjectId, error) {
	response, err := o.createVniPool(ctx, in)
	if err != nil {
		return "", fmt.Errorf("error creating Vni pool - %w", err)
	}
	return response, nil
}

// GetVniPool returns, by ObjectId, a specific Vni pool
func (o *Client) GetVniPool(ctx context.Context, in ObjectId) (*VniPool, error) {
	return o.getVniPool(ctx, in)
}

// GetVniPoolByName returns *VniPool for the specified VNI pool name
func (o *Client) GetVniPoolByName(ctx context.Context, name string) (*VniPool, error) {
	return o.getVniPoolByName(ctx, name)
}

// DeleteVniPool deletes a VNI pool, by ObjectId from Apstra
func (o *Client) DeleteVniPool(ctx context.Context, in ObjectId) error {
	return o.deleteVniPool(ctx, in)
}

// UpdateVniPool updates a VNI pool by ObjectId with new Vni pool config
func (o *Client) UpdateVniPool(ctx context.Context, id ObjectId, cfg *VniPoolRequest) error {
	return o.updateVniPool(ctx, id, cfg)
}

// CreateAgentProfile creates a new Agent Profile identified by 'cfg'
func (o *Client) CreateAgentProfile(ctx context.Context, cfg *AgentProfileConfig) (ObjectId, error) {
	return o.createAgentProfile(ctx, cfg)
}

// ListAgentProfileIds returns a []ObjectId representing Agent Profiles
func (o *Client) ListAgentProfileIds(ctx context.Context) ([]ObjectId, error) {
	return o.listAgentProfileIds(ctx)
}

// GetAgentProfile returns the AgentProfile identified by 'id'
func (o *Client) GetAgentProfile(ctx context.Context, id ObjectId) (*AgentProfile, error) {
	return o.getAgentProfile(ctx, id)
}

// GetAllAgentProfiles returns the []AgentProfileId representing all
// Agent Profiles
func (o *Client) GetAllAgentProfiles(ctx context.Context) ([]AgentProfile, error) {
	return o.getAllAgentProfiles(ctx)
}

// UpdateAgentProfile updates an Agent Profile identified by 'cfg'
func (o *Client) UpdateAgentProfile(ctx context.Context, id ObjectId, cfg *AgentProfileConfig) error {
	return o.updateAgentProfile(ctx, id, cfg)
}

// DeleteAgentProfile deletes the Agent Profile 'id'
func (o *Client) DeleteAgentProfile(ctx context.Context, id ObjectId) error {
	return o.deleteAgentProfile(ctx, id)
}

// GetAgentProfileByLabel returns the Agent Profile with the given
// label. Apstra doesn't allow label collisions, so this should be a unique
// match. If no match, a ClientErr with Type ErrNotfound is returned.
func (o *Client) GetAgentProfileByLabel(ctx context.Context, label string) (*AgentProfile, error) {
	return o.getAgentProfileByLabel(ctx, label)
}

// CreateSystemAgent creates an Apstra System Agent and returns its ID
func (o *Client) CreateSystemAgent(ctx context.Context, request *SystemAgentRequest) (ObjectId, error) {
	return o.createSystemAgent(ctx, request)
}

// ListSystemAgents returns []ObjectId representing all Managed Device System Agents
func (o *Client) ListSystemAgents(ctx context.Context) ([]ObjectId, error) {
	return o.listSystemAgents(ctx)
}

// GetAllSystemAgents returns a SystemAgent structure representing the supplied ID
func (o *Client) GetAllSystemAgents(ctx context.Context) ([]SystemAgent, error) {
	return o.getAllSystemAgents(ctx)
}

// GetSystemAgent returns a SystemAgent structure representing the supplied ID
func (o *Client) GetSystemAgent(ctx context.Context, id ObjectId) (*SystemAgent, error) {
	return o.getSystemAgent(ctx, id)
}

// GetSystemAgentByManagementIp returns *SystemAgent representing the
// Agent with the given "Management Ip" (which in Apstra terms can also
// be a hostname). Apstra doesn't allow management IP collisions, so this should
// be a unique match. If no match, a ClientErr with type ErrNotfound is
// returned.
func (o *Client) GetSystemAgentByManagementIp(ctx context.Context, ip string) (*SystemAgent, error) {
	return o.getSystemAgentByManagementIp(ctx, ip)
}

// UpdateSystemAgent creates an Apstra Agent and returns its ID
func (o *Client) UpdateSystemAgent(ctx context.Context, id ObjectId, request *SystemAgentRequest) error {
	return o.updateSystemAgent(ctx, id, request)
}

// DeleteSystemAgent creates an Apstra Agent and returns its ID
func (o *Client) DeleteSystemAgent(ctx context.Context, id ObjectId) error {
	return o.deleteSystemAgent(ctx, id)
}

// SystemAgentRunJob requests a job be started on the Agent, returns the
// resulting JobId
func (o *Client) SystemAgentRunJob(ctx context.Context, agentId ObjectId, jobType AgentJobType) (*AgentJobStatus, error) {
	jobId, err := o.systemAgentStartJob(ctx, agentId, jobType)
	if err != nil {
		return nil, err
	}

	err = o.systemAgentWaitForJobToExist(ctx, agentId, jobId)
	if err != nil {
		return nil, err
	}

	err = o.systemAgentWaitForJobTermination(ctx, agentId, jobId)
	if err != nil {
		return nil, err
	}

	switch jobType {
	case AgentJobTypeInstall:
		err = o.systemAgentWaitForConnection(ctx, agentId) // todo: this might be a bit much, perhaps we can release this wait sooner?
		if err != nil {
			return nil, err
		}
	default:
	}

	return o.GetSystemAgentJobStatus(ctx, agentId, jobId)
}

// GetSystemAgentJobHistory returns []AgentJobStatus representing all jobs executed by the agent
func (o *Client) GetSystemAgentJobHistory(ctx context.Context, id ObjectId) ([]AgentJobStatus, error) {
	return o.getSystemAgentJobHistory(ctx, id)
}

// GetSystemAgentJobStatus returns *AgentJobStatus for the given agent and job
func (o *Client) GetSystemAgentJobStatus(ctx context.Context, agentId ObjectId, jobId JobId) (*AgentJobStatus, error) {
	return o.getSystemAgentJobStatus(ctx, agentId, jobId)
}

// ListSystems returns []SystemId representing systems configured on the Apstra
// server.
func (o *Client) ListSystems(ctx context.Context) ([]SystemId, error) {
	return o.listSystems(ctx)
}

// GetAllSystemsInfo returns []ManagedSystemInfo representing all systems
// configured on the Apstra server.
func (o *Client) GetAllSystemsInfo(ctx context.Context) ([]ManagedSystemInfo, error) {
	rawSlice, err := o.getAllSystemsInfo(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]ManagedSystemInfo, len(rawSlice))
	for i, raw := range rawSlice {
		polished, err := raw.polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}

	return result, nil
}

// GetSystemInfo returns a *ManagedSystemInfo representing the requested SystemId
func (o *Client) GetSystemInfo(ctx context.Context, id SystemId) (*ManagedSystemInfo, error) {
	raw, err := o.getSystemInfo(ctx, id)
	if err != nil {
		return nil, err
	}

	return raw.polish()
}

// UpdateSystem updates the supplied SystemId
func (o *Client) UpdateSystem(ctx context.Context, id SystemId, cfg *SystemUserConfig) error {
	return o.updateSystem(ctx, id, cfg)
}

// DeleteSystem deletes the specified SystemId
func (o *Client) DeleteSystem(ctx context.Context, id SystemId) error {
	return o.deleteSystem(ctx, id)
}

// UpdateManagedDevice sets the UserConfig info for a managed system
func (o *Client) UpdateManagedDevice(ctx context.Context, id SystemId, cfg *SystemUserConfig) error {
	return o.updateSystem(ctx, id, cfg)
}

// UpdateManagedDeviceByAgentId sets the UserConfig info for a managed system
func (o *Client) UpdateManagedDeviceByAgentId(ctx context.Context, id ObjectId, cfg *SystemUserConfig) error {
	return o.updateSystemByAgentId(ctx, id, cfg)
}

// ListAllBlueprintIds returns []ObjectId representing all blueprints
func (o *Client) ListAllBlueprintIds(ctx context.Context) ([]ObjectId, error) {
	return o.listAllBlueprintIds(ctx)
}

// GetAllBlueprintStatus returns []BlueprintStatus summarizing blueprints configured on Apstra
func (o *Client) GetAllBlueprintStatus(ctx context.Context) ([]BlueprintStatus, error) {
	rawBpStatuses, err := o.getAllBlueprintStatus(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]BlueprintStatus, len(rawBpStatuses))
	for i, bps := range rawBpStatuses {
		polished, err := bps.polish()
		if err != nil {
			return nil, fmt.Errorf("error polishing blueprint status - %w", err)
		}
		result[i] = *polished
	}
	return result, nil
}

// CreateBlueprintFromTemplate creates a blueprint using the supplied reference design and template
func (o *Client) CreateBlueprintFromTemplate(ctx context.Context, req *CreateBlueprintFromTemplateRequest) (ObjectId, error) {
	if req.FabricSettings != nil && req.FabricSettings.Ipv6Enabled != nil && *req.FabricSettings.Ipv6Enabled {
		return "", errors.New("IPv6 cannot be enabled during blueprint creation")
	}

	var id ObjectId
	var err error
	switch {
	case compatibility.EqApstra420.Check(o.apiVersion):
		id, err = o.createBlueprintFromTemplate420(ctx, req.raw420())
		if err != nil {
			return id, fmt.Errorf("failed while creating new blueprint - %w", err)
		}
	default:
		id, err = o.createBlueprintFromTemplate(ctx, req.raw())
		if err != nil {
			return id, err
		}
	}

	if req.SkipCablingReadinessCheck {
		return id, nil
	}

	err = o.waitForBlueprintCabling(ctx, id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// GetBlueprintStatus returns *BlueprintStatus for the specified blueprint ID
func (o *Client) GetBlueprintStatus(ctx context.Context, id ObjectId) (*BlueprintStatus, error) {
	raw, err := o.getBlueprintStatus(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching blueprint status - %w", err)
	}
	return raw.polish()
}

// GetBlueprintStatusByName returns *BlueprintStatus for the specified blueprint name
func (o *Client) GetBlueprintStatusByName(ctx context.Context, name string) (*BlueprintStatus, error) {
	raw, err := o.getBlueprintStatusByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error fetching blueprint status by name - %w", err)
	}
	return raw.polish()
}

// DeleteBlueprint deletes the specified blueprint
func (o *Client) DeleteBlueprint(ctx context.Context, id ObjectId) error {
	err := o.deleteBlueprint(ctx, id)
	if err != nil {
		return err
	}

	t := immediateTicker(clientPollingIntervalMs)
	defer t.Stop()
	for {
		<-t.C
		ids, err := o.listAllBlueprintIds(ctx)
		if err != nil {
			return err
		}
		if !itemInSlice(id, ids) {
			break
		}
	}
	return nil
}

// CreateIp4Pool creates an IPv4 resource pool
func (o *Client) CreateIp4Pool(ctx context.Context, in *NewIpPoolRequest) (ObjectId, error) {
	return o.createIpPool(ctx, false, in)
}

// ListIp4PoolIds returns []ObjectId representing all IPv4 resource pools
func (o *Client) ListIp4PoolIds(ctx context.Context) ([]ObjectId, error) {
	return o.listIpPoolIds(ctx, apiUrlResourcesIp4Pools)
}

// GetIp4Pools returns all IPv4 pools configured on Apstra
func (o *Client) GetIp4Pools(ctx context.Context) ([]IpPool, error) {
	pools, err := o.getIpPools(ctx, apiUrlResourcesIp4Pools)
	if err != nil {
		return nil, err
	}
	polished := make([]IpPool, len(pools))
	for i, raw := range pools {
		p, err := raw.polish()
		if err != nil {
			return nil, err
		}
		polished[i] = *p
	}
	return polished, nil
}

// GetIp4Pool returns an IPv4 resource pool
func (o *Client) GetIp4Pool(ctx context.Context, poolId ObjectId) (*IpPool, error) {
	raw, err := o.getIpPool(ctx, fmt.Sprintf(apiUrlResourcesIp4PoolById, poolId))
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// GetIp4PoolByName returns an IPv4 resource pool
func (o *Client) GetIp4PoolByName(ctx context.Context, desiredName string) (*IpPool, error) {
	raw, err := o.getIpPoolByName(ctx, apiUrlResourcesIp4Pools, desiredName)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// DeleteIp4Pool deletes the specified IPv4 resource pool
func (o *Client) DeleteIp4Pool(ctx context.Context, id ObjectId) error {
	return o.deleteIpPool(ctx, apiUrlResourcesIp4PoolById, id)
}

// UpdateIp4Pool updates (full replace) an existing IPv4 address pool using a NewIpPoolRequest object
func (o *Client) UpdateIp4Pool(ctx context.Context, poolId ObjectId, request *NewIpPoolRequest) error {
	return o.updateIpPool(ctx, fmt.Sprintf(apiUrlResourcesIp4PoolById, poolId), request)
}

// CreateIp6Pool creates an IPv6 resource pool
func (o *Client) CreateIp6Pool(ctx context.Context, in *NewIpPoolRequest) (ObjectId, error) {
	return o.createIpPool(ctx, true, in)
}

// ListIp6PoolIds returns []ObjectId representing all IPv6 resource pools
func (o *Client) ListIp6PoolIds(ctx context.Context) ([]ObjectId, error) {
	return o.listIpPoolIds(ctx, apiUrlResourcesIp6Pools)
}

// GetIp6Pools returns all IPv6 pools configured on Apstra
func (o *Client) GetIp6Pools(ctx context.Context) ([]IpPool, error) {
	pools, err := o.getIpPools(ctx, apiUrlResourcesIp6Pools)
	if err != nil {
		return nil, err
	}
	polished := make([]IpPool, len(pools))
	for i, raw := range pools {
		p, err := raw.polish()
		if err != nil {
			return nil, err
		}
		polished[i] = *p
	}
	return polished, nil
}

// GetIp6Pool returns an IPv6 resource pool
func (o *Client) GetIp6Pool(ctx context.Context, poolId ObjectId) (*IpPool, error) {
	raw, err := o.getIpPool(ctx, fmt.Sprintf(apiUrlResourcesIp6PoolById, poolId))
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// GetIp6PoolByName returns an IPv6 resource pool
func (o *Client) GetIp6PoolByName(ctx context.Context, desiredName string) (*IpPool, error) {
	raw, err := o.getIpPoolByName(ctx, apiUrlResourcesIp6Pools, desiredName)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// DeleteIp6Pool deletes the specified IPv6 resource pool
func (o *Client) DeleteIp6Pool(ctx context.Context, id ObjectId) error {
	return o.deleteIpPool(ctx, apiUrlResourcesIp6PoolById, id)
}

// UpdateIp6Pool updates (full replace) an existing IPv6 address pool using a NewIpPoolRequest object
func (o *Client) UpdateIp6Pool(ctx context.Context, poolId ObjectId, request *NewIpPoolRequest) error {
	return o.updateIpPool(ctx, fmt.Sprintf(apiUrlResourcesIp6PoolById, poolId), request)
}

// ListLogicalDeviceIds returns a list of logical device IDs configured in Apstra
func (o *Client) ListLogicalDeviceIds(ctx context.Context) ([]ObjectId, error) {
	return o.listRackTypeIds(ctx)
}

// GetLogicalDevice returns the requested *LogicalDevice
func (o *Client) GetLogicalDevice(ctx context.Context, id ObjectId) (*LogicalDevice, error) {
	logicalDevice, err := o.getLogicalDevice(ctx, id)
	if err != nil {
		return nil, err
	}
	return logicalDevice.polish()
}

// GetLogicalDeviceByName returns *LogicalDevice matching name if exactly one
// logical device uses that name. No match or multiple match conditions produce
// and error.
func (o *Client) GetLogicalDeviceByName(ctx context.Context, name string) (*LogicalDevice, error) {
	return o.getLogicalDeviceByName(ctx, name)
}

// CreateLogicalDevice creates a new logical device, returns its ObjectId
func (o *Client) CreateLogicalDevice(ctx context.Context, in *LogicalDeviceData) (ObjectId, error) {
	for _, panel := range in.Panels {
		for _, portGroup := range panel.PortGroups {
			err := portGroup.Roles.Validate()
			if err != nil {
				return "", err
			}
		}
	}

	return o.createLogicalDevice(ctx, in.raw())
}

// UpdateLogicalDevice replaces the whole logical device configuration specified
// by id with the supplied details.
func (o *Client) UpdateLogicalDevice(ctx context.Context, id ObjectId, in *LogicalDeviceData) error {
	for _, panel := range in.Panels {
		for _, portGroup := range panel.PortGroups {
			err := portGroup.Roles.Validate()
			if err != nil {
				return err
			}
		}
	}

	return o.updateLogicalDevice(ctx, id, in.raw())
}

// DeleteLogicalDevice deletes the specified logical device
func (o *Client) DeleteLogicalDevice(ctx context.Context, id ObjectId) error {
	return o.deleteLogicalDevice(ctx, id)
}

// ListAllTags returns []ObjectId representing all DesignTag objects
func (o *Client) ListAllTags(ctx context.Context) ([]ObjectId, error) {
	return o.listAllTags(ctx)
}

// GetTag returns *DesignTag describing the specified ObjectId
func (o *Client) GetTag(ctx context.Context, id ObjectId) (*DesignTag, error) {
	raw, err := o.getTag(ctx, id)
	if err != nil {
		return nil, err
	}
	return raw.polish(), nil
}

// GetTagByLabel returns a *DesignTag matching the supplied DesignTag.Label
// string ("Name" in the web UI). This is a case-sensitive search even though
// apstra enforces uniqueness in a case-insensitive manner. An error is returned
// if no DesignTag objects match the supplied DesignTag.Label.
func (o *Client) GetTagByLabel(ctx context.Context, label string) (*DesignTag, error) {
	raw, err := o.getTagByLabel(ctx, label)
	if err != nil {
		return nil, err
	}
	return raw.polish(), nil
}

// GetTagsByLabels returns []DesignTag matching the supplied slice of labels
// which may not contain duplicates. If any tag does not exist, an error is
// returned.
func (o *Client) GetTagsByLabels(ctx context.Context, labels []string) ([]DesignTag, error) {
	raw, err := o.getTagsByLabels(ctx, labels)
	if err != nil {
		return nil, err
	}
	result := make([]DesignTag, len(raw))
	for i, tag := range raw {
		result[i] = *tag.polish()
	}
	return result, nil
}

// GetAllTags returns []DesignTag describing all DesignTag objects
func (o *Client) GetAllTags(ctx context.Context) ([]DesignTag, error) {
	rawTags, err := o.getAllTags(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]DesignTag, len(rawTags))
	for i, rawTag := range rawTags {
		result[i] = *rawTag.polish()
	}
	return result, nil
}

// CreateTag creates a DesignTag and returns its ObjectId. Note that the
// DesignTag.Label field across all tags is required to be unique and case
// is not considered when making that comparison.
func (o *Client) CreateTag(ctx context.Context, in *DesignTagRequest) (ObjectId, error) {
	return o.createTag(ctx, in)
}

// UpdateTag updates a DesignTag by ObjectId. Note that the DesignTag.Label
// is required, but cannot be changed, so it's really just DesignTag.Description
// that we're allowed to monkey around with.
func (o *Client) UpdateTag(ctx context.Context, id ObjectId, in *DesignTagRequest) error {
	return o.updateTag(ctx, id, in)
}

// DeleteTag deletes the specified DesignTag by its ObjectId
func (o *Client) DeleteTag(ctx context.Context, id ObjectId) error {
	return o.deleteTag(ctx, id)
}

// CreateConfiglet creates a Configlet and returns its ObjectId.
func (o *Client) CreateConfiglet(ctx context.Context, in *ConfigletData) (ObjectId, error) {
	for i, generator := range in.Generators {
		if generator.SectionCondition != "" {
			return "", fmt.Errorf("SectionCondition not supported in Design Catalog - generator[%d] has non-empty SectionCondition", i)
		}
	}

	return o.createConfiglet(ctx, in)
}

// DeleteConfiglet deletes a configlet.
func (o *Client) DeleteConfiglet(ctx context.Context, in ObjectId) error {
	return o.deleteConfiglet(ctx, in)
}

// GetConfiglet Accepts an ID and returns the Configlet object
func (o *Client) GetConfiglet(ctx context.Context, in ObjectId) (*Configlet, error) {
	return o.getConfiglet(ctx, in)
}

// UpdateConfiglet updates a configlet
func (o *Client) UpdateConfiglet(ctx context.Context, id ObjectId, in *ConfigletData) error {
	for i, generator := range in.Generators {
		if generator.SectionCondition != "" {
			return fmt.Errorf("SectionCondition not supported in Design Catalog - generator[%d] has non-empty SectionCondition", i)
		}
	}

	return o.updateConfiglet(ctx, id, in)
}

// GetAllConfiglets returns []Configlet representing all configlets
func (o *Client) GetAllConfiglets(ctx context.Context) ([]Configlet, error) {
	return o.getAllConfiglets(ctx)
}

// ListAllConfiglets gets the List of All configlet IDs
func (o *Client) ListAllConfiglets(ctx context.Context) ([]ObjectId, error) {
	return o.listAllConfiglets(ctx)
}

// GetConfigletByName gets a configlet by Name
func (o *Client) GetConfigletByName(ctx context.Context, Name string) (*Configlet, error) {
	return o.getConfigletByName(ctx, Name)
}

// ListAllTemplateIds returns []ObjectId representing all blueprint templates
func (o *Client) ListAllTemplateIds(ctx context.Context) ([]ObjectId, error) {
	return o.listAllTemplateIds(ctx)
}

// GetAllTemplates returns []Template where each element
// is one of these:
//
//	TemplateRackBased
//	TemplatePodBased
//	TemplateL3Collapsed
func (o *Client) GetAllTemplates(ctx context.Context) ([]Template, error) {
	templates, err := o.getAllTemplates(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]Template, len(templates))
	for i, raw := range templates {
		polished, err := raw.polish()
		if err != nil {
			return nil, err
		}
		result[i] = polished
	}
	return result, nil
}

// GetRackBasedTemplate returns *TemplateRackBased represented by `id`
func (o *Client) GetRackBasedTemplate(ctx context.Context, id ObjectId) (*TemplateRackBased, error) {
	raw, err := o.getRackBasedTemplate(ctx, id)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// GetAllRackBasedTemplates returns []TemplateRackBased representing all rack_based templates
func (o *Client) GetAllRackBasedTemplates(ctx context.Context) ([]TemplateRackBased, error) {
	rawTemplates, err := o.getAllRackBasedTemplates(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]TemplateRackBased, len(rawTemplates))
	for i, rawTemplate := range rawTemplates {
		polished, err := rawTemplate.polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}
	return result, nil
}

// GetRackBasedTemplateByName returns *RackBasedTemplate if exactly one pod_based template uses the
// specified name. If zero templates or more than one template uses the name, an error is returned.
func (o *Client) GetRackBasedTemplateByName(ctx context.Context, name string) (*TemplateRackBased, error) {
	t, err := o.getTemplateByTypeAndName(ctx, templateTypeRackBased, name)
	if err != nil {
		return nil, err
	}
	result := &rawTemplateRackBased{}
	err = json.Unmarshal(*t, result)
	if err != nil {
		return nil, err
	}
	return result.polish()
}

// GetPodBasedTemplate returns *TemplatePodBased represented by `id`
func (o *Client) GetPodBasedTemplate(ctx context.Context, id ObjectId) (*TemplatePodBased, error) {
	raw, err := o.getPodBasedTemplate(ctx, id)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// GetAllPodBasedTemplates returns []TemplatePodBased representing all pod_based templates
func (o *Client) GetAllPodBasedTemplates(ctx context.Context) ([]TemplatePodBased, error) {
	rawTemplates, err := o.getAllPodBasedTemplates(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]TemplatePodBased, len(rawTemplates))
	for i, rawTemplate := range rawTemplates {
		polished, err := rawTemplate.polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}
	return result, nil
}

// GetPodBasedTemplateByName returns *PodBasedTemplate if exactly one pod_based template uses the
// specified name. If zero templates or more than one template uses the name, an error is returned.
func (o *Client) GetPodBasedTemplateByName(ctx context.Context, name string) (*TemplatePodBased, error) {
	t, err := o.getTemplateByTypeAndName(ctx, templateTypePodBased, name)
	if err != nil {
		return nil, err
	}
	result := &rawTemplatePodBased{}
	err = json.Unmarshal(*t, result)
	if err != nil {
		return nil, err
	}
	return result.polish()
}

// GetL3CollapsedTemplate returns *TemplateL3Collapsed represented by `id`
func (o *Client) GetL3CollapsedTemplate(ctx context.Context, id ObjectId) (*TemplateL3Collapsed, error) {
	raw, err := o.getL3CollapsedTemplate(ctx, id)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// GetAllL3CollapsedTemplates returns []TemplateL3Collapsed representing all l3_collapsed templates
func (o *Client) GetAllL3CollapsedTemplates(ctx context.Context) ([]TemplateL3Collapsed, error) {
	rawTemplates, err := o.getAllL3CollapsedTemplates(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]TemplateL3Collapsed, len(rawTemplates))
	for i, rawTemplate := range rawTemplates {
		polished, err := rawTemplate.polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}
	return result, nil
}

// GetL3CollapsedTemplateByName returns *L3CollapsedTemplate if exactly one pod_based template uses the
// specified name. If zero templates or more than one template uses the name, an error is returned.
func (o *Client) GetL3CollapsedTemplateByName(ctx context.Context, name string) (*TemplateL3Collapsed, error) {
	t, err := o.getTemplateByTypeAndName(ctx, templateTypeL3Collapsed, name)
	if err != nil {
		return nil, err
	}
	result := &rawTemplateL3Collapsed{}
	err = json.Unmarshal(*t, result)
	if err != nil {
		return nil, err
	}
	return result.polish()
}

// CreateRackBasedTemplate creates a template based on the supplied CreateRackBasedTempalteRequest
func (o *Client) CreateRackBasedTemplate(ctx context.Context, in *CreateRackBasedTemplateRequest) (ObjectId, error) {
	raw, err := in.raw(ctx, o)
	if err != nil {
		return "", fmt.Errorf("error preparing template request - %w", err)
	}
	return o.createRackBasedTemplate(ctx, raw)
}

// UpdateRackBasedTemplate updates a template based on the supplied CreateRackBasedTempalteRequest
func (o *Client) UpdateRackBasedTemplate(ctx context.Context, id ObjectId, in *CreateRackBasedTemplateRequest) error {
	return o.updateRackBasedTemplate(ctx, id, in)
}

// CreatePodBasedTemplate creates a template based on the supplied CreatePodBasedTempalteRequest
func (o *Client) CreatePodBasedTemplate(ctx context.Context, in *CreatePodBasedTemplateRequest) (ObjectId, error) {
	raw, err := in.raw(ctx, o)
	if err != nil {
		return "", fmt.Errorf("error preparing template request - %w", err)
	}
	return o.createPodBasedTemplate(ctx, raw)
}

// UpdatePodBasedTemplate updates a template based on the supplied CreatePodBasedTempalteRequest
func (o *Client) UpdatePodBasedTemplate(ctx context.Context, id ObjectId, in *CreatePodBasedTemplateRequest) error {
	return o.updatePodBasedTemplate(ctx, id, in)
}

// CreateL3CollapsedTemplate creates a template based on the supplied CreateL3CollapsedTemplateRequest
func (o *Client) CreateL3CollapsedTemplate(ctx context.Context, in *CreateL3CollapsedTemplateRequest) (ObjectId, error) {
	raw, err := in.raw(ctx, o)
	if err != nil {
		return "", fmt.Errorf("error preparing template request - %w", err)
	}
	return o.createL3CollapsedTemplate(ctx, raw)
}

// UpdateL3CollapsedTemplate updates a template based on the supplied CreatePodBasedTempalteRequest
func (o *Client) UpdateL3CollapsedTemplate(ctx context.Context, id ObjectId, in *CreateL3CollapsedTemplateRequest) error {
	return o.updateL3CollapsedTemplate(ctx, id, in)
}

// DeleteTemplate deletes the template specified by id
func (o *Client) DeleteTemplate(ctx context.Context, id ObjectId) error {
	return o.deleteTemplate(ctx, id)
}

// ListAllInterfaceMapIds returns []ObjectId representing all interface maps
func (o *Client) ListAllInterfaceMapIds(ctx context.Context) ([]ObjectId, error) {
	return o.listAllInterfaceMapIds(ctx)
}

// GetInterfaceMap returns *InterfaceMap representing the interface map identified by id
func (o *Client) GetInterfaceMap(ctx context.Context, id ObjectId) (*InterfaceMap, error) {
	interfaceMap, err := o.getInterfaceMap(ctx, id)
	if err != nil {
		return nil, err
	}
	return interfaceMap.polish()
}

// GetAllInterfaceMaps returns []InterfaceMap representing all interface maps
// configured on Apstra
func (o *Client) GetAllInterfaceMaps(ctx context.Context) ([]InterfaceMap, error) {
	interfaceMaps, err := o.getAllInterfaceMaps(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]InterfaceMap, len(interfaceMaps))
	for i := range interfaceMaps {
		polished, err := interfaceMaps[i].polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}
	return result, nil
}

// GetInterfaceMapByName returns *Interface map where exactly one interface map
// uses the desired name.
func (o *Client) GetInterfaceMapByName(ctx context.Context, desired string) (*InterfaceMap, error) {
	raw, err := o.getInterfaceMapByName(ctx, desired)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// CreateInterfaceMap creates an interface map, returns its ObjectId
func (o *Client) CreateInterfaceMap(ctx context.Context, in *InterfaceMapData) (ObjectId, error) {
	return o.createInterfaceMap(ctx, in)
}

// UpdateInterfaceMap updates the interface map represented by id, with the details in ifMap
func (o *Client) UpdateInterfaceMap(ctx context.Context, id ObjectId, ifMap *InterfaceMapData) error {
	return o.updateInterfaceMap(ctx, id, ifMap)
}

// DeleteInterfaceMap deletes the interface map identified by id
func (o *Client) DeleteInterfaceMap(ctx context.Context, id ObjectId) error {
	return o.deleteInterfaceMap(ctx, id)
}

// GetNode fetches the specified node and unpacks it into target
func (o *Client) GetNode(ctx context.Context, blueprint ObjectId, nodeId ObjectId, target interface{}) error {
	return o.getNode(ctx, blueprint, nodeId, target)
}

// GetNodes fetches the node of the specified type, unpacks the API response
// into 'response'
func (o *Client) GetNodes(ctx context.Context, blueprint ObjectId, nodeType NodeType, response interface{}) error {
	return o.getNodes(ctx, blueprint, nodeType, response)
}

// PatchNode patches (only submitted fields are changed) the specified node
// using the contents of 'request'. The server's response (whole node info
// without map wrapper?) is returned in 'response'
func (o *Client) PatchNode(ctx context.Context, blueprint ObjectId, node ObjectId, request interface{}, response interface{}) error {
	return o.patchNode(ctx, blueprint, node, request, response, false)
}

// PatchNodeUnsafe patches (only submitted fields are changed) the specified node
// using the contents of 'request', and the "allow_unsafe=true" request parameter
// required by Apstra 5.0.0. The server's response (whole node info
// without map wrapper?) is returned in 'response'
func (o *Client) PatchNodeUnsafe(ctx context.Context, blueprint ObjectId, node ObjectId, request interface{}, response interface{}) error {
	return o.patchNode(ctx, blueprint, node, request, response, true)
}

// PatchNodes patches (only submitted fields are changed) nodes described
// using the contents of 'request'.
func (o *Client) PatchNodes(ctx context.Context, blueprint ObjectId, request []interface{}) error {
	return o.patchNodes(ctx, blueprint, request)
}

// CreateRackType creates an Apstra Rack Type based on the contents of the
// supplied RackTypeRequest.
// Consistent with the Apstra UI and documentation, logical devices (switches,
// generic systems) and tags cloned within the rack are specified by referencing
// items found in the global catalog. Changes to global catalog items will not
// propagate into previously-created rack types.
func (o *Client) CreateRackType(ctx context.Context, request *RackTypeRequest) (ObjectId, error) {
	raw, err := request.raw(ctx, o)
	if err != nil {
		return "", err
	}
	return o.createRackType(ctx, raw)
}

// UpdateRackType updates the Apstra Rack Type identified by id, based on the
// contents of the supplied RackTypeRequest.
// Consistent with the Apstra UI and documentation, logical devices (switches,
// generic systems) and tags cloned within the rack are specified by referencing
// items found in the global catalog. Changes to global catalog items will not
// propagate into previously-created rack types.
func (o *Client) UpdateRackType(ctx context.Context, id ObjectId, request *RackTypeRequest) error {
	return o.updateRackType(ctx, id, request)
}

// ListRackTypeIds returns a []ObjectId representing all rack types configured
// on Apstra.
func (o *Client) ListRackTypeIds(ctx context.Context) ([]ObjectId, error) {
	return o.listRackTypeIds(ctx)
}

// GetRackType returns *RackType detailing the rack type identified by id.
func (o *Client) GetRackType(ctx context.Context, id ObjectId) (*RackType, error) {
	rt, err := o.getRackType(ctx, id)
	if err != nil {
		return nil, err
	}
	return rt.polish()
}

// GetAllRackTypes returns []RackType representing all rack types configured
// on Apstra.
func (o *Client) GetAllRackTypes(ctx context.Context) ([]RackType, error) {
	return o.getAllRackTypes(ctx)
}

// GetRackTypeByName returns *RackType detailing the rack type identified by name.
func (o *Client) GetRackTypeByName(ctx context.Context, name string) (*RackType, error) {
	return o.getRackTypeByName(ctx, name)
}

// DeleteRackType deletes the rack type identified by id.
func (o *Client) DeleteRackType(ctx context.Context, id ObjectId) error {
	return o.deleteRackType(ctx, id)
}

// Log causes the message to be logged according to the policy for the selected msgLevel
func (o *Client) Log(msgLevel int, msg string) {
	o.logStr(msgLevel, msg)
}

// Logf causes the message to be logged according to the policy for the selected msgLevel
func (o *Client) Logf(msgLevel int, msg string, a ...any) {
	o.logStrf(msgLevel, msg, a...)
}

// ApiVersion returns the version string reported by the Apstra API
func (o *Client) ApiVersion() string {
	return o.apiVersion.String()
}

// GetDeviceProfile returns device profile
func (o *Client) GetDeviceProfile(ctx context.Context, id ObjectId) (*DeviceProfile, error) {
	raw, err := o.getDeviceProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	return raw.polish(), nil
}

// GetAllDeviceProfiles returns []DeviceProfile
func (o *Client) GetAllDeviceProfiles(ctx context.Context) ([]DeviceProfile, error) {
	raw, err := o.getAllDeviceProfiles(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]DeviceProfile, len(raw))
	for i := range raw {
		result[i] = *raw[i].polish()
	}
	return result, nil
}

// GetDeviceProfilesByName returns []DeviceProfile including all profiles using the desired name
func (o *Client) GetDeviceProfilesByName(ctx context.Context, desired string) ([]DeviceProfile, error) {
	raw, err := o.getDeviceProfilesByName(ctx, desired)
	if err != nil {
		return nil, err
	}
	result := make([]DeviceProfile, len(raw))
	for i := range raw {
		result[i] = *raw[i].polish()
	}
	return result, nil
}

// GetDeviceProfileByName returns *DeviceProfile indicating the device profile which uses the
// desired name, or an error if 0 or > 1 device profiles match.
func (o *Client) GetDeviceProfileByName(ctx context.Context, desired string) (*DeviceProfile, error) {
	raw, err := o.getDeviceProfileByName(ctx, desired)
	if err != nil {
		return nil, err
	}
	return raw.polish(), nil
}

// CreateDeviceProfile creates device profile
func (o *Client) CreateDeviceProfile(ctx context.Context, profile *DeviceProfileData) (ObjectId, error) {
	return o.createDeviceProfile(ctx, profile.raw())
}

// UpdateDeviceProfile updates existing device profile
func (o *Client) UpdateDeviceProfile(ctx context.Context, id ObjectId, profile *DeviceProfileData) error {
	return o.updateDeviceProfile(ctx, id, profile.raw())
}

// DeleteDeviceProfile deletes existing device profile
func (o *Client) DeleteDeviceProfile(ctx context.Context, id ObjectId) error {
	return o.deleteDeviceProfile(ctx, id)
}

// ServerName returns the hostname (or IP address string) by which the client
// knows the Apstra server. It's mostly useful when setting up streaming event
// receivers.
func (o *Client) ServerName() string {
	return o.baseUrl.Host
}

// GetTemplateType returns the TemplateType of the template known by id
func (o *Client) GetTemplateType(ctx context.Context, id ObjectId) (TemplateType, error) {
	t, err := o.getTemplateType(ctx, id)
	if err != nil {
		return -1, err
	}
	T, err := t.parse()
	return TemplateType(T), err
}

// GetTemplateIdsTypesByName returns map[ObjectId]TemplateType including all
// templates with the desired name found in the apstra global catalog.
func (o *Client) GetTemplateIdsTypesByName(ctx context.Context, desired string) (map[ObjectId]TemplateType, error) {
	return o.getTemplateIdsTypesByName(ctx, desired)
}

// GetTemplateIdTypeByName returns the ObjectId and TemplateType of the single
// template in the apstra global catalog which uses the name 'desired'. If
// zero templates or more than 1 templates use the name, an error is returned.
func (o *Client) GetTemplateIdTypeByName(ctx context.Context, desired string) (ObjectId, TemplateType, error) {
	return o.getTemplateIdTypeByName(ctx, desired)
}

// GetSystemAgentManagerConfig returns *SystemAgentManagerConfig representing the Advanced Settings
// found on the Managed Devices page of the Web UI.
func (o *Client) GetSystemAgentManagerConfig(ctx context.Context) (*SystemAgentManagerConfig, error) {
	return o.getSystemAgentManagerConfig(ctx)
}

// SetSystemAgentManagerConfig uses a *SystemAgentManagerConfig object to configure the Advanced Settings
// found on the Managed Devices page of the Web UI.
func (o *Client) SetSystemAgentManagerConfig(ctx context.Context, cfg *SystemAgentManagerConfig) error {
	if compatibility.HasDeviceOsImageDownloadTimeout.Check(o.apiVersion) != (cfg.DeviceOsImageDownloadTimeout != nil) {
		return fmt.Errorf("DeviceOsImageDownloadTimeout is required with apstra %s, and must not be used with other versions", compatibility.HasDeviceOsImageDownloadTimeout)
	}
	if !compatibility.SystemManagerHasSkipInterfaceShutdownOnUpgrade.Check(o.apiVersion) && cfg.SkipInterfaceShutdownOnUpgrade {
		return fmt.Errorf("SkipInterfaceShutdownOnUpgrade may only be used with apstra %s", compatibility.SystemManagerHasSkipInterfaceShutdownOnUpgrade)
	}

	return o.setSystemAgentManagerConfig(ctx, cfg)
}

// GetInterfaceMapDigest returns *InterfaceMapDigest representing the supplied ObjectId
func (o *Client) GetInterfaceMapDigest(ctx context.Context, id ObjectId) (*InterfaceMapDigest, error) {
	return o.getInterfaceMapDigest(ctx, id)
}

// GetInterfaceMapDigests returns InterfaceMapDigests representing all interface maps
func (o *Client) GetInterfaceMapDigests(ctx context.Context) (InterfaceMapDigests, error) {
	return o.getInterfaceMapDigests(ctx)
}

// GetInterfaceMapDigestsByDeviceProfile returns InterfaceMapDigests
// representing all interface maps which reference the desired DeviceProfile ID
func (o *Client) GetInterfaceMapDigestsByDeviceProfile(ctx context.Context, desired ObjectId) (InterfaceMapDigests, error) {
	return o.getInterfaceMapDigestsByDeviceProfile(ctx, desired)
}

// GetInterfaceMapDigestsByLogicalDevice returns InterfaceMapDigests
// representing all interface maps which reference the desired LogicalDevice ID
func (o *Client) GetInterfaceMapDigestsByLogicalDevice(ctx context.Context, desired ObjectId) (InterfaceMapDigests, error) {
	return o.getInterfaceMapDigestsByLogicalDevice(ctx, desired)
}

// GetInterfaceMapDigestsLogicalDeviceAndDeviceProfile returns InterfaceMapDigests
// representing all interface maps which reference the desired LogicalDevice ID and DeviceProfile ID
func (o *Client) GetInterfaceMapDigestsLogicalDeviceAndDeviceProfile(ctx context.Context, ldId ObjectId, dpId ObjectId) (InterfaceMapDigests, error) {
	return o.getInterfaceMapDigestsLogicalDeviceAndDeviceProfile(ctx, ldId, dpId)
}

// AssignAgentProfile assigns an AgentProfile to each SystemAgent enumerated in AssignAgentProfileRequest
func (o *Client) AssignAgentProfile(ctx context.Context, req *AssignAgentProfileRequest) error {
	return o.assignAgentProfile(ctx, req)
}

// Ready returns an error if the Apstra service isn't ready to go
func (o *Client) Ready(ctx context.Context) error {
	return o.ready(ctx)
}

// WaitUntilReady blocks until the Apstra service is ready to go, the context
// is cancelled, or a non-retryable error occurs.
func (o *Client) WaitUntilReady(ctx context.Context) error {
	return o.waitUntilReady(ctx)
}

// GetAuditConfig returns current Audit Configuration as *AuditConfig
func (o *Client) GetAuditConfig(ctx context.Context) (*AuditConfig, error) {
	return o.getAuditConfig(ctx)
}

// PutAuditConfig sets Audit Configuration according to passed *AuditConfig
func (o *Client) PutAuditConfig(ctx context.Context, cfg *AuditConfig) error {
	return o.putAuditConfig(ctx, cfg)
}

// ListAllPropertySets returns []ObjectId representing all property sets configured on Apstra
func (o *Client) ListAllPropertySets(ctx context.Context) ([]ObjectId, error) {
	return o.listAllPropertySets(ctx)
}

// GetAllPropertySets returns []PropertySet representing all property sets configured on Apstra
func (o *Client) GetAllPropertySets(ctx context.Context) ([]PropertySet, error) {
	ps, err := o.getAllPropertySets(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]PropertySet, len(ps))
	for i := range ps {
		polished, err := ps[i].polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}
	return result, nil
}

// GetPropertySet returns *PropertySet representing the property set with the given ID
func (o *Client) GetPropertySet(ctx context.Context, id ObjectId) (*PropertySet, error) {
	ps, err := o.getPropertySet(ctx, id)
	if err != nil {
		return nil, err
	}
	return ps.polish()
}

// GetPropertySetByLabel returns *PropertySet representing the only property set with the given label, or an error if multiple property sets share the label.
func (o *Client) GetPropertySetByLabel(ctx context.Context, in string) (*PropertySet, error) {
	ps, err := o.getPropertySetByLabel(ctx, in)
	if err != nil {
		return nil, err
	}
	return ps.polish()
}

// CreatePropertySet creates a property set with the data in PropertySetData. On success, it returns the id of the new property set that was created
func (o *Client) CreatePropertySet(ctx context.Context, in *PropertySetData) (ObjectId, error) {
	return o.createPropertySet(ctx, in)
}

// UpdatePropertySet updates a property set identified by an id with the new set of data
func (o *Client) UpdatePropertySet(ctx context.Context, id ObjectId, in *PropertySetData) error {
	return o.updatePropertySet(ctx, id, in)
}

// DeletePropertySet deletes a property given the id
func (o *Client) DeletePropertySet(ctx context.Context, id ObjectId) error {
	return o.deletePropertySet(ctx, id)
}

// Private method added for Client.ready(), public wrapper not currently needed.
// // GetTelemetryQuery returns *TelemetryQuery
// func (o *Client) GetTelemetryQuery(ctx context.Context) (*TelemetryQueryResponse, error){
//	return o.getTelemetryQuery(ctx)
// }

// DeployBlueprint commits the staging blueprint to the active blueprint
func (o *Client) DeployBlueprint(ctx context.Context, in *BlueprintDeployRequest) (*BlueprintDeployResponse, error) {
	response, err := o.deployBlueprint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.polish()
}

// GetRevisions returns []BlueprintRevision of blueprint 'id' representing
// recent revisions available for rollback
func (o *Client) GetRevisions(ctx context.Context, id ObjectId) ([]BlueprintRevision, error) {
	raw, err := o.getBlueprintRevisions(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]BlueprintRevision, len(raw))
	for i := range raw {
		polished, err := raw[i].polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}
	return result, nil
}

// GetRevision returns *BlueprintRevision representing a specific
// recent blueprint revision number 'rev' of blueprint 'id'
func (o *Client) GetRevision(ctx context.Context, id ObjectId, rev int) (*BlueprintRevision, error) {
	revisions, err := o.getBlueprintRevisions(ctx, id)
	if err != nil {
		return nil, err
	}

	for i := range revisions {
		polished, err := revisions[i].polish()
		if err != nil {
			return nil, err
		}

		if polished.RevisionId == rev {
			return polished, nil
		}
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("blueprint %q revision %d not available in rollback history", id, rev),
	}
}

// GetLastDeployedRevision returns *BlueprintRevision representing the most
// recent deployment of blueprint 'id'
func (o *Client) GetLastDeployedRevision(ctx context.Context, id ObjectId) (*BlueprintRevision, error) {
	revisions, err := o.getBlueprintRevisions(ctx, id)
	if err != nil {
		return nil, err
	}

	highestRevNum := -1
	var highestRevPtr *BlueprintRevision
	for i := range revisions {
		polished, err := revisions[i].polish()
		if err != nil {
			return nil, err
		}
		if polished.RevisionId > highestRevNum {
			highestRevPtr = polished
		}
	}

	if highestRevPtr == nil {
		err = ClientErr{
			errType: ErrUncommitted,
			err:     fmt.Errorf("no commits/deployments of blueprint %q found", id),
		}
	}

	return highestRevPtr, err
}

func (o *Client) BlueprintOverlayControlProtocol(ctx context.Context, id ObjectId) (OverlayControlProtocol, error) {
	nodeAttributes := []QEEAttribute{{"name", QEStringVal("node")}}
	switch {
	case compatibility.BpHasVirtualNetworkPolicyNode.Check(o.apiVersion):
		nodeAttributes = append(nodeAttributes, NodeTypeVirtualNetworkPolicy.QEEAttribute())
	default:
		nodeAttributes = append(nodeAttributes, NodeTypeFabricPolicy.QEEAttribute())
	}

	query := new(PathQuery).
		SetBlueprintId(id).
		SetClient(o).
		Node(nodeAttributes)

	var queryResult struct {
		Items []struct {
			VirtualNetworkPolicy struct {
				OverlayControlProtocol overlayControlProtocol `json:"overlay_control_protocol"`
			} `json:"node"`
		} `json:"items"`
	}

	err := query.Do(ctx, &queryResult)
	if err != nil {
		return 0, fmt.Errorf("error querying blueprint virtual network policy - %w", err)
	}

	if len(queryResult.Items) != 1 {
		return 0, fmt.Errorf("expected 1 overlay_control_protocol node, got %d", len(queryResult.Items))
	}

	ocp, err := queryResult.Items[0].VirtualNetworkPolicy.OverlayControlProtocol.parse()
	if err != nil {
		return 0, fmt.Errorf("error parsing overlay control protocol %q - %w",
			queryResult.Items[0].VirtualNetworkPolicy.OverlayControlProtocol, err)
	}

	return OverlayControlProtocol(ocp), nil
}

// CreateModularDeviceProfile creates a ModularDeviceProfile in Apstra based
// on the supplied object, and returns its ID.
func (o *Client) CreateModularDeviceProfile(ctx context.Context, in *ModularDeviceProfile) (ObjectId, error) {
	return o.createModularDeviceProfile(ctx, in.raw())
}

// GetModularDeviceProfile returns *ModularDeviceProfile found in Apstra with the supplied ID.
func (o *Client) GetModularDeviceProfile(ctx context.Context, id ObjectId) (*ModularDeviceProfile, error) {
	raw, err := o.getModularDeviceProfile(ctx, id)
	if err != nil {
		return nil, err
	}

	return raw.polish(), nil
}

// UpdateModularDeviceProfile updates a ModularDeviceProfile identified by id
// using the supplied ModularDeviceProfile.
func (o *Client) UpdateModularDeviceProfile(ctx context.Context, id ObjectId, cfg *ModularDeviceProfile) error {
	return o.updateModularDeviceProfile(ctx, id, cfg.raw())
}

// DeleteModularDeviceProfile deletes the ModularDeviceProfile identified by id
func (o *Client) DeleteModularDeviceProfile(ctx context.Context, id ObjectId) error {
	return o.deleteModularDeviceProfile(ctx, id)
}
