package aosSdk

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	EnvApstraUser          = "APSTRA_USER"
	EnvApstraPass          = "APSTRA_PASS"
	EnvApstraHost          = "APSTRA_HOST"
	EnvApstraPort          = "APSTRA_PORT"
	EnvApstraScheme        = "APSTRA_SCHEME"
	EnvApstraApiKeyLogFile = "APSTRA_API_TLS_KEYFILE"

	defaultTimeout = 10 * time.Second

	httpMethodGet    = httpMethod("GET")
	httpMethodPost   = httpMethod("POST")
	httpMethodDelete = httpMethod("DELETE")

	aosAuthHeader = "Authtoken"
)

type httpMethod string

// ClientCfg passed to NewClient() when instantiating a new Client{}
type ClientCfg struct {
	Scheme    string
	Host      string
	Port      uint16
	User      string
	Pass      string
	TlsConfig tls.Config
	Ctx       context.Context
	Timeout   time.Duration
	cancel    func()
	errChan   chan<- error
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
	httpClient  *http.Client
	httpHeaders map[string]string      // default set of http headers
	tmQuit      chan struct{}          // task monitor exit trigger
	taskMonChan chan taskMontiorMonReq // send tasks for monitoring here
}

// NewClient creates a Client object
func NewClient(cfg *ClientCfg) (*Client, error) {
	if cfg.Ctx == nil {
		cfg.Ctx = context.TODO() // default context
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout // default timeout
	}
	ctxCancel, cancel := context.WithCancel(cfg.Ctx)
	cfg.Ctx = ctxCancel
	cfg.cancel = cancel

	baseUrlString := fmt.Sprintf("%s://%s:%d", cfg.Scheme, cfg.Host, cfg.Port)
	baseUrl, err := url.Parse(baseUrlString)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", baseUrlString, err)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &cfg.TlsConfig,
		},
	}

	aosClient := &Client{
		cfg:         cfg,
		baseUrl:     baseUrl,
		httpClient:  httpClient,
		httpHeaders: map[string]string{"Accept": "application/json"},
		tmQuit:      make(chan struct{}),
		taskMonChan: make(chan taskMontiorMonReq),
	}

	newTaskMonitor(aosClient).start(aosClient.tmQuit)
	return aosClient, nil
}

// todo: add context input to all public Client methods

// ServerName returns the name of the AOS server this client has been configured to use
func (o Client) ServerName() string {
	return o.cfg.Host
}

// Login submits username and password from the ClientCfg (Client.cfg) to the
// Apstra API, retrieves an authorization token. It is optional. If the client
// is not already logged in, Apstra will send HTTP 401. The client will log
// itself in and resubmit the request.
func (o *Client) Login() error {
	return o.login()
}

func (o *Client) login() error {
	aosUrl, err := url.Parse(apiUrlUserLogin)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlUserLogin, err)
	}
	response := &userLoginResponse{}
	err = o.talkToAos(&talkToAosIn{
		method: httpMethodPost,
		url:    aosUrl,
		apiInput: &userLoginRequest{
			Username: o.cfg.User,
			Password: o.cfg.Pass,
		},
		apiResponse: response,
	})
	if err != nil {
		return fmt.Errorf("error talking to AOS in Login - %w", err)
	}

	// stash auth token in client's default set of aos http httpHeaders
	o.httpHeaders[aosAuthHeader] = response.Token

	return nil
}

// Logout invalidates the Apstra API token held by Client
func (o Client) Logout() error {
	return o.logout()
}

func (o Client) logout() error {
	defer close(o.tmQuit) // shut down the task monitor gothread

	aosUrl, err := url.Parse(apiUrlUserLogout)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlUserLogout, err)
	}
	err = o.talkToAos(&talkToAosIn{
		method: httpMethodPost,
		url:    aosUrl,
	})
	if err != nil {
		return fmt.Errorf("error calling '%s' - %w", apiUrlUserLogout, err)
	}
	delete(o.httpHeaders, aosAuthHeader)
	return nil
}

// functions below here are implemented in other files.

// GetAllBlueprintIds returns a slice of IDs representing all blueprints
func (o Client) GetAllBlueprintIds() ([]ObjectId, error) {
	return o.getAllBlueprintIds()
}

// GetBlueprint returns *GetBlueprintResponse detailing the requested blueprint
func (o Client) GetBlueprint(in ObjectId) (*GetBlueprintResponse, error) {
	return o.getBlueprint(in)
}

// GetStreamingConfig returns a slice of *StreamingConfigInfo representing
// the requested Apstra streaming configs / receivers
func (o Client) GetStreamingConfig(id ObjectId) (*StreamingConfigInfo, error) {
	return o.getStreamingConfig(id)
}

// NewStreamingConfig creates a StreamingConfig (Streaming Receiver) on the
// Apstra server.
func (o Client) NewStreamingConfig(cfg *StreamingConfigParams) (ObjectId, error) {
	response, err := o.newStreamingConfig(cfg)
	return response.Id, err
}

// DeleteStreamingConfig deletes the specified streaming config / receiver from
// the Apstra server configuration.
func (o Client) DeleteStreamingConfig(id ObjectId) error {
	return o.deleteStreamingConfig(id)
}

// todo restore this function
//// GetVersion calls apiUrlVersion, returns the Apstra server version as a
//// VersionResponse
//func (o Client) GetVersion() (*VersionResponse, error) {
//	return o.getVersion()
//}

// CreateRoutingZone creates an Apstra Routing Zone / Security Zone / VRF
func (o Client) CreateRoutingZone(cfg *CreateRoutingZoneCfg) (ObjectId, error) {
	response, err := o.createRoutingZone(cfg)
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// DeleteRoutingZone deletes an Apstra Routing Zone / Security Zone / VRF
func (o Client) DeleteRoutingZone(blueprintId ObjectId, zoneId ObjectId) error {
	return o.deleteRoutingZone(blueprintId, zoneId)
}

// GetRoutingZones returns all Apstra Routing Zones / Security Zones / VRFs
// associated with the specified blueprint
func (o Client) GetRoutingZones(blueprintId ObjectId) ([]SecurityZone, error) {
	return o.getAllRoutingZones(blueprintId)
}
