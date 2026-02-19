// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/internal"
)

const (
	apiUrlFfAggLinks    = apiUrlBlueprintById + apiUrlPathDelim + "aggregate-links"
	apiUrlFfAggLinkById = apiUrlFfAggLinks + apiUrlPathDelim + "%s"
)

var (
	_ internal.IDer    = (*FreeformAggregateLink)(nil)
	_ json.Marshaler   = (*FreeformAggregateLink)(nil)
	_ json.Unmarshaler = (*FreeformAggregateLink)(nil)
)

type FreeformAggregateLink struct {
	Label          string
	MemberLinkIds  []string
	EndpointGroups [2]FreeformAggregateLinkEndpointGroup
	Tags           []string

	id string
}

func (o FreeformAggregateLink) ID() *string {
	if o.id == "" {
		return nil
	}
	return &o.id
}

func (o FreeformAggregateLink) MarshalJSON() ([]byte, error) {
	// set the endpoint group numbers (0 or 1) on the EndpointGroups and the Endpoints contained in each group
	for i := range o.EndpointGroups {
		o.EndpointGroups[i].endpointGroupNumber = i
		for j := range o.EndpointGroups[i].Endpoints {
			o.EndpointGroups[i].Endpoints[j].endpointGroup = i
		}
	}

	raw := struct {
		Label          string                                     `json:"label"`
		MemberLinkIDs  []string                                   `json:"member_link_ids"`
		Endpoints      []FreeformAggregateLinkEndpoint            `json:"endpoints"`
		EndpointGroups map[int]FreeformAggregateLinkEndpointGroup `json:"endpoint_groups"`
		Tags           []string                                   `json:"tags"`
	}{
		Label:          o.Label,
		MemberLinkIDs:  o.MemberLinkIds,
		Endpoints:      make([]FreeformAggregateLinkEndpoint, 0, len(o.EndpointGroups[0].Endpoints)+len(o.EndpointGroups[1].Endpoints)),
		EndpointGroups: make(map[int]FreeformAggregateLinkEndpointGroup, 2),
		Tags:           o.Tags,
	}

	for _, endpoint := range o.EndpointGroups[0].Endpoints {
		raw.Endpoints = append(raw.Endpoints, endpoint)
	}
	for _, endpoint := range o.EndpointGroups[1].Endpoints {
		raw.Endpoints = append(raw.Endpoints, endpoint)
	}

	raw.EndpointGroups[0] = o.EndpointGroups[0]
	raw.EndpointGroups[1] = o.EndpointGroups[1]

	return json.Marshal(raw)
}

func (o *FreeformAggregateLink) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID             string                                     `json:"id"`
		Label          string                                     `json:"label"`
		MemberLinkIDs  []string                                   `json:"member_link_ids"`
		Endpoints      []FreeformAggregateLinkEndpoint            `json:"endpoints"`
		EndpointGroups map[int]FreeformAggregateLinkEndpointGroup `json:"endpoint_groups"`
		Tags           []string                                   `json:"tags"`
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}
	if len(raw.EndpointGroups) != 2 {
		return fmt.Errorf("expected 2 endpoint groups, got %d", len(raw.EndpointGroups))
	}
	var ok bool
	if _, ok = raw.EndpointGroups[0]; !ok {
		return fmt.Errorf("endpoint group 0 not found")
	}
	if _, ok = raw.EndpointGroups[1]; !ok {
		return fmt.Errorf("endpoint group 1 not found")
	}

	o.id = raw.ID
	o.Label = raw.Label
	o.MemberLinkIds = raw.MemberLinkIDs
	o.Tags = raw.Tags
	o.EndpointGroups[0] = raw.EndpointGroups[0]
	o.EndpointGroups[1] = raw.EndpointGroups[1]

	for i, ep := range raw.Endpoints {
		if ep.endpointGroup != 0 && ep.endpointGroup != 1 {
			return fmt.Errorf("unexpected endpoint group %d at index %d", ep.endpointGroup, i)
		}
		o.EndpointGroups[ep.endpointGroup].Endpoints = append(o.EndpointGroups[ep.endpointGroup].Endpoints, ep)
	}

	return nil
}

func (o FreeformClient) CreateAggregateLink(ctx context.Context, in FreeformAggregateLink) (string, error) {
	var response struct {
		ID string `json:"id"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFfAggLinks, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.ID, nil
}

func (o FreeformClient) GetAggregateLink(ctx context.Context, id string) (FreeformAggregateLink, error) {
	var response FreeformAggregateLink

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfAggLinkById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return response, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o FreeformClient) GetAggregateLinks(ctx context.Context) ([]FreeformAggregateLink, error) {
	var response struct {
		Items []FreeformAggregateLink `json:"items"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfAggLinks, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o FreeformClient) GetAggregateLinkByLabel(ctx context.Context, label string) (FreeformAggregateLink, error) {
	var result []FreeformAggregateLink

	all, err := o.GetAggregateLinks(ctx)
	if err != nil {
		return FreeformAggregateLink{}, err
	}

	for _, item := range all {
		if item.Label == label {
			result = append(result, item)
		}
	}

	switch len(result) {
	case 0:
		return zero.SliceItem(result), ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("%T with label %s not found", zero.SliceItem(result), label),
		}
	case 1:
		return result[0], nil
	default: // len(result) > 1
		return zero.SliceItem(result), ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("found multiple candidate %T with label %s", zero.SliceItem(result), label),
		}
	}
}

func (o FreeformClient) ListAggregateLinks(ctx context.Context) ([]string, error) {
	all, err := o.GetAggregateLinks(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(all))
	for i, item := range all {
		result[i] = item.id
	}

	return result, nil
}

func (o *FreeformClient) UpdateAggregateLink(ctx context.Context, in FreeformAggregateLink) error {
	if in.ID() == nil {
		return fmt.Errorf("id should not be nil")
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFfAggLinkById, o.blueprintId, *in.ID()),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) DeleteAggregateLink(ctx context.Context, id string) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfAggLinkById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
