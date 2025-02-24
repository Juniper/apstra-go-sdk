package apstra

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func (o *Client) talkToApiOps(id string, apstraReq *http.Request) (*http.Response, error) {
	if id == "" {
		return nil, ClientErr{
			errType: ErrInvalidId,
			err:     fmt.Errorf("API-ops ID must not be empty"),
		}
	}

	// Based on:
	// https://ssd-git.juniper.net/aide-jcloud/apstra-marvis/api-ops/-/blob/develop/proto/aospayload.pb.go#L106-125
	var proxyBody struct {
		// HTTP Methos GET, PUT, POST ... if not set correctly then request would fail
		Method string `json:"method"`
		// url path is the AOS API URL, must not contain protocol, scheme IP/Host and port
		UrlPath string `json:"urlPath"`
		// request body for POST, PUT, PATCH requests
		Body []byte `json:"body,omitempty"`
		// request parameters
		Params map[string]string `json:"params,omitempty"`
		// request headers
		Headers map[string]string `json:"headers"`
	}

	params := apstraReq.URL.Query()

	var err error

	proxyBody.Method = apstraReq.Method
	proxyBody.UrlPath = apstraReq.URL.Path
	proxyBody.Body, err = io.ReadAll(apstraReq.Body)
	if err != nil {
		return nil, fmt.Errorf("error prepping proxy request body from apstraReq: %w", err)
	}
	proxyBody.Params = make(map[string]string, len(params))
	for k, v := range params {
		switch len(v) {
		case 0:
		case 1:
			proxyBody.Params[k] = v[0]
		default:
			return nil, fmt.Errorf("cannot format query string param %q for the proxy: only one string supported per param, got %d strings", k, len(v))
		}
	}
	proxyBody.Headers = make(map[string]string, len(apstraReq.Header))
	for k, v := range apstraReq.Header {
		switch len(v) {
		case 0:
		case 1:
			proxyBody.Headers[k] = v[0]
		default:
			return nil, fmt.Errorf("cannot format query string header %q for the proxy: only one string supported per header, got %d strings", k, len(v))
		}
	}

	proxyRequest, err := http.NewRequest(apstraReq.Method, o.baseUrl.String(), bytes.NewReader(proxyBody.Body))
	if err != nil {
		return nil, fmt.Errorf("error prepping proxyRequest from apstraReq: %w", err)
	}

	proxyRequest = proxyRequest.WithContext(apstraReq.Context())

	return o.httpClient.Do(proxyRequest)
}
