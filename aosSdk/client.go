package aosSdk

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
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
	baseUrl   string
	cfg       *ClientCfg
	login     *userLoginResponse // remove this? token now in defHdrs and session ID does nothing
	loginTime time.Time
	client    *http.Client
	defHdrs   []aosHttpHeader
}

// NewClient creates a Client object
func NewClient(cfg ClientCfg) *Client {
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

	return &Client{cfg: &cfg, baseUrl: baseUrl, client: client, defHdrs: defHdrs, login: &userLoginResponse{}}
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

	if o.login.Token == "" && in.url != apiUrlUserLogin {
		return errors.New("cannot interact with AOS API without token")
	}

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

	// set request headers
	if in.toServerPtr != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for i := range o.defHdrs {
		req.Header.Set(o.defHdrs[i].key, o.defHdrs[i].val)
	}

	// talk to the server
	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http.client.Do - %v", err)
	}
	defer resp.Body.Close()

	// response not okay?
	if resp.StatusCode/100 != 2 {
		// trim authentication token from request
		req.Header.Del("Authtoken")

		// limit response details for URLs known to deal in credentials
		if in.url != apiUrlUserLogin {
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

// todo: need smarter handling of response codes, errors, errors in response body
func (o Client) get(url string, expectedResponseCodes []int, jsonPtr interface{}) error {
	if o.login.Token == "" {
		return errors.New("cannot interact with AOS API without token")
	}

	ctx, cancel := context.WithTimeout(o.cfg.Ctx, o.cfg.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating http Request - %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authtoken", o.login.Token)

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http.client.Do - %v", err)
	}
	defer resp.Body.Close()

	if !intSliceContains(expectedResponseCodes, resp.StatusCode) {
		dump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("unexpected http response code '%d' (permitted: '%s') at '%s' (http dump follows)\n%s",
			resp.StatusCode, strings.Join(intSliceToStringSlice(expectedResponseCodes), ","), url, string(dump))
	}

	if jsonPtr != nil {
		return json.NewDecoder(resp.Body).Decode(jsonPtr)
	}

	return nil
}

// todo: need smarter handling of response codes, errors, errors in response body
func (o *Client) post(url string, payload []byte, expectedResponseCodes []int, jsonPtr interface{}) error {
	if o.login.Token == "" && url != o.baseUrl+apiUrlUserLogin {
		return fmt.Errorf("cannot interact with AOS API without token")
	}

	ctx, cancel := context.WithTimeout(o.cfg.Ctx, o.cfg.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating http Request - %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authtoken", o.login.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http.client.Do - %v", err)
	}
	defer resp.Body.Close()

	if !intSliceContains(expectedResponseCodes, resp.StatusCode) {
		dump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("unexpected http response code '%d' (permitted: '%s') at '%s' (http dump follows)\n%s",
			resp.StatusCode, strings.Join(intSliceToStringSlice(expectedResponseCodes), ","), url, string(dump))
	}

	if jsonPtr != nil {
		return json.NewDecoder(resp.Body).Decode(jsonPtr)
	}

	return nil
}

func (o *Client) Login() error {
	err := o.talkToAos(&talkToAosIn{
		method: httpMethodPost,
		url:    apiUrlUserLogin,
		toServerPtr: &userLoginRequest{
			Username: o.cfg.User,
			Password: o.cfg.Pass,
		},
		fromServerPtr: o.login,
	})
	if err != nil {
		return fmt.Errorf("error in Login - %v", err)
	}

	// stash auth token in client's default set of aos http headers
	o.defHdrs = append(o.defHdrs, aosHttpHeader{
		key: "Authtoken",
		val: o.login.Token,
	})

	// save login time for future token refresh function
	o.loginTime = time.Now()

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

func (o Client) SessionId() string {
	return o.login.Id
}
