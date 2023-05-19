package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiUrlBlueprintPropertySets = apiUrlBlueprintById + apiUrlPathDelim + "property-sets"
	apiUrlBlueprintPropertySet  = apiUrlBlueprintById + apiUrlPathDelim + "property-sets" + apiUrlPathDelim + "%s"
)

type ImportedPropertySet struct {
	Id         ObjectId        `json:"id"`
	Label      string          `json:"label"`
	Stale      bool            `json:"stale"`
	Values     json.RawMessage `json:"values"`
	ValuesYaml string          `json:"values_yaml"`
}

type ImportPropertySetRequest struct {
	Id   ObjectId `json:"id"`
	Keys []string `json:"keys"`
}

type ImportPropertySetResponse struct {
	Id ObjectId `json:"id"`
}

func (o *Client) importPropertySet(ctx context.Context, bpid ObjectId, psid ObjectId, keys ...string) (ObjectId, error) {

	importPropertySetUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintPropertySets, bpid.String()))
	if err != nil {
		return "", err
	}

	response := &ImportPropertySetResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		url:         importPropertySetUrl,
		apiInput:    ImportPropertySetRequest{Id: psid, Keys: keys},
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, convertTtaeToAceWherePossible(err)
}

func (o *Client) getAllImportedPropertySets(ctx context.Context, bpid ObjectId) ([]ImportedPropertySet, error) {
	result := &struct {
		Items []ImportedPropertySet `json:"items"`
	}{}
	allImportedPropertySetUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintPropertySets, bpid.String()))

	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         allImportedPropertySetUrl,
		apiResponse: result,
	})
	return result.Items, convertTtaeToAceWherePossible(err)
}

func (o *Client) getImportedPropertySetByName(ctx context.Context, bpid ObjectId, name string) (*ImportedPropertySet, error) {
	allps, err := o.getAllImportedPropertySets(ctx, bpid)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	for _, t := range allps {
		if t.Label == name {
			return &t, nil
		}
	}

	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf(" Property Set with name '%s' not found", name),
	}
}

func (o *Client) getImportedPropertySet(ctx context.Context, bpid ObjectId, psid ObjectId) (*ImportedPropertySet, error) {
	result := &ImportedPropertySet{}

	getImportedPropertySetUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintPropertySet, bpid.String(), psid.String()))
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         getImportedPropertySetUrl,
		apiResponse: result,
	})
	return result, convertTtaeToAceWherePossible(err)
}

func (o *Client) updateImportedPropertySet(ctx context.Context, bpid ObjectId, psid ObjectId, keys ...string) error {

	updateImportedPropertySetUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintPropertySet, bpid.String(), psid.String()))
	if err != nil {
		return err
	}

	err = o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPut,
		url:    updateImportedPropertySetUrl,
		apiInput: ImportPropertySetRequest{
			Id:   psid,
			Keys: keys,
		},
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) deleteImportedPropertySet(ctx context.Context, bpid ObjectId, pid ObjectId) error {
	deleteImportedPropertySetUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintPropertySet, bpid.String(), pid.String()))
	if err != nil {
		return err
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		url:    deleteImportedPropertySetUrl,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
