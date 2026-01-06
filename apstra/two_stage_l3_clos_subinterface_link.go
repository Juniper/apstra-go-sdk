// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlWebExpSubinterfaces = apiUrlBlueprintByIdPrefix + "experience/web/subinterfaces"
)

var _ json.Unmarshaler = (*TwoStageL3ClosSubinterfaceLink)(nil)

type TwoStageL3ClosSubinterfaceLinkEndpoint struct {
	System struct {
		Id    ObjectId
		Label string
		Role  SystemRole
	}
	InterfaceId    ObjectId
	SubinterfaceId ObjectId
	Subinterface   TwoStageL3ClosSubinterface
}

type TwoStageL3ClosSubinterfaceLink struct {
	Id        ObjectId // logical link ID
	VlanId    *VLAN
	Endpoints []TwoStageL3ClosSubinterfaceLinkEndpoint
	SzId      ObjectId
	SzLabel   string
}

func (o *TwoStageL3ClosSubinterfaceLink) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		LinkId    ObjectId `json:"link_id"`
		VlanId    *VLAN    `json:"vlan_id"`
		Endpoints []struct {
			System struct {
				Id    ObjectId   `json:"id"`
				Label string     `json:"label"`
				Role  systemRole `json:"role"`
			} `json:"system"`
			Interface struct {
				Id ObjectId `json:"id"`
			} `json:"interface"`
			Subinterface json.RawMessage `json:"subinterface"`
		} `json:"endpoints"`
		SzId    ObjectId `json:"sz_id"`
		SzLabel string   `json:"sz_label"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.LinkId
	o.VlanId = raw.VlanId
	o.SzId = raw.SzId
	o.SzLabel = raw.SzLabel
	o.Endpoints = make([]TwoStageL3ClosSubinterfaceLinkEndpoint, len(raw.Endpoints))
	for i, rep := range raw.Endpoints {
		var rawSubinterface struct {
			Id ObjectId `json:"id"`
		}

		err = json.Unmarshal(rep.Subinterface, &rawSubinterface)
		if err != nil {
			return err
		}
		subinterfaceId := rawSubinterface.Id

		var subinterface TwoStageL3ClosSubinterface
		err = json.Unmarshal(rep.Subinterface, &subinterface)
		if err != nil {
			return err
		}

		// o.Endpoints[i].Subinterface.Ipv4AddrType = new(I)
		sysRole, err := rep.System.Role.parse()
		if err != nil {
			return fmt.Errorf("failed to parse system role %q while unmarshaling TwoStageL3ClosSubinterfaceLink - %w", rep.System.Role, err)
		}

		o.Endpoints[i] = TwoStageL3ClosSubinterfaceLinkEndpoint{
			System: struct {
				Id    ObjectId
				Label string
				Role  SystemRole
			}{
				Id:    rep.System.Id,
				Label: rep.System.Label,
				Role:  SystemRole(sysRole),
			},
			InterfaceId:    rep.Interface.Id,
			Subinterface:   subinterface,
			SubinterfaceId: subinterfaceId,
		}
	}

	return nil
}

func (o *TwoStageL3ClosClient) GetSubinterfaceLink(ctx context.Context, id ObjectId) (*TwoStageL3ClosSubinterfaceLink, error) {
	all, err := o.GetAllSubinterfaceLinks(ctx)
	if err != nil {
		return nil, err
	}

	for _, link := range all {
		if link.Id == id {
			return &link, nil
		}
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("link %q not found", id),
	}
}

func (o *TwoStageL3ClosClient) GetAllSubinterfaceLinks(ctx context.Context) ([]TwoStageL3ClosSubinterfaceLink, error) {
	var response struct {
		Links []TwoStageL3ClosSubinterfaceLink `json:"subinterfaces"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlWebExpSubinterfaces, o.Id()),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Links, nil
}
