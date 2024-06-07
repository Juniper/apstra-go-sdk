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

var _ json.Unmarshaler = new(FreeformRaLocalPools)

type FreeformRaLocalPools struct {
	Id           ObjectId
	PoolType     string // todo do i need an enum or something here?
	ResourceType string
	Label        ObjectId
	// todo implement "one of the next values " chunks
}

func (o *FreeformRaLocalPools) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Label        ObjectId `json:"label"`
		PoolType     string   `json:"pool_type"`
		ResourceType string   `json:"resource_type"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	o.Label = raw.Label
	o.PoolType = raw.PoolType
	o.ResourceType = raw.ResourceType
	return err
}

var _ json.Marshaler = new(FreeformRaLocalPools)

func (o FreeformRaLocalPools) MarshalJSON() ([]byte, error) {
	var raw struct {
		Label        string `json:"label"`
		PoolType     string `json:"pool_type"`
		ResourceType string `json:"resource_type"`
	}
	raw.Label = o.Label.String()
	raw.PoolType = o.PoolType
	raw.ResourceType = o.ResourceType
	return json.Marshal(&raw)
}

func (o *FreeformClient) GetAllFreeformLocalPools(ctx context.Context) ([]FreeformRaLocalPools, error) {
	var response struct {
		Items []FreeformRaLocalPools `json:"items"`
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

func (o *FreeformClient) GetFreeformRaLocalPool(ctx context.Context, id ObjectId) (*FreeformRaLocalPools, error) {
	response := new(FreeformRaLocalPools)
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

func (o *FreeformClient) UpdateFreeformRaLocalPool(ctx context.Context, id ObjectId, in *FreeformRaLocalPools) error {
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

func (o *FreeformClient) DeleteFreeformRaLocalPool(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFFRaLocalPoolById, o.blueprintId, id),
	})
}
