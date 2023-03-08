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

type BlueprintDeployRequest struct {
	Id          ObjectId
	Description string
	Version     int
}

type BlueprintDeployResponse struct {
	State   string `json:"state"`
	Version int    `json:"version"`
	Error   string `json:"error"`
}

func (o *Client) deployBlueprint(ctx context.Context, in *BlueprintDeployRequest) (*BlueprintDeployResponse, error) {
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

	response := &BlueprintDeployResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         url,
		apiResponse: &response,
	})

	return response, convertTtaeToAceWherePossible(err)
}
