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

func (o *TwoStageL3ClosClient) setAntiAffinityPolicy420(ctx context.Context, in *rawAntiAffinityPolicy) error {
	if in == nil {
		return nil
	}

	if !version.MustConstraints(version.NewConstraint("<=" + apstra420)).Check(o.client.apiVersion) {
		return fmt.Errorf("apstra %s does not support %q", o.client.apiVersion, apiUrlBlueprintAntiAffinityPolicy)
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintAntiAffinityPolicy, o.blueprintId),
		apiInput: &in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
