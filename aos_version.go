package apstraTelemetry

import (
	"fmt"
)

const (
	aosApiVersion = "/api/version"
)

type AosVersionResponse struct {
	Major   string `json:"major"`
	Version string `json:"version"`
	Build   string `json:"build"`
	Minor   string `json:"minor"`
}

func (o AosClient) GetVersion() (*AosVersionResponse, error) {
	var versionResponse AosVersionResponse
	url := o.baseUrl + aosApiVersion
	err := o.get(url, []int{200}, &versionResponse)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &versionResponse, nil
}
