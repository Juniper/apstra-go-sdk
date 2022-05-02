package aosSdk

import (
	"fmt"
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

func (o Client) getVersionsAosdi() (*versionsAosdiResponse, error) {
	var response versionsAosdiResponse
	err := o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           apiUrlVersionsAosdi,
		fromServerPtr: &response,
	})
	return &response, err
}

func (o Client) getVersionsApi() (*versionsApiResponse, error) {
	var response versionsApiResponse
	err := o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           apiUrlVersionsApi,
		fromServerPtr: &response,
	})
	return &response, err
}

func (o Client) getVersionsBuild() (*versionsBuildResponse, error) {
	var response versionsBuildResponse
	err := o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           apiUrlVersionsBuild,
		fromServerPtr: &response,
	})
	return &response, err
}

func (o Client) postVersionsDevice(request *versionsDeviceRequest) (*versionsDeviceResponse, error) {
	var response versionsDeviceResponse
	err := o.talkToAos(&talkToAosIn{
		method:      httpMethodPost,
		url:         apiUrlVersionsDevice,
		toServerPtr: request,
	})
	return &response, err
}

func (o Client) postVersionsIba(request *versionsIbaRequest) (*versionsIbaResponse, error) {
	var response versionsIbaResponse
	err := o.talkToAos(&talkToAosIn{
		method:        httpMethodPost,
		url:           apiUrlVersionsIba,
		toServerPtr:   request,
		fromServerPtr: &response,
	})
	return &response, err
}

func (o Client) postVersionsNode(request *versionsNodeRequest) (*versionsNodeResponse, error) {
	var response versionsNodeResponse
	err := o.talkToAos(&talkToAosIn{
		method:        httpMethodPost,
		url:           apiUrlVersionsNode,
		toServerPtr:   request,
		fromServerPtr: &response,
	})
	return &response, err
}

func (o Client) getVersionsServer() (*versionsServerResponse, error) {
	var response versionsServerResponse
	err := o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           apiUrlVersionsServer,
		fromServerPtr: &response,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", apiUrlVersionsServer, err)
	}
	return &response, nil
}
