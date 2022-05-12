package apstra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiUrlVersionsPrefix = "/api/versions/"
	apiUrlVersionsAosdi  = apiUrlVersionsPrefix + "aosdi"
	apiUrlVersionsApi    = apiUrlVersionsPrefix + "api"
	apiUrlVersionsBuild  = apiUrlVersionsPrefix + "build"
	apiUrlVersionsDevice = apiUrlVersionsPrefix + "device"
	apiUrlVersionsIba    = apiUrlVersionsPrefix + "iba"
	apiUrlVersionsNode   = apiUrlVersionsPrefix + "node"
	apiUrlVersionsServer = apiUrlVersionsPrefix + "server"
)

type versionsAosdiResponse struct {
	Version       string `json:"version"`
	BuildDateTime string `json:"build_datetime"`
}

type versionsApiResponse struct {
	Major   string `json:"major"`
	Version string `json:"version"`
	Build   string `json:"build"`
	Minor   string `json:"minor"`
}

type versionsBuildResponse struct {
	Version       string `json:"version"`
	BuildDateTime string `json:"build_datetime"`
}

type versionsDeviceRequest struct {
	SerialNumber string `json:"serial_number"`
	Version      string `json:"version"`
	Platform     string `json:"platform"`
}

type versionsDeviceResponse struct {
	Status       string `json:"status"`
	Url          string `json:"url"`
	RetryTimeout int    `json:"retry_timeout"`
	Cksum        string `json:"cksum"`
}

type versionsIbaRequest struct {
	Version  string `json:"version"`
	SystemId string `json:"system_id"`
}

type versionsIbaResponse struct {
	Status       string `json:"status"`
	Url          string `json:"url"`
	RetryTimeout int    `json:"retry_timeout"`
	Cksum        string `json:"cksum"`
}

type versionsNodeRequest struct {
	IpAddress string `json:"ip_address"`
	Version   string `json:"version"`
	SystemId  string `json:"system_id"`
}

type versionsNodeResponse struct {
	Status       string `json:"status"`
	RetryTimeout int    `json:"retry_timeout"`
}

type versionsServerResponse struct {
	Version       string `json:"version"`
	BuildDateTime string `json:"build_datetime"`
}

func (o Client) getVersionsAosdi(ctx context.Context) (*versionsAosdiResponse, error) {
	apstraUrl, err := url.Parse(apiUrlVersionsAosdi)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersionsAosdi, err)
	}
	response := &versionsAosdiResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
}

func (o Client) getVersionsApi(ctx context.Context) (*versionsApiResponse, error) {
	apstraUrl, err := url.Parse(apiUrlVersionsApi)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersionsApi, err)
	}
	response := &versionsApiResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
}

func (o Client) getVersionsBuild(ctx context.Context) (*versionsBuildResponse, error) {
	apstraUrl, err := url.Parse(apiUrlVersionsBuild)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersionsBuild, err)
	}
	response := &versionsBuildResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
}

func (o Client) postVersionsDevice(ctx context.Context, request *versionsDeviceRequest) (*versionsDeviceResponse, error) {
	apstraUrl, err := url.Parse(apiUrlVersionsDevice)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersionsDevice, err)
	}
	response := &versionsDeviceResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		url:         apstraUrl,
		apiInput:    request,
		apiResponse: response,
	})
}

func (o Client) postVersionsIba(ctx context.Context, request *versionsIbaRequest) (*versionsIbaResponse, error) {
	apstraUrl, err := url.Parse(apiUrlVersionsIba)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersionsIba, err)
	}
	response := &versionsIbaResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		url:         apstraUrl,
		apiInput:    request,
		apiResponse: response,
	})
}

func (o Client) postVersionsNode(ctx context.Context, request *versionsNodeRequest) (*versionsNodeResponse, error) {
	apstraUrl, err := url.Parse(apiUrlVersionsNode)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersionsNode, err)
	}
	response := &versionsNodeResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		url:         apstraUrl,
		apiInput:    request,
		apiResponse: response,
	})
}

func (o Client) getVersionsServer(ctx context.Context) (*versionsServerResponse, error) {
	apstraUrl, err := url.Parse(apiUrlVersionsServer)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersionsServer, err)
	}
	response := &versionsServerResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
}
