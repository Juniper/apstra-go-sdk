package aosSdk

import (
	"encoding/json"
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
	Version  string `json:"version""`
	SystemId string `json:"system_id""`
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
	url := apiUrlVersionsAosdi
	err := o.get(url, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o Client) getVersionsApi() (*versionsApiResponse, error) {
	var response versionsApiResponse
	url := apiUrlVersionsApi
	err := o.get(url, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o Client) getVersionsBuild() (*versionsBuildResponse, error) {
	var response versionsBuildResponse
	url := apiUrlVersionsApi
	err := o.get(url, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o Client) postVersionsDevice(request *versionsDeviceRequest) (*versionsDeviceResponse, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling versionsDeviceRequest object - %v", err)
	}
	var response versionsDeviceResponse
	url := apiUrlVersionsDevice
	err = o.post(url, payload, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o Client) postVersionsIba(request *versionsIbaRequest) (*versionsIbaResponse, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling versionsDeviceRequest object - %v", err)
	}
	var response versionsIbaResponse
	url := apiUrlVersionsIba
	err = o.post(url, payload, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o Client) postVersionsNode(request *versionsNodeRequest) (*versionsNodeResponse, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling versionsDeviceRequest object - %v", err)
	}
	var response versionsNodeResponse
	url := apiUrlVersionsIba
	err = o.post(url, payload, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}

func (o Client) getVersionsServer() (*versionsServerResponse, error) {
	var response versionsServerResponse
	url := apiUrlVersionsIba
	err := o.get(url, []int{200}, &response)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &response, nil
}
