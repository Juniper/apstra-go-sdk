// Copyright (c) Juniper Networks, Inc., 2022-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Juniper/apstra-go-sdk/enum"
	resourcegroup "github.com/Juniper/apstra-go-sdk/resource_group"
)

const (
	apiUrlBlueprintResourceGroups        = apiUrlBlueprintById + apiUrlPathDelim + "resource_groups"
	apiUrlResourceGroupsPrefix           = apiUrlBlueprintResourceGroups + apiUrlPathDelim
	apiUrlBlueprintResourceGroupTypeName = apiUrlResourceGroupsPrefix + "%s" + apiUrlPathDelim + "%s"

	resourceGroupNameWithOwner     = "%s:%s,%s"
	resourceGroupOwnerSecurityZone = "sz"
)

type ResourceGroup struct {
	Type           enum.ResourceType
	Name           resourcegroup.ResourceGroup
	SecurityZoneId *string
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
	PoolIds       []string `json:"pool_ids"`
}

func (o *ResourceGroupAllocation) raw() *rawResourceGroupAllocation {
	var poolIds []string
	if o.PoolIds == nil {
		poolIds = make([]string, 0)
	} else {
		poolIds = o.PoolIds
	}

	// 'name' here can need a value like "leaf_loopback_ips" or
	// "sz:ISKtui8i80vl0ljsdJQ,leaf_loopback_ips", depending on whether the
	// resource group belongs to a blueprint or to a child object within a
	// blueprint.
	name := o.ResourceGroup.Name.String()
	if o.ResourceGroup.SecurityZoneId != nil {
		name = fmt.Sprintf(resourceGroupNameWithOwner, resourceGroupOwnerSecurityZone, *o.ResourceGroup.SecurityZoneId, name)
	}

	return &rawResourceGroupAllocation{
		Type:    o.ResourceGroup.Type,
		Name:    name,
		PoolIds: poolIds,
	}
}

func (o *ResourceGroupAllocation) IsEmpty() bool {
	return len(o.PoolIds) == 0
}

type rawResourceGroupAllocation struct {
	Type    enum.ResourceType `json:"type,omitempty"`
	Name    string            `json:"name,omitempty"`
	PoolIds []string          `json:"pool_ids"`
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
	rga := &ResourceGroupAllocation{
		PoolIds: o.PoolIds,
		ResourceGroup: ResourceGroup{
			Type: o.Type,
		},
	}

	switch {
	case strings.HasPrefix(o.Name, resourceGroupOwnerSecurityZone+":"):
		fields := strings.Split(o.Name, ",")
		if len(fields) != 2 {
			return nil, fmt.Errorf(
				"error processing resource group name %q, expected split on ',' to produce 2 results, got %d",
				o.Name, len(fields),
			)
		}
		err := rga.ResourceGroup.Name.UnmarshalText([]byte(fields[1]))
		if err != nil {
			return nil, err
		}
		szId := strings.TrimPrefix(fields[0], resourceGroupOwnerSecurityZone+":")
		rga.ResourceGroup.SecurityZoneId = &szId
	default:
		err := rga.ResourceGroup.Name.UnmarshalText([]byte(o.Name))
		if err != nil {
			return nil, err
		}
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
