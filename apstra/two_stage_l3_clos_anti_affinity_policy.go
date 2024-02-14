package apstra

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-version"
	"net/http"
)

const (
	apiUrlBlueprintAntiAffinityPolicy = apiUrlBlueprintByIdPrefix + "anti-affinity-policy"
)

func (o *TwoStageL3ClosClient) getAntiAffinityPolicy420(ctx context.Context) (*rawAntiAffinityPolicy, error) {
	if !version.MustConstraints(version.NewConstraint("<=" + apstra420)).Check(o.client.apiVersion) {
		return nil, fmt.Errorf("apstra %s does not support %q", o.client.apiVersion, apiUrlBlueprintAntiAffinityPolicy)
	}

	var result rawAntiAffinityPolicy
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintAntiAffinityPolicy, o.blueprintId),
		apiResponse: &result,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &result, nil
}
