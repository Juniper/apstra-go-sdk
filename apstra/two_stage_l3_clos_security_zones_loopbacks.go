// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
)

const apiUrlBlueprintSecurityZoneLoopbacksById = apiUrlBlueprintSecurityZoneById + apiUrlPathDelim + "loopbacks"

var _ json.Marshaler = (*SecurityZoneLoopback)(nil)

// SecurityZoneLoopback is intended to be used with the SetSecurityZoneLoopbacks() method
// and the apiUrlBlueprintSecurityZoneLoopbacksById API endpoint. It is possible to produce
// three different outcomes in the rendered JSON for both IPv4Addr and IPv6Addr elements:
//
//  1. When the element and its IP and Mask elements are non-nil, a string will be produced
//     when the struct is marshaled as JSON.
//  2. When the element is non-nil but contains a nil IP or Mask, a `null` will be produced
//     when the struct is marshaled as JSON.
//  3. When the element is nil, no output will be produced for that element when the struct
//     is marshaled as JSON.
//
// Example:
//
//	aVal := netip.MustParsePrefix("192.0.2.0/32")
//	a := apstra.SecurityZoneLoopback{IPv4Addr: &aVal}
//	b := apstra.SecurityZoneLoopback{IPv4Addr: &netip.Prefix{}}
//	c := apstra.SecurityZoneLoopback{IPv4Addr: nil}
//
//	aJson, _ := json.Marshal(a)
//	bJson, _ := json.Marshal(b)
//	cJson, _ := json.Marshal(c)
//
//	fmt.Print(string(aJson) + "\n" + string(bJson) + "\n" + string(cJson) + "\n")
//
// Output:
//
//	{"ipv4_addr":"192.0.2.0/32"}
//	{"ipv4_addr":null}
//	{}
//
// Note that the API is only concerned with the IPv4Addr and IPv6Addr elements.
// We include those other elements in various GET methods for convenience.
type SecurityZoneLoopback struct {
	Id             ObjectId
	IPv4Addr       *netip.Prefix
	IPv6Addr       *netip.Prefix
	LoopbackId     int
	SecurityZoneId ObjectId
}

func (o *SecurityZoneLoopback) loadIpsFromStringPointers(ipv4Addr, ipv6Addr *string) error {
	if ipv4Addr != nil {
		ip, err := netip.ParsePrefix(*ipv4Addr)
		if err != nil {
			return fmt.Errorf("failed parsing ipv4_addr value %q - %w", *ipv4Addr, err)
		}
		o.IPv4Addr = &ip
	}

	if ipv6Addr != nil {
		ip, err := netip.ParsePrefix(*ipv6Addr)
		if err != nil {
			return fmt.Errorf("failed parsing ipv6_addr value %q - %w", *ipv6Addr, err)
		}
		o.IPv6Addr = &ip
	}

	return nil
}

// MarshalJSON only concerns itself with the IPv4Addr and IPv6Addr elements.
// All other elements are ignored because the API doesn't want them. We include
// those other elements in various GET methods for convenience.
func (o SecurityZoneLoopback) MarshalJSON() ([]byte, error) {
	ipInfo := make(map[string]*string)

	if o.IPv4Addr != nil {
		if o.IPv4Addr.IsValid() {
			ipInfo["ipv4_addr"] = toPtr(o.IPv4Addr.String())
		} else {
			ipInfo["ipv4_addr"] = nil
		}
	}

	if o.IPv6Addr != nil {
		if o.IPv6Addr.IsValid() {
			ipInfo["ipv6_addr"] = toPtr(o.IPv6Addr.String())
		} else {
			ipInfo["ipv6_addr"] = nil
		}
	}

	return json.Marshal(ipInfo)
}

