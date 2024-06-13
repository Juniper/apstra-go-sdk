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

var _ json.Unmarshaler = new(FreeformRaLocalPoolGeneratorData)

type FreeformRaLocalPoolGenerator struct {
	Id   ObjectId
	Data *FreeformRaLocalPoolGeneratorData
}
type FreeformRaLocalPoolGeneratorData struct {
	PoolType     string
	ResourceType FFResourceType
	Label        string
	Scope        string
	Chunks       []FFLocalIntPoolChunk
	version      int
}

func (o *FreeformRaLocalPoolGeneratorData) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Label        ObjectId       `json:"label"`
		Scope        string         `json:"scope"`
		PoolType     string         `json:"pool_type"`
		ResourceType FFResourceType `json:"resource_type"`
		Definition   struct {
			Chunks []FFLocalIntPoolChunk `json:"chunks"`
		} `json:"definition"`
		Version string `json:"version"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	o.Label = string(raw.Label)
	o.Scope = raw.Scope
	o.PoolType = raw.PoolType
	o.ResourceType = raw.ResourceType
	return err
}

var _ json.Marshaler = new(FreeformRaLocalPoolGenerator)

func (o FreeformRaLocalPoolGenerator) MarshalJSON() ([]byte, error) {
	var raw struct {
		Label        string         `json:"label"`
		Scope        string         `json:"scope"`
		PoolType     string         `json:"pool_type"`
		ResourceType FFResourceType `json:"resource_type"`
	}
	raw.Label = o.Data.Label
	raw.Scope = o.Data.Scope
	raw.PoolType = o.Data.PoolType
	raw.ResourceType = o.Data.ResourceType
	return json.Marshal(&raw)
}

func (o *FreeformClient) CreateLocalPoolGenerator(ctx context.Context, in *FreeformRaLocalPoolGenerator) (ObjectId, error) {
	var response objectIdResponse
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFFRaLocalPoolGenerators, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *FreeformClient) GetAllLocalPoolGenerators(ctx context.Context) ([]FreeformRaLocalPoolGenerator, error) {
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

func (o *FreeformClient) GetRaLocalPoolGenerator(ctx context.Context, id ObjectId) (*FreeformRaLocalPoolGenerator, error) {
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

func (o *FreeformClient) UpdateRaLocalPoolGenerator(ctx context.Context, id ObjectId, in *FreeformRaLocalPoolGenerator) error {
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

func (o *FreeformClient) DeleteRaLocalPoolGenerator(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFFRaLocalPoolGeneratorById, o.blueprintId, id),
	})
}
