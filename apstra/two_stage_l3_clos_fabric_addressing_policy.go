package apstra

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	if in.FabricL3Mtu != nil && fabricL3MtuForbidden().Includes(o.client.apiVersion) {
		return ClientErr{
			errType: ErrCompatibility,
			err:     errors.New(fabricL3MtuForbiddenError),
		}
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
