package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlFfRaGroupGenerators    = apiUrlBlueprintById + apiUrlPathDelim + "ra-group-generators"
	apiUrlFfRaGroupGeneratorById = apiUrlFfRaGroupGenerators + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformRaGroupGenerator)

type FreeformRaGroupGenerator struct {
	Id   ObjectId
	Data *FreeformRaGroupGeneratorData
}

func (o *FreeformRaGroupGenerator) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id       ObjectId  `json:"id"`
		ParentId *ObjectId `json:"parent_id"`
		Label    string    `json:"label"`
		Scope    string    `json:"scope"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(FreeformRaGroupGeneratorData)
	o.Data.ParentId = raw.ParentId
	o.Data.Label = raw.Label
	o.Data.Scope = raw.Scope

	return err
}

type FreeformRaGroupGeneratorData struct {
	ParentId *ObjectId `json:"parent_id"`
	Label    string    `json:"label"`
	Scope    string    `json:"scope"`
}

func (o *FreeformClient) CreateRaGroupGenerator(ctx context.Context, in *FreeformRaGroupGeneratorData) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFfRaGroupGenerators, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *FreeformClient) GetAllRaGroupGenerators(ctx context.Context) ([]FreeformRaGroupGenerator, error) {
	var response struct {
		Items []FreeformRaGroupGenerator `json:"items"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfRaGroupGenerators, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) GetRaGroupGenerator(ctx context.Context, id ObjectId) (*FreeformRaGroupGenerator, error) {
	var response FreeformRaGroupGenerator

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfRaGroupGeneratorById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) UpdateRaGroupGenerator(ctx context.Context, id ObjectId, in *FreeformRaGroupGeneratorData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFfRaGroupGeneratorById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) DeleteRaGroupGenerator(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfRaGroupGeneratorById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
