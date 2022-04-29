package aosSdk

import (
	"encoding/json"
	"fmt"
)

const (
	aosApiVersionsPrefix = "/api/versions/"
	aosApiVersionsAosdi  = aosApiVersionsPrefix + "aosdi"
	aosApiVersionsApi    = aosApiVersionsPrefix + "api"
	aosApiVersionsBuild  = aosApiVersionsPrefix + "build"
	aosApiVersionsDevice = aosApiVersionsPrefix + "device"
	aosApiVersionsIba    = aosApiVersionsPrefix + "iba"
	aosApiVersionsNode   = aosApiVersionsPrefix + "node"
	aosApiVersionsServer = aosApiVersionsPrefix + "server"
)

type aosApiVersionsAosdiResponse struct {
	Version       string `json:"version"`
	BuildDateTime string `json:"build_datetime"`
}

type aosApiVersionsApiResponse struct {
	Major   string `json:"major"`
	Version string `json:"version"`
	Build   string `json:"build"`
	Minor   string `json:"minor"`
}

type aosApiVersionsBuildResponse struct {
	Version       string `json:"version"`
	BuildDateTime string `json:"build_datetime"`
}

type aosApiVersionsDeviceRequest struct {
	SerialNumber string `json:"serial_number"`
	Version      string `json:"version"`
	Platform     string `json:"platform"`
}

type aosApiVersionsDeviceResponse struct {
	Status       string `json:"status"`
	Url          string `json:"url"`
	RetryTimeout int    `json:"retry_timeout"`
	Cksum        string `json:"cksum"`
}

type aosApiVersionsIbaRequest struct {
	Version  string `json:"version""`
	SystemId string `json:"system_id""`
}

type aosApiVersionsIbaResponse struct {
	Status       string `json:"status"`
	Url          string `json:"url"`
	RetryTimeout int    `json:"retry_timeout"`
	Cksum        string `json:"cksum"`
}

type aosApiVersionsNodeRequest struct {
	IpAddress string `json:"ip_address"`
	Version   string `json:"version"`
	SystemId  string `json:"system_id"`
}

type aosApiVersionsNodeResponse struct {
	Status       string `json:"status"`
	RetryTimeout int    `json:"retry_timeout"`
}

type aosApiVersionsServerResponse struct {
	Version       string `json:"version"`
	BuildDateTime string `json:"build_datetime"`
}

func (o AosClient) getVersionsAosdi() (*aosApiVersionsAosdiResponse, error) {
	var response aosApiVersionsAosdiResponse
	url := aosApiVersionsAosdi
	err := o.get(url, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o AosClient) getVersionsApi() (*aosApiVersionsApiResponse, error) {
	var response aosApiVersionsApiResponse
	url := aosApiVersionsApi
	err := o.get(url, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o AosClient) getVersionsBuild() (*aosApiVersionsBuildResponse, error) {
	var response aosApiVersionsBuildResponse
	url := aosApiVersionsApi
	err := o.get(url, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o AosClient) postVersionsDevice(request *aosApiVersionsDeviceRequest) (*aosApiVersionsDeviceResponse, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling aosApiVersionsDeviceRequest object - %v", err)
	}
	var response aosApiVersionsDeviceResponse
	url := aosApiVersionsDevice
	err = o.post(url, payload, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o AosClient) postVersionsIba(request *aosApiVersionsIbaRequest) (*aosApiVersionsIbaResponse, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling aosApiVersionsDeviceRequest object - %v", err)
	}
	var response aosApiVersionsIbaResponse
	url := aosApiVersionsIba
	err = o.post(url, payload, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o AosClient) postVersionsNode(request *aosApiVersionsNodeRequest) (*aosApiVersionsNodeResponse, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling aosApiVersionsDeviceRequest object - %v", err)
	}
	var response aosApiVersionsNodeResponse
	url := aosApiVersionsIba
	err = o.post(url, payload, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o AosClient) getVersionsServer() (*aosApiVersionsServerResponse, error) {
	var response aosApiVersionsServerResponse
	url := aosApiVersionsIba
	err := o.get(url, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}
