package goapstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintResourceGroups        = apiUrlBlueprintById + apiUrlPathDelim + "resource_groups"
	apiUrlResourceGroupsPrefix           = apiUrlBlueprintResourceGroups + apiUrlPathDelim
	apiUrlBlueprintResourceGroupTypeName = apiUrlResourceGroupsPrefix + apiUrlPathDelim + "%s" + apiUrlPathDelim + "%s"
)

const (
	ResourceTypeAsnPool = ResourceType(iota)
	ResourceTypeIp4Pool
	ResourceTypeIp6Pool
	ResourceTypeVniPool
	ResourceTypeUnknown

	resourceTypeAsnPool = "asn"
	resourceTypeIp4Pool = "x" // todo
	resourceTypeIp6Pool = "y" // todo
	resourceTypeVniPool = "z" // todo
	resourceTypeUnknown = "resource type %d unknown"
)

const (
	ResourceGroupNameSpineAsn = ResourceGroupName(iota)
	ResourceGroupNameUnknown

	resourceGroupNameSpineAsn = resourceGroupName("spine_asns")
	resourceGroupNameUnknown  = "group name %d unknown"
)

type ResourceGroupName int

func (o ResourceGroupName) String() string {
	return string(o.raw())
}

func (o ResourceGroupName) raw() resourceGroupName {
	switch o {
	case ResourceGroupNameSpineAsn:
		return resourceGroupNameSpineAsn
	default:
		return resourceGroupName(fmt.Sprintf(resourceGroupNameUnknown, o))
	}
}

type resourceGroupName string

func (o resourceGroupName) parse() (ResourceGroupName, error) {
	switch o {
	case resourceGroupNameSpineAsn:
		return ResourceGroupNameSpineAsn, nil
	default:
		return ResourceGroupNameUnknown, fmt.Errorf("unknown group name '%s'", o)
	}
}

type ResourceType int

func (o ResourceType) String() string {
	return string(o.raw())
}

func (o ResourceType) raw() resourceType {
	switch o {
	case ResourceTypeAsnPool:
		return resourceTypeAsnPool
	case ResourceTypeIp4Pool:
		return resourceTypeIp4Pool
	case ResourceTypeIp6Pool:
		return resourceTypeIp6Pool
	case ResourceTypeVniPool:
		return resourceTypeVniPool
	default:
		return resourceType(fmt.Sprintf(resourceTypeUnknown, o))
	}
}

type resourceType string

func (o resourceType) parse() (ResourceType, error) {
	switch o {
	case resourceTypeAsnPool:
		return ResourceTypeAsnPool, nil
	case resourceTypeIp4Pool:
		return ResourceTypeIp4Pool, nil
	case resourceTypeIp6Pool:
		return ResourceTypeIp6Pool, nil
	default:
		return ResourceTypeUnknown, fmt.Errorf("unknown resource type '%s'", o)
	}
}

type ResourceGroupAllocation struct {
	Type    ResourceType      `json:"type"`
	Name    ResourceGroupName `json:"name"`
	PoolIds []ObjectId        `json:"pool_ids"`
}

func (o *ResourceGroupAllocation) raw() *rawResourceGroupAllocation {
	return &rawResourceGroupAllocation{
		Type:    o.Type.raw(),
		Name:    o.Name.raw(),
		PoolIds: o.PoolIds,
	}
}

type rawResourceGroupAllocation struct {
	Type    resourceType      `json:"type,omitempty"`
	Name    resourceGroupName `json:"name,omitempty"`
	PoolIds []ObjectId        `json:"pool_ids"`
}

func (o *Blueprint) getResourceAllocation(ctx context.Context, in *ResourceGroupAllocation) (*ResourceGroupAllocation, error) {
	response := &ResourceGroupAllocation{}
	return response, o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintResourceGroupTypeName, o.Id, in.Type, in.Name),
		apiResponse: response,
	})
}

func (o *Blueprint) setResourceAllocation(ctx context.Context, in *ResourceGroupAllocation) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintResourceGroupTypeName, o.Id, in.Type, in.Name),
		apiInput: in.raw(),
	})
}
