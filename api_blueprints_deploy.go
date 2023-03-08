package goapstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintDeploy = apiUrlBlueprintById + apiUrlPathDelim + "deploy"
)

type BlueprintDeploy struct {
	Id          ObjectId
	Description string
	Version     int
}

func (o *Client) deployBlueprint(ctx context.Context, in *BlueprintDeploy) error {
	deploy := &struct {
		Description string `json:"description"`
		Version     int    `json:"version"`
	}{
		Description: in.Description,
		Version:     in.Version,
	}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintDeploy, in.Id),
		apiInput: deploy,
	})
	return convertTtaeToAceWherePossible(err)
}
