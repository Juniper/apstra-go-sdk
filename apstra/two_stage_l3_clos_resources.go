// Copyright (c) Juniper Networks, Inc., 2022-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"bytes"
	"context"
	"encoding"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"strings"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/urls"
)

const (
	resourceGroupNameWithOwner          = "%s%s,%s"
	resourceGroupOwnerRoutingZonePrefix = "sz:"
)

var (
	_ encoding.TextMarshaler   = (*ResourceGroup)(nil)
	_ encoding.TextUnmarshaler = (*ResourceGroup)(nil)
)

type ResourceGroup struct {
	Name           enum.ResourceGroup
	SecurityZoneID *string
}

func (o ResourceGroup) String() string {
	s, _ := o.MarshalText() // cannot error
	return string(s)
}

// MarshalText returns a string representation of the resource group, which may be either a
// simple string (e.g. "leaf_loopback_ips") or a string prefixed with the owner type and id
// (e.g. "sz:ISKtui8i80vl0ljsdJQ,leaf_loopback_ips"), depending on whether the resource
// group belongs directly to a blueprint or to an object within the blueprint.
func (o ResourceGroup) MarshalText() (text []byte, err error) {
	// 'name' here can need a value like "leaf_loopback_ips" or
	// "sz:ISKtui8i80vl0ljsdJQ,leaf_loopback_ips"
	switch { // only one case (so far?)
	case o.SecurityZoneID != nil:
		return []byte(fmt.Sprintf(resourceGroupNameWithOwner, resourceGroupOwnerRoutingZonePrefix, *o.SecurityZoneID, o.Name)), nil
	}

	return []byte(o.Name.String()), nil
}

func (o *ResourceGroup) UnmarshalText(b []byte) error {
	switch {
	case bytes.HasPrefix(b, []byte(resourceGroupOwnerRoutingZonePrefix)):
		fields := bytes.Split(b, []byte(","))
		if len(fields) != 2 {
			return fmt.Errorf("parsing resource group name %q: expected split on ',' to produce 2 results, got %d", string(b), len(fields))
		}

		err := o.Name.FromString(string(fields[1]))
		if err != nil {
			return fmt.Errorf("parsing resource group name %q: %w", string(b), err)
		}
		o.SecurityZoneID = pointer.To(strings.TrimPrefix(string(fields[0]), resourceGroupOwnerRoutingZonePrefix))
	default:
		err := o.Name.FromString(string(b))
		if err != nil {
			return fmt.Errorf("parsing resource group name %q: %w", string(b), err)
		}
		o.SecurityZoneID = nil
	}

	return nil
}

type ResourceGroupAllocations []ResourceGroupAllocation

// Get returns the ResourceGroupAllocation for the requested ResourceGroup, or nil
// if no matching ResourceGroupAllocation exists in this ResourceGroupAllocations
func (o ResourceGroupAllocations) Get(requested ResourceGroup) ResourceGroupAllocation {
	result := ResourceGroupAllocation{ResourceGroup: requested}

	for _, rg := range o {
		if rg.ResourceGroup.Name != requested.Name {
			continue
		}
		if (rg.ResourceGroup.SecurityZoneID != nil && requested.SecurityZoneID != nil) &&
			*rg.ResourceGroup.SecurityZoneID != *requested.SecurityZoneID {
			continue
		}
		result.PoolIds = append(result.PoolIds, rg.PoolIds...)
	}

	return result
}

var _ json.Marshaler = (*ResourceGroupAllocation)(nil)

type ResourceGroupAllocation struct {
	ResourceGroup ResourceGroup `json:"name"`
	PoolIds       []string      `json:"pool_ids"`
}

func (o ResourceGroupAllocation) MarshalJSON() ([]byte, error) {
	// Send empty pool_ids rather than null.
	if o.PoolIds == nil {
		o.PoolIds = []string{}
	}

	groupType := resourceGroupType(o.ResourceGroup.Name)
	if groupType == nil {
		return nil, fmt.Errorf("unable to determine resource type for resource group %q", o.ResourceGroup.Name)
	}

	type Alias ResourceGroupAllocation // Does not implement json.Marshaler, so no infinite recursion.

	return json.Marshal(struct {
		Alias
		Type *enum.ResourceType `json:"type"`
	}{
		Alias: Alias(o),
		Type:  groupType,
	})
}

func (o *ResourceGroupAllocation) IsEmpty() bool {
	return len(o.PoolIds) == 0
}

func (o *TwoStageL3ClosClient) getResourceAllocations(ctx context.Context) ([]ResourceGroupAllocation, error) {
	response := &struct {
		Items []ResourceGroupAllocation `json:"items"`
	}{}
	return response.Items, o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(urls.DatacenterResourceGroups, o.blueprintId),
		apiResponse: response,
	})
}

func (o *TwoStageL3ClosClient) getResourceAllocation(ctx context.Context, rg ResourceGroup) (ResourceGroupAllocation, error) {
	var response ResourceGroupAllocation

	rgType := resourceGroupType(rg.Name)
	if rgType == nil {
		return response, fmt.Errorf("unable to determine resource type for resource group %q", rg.Name)
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(urls.DatacenterResourceGroupsByTypeName, o.blueprintId, *rgType, rg),
		apiResponse: &response,
	})
	if err != nil {
		return response, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *TwoStageL3ClosClient) setResourceAllocation(ctx context.Context, rga ResourceGroupAllocation) error {
	rgType := resourceGroupType(rga.ResourceGroup.Name)
	if rgType == nil {
		return fmt.Errorf("unable to determine resource type for resource group %q", rga.ResourceGroup.Name)
	}

	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(urls.DatacenterResourceGroupsByTypeName, o.blueprintId, *rgType, rga.ResourceGroup.String()),
		apiInput: rga,
	})
}

func resourceGroupType(rg enum.ResourceGroup) *enum.ResourceType {
	switch { // Order matters: IPv6 must be checked before IPv4
	case strings.HasSuffix(rg.Value, "_asns"):
		return &enum.ResourceTypeASN
	case strings.HasPrefix(rg.Value, "ipv6_") || strings.HasSuffix(rg.Value, "_ipv6"):
		return &enum.ResourceTypeIPv6
	case strings.HasSuffix(rg.Value, "_ips") || strings.HasSuffix(rg.Value, "_subnets"):
		return &enum.ResourceTypeIPv4
	case strings.HasSuffix(rg.Value, "_vnis") || strings.HasSuffix(rg.Value, "_vn_ids"):
		return &enum.ResourceTypeVNI
	}

	return nil
}