// SetSecurityZoneLoopbacks takes a map of SecurityZoneLoopback keyed by the loopback interface graph node ID.
// Only IP information is required in the `SecurityZoneLoopback` objects passed to this function. The other
// elements will not be rendered to JSON and would be ignored by the API. See the SecurityZoneLoopback for an
// explanation of how to set and clear addresses (Go nil vs. JSON null, etc...)
func (o *TwoStageL3ClosClient) SetSecurityZoneLoopbacks(ctx context.Context, szId ObjectId, loopbacks map[ObjectId]SecurityZoneLoopback) error {
	if !compatibility.SecurityZoneLoopbackApiSupported.Check(o.client.apiVersion) {
		return fmt.Errorf("SetSecurityZoneLoopbacks requires Apstra version %s, have version %s",
			compatibility.SecurityZoneLoopbackApiSupported, o.client.apiVersion,
		)
	}

	var apiInput struct {
		Loopbacks map[ObjectId]json.RawMessage `json:"loopbacks"`
	}
	apiInput.Loopbacks = make(map[ObjectId]json.RawMessage)

	for k, v := range loopbacks {
		rawJson, err := json.Marshal(v)
		if err != nil {
			return err
		}
		apiInput.Loopbacks[k] = rawJson
	}

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintSecurityZoneLoopbacksById, o.blueprintId, szId),
		apiInput: apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

// GetSecurityZoneLoopbacks returns a map keyed by loopback interface node ID.
// This is the format used by the API (PATCH), and by our SetSecurityZoneLoopbacks
// function.
func (o *TwoStageL3ClosClient) GetSecurityZoneLoopbacks(ctx context.Context, szId ObjectId) (map[ObjectId]SecurityZoneLoopback, error) {
	query := new(PathQuery).
		SetBlueprintId(o.blueprintId).
		SetClient(o.client).
		Node([]QEEAttribute{
			NodeTypeSecurityZone.QEEAttribute(),
			{Key: "id", Value: QEStringVal(szId)},
		}).
		Out([]QEEAttribute{RelationshipTypeInstantiatedBy.QEEAttribute()}).
		Node([]QEEAttribute{NodeTypeSecurityZoneInstance.QEEAttribute()}).
		Out([]QEEAttribute{RelationshipTypeMemberInterfaces.QEEAttribute()}).
		Node([]QEEAttribute{
			NodeTypeInterface.QEEAttribute(),
			{Key: "if_type", Value: QEStringVal("loopback")},
			{Key: "name", Value: QEStringVal("n_interface")},
		})

	var queryResponse struct {
		Items []struct {
			Interface struct {
				Id         ObjectId `json:"id"`
				LoopbackId int      `json:"loopback_id"`
				IPv4Addr   *string  `json:"ipv4_addr"`
				IPv6Addr   *string  `json:"ipv6_addr"`
			} `json:"n_interface"`
		} `json:"items"`
	}

	err := query.Do(ctx, &queryResponse)
	if err != nil {
		return nil, err
	}
	if len(queryResponse.Items) == 0 {
		return nil, nil
	}

	result := make(map[ObjectId]SecurityZoneLoopback, len(queryResponse.Items))
	for _, item := range queryResponse.Items {
		szl := SecurityZoneLoopback{
			Id:             item.Interface.Id,
			LoopbackId:     item.Interface.LoopbackId,
			SecurityZoneId: szId,
		}

		err = szl.loadIpsFromStringPointers(item.Interface.IPv4Addr, item.Interface.IPv6Addr)
		if err != nil {
			return nil, fmt.Errorf("failed loading node %q IP info - %w", item.Interface.Id, err)
		}

		result[item.Interface.Id] = szl
	}

	return result, nil
}

