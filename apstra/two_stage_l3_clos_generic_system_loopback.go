// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const (
	apiUrlBlueprintSystems          = apiUrlBlueprintById + apiUrlPathDelim + "systems"
	apiUrlBlueprintSystemById       = apiUrlBlueprintSystems + apiUrlPathDelim + "%s"
	apiUrlBlueprintSystemByIdPrefix = apiUrlBlueprintSystemById + apiUrlPathDelim
	apiUrlGenericSystemLoopback     = apiUrlBlueprintSystemById + apiUrlPathDelim + "loopback"
	apiUrlGenericSystemLoopbackById = apiUrlGenericSystemLoopback + apiUrlPathDelim + "%d"
)

type GenericSystemLoopback struct {
	Ipv4Addr       *net.IPNet
	Ipv6Addr       *net.IPNet
	LoopbackNodeId ObjectId
	LoopbackId     uint32
	SecurityZoneId ObjectId
}

func (o GenericSystemLoopback) raw() *rawGenericSystemLoopback {
	var ipv4Addr, ipv6Addr *string
	if o.Ipv4Addr != nil {
		s := o.Ipv4Addr.String()
		ipv4Addr = &s
	}
	if o.Ipv6Addr != nil {
		s := o.Ipv6Addr.String()
		ipv6Addr = &s
	}

	return &rawGenericSystemLoopback{
		Ipv4Addr:       ipv4Addr,
		Ipv6Addr:       ipv6Addr,
		SecurityZoneId: o.SecurityZoneId,
	}
}

type rawGenericSystemLoopback struct {
	Ipv4Addr       *string  `json:"ipv4_addr"` // do not omitempty - we must send `null` to remove the address
	Ipv6Addr       *string  `json:"ipv6_addr"` // do not omitempty - we must send `null` to remove the address
	LoopbackNodeId ObjectId `json:"loopback_node_id,omitempty"`
	SecurityZoneId ObjectId `json:"sz_id,omitempty"`
}

func (o rawGenericSystemLoopback) polish() (*GenericSystemLoopback, error) {
	var err error

	var ipv4Addr *net.IPNet
	if o.Ipv4Addr != nil {
		_, ipv4Addr, err = net.ParseCIDR(*o.Ipv4Addr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing loopback 'ipv4_addr' from API: %q - %w", *o.Ipv4Addr, err)
		}
	}

	var ipv6Addr *net.IPNet
	if o.Ipv6Addr != nil {
		_, ipv6Addr, err = net.ParseCIDR(*o.Ipv6Addr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing loopback 'ipv6_addr' from API: %q - %w", *o.Ipv6Addr, err)
		}
	}

	return &GenericSystemLoopback{
		Ipv4Addr:       ipv4Addr,
		Ipv6Addr:       ipv6Addr,
		LoopbackNodeId: o.LoopbackNodeId,
		SecurityZoneId: o.SecurityZoneId,
	}, nil
}

func (o *TwoStageL3ClosClient) GetGenericSystemLoopback(ctx context.Context, nodeId ObjectId, loopbackId int) (*GenericSystemLoopback, error) {
	raw, err := o.getGenericSystemLoopback(ctx, nodeId, loopbackId)
	if err != nil {
		return nil, err
	}

	return raw.polish()
}

func (o *TwoStageL3ClosClient) getGenericSystemLoopback(ctx context.Context, nodeId ObjectId, loopbackId int) (*rawGenericSystemLoopback, error) {
	var response rawGenericSystemLoopback
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlGenericSystemLoopbackById, o.blueprintId, nodeId, loopbackId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *TwoStageL3ClosClient) GetGenericSystemLoopbacks(ctx context.Context, nodeId ObjectId) (map[int]GenericSystemLoopback, error) {
	raw, err := o.getGenericSystemLoopbacks(ctx, nodeId)
	if err != nil {
		return nil, err
	}

	result := make(map[int]GenericSystemLoopback, len(raw))
	for k, v := range raw {
		p, err := v.polish()
		if err != nil {
			return nil, err
		}
		result[k] = *p
	}

	return result, nil
}

func (o *TwoStageL3ClosClient) getGenericSystemLoopbacks(ctx context.Context, nodeId ObjectId) (map[int]rawGenericSystemLoopback, error) {
	// prep a graph query which finds all loopback interfaces attached to the given node
	query := new(PathQuery).
		SetBlueprintId(o.blueprintId).
		SetClient(o.client).
		Node([]QEEAttribute{
			NodeTypeSystem.QEEAttribute(),
			{Key: "id", Value: QEStringVal(nodeId.String())},
		}).
		Out([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
		Node([]QEEAttribute{
			NodeTypeInterface.QEEAttribute(),
			{Key: "if_type", Value: QEStringVal("loopback")},
			{Key: "name", Value: QEStringVal("n_interface")},
		})

	// we only need one attribute: the loopback id (an integer)
	var result struct {
		Items []struct {
			LoopbackId int `json:"loopback_id"`
		} `json:"items"`
	}

	s := query.String()
	_ = s

	// run the query
	err := query.Do(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("failed executing graph query %q - %w", query.String(), err)
	}

	// prepare the result map
	resultMap := make(map[int]rawGenericSystemLoopback, len(result.Items))
	for _, item := range result.Items {
		loopback, err := o.getGenericSystemLoopback(ctx, nodeId, item.LoopbackId)
		if err != nil {
			return nil, fmt.Errorf("failed fetching blueprint %q node %q loopback %d info - %w",
				o.blueprintId, nodeId, item.LoopbackId, err)
		}

		resultMap[item.LoopbackId] = *loopback
	}

	return resultMap, nil
}

func (o *TwoStageL3ClosClient) SetGenericSystemLoopback(ctx context.Context, nodeId ObjectId, loopbackId int, in *GenericSystemLoopback) error {
	return o.setGenericSystemLoopback(ctx, nodeId, loopbackId, in.raw())
}

func (o *TwoStageL3ClosClient) setGenericSystemLoopback(ctx context.Context, nodeId ObjectId, loopbackId int, in *rawGenericSystemLoopback) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlGenericSystemLoopbackById, o.blueprintId, nodeId, loopbackId),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
