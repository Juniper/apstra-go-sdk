package goapstra

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	EnvApstraUser             = "APSTRA_USER"
	EnvApstraPass             = "APSTRA_PASS"
	EnvApstraHost             = "APSTRA_HOST"
	EnvApstraPort             = "APSTRA_PORT"
	EnvApstraScheme           = "APSTRA_SCHEME"
	EnvApstraApiKeyLogFile    = "APSTRA_API_TLS_LOGFILE"
	EnvApstraStreamKeyLogFile = "APSTRA_STREAM_TLS_LOGFILE"

	defaultTimeout = 10 * time.Second
	defaultScheme  = "https"
	insecureScheme = "http"

	apstraAuthHeader = "Authtoken"

	ErrUnknown = iota
	ErrAsnRangeOverlap
	ErrAsnOutOfRange
	ErrNotfound
	ErrExists
	ErrConflict
	ErrAuthFail
	ErrInUse
	ErrMultipleMatch

	clientPollingIntervalMs = 500

	clientAuthTokenMutex = iota
	clientApiResourceAsnPoolRangeMutex
	clientApiResourceIp4PoolRangeMutex
	clientApiResourceIp6PoolRangeMutex
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

// ClientCfg passed to NewClient() when instantiating a new Client{}
type ClientCfg struct {
	Scheme    string          // "https", probably
	Host      string          // "apstra.company.com" or "192.168.10.10"
	Port      uint16          // zero value for default httpClient behavior
	User      string          // Apstra API/UI username
	Pass      string          // Apstra API/UI password
	TlsConfig *tls.Config     // optional, used with https transactions
	Timeout   time.Duration   // <0 = infinite; 0 = defaultTimeout; >0 = this value is used
	ErrChan   chan<- error    // async client errors (apstra task polling, etc) sent here
	ctx       context.Context // used for async operations (apstra task polling, etc)
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

// Client interacts with an AOS API server
type Client struct {
	baseUrl     *url.URL
	cfg         *ClientCfg
	httpClient  apstraHttpClient
	httpHeaders map[string]string       // default set of http headers
	tmQuit      chan struct{}           // task monitor exit trigger
	taskMonChan chan *taskMonitorMonReq // send tasks for monitoring here
	ctx         context.Context         // copied from ClientCfg, for async operations
	sync        map[int]*sync.Mutex     // some client operations are not concurrency safe. Their locks live here.
	syncLock    sync.Mutex              // control access to the 'sync' map
}

// pullFromEnv tries to pull missing config elements from the environment
func (o *ClientCfg) pullFromEnv() error {
	if o.Scheme == "" {
		o.Scheme = os.Getenv(EnvApstraScheme)
	}
	if o.User == "" {
		o.User = os.Getenv(EnvApstraUser)
	}
	if o.Pass == "" {
		o.Pass = os.Getenv(EnvApstraPass)
	}
	if o.Host == "" {
		o.Host = os.Getenv(EnvApstraHost)
	}
	if o.Port == 0 {
		if portStr, found := os.LookupEnv(EnvApstraPort); found {
			port, err := strconv.ParseUint(portStr, 10, 16)
			if err != nil {
				return fmt.Errorf("error parsing Apstra port - %w", err)
			}
			o.Port = uint16(port)
		}
	}
	return nil
}

func (o *Client) NewTwoStageL3ClosClient(ctx context.Context, blueprintId ObjectId) (*TwoStageLThreeClosClient, error) {
	bp, err := o.getBlueprintStatus(ctx, blueprintId)
	if err != nil {
		return nil, err
	}
	if bp.Design != RefDesignTwoStageL3Clos {
		return nil, fmt.Errorf("cannot create '%s' client for nonexistent blueprint '%s' (type '%s')",
			RefDesignTwoStageL3Clos.String(), blueprintId, bp.Design.String())
	}
	return &TwoStageLThreeClosClient{
		client:      o,
		blueprintId: blueprintId,
	}, nil
}

// applyDefaults sets config elements which have default values
func (o *ClientCfg) applyDefaults() {
	if o.Scheme == "" {
		o.Scheme = defaultScheme
	}
}

func (o ClientCfg) validate() error {
	switch {
	case strings.ToLower(o.Scheme) != defaultScheme && strings.ToLower(o.Scheme) != insecureScheme:
		return fmt.Errorf("error invalid URL scheme for Apstra service '%s'", o.Scheme)
	case o.Host == "":
		return errors.New("error hostname for Apstra service cannot be empty")
	case o.User == "":
		return errors.New("error username for Apstra service cannot be empty")
	case o.Pass == "":
		return errors.New("error password for Apstra service cannot be empty")
	}

	return nil
}

// NewClient creates a Client object
func NewClient(cfg *ClientCfg) (*Client, error) {
	err := cfg.pullFromEnv()
	if err != nil {
		return nil, err
	}

	cfg.applyDefaults()

	err = cfg.validate()
	if err != nil {
		return nil, err
	}

	var portStr string
	if cfg.Port > 0 { // Go default == "unset" for our purposes; this should be safe b/c rfc6335
		portStr = fmt.Sprintf(":%d", cfg.Port)
	}
	baseUrlString := fmt.Sprintf("%s://%s%s", cfg.Scheme, cfg.Host, portStr)
	baseUrl, err := url.Parse(baseUrlString)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", baseUrlString, err)
	}

	tlsCfg := cfg.TlsConfig
	if tlsCfg != nil {
		if tlsCfg.InsecureSkipVerify {
			debugStr(1, "TLS certificate verification disabled")
		}
		klw, err := keyLogWriter(EnvApstraApiKeyLogFile)
		if err != nil {
			return nil, fmt.Errorf("error prepping TLS key log from env var '%s' - %w", EnvApstraApiKeyLogFile, err)
		}
		if klw != nil {
			tlsCfg.KeyLogWriter = klw
			debugStr(1, fmt.Sprintf("TLS session keys being logged to %s", os.Getenv(EnvApstraApiKeyLogFile)))
		}
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: cfg.TlsConfig,
		},
	}

	var ctx context.Context
	if cfg.ctx == nil {
		ctx = context.TODO()
	} else {
		ctx = cfg.ctx
	}
	c := &Client{
		cfg:         cfg,
		baseUrl:     baseUrl,
		httpClient:  httpClient,
		httpHeaders: map[string]string{"Accept": "application/json"},
		tmQuit:      make(chan struct{}),
		taskMonChan: make(chan *taskMonitorMonReq),
		ctx:         ctx,
		sync:        make(map[int]*sync.Mutex),
	}

	newTaskMonitor(c).start()

	debugStr(1, fmt.Sprintf("Apstra client for %s created", c.baseUrl.String()))

	return c, nil
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

