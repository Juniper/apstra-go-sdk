package apstra

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const (
	apiUrlBlueprintResourceGroups        = apiUrlBlueprintById + apiUrlPathDelim + "resource_groups"
	apiUrlResourceGroupsPrefix           = apiUrlBlueprintResourceGroups + apiUrlPathDelim
	apiUrlBlueprintResourceGroupTypeName = apiUrlResourceGroupsPrefix + "%s" + apiUrlPathDelim + "%s"

	resourceGroupNameWithOwner    = "%s:%s,%s"
	resourceGroupOwnerecurityZone = "sz"
)

const (
	ResourceTypeNone = ResourceType(iota)
	ResourceTypeAsnPool
	ResourceTypeIp4Pool
	ResourceTypeIp6Pool
	ResourceTypeVniPool
	ResourceTypeUnknown

	// .../aos/reference_design/extension/resource_allocation/__init__.py says:
	// RESOURCE_TYPES = ['ip', 'ipv6', 'asn', 'vlan', 'vni']
	resourceTypeNone    = resourceType("")
	resourceTypeAsnPool = resourceType("asn")
	resourceTypeIp4Pool = resourceType("ip")
	resourceTypeIp6Pool = resourceType("ipv6")
	resourceTypeVniPool = resourceType("vni")
	resourceTypeUnknown = "resource type %d unknown"
)

const (
	ResourceGroupNameNone = ResourceGroupName(iota)
	ResourceGroupNameSuperspineAsn
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
	ResourceGroupNameAccessAccessIp4
	ResourceGroupNameLeafLeafIp4
	ResourceGroupNameLeafL3PeerLinkLinkIps
	ResourceGroupNameMlagDomainIp4
	ResourceGroupNameVtepIp4
	ResourceGroupNameEvpnL3Vni
	ResourceGroupNameVirtualNetworkSviIpv4
	ResourceGroupNameVirtualNetworkSviIpv6
	ResourceGroupNameVxlanVnIds
	ResourceGroupNameUnknown

	resourceGroupNameNone                  = resourceGroupName("")
	resourceGroupNameSuperspineAsn         = resourceGroupName("superspine_asns")
	resourceGroupNameSpineAsn              = resourceGroupName("spine_asns")
	resourceGroupNameLeafAsn               = resourceGroupName("leaf_asns")
	resourceGroupNameAccessAsn             = resourceGroupName("access_asns")
	resourceGroupNameSuperspineIp4         = resourceGroupName("superspine_loopback_ips")
	resourceGroupNameSpineIp4              = resourceGroupName("spine_loopback_ips")
	resourceGroupNameLeafIp4               = resourceGroupName("leaf_loopback_ips")
	resourceGroupNameAccessIp4             = resourceGroupName("access_loopback_ips")
	resourceGroupNameSuperspineSpineIp4    = resourceGroupName("spine_superspine_link_ips")
	resourceGroupNameSuperspineSpineIp6    = resourceGroupName("ipv6_spine_superspine_link_ips")
	resourceGroupNameSpineLeafIp4          = resourceGroupName("spine_leaf_link_ips")
	resourceGroupNameSpineLeafIp6          = resourceGroupName("ipv6_spine_leaf_link_ips")
	resourceGroupNameLeafLeafIp4           = resourceGroupName("leaf_leaf_link_ips")
	resourceGroupNameLeafL3PeerLinkLinkIps = resourceGroupName("leaf_l3_peer_link_link_ips")
	resourceGroupNameMlagDomainSviIp4      = resourceGroupName("mlag_domain_svi_subnets")
	resourceGroupNameAccessAccessIp4       = resourceGroupName("access_l3_peer_link_link_ips")
	resourceGroupNameVtepIp4               = resourceGroupName("vtep_ips")
	resourceGroupNameEvpnL3Vni             = resourceGroupName("evpn_l3_vnis")
	resourceGroupNameVirtualNetworkSviIpv4 = resourceGroupName("virtual_network_svi_subnets")
	resourceGroupNameVirtualNetworkSviIpv6 = resourceGroupName("virtual_network_svi_subnets_ipv6")
	resourceGroupNameVxlanVnIds            = resourceGroupName("vxlan_vn_ids")
	resourceGroupNameUnknown               = "group name %d unknown"
)

type ResourceGroupName int

func (o ResourceGroupName) String() string {
	return string(o.raw())
}

func (o ResourceGroupName) Int() int {
	return int(o)
}

