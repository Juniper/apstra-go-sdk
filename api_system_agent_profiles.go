package goapstra

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiUrlSystemAgentProfiles     = "/api/system-agent-profiles"
	apiUrlSystemAgentProfilesById = apiUrlSystemAgentProfiles + apiUrlPathDelim + "%s"

	apstraAgentPlatformJunos = "junos"
	apstraAgentPlatformEOS   = "eos"
	apstraAgentPlatformNXOS  = "nxos"
)

type optionsSystemAgentProfilesResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type SystemAgentProfileConfig struct {
	Label       string            `json:"label"`
	Username    string            `json:"username,omitempty""`
	Password    string            `json:"password,omitempty"`
	Platform    string            `json:"platform,omitempty"`
	Packages    []string          `json:"packages"`
	OpenOptions map[string]string `json:"open_options"`
}

type getSystemAgentProfilesResponse struct {
	Items []SystemAgentProfile `json:"items"`
}

type SystemAgentProfile struct {
	Label       string            `json:"label"`
	HasUsername bool              `json:"has_username"`
	HasPassword bool              `json:"has_password"`
	Platform    string            `json:"platform"`
	Packages    []string          `json:"packages"`
	Id          ObjectId          `json:"id"`
	OpenOptions map[string]string `json:"open_options"`
}

func (o *Client) listSystemAgentProfileIds(ctx context.Context) ([]ObjectId, error) {
	method := http.MethodOptions
	urlStr := apiUrlSystemAgentProfiles
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &optionsSystemAgentProfilesResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' at '%s'", method, apstraUrl.String())
	}
	return response.Items, nil
}

func (o *Client) createSystemAgentProfile(ctx context.Context, in *SystemAgentProfileConfig) (ObjectId, error) {
	method := http.MethodPost
	urlStr := apiUrlSystemAgentProfiles
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("error parsing url '%s' - %w", apiUrlSystemAgentProfiles, err)
	}
	response := &objectIdResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", fmt.Errorf("error calling '%s' at '%s'", method, apstraUrl.String())
	}
	return response.Id, nil
}

func (o *Client) getSystemAgentProfile(ctx context.Context, id ObjectId) (*SystemAgentProfile, error) {
	method := http.MethodGet
	urlStr := fmt.Sprintf(apiUrlSystemAgentProfilesById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlSystemAgentProfiles, err)
	}
	response := &SystemAgentProfile{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		var ttae TalkToApstraErr
		if errors.As(err, &ttae) && ttae.Response.StatusCode == http.StatusNotFound {
			return nil, ApstraClientErr{
				errType: ErrNotfound,
				err:     err,
			}
		}
		return nil, fmt.Errorf("error calling '%s' at '%s'", method, apstraUrl.String())
	}
	return response, nil
}

func (o *Client) getAllSystemAgentProfiles(ctx context.Context) ([]SystemAgentProfile, error) {
	method := http.MethodGet
	urlStr := apiUrlSystemAgentProfiles
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &getSystemAgentProfilesResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' at '%s'", method, apstraUrl.String())
	}
	return response.Items, nil
}

func (o *Client) updateSystemAgentProfile(ctx context.Context, id ObjectId, in *SystemAgentProfileConfig) error {
	method := http.MethodPut
	urlStr := fmt.Sprintf(apiUrlSystemAgentProfilesById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:   method,
		url:      apstraUrl,
		apiInput: in,
	})
	if err != nil {
		return fmt.Errorf("error calling '%s' at '%s'", method, apstraUrl.String())
	}
	return nil
}

func (o *Client) deleteSystemAgentProfile(ctx context.Context, id ObjectId) error {
	method := http.MethodDelete
	urlStr := fmt.Sprintf(apiUrlSystemAgentProfilesById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlSystemAgentProfiles, err)
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method: method,
		url:    apstraUrl,
	})
	if err != nil {
		return fmt.Errorf("error calling '%s' at '%s'", method, apstraUrl.String())
	}
	return nil
}

func (o *Client) getSystemAgentProfileByLabel(ctx context.Context, label string) (*SystemAgentProfile, error) {
	method := http.MethodGet
	urlStr := apiUrlSystemAgentProfiles
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &getSystemAgentProfilesResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' at '%s'", method, apstraUrl.String())
	}

	found := -1 //slice index where the matching System Agent Profile can be found
	for i, sap := range response.Items {
		if sap.Label == label {
			if found >= 0 {
				return nil, fmt.Errorf("multiple matches for System Agent Profile with label '%s'", label)
			}
			found = i
		}
	}

	if found < 0 {
		return nil, ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no System Agent Profile with label '%s' found", label),
		}
	}
	return &response.Items[found], nil
}
