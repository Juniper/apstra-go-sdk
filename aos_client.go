package apstraTelemetry

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
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
