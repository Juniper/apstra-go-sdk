// Copyright (c) Juniper Networks, Inc., 2022-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/str"
	"github.com/Juniper/apstra-go-sdk/internal/urls"
	"github.com/hashicorp/go-version"
)

type (
	SystemRole int
	systemRole string
)

const (
	SystemRoleNone = SystemRole(iota)
	SystemRoleAccess
	SystemRoleGeneric
	SystemRoleLeaf
	SystemRoleSpine
	SystemRoleSuperSpine
	SystemRoleRedundancyGroup
	SystemRoleUnknown = "unknown System Role '%s'"

	systemRoleNone            = systemRole("")
	systemRoleAccess          = systemRole("access")
	systemRoleGeneric         = systemRole("generic")
	systemRoleLeaf            = systemRole("leaf")
	systemRoleSpine           = systemRole("spine")
	systemRoleSuperSpine      = systemRole("superspine")
	systemRoleRedundancyGroup = systemRole("redundancy_group")
	systemRoleUnknown         = "unknown System Role '%d'"
)

func (o SystemRole) String() string {
	return string(o.raw())
}

func (o SystemRole) int() int {
	return int(o)
}

func (o SystemRole) raw() systemRole {
	switch o {
	case SystemRoleNone:
		return systemRoleNone
	case SystemRoleAccess:
		return systemRoleAccess
	case SystemRoleGeneric:
		return systemRoleGeneric
	case SystemRoleLeaf:
		return systemRoleLeaf
	case SystemRoleSpine:
		return systemRoleSpine
	case SystemRoleSuperSpine:
		return systemRoleSuperSpine
	case SystemRoleRedundancyGroup:
		return systemRoleRedundancyGroup
	default:
		return systemRole(fmt.Sprintf(systemRoleUnknown, o))
	}
}

func (o *SystemRole) FromString(in string) error {
	i, err := systemRole(in).parse()
	if err != nil {
		return err
	}
	*o = SystemRole(i)
	return nil
}

func (o systemRole) string() string {
	return string(o)
}

func (o systemRole) parse() (int, error) {
	switch o {
	case systemRoleNone:
		return int(SystemRoleNone), nil
	case systemRoleAccess:
		return int(SystemRoleAccess), nil
	case systemRoleGeneric:
		return int(SystemRoleGeneric), nil
	case systemRoleLeaf:
		return int(SystemRoleLeaf), nil
	case systemRoleSpine:
		return int(SystemRoleSpine), nil
	case systemRoleSuperSpine:
		return int(SystemRoleSuperSpine), nil
	case systemRoleRedundancyGroup:
		return int(SystemRoleRedundancyGroup), nil
	default:
		return 0, fmt.Errorf(SystemRoleUnknown, o)
	}
}

type Endpoint struct {
	InterfaceId ObjectId `json:"interface_id"`
	TagType     string   `json:"tag_type"`
	Label       string   `json:"label"`
}

