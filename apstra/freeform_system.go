package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlFFGenericSystems     = apiUrlBlueprintById + apiUrlPathDelim + "generic-systems"
	apiUrlFFGenericSystemsById = apiUrlFFGenericSystems + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformSystem)

type FreeformSystem struct {
	Id   ObjectId
	Data *FreeformSystemData
}

func (o *FreeformSystem) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id ObjectId `json:"id"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(FreeformSystemData)
	err = json.Unmarshal(bytes, o.Data)
	if err != nil {
		return err
	}
	return nil
}

var (
	_ json.Marshaler   = new(FreeformSystemData)
	_ json.Unmarshaler = new(FreeformSystemData)
)

type FreeformSystemData struct {
	SystemId        *ObjectId
	Type            SystemType
	Label           string
	Hostname        string
	Tags            []ObjectId
	DeviceProfileId ObjectId
}

func (o FreeformSystemData) MarshalJSON() ([]byte, error) {
	var raw struct {
		SystemId        ObjectId   `json:"system_id,omitempty"`
		SystemType      string     `json:"system_type"`
		Label           string     `json:"label"`
		Hostname        string     `json:"hostname,omitempty"`
		Tags            []ObjectId `json:"tags"`
		DeviceProfileId ObjectId   `json:"device_profile_id"`
	}
	if o.SystemId != nil {
		raw.SystemId = *o.SystemId
	}
	raw.SystemType = o.Type.String()
	raw.Label = o.Label
	raw.Hostname = o.Hostname
	raw.Tags = o.Tags
	raw.DeviceProfileId = o.DeviceProfileId
	return json.Marshal(&raw)
}

func (o *FreeformSystemData) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		SystemId      *ObjectId  `json:"system_id"`
		SystemType    systemType `json:"system_type"`
		Label         string     `json:"label"`
		Hostname      string     `json:"hostname,omitempty"`
		Tags          []ObjectId `json:"tags"`
		DeviceProfile struct {
			Id ObjectId `json:"id"`
		} `json:"device_profile"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	st, err := raw.SystemType.parse()
	if err != nil {
		return err
	}

	o.SystemId = raw.SystemId
	o.Type = SystemType(st)
	o.Label = raw.Label
	o.Hostname = raw.Hostname
	o.Tags = raw.Tags
	o.DeviceProfileId = raw.DeviceProfile.Id

	return nil
}

func (o *FreeformClient) CreateSystem(ctx context.Context, in *FreeformSystemData) (ObjectId, error) {
	response := new(objectIdResponse)
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFFGenericSystems, o.blueprintId),
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
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
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFFGenericSystemsById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
