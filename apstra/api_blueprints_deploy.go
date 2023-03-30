package apstra

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	apiUrlBlueprintDeploy    = apiUrlBlueprintById + apiUrlPathDelim + "deploy"
	apiUrlBlueprintRevisions = apiUrlBlueprintById + apiUrlPathDelim + "revisions"
)

type DeployStatus int
type deployStatus string

const (
	DeployStatusSuccess = DeployStatus(iota)
	DeployStatusFailure
	DeployStatusUnknown = "unknown deploy status '%s'"

	deployStatusSuccess = deployStatus("success")
	deployStatusFailure = deployStatus("failure")
	deployStatusUnknown = "unknown redundancy protocol '%d'"
)

func (o DeployStatus) Int() int {
	return int(o)
}

func (o DeployStatus) String() string {
	switch o {
	case DeployStatusSuccess:
		return string(deployStatusSuccess)
	case DeployStatusFailure:
		return string(deployStatusFailure)
	default:
		return fmt.Sprintf(deployStatusUnknown, o)
	}
}

func (o deployStatus) string() string {
	return string(o)
}

func (o deployStatus) parse() (int, error) {
	switch o {
	case deployStatusSuccess:
		return int(DeployStatusSuccess), nil
	case deployStatusFailure:
		return int(DeployStatusFailure), nil
	default:
		return 0, fmt.Errorf(DeployStatusUnknown, o)
	}
}

type BlueprintDeployRequest struct {
	Id          ObjectId
	Description string
	Version     int
}

type rawBlueprintDeployResponse struct {
	Status  deployStatus `json:"state"`
	Version int          `json:"version"`
	Error   *string      `json:"error,omitempty"`
}

func (o *rawBlueprintDeployResponse) polish() (*BlueprintDeployResponse, error) {
	status, err := o.Status.parse()
	if err != nil {
		return nil, err
	}

	return &BlueprintDeployResponse{
		Status:  DeployStatus(status),
		Version: o.Version,
		Error:   o.Error,
	}, nil
}

type BlueprintDeployResponse struct {
	Status  DeployStatus `json:"state"`
	Version int          `json:"version"`
	Error   *string      `json:"error,omitempty"`
}

func (o *Client) deployBlueprint(ctx context.Context, in *BlueprintDeployRequest) (*rawBlueprintDeployResponse, error) {
	request := &struct {
		Description string `json:"description"`
		Version     int    `json:"version"`
	}{
		Description: in.Description,
		Version:     in.Version,
	}

	url, err := url.Parse(fmt.Sprintf(apiUrlBlueprintDeploy, in.Id))
	if err != nil {
		return nil, err
	}

	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		url:      url,
		apiInput: request,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	response := &rawBlueprintDeployResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         url,
		apiResponse: &response,
	})

	return response, convertTtaeToAceWherePossible(err)
}

type rawBlueprintRevision struct {
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	AutoSaved   bool      `json:"auto_saved"`
	UserSaved   bool      `json:"user_saved"`
	UserIp      string    `json:"user_ip"`
	User        string    `json:"user"`
	RevisionId  string    `json:"revision_id"`
}

func (o *rawBlueprintRevision) polish() (*BlueprintRevision, error) {
	revisionId, err := strconv.Atoi(o.RevisionId)
	if err != nil {
		return nil, fmt.Errorf("error parsing blueprint revision %q to integer - %w",
			o.RevisionId, err)
	}

	return &BlueprintRevision{
		Description: o.Description,
		CreatedAt:   o.CreatedAt,
		AutoSaved:   o.AutoSaved,
		UserSaved:   o.UserSaved,
		UserIp:      net.ParseIP(o.UserIp),
		User:        o.User,
		RevisionId:  revisionId,
	}, nil
}

type BlueprintRevision struct {
	Description string
	CreatedAt   time.Time
	AutoSaved   bool
	UserSaved   bool
	UserIp      net.IP
	User        string
	RevisionId  int
}

func (o *Client) getBlueprintRevisions(ctx context.Context, id ObjectId) ([]rawBlueprintRevision, error) {
	result := &struct {
		Items []rawBlueprintRevision `json:"items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintRevisions, id),
		apiResponse: result,
	})
	return result.Items, convertTtaeToAceWherePossible(err)
}