func (o *TwoStageL3ClosClient) ListVirtualNetworks(ctx context.Context) ([]string, error) {
	var response struct {
		VirtualNetworks map[string]datacenter.VirtualNetwork `json:"virtual_networks"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodGet,
		urlStr:         fmt.Sprintf(urls.DatacenterVirtualNetworks, o.blueprintId),
		apiResponse:    &response,
		unsynchronized: true,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return slices.Collect(maps.Keys(response.VirtualNetworks)), nil
}

func (o *TwoStageL3ClosClient) GetVirtualNetwork(ctx context.Context, id string) (datacenter.VirtualNetwork, error) {
	var response datacenter.VirtualNetwork

	apstraUrl, err := o.urlWithParam(fmt.Sprintf(urls.DatacenterVirtualNetworkByID, o.blueprintId, id))
	if err != nil {
		return response, err
	}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodGet,
		url:            apstraUrl,
		apiResponse:    &response,
		unsynchronized: true,
	})
	if err != nil {
		return response, convertTtaeToAceWherePossible(err)
	}

	vns := []datacenter.VirtualNetwork{response}
	if compatibility.VirtualNetworkAddressesInActiveGraphOnly.Check(version.Must(version.NewVersion(o.client.ApiVersion()))) {
		err = o.getVirtualNetworkAddressingFromActiveGraph(ctx, vns)
	}
	return vns[0], nil
}

func (o *TwoStageL3ClosClient) GetVirtualNetworkByLabel(ctx context.Context, label string) (datacenter.VirtualNetwork, error) {
	vns, err := o.GetVirtualNetworks(ctx)
	if err != nil {
		return datacenter.VirtualNetwork{}, err
	}

	var result *datacenter.VirtualNetwork
	for _, vn := range vns {
		if vn.Label == label {
			if result == nil {
				result = &vn
			} else {
				return datacenter.VirtualNetwork{}, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple matches for virtual network with label %q", label),
				}
			}
		}
	}

	if result == nil {
		return datacenter.VirtualNetwork{}, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("virtual network with label %q not found", label),
		}
	}

	return *result, nil
}

func (o *TwoStageL3ClosClient) GetVirtualNetworks(ctx context.Context) ([]datacenter.VirtualNetwork, error) {
	var response struct {
		VirtualNetworks map[string]datacenter.VirtualNetwork `json:"virtual_networks"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(urls.DatacenterVirtualNetworks, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	vns := slices.Collect(maps.Values(response.VirtualNetworks))
	if compatibility.VirtualNetworkAddressesInActiveGraphOnly.Check(version.Must(version.NewVersion(o.client.ApiVersion()))) {
		err = o.getVirtualNetworkAddressingFromActiveGraph(ctx, vns)
	}
	return vns, err
}

func (o *TwoStageL3ClosClient) CreateVirtualNetwork(ctx context.Context, vn datacenter.VirtualNetwork) (string, error) {
	if vn.Tags != nil {
		return "", ClientErr{
			errType: ErrNotSupported,
			err:     fmt.Errorf("tags must be nil when creating virtual network with Apstra %s", o.client.ApiVersion()),
		}
	}

	if vn.Type == enum.VnTypeVlan && len(vn.Bindings) > 1 {
		return "", ClientErr{
			errType: ErrInvalidId,
			err:     fmt.Errorf("virtual network of type %q cannot have more than one binding", enum.VnTypeVlan),
		}
	}

	var response struct {
		Id string `json:"id"`
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(urls.DatacenterVirtualNetworks, o.blueprintId),
		apiInput:    vn,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *TwoStageL3ClosClient) UpdateVirtualNetwork(ctx context.Context, vn datacenter.VirtualNetwork) error {
	if vn.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	if vn.Tags != nil {
		return ClientErr{
			errType: ErrNotSupported,
			err:     fmt.Errorf("tags must be nil when updating virtual network with Apstra %s", o.client.ApiVersion()),
		}
	}

	if vn.Bindings == nil {
		vn.Bindings = make([]datacenter.VNBinding, 0) // convert nil slice to empty slice because API doesn't implement PUT semantics correctly
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(urls.DatacenterVirtualNetworkByID, o.blueprintId, *vn.ID()),
		apiInput: vn,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) DeleteVirtualNetwork(ctx context.Context, id string) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(urls.DatacenterVirtualNetworkByID, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

// getVirtualNetworkAddressingFromActiveGraph is a workaround for the fact that the API
// doesn't return subnet/gateway information when there are no bindings, even though it
// could be populated. It queries the active graph for any VNs that are missing this
// information and populates it from there. This is only necessary for VNs with no bindings
// because when there are bindings the API does return this information.
func (o *TwoStageL3ClosClient) getVirtualNetworkAddressingFromActiveGraph(ctx context.Context, vns []datacenter.VirtualNetwork) error {
	var ids []string // suspect IDs we want to retrieve
	for _, vn := range vns {
		if vn.ID() == nil {
			return fmt.Errorf("VNs must have non-nil id in %s", str.FuncName())
		}

		// The condition we're fixing only happens with zero bindings AND
		// a nil subnet or gateway value where one could reasonably be populated.
		if len(vn.Bindings) == 0 &&
			(vn.IPv4Enabled && vn.IPv4Subnet == nil ||
				vn.IPv6Enabled && vn.IPv6Subnet == nil ||
				vn.VirtualGatewayIPv4Enabled && vn.VirtualGatewayIPv4 == nil ||
				vn.VirtualGatewayIPv6Enabled && vn.VirtualGatewayIPv6 == nil) {
			ids = append(ids, *vn.ID())
		}
	}

	if len(ids) == 0 {
		return nil // nothing to do because no potentially omitted subnet/gateway values
	}

	// Query the active graph for interesting VNs.
	vnMap, err := o.vnsFromActiveGraph(ctx, ids)
	if err != nil {
		return err
	}

	// Populate missing values of each VN with values from the active graph.
	for i, vn := range vns {
		if vnFromMap, ok := vnMap[*vn.ID()]; ok {
			if vns[i].IPv4Subnet == nil {
				vns[i].IPv4Subnet = pointer.ToCopyOfValue(vnFromMap.IPv4Subnet)
			}
			if vns[i].IPv6Subnet == nil {
				vns[i].IPv6Subnet = pointer.ToCopyOfValue(vnFromMap.IPv6Subnet)
			}
			if vns[i].VirtualGatewayIPv4 == nil {
				copy(vns[i].VirtualGatewayIPv4, vnFromMap.VirtualGatewayIPv4)
			}
			if vns[i].VirtualGatewayIPv6 == nil {
				copy(vns[i].VirtualGatewayIPv6, vnFromMap.VirtualGatewayIPv6)
			}
		}
	}

	return nil
}

// vnsFromActiveGraph queries the active graph for the given VN IDs and returns a map
// of VN ID to VirtualNetwork object. This is used as a workaround for the fact that
// the API doesn't return subnet/gateway information for VNs with no bindings.
func (o *TwoStageL3ClosClient) vnsFromActiveGraph(ctx context.Context, ids []string) (map[string]datacenter.VirtualNetwork, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	q := new(PathQuery).
		SetClient(o.Client()).
		SetBlueprintId(o.Id()).
		SetBlueprintType(BlueprintTypeConfig).
		Node([]QEEAttribute{
			NodeTypeVirtualNetwork.QEEAttribute(),
			{Key: "id", Value: QEStringValIsIn(ids)},
			{Key: "name", Value: QEStringVal("n_virtual_network")},
		})

	var t struct {
		Items []struct {
			VN datacenter.VirtualNetwork `json:"n_virtual_network"`
		}
	}

	err := q.Do(ctx, &t)
	if err != nil {
		return nil, err
	}

	result := make(map[string]datacenter.VirtualNetwork, len(t.Items))
	for i, item := range t.Items {
		id := item.VN.ID()
		if id == nil {
			return nil, fmt.Errorf("invalid virtual network ID at item %d", i)
		}
		result[*id] = item.VN
	}

	return result, nil
}
