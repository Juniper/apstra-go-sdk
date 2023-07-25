package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintConfiglets       = apiUrlBlueprintById + apiUrlPathDelim + "configlets"
	apiUrlBlueprintConfigletsPrefix = apiUrlBlueprintConfiglets + apiUrlPathDelim
	apiUrlBlueprintConfigletsById   = apiUrlBlueprintConfigletsPrefix + "%s"
)

type TwoStageL3ClosConfiglet struct {
	Configlet rawConfigletRequest `json:"configlet"`
	Id        string              `json:"id"`
	Condition string              `json:"condition"`
	Label     string              `json:"label"`
}

func (o *TwoStageL3ClosClient) getAllConfiglets(ctx context.Context) ([]TwoStageL3ClosConfiglet, error) {
	blueprintconfigleturl := fmt.Sprintf(apiUrlBlueprintConfiglets, o.blueprintId.String())
	response := &struct {
		Items []TwoStageL3ClosConfiglet `json:"items"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      blueprintconfigleturl,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *TwoStageL3ClosClient) getAllConfigletIds(ctx context.Context) ([]ObjectId, error) {
	configlets, err := o.getAllConfiglets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	ids := make([]ObjectId, len(configlets))
	for i, c := range configlets {
		ids[i] = ObjectId(c.Id)
	}
	return ids, nil
}

func (o *TwoStageL3ClosClient) getConfiglet(ctx context.Context, id ObjectId) (*TwoStageL3ClosConfiglet, error) {
	blueprintconfigleturl := fmt.Sprintf(apiUrlBlueprintConfigletsById, o.blueprintId.String(), id.String())
	response := &TwoStageL3ClosConfiglet{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      blueprintconfigleturl,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *TwoStageL3ClosClient) getConfigletByName(ctx context.Context, name string) (*TwoStageL3ClosConfiglet, error) {
	cgs, err := o.getAllConfiglets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	for _, t := range cgs {
		if t.Label == name {
			return &t, nil
		}
	}
	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf(" Configlet with name '%s' not found", name),
	}
}

func (o *TwoStageL3ClosClient) importConfiglet(ctx context.Context, c ConfigletData, condition string, label string) (ObjectId, error) {
	blueprintconfigleturl := fmt.Sprintf(apiUrlBlueprintConfiglets, o.blueprintId.String())
	response := &objectIdResponse{}
	raw := (*ConfigletRequest)(&c).raw()
	in := TwoStageL3ClosConfiglet{
		Configlet: *raw,
		Condition: condition,
		Label:     label,
	}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      blueprintconfigleturl,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}
func (o *TwoStageL3ClosClient) importConfigletByID(ctx context.Context, id ObjectId, condition string, label string) (ObjectId, error) {
	blueprintconfigleturl := fmt.Sprintf(apiUrlBlueprintConfiglets, o.blueprintId.String())
	response := &objectIdResponse{}
	cfglet, err := o.client.GetConfiglet(ctx, id)
	cr := (*ConfigletRequest)(cfglet.Data).raw()
	if len(label) == 0 {
		label = cr.DisplayName
	}
	in := TwoStageL3ClosConfiglet{
		Configlet: *cr,
		Condition: condition,
		Label:     label,
	}
	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      blueprintconfigleturl,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *TwoStageL3ClosClient) updateConfiglet(ctx context.Context, in *TwoStageL3ClosConfiglet) error {
	blueprintconfigleturl := fmt.Sprintf(apiUrlBlueprintConfigletsById, o.blueprintId.String(), in.Id)
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   blueprintconfigleturl,
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *TwoStageL3ClosClient) deleteConfiglet(ctx context.Context, id ObjectId) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintConfigletsById, o.blueprintId.String(), id.String()),
	})
}
