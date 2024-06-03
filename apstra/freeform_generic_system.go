package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlFFGenericSystems     = apiUrlBlueprintById + apiUrlPathDelim + "generic_systems"
	apiUrlFFGenericSystemsById = apiUrlFFGenericSystems + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformGenericSystem)
var _ json.Marshaler = new(FreeformGenericSystem)

type FreeformGenericSystem struct {
	SystemId   ObjectId
	SystemType SystemType
}

func (o FreeformGenericSystem) MarshalJSON() ([]byte, error) {
	var raw struct {
		SystemId   ObjectId `json:"system_id"`
		SystemType string   `json:"system_type"`
	}
	raw.SystemId = o.SystemId
	raw.SystemType = string(o.SystemType)
	return json.Marshal(&raw)
}

func (o *FreeformGenericSystem) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		SystemId   ObjectId   `json:"system_id"`
		SystemType SystemType `json:"system_type"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	o.SystemId = raw.SystemId
	o.SystemType = raw.SystemType
	return err
}

func (o *FreeformClient) GetFreeformGenericSystem(ctx context.Context, systemId ObjectId) (*FreeformGenericSystem, error) {
	response := new(FreeformGenericSystem)
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFGenericSystemsById, o.blueprintId, systemId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *FreeformClient) GetAllFreeformGenericSystems(ctx context.Context) (map[int]FreeformGenericSystem, error) {
	var response []FreeformGenericSystem
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFGenericSystems, o.blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}
func (o *FreeformClient) CreateFreeformGenericSystem(ctx context.Context, in *FreeformGenericSystem) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlFFGenericSystems,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *FreeformClient) UpdateFreeformGenericSystem(ctx context.Context, id ObjectId, in *FreeformGenericSystem) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFFGenericSystemsById, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