// ServerName returns the name of the AOS server this client has been configured to use
func (o *Client) ServerName() string {
	return o.cfg.Host
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

// ListAsnPoolIds returns ASN pools configured on Apstra
func (o *Client) ListAsnPoolIds(ctx context.Context) ([]ObjectId, error) {
	return o.listAsnPoolIds(ctx)
}

// CreateAsnPool adds an ASN pool to Apstra
func (o *Client) CreateAsnPool(ctx context.Context, in *AsnPool) (ObjectId, error) {
	response, err := o.createAsnPool(ctx, in)
	if err != nil {
		return "", fmt.Errorf("error creating ASN pool - %w", err)
	}
	return response.Id, nil
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
func (o *Client) UpdateAsnPool(ctx context.Context, id ObjectId, cfg *AsnPool) error {
	// AsnPool "write" operations are not concurrency safe
	// It is important that this lock is performed in the public method, rather than the private
	// one below, because other callers of the private method implement their own locking.
	o.lock(clientApiResourceAsnPoolRangeMutex)
	defer o.unlock(clientApiResourceAsnPoolRangeMutex)

	return o.updateAsnPool(ctx, id, cfg)
}

// CreateAsnPoolRange updates an ASN pool by adding a new AsnRange
func (o *Client) CreateAsnPoolRange(ctx context.Context, poolId ObjectId, newRange *AsnRange) error {
	return o.createAsnPoolRange(ctx, poolId, newRange)
}

// AsnPoolRangeExists reports whether an exact match range (first and last ASN)
// exists in ASN pool poolId
func (o *Client) AsnPoolRangeExists(ctx context.Context, poolId ObjectId, asnRange *AsnRange) (bool, error) {
	return o.asnPoolRangeExists(ctx, poolId, asnRange)
}

// DeleteAsnPoolRange updates an ASN pool by adding a new AsnRange
func (o *Client) DeleteAsnPoolRange(ctx context.Context, poolId ObjectId, deleteme *AsnRange) error {
	return o.deleteAsnPoolRange(ctx, poolId, deleteme)
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

// UpdateAgentProfile updates a Agent Profile identified by 'cfg'
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

// CreateBlueprintFromTemplate creates a blueprint using the supplied reference design and template
func (o *Client) CreateBlueprintFromTemplate(ctx context.Context, cfg *CreateBluePrintFromTemplate) (ObjectId, error) {
	return o.createBlueprintFromTemplate(ctx, cfg)
}

// GetBlueprintStatus returns *BlueprintStatus for the specified blueprint ID
func (o *Client) GetBlueprintStatus(ctx context.Context, id ObjectId) (*BlueprintStatus, error) {
	return o.getBlueprintStatus(ctx, id)
}

// GetBlueprintStatusByName returns *BlueprintStatus for the specified blueprint name
func (o *Client) GetBlueprintStatusByName(ctx context.Context, name string) (*BlueprintStatus, error) {
	return o.getBlueprintStatusByName(ctx, name)
}

// DeleteBlueprint deletes the specified blueprint
func (o *Client) DeleteBlueprint(ctx context.Context, id ObjectId) error {
	return o.deleteBlueprint(ctx, id)
}

// CreateIp4Pool creates an IPv4 resource pool
func (o *Client) CreateIp4Pool(ctx context.Context, in *NewIp4PoolRequest) (ObjectId, error) {
	return o.createIp4Pool(ctx, in)
}

// ListIp4PoolIds returns []ObjectId representing all IPv4 resource pools
func (o *Client) ListIp4PoolIds(ctx context.Context) ([]ObjectId, error) {
	return o.listIp4PoolIds(ctx)
}

// GetIp4Pools returns all IPv4 pools configured on Apstra
func (o *Client) GetIp4Pools(ctx context.Context) ([]Ip4Pool, error) {
	return o.getIp4Pools(ctx)
}

// GetIp4Pool returns an IPv4 resource pool
func (o *Client) GetIp4Pool(ctx context.Context, poolId ObjectId) (*Ip4Pool, error) {
	return o.getIp4Pool(ctx, poolId)
}

// GetIp4PoolByName returns an IPv4 resource pool
func (o *Client) GetIp4PoolByName(ctx context.Context, desiredName string) (*Ip4Pool, error) {
	return o.getIp4PoolByName(ctx, desiredName)
}

// DeleteIp4Pool deletes the specified IPv4 resource pool
func (o *Client) DeleteIp4Pool(ctx context.Context, id ObjectId) error {
	return o.deleteIp4Pool(ctx, id)
}

// UpdateIp4Pool updates (full replace) an existing IPv4 address pool using a NewIp4PoolRequest object
func (o *Client) UpdateIp4Pool(ctx context.Context, poolId ObjectId, request *NewIp4PoolRequest) error {
	// Ip4Pool "write" operations are not concurrency safe.
	// It is important that this lock is performed in the public method, rather than the private
	// one below, because other callers of the private method implement their own locking.
	o.lock(clientApiResourceIp4PoolRangeMutex)
	defer o.unlock(clientApiResourceIp4PoolRangeMutex)
	return o.updateIp4Pool(ctx, poolId, request)
}

// AddSubnetToIp4Pool adds a subnet to an IPv4 resource pool. Overlap with an existing subnet will
// produce an error
func (o *Client) AddSubnetToIp4Pool(ctx context.Context, poolId ObjectId, new *net.IPNet) error {
	return o.addSubnetToIp4Pool(ctx, poolId, new)
}

// DeleteSubnetFromIp4Pool deletes a subnet from an IPv4 resource pool. If the subnet does not exist,
// an ApstraClientErr with type ErrNotfound will be returned.
func (o *Client) DeleteSubnetFromIp4Pool(ctx context.Context, poolId ObjectId, target *net.IPNet) error {
	return o.deleteSubnetFromIp4Pool(ctx, poolId, target)
}

// ListLogicalDeviceIds returns a list of logical device IDs configured in Apstra
func (o *Client) ListLogicalDeviceIds(ctx context.Context) ([]ObjectId, error) {
	return o.listRackTypeIds(ctx)
}

// GetLogicalDevice returns the requested *LogicalDevice
func (o *Client) GetLogicalDevice(ctx context.Context, id ObjectId) (*LogicalDevice, error) {
	return o.getLogicalDevice(ctx, id)
}

// GetLogicalDeviceByName returns *LogicalDevice matching name if exactly one
// logical device uses that name. No match or multiple match conditions produce
// and error.
func (o *Client) GetLogicalDeviceByName(ctx context.Context, name string) (*LogicalDevice, error) {
	return o.getLogicalDeviceByName(ctx, name)
}

// CreateLogicalDevice creates a new logical device, returns its ObjectId
func (o *Client) CreateLogicalDevice(ctx context.Context, in *LogicalDevice) (ObjectId, error) {
	return o.createLogicalDevice(ctx, in)
}

// UpdateLogicalDevice replaces the whole logical device configuration specified
// by id with the supplied details.
func (o *Client) UpdateLogicalDevice(ctx context.Context, id ObjectId, in *LogicalDevice) error {
	return o.updateLogicalDevice(ctx, id, in)
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
	return o.getTag(ctx, id)
}

// GetTagByLabel returns a *DesignTag matching the supplied DesignTag.Label
// string ("Name" in the web UI). This is a case-insensitive search because
// apstra enforces uniqueness in a case-insensitive manner. An error is returned
// if no DesignTag objects match the supplied DesignTag.Label.
func (o *Client) GetTagByLabel(ctx context.Context, label TagLabel) (*DesignTag, error) {
	return o.getTagByLabel(ctx, label)
}

// GetAllTags returns []DesignTag describing all DesignTag objects
func (o *Client) GetAllTags(ctx context.Context) ([]DesignTag, error) {
	return o.getAllTags(ctx)
}

// CreateTag creates a DesignTag and returns its ObjectId. Note that the
// DesignTag.Label field across all tags is required to be unique and case
// is not considered when making that comparison.
func (o *Client) CreateTag(ctx context.Context, in *DesignTag) (ObjectId, error) {
	return o.createTag(ctx, in)
}

// UpdateTag updates a DesignTag by ObjectId. Note that the DesignTag.Label
// is required, but cannot be changed, so it's really just DesignTag.Description
// that we're allowed to monkey around with.
func (o *Client) UpdateTag(ctx context.Context, id ObjectId, in *DesignTag) (ObjectId, error) {
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

// GetAllTemplates returns map[TemplateType][]interface{} where each element
// is one of these:
//   []TemplateRackBased
//   []TemplatePodBased
//   []TemplateL3Collapsed
func (o *Client) GetAllTemplates(ctx context.Context) (map[TemplateType][]interface{}, error) {
	return o.getAllTemplates(ctx)
}

// GetRackBasedTemplate returns *TemplateRackBased represented by `id`
func (o *Client) GetRackBasedTemplate(ctx context.Context, id ObjectId) (*TemplateRackBased, error) {
	return o.getRackBasedTemplate(ctx, id)
}

// GetAllRackBasedTemplates returns []TemplateRackBased representing all rack_based templates
func (o *Client) GetAllRackBasedTemplates(ctx context.Context) ([]TemplateRackBased, error) {
	return o.getAllRackBasedTemplates(ctx)
}

// GetPodBasedTemplate returns *TemplatePodBased represented by `id`
func (o *Client) GetPodBasedTemplate(ctx context.Context, id ObjectId) (*TemplatePodBased, error) {
	return o.getPodBasedTemplate(ctx, id)
}

// GetAllPodBasedTemplates returns []TemplatePodBased representing all pod_based templates
func (o *Client) GetAllPodBasedTemplates(ctx context.Context) ([]TemplatePodBased, error) {
	return o.getAllPodBasedTemplates(ctx)
}

// GetL3CollapsedTemplate returns *TemplateL3Collapsed represented by `id`
func (o *Client) GetL3CollapsedTemplate(ctx context.Context, id ObjectId) (*TemplateL3Collapsed, error) {
	return o.getL3CollapsedTemplate(ctx, id)
}

// GetAllL3CollapsedTemplates returns []TemplateL3Collapsed representing all l3_collapsed templates
func (o *Client) GetAllL3CollapsedTemplates(ctx context.Context) ([]TemplateL3Collapsed, error) {
	return o.getAllL3CollapsedTemplates(ctx)
}

// GetTemplateAndType returns the TemplateType and template object (*TemplateTypeRackBased, *TemplateTypePodBased,
// *TemplateTypeL3Collapsed) associated with the specified template id.
func (o *Client) GetTemplateAndType(ctx context.Context, id ObjectId) (TemplateType, interface{}, error) {
	return o.getTemplateAndType(ctx, id)
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
	return o.getInterfaceMap(ctx, id)
}

// CreateInterfaceMap creates an interface map, returns its ObjectId
func (o *Client) CreateInterfaceMap(ctx context.Context, in *InterfaceMap) (ObjectId, error) {
	return o.createInterfaceMap(ctx, in)
}

// UpdateInterfaceMap updates the interface map represented by id, with the details in ifMap
func (o *Client) UpdateInterfaceMap(ctx context.Context, id ObjectId, ifMap *InterfaceMap) error {
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
	return o.createRackType(ctx, request)
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
	return o.getRackType(ctx, id)
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
