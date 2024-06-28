package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlFFPropertySets    = apiUrlBlueprintById + apiUrlPathDelim + "property-sets"
	apiUrlFFPropertySetById = apiUrlFFPropertySets + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformPropertySet)

type FreeformPropertySet struct {
	Id   ObjectId
	Data *FreeformPropertySetData
}

func (o *FreeformPropertySet) UnmarshalJSON(bytes []byte) error {
	if o.Data == nil {
		o.Data = new(FreeformPropertySetData)
	}

	var raw struct {
		Id       ObjectId        `json:"property_set_id"`
		SystemId *ObjectId       `json:"system_id"`
		Label    string          `json:"label,omitempty"`
		Values   json.RawMessage `json:"values,omitempty"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data.SystemId = raw.SystemId
	o.Data.Label = raw.Label
	o.Data.Values = raw.Values

	return err
}

type FreeformPropertySetData struct {
	SystemId *ObjectId       `json:"system_id"`
	Label    string          `json:"label"`
	Values   json.RawMessage `json:"values,omitempty"`
}

func (o *FreeformClient) CreatePropertySet(ctx context.Context, in *FreeformPropertySetData) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFFPropertySets, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *FreeformClient) GetPropertySet(ctx context.Context, id ObjectId) (*FreeformPropertySet, error) {
	var response FreeformPropertySet

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFPropertySetById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) GetAllPropertySets(ctx context.Context) ([]FreeformPropertySet, error) {
	var response struct {
		Items []FreeformPropertySet `json:"items"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFPropertySets, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) UpdatePropertySet(ctx context.Context, id ObjectId, in *FreeformPropertySetData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFFPropertySetById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) DeletePropertySet(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFFPropertySetById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
