package aosSdk

import (
	"fmt"
	"net/url"
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
	url, err := url.Parse(apiUrlVersion)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersion, err)
	}
	var response VersionResponse
	_, err = o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           url,
		fromServerPtr: &response,
	})
	return &response, err
}
