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
	err := o.get(o.baseUrl+aosApiVersion, []int{200}, &versionResponse)
	if err != nil {
		return nil, fmt.Errorf("error calling AosClient.get() - %v", err)
	}
	return &versionResponse, nil
}
