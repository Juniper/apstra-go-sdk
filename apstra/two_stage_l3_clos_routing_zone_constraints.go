package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlRoutingZoneConstraints    = apiUrlBlueprintByIdPrefix + "routing-zone-constraints"
	apiUrlRoutingZoneConstraintById = apiUrlRoutingZoneConstraints + apiUrlPathDelim + "%s"
)

var _ json.Marshaler = (*RoutingZoneConstraintData)(nil)

type RoutingZoneConstraintData struct {
	Label           string
	Mode            RoutingZoneConstraintMode
	MaxRoutingZones *int
	RoutingZoneIds  []ObjectId
}

func (o RoutingZoneConstraintData) MarshalJSON() ([]byte, error) {
	var raw struct {
		Label                      string     `json:"label"`
		RoutingZonesListConstraint string     `json:"routing_zones_list_constraint,omitempty"`
		MaxCountConstraint         *int       `json:"max_count_constraint"`
		Constraints                []ObjectId `json:"constraints"`
	}

	raw.Label = o.Label
	raw.RoutingZonesListConstraint = o.Mode.String()
	raw.MaxCountConstraint = o.MaxRoutingZones
	if o.RoutingZoneIds != nil {
		raw.Constraints = o.RoutingZoneIds
	} else {
		raw.Constraints = []ObjectId{}
	}

	return json.Marshal(raw)
}

var _ json.Unmarshaler = (*RoutingZoneConstraint)(nil)

type RoutingZoneConstraint struct {
	Id   ObjectId
	Data *RoutingZoneConstraintData
}

func (o *RoutingZoneConstraint) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id                         ObjectId   `json:"id"`
		Label                      string     `json:"label"`
		RoutingZonesListConstraint string     `json:"routing_zones_list_constraint"`
		MaxCountConstraint         *int       `json:"max_count_constraint"`
		Constraints                []ObjectId `json:"constraints,omitempty"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(RoutingZoneConstraintData)
	o.Data.Label = raw.Label
	err = o.Data.Mode.FromString(raw.RoutingZonesListConstraint)
	if err != nil {
		return err
	}
	o.Data.MaxRoutingZones = raw.MaxCountConstraint
	o.Data.RoutingZoneIds = raw.Constraints

	return nil
}

func (o *TwoStageL3ClosClient) CreateRoutingZoneConstraint(ctx context.Context, in *RoutingZoneConstraintData) (ObjectId, error) {
	var response objectIdResponse
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlRoutingZoneConstraints, o.Id()),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *TwoStageL3ClosClient) GetRoutingZoneConstraint(ctx context.Context, id ObjectId) (*RoutingZoneConstraint, error) {
	var response RoutingZoneConstraint
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlRoutingZoneConstraintById, o.Id(), id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	response.Id = id // id field not returned by /api/blueprints/{blueprint_id}/routing-zone-constraints/{routing_zone_constraint_id}
	return &response, nil
}

func (o *TwoStageL3ClosClient) GetAllRoutingZoneConstraints(ctx context.Context) ([]RoutingZoneConstraint, error) {
	var response struct {
		Items []RoutingZoneConstraint `json:"items"`
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlRoutingZoneConstraints, o.Id()),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *TwoStageL3ClosClient) GetRoutingZoneConstraintByName(ctx context.Context, name string) (*RoutingZoneConstraint, error) {
	all, err := o.GetAllRoutingZoneConstraints(ctx)
	if err != nil {
		return nil, err
	}

	var result *RoutingZoneConstraint
	for _, routingZoneConstraint := range all {
		routingZoneConstraint := routingZoneConstraint
		if routingZoneConstraint.Data.Label == name {
			if result != nil {
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("blueprint %q has multiple routing zone constrains with label %q", o.Id(), name),
				}
			}

			result = &routingZoneConstraint
		}
	}

	if result == nil {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("blueprint %q has no routing zone constrains with label %q", o.Id(), name),
		}
	}

	return result, nil
}

func (o *TwoStageL3ClosClient) UpdateRoutingZoneConstraint(ctx context.Context, id ObjectId, in *RoutingZoneConstraintData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlRoutingZoneConstraintById, o.Id(), id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) DeleteRoutingZoneConstraint(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlRoutingZoneConstraintById, o.Id(), id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
