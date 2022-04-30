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
)

// ClientCfg passed to NewClient() when instantiating a new Client{}
type ClientCfg struct {
	Scheme    string
	Host      string
	Port      uint16
	User      string
	Pass      string
	TlsConfig *tls.Config
	Ctx       context.Context
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
func NewClient(cfg *ClientCfg) (*Client, error) {
	// todo move hxxps insecure bs out of this library
	if cfg.Ctx == nil {
		return nil, fmt.Errorf("refusing to create client without context")
	}
	if cfg.TlsConfig == nil {
		cfg.TlsConfig = &tls.Config{}
	}
	var baseUrl string
	switch cfg.Scheme {
	case schemeHttp:
		baseUrl = fmt.Sprintf("%s://%s:%d", schemeHttp, cfg.Host, cfg.Port)
	case schemeHttps:
		baseUrl = fmt.Sprintf("%s://%s:%d", schemeHttps, cfg.Host, cfg.Port)
	case schemeHttpsUnsafe:
		baseUrl = fmt.Sprintf("%s://%s:%d", schemeHttps, cfg.Host, cfg.Port)
		cfg.TlsConfig.InsecureSkipVerify = true
	default:
		return nil, fmt.Errorf("scheme '%s' is not supported", cfg.Scheme)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: cfg.TlsConfig,
		},
	}

	return &Client{cfg: cfg, baseUrl: baseUrl, client: client, login: &userLoginResponse{}}, nil
}

// todo: need smarter handling of response codes, errors, errors in response body
func (o Client) get(url string, expectedResponseCodes []int, jsonPtr interface{}) error {
	if o.login.Token == "" {
		return fmt.Errorf("cannot interact with AOS API without token")
	}

	req, err := http.NewRequestWithContext(o.cfg.Ctx, "GET", url, nil)
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

	req, err := http.NewRequestWithContext(o.cfg.Ctx, "POST", url, bytes.NewBuffer(payload))
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
	return o.userLogin()
}

func (o Client) Logout() error {
	return o.userLogout()
}

func (o Client) GetStreamingConfigs() ([]GetStreamingConfigResponse, error) {
	return o.getAllStreamingConfigs()
}

func (o Client) GetVersion() (*VersionResponse, error) {
	return o.getVersion()
}
