package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintFabricAddressingPolicy = apiUrlBlueprintById + apiUrlPathDelim + "fabric-addressing-policy"
)

type TwoStageL3ClosFabricAddressingPolicy struct {
	Ipv6Enabled bool  `json:"ipv6_enabled"`
	EsiMacMsb   uint8 `json:"esi_mac_msb"`
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
	if in.EsiMacMsb%2 != 0 {
		return fmt.Errorf("fabric addressing policy esi mac msb must be even, got %d", in.EsiMacMsb)
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
