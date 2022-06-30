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

	apstraSystemAgentPlatformStringSep = "=="
)

type optionsAgentProfilesResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

// AgentProfileConfig is used when creating or updating an Agent Profile
type AgentProfileConfig struct {
	Label       string            `json:"label"`
	Username    string            `json:"username,omitempty""`
	Password    string            `json:"password,omitempty"`
	Platform    string            `json:"platform,omitempty"`
	Packages    AgentPackages     `json:"packages"`
	OpenOptions map[string]string `json:"open_options"`
}

// raw turns a *AgentProfile (from our caller) into a rawAgentProfile
func (o *AgentProfileConfig) raw() *rawAgentProfileConfig {
	result := &rawAgentProfileConfig{
		Label:       o.Label,
		Username:    o.Username,
		Password:    o.Password,
		Platform:    o.Platform,
		Packages:    o.Packages.raw(),
		OpenOptions: o.OpenOptions,
	}
	if result.OpenOptions == nil { // this would result in 'null' in JSON payload
		result.OpenOptions = make(map[string]string) // send '{}' instead
	}
	return result
}

// rawAgentProfileConfig is the nasty type expected by the API. Element
// Packages is really a map, but k,v are string-joined with "==" here.
type rawAgentProfileConfig struct {
	Label       string            `json:"label"`
	Username    string            `json:"username,omitempty""`
	Password    string            `json:"password,omitempty"`
	Platform    string            `json:"platform,omitempty"`
	Packages    rawAgentPackages  `json:"packages"`
	OpenOptions map[string]string `json:"open_options"`
}

type getAgentProfilesResponse struct {
	Items []rawAgentProfile `json:"items"`
}

// AgentProfile describes an Agent Profile to our callers.
// It has the Packages element presented sensibly as a map.
type AgentProfile struct {
	Label       string            `json:"label"`
	HasUsername bool              `json:"has_username"`
	HasPassword bool              `json:"has_password"`
	Platform    string            `json:"platform"`
	Packages    AgentPackages     `json:"packages"`
	Id          ObjectId          `json:"id"`
	OpenOptions map[string]string `json:"open_options"`
}

// rawAgentProfile represents the API's description of an Agent Profile.
// The Packages element is really a map, but has k, v string-joined with "==".
type rawAgentProfile struct {
	Label       string            `json:"label"`
	HasUsername bool              `json:"has_username"`
	HasPassword bool              `json:"has_password"`
	Platform    string            `json:"platform"`
	Packages    rawAgentPackages  `json:"packages"`
	Id          ObjectId          `json:"id"`
	OpenOptions map[string]string `json:"open_options"`
}

// polish turns a rawAgentProfile (from the API) into a AgentProfile
func (o *rawAgentProfile) polish() *AgentProfile {
	return &AgentProfile{
		Label:       o.Label,
		HasUsername: o.HasUsername,
		HasPassword: o.HasPassword,
		Platform:    o.Platform,
		Packages:    o.Packages.polish(),
		Id:          o.Id,
		OpenOptions: o.OpenOptions,
	}
}

func (o *Client) listAgentProfileIds(ctx context.Context) ([]ObjectId, error) {
	method := http.MethodOptions
	urlStr := apiUrlSystemAgentProfiles
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &optionsAgentProfilesResponse{}
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

func (o *Client) createAgentProfile(ctx context.Context, in *AgentProfileConfig) (ObjectId, error) {
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
		apiInput:    in.raw(),
		apiResponse: response,
	})
	if err != nil {
		var ttae TalkToApstraErr
		if errors.As(err, &ttae) && ttae.Response.StatusCode == http.StatusConflict {
			return "", ApstraClientErr{
				errType: ErrConflict,
				err:     fmt.Errorf("error Agent Profile '%s' likely already exists - %w", in.Label, err),
			}
		}
		return "", fmt.Errorf("error calling '%s' at '%s'", method, apstraUrl.String())
	}
	return response.Id, nil
}

func (o *Client) getAgentProfile(ctx context.Context, id ObjectId) (*AgentProfile, error) {
	method := http.MethodGet
	urlStr := fmt.Sprintf(apiUrlSystemAgentProfilesById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlSystemAgentProfiles, err)
	}
	response := &rawAgentProfile{}
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
	return response.polish(), nil
}

func (o *Client) getAllAgentProfiles(ctx context.Context) ([]AgentProfile, error) {
	method := http.MethodGet
	urlStr := apiUrlSystemAgentProfiles
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &getAgentProfilesResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' at '%s'", method, apstraUrl.String())
	}

	var out []AgentProfile
	for _, sap := range response.Items {
		out = append(out, *sap.polish())
	}
	return out, nil
}

func (o *Client) updateAgentProfile(ctx context.Context, id ObjectId, in *AgentProfileConfig) error {
	method := http.MethodPatch
	urlStr := fmt.Sprintf(apiUrlSystemAgentProfilesById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}

	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:   method,
		url:      apstraUrl,
		apiInput: in.raw(),
	})
	if err != nil {
		return fmt.Errorf("error calling '%s' at '%s'", method, apstraUrl.String())
	}
	return nil
}

func (o *Client) deleteAgentProfile(ctx context.Context, id ObjectId) error {
	method := http.MethodDelete
	urlStr := fmt.Sprintf(apiUrlSystemAgentProfilesById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
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

func (o *Client) getAgentProfileByLabel(ctx context.Context, label string) (*AgentProfile, error) {
	method := http.MethodGet
	urlStr := apiUrlSystemAgentProfiles
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &getAgentProfilesResponse{}
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
	return response.Items[found].polish(), nil
}
