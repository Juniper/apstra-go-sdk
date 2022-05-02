package aosSdk

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
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
)

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
	login     *userLoginResponse
	loginTime time.Time
	client    *http.Client
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

	return &Client{cfg: &cfg, baseUrl: baseUrl, client: client, login: &userLoginResponse{}}
}

// todo: need smarter handling of response codes, errors, errors in response body
func (o Client) get(url string, expectedResponseCodes []int, jsonPtr interface{}) error {
	if o.login.Token == "" {
		return fmt.Errorf("cannot interact with AOS API without token")
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

// todo: need smarter handling of response codes, errors, errors in response body
func (o *Client) delete(url string, expectedResponseCodes []int) error {
	if o.login.Token == "" {
		return fmt.Errorf("cannot interact with AOS API without token")
	}

	ctx, cancel := context.WithTimeout(o.cfg.Ctx, o.cfg.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
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

	return nil
}

func (o *Client) Login() error {
	return o.userLogin()
}

func (o Client) Logout() error {
	return o.userLogout()
}

func (o Client) GetStreamingConfigs() ([]StreamingConfigCfg, error) {
	return o.getAllStreamingConfigs()
}

func (o Client) GetVersion() (*VersionResponse, error) {
	return o.getVersion()
}
