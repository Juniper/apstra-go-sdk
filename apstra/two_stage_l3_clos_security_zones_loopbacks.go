// Copyright (c) Juniper Networks, Inc., 2024-2024.
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
type SecurityZoneLoopback struct {
	IPv4Addr *netip.Prefix
	IPv6Addr *netip.Prefix
}

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

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintSecurityZoneLoopbacksById, o.blueprintId, szId),
		apiInput: apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

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
				Id       ObjectId `json:"id"`
				IPv4Addr *string  `json:"ipv4_addr"`
				IPv6Addr *string  `json:"ipv6_addr"`
			} `json:"n_interface""`
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
		var ipv4Addr *netip.Prefix
		if item.Interface.IPv4Addr != nil {
			ip, err := netip.ParsePrefix(*item.Interface.IPv4Addr)
			if err != nil {
				return nil, fmt.Errorf("failed parsing node %q ipv4_addr value %q - %w", item.Interface.Id, *item.Interface.IPv4Addr, err)
			}
			ipv4Addr = &ip
		}

		var ipv6Addr *netip.Prefix
		if item.Interface.IPv6Addr != nil {
			ip, err := netip.ParsePrefix(*item.Interface.IPv6Addr)
			if err != nil {
				return nil, fmt.Errorf("failed parsing node %q ipv6_addr value %q - %w", item.Interface.Id, *item.Interface.IPv6Addr, err)
			}
			ipv6Addr = &ip
		}

		result[item.Interface.Id] = SecurityZoneLoopback{
			IPv4Addr: ipv4Addr,
			IPv6Addr: ipv6Addr,
		}
	}

	return result, nil
}

func (o *TwoStageL3ClosClient) GetSecurityZoneLoopbackByInterfaceId(ctx context.Context, id ObjectId) (*SecurityZoneLoopback, error) {
	var target struct {
		IPv4Addr *string `json:"ipv4_addr"`
		IPv6Addr *string `json:"ipv6_addr"`
	}

	err := o.client.GetNode(ctx, o.blueprintId, id, &target)
	if err != nil {
		return nil, fmt.Errorf("failed fetching ndoe %q from blueprint %q - %w", id, o.blueprintId, err)
	}

	var ipv4Addr *netip.Prefix
	if target.IPv4Addr != nil {
		ip, err := netip.ParsePrefix(*target.IPv4Addr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing node %q ipv4_addr value %q - %w", id, *target.IPv4Addr, err)
		}
		ipv4Addr = &ip
	}

	var ipv6Addr *netip.Prefix
	if target.IPv6Addr != nil {
		ip, err := netip.ParsePrefix(*target.IPv6Addr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing node %q ipv6_addr value %q - %w", id, *target.IPv6Addr, err)
		}
		ipv6Addr = &ip
	}

	return &SecurityZoneLoopback{
		IPv4Addr: ipv4Addr,
		IPv6Addr: ipv6Addr,
	}, nil
}