func (o *ResourceGroupName) FromString(in string) error {
	i, err := resourceGroupName(in).parse()
	if err != nil {
		return err
	}
	*o = ResourceGroupName(i)
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
	case ResourceGroupNameAccessAccessIp4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameLeafLeafIp4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameLeafL3PeerLinkLinkIps:
		return ResourceTypeIp4Pool
	case ResourceGroupNameMlagDomainIp4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameVtepIp4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameEvpnL3Vni:
		return ResourceTypeVniPool
	case ResourceGroupNameVirtualNetworkSviIpv4:
		return ResourceTypeIp4Pool
	case ResourceGroupNameVirtualNetworkSviIpv6:
		return ResourceTypeIp6Pool
	case ResourceGroupNameVxlanVnIds:
		return ResourceTypeVniPool
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
	case ResourceGroupNameNone:
		return resourceGroupNameNone
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
	case ResourceGroupNameAccessAccessIp4:
		return resourceGroupNameAccessAccessIp4
	case ResourceGroupNameLeafLeafIp4:
		return resourceGroupNameLeafLeafIp4
	case ResourceGroupNameLeafL3PeerLinkLinkIps:
		return resourceGroupNameLeafL3PeerLinkLinkIps
	case ResourceGroupNameMlagDomainIp4:
		return resourceGroupNameMlagDomainSviIp4
	case ResourceGroupNameVtepIp4:
		return resourceGroupNameVtepIp4
	case ResourceGroupNameEvpnL3Vni:
		return resourceGroupNameEvpnL3Vni
	case ResourceGroupNameVirtualNetworkSviIpv4:
		return resourceGroupNameVirtualNetworkSviIpv4
	case ResourceGroupNameVirtualNetworkSviIpv6:
		return resourceGroupNameVirtualNetworkSviIpv6
	case ResourceGroupNameVxlanVnIds:
		return resourceGroupNameVxlanVnIds
	default:
		return resourceGroupName(fmt.Sprintf(resourceGroupNameUnknown, o))
	}
}

type resourceGroupName string

func (o resourceGroupName) parse() (int, error) {
	switch o {
	case resourceGroupNameNone:
		return int(ResourceGroupNameNone), nil
	case resourceGroupNameSuperspineAsn:
		return int(ResourceGroupNameSuperspineAsn), nil
	case resourceGroupNameSpineAsn:
		return int(ResourceGroupNameSpineAsn), nil
	case resourceGroupNameLeafAsn:
		return int(ResourceGroupNameLeafAsn), nil
	case resourceGroupNameAccessAsn:
		return int(ResourceGroupNameAccessAsn), nil
	case resourceGroupNameSuperspineIp4:
		return int(ResourceGroupNameSuperspineIp4), nil
	case resourceGroupNameSpineIp4:
		return int(ResourceGroupNameSpineIp4), nil
	case resourceGroupNameLeafIp4:
		return int(ResourceGroupNameLeafIp4), nil
	case resourceGroupNameAccessIp4:
		return int(ResourceGroupNameAccessIp4), nil
	case resourceGroupNameSuperspineSpineIp4:
		return int(ResourceGroupNameSuperspineSpineIp4), nil
	case resourceGroupNameSuperspineSpineIp6:
		return int(ResourceGroupNameSuperspineSpineIp6), nil
	case resourceGroupNameSpineLeafIp4:
		return int(ResourceGroupNameSpineLeafIp4), nil
	case resourceGroupNameSpineLeafIp6:
		return int(ResourceGroupNameSpineLeafIp6), nil
	case resourceGroupNameAccessAccessIp4:
		return int(ResourceGroupNameAccessAccessIp4), nil
	case resourceGroupNameLeafLeafIp4:
		return int(ResourceGroupNameLeafLeafIp4), nil
	case resourceGroupNameLeafL3PeerLinkLinkIps:
		return int(ResourceGroupNameLeafL3PeerLinkLinkIps), nil
	case resourceGroupNameMlagDomainSviIp4:
		return int(ResourceGroupNameMlagDomainIp4), nil
	case resourceGroupNameVtepIp4:
		return int(ResourceGroupNameVtepIp4), nil
	case resourceGroupNameEvpnL3Vni:
		return int(ResourceGroupNameEvpnL3Vni), nil
	case resourceGroupNameVirtualNetworkSviIpv4:
		return int(ResourceGroupNameVirtualNetworkSviIpv4), nil
	case resourceGroupNameVirtualNetworkSviIpv6:
		return int(ResourceGroupNameVirtualNetworkSviIpv6), nil
	case resourceGroupNameVxlanVnIds:
		return int(ResourceGroupNameVxlanVnIds), nil
	default:
		return int(ResourceGroupNameUnknown), fmt.Errorf("unknown group name '%s'", o)
	}
}

func (o resourceGroupName) string() string {
	return string(o)
}

type ResourceType int

func (o ResourceType) String() string {
	return string(o.raw())
}

func (o ResourceType) Int() int {
	return int(o)
}

