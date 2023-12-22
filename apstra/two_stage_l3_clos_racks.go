package apstra

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintAddRacks    = apiUrlBlueprintById + apiUrlPathDelim + "add-racks"
	apiUrlBlueprintDeleteRacks = apiUrlBlueprintById + apiUrlPathDelim + "delete-racks"
)

type rawTwoStageL3ClosRacksRequest struct {
	PodId          ObjectId             `json:"pod_id,omitempty"`
	RackTypeCounts map[ObjectId]int     `json:"rack_type_counts"`
	RackTypes      []rawRackTypeRequest `json:"rack_types"`
}

type TwoStageL3ClosRackRequest struct {
	PodId      ObjectId
	RackTypeId ObjectId
	RackType   *RackTypeRequest
}

func (o *TwoStageL3ClosRackRequest) raw(ctx context.Context, client *Client) (*rawTwoStageL3ClosRacksRequest, error) {
	var rawRTR *rawRackTypeRequest
	var err error

	if o.RackType != nil {
		rawRTR, err = o.RackType.raw(ctx, client)
		if err != nil {
			return nil, err
		}
	} else {
		// fetch the raw rack type from the API
		rawRT, err := client.getRackType(ctx, o.RackTypeId)
		if err != nil {
			return nil, err
		}

		rawRTR = rawRT.request()
	}

	if rawRTR.Id == "" {
		rawRTR.Id, err = uuid1AsObjectId()
		if err != nil {
			return nil, err
		}
	}

	return &rawTwoStageL3ClosRacksRequest{
		PodId:          o.PodId,
		RackTypeCounts: map[ObjectId]int{rawRTR.Id: 1},
		RackTypes:      []rawRackTypeRequest{*rawRTR},
	}, nil
}

func (o *TwoStageL3ClosClient) createRacks(ctx context.Context, request *rawTwoStageL3ClosRacksRequest) ([]ObjectId, error) {
	var apiResponse struct {
		RackIds []ObjectId `json:"rack_ids"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintAddRacks, o.Id()),
		apiInput:    &request,
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return apiResponse.RackIds, nil
}

func (o *TwoStageL3ClosClient) CreateRack(ctx context.Context, in *TwoStageL3ClosRackRequest) (ObjectId, error) {
	if in.RackTypeId != "" && in.RackType != nil {
		return "", errors.New("TwoStageL3ClosRackRequest passed to CreateRack must not have both a RackType and a RackTypeId")
	}

	request, err := in.raw(ctx, o.client)
	if err != nil {
		return "", err
	}

	ids, err := o.createRacks(ctx, request)
	if err != nil {
		return "", err
	}

	if len(ids) != 1 {
		return "", fmt.Errorf("creating a new rack should yield exactly 1 ID, got %d IDs: %s", len(ids), ids)
	}

	return ids[0], nil
}

func (o *TwoStageL3ClosClient) DeleteRack(ctx context.Context, id ObjectId) error {
	type request struct {
		RacksToDelete []ObjectId `json:"racks_to_delete"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPost,
		urlStr:   fmt.Sprintf(apiUrlBlueprintDeleteRacks, o.Id()),
		apiInput: &request{RacksToDelete: []ObjectId{id}},
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
