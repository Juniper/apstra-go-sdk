package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlConfigTemplates    = apiUrlBlueprintById + apiUrlPathDelim + "config-templates"
	apiUrlConfigTemplateById = apiUrlConfigTemplates + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(ConfigTemplate)
var _ json.Marshaler = new(ConfigTemplate)

type ConfigTemplate struct {
	Id   ObjectId
	Data *ConfigTemplateData
}

type ConfigTemplateData struct {
	Label string
	//TemplateId string    //we think this is not useful. TBD.//
	Text string
}

func (o ConfigTemplate) MarshalJSON() ([]byte, error) {
	var raw struct {
		Id    ObjectId `json:"id"`
		Label string   `json:"label,omitempty"`
		Text  string   `json:"text,omitempty"`
	}
	raw.Id = o.Id
	if o.Data != nil {
		raw.Label = o.Data.Label
		raw.Text = o.Data.Text
	}
	return json.Marshal(&raw)
}

func (o *ConfigTemplate) UnmarshalJSON(bytes []byte) error {
	if o.Data == nil {
		o.Data = new(ConfigTemplateData)
	}
	var raw struct {
		Id    ObjectId `json:"id"`
		Label string   `json:"label,omitempty"`
		Text  string   `json:"text,omitempty"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	o.Id = raw.Id
	o.Data.Label = raw.Label
	o.Data.Text = raw.Text
	return err
}

func (o *FreeformClient) GetConfigTemplate(ctx context.Context, id ObjectId) (*ConfigTemplate, error) {
	response := new(ConfigTemplate)
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlConfigTemplateById, o.blueprintId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *FreeformClient) GetAllConfigTemplates(ctx context.Context, label string) ([]ConfigTemplate, error) {
	var response []ConfigTemplate
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlConfigTemplates, o.blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *FreeformClient) CreateConfigTemplate(ctx context.Context, in *ConfigTemplate) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlConfigTemplates,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *FreeformClient) UpdateConfigTemplate(ctx context.Context, id ObjectId, in *ConfigTemplate) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlConfigTemplates, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *FreeformClient) DeleteConfigTemplate(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlConfigTemplateById, o.blueprintId, id),
	})
}
