package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlFFPreferences = apiUrlBlueprintById + apiUrlPathDelim + "ra-groups"
)

var _ json.Unmarshaler = new(FreeformPreferences)

type FreeformPreferences struct {
	UserData []FreeformPreferencesData
}

func (o *FreeformPreferences) UnmarshalJSON(bytes []byte) error {
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

var _ json.Marshaler = new(FreeformPreferencesData)

type FreeformPreferencesData struct {
	Id       string
	Position FFPosition
}

type FFPosition [3]string

func (o FreeformPreferencesData) MarshalJSON() ([]byte, error) {
	var raw struct {
		Preferences string `json:"preferences"`
	}
	raw.Preferences = o.Type.String()
	return json.Marshal(&raw)
}

func (o *FreeformClient) GetFreeformPreferences(ctx context.Context) ([]FreeformSystem, error) {
	var response struct {
		UserData []FreeformPreferencesData `json:"items"`
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFPreferences, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *FreeformClient) UpdateFreeformPreferences(ctx context.Context, id ObjectId, in *FreeformSystemData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFFPreferences, o.blueprintId),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *FreeformClient) DeleteFreeformPreferences(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFFPreferences, o.blueprintId, id),
	})
}
