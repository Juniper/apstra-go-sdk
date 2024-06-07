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
	Id       ObjectId
	Scope    string
	Label    ObjectId
	ParentId ObjectId
}

func (o *FreeformRaGroupGenerator) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ParentId ObjectId `json:"parent_id"`
		Label    ObjectId `json:"label"`
		Scope    string   `json:"scope"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	o.ParentId = raw.ParentId
	o.Label = raw.Label
	o.Scope = raw.Scope
	return err
}

var _ json.Marshaler = new(FreeformRaGroupGenerator)

func (o FreeformRaGroupGenerator) MarshalJSON() ([]byte, error) {
	var raw struct {
		ParentID string `json:"parent_id"`
		Label    string `json:"label"`
		Scope    string `json:"scope"`
	}
	raw.ParentID = o.ParentId.String()
	raw.Label = o.Label.String()
	raw.Scope = o.Scope
	return json.Marshal(&raw)
}

func (o *FreeformClient) GetAllFreeformGroupGenerators(ctx context.Context) ([]FreeformRaGroupGenerator, error) {
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

func (o *FreeformClient) GetFreeformRaGroupGenerator(ctx context.Context, id ObjectId) (*FreeformRaGroupGenerator, error) {
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

func (o *FreeformClient) UpdateFreeformRaGroupGenerator(ctx context.Context, id ObjectId, in *FreeformRaGroupGenerator) error {
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

func (o *FreeformClient) DeleteFreeformRaGroupGenerator(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFFRaGroupGeneratorById, o.blueprintId, id),
	})
}
