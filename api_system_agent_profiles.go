package goapstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	apiUrlSystemAgentProfiles           = "/api/system-agent-profiles"
	apiUrlSystemAgentProfilesById       = apiUrlSystemAgentProfiles + apiUrlPathDelim + "%s"
	apiUrlSystemAgentProfilesAssignById = apiUrlSystemAgentProfiles + apiUrlPathDelim + "%s" + "/assign"

	apstraAgentPlatformJunos = "junos"
	apstraAgentPlatformEOS   = "eos"
	apstraAgentPlatformNXOS  = "nxos"

	apstraSystemAgentPlatformStringSep = "=="

	apstraErrAgentProfileInUse = "Profile is in use"
)

type optionsAgentProfilesResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type AssignAgentProfileRequest struct {
	SystemAgents     []ObjectId `json:"system_agents"`
	ProfileId        ObjectId   `json:"profile_id"`
	ClearPackages    bool       `json:"clear_packages"`
	ClearOpenOptions bool       `json:"clear_open_options"`
}

// AgentProfileConfig is used when creating or updating an Agent Profile
type AgentProfileConfig struct {
	Label       string
	Username    string
	Password    string
	Platform    string
	Packages    AgentPackages
	OpenOptions map[string]string
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
	Username    string            `json:"username,omitempty"`
	Password    string            `json:"password,omitempty"`
	Platform    string            `json:"platform"`
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
	response := &optionsAgentProfilesResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlSystemAgentProfiles,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) createAgentProfile(ctx context.Context, in *AgentProfileConfig) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlSystemAgentProfiles,
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
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) getAgentProfile(ctx context.Context, id ObjectId) (*AgentProfile, error) {
	response := &rawAgentProfile{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlSystemAgentProfilesById, id),
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
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.polish(), nil
}

func (o *Client) getAllAgentProfiles(ctx context.Context) ([]AgentProfile, error) {
	response := &getAgentProfilesResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlSystemAgentProfiles,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	var out []AgentProfile
	for _, sap := range response.Items {
		out = append(out, *sap.polish())
	}
	return out, nil
}

func (o *Client) updateAgentProfile(ctx context.Context, id ObjectId, in *AgentProfileConfig) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlSystemAgentProfilesById, id),
		apiInput: in.raw(),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) deleteAgentProfile(ctx context.Context, id ObjectId) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlSystemAgentProfilesById, id),
	})
	if err != nil {
		var ttae TalkToApstraErr
		if errors.As(err, &ttae) && ttae.Response.StatusCode == http.StatusUnprocessableEntity {
			body, _ := io.ReadAll(ttae.Response.Body)
			var ae apstraErr
			_ = json.Unmarshal(body, &apstraErr{})
			if ae.Errors == apstraErrAgentProfileInUse {
				return ApstraClientErr{
					errType: ErrInUse,
					err:     fmt.Errorf("agent profile '%s' is in use, cannot delete", id),
				}
			}
		}
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) getAgentProfileByLabel(ctx context.Context, label string) (*AgentProfile, error) {
	response := &getAgentProfilesResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlSystemAgentProfiles,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	found := -1 //slice index where the matching System Agent Profile can be found
	for i, sap := range response.Items {
		if sap.Label == label {
			if found >= 0 {
				return nil, ApstraClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple matches for System Agent Profile with label '%s'", label),
				}
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

func (o *Client) assignAgentProfile(ctx context.Context, req *AssignAgentProfileRequest) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPost,
		urlStr:   apiUrlSystemAgentProfiles,
		apiInput: req,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
