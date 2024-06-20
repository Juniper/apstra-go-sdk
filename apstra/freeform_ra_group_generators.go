package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlFFRaGroupGenerators    = apiUrlBlueprintById + apiUrlPathDelim + "ra-group-generators"
	apiUrlFFRaGroupGeneratorById = apiUrlFFRaGroupGenerators + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformRaGroupGenerator)

type FreeformRaGroupGenerator struct {
	Id   ObjectId
	Data *FreeformRaGroupGeneratorData
}

type FreeformRaGroupGeneratorData struct {
	ParentId *ObjectId
	Label    string
	Scope    string
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
	o.Data.ParentId = raw.ParentId
	o.Data.Label = raw.Label
	o.Data.Scope = raw.Scope
	return err
}

func (o *FreeformClient) CreateRaGenerator(ctx context.Context, in *FreeformRaGroupGeneratorData) (ObjectId, error) {
	var response objectIdResponse
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFFRaGroupGenerators, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *FreeformClient) GetAllGroupGenerators(ctx context.Context) ([]FreeformRaGroupGenerator, error) {
	var response struct {
		Items []FreeformRaGroupGenerator `json:"items"`
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFRaGroupGenerators, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *FreeformClient) GetRaGroupGenerator(ctx context.Context, id ObjectId) (*FreeformRaGroupGenerator, error) {
	response := new(FreeformRaGroupGenerator)
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFRaGroupGeneratorById, o.blueprintId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *FreeformClient) UpdateRaGroupGenerator(ctx context.Context, id ObjectId, in *FreeformRaGroupGenerator) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFFRaGroupGeneratorById, o.blueprintId, id),
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
		urlStr: fmt.Sprintf(apiUrlFFRaGroupGeneratorById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
