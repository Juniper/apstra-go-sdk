package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlFFRaLocalPools    = apiUrlBlueprintById + apiUrlPathDelim + "ra-local-pools"
	apiUrlFFRaLocalPoolById = apiUrlFFRaLocalPools + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformRaLocalIntPool)

type FreeformRaLocalIntPool struct {
	Id   ObjectId
	Data *FreeformRaLocalIntPoolData
}

func (o *FreeformRaLocalIntPool) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id           ObjectId       `json:"id"`
		Label        string         `json:"label"`
		PoolType     string         `json:"pool_type"`
		ResourceType FFResourceType `json:"resource_type"`
		OwnerId      ObjectId       `json:"owner_id"`
		GeneratorId  *ObjectId      `json:"generator_id"`
		Definition   struct {
			Chunks []FFLocalIntPoolChunk `json:"chunks"`
		} `json:"definition"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	o.Id = raw.Id
	o.Data.Label = raw.Label
	o.Data.ResourceType = raw.ResourceType
	o.Data.OwnerId = raw.OwnerId
	o.Data.GeneratorId = raw.GeneratorId
	o.Data.Chunks = raw.Definition.Chunks

	return err
}

var _ json.Marshaler = new(FreeformRaLocalIntPoolData)
var _ json.Unmarshaler = new(FreeformRaLocalIntPoolData)

type FreeformRaLocalIntPoolData struct {
	ResourceType FFResourceType
	Label        string
	OwnerId      ObjectId
	GeneratorId  *ObjectId
	Chunks       []FFLocalIntPoolChunk
}

type FFLocalIntPoolChunk struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// todo is o below a pointer or just normal?
func (o *FreeformRaLocalIntPoolData) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ResourceType FFResourceType `json:"resource_type"`
		Label        string         `json:"label"`
		OwnerId      ObjectId       `json:"owner_id"`
		GeneratorId  *ObjectId      `json:"generator_id"`
		Definition   struct {
			Chunks []FFLocalIntPoolChunk `json:"chunks"`
		} `json:"definition"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	o.Label = raw.Label
	o.ResourceType = raw.ResourceType
	o.OwnerId = raw.OwnerId
	o.GeneratorId = raw.GeneratorId
	o.Chunks = raw.Definition.Chunks

	return err
}

func (o FreeformRaLocalIntPoolData) MarshalJSON() ([]byte, error) {
	var raw struct {
		Label        string                `json:"label"`
		OwnerId      string                `json:"owner_id"`
		GeneratorId  *ObjectId             `json:"generator_id"`
		ResourceType FFResourceType        `json:"resource_type"`
		Chunks       []FFLocalIntPoolChunk `json:"chunks"`
	}
	raw.Label = o.Label
	raw.OwnerId = string(o.OwnerId)
	raw.GeneratorId = o.GeneratorId
	raw.ResourceType = o.ResourceType
	raw.Chunks = o.Chunks
	return json.Marshal(&raw)
}

func (o *FreeformClient) CreateLocalIntPool(ctx context.Context, in *FreeformRaLocalIntPool) (ObjectId, error) {
	var response objectIdResponse
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFFRaLocalPools, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *FreeformClient) GetAllLocalIntPools(ctx context.Context) ([]FreeformRaLocalIntPool, error) {
	var response struct {
		Items []FreeformRaLocalIntPool `json:"items"`
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFRaLocalPools, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *FreeformClient) GetRaLocalIntPool(ctx context.Context, id ObjectId) (*FreeformRaLocalIntPool, error) {
	response := new(FreeformRaLocalIntPool)
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFRaLocalPoolById, o.blueprintId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *FreeformClient) UpdateRaLocalIntPool(ctx context.Context, id ObjectId, in *FreeformRaLocalIntPool) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFFRaLocalPoolById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *FreeformClient) DeleteRaLocalIntPool(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFFRaLocalPoolById, o.blueprintId, id),
	})
}
