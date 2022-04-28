package apstraTelemetry

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	schemeHttp        = "http"
	schemeHttps       = "https"
	schemeHttpsUnsafe = "hxxps"

	aosApiLogin  = "/api/user/login"
	aosApiLogout = "/api/user/logout"
)

// AosClientCfg passed to NewAosClient() when instantiating a new AosClient{}
type AosClientCfg struct {
	Scheme string
	Host   string
	Port   uint16
	User   string
	Pass   string
}

// aosLoginReq payload to the aosApiLogin API endpoint
type aosLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// aosLoginResp payload returned by the aosApiLogin API endpoint
type aosLoginResp struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}

// AosClient interacts with an AOS API server
type AosClient struct {
	baseUrl string
	cfg     *AosClientCfg
	token   string
	client  *http.Client
}

// Login submits credentials to an API server, collects a login token
// todo - need to handle token timeout
func (o *AosClient) Login() (err error) {
	msg, err := json.Marshal(aosLoginReq{
		Username: o.cfg.User,
		Password: o.cfg.Pass,
	})
	if err != nil {
		return fmt.Errorf("error marshaling aosLogin object - %v", err)
	}

	req, err := http.NewRequest("POST", o.baseUrl+aosApiLogin, bytes.NewBuffer(msg))
	if err != nil {
		return fmt.Errorf("error creating http Request - %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http.client.Do - %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("http response code is not '%d' got '%d' at '%s'", 201, resp.StatusCode, aosApiLogin)
	}

	var loginResp *aosLoginResp
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return fmt.Errorf("error decoding aosLoginResp JSON - %v", err)
	}

	o.token = loginResp.Token

	return nil
}

func (o AosClient) Logout() error {
	req, err := http.NewRequest("POST", o.baseUrl+aosApiLogout, nil)
	if err != nil {
		return fmt.Errorf("error creating http Request - %v", err)
	}
	req.Header.Set("Authtoken", o.token)

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http.client.Do - %v", err)
	}
	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing logout http response body - %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("http response code is not '%d' got '%d' at '%s'", 200, resp.StatusCode, aosApiLogout)
	}

	return nil
}

// NewAosClient creates an AosClient object
func NewAosClient(cfg AosClientCfg) (*AosClient, error) {
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

	return &AosClient{cfg: &cfg, baseUrl: baseUrl, client: client}, nil
}
