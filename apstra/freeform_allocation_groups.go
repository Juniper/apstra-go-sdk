package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	apiUrlFfAllocGroups          = apiUrlBlueprintById + apiUrlPathDelim + "resource_groups"
	apiUrlFfAllocGroupByType     = apiUrlFfAllocGroups + apiUrlPathDelim + "%s"
	apiUrlFfAllocGroupByTypeName = apiUrlFfAllocGroupByType + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformAllocGroup)

type FreeformAllocGroup struct {
	Id   ObjectId
	Data *FreeformAllocGroupData
}

func (o *FreeformAllocGroup) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id      ObjectId   `json:"id"`
		Name    string     `json:"name"`
		Type    string     `json:"type"`
		PoolIds []ObjectId `json:"pool_ids"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	if o.Data == nil {
		o.Data = new(FreeformAllocGroupData)
	}

	o.Id = raw.Id
	o.Data.Name = raw.Name
	o.Data.PoolIds = raw.PoolIds

	return o.Data.Type.FromString(raw.Type)
}

type FreeformAllocGroupData struct {
	Name    string           `json:"group_name"`
	Type    ResourcePoolType `json:"-"`
	PoolIds []ObjectId       `json:"pool_ids"`
}

func (o *FreeformClient) CreateAllocGroup(ctx context.Context, in *FreeformAllocGroupData) (ObjectId, error) {
	return o.createAllocGroup(ctx, in, in.Type)
}

func (o *FreeformClient) createAllocGroup(ctx context.Context, in *FreeformAllocGroupData, groupType ResourcePoolType) (ObjectId, error) {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPost,
		urlStr:   fmt.Sprintf(apiUrlFfAllocGroupByType, o.blueprintId, groupType),
		apiInput: in,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return ObjectId(fmt.Sprintf("rag_%s_%s", groupType, in.Name)), nil
}

func (o *FreeformClient) GetAllAllocGroups(ctx context.Context) ([]FreeformAllocGroup, error) {
	var response struct {
		Items []FreeformAllocGroup `json:"items"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfAllocGroups, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) GetAllocGroup(ctx context.Context, id ObjectId) (*FreeformAllocGroup, error) {
	parts := strings.SplitN(id.String(), "_", 3)
	if parts[0] != "rag" || len(parts) != 3 {
		return nil, ClientErr{
			errType: ErrInvalidId,
			err:     fmt.Errorf("freeform resource groups groupids must take the form rag_xxx_yyy, got : %q", id),
		}
	}

	var response FreeformAllocGroup
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfAllocGroupByTypeName, o.blueprintId, parts[1], parts[2]),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) GetAllocGroupByName(ctx context.Context, name string) (*FreeformAllocGroup, error) {
	all, err := o.GetAllAllocGroups(ctx)
	if err != nil {
		return nil, err
	}

	var result *FreeformAllocGroup
	for _, ffrag := range all {
		ffrag := ffrag
		if ffrag.Data.Name == name {
			if result != nil {
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple freeform allocation groups in blueprint %q have name %q", o.client.id, name),
				}
			}

			result = &ffrag
		}
	}

	if result == nil {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no freeform allocation group  in blueprint %q has name %q", o.client.id, name),
		}
	}

	return result, nil
}

func (o *FreeformClient) UpdateAllocGroup(ctx context.Context, id ObjectId, in *FreeformAllocGroupData) error {
	parts := strings.SplitN(id.String(), "_", 3)
	if parts[0] != "rag" || len(parts) != 3 {
		return ClientErr{
			errType: ErrInvalidId,
			err:     fmt.Errorf("freeform resource groups groupids must take the form rag_xxx_yyy, got : %q", id),
		}
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlFfAllocGroupByTypeName, o.blueprintId, parts[1], parts[2]),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

//func (o *FreeformClient) UpdateAllocGroups(ctx context.Context, in []FreeformAllocGroupData) error {
//
//	var apiInput struct {
//		ResourceGroups []json.RawMessage `json:"resource_groups"`
//	}
//	apiInput.ResourceGroups = make([]json.RawMessage, len(in))
//	var err error
//	for i, groupData := range in {
//		apiInput.ResourceGroups[i], err = json.Marshal(groupData)
//		if err != nil {
//			return err
//		}
//	}
//
//	err = o.client.talkToApstra(ctx, &talkToApstraIn{
//		method:   http.MethodPatch,
//		urlStr:   fmt.Sprintf(apiUrlFfAllocGroups, o.blueprintId),
//		apiInput: in,
//	})
//	if err != nil {
//		return convertTtaeToAceWherePossible(err)
//	}
//
//	return nil
//}

func (o *FreeformClient) DeleteAllocGroup(ctx context.Context, id ObjectId) error {
	parts := strings.SplitN(id.String(), "_", 3)
	if parts[0] != "rag" || len(parts) != 3 {
		return ClientErr{
			errType: ErrInvalidId,
			err:     fmt.Errorf("freeform resource groups groupids must take the form rag_xxx_yyy, got : %q", id),
		}
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfAllocGroupByTypeName, o.blueprintId, parts[1], parts[2]),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
