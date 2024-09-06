package apstra

import (
	"context"
	"fmt"
	"net/http"
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

func (o *Client) getVersion(ctx context.Context) (*VersionResponse, error) {
	apstraUrl, err := url.Parse(apiUrlVersion)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersion, err)
	}

	response := &VersionResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
}
