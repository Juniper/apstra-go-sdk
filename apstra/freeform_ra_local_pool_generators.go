package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlFFRaLocalPoolGenerators    = apiUrlBlueprintById + apiUrlPathDelim + "ra-local-pool-generators"
	apiUrlFFRaLocalPoolGeneratorById = apiUrlFFRaLocalPoolGenerators + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformRaLocalPoolGenerator)

type FreeformRaLocalPoolGenerator struct {
	Id           ObjectId
	Scope        string
	PoolType     int // todo do i need an enum or something here?
	ResourceType string
	Label        ObjectId
	// todo implement "one of the next values " chunks
}

func (o *FreeformRaLocalPoolGenerator) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Label        ObjectId `json:"label"`
		Scope        string   `json:"scope"`
		PoolType     int      `json:"pool_type"`
		ResourceType string   `json:"resource_type"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	o.Label = raw.Label
	o.Scope = raw.Scope
	o.PoolType = raw.PoolType
	o.ResourceType = raw.ResourceType
	return err
}

var _ json.Marshaler = new(FreeformRaLocalPoolGenerator)

func (o FreeformRaLocalPoolGenerator) MarshalJSON() ([]byte, error) {
	var raw struct {
		Label        string `json:"label"`
		Scope        string `json:"scope"`
		PoolType     int    `json:"pool_type"`
		ResourceType string `json:"resource_type"`
	}
	raw.Label = o.Label.String()
	raw.Scope = o.Scope
	raw.PoolType = o.PoolType
	raw.ResourceType = o.ResourceType
	return json.Marshal(&raw)
}

func (o *FreeformClient) GetAllFreeformLocalPoolGenerators(ctx context.Context) ([]FreeformRaLocalPoolGenerator, error) {
	var response struct {
		Items []FreeformRaLocalPoolGenerator `json:"items"`
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFRaLocalPoolGenerators, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *FreeformClient) GetFreeformRaLocalPoolGenerator(ctx context.Context, id ObjectId) (*FreeformRaLocalPoolGenerator, error) {
	response := new(FreeformRaLocalPoolGenerator)
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFRaLocalPoolGeneratorById, o.blueprintId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *FreeformClient) UpdateFreeformRaLocalPoolGenerator(ctx context.Context, id ObjectId, in *FreeformRaLocalPoolGenerator) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFFRaLocalPoolGeneratorById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *FreeformClient) DeleteFreeformRaLocalPoolGenerator(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFFRaLocalPoolGeneratorById, o.blueprintId, id),
	})
}
