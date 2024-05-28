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
var _ json.Marshaler = new(FreeformPropertySet)

type FreeformPropertySet struct {
	Id   ObjectId
	Data *FFPropertySetData
}

type FFPropertySetData struct {
	SystemId ObjectId
	Label    string
	//TemplateId string    //we think this is not useful. TBD.//
	Values     string
	ValuesYaml string
}

func (o FreeformPropertySet) MarshalJSON() ([]byte, error) {
	var raw struct {
		SystemId   ObjectId `json:"system_id"`
		Id         ObjectId `json:"property_set_id"`
		Label      string   `json:"label,omitempty"`
		Values     string   `json:"values,omitempty"`
		ValuesYaml string   `json:"values_yaml,omitempty"`
	}
	raw.Id = o.Id
	if o.Data != nil {
		raw.SystemId = o.Data.SystemId
		raw.Label = o.Data.Label
		raw.Values = o.Data.Values
		raw.ValuesYaml = o.Data.ValuesYaml
	}
	return json.Marshal(&raw)
}

func (o *FreeformPropertySet) UnmarshalJSON(bytes []byte) error {
	if o.Data == nil {
		o.Data = new(FFPropertySetData)
	}
	var raw struct {
		Id         ObjectId `json:"property_set_id"`
		SystemId   ObjectId `json:"system_id"`
		Label      string   `json:"label,omitempty"`
		Values     string   `json:"values,omitempty"`
		ValuesYaml string   `json:"value_yaml,omitempty"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	o.Id = raw.Id
	o.Data.SystemId = raw.SystemId
	o.Data.Label = raw.Label
	o.Data.Values = raw.Values
	o.Data.ValuesYaml = raw.ValuesYaml
	return err
}

func (o *FreeformClient) GetFreeformPropertySet(ctx context.Context, id ObjectId) (*FreeformPropertySet, error) {
	response := new(FreeformPropertySet)
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFPropertySetById, o.blueprintId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *FreeformClient) GetAllFreeformPropertySets(ctx context.Context, label string) ([]FreeformPropertySet, error) {
	var response []FreeformPropertySet
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFPropertySets, o.blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *FreeformClient) CreateFreeformPropertySet(ctx context.Context, in *FreeformPropertySet) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlFFPropertySets,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *FreeformClient) UpdateFreeformPropertySet(ctx context.Context, id ObjectId, in *FreeformPropertySet) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlFFPropertySets, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *FreeformClient) DeleteFreeformPropertySet(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFFPropertySetById, o.blueprintId, id),
	})
}
