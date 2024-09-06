package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	apiUrlFfResourceGenerators    = apiUrlBlueprintById + apiUrlPathDelim + "ra-resource-generators"
	apiUrlFfResourceGeneratorById = apiUrlFfResourceGenerators + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformResourceGenerator)

type FreeformResourceGenerator struct {
	Id   ObjectId
	Data *FreeformResourceGeneratorData
}

func (o *FreeformResourceGenerator) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id                 ObjectId  `json:"id"`
		ResourceType       string    `json:"resource_type"`
		Label              string    `json:"label"`
		Scope              string    `json:"scope"`
		AllocatedFrom      *ObjectId `json:"allocated_from"`
		ScopeNodePoolLabel *string   `json:"scope_node_pool_label"`
		ContainerId        ObjectId  `json:"container_id"`
		SubnetPrefixLen    *int      `json:"subnet_prefix_len"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(FreeformResourceGeneratorData)
	o.Data.Label = raw.Label
	o.Data.Scope = raw.Scope
	o.Data.AllocatedFrom = raw.AllocatedFrom
	o.Data.ScopeNodePoolLabel = raw.ScopeNodePoolLabel
	o.Data.ContainerId = raw.ContainerId
	o.Data.SubnetPrefixLen = raw.SubnetPrefixLen
	err = o.Data.ResourceType.FromString(raw.ResourceType)
	if err != nil {
		return err
	}

	return nil
}

var _ json.Marshaler = new(FreeformResourceGeneratorData)

type FreeformResourceGeneratorData struct {
	ResourceType       enum.FFResourceType
	Label              string
	Scope              string
	AllocatedFrom      *ObjectId
	ScopeNodePoolLabel *string
	ContainerId        ObjectId
	SubnetPrefixLen    *int
}

func (o FreeformResourceGeneratorData) MarshalJSON() ([]byte, error) {
	var raw struct {
		ResourceType       string    `json:"resource_type,omitempty"`
		Label              string    `json:"label,omitempty"`
		Scope              string    `json:"scope,omitempty"`
		AllocatedFrom      *ObjectId `json:"allocated_from,omitempty"`
		ScopeNodePoolLabel *string   `json:"scope_node_pool_label,omitempty"`
		ContainerId        ObjectId  `json:"container_id,omitempty"`
		SubnetPrefixLen    *int      `json:"subnet_prefix_len,omitempty"`
	}

	raw.ResourceType = o.ResourceType.String()
	raw.Label = o.Label
	raw.Scope = o.Scope
	raw.AllocatedFrom = o.AllocatedFrom
	raw.ScopeNodePoolLabel = o.ScopeNodePoolLabel
	raw.ContainerId = o.ContainerId
	raw.SubnetPrefixLen = o.SubnetPrefixLen

	return json.Marshal(&raw)
}

func (o *FreeformClient) CreateResourceGenerator(ctx context.Context, in *FreeformResourceGeneratorData) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFfResourceGenerators, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *FreeformClient) GetAllResourceGenerators(ctx context.Context) ([]FreeformResourceGenerator, error) {
	var response struct {
		Items []FreeformResourceGenerator `json:"items"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfResourceGenerators, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) GetResourceGenerator(ctx context.Context, id ObjectId) (*FreeformResourceGenerator, error) {
	var response FreeformResourceGenerator

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfResourceGeneratorById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) GetResourceGeneratorByName(ctx context.Context, name string) (*FreeformResourceGenerator, error) {
	all, err := o.GetAllResourceGenerators(ctx)
	if err != nil {
		return nil, err
	}

	var result *FreeformResourceGenerator
	for _, ffResGen := range all {
		ffResGen := ffResGen
		if ffResGen.Data.Label == name {
			if result != nil {
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple freeform resource generators in blueprint %q have name %q", o.client.id, name),
				}
			}

			result = &ffResGen
		}
	}

	if result == nil {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no freeform resource generator in blueprint %q has name %q", o.client.id, name),
		}
	}

	return result, nil
}

func (o *FreeformClient) UpdateResourceGenerator(ctx context.Context, id ObjectId, in *FreeformResourceGeneratorData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFfResourceGeneratorById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) DeleteResourceGenerator(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfResourceGeneratorById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
