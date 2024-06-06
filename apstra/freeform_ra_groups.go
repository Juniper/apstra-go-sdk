package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlFFRaGroups    = apiUrlBlueprintById + apiUrlPathDelim + "ra-groups"
	apiUrlFFRaGroupById = apiUrlFFRaGroups + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformRaGroup)

type FreeformRaGroup struct {
	Id       ObjectId
	ParentId ObjectId
	Label    ObjectId
	Tags     []ObjectId
	Data     FreeformRaGroupData
	// todo add the data key value mapping
}

func (o *FreeformRaGroup) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ParentId ObjectId `json:"parent_id"`
		Label    ObjectId `json:"label"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	o.ParentId = raw.ParentId
	o.Label = raw.Label
	return err
}

var _ json.Marshaler = new(FreeformRaGroup)

type FreeformRaGroupData struct {
	Key   string
	Value string
}

func (o FreeformRaGroup) MarshalJSON() ([]byte, error) {
	var raw struct {
		ParentID string `json:"parent_id"`
		Label    string `json:"label"`
	}
	raw.ParentID = o.ParentId.String()
	raw.Label = o.Label.String()
	return json.Marshal(&raw)
}

func (o *FreeformClient) GetAllFreeformGroups(ctx context.Context) ([]FreeformRaGroup, error) {
	var response struct {
		Items []FreeformRaGroup `json:"items"`
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFRaGroups, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *FreeformClient) GetFreeformRaGroup(ctx context.Context, id ObjectId) (*FreeformRaGroup, error) {
	response := new(FreeformRaGroup)
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFFRaGroupById, o.blueprintId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *FreeformClient) UpdateFreeformRaGroup(ctx context.Context, id ObjectId, in *FreeformRaGroup) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFFRaGroupById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *FreeformClient) DeleteFreeformRaGroup(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFFRaGroupById, o.blueprintId, id),
	})
}
