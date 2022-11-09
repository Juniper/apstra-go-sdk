// Package goapstra implements API client for Juniper Apstra
package goapstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	DefaultTimeout   = 10 * time.Second
	apstraAuthHeader = "Authtoken"
	ErrUnknown       = iota
	ErrAsnOutOfRange
	ErrAsnRangeOverlap
	ErrRangeOverlap
	ErrAuthFail
	ErrCompatibility
	ErrConflict
	ErrExists
	ErrInUse
	ErrMultipleMatch
	ErrNotfound

	clientPollingIntervalMs = 1000

	clientHttpHeadersMutex = iota
	clientApiResourceAsnPoolRangeMutex
	clientApiResourceVniPoolRangeMutex
	clientApiResourceIpPoolRangeMutex
)

type ApstraClientErr struct {
	errType int
	err     error
}

func (o ApstraClientErr) Error() string {
	return o.err.Error()
}

func (o ApstraClientErr) Type() int {
	return o.errType
}

type apstraHttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// ClientCfg is passed to NewClient() when instantiating a new goapstra Client.
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
	Url          string          // URL to access Apstra
	User         string          // Apstra API/UI username
	Pass         string          // Apstra API/UI password
	LogLevel     int             // set < 0 for no logging
	Logger       Logger          // optional caller-created logger sorted by increasing verbosity
	HttpClient   *http.Client    // optional
	Timeout      time.Duration   // <0 = infinite; 0 = DefaultTimeout; >0 = this value is used
	ErrChan      chan<- error    // async client errors (apstra task polling, etc) sent here
	ctx          context.Context // used for async operations (apstra task polling, etc.)
	Experimental bool            // used to enable experimental features
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
	apiVersion  string                  // as reported by apstra API
	baseUrl     *url.URL                // everything up to the file path, generated based on env and cfg
	cfg         ClientCfg               // passed by the caller when creating Client
	httpClient  apstraHttpClient        // used when talking to apstra
	httpHeaders map[string]string       // default set of http headers
	tmQuit      chan struct{}           // task monitor exit trigger
	taskMonChan chan *taskMonitorMonReq // send tasks for monitoring here
	ctx         context.Context         // copied from ClientCfg, for async operations
	logger      Logger                  // logs sent here
	sync        map[int]*sync.Mutex     // some client operations are not concurrency safe. Their locks live here.
	syncLock    sync.Mutex              // control access to the 'sync' map
}

func (o *Client) NewTwoStageL3ClosClient(ctx context.Context, blueprintId ObjectId) (*TwoStageL3ClosClient, error) {
	bp, err := o.getBlueprintStatus(ctx, blueprintId)
	if err != nil {
		return nil, err
	}
	if bp.Design != refDesignDatacenter {
		return nil, fmt.Errorf("cannot create '%s' client for blueprint '%s' (type '%s')",
			RefDesignTwoStageL3Clos.String(), blueprintId, bp.Design)
	}
	result := &TwoStageL3ClosClient{
		client:      o,
		blueprintId: blueprintId,
	}
	result.mutex = &TwoStageL3ClosMutex{client: result}

	return result, nil
}

func (o ClientCfg) validate() error {
	switch {
	case o.Url == "":
		return errors.New("error Url for Apstra Service cannot be empty")
	case o.User == "":
		return errors.New("error username for Apstra service cannot be empty")
	case o.Pass == "":
		return errors.New("error password for Apstra service cannot be empty")
	}
	return nil
}

// NewClient creates a Client object
func (o ClientCfg) NewClient() (*Client, error) {
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

	httpClient := o.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	c := &Client{
		cfg:         o,
		baseUrl:     baseUrl,
		httpClient:  httpClient,
		httpHeaders: map[string]string{"Accept": "application/json"},
		logger:      logger,
		taskMonChan: make(chan *taskMonitorMonReq),
		ctx:         o.ctx,
		sync:        make(map[int]*sync.Mutex),
	}

	// set default context if necessary
	if c.ctx == nil {
		c.ctx = context.Background()
	}

	v, err := c.getApiVersion(c.ctx)
	if err != nil {
		return nil, err
	}

	if !apstraSupportedApi().Includes(v) {
		msg := fmt.Sprintf("unsupported API version: '%s'", c.apiVersion)
		c.logStr(0, msg)
		if !c.cfg.Experimental {
			return nil, errors.New(msg)
		}
	}

	c.logStr(1, fmt.Sprintf("Apstra client for %s created", c.baseUrl.String()))

	return c, nil
}

