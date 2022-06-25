package goapstra

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	apiUrlSystemAgentProfiles     = "/api/system-agent-profiles"
	apiUrlSystemAgentProfilesById = apiUrlSystemAgentProfiles + apiUrlPathDelim + "%s"

	apstraAgentPlatformJunos = "junos"
	apstraAgentPlatformEOS   = "eos"
	apstraAgentPlatformNXOS  = "nxos"

	apstraSystemAgentPlatformStringSep = "=="
)

type optionsSystemAgentProfilesResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

// SystemAgentProfileConfig is used when creating or updating an Agent Profile
type SystemAgentProfileConfig struct {
	Label       string            `json:"label"`
	Username    string            `json:"username,omitempty""`
	Password    string            `json:"password,omitempty"`
	Platform    string            `json:"platform,omitempty"`
	Packages    map[string]string `json:"packages"`
	OpenOptions map[string]string `json:"open_options"`
}

// raw turns a SystemAgentProfile (from our caller) into a rawSystemAgentProfile
func (o *SystemAgentProfileConfig) raw(id ObjectId) *rawSystemAgentProfileConfig {
	//goland:noinspection GoPreferNilSlice
	packages := []string{}
	for k, v := range o.Packages {
		packages = append(packages, k+apstraSystemAgentPlatformStringSep+v)
	}
	result := &rawSystemAgentProfileConfig{
		Id:          string(id),
		Label:       o.Label,
		Username:    o.Username,
		Password:    o.Password,
		Platform:    o.Platform,
		Packages:    packages,
		OpenOptions: o.OpenOptions,
	}
	if result.OpenOptions == nil { // this would result in 'null' in JSON payload
		result.OpenOptions = make(map[string]string) // send '{}' instead
	}
	return result
}

// rawSystemAgentProfileConfig is the nasty type expected by the API. Element
// Packages is really a map, but k,v are string-joined with "==" here.
type rawSystemAgentProfileConfig struct {
	Id          string            `json:"id,omitempty"`
	Label       string            `json:"label"`
	Username    string            `json:"username,omitempty""`
	Password    string            `json:"password,omitempty"`
	Platform    string            `json:"platform,omitempty"`
	Packages    []string          `json:"packages"`
	OpenOptions map[string]string `json:"open_options"`
	Profile     interface{}       `json:"profile"` // this exists to be 'null' json - don't know why
}

type getSystemAgentProfilesResponse struct {
	Items []rawSystemAgentProfile `json:"items"`
}

// SystemAgentProfile describes an Agent Profile to our callers.
// It has the Packages element presented sensibly as a map.
type SystemAgentProfile struct {
	Label       string            `json:"label"`
	HasUsername bool              `json:"has_username"`
	HasPassword bool              `json:"has_password"`
	Platform    string            `json:"platform"`
	Packages    map[string]string `json:"packages"`
	Id          ObjectId          `json:"id"`
	OpenOptions map[string]string `json:"open_options"`
}

// rawSystemAgentProfile represents the API's description of an Agent Profile.
// The Packages element is really a map, but has k, v string-joined with "==".
type rawSystemAgentProfile struct {
	Label       string            `json:"label"`
	HasUsername bool              `json:"has_username"`
	HasPassword bool              `json:"has_password"`
	Platform    string            `json:"platform"`
	Packages    []string          `json:"packages"`
	Id          ObjectId          `json:"id"`
	OpenOptions map[string]string `json:"open_options"`
}

// polish turns a rawSystemAgentProfile (from the API) into a SystemAgentProfile
func (o *rawSystemAgentProfile) polish() *SystemAgentProfile {
	var packages map[string]string
	if len(o.Packages) > 0 {
		packages = make(map[string]string)
	}
	for _, s := range o.Packages {
		kv := strings.SplitN(s, apstraSystemAgentPlatformStringSep, 2)
		switch len(kv) {
		case 2:
			packages[kv[0]] = kv[1]
		case 1:
			packages[kv[0]] = ""
		}
	}
	return &SystemAgentProfile{
		Label:       o.Label,
		HasUsername: o.HasUsername,
		HasPassword: o.HasPassword,
		Platform:    o.Platform,
		Packages:    packages,
		Id:          o.Id,
		OpenOptions: o.OpenOptions,
	}
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
		apiInput:    in.raw(""),
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

func (o *Client) getSystemAgentProfile(ctx context.Context, id ObjectId) (*SystemAgentProfile, error) {
	method := http.MethodGet
	urlStr := fmt.Sprintf(apiUrlSystemAgentProfilesById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlSystemAgentProfiles, err)
	}
	response := &rawSystemAgentProfile{}
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

	var out []SystemAgentProfile
	for _, sap := range response.Items {
		out = append(out, *sap.polish())
	}
	return out, nil
}

func (o *Client) updateSystemAgentProfile(ctx context.Context, id ObjectId, in *SystemAgentProfileConfig) error {
	method := http.MethodPatch
	urlStr := fmt.Sprintf(apiUrlSystemAgentProfilesById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}

	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:   method,
		url:      apstraUrl,
		apiInput: in.raw(id),
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
	return response.Items[found].polish(), nil
}
