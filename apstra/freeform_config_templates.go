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

//var _ json.Marshaler = new(ConfigTemplate) // todo delete me?

type ConfigTemplate struct {
	Id   ObjectId
	Data *ConfigTemplateData
}

type ConfigTemplateData struct {
	Label string `json:"label"`
	Text  string `json:"text"`
}

//func (o ConfigTemplate) MarshalJSON() ([]byte, error) { // todo delete me?
//	var raw struct {
//		Id    ObjectId `json:"id"`
//		Label string   `json:"label,omitempty"`
//		Text  string   `json:"text,omitempty"`
//	}
//	raw.Id = o.Id
//	if o.Data != nil {
//		raw.Label = o.Data.Label
//		raw.Text = o.Data.Text
//	}
//	return json.Marshal(&raw)
//}

func (o *ConfigTemplate) UnmarshalJSON(bytes []byte) error {
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
	if o.Data == nil {
		o.Data = new(ConfigTemplateData)
	}
	o.Data.Label = raw.Label
	o.Data.Text = raw.Text

	return err
}

func (o *FreeformClient) CreateConfigTemplate(ctx context.Context, in *ConfigTemplateData) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlConfigTemplates, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *FreeformClient) GetConfigTemplate(ctx context.Context, id ObjectId) (*ConfigTemplate, error) {
	var response ConfigTemplate

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlConfigTemplateById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) GetAllConfigTemplates(ctx context.Context) ([]ConfigTemplate, error) {
	var response struct {
		Items []ConfigTemplate `json:"items"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlConfigTemplates, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) UpdateConfigTemplate(ctx context.Context, id ObjectId, in *ConfigTemplateData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlConfigTemplateById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) DeleteConfigTemplate(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlConfigTemplateById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
