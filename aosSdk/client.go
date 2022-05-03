package aosSdk

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	EnvApstraUser   = "APSTRA_USER"
	EnvApstraPass   = "APSTRA_PASS"
	EnvApstraHost   = "APSTRA_HOST"
	EnvApstraPort   = "APSTRA_PORT"
	EnvApstraScheme = "APSTRA_SCHEME"

	defaultTimeout = 10 * time.Second

	errResponseLimit = 4096

	httpMethodGet    = httpMethod("GET")
	httpMethodPost   = httpMethod("POST")
	httpMethodDelete = httpMethod("DELETE")
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
}

// Client interacts with an AOS API server
type Client struct {
	baseUrl     string
	cfg         *ClientCfg
	client      *http.Client
	httpHeaders []aosHttpHeader // todo: maybe this should be a dictionary?
	token       *jwt
}

// NewClient creates a Client object
func NewClient(cfg *ClientCfg) *Client {
	if cfg.Ctx == nil {
		cfg.Ctx = context.TODO() // default context
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout // default timeout
	}
	ctxCancel, cancel := context.WithCancel(cfg.Ctx)
	cfg.Ctx = ctxCancel
	cfg.cancel = cancel

	baseUrl := fmt.Sprintf("%s://%s:%d", cfg.Scheme, cfg.Host, cfg.Port)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &cfg.TlsConfig,
		},
	}

	defHdrs := []aosHttpHeader{
		{
			key: "Accept",
			val: "application/json",
		},
	}

	return &Client{cfg: cfg, baseUrl: baseUrl, client: client, httpHeaders: defHdrs}
}

type aosHttpHeader struct {
	key string
	val string
}

type talkToAosIn struct {
	method        httpMethod
	url           string
	toServerPtr   interface{}
	fromServerPtr interface{}
}

type talkToAosErr struct {
	url        string
	request    *http.Request
	response   *bytes.Buffer
	error      string
	statusCode int
}

func (o talkToAosErr) Error() string {
	if o.error == "" {
		return fmt.Sprintf("http response code %d at %s", o.statusCode, o.url)
	}
	return o.error
}

func (o Client) talkToAos(in *talkToAosIn) error {
	var err error
	var body []byte

	// are we sending data to the server?
	if in.toServerPtr != nil {
		body, err = json.Marshal(in.toServerPtr)
		if err != nil {
			return fmt.Errorf("error marshaling payload in talkToAos - %v", err)
		}
	}

	// wrap context with timeout
	ctx, cancel := context.WithTimeout(o.cfg.Ctx, o.cfg.Timeout)
	defer cancel()

	// create request
	req, err := http.NewRequestWithContext(ctx, string(in.method), o.baseUrl+in.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating http Request - %v", err)
	}

	// set request httpHeaders
	if in.toServerPtr != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for i := range o.httpHeaders {
		req.Header.Set(o.httpHeaders[i].key, o.httpHeaders[i].val)
	}

	// talk to the server
	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http.client.Do - %v", err)
	}
	// noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// response not okay?
	if resp.StatusCode/100 != 2 {
		// trim authentication token from request
		req.Header.Del("Authtoken")

		// Auth fail?
		if resp.StatusCode == 401 {
			// Auth fail at login API is fatal
			if in.url == apiUrlUserLogin {
				return talkToAosErr{
					url:        in.url,
					statusCode: resp.StatusCode,
					error:      fmt.Sprintf("HTTP %d at %s - check username/password", resp.StatusCode, apiUrlUserLogin),
				}
			}

			// Try logging in
			err := o.Login()
			if err != nil {
				return err
			}

			// Try the request again
			return o.talkToAos(in)
		}

		// limit response details for URLs known to deal in credentials
		if in.url == apiUrlUserLogin {
			return talkToAosErr{
				url:        in.url,
				statusCode: resp.StatusCode,
			}
		}

		// big error response for "safe" URLs
		return talkToAosErr{
			url:        in.url,
			request:    req,
			response:   bytes.NewBuffer(make([]byte, errResponseLimit)),
			statusCode: resp.StatusCode,
		}
	}

	// caller not expecting any response?
	if in.fromServerPtr == nil {
		return nil
	}

	// decode response body into the caller-specified structure
	return json.NewDecoder(resp.Body).Decode(in.fromServerPtr)
}

func (o *Client) Login() error {
	var response userLoginResponse
	err := o.talkToAos(&talkToAosIn{
		method: httpMethodPost,
		url:    apiUrlUserLogin,
		toServerPtr: &userLoginRequest{
			Username: o.cfg.User,
			Password: o.cfg.Pass,
		},
		fromServerPtr: &response,
	})
	if err != nil {
		return fmt.Errorf("error talking to AOS in Login - %v", err)
	}

	o.token, err = newJwt(response.Token)
	if err != nil {
		return fmt.Errorf("error parsing JWT in Login - %v", err)
	}

	// stash auth token in client's default set of aos http httpHeaders
	// todo: multiple calls to Login() will cause many Authtoken httpHeaders to be saved
	//  convert to a dictionary or seek/destroy duplicates here
	o.httpHeaders = append(o.httpHeaders, aosHttpHeader{
		key: "Authtoken",
		val: o.token.Raw(),
	})

	return nil
}

func (o Client) Logout() error {
	return o.talkToAos(&talkToAosIn{
		method: httpMethodPost,
		url:    apiUrlUserLogout,
	})
}

func (o Client) GetStreamingConfigs() ([]StreamingConfigCfg, error) {
	return o.getAllStreamingConfigs()
}

func (o Client) GetVersion() (*VersionResponse, error) {
	return o.getVersion()
}
