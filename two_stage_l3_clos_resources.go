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

	// .../aos/reference_design/extension/resource_allocation/__init__.py says:
	// RESOURCE_TYPES = ['ip', 'ipv6', 'asn', 'vlan', 'vni']
	resourceTypeAsnPool = "asn"
	resourceTypeIp4Pool = "ip"
	resourceTypeIp6Pool = "ipv6"
	resourceTypeVniPool = "vni"
	resourceTypeUnknown = "resource type %d unknown"
)

const (
	ResourceGroupNameSuperspineAsn = ResourceGroupName(iota)
	ResourceGroupNameSpineAsn
	ResourceGroupNameLeafAsn
	ResourceGroupNameAccessAsn
	ResourceGroupNameSuperspineIp4
	ResourceGroupNameSpineIp4
	ResourceGroupNameLeafIp4
	ResourceGroupNameAccessIp4
	ResourceGroupNameSuperspineSpineIp4
	ResourceGroupNameSuperspineSpineIp6
	ResourceGroupNameSpineLeafIp4
	ResourceGroupNameSpineLeafIp6
	ResourceGroupNameAccessAccessIps
	ResourceGroupNameLeafLeafIp4
	ResourceGroupNameMlagDomainSviSubnets
	ResourceGroupNameVtepIps
	ResourceGroupNameUnknown

	resourceGroupNameSuperspineAsn        = resourceGroupName("superspine_asns")
	resourceGroupNameSpineAsn             = resourceGroupName("spine_asns")
	resourceGroupNameLeafAsn              = resourceGroupName("leaf_asns")
	resourceGroupNameAccessAsn            = resourceGroupName("access_asns")
	resourceGroupNameSuperspineIp4        = resourceGroupName("superspine_loopback_ips")
	resourceGroupNameSpineIp4             = resourceGroupName("spine_loopback_ips")
	resourceGroupNameLeafIp4              = resourceGroupName("leaf_loopback_ips")
	resourceGroupNameAccessIp4            = resourceGroupName("access_loopback_ips")
	resourceGroupNameSuperspineSpineIp4   = resourceGroupName("spine_superspine_link_ips")
	resourceGroupNameSuperspineSpineIp6   = resourceGroupName("ipv6_spine_superspine_link_ips")
	resourceGroupNameSpineLeafIp4         = resourceGroupName("spine_leaf_link_ips")
	resourceGroupNameSpineLeafIp6         = resourceGroupName("ipv6_spine_leaf_link_ips")
	resourceGroupNameLeafLeafIp4          = resourceGroupName("leaf_leaf_link_ips")
	resourceGroupNameMlagDomainSviSubnets = resourceGroupName("mlag_domain_svi_subnets")
	resourceGroupNameAccessAccessIps      = resourceGroupName("access_l3_peer_link_link_ips")
	resourceGroupNameVtepIps              = resourceGroupName("vtep_ips")
	resourceGroupNameUnknown              = "group name %d unknown"
)

type ResourceGroupName int

func (o ResourceGroupName) String() string {
	return string(o.raw())
}

func (o *ResourceGroupName) FromString(in string) error {
	i, err := resourceGroupName(in).parse()
	if err != nil {
		return err
	}
	*o = i
	return nil
}

func (o *ResourceGroupName) Type() ResourceType {
	switch *o {
	case ResourceGroupNameSuperspineAsn:
		return ResourceTypeAsnPool
	case ResourceGroupNameSpineAsn:
		return ResourceTypeAsnPool
	case ResourceGroupNameLeafAsn:
		return ResourceTypeAsnPool
	case ResourceGroupNameAccessAsn:
		return ResourceTypeAsnPool
	case ResourceGroupNameSuperspineIp4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameSpineIp4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameLeafIp4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameAccessIp4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameSuperspineSpineIp4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameSuperspineSpineIp6:
		return ResourceTypeIp6Pool
	case ResourceGroupNameSpineLeafIp4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameSpineLeafIp6:
		return ResourceTypeIp6Pool
	case ResourceGroupNameAccessAccessIps:
		return ResourceTypeIp4Pool
	case ResourceGroupNameLeafLeafIp4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameMlagDomainSviSubnets:
		return ResourceTypeIp4Pool
	case ResourceGroupNameVtepIps:
		return ResourceTypeIp4Pool
	}
	return ResourceTypeUnknown
}

// AllResourceGroupNames returns the []ResourceGroupName representing
// all supported ResourceGroupName
func AllResourceGroupNames() []ResourceGroupName {
	i := 0
	var result []ResourceGroupName
	for {
		var rgn ResourceGroupName
		err := rgn.FromString(ResourceGroupName(i).String())
		if err != nil {
			return result[:i]
		}
		result = append(result, rgn)
		i++
	}
}