func (o *Client) getApiVersion(ctx context.Context) (string, error) {
	if o.apiVersion != "" {
		return o.apiVersion, nil
	}
	apiVersion, err := o.getVersionsApi(ctx)
	if err != nil {
		return "", err
	}
	o.apiVersion = apiVersion.Version
	return o.apiVersion, nil
}

// lock creates (if necessary) a *sync.Mutex in Client.sync, and then locks it.
func (o *Client) lock(id int) {

	o.syncLock.Lock() // lock the map of locks - no defer unlock here, we unlock aggressively in the 'found' case below.
	if mu, found := o.sync[id]; found {
		o.syncLock.Unlock()

		mu.Lock()
	} else {
		mu := &sync.Mutex{}
		mu.Lock()
		o.sync[id] = mu

		o.syncLock.Unlock()
	}
}

// unlock releases the named *sync.Mutex in Client.sync
func (o *Client) unlock(id int) {
	o.sync[id].Unlock()
}

// Login submits username and password from the ClientCfg (Client.cfg) to the
// Apstra API, retrieves an authorization token. It is optional. If the client
// is not already logged in, Apstra will send HTTP 401. The client will log
// itself in and resubmit the request.
func (o *Client) Login(ctx context.Context) error {
	return o.login(ctx)
}

// Logout invalidates the Apstra API token held by Client
func (o *Client) Logout(ctx context.Context) error {
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

// GetAnomalies is limited to 10k response items // todo: pagination?
func (o *Client) GetAnomalies(ctx context.Context) ([]Anomaly, error) {
	result, err := o.getAnomalies(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
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
	// AsnPool "write" operations are not concurrency safe
	// It is important that this lock is performed in the public method, rather than the private
	// one below, because other callers of the private method implement their own locking.
	o.lock(clientApiResourceAsnPoolRangeMutex)
	defer o.unlock(clientApiResourceAsnPoolRangeMutex)

	return o.updateAsnPool(ctx, id, cfg)
}

// CreateAsnPoolRange updates an ASN pool by adding a new AsnRange
func (o *Client) CreateAsnPoolRange(ctx context.Context, poolId ObjectId, newRange IntfIntRange) error {
	return o.createAsnPoolRange(ctx, poolId, newRange)
}

// AsnPoolRangeExists reports whether an exact match range (first and last ASN)
// exists in ASN pool poolId
func (o *Client) AsnPoolRangeExists(ctx context.Context, poolId ObjectId, asnRange IntfIntRange) (bool, error) {
	return o.asnPoolRangeExists(ctx, poolId, asnRange)
}

// DeleteAsnPoolRange updates an ASN pool by adding a new AsnRange
func (o *Client) DeleteAsnPoolRange(ctx context.Context, poolId ObjectId, deleteme IntfIntRange) error {
	return o.deleteAsnPoolRange(ctx, poolId, deleteme)
}

// GetVniPools returns Vni pools configured on Apstra
func (o *Client) GetVniPools(ctx context.Context) ([]VniPool, error) {
	return o.getVniPools(ctx)
}

// ListVniPoolIds returns Vni pools configured on Apstra
func (o *Client) ListVniPoolIds(ctx context.Context) ([]ObjectId, error) {
	return o.listVniPoolIds(ctx)
}

// CreateVniPool adds an Vni pool to Apstra
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

// DeleteVniPool deletes an Vni pool, by ObjectId from Apstra
func (o *Client) DeleteVniPool(ctx context.Context, in ObjectId) error {
	return o.deleteVniPool(ctx, in)
}

// UpdateVniPool updates an Vni pool by ObjectId with new Vni pool config
func (o *Client) UpdateVniPool(ctx context.Context, id ObjectId, cfg *VniPoolRequest) error {
	// VniPool "write" operations are not concurrency safe
	// It is important that this lock is performed in the public method, rather than the private
	// one below, because other callers of the private method implement their own locking.
	o.lock(clientApiResourceVniPoolRangeMutex)
	defer o.unlock(clientApiResourceVniPoolRangeMutex)

	return o.updateVniPool(ctx, id, cfg)
}

// CreateVniPoolRange updates an Vni pool by adding a new VniRange
func (o *Client) CreateVniPoolRange(ctx context.Context, poolId ObjectId, newRange IntfIntRange) error {
	return o.createVniPoolRange(ctx, poolId, newRange)
}

// VniPoolRangeExists reports whether an exact match range (first and last Vni)
// exists in Vni pool poolId
func (o *Client) VniPoolRangeExists(ctx context.Context, poolId ObjectId, VniRange IntfIntRange) (bool, error) {
	return o.vniPoolRangeExists(ctx, poolId, VniRange)
}

// DeleteVniPoolRange updates an Vni pool by adding a new VniRange
func (o *Client) DeleteVniPoolRange(ctx context.Context, poolId ObjectId, deleteme IntfIntRange) error {
	return o.deleteVniPoolRange(ctx, poolId, deleteme)
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
// match. If no match, an ApstraClientErr with Type ErrNotfound is returned.
func (o *Client) GetAgentProfileByLabel(ctx context.Context, label string) (*AgentProfile, error) {
	return o.getAgentProfileByLabel(ctx, label)
}

// CreateAgent creates an Apstra Agent and returns its ID
func (o *Client) CreateAgent(ctx context.Context, request *SystemAgentRequest) (ObjectId, error) {
	return o.createAgent(ctx, request)
}

// GetSystemAgent returns a SystemAgent structure representing the supplied ID
func (o *Client) GetSystemAgent(ctx context.Context, id ObjectId) (*SystemAgent, error) {
	return o.getSystemAgent(ctx, id)
}

// GetSystemAgentByManagementIp returns *SystemAgent representing the
// Agent with the given "Management Ip" (which in Apstra terms can also
// be a hostname). Apstra doesn't allow management IP collisions, so this should
// be a unique match. If no match, an ApstraClientErr with type ErrNotfound is
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
	return o.getAllSystemsInfo(ctx)
}

// GetSystemInfo returns a *ManagedSystemInfo representing the requested SystemId
func (o *Client) GetSystemInfo(ctx context.Context, id SystemId) (*ManagedSystemInfo, error) {
	return o.getSystemInfo(ctx, id)
}

// UpdateSystem deletes the supplied SystemId
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
	return o.createBlueprintFromTemplate(ctx, req.raw())
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
	return o.deleteBlueprint(ctx, id)
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
	// Ip4Pool "write" operations are not concurrency safe.
	// It is important that this lock is performed in the public method, rather than the private
	// one below, because other callers of the private method implement their own locking.
	o.lock(clientApiResourceIpPoolRangeMutex)
	defer o.unlock(clientApiResourceIpPoolRangeMutex)
	return o.updateIpPool(ctx, fmt.Sprintf(apiUrlResourcesIp4PoolById, poolId), request)
}

// AddSubnetToIp4Pool adds a subnet to an IPv4 resource pool. Overlap with an existing subnet will
// produce an error
func (o *Client) AddSubnetToIp4Pool(ctx context.Context, poolId ObjectId, new *net.IPNet) error {
	return o.addSubnetToIpPool(ctx, poolId, new)
}

// DeleteSubnetFromIp4Pool deletes a subnet from an IPv4 resource pool. If the subnet does not exist,
// an ApstraClientErr with type ErrNotfound will be returned.
func (o *Client) DeleteSubnetFromIp4Pool(ctx context.Context, poolId ObjectId, target *net.IPNet) error {
	return o.deleteSubnetFromIpPool(ctx, poolId, target)
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
	// Ip6Pool "write" operations are not concurrency safe.
	// It is important that this lock is performed in the public method, rather than the private
	// one below, because other callers of the private method implement their own locking.
	o.lock(clientApiResourceIpPoolRangeMutex)
	defer o.unlock(clientApiResourceIpPoolRangeMutex)
	return o.updateIpPool(ctx, fmt.Sprintf(apiUrlResourcesIp4PoolById, poolId), request)
}

// AddSubnetToIp6Pool adds a subnet to an IPv6 resource pool. Overlap with an existing subnet will
// produce an error
func (o *Client) AddSubnetToIp6Pool(ctx context.Context, poolId ObjectId, new *net.IPNet) error {
	return o.addSubnetToIpPool(ctx, poolId, new)
}

// DeleteSubnetFromIp6Pool deletes a subnet from an IPv6 resource pool. If the subnet does not exist,
// an ApstraClientErr with type ErrNotfound will be returned.
func (o *Client) DeleteSubnetFromIp6Pool(ctx context.Context, poolId ObjectId, target *net.IPNet) error {
	return o.deleteSubnetFromIpPool(ctx, poolId, target)
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
	return o.createLogicalDevice(ctx, in.raw())
}

// UpdateLogicalDevice replaces the whole logical device configuration specified
// by id with the supplied details.
func (o *Client) UpdateLogicalDevice(ctx context.Context, id ObjectId, in *LogicalDeviceData) error {
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

// ListAllTemplateIds returns []ObjectId representing all blueprint templates
func (o *Client) ListAllTemplateIds(ctx context.Context) ([]ObjectId, error) {
	return o.listAllTemplateIds(ctx)
}

// GetAllTemplates returns []Template where each element
// is one of these:
//   TemplateRackBased
//   TemplatePodBased
//   TemplateL3Collapsed
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
// specified name. If zero or more than one templates use the name, an error is returned.
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
// specified name. If zero or more than one templates use the name, an error is returned.
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
// specified name. If zero or more than one templates use the name, an error is returned.
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
func (o *Client) UpdateRackBasedTemplate(ctx context.Context, id ObjectId, in *CreateRackBasedTemplateRequest) (ObjectId, error) {
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
func (o *Client) UpdatePodBasedTemplate(ctx context.Context, id ObjectId, in *CreatePodBasedTemplateRequest) (ObjectId, error) {
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
func (o *Client) UpdateL3CollapsedTemplate(ctx context.Context, id ObjectId, in *CreateL3CollapsedTemplateRequest) (ObjectId, error) {
	return o.updateL3CollapsedTemplate(ctx, id, in)
}

// DeleteTemplate deletes the template specified by id
func (o *Client) DeleteTemplate(ctx context.Context, id ObjectId) error {
	return o.deleteTemplate(ctx, id)
}

// NewQuery returns a *QEQuery with embedded *Client
func (o *Client) NewQuery(blueprint ObjectId) *QEQuery {
	return o.newQuery(blueprint)
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

// GetNodes fetches the node of the specified type, unpacks the API response
// into 'response'
func (o *Client) GetNodes(ctx context.Context, blueprint ObjectId, nodeType NodeType, response interface{}) error {
	return o.getNodes(ctx, blueprint, nodeType, response)
}

// PatchNode patches (only submitted fields are changed) the specified node
// using the contents of 'request', the server's response (whole node info
// without map wrapper?) is returned in 'response'
func (o *Client) PatchNode(ctx context.Context, blueprint ObjectId, node ObjectId, request interface{}, response interface{}) error {
	return o.patchNode(ctx, blueprint, node, request, response)
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
	return o.apiVersion
}

// GetDeviceProfile returns device profile
func (o *Client) GetDeviceProfile(ctx context.Context, id ObjectId) (*DeviceProfile, error) {
	return o.getDeviceProfile(ctx, id)
}

// GetAllDeviceProfiles returns []DeviceProfile
func (o *Client) GetAllDeviceProfiles(ctx context.Context) ([]DeviceProfile, error) {
	return o.getAllDeviceProfiles(ctx)
}

// GetDeviceProfilesByName returns []DeviceProfile including all profiles using the desired name
func (o *Client) GetDeviceProfilesByName(ctx context.Context, desired string) ([]DeviceProfile, error) {
	return o.getDeviceProfilesByName(ctx, desired)
}

// GetDeviceProfileByName returns *DeviceProfile indicating the device profile which uses the
// desired name, or an error if 0 or > 1 device profiles match.
func (o *Client) GetDeviceProfileByName(ctx context.Context, desired string) (*DeviceProfile, error) {
	return o.getDeviceProfileByName(ctx, desired)
}

// CreateDeviceProfile creates device profile
func (o *Client) CreateDeviceProfile(ctx context.Context, profile DeviceProfile) (ObjectId, error) {
	return o.createDeviceProfile(ctx, profile)
}

// UpdateDeviceProfile updates existing device profile
func (o *Client) UpdateDeviceProfile(ctx context.Context, id ObjectId, profile DeviceProfile) error {
	return o.updateDeviceProfile(ctx, id, profile)
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