func (o ResourceType) raw() resourceType {
	switch o {
	case ResourceTypeNone:
		return resourceTypeNone
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
	*o = ResourceType(i)
	return nil
}

// AllResourceTypes returns the []ResourceType representing
// all supported ResourceType
func AllResourceTypes() []ResourceType {
	i := 0
	var result []ResourceType
	for {
		var rt ResourceType
		err := rt.FromString(ResourceType(i).String())
		if err != nil {
			return result[:i]
		}
		result = append(result, rt)
		i++
	}
}

type resourceType string

func (o resourceType) string() string {
	return string(o)
}

func (o resourceType) parse() (int, error) {
	switch o {
	case resourceTypeNone:
		return int(ResourceTypeNone), nil
	case resourceTypeAsnPool:
		return int(ResourceTypeAsnPool), nil
	case resourceTypeIp4Pool:
		return int(ResourceTypeIp4Pool), nil
	case resourceTypeIp6Pool:
		return int(ResourceTypeIp6Pool), nil
	case resourceTypeVniPool:
		return int(ResourceTypeVniPool), nil
	default:
		return int(ResourceTypeUnknown), fmt.Errorf("unknown resource type '%s'", o)
	}
}

type ResourceGroup struct {
	Type           ResourceType
	Name           ResourceGroupName
	SecurityZoneId *ObjectId
}

type ResourceGroupAllocations []ResourceGroupAllocation

// Get returns the ResourceGroupAllocation for the requested ResourceGroup, or nil
// if no matching ResourceGroupAllocation exists in this ResourceGroupAllocations
func (o ResourceGroupAllocations) Get(requested *ResourceGroup) *ResourceGroupAllocation {
	for _, rg := range o {
		if rg.ResourceGroup.Type != requested.Type {
			continue
		}
		if rg.ResourceGroup.Name != requested.Name {
			continue
		}
		if (rg.ResourceGroup.SecurityZoneId != nil && requested.SecurityZoneId != nil) &&
			*rg.ResourceGroup.SecurityZoneId != *requested.SecurityZoneId {
			continue
		}
		return &rg
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

	// 'name' here can need a value like "leaf_loopback_ips" or
	// "sz:ISKtui8i80vl0ljsdJQ,leaf_loopback_ips", depending on whether the
	// resource group belongs to a blueprint or to a child object within a
	// blueprint.
	name := o.ResourceGroup.Name.raw()
	switch {
	case o.ResourceGroup.SecurityZoneId != nil:
		name = resourceGroupName(fmt.Sprintf(
			resourceGroupNameWithOwner, resourceGroupOwnerecurityZone,
			*o.ResourceGroup.SecurityZoneId, name))
	}

	return &rawResourceGroupAllocation{
		Type:    o.ResourceGroup.Type.raw(),
		Name:    name,
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

// polish leans on some apstra code which determines whether a resource group
// name indicates that it's owned by a child object within a blueprint (security
// zone may be the only case of this):
//
//	def parse_resource_group_name(resource_group_name):
//	   """ Parse passed resource_group_name and return namedtuple with
//	       sz_id and pure resource_group_name.
//	   """
//	   sz_id = None
//	   if resource_group_name.startswith('sz:'):
//	       fields = resource_group_name.split(',', 1)
//	       if len(fields) != 2:
//	           return None
//	       sz, resource_group_name = fields
//	       sz_id = sz[len('sz:'):]
//	   return ParsedResourceGroupName(sz_id=sz_id, rg_name=resource_group_name)
func (o *rawResourceGroupAllocation) polish() (*ResourceGroupAllocation, error) {
	t, err := o.Type.parse()
	if err != nil {
		return nil, err
	}

	rga := &ResourceGroupAllocation{
		PoolIds: o.PoolIds,
		ResourceGroup: ResourceGroup{
			Type: ResourceType(t),
		},
	}

	switch {
	case strings.HasPrefix(string(o.Name), resourceGroupOwnerecurityZone+":"):
		fields := strings.Split(string(o.Name), ",")
		if len(fields) != 2 {
			return nil, fmt.Errorf(
				"error processing resource group name %q, expected split on ',' to produce 2 results, got %d",
				o.Name, len(fields))
		}
		err = rga.ResourceGroup.Name.FromString(fields[1])
		if err != nil {
			return nil, err
		}
		szId := ObjectId(strings.TrimPrefix(fields[0], resourceGroupOwnerecurityZone+":"))
		rga.ResourceGroup.SecurityZoneId = &szId
	default:
		name, err := o.Name.parse()
		if err != nil {
			return nil, err
		}
		rga.ResourceGroup.Name = ResourceGroupName(name)
	}

	return rga, nil
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
	rga := ResourceGroupAllocation{
		ResourceGroup: *rg,
	}
	response := rga.raw()
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintResourceGroupTypeName, o.blueprintId, response.Type, response.Name),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *TwoStageL3ClosClient) setResourceAllocation(ctx context.Context, rga *rawResourceGroupAllocation) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintResourceGroupTypeName, o.blueprintId, rga.Type, rga.Name),
		apiInput: rga,
	})
}
