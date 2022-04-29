package aosSdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
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

// AosClientCfg passed to NewAosClient() when instantiating a new AosClient{}
type AosClientCfg struct {
	Scheme string
	Host   string
	Port   uint16
	User   string
	Pass   string
}

// AosClient interacts with an AOS API server
type AosClient struct {
	baseUrl string
	cfg     *AosClientCfg
	token   string
	client  *http.Client
}

// NewAosClient creates an AosClient object
func NewAosClient(cfg *AosClientCfg) (*AosClient, error) {
	tlsConfig := &tls.Config{}
	var baseUrl string
	switch cfg.Scheme {
	case schemeHttp:
		baseUrl = fmt.Sprintf("%s://%s:%d", schemeHttp, cfg.Host, cfg.Port)
	case schemeHttps:
		baseUrl = fmt.Sprintf("%s://%s:%d", schemeHttps, cfg.Host, cfg.Port)
	case schemeHttpsUnsafe:
		baseUrl = fmt.Sprintf("%s://%s:%d", schemeHttps, cfg.Host, cfg.Port)
		tlsConfig.InsecureSkipVerify = true
	default:
		return nil, fmt.Errorf("scheme '%s' is not supported", cfg.Scheme)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return &AosClient{cfg: cfg, baseUrl: baseUrl, client: client}, nil
}

func (o AosClient) get(url string, expectedResponseCodes []int, jsonPtr interface{}) error {
	if o.token == "" {
		return fmt.Errorf("cannot interact with AOS API without token")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating http Request - %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authtoken", o.token)

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http.client.Do - %v", err)
	}
	defer resp.Body.Close()

	if !intSliceContains(expectedResponseCodes, resp.StatusCode) {
		return fmt.Errorf("unexpected http response code '%d' (permitted: '%s') at '%s'",
			resp.StatusCode, strings.Join(intSliceToStringSlice(expectedResponseCodes), ","), url)
	}

	if jsonPtr != nil {
		return json.NewDecoder(resp.Body).Decode(jsonPtr)
	}

	return nil
}

func (o AosClient) post(url string, payload []byte, expectedResponseCodes []int, jsonPtr interface{}) error {
	if o.token == "" {
		return fmt.Errorf("cannot interact with AOS API without token")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating http Request - %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authtoken", o.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http.client.Do - %v", err)
	}
	defer resp.Body.Close()

	if !intSliceContains(expectedResponseCodes, resp.StatusCode) {
		return fmt.Errorf("unexpected http response code '%d' (permitted: '%s') at '%s'",
			resp.StatusCode, strings.Join(intSliceToStringSlice(expectedResponseCodes), ","), url)
	}

	if jsonPtr != nil {
		return json.NewDecoder(resp.Body).Decode(jsonPtr)
	}

	return nil
}

func (o *AosClient) Login() error {
	return o.userLogin()
}

func (o AosClient) Logout() error {
	return o.userLogout()
}

func (o AosClient) GetStreamingConfigs() ([]*AosStreamingConfig, error) {
	return o.getAllStreamingConfigs()
}

func (o AosClient) GetVersion() (*AosVersionResponse, error) {
	return o.getVersion()
}
