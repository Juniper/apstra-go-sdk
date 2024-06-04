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

var _ json.Unmarshaler = new(FreeformSystem)

type FreeformSystem struct {
	Id   ObjectId
	Data *FreeformSystemData
}

func (o *FreeformSystem) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		SystemId   ObjectId   `json:"system_id"`
		SystemType SystemType `json:"system_type"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	o.Id = raw.SystemId
	o.Data.Type = raw.SystemType
	return err
}

var _ json.Marshaler = new(FreeformSystemData)

type FreeformSystemData struct {
	Type            SystemType
	Label           string
	Hostname        string
	Tags            []ObjectId
	DeviceProfileId DeviceProfile
}

func (o FreeformSystemData) MarshalJSON() ([]byte, error) {
	var raw struct {
		SystemType string `json:"system_type"`
	}
	raw.SystemType = o.Type.String()
	return json.Marshal(&raw)
}

func (o *FreeformClient) GetFreeformSystem(ctx context.Context, systemId ObjectId) (*FreeformSystem, error) {
	response := new(FreeformSystem)
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

func (o *FreeformClient) GetAllFreeformSystems(ctx context.Context) ([]FreeformSystem, error) {
	var response struct {
		Items []FreeformSystem `json:"items"`
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFGenericSystems, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *FreeformClient) CreateFreeformSystem(ctx context.Context, in *FreeformSystemData) (ObjectId, error) {
	response := new(objectIdResponse)
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

func (o *FreeformClient) UpdateFreeformSystem(ctx context.Context, id ObjectId, in *FreeformSystemData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFFGenericSystemsById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *FreeformClient) DeleteFreeformSystem(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlConfigTemplateById, o.blueprintId, id),
	})
}
