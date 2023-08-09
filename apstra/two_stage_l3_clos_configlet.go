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

type rawTwoStageL3ClosConfigletData struct {
	Data      rawConfigletData `json:"configlet"`
	Condition string           `json:"condition"`
	Label     string           `json:"label"`
}

type TwoStageL3ClosConfigletData struct {
	Data      ConfigletData
	Condition string
	Label     string
}

type rawTwoStageL3ClosConfiglet struct {
	Data      rawConfigletData `json:"configlet"`
	Id        ObjectId         `json:"id"`
	Condition string           `json:"condition"`
	Label     string           `json:"label"`
}

type TwoStageL3ClosConfiglet struct {
	Data TwoStageL3ClosConfigletData
	Id   ObjectId
}

func (o *TwoStageL3ClosConfigletData) raw() *rawTwoStageL3ClosConfigletData {
	rawtc := rawTwoStageL3ClosConfigletData{}
	rawtc.Data = *o.Data.raw()
	rawtc.Condition = o.Condition
	rawtc.Label = o.Label
	return &rawtc
}

func (o *TwoStageL3ClosConfiglet) raw() *rawTwoStageL3ClosConfiglet {
	d := o.Data.raw()
	return &rawTwoStageL3ClosConfiglet{
		Data: rawConfigletData{
			RefArchs:    d.Data.RefArchs,
			Generators:  d.Data.Generators,
			DisplayName: d.Data.DisplayName,
		},
		Id:        o.Id,
		Condition: d.Condition,
		Label:     d.Label,
	}
}

func (o *rawTwoStageL3ClosConfiglet) polish() (*TwoStageL3ClosConfiglet, error) {
	c := TwoStageL3ClosConfiglet{}
	c.Id = o.Id
	d, err := o.Data.polish()
	if err != nil {
		return nil, err
	}
	c.Data = TwoStageL3ClosConfigletData{
		Data:      *d,
		Condition: o.Condition,
		Label:     o.Label,
	}
	return &c, err
}

func (o *TwoStageL3ClosClient) getAllConfiglets(ctx context.Context) ([]rawTwoStageL3ClosConfiglet, error) {
	response := &struct {
		Items []rawTwoStageL3ClosConfiglet `json:"items"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintConfiglets, o.blueprintId.String()),
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
		ids[i] = c.Id
	}
	return ids, nil
}

func (o *TwoStageL3ClosClient) getConfiglet(ctx context.Context, id ObjectId) (*rawTwoStageL3ClosConfiglet, error) {
	response := &rawTwoStageL3ClosConfiglet{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintConfigletsById, o.blueprintId.String(), id.String()),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *TwoStageL3ClosClient) getConfigletByName(ctx context.Context, name string) (*rawTwoStageL3ClosConfiglet,
	error) {
	cgs, err := o.getAllConfiglets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	idx := -1
	for i, t := range cgs {
		if t.Label == name {
			if idx == -1 {
				idx = i
			} else { // This is clearly the second occurrence
				return nil, ApstraClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("name '%s' does not uniquely identify a configlet", name),
				}
			}
		}
	}
	if idx != -1 {
		return &cgs[idx], nil
	}
	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no Configlet with name '%s' found", name),
	}
}

func (o *TwoStageL3ClosClient) createConfiglet(ctx context.Context, in *rawTwoStageL3ClosConfigletData) (ObjectId,
	error) {
	response := &objectIdResponse{}
	if len(in.Label) == 0 {
		in.Label = in.Data.DisplayName
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintConfiglets, o.blueprintId.String()),
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *TwoStageL3ClosClient) updateConfiglet(ctx context.Context, id ObjectId,
	in *rawTwoStageL3ClosConfigletData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintConfigletsById, o.blueprintId.String(), id),
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
