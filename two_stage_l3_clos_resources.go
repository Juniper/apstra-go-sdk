package goapstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintResourceGroups        = apiUrlBlueprintById + apiUrlPathDelim + "resource_groups"
	apiUrlResourceGroupsPrefix           = apiUrlBlueprintResourceGroups + apiUrlPathDelim
	apiUrlBlueprintResourceGroupTypeName = apiUrlResourceGroupsPrefix + "%s" + apiUrlPathDelim + "%s"
)

const (
	ResourceTypeAsnPool = ResourceType(iota)
	ResourceTypeIp4Pool
	ResourceTypeIp6Pool
	ResourceTypeVniPool
	ResourceTypeUnknown

	resourceTypeAsnPool = "asn"
	resourceTypeIp4Pool = "ip"
	resourceTypeIp6Pool = "y" // todo
	resourceTypeVniPool = "z" // todo
	resourceTypeUnknown = "resource type %d unknown"
)

const (
	ResourceGroupNameSuperspineAsn = ResourceGroupName(iota)
	ResourceGroupNameSpineAsn
	ResourceGroupNameLeafAsn
	ResourceGroupNameSuperspineIps
	ResourceGroupNameSpineIps
	ResourceGroupNameLeafIps
	ResourceGroupNameSuperspineSpineIps
	ResourceGroupNameSpineLeafIps
	ResourceGroupNameMlagDomainSviSubnets
	ResourceGroupNameUnknown

	resourceGroupNameSuperspineAsn        = resourceGroupName("superspine_asns")
	resourceGroupNameSpineAsn             = resourceGroupName("spine_asns")
	resourceGroupNameLeafAsn              = resourceGroupName("leaf_asns")
	resourceGroupNameSuperspineIps        = resourceGroupName("superspine_loopback_ips")
	resourceGroupNameSpineIps             = resourceGroupName("spine_loopback_ips")
	resourceGroupNameLeafIps              = resourceGroupName("leaf_loopback_ips")
	resourceGroupNameSuperspineSpineIps   = resourceGroupName("spine_superspine_link_ips")
	resourceGroupNameSpineLeafIps         = resourceGroupName("spine_leaf_link_ips")
	resourceGroupNameMlagDomainSviSubnets = resourceGroupName("mlag_domain_svi_subnets")
	resourceGroupNameUnknown              = "group name %d unknown"
)

type ResourceGroupName int

func (o ResourceGroupName) String() string {
	return string(o.raw())
}

func (o ResourceGroupName) raw() resourceGroupName {
	switch o {
	case ResourceGroupNameSuperspineAsn:
		return resourceGroupNameSuperspineAsn
	case ResourceGroupNameSpineAsn:
		return resourceGroupNameSpineAsn
	case ResourceGroupNameLeafAsn:
		return resourceGroupNameLeafAsn
	case ResourceGroupNameSuperspineIps:
		return resourceGroupNameSuperspineIps
	case ResourceGroupNameSpineIps:
		return resourceGroupNameSpineIps
	case ResourceGroupNameLeafIps:
		return resourceGroupNameLeafIps
	case ResourceGroupNameSuperspineSpineIps:
		return resourceGroupNameSuperspineSpineIps
	case ResourceGroupNameSpineLeafIps:
		return resourceGroupNameSpineLeafIps
	case ResourceGroupNameMlagDomainSviSubnets:
		return resourceGroupNameMlagDomainSviSubnets
	default:
		return resourceGroupName(fmt.Sprintf(resourceGroupNameUnknown, o))
	}
}

type resourceGroupName string

func (o resourceGroupName) parse() (ResourceGroupName, error) {
	switch o {
	case resourceGroupNameSuperspineAsn:
		return ResourceGroupNameSuperspineAsn, nil
	case resourceGroupNameSpineAsn:
		return ResourceGroupNameSpineAsn, nil
	case resourceGroupNameLeafAsn:
		return ResourceGroupNameLeafAsn, nil
	case resourceGroupNameSuperspineIps:
		return ResourceGroupNameSuperspineIps, nil
	case resourceGroupNameSpineIps:
		return ResourceGroupNameSpineIps, nil
	case resourceGroupNameLeafIps:
		return ResourceGroupNameLeafIps, nil
	case resourceGroupNameSuperspineSpineIps:
		return ResourceGroupNameSuperspineSpineIps, nil
	case resourceGroupNameSpineLeafIps:
		return ResourceGroupNameSpineLeafIps, nil
	case resourceGroupNameMlagDomainSviSubnets:
		return ResourceGroupNameMlagDomainSviSubnets, nil
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

func (o *rawResourceGroupAllocation) polish() (*ResourceGroupAllocation, error) {
	t, err := o.Type.parse()
	if err != nil {
		return nil, err
	}

	n, err := o.Name.parse()
	if err != nil {
		return nil, err
	}

	return &ResourceGroupAllocation{
		Type:    t,
		Name:    n,
		PoolIds: o.PoolIds,
	}, nil
}

func (o *TwoStageL3ClosClient) getResourceAllocation(ctx context.Context, rga *ResourceGroupAllocation) (*ResourceGroupAllocation, error) {
	response := &rawResourceGroupAllocation{}
	ttii := talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintResourceGroupTypeName, o.blueprintId, rga.Type.String(), rga.Name.String()),
		apiResponse: response,
	}
	err := o.client.talkToApstra(ctx, &ttii)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.polish()
}

func (o *TwoStageL3ClosClient) setResourceAllocation(ctx context.Context, rga *ResourceGroupAllocation) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintResourceGroupTypeName, o.blueprintId, rga.Type.String(), rga.Name.String()),
		apiInput: rga.raw(),
	})
}
