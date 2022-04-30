package aosSdk

import (
	"fmt"
)

const (
	apiUrlVersion = "/api/version"
)

type VersionResponse struct {
	Major   string `json:"major"`
	Version string `json:"version"`
	Build   string `json:"build"`
	Minor   string `json:"minor"`
}

func (o Client) getVersion() (*VersionResponse, error) {
	var versionResponse VersionResponse
	url := o.baseUrl + apiUrlVersion
	err := o.get(url, []int{200}, &versionResponse)
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' - %v", url, err)
	}
	return &versionResponse, nil
}