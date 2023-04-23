package apstra

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlPropertySets       = "/api/property-sets"
	apiUrlPropertySetsPrefix = apiUrlPropertySets + apiUrlPathDelim
	apiUrlPropertySetById    = apiUrlPropertySetsPrefix + "%s"
)

type PropertySet struct {
	Id        ObjectId
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      *PropertySetData
}

// We don't really need a "raw" PropertySetData because there are no iota etc that need translation
type PropertySetData struct {
	Label      string            `json:"label"`
	Values     map[string]string `json:"values"`
	Blueprints []ObjectId        `json:"blueprints,omitempty"`
}

type rawPropertySet struct {
	Id         ObjectId          `json:"id,omitempty"`
	Label      string            `json:"label"`
	Values     map[string]string `json:"values"`
	Blueprints []ObjectId        `json:"blueprints,omitempty"`
	CreatedAt  string            `json:"created_at,omitempty"`
	UpdatedAt  string            `json:"updated_at,omitempty"`
}

type PropertySetRequest PropertySetData

func (o *rawPropertySet) polish() (*PropertySet, error) {
	created, err := time.Parse("2006-01-02T15:04:05.000000+0000", o.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing create time %s - %w", o.CreatedAt, err)
	}
	updated, err := time.Parse("2006-01-02T15:04:05.000000+0000", o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing update time %s - %w", o.UpdatedAt, err)
	}
	return &PropertySet{
		Id:        o.Id,
		CreatedAt: created,
		UpdatedAt: updated,
		Data: &PropertySetData{
			Label:      o.Label,
			Values:     o.Values,
			Blueprints: o.Blueprints,
		},
	}, nil
}

func (o *Client) listAllPropertySets(ctx context.Context) ([]ObjectId, error) {
	response := &struct {
		Items []ObjectId `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlPropertySets,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getPropertySet(ctx context.Context, id ObjectId) (*rawPropertySet, error) {
	response := &rawPropertySet{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlPropertySetById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getPropertySetByLabel(ctx context.Context, label string) (*rawPropertySet, error) {
	propertySets, err := o.getPropertySetsByLabel(ctx, label)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	if len(propertySets) > 1 {
		return nil, ApstraClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("found multiple (%d) property sets with label %q", len(propertySets), label),
		}
	}

	return &propertySets[0], nil
}

func (o *Client) getPropertySetsByLabel(ctx context.Context, label string) ([]rawPropertySet, error) {
	allPropertySets, err := o.getAllPropertySets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	var result []rawPropertySet
	for _, ps := range allPropertySets {
		if ps.Label == label {
			result = append(result, ps)
		}
	}

	if len(result) == 0 {
		return nil, ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("property set with label '%s' not found", label),
		}
	}
	return result, nil
}

func (o *Client) getAllPropertySets(ctx context.Context) ([]rawPropertySet, error) {
	ids, err := o.listAllPropertySets(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]rawPropertySet, len(ids))
	for i := range ids {
		ps, err := o.getPropertySet(ctx, ids[i])
		if err != nil {
			return nil, err
		}
		result[i] = *ps
	}
	return result, nil
}

func (o *Client) createPropertySet(ctx context.Context, in *PropertySetRequest) (ObjectId, error) {
	if len(in.Blueprints) != 0 {
		return "", errors.New("blueprints field must be empty when creating property set")
	}
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlPropertySets,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) updatePropertySet(ctx context.Context, id ObjectId, in *PropertySetRequest) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlPropertySetById, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) deletePropertySet(ctx context.Context, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlPropertySetById, id),
	})
}
