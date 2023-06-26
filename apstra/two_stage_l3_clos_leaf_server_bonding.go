package apstra

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

const (
	apiUrlLeafServerLinkLabels = apiUrlBlueprintById + apiUrlPathDelim + "leaf-server-link-labels"
)

type LinkLagParams struct {
	GroupLabel string
	LagMode    RackLinkLagMode
	Tags       []string
}

func (o *LinkLagParams) raw() (*rawLinkLagParams, error) {
	groupLabel := o.GroupLabel
	if groupLabel == "" {
		initUUID()
		uuid1, err := uuid.NewUUID()
		if err != nil {
			return nil, fmt.Errorf("error generating type 1 uuid - %w", err)
		}
		groupLabel = uuid1.String()
	}

	return &rawLinkLagParams{
		GroupLabel: groupLabel,
		LagMode:    rackLinkLagMode(o.LagMode.String()),
		Tags:       o.Tags,
	}, nil
}

type rawLinkLagParams struct {
	GroupLabel string          `json:"group_label"`
	LagMode    rackLinkLagMode `json:"lag_mode,omitempty"`
	Tags       []string        `json:"tags"`
}

// SetLinkLagParamsRequest is a map of LAG parameters keyed by link node ID
type SetLinkLagParamsRequest map[ObjectId]LinkLagParams

// SetLinkLagParams configures the links identified in the request
// Links with no supplied GroupLabel will be given a unique random label making
// them the only members of their own group.
func (o *TwoStageL3ClosClient) SetLinkLagParams(ctx context.Context, req *SetLinkLagParamsRequest) error {
	var apiInput struct {
		Requests map[ObjectId]rawLinkLagParams `json:"links"`
	}
	apiInput.Requests = make(map[ObjectId]rawLinkLagParams)
	for k, v := range *req {
		raw, err := v.raw()
		if err != nil {
			return err
		}
		apiInput.Requests[k] = *raw
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlLeafServerLinkLabels, o.blueprintId),
		apiInput: &apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