// GetSecurityZoneLoopbackByInterfaceId returns a single *SecurityZoneLoopback
// representing the specified loopback graph node ID.
func (o *TwoStageL3ClosClient) GetSecurityZoneLoopbackByInterfaceId(ctx context.Context, ifId ObjectId) (*SecurityZoneLoopback, error) {
	query := new(PathQuery).
		SetBlueprintId(o.blueprintId).
		SetClient(o.client).
		Node([]QEEAttribute{
			NodeTypeInterface.QEEAttribute(),
			{Key: "id", Value: QEStringVal(ifId)},
			{Key: "if_type", Value: QEStringVal("loopback")},
			{Key: "name", Value: QEStringVal("n_interface")},
		}).
		In([]QEEAttribute{RelationshipTypeMemberInterfaces.QEEAttribute()}).
		Node([]QEEAttribute{NodeTypeSecurityZoneInstance.QEEAttribute()}).
		In([]QEEAttribute{RelationshipTypeInstantiatedBy.QEEAttribute()}).
		Node([]QEEAttribute{
			NodeTypeSecurityZone.QEEAttribute(),
			{Key: "name", Value: QEStringVal("n_security_zone")},
		})

	var queryResult struct {
		Items []struct {
			Interface struct {
				Id         ObjectId `json:"id"`
				Ipv4Addr   *string  `json:"ipv4_addr"`
				Ipv6Addr   *string  `json:"ipv6_addr"`
				LoopbackId int      `json:"loopback_id"`
			} `json:"n_interface"`
			SecurityZone struct {
				Id ObjectId `json:"id"`
			} `json:"n_security_zone"`
		} `json:"items"`
	}

	err := query.Do(ctx, &queryResult)
	if err != nil {
		return nil, fmt.Errorf("failed while querying for loopback interface %q - %w", ifId, err)
	}

	switch len(queryResult.Items) {
	case 0:
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("loopback interface %q not found with query %q", ifId, query),
		}
	case 1:
	default:
		return nil, ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("found %d matches while looking for loopback interface %q with query %q", len(queryResult.Items), ifId, query),
		}
	}

	szl := SecurityZoneLoopback{
		Id:             queryResult.Items[0].Interface.Id,
		LoopbackId:     queryResult.Items[0].Interface.LoopbackId,
		SecurityZoneId: queryResult.Items[0].SecurityZone.Id,
	}

	err = szl.loadIpsFromStringPointers(queryResult.Items[0].Interface.Ipv4Addr, queryResult.Items[0].Interface.Ipv6Addr)
	if err != nil {
		return nil, fmt.Errorf("failed loading node %q IP info - %w", queryResult.Items[0].Interface.Id, err)
	}

	return &szl, nil
}

// GetSecurityZoneLoopbacksBySystemId returns []SecurityZoneLoopback representing
// all of the loopback interfaces belonging to the specified system node.
func (o *TwoStageL3ClosClient) GetSecurityZoneLoopbacksBySystemId(ctx context.Context, sysId ObjectId) ([]SecurityZoneLoopback, error) {
	query := new(PathQuery).
		SetBlueprintId(o.blueprintId).
		SetClient(o.client).
		Node([]QEEAttribute{
			NodeTypeSystem.QEEAttribute(),
			{Key: "id", Value: QEStringVal(sysId)},
		}).
		Out([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
		Node([]QEEAttribute{
			NodeTypeInterface.QEEAttribute(),
			{Key: "if_type", Value: QEStringVal("loopback")},
			{Key: "name", Value: QEStringVal("n_interface")},
		}).
		In([]QEEAttribute{RelationshipTypeMemberInterfaces.QEEAttribute()}).
		Node([]QEEAttribute{NodeTypeSecurityZoneInstance.QEEAttribute()}).
		In([]QEEAttribute{RelationshipTypeInstantiatedBy.QEEAttribute()}).
		Node([]QEEAttribute{
			NodeTypeSecurityZone.QEEAttribute(),
			{Key: "name", Value: QEStringVal("n_security_zone")},
		})

	var queryResult struct {
		Items []struct {
			Interface struct {
				Id         ObjectId `json:"id"`
				Ipv4Addr   *string  `json:"ipv4_addr"`
				Ipv6Addr   *string  `json:"ipv6_addr"`
				LoopbackId int      `json:"loopback_id"`
			} `json:"n_interface"`
			SecurityZone struct {
				Id ObjectId `json:"id"`
			} `json:"n_security_zone"`
		} `json:"items"`
	}

	err := query.Do(ctx, &queryResult)
	if err != nil {
		return nil, fmt.Errorf("failed while querying for system %q loopback", sysId)
	}

	result := make([]SecurityZoneLoopback, len(queryResult.Items))
	for i, item := range queryResult.Items {
		szl := SecurityZoneLoopback{
			Id:             item.Interface.Id,
			LoopbackId:     item.Interface.LoopbackId,
			SecurityZoneId: item.SecurityZone.Id,
		}

		err = szl.loadIpsFromStringPointers(item.Interface.Ipv4Addr, item.Interface.Ipv6Addr)
		if err != nil {
			return nil, fmt.Errorf("failed loading node %q IP info - %w", item.Interface.Id, err)
		}

		result[i] = szl
	}

	return result, nil
}
