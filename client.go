package goapstra

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

// taskIdResponse data structure is returned by Apstra for *some* operations, when the
// URL Query String includes `async=full`
type taskIdResponse struct {
	TaskId TaskId `json:"task_id"`
}

// objectIdResponse is returned by various calls which create an Apstra object
type objectIdResponse struct {
	Id ObjectId `json:"id"`
}

// ObjectId known to Apstra for various objects/resources
type ObjectId string

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

// applyDefaults sets config elements which have default values
func (o *ClientCfg) applyDefaults() {
	if o.Scheme == "" {
		o.Scheme = defaultScheme
	}
}

func (o ClientCfg) validate() error {
	switch {
	case o.Scheme != defaultScheme && o.Scheme != insecureScheme:
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

	debugStr(1, fmt.Sprintf("Apstra client for %s created", c.baseUrl.String()))

	return c, nil
}

// lock creates (if necessary) a *sync.Mutex in Client.sync, and then locks it.
func (o *Client) lock(id int) {
	o.syncLock.Lock() // lock the map of locks
	defer o.syncLock.Unlock()
	if _, found := o.sync[id]; !found {
		o.sync[id] = &sync.Mutex{}
	}
	o.sync[id].Lock()
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

// GetAllBlueprintIds returns a slice of IDs representing all blueprints
func (o *Client) GetAllBlueprintIds(ctx context.Context) ([]ObjectId, error) {
	return o.getAllBlueprintIds(ctx)
}

// GetBlueprint returns *GetBlueprintResponse detailing the requested blueprint
func (o *Client) GetBlueprint(ctx context.Context, in ObjectId) (*GetBlueprintResponse, error) {
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

// CreateRoutingZone creates an Apstra Routing Zone / Security Zone / VRF
func (o *Client) CreateRoutingZone(ctx context.Context, cfg *CreateRoutingZoneCfg) (ObjectId, error) {
	response, err := o.createRoutingZone(ctx, cfg)
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// DeleteRoutingZone deletes an Apstra Routing Zone / Security Zone / VRF
func (o *Client) DeleteRoutingZone(ctx context.Context, blueprintId ObjectId, zoneId ObjectId) error {
	return o.deleteRoutingZone(ctx, blueprintId, zoneId)
}

// GetRoutingZones returns all Apstra Routing Zones / Security Zones / VRFs
// associated with the specified blueprint
func (o *Client) GetRoutingZones(ctx context.Context, blueprintId ObjectId) ([]SecurityZone, error) {
	return o.getAllRoutingZones(ctx, blueprintId)
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
func (o *Client) GetAnomalies(ctx context.Context) ([]*Anomaly, error) {
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

// GetAsnPool returns, by ObjectId, a specifc ASN pool
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
	o.lock(clientApiPoolRangeMutex)
	defer o.unlock(clientApiPoolRangeMutex)

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

// CreateSystemAgentProfile creates a new System Agent Profile identified by 'cfg'
func (o *Client) CreateSystemAgentProfile(ctx context.Context, cfg *SystemAgentProfileConfig) (ObjectId, error) {
	return o.createSystemAgentProfile(ctx, cfg)
}

// ListSystemAgentProfileIds returns a []ObjectId representing System Agent Profiles
func (o *Client) ListSystemAgentProfileIds(ctx context.Context) ([]ObjectId, error) {
	return o.listSystemAgentProfileIds(ctx)
}

// GetSystemAgentProfile returns the SystemAgentProfile identified by 'id'
func (o *Client) GetSystemAgentProfile(ctx context.Context, id ObjectId) (*SystemAgentProfile, error) {
	return o.getSystemAgentProfile(ctx, id)
}

// GetAllSystemAgentProfiles returns the []SystemAgentProfile representing all
// System Agent Profiles
func (o *Client) GetAllSystemAgentProfiles(ctx context.Context) ([]SystemAgentProfile, error) {
	return o.getAllSystemAgentProfiles(ctx)
}

// UpdateSystemAgentProfile updates a System Agent Profile identified by 'cfg'
func (o *Client) UpdateSystemAgentProfile(ctx context.Context, id ObjectId, cfg *SystemAgentProfileConfig) error {
	return o.updateSystemAgentProfile(ctx, id, cfg)
}

// DeleteSystemAgentProfile deletes the System Agent Profile 'id'
func (o *Client) DeleteSystemAgentProfile(ctx context.Context, id ObjectId) error {
	return o.deleteSystemAgentProfile(ctx, id)
}

// GetSystemAgentProfileByLabel returns the System Agent Profile with the given
// label. Apstra doesn't allow label collisions, so this should be a unique
// match. If no match, an ApstraClientErr with Type ErrNotfound is returned.
func (o *Client) GetSystemAgentProfileByLabel(ctx context.Context, label string) (*SystemAgentProfile, error) {
	return o.getSystemAgentProfileByLabel(ctx, label)
}
