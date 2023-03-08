package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiUrlBlueprintDeploy = apiUrlBlueprintById + apiUrlPathDelim + "deploy"
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

func (o DeployStatus) raw() deployStatus {
	return deployStatus(o.String())
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

//func (o *DeployStatus) FromString(in string) error {
//	i, err := deployStatus(in).parse()
//	if err != nil {
//		return err
//	}
//	*o = DeployStatus(i)
//	return nil
//}

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
