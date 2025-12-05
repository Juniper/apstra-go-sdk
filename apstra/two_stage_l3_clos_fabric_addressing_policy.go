// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/compatibility"
)

const (
	apiUrlBlueprintFabricAddressingPolicy = apiUrlBlueprintById + apiUrlPathDelim + "fabric-addressing-policy"
)

type TwoStageL3ClosFabricAddressingPolicy struct {
	Ipv6Enabled *bool   `json:"ipv6_enabled,omitempty"`
	EsiMacMsb   *uint8  `json:"esi_mac_msb,omitempty"`
	FabricL3Mtu *uint16 `json:"fabric_l3_mtu,omitempty"`
}

func (o *TwoStageL3ClosClient) GetFabricAddressingPolicy(ctx context.Context) (*TwoStageL3ClosFabricAddressingPolicy, error) {
	var result TwoStageL3ClosFabricAddressingPolicy
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintFabricAddressingPolicy, o.blueprintId),
		apiResponse: &result,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &result, nil
}

func (o *TwoStageL3ClosClient) SetFabricAddressingPolicy(ctx context.Context, in *TwoStageL3ClosFabricAddressingPolicy) error {
	if !compatibility.EqApstra420.Check(o.client.apiVersion) {
		return fmt.Errorf("SetFabricAddressingPolicy only for use with Apstra %s", compatibility.EqApstra420)
	}

	if in.Ipv6Enabled == nil &&
		in.EsiMacMsb == nil &&
		in.FabricL3Mtu == nil {
		return nil // nothing to do if all relevant input fields are nil
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintFabricAddressingPolicy, o.blueprintId),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