func (o ResourceGroupName) raw() resourceGroupName {
	switch o {
	case ResourceGroupNameSuperspineAsn:
		return resourceGroupNameSuperspineAsn
	case ResourceGroupNameSpineAsn:
		return resourceGroupNameSpineAsn
	case ResourceGroupNameLeafAsn:
		return resourceGroupNameLeafAsn
	case ResourceGroupNameAccessAsn:
		return resourceGroupNameAccessAsn
	case ResourceGroupNameSuperspineIp4:
		return resourceGroupNameSuperspineIp4
	case ResourceGroupNameSpineIp4:
		return resourceGroupNameSpineIp4
	case ResourceGroupNameLeafIp4:
		return resourceGroupNameLeafIp4
	case ResourceGroupNameAccessIp4:
		return resourceGroupNameAccessIp4
	case ResourceGroupNameSuperspineSpineIp4:
		return resourceGroupNameSuperspineSpineIp4
	case ResourceGroupNameSuperspineSpineIp6:
		return resourceGroupNameSuperspineSpineIp6
	case ResourceGroupNameSpineLeafIp4:
		return resourceGroupNameSpineLeafIp4
	case ResourceGroupNameSpineLeafIp6:
		return resourceGroupNameSpineLeafIp6
	case ResourceGroupNameAccessAccessIps:
		return resourceGroupNameAccessAccessIps
	case ResourceGroupNameLeafLeafIp4:
		return resourceGroupNameLeafLeafIp4
	case ResourceGroupNameMlagDomainSviSubnets:
		return resourceGroupNameMlagDomainSviSubnets
	case ResourceGroupNameVtepIps:
		return resourceGroupNameVtepIps
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
	case resourceGroupNameAccessAsn:
		return ResourceGroupNameAccessAsn, nil
	case resourceGroupNameSuperspineIp4:
		return ResourceGroupNameSuperspineIp4, nil
	case resourceGroupNameSpineIp4:
		return ResourceGroupNameSpineIp4, nil
	case resourceGroupNameLeafIp4:
		return ResourceGroupNameLeafIp4, nil
	case resourceGroupNameAccessIp4:
		return ResourceGroupNameAccessIp4, nil
	case resourceGroupNameSuperspineSpineIp4:
		return ResourceGroupNameSuperspineSpineIp4, nil
	case resourceGroupNameSuperspineSpineIp6:
		return ResourceGroupNameSuperspineSpineIp6, nil
	case resourceGroupNameSpineLeafIp4:
		return ResourceGroupNameSpineLeafIp4, nil
	case resourceGroupNameSpineLeafIp6:
		return ResourceGroupNameSpineLeafIp6, nil
	case resourceGroupNameAccessAccessIps:
		return ResourceGroupNameAccessAccessIps, nil
	case resourceGroupNameLeafLeafIp4:
		return ResourceGroupNameLeafLeafIp4, nil
	case resourceGroupNameMlagDomainSviSubnets:
		return ResourceGroupNameMlagDomainSviSubnets, nil
	case resourceGroupNameVtepIps:
		return ResourceGroupNameVtepIps, nil
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

func (o *ResourceType) FromString(in string) error {
	i, err := resourceType(in).parse()
	if err != nil {
		return err
	}
	*o = i
	return nil
}

// AllResourceTypes returns the []ResourceType representing
// all supported ResourceType
func AllResourceTypes() []ResourceType {
	i := 0
	var result []ResourceType
	for {
		var rgn ResourceType
		err := rgn.FromString(ResourceType(i).String())
		if err != nil {
			return result[:i]
		}
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

type ResourceGroup struct {
	Type ResourceType      `json:"type"`
	Name ResourceGroupName `json:"name"`
}

type ResourceGroupAllocations []ResourceGroupAllocation

// Get returns the ResourceGroupAllocation for the requested ResourceGroup, or nil
// if no matching ResourceGroupAllocation exists in this ResourceGroupAllocations
func (o ResourceGroupAllocations) Get(requested *ResourceGroup) *ResourceGroupAllocation {
	for _, rg := range o {
		if rg.ResourceGroup.Name == requested.Name && rg.ResourceGroup.Type == requested.Type {
			return &rg
		}
	}
	return nil
}

type ResourceGroupAllocation struct {
	ResourceGroup ResourceGroup
	PoolIds       []ObjectId `json:"pool_ids"`
}

func (o *ResourceGroupAllocation) raw() *rawResourceGroupAllocation {
	var poolIds []ObjectId
	if o.PoolIds == nil {
		poolIds = make([]ObjectId, 0)
	} else {
		poolIds = o.PoolIds
	}

	return &rawResourceGroupAllocation{
		Type:    o.ResourceGroup.Type.raw(),
		Name:    o.ResourceGroup.Name.raw(),
		PoolIds: poolIds,
	}
}

func (o *ResourceGroupAllocation) IsEmpty() bool {
	return len(o.PoolIds) == 0
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
		ResourceGroup: ResourceGroup{
			Type: t,
			Name: n,
		},
		PoolIds: o.PoolIds,
	}, nil
}

func (o *TwoStageL3ClosClient) getAllResourceAllocations(ctx context.Context) ([]rawResourceGroupAllocation, error) {
	response := &struct {
		Items []rawResourceGroupAllocation `json:"items"`
	}{}
	return response.Items, o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintResourceGroups, o.blueprintId),
		apiResponse: response,
	})
}

func (o *TwoStageL3ClosClient) getResourceAllocation(ctx context.Context, rg *ResourceGroup) (*rawResourceGroupAllocation, error) {
	response := &rawResourceGroupAllocation{}
	ttii := talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintResourceGroupTypeName, o.blueprintId, rg.Type.String(), rg.Name.String()),
		apiResponse: response,
	}
	err := o.client.talkToApstra(ctx, &ttii)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *TwoStageL3ClosClient) setResourceAllocation(ctx context.Context, rga *ResourceGroupAllocation) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintResourceGroupTypeName, o.blueprintId, rga.ResourceGroup.Type.String(), rga.ResourceGroup.Name.String()),
		apiInput: rga.raw(),
	})
}
